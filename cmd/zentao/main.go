package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"sort"
	"strconv"
	"strings"
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

// 项目描述信息
type TargetDesc struct {
	ID   int    `json:"id"`   // 项目ID
	Name string `json:"name"` // 项目名称
}




// ProjectWeekly 项目周报统计
type ProjectWeekly struct {
	ID   int
	Name string // 关注的项目信息

	// Bug相关统计
	WeekNewBugs      int // 本周新增的bug数量
	WeekResolvedBugs int // 本周解决的bug数量
	WeekCloseBugs    int // 本周关闭的bug数量
	ToBeClosedBugs   int // 待解决的bug数量
	ClosedBugs       int // 已解决的bug总数

	// Story相关统计
	WeekNewStories      int // 本周新增的需求数量
	WeekResolvedStories int // 本周解决的需求数量
	WeekCloseStories    int // 本周关闭的需求数量
	ToBeClosedStories   int // 待完成的需求数量
	ClosedStories       int // 已完成的需求总数

	BugList        []Event          // bug列表
	BugToDoMetric  map[string]int   // bug待办统计，key为指派人，value为待办数量
	StoryList      []Event          // 需求列表
	StoryToDoMetric map[string]int  // 需求待办统计，key为指派人，value为待办数量

	ReportText  string           // 报告文本内容
	SubProjects []*ProjectWeekly // 子项目列表

}

type Users struct {
	UserMap map[string]string // 用户信息 account -> realname
}	

// Event 事件信息（Bug或Story）
type Event struct {
	ID         int    // 事件ID
	Title      string // 事件标题
	AssignName string // 指派人
}

// MapToSortedKV 将map转换为按Value降序排列的KV数组
type KV struct {
	Key   string
	Value int
}

func NewUsers() *Users {
	users := new(Users)
	users.UserMap = make(map[string]string, 256)
	return users	
}

func (u *Users) InitUserMap(zt *zentao.Client) {
		users, _, err := zt.Users.List("1000", "1")
		if err!= nil {
			panic(err)
		}

		for _, user := range users.Users {
			u.UserMap[user.Account] =user.Realname 
		}
}

func (u *Users) getUserRealName(acc string) string {
	if real, ok := u.UserMap[acc]; ok {
		return real
	}
	return  ""
}

// NewProject 创建新的项目周报实例
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

// Parse 解析项目数据，包括Bug和Story的统计
func (pk *ProjectWeekly) Parse(zt *zentao.Client) {

	// page size
	// pk.InitUserMap(zt)

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

// isCurrentWeek 判断给定日期是否在本周
func (pk ProjectWeekly) isCurrentWeek(dt string) bool {

	if dt == "" {
		return false
	}

	parsedTime, _ := time.Parse(time.RFC3339, dt)
	y1, w1 := parsedTime.ISOWeek()
	y2, w2 := time.Now().ISOWeek()

	// 手动指定目标日期
	if TargetDate != "yyyy-mm-dd" {
		targetDt,_ := time.Parse(time.DateOnly, TargetDate)
		y2, w2 = targetDt.ISOWeek()
	}

	if y1 == y2 && w1 == w2 {
		return true
	}

	return false
}
func (pk *ProjectWeekly) incBugWeekNew(bug zentao.BugBody) {
	if bug.Status != zentao.ResolvedBugStatus && bug.Status!= zentao.ClosedBugStatus {
		realname := ""
		realname = bug.AssignedTo.(string)
		pk.BugList = append(pk.BugList, Event{ID: bug.ID, Title: bug.Title, AssignName: users.getUserRealName(realname)})
		pk.ToBeClosedBugs++
	}
	if bug.OpenedDate != "" && pk.isCurrentWeek(bug.OpenedDate) {
		pk.WeekNewBugs++

		fmt.Println("bug:", bug.ID, bug.Title, bug.OpenedDate)
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

	if st.Status != zentao.StatusClosed {
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

func (pk *ProjectWeekly)mapToSortedKV(data map[string]int) []KV {
    var sorted []KV
    for k, v := range data {
        sorted = append(sorted,KV{k, v})
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
* {{.Key}} : {{.Value}}
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
func (pk *ProjectWeekly) ReportBugList() string {
    // 统计数据
    tplText := `
** 待办BUG列表 **
{{- range .}}
* ID:{{.ID}} 指派:{{.AssignName}} 标题:{{.Title}}
{{- end }}
`
    tmpl, err := template.New("example").Parse(tplText)
    if err != nil {
        panic(err)
    }

    sort.Slice(pk.BugList, func(i, j int) bool {
        return pk.BugList[i].AssignName > pk.BugList[j].AssignName
    })

    var buf bytes.Buffer
    err = tmpl.Execute(&buf, pk.BugList)
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
* {{.Key}} : {{.Value}}
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

func (pk *ProjectWeekly) ReportStoryList() string {
    // 统计数据
    tplText := `
** 待做story列表 **
{{- range .}}
* ID:{{.ID}} 指派:{{.AssignName}} 标题:{{.Title}}
{{- end }}
`
    tmpl, err := template.New("example").Parse(tplText)
    if err != nil {
        panic(err)
    }

    sort.Slice(pk.StoryList, func(i, j int) bool {
        return pk.StoryList[i].AssignName > pk.StoryList[j].AssignName
    })

    var buf bytes.Buffer
    err = tmpl.Execute(&buf, pk.StoryList)
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
var TargetDate string

var users *Users

//平台群参数
// SECa8df6b21c98bec8f46f3b6681bf40ae19570092f7a512f912d7614eb8150b740
// https://oapi.dingtalk.com/robot/send?access_token=14c092e8a82ef5197797f1275f043420c0bb83cd5f4e3bfeaf3e03fe798503a9
// 项目技术负责人[普]
// flag.StringVar(&AccessToken, "token", "6fb261b244f9ef169b001cbe967b210576607bcee0873885436f5cfe54581d36", "钉钉token")
// flag.StringVar(&AccessSecret, "secret", "SEC67141ab5e29e2dce7196d46a7e0dedf8a7bb96887e880137d8ad6817e64bdc8b", "钉钉sec")

func init() {
	flag.StringVar(&WorkMode, "m", "all", "模式, todo, all, platform,test")
	flag.StringVar(&AccessToken, "token", "14c092e8a82ef5197797f1275f043420c0bb83cd5f4e3bfeaf3e03fe798503a9", "钉钉token")
	flag.StringVar(&AccessSecret, "secret", "SECa8df6b21c98bec8f46f3b6681bf40ae19570092f7a512f912d7614eb8150b740", "钉钉sec")
	flag.StringVar(&TargetCfgPath, "p", "target_platform.json", "目标平台清单")
	flag.StringVar(&TargetDate, "d", "yyyy-mm-dd", "目标日期,yyyy-mm-dd")
	flag.Parse()

	users = NewUsers()
}

// SplitStringByMaxLineSize 按最大字节数切割字符串，保证每块都是完整的行
func SplitStringByMaxLineSize(str string, maxBytes int) []string {
    var chunks []string
    var buffer bytes.Buffer
    scanner := bufio.NewScanner(strings.NewReader(str))
    
    for scanner.Scan() {
        line := scanner.Text() + "\n"  // 补回被Scanner去除的换行符
        lineBytes := []byte(line)
        
        // 如果当前行单独超过限制
        if len(lineBytes) > maxBytes {
            if buffer.Len() > 0 {
                chunks = append(chunks, buffer.String())
                buffer.Reset()
            }
            chunks = append(chunks, string(lineBytes))
            continue
        }
        
        // 检查添加后是否超限
        if buffer.Len() + len(lineBytes) > maxBytes {
            chunks = append(chunks, buffer.String())
            buffer.Reset()
        }
        
        buffer.Write(lineBytes)
    }
    
    // 添加最后剩余内容
    if buffer.Len() > 0 {
        chunks = append(chunks, buffer.String())
    }
    
    return chunks
}

func sendMsg(msg string) {
	// send message to dingding
	if AccessToken != "" && AccessSecret != "" && WorkMode != "beta" && WorkMode!= "test"{
		var receiver dd.Robot
		receiver.AccessToken = AccessToken
		receiver.Secret = AccessSecret
		sign := receiver.Signature()
		msgs := SplitStringByMaxLineSize(msg, 1024*16)
		for _, m := range msgs {
			params := receiver.SendMarkdown("钉助理", m, []string{}, []string{}, false)
			dd.SendRequest(sign, params)
			dd.SendRequest(sign, params)
		}
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

	// 初始化用户信息
	users.InitUserMap(zt)
	comProjects := NewProject(0, "中科方寸公司")
	reportText := ""

	if WorkMode == "todo" {

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

		bugsToDo := comProjects.ReportBugList()
		fmt.Println(bugsToDo)
		sendMsg(bugsToDo)

		storiesToDo := comProjects.ReportStoryList()
		fmt.Println(storiesToDo)
		sendMsg(storiesToDo)

	}
	if WorkMode == "all" || WorkMode == "beta"  {
		pros, _, err := zt.Projects.List("9999", "1")
		if err != nil {
			panic(err)
		}

		for _, pro := range pros.Projects {
			fmt.Println(pro.ID, pro.Name)
			pk := NewProject(pro.ID, pro.Name)
			pk.Parse(zt)
			comProjects.Add(pk)
		}
		reportText = comProjects.report()

		bugsToDo := comProjects.ReportBugToDoMetric()
		fmt.Println(bugsToDo)
		sendMsg(bugsToDo)
		bugsToDo = comProjects.ReportBugList()
		fmt.Println(bugsToDo)
		sendMsg(bugsToDo)

		storiesToDo := comProjects.ReportStoryToDoMetric()
		fmt.Println(storiesToDo)
		sendMsg(storiesToDo)

		storiesToDo = comProjects.ReportStoryList()
		fmt.Println(storiesToDo)
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

			pk := NewProject(458, "山东平台")
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
