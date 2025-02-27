package dd

import (
	"bytes"
	"fmt"
	"text/template"
)

type SimpleNotice struct {
	Title    string   `json:"title"`
	ImageUrl string   `json:"imgeurl"`
	UrlName  string   `json:"urlname"`
	Url      string   `json:"url"`
	Project  string   `json:"project"`
	Tag      string   `json:"tag"`
	Author   string   `json:"author"`
	Items    []string `json:"items"`
}

func (sn *SimpleNotice) Append(item string) {
	sn.Items = append(sn.Items, item)
}
func (sn *SimpleNotice) ToString() string {
	tplText := `
**项目** : {{.Project}}

{{if .ImageUrl}}![image]({{.ImageUrl}}) {{end}}

**标题**: {{.Title}}

{{if .Tag}}**标签**: {{.Tag}}{{end}}

**作者**: {{.Author}}

**详情**:

{{if .Url}}[链接]({{.Url}}) {{end}}

{{- range $i, $e := .Items }}
* {{$e}}
{{- end }}
`
	tpl, err := template.New("note").Parse(tplText)
	if err != nil {
		fmt.Printf("failed parse tpltext,err:%s\n", err.Error())
		return ""
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, sn)
	if err != nil {
		fmt.Printf("failed execute tpltext,err:%s\n", err.Error())
		return ""
	}
	return buf.String()
}

func HookSendMarkdownDingGroupWithConf(notice *SimpleNotice, token string, secret string) {

	if text := notice.ToString(); len(text) > 0 {
		var receiver Robot
		receiver.AccessToken = token
		receiver.Secret = secret
		webHookUrl := receiver.Signature()
		params := receiver.SendMarkdown(notice.Title, text, []string{}, []string{}, false)
		SendRequest(webHookUrl, params)
	}
}
