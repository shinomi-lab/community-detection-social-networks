package main

import (
	// "encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"

	// "log"
	diff "m/difftools/diffusion"
	opt "m/difftools/optimization"

	// "os"
	// "strconv"
	// "strings"
	// "time"
	"math/rand"
	// "reflect"
)

func Make_adj_interest_assum(adjFilePath string, seed int64) ([][]int, [][]int, [][]int) {
	bytes, err := ioutil.ReadFile(adjFilePath)
	if err != nil {
		panic(err)
	}

	// fmt.Println(string(bytes))

	var dataJson string = string(bytes)

	arr := make(map[int]map[int]int)
	// var arr []string
	_ = json.Unmarshal([]byte(dataJson), &arr)
	// fmt.Println(arr)

	// fmt.Println(arr[0][1])

	n := len(arr)

	var interest_list [][]int = diff.Make_interest_list(n, seed)

	var assum_list [][]int = diff.Make_assum_list(n, seed)
	var adj [][]int = make([][]int, n)

	for i := 0; i < n; i++ {
		adj[i] = make([]int, n)
		for j := 0; j < n; j++ {
			adj[i][j] = arr[j][i]
		}
	}
	return adj, interest_list, assum_list
}

func user_same(adj [][]int, interest_list [][]int, assum_list [][]int, exit_f bool, use_cost_infl bool) {
	var pop_list [2]int

	pop_list[0] = diff.Pop_high
	pop_list[1] = diff.Pop_high

	var seq [16]float64 = diff.Make_probability()

	var prob_map [2][2][2][2]float64 = diff.Map_probagbility(seq)
	non_use_list := make([]int, 1)
	max_user := 0 //最もフォロワ数が多いユーザ名
	max_user_num := 0
	user_num_counter := 0
	for i := 0; i < len(adj); i++ {
		user_num_counter = 0
		for j := 0; j < len(adj); j++ {
			if adj[i][j] == 1 {
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
		opt.SameImpressionCostInfl(0, 100, adj, non_use_list, prob_map, pop_list, interest_list, assum_list, true, exit_f)

	} else {
		opt.SameImpressionCost(0, 100, adj, non_use_list, prob_map, pop_list, interest_list, assum_list, true, exit_f)
	}
}

func follower_same(adj [][]int, interest_list [][]int, assum_list [][]int, exit_f bool, use_cost_infl bool) {
	var pop_list [2]int

	pop_list[0] = diff.Pop_high
	pop_list[1] = diff.Pop_high

	var seq [16]float64 = diff.Make_probability()

	var prob_map [2][2][2][2]float64 = diff.Map_probagbility(seq)
	SeedSet_F := make([]int, len(adj))
	max_user := 0 //最もフォロワ数が多いユーザ名
	max_user_num := 0
	user_num_counter := 0
	for i := 0; i < len(adj); i++ {
		user_num_counter = 0
		for j := 0; j < len(adj); j++ {
			if adj[i][j] == 1 {
				user_num_counter++
			}
		}
		if max_user_num < user_num_counter {
			max_user = i
			max_user_num = user_num_counter
		}
	}
	SeedSet_F[max_user] = 1

	opt.SameImpressionCostFollower(100, adj, SeedSet_F, prob_map, pop_list, interest_list, assum_list, 15, 16, exit_f, use_cost_infl)
}

func main() {
	use_cost_infl := true
	fmt.Println("start cal_sameimpression_cost.go")
	rand.Seed(int64(1))
	var seed int64 = 1

	adjFilePath := "Graphs/adj_json1000node.txt"
	adjFilePath = "adj_jsonTwitterInteractionUCongress.txt"
	adj, interest_list, assum_list := Make_adj_interest_assum(adjFilePath, seed)
	// fmt.Println(len(adj))
	// os.Exit(0)
	user_same(adj, interest_list, assum_list, true, use_cost_infl)
	// follower_same(adj,interest_list,assum_list,true)

}
