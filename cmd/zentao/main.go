package main

import (
	"bytes"
	"flag"
	"fmt"
	"sort"
	"strconv"
	"text/template"
	"time"

	"github.com/taerc/ezgo/dd"
	ezgo "github.com/taerc/ezgo/pkg"
	"github.com/taerc/go-zentao/v21/zentao"
)

const (
	UserName = "wangfangming"
	Passwd   = "aship@2021"
	BaseUrl  = "http://172.16.10.21:9000/zentao"
)

type TargetDesc struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProjectWeekly struct {
	ID   int
	Name string // 关注的项目信息

	// 关注的  bug 工作信息
	WeekNewBugs    int
	WeekResolvedBugs  int
	WeekCloseBugs  int
	ToBeClosedBugs int
	ClosedBugs     int
	// bug list

	// 关注的Story
	WeekNewStories    int
	WeekResolvedStories int
	WeekCloseStories  int
	ToBeClosedStories int
	ClosedStories     int

	BugList []Event
	BugToDoMetric map[string]int

	StoryList []Event
	StoryToDoMetric map[string]int

	// 报告文本模版
	ReportText  string
	SubProjects []*ProjectWeekly
}

type StaffWeekly struct {
	NewBugs    int
	CloseBugs  int
	ToBeSolved int
}

type kv struct {
	Key   string
	Value int
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
	pk.SubProjects = make([]*ProjectWeekly, 0)
	pk.BugToDoMetric = make(map[string]int,256)
	pk.StoryList = make([]Event, 0)
	pk.StoryToDoMetric = make(map[string]int, 256)
			
	return pk
}

func (pk *ProjectWeekly) Parse(zt *zentao.Client) {
	// page size

	// 测试 bug
	total := 50
	limit := 50
	page := 1
	for total > 0 && limit > 0 && page > 0 && limit*page <= (total/limit+1)*limit {
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
	for total > 0 && limit > 0 && page > 0 && limit*page <= (total/limit+1)*limit {
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

func (pk ProjectWeekly) isCurrentWeek(dt string) bool {

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
	if bug.Status != zentao.ResolvedBugStatus && bug.Status!= zentao.ClosedBugStatus {
		realname := ""
		realname = bug.AssignedTo.(string)
		pk.BugList = append(pk.BugList, Event{ID: bug.ID, Title: bug.Title, AssignName: realname})
		pk.ToBeClosedBugs++
	}
	if bug.OpenedDate != "" && pk.isCurrentWeek(bug.OpenedDate) {
		pk.WeekNewBugs++
	}
}

func (pk *ProjectWeekly) incBugWeekClosed(bug zentao.BugBody) {
	if bug.Status == zentao.ResolvedBugStatus || bug.Status == zentao.ClosedBugStatus{
		pk.ClosedBugs++
	} 

	if bug.Status == zentao.ResolvedBugStatus || bug.Status == zentao.ClosedBugStatus {
		if resolvedDate, ok := bug.ResolvedDate.(string); ok && resolvedDate != "" {
			if pk.isCurrentWeek(resolvedDate) {
				pk.WeekResolvedBugs++
			}
		}
	} 
	if bug.Status == zentao.ClosedBugStatus && bug.ClosedDate != "" && pk.isCurrentWeek(bug.ClosedDate) {
		pk.WeekCloseBugs++
	}
}

func (pk *ProjectWeekly) incStoryWeekNew(st zentao.StoriesBody) {

	if st.Status == zentao.StatusActive {
		pk.ToBeClosedStories++

		realname := ""
		if st.Assignedto != "" {
			if assign, ok := st.Assignedto.(map[string]interface{}); ok {
				realname = assign["realname"].(string)
			}
			pk.StoryList = append(pk.StoryList, Event{ID: st.ID, Title: st.Title, AssignName: realname})
		}
	}
	if st.Status == zentao.StatusActive && pk.isCurrentWeek(st.Openeddate) {
		pk.WeekNewStories++
	}

}

func (pk *ProjectWeekly) incStoryWeekClosed(st zentao.StoriesBody) {

	if st.Status == zentao.StatusClosed {
		pk.ClosedStories++
	}
	if st.Status == zentao.StatusClosed && pk.isCurrentWeek(st.Closeddate) {
		pk.WeekCloseStories++
	}

}

func (pk *ProjectWeekly) report() string {
	tplText := `
**{{.Name}}** 
* 编号 {{.ID}}
* 本周新增BUG {{.WeekNewBugs}}
* 本周解决 {{.WeekResolvedBugs}}
* 本周关闭 {{.WeekCloseBugs}}
* 待解决 {{.ToBeClosedBugs}}
* 共解决 {{.ClosedBugs}}

* 本周新增STORY {{.WeekNewStories}}
* 本周完成 {{.WeekCloseStories}}
* 待完成 {{.ToBeClosedStories}}
* 共完成 {{.ClosedStories}}
`
	// 解析模板字符串为模板对象
	tmpl, err := template.New("example").Parse(tplText)
	if err != nil {
		panic(err)
	}

	// 将数据传递给模板并执行，将结果输出到标准输出
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, pk)
	if err != nil {
		fmt.Printf("failed execute tpltext,err:%s\n", err.Error())
		return ""
	}
	pk.ReportText = buf.String()
	return pk.ReportText
}

// operation between project
func (pk *ProjectWeekly) Add(pw *ProjectWeekly) {

	pk.WeekNewBugs += pw.WeekNewBugs
	pk.WeekResolvedBugs += pw.WeekResolvedBugs
	pk.WeekCloseBugs += pw.WeekCloseBugs
	pk.ToBeClosedBugs += pw.ToBeClosedBugs
	pk.ClosedBugs += pw.ClosedBugs

	pk.WeekNewStories += pw.WeekNewStories
	pk.WeekCloseStories += pw.WeekCloseStories
	pk.ToBeClosedStories += pw.ToBeClosedStories
	pk.ClosedStories += pw.ClosedStories
	pk.BugList = append(pk.BugList, pw.BugList...)
	pk.StoryList = append(pk.StoryList, pw.StoryList...)
}

func (pk *ProjectWeekly) AppendReport(pw *ProjectWeekly) {
	pk.SubProjects = append(pk.SubProjects, pw)
}

// Merge Sub Projects Reportext
func (pk *ProjectWeekly) MergeSubReportText() string {
	tplText := `
*重点项目需求与BUG报告* 

{{- range $i, $e := .SubProjects}}
{{$e.ReportText}}
{{- end }}

`
	// 解析模板字符串为模板对象
	tmpl, err := template.New("example").Parse(tplText)
	if err != nil {
		panic(err)
	}

	// 将数据传递给模板并执行，将结果输出到标准输出
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, pk)
	if err != nil {
		fmt.Printf("failed execute tpltext,err:%s\n", err.Error())
		return ""
	}
	return buf.String()
}

func (pk *ProjectWeekly)mapToSortedKV(data map[string]int) []kv {
    var sorted []kv
    for k, v := range data {
        sorted = append(sorted, kv{k, v})
    }
    // 按 Value 降序排序
    sort.Slice(sorted, func(i, j int) bool {
        return sorted[i].Value > sorted[j].Value
    })

	return sorted
}

func (pk *ProjectWeekly) ReportBugToDoMetric() string {
    // 统计数据
    for _, bug := range pk.BugList {
        if _, ok := pk.BugToDoMetric[bug.AssignName]; ok {
            pk.BugToDoMetric[bug.AssignName] += 1
        } else {
            pk.BugToDoMetric[bug.AssignName] = 1
        }
    }

    tplText := `
** 待办BUG **
{{- range .}}
{{.Key}} : {{.Value}}
{{- end }}
`
    tmpl, err := template.New("example").Parse(tplText)
    if err != nil {
        panic(err)
    }

	sorted := pk.mapToSortedKV(pk.BugToDoMetric)
    var buf bytes.Buffer
    err = tmpl.Execute(&buf, sorted)
    if err != nil {
        fmt.Printf("failed execute tpltext,err:%s\n", err.Error())
        return ""
    }
    return buf.String()
}

func (pk *ProjectWeekly) ReportStoryToDoMetric() string {
    // 统计数据
    for _, story := range pk.StoryList {
        if _, ok := pk.StoryToDoMetric[story.AssignName]; ok {
            pk.StoryToDoMetric[story.AssignName] += 1
        } else {
            pk.StoryToDoMetric[story.AssignName] = 1
        }
    }
    tplText := `
** 待办Story **
{{- range .}}
{{.Key}} : {{.Value}}
{{- end }}
`
    tmpl, err := template.New("example").Parse(tplText)
    if err != nil {
        panic(err)
    }
	sorted := pk.mapToSortedKV(pk.StoryToDoMetric)
    var buf bytes.Buffer
    err = tmpl.Execute(&buf, sorted)
    if err != nil {
        fmt.Printf("failed execute tpltext,err:%s\n", err.Error())
        return ""
    }
    return buf.String()
}



var WorkMode string
var AccessToken string
var AccessSecret string
var TargetCfgPath string

func init() {
	flag.StringVar(&WorkMode, "m", "all", "模式, all, platform,test")
	flag.StringVar(&AccessToken, "token", "6fb261b244f9ef169b001cbe967b210576607bcee0873885436f5cfe54581d36", "钉钉token")
	flag.StringVar(&AccessSecret, "secret", "SEC67141ab5e29e2dce7196d46a7e0dedf8a7bb96887e880137d8ad6817e64bdc8b", "钉钉sec")
	flag.StringVar(&TargetCfgPath, "p", "target_platform.json", "目标平台清单")
	flag.Parse()
}

func sendMsg(msg string) {
	// send message to dingding
	if AccessToken != "" && AccessSecret != "" && WorkMode != "test"{
		var receiver dd.Robot
		receiver.AccessToken = AccessToken
		receiver.Secret = AccessSecret
		sign := receiver.Signature()
		params := receiver.SendMarkdown("钉助理", msg, []string{}, []string{}, false)
		dd.SendRequest(sign, params)
	}

}

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
		reportText = comProjects.report()

		bugsToDo := comProjects.ReportBugToDoMetric()
		sendMsg(bugsToDo)

		storiesToDo := comProjects.ReportStoryToDoMetric()
		sendMsg(storiesToDo)
	}

	if WorkMode == "platform" {

		platforms := make([]TargetDesc, 0)
		if ezgo.PathExists(TargetCfgPath) {
			ezgo.LoadJson(TargetCfgPath, &platforms)
		}

		comProjects := NewProject(0, "平台组项目")
		for _, pro := range platforms {
			pk := NewProject(pro.ID, pro.Name)
			pk.Parse(zt)
			comProjects.Add(pk)
			comProjects.AppendReport(pk)
		}
		comProjects.report()
		comProjects.AppendReport(comProjects)
		reportText = comProjects.MergeSubReportText()
	}

	if WorkMode == "test" {

			pk := NewProject(152, "冀北项目")
			pk.Parse(zt)
			comProjects.Add(pk)
			reportText = comProjects.report()

			txt := comProjects.ReportBugToDoMetric()
			fmt.Println(txt)
			txt = comProjects.ReportStoryToDoMetric()
			fmt.Println(txt)

	}

	fmt.Println(reportText)
	sendMsg(reportText)

}
