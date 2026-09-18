// single_dp.go
package main

import (
	diff "difftools/diffusion"
	exp "difftools/experiment"
	"difftools/network"
	opt "difftools/optimization"
	"encoding/json"
	"fmt"
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
	AdjFilePath string `json:"adj_file_path"`
	SFType      int    `json:"s_f_type"`
	// UseUser           bool    `json:"use_user"`
	// UseInfl           bool    `json:"use_infl"`
	UseCongress bool `json:"use_congress"`
	// UseKaiki          bool    `json:"use_kaiki"`
	// UseFollower       bool    `json:"use_follower"`
	CostFunc          opt.CostFunc `json:"cost_func"`
	NumPickUsers      int          `json:"num_pick_users"`
	UserWeightInitial float64      `json:"user_weight_initial"`
}

type Task struct {
	Seed     int64   `json:"seed"`
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
		// rand.Seed(int64(task.Seed))

		// イテレーションごとに初期シードを変えるべきか？
		r := rand.New(rand.NewSource(task.Seed))

		// 2. データの生成 (Make_adj...)
		// adj, interest_list, assum_list := Make_adj_interest_assum(settings.AdjFilePath, int64(task.Seed))
		net := network.ReadJson(settings.AdjFilePath)
		interestList := diff.MakeInterestList(net.N, r)
		assumList := diff.MakeAssumList(net.N, r)
		probTable := diff.GetUserProbTable()

		// 3. 初回のみ実行する処理 (cal_max_users)
		// 元のコードの if i == 0 に相当
		if i == 0 {
			exp.CalMaxUsers(net, settings.NumPickUsers)
		}

		// 4. 計算実行 (compute_maximization_DP)
		// JSONから読み込んだ値を渡す
		exp.ComputeMaximizationDP(
			net,
			interestList,
			assumList,
			probTable,
			settings.UserWeightInitial,
			task.Capacity, // JSONの計算済み値
			settings.CostFunc,
			// settings.UseKaiki,
			// settings.UseUser,
			// settings.UseInfl,
			// settings.UseFollower,
			1, // 元コードの固定値
			settings.SFType,
			false, // 元コードの固定値
			r,
		)
	}

	fmt.Printf("<< FINISHED: %s\n", batch.GroupID)
}
