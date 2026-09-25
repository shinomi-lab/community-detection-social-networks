package sample

import (
	diff "difftools/diffusion"
	exp "difftools/experiment"
	"difftools/network"
	opt "difftools/optimization"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func partial(N int, adjFilePath string) network.Network {
	net := network.ReadAdjJson(adjFilePath)

	adjPart := make([][]int, N)
	for i := range N {
		adjPart[i] = make([]int, N)
		for j := range N {
			adjPart[i][j] = net.Adj[i][j]
		}
	}
	FollowerNums := network.GetFollowerNums(adjPart)

	return network.Network{Adj: adjPart, N: N, Followers: FollowerNums}
}

func sample1() {
	var n int = 100
	var seed int64 = 1
	var K_F int = 5
	var K_T int = 10
	var sample_size int = 1000
	var pop_list = diff.MakePopList(
		// pop_list[0] = diff.PopHigh
		diff.PopHigh,
		// pop_list[1] = diff.PopHigh
		diff.PopHigh,
	)
	// [2]int

	fmt.Println(K_T, K_F, diff.InfoType_F, sample_size, pop_list)
	adjFilePath := "adj_jsonTwitterInteractionUCongress.txt"
	result_Path := "Twitter_Data/"
	// bytes, err := os.ReadFile(adjFilePath)
	// if err != nil {
	// 	panic(err)
	// }

	// // fmt.Println(string(bytes))

	// var dataJson string = string(bytes)

	// arr := make(map[int]map[int]int)
	// // var arr []string
	// _ = json.Unmarshal([]byte(dataJson), &arr)
	// // fmt.Println(arr)

	// // fmt.Println(arr[0][1])

	// var adj [][]int = make([][]int, n)

	// //make adj
	// for i := 0; i < n; i++ {
	// 	adj[i] = make([]int, n)
	// 	for j := 0; j < n; j++ {
	// 		adj[i][j] = arr[j][i]
	// 	}
	// }

	net := partial(n, adjFilePath)

	fmt.Println(net.Adj)
	fmt.Println("--------------------")

	r := rand.New(rand.NewSource(seed))

	SeedSet_F, _ := diff.MakeSeedSetF(net, 1, r)
	interest_list := diff.MakeInterestList(n, r)
	assum_list := diff.MakeAssumList(n, r)
	prob_map := diff.GetUserProbTable()

	fmt.Println("Seedsetf")
	fmt.Println(SeedSet_F)

	fmt.Println("prob_map")
	fmt.Println(prob_map)

	SeedSet_F_strong2 := make([]diff.SeedInfo, net.N)
	SeedSet_F_strong2[0] = 1

	//pythonのやつと実行結果が違う理由を確かめるために使った部分
	// python_ans,python_node_list := opt.PythonSuppression(adj, SeedSet_F_strong2,  prob_map, pop_list, interest_list, assum_list, true)
	//
	// fmt.Println(python_node_list)
	//
	// for i:=0;i<len(python_ans);i++{
	// 	fmt.Println(python_ans[i])
	// }

	// os.Exit(0)
	// node_num := 5// it mean node num

	//人数を流動的にして拡散を調べている
	//	総フォロワー数を固定できていない
	//拡散可能な人数を調べている
	infler_num := 0
	OnlyInfler := true
	for j := 0; j < net.N; j++ {
		for k := 0; k < net.N; k++ {
			if net.Adj[j][k] != 0 {
				infler_num += 1
				break
			}
		}
	}
	greedy_ans, _, greedy_ans_v := opt.Greedy(100, net, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list, infler_num, true, 1000, r)

	fmt.Println(greedy_ans)
	fmt.Println(greedy_ans_v)

	os.Exit(0)

	kurikaesi := 100 //it mean loop num nearly sample_size

	fmt.Println("infler_num:", infler_num)
	ans3 := make([][]float64, 0) //values?
	ans5 := make([][]int, 0)     //nodes

	for i := 1; i < infler_num; i++ {
		// kurikaesi = i * net.N
		_, ans2, ans4 := opt.RandomSuppression(net, i, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list, kurikaesi, OnlyInfler, r)

		ans3 = append(ans3, ans2...)
		ans5 = append(ans5, ans4...)
		//

		file1, err := os.Create(result_Path + strconv.Itoa(i) + "kurikasi" + strconv.Itoa(kurikaesi) + ".csv")
		if err != nil {
			panic(err)
		}
		defer file1.Close()

		// Writerを作成
		writer := csv.NewWriter(file1)

		// データを書き込み
		for _, row := range ans2 {
			stringRow := make([]string, len(row))
			for i, v := range row {
				stringRow[i] = strconv.FormatFloat(v, 'f', -1, 64)
			}
			writer.Write(stringRow)
		}

		// バッファに残っているデータを書き込み
		writer.Flush()
	}

	file1, err := os.Create(result_Path + "allkurikasi" + strconv.Itoa(kurikaesi) + ".csv")
	if err != nil {
		panic(err)
	}
	defer file1.Close()

	// Writerを作成
	writer := csv.NewWriter(file1)

	// データを書き込み
	for _, row := range ans3 {
		stringRow := make([]string, len(row))
		for i, v := range row {
			stringRow[i] = strconv.FormatFloat(v, 'f', -1, 64)
		}
		writer.Write(stringRow)
	}

	// バッファに残っているデータを書き込み
	writer.Flush()

	file1, err = os.Create(result_Path + "Nodeallkurikasi" + strconv.Itoa(kurikaesi) + ".csv")
	if err != nil {
		panic(err)
	}
	defer file1.Close()

	// Writerを作成
	writer = csv.NewWriter(file1)
	// データを書き込み
	for _, row := range ans5 {
		stringRow := make([]string, len(row))
		for i, v := range row {
			stringRow[i] = strconv.Itoa(v)
		}
		writer.Write(stringRow)
	}

	// バッファに残っているデータを書き込み
	writer.Flush()
	//
	// opt.CallKumiawase(adj, 11,13, SeedSet_F_strong2)
	os.Exit(0)

	//なぜか同じ拡散力のやつで同じ拡散を調べている
	for j := 0; j < 50; j++ {

		selected_list := opt.CallKumiawase_Impression(
			net, j, j+4, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list)

		opt.Selected_Suppression(net, selected_list, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list, r)
	}

	node_num := 100
	for i := 0; i < 20; i++ {
		adjFilePath = "adj_json" + strconv.Itoa(node_num) + "node" + strconv.Itoa(i) + "seed.txt"
		net = partial(node_num, adjFilePath)
		// bytes, err = os.ReadFile(adjFilePath)
		// if err != nil {
		// 	panic(err)
		// }

		// // fmt.Println(string(bytes))

		// dataJson = string(bytes)

		// arr = make(map[int]map[int]int)
		// // var arr []string
		// _ = json.Unmarshal([]byte(dataJson), &arr)
		// // fmt.Println(arr)

		// // fmt.Println(arr[0][1])

		// adj = make([][]int, n)

		// for i := 0; i < n; i++ {
		// 	adj[i] = make([]int, n)
		// 	for j := 0; j < n; j++ {
		// 		adj[i][j] = arr[j][i]
		// 	}
		// }

		for j := 0; j < node_num-1; j++ {

			selected_list := opt.CallKumiawase_Impression(
				net, j, j+4, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list)
			fmt.Print(j, ":")

			opt.Selected_Suppression(
				net, selected_list, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list, r)
		}
		fmt.Println("-----------------------------------------")
	}

	os.Exit(0)

	// fmt.Println((seq))
	//
	// fmt.Println(prob_map)
	seedsetf_pram := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if SeedSet_F[i] == 1 {
			seedsetf_pram = append(seedsetf_pram, i)
		}
	}
	pram_data := Parameter{
		GraphPath:        adjFilePath,
		Node_n:           n,
		Random_seed:      seed,
		K_F:              K_F,
		K_T:              K_T,
		Mont_sample_size: sample_size,
		Pop_list:         pop_list,
		SeedSet_F:        seedsetf_pram,
		InterestList:     interest_list,
		AssumList:        assum_list,
		ProbMap:          prob_map,
	}

	fmt.Println(pram_data)

	// jsonData, err := json.Marshal(pram_data)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// pram_file, err := os.Create("result/param_json.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//  _, err1 := pram_file.WriteString(jsonData)
	//  if err1 != nil {
	// 	log.Fatal(err)
	// }
	const layout2 = "2006-01-02 15:04:05"
	str := strings.ReplaceAll(time.Now().Format(layout2), ":", "-")
	folder_path := "result/" + strconv.Itoa(n) + "node" + str
	err = os.Mkdir(folder_path, os.ModePerm)
	if err != nil {
		fmt.Println("error create" + folder_path)
		log.Fatal(err)
	}

	file, _ := json.MarshalIndent(pram_data, "", " ")
	_ = os.WriteFile(folder_path+"/param_json.txt", file, 0644)

	//test part

	SeedSet_F_strong := make([]diff.SeedInfo, net.N)
	SeedSet_F_strong[0] = 1 //here
	// SeedSet_Greedy[1] = 1
	//偽情報の発信源を色々と
	// greedy_ans1, greedy_value1, greedy_value21 := opt.Greedy(1, 100, adj, SeedSet_F_strong, prob_map, pop_list, interest_list, assum_list, 3, true, 1000)
	// fmt.Println(greedy_ans1, greedy_value1, greedy_value21)
	// os.Exit(0)

	sample_size = 1000000
	// var S []int
	// var hist [][]float64
	S, hist := exp.SimSubmod(
		sample_size, net, pop_list, interest_list, assum_list, SeedSet_F_strong, K_T, prob_map, folder_path, r)

	fmt.Println("End Check_submod")
	os.Exit(0)

	// new_folder_path := folder_path + "/FocusLoop"
	// err = os.Mkdir(new_folder_path, os.ModePerm)
	// if err != nil {
	// 	fmt.Println("error create FocusLoop")
	// 	log.Fatal(err)
	// }
	// //[0 1 2 6 8 15 18 20 21 22 28 32 33 37 48 49 61 67 76 93]
	// loop_n := 1000
	// sample_size = 1000
	// list1 := []int{0, 6, 8}
	// list2 := []int{0, 20}
	// opt.FocusLoop(loop_n, list1, list2, SeedSet_F, 1, sample_size, adj, prob_map, pop_list, interest_list, assum_list, new_folder_path)

	// time.Sleep(time.Second * 2)
	// list1 = []int{68, 48, 2, 8}
	// list2 = []int{8, 48, 37, 93}
	// opt.FocusLoop(loop_n, list1, list2, SeedSet_F, 1, sample_size, adj, prob_map, pop_list, interest_list, assum_list, new_folder_path)

	// time.Sleep(time.Second * 2)
	// list1 = []int{20, 15, 6}
	// list2 = []int{48, 0, 18}
	// opt.FocusLoop(loop_n, list1, list2, SeedSet_F, 1, sample_size, adj, prob_map, pop_list, interest_list, assum_list, new_folder_path)

	// fmt.Println("End FocusLoop")

	// filename := folder_path + "/GreedyAndStrict2.csv"
	// f, err2 := os.Create(filename)
	// if err2 != nil {
	// 	fmt.Println("error create" + filename)
	// 	log.Fatal(err)
	// }
	//
	// w := csv.NewWriter(f)
	//
	// colmns := []string{"greedy_ans", "strict_ans", "greedy_value", "greedy_value2", "strict_value", "strict_value2", "kinjiritu", "random_seed"}
	// w.Write(colmns)
	//
	// //start loop
	// sample_size = 1000
	// sample_size2 := 1000
	// var random_seed int64
	// random_seed = 0
	//
	// //選ばれうる0 1 2 6 8 15 18 20 37 48
	//
	// seedsetfs := []int{0, 1, 2, 6}
	// // var seedsetfs []int = make([]int, len(diff.Set))
	// // _ = copy(seedsetfs, diff.Set)
	// // fmt.Println(seedsetfs)
	// // os.Exit(0)
	// fmt.Println((len(adj)))
	// for i := 0; i < 2; i++ {
	// 	SeedSet_Greedy := make([]int, len(adj))
	// 	SeedSet_Greedy[seedsetfs[i]] = 1 //
	// 	//偽情報の発信源を色々と
	// 	for random_seed = 0; random_seed < 10; random_seed++ {
	// 		greedy_ans, greedy_value, greedy_value2 := opt.Greedy(random_seed, sample_size, adj, SeedSet_Greedy, prob_map, pop_list, interest_list, assum_list, 3, true, sample_size2)
	//
	// 		fmt.Println("greedy_ans")
	// 		fmt.Println(greedy_ans, greedy_value, greedy_value2)
	//
	// 		strict_ans, strict_value, strict_value2 := opt.Strict(random_seed, sample_size, adj, SeedSet_Greedy, prob_map, pop_list, interest_list, assum_list, 3, true, sample_size2)
	//
	// 		fmt.Println("strict_ans")
	// 		fmt.Println(strict_ans, strict_value, strict_value2)
	//
	// 		fmt.Println("近似率")
	// 		fmt.Println(greedy_value2 / strict_value2)
	//
	// 		Sets_string := make([][]string, 2)
	// 		Sets_string[0] = opt.Int_to_String(greedy_ans)
	// 		Sets_string[1] = opt.Int_to_String(strict_ans)
	//
	// 		part0 := []string{strings.Join(Sets_string[0], "-"), strings.Join(Sets_string[1], "-")} //here
	//
	// 		a := []float64{greedy_value, greedy_value2, strict_value, strict_value2, greedy_value2 / strict_value2, float64(random_seed)}
	//
	// 		part1 := opt.Float_to_String(a)
	//
	// 		retu := append(part0, part1...)
	//
	// 		w.Write(retu)
	// 	}
	// }
	//
	// //loop end
	//
	// w.Flush()
	//
	// if err := w.Error(); err != nil {
	// 	log.Fatal(err)
	// }
	// //SeedSet_Greedy := make([]int,len(adj))
	// //SeedSet_Greedy[15] = 2

	fmt.Println(S, hist)
}
