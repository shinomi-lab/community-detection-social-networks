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
	seedSet []int,
	userProbTable diff.UserProbTable,
	pop [2]int,
	interestList [][]int,
	assumList [][]int,
	r *rand.Rand,
) []float64 {
	// n := len(adj)

	// dist := make([][]int, diff.InfoTypes_n)
	ans := make([]float64, diff.InfoTypes_n)

	for i := 0; i < sampleSize; i++ {
		var dist = diff.Diffuse(net, seedSet, userProbTable, pop, interestList, assumList, r) //-1 is correct?
		ans[diff.InfoType_F] += float64(len(dist[diff.InfoType_F])) / float64(sampleSize)
		ans[diff.InfoType_T] += float64(len(dist[diff.InfoType_T])) / float64(sampleSize)
	}
	return ans
}
