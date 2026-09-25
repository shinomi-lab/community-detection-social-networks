// cal_sameimpression_cost.go
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

func user_same(
	// adj [][]int,
	net network.Network,
	interest_list diff.InterestList,
	assum_list diff.AssumList, exit_f bool, use_cost_infl bool,
	r *rand.Rand,
) {
	var pop_list = diff.MakePopList(diff.PopHigh, diff.PopHigh)
	// var pop_list [2]int
	// pop_list[0] = diff.PopHigh
	// pop_list[1] = diff.PopHigh

	// var seq [16]float64 = diff.Make_probability()
	// var prob_map diff.UserProbTable = diff.Map_probagbility(seq)
	prob_map := diff.GetUserProbTable()

	non_use_list := make([]int, 1)
	max_user := 0 //最もフォロワ数が多いユーザ名
	max_user_num := 0
	user_num_counter := 0
	for i := 0; i < net.N; i++ {
		user_num_counter = 0
		for j := 0; j < net.N; j++ {
			if net.Adj[i][j] == 1 {
				user_num_counter++
			}
		}
		if max_user_num < user_num_counter {
			max_user = i
			max_user_num = user_num_counter
		}
	}
	non_use_list[0] = max_user

	if use_cost_infl {
		opt.SameImpressionCostInfl(
			100, net, non_use_list, prob_map, pop_list, interest_list, assum_list, true, exit_f, r)

	} else {
		opt.SameImpressionCost(
			100, net, non_use_list, prob_map, pop_list, interest_list, assum_list, true, exit_f, r)
	}
}

func follower_same(
	// adj [][]int,
	net network.Network,
	interest_list diff.InterestList, assum_list diff.AssumList, exit_f bool, use_cost_infl bool,
	r *rand.Rand,
) {
	var pop_list = diff.MakePopList(diff.PopHigh, diff.PopHigh)
	// var pop_list [2]int
	// pop_list[0] = diff.PopHigh
	// pop_list[1] = diff.PopHigh

	// var seq [16]float64 = diff.Make_probability()
	// var prob_map diff.UserProbTable = diff.Map_probagbility(seq)
	prob_map := diff.GetUserProbTable()

	SeedSet_F := make([]diff.SeedInfo, net.N)
	max_user := 0 //最もフォロワ数が多いユーザ名
	max_user_num := 0
	user_num_counter := 0
	for i := 0; i < net.N; i++ {
		user_num_counter = 0
		for j := 0; j < net.N; j++ {
			if net.Adj[i][j] == 1 {
				user_num_counter++
			}
		}
		if max_user_num < user_num_counter {
			max_user = i
			max_user_num = user_num_counter
		}
	}
	SeedSet_F[max_user] = diff.SeedInfoF // 1

	opt.SameImpressionCostFollower(
		100, net, SeedSet_F, prob_map, pop_list, interest_list, assum_list, 15, 16, exit_f, use_cost_infl, r)
}

func main() {
	use_cost_infl := true
	fmt.Println("start cal_sameimpression_cost.go")
	// rand.Seed(int64(1))
	var seed int64 = 1
	r := rand.New(rand.NewSource(seed))

	adjFilePath := "Graphs/adj_json1000node.txt"
	adjFilePath = "adj_jsonTwitterInteractionUCongress.txt"
	net := network.ReadAdjJson(adjFilePath)
	interestList := diff.MakeInterestList(net.N, r)
	assumList := diff.MakeAssumList(net.N, r)
	// adj, interest_list, assum_list := Make_adj_interest_assum(adjFilePath, seed)
	// fmt.Println(net.N)
	// os.Exit(0)
	user_same(net, interestList, assumList, true, use_cost_infl, r)
	// follower_same(adj,interest_list,assum_list,true)

}
