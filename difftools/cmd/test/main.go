package main

import (
	exp "difftools/experiment"
	opt "difftools/optimization"
	"path/filepath"
	"runtime"
)

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("ソースファイルのパスを取得できませんでした")
	}
	mainDir := filepath.Dir(filename)

	trueOnly, trueWithFalse := exp.EvalInfmaxTF(
		100,
		filepath.Join(mainDir, "sample_adj.json"),
		// 影響関数計算用
		1000,
		// 偽情報発信源
		[]int{},
		// DP用
		1,
		500,
		opt.CostUser{},
	)
	println(trueOnly, trueWithFalse)

}
