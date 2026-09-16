package diffusion

import (
	"difftools/funcs"
	"difftools/network"
	"math/rand"
	// "fmt"
)

// 乱数生成器を変数`r`としてパラメータ化
//
// return: infoごとに受け取ったノードたち(index)
func Diffuse(
	net network.Network,
	// adj [][]int, // 隣接行列
	seedSet []int, // シードセットの初期状態
	userProbTable UserProbTable,
	popList [2]int,
	interestList [][]int,
	assumList [][]int,
	r *rand.Rand,
) [][]int {
	// var n int //ネットワークのユーザ数
	// n = len(adj)
	receivedList := make([][]int, InfoTypes_n) //最終的に情報を受け取ったユーザ群を保存するリスト
	for i := 0; i < InfoTypes_n; i++ {
		receivedList[i] = make([]int, 0, net.N)
	}

	current := make([][]int, InfoTypes_n) //次に情報発信をするユーザ群（これがなくなったら終了）
	for i := 0; i < InfoTypes_n; i++ {
		current[i] = make([]int, 0, net.N)
	}
	var infotypes []int = []int{InfoType_F, InfoType_T}

	//初期設定(発信源を設定している)
	for j := 0; j < net.N; j++ {
		for _, info := range infotypes {
			// if seedSet[j] == info+1 {
			if seedSet[j] == InfoToSeed(info) {
				current[info] = append(current[info], j)
				receivedList[info] = append(receivedList[info], j)
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
				pop := popList[info]
				interest := interestList[s_node][pop]
				assum := assumList[s_node][info]
				p := userProbTable[pop][info][interest][assum]
				if counter == 1 {
					p = p * 2
				}

				for j := 0; j < net.N; j++ {
					if net.Adj[s_node][j] == 0 || funcs.IsInSet(receivedList[InfoType_F], j) || funcs.IsInSet(receivedList[InfoType_T], j) || funcs.IsInSet(next[InfoType_F], j) || funcs.IsInSet(next[info], j) {
						//道がないorすでに情報を受け取っているor次に偽の情報または同じ種類の情報を受け取ろうとしている
						continue
					}

					//randp := rand.Float64()
					// fmt.Println(randp)

					if p == 1 || p > r.Float64() {
						// 確率的な情報拡散（確率pで情報発信）
						// fmt.Println(s_node,"to",j,"\t",info, "complete",p,randp)
						next[info] = append(next[info], j)
						if info == InfoType_F && funcs.IsInSet(next[InfoType_T], j) {
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
			receivedList[info] = funcs.UnionSets(receivedList[info], next[info])
		}
	}
	// fmt.Println("recieved_list:",recieved_list)
	return receivedList
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
