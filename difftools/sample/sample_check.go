package sample

import (
	diff "difftools/diffusion"
	"difftools/network"
	opt "difftools/optimization"
	"fmt"
	"math/rand"
	"time"
)

func Make_adj_interest_assum(adjFilePath string, r *rand.Rand) (network.Network, [][]int, [][]int) {
	// bytes, err := os.ReadFile(adjFilePath)
	// if err != nil {
	// 	panic(err)
	// }

	// // fmt.Println(string(bytes))

	// var dataJson string = string(bytes)

	// arr := make(map[int]map[int]int)
	// // var arr []string
	// _ = json.Unmarshal([]byte(dataJson), &arr)
	// // fmt.Println(arr)

	// // fmt.Println(arr[0][1])

	net := network.ReadJson(adjFilePath)
	n := net.N

	var interest_list [][]int = diff.MakeInterestList(n, r)
	var assum_list [][]int = diff.MakeAssumList(n, r)
	// var adj [][]int = make([][]int, n)

	// for i := 0; i < n; i++ {
	// 	adj[i] = make([]int, n)
	// 	for j := 0; j < n; j++ {
	// 		adj[i][j] = arr[j][i]
	// 	}
	// }
	return net, interest_list, assum_list
}

func use_strict(
	// adj [][]int,
	net network.Network,
	interest_list [][]int,
	assum_list [][]int,
	user_weight float64,
) ([][]int, []int, diff.UserProbTable, [2]int) {

	// var n int = 50
	// var seed int64 = 1
	// var K_F int = 5
	// var K_T int = 10
	// var sample_size int = 1000
	var pop_list [2]int
	pop_list[0] = diff.Pop_high
	pop_list[1] = diff.Pop_high

	// fmt.Println(K_T, K_F, diff.InfoType_F, sample_size, pop_list)
	// adjFilePath := "adj_jsonTwitterInteractionUCongress.txt"
	// adjFilePath := "Graphs/adj_json1000node.txt"
	// result_Path := "Twitter_Data/"

	// var SeedSet_F []int = diff.Make_seedSet_F(n, 1, seed, adj)

	// var interest_list [][]int = diff.Make_interest_list(n, seed)
	//
	// var assum_list [][]int = diff.Make_assum_list(n, seed)

	// var seq [16]float64 = diff.Make_probability()
	var prob_map diff.UserProbTable = diff.GetUserProbTable()

	// fmt.Println("Seedsetf")
	// fmt.Println(SeedSet_F)
	//
	// fmt.Println("prob_map")
	// fmt.Println(prob_map)

	SeedSet_F_strong2 := make([]int, net.N)
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
	SeedSet_F_strong2[max_user] = 1

	//人数を流動的にして拡散を調べている
	//	総フォロワー数を固定できていない
	//拡散可能な人数を調べている
	infler_num := 0
	// OnlyInfler := true
	for j := 0; j < net.N; j++ {
		for k := 0; k < net.N; k++ {
			if net.Adj[j][k] != 0 {
				infler_num += 1
				break
			}
		}
	}

	fmt.Println("start_strict")
	// slice_test := [][]int{{1, 15, 18}, {1, 15, 18},{1, 15, 18},{1, 15, 18},{1, 15, 18}}
	// opt.Selected_Suppression_Maximum(adj,slice_test,SeedSet_F_strong2,  prob_map , pop_list, interest_list, assum_list)
	// os.Exit(0)
	s := time.Now()
	under := 0.0
	upper := 2.0
	selected_list := opt.CallKumiawase2(net, under, upper, SeedSet_F_strong2, true, max_user, user_weight)

	// fmt.Println("end all kumiawase",selected_list)
	cost_sum := 0.0
	for k := 0; k < len(selected_list); k++ {
		selecte := selected_list[k]
		cost_sum = 0
		for l := 0; l < len(selecte); l++ {
			cost_sum += opt.Cal_cost(user_weight, 1.0-user_weight, net, selecte[l], max_user)
		}
		if cost_sum < under || cost_sum > upper {
			fmt.Println("えらーcost_sum:", cost_sum)
		}

	}
	// os.Exit(0)
	r := rand.New(rand.NewSource(0))
	strict_ans, strict_ans_v, strict_ans_fv := opt.SelectedSuppressionMaximum(
		net, selected_list, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list, r)

	fmt.Println("strict_time", time.Since(s))
	cost_sum = 0
	for j := 0; j < len(strict_ans); j++ {
		cost_sum += opt.Cal_cost(0.5, 0.5, net, strict_ans[j], max_user)
	}

	fmt.Println(strict_ans)
	fmt.Println(strict_ans_v)
	fmt.Println(strict_ans_fv)
	fmt.Println("cost_sum:", cost_sum)

	return selected_list, SeedSet_F_strong2, prob_map, pop_list
}

func use_greedy(
	// adj [][]int,
	net network.Network,
	interest_list [][]int, assum_list [][]int,
	user_weight float64, capacity float64,
	) ([]int, diff.UserProbTable, [2]int) {

	// var n int = 50
	// var seesd int64 = 1
	// var K_F int = 5
	// var K_T int = 10
	// var sample_size int = 1000
	var pop_list [2]int
	pop_list[0] = diff.Pop_high
	pop_list[1] = diff.Pop_high

	// fmt.Println(string(bytes))

	fmt.Println("--------------------")

	// var SeedSet_F []int = diff.Make_seedSet_F(n, 1, seed, adj)

	// var interest_list [][]int = diff.Make_interest_list(n, seed)
	//
	// var assum_list [][]int = diff.Make_assum_list(n, seed)

	// var seq [16]float64 = diff.Make_probability()

	prob_map := diff.GetUserProbTable()

	// fmt.Println("Seedsetf")
	// fmt.Println(SeedSet_F)
	//
	// fmt.Println("prob_map")
	// fmt.Println(prob_map)

	SeedSet_F_strong2 := make([]int, net.N)
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
	SeedSet_F_strong2[max_user] = 1

	//人数を流動的にして拡散を調べている
	//	総フォロワー数を固定できていない
	//拡散可能な人数を調べている
	infler_num := 0
	// OnlyInfler := true
	for j := 0; j < net.N; j++ {
		for k := 0; k < net.N; k++ {
			if net.Adj[j][k] != 0 {
				infler_num += 1
				break
			}
		}
	}

	r := rand.New(rand.NewSource(0))

	fmt.Println("start_greedy")
	// greedy_ans1, _, _ := opt.Greedy(0,100,adj,SeedSet_F_strong2, prob_map,pop_list,interest_list,assum_list,5,true,1000)

	cost_sum := 0.0
	// for j:=0;j<len(greedy_ans1);j++{
	// 	cost_sum += opt.Cal_cost_kaiki(user_weight,1-user_weight,adj, greedy_ans1[j], max_user)
	// }
	// fmt.Println("cost_sum",cost_sum)
	// cost_sum = 0
	// os.Exit(0)

	s := time.Now()
	//虚偽情報アリの影響最大化問題の解を求める
	greedy_ans, _ := opt.Greedy_exp(
		100, net, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list,
		infler_num, true, capacity, max_user, true, user_weight, true, r)
	fmt.Println("greedy_time:", time.Since(s))

	for j := 0; j < len(greedy_ans); j++ {
		cost_sum += opt.Cal_cost_kaiki(user_weight, 1-user_weight, net, greedy_ans[j], max_user)
	}
	greedy_ans2 := make([][]int, 0)
	greedy_ans2 = append(greedy_ans2, greedy_ans)
	SeedSet_F_strong2 = make([]int, net.N)
	SeedSet_F_strong2[max_user] = 1
	_, test_greedy_ans_v, test_greedy_ans_fv := opt.SelectedSuppressionMaximum(
		net, greedy_ans2, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list, r)
	SeedSet_F_strong2 = make([]int, net.N) //念のため初期化
	SeedSet_F_strong2[max_user] = 1

	fmt.Println("虚偽情報アリの解", greedy_ans, test_greedy_ans_v, test_greedy_ans_fv)
	nonF_SeedSet := make([]int, net.N)

	_, test_greedy_ans_v, test_greedy_ans_fv = opt.SelectedSuppressionMaximum(
		net, greedy_ans2, nonF_SeedSet, prob_map, pop_list, interest_list, assum_list, r)

	fmt.Println("虚偽情報アリの解を無しに使ってみたら...", test_greedy_ans_v, test_greedy_ans_fv)
	// fmt.Println(greedy_ans_v)
	fmt.Println("cost_sum:", cost_sum)
	nonF_SeedSet = make([]int, net.N) //念のため初期化
	greedy_ans, _ = opt.Greedy_exp(
		100, net, nonF_SeedSet, prob_map, pop_list, interest_list, assum_list,
		infler_num, true, capacity, max_user, true, user_weight, true, r)
	fmt.Println("greedy_time:", time.Since(s))

	cost_sum = 0
	for j := 0; j < len(greedy_ans); j++ {
		cost_sum += opt.Cal_cost_kaiki(user_weight, 1-user_weight, net, greedy_ans[j], max_user)
	}
	greedy_ans2 = make([][]int, 0)
	greedy_ans2 = append(greedy_ans2, greedy_ans)

	nonF_SeedSet = make([]int, net.N) //念のため初期化
	_, test_greedy_ans_v, test_greedy_ans_fv = opt.SelectedSuppressionMaximum(
		net, greedy_ans2, nonF_SeedSet, prob_map, pop_list, interest_list, assum_list, r)

	fmt.Println("虚偽情報なしの解", greedy_ans, test_greedy_ans_v, test_greedy_ans_fv)
	// fmt.Println(greedy_ans_v)
	// fmt.Println(test_greedy_ans_v)
	fmt.Println("cost_sum:", cost_sum)

	_, test_greedy_ans_v, test_greedy_ans_fv = opt.SelectedSuppressionMaximum(
		net, greedy_ans2, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list, r)

	fmt.Println("虚偽情報なしの解をアリに使ってみたら...", test_greedy_ans_v, test_greedy_ans_fv)

	return SeedSet_F_strong2, prob_map, pop_list
}

type Parameter struct {
	GraphPath        string
	Node_n           int
	Random_seed      int64
	K_F              int
	K_T              int
	Mont_sample_size int
	Pop_list         [2]int
	SeedSet_F        []int
	ProbMap          diff.UserProbTable
	InterestList     [][]int
	AssumList        [][]int
	// Seq              [16]float64
}
