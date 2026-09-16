package experiment

import (
	diff "difftools/diffusion"
	"difftools/network"
	opt "difftools/optimization"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func MakeSeedSetFString2(
	// adj [][]int, nNodes int,
	net network.Network,
	S_f_type int,
) ([]int, []int, int) {
	// ユーザの初期状態
	// 偽情報の発信源の変数
	SeedSet_F_strong2 := make([]int, net.N)
	// 虚偽情報の発信源が与えられたユーザーIDのリスト
	// 虚偽情報の発信源を選択されないようにする(単一情報で)
	// 単一情報用に複数情報で偽情報の発信源が発信源にならないように
	usersFSeeded := make([]int, 1)
	// 最もフォロワ数が多いユーザ名
	mostFollowedUser := 0

	//虚偽情報の発信源を定義
	switch S_f_type {
	case 1:
		//単独ユーザの場合

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
				mostFollowedUser = i
				max_user_num = user_num_counter
			}
		}
		// SeedSet_F_strong2[mostFollowedUser] = diff.InfoType_F + 1 //虚偽情報の発信源を定義
		SeedSet_F_strong2[mostFollowedUser] = diff.SeedInfoF //虚偽情報の発信源を定義
		usersFSeeded[0] = mostFollowedUser

	case 2:
		//複数ユーザの場合

		num2 := 0
		num3 := 0
		for focus_user, slice := range net.Adj {
			num := 0
			for _, edge := range slice {
				num += edge
				if edge > 1 {
					//多重辺がない設定ではこれはエラー
					fmt.Println("error")
					os.Exit(0)
				}
			}

			if num > 20 && num < 30 {
				// if num2%2 == 0 { //個数調整 ego-twitter用
				if num2%20 == 0 { //個数調整 congress用
					SeedSet_F_strong2[focus_user] = 1 //虚偽情報の発信源を定義
					if num3 == 0 {
						usersFSeeded[0] = focus_user
					} else {
						usersFSeeded = append(usersFSeeded, focus_user)
					}
					num3++
				}
				num2++
			}
		}
	}
	return SeedSet_F_strong2, usersFSeeded, mostFollowedUser
}

func ComputeMaximizationDP(
	// adj [][]int,
	// nNodes int,
	net network.Network,
	interestList [][]int,
	assumList [][]int,
	probTable diff.UserProbTable,
	user_weight float64,
	capacity float64,
	use_kaiki bool,
	use_user bool,
	use_infl bool,
	use_follower bool,
	nick int,
	S_f_type int,
	only_last bool,
	r *rand.Rand,
) ([]int, [2]int) {
	// var n int = 50
	// var seesd int64 = 1
	// var K_F int = 5
	// var K_T int = 10
	// var sample_size int = 1000
	//初期化
	var pop_list [2]int
	pop_list[0] = diff.Pop_high
	pop_list[1] = diff.Pop_high

	// fmt.Println(string(bytes))

	fmt.Println("--------------------")

	// var SeedSet_F []int = diff.Make_seedSet_F(n, 1, seed, adj)
	//確率マッピングの作成
	// var seq [16]float64 = diff.MakeProbability()
	// var prob_map diff.UserProbTable = diff.GetUserProbTable()

	// //初期設定
	// //ユーザの初期状態　偽情報の発信源の変数
	// SeedSet_F_strong2 := make([]int, len(adj))
	// //虚偽情報の発信源を選択されないようにする(単一情報で)　単一情報用に複数情報で偽情報の発信源が発信源にならないように
	// non_use_list := make([]int, 1)
	// //最もフォロワ数が多いユーザ名
	// max_user := 0

	// //虚偽情報の発信源を定義
	// switch S_f_type {
	// case 1:
	// 	//単独ユーザの場合

	// 	max_user_num := 0
	// 	user_num_counter := 0
	// 	for i := 0; i < len(adj); i++ {
	// 		user_num_counter = 0
	// 		for j := 0; j < len(adj); j++ {
	// 			if adj[i][j] == 1 {
	// 				user_num_counter++
	// 			}
	// 		}
	// 		if max_user_num < user_num_counter {
	// 			max_user = i
	// 			max_user_num = user_num_counter
	// 		}
	// 	}
	// 	SeedSet_F_strong2[max_user] = 1 //虚偽情報の発信源を定義
	// 	non_use_list[0] = max_user

	// case 2:
	// 	//複数ユーザの場合

	// 	num2 := 0
	// 	num3 := 0
	// 	for focus_user, slice := range adj {
	// 		num := 0
	// 		for _, edge := range slice {
	// 			num += edge
	// 			if edge > 1 {
	// 				//多重辺がない設定ではこれはエラー
	// 				fmt.Println("error")
	// 				os.Exit(0)
	// 			}
	// 		}

	// 		if num > 20 && num < 30 {
	// 			// if num2%2 == 0 { //個数調整 ego-twitter用
	// 			if num2%20 == 0 { //個数調整 congress用
	// 				SeedSet_F_strong2[focus_user] = 1 //虚偽情報の発信源を定義
	// 				if num3 == 0 {
	// 					non_use_list[0] = focus_user
	// 				} else {
	// 					non_use_list = append(non_use_list, focus_user)
	// 				}
	// 				num3++
	// 			}
	// 			num2++
	// 		}
	// 	}
	// }

	SeedSet, usersFSeeded, mostFollowedUser := MakeSeedSetFString2(net, S_f_type)

	//人数を流動的にして拡散を調べている
	//	総フォロワー数を固定できていない

	//拡散可能なユーザ数を調べている
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

	fmt.Println("start_DP")
	// greedy_ans1, _, _ := opt.Greedy(0,100,adj,SeedSet_F_strong2, prob_map,pop_list,interest_list,assum_list,5,true,1000)

	cost_sum := 0
	// for j:=0;j<len(greedy_ans1);j++{
	// 	cost_sum += opt.Cal_cost_kaiki(user_weight,1-user_weight,adj, greedy_ans1[j], max_user)
	// }
	// fmt.Println("cost_sum",cost_sum)
	// cost_sum = 0
	// os.Exit(0)

	s := time.Now()
	//虚偽情報アリの影響最大化問題の解を求める
	//複数情報の影響最大化問題をとく
	// DP_ans2 := make([][]int,0)
	if !only_last {
		//only_last の場合複数情報の影響最大化問題は求めない
		DP_ans, _ := opt.DP(
			100,
			net,
			SeedSet,
			probTable,
			pop_list,
			interestList,
			assumList,
			infler_num,
			true,
			capacity,
			mostFollowedUser,
			true,
			user_weight,
			use_kaiki,
			use_follower,
			nick,
			usersFSeeded,
			use_user,
			use_infl,
			r,
		)
		fmt.Println("DP_time:", time.Since(s))

		// DP_ans := DP_user_infl.Users
		//コストの算出
		for j := 0; j < len(DP_ans); j++ {
			if use_infl {
				cost_sum += opt.Cal_cost_infl_int(net, DP_ans[j], probTable, pop_list, interestList, assumList)
			} else if use_follower {
				cost_sum += opt.Cal_cost_follower_int(
					user_weight, 1-user_weight, net, DP_ans[j], mostFollowedUser)
			} else {
				cost_sum += opt.Cal_cost_infl_int(net, DP_ans[j], probTable, pop_list, interestList, assumList)
			}
		}
		DP_ans2 := make([][]int, 0)
		DP_ans2 = append(DP_ans2, DP_ans)
		// SeedSet_F_strong2 = make([]int, len(adj))
		// SeedSet_F_strong2[max_user] = 1
		_, test_DP_ans_v, test_DP_ans_fv := opt.SelectedSuppressionMaximum(
			net,
			DP_ans2, SeedSet, probTable, pop_list, interestList, assumList, r)
		// SeedSet_F_strong2 = make([]int, len(adj))//念のため初期化
		// SeedSet_F_strong2[max_user] = 1

		fmt.Println("虚偽情報アリの解", DP_ans, test_DP_ans_v, test_DP_ans_fv)
		nonF_SeedSet := make([]int, net.N)

		_, test_DP_ans_v, test_DP_ans_fv = opt.SelectedSuppressionMaximum(
			net,
			DP_ans2, nonF_SeedSet, probTable, pop_list, interestList, assumList, r)

		fmt.Println("虚偽情報アリの解を無しに使ってみたら...", test_DP_ans_v, test_DP_ans_fv)
	}

	// fmt.Println(greedy_ans_v)
	fmt.Println("cost_sum:", cost_sum)

	//単一情報の影響最大化問題の解を求める
	nonF_SeedSet := make([]int, net.N) //念のため初期化　偽情報の発信源が無いとき用のからのリスト
	DP_ans, _ := opt.DP(
		100,
		net,
		nonF_SeedSet,
		probTable,
		pop_list,
		interestList,
		assumList,
		infler_num,
		true,
		capacity,
		mostFollowedUser,
		true,
		user_weight,
		use_kaiki,
		use_follower,
		nick,
		usersFSeeded,
		use_user,
		use_infl,
		r,
	)

	fmt.Println("DP_time:", time.Since(s))

	// DP_ans := DP_user_infl.Users
	//コストの算出
	cost_sum = 0
	for j := 0; j < len(DP_ans); j++ {
		cost_sum += opt.Cal_cost_infl_int(net, DP_ans[j], probTable, pop_list, interestList, assumList)
	}
	DP_ans2 := make([][]int, 0)
	DP_ans2 = append(DP_ans2, DP_ans)

	nonF_SeedSet = make([]int, net.N) //念のため初期化
	_, test_DP_ans_v, test_DP_ans_fv := opt.SelectedSuppressionMaximum(
		net, DP_ans2, nonF_SeedSet, probTable, pop_list, interestList, assumList, r)

	fmt.Println("虚偽情報なしの解", DP_ans, test_DP_ans_v, test_DP_ans_fv)
	// fmt.Println(greedy_ans_v)
	// fmt.Println(test_greedy_ans_v)
	fmt.Println("cost_sum:", cost_sum)

	_, test_DP_ans_v, test_DP_ans_fv = opt.SelectedSuppressionMaximum(
		net, DP_ans2, SeedSet, probTable, pop_list, interestList, assumList, r)

	fmt.Println("虚偽情報なしの解をアリに使ってみたら...", test_DP_ans_v, test_DP_ans_fv)

	return SeedSet, pop_list
}

// use_DP in sample_check_submodu_sub.go
func ComputeDP(
	// adj [][]int,
	net network.Network,
	interest_list [][]int,
	assum_list [][]int,
	probTable diff.UserProbTable,
	user_weight float64,
	capacity float64,
	use_kaiki bool,
	use_user bool,
	use_infl bool,
	nick int,
	S_f_type int,
	only_last bool,
	DP_ans_d []int,
	DP_ans_s []int,
	r *rand.Rand,
) ([]int, [2]int) {

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

	// var seq [16]float64 = diff.MakeProbability()

	var prob_map diff.UserProbTable = diff.GetUserProbTable()

	SeedSet_F_strong2 := make([]int, net.N) //ユーザの初期状態
	non_use_list := make([]int, 1)          //虚偽情報の発信源を選択されないようにする(単一情報で)
	max_user := 0                           //最もフォロワ数が多いユーザ名

	//虚偽情報の発信源を定義
	switch S_f_type {
	case 1:

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
		SeedSet_F_strong2[max_user] = 1 //虚偽情報の発信源を定義
		non_use_list[0] = max_user
	case 2:
		num2 := 0
		num3 := 0
		for focus_user, slice := range net.Adj {
			num := 0
			for _, edge := range slice {
				num += edge
				if edge > 1 {
					fmt.Println("error")
					os.Exit(0)
				}
			}
			if num > 20 && num < 30 {
				if num2%20 == 0 { //個数調整
					SeedSet_F_strong2[focus_user] = 1 //虚偽情報の発信源を定義
					if num3 == 0 {
						non_use_list[0] = focus_user
					} else {
						non_use_list = append(non_use_list, focus_user)
					}
					num3++
				}
				num2++
			}
		}
	}

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

	fmt.Println("start_DP")
	// greedy_ans1, _, _ := opt.Greedy(0,100,adj,SeedSet_F_strong2, prob_map,pop_list,interest_list,assum_list,5,true,1000)

	cost_sum := 0
	// for j:=0;j<len(greedy_ans1);j++{
	// 	cost_sum += opt.Cal_cost_kaiki(user_weight,1-user_weight,adj, greedy_ans1[j], max_user)
	// }
	// fmt.Println("cost_sum",cost_sum)
	// cost_sum = 0
	// os.Exit(0)

	s := time.Now()
	//虚偽情報アリの影響最大化問題の解を求める
	// DP_ans2 := make([][]int,0)
	if !only_last {

		// DP_ans := DP_user_infl.Users

		for j := 0; j < len(DP_ans_d); j++ {
			cost_sum += opt.Cal_cost_infl_int(net, DP_ans_d[j], prob_map, pop_list, interest_list, assum_list)
		}
		DP_ans2 := make([][]int, 0)
		DP_ans2 = append(DP_ans2, DP_ans_d)
		// SeedSet_F_strong2 = make([]int, len(adj))
		// SeedSet_F_strong2[max_user] = 1
		_, test_DP_ans_v, test_DP_ans_fv := opt.SelectedSuppressionMaximum(
			net, DP_ans2, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list, r)
		// SeedSet_F_strong2 = make([]int, len(adj))//念のため初期化
		// SeedSet_F_strong2[max_user] = 1

		fmt.Println("虚偽情報アリの解", DP_ans_d, test_DP_ans_v, test_DP_ans_fv)
		nonF_SeedSet := make([]int, net.N)

		_, test_DP_ans_v, test_DP_ans_fv = opt.SelectedSuppressionMaximum(
			net, DP_ans2, nonF_SeedSet, prob_map, pop_list, interest_list, assum_list, r)

		fmt.Println("虚偽情報アリの解を無しに使ってみたら...", test_DP_ans_v, test_DP_ans_fv)
	}

	// fmt.Println(greedy_ans_v)
	fmt.Println("cost_sum:", cost_sum)
	nonF_SeedSet := make([]int, net.N) //念のため初期化

	fmt.Println("DP_time:", time.Since(s))

	// DP_ans := DP_user_infl.Users

	cost_sum = 0
	for j := 0; j < len(DP_ans_s); j++ {
		cost_sum += opt.Cal_cost_infl_int(net, DP_ans_s[j], prob_map, pop_list, interest_list, assum_list)

	}
	DP_ans2 := make([][]int, 0)
	DP_ans2 = append(DP_ans2, DP_ans_s)

	nonF_SeedSet = make([]int, net.N) //念のため初期化
	_, test_DP_ans_v, test_DP_ans_fv := opt.SelectedSuppressionMaximum(
		net, DP_ans2, nonF_SeedSet, prob_map, pop_list, interest_list, assum_list, r)

	fmt.Println("虚偽情報なしの解", DP_ans_s, test_DP_ans_v, test_DP_ans_fv)
	// fmt.Println(greedy_ans_v)
	// fmt.Println(test_greedy_ans_v)
	fmt.Println("cost_sum:", cost_sum)

	_, test_DP_ans_v, test_DP_ans_fv = opt.SelectedSuppressionMaximum(
		net, DP_ans2, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list, r)

	fmt.Println("虚偽情報なしの解をアリに使ってみたら...", test_DP_ans_v, test_DP_ans_fv)

	return SeedSet_F_strong2, pop_list
}

func CalMaxUsers(
	// adj [][]int,
	net network.Network,
	nPicks int,
) {
	// max_user := 0 //最もフォロワ数が多いユーザ名
	m := nPicks - 1
	max_users := make([]int, nPicks)
	// max_user_num := 0
	user_num_counter := 0
	max_user_nums := make([]int, nPicks)
	for i := 0; i < net.N; i++ {
		user_num_counter = 0
		for l := 0; l < net.N; l++ {
			if net.Adj[i][l] == 1 {
				user_num_counter++
			}
		}

		for j := 0; j < nPicks; j++ {
			if max_user_nums[j] < user_num_counter {

				for k := j; k < m; k++ {
					max_users[m-k+j] = max_users[m-k+j-1]
					max_user_nums[m-k+j] = max_user_nums[m-k+j-1]
					// fmt.Println(i,max_user_nums)
				}
				max_users[j] = i
				max_user_nums[j] = user_num_counter
				// fmt.Println(i,max_user_nums)

				break
			}
		}
	}
	fmt.Println(max_users, max_user_nums)
}

func CalMaxUsersFixed(
	// adj [][]int,
	net network.Network,
	nPicks int,
) {
	if nPicks == 0 || net.N != nPicks {
		fmt.Println("Invalid input: n or adj dimensions are incorrect.")
		return
	}

	// ユーザーID (index) を格納するスライス
	// 最終的にフォロワー数が多い順にソートされます
	max_users := make([]int, nPicks)
	// 各ユーザーのフォロワー数を格納するスライス
	// 最終的にフォロワー数の降順でソートされます
	max_user_nums := make([]int, nPicks)

	// 各ユーザー (i) についてループ
	for i := 0; i < nPicks; i++ {
		follower_count := 0 // user i のフォロワー数

		// [修正点 1] フォロワー数を計算 (in-degree)
		// adj[l][i] == 1 は、ユーザー l が ユーザー i をフォローしていることを意味します
		for l := 0; l < nPicks; l++ {
			// adj が n x n の正方行列であると仮定
			if len(net.Adj[l]) == nPicks && net.Adj[l][i] == 1 {
				follower_count++
			}
		}

		// 挿入ソート: follower_count を
		// max_user_nums の正しい位置に挿入します
		for j := 0; j < nPicks; j++ {
			// ユーザー i のフォロワー数が、
			// 現在 j 番目のユーザーのフォロワー数より多い場合
			if max_user_nums[j] < follower_count {

				// [修正点 2] 挿入ロジックを明確化
				// ユーザー i を j 番目に挿入するため、
				// j 番目以降の要素を
				// 後ろに1つずつずらします (n-1 から j+1 まで)
				for k := nPicks - 1; k > j; k-- {
					max_users[k] = max_users[k-1]
					max_user_nums[k] = max_user_nums[k-1]
				}

				// j 番目に user i とそのフォロワー数を挿入
				max_users[j] = i
				max_user_nums[j] = follower_count
				break // 挿入したので内側の j ループを抜ける
			}
		}
	}

	fmt.Println("Ranking (User ID):", max_users)
	fmt.Println("Ranking (Follower #):", max_user_nums)
}

func SimSubmod(
	sample_size int,
	// adj [][]int,
	net network.Network,
	pop_list [2]int, interest_list [][]int,
	assum_list [][]int, SeedSet_F []int,
	K_T int,
	prob_map diff.UserProbTable,
	folder_path string,
	r *rand.Rand,
) ([]int, [][]float64) {
	var S []int
	var hist [][]float64
	S, hist = opt.CheckSubmod(
		K_T, sample_size, net, SeedSet_F, prob_map, pop_list, interest_list, assum_list, folder_path, r)

	return S, hist
}
