package utils

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ApplyPatch applies an OpenCode/opencode-style patch to files under the
// current working directory.
func ApplyPatch(patchText string) error {
	return ApplyPatchInDir(".", patchText)
}

// ApplyPatchInDir applies an OpenCode/opencode-style patch to files under baseDir.
//
// Supported operations:
// - *** Add File: <path>
// - *** Update File: <path> (optional: *** Move to: <path>)
// - *** Delete File: <path>
//
// For Add File: content lines must start with '+'.
// For Update File: use hunks starting with '@@ <anchor>' and then lines prefixed
// with ' ', '+', '-' (unified-diff-like, without line numbers).
func ApplyPatchInDir(baseDir string, patchText string) error {
	ops, err := parseApplyPatch(patchText)
	if err != nil {
		return err
	}
	for _, op := range ops {
		if err := op.apply(baseDir); err != nil {
			return err
		}
	}
	return nil
}

// ApplyPatchContext keeps a context-aware entrypoint.
func ApplyPatchContext(ctx context.Context, baseDir string, patchText string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return ApplyPatchInDir(baseDir, patchText)
}

type patchOpKind int

const (
	patchOpAdd patchOpKind = iota + 1
	patchOpUpdate
	patchOpDelete
)

type patchOp struct {
	kind patchOpKind
	path string
	move string

	addLines []string
	hunks    []patchHunk
}

type patchHunk struct {
	anchor string
	lines  []hunkLine
}

type hunkLine struct {
	op   byte // ' ', '+', '-'
	text string
}

func (op patchOp) apply(baseDir string) error {
	path, err := cleanRelPath(op.path)
	if err != nil {
		return err
	}
	full := filepath.Join(baseDir, path)

	switch op.kind {
	case patchOpAdd:
		if Exists(full) {
			return fmt.Errorf("%w: add file already exists: %s", ErrInvalidPatch, path)
		}
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			return fmt.Errorf("mkdir failed: %w", err)
		}
		content := strings.Join(op.addLines, "\n")
		return writeFileAtomic(full, []byte(content), 0644)

	case patchOpDelete:
		if !Exists(full) {
			return fmt.Errorf("%w: delete file not found: %s", ErrInvalidPatch, path)
		}
		if err := os.Remove(full); err != nil {
			return fmt.Errorf("delete file failed: %w", err)
		}
		return nil

	case patchOpUpdate:
		data, err := os.ReadFile(full)
		if err != nil {
			return fmt.Errorf("read file failed: %w", err)
		}
		updated, err := applyHunks(normalizeNewlines(string(data)), op.hunks)
		if err != nil {
			return fmt.Errorf("%w: update %s: %v", ErrInvalidPatch, path, err)
		}
		dstFull := full
		if strings.TrimSpace(op.move) != "" {
			movePath, err := cleanRelPath(op.move)
			if err != nil {
				return err
			}
			dstFull = filepath.Join(baseDir, movePath)
			if movePath != path {
				if Exists(dstFull) {
					return fmt.Errorf("%w: move target exists: %s", ErrInvalidPatch, movePath)
				}
				if err := os.MkdirAll(filepath.Dir(dstFull), 0755); err != nil {
					return fmt.Errorf("mkdir failed: %w", err)
				}
			}
		}

		if err := writeFileAtomic(dstFull, []byte(updated), 0644); err != nil {
			return err
		}
		if dstFull != full {
			if err := os.Remove(full); err != nil {
				return fmt.Errorf("remove old file after move failed: %w", err)
			}
		}
		return nil

	default:
		return fmt.Errorf("%w: unknown patch operation", ErrInvalidPatch)
	}
}

func parseApplyPatch(patchText string) ([]patchOp, error) {
	patchText = normalizeNewlines(patchText)

	lines, err := readAllLines(patchText)
	if err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("%w: empty patch", ErrInvalidPatch)
	}
	if strings.TrimSpace(lines[0]) != "*** Begin Patch" {
		return nil, fmt.Errorf("%w: missing '*** Begin Patch'", ErrInvalidPatch)
	}

	var ops []patchOp
	for i := 1; i < len(lines); {
		l := lines[i]
		if strings.TrimSpace(l) == "" {
			i++
			continue
		}
		if strings.TrimSpace(l) == "*** End Patch" {
			return ops, nil
		}

		switch {
		case strings.HasPrefix(l, "*** Add File: "):
			op := patchOp{kind: patchOpAdd, path: strings.TrimSpace(strings.TrimPrefix(l, "*** Add File: "))}
			i++
			for i < len(lines) {
				if strings.HasPrefix(lines[i], "*** ") {
					break
				}
				ln := lines[i]
				if len(ln) == 0 || ln[0] != '+' {
					return nil, fmt.Errorf("%w: add file line must start with '+': %q", ErrInvalidPatch, ln)
				}
				op.addLines = append(op.addLines, ln[1:])
				i++
			}
			ops = append(ops, op)
			continue

		case strings.HasPrefix(l, "*** Delete File: "):
			op := patchOp{kind: patchOpDelete, path: strings.TrimSpace(strings.TrimPrefix(l, "*** Delete File: "))}
			ops = append(ops, op)
			i++
			continue

		case strings.HasPrefix(l, "*** Update File: "):
			op := patchOp{kind: patchOpUpdate, path: strings.TrimSpace(strings.TrimPrefix(l, "*** Update File: "))}
			i++
			for i < len(lines) {
				ln := lines[i]
				if strings.TrimSpace(ln) == "" {
					i++
					continue
				}
				if strings.HasPrefix(ln, "*** Move to: ") {
					op.move = strings.TrimSpace(strings.TrimPrefix(ln, "*** Move to: "))
					i++
					continue
				}
				if strings.HasPrefix(ln, "*** ") {
					break
				}
				if strings.TrimSpace(ln) == "*** End Patch" {
					break
				}
				if !strings.HasPrefix(ln, "@@") {
					return nil, fmt.Errorf("%w: expected hunk header '@@', got: %q", ErrInvalidPatch, ln)
				}

				h := patchHunk{anchor: strings.TrimSpace(strings.TrimPrefix(ln, "@@"))}
				i++
				for i < len(lines) {
					pl := lines[i]
					if strings.HasPrefix(pl, "@@") || strings.HasPrefix(pl, "*** ") || strings.TrimSpace(pl) == "*** End Patch" {
						break
					}
					if pl == "" {
						h.lines = append(h.lines, hunkLine{op: ' ', text: ""})
						i++
						continue
					}
					c := pl[0]
					if c != ' ' && c != '+' && c != '-' {
						return nil, fmt.Errorf("%w: invalid hunk line prefix: %q", ErrInvalidPatch, pl)
					}
					h.lines = append(h.lines, hunkLine{op: c, text: pl[1:]})
					i++
				}
				op.hunks = append(op.hunks, h)
				continue
			}
			ops = append(ops, op)
			continue

		default:
			return nil, fmt.Errorf("%w: unknown patch header: %q", ErrInvalidPatch, l)
		}
	}

	return nil, fmt.Errorf("%w: missing '*** End Patch'", ErrInvalidPatch)
}

func applyHunks(content string, hunks []patchHunk) (string, error) {
	if len(hunks) == 0 {
		return content, nil
	}
	lines := strings.Split(content, "\n")

	for _, h := range hunks {
		start := 0
		if strings.TrimSpace(h.anchor) != "" {
			idx := indexOfLine(lines, h.anchor)
			if idx < 0 {
				return "", fmt.Errorf("anchor not found: %q", h.anchor)
			}
			start = idx
		}

		oldSeq, newSeq, err := hunkSequences(h.lines)
		if err != nil {
			return "", err
		}
		if len(oldSeq) == 0 {
			insAt := start
			if strings.TrimSpace(h.anchor) != "" {
				insAt = start + 1
			}
			lines = append(lines[:insAt], append(newSeq, lines[insAt:]...)...)
			continue
		}
		pos := indexOfSequence(lines, oldSeq, start)
		if pos < 0 {
			return "", fmt.Errorf("hunk target not found near anchor %q", h.anchor)
		}
		lines = append(lines[:pos], append(newSeq, lines[pos+len(oldSeq):]...)...)
	}

	return strings.Join(lines, "\n"), nil
}

func hunkSequences(hunkLines []hunkLine) (oldSeq []string, newSeq []string, err error) {
	for _, hl := range hunkLines {
		switch hl.op {
		case ' ':
			oldSeq = append(oldSeq, hl.text)
			newSeq = append(newSeq, hl.text)
		case '-':
			oldSeq = append(oldSeq, hl.text)
		case '+':
			newSeq = append(newSeq, hl.text)
		default:
			return nil, nil, fmt.Errorf("invalid hunk op: %q", hl.op)
		}
	}
	return oldSeq, newSeq, nil
}

func indexOfLine(lines []string, needle string) int {
	for i, l := range lines {
		if l == needle {
			return i
		}
	}
	return -1
}

func indexOfSequence(lines []string, seq []string, start int) int {
	if len(seq) == 0 {
		return start
	}
	if start < 0 {
		start = 0
	}
	for i := start; i+len(seq) <= len(lines); i++ {
		ok := true
		for j := 0; j < len(seq); j++ {
			if lines[i+j] != seq[j] {
				ok = false
				break
			}
		}
		if ok {
			return i
		}
	}
	return -1
}

func cleanRelPath(p string) (string, error) {
	p = strings.TrimSpace(p)
	if p == "" {
		return "", fmt.Errorf("%w: empty path", ErrInvalidPatch)
	}
	if filepath.IsAbs(p) {
		return "", fmt.Errorf("%w: absolute path not allowed: %s", ErrInvalidPatch, p)
	}
	clean := filepath.Clean(p)
	if clean == "." {
		return "", fmt.Errorf("%w: invalid path: %s", ErrInvalidPatch, p)
	}
	if strings.HasPrefix(clean, ".."+string(filepath.Separator)) || clean == ".." {
		return "", fmt.Errorf("%w: path escapes base dir: %s", ErrInvalidPatch, p)
	}
	return clean, nil
}

func normalizeNewlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

func readAllLines(text string) ([]string, error) {
	s := bufio.NewScanner(strings.NewReader(text))
	s.Buffer(make([]byte, 64*1024), 4*1024*1024)
	var out []string
	for s.Scan() {
		out = append(out, s.Text())
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("%w: scan patch failed: %v", ErrInvalidPatch, err)
	}
	return out, nil
}
