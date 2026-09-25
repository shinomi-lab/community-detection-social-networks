package optimization

// 影響関数を使って影響最大化問題を解いてる

import (
	"bufio"
	"bytes"
	diff "difftools/diffusion"
	"difftools/network"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
)

func Greedy(
	sample_size int,
	// adj [][]int,
	net network.Network,
	Seed_set []diff.SeedInfo,
	prob_map diff.UserProbTable,
	pop diff.PopList,
	interest_list diff.InterestList,
	assum_list diff.AssumList,
	ans_len int,
	Count_true bool,
	sample_size2 int,
	r *rand.Rand,
) ([]int, float64, []float64) {
	//sample_size2はグリーディで求めた解をより詳しくやる
	var n int = net.N // len(adj)
	var max float64 = 0
	var result float64
	var index int
	var ans []int
	var ans_v []float64

	ans = make([]int, 0, ans_len)
	S := make([]diff.SeedInfo, len(Seed_set))
	_ = copy(S, Seed_set)
	S_test := make([]diff.SeedInfo, len(Seed_set))
	_ = copy(S_test, Seed_set)

	var info_num diff.SeedInfo = diff.BoolToSeed(Count_true)
	// if Count_true {
	// 	info_num = diff.SeedInfoT // 2
	// } else {
	// 	info_num = diff.SeedInfoF // 1
	// }

	for i := 0; i < ans_len; i++ {
		fmt.Println(i)
		max = 0
		for j := 0; j < n; j++ {
			if (j+1)%100 == 0 {
				fmt.Println(i, "-", (j+1)/100)
			}
			_ = copy(S_test, S)
			if S_test[j] != 0 { //すでに発信源のユーザだったら
				continue
			}
			S_test[j] = info_num

			dist := RunInflProp(sample_size, net, S_test, prob_map, pop, interest_list, assum_list, r)
			if Count_true {
				result = dist[diff.InfoType_T]
			} else {
				result = dist[diff.InfoType_F]
			}

			if result > max {
				max = result
				index = j
			}
		} //subloop end

		ans = append(ans, index)
		ans_v = append(ans_v, max)
		S[index] = info_num

	} //mainloop end

	// var max_2 float64
	// dist2 := Infl_prop_exp(seed, sample_size2, adj, S, prob_map, pop, interest_list, assum_list)
	// if Count_true {
	// 	max_2 = dist2[diff.InfoType_T]
	// } else {
	// 	max_2 = dist2[diff.InfoType_F]
	// }
	return ans, max, ans_v
}

type User_Infl struct {
	users []int
	infl  float64
}

type ByInfl []User_Infl

func (a ByInfl) Len() int           { return len(a) }
func (a ByInfl) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByInfl) Less(i, j int) bool { return a[i].infl < a[j].infl }

func Greedy_exp(
	sample_size int,
	// adj [][]int,
	net network.Network,
	Seed_set []diff.SeedInfo,
	prob_map diff.UserProbTable,
	pop diff.PopList,
	interest_list diff.InterestList,
	assum_list diff.AssumList,
	ans_len int,
	Count_true bool,
	capacity float64,
	max_user int,
	OnlyInfler bool,
	user_weight float64,
	use_kaiki bool,
	r *rand.Rand,
) ([]int, float64) {

	// var costcal func(float64, float64,[][]int,int,int) float64
	// if use_kaiki{
	// 	costcal = Cal_cost_kaiki
	// }else{
	// 	costcal = Cal_cost
	// }
	var n int = net.N // len(adj)
	var max float64 = 0
	var result float64
	var index int
	var ans []int

	ans = make([]int, 0, ans_len)
	S := make([]diff.SeedInfo, len(Seed_set))
	_ = copy(S, Seed_set)
	S_test := make([]diff.SeedInfo, len(Seed_set))
	_ = copy(S_test, Seed_set)
	cap_use := capacity
	pre_infl := 0.0

	var info_num diff.SeedInfo = diff.BoolToSeed(Count_true)
	// if Count_true {
	// 	info_num = diff.SeedInfoT // 2
	// } else {
	// 	info_num = diff.SeedInfoF // 1
	// }

	for {
		max = -1
		for j := 0; j < n; j++ {

			_ = copy(S_test, S)             //初期化
			for i := 0; i < len(ans); i++ { //初期設定
				S_test[ans[i]] = info_num
			}
			if S_test[j] != 0 { //すでに発信源のユーザだったら
				continue
			}
			if OnlyInfler {
				// if FolowerSize(adj, j) == 0 {
				if net.Followers[j] == 0 {
					continue
				}
			}

			// cost := costcal(user_weight, 1-user_weight, adj, j, max_user)
			cost := cal_cost_infl(net, j, prob_map, pop, interest_list, assum_list)
			if cost > cap_use { //コストが大きすぎるユーザなら
				continue
			}
			S_test[j] = info_num
			// rand.Seed(100) //おそらく後で消す　重要
			dist := RunInflProp(sample_size, net, S_test, prob_map, pop, interest_list, assum_list, r)
			if Count_true {
				result = (dist[diff.InfoType_T] - pre_infl) / cost
			} else {
				result = (dist[diff.InfoType_F] - pre_infl) / cost
			}

			if result > max {
				pre_infl = dist[diff.InfoType_T]
				max = result
				index = j
			}
		} //subloop end

		if max == -1 {
			break
		}
		ans = append(ans, index)
		cap_use -= cal_cost_infl(net, index, prob_map, pop, interest_list, assum_list)
		// cap_use -= costcal(user_weight,1-user_weight,adj,index,max_user)

		S[index] = info_num
	}
	return ans, max
}

func Cal_cost(
	u_weight float64,
	f_wight float64,
	// adj [][]int,
	net network.Network,
	node int,
	max_user int,
) float64 {
	// f := FolowerSize(adj, node)
	f := net.Followers[node]
	if f == 0 {
		f = 10000000000000
	}
	f_max := net.Followers[max_user]
	// u_max := net.N // len(adj)
	// f_max := 0
	// for i := 0; i < u_max; i++ {
	// 	if net.Adj[max_user][i] == 1 {
	// 		f_max++
	// 	}
	// }
	// fmt.Println("folowersize",f)
	// fmt.Println("f:",f,"f_wight:",f_wight,"f_max:",f_max,"u_weight:",u_weight)
	return float64(f)*f_wight/float64(f_max) + u_weight/2
}

func Cal_cost_kaiki(
	u_weight float64,
	f_wight float64,
	// adj [][]int,
	net network.Network,
	node int,
	max_user int,
) float64 {
	// return 100.0
	// f := FolowerSize(adj, node)
	f := net.Followers[node]
	if f == 0 {
		f = 10000000000000
	}
	file, err := os.Open("kaiki.txt")
	if err != nil {
		fmt.Println("ファイルを開く際にエラーが発生しました:", err)
		return -1
	}
	defer file.Close()

	// ファイルを読み込む
	scanner := bufio.NewScanner(file)
	var line string
	if scanner.Scan() {
		line = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("ファイルを読み取る際にエラーが発生しました:", err)
		return -1
	}

	// コンマで分割して整数に変換する
	parts := strings.Split(line, ",")
	if len(parts) != 2 {
		fmt.Println("予期しないデータ形式です")
		return -1
	}

	intercept, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		fmt.Println("最初の部分を整数に変換する際にエラーが発生しました:", err)
		return -1
	}

	slope, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		fmt.Println("二番目の部分を整数に変換する際にエラーが発生しました:", err)
		return -1
	}
	// fmt.Println(math.Log(float64(f))*slope +intercept)
	return math.Log(float64(f))*slope + intercept
}

func Cal_cost_kaiki_int(
	u_weight float64,
	f_wight float64,
	// adj [][]int,
	net network.Network,
	node int,
	max_user int,
) int {
	return int(math.Round(Cal_cost_kaiki(u_weight, f_wight, net, node, max_user)))
}

func cal_cost_infl(
	// adj [][]int,
	net network.Network,
	node int,
	prob_map diff.UserProbTable,
	pop diff.PopList,
	interest_list diff.InterestList,
	assum_list diff.AssumList,
) float64 {

	S_test := make([]diff.SeedInfo, net.N)
	// S_test[node] = 2
	S_test[node] = diff.SeedInfoT
	r := rand.New(rand.NewSource((100))) // [todo]
	dist := RunInflProp(1000, net, S_test, prob_map, pop, interest_list, assum_list, r)

	return dist[diff.InfoType_T]
}

func Cal_cost_infl_int(
	// adj [][]int,
	net network.Network,
	node int,
	prob_map diff.UserProbTable,
	pop diff.PopList,
	interest_list diff.InterestList,
	assum_list diff.AssumList,
) int {
	cost := cal_cost_infl(net, node, prob_map, pop, interest_list, assum_list)
	return int(math.Round(cost))
}

func Cal_cost_user(
	// u_weight float64, f_wight float64, adj [][]int, node int, max_user int,
	u_weight float64,
	f_wight float64,
	net network.Network,
	node int,
	max_user int,
) float64 {
	return 1.0
}

func Cal_cost_user_int(
	// u_weight float64, f_wight float64, adj [][]int, node int, max_user int,
	u_weight float64,
	f_wight float64,
	net network.Network,
	node int,
	max_user int,
) int {
	return 1
}
func Cal_cost_follower(
	// u_weight float64, f_wight float64, adj [][]int, node int, max_user int,
	u_weight float64,
	f_wight float64,
	net network.Network,
	node int,
	max_user int,
) float64 {
	// return float64(FolowerSize(adj, node))
	return float64(net.Followers[node])
}

func Cal_cost_follower_int(
	// u_weight float64, f_wight float64, adj [][]int, node int, max_user int,
	u_weight float64,
	f_wight float64,
	net network.Network,
	node int,
	max_user int,
) int {
	// return FolowerSize(adj, node)
	return net.Followers[node]
}

// 付録Cはここから
type Users_infl struct {
	Infl  float64
	Users []int
} //構造体を定義

func (ui *Users_infl) AddUser(user int) {
	ui.Users = append(ui.Users, user)
}

func (ui *Users_infl) CopyUsers(users []int) {
	ui.Users = make([]int, len(users))
	copy(ui.Users, users)
}

type CostFunc int

const (
	CostDefault CostFunc = iota
	CostKaiki
	CostUser
	CostFollower
	CostInfluence
)

var costFuncToString = map[CostFunc]string{
	CostDefault:   "default",
	CostKaiki:     "kaiki",
	CostUser:      "user",
	CostFollower:  "follower",
	CostInfluence: "influence",
}

var stringToCostFunc = map[string]CostFunc{
	"default":   CostDefault,
	"kaiki":     CostKaiki,
	"user":      CostUser,
	"follower":  CostFollower,
	"influence": CostInfluence,
}

func (f CostFunc) String() string {
	if str, ok := costFuncToString[f]; ok {
		return str
	}
	return "default"
}

func (s CostFunc) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (s *CostFunc) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	if val, ok := stringToCostFunc[str]; ok {
		*s = val
	} else {
		*s = CostDefault // 定義外の文字列はCostDefault
	}
	return nil
}

// countTure: 新情報を扱うかどうか
//
// nick: 刻み幅，大きいほど計算時間が短く精度が悪くなる
//
// OnlyInfler: 次数を持たないユーザも扱うかどうか
//
// costFunc: CostKaiki | CostFollower | CostUser | CostInfluence | CostWeight
func DP(
	sampleSize int,
	net network.Network,
	seedSet []diff.SeedInfo,
	probTable diff.UserProbTable,
	pop diff.PopList,
	interestList diff.InterestList,
	assumList diff.AssumList,
	countTrue bool,
	capacity float64,
	maxUser int,
	onlyInfler bool,
	userWeight float64,
	costFunc CostFunc,
	nick int,
	nonUseList []int,
	r *rand.Rand,
	// ansLen int,
	// useKaiki bool,
	// useFollower bool,
	// useUser bool,
	// useInfl bool,
) ([]int, float64) {
	var info_num diff.SeedInfo = diff.BoolToSeed(countTrue)

	//制約関数の設定
	var costInflFunc func(focusUser int) int

	// use_int_cost := false
	// var costcal func(float64, float64, network.Network, int, int) float64
	// var costcal_int func(float64, float64, network.Network, int, int) int
	// if useKaiki {
	// 	// use_int_cost = true
	// 	// costcal = Cal_cost_kaiki
	// 	// costcal_int = Cal_cost_kaiki_int
	// } else if useUser {
	// 	use_int_cost = true
	// 	// costcal = Cal_cost_user
	// 	// costcal_int = Cal_cost_user_int
	// } else if useFollower {
	// 	// use_int_cost = true
	// 	// costcal = Cal_cost_follower
	// 	// costcal_int = Cal_cost_follower_int
	// } else {
	// 	use_int_cost = false
	// 	costcal = Cal_cost
	// }

	switch costFunc {
	case CostDefault:
		costInflFunc = func(focusUser int) int {
			return int(Cal_cost(userWeight, 1-userWeight, net, focusUser, maxUser))
		}
	case CostKaiki:
		costInflFunc = func(focusUser int) int {
			return Cal_cost_kaiki_int(userWeight, 1-userWeight, net, focusUser, maxUser)
		}
	case CostUser:
		costInflFunc = func(focusUser int) int {
			return 1
		}
	case CostFollower:
		costInflFunc = func(focusUser int) int {
			return net.Followers[focusUser]
		}
	case CostInfluence:
		costInflFunc = func(focusUser int) int {
			return Cal_cost_infl_int(net, focusUser, probTable, pop, interestList, assumList)
		}
	}

	S := make([]diff.SeedInfo, len(seedSet))
	_ = copy(S, seedSet)

	onlyiflerlist := OnlyInflerlist(net, nonUseList)
	onlyinfler_num := len(onlyiflerlist)

	n := onlyinfler_num
	l_list := int(int(capacity)/nick) + 1

	// dp := make([][]float64,n+1) // <- dpは構造体にして拡散量と既に選ばれているユーザ集合を入れる
	dp := make([][]Users_infl, n+1)

	for i := 0; i < n+1; i++ {
		// dp[i] = make([]float64,l_list)
		dp[i] = make([]Users_infl, l_list)
	}

	for w := range l_list {
		dp[0][w].Infl = 0
	}
	for i := range n {
		focusUser := onlyiflerlist[i]
		// if !use_int_cost { //cost_inflだけ引数の数が違うため別処理
		// 	if useInfl {
		// 		cost_i = cal_cost_infl(net, focus_user, probMap, pop, interestList, assumList)
		// 	} else {
		// 		cost_i = costcal(userWeight, 1-userWeight, net, focus_user, maxUser)
		// 	}
		// 	cost_i_int = int(cost_i)
		// } else {
		// 	if useInfl {
		// 		cost_i_int = Cal_cost_infl_int(net, focus_user, probMap, pop, interestList, assumList)
		// 	} else {
		// 		cost_i_int = costcal_int(userWeight, 1-userWeight, net, focus_user, maxUser)
		// 	}
		// }
		cost_i_int := costInflFunc(focusUser)

		for j := range l_list {
			// cost_w := j*nick
			//dp[i+1][j]に代入していく i番目までを選べるコストj*nick以下
			if j < cost_i_int/nick { //大きすぎると不可能
				dp[i+1][j].Infl = dp[i][j].Infl
				dp[i+1][j].CopyUsers(dp[i][j].Users)
				continue
			}
			_ = copy(S, seedSet) //初期化
			last_cost := j - cost_i_int/nick
			for k := 0; k < len(dp[i][last_cost].Users); k++ {
				S[dp[i][last_cost].Users[k]] = info_num
			}
			S[focusUser] = info_num
			// rand.Seed(100) //おそらく後で消す　重要
			dist := RunInflProp(sampleSize, net, S, probTable, pop, interestList, assumList, r)
			// if countTrue {
			// 	result = dist[diff.InfoType_T]
			// } else {
			// 	result = dist[diff.InfoType_F]
			// }
			result := dist[diff.BoolToInfo(countTrue)]

			//resultをdp[i][j].users+onlyiflerlist[i]での拡散にする複数回同じ拡散を調べたくないけど，一旦後回し？
			if dp[i][j].Infl < result {

				dp[i+1][j].Infl = roundTo(result, 4)
				dp[i+1][j].CopyUsers(dp[i][last_cost].Users)
				dp[i+1][j].AddUser(focusUser)
			} else {
				dp[i+1][j].Infl = dp[i][j].Infl
				dp[i+1][j].CopyUsers(dp[i][j].Users)
			}
		}
	}
	return dp[n][l_list-1].Users, dp[n][l_list-1].Infl
}

// func PrintDp(dp [][]Users_infl) {
// 	for i := 0; i < len(dp); i++ {
// 		fmt.Println(dp[i])
// 	}
// }

func OnlyInflerlist(
	// adj [][]int,
	net network.Network,
	non_use_list []int) []int {
	// n := len(adj)
	ans := make([]int, 0, net.N)

	for i := 0; i < net.N; i++ {
		// if FolowerSize(adj, i) != 0 {
		if net.Followers[i] != 0 {
			// if !isInList(i, non_use_list) {
			if !slices.Contains(non_use_list, i) {
				ans = append(ans, i)
			}
		}
	}
	// fmt.Println(ans)
	return ans
}

func isInList(n int, list []int) bool {
	return slices.Contains(list, n)
	// for _, k := range list {
	// 	if k == n {
	// 		return true
	// 	}
	// }
	// return false
}

func roundTo(n float64, precision int) float64 {
	scale := math.Pow(10, float64(precision))
	return math.Round(n*scale) / scale
}
