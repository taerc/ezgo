package gitlab

import (
	"github.com/gin-gonic/gin"
	"github.com/taerc/ezgo"
	"github.com/taerc/ezgo/conf"
	"github.com/taerc/ezgo/notify"
)

type TagEventsLoad struct {
	*ezgo.GinFlow `json:"-,omitempty"`
	ObjectKind    string `json:"object_kind"`
	EventName     string `json:"event_name"`
	Before        string `json:"before"`
	After         string `json:"after"`
	Ref           string `json:"ref"`
	Message       string `json:"message"`
	CheckoutSha   string `json:"checkout_sha"`
	UserID        int    `json:"user_id"`
	UserName      string `json:"user_name"`
	UserAvatar    string `json:"user_avatar"`
	ProjectID     int    `json:"project_id"`
	Project       struct {
		ID                int         `json:"id"`
		Name              string      `json:"name"`
		Description       string      `json:"description"`
		WebURL            string      `json:"web_url"`
		AvatarURL         interface{} `json:"avatar_url"`
		GitSSHURL         string      `json:"git_ssh_url"`
		GitHTTPURL        string      `json:"git_http_url"`
		Namespace         string      `json:"namespace"`
		VisibilityLevel   int         `json:"visibility_level"`
		PathWithNamespace string      `json:"path_with_namespace"`
		DefaultBranch     string      `json:"default_branch"`
		Homepage          string      `json:"homepage"`
		URL               string      `json:"url"`
		SSHURL            string      `json:"ssh_url"`
		HTTPURL           string      `json:"http_url"`
	} `json:"project"`
	Repository struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		Description     string `json:"description"`
		Homepage        string `json:"homepage"`
		GitHTTPURL      string `json:"git_http_url"`
		GitSSHURL       string `json:"git_ssh_url"`
		VisibilityLevel int    `json:"visibility_level"`
	} `json:"repository"`
	Commits           []interface{} `json:"commits"`
	TotalCommitsCount int           `json:"total_commits_count"`
}

func (tel *TagEventsLoad) Proc(ctx *gin.Context) {

	if e := tel.BindJson(ctx, tel); e != ezgo.Success {
		return
	}

	if tel.ObjectKind == "tag_push" && conf.Config.Token != "" && conf.Config.Secret != "" {

		sn := &notify.SimpleNotice{}
		sn.Title = "发布"
		sn.Project = tel.Project.Name
		sn.Author = tel.UserName
		sn.Tag = tel.Ref
		sn.Items = ezgo.StringSplits(tel.Message, []string{",", "，"})
		notify.HookSendMarkdownDingGroupWithConf(sn, conf.Config.Token, conf.Config.Secret)
	}
	return
}
