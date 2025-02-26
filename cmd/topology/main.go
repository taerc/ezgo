package main

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
func GenerateGeoForest(points []*GeoPoint) [][]Edge {
	// 阶段1：计算有效边集
	edges := computeValidEdges(points)

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
func computeValidEdges(points []*GeoPoint) []Edge {
	const (
		maxDistance = 100 // 米
		earthRadius = 6371000
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

// 示例运行
func main() {

	// 数据使用10.28 上
	// 10kV81516干线
	points := []*GeoPoint{
		{"2830954", "#003", 124.98170461751867, 45.80236086240778},
		{"2830955", "#002", 124.9811009852985, 45.80251331291197},
		{"2830956", "#001", 124.98049131536185, 45.802660019385385},
		{"2830957", "#037+1", 124.96262655509926, 45.80481884550502},
		{"2830958", "#037", 124.96264903487557, 45.80481208718511},
		{"2830959", "#036+3", 124.9629293083139, 45.804749345260134},
		{"2830960", "#036+2", 124.9629568141383, 45.80474105464682},
		{"2830961", "#036+1", 124.96317932046965, 45.80468284243828},
		{"2830962", "#036", 124.96320675944217, 45.80467609227396},
		{"2830963", "#035", 124.96392538560282, 45.80448746142196},
		{"2830964", "#034", 124.96456271872479, 45.80432415840722},
		{"2830965", "#033", 124.96515400043134, 45.804166525659355},
		{"2830967", "#032", 124.96582028987044, 45.803991055688066},
		{"2830968", "#031", 124.96651107417632, 45.80381481913317},
		{"2830969", "#030", 124.96710331795411, 45.80365906701797},
		{"2830970", "#029", 124.96785385652899, 45.803455961190544},
		{"2830971", "#028", 124.96854156815007, 45.80329040452296},
		{"2830973", "#027", 124.96919052526644, 45.80311717946639},
		{"2830974", "#026", 124.96984966533495, 45.80294558265144},
		{"2830975", "#025+1", 124.97051608041544, 45.80277175082501},
		{"2830976", "#025", 124.97054115226938, 45.802764063520485},
		{"2830977", "#024", 124.97118899801376, 45.80259839790871},
		{"2830978", "#023", 124.9718375181787, 45.80243252699774},
		{"2830979", "#022", 124.97248155002086, 45.802262485696986},
		{"2830980", "#021", 124.97312683980456, 45.80209539178877},
		{"2830981", "#020", 124.97379128144297, 45.801923599381446},
		{"2830982", "#019", 124.97446075685214, 45.80174705162421},
		{"2830983", "#018", 124.97510157665116, 45.80159087080445},
		{"2830984", "#017", 124.97585254883082, 45.80139622467642},
		{"2830985", "#016", 124.9764510385523, 45.80124566467971},
		{"2830986", "#015", 124.97710996058507, 45.80106583576025},
		{"2830987", "#014+1", 124.9777156485494, 45.80091231893028},
		{"2830988", "#014", 124.97774046056074, 45.80090555786951},
		{"2830989", "#013", 124.97835302011674, 45.80074375255143},
		{"2830990", "#012", 124.97886379352586, 45.80061430831068},
		{"2830991", "#011", 124.97945873992775, 45.80045566263645},
		{"2830992", "#010", 124.98004849155056, 45.800302972282545},
		{"2830993", "#009", 124.98060520227179, 45.80015474570382},
		{"2830994", "#008", 124.9811434588459, 45.800016060422436},
		{"2830995", "#007", 124.98183956292321, 45.79984019136491},
		{"2830996", "#006", 124.98247617560507, 45.7996594002319},
		{"2830997", "#005", 124.98301720863738, 45.79952675393194},
		{"2830998", "#004", 124.98310824309493, 45.79921028487757},
		{"2830999", "#003", 124.98318173570274, 45.79897258297864},
		{"2831000", "#002", 124.9834179949487, 45.799014694295884},
		{"2831001", "#001", 124.9835268884542, 45.79903155649439},
		{"2834672", "#004", 124.98232683124488, 45.802204808250025},
		{"2834713", "#005", 124.98294811181304, 45.80204668089125},
		{"2834793", "#006", 124.98291288659641, 45.80245689057165},
		{"2834865", "#006+1", 124.98291239111734, 45.80247661084325},
		{"2834927", "#007", 124.98357048381142, 45.80189837599144},
		{"2834954", "#008", 124.98424106238043, 45.801730201449274},
		{"2834978", "#009", 124.9849024008554, 45.80156990994417},
		{"2834997", "#010", 124.9855426477759, 45.80140653809005},
		{"2835138", "#011", 124.98615378253675, 45.80125947344667},
		{"2835191", "#012", 124.98683656232319, 45.801086271563925},
		{"2835228", "#013", 124.98734814784358, 45.800958263758474},
		{"2835242", "#014", 124.98794168185684, 45.8008079087909},
		{"2835280", "#015", 124.98854366388541, 45.8006503936192},
		{"2835322", "#016", 124.98908285043797, 45.800518577491935},
		{"2835345", "#017", 124.98947466639258, 45.80042659426682},
		{"2835407", "#018", 124.98969441212063, 45.80070559043525},
		{"2835432", "#019", 124.98970772005781, 45.80072460531562},
		{"2835673", "#038", 124.9789776220626, 45.80088254679304},
		{"2835685", "#039", 124.97912567292099, 45.801220833525},
		{"2835733", "#040", 124.97935535245483, 45.80176660258272},
		{"2835742", "#041", 124.97952844864294, 45.80216289907515},
		{"2835776", "#042", 124.97968491398325, 45.802556112288336},
		{"2835799", "#043", 124.97980877314899, 45.80283797306536},
		{"2835809", "#044", 124.97995458597757, 45.80318929535825},
		{"2835821", "#045", 124.98008592032453, 45.80349213566148},
		{"2835827", "#046", 124.98031785227796, 45.80404330821498},
		{"2835884", "#001", 124.97956164916114, 45.80289078140117},
		{"2835921", "#001+1", 124.97951617209338, 45.80290185787053},
		{"2835932", "#001+2", 124.97909308548195, 45.80301578186503},
		{"2835950", "#003", 124.97842158052718, 45.803180965447865},
		{"2835968", "#004", 124.97869883192907, 45.80352306029723},
		{"2835986", "#005", 124.97889739336418, 45.80376497484512},
		{"2835997", "#006", 124.97890918255602, 45.803782737604},
		{"2836006", "#007", 124.97782711662266, 45.803323434520685},
		{"2836013", "#008", 124.97722540679227, 45.80347181009804},
		{"2836017", "#009", 124.97593901958963, 45.803794484459154},
		{"2836157", "#010", 124.97582954979198, 45.80424986245548},
		{"2836166", "#010+1", 124.97582536779797, 45.804268367680294},
		{"2836190", "#011", 124.97535517893586, 45.80393205617139},
		{"2836193", "#012", 124.97532840505103, 45.803940074085766},
		{"2836204", "#012+1", 124.9746675963847, 45.804113293017934},
		{"2836211", "#013", 124.97404925335888, 45.80426902758787},
		{"2836219", "#014", 124.97336129155713, 45.80443745792135},
		{"2836229", "#015", 124.97274235212174, 45.804597272696405},
		{"2836236", "#016", 124.97213733346281, 45.80474821201976},
		{"2836274", "#017", 124.97155589876526, 45.80489240147396},
		{"2836299", "#018", 124.9717313701247, 45.805397181992745},
		{"2836305", "#018+1", 124.97173855043098, 45.80541485031222},
		{"2930559", "#047", 124.98049663280429, 45.80446317154368},
		{"2930578", "#048", 124.98068357233646, 45.80490340325937},
		{"2930593", "#049", 124.98086445160655, 45.80534073591543},
		{"2930610", "#050", 124.98105337775254, 45.8057814806972},
		{"2930621", "#051", 124.9812348434528, 45.806223079966244},
		{"2930636", "#052", 124.98139862629236, 45.806608553906486},
		{"2930650", "#053", 124.98157495450856, 45.80702736163202},
		{"2930664", "#054", 124.98175357565374, 45.80744485541607},
		{"2930675", "#055", 124.9819228321018, 45.807848514835044},
		{"2930693", "#056", 124.98214128341148, 45.808372318162604},
		{"2930714", "#057", 124.98135189394593, 45.808579484650174},
		{"2930746", "#058", 124.98090037097576, 45.80869358135301},
		{"2930762", "#059", 124.98017577972114, 45.80887769650151},
		{"2930923", "#061", 124.97895802210104, 45.8092034287149},
		{"2930940", "#062", 124.97862236954327, 45.80947396368792},
		{"2930962", "#062+1", 124.97799710536806, 45.80943926291022},
		{"2930973", "#062+2", 124.97786044165467, 45.80945871991101},
		{"2930994", "#063", 124.97746029142922, 45.80956084612166},
		{"2931028", "#064", 124.97686398287023, 45.8097218343842},
		{"2931048", "#065", 124.97626910982147, 45.80987687896167},
		{"2931059", "#066", 124.97661886813464, 45.81030361553426},
		{"2931071", "#067", 124.97693728316511, 45.81067275206956},
		{"2931085", "#068", 124.97726843530037, 45.811067881280714},
		{"2931095", "#069", 124.97761699368449, 45.8114935679161},
		{"2931196", "#070", 124.97791751082846, 45.81185238110953},
		{"2931232", "#084", 124.97817576509699, 45.81216459766887},
		{"2931676", "#094", 124.97759914804594, 45.81538219631272},
		{"2931682", "#095", 124.97716622287463, 45.81548314949305},
		{"2931685", "#096", 124.97652413840868, 45.81562799392427},
		{"2931688", "#097", 124.97586411815332, 45.815776816232706},
		{"2931690", "#098", 124.97522197586679, 45.815928513835765},
		{"2931694", "#099", 124.97465527179814, 45.81605923125822},
		{"2931700", "#99+1", 124.97396168167569, 45.816216643950234},
		{"2931710", "#100", 124.97338022976761, 45.81635237970595},
		{"2931716", "#100+1", 124.97335065483132, 45.81635870283243},
		{"2931720", "#102", 124.9726738535758, 45.81651828570229},
		{"2931724", "#102+1", 124.97197592228748, 45.81666827379906},
		{"2931763", "#103", 124.97309422598158, 45.81686426402857},
		{"2931768", "#104", 124.97354126455953, 45.8172166754194},
		{"2931773", "#105", 124.97396070013147, 45.81757663450883},
		{"2931778", "#106", 124.97442553233078, 45.817954546539475},
		{"2931782", "#107", 124.97488693234686, 45.81832621093554},
		{"2931790", "#108", 124.97535752969362, 45.81870964763927},
		{"2931797", "#109", 124.97585887764033, 45.818639335794465},
		{"2931817", "#111", 124.97643092387207, 45.81856991132475},
		{"2931837", "#111+1", 124.97656387465561, 45.81855361457002},
	}

	forest := GenerateGeoForest(points)

	// for i, tree := range forest {
	// 	if len(tree) == 0 {
	// 		fmt.Printf("孤立树 %d: 单独节点\n", i+1)
	// 		continue
	// 	}
	// 	fmt.Printf("生成树 %d (包含%d个连接):\n", i+1, len(tree))
	// 	for _, e := range tree {
	// 		// fmt.Printf("  %s ←[%.1fm]→ %s\n", e.From.ID, e.Distance, e.To.ID)
	// 		// fmt.Printf("  %s ←[%.1fm]→ %s\n", e.From.Name,  e.Distance, e.To.Name)
	// 		fmt.Printf("  %s → %s [%1.fm]\n", e.From.Name, e.To.Name, e.Distance)
	// 	}
	// 	fmt.Println()
	// }

	dot := ExportAsDOT(forest)
	fmt.Println(dot)
}

// 生成Graphviz格式的可视化文件
func ExportAsDOT(forest [][]Edge) string {
	builder := strings.Builder{}
	builder.WriteString("graph G {\n")
	for _, tree := range forest {
		for _, e := range tree {
			// builder.WriteString(fmt.Sprintf("  \"%s\" -- \"%s\" [label=%.1f];\n",
			// 	e.From.Name, e.To.Name, e.Distance)) // 杆塔号有重名的情况，需要重新分析
			builder.WriteString(fmt.Sprintf("  \"%s\" -- \"%s\" [label=%.1f];\n",
				e.From.ID + "-" + e.From.Name , e.To.ID + "-" + e.To.Name, e.Distance))
		}
	}
	builder.WriteString("}")
	os.WriteFile("graph.dot", []byte(builder.String()), 0644)
	return builder.String()
}
