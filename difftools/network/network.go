package network

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FollowerSize(adj [][]int, node int) int {
	ans := 0
	for _, isEdge := range adj[node] {
		ans += isEdge
	}

	return ans
}

func GetFollowerNums(adj [][]int) []int {
	n := len(adj)
	counter := make([]int, n)
	for i, arr := range adj {
		for _, a := range arr {
			counter[i] += a
		}
	}
	return counter
}

type Network struct {
	// 隣接行列
	Adj [][]int
	// ノード数
	N int
	// フォロワー数の配列(index: userId)
	Followers []int
}

func ReadAdjJson(adjFilePath string) Network {
	bytes, err := os.ReadFile(adjFilePath)
	if err != nil {
		panic(err)
	}

	var dataJson string = string(bytes)
	arr := make(map[int]map[int]int)
	_ = json.Unmarshal([]byte(dataJson), &arr)

	n := len(arr) // number of nodes
	adj := make([][]int, n)

	for i := range n {
		adj[i] = make([]int, n)
		for j := range n {
			adj[i][j] = arr[j][i]
		}
	}
	FollowerNums := GetFollowerNums(adj)

	return Network{adj, n, FollowerNums}
}

type edge struct {
	from uint
	to   uint
}

func scanEdgelist(scanner *bufio.Scanner) Network {
	edges := make([]edge, 0)
	var n uint = 0

	for scanner.Scan() {
		var edge edge
		t := scanner.Text()
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		_, err := fmt.Sscanf(t, "%d %d", &edge.from, &edge.to)
		if err != nil {
			panic(err)
		}
		n = max(n, edge.from, edge.to)
		edges = append(edges, edge)
	}

	// max id = number of nodes - 1
	n++

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	adj := make([][]int, n)
	for i := range n {
		adj[i] = make([]int, n)
	}

	// ReadAdjJsonに合わせて転置
	for _, e := range edges {
		adj[e.to][e.from] = 1
	}

	return Network{adj, int(n), GetFollowerNums(adj)}
}

func ReadEdgelist(filePath string) Network {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".txt", "":
		f, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		return scanEdgelist(scanner)
	case ".gz":
		f, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		gzReader, err := gzip.NewReader(f)
		if err != nil {
			panic(err)
		}
		defer gzReader.Close()

		scanner := bufio.NewScanner(gzReader)
		return scanEdgelist(scanner)
	default:
		panic(fmt.Sprintf("Invalid type: %s", filePath))
	}
}
