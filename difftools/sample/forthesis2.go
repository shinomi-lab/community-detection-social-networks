func compute_maximization_DP(adj [][]int, interest_list [][]int, assum_list [][]int, user_weight float64, capacity float64, use_kaiki bool, use_user bool, use_infl bool, use_follower bool, nick int, S_f_type int, only_last bool) ([][]int, []int, [2][2][2][2]float64, [2]int, [][]int, [][]int) {
    // 初期化
    var pop_list [2]int
    pop_list[0] = diff.Pop_high
    pop_list[1] = diff.Pop_high

    fmt.Println("--------------------")

    // 確率マッピングの作成
    var seq [16]float64 = diff.Make_probability()
    var prob_map [2][2][2][2]float64 = diff.Map_probagbility(seq)

    // 初期設定
    SeedSet_F_strong2 := make([]int, len(adj)) // 虚偽情報の発信源の変数
    non_use_list := make([]int, 1)            // 単一情報用に複数情報で虚偽情報の発信源が発信源にならないように
    max_user := 0                             // 最大フォロワ数を持つユーザ

    // 虚偽情報の発信源の定義
    if S_f_type == 1 { // 単独ユーザの場合
        max_user_num := 0
        for i := 0; i < len(adj); i++ {
            user_num_counter := 0
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
        SeedSet_F_strong2[max_user] = 1
        non_use_list[0] = max_user
    } else if S_f_type == 2 { // 複数ユーザの場合
        num2, num3 := 0, 0
        for focus_user, slice := range adj {
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
                if num2%20 == 0 {//個数調整 congress用
                    SeedSet_F_strong2[focus_user] = 1
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

    //拡散可能な(出次数が0より大きい)ユーザ数の計算
    infler_num := 0
    for j := 0; j < len(adj); j++ {
        for k := 0; k < len(adj); k++ {
            if adj[j][k] != 0 {
                infler_num++
                break
            }
        }
    }

    fmt.Println("start_DP")

    cost_sum := 0
    s := time.Now()

    // 複数情報の影響最大化問題の解を求める
    if !only_last {//only_lastの場合複数情報の影響最大化問題は求めない
        DP_ans, _ := opt.DP(0, 100, adj, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list, infler_num, true, capacity, max_user, true, user_weight, use_kaiki, use_follower, nick, non_use_list, use_user, use_infl)
        fmt.Println("DP_time:", time.Since(s))
        //コストの算出
        for j := 0; j < len(DP_ans); j++ {
            if use_infl {
                cost_sum += opt.Cal_cost_infl_int(adj, DP_ans[j], prob_map, pop_list, interest_list, assum_list)
            } else if use_follower {
                cost_sum += opt.Cal_cost_follower_int(user_weight, 1-user_weight, adj, DP_ans[j], max_user)
            } else {
                cost_sum += opt.Cal_cost_infl_int(adj, DP_ans[j], prob_map, pop_list, interest_list, assum_list)
            }
        }
        DP_ans2 := make([][]int, 0)
        DP_ans2 = append(DP_ans2, DP_ans)

        _, test_DP_ans_v, test_DP_ans_fv := opt.Selected_Suppression_Maximum(adj, DP_ans2, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list)
        fmt.Println("虚偽情報アリの解", DP_ans, test_DP_ans_v, test_DP_ans_fv)

        nonF_SeedSet := make([]int, len(adj))nonF_SeedSet := make([]int, len(adj)) //虚偽情報の発信源が無いとき用のからのリスト
        _, test_DP_ans_v, test_DP_ans_fv = opt.Selected_Suppression_Maximum(adj, DP_ans2, nonF_SeedSet, prob_map, pop_list, interest_list, assum_list)
        fmt.Println("虚偽情報アリの解を無しに使ってみたら...", test_DP_ans_v, test_DP_ans_fv)
        fmt.Println("cost_sum:", cost_sum)
    }

    // 単一情報の影響最大化問題の解を求める
    nonF_SeedSet := make([]int, len(adj))nonF_SeedSet := make([]int, len(adj)) //虚偽情報の発信源が無いとき用のからのリスト
    DP_ans, _ := opt.DP(0, 100, adj, nonF_SeedSet, prob_map, pop_list, interest_list, assum_list, infler_num, true, capacity, max_user, true, user_weight, use_kaiki, use_follower, nick, non_use_list, use_user, use_infl)
    fmt.Println("DP_time:", time.Since(s))

    //コストの算出
    cost_sum = 0
    for j := 0; j < len(DP_ans); j++ {
        if use_infl {
            cost_sum += opt.Cal_cost_infl_int(adj, DP_ans[j], prob_map, pop_list, interest_list, assum_list)
        } else if use_follower {
            cost_sum += opt.Cal_cost_follower_int(user_weight, 1-user_weight, adj, DP_ans[j], max_user)
        } else {
            cost_sum += opt.Cal_cost_infl_int(adj, DP_ans[j], prob_map, pop_list, interest_list, assum_list)
        }
    }
    DP_ans2 := make([][]int, 0)
    DP_ans2 = append(DP_ans2, DP_ans)

    _, test_DP_ans_v, test_DP_ans_fv := opt.Selected_Suppression_Maximum(adj, DP_ans2, nonF_SeedSet, prob_map, pop_list, interest_list, assum_list)
    fmt.Println("虚偽情報なしの解", DP_ans, test_DP_ans_v, test_DP_ans_fv)

    _, test_DP_ans_v, test_DP_ans_fv = opt.Selected_Suppression_Maximum(adj, DP_ans2, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list)
    fmt.Println("虚偽情報なしの解をアリに使ってみたら...", test_DP_ans_v, test_DP_ans_fv)

}
