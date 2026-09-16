package network

import (
	"encoding/json"
	"os"
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

func ReadJson(adjFilePath string) Network {
	Adj, N := ReadAdjMatJson(adjFilePath)
	FollowerNums := GetFollowerNums(Adj)

	return Network{Adj, N, FollowerNums}
}

func ReadAdjMatJson(adjFilePath string) ([][]int, int) {
	bytes, err := os.ReadFile(adjFilePath)
	if err != nil {
		panic(err)
	}

	var dataJson string = string(bytes)
	arr := make(map[int]map[int]int)
	_ = json.Unmarshal([]byte(dataJson), &arr)

	nNodes := len(arr) // number of nodes
	adj := make([][]int, nNodes)

	for i := 0; i < nNodes; i++ {
		adj[i] = make([]int, nNodes)
		for j := 0; j < nNodes; j++ {
			adj[i][j] = arr[j][i]
		}
	}
	return adj, nNodes
}
