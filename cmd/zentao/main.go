package main

import (
	"bytes"
	"flag"
	"fmt"
	"strconv"
	"text/template"
	"time"

	"github.com/taerc/ezgo/dd"
	"github.com/taerc/go-zentao/v21/zentao"
)

const (
	UserName = "wangfangming"
	Passwd   = "aship@2021"
	BaseUrl  = "http://172.16.10.21:9000/zentao"
)

const (
	StatusKey      = "status"
	StatusActive   = "active"
	StatusResolved = "resolved"
)

type TargetDesc struct {
	ID   int
	Name string
}

type ProjectWeekly struct {
	ID   int
	Name string // 关注的项目信息

	// 关注的  bug 工作信息
	WeekNewBugs    int
	WeekCloseBugs  int
	ToBeClosedBugs int
	ClosedBugs     int
	// bug list

	// 关注的Story
	WeekNewStories    int
	WeekCloseSories   int
	ToBeClosedStories int
	ClosedStories     int

	BugList []Event

	// 报告文本模版
	ReportItem []string
}

type StaffWeekly struct {
	NewBugs    int
	CloseBugs  int
	ToBeSolved int
}

type Event struct {
	ID         int
	Title      string
	AssignName string
}

func NewProject(id int, name string) *ProjectWeekly {

	pk := new(ProjectWeekly)
	pk.ID = id
	pk.Name = name
	pk.BugList = make([]Event, 0)
	pk.ReportItem = make([]string, 0)
	return pk
}

func (pk *ProjectWeekly) Parse(zt *zentao.Client) {
	// page size

	// 测试 bug
	total := 50
	limit := 50
	page := 1
	for total > 0 && limit > 0 && page > 0 && limit*page < (total/limit+1)*limit {
		p1, _, err := zt.Bugs.ListByProjects(int64(pk.ID), zentao.ListOptions{
			Page:  page,
			Limit: limit,
		})
		if err == nil {
			for _, bug := range p1.Bugs {
				pk.incBugWeekNew(bug)
				pk.incBugWeekClosed(bug)
			}
		}
		total, limit, page = p1.Total, p1.Limit, p1.Page
		page += 1
	}

	// 用户故事
	total = 50
	limit = 50
	page = 1
	for total > 0 && limit > 0 && page > 0 && limit*page < (total/limit+1)*limit {
		st, _, err := zt.Stories.ProjectsList(pk.ID, strconv.Itoa(limit), strconv.Itoa(page), "all")
		if err == nil {
			for _, s := range st.Stories {
				pk.incStoryWeekNew(s)
				pk.incStoryWeekClosed(s)
			}
		}
		total, limit, page = st.Total, st.Limit, st.Page
		page += 1
	}

	pk.report()

}

func (pk ProjectWeekly) weekCompare(dt string) bool {

	if dt == "" {
		return false
	}

	parsedTime, _ := time.Parse(time.RFC3339, dt)
	y1, w1 := parsedTime.ISOWeek()
	y2, w2 := time.Now().ISOWeek()

	if y1 == y2 && w1 == w2 {
		return true
	}

	return false
}
func (pk *ProjectWeekly) incBugWeekNew(bug zentao.BugBody) {
	if bug.Status == StatusActive {
		ast, _ := (bug.AssignedTo).(map[string]interface{})
		realname := ""
		if ast["realname"] != nil {
			realname = ast["realname"].(string)
		}
		pk.BugList = append(pk.BugList, Event{ID: bug.ID, Title: bug.Title, AssignName: realname})
		pk.ToBeClosedBugs++
	}
	if bug.Status == StatusActive && pk.weekCompare(bug.OpenedDate) {
		pk.WeekNewBugs++
	}
}

func (pk *ProjectWeekly) incBugWeekClosed(bug zentao.BugBody) {
	if bug.Status == StatusResolved {
		pk.ClosedBugs++
	}

	if bug.Status == StatusResolved && bug.ResolvedDate != "" && pk.weekCompare((bug.ResolvedDate).(string)) {
		pk.WeekCloseBugs++
	}
}

func (pk *ProjectWeekly) incStoryWeekNew(st zentao.StoriesBody) {

	if st.Status == zentao.StatusActive {
		pk.ToBeClosedStories++
	}
	if st.Status == zentao.StatusActive && pk.weekCompare(st.Openeddate) {
		pk.WeekNewStories++
	}

}

func (pk *ProjectWeekly) incStoryWeekClosed(st zentao.StoriesBody) {

	if st.Status == zentao.StatusClosed {
		pk.ClosedStories++
	}
	if st.Status == zentao.StatusClosed && pk.weekCompare(st.Closeddate) {
		pk.WeekNewStories++
	}

}

func (pk *ProjectWeekly) report() {

	rpt := fmt.Sprintf("[%v-%v] 本周新增BUG: %v 本周解决: %v  待解决: %v 共解决: %v ",
		pk.Name, pk.ID,
		pk.WeekNewBugs, pk.WeekCloseBugs, pk.ToBeClosedBugs, pk.ClosedBugs)
	rpt += fmt.Sprintf("本周新增STORY: %v 本周完成: %v 待完成: %v 共完成: %v",
		pk.WeekNewStories, pk.WeekCloseSories, pk.ToBeClosedStories, pk.ClosedStories)
	pk.ReportItem = append(pk.ReportItem, rpt)
}

// operation between project
func (pk *ProjectWeekly) Add(pw *ProjectWeekly) {

	pk.WeekNewBugs += pw.WeekNewBugs
	pk.WeekCloseBugs += pw.WeekCloseBugs
	pk.ToBeClosedBugs += pw.ToBeClosedBugs
	pk.ClosedBugs += pw.ClosedBugs

	pk.WeekNewStories += pw.WeekNewStories
	pk.WeekCloseSories += pw.WeekCloseSories
	pk.ToBeClosedStories += pw.ToBeClosedStories
	pk.ClosedStories += pw.ClosedStories
}

func (pk *ProjectWeekly) AppendReport(pw *ProjectWeekly) {
	pk.ReportItem = append(pk.ReportItem, pw.ReportItem...)
}

func (pk *ProjectWeekly) MarkdownReport() string {
	tplText := `
**重点项目需求与BUG报告** 

{{- range .}}
* {{.}}
{{- end }}

`

	// 解析模板字符串为模板对象
	tmpl, err := template.New("example").Parse(tplText)
	if err != nil {
		panic(err)
	}

	// 将数据传递给模板并执行，将结果输出到标准输出
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, pk.ReportItem)
	if err != nil {
		fmt.Printf("failed execute tpltext,err:%s\n", err.Error())
		return ""
	}
	return buf.String()
}

var WorkMode string
var AccessToken string
var AccessSecret string

func init() {
	flag.StringVar(&WorkMode, "m", "all", "方寸所有项目监视, all, platform")
	flag.StringVar(&AccessToken, "token", "6fb261b244f9ef169b001cbe967b210576607bcee0873885436f5cfe54581d36", "token of ding")
	flag.StringVar(&AccessSecret, "secret", "SEC67141ab5e29e2dce7196d46a7e0dedf8a7bb96887e880137d8ad6817e64bdc8b", "sec of ding")
	flag.Parse()
}

// https://oapi.dingtalk.com/robot/send?access_token=6fb261b244f9ef169b001cbe967b210576607bcee0873885436f5cfe54581d36

func main() {
	zt, err := zentao.NewBasicAuthClient(
		UserName,
		Passwd,
		zentao.WithBaseURL(BaseUrl),
		zentao.WithoutProxy(),
	)
	if err != nil {
		panic(err)
	}
	comProjects := NewProject(0, "中科方寸公司")
	reportText := ""
	if WorkMode == "all" {
		pros, _, err := zt.Projects.List("9999", "1")
		if err != nil {
			panic(err)
		}

		for _, pro := range pros.Projects {
			pk := NewProject(pro.ID, pro.Name)
			pk.Parse(zt)
			comProjects.Add(pk)
		}
		comProjects.report()
		reportText = comProjects.MarkdownReport()
	}

	if WorkMode == "platform" {

		platforms := []TargetDesc{
			{
				247, "新疆信产配网平台",
			},
			{
				442, "新疆信产25年配网无人机平台二期",
			},
			{
				114, "蒙东电网无人机自适应项目",
			},
			{
				370, "蒙东电网无人机工程验收",
			},
			{
				444, "蒙东赤峰专网深化应用",
			},
			{
				447, "蒙东电网无人机自适应项目-二期",
			},
			{
				152, "冀北配网平台项目",
			},
			{
				437, "吉林省配网无人机平台项目",
			},
			{
				430, "配网无人机管理平台2.0项目",
			},
		}

		comProjects := NewProject(0, "平台组项目")
		for _, pro := range platforms {
			pk := NewProject(pro.ID, pro.Name)
			pk.Parse(zt)
			comProjects.Add(pk)
			comProjects.AppendReport(pk)
		}
		comProjects.report()
		reportText = comProjects.MarkdownReport()
	}

	// send message to dingding
	if AccessToken != "" && AccessSecret != "" {
		var receiver dd.Robot
		receiver.AccessToken = AccessToken
		receiver.Secret = AccessSecret
		sign := receiver.Signature()
		params := receiver.SendMarkdown("钉助理", reportText, []string{}, []string{}, false)
		dd.SendRequest(sign, params)
	}
}
