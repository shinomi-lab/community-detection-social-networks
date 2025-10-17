package diffusion

import (
	"fmt"
	"math/rand"
)

var InfoType_F int = 0
var InfoType_T int = 1
var InfoTypes_n int = 2
var Pop_low int = 0
var Pop_high int = 1
var Pops_n int = 2
var Set []int

// func make_Info(pop int) {
// 	var a [InfoTypes_n][pops_n]int
// }

// func make_InfoTypes() [2]int{
// 	a [InfoTypes_n]int := [InfoType_F,InfoType_T]
// 	return a
// }

func Make_seedSet_F(n int, k int, seed int64, adj [][]int) []int {
	//n:ノード数,k:SeedSetFの個数
	// rand.Seed(seed)
	_ = seed
	Fs := make([]int, n)//各ノードの初期状態を格納する配列
	Set = make([]int,0,n)//選ばれる可能性があるノードたち(出次数が1以上)

	for i:=0;i<n;i++{
		for j:=0;j<n;j++{
			if adj[i][j] > 0{
				Set = append(Set,i)
				break
			}
		}
	}
	fmt.Println("選ばれうる（F)")
	fmt.Println(Set)
	if len(Set) < k{
		k = len(Set)
		fmt.Println("十分な数の候補がありません")
	}
	for i := 0; i < k; {
		r := Set[rand.Intn(len(Set))]
		if Fs[r] == 0 {
			Fs[r] = 1
			i++
		}
	}

	return Fs
}
