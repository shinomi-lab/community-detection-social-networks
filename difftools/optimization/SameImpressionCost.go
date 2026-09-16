package optimization

import (
	"bufio"
	diff "difftools/diffusion"
	"difftools/network"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func SameImpressionCost(
	sampleSize int,
	// adj [][]int,
	net network.Network,
	nonUseList []int,
	probMap diff.UserProbTable,
	pop [2]int,
	interestList [][]int,
	assumList [][]int,
	OnlyInfler bool,
	exit_f bool,
	r *rand.Rand,
) {
	//sample_size2はグリーディで求めた解をより詳しくやる
	var n int = net.N
	var result float64

	S := make([]int, net.N)
	S_test := make([]int, net.N)
	if exit_f {
		max_user := 0 //最もフォロワ数が多いユーザ名
		max_user_num := 0
		user_num_counter := 0
		for i := 0; i < net.N; i++ {
			user_num_counter = 0
			for j := 0; j < net.N; j++ {
				if net.Adj[i][j] == 1 {
					user_num_counter++
				}
			}
			if max_user_num < user_num_counter {
				max_user = i
				max_user_num = user_num_counter
			}
		}
		S[max_user] = 1
	}

	var info_num int

	info_num = 2

	file, err := os.Create("SameImporessionCost.csv")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	// 関数の終了時にファイルを閉じる
	defer file.Close()
	writer := bufio.NewWriter(file)

	for j := 0; j < n; j++ {

		_ = copy(S_test, S) //初期化

		if S_test[j] == 2 { //すでに発信源のユーザだったら
			fmt.Println("error in SameImpressionCost.go")
			os.Exit(0)
		}
		if S_test[j] == 1 { //虚偽情報の発信源だったら
			continue
		}
		if contains(nonUseList, j) {
			continue
		}
		if OnlyInfler {
			// if FolowerSize(adj, j) == 0 {
			if net.Followers[j] == 0 {
				continue
			}
		}

		S_test[j] = info_num

		dist := RunInflProp(sampleSize, net, S_test, probMap, pop, interestList, assumList, r)

		result = dist[diff.InfoType_T]

		// ファイルを作成または開く

		// ファイルに書き込む
		_, err = writer.WriteString(
			fmt.Sprintf("%d", j) + "," +
				// fmt.Sprintf("%d", FolowerSize(adj, j)) + "," +
				fmt.Sprintf("%d", net.Followers[j]) + "," +
				fmt.Sprintf("%f", result) + "\n")
		// _, err = file.WriteString("aaaaa\n")
		// fmt.Println("aa")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		// fmt.Println("File written successfully")
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing buffer:", err)
		return
	}
}

func SameImpressionCostInfl(
	sampleSize int,
	// adj [][]int,
	net network.Network,
	nonUseList []int,
	probMap diff.UserProbTable,
	pop [2]int,
	interestList [][]int,
	assumList [][]int,
	OnlyInfler bool,
	exit_f bool,
	r *rand.Rand,
) {
	//sample_size2はグリーディで求めた解をより詳しくやる
	var n int = net.N
	var result float64

	S := make([]int, net.N)
	S_test := make([]int, net.N)
	if exit_f {
		max_user := 0 //最もフォロワ数が多いユーザ名
		max_user_num := 0
		user_num_counter := 0
		for i := 0; i < net.N; i++ {
			user_num_counter = 0
			for j := 0; j < net.N; j++ {
				if net.Adj[i][j] == 1 {
					user_num_counter++
				}
			}
			if max_user_num < user_num_counter {
				max_user = i
				max_user_num = user_num_counter
			}
		}
		S[max_user] = 1
	}

	var info_num int

	info_num = 2

	file, err := os.Create("SameImporessionCost.csv")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	// 関数の終了時にファイルを閉じる
	defer file.Close()
	writer := bufio.NewWriter(file)

	for j := 0; j < n; j++ {
		fmt.Println("first loop")
		j = j + r.Intn(60)
		if j >= n {
			break
		}
		_ = copy(S_test, S) //初期化

		if S_test[j] == 2 { //すでに発信源のユーザだったら
			continue
		}
		if S_test[j] == 1 { //虚偽情報の発信源だったら
			continue
		}
		if contains(nonUseList, j) {
			continue
		}
		if OnlyInfler {
			// if FolowerSize(adj, j) == 0 {
			if net.Followers[j] == 0 {
				continue
			}
		}

		S_test[j] = info_num
		for i := 0; i < n; i++ {
			currentTime := time.Now()

			// 日付と時間（時分）のフォーマット
			formattedTime := currentTime.Format("01-02 15:04") // 月-日 時:分

			// フォーマットした時間を表示

			fmt.Println("second loop", formattedTime)
			i = i + r.Intn(60)
			if i == j {
				continue
			}
			if i >= n {
				break
			}
			_ = copy(S_test, S) //初期化

			if S_test[i] == 2 { //すでに発信源のユーザだったら
				continue
			}
			if S_test[i] == 1 { //虚偽情報の発信源だったら
				continue
			}
			if contains(nonUseList, i) {
				continue
			}
			if OnlyInfler {
				// if FolowerSize(adj, i) == 0 {
				if net.Followers[i] == 0 {
					continue
				}
			}

			S_test[i] = info_num
			for k := 0; k < n; k++ {
				k = k + r.Intn(60)
				if k == j || k == i {
					continue
				}
				if k >= n {
					break
				}
				_ = copy(S_test, S) //初期化

				if S_test[k] == 2 { //すでに発信源のユーザだったら
					continue
				}
				if S_test[k] == 1 { //虚偽情報の発信源だったら
					continue
				}
				if contains(nonUseList, k) {
					continue
				}
				if OnlyInfler {
					// if FolowerSize(adj, k) == 0 {
					if net.Followers[k] == 0 {
						continue
					}
				}

				S_test[k] = info_num

				dist := RunInflProp(sampleSize, net, S_test, probMap, pop, interestList, assumList, r)

				result = dist[diff.InfoType_T]

				// ファイルを作成または開く

				// ファイルに書き込む
				_, err = writer.WriteString(
					fmt.Sprintf("%d", j) + "," +
						fmt.Sprintf("%d", Cal_cost_infl_int(net, j, probMap, pop, interestList, assumList)+
							Cal_cost_infl_int(net, i, probMap, pop, interestList, assumList)+
							Cal_cost_infl_int(net, k, probMap, pop, interestList, assumList)) + "," +
						fmt.Sprintf("%f", result) + "\n")
				// _, err = file.WriteString("aaaaa\n")
				// fmt.Println("aa")
				if err != nil {
					fmt.Println("Error writing to file:", err)
					return
				}
			}
		}

		// fmt.Println("File written successfully")
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing buffer:", err)
		return
	}
}

func contains(slice []int, target int) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}
