package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ---------------------------------------------------------
// 1. JSON出力用の構造体定義
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
// 2. ループ処理用の設定構造体
// ---------------------------------------------------------

type GenMode struct {
	Label   string // ファイル名に使う "User" や "Infl"
	UseUser bool
	UseInfl bool
}

// ---------------------------------------------------------
// 3. メイン処理
// ---------------------------------------------------------

func main() {
	// 保存先ディレクトリ
	outputDir := "./settings"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}

	// ==========================================
	// ★ 変動パターンの定義
	// ==========================================

	// 入力ファイルリスト
	adjFiles := []string{
		"community_20_adjmat.txt", // 必要に応じて追加
	}

	// SFタイプのリスト
	sfTypes := []int{1, 2}

	// モード定義: 排他制御 (User=True/Infl=False) と (User=False/Infl=True)
	modes := []GenMode{
		{Label: "User", UseUser: true, UseInfl: false},
		{Label: "Infl", UseUser: false, UseInfl: true},
	}

	// 内部ループ用 (Seed, Scale)
	seeds := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	scales := []float64{1.0, 2.0, 3.0, 4.0}

	// 固定値
	useCongress := true
	useFollower := false
	useKaiki := false
	numPickUsers := 7
	userWeight := 0.0

	count := 0

	// ==========================================
	// ★ 生成ループ
	// ==========================================

	for _, adjPath := range adjFiles {
		for _, sf := range sfTypes {
			// 変数名を genMode に統一してエラー回避
			for _, genMode := range modes {

				// 1. タスクリスト(Tasks)の作成
				var tasks []Task
				for _, seed := range seeds {
					for _, j := range scales {

						// Capacityの計算
						var capacity float64
						// genMode.UseInfl を使って判定
						if genMode.UseInfl && useCongress {
							capacity = j * 100
						} else if genMode.UseInfl && !useCongress {
							capacity = j * 5
						} else {
							// Userモードの場合 (Infl=false)
							capacity = j
						}

						tasks = append(tasks, Task{
							Seed:     seed,
							ScaleJ:   j,
							Capacity: capacity,
						})
					}
				}

				// 2. ファイル名の生成
				// 拡張子を除去 (例: community_3_adjmat)
				baseName := strings.TrimSuffix(filepath.Base(adjPath), filepath.Ext(adjPath))

				// 指定形式: [ファイル名]_[User/Infl]_sf[数字]
				// 例: community_3_adjmat_User_sf1
				fileID := fmt.Sprintf("%s_%s_sf%d", baseName, genMode.Label, sf)
				fileName := fmt.Sprintf("%s.json", fileID)

				// 3. データ構造の結合
				batchConfig := ExperimentBatch{
					GroupID: fileID,
					Settings: GlobalSettings{
						AdjFilePath:       adjPath,
						SFType:            sf,
						UseUser:           genMode.UseUser,
						UseInfl:           genMode.UseInfl,
						UseCongress:       useCongress,
						UseKaiki:          useKaiki,
						UseFollower:       useFollower,
						NumPickUsers:      numPickUsers,
						UserWeightInitial: userWeight,
					},
					Tasks: tasks,
				}

				// 4. ファイル書き出し
				writeJSON(filepath.Join(outputDir, fileName), batchConfig)
				count++
			}
		}
	}

	fmt.Printf("生成完了: %d 個のJSONファイルを %s に作成しました。\n", count, outputDir)
}

func writeJSON(path string, data interface{}) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		fmt.Println("Error encoding:", err)
	}
}
