package main

import (
	"bytes"
	"flag"
	"fmt"
	"text/template"

	"github.com/taerc/ezgo/dd"
)

// 基于钉钉的简单任务同步助手
// 文本消息和静态链接

var AccessToken string
var AccessSecret string
var NoteText string
var NoteLink string

// https://oapi.dingtalk.com/robot/send?access_token=6587d40230371eb38fff496113ffdc4500b0100dd7208ef1e779313573f3c430

func init() {
	flag.StringVar(&AccessToken, "token", "6587d40230371eb38fff496113ffdc4500b0100dd7208ef1e779313573f3c430", "path of configure file.")
	flag.StringVar(&AccessSecret, "secret", "SEC770b9531b28ba60150c930e964865fbd8d8649e1e9409e65aeb6c4e00aa06bf8", "path of configure file.")
	flag.StringVar(&NoteText, "text", "", "path of configure file.")
	flag.StringVar(&NoteLink, "link", "", "path of configure file.")
	flag.Parse()
}

type NoteMe struct {
	Text string
	Link string
}

func (nm *NoteMe) DingMessage() string {

	ddText := `
**TIPS**:

{{.Text}}

{{if .Link}}[在这里]({{.Link}}) {{end}}

`
	tpl, err := template.New("note").Parse(ddText)
	if err != nil {
		fmt.Printf("failed parse tpltext,err:%s\n", err.Error())
		return ""
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, nm)
	if err != nil {
		fmt.Printf("failed execute tpltext,err:%s\n", err.Error())
		return ""
	}

	text := buf.String()

	if AccessToken != "" && AccessSecret != "" {
		var receiver dd.Robot
		receiver.AccessToken = AccessToken
		receiver.Secret = AccessSecret
		sign := receiver.Signature()
		params := receiver.SendMarkdown("钉钉助理", text, []string{}, []string{}, false)
		dd.SendRequest(sign, params)
	}
	return text
}

func main() {

	if AccessToken == "" || AccessSecret == "" || NoteText == "" {
		fmt.Println("invalid input args!")
		return
	}

	nm := NoteMe{
		Text: NoteText,
		Link: NoteLink,
	}
	nm.DingMessage()

}
