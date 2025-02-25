package main

import (
	"fmt"
	"math"
	"sort"
)

// 空间坐标结构（2025标准）
type GeoPoint struct {
	ID  int     `json:"id"`
	Tag string  `json:"tag"`
	Lon float64 `json:"lon"` // 经度（WGS84）
	Lat float64 `json:"lat"` // 纬度（WGS84）
}

// 空间连接关系
type SpatialLink struct {
	From     int     `json:"from"`
	To       int     `json:"to"`
	Distance float64 `json:"distance"` // 单位：米
}

// 地理森林结构
type GeoForest struct {
	Clusters map[int][]int         `json:"clusters"` // 簇ID -> 成员节点
	Edges    map[int][]SpatialLink `json:"edges"`    // 簇ID -> 连接关系
}

// 高精度测距（2025新版Haversine优化算法）
func spatialDistance(lon1, lat1, lon2, lat2 float64) float64 {
	const (
		R = 6371e3 // 地球半径（米）
		π = math.Pi
	)

	th1 := lat1 * π / 180
	th2 := lat2 * π / 180
	delt1 := (lat2 - lat1) * π / 180
	delt2 := (lon2 - lon1) * π / 180

	a := math.Sin(delt1/2)*math.Sin(delt1/2) + math.Cos(th1)*math.Cos(th2)*math.Sin(delt2/2)*math.Sin(delt2/2)
	return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// 并查集（支持动态路径压缩）
type UnionSet struct {
	parent map[int]int
	rank   map[int]int
}

func NewUnionSet(points []GeoPoint) *UnionSet {
	us := &UnionSet{
		parent: make(map[int]int),
		rank:   make(map[int]int),
	}
	for _, p := range points {
		us.parent[p.ID] = p.ID
		us.rank[p.ID] = 0
	}
	return us
}

func (us *UnionSet) Find(id int) int {
	if us.parent[id] != id {
		us.parent[id] = us.Find(us.parent[id])
	}
	return us.parent[id]
}

func (us *UnionSet) Union(x, y int) {
	xRoot := us.Find(x)
	yRoot := us.Find(y)
	if xRoot == yRoot {
		return
	}

	if us.rank[xRoot] < us.rank[yRoot] {
		us.parent[xRoot] = yRoot
	} else {
		us.parent[yRoot] = xRoot
		if us.rank[xRoot] == us.rank[yRoot] {
			us.rank[xRoot]++
		}
	}
}

func BuildGeoForest(points []GeoPoint) GeoForest {
	// 生成候选边
	var candidateEdges []SpatialLink
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			d := spatialDistance(points[i].Lon, points[i].Lat,
				points[j].Lon, points[j].Lat)
			if d <= 100 {
				candidateEdges = append(candidateEdges,
					SpatialLink{points[i].ID, points[j].ID, d})
			}
		}
	}

	// Kruskal算法核心
	sort.Slice(candidateEdges, func(i, j int) bool {
		return candidateEdges[i].Distance < candidateEdges[j].Distance
	})

	us := NewUnionSet(points)
	forest := GeoForest{
		Clusters: make(map[int][]int),
		Edges:    make(map[int][]SpatialLink),
	}

	// 构建连接关系
	for _, edge := range candidateEdges {
		if us.Find(edge.From) != us.Find(edge.To) {
			root := us.Find(edge.From)
			forest.Edges[root] = append(forest.Edges[root], edge)
			us.Union(edge.From, edge.To)
		}
	}

	// 生成簇结构
	for _, p := range points {
		root := us.Find(p.ID)
		forest.Clusters[root] = append(forest.Clusters[root], p.ID)
	}

	return forest
}

func main() {
	// 2025测试数据集（时空坐标）
	points := []GeoPoint{
		{101, "气象站A", 116.404177, 39.909652},
		{102, "5G基站B", 116.404181, 39.909655},
		{103, "无人机C", 116.405002, 39.910112},
		{104, "智能路灯D", 116.404190, 39.909658},
	}

	forest := BuildGeoForest(points)

	// 三维可视化输出
	fmt.Println("🌐 空间拓扑结构（2025时空网格标准）")
	for cluster, members := range forest.Clusters {
		fmt.Printf("\n▣ 簇 %d [成员%d个]:\n", cluster, len(members))
		fmt.Printf("   📌 节点ID: %v\n", members)

		if edges, ok := forest.Edges[cluster]; ok {
			fmt.Println("   ⛓️ 连接关系:")
			total := 0.0
			for _, e := range edges {
				fmt.Printf("     ▸ %d ↔ %d (%.2fm)\n", e.From, e.To, e.Distance)
				total += e.Distance
			}
			fmt.Printf("   📏 总链路: %.2f米\n", total)
		}
	}
}
