type Users_infl struct {
    Infl float64
    Users  []int
}//構造体を定義

func (ui *Users_infl) AddUser(user int) {
    ui.Users = append(ui.Users,user)
}

func (ui *Users_infl) CopyUsers(users []int){
	ui.Users = make([]int,len(users))
	copy(ui.Users,users)
}

func DP(seed int64, sample_size int, adj [][]int, Seed_set []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int, ans_len int, Count_true bool, capacity float64, max_user int, OnlyInfler bool, user_weight float64, use_kaiki bool, use_follower bool, nick int, non_use_list []int, use_user bool, use_infl bool) ([]int, float64) {//nick は刻み幅，大きいほど計算時間が短く精度が悪くなる
	var info_num, cost_i_int int
	var result, cost_i float64
	use_int_cost := false

	if Count_true {
		info_num = 2
	} else {
		info_num = 1
	}

	var costcal func(float64, float64, [][]int, int, int) float64
	var costcal_int func(float64, float64, [][]int, int, int) int
  // 制約関数の設定
	if use_kaiki {
		use_int_cost = true
		costcal = Cal_cost_kaiki
		costcal_int = Cal_cost_kaiki_int
	} else if use_user {
		use_int_cost = true
		costcal = Cal_cost_user
		costcal_int = Cal_cost_user_int
	} else if use_follower {
		use_int_cost = true
		costcal = Cal_cost_follower
		costcal_int = Cal_cost_follower_int
	} else {
		use_int_cost = false
		costcal = Cal_cost
	}

	S := make([]int, len(Seed_set))
	copy(S, Seed_set)

	onlyiflerlist := OnlyInflerlist(adj, non_use_list)
	onlyinfler_num := len(onlyiflerlist)
	n := onlyinfler_num
	l_list := int(capacity/float64(nick)) + 1

	dp := make([][]Users_infl, n+1)
	for i := range dp {
		dp[i] = make([]Users_infl, l_list)
	}
	for w := 0; w < l_list; w++ {
		dp[0][w].Infl = 0
	}

	for i := 0; i < n; i++ {
		focus_user := onlyiflerlist[i]
		if !use_int_cost {
			if use_infl {//cost_inflだけ引数の数が違うので別処理
				cost_i = Cal_cost_infl(adj, focus_user, prob_map, pop, interest_list, assum_list)
			} else {
				cost_i = costcal(user_weight, 1-user_weight, adj, focus_user, max_user)
			}
			cost_i_int = int(cost_i)
		} else {
			if use_infl {
				cost_i_int = Cal_cost_infl_int(adj, focus_user, prob_map, pop, interest_list, assum_list)
			} else {
				cost_i_int = costcal_int(user_weight, 1-user_weight, adj, focus_user, max_user)
			}
		}

		for j := 0; j < l_list; j++ {
			if j < cost_i_int/nick {//大きすぎると不可能
				dp[i+1][j].Infl = dp[i][j].Infl
				dp[i+1][j].CopyUsers(dp[i][j].Users)
				continue
			}

			copy(S, Seed_set)
			last_cost := j - cost_i_int/nick
			for _, u := range dp[i][last_cost].Users {
				S[u] = info_num
			}
			S[focus_user] = info_num
			rand.Seed(100)
			dist := Infl_prop_exp(seed, sample_size, adj, S, prob_map, pop, interest_list, assum_list)

			if Count_true {
				result = dist[diff.InfoType_T]
			} else {
				result = dist[diff.InfoType_F]
			}

			if dp[i][j].Infl < result {
				dp[i+1][j].Infl = roundTo(result, 4)
				dp[i+1][j].CopyUsers(dp[i][last_cost].Users)
				dp[i+1][j].AddUser(focus_user)
			} else {
				dp[i+1][j].Infl = dp[i][j].Infl
				dp[i+1][j].CopyUsers(dp[i][j].Users)
			}
		}
	}
	return dp[n][l_list-1].Users, dp[n][l_list-1].Infl
}
