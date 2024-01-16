package plugins

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"strings"

	"cmds/sugar"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
	"golang.design/x/clipboard"
)

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
