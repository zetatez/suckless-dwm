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

	"cmds/sugar"

	"golang.design/x/clipboard"
	// ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
)

func GetHostName() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	cmd := "hostname"
	stdout, _, err := sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	content := stdout
	sugar.Notify(fmt.Sprintf("get success: %s", content))
	changed := clipboard.Write(clipboard.FmtText, []byte(content))
	<-changed
	sugar.Notify("previous clipboard expired")
}

func GetIPAddress() {
	interfaceName := "wlan0"
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		sugar.Notify(err)
		return
	}

	addrs, err := iface.Addrs()
	if err != nil {
		sugar.Notify(err)
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
			sugar.Notify(fmt.Sprintf("get success: %s", content))
			changed := clipboard.Write(clipboard.FmtText, []byte(content))
			<-changed
			sugar.Notify("previous clipboard expired")
		}
	}
	return
}

func GetCurrentDatetime() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	content := time.Now().Format(time.DateTime)
	sugar.Notify(fmt.Sprintf("get success: %s", content))
	changed := clipboard.Write(clipboard.FmtText, []byte(content))
	<-changed
	sugar.Notify("previous clipboard expired")
}

func GetCurrentUnixSec() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	content := fmt.Sprintf("%d", time.Now().Unix())
	sugar.Notify(fmt.Sprintf("get success: %s", content))
	changed := clipboard.Write(clipboard.FmtText, []byte(content))
	<-changed
	sugar.Notify("previous clipboard expired")
}

func TransformDatetime2UnixSec() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)

	t, err := time.Parse(time.DateTime, strings.TrimSpace(string(text)))
	if err != nil {
		sugar.Notify(err)
		return
	}
	formatedText := fmt.Sprintf("%d", t.Unix())
	sugar.Notify(fmt.Sprintf("tranfer success: \n%s", formatedText))
	changed := clipboard.Write(clipboard.FmtText, []byte(formatedText))
	<-changed
	sugar.Notify("previous clipboard expired")
}

func TransformUnixSec2DateTime() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)
	unix, err := strconv.ParseInt(strings.TrimSpace(string(text)), 10, 64)
	if err != nil {
		sugar.Notify(err)
		return
	}
	datetime := time.Unix(unix, 0).Format(time.DateTime)
	sugar.Notify(fmt.Sprintf("tranfer success: \n%s", datetime))
	changed := clipboard.Write(clipboard.FmtText, []byte(datetime))
	<-changed
	sugar.Notify("previous clipboard expired")
}

func LazyOpenSearchFile() {
	cmd := `st -e lazy-open-search-file`
	sugar.NewExecService().RunScriptShell(cmd)
}

func LazyOpenSearchBook() {
	cmd := `st -e lazy-open-search-book`
	sugar.NewExecService().RunScriptShell(cmd)
}

func LazyOpenSearchWiki() {
	cmd := `st -e lazy-open-search-wiki`
	sugar.NewExecService().RunScriptShell(cmd)
}

func LazyOpenSearchMedia() {
	cmd := `st -e lazy-open-search-media`
	sugar.NewExecService().RunScriptShell(cmd)
}

func LazyOpenSearchFileContent() {
	cmd := `st -e lazy-open-search-file-content`
	sugar.NewExecService().RunScriptShell(cmd)
}

func SearchFromWeb(content string) {
	sugar.NewExecService().RunScriptShell(
		fmt.Sprintf(
			// "chrome --proxy-server=%s https://www.google.com/search?q='%s'",
			"qutebrowser --set content.proxy %s https://www.google.com/search?q='%s'",
			ProxyServer,
			content,
		),
	)
}

func SearchBooksOnline() {
	content, err := sugar.GetInput("search books online: ")
	if err != nil {
		sugar.Notify(err)
		return
	}
	urls := []string{
		"https://libgen.is/search.php?req='%s'",
		"https://openlibrary.org/search?q='%s'",
		"https://z-lib.id/s?q='%s'",
	}
	wg := sync.WaitGroup{}
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			sugar.NewExecService().RunScriptShell(
				fmt.Sprintf(
					// "chrome --proxy-server=%s %s",
					"qutebrowser --set content.proxy %s %s",
					ProxyServer,
					fmt.Sprintf(url, content),
				),
			)
		}(url)
	}
	wg.Wait()
}

func SearchVideosOnline() {
	content, err := sugar.GetInput("search videos online: ")
	if err != nil {
		sugar.Notify(err)
		return
	}
	urls := []string{
		"https://search.bilibili.com/all?keyword='%s'",
		"https://www.youtube.com/results?search_query='%s'",
	}
	wg := sync.WaitGroup{}
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			sugar.NewExecService().RunScriptShell(
				fmt.Sprintf(
					// "chrome --proxy-server=%s %s",
					"qutebrowser --set content.proxy %s %s",
					ProxyServer,
					fmt.Sprintf(url, content),
				),
			)
		}(url)
	}
	wg.Wait()
}

func NoteDiary() {
	dateStr := time.Now().Format(time.DateOnly)
	fileDir := path.Join(os.Getenv("HOME"), "github", "obsidian", "diary")
	filePath := path.Join(fileDir, dateStr+".md")
	if !sugar.IsDirExists(fileDir) {
		if err := os.Mkdir(fileDir, 0o755); err != nil {
			sugar.Notify(err)
			return
		}
	}
	if !sugar.IsFileExists(filePath) {
		f, err := os.Create(filePath)
		if err != nil {
			sugar.Notify(err)
			return
		}
		fmt.Fprintf(f, "\n### %s\n\n", dateStr)
		f.Close()
	}
	sugar.Toggle(fmt.Sprintf("st -e nvim +$ '%s'", filePath))
}

func NoteTimeline() {
	t := time.Now()
	dateStr := t.Format(time.DateOnly)
	datetimeStr := t.Format(time.DateTime)
	fileDir := path.Join(os.Getenv("HOME"), "github", "obsidian", "timeline")
	filePath := path.Join(fileDir, dateStr+".md")
	if !sugar.IsDirExists(fileDir) {
		if err := os.Mkdir(fileDir, 0o755); err != nil {
			sugar.Notify(err)
			return
		}
	}
	if !sugar.IsFileExists(filePath) {
		f, err := os.Create(filePath)
		if err != nil {
			sugar.Notify(err)
			return
		}
		fmt.Fprintf(f, "\n## %s\n\n", dateStr)
		f.Close()
	}
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0o755)
	if err != nil {
		sugar.Notify(err)
		return
	}
	fmt.Fprintf(f, "\n### %s\n\n", datetimeStr)
	f.Close()
	sugar.Toggle(
		fmt.Sprintf("st -e nvim +$ '%s'", filePath),
	)
}

func NoteFlashCard() {
	t := time.Now()
	fileDir := path.Join(os.Getenv("HOME"), "github", "obsidian", "flash-card")
	filePath := path.Join(
		fileDir,
		t.Format("2006-01-02.15.04.05.000000000")+".md",
	)
	if !sugar.IsDirExists(fileDir) {
		if err := os.Mkdir(fileDir, 0o755); err != nil {
			sugar.Notify(err)
			return
		}
	}
	if !sugar.IsFileExists(filePath) {
		f, err := os.Create(filePath)
		if err != nil {
			sugar.Notify(err)
			return
		}
		fmt.Fprintf(f, "### %s\n\n", t.Format(time.DateTime))
		f.Close()
	}
	_, _, err := sugar.NewExecService().RunScriptShell(fmt.Sprintf("st -e nvim +$ '%s'", filePath))
	if err != nil {
		sugar.Notify(err)
	}
}

func HandleCopied() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)
	content := strings.TrimSpace(string(text))
	switch {
	case sugar.Exists(content) && sugar.IsFile(content):
		sugar.Lazy("open", content)
		return
	case sugar.IsUrl(content):
		ChromeOpenUrl("--proxy-server="+ProxyServer, content)()
		return
	default:
		SearchFromWeb(content)
	}
}

func WifiConnect() {
	cmd := "nmcli device wifi list|sed '1d'|sed '/--/ d'|awk '{print $2}'|sort|uniq"
	stdout, _, err := sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	cmd = fmt.Sprintf("echo '%s'|dmenu -p 'connect to wifi'", stdout)
	stdout, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	essid := strings.TrimSpace(stdout)
	if essid == "" {
		return
	}
	cmd = "dmenu < /dev/null -p 'password'"
	stdout, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	password := strings.TrimSpace(stdout)
	cmd = fmt.Sprintf("nmcli device wifi connect %s password %s", essid, password)
	_, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	sugar.Notify("wifi connect success")
}

func JumpToCodeFromLog() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	textbyte := clipboard.Read(clipboard.FmtText)
	text := string(textbyte)
	regex := `(?P<filepath>/[^\:]+):(?P<row>\d+)\s+`
	r := regexp.MustCompile(regex)
	match := r.FindStringSubmatch(text)
	if len(match) < 3 {
		sugar.Notify("not match")
		return
	}
	filepath := match[1]
	row := match[2]
	cmd := fmt.Sprintf(
		"st -e nvim +%s %s",
		row,
		filepath,
	)
	_, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
}

func SshTo() {
	mysshListFilePath := path.Join(os.Getenv("HOME"), ".ssh/my.ssh.list")
	if !sugar.IsFileExists(mysshListFilePath) {
		f, err := os.Create(mysshListFilePath)
		if err != nil {
			sugar.Notify(err)
			return
		}
		f.Close()
	}

	// read from to ~/.ssh/my.ssh.list
	b, err := os.ReadFile(mysshListFilePath)
	if err != nil {
		sugar.Notify(err)
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
	knownHosts, err := sugar.GetKnownHosts()
	if err != nil {
		sugar.Notify(err)
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
	chioce, err := sugar.Choose("ssh to: ", promptList)
	if err != nil {
		sugar.Notify(err)
		return
	}
	chioce = strings.TrimSpace(chioce)
	slice := regexp.MustCompile(`[ \r\t\s]+`).Split(chioce, -1)

	switch {
	case len(slice) > 3:
		host := strings.TrimSpace(slice[0])
		user := strings.TrimSpace(slice[1])
		password := strings.TrimSpace(slice[2])
		err = sugar.SSH(host, 22, user, password)
		if err != nil {
			sugar.Notify(err)
			return
		}
		return
	default:
		host := strings.TrimSpace(slice[0])
		user, err := sugar.GetInput("user: ")
		if err != nil {
			sugar.Notify(err)
			return
		}
		password, err := sugar.GetInput("password: ")
		if err != nil {
			sugar.Notify(err)
			return
		}
		desc, err := sugar.GetInput("desc: ")
		if err != nil {
			sugar.Notify(err)
			return
		}

		err = sugar.SSH(host, 22, user, password)
		if err != nil {
			sugar.Notify(err)
			return
		}

		// append to ~/.ssh/my.ssh.list
		file, err := os.OpenFile(
			mysshListFilePath,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0o666,
		)
		if err != nil {
			sugar.Notify(err)
		}
		defer file.Close()
		writer := bufio.NewWriter(file)
		fmt.Fprintf(writer, "%-20s %-20s %-20s # %s\r\n", host, user, password, desc)
		writer.Flush()
	}
}

/*
func OCR() {
	// --- screenshot
	_, stderr, err := sugar.NewExecService().RunScriptShell("flameshot gui -p /tmp")
	if err != nil {
		sugar.Notify(err)
		return
	}
	// stderr: flameshot: info: Capture saved as /tmp/20240120132854.png
	slice := strings.Split(strings.TrimSpace(stderr), " ")
	if len(slice) == 0 {
		sugar.Notify(
			fmt.Errorf("ocr failed: flameshot output %s", stderr),
		)
		return
	}
	filepath := slice[len(slice)-1]

	// --- read screenshot
	b, err := os.ReadFile(filepath)
	if err != nil {
		sugar.Notify(err)
		return
	}
	base64str := base64.StdEncoding.EncodeToString(b)
	defer os.Remove(filepath)

	// --- ocr
	// 实例化一个认证对象，入参需要传入腾讯云账户 SecretId 和 SecretKey，此处还需注意密钥对的保密
	// 代码泄露可能会导致 SecretId 和 SecretKey 泄露，并威胁账号下所有资源的安全性。
	// 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取
	b, err = os.ReadFile(
		path.Join(os.Getenv("HOME"), TencentApiSecretKey),
	)
	if err != nil {
		sugar.Notify(err)
		return
	}
	slice = strings.Split(strings.TrimSpace(string(b)), " ")
	if len(slice) != 2 {
		sugar.Notify(fmt.Errorf("ocr failed: read secret key failed"))
		return
	}
	secretId := strings.TrimSpace(strings.TrimPrefix(slice[0], "SecretId:"))
	secretKey := strings.TrimSpace(strings.TrimPrefix(slice[1], "SecretKey:"))

	credential := common.NewCredential(
		secretId,
		secretKey,
	)
	cp := profile.NewClientProfile()
	cp.HttpProfile.Endpoint = "ocr.tencentcloudapi.com"
	client, err := ocr.NewClient(credential, "ap-shanghai", cp)
	if err != nil {
		sugar.Notify(err)
		return
	}
	request := ocr.NewGeneralBasicOCRRequest()
	request.ImageBase64 = common.StringPtr(base64str)
	response, err := client.GeneralBasicOCR(request)
	if err != nil {
		sugar.Notify(err)
		return
	}

	// --- simple format
	doc := [][]*ocr.TextDetection{}
	row := []*ocr.TextDetection{}
	var eachLineFirstItem *ocr.TextDetection
	for _, item := range response.Response.TextDetections {
		if len(row) == 0 {
			eachLineFirstItem = item
			row = append(row, item)
			continue
		}
		if *eachLineFirstItem.Polygon[0].X-int64(3) <= *item.Polygon[0].X && *item.Polygon[0].X <= *eachLineFirstItem.Polygon[0].X+int64(3) {
			row = append(row, item)
		} else {
			doc = append(doc, row)
			row = []*ocr.TextDetection{}
			row = append(row, item)
			eachLineFirstItem = item
		}
	}
	docstr := ""
	for _, row := range doc {
		rowstr := ""
		for _, item := range row {
			rowstr += *item.DetectedText
		}
		docstr += rowstr + "\n"
	}
	sugar.Notify(fmt.Sprintf("ocr result:\n%s", docstr))

	// write to clipboard
	err = clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	changed := clipboard.Write(clipboard.FmtText, []byte(docstr))
	select {
	case <-changed:
		sugar.Notify("previous clipboard expired")
	}
}
*/
