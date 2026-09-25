// outputTest.go
package main

import (
	// "encoding/csv"

	"fmt"

	// "log"
	diff "difftools/diffusion"
	"difftools/network"
	opt "difftools/optimization"

	// "os"
	// "strconv"
	// "strings"
	// "time"
	"math/rand"
	// "reflect"
)

func main() {
	var seed int64 = 1
	r := rand.New(rand.NewSource(seed))

	adjFilePath := "Graphs/adj_json50node.txt"
	// adj, interest_list, assum_list := Make_adj_interest_assum(adjFilePath, seed)
	net := network.ReadAdjJson(adjFilePath)
	interestList := diff.MakeInterestList(net.N, r)
	assumList := diff.MakeAssumList(net.N, r)

	var pop_list = diff.MakePopList(diff.PopHigh, diff.PopHigh)
	// var pop_list [2]int
	// pop_list[0] = diff.PopHigh
	// pop_list[1] = diff.PopHigh
	prob_map := diff.GetUserProbTable()

	InflTest(net, prob_map, pop_list, interestList, assumList)
	Selected_Suppression_MaximumTest(net, prob_map, pop_list, interestList, assumList, r)
}

func InflTest(
	// adj [][]int,
	net network.Network,
	prob_map diff.UserProbTable, pop_list diff.PopList, interest_list diff.InterestList, assum_list diff.AssumList) {

	SeedSetF := make([]diff.SeedInfo, net.N)
	SeedSetF[2] = diff.SeedInfoF // 1
	SeedSetF[0] = diff.SeedInfoT // 2
	// rand.Seed(100)
	r := rand.New(rand.NewSource(100))
	hist := opt.RunInflProp(1000, net, SeedSetF, prob_map, pop_list, interest_list, assum_list, r)

	fmt.Println(hist)
}

func Selected_Suppression_MaximumTest(
	// adj [][]int,
	net network.Network,
	prob_map diff.UserProbTable, pop_list diff.PopList, interest_list diff.InterestList, assum_list diff.AssumList,
	r *rand.Rand,
) {
	SeedSetF := make([]diff.SeedInfo, net.N)
	SeedSetF[2] = diff.SeedInfoF // 1
	greedy_ans := []int{0}
	greedy_ans2 := make([][]int, 0)
	greedy_ans2 = append(greedy_ans2, greedy_ans)

	_, test_greedy_ans_v, _ := opt.SelectedSuppressionMaximum(
		net, greedy_ans2, SeedSetF, prob_map, pop_list, interest_list, assum_list, r)

	fmt.Println(test_greedy_ans_v)
}

// func Make_adj_interest_assum(adjFilePath string, seed int64) ([][]int, [][]int, [][]int) {
// 	bytes, err := ioutil.ReadFile(adjFilePath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// fmt.Println(string(bytes))

// 	var dataJson string = string(bytes)

// 	arr := make(map[int]map[int]int)
// 	// var arr []string
// 	_ = json.Unmarshal([]byte(dataJson), &arr)
// 	// fmt.Println(arr)

// 	// fmt.Println(arr[0][1])

// 	n := len(arr)

// 	var interest_list diff.InterestList = diff.Make_interest_list(n, seed)

// 	var assum_list diff.AssumList = diff.Make_assum_list(n, seed)
// 	var adj [][]int = make([][]int, n)

// 	for i := 0; i < n; i++ {
// 		adj[i] = make([]int, n)
// 		for j := 0; j < n; j++ {
// 			adj[i][j] = arr[j][i]
// 		}
// 	}
// 	return adj, interest_list, assum_list
// }
