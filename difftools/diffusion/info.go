package diffusion

import (
	"difftools/network"
	"fmt"
	"math/rand"
)

type InfoType int

const (
	InfoType_F InfoType = 0
	InfoType_T InfoType = 1
)
const InfoTypes_n = 2

var AllInfoTypes = []InfoType{
	InfoType_F,
	InfoType_T,
}

type SeedInfo int

const (
	SeedInfoF SeedInfo = SeedInfo(InfoType_F + 1)
	SeedInfoT SeedInfo = SeedInfo(InfoType_T + 1)
)

func InfoToSeed(info InfoType) SeedInfo {
	return SeedInfo(info + 1)
}

func BoolToSeed(b bool) SeedInfo {
	if b {
		return SeedInfoT
	} else {
		return SeedInfoF
	}
}

// var InfoType_F int = 0
// var InfoType_T int = 1
// var InfoTypes_n int = 2

type PopType int

const (
	PopLow  PopType = 0
	PopHigh PopType = 1
)

const Pops_n int = 2

type PopList map[InfoType]PopType

func MakePopList(popInfoF PopType, popInfoT PopType) PopList {
	list := make(PopList, 2)
	list[InfoType_F] = popInfoF
	list[InfoType_T] = popInfoT
	return list
}

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
