package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Graph は隣接リストによる無向グラフ表現
type Graph struct {
	Adj map[int]map[int]bool
}

func NewGraph() *Graph {
	return &Graph{Adj: make(map[int]map[int]bool)}
}

func (g *Graph) AddEdge(u, v int) {
	if g.Adj[u] == nil {
		g.Adj[u] = make(map[int]bool)
	}
	if g.Adj[v] == nil {
		g.Adj[v] = make(map[int]bool)
	}
	g.Adj[u][v] = true
	g.Adj[v][u] = true
}

// 1. 隣接行列ファイル (LFR_n5000_mu0.8_adjmat.txt) の読み込み
func loadAdjMatrix(filepath string) (*Graph, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	g := NewGraph()
	scanner := bufio.NewScanner(file)
	row := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// スペースまたはカンマ区切りの数値を読み込む
		fields := strings.Fields(line)
		for col, valStr := range fields {
			val, _ := strconv.Atoi(valStr)
			if val == 1 && row < col { // 無向グラフのため重複登録を防止
				g.AddEdge(row, col)
			}
		}
		row++
	}
	return g, scanner.Err()
}

// 2. コミュニティ割り当てファイル (LFR10_communities) の読み込み
func loadCommunities(filepath string) (map[int][]int, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	communities := make(map[int][]int)
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		// フォーマット対応: 「ノードID コミュニティID」形式の場合
		if len(fields) == 2 {
			nodeID, _ := strconv.Atoi(fields[0])
			commID, _ := strconv.Atoi(fields[1])
			communities[commID] = append(communities[commID], nodeID)
		} else {
			// 行番号がコミュニティIDで、その行に属するノードが並んでいる形式の場合
			for _, f := range fields {
				nodeID, _ := strconv.Atoi(f)
				communities[lineNum] = append(communities[lineNum], nodeID)
			}
		}
		lineNum++
	}
	return communities, scanner.Err()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// ファイルパス（環境に合わせて書き換えてください）
	adjMatrixPath := "LFR_n5000_mu0.8_adjmat.txt"
	commPath := "Community verification/LFR10_communities"

	fmt.Println("データを読み込み中...")
	G, err := loadAdjMatrix(adjMatrixPath)
	if err != nil {
		log.Fatalf("隣接行列の読み込みエラー: %v", err)
	}

	communities, err := loadCommunities(commPath)
	if err != nil {
		log.Fatalf("コミュニティデータの読み込みエラー: %v", err)
	}

	opponentID := 1 // 虚偽情報コミュニティ (例: ID 1)
	targetID := 2   // 自陣コミュニティ (例: ID 2)

	opponentNodes := communities[opponentID]
	targetNodes := communities[targetID]

	// ルックアップ用のセット作成
	oppSet := make(map[int]bool)
	for _, n := range opponentNodes {
		oppSet[n] = true
	}
	targetSet := make(map[int]bool)
	for _, n := range targetNodes {
		targetSet[n] = true
	}

	// --- 境界エッジ (u, v) の抽出 ---
	// u: Target側 (Community 2), v: Opponent側 (Community 1)
	type Edge struct{ U, V int }
	var interEdges []Edge
	oppBridgeCounts := make(map[int]int)
	targetBridgeCounts := make(map[int]int)

	for _, u := range targetNodes {
		for v := range G.Adj[u] {
			if oppSet[v] {
				interEdges = append(interEdges, Edge{U: u, V: v})
				targetBridgeCounts[u]++
				oppBridgeCounts[v]++
			}
		}
	}

	fmt.Printf("--- データの読み込み完了: 境界エッジ数 = %d 本 ---\n\n", len(interEdges))

	// ==========================================
	// 1. ブリッジノードを起源とするシナリオ (3種)
	// ==========================================

	// 【シナリオ①】相手国側のノード v だけ
	topV, maxVCount := -1, -1
	for v, count := range oppBridgeCounts {
		if count > maxVCount {
			maxVCount = count
			topV = v
		}
	}
	scenario1 := []int{topV}
	fmt.Printf("【シナリオ①】相手国側の最前線ハブ v: %v (境界接続数: %d)\n", scenario1, maxVCount)

	// 【シナリオ②】自陣側のノード u だけ
	topU, maxUCount := -1, -1
	for u, count := range targetBridgeCounts {
		if count > maxUCount {
			maxUCount = count
			topU = u
		}
	}
	scenario2 := []int{topU}
	fmt.Printf("【シナリオ②】自陣側のノード u: %v (境界接続数: %d)\n", scenario2, maxUCount)

	// 【シナリオ③】両方のノード u と v を同時
	scenario3 := []int{topV, topU}
	fmt.Printf("【シナリオ③】両方のノード u, v 同時: %v\n", scenario3)

	// ==========================================
	// 2. コミュニティの中心を起源とするシナリオ (1種)
	// ==========================================

	// 【シナリオ④】コミュニティ内部次数（Internal Degree）が最大のノード
	topInfluencer, maxIntDeg := -1, -1
	for _, node := range opponentNodes {
		intDeg := 0
		for neighbor := range G.Adj[node] {
			if oppSet[neighbor] {
				intDeg++
			}
		}
		if intDeg > maxIntDeg {
			maxIntDeg = intDeg
			topInfluencer = node
		}
	}
	scenario4 := []int{topInfluencer}
	fmt.Printf("【シナリオ④】コミュニティ最高インフルエンサー: %v (内部次数: %d)\n", scenario4, maxIntDeg)

	// ==========================================
	// 3. コミュニティの端寄りを起源とするシナリオ (2種)
	// ==========================================

	// 内部次数が 1 のノード（末端ノード）を抽出
	var peripheryCandidates []int
	for _, node := range opponentNodes {
		intDeg := 0
		for neighbor := range G.Adj[node] {
			if oppSet[neighbor] {
				intDeg++
			}
		}
		if intDeg == 1 {
			peripheryCandidates = append(peripheryCandidates, node)
		}
	}

	if len(peripheryCandidates) == 0 {
		log.Fatal("内部次数が1の端寄りノードが見つかりませんでした。")
	}

	// 【シナリオ⑤】端寄りの「1つのユーザー」のみ
	randSingle := peripheryCandidates[rand.Intn(len(peripheryCandidates))]
	scenario5 := []int{randSingle}
	fmt.Printf("【シナリオ⑤】端寄りの1ユーザー: %v\n", scenario5)

	// 【シナリオ⑥】端寄りの「ランダムな複数ユーザー」（例: 5人）
	numSeeds := 5
	if len(peripheryCandidates) < numSeeds {
		numSeeds = len(peripheryCandidates)
	}

	// シャッフルして指定人数を抽出
	shuffled := make([]int, len(peripheryCandidates))
	copy(shuffled, peripheryCandidates)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	scenario6 := shuffled[:numSeeds]
	fmt.Printf("【シナリオ⑥】端寄りの複数ユーザー (%d人): %v\n", numSeeds, scenario6)
}
