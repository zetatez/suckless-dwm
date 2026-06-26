package filebrowser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"assistant/internal/bootstrap/psl"
)

var (
	ErrOutsideRoot = errors.New("path outside root")
	ErrDenied      = errors.New("path denied")
	ErrNotFound    = errors.New("not found")
	ErrTooLarge    = errors.New("file too large")
)

type Entry struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir"`
	ModTime int64  `json:"mtime"`
	Mode    string `json:"mode"`
}

type ListResult struct {
	Root    string  `json:"root"`
	Path    string  `json:"path"`
	Parent  string  `json:"parent"`
	Entries []Entry `json:"entries"`
}

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) Root() string { return psl.GetConfig().FileBrowser.Root }

// pathMatch 判断 rel 是否等于 pat 或在 pat 之下。
func pathMatch(rel, pat string) bool {
	if rel == pat {
		return true
	}
	return strings.HasPrefix(rel, pat+string(filepath.Separator))
}

// isAllowed 判定 normRel 是否可访问/可见。
// 规则：
//   - root("") 始终允许（否则白名单下连入口都进不去）
//   - allow 为空：放行
//   - allow 非空：normRel 必须命中其中一项；命中含两种语义：
//     a) normRel 在某个 allow 之下（含等于）
//     b) normRel 是某个 allow 的祖先（让用户能穿越进入）
//   - deny 优先于 allow：只要 normRel 在某个 deny 之下（含等于）即拒绝
func isAllowed(normRel string, allow, deny []string) bool {
	if normRel == "" {
		// 根目录：仅检查 deny 中是否包含 "" / "."（极少见，但容错）
		for _, d := range deny {
			d = strings.TrimSpace(d)
			if d == "" {
				continue
			}
			if filepath.Clean(d) == "." {
				return false
			}
		}
		return true
	}
	// deny 优先
	for _, d := range deny {
		d = strings.TrimSpace(d)
		if d == "" {
			continue
		}
		if pathMatch(normRel, filepath.Clean(d)) {
			return false
		}
	}
	// allow（白名单）
	hasAllow := false
	for _, a := range allow {
		if strings.TrimSpace(a) != "" {
			hasAllow = true
			break
		}
	}
	if !hasAllow {
		return true
	}
	for _, a := range allow {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		aClean := filepath.Clean(a)
		// rel 在 allow 之下（含等于）或 rel 是 allow 的祖先
		if pathMatch(normRel, aClean) || pathMatch(aClean, normRel) {
			return true
		}
	}
	return false
}

// resolve 把外部传入的相对(基于 root)路径，校验后转为绝对路径。
// 同时返回相对于 root 的规范化路径（用于回显/拼链接）。
func (s *Service) resolve(rel string) (abs string, normRel string, err error) {
	cfg := psl.GetConfig().FileBrowser
	root, err := filepath.Abs(cfg.Root)
	if err != nil {
		return "", "", err
	}
	rel = strings.TrimSpace(rel)
	if rel == "" {
		rel = "."
	}
	// 不允许绝对路径，强制相对 root
	if filepath.IsAbs(rel) {
		return "", "", ErrOutsideRoot
	}
	joined := filepath.Join(root, rel)
	clean := filepath.Clean(joined)
	// 防穿越：clean 必须以 root 为前缀
	rootWithSep := root + string(filepath.Separator)
	if clean != root && !strings.HasPrefix(clean, rootWithSep) {
		return "", "", ErrOutsideRoot
	}

	normRel, err = filepath.Rel(root, clean)
	if err != nil {
		return "", "", err
	}
	if normRel == "." {
		normRel = ""
	}
	if !isAllowed(normRel, cfg.Allow, cfg.Deny) {
		return "", "", ErrDenied
	}
	return clean, normRel, nil
}

func (s *Service) ListDir(rel string) (*ListResult, error) {
	abs, normRel, err := s.resolve(rel)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("not a directory")
	}
	dirents, err := os.ReadDir(abs)
	if err != nil {
		return nil, err
	}
	cfg := psl.GetConfig().FileBrowser
	entries := make([]Entry, 0, len(dirents))
	for _, de := range dirents {
		name := de.Name()
		childRel := filepath.Join(normRel, name)
		if !isAllowed(childRel, cfg.Allow, cfg.Deny) {
			continue
		}
		fi, err := de.Info()
		if err != nil {
			continue
		}
		entries = append(entries, Entry{
			Name:    name,
			Path:    childRel,
			Size:    fi.Size(),
			IsDir:   fi.IsDir(),
			ModTime: fi.ModTime().Unix(),
			Mode:    fi.Mode().String(),
		})
	}
	parent := ""
	if normRel != "" {
		parent = filepath.Dir(normRel)
		if parent == "." {
			parent = ""
		}
	}
	return &ListResult{
		Root:    s.Root(),
		Path:    normRel,
		Parent:  parent,
		Entries: entries,
	}, nil
}

// ResolveFile 校验后返回可被 Gin 用作 c.File / c.FileAttachment 的绝对路径。
func (s *Service) ResolveFile(rel string) (string, os.FileInfo, error) {
	abs, _, err := s.resolve(rel)
	if err != nil {
		return "", nil, err
	}
	info, err := os.Stat(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil, ErrNotFound
		}
		return "", nil, err
	}
	if info.IsDir() {
		return "", nil, fmt.Errorf("is a directory")
	}
	return abs, info, nil
}

// ResolveRaw 类似 ResolveFile，但额外检查大小上限。
func (s *Service) ResolveRaw(rel string) (string, os.FileInfo, error) {
	abs, info, err := s.ResolveFile(rel)
	if err != nil {
		return "", nil, err
	}
	if info.Size() > psl.GetConfig().FileBrowser.MaxRawBytes {
		return "", nil, ErrTooLarge
	}
	return abs, info, nil
}
