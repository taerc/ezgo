package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// UAVRoute 结构体映射 XML 的根元素
type UAVRoute struct {
	XMLName      xml.Name `xml:"UAVRoute"`
	Version      string   `xml:"version,attr"`
	CreatedTime  string   `xml:"CreatedTime"`
	Creator      string   `xml:"Creator"`
	Corporation  string   `xml:"Corporation"`
	DroneInfo    DroneInfo
	SafetyDis    float64 `xml:"SafetyDis"`
	ShootingDis  float64 `xml:"ShootingDis"`
	RouteType    int     `xml:"RouteType"`
	BaseLocation string  `xml:"BaseLocation"`
	LineName     string  `xml:"LineName"`
	Voltage      string  `xml:"Voltage"`
	TowerList    []Tower `xml:"TowerList>Tower"`
}

// DroneInfo 结构体映射无人机信息
type DroneInfo struct {
	DroneType string `xml:"DroneType"`
}

// Tower 结构体映射杆塔信息
type Tower struct {
	Name         string       `xml:"Name,attr"`
	TowerCode    string       `xml:"TowerCode"`
	WayPointList []HoverPoint `xml:"WayPointList>HoverPoint"`
}

// HoverPoint 结构体映射悬停点信息
type HoverPoint struct {
	Type              string       `xml:"Type,attr"`
	Number            int          `xml:"Number"`
	HoverLocation     string       `xml:"HoverLocation"`
	UAVHeading        float64      `xml:"UAVHeading"`
	GimbalPitch       float64      `xml:"GimbalPitch"`
	RotationDirection string       `xml:"RotationDirection"`
	Speed             float64      `xml:"Speed,omitempty"`
	ActionList        []Action     `xml:"ActionList>Action"`
	MediaPointList    []MediaPoint `xml:"MediaPointList>MediaPoint"`
}

// Action 结构体映射动作信息
type Action struct {
	Order int    `xml:"order,attr"`
	Type  int    `xml:"Type"`
	Param string `xml:"Param"`
}

// MediaPoint 结构体映射拍照点信息
type MediaPoint struct {
	Code               string `xml:"Code"`
	PhotoPointLocation string `xml:"PhotoPointLocation"`
	DisplayName        string `xml:"DisplayName"`
	MediaName          string `xml:"MediaName"`
	LineName           string `xml:"LineName"`
	LineTowerName      string `xml:"LineTowerName"`
	CircuitName        string `xml:"CircuitName"`
	PhaseName          string `xml:"PhaseName"`
	PointName          string `xml:"PointName"`
}

func main() {
	// 读取 XML 文件
	xmlFile, err := os.ReadFile("D:\\wkspace\\doc\\23-山东配网\\01-国网智能材料\\trl航线文件\\杨集四线周塘支线#1-#9.trl")
	if err != nil {
		fmt.Println("Error reading XML file:", err)
		return
	}

	// 替换 XML 版本声明为 1.0
	xmlData := strings.Replace(string(xmlFile), `version="2.0"`, `version="1.0"`, 1)

	// 解析 XML 数据
	var route UAVRoute
	err = xml.Unmarshal([]byte(xmlData), &route)
	if err != nil {
		fmt.Println("Error unmarshalling XML:", err)
		return
	}

	// 打印解析后的数据
	fmt.Printf("UAVRoute Version: %s\n", route.Version)
	fmt.Printf("Created Time: %s\n", route.CreatedTime)
	fmt.Printf("Line Name: %s\n", route.LineName)
	fmt.Printf("Voltage: %s\n", route.Voltage)
	fmt.Printf("Drone Type: %s\n", route.DroneInfo.DroneType)
	fmt.Printf("Safety Distance: %.3f\n", route.SafetyDis)
	fmt.Printf("Shooting Distance: %.3f\n", route.ShootingDis)

	// for _, tower := range route.TowerList {
	// 	fmt.Printf("\nTower Name: %s\n", tower.Name)
	// 	fmt.Printf("Tower Code: %s\n", tower.TowerCode)
	// 	for _, point := range tower.WayPointList {
	// 		fmt.Printf("  Hover Point Type: %s, Number: %d\n", point.Type, point.Number)
	// 		fmt.Printf("  Hover Location: %s\n", point.HoverLocation)
	// 		fmt.Printf("  UAV Heading: %.2f\n", point.UAVHeading)
	// 		fmt.Printf("  Gimbal Pitch: %.2f\n", point.GimbalPitch)
	// 		fmt.Printf("  Rotation Direction: %s\n", point.RotationDirection)
	// 		if point.Speed > 0 {
	// 			fmt.Printf("  Speed: %.2f\n", point.Speed)
	// 		}
	// 		for _, action := range point.ActionList {
	// 			fmt.Printf("    Action Order: %d, Type: %d, Param: %s\n", action.Order, action.Type, action.Param)
	// 		}
	// 		for _, media := range point.MediaPointList {
	// 			fmt.Printf("    Media Point Code: %s\n", media.Code)
	// 			fmt.Printf("    Photo Point Location: %s\n", media.PhotoPointLocation)
	// 			fmt.Printf("    Display Name: %s\n", media.DisplayName)
	// 			fmt.Printf("    Media Name: %s\n", media.MediaName)
	// 		}
	// 	}
	// }
	for _, tower := range route.TowerList {
		for _, point := range tower.WayPointList {
			p := strings.Split(point.HoverLocation, " ")
			fmt.Printf("{\"%v\",\"%v\",%v, %v},\n", tower.Name+"-"+strconv.Itoa(point.Number), tower.Name, p[0], p[1])
		}
	}
}
