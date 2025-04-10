package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"os"
)

func load(fd string) map[string]int {

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

var DataPath string
var workerList []string
var dictMap map[string]string
var workerGitMetre map[string]int

func addWorkerGitMetre(worker string, num int) {

	if name, ok := dictMap[worker]; !ok {
		workerGitMetre[worker] = 0
		fmt.Println(worker, "not in dict")
	} else {
		if _, ok := workerGitMetre[name]; !ok {
			workerGitMetre[name] = num
		} else {
			workerGitMetre[name] += num
		}
	}
}

func init() {
	flag.StringVar(&DataPath, "p", "xxx.gob", "gob")
	flag.Parse()

	workerList = []string{
		"张桐馨", "汪峰", "包容", "朱水龙", "胡叶涛",
		"刘昂", "曲洪良", "王子岳", "刘涛", "大李扬",
		"叶鑫", "小李扬", "罗健键", "董爽", "刑康",
		"曹梦玄", "韩振豪", "赵加兴", "刘凌云", "孙光勇",
		"魏龙飞", "师伟恒", "姚文杰", "王晓旺", "张一凡",
		"袁永剑", "田浩", "李向阳", "陈浩楠", "李裕龙",
		"吴宸如", "余震",
	}
	workerGitMetre = make(map[string]int, len(workerList))

	dictMap = map[string]string{
		"zhangtongxin": "张桐馨",
		"wangfeng":     "汪峰",
		"BR":           "包容",
		"zhushuilong":  "朱水龙",
		"shuilong":     "朱水龙",
		"Yetao":        "胡叶涛",
		"yetao":        "胡叶涛",
		"liuang":       "刘昂",
		"ang":          "刘昂",
		"hongliang qu": "曲洪良",
		"quhongliang":  "曲洪良",
		"WangZiyue":    "王子岳",
		"snowlab":      "王子岳",
		"liutao":       "刘涛",
		"liyang":       "大李扬",
		"yexin":        "叶鑫",
		"liyang-q":     "小李扬",
		"luojianjian":  "罗健键",
		"dongshuang":   "董爽",
		"xingkang":     "刑康",
		"caomengxuan":  "曹梦玄",
		"韩振豪":          "韩振豪",
		"赵加兴":          "赵加兴",
		"zhaojiaxing":  "赵加兴",
		"刘凌云":          "刘凌云",
		"liulingyun":   "刘凌云",
		"sunguangyong": "孙光勇",
		"weilongfei":   "魏龙飞",
		"swh":          "师伟恒",
		"shiweiheng":   "师伟恒",
		"yaowenjie":    "姚文杰",
		"xiaowang":     "王晓旺",
		"karsa":        "王晓旺",
		"zhangyifan":   "张一凡",
		"yvan":         "张一凡",
		"yuanyongjian": "袁永剑",
		"田浩":           "田浩",
		"tianh":        "田浩",
		"lixiangyang":  "李向阳",
		"haonan":       "陈浩楠",
		"李裕龙":          "李裕龙",
		"wangfangming": "王访明",
		"wcr":          "吴宸如",
		"yuzhen":       "余震",
		"wangzhaoren":  "王肈人",
		"wulei":        "吴磊",
		"ChuckiePan":   "潘小功",
		"yanpengpeng":  "宴朋朋",
	}

}
func main() {

	data := load(DataPath)

	for k, v := range data {
		fmt.Println(k, v)
		addWorkerGitMetre(k, v)
	}

	fmt.Println("=======================================")
	fmt.Println("Git提交操作记录:")
	fmt.Println("=======================================")
	for _, n := range workerList {
		fmt.Println(n, ":", workerGitMetre[n])
	}

}
