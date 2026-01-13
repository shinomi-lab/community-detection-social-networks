package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	diff "m/difftools/diffusion"
	opt "m/difftools/optimization"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ---------------------------------------------------------
// 1. 設定読み込み用の構造体 (生成プログラムと同じもの)
// ---------------------------------------------------------

type GlobalSettings struct {
	AdjFilePath       string  `json:"adj_file_path"`
	SFType            int     `json:"s_f_type"`
	UseUser           bool    `json:"use_user"`
	UseInfl           bool    `json:"use_infl"`
	UseCongress       bool    `json:"use_congress"`
	UseKaiki          bool    `json:"use_kaiki"`
	UseFollower       bool    `json:"use_follower"`
	NumPickUsers      int     `json:"num_pick_users"`
	UserWeightInitial float64 `json:"user_weight_initial"`
}

type Task struct {
	Seed     int     `json:"seed"`
	ScaleJ   float64 `json:"scale_j"`
	Capacity float64 `json:"capacity"`
}

type ExperimentBatch struct {
	GroupID  string         `json:"experiment_group_id"`
	Settings GlobalSettings `json:"settings"`
	Tasks    []Task         `json:"tasks"`
}

// ---------------------------------------------------------
// 2. メイン実行関数 (Runner)
// ---------------------------------------------------------

func main() {
	// JSONファイルが保存されているディレクトリ
	configDir := "./configs_final"

	// ディレクトリ内の全JSONファイルを取得
	files, err := filepath.Glob(filepath.Join(configDir, "*.json"))
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	if len(files) == 0 {
		fmt.Println("JSON file not found in", configDir)
		return
	}

	fmt.Printf("Found %d config files. Starting simulation...\n", len(files))

	// --- 並列実行の設定 ---
	var wg sync.WaitGroup
	// ★ 同時実行数を制限するためのチャネル (例: 同時に4ファイルまで)
	// PCのスペックに合わせて数字を変更してください (4~8程度推奨)
	maxConcurrent := 4
	sem := make(chan struct{}, maxConcurrent)

	startTotal := time.Now()

	for _, file := range files {
		wg.Add(1)

		// Goroutine起動
		go func(fPath string) {
			defer wg.Done()

			// セマフォを取得（満員ならここで待機）
			sem <- struct{}{}

			// 処理実行
			runExperimentBatch(fPath)

			// セマフォを解放
			<-sem
		}(file)
	}

	// 全ての処理が終わるのを待つ
	wg.Wait()

	fmt.Printf("\nAll experiments completed in %v\n", time.Since(startTotal))
}

// ---------------------------------------------------------
// 3. ファイルごとの処理ロジック
// ---------------------------------------------------------

func runExperimentBatch(filePath string) {
	// JSON読み込み
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return
	}

	var batch ExperimentBatch
	if err := json.Unmarshal(bytes, &batch); err != nil {
		fmt.Printf("Error parsing JSON %s: %v\n", filePath, err)
		return
	}

	settings := batch.Settings
	fmt.Printf(">> START: %s (Tasks: %d)\n", batch.GroupID, len(batch.Tasks))

	// --- タスク（Seed/Scale）のループ ---
	for i, task := range batch.Tasks {
		// 進捗表示 (省略可)
		// fmt.Printf("   [%s] processing seed:%d scale:%.1f\n", batch.GroupID, task.Seed, task.ScaleJ)

		// 1. シード値の設定
		// 注意: 並列処理の場合、rand.Seed(Global)は他の並列処理に影響する可能性があります。
		// 厳密に行う場合は rand.New(rand.NewSource(...)) を各関数に渡すべきですが、
		// ここでは元のコードの仕様に合わせています。
		rand.Seed(int64(task.Seed))

		// 2. データの生成 (Make_adj...)
		adj, interest_list, assum_list := Make_adj_interest_assum(settings.AdjFilePath, int64(task.Seed))

		// 3. 初回のみ実行する処理 (cal_max_users)
		// 元のコードの if i == 0 に相当
		if i == 0 {
			cal_max_users(adj, settings.NumPickUsers)
		}

		// 4. 計算実行 (compute_maximization_DP)
		// JSONから読み込んだ値を渡す
		compute_maximization_DP(
			adj,
			interest_list,
			assum_list,
			settings.UserWeightInitial,
			task.Capacity, // JSONの計算済み値
			settings.UseKaiki,
			settings.UseUser,
			settings.UseInfl,
			settings.UseFollower,
			1, // 元コードの固定値
			settings.SFType,
			false, // 元コードの固定値
		)
	}

	fmt.Printf("<< FINISHED: %s\n", batch.GroupID)
}

// ---------------------------------------------------------
// 4. 既存関数のプレースホルダー (ここにあなたの関数を貼り付けてください)
// ---------------------------------------------------------

// ダミー構造体（型合わせ用）
// あなたのコードに合わせて変更してください
type AdjType []int  // 仮
type ListType []int // 仮

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

func cal_max_users(adj [][]int, n int) {
	// max_user := 0 //最もフォロワ数が多いユーザ名
	m := n - 1
	max_users := make([]int, n)
	// max_user_num := 0
	user_num_counter := 0
	max_user_nums := make([]int, n)
	for i := 0; i < len(adj); i++ {
		user_num_counter = 0
		for l := 0; l < len(adj); l++ {
			if adj[i][l] == 1 {
				user_num_counter++
			}
		}

		for j := 0; j < n; j++ {
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

func compute_maximization_DP(adj [][]int, interest_list [][]int, assum_list [][]int, user_weight float64, capacity float64, use_kaiki bool, use_user bool, use_infl bool, use_follower bool, nick int, S_f_type int, only_last bool) ([][]int, []int, [2][2][2][2]float64, [2]int, [][]int, [][]int) {
	// ★ここに実際の compute_maximization_DP の中身を実装してください
	//初期化
	var pop_list [2]int
	pop_list[0] = diff.Pop_high
	pop_list[1] = diff.Pop_high

	fmt.Println("--------------------")

	//確率マッピングの作成
	var seq [16]float64 = diff.Make_probability()
	var prob_map [2][2][2][2]float64 = diff.Map_probagbility(seq)

	//初期設定
	SeedSet_F_strong2 := make([]int, len(adj)) //ユーザの初期状態　偽情報の発信源の変数
	non_use_list := make([]int, 1)             //虚偽情報の発信源を選択されないようにする(単一情報で)　単一情報用に複数情報で偽情報の発信源が発信源にならないように
	max_user := 0                              //最もフォロワ数が多いユーザ名

	//虚偽情報の発信源を定義
	if S_f_type == 1 {
		//単独ユーザの場合

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
		SeedSet_F_strong2[max_user] = 1 //虚偽情報の発信源を定義
		non_use_list[0] = max_user

	} else if S_f_type == 2 {
		//複数ユーザの場合

		num2 := 0
		num3 := 0
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
				if num2%20 == 0 { //個数調整 congress用
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

	//拡散可能なユーザ数を調べている
	infler_num := 0
	// OnlyInfler := true
	for j := 0; j < len(adj); j++ {
		for k := 0; k < len(adj); k++ {
			if adj[j][k] != 0 {
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
		DP_ans, _ := opt.DP(0, 100, adj, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list, infler_num, true, capacity, max_user, true, user_weight, use_kaiki, use_follower, nick, non_use_list, use_user, use_infl)
		fmt.Println("DP_time:", time.Since(s))

		// DP_ans := DP_user_infl.Users
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
		// SeedSet_F_strong2 = make([]int, len(adj))
		// SeedSet_F_strong2[max_user] = 1
		_, test_DP_ans_v, test_DP_ans_fv := opt.Selected_Suppression_Maximum(adj, DP_ans2, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list)
		// SeedSet_F_strong2 = make([]int, len(adj))//念のため初期化
		// SeedSet_F_strong2[max_user] = 1

		fmt.Println("虚偽情報アリの解", DP_ans, test_DP_ans_v, test_DP_ans_fv)
		nonF_SeedSet := make([]int, len(adj))

		_, test_DP_ans_v, test_DP_ans_fv = opt.Selected_Suppression_Maximum(adj, DP_ans2, nonF_SeedSet, prob_map, pop_list, interest_list, assum_list)

		fmt.Println("虚偽情報アリの解を無しに使ってみたら...", test_DP_ans_v, test_DP_ans_fv)
	}

	// fmt.Println(greedy_ans_v)
	fmt.Println("cost_sum:", cost_sum)

	//単一情報の影響最大化問題の解を求める
	nonF_SeedSet := make([]int, len(adj)) //念のため初期化　偽情報の発信源が無いとき用のからのリスト
	DP_ans, _ := opt.DP(0, 100, adj, nonF_SeedSet, prob_map, pop_list, interest_list, assum_list, infler_num, true, capacity, max_user, true, user_weight, use_kaiki, use_follower, nick, non_use_list, use_user, use_infl)

	fmt.Println("DP_time:", time.Since(s))

	// DP_ans := DP_user_infl.Users
	//コストの算出
	cost_sum = 0
	for j := 0; j < len(DP_ans); j++ {
		cost_sum += opt.Cal_cost_infl_int(adj, DP_ans[j], prob_map, pop_list, interest_list, assum_list)
	}
	DP_ans2 := make([][]int, 0)
	DP_ans2 = append(DP_ans2, DP_ans)

	nonF_SeedSet = make([]int, len(adj)) //念のため初期化
	_, test_DP_ans_v, test_DP_ans_fv := opt.Selected_Suppression_Maximum(adj, DP_ans2, nonF_SeedSet, prob_map, pop_list, interest_list, assum_list)

	fmt.Println("虚偽情報なしの解", DP_ans, test_DP_ans_v, test_DP_ans_fv)
	// fmt.Println(greedy_ans_v)
	// fmt.Println(test_greedy_ans_v)
	fmt.Println("cost_sum:", cost_sum)

	_, test_DP_ans_v, test_DP_ans_fv = opt.Selected_Suppression_Maximum(adj, DP_ans2, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list)

	fmt.Println("虚偽情報なしの解をアリに使ってみたら...", test_DP_ans_v, test_DP_ans_fv)

	return adj, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list
	// 重い処理のシミュレーション（動作確認用のWait）
	// time.Sleep(100 * time.Millisecond)
}
