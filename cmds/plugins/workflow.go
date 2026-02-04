package plugins

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"cmds/utils"

	"golang.design/x/clipboard"
)

func SnipFzf() error {
	snipDir := os.ExpandEnv("$HOME/share/github/obsidian/.snippets")
	if _, err := os.Stat(snipDir); err != nil {
		return fmt.Errorf("snippet dir not found: %s", snipDir)
	}

	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		return err
	}
	defer readPipe.Close()

	// shell 脚本: 选择 snip
	script := fmt.Sprintf(`
cd %s || exit 1
selected=$(find . -type f | sed 's|^\./||' |
fzf \
  --prompt="Snip> " \
  --height=100%% \
  --border \
  --preview='bat --style=plain --color=always {} 2>/dev/null || cat {}' \
  --preview-window=right:60%%)
printf '%%s' "$selected" >&3
`, snipDir)

	cmd := exec.Command(utils.GetOSDefaultTerminal(), "-e", "sh", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{writePipe}

	if err := cmd.Start(); err != nil {
		writePipe.Close()
		return err
	}
	writePipe.Close()

	if err := cmd.Wait(); err != nil {
		return nil // 用户 ESC 退出时，fzf 返回非 0，直接忽略
	}

	data, err := io.ReadAll(readPipe)
	if err != nil {
		return err
	}

	file := strings.TrimSpace(string(data))
	if file == "" {
		return nil
	}

	content, err := os.ReadFile(filepath.Join(snipDir, file))
	if err != nil {
		return err
	}

	if err := clipboard.Init(); err != nil {
		return err
	}
	utils.Notify(fmt.Sprintf("Snip copied:\n%s", file))
	clipboard.Write(clipboard.FmtText, content)
	time.Sleep(30 * time.Second)
	utils.Notify("previous clipboard expired")
	return nil
}

func GetIPAddress() {
	interfaceName := "wlan0"
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		utils.Notify(err)
		return
	}

	addrs, err := iface.Addrs()
	if err != nil {
		utils.Notify(err)
		return
	}

	for _, addr := range addrs {
		var ip net.IP

		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if ip.IsLoopback() {
			continue
		}

		if ip.To4() != nil {
			content := ip.String()
			utils.Notify(fmt.Sprintf("get success: %s", content))
			clipboard.Write(clipboard.FmtText, []byte(content))
			time.Sleep(30 * time.Second)
			utils.Notify("previous clipboard expired")
		}
	}
}

func GetCurDatetime() {
	err := clipboard.Init()
	if err != nil {
		utils.Notify(err)
		return
	}
	content := time.Now().Format(time.DateTime)
	utils.Notify(fmt.Sprintf("get success: %s", content))
	clipboard.Write(clipboard.FmtText, []byte(content))
	time.Sleep(30 * time.Second)
	utils.Notify("previous clipboard expired")
}

func GetCurUnixSec() {
	err := clipboard.Init()
	if err != nil {
		utils.Notify(err)
		return
	}
	content := fmt.Sprintf("%d", time.Now().Unix())
	utils.Notify(fmt.Sprintf("get success: %s", content))
	clipboard.Write(clipboard.FmtText, []byte(content))
	time.Sleep(30 * time.Second)
	utils.Notify("previous clipboard expired")
}

func TransformDatetime2UnixSec() {
	err := clipboard.Init()
	if err != nil {
		utils.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)

	t, err := time.Parse(time.DateTime, strings.TrimSpace(string(text)))
	if err != nil {
		utils.Notify(err)
		return
	}
	formatedText := fmt.Sprintf("%d", t.Unix())
	utils.Notify(fmt.Sprintf("tranfer success: \n%s", formatedText))
	clipboard.Write(clipboard.FmtText, []byte(formatedText))
	time.Sleep(30 * time.Second)
	utils.Notify("previous clipboard expired")
}

func TransformUnixSec2DateTime() {
	err := clipboard.Init()
	if err != nil {
		utils.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)
	unix, err := strconv.ParseInt(strings.TrimSpace(string(text)), 10, 64)
	if err != nil {
		utils.Notify(err)
		return
	}
	datetime := time.Unix(unix, 0).Format(time.DateTime)
	utils.Notify(fmt.Sprintf("tranfer success: \n%s", datetime))
	clipboard.Write(clipboard.FmtText, []byte(datetime))
	time.Sleep(30 * time.Second)
	utils.Notify("previous clipboard expired")
}

func LazyOpenSearchFile() {
	utils.RunScript("bash", fmt.Sprintf("%s -e lazy-open-search-file", utils.GetOSDefaultTerminal()))
}

func LazyOpenSearchBook() {
	utils.RunScript("bash", fmt.Sprintf("%s -e lazy-open-search-book", utils.GetOSDefaultTerminal()))
}

func LazyOpenSearchWiki() {
	utils.RunScript("bash", fmt.Sprintf("%s -e lazy-open-search-wiki", utils.GetOSDefaultTerminal()))
}

func LazyOpenSearchMedia() {
	utils.RunScript("bash", fmt.Sprintf("%s -e lazy-open-search-media", utils.GetOSDefaultTerminal()))
}

func LazyOpenSearchFileContent() {
	utils.RunScript("bash", fmt.Sprintf("%s -e lazy-open-search-file-content", utils.GetOSDefaultTerminal()))
}

func SearchFromWeb(content string) {
	q := url.QueryEscape(content)
	u := "https://www.google.com/search?q=" + q
	OpenUrlWithQutebrowser(u)()
}

func SearchBooksOnline() {
	content, err := utils.GetInput("search books online: ")
	if err != nil {
		utils.Notify(err)
		return
	}
	q := url.QueryEscape(content)
	urls := []string{
		"https://openlibrary.org/search?q=" + q,
		"https://z-lib.id/s?q=" + q,
	}
	wg := sync.WaitGroup{}
	for _, u := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			OpenUrlWithQutebrowser(u)()
		}(u)
	}
	wg.Wait()
}

func SearchVideosOnline() {
	content, err := utils.GetInput("search videos online: ")
	if err != nil {
		utils.Notify(err)
		return
	}
	q := url.QueryEscape(content)
	urls := []string{
		"https://search.bilibili.com/all?keyword=" + q,
		"https://www.youtube.com/results?search_query=" + q,
	}
	wg := sync.WaitGroup{}
	for _, u := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			OpenUrlWithQutebrowser(u)()
		}(u)
	}
	wg.Wait()
}

func NoteToDo() {
	fileDir := path.Join(os.Getenv("HOME"), GithubPath, "obsidian", "working")
	filePath := path.Join(fileDir, "TODO.md")
	if !utils.IsDirExists(fileDir) {
		if err := os.Mkdir(fileDir, 0o755); err != nil {
			utils.Notify(err)
			return
		}
	}
	if !utils.IsFileExists(filePath) {
		f, err := os.Create(filePath)
		if err != nil {
			utils.Notify(err)
			return
		}
		_, _ = fmt.Fprintf(f, "\n## ToDo\n\n")
		_ = f.Close()
	}
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0o755)
	if err != nil {
		utils.Notify(err)
		return
	}
	_, _ = fmt.Fprintf(f, "\n- [ ] %s", time.Now().Format(time.DateTime))
	_ = f.Close()
	_, _, _ = utils.RunScript("bash", fmt.Sprintf("%s -e nvim +$ '%s'", utils.GetOSDefaultTerminal(), filePath))
}

func NoteScripts() {
	fileDir := path.Join(os.Getenv("HOME"), GithubPath, "obsidian", "working")
	filePath := path.Join(fileDir, "scripts.md")
	if !utils.IsDirExists(fileDir) {
		if err := os.Mkdir(fileDir, 0o755); err != nil {
			utils.Notify(err)
			return
		}
	}
	if !utils.IsFileExists(filePath) {
		f, err := os.Create(filePath)
		if err != nil {
			utils.Notify(err)
			return
		}
		_, _ = fmt.Fprintf(f, "\n## Scripts\n\n")
		_ = f.Close()
	}
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0o755)
	if err != nil {
		utils.Notify(err)
		return
	}
	_, _ = fmt.Fprintf(f, "\n\n###")
	_ = f.Close()
	_, _, _ = utils.RunScript(
		"bash",
		fmt.Sprintf("%s -e nvim +$ '%s'", utils.GetOSDefaultTerminal(), filePath),
	)
}

func NoteMonthlyWork() {
	t := time.Now()
	dateStr := t.Format("2006-01")
	fileDir := path.Join(os.Getenv("HOME"), GithubPath, "obsidian", "working", "monthly.work")
	filePath := path.Join(fileDir, dateStr+".md")
	if !utils.IsDirExists(fileDir) {
		if err := os.Mkdir(fileDir, 0o755); err != nil {
			utils.Notify(err)
			return
		}
	}
	if !utils.IsFileExists(filePath) {
		f, err := os.Create(filePath)
		if err != nil {
			utils.Notify(err)
			return
		}
		_, _ = fmt.Fprintf(f, "\n## %s\n\n", dateStr)
		_ = f.Close()
	}
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0o755)
	if err != nil {
		utils.Notify(err)
		return
	}
	_, _ = fmt.Fprintf(f, "\n### %s\n\n", time.Now().Format(time.DateTime))
	_ = f.Close()
	_, _, _ = utils.RunScript("bash", fmt.Sprintf("%s -e nvim +$ '%s'", utils.GetOSDefaultTerminal(), filePath))
}

func HandleCopied() {
	if err := clipboard.Init(); err != nil {
		utils.Notify(err)
		return
	}

	text := strings.TrimSpace(string(clipboard.Read(clipboard.FmtText)))
	if text == "" {
		return
	}

	// 1) log/stacktrace -> file:line(:col)
	if file, line, col, ok := extractFirstExistingFileLocation(text); ok {
		openFileAt(file, line, col)
		return
	}

	// 2) direct path
	if utils.Exists(text) && utils.IsFile(text) {
		utils.Lazy("open", text)
		return
	}

	// 3) markdown link: [x](url)
	if url, ok := extractMarkdownURL(text); ok {
		OpenUrlWithQutebrowser(url)()
		return
	}

	// 4) url
	if utils.IsURL(text) {
		OpenUrlWithQutebrowser(text)()
		return
	}

	// 5) default: web search
	SearchFromWeb(text)
}

type fileLocationPattern struct {
	re   *regexp.Regexp
	file int
	line int
	col  int
}

var fileLocationPatterns = []fileLocationPattern{
	// /abs/path/file.ext:123 or /abs/path/file.ext:123:45
	{re: regexp.MustCompile(`(?m)(/[^:\s]+):(\d+)(?::(\d+))?`), file: 1, line: 2, col: 3},
	// relative/path/file.ext:123 or relative/path/file.ext:123:45
	{re: regexp.MustCompile(`(?m)([A-Za-z0-9_./\-~]+\.[A-Za-z0-9]+):(\d+)(?::(\d+))?`), file: 1, line: 2, col: 3},
	// python: File "...", line 123
	{re: regexp.MustCompile(`(?m)File\s+"([^"]+)",\s+line\s+(\d+)`), file: 1, line: 2, col: 0},
	// node: at ... (/path/file.js:12:34)
	{re: regexp.MustCompile(`(?m)\((/[^:()]+):(\d+):(\d+)\)`), file: 1, line: 2, col: 3},
	{re: regexp.MustCompile(`(?m)\s+at\s+(/[^:\s]+):(\d+):(\d+)`), file: 1, line: 2, col: 3},
	// rust: --> /path/file.rs:12:34
	{re: regexp.MustCompile(`(?m)-->\s+(/[^:\s]+):(\d+):(\d+)`), file: 1, line: 2, col: 3},
}

func extractFirstExistingFileLocation(text string) (file string, line, col int, ok bool) {
	for _, p := range fileLocationPatterns {
		m := p.re.FindStringSubmatch(text)
		if len(m) == 0 {
			continue
		}

		candidate := strings.TrimSpace(m[p.file])
		candidate = strings.TrimSuffix(candidate, ")")
		candidate = strings.TrimSuffix(candidate, ":")

		l, err := strconv.Atoi(m[p.line])
		if err != nil || l <= 0 {
			continue
		}
		c := 0
		if p.col > 0 && p.col < len(m) {
			if m[p.col] != "" {
				if x, err := strconv.Atoi(m[p.col]); err == nil {
					c = x
				}
			}
		}

		if !filepath.IsAbs(candidate) {
			if abs, err := filepath.Abs(candidate); err == nil {
				candidate = abs
			}
		}

		if utils.Exists(candidate) {
			return candidate, l, c, true
		}
	}
	return "", 0, 0, false
}

func openFileAt(file string, line, col int) {
	term := utils.GetOSDefaultTerminal()
	fileQ := shellSingleQuote(file)
	if col > 0 {
		cmd := fmt.Sprintf(
			"%s -e nvim +'%s' %s",
			term,
			fmt.Sprintf("call cursor(%d,%d)", line, col),
			fileQ,
		)
		_, _, err := utils.RunScript("bash", cmd)
		if err != nil {
			utils.Notify(err)
		}
		return
	}

	cmd := fmt.Sprintf("%s -e nvim +%d %s", term, line, fileQ)
	_, _, err := utils.RunScript("bash", cmd)
	if err != nil {
		utils.Notify(err)
	}
}

func shellSingleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func extractMarkdownURL(text string) (url string, ok bool) {
	m := regexp.MustCompile(`\[[^\]]*\]\((https?://[^\s)]+)\)`).FindStringSubmatch(strings.TrimSpace(text))
	if len(m) == 2 {
		return m[1], true
	}
	return "", false
}

func SshTo() {
	mysshListFilePath := path.Join(os.Getenv("HOME"), ".ssh/my.ssh.list")
	if !utils.IsFileExists(mysshListFilePath) {
		f, err := os.Create(mysshListFilePath)
		if err != nil {
			utils.Notify(err)
			return
		}
		_ = f.Close()
	}

	// read from to ~/.ssh/my.ssh.list
	b, err := os.ReadFile(mysshListFilePath)
	if err != nil {
		utils.Notify(err)
		return
	}
	mySshList := []map[string]string{}
	slice1 := strings.Split(string(b), "\n")
	for _, x := range slice1 {
		x = strings.TrimSpace(x)
		slice2 := regexp.MustCompile(`[ \r\t\s]+`).Split(x, -1)
		if len(slice2) < 3 {
			continue
		}
		host := strings.TrimSpace(slice2[0])
		user := strings.TrimSpace(slice2[1])
		password := strings.TrimSpace(slice2[2])
		slice3 := strings.Split(x, "#")
		if len(slice3) != 2 {
			continue
		}
		desc := strings.TrimSpace(slice3[1])
		mySshList = append(
			mySshList,
			map[string]string{"host": host, "user": user, "password": password, "desc": desc},
		)
	}

	// read from ~/.ssh/known_hosts
	knownHosts, err := utils.GetKnownHosts()
	if err != nil {
		utils.Notify(err)
		return
	}

	// prompt
	promptList := []string{}
	for _, x := range mySshList {
		promptList = append(promptList, fmt.Sprintf("%-20s %-20s %-20s # %s", x["host"], x["user"], x["password"], x["desc"]))
	}
	for host := range knownHosts {
		promptList = append(promptList, fmt.Sprintf("%-20s", knownHosts[host]))
	}

	// choose
	chioce, err := utils.Choose("ssh to: ", promptList)
	if err != nil {
		utils.Notify(err)
		return
	}
	chioce = strings.TrimSpace(chioce)
	slice := regexp.MustCompile(`[ \r\t\s]+`).Split(chioce, -1)

	switch {
	case len(slice) > 3:
		host := strings.TrimSpace(slice[0])
		user := strings.TrimSpace(slice[1])
		password := strings.TrimSpace(slice[2])
		err = utils.SSH(host, 22, user, password)
		if err != nil {
			utils.Notify(err)
			return
		}
		return
	default:
		host := strings.TrimSpace(slice[0])
		user, err := utils.GetInput("user: ")
		if err != nil {
			utils.Notify(err)
			return
		}
		password, err := utils.GetInput("password: ")
		if err != nil {
			utils.Notify(err)
			return
		}
		desc, err := utils.GetInput("desc: ")
		if err != nil {
			utils.Notify(err)
			return
		}

		err = utils.SSH(host, 22, user, password)
		if err != nil {
			utils.Notify(err)
			return
		}

		// append to ~/.ssh/my.ssh.list
		file, err := os.OpenFile(
			mysshListFilePath,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0o666,
		)
		if err != nil {
			utils.Notify(err)
		}
		defer file.Close()
		writer := bufio.NewWriter(file)
		_, _ = fmt.Fprintf(writer, "%-20s %-20s %-20s # %s\r\n", host, user, password, desc)
		_ = writer.Flush()
	}
}
