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
	net := network.ReadAdjJson(filepath.Join(mainDir, "sample_adj.json"))

	// 偽情報発信源
	falseUsers := []int{}

	// 影響関数計算用
	sampleSize := 1000

	// DP用
	nick := 1
	capacity := 500.0

	// 制約関数依存
	costFunc := opt.CostUser
	userWeight := 1.0
	mostFollowedUser := 0 // exp.FindMostFollowedUser(net)

	// nFollowersFrom := 21
	// nFollowersTo := 31
	// skip := 20
	// falseUsers := exp.ChooseMultipleUsers(net, nFollowersFrom, nFollowersTo, skip)

	// ここから共通処理
	popList := make(diff.PopList, diff.Pops_n)
	popList[diff.InfoType_F] = diff.PopHigh
	popList[diff.InfoType_T] = diff.PopHigh

	interestList := diff.MakeInterestList(net.N, r)
	assumList := diff.MakeAssumList(net.N, r)
	probTable := diff.GetUserProbTable()

	falseSeedSet := exp.UsersToSeedSet(falseUsers, net, diff.SeedInfoF)

	trueUsers, _ := opt.DP(
		sampleSize,
		net,
		falseSeedSet,
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
		falseUsers,
		r)

	testTrueSeedSet := exp.UsersToSeedSet(trueUsers, net, diff.SeedInfoT)
	dist := opt.RunInflProp(
		sampleSize, net, testTrueSeedSet, probTable, popList, interestList, assumList, r,
	)
	// test_DP_ans_tv, test_DP_ans_fv
	println(dist[diff.InfoType_T], dist[diff.InfoType_F])
}
