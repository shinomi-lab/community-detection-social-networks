package optimization

// diffuse.goを使ってモンテカルロ法で影響関数を求めている

import (
	diff "difftools/diffusion"
	"difftools/network"
	"math/rand"
)

// 乱数シード値の代わりに乱数生成器`r`を渡す
//
// return: result of mont (配列で影響関数の答えがinfoごとにある)
func RunInflProp(
	sampleSize int,
	net network.Network,
	// adj [][]int,
	seedSet []diff.SeedInfo,
	userProbTable diff.UserProbTable,
	pop diff.PopList,
	interestList diff.InterestList,
	assumList diff.AssumList,
	r *rand.Rand,
) map[diff.InfoType]float64 {

	// n := len(adj)

	// dist := make([][]int, diff.InfoTypes_n)
	ans := make(map[diff.InfoType]float64, diff.InfoTypes_n)
	sumF := 0
	sumT := 0

	for i := 0; i < sampleSize; i++ {
		dist := diff.Diffuse(net, seedSet, userProbTable, pop, interestList, assumList, r)
		sumF += len(dist[diff.InfoType_F])
		sumT += len(dist[diff.InfoType_T])
	}
	ans[diff.InfoType_F] = float64(sumF) / float64(sampleSize)
	ans[diff.InfoType_T] = float64(sumT) / float64(sampleSize)

	/*
		for i := 0; i < sampleSize; i++ {
			var dist = diff.Diffuse(net, seedSet, userProbTable, pop, interestList, assumList, r) //-1 is correct?
			ans[diff.InfoType_F] += float64(len(dist[diff.InfoType_F])) / float64(sampleSize)
			ans[diff.InfoType_T] += float64(len(dist[diff.InfoType_T])) / float64(sampleSize)
		}
	*/

	return ans
}
