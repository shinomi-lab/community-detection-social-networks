package diffusion

import (
	"difftools/network"
	"fmt"
	"math/rand"
)

const (
	InfoType_F  = 0
	InfoType_T  = 1
	InfoTypes_n = 2
)

const (
	SeedInfoF = InfoType_F + 1
	SeedInfoT = InfoType_T + 1
)

func InfoToSeed(info int) int {
	return info + 1
}

// var InfoType_F int = 0
// var InfoType_T int = 1
// var InfoTypes_n int = 2
const (
	Pop_low  = 0
	Pop_high = 1
)

// var Pops_n int = 2
// func make_Info(pop int) {
// 	var a [InfoTypes_n][pops_n]int
// }

// func make_InfoTypes() [2]int{
// 	a [InfoTypes_n]int := [InfoType_F,InfoType_T]
// 	return a
// }

func MakeSeedSetF(
	// adj [][]int,
	// n int,
	net network.Network,
	k int,
	r *rand.Rand) ([]int, []int) {
	//n:ノード数,k:SeedSetFの個数
	Fs := make([]int, net.N)              //各ノードの初期状態を格納する配列
	CandidateSet := make([]int, 0, net.N) //選ばれる可能性があるノードたち(出次数が1以上)

	for i := 0; i < net.N; i++ {
		for j := 0; j < net.N; j++ {
			if net.Adj[i][j] > 0 {
				CandidateSet = append(CandidateSet, i)
				break
			}
		}
	}
	fmt.Println("選ばれうる（F)")
	fmt.Println(CandidateSet)
	if len(CandidateSet) < k {
		k = len(CandidateSet)
		fmt.Println("十分な数の候補がありません")
	}
	for i := 0; i < k; {
		r := CandidateSet[r.Intn(len(CandidateSet))]
		if Fs[r] == 0 {
			Fs[r] = 1
			i++
		}
	}

	return Fs, CandidateSet
}
