package filebrowser

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"assistant/internal/bootstrap/psl"
)

var (
	ErrOutsideRoot = errors.New("path outside root")
	ErrDenied      = errors.New("path denied")
	ErrNotFound    = errors.New("not found")
	ErrTooLarge    = errors.New("file too large")
	ErrExists      = errors.New("file already exists")
	ErrBadName     = errors.New("invalid filename")
	ErrEmpty       = errors.New("no files to download")
	ErrIsRoot      = errors.New("cannot operate on root")
)

const (
	MaxRawBytes    int64 = 248 * 1024 * 1024       // 在线预览最大 248MB
	MaxUploadBytes int64 = 32 * 1024 * 1024 * 1024 // 上传最大 32GB
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

// IsPublicPath 判断给定的 path(用户原始入参，相对 root) 是否落在 public 列表之下(含等于)。
// 入参做与 resolve 相同的标准化，但不依赖文件系统。
func IsPublicPath(raw string) bool {
	cfg := psl.GetConfig().FileBrowser
	if len(cfg.Public) == 0 {
		return false
	}
	raw = strings.TrimSpace(raw)
	if raw == "" || filepath.IsAbs(raw) {
		return false
	}
	clean := filepath.Clean(raw)
	if clean == "." || strings.HasPrefix(clean, "..") {
		return false
	}
	for _, p := range cfg.Public {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if pathMatch(clean, filepath.Clean(p)) {
			return true
		}
	}
	return false
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
	if info.Size() > MaxRawBytes {
		return "", nil, ErrTooLarge
	}
	return abs, info, nil
}

// CreateUploadFile 在 dirRel 目录下以 name 为文件名创建文件用于上传。
// overwrite=true 时覆盖已有文件；否则 O_EXCL（失败返回 ErrExists）。
// 调用方负责把上传内容写入返回的 *os.File 后 Close。
// 返回值: 已打开的文件句柄、相对 root 的最终路径、错误。
func (s *Service) CreateUploadFile(dirRel, name string, overwrite bool) (*os.File, string, error) {
	// 文件名校验：禁止空、含分隔符、含 ".." 段
	name = strings.TrimSpace(name)
	if name == "" || name == "." || name == ".." {
		return nil, "", ErrBadName
	}
	if strings.ContainsAny(name, `/\`) {
		return nil, "", ErrBadName
	}
	if name != filepath.Base(name) {
		return nil, "", ErrBadName
	}

	dirAbs, dirNorm, err := s.resolve(dirRel)
	if err != nil {
		return nil, "", err
	}
	info, err := os.Stat(dirAbs)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", ErrNotFound
		}
		return nil, "", err
	}
	if !info.IsDir() {
		return nil, "", fmt.Errorf("not a directory")
	}

	// 用 resolve 再次校验目标路径(allow/deny/穿越)
	targetRel := filepath.Join(dirNorm, name)
	targetAbs, finalRel, err := s.resolve(targetRel)
	if err != nil {
		return nil, "", err
	}

	flags := os.O_WRONLY | os.O_CREATE
	if !overwrite {
		flags |= os.O_EXCL
	}
	f, err := os.OpenFile(targetAbs, flags, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, "", ErrExists
		}
		return nil, "", err
	}
	return f, finalRel, nil
}

// CreateTarGz 将指定路径列表(文件或目录)打包为 tar.gz，返回临时文件路径。
// 调用方负责在读取后删除该文件。目录会递归添加其下所有文件。
func (s *Service) CreateTarGz(paths []string) (string, error) {
	entries := make(map[string]string) // archivePath → absPath
	for _, rel := range paths {
		abs, normRel, err := s.resolve(rel)
		if err != nil {
			continue
		}
		if err := s.collectTarEntries(abs, normRel, entries); err != nil {
			continue
		}
	}
	if len(entries) == 0 {
		return "", ErrEmpty
	}

	tmp, err := os.CreateTemp("", "assistant-download-*.tar.gz")
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	gw := gzip.NewWriter(tmp)
	tw := tar.NewWriter(gw)

	for arcPath, abs := range entries {
		f, err := os.Open(abs)
		if err != nil {
			continue
		}
		info, _ := f.Stat()
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			f.Close()
			continue
		}
		hdr.Name = arcPath
		if err := tw.WriteHeader(hdr); err != nil {
			f.Close()
			continue
		}
		io.Copy(tw, f)
		f.Close()
	}
	if err := tw.Close(); err != nil {
		os.Remove(tmp.Name())
		return "", err
	}
	if err := gw.Close(); err != nil {
		os.Remove(tmp.Name())
		return "", err
	}
	return tmp.Name(), nil
}

// Mkdir 在 root 下新建目录。rel 包含新目录名。
func (s *Service) Mkdir(rel string) error {
	name := filepath.Base(rel)
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
		return ErrBadName
	}
	targetAbs, _, err := s.resolve(rel)
	if err != nil {
		return err
	}
	if _, err := os.Stat(targetAbs); err == nil {
		return ErrExists
	} else if !os.IsNotExist(err) {
		return err
	}
	parentAbs := filepath.Dir(targetAbs)
	pInfo, err := os.Stat(parentAbs)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNotFound
		}
		return err
	}
	if !pInfo.IsDir() {
		return fmt.Errorf("parent is not a directory")
	}
	return os.Mkdir(targetAbs, 0o755)
}

// Touch 创建新文件（不覆盖）。
func (s *Service) Touch(rel string) error {
	name := filepath.Base(rel)
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
		return ErrBadName
	}
	targetAbs, _, err := s.resolve(rel)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(targetAbs, os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if os.IsExist(err) {
			return ErrExists
		}
		return err
	}
	return f.Close()
}

// Rename 重命名文件或目录。
func (s *Service) Rename(rel, newName string) error {
	newName = strings.TrimSpace(newName)
	if newName == "" || newName == "." || newName == ".." || strings.ContainsAny(newName, `/\`) {
		return ErrBadName
	}
	srcAbs, normRel, err := s.resolve(rel)
	if err != nil {
		return err
	}
	if normRel == "" {
		return ErrIsRoot
	}
	parentAbs := filepath.Dir(srcAbs)
	dstAbs := filepath.Join(parentAbs, newName)
	cfg := psl.GetConfig().FileBrowser
	dstRel := filepath.Join(filepath.Dir(normRel), newName)
	if !isAllowed(dstRel, cfg.Allow, cfg.Deny) {
		return ErrDenied
	}
	if _, err := os.Stat(dstAbs); err == nil {
		return ErrExists
	} else if !os.IsNotExist(err) {
		return err
	}
	return os.Rename(srcAbs, dstAbs)
}

// Move 批量移动文件/目录到目标目录。
func (s *Service) Move(paths []string, destRel string) (int, error) {
	if len(paths) == 0 {
		return 0, fmt.Errorf("no paths provided")
	}
	destAbs, destNorm, err := s.resolve(destRel)
	if err != nil {
		return 0, err
	}
	dInfo, err := os.Stat(destAbs)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, ErrNotFound
		}
		return 0, err
	}
	if !dInfo.IsDir() {
		return 0, fmt.Errorf("destination is not a directory")
	}
	moved := 0
	cfg := psl.GetConfig().FileBrowser
	for _, rel := range paths {
		srcAbs, _, err := s.resolve(rel)
		if err != nil {
			continue
		}
		name := filepath.Base(rel)
		dstAbs := filepath.Join(destAbs, name)
		dstRel := filepath.Join(destNorm, name)
		if !isAllowed(dstRel, cfg.Allow, cfg.Deny) {
			continue
		}
		if err := os.Rename(srcAbs, dstAbs); err == nil {
			moved++
		}
	}
	if moved == 0 {
		return 0, fmt.Errorf("nothing was moved")
	}
	return moved, nil
}

// Copy 批量复制文件/目录到目标目录。
func (s *Service) Copy(paths []string, destRel string) (int, error) {
	if len(paths) == 0 {
		return 0, fmt.Errorf("no paths provided")
	}
	destAbs, destNorm, err := s.resolve(destRel)
	if err != nil {
		return 0, err
	}
	dInfo, err := os.Stat(destAbs)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, ErrNotFound
		}
		return 0, err
	}
	if !dInfo.IsDir() {
		return 0, fmt.Errorf("destination is not a directory")
	}
	copied := 0
	cfg := psl.GetConfig().FileBrowser
	for _, rel := range paths {
		srcAbs, _, err := s.resolve(rel)
		if err != nil {
			continue
		}
		name := filepath.Base(rel)
		dstAbs := filepath.Join(destAbs, name)
		dstRel := filepath.Join(destNorm, name)
		if !isAllowed(dstRel, cfg.Allow, cfg.Deny) {
			continue
		}
		if err := s.copyOne(srcAbs, dstAbs); err == nil {
			copied++
		}
	}
	if copied == 0 {
		return 0, fmt.Errorf("nothing was copied")
	}
	return copied, nil
}

func (s *Service) copyOne(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return s.copyDir(src, dst)
	}
	return s.copyFile(src, dst)
}

func (s *Service) copyFile(src, dst string) error {
	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()
	dstF, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if _, err := io.Copy(dstF, srcF); err != nil {
		dstF.Close()
		os.Remove(dst)
		return err
	}
	if err := dstF.Close(); err != nil {
		os.Remove(dst)
		return err
	}
	return nil
}

func (s *Service) copyDir(src, dst string) error {
	if err := os.Mkdir(dst, 0o755); err != nil {
		return err
	}
	dirents, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, de := range dirents {
		childSrc := filepath.Join(src, de.Name())
		childDst := filepath.Join(dst, de.Name())
		if err := s.copyOne(childSrc, childDst); err != nil {
			return err
		}
	}
	return nil
}

const trashRel = ".local/share/Trash"

func (s *Service) trashDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, trashRel), nil
}

func (s *Service) ensureTrash() (filesDir, infoDir string, err error) {
	root, err := s.trashDir()
	if err != nil {
		return "", "", err
	}
	filesDir = filepath.Join(root, "files")
	infoDir = filepath.Join(root, "info")
	for _, d := range []string{filesDir, infoDir} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			return "", "", err
		}
	}
	return
}

func uniqueTrashName(filesDir, base string) string {
	candidate := base
	for i := 1; ; i++ {
		if _, err := os.Stat(filepath.Join(filesDir, candidate)); os.IsNotExist(err) {
			return candidate
		}
		candidate = fmt.Sprintf("%s (%d)", base, i)
	}
}

func writeTrashInfo(infoDir, trashName, origAbs string) error {
	path := filepath.Join(infoDir, trashName+".trashinfo")
	content := fmt.Sprintf("[Trash Info]\nPath=%s\nDeletionDate=%s\n", origAbs, time.Now().Format(time.RFC3339))
	return os.WriteFile(path, []byte(content), 0o644)
}

func readTrashInfo(infoPath string) (origPath, delDate string, err error) {
	data, err := os.ReadFile(infoPath)
	if err != nil {
		return "", "", err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Path=") {
			origPath = line[5:]
		} else if strings.HasPrefix(line, "DeletionDate=") {
			delDate = line[13:]
		}
	}
	if origPath == "" {
		return "", "", fmt.Errorf("invalid trashinfo: missing Path")
	}
	return
}

// TrashEntry 表示回收站中的一项。
type TrashEntry struct {
	TrashName    string `json:"trash_name"`
	OriginalPath string `json:"original_path"`
	DeletionDate string `json:"deletion_date"`
	Size         int64  `json:"size"`
	IsDir        bool   `json:"is_dir"`
}

// Delete 将文件/目录移动到回收站 ~/.local/share/Trash。
func (s *Service) Delete(paths []string) (int, error) {
	if len(paths) == 0 {
		return 0, fmt.Errorf("no paths provided")
	}
	for _, rel := range paths {
		_, normRel, err := s.resolve(rel)
		if err != nil {
			continue
		}
		if normRel == "" {
			return 0, ErrIsRoot
		}
	}
	filesDir, infoDir, err := s.ensureTrash()
	if err != nil {
		return 0, err
	}
	deleted := 0
	for _, rel := range paths {
		abs, _, err := s.resolve(rel)
		if err != nil {
			continue
		}
		if _, err := os.Stat(abs); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			continue
		}
		base := filepath.Base(abs)
		trashName := uniqueTrashName(filesDir, base)
		dst := filepath.Join(filesDir, trashName)
		if err := s.moveFile(abs, dst); err != nil {
			continue
		}
		if err := writeTrashInfo(infoDir, trashName, abs); err != nil {
			os.Rename(dst, abs)
			continue
		}
		deleted++
	}
	if deleted == 0 {
		return 0, ErrNotFound
	}
	return deleted, nil
}

// moveFile 将 src 移到 dst（先尝试 rename，跨文件系统则 copy+delete）。
func (s *Service) moveFile(src, dst string) error {
	err := os.Rename(src, dst)
	if err == nil {
		return nil
	}
	if errors.Is(err, syscall.EXDEV) {
		if err := s.copyOne(src, dst); err != nil {
			os.RemoveAll(dst)
			return err
		}
		if err := os.RemoveAll(src); err != nil {
			os.RemoveAll(dst)
			return err
		}
		return nil
	}
	return err
}

// ListTrash 列出回收站内容。
func (s *Service) ListTrash() ([]TrashEntry, error) {
	root, err := s.trashDir()
	if err != nil {
		return nil, err
	}
	infoDir := filepath.Join(root, "info")
	filesDir := filepath.Join(root, "files")
	entries, err := os.ReadDir(infoDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []TrashEntry{}, nil
		}
		return nil, err
	}
	result := make([]TrashEntry, 0, len(entries))
	for _, de := range entries {
		if de.IsDir() || !strings.HasSuffix(de.Name(), ".trashinfo") {
			continue
		}
		trashName := strings.TrimSuffix(de.Name(), ".trashinfo")
		origPath, delDate, err := readTrashInfo(filepath.Join(infoDir, de.Name()))
		if err != nil {
			continue
		}
		trashAbs := filepath.Join(filesDir, trashName)
		fi, err := os.Stat(trashAbs)
		if err != nil {
			continue
		}
		result = append(result, TrashEntry{
			TrashName:    trashName,
			OriginalPath: origPath,
			DeletionDate: delDate,
			Size:         fi.Size(),
			IsDir:        fi.IsDir(),
		})
	}
	return result, nil
}

// RestoreTrash 从回收站恢复到原位置。
func (s *Service) RestoreTrash(trashNames []string) (int, error) {
	root, err := s.trashDir()
	if err != nil {
		return 0, err
	}
	filesDir := filepath.Join(root, "files")
	infoDir := filepath.Join(root, "info")
	restored := 0
	for _, name := range trashNames {
		infoPath := filepath.Join(infoDir, name+".trashinfo")
		origPath, _, err := readTrashInfo(infoPath)
		if err != nil {
			continue
		}
		trashAbs := filepath.Join(filesDir, name)
		if _, err := os.Stat(trashAbs); err != nil {
			continue
		}
		if _, err := os.Stat(origPath); err == nil {
			continue
		}
		parent := filepath.Dir(origPath)
		if err := os.MkdirAll(parent, 0o755); err != nil {
			continue
		}
		if err := s.moveFile(trashAbs, origPath); err != nil {
			continue
		}
		os.Remove(infoPath)
		restored++
	}
	if restored == 0 {
		return 0, fmt.Errorf("nothing was restored")
	}
	return restored, nil
}

// PermanentDelete 从回收站永久删除。
func (s *Service) PermanentDelete(trashNames []string) (int, error) {
	root, err := s.trashDir()
	if err != nil {
		return 0, err
	}
	filesDir := filepath.Join(root, "files")
	infoDir := filepath.Join(root, "info")
	deleted := 0
	for _, name := range trashNames {
		if name == "" {
			continue
		}
		if err := os.RemoveAll(filepath.Join(filesDir, name)); err == nil {
			os.Remove(filepath.Join(infoDir, name+".trashinfo"))
			deleted++
		}
	}
	if deleted == 0 {
		return 0, fmt.Errorf("nothing was permanently deleted")
	}
	return deleted, nil
}

// collectTarEntries 递归收集文件路径，供 CreateTarGz 使用。
func (s *Service) collectTarEntries(abs, rel string, entries map[string]string) error {
	cfg := psl.GetConfig().FileBrowser
	info, err := os.Stat(abs)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		if isAllowed(rel, cfg.Allow, cfg.Deny) {
			entries[rel] = abs
		}
		return nil
	}
	dirents, err := os.ReadDir(abs)
	if err != nil {
		return err
	}
	for _, de := range dirents {
		childAbs := filepath.Join(abs, de.Name())
		childRel := filepath.Join(rel, de.Name())
		if !isAllowed(childRel, cfg.Allow, cfg.Deny) {
			continue
		}
		s.collectTarEntries(childAbs, childRel, entries)
	}
	return nil
}
