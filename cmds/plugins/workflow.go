package plugins

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"cmds/utils"

	"golang.design/x/clipboard"
)

func GetHostName() {
	err := clipboard.Init()
	if err != nil {
		utils.Notify(err)
		return
	}
	cmd := "hostname"
	stdout, _, err := utils.RunScript("bash", cmd)
	if err != nil {
		utils.Notify(err)
		return
	}
	content := stdout
	utils.Notify(fmt.Sprintf("get success: %s", content))
	changed := clipboard.Write(clipboard.FmtText, []byte(content))
	<-changed
	utils.Notify("previous clipboard expired")
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
			changed := clipboard.Write(clipboard.FmtText, []byte(content))
			<-changed
			utils.Notify("previous clipboard expired")
		}
	}
}

func GetCurrentDatetime() {
	err := clipboard.Init()
	if err != nil {
		utils.Notify(err)
		return
	}
	content := time.Now().Format(time.DateTime)
	utils.Notify(fmt.Sprintf("get success: %s", content))
	changed := clipboard.Write(clipboard.FmtText, []byte(content))
	<-changed
	utils.Notify("previous clipboard expired")
}

func GetCurrentUnixSec() {
	err := clipboard.Init()
	if err != nil {
		utils.Notify(err)
		return
	}
	content := fmt.Sprintf("%d", time.Now().Unix())
	utils.Notify(fmt.Sprintf("get success: %s", content))
	changed := clipboard.Write(clipboard.FmtText, []byte(content))
	<-changed
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
	changed := clipboard.Write(clipboard.FmtText, []byte(formatedText))
	<-changed
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
	changed := clipboard.Write(clipboard.FmtText, []byte(datetime))
	<-changed
	utils.Notify("previous clipboard expired")
}

func LazyOpenSearchFile() {
	_, _, _ = utils.RunScript("bash", fmt.Sprintf("%s -e lazy-open-search-file", utils.GetOSDefaultTerminal()))
}

func LazyOpenSearchBook() {
	_, _, _ = utils.RunScript("bash", fmt.Sprintf("%s -e lazy-open-search-book", utils.GetOSDefaultTerminal()))
}

func LazyOpenSearchWiki() {
	_, _, _ = utils.RunScript("bash", fmt.Sprintf("%s -e lazy-open-search-wiki", utils.GetOSDefaultTerminal()))
}

func LazyOpenSearchMedia() {
	_, _, _ = utils.RunScript("bash", fmt.Sprintf("%s -e lazy-open-search-media", utils.GetOSDefaultTerminal()))
}

func LazyOpenSearchFileContent() {
	_, _, _ = utils.RunScript("bash", fmt.Sprintf("%s -e lazy-open-search-file-content", utils.GetOSDefaultTerminal()))
}

func SearchFromWeb(content string) {
	url := fmt.Sprintf("https://www.google.com/search?q='%s'", content)
	OpenUrlWithQutebrowser(url)()
}

func SearchBooksOnline() {
	content, err := utils.GetInput("search books online: ")
	if err != nil {
		utils.Notify(err)
		return
	}
	urls := []string{
		"https://openlibrary.org/search?q='%s'",
		"https://z-lib.id/s?q='%s'",
	}
	wg := sync.WaitGroup{}
	for _, urlTpl := range urls {
		wg.Add(1)
		go func(urlTpl string) {
			defer wg.Done()
			url := fmt.Sprintf(urlTpl, content)
			OpenUrlWithQutebrowser(url)()
		}(urlTpl)
	}
	wg.Wait()
}

func SearchVideosOnline() {
	content, err := utils.GetInput("search videos online: ")
	if err != nil {
		utils.Notify(err)
		return
	}
	urls := []string{
		"https://search.bilibili.com/all?keyword='%s'",
		"https://www.youtube.com/results?search_query='%s'",
	}
	wg := sync.WaitGroup{}
	for _, urlTpl := range urls {
		wg.Add(1)
		go func(urlTpl string) {
			defer wg.Done()
			url := fmt.Sprintf(urlTpl, content)
			OpenUrlWithQutebrowser(url)()
		}(urlTpl)
	}
	wg.Wait()
}

func NoteScripts() {
	fileDir := path.Join(os.Getenv("HOME"), GithubPath, "obsidian")
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

func NoteToDo() {
	fileDir := path.Join(os.Getenv("HOME"), GithubPath, "obsidian")
	filePath := path.Join(fileDir, "ToDo.md")
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
	_, _ = fmt.Fprintf(f, "\n- [ ] %s", time.Now().Format(time.DateOnly))
	_ = f.Close()
	_, _, _ = utils.RunScript("bash", fmt.Sprintf("%s -e nvim +$ '%s'", utils.GetOSDefaultTerminal(), filePath))
}

func NoteDiary() {
	dateStr := time.Now().Format(time.DateOnly)
	fileDir := path.Join(os.Getenv("HOME"), GithubPath, "obsidian", "diary")
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
		_, _ = fmt.Fprintf(f, "\n### Diary %s\n\n", dateStr)
		_ = f.Close()
	}
	_, _, _ = utils.RunScript("bash", fmt.Sprintf("%s -e nvim +$ +$ '%s'", utils.GetOSDefaultTerminal(), filePath))
}

func NoteTimeline() {
	t := time.Now()
	dateStr := t.Format(time.DateOnly)
	datetimeStr := t.Format(time.DateTime)
	fileDir := path.Join(os.Getenv("HOME"), GithubPath, "obsidian", "timeline")
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
		_, _ = fmt.Fprintf(f, "\n## Time Line %s\n\n", dateStr)
		_ = f.Close()
	}
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0o755)
	if err != nil {
		utils.Notify(err)
		return
	}
	_, _ = fmt.Fprintf(f, "\n### %s\n\n", datetimeStr)
	_ = f.Close()
	_, _, _ = utils.RunScript("bash", fmt.Sprintf("%s -e nvim +$ '%s'", utils.GetOSDefaultTerminal(), filePath))
}

func NoteFlashCard() {
	t := time.Now()
	fileDir := path.Join(os.Getenv("HOME"), GithubPath, "obsidian", "flash-card")
	filePath := path.Join(fileDir, t.Format("2006-01-02.15.04.05.000000000")+".md")
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
		_, _ = fmt.Fprintf(f, "### Flash Card %s\n\n", t.Format(time.DateTime))
		_ = f.Close()
	}
	_, _, _ = utils.RunScript("bash", fmt.Sprintf("%s -e nvim +$ '%s'", utils.GetOSDefaultTerminal(), filePath))
}

func HandleCopied() {
	err := clipboard.Init()
	if err != nil {
		utils.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)
	content := strings.TrimSpace(string(text))
	switch {
	case utils.Exists(content) && utils.IsFile(content):
		utils.Lazy("open", content)
		return
	case utils.IsURL(content):
		url := content
		OpenUrlWithChrome(url)()
		return
	default:
		SearchFromWeb(content)
	}
}

func JumpToCodeFromLog() {
	err := clipboard.Init()
	if err != nil {
		utils.Notify(err)
		return
	}
	textbyte := clipboard.Read(clipboard.FmtText)
	text := string(textbyte)
	regex := `(?P<filepath>/[^\:]+):(?P<row>\d+)\s+`
	r := regexp.MustCompile(regex)
	match := r.FindStringSubmatch(text)
	if len(match) < 3 {
		utils.Notify("not match")
		return
	}
	filepath := match[1]
	row := match[2]
	_, _, err = utils.RunScript("bash", fmt.Sprintf("%s -e nvim +%s %s", utils.GetOSDefaultTerminal(), row, filepath))
	if err != nil {
		utils.Notify(err)
		return
	}
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
