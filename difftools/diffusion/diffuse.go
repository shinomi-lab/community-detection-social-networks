package diffusion

import (
	"m/difftools/funcs"
	"math/rand"
	// "fmt"
)

func Adjmat(adj [][]int, SeedSet []int, seed int64, prob_map [2][2][2][2]float64, pop_list [2]int, interest_list [][]int, assum_list [][]int) [][]int {
	//return value is infoごとに受け取ったノードたち(index)
	//adj: 隣接行列
	//SeedSet: シードセットの初期状態
	//seed: 乱数シード値他引数
	var n int //ネットワークのユーザ数
	n = len(adj)
	recieved_list := make([][]int, InfoTypes_n) //最終的に情報を受け取ったユーザ群を保存するリスト
	for i := 0; i < InfoTypes_n; i++ {
		recieved_list[i] = make([]int, 0, n)
	}

	_ = seed

	current := make([][]int, InfoTypes_n) //次に情報発信をするユーザ群（これがなくなったら終了）
	for i := 0; i < InfoTypes_n; i++ {
		current[i] = make([]int, 0, n)
	}
	var infotypes []int = []int{InfoType_F, InfoType_T}

	//初期設定(発信源を設定している)
	for j := 0; j < n; j++ {
		for _, info := range infotypes {
			if SeedSet[j] == info+1 {
				current[info] = append(current[info], j)
				recieved_list[info] = append(recieved_list[info], j)
			}

		}
	}

	//main loop
	counter := 0
	for len(current[InfoType_F]) > 0 || len(current[InfoType_T]) > 0 {
		counter = counter + 1
		// fmt.Println("current",current)
		// fmt.Println("recieved_list",recieved_list)
		next := make([][]int, InfoTypes_n) //次のステップで情報を拡散するノードを格納する配列
		for info, set := range current {
			for _, s_node := range set {
				pop := pop_list[info]
				interest := interest_list[s_node][pop]
				assum := assum_list[s_node][info]
				p := prob_map[pop][info][interest][assum]
				if counter == 1 {
					p = p * 2
				}

				for j := 0; j < n; j++ {
					if adj[s_node][j] == 0 || funcs.Set_Has(recieved_list[InfoType_F], j) || funcs.Set_Has(recieved_list[InfoType_T], j) || funcs.Set_Has(next[InfoType_F], j) || funcs.Set_Has(next[info], j) {
						//道がないorすでに情報を受け取っているor次に偽の情報または同じ種類の情報を受け取ろうとしている
						continue
					}

					//randp := rand.Float64()
					// fmt.Println(randp)

					if p == 1 || p > rand.Float64() {
						// 確率的な情報拡散（確率pで情報発信）
						// fmt.Println(s_node,"to",j,"\t",info, "complete",p,randp)
						next[info] = append(next[info], j)
						if info == InfoType_F && funcs.Set_Has(next[InfoType_T], j) {
							remove(next[InfoType_T], j)
						}
					} else {
						// fmt.Println(s_node,"to",j,"\t",info, "defete",p,randp)

					}
				}
			}
		}
		current = make([][]int, InfoTypes_n) //currentリセット
		_ = copy(current, next)              //nextにcuurentに代入
		for _, info := range infotypes {
			recieved_list[info] = funcs.Set_Sum(recieved_list[info], next[info])
		}
	}
	// fmt.Println("recieved_list:",recieved_list)
	return recieved_list
}

// リストから要素を消す関数
func remove(ints []int, search int) []int {
	result := []int{}
	for _, v := range ints {
		if v != search {
			result = append(result, v)
		}
	}
	return result
}
