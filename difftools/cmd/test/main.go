package main

import (
	diff "difftools/diffusion"
	exp "difftools/experiment"
	"difftools/network"
	opt "difftools/optimization"
	"math/rand"
	"path/filepath"
	"runtime"
)

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("ソースファイルのパスを取得できませんでした")
	}
	mainDir := filepath.Dir(filename)

	r := rand.New(rand.NewSource(int64(100)))

	net := network.ReadJson(filepath.Join(mainDir, "sample_adj.json"))
	interestList := diff.MakeInterestList(net.N, r)
	assumList := diff.MakeAssumList(net.N, r)
	probTable := diff.GetUserProbTable()
	nick := 1
	S_f_type := 2

	costFunc := opt.CostDefault
	capacity := 500.0
	userWeight := 1.0

	SeedSet, usersFSeeded, mostFollowedUser := exp.MakeSeedSetFString2(net, S_f_type)

	popList := make(diff.PopList, diff.Pops_n)
	popList[diff.InfoType_F] = diff.PopHigh
	popList[diff.InfoType_T] = diff.PopHigh

	DP_ans, _ := opt.DP(
		100,
		net,
		SeedSet,
		probTable,
		popList,
		interestList,
		assumList,
		true,
		capacity,
		mostFollowedUser,
		true,
		userWeight,
		costFunc,
		nick,
		usersFSeeded,
		r)

	DP_ans2 := make([][]int, 0)
	DP_ans2 = append(DP_ans2, DP_ans)
	_, test_DP_ans_v, test_DP_ans_fv := opt.SelectedSuppressionMaximum(
		net,
		DP_ans2,
		[]diff.SeedInfo{},
		probTable,
		popList,
		interestList,
		assumList,
		r)

	println(test_DP_ans_v, test_DP_ans_fv)
}
