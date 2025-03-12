package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	ezgo "github.com/taerc/ezgo/pkg"
)

type TowerPointData struct {
	TowerName     string `json:"towerName"`
	LineName      string `json:"lineName"`
	CreatedTime   string `json:"createdTime"`
	WayPointList  []PointSimpleData `json:"wayPointList"`
	Filled        bool `json:"filled"`
}

type PointSimpleData struct {
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	Alt           float64 `json:"alt"`
	Head          float64 `json:"head"`
	Shoot         bool `json:"shoot"`
	DevicePartName string `json:"devicePartName"`
	MediaName     string `json:"mediaName"`
	Index         int `json:"index"`
}

func replaceXMLVersion(xmlContent string) string {
    return strings.Replace(xmlContent, `version="2.0"`, `version="1.0"`, 1)
}

func parseTrl(filePath string) ([]TowerPointData, error) {
	var list []TowerPointData

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return list, fmt.Errorf("invalid file name: %v", err)
	}

	if fileInfo.IsDir() {
		files, err := os.ReadDir(filePath)
		if err != nil {
			return list, fmt.Errorf("failed to read directory: %v", err)
		}

		for _, file := range files {
			childFilePath := filepath.Join(filePath, file.Name())
			childList, err := parseTrl(childFilePath)
			if err != nil {
				fmt.Printf("Error parsing file %s: %v\n", file.Name(), err)
				continue
			}
			list = append(list, childList...)
		}
		return list, nil
	}

	if !strings.Contains(fileInfo.Name(), "trl") {
		return list, nil
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return list, fmt.Errorf("failed to open file: %v", err)
	}
	// defer file.Close()

	// 替换 XML 版本
    xmlContent := replaceXMLVersion(string(fileContent))

    // 使用 strings.NewReader 将字符串转换为 io.Reader
    decoder := xml.NewDecoder(strings.NewReader(xmlContent))

	// decoder := xml.NewDecoder(file)
	var towerPointData *TowerPointData
	var pointSimpleData *PointSimpleData
	var lineName, createdTime string

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return list, fmt.Errorf("error decoding XML: %v", err)
		}

				var actionType string
		switch se := token.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "LineName":
				if lineName == "" {
					err := decoder.DecodeElement(&lineName, &se)
					if err != nil {
						return list, fmt.Errorf("error decoding LineName: %v", err)
					}
					if strings.Contains(lineName, "线") {
						fmt.Printf("加载轨迹线路名称原始=%s\n", lineName)
						if strings.Contains(lineName, "线") && !strings.HasSuffix(lineName, "线") {
							lineName = lineName[:strings.LastIndex(lineName, "线")+1]
						}
						fmt.Printf("加载轨迹线路名称修改后=%s\n", lineName)
						if towerPointData != nil && towerPointData.LineName == "" {
							towerPointData.LineName = lineName
						}
					}
				}
			case "CreatedTime":
				err := decoder.DecodeElement(&createdTime, &se)
				if err != nil {
					return list, fmt.Errorf("error decoding CreatedTime: %v", err)
				}
			case "Tower":
				var towerName string
				for _, attr := range se.Attr {
					if attr.Name.Local == "name" {
						towerName = attr.Value
						break
					}
				}
				if !strings.HasPrefix(towerName, "#") {
					towerName = "#" + strings.ReplaceAll(towerName, "#", "")
				}
				towerPointData = &TowerPointData{}
				if strings.Contains(lineName, "线") {
					towerPointData.LineName = lineName
					towerPointData.TowerName = lineName + "_" + towerName
				}
				towerPointData.TowerName = towerName
				towerPointData.CreatedTime = createdTime
			case "WayPointList":
				towerPointData.WayPointList = []PointSimpleData{}
			case "HoverPoint":
				pointSimpleData = &PointSimpleData{}
			case "Number":
				var indexString string
				err := decoder.DecodeElement(&indexString, &se)
				if err != nil {
					return list, fmt.Errorf("error decoding Number: %v", err)
				}
				index, err := strconv.Atoi(indexString)
				if err != nil {
					fmt.Printf("pull2xml: index %s\n", indexString)
				} else {
					pointSimpleData.Index = index
				}
			case "HoverLocation":
				var location string
				err := decoder.DecodeElement(&location, &se)
				if err != nil {
					return list, fmt.Errorf("error decoding HoverLocation: %v", err)
				}
				if location != "" {
					coords := strings.Split(location, " ")
					if len(coords) == 3 {
						lon, err := strconv.ParseFloat(coords[0], 64)
						if err != nil {
							fmt.Printf("pull2xml: location %s\n", location)
							continue
						}
						lat, err := strconv.ParseFloat(coords[1], 64)
						if err != nil {
							fmt.Printf("pull2xml: location %s\n", location)
							continue
						}
						alt, err := strconv.ParseFloat(coords[2], 64)
						if err != nil {
							fmt.Printf("pull2xml: location %s\n", location)
							continue
						}
						pointSimpleData.Lat = lat
						pointSimpleData.Lon = lon
						pointSimpleData.Alt = alt
					}
				}
			case "UAVHeading":
				var headString string
				err := decoder.DecodeElement(&headString, &se)
				if err != nil {
					return list, fmt.Errorf("error decoding UAVHeading: %v", err)
				}
				head, err := strconv.ParseFloat(headString, 64)
				if err != nil {
					fmt.Printf("pull2xml: head %s\n", headString)
				} else {
					pointSimpleData.Head = head
				}
			case "Type":
				err := decoder.DecodeElement(&actionType, &se)
				if err != nil {
					return list, fmt.Errorf("error decoding Type: %v", err)
				}
				if actionType == "3" {
					pointSimpleData.Shoot = true
				}
			case "Param":
				var param string
				err := decoder.DecodeElement(&param, &se)
				if err != nil {
					return list, fmt.Errorf("error decoding Param: %v", err)
				}
				if actionType == "1" {
					gimbalPitch, err := strconv.ParseFloat(param, 64)
					if err != nil {
						fmt.Printf("pull2xml: gimbalPitch %s\n", param)
					} else {
						pointSimpleData.Head = gimbalPitch
					}
				} else if actionType == "5" {
					focalLength, err := strconv.Atoi(param)
					if err != nil {
						fmt.Printf("pull2xml: focalLength %s\n", param)
					} else {
						pointSimpleData.Index = focalLength
					}
				}
			case "MediaName":
				err := decoder.DecodeElement(&pointSimpleData.MediaName, &se)
				if err != nil {
					return list, fmt.Errorf("error decoding MediaName: %v", err)
				}
			case "DisplayName":
				err := decoder.DecodeElement(&pointSimpleData.DevicePartName, &se)
				if err != nil {
					return list, fmt.Errorf("error decoding DisplayName: %v", err)
				}
			}
		case xml.EndElement:
			switch se.Name.Local {
			case "Tower":
				towerPointData.Filled = true
				if towerPointData.LineName == "" {
					towerPointData.LineName = "未知kV未知线"
					towerPointData.TowerName = towerPointData.LineName + "_" + towerPointData.TowerName
				}
				fmt.Printf("解析到了一个塔,全称=%s, 线路名称=%s, 塔号=%s, 拍照点数=%d, 总航点数=%d\n",
					towerPointData.TowerName, towerPointData.LineName, towerPointData.TowerName, len(towerPointData.WayPointList), len(towerPointData.WayPointList))
				list = append(list, *towerPointData)
			case "HoverPoint":
				if !pointSimpleData.Shoot {
					pointSimpleData.DevicePartName = "安全点"
				} else {
					mediaName := towerPointData.LineName + "_" + towerPointData.TowerName + "_" + pointSimpleData.DevicePartName + ".JPG"
					pointSimpleData.MediaName = mediaName
				}
				towerPointData.WayPointList = append(towerPointData.WayPointList, *pointSimpleData)
			}
		}
	}

	return list, nil
}

func main() {
	filePath := "D:\\wkspace\\doc\\23-山东配网\\01-国网智能材料\\trl航线文件\\正东线#2\\正东线#2.trl"
	list, err := parseTrl(filePath)
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		return
	}
	fmt.Printf("Parsed data: %+v\n", list)

	ezgo.SaveJson("1.json", list)
}