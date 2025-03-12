package topology

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

// 地理坐标点结构
type GeoPoint struct {
	ID   string
	Name string
	Lon  float64 // 经度
	Lat  float64 // 纬度
}

// 连接边结构
type Edge struct {
	From, To *GeoPoint
	Distance float64 // 单位：米
}

// 生成地理连接森林
func GenerateGeoForest(points []*GeoPoint, maxDistance float64) [][]Edge {
	// 阶段1：计算有效边集
	edges := computeValidEdges(points, maxDistance)

	// 阶段2：按距离排序
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].Distance < edges[j].Distance
	})

	// 阶段3：构建并查集
	uf := NewStringUnionFind()
	for _, p := range points {
		uf.Add(p.ID)
	}

	// 阶段4：Kruskal算法生成森林
	var forest [][]Edge
	edgeMap := make(map[string][]Edge) // 根节点ID到边的映射

	for _, e := range edges {
		if rootFrom, rootTo := uf.Find(e.From.ID), uf.Find(e.To.ID); rootFrom != rootTo {
			uf.Union(e.From.ID, e.To.ID)
			newRoot := uf.Find(e.From.ID)

			// 合并两个连通分量的边
			combined := append(edgeMap[rootFrom], edgeMap[rootTo]...)
			combined = append(combined, e)
			delete(edgeMap, rootFrom)
			delete(edgeMap, rootTo)
			edgeMap[newRoot] = combined
		}
	}

	// 阶段5：收集结果
	for _, v := range edgeMap {
		forest = append(forest, v)
	}

	// 处理孤立点（单独作为树）
	for _, p := range points {
		if uf.Find(p.ID) == p.ID && len(edgeMap[p.ID]) == 0 {
			forest = append(forest, []Edge{})
		}
	}

	return forest
}

// 计算有效边（Haversine公式）
func computeValidEdges(points []*GeoPoint, maxDistance float64) []Edge {
	const (
		earthRadius = 6371008.7714 // 地球半径（米）
	)

	var edges []Edge
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			p1, p2 := points[i], points[j]

			// 坐标转弧度
			lat1 := p1.Lat * math.Pi / 180
			lon1 := p1.Lon * math.Pi / 180
			lat2 := p2.Lat * math.Pi / 180
			lon2 := p2.Lon * math.Pi / 180

			// 计算差值
			dlat := lat2 - lat1
			dlon := lon2 - lon1

			// Haversine公式
			a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dlon/2), 2)
			c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
			distance := earthRadius * c

			if distance <= maxDistance {
				edges = append(edges, Edge{
					From:     p1,
					To:       p2,
					Distance: distance,
				})
			}
		}
	}
	return edges
}

// 字符串专用并查集
type StringUnionFind struct {
	parent map[string]string
	rank   map[string]int
}

func NewStringUnionFind() *StringUnionFind {
	return &StringUnionFind{
		parent: make(map[string]string),
		rank:   make(map[string]int),
	}
}

func (suf *StringUnionFind) Add(x string) {
	if _, exists := suf.parent[x]; !exists {
		suf.parent[x] = x
		suf.rank[x] = 0
	}
}

func (suf *StringUnionFind) Find(x string) string {
	if suf.parent[x] != x {
		suf.parent[x] = suf.Find(suf.parent[x])
	}
	return suf.parent[x]
}

func (suf *StringUnionFind) Union(x, y string) {
	xRoot := suf.Find(x)
	yRoot := suf.Find(y)

	if xRoot == yRoot {
		return
	}

	if suf.rank[xRoot] < suf.rank[yRoot] {
		suf.parent[xRoot] = yRoot
	} else {
		suf.parent[yRoot] = xRoot
		if suf.rank[xRoot] == suf.rank[yRoot] {
			suf.rank[xRoot]++
		}
	}
}

// 生成Graphviz格式的可视化文件
func ExportAsDOT(forest [][]Edge, dotPath string) string {
	builder := strings.Builder{}
	builder.WriteString("graph G {\n")
	for _, tree := range forest {
		for _, e := range tree {
			builder.WriteString(fmt.Sprintf("  \"%s\" -- \"%s\" [label=%.1f];\n",
				e.From.ID+"-"+e.From.Name, e.To.ID+"-"+e.To.Name, e.Distance))
		}
	}
	builder.WriteString("}")
	os.WriteFile(dotPath, []byte(builder.String()), 0644)
	return builder.String()
}
