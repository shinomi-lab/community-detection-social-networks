package main

import (
	exp "difftools/experiment"
	opt "difftools/optimization"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
	strconv "strconv"
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
		log.Fatal("Usage: go run main.go <input_file_path> <output_dir>")
	}

	// read input paramter JSON file into param object
	// JSON Format: list of EvalParameter objects
	var param EvalParameter
	configFilePath := os.Args[1]
	paramFile, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(paramFile, &param)

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

	// job worker
	numWorkers := runtime.NumCPU()
	results := make(chan result, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, iter := range param.Iterations {
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
					iter.Seed,
					param.GraphPath,
					// 影響関数計算用
					param.SampleSize,
					// 偽情報発信源
					param.FalseUsers,
					// DP用
					1,
					iter.Capacity,
					costFunc,
					//targetCommunityNodes, // 抽出したコミュニティノードリストを渡す
				)
				results <- result{
					trueOnly:      trueOnly,
					trueWithFalse: trueWithFalse,
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		row := []string{
			strconv.FormatFloat(result.trueOnly, 'f', -1, 64),
			strconv.FormatFloat(result.trueWithFalse, 'f', -1, 64),
		}
		err = writer.Write(row)
		if err != nil {
			log.Fatal(err)
		}
	}

	logFile.WriteString("Input parameter file: " + configFilePath + "\n")
	logFile.WriteString("Output results file: " + csvFilePath + "\n")
}
