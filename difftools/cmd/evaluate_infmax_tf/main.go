package main

import (
	exp "difftools/experiment"
	opt "difftools/optimization"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	strconv "strconv"
	"strings"
	"sync"
	"time"
)

type EvalParameter struct {
	GraphPath string `json:"graph_path"`
	//CommunityPath     string      `json:"community_path"`      // communities.txt (.gz) へのパス
	//TargetCommunityID int         `json:"target_community_id"` // 対象コミュニティID
	FalseUsers []int       `json:"false_users"`
	CostFunc   string      `json:"cost_func"` // "user" | "influence"
	SampleSize int         `json:"sample_size"`
	Iterations []iteration `json:"iterations"`
}

type iteration struct {
	Seed     int64 `json:"seed"`
	Capacity int   `json:"capacity"`
}

type result struct {
	trueOnly      float64
	trueWithFalse float64
}

// どのJSONの、どのIterationかを特定するためのジョブ構造体
type job struct {
	jsonIndex int
	iterIndex int
	param     EvalParameter
	iter      iteration
}

// ワーカーからの計算結果
type indexedResult struct {
	jsonIndex int
	iterIndex int
	res       result
}

// JSON Format Example:
// {
//   "graph_path": "path/to/graph.json",
//   "false_users": [1, 2, 3],
//   "cost_func": "influence",
//   "sample_size": 1000,
//   "iterations": [
//     {
//       "seed": 42,
//       "capacity": 500,
//     },
//     {
//       "seed": 43,
//       "capacity": 800,
//     },
//   ]
// }

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go <input_file_or_dir_path> <output_dir>")
	}
	inputPath := os.Args[1]
	outputDir := os.Args[2]

	// read input paramter JSON file into param object
	// JSON Format: list of EvalParameter objects
	//var param EvalParameter
	var jsonFiles []string
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		log.Fatalf("入力パスの取得に失敗しました: %v", err)
	}
	if fileInfo.IsDir() {
		// フォルダが指定された場合、内部の *.json ファイルを取得
		files, err := filepath.Glob(filepath.Join(inputPath, "*.json"))
		if err != nil {
			log.Fatalf("JSONファイルの探索に失敗しました: %v", err)
		}
		jsonFiles = files
	} else {
		// 単一のファイルが指定された場合
		jsonFiles = []string{inputPath}
	}

	if len(jsonFiles) == 0 {
		log.Fatal("実行対象のJSONファイルが見つかりませんでした。")
	}

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("出力ディレクトリの作成に失敗しました: %v", err)
	}

	// 2. すべての JSON ファイルから全タスクを集約
	var allJobs []job
	jsonIterCounts := make([]int, len(jsonFiles))

	for jsonIdx, jsonPath := range jsonFiles {
		paramFile, err := os.ReadFile(jsonPath)
		if err != nil {
			log.Fatalf("JSONファイルの読み込みエラー (%s): %v", jsonPath, err)
		}

		var param EvalParameter
		if err := json.Unmarshal(paramFile, &param); err != nil {
			log.Fatalf("JSON Unmarshal エラー (%s): %v", jsonPath, err)
		}

		jsonIterCounts[jsonIdx] = len(param.Iterations)

		for iterIdx, iter := range param.Iterations {
			allJobs = append(allJobs, job{
				jsonIndex: jsonIdx,
				iterIndex: iterIdx,
				param:     param,
				iter:      iter,
			})
		}
	}

	/*
		configFilePath := os.Args[1]
		paramFile, err := os.ReadFile(configFilePath)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(paramFile, &param)
	*/

	/*
		if err != nil {
			log.Fatal("JSON Unmarshal Error: ", err)
		}

		// コミュニティ属性ファイルから指定コミュニティに属するノードリストを取得
		targetCommunityNodes, err := network.ReadCommunityNodes(param.CommunityPath, param.TargetCommunityID)
		if err != nil {
			log.Fatalf("コミュニティファイルの読み込みに失敗しました: %v", err)
		}
		log.Printf("コミュニティ %d の対象ノード数: %d 人\n", param.TargetCommunityID, len(targetCommunityNodes))
	*/

	/*
		// create csv file to output results
		outputDir := os.Args[2]
		date := time.Now().Format("20060102_150405")
		csvFilePath := filepath.Join(outputDir, date+".csv")
		csvFile, err := os.Create(csvFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer csvFile.Close()
		writer := csv.NewWriter(csvFile)
		defer writer.Flush()

		// create log file to output mapping of input parameter file and output results
		logFile, err := os.Create(filepath.Join(outputDir, date+".log"))
		if err != nil {
			log.Fatal(err)
		}
		defer logFile.Close()

		header := []string{"true_only", "true_with_false"}
		if err := writer.Write(header); err != nil {
			log.Fatal(err)
		}
	*/

	// job worker
	//numWorkers := runtime.NumCPU()
	numJobs := len(allJobs)
	jobs := make(chan job, numJobs)
	results := make(chan indexedResult, numJobs)
	for _, j := range allJobs {
		jobs <- j
	}
	close(jobs) // 投入完了後にクローズ

	numWorkers := runtime.NumCPU()
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				var costFunc opt.CostFunc
				switch j.param.CostFunc {
				case "user":
					costFunc = opt.CostUser{}
				case "influence":
					costFunc = opt.CostInfluence{}
				default:
					log.Fatalf("Unknown cost function: %s (JSON: %s)", j.param.CostFunc, jsonFiles[j.jsonIndex])
				}

				trueOnly, trueWithFalse := exp.EvalInfmaxTF(
					j.iter.Seed,
					j.param.GraphPath,
					j.param.SampleSize,
					j.param.FalseUsers,
					1,
					j.iter.Capacity,
					costFunc,
				)

				results <- indexedResult{
					jsonIndex: j.jsonIndex,
					iterIndex: j.iterIndex,
					res: result{
						trueOnly:      trueOnly,
						trueWithFalse: trueWithFalse,
					},
				}
			}
		}()
	}
	/*
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := range jobs {
					var costFunc opt.CostFunc
					switch param.CostFunc {
					case "user":
						costFunc = opt.CostUser{}
					case "influence":
						costFunc = opt.CostInfluence{}
					default:
						log.Fatal("Unknown cost function: " + param.CostFunc)
					}
					trueOnly, trueWithFalse := exp.EvalInfmaxTF(
						j.iter.Seed,
						param.GraphPath,
						// 影響関数計算用
						param.SampleSize,
						// 偽情報発信源
						param.FalseUsers,
						// DP用
						1,
						j.iter.Capacity,
						costFunc,
						//targetCommunityNodes, // 抽出したコミュニティノードリストを渡す
					)
					resultsChan <- indexedResult{
						index: j.index,
						res:   result{trueOnly: trueOnly, trueWithFalse: trueWithFalse},
					}
				}
			}()
		}
	*/

	go func() {
		wg.Wait()
		close(results)
	}()

	// 3. 計算結果を受け取り、元の Iterations の順番どおりに並べ直す
	finalResults := make(map[int]map[int]result)
	for r := range results {
		if _, exists := finalResults[r.jsonIndex]; !exists {
			finalResults[r.jsonIndex] = make(map[int]result)
		}
		finalResults[r.jsonIndex][r.iterIndex] = r.res
	}

	// 6. JSONファイルごとに CSV とログを出力
	date := time.Now().Format("20060102_150405")

	for jsonIdx, jsonPath := range jsonFiles {
		baseName := filepath.Base(jsonPath)
		ext := filepath.Ext(baseName)
		fileNameWithoutExt := strings.TrimSuffix(baseName, ext)

		// CSV ファイル作成 (例: bridge_v_only_20261031_153000.csv)
		csvFileName := fmt.Sprintf("%s_%s.csv", fileNameWithoutExt, date)
		csvFilePath := filepath.Join(outputDir, csvFileName)
		csvFile, err := os.Create(csvFilePath)
		if err != nil {
			log.Fatalf("CSVファイル作成エラー: %v", err)
		}
		writer := csv.NewWriter(csvFile)

		header := []string{"true_only", "true_with_false"}
		if err := writer.Write(header); err != nil {
			log.Fatalf("CSVヘッダー書き込みエラー: %v", err)
		}

		// 元の Iterations の順序通りに書き出し
		iterCount := jsonIterCounts[jsonIdx]
		for iterIdx := 0; iterIdx < iterCount; iterIdx++ {
			res := finalResults[jsonIdx][iterIdx]
			row := []string{
				strconv.FormatFloat(res.trueOnly, 'f', -1, 64),
				strconv.FormatFloat(res.trueWithFalse, 'f', -1, 64),
			}
			if err := writer.Write(row); err != nil {
				log.Fatalf("CSVデータ書き込みエラー: %v", err)
			}
		}
		writer.Flush()
		csvFile.Close()
		/*
			for _, res := range finalResults {
				row := []string{
					strconv.FormatFloat(res.trueOnly, 'f', -1, 64),
					strconv.FormatFloat(res.trueWithFalse, 'f', -1, 64),
				}
				if err := writer.Write(row); err != nil {
					log.Fatal(err)
				}
			}
		*/

		//logFile.WriteString("Input parameter file: " + configFilePath + "\n")
		//logFile.WriteString("Output results file: " + csvFilePath + "\n")
		logFilePath := filepath.Join(outputDir, fmt.Sprintf("%s_%s.log", fileNameWithoutExt, date))
		logFile, err := os.Create(logFilePath)
		if err == nil {
			logFile.WriteString("Input parameter file: " + jsonPath + "\n")
			logFile.WriteString("Output results file: " + csvFilePath + "\n")
			logFile.Close()
		}

		log.Printf("完了: %s -> %s", baseName, csvFilePath)
	}
}
