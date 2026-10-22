package main

import (
	exp "difftools/experiment"
	opt "difftools/optimization"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	strconv "strconv"
	"time"
)

type EvalParameter struct {
	GraphPath  string      `json:"graph_path"`
	FalseUsers []int       `json:"false_users"`
	CostFunc   string      `json:"cost_func"` // "follower" | "influence"
	SampleSize int         `json:"sample_size"`
	Iterations []iteration `json:"iterations"`
}

type iteration struct {
	Seed     int64 `json:"seed"`
	Capacity int   `json:"capacity"`
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
	json.Unmarshal(paramFile, &param)

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

	for _, iter := range param.Iterations {
		var costFunc opt.CostFunc
		switch param.CostFunc {
		case "follower":
			costFunc = opt.CostFollower{}
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
		)
		row := []string{
			strconv.FormatFloat(trueOnly, 'f', -1, 64),
			strconv.FormatFloat(trueWithFalse, 'f', -1, 64),
		}
		err = writer.Write(row)
		if err != nil {
			log.Fatal(err)
		}
	}
	logFile.WriteString("Input parameter file: " + configFilePath + "\n")
	logFile.WriteString("Output results file: " + csvFilePath + "\n")
}
