package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	ErrNotGitRepo   = errors.New("not a git repository")
	ErrNoRemote     = errors.New("remote not found")
	ErrNoUpstream   = errors.New("no upstream configured")
	ErrDetachedHead = errors.New("detached HEAD")
)

type GitCommandResult struct {
	Args     []string
	Stdout   string
	Stderr   string
	ExitCode int
}

func (r GitCommandResult) String() string {
	return strings.Join(append([]string{"git"}, r.Args...), " ")
}

type GitError struct {
	Result GitCommandResult
	Kind   error
	Cause  error
}

func (e *GitError) Error() string {
	msg := strings.TrimSpace(e.Result.Stderr)
	if msg == "" {
		msg = strings.TrimSpace(e.Result.Stdout)
	}
	if msg != "" {
		return fmt.Sprintf("git command failed (exit=%d): %s", e.Result.ExitCode, msg)
	}
	return fmt.Sprintf("git command failed (exit=%d)", e.Result.ExitCode)
}

func (e *GitError) Unwrap() error { return e.Cause }

func (e *GitError) Is(target error) bool {
	if target == nil {
		return false
	}
	return e.Kind == target
}

type GitRunner struct {
	WorkDir string
	Binary  string
}

func NewGitRunner(workDir string) *GitRunner {
	return &GitRunner{WorkDir: workDir, Binary: "git"}
}

func (g *GitRunner) Run(ctx context.Context, args ...string) (GitCommandResult, error) {
	bin := g.Binary
	if strings.TrimSpace(bin) == "" {
		bin = "git"
	}

	// Prefix options must appear before the subcommand.
	prefixed := make([]string, 0, len(args)+4)
	prefixed = append(prefixed, "--no-pager", "-c", "color.ui=false")
	prefixed = append(prefixed, args...)

	cmd := exec.CommandContext(ctx, bin, prefixed...)
	if strings.TrimSpace(g.WorkDir) != "" {
		cmd.Dir = g.WorkDir
	}

	// Keep git invocations predictable and non-interactive.
	cmd.Env = append(os.Environ(),
		"GIT_PAGER=cat",
		"GIT_TERMINAL_PROMPT=0",
		"GCM_INTERACTIVE=never",
		"GIT_OPTIONAL_LOCKS=0",
	)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	res := GitCommandResult{
		Args:     prefixed,
		Stdout:   outBuf.String(),
		Stderr:   errBuf.String(),
		ExitCode: 0,
	}
	if err == nil {
		return res, nil
	}

	// If context was cancelled, prefer returning ctx error.
	if ctxErr := ctx.Err(); ctxErr != nil {
		return res, ctxErr
	}

	var ee *exec.ExitError
	if errors.As(err, &ee) {
		res.ExitCode = ee.ExitCode()
		kind := classifyGitFailure(res)
		return res, &GitError{Result: res, Kind: kind, Cause: err}
	}

	return res, fmt.Errorf("git command failed: %w", err)
}

func classifyGitFailure(res GitCommandResult) error {
	msg := strings.ToLower(strings.TrimSpace(res.Stderr))
	if msg == "" {
		msg = strings.ToLower(strings.TrimSpace(res.Stdout))
	}

	switch {
	case strings.Contains(msg, "not a git repository") || strings.Contains(msg, "must be run in a work tree"):
		return ErrNotGitRepo
	case strings.Contains(msg, "no such remote"):
		return ErrNoRemote
	case strings.Contains(msg, "no upstream configured") || (strings.Contains(msg, "no such branch") && strings.Contains(msg, "@{u}")):
		return ErrNoUpstream
	case strings.Contains(msg, "ref head is not a symbolic ref") || strings.Contains(msg, "not a symbolic ref: head"):
		return ErrDetachedHead
	default:
		return nil
	}
}

// Status (porcelain parsing)

type GitStatusEntry struct {
	X byte
	Y byte

	Path     string
	OrigPath string // for renames/copies (when available)
}

func (e GitStatusEntry) IsUntracked() bool { return e.X == '?' && e.Y == '?' }

func (e GitStatusEntry) IsIgnored() bool { return e.X == '!' && e.Y == '!' }

func (e GitStatusEntry) IsStaged() bool { return e.X != ' ' && !e.IsUntracked() && !e.IsIgnored() }

func (e GitStatusEntry) IsUnstaged() bool { return e.Y != ' ' && !e.IsUntracked() && !e.IsIgnored() }

type GitStatus struct {
	Entries []GitStatusEntry
}

func (s GitStatus) HasChanges() bool {
	for _, e := range s.Entries {
		if !e.IsIgnored() {
			return true
		}
	}
	return false
}

func (s GitStatus) UntrackedPaths() []string {
	var out []string
	for _, e := range s.Entries {
		if e.IsUntracked() {
			out = append(out, e.Path)
		}
	}
	return out
}

func GitStatusPorcelain(ctx context.Context, dir string) (GitStatus, error) {
	r := NewGitRunner(dir)
	res, err := r.Run(ctx, "status", "--porcelain=v1", "-z")
	if err != nil {
		if errors.Is(err, ErrNotGitRepo) {
			return GitStatus{}, ErrNotGitRepo
		}
		return GitStatus{}, err
	}
	entries, err := parseGitStatusPorcelainZ(res.Stdout)
	if err != nil {
		return GitStatus{}, err
	}
	return GitStatus{Entries: entries}, nil
}

func parseGitStatusPorcelainZ(out string) ([]GitStatusEntry, error) {
	if out == "" {
		return nil, nil
	}
	parts := strings.Split(out, "\x00")
	// Trailing NUL produces an extra empty element.
	if len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}

	var entries []GitStatusEntry
	for i := 0; i < len(parts); i++ {
		p := parts[i]
		if len(p) < 3 {
			return nil, fmt.Errorf("invalid porcelain record: %q", p)
		}
		x, y := p[0], p[1]
		if p[2] != ' ' {
			return nil, fmt.Errorf("invalid porcelain record (missing space): %q", p)
		}
		path := p[3:]
		entry := GitStatusEntry{X: x, Y: y, Path: path}

		// In -z format, renames/copies are encoded as: "R  src\0dst\0".
		if x == 'R' || x == 'C' || y == 'R' || y == 'C' {
			if i+1 >= len(parts) {
				return nil, fmt.Errorf("invalid rename/copy record: %q", p)
			}
			entry.OrigPath = path
			entry.Path = parts[i+1]
			i++
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// Repo

type GitRepo struct {
	Root   string
	runner *GitRunner
}

// DefaultBranchRef attempts to detect the repository's default branch.
// It prefers the remote HEAD symbolic ref (e.g. "origin/main") and falls back
// to common branch names.
func (r *GitRepo) DefaultBranchRef(ctx context.Context) (string, error) {
	return r.DefaultBranchRefForRemote(ctx, "origin")
}

// DefaultBranchName is like DefaultBranchRef but returns just the branch name
// (e.g. "main").
func (r *GitRepo) DefaultBranchName(ctx context.Context) (string, error) {
	ref, err := r.DefaultBranchRef(ctx)
	if err != nil {
		return "", err
	}
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return "", fmt.Errorf("default branch not found")
	}
	if i := strings.LastIndex(ref, "/"); i >= 0 {
		return ref[i+1:], nil
	}
	return ref, nil
}

func (r *GitRepo) DefaultBranchRefForRemote(ctx context.Context, remote string) (string, error) {
	remote = strings.TrimSpace(remote)
	if remote == "" {
		remote = "origin"
	}

	// 1) Prefer "refs/remotes/<remote>/HEAD" when configured.
	remoteHeadRef := fmt.Sprintf("refs/remotes/%s/HEAD", remote)
	res, err := r.runner.Run(ctx, "symbolic-ref", "--quiet", "--short", remoteHeadRef)
	if err == nil {
		ref := strings.TrimSpace(res.Stdout)
		if ref != "" {
			return ref, nil
		}
	}
	if err != nil {
		if errors.Is(err, ErrNotGitRepo) {
			return "", ErrNotGitRepo
		}
		// Ignore missing/invalid remote HEAD and fall through.
	}

	// 2) Parse "git remote show <remote>" for "HEAD branch: <name>".
	res, err = r.runner.Run(ctx, "remote", "show", remote)
	if err == nil {
		if head := parseRemoteShowHeadBranch(res.Stdout); head != "" {
			return remote + "/" + head, nil
		}
	}
	if err != nil {
		if errors.Is(err, ErrNotGitRepo) {
			return "", ErrNotGitRepo
		}
		if errors.Is(err, ErrNoRemote) {
			// No remote configured; fall back to local branches.
		}
	}

	// 3) Heuristic fallbacks.
	remoteCandidates := []string{remote + "/main", remote + "/master"}
	for _, c := range remoteCandidates {
		exists, exErr := r.refExists(ctx, "refs/remotes/"+c)
		if exErr != nil {
			return "", exErr
		}
		if exists {
			return c, nil
		}
	}

	localCandidates := []string{"main", "master"}
	for _, c := range localCandidates {
		exists, exErr := r.refExists(ctx, "refs/heads/"+c)
		if exErr != nil {
			return "", exErr
		}
		if exists {
			return c, nil
		}
	}

	return "", fmt.Errorf("default branch not found")
}

func (r *GitRepo) refExists(ctx context.Context, fullRef string) (bool, error) {
	_, err := r.runner.Run(ctx, "show-ref", "--verify", "--quiet", fullRef)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, ErrNotGitRepo) {
		return false, ErrNotGitRepo
	}
	var ge *GitError
	if errors.As(err, &ge) {
		// show-ref --verify exits with 1 when a ref doesn't exist.
		if ge.Result.ExitCode == 1 {
			return false, nil
		}
	}
	return false, err
}

func parseRemoteShowHeadBranch(out string) string {
	for _, line := range strings.Split(out, "\n") {
		l := strings.TrimSpace(line)
		if !strings.HasPrefix(l, "HEAD branch:") {
			continue
		}
		v := strings.TrimSpace(strings.TrimPrefix(l, "HEAD branch:"))
		v = strings.Trim(v, "\"'")
		if v == "" || v == "(unknown)" {
			return ""
		}
		return v
	}
	return ""
}

func IsGitRepo(ctx context.Context, dir string) (bool, error) {
	r := NewGitRunner(dir)
	res, err := r.Run(ctx, "rev-parse", "--is-inside-work-tree")
	if err != nil {
		if errors.Is(err, ErrNotGitRepo) {
			return false, nil
		}
		return false, err
	}
	return strings.TrimSpace(res.Stdout) == "true", nil
}

func OpenGitRepo(ctx context.Context, startDir string) (*GitRepo, error) {
	r := NewGitRunner(startDir)
	res, err := r.Run(ctx, "rev-parse", "--show-toplevel")
	if err != nil {
		if errors.Is(err, ErrNotGitRepo) {
			return nil, ErrNotGitRepo
		}
		return nil, err
	}
	root := strings.TrimSpace(res.Stdout)
	if root == "" {
		return nil, fmt.Errorf("git repo root not found")
	}
	return &GitRepo{Root: root, runner: NewGitRunner(root)}, nil
}

func (r *GitRepo) Status(ctx context.Context) (GitStatus, error) {
	res, err := r.runner.Run(ctx, "status", "--porcelain=v1", "-z")
	if err != nil {
		if errors.Is(err, ErrNotGitRepo) {
			return GitStatus{}, ErrNotGitRepo
		}
		return GitStatus{}, err
	}
	entries, err := parseGitStatusPorcelainZ(res.Stdout)
	if err != nil {
		return GitStatus{}, err
	}
	return GitStatus{Entries: entries}, nil
}

// StatusRawPorcelainZ returns the raw output of `git status --porcelain=v1 -z`.
func (r *GitRepo) StatusRawPorcelainZ(ctx context.Context) (string, error) {
	res, err := r.runner.Run(ctx, "status", "--porcelain=v1", "-z")
	if err != nil {
		if errors.Is(err, ErrNotGitRepo) {
			return "", ErrNotGitRepo
		}
		return "", err
	}
	return res.Stdout, nil
}

func (r *GitRepo) Diff(ctx context.Context, staged bool, paths ...string) (string, error) {
	args := []string{"diff", "--no-color", "--no-ext-diff"}
	if staged {
		args = append(args, "--cached")
	}
	if len(paths) > 0 {
		args = append(args, "--")
		args = append(args, paths...)
	}
	res, err := r.runner.Run(ctx, args...)
	if err != nil {
		if errors.Is(err, ErrNotGitRepo) {
			return "", ErrNotGitRepo
		}
		return "", err
	}
	return res.Stdout, nil
}

// DiffRange returns `git diff <rangeSpec>` output (e.g. "main...HEAD").
func (r *GitRepo) DiffRange(ctx context.Context, rangeSpec string, paths ...string) (string, error) {
	if strings.TrimSpace(rangeSpec) == "" {
		return "", fmt.Errorf("rangeSpec is empty")
	}
	args := []string{"diff", "--no-color", "--no-ext-diff", rangeSpec}
	if len(paths) > 0 {
		args = append(args, "--")
		args = append(args, paths...)
	}
	res, err := r.runner.Run(ctx, args...)
	if err != nil {
		if errors.Is(err, ErrNotGitRepo) {
			return "", ErrNotGitRepo
		}
		return "", err
	}
	return res.Stdout, nil
}

func (r *GitRepo) Add(ctx context.Context, paths ...string) error {
	args := []string{"add"}
	if len(paths) == 0 {
		args = append(args, "-A")
		_, err := r.runner.Run(ctx, args...)
		if err != nil && errors.Is(err, ErrNotGitRepo) {
			return ErrNotGitRepo
		}
		return err
	}
	args = append(args, "--")
	args = append(args, paths...)
	_, err := r.runner.Run(ctx, args...)
	if err != nil && errors.Is(err, ErrNotGitRepo) {
		return ErrNotGitRepo
	}
	return err
}

func (r *GitRepo) Commit(ctx context.Context, message string) error {
	if strings.TrimSpace(message) == "" {
		return fmt.Errorf("commit message is empty")
	}
	_, err := r.runner.Run(ctx, "commit", "-m", message)
	if err != nil && errors.Is(err, ErrNotGitRepo) {
		return ErrNotGitRepo
	}
	return err
}

func (r *GitRepo) CurrentBranch(ctx context.Context) (string, error) {
	res, err := r.runner.Run(ctx, "symbolic-ref", "--quiet", "--short", "HEAD")
	if err != nil {
		if errors.Is(err, ErrNotGitRepo) {
			return "", ErrNotGitRepo
		}
		if errors.Is(err, ErrDetachedHead) {
			return "", ErrDetachedHead
		}
		return "", err
	}
	return strings.TrimSpace(res.Stdout), nil
}

// UpstreamRef returns upstream ref name like "origin/main".
func (r *GitRepo) UpstreamRef(ctx context.Context) (string, error) {
	res, err := r.runner.Run(ctx, "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err != nil {
		if errors.Is(err, ErrNoUpstream) {
			return "", ErrNoUpstream
		}
		if errors.Is(err, ErrNotGitRepo) {
			return "", ErrNotGitRepo
		}
		return "", err
	}
	return strings.TrimSpace(res.Stdout), nil
}

func (r *GitRepo) LogSubjects(ctx context.Context, n int) ([]string, error) {
	if n <= 0 {
		n = 10
	}
	res, err := r.runner.Run(ctx, "log", "-n", strconv.Itoa(n), "--pretty=format:%s")
	if err != nil {
		if errors.Is(err, ErrNotGitRepo) {
			return nil, ErrNotGitRepo
		}
		return nil, err
	}
	out := strings.TrimRight(res.Stdout, "\n")
	if strings.TrimSpace(out) == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

func (r *GitRepo) AheadBehindUpstream(ctx context.Context) (ahead int, behind int, err error) {
	res, err := r.runner.Run(ctx, "rev-list", "--left-right", "--count", "@{u}...HEAD")
	if err != nil {
		if errors.Is(err, ErrNoUpstream) {
			return 0, 0, ErrNoUpstream
		}
		if errors.Is(err, ErrNotGitRepo) {
			return 0, 0, ErrNotGitRepo
		}
		return 0, 0, err
	}
	fields := strings.Fields(strings.TrimSpace(res.Stdout))
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("unexpected rev-list output: %q", strings.TrimSpace(res.Stdout))
	}
	behind, err = strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, fmt.Errorf("parse behind failed: %w", err)
	}
	ahead, err = strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, fmt.Errorf("parse ahead failed: %w", err)
	}
	return ahead, behind, nil
}

func (r *GitRepo) HeadSHA(ctx context.Context) (string, error) {
	res, err := r.runner.Run(ctx, "rev-parse", "HEAD")
	if err != nil {
		if errors.Is(err, ErrNotGitRepo) {
			return "", ErrNotGitRepo
		}
		return "", err
	}
	return strings.TrimSpace(res.Stdout), nil
}

func (r *GitRepo) TrackedFiles(ctx context.Context) ([]string, error) {
	res, err := r.runner.Run(ctx, "ls-files", "-z")
	if err != nil {
		if errors.Is(err, ErrNotGitRepo) {
			return nil, ErrNotGitRepo
		}
		return nil, err
	}
	if res.Stdout == "" {
		return nil, nil
	}
	parts := strings.Split(res.Stdout, "\x00")
	if len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return parts, nil
}

func (r *GitRepo) ChangedPaths(ctx context.Context, staged bool, paths ...string) ([]string, error) {
	args := []string{"diff", "--name-only", "-z"}
	if staged {
		args = append(args, "--cached")
	}
	if len(paths) > 0 {
		args = append(args, "--")
		args = append(args, paths...)
	}
	res, err := r.runner.Run(ctx, args...)
	if err != nil {
		if errors.Is(err, ErrNotGitRepo) {
			return nil, ErrNotGitRepo
		}
		return nil, err
	}
	return splitNulList(res.Stdout), nil
}

func (r *GitRepo) ChangedPathsRange(ctx context.Context, rangeSpec string, paths ...string) ([]string, error) {
	if strings.TrimSpace(rangeSpec) == "" {
		return nil, fmt.Errorf("rangeSpec is empty")
	}
	args := []string{"diff", "--name-only", "-z", rangeSpec}
	if len(paths) > 0 {
		args = append(args, "--")
		args = append(args, paths...)
	}
	res, err := r.runner.Run(ctx, args...)
	if err != nil {
		if errors.Is(err, ErrNotGitRepo) {
			return nil, ErrNotGitRepo
		}
		return nil, err
	}
	return splitNulList(res.Stdout), nil
}

func splitNulList(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, "\x00")
	if len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	if len(parts) == 1 && parts[0] == "" {
		return nil
	}
	return parts
}
