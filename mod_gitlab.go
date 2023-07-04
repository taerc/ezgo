package ezgo

import (
	"github.com/gin-gonic/gin"
	"github.com/taerc/ezgo/conf"
	"sync"
	"time"
)

type pushEventPayload struct {
	*GinFlow     `json:"-,omitempty"`
	ObjectKind   string      `json:"object_kind"`
	EventName    string      `json:"event_name"`
	Before       string      `json:"before"`
	After        string      `json:"after"`
	Ref          string      `json:"ref"`
	CheckoutSha  string      `json:"checkout_sha"`
	Message      interface{} `json:"message"`
	UserID       int         `json:"user_id"`
	UserName     string      `json:"user_name"`
	UserUsername string      `json:"user_username"`
	UserEmail    string      `json:"user_email"`
	UserAvatar   string      `json:"user_avatar"`
	ProjectID    int         `json:"project_id"`
	Project      struct {
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
		CiConfigPath      string      `json:"ci_config_path"`
		Homepage          string      `json:"homepage"`
		URL               string      `json:"url"`
		SSHURL            string      `json:"ssh_url"`
		HTTPURL           string      `json:"http_url"`
	} `json:"project"`
	Commits []struct {
		ID        string    `json:"id"`
		Message   string    `json:"message"`
		Title     string    `json:"title"`
		Timestamp time.Time `json:"timestamp"`
		URL       string    `json:"url"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Added    []interface{} `json:"added"`
		Modified []string      `json:"modified"`
		Removed  []interface{} `json:"removed"`
	} `json:"commits"`
	TotalCommitsCount int `json:"total_commits_count"`
	PushOptions       struct {
	} `json:"push_options"`
	Repository struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		Description     string `json:"description"`
		Homepage        string `json:"homepage"`
		GitHTTPURL      string `json:"git_http_url"`
		GitSSHURL       string `json:"git_ssh_url"`
		VisibilityLevel int    `json:"visibility_level"`
	} `json:"repository"`
}

func (pep *pushEventPayload) Proc(ctx *gin.Context) {

	if e := pep.BindJson(ctx, pep); e != Success {
		return
	}

	if pep.ObjectKind == "push" && conf.Config.Token != "" && conf.Config.Secret != "" {

		sn := &SimpleNotice{}
		sn.Title = "更新"
		sn.Project = pep.Project.Name
		sn.Author = pep.UserName
		for _, msg := range pep.Commits {
			sn.Append(msg.Message)
		}

		HookSendMarkdownDingGroupWithConf(sn, conf.Config.Token, conf.Config.Secret)
	}
	return

}

type tagEventsLoad struct {
	*GinFlow    `json:"-,omitempty"`
	ObjectKind  string `json:"object_kind"`
	EventName   string `json:"event_name"`
	Before      string `json:"before"`
	After       string `json:"after"`
	Ref         string `json:"ref"`
	Message     string `json:"message"`
	CheckoutSha string `json:"checkout_sha"`
	UserID      int    `json:"user_id"`
	UserName    string `json:"user_name"`
	UserAvatar  string `json:"user_avatar"`
	ProjectID   int    `json:"project_id"`
	Project     struct {
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

func (tel *tagEventsLoad) Proc(ctx *gin.Context) {

	if e := tel.BindJson(ctx, tel); e != Success {
		return
	}

	if tel.ObjectKind == "tag_push" && conf.Config.Token != "" && conf.Config.Secret != "" {

		sn := &SimpleNotice{}
		sn.Title = "发布"
		sn.Project = tel.Project.Name
		sn.Author = tel.UserName
		sn.Tag = tel.Ref
		sn.Items = StringSplits(tel.Message, []string{",", "，"})
		HookSendMarkdownDingGroupWithConf(sn, conf.Config.Token, conf.Config.Secret)
	}
	return
}

func WithModuleGitLab() func(wg *sync.WaitGroup) {
	return func(wg *sync.WaitGroup) {
		wg.Done()
		route := Group("/gitlab/hook/event/")
		SetPostProc(route, "/push", &pushEventPayload{})
		SetPostProc(route, "/push_tag", &tagEventsLoad{})
		Info(nil, M, "Load finished!")
	}
}
