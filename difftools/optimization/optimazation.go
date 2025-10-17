package optimization

// diffuse.goを使ってモンテカルロ法で影響関数を求めている

import (
	diff "m/difftools/diffusion"
)

func Infl_prop_exp(seed int64, sample_size int, adj [][]int, Seed_set []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int) []float64 {
	// return value is result of mont (配列で影響関数の答えがinfoごとにある)
	// n := len(adj)

	_ = seed
	dist := make([][]int, diff.InfoTypes_n)
	ans := make([]float64, diff.InfoTypes_n)

	for i := 0; i < sample_size; i++ {
		dist = diff.Adjmat(adj, Seed_set, -1, prob_map, pop, interest_list, assum_list) //-1 is correct?
		ans[diff.InfoType_F] += float64(len(dist[diff.InfoType_F])) / float64(sample_size)
		ans[diff.InfoType_T] += float64(len(dist[diff.InfoType_T])) / float64(sample_size)
	}
	return ans
}
