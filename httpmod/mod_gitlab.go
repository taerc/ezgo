package httpmod

import (
	"encoding/gob"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/taerc/ezgo/conf"
	"github.com/taerc/ezgo/dd"
	ezgo "github.com/taerc/ezgo/pkg"
)

var gitOpChan = make(chan string, 100)

type gitlabService struct {
}

type pushEventPayload struct {
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

func (g *gitlabService) Update(ctx *gin.Context) {

	pep := &pushEventPayload{}
	if e := JsonBind(ctx, pep); e != nil {
		ErrorResponse(ctx, e)
		return
	}

	if pep.ObjectKind == "push" && conf.Config.Ding.Token != "" && conf.Config.Ding.Secret != "" {

		sn := &dd.SimpleNotice{}
		sn.Title = "更新"
		sn.Project = pep.Project.Name
		sn.Author = pep.UserName
		for _, msg := range pep.Commits {
			sn.Append(msg.Message)
		}

		gitOpChan <- sn.Author
		dd.HookSendMarkdownDingGroupWithConf(sn, conf.Config.Ding.Token, conf.Config.Ding.Secret)
	}
	return

}

type tagEventsLoad struct {
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

func (g *gitlabService) Publish(ctx *gin.Context) {

	tel := &tagEventsLoad{}

	if e := JsonBind(ctx, tel); e != nil {
		ErrorResponse(ctx, e)
		return
	}

	// Info(ctx, M, fmt.Sprintf("ObjectKind [%s]", tel.ObjectKind))
	if tel.ObjectKind == "push" && conf.Config.Ding.Token != "" && conf.Config.Ding.Secret != "" {
		sn := &dd.SimpleNotice{}
		sn.Title = "发布"
		sn.Project = tel.Project.Name
		sn.Author = tel.UserName
		sn.Tag = tel.Ref
		sn.Items = ezgo.StringSplits(tel.Message, []string{",", "，"})
		gitOpChan <- sn.Author
		dd.HookSendMarkdownDingGroupWithConf(sn, conf.Config.Ding.Token, conf.Config.Ding.Secret)
	}
}

// Merge request
type tagEventsMerge struct {
	ObjectKind string `json:"object_kind"`
	User       struct {
		UserName string `json:"username"`
	} `json:"user"`
	Project struct {
		Name              string `json:"name"`
		Description       string `json:"description"`
		WebURL            string `json:"web_url"`
		GitSSHURL         string `json:"git_ssh_url"`
		GitHTTPURL        string `json:"git_http_url"`
		Namespace         string `json:"namespace"`
		VisibilityLevel   int    `json:"visibility_level"`
		PathWithNamespace string `json:"path_with_namespace"`
		DefaultBranch     string `json:"default_branch"`
		Homepage          string `json:"homepage"`
		URL               string `json:"url"`
		SSHURL            string `json:"ssh_url"`
		HTTPURL           string `json:"http_url"`
	} `json:"project"`
}

func (g *gitlabService) Merge(ctx *gin.Context) {

	tel := &tagEventsMerge{}

	if e := JsonBind(ctx, tel); e != nil {
		ErrorResponse(ctx, e)
		return
	}

	// Info(ctx, M, fmt.Sprintf("ObjectKind [%s]", tel.ObjectKind))
	if tel.ObjectKind == "merge_request" && conf.Config.Ding.Token != "" && conf.Config.Ding.Secret != "" {
		sn := &dd.SimpleNotice{}
		sn.Title = "合并代码"
		sn.Project = tel.Project.Name
		sn.Author = tel.User.UserName
		sn.Tag = "-"
		sn.Items = []string{}
		gitOpChan <- sn.Author
		dd.HookSendMarkdownDingGroupWithConf(sn, conf.Config.Ding.Token, conf.Config.Ding.Secret)
	}
}

func (gs *gitlabService) load(fd string) map[string]int {

	// 打开文件
	file, err := os.Open(fd)
	if err != nil {
		fmt.Println("Error opening file:", err)
		data := make(map[string]int, 256)
		return data
	}
	defer file.Close()

	// 创建 Gob 解码器
	decoder := gob.NewDecoder(file)

	// 创建一个空的 map 对象，用于存储解码后的数据
	var data map[string]int

	// 解码数据
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding data:", err)
		data := make(map[string]int, 256)
		return data
	}
	return data
}

func (gs *gitlabService) filename() string {
	y, d := time.Now().ISOWeek()
	return fmt.Sprintf("gitops-count-%v-%v.gob", y, d)
}

func (gs *gitlabService) save(fd string, data map[string]int) {

	// 创建文件
	file, err := os.Create(fd)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// 创建 Gob 编码器
	encoder := gob.NewEncoder(file)

	// 将 map 编码并写入文件
	err = encoder.Encode(data)
	if err != nil {
		fmt.Println("Error encoding map:", err)
		return
	}
}

func (gs *gitlabService) GitOpsCount(sn <-chan string) {

	fd := gs.filename()
	for name := range sn {
		db := gs.load(fd)

		_, ok := db[name]
		if ok {
			db[name] += 1
		} else {
			db[name] = 1
		}

		fmt.Println(name, db[name])
		gs.save(fd, db)
	}

}

func WithModuleGitLab() func(wg *sync.WaitGroup) {
	return func(wg *sync.WaitGroup) {
		defer wg.Done()
		s := new(gitlabService)
		route := Group("/gitlab/hook/event/")
		POST(route, "/push", s.Update)
		POST(route, "/pushtag", s.Publish)
		POST(route, "/merge", s.Merge)

		go s.GitOpsCount(gitOpChan)
	}
}
