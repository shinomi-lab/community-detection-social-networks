package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	diff "m/difftools/diffusion"
	opt "m/difftools/optimization"
	"os"
	"strconv"
	"strings"
	"time"
	"math/rand"
	// "reflect"
)

type Parameter struct {
	GraphPath        string
	Node_n           int
	Random_seed      int64
	K_F              int
	K_T              int
	Mont_sample_size int
	Prob_map         [2][2][2][2]float64
	Seq              [16]float64
	Pop_list         [2]int
	SeedSet_F        []int
	Interest_list    [][]int
	Assum_list       [][]int
}

func Make_adj_interest_assum(adjFilePath string, seed int64)([][]int,[][]int,[][]int){
	bytes, err := ioutil.ReadFile(adjFilePath)
	if err != nil {
		panic(err)
	}

	// fmt.Println(string(bytes))

	var dataJson string = string(bytes)

	arr := make(map[int]map[int]int)
	// var arr []string
	_ = json.Unmarshal([]byte(dataJson), &arr)
	// fmt.Println(arr)

	// fmt.Println(arr[0][1])

	n := len(arr)

	var interest_list [][]int = diff.Make_interest_list(n, seed)

	var assum_list [][]int = diff.Make_assum_list(n, seed)
	var adj [][]int = make([][]int, n)

	for i := 0; i < n; i++ {
		adj[i] = make([]int, n)
		for j := 0; j < n; j++ {
			adj[i][j] = arr[j][i]
		}
	}
	return adj,interest_list,assum_list
}

func use_strict(adj [][]int, interest_list [][]int,assum_list [][]int, user_weight float64)([][]int, [][]int, []int, [2][2][2][2]float64, [2]int, [][]int, [][]int){


			// var n int = 50
			// var seed int64 = 1
			// var K_F int = 5
			// var K_T int = 10
			// var sample_size int = 1000
			var pop_list [2]int
			pop_list[0] = diff.Pop_high
			pop_list[1] = diff.Pop_high

			// fmt.Println(K_T, K_F, diff.InfoType_F, sample_size, pop_list)
			// adjFilePath := "adj_jsonTwitterInteractionUCongress.txt"
			// adjFilePath := "Graphs/adj_json1000node.txt"
			// result_Path := "Twitter_Data/"


			// var SeedSet_F []int = diff.Make_seedSet_F(n, 1, seed, adj)

			// var interest_list [][]int = diff.Make_interest_list(n, seed)
			//
			// var assum_list [][]int = diff.Make_assum_list(n, seed)

			var seq [16]float64 = diff.Make_probability()

			var prob_map [2][2][2][2]float64 = diff.Map_probagbility(seq)

			// fmt.Println("Seedsetf")
			// fmt.Println(SeedSet_F)
			//
			// fmt.Println("prob_map")
			// fmt.Println(prob_map)

			SeedSet_F_strong2 := make([]int, len(adj))
			max_user := 0 //最もフォロワ数が多いユーザ名
			max_user_num := 0
			user_num_counter := 0
			for i:=0; i<len(adj); i++{
				user_num_counter = 0
				for j:=0; j<len(adj); j++{
					if(adj[i][j] == 1){
						user_num_counter ++
					}
				}
				if(max_user_num < user_num_counter){
					max_user = i
					max_user_num = user_num_counter
				}
			}
			SeedSet_F_strong2[max_user] = 1



			//人数を流動的にして拡散を調べている
				//	総フォロワー数を固定できていない
			//拡散可能な人数を調べている
			infler_num := 0
			// OnlyInfler := true
			for j:=0;j<len(adj);j++{
				for k:=0;k<len(adj);k++{
					if(adj[j][k] != 0){
						infler_num += 1
						break
					}
				}
			}

			fmt.Println("start_strict")
			// slice_test := [][]int{{1, 15, 18}, {1, 15, 18},{1, 15, 18},{1, 15, 18},{1, 15, 18}}
			// opt.Selected_Suppression_Maximum(adj,slice_test,SeedSet_F_strong2,  prob_map , pop_list, interest_list, assum_list)
			// os.Exit(0)
			s := time.Now()
			under := 0.0
			upper := 2.0
			selected_list := opt.CallKumiawase2(adj, under, upper, SeedSet_F_strong2, true, max_user, user_weight)

			// fmt.Println("end all kumiawase",selected_list)
			cost_sum := 0.0
			for k:=0;k<len(selected_list);k++{
				selecte := selected_list[k]
				cost_sum = 0
				for l:=0;l<len(selecte);l++{
					cost_sum += opt.Cal_cost(user_weight,1.0-user_weight,adj, selecte[l], max_user)
				}
				if(cost_sum < under || cost_sum >upper){
					fmt.Println("えらーcost_sum:",cost_sum)
				}

			}
			// os.Exit(0)
			strict_ans,strict_ans_v,strict_ans_fv := opt.Selected_Suppression_Maximum(adj, selected_list, SeedSet_F_strong2,  prob_map , pop_list, interest_list, assum_list)

			fmt.Println("strict_time",time.Since(s))
			cost_sum = 0
			for j:=0;j<len(strict_ans);j++{
				cost_sum += opt.Cal_cost(0.5,0.5,adj, strict_ans[j], max_user)
			}

			fmt.Println(strict_ans)
			fmt.Println(strict_ans_v)
			fmt.Println(strict_ans_fv)
			fmt.Println("cost_sum:",cost_sum)


			return adj, selected_list, SeedSet_F_strong2,  prob_map , pop_list, interest_list, assum_list
}

func use_greedy(adj [][]int, interest_list [][]int,assum_list [][]int, user_weight float64, capacity float64)([][]int, []int, [2][2][2][2]float64, [2]int, [][]int, [][]int){

		// var n int = 50
		// var seesd int64 = 1
		// var K_F int = 5
		// var K_T int = 10
		// var sample_size int = 1000
		var pop_list [2]int
		pop_list[0] = diff.Pop_high
		pop_list[1] = diff.Pop_high




		// fmt.Println(string(bytes))






		fmt.Println("--------------------")

		// var SeedSet_F []int = diff.Make_seedSet_F(n, 1, seed, adj)

		// var interest_list [][]int = diff.Make_interest_list(n, seed)
		//
		// var assum_list [][]int = diff.Make_assum_list(n, seed)

		var seq [16]float64 = diff.Make_probability()

		var prob_map [2][2][2][2]float64 = diff.Map_probagbility(seq)

		// fmt.Println("Seedsetf")
		// fmt.Println(SeedSet_F)
		//
		// fmt.Println("prob_map")
		// fmt.Println(prob_map)

		SeedSet_F_strong2 := make([]int, len(adj))
		max_user := 0 //最もフォロワ数が多いユーザ名
		max_user_num := 0
		user_num_counter := 0
		for i:=0; i<len(adj); i++{
			user_num_counter = 0
			for j:=0; j<len(adj); j++{
				if(adj[i][j] == 1){
					user_num_counter ++
				}
			}
			if(max_user_num < user_num_counter){
				max_user = i
				max_user_num = user_num_counter
			}
		}
		SeedSet_F_strong2[max_user] = 1



		//人数を流動的にして拡散を調べている
			//	総フォロワー数を固定できていない
		//拡散可能な人数を調べている
		infler_num := 0
		// OnlyInfler := true
		for j:=0;j<len(adj);j++{
			for k:=0;k<len(adj);k++{
				if(adj[j][k] != 0){
					infler_num += 1
					break
				}
			}
		}

		fmt.Println("start_greedy")
		// greedy_ans1, _, _ := opt.Greedy(0,100,adj,SeedSet_F_strong2, prob_map,pop_list,interest_list,assum_list,5,true,1000)

		cost_sum := 0.0
		// for j:=0;j<len(greedy_ans1);j++{
		// 	cost_sum += opt.Cal_cost_kaiki(user_weight,1-user_weight,adj, greedy_ans1[j], max_user)
		// }
		// fmt.Println("cost_sum",cost_sum)
		// cost_sum = 0
		// os.Exit(0)

		s := time.Now()
		//虚偽情報アリの影響最大化問題の解を求める
		greedy_ans, _ := opt.Greedy_exp(0,100,adj,SeedSet_F_strong2, prob_map,pop_list,interest_list,assum_list,infler_num,true,capacity,max_user,true, user_weight,true)
		fmt.Println("greedy_time:",time.Since(s))

		for j:=0;j<len(greedy_ans);j++{
			cost_sum += opt.Cal_cost_kaiki(user_weight,1-user_weight,adj, greedy_ans[j], max_user)
		}
		greedy_ans2 := make([][]int,0)
		greedy_ans2 = append(greedy_ans2,greedy_ans)
		SeedSet_F_strong2 = make([]int, len(adj))
		SeedSet_F_strong2[max_user] = 1
		_,test_greedy_ans_v,test_greedy_ans_fv := opt.Selected_Suppression_Maximum(adj,greedy_ans2, SeedSet_F_strong2,  prob_map , pop_list, interest_list, assum_list)
		SeedSet_F_strong2 = make([]int, len(adj))//念のため初期化
		SeedSet_F_strong2[max_user] = 1

		fmt.Println("虚偽情報アリの解",greedy_ans,test_greedy_ans_v,test_greedy_ans_fv)
		nonF_SeedSet := make([]int, len(adj))

		_,test_greedy_ans_v,test_greedy_ans_fv = opt.Selected_Suppression_Maximum(adj,greedy_ans2, nonF_SeedSet,  prob_map , pop_list, interest_list, assum_list)

		fmt.Println("虚偽情報アリの解を無しに使ってみたら...",test_greedy_ans_v,test_greedy_ans_fv)
		// fmt.Println(greedy_ans_v)
		fmt.Println("cost_sum:",cost_sum)
		nonF_SeedSet = make([]int, len(adj))//念のため初期化
		greedy_ans, _ = opt.Greedy_exp(0,100,adj,nonF_SeedSet, prob_map,pop_list,interest_list,assum_list,infler_num,true,capacity,max_user,true, user_weight,true)
		fmt.Println("greedy_time:",time.Since(s))

		cost_sum = 0
		for j:=0;j<len(greedy_ans);j++{
			cost_sum += opt.Cal_cost_kaiki(user_weight,1-user_weight,adj, greedy_ans[j], max_user)
		}
		greedy_ans2 = make([][]int,0)
		greedy_ans2 = append(greedy_ans2,greedy_ans)

		nonF_SeedSet = make([]int, len(adj))//念のため初期化
		_,test_greedy_ans_v,test_greedy_ans_fv = opt.Selected_Suppression_Maximum(adj,greedy_ans2, nonF_SeedSet,  prob_map , pop_list, interest_list, assum_list)

		fmt.Println("虚偽情報なしの解",greedy_ans,test_greedy_ans_v,test_greedy_ans_fv)
		// fmt.Println(greedy_ans_v)
		// fmt.Println(test_greedy_ans_v)
		fmt.Println("cost_sum:",cost_sum)


		_,test_greedy_ans_v,test_greedy_ans_fv = opt.Selected_Suppression_Maximum(adj,greedy_ans2, SeedSet_F_strong2,  prob_map , pop_list, interest_list, assum_list)

		fmt.Println("虚偽情報なしの解をアリに使ってみたら...",test_greedy_ans_v,test_greedy_ans_fv)



		return adj, SeedSet_F_strong2,  prob_map , pop_list, interest_list, assum_list
}

func use_DP(adj [][]int, interest_list [][]int,assum_list [][]int, user_weight float64, capacity float64,use_kaiki bool, use_user bool, use_infl bool, nick int, S_f_type int, only_last bool, DP_ans_d []int, DP_ans_s []int)([][]int, []int, [2][2][2][2]float64, [2]int, [][]int, [][]int){

		// var n int = 50
		// var seesd int64 = 1
		// var K_F int = 5
		// var K_T int = 10
		// var sample_size int = 1000
		var pop_list [2]int
		pop_list[0] = diff.Pop_high
		pop_list[1] = diff.Pop_high




		// fmt.Println(string(bytes))






		fmt.Println("--------------------")

		// var SeedSet_F []int = diff.Make_seedSet_F(n, 1, seed, adj)

		// var interest_list [][]int = diff.Make_interest_list(n, seed)
		//
		// var assum_list [][]int = diff.Make_assum_list(n, seed)

		var seq [16]float64 = diff.Make_probability()

		var prob_map [2][2][2][2]float64 = diff.Map_probagbility(seq)



		SeedSet_F_strong2 := make([]int, len(adj))//ユーザの初期状態
		non_use_list := make([]int,1)//虚偽情報の発信源を選択されないようにする(単一情報で)
		max_user := 0 //最もフォロワ数が多いユーザ名

		//虚偽情報の発信源を定義
		if(S_f_type == 1){

			max_user_num := 0
			user_num_counter := 0
			for i:=0; i<len(adj); i++{
				user_num_counter = 0
				for j:=0; j<len(adj); j++{
					if(adj[i][j] == 1){
						user_num_counter ++
					}
				}
				if(max_user_num < user_num_counter){
					max_user = i
					max_user_num = user_num_counter
				}
			}
			SeedSet_F_strong2[max_user] = 1//虚偽情報の発信源を定義
			non_use_list[0] = max_user
		}else if (S_f_type == 2){
			num2 := 0
			num3 := 0
			for focus_user,slice := range adj{
				num := 0
				for _,edge := range slice{
					num += edge
					if edge >1 {
						fmt.Println("error")
						os.Exit(0)
					}
				}
				if num >20 && num < 30 {
					if num2 % 20 == 0{//個数調整
						SeedSet_F_strong2[focus_user] = 1//虚偽情報の発信源を定義
						if num3 == 0{
							non_use_list[0] = focus_user
						}else{
							non_use_list = append(non_use_list,focus_user)
						}
						num3 ++
					}
					num2 ++
				}
			}
		}




		//人数を流動的にして拡散を調べている
			//	総フォロワー数を固定できていない
		//拡散可能な人数を調べている
		infler_num := 0
		// OnlyInfler := true
		for j:=0;j<len(adj);j++{
			for k:=0;k<len(adj);k++{
				if(adj[j][k] != 0){
					infler_num += 1
					break
				}
			}
		}

		fmt.Println("start_DP")
		// greedy_ans1, _, _ := opt.Greedy(0,100,adj,SeedSet_F_strong2, prob_map,pop_list,interest_list,assum_list,5,true,1000)

		cost_sum := 0
		// for j:=0;j<len(greedy_ans1);j++{
		// 	cost_sum += opt.Cal_cost_kaiki(user_weight,1-user_weight,adj, greedy_ans1[j], max_user)
		// }
		// fmt.Println("cost_sum",cost_sum)
		// cost_sum = 0
		// os.Exit(0)

		s := time.Now()
		//虚偽情報アリの影響最大化問題の解を求める
		// DP_ans2 := make([][]int,0)
		if(!only_last){

			// DP_ans := DP_user_infl.Users

			for j:=0;j<len(DP_ans_d);j++{
				cost_sum += opt.Cal_cost_infl_int(adj,DP_ans_d[j],prob_map,pop_list,interest_list,assum_list)
			}
			DP_ans2 := make([][]int,0)
			DP_ans2 = append(DP_ans2,DP_ans_d)
			// SeedSet_F_strong2 = make([]int, len(adj))
			// SeedSet_F_strong2[max_user] = 1
			_,test_DP_ans_v,test_DP_ans_fv := opt.Selected_Suppression_Maximum(adj,DP_ans2, SeedSet_F_strong2,  prob_map , pop_list, interest_list, assum_list)
			// SeedSet_F_strong2 = make([]int, len(adj))//念のため初期化
			// SeedSet_F_strong2[max_user] = 1

			fmt.Println("虚偽情報アリの解",DP_ans_d,test_DP_ans_v,test_DP_ans_fv)
			nonF_SeedSet := make([]int, len(adj))

			_,test_DP_ans_v,test_DP_ans_fv = opt.Selected_Suppression_Maximum(adj,DP_ans2, nonF_SeedSet,  prob_map , pop_list, interest_list, assum_list)

			fmt.Println("虚偽情報アリの解を無しに使ってみたら...",test_DP_ans_v,test_DP_ans_fv)
		}

		// fmt.Println(greedy_ans_v)
		fmt.Println("cost_sum:",cost_sum)
		nonF_SeedSet := make([]int, len(adj))//念のため初期化




		fmt.Println("DP_time:",time.Since(s))

		// DP_ans := DP_user_infl.Users

		cost_sum = 0
		for j:=0;j<len(DP_ans_s);j++{
			cost_sum += opt.Cal_cost_infl_int(adj,DP_ans_s[j],prob_map,pop_list,interest_list,assum_list)

		}
		DP_ans2 := make([][]int,0)
		DP_ans2 = append(DP_ans2,DP_ans_s)

		nonF_SeedSet = make([]int, len(adj))//念のため初期化
		_,test_DP_ans_v,test_DP_ans_fv := opt.Selected_Suppression_Maximum(adj,DP_ans2, nonF_SeedSet,  prob_map , pop_list, interest_list, assum_list)

		fmt.Println("虚偽情報なしの解",DP_ans_s,test_DP_ans_v,test_DP_ans_fv)
		// fmt.Println(greedy_ans_v)
		// fmt.Println(test_greedy_ans_v)
		fmt.Println("cost_sum:",cost_sum)


		_,test_DP_ans_v,test_DP_ans_fv = opt.Selected_Suppression_Maximum(adj,DP_ans2, SeedSet_F_strong2,  prob_map , pop_list, interest_list, assum_list)

		fmt.Println("虚偽情報なしの解をアリに使ってみたら...",test_DP_ans_v,test_DP_ans_fv)



		return adj, SeedSet_F_strong2,  prob_map , pop_list, interest_list, assum_list
}

func cal_max_users(adj [][]int, n int){
	// max_user := 0 //最もフォロワ数が多いユーザ名
	m := n-1
	max_users := make([]int,n)
	// max_user_num := 0
	user_num_counter := 0
	max_user_nums := make([]int,n)
	for i:=0; i<len(adj); i++{
		user_num_counter = 0
		for l:=0; l<len(adj); l++{
			if(adj[i][l] == 1){
				user_num_counter ++
			}
		}

		for j:=0;j<n;j++{
      if(max_user_nums[j] < user_num_counter){

        for k:=j;k<m;k++{
          max_users[m-k+j] = max_users[m-k+j-1]
          max_user_nums[m-k+j] = max_user_nums[m-k+j-1]
          // fmt.Println(i,max_user_nums)
        }
        max_users[j] = i
        max_user_nums[j] = user_num_counter
        // fmt.Println(i,max_user_nums)

        break
      }
  	}
	}
	fmt.Println(max_users, max_user_nums)
}

func sample1() {
	var n int = 100
	var seed int64 = 1
	var K_F int = 5
	var K_T int = 10
	var sample_size int = 1000
	var pop_list [2]int
	pop_list[0] = diff.Pop_high
	pop_list[1] = diff.Pop_high

	fmt.Println(K_T, K_F, diff.InfoType_F, sample_size, pop_list)
	adjFilePath := "adj_jsonTwitterInteractionUCongress.txt"
	result_Path := "Twitter_Data/"
	bytes, err := ioutil.ReadFile(adjFilePath)
	if err != nil {
		panic(err)
	}

	// fmt.Println(string(bytes))

	var dataJson string = string(bytes)

	arr := make(map[int]map[int]int)
	// var arr []string
	_ = json.Unmarshal([]byte(dataJson), &arr)
	// fmt.Println(arr)

	// fmt.Println(arr[0][1])

	var adj [][]int = make([][]int, n)

	//make adj
	for i := 0; i < n; i++ {
		adj[i] = make([]int, n)
		for j := 0; j < n; j++ {
			adj[i][j] = arr[j][i]
		}
	}

	fmt.Println(adj)
	fmt.Println("--------------------")

	var SeedSet_F []int = diff.Make_seedSet_F(n, 1, seed, adj)

	var interest_list [][]int = diff.Make_interest_list(n, seed)

	var assum_list [][]int = diff.Make_assum_list(n, seed)

	var seq [16]float64 = diff.Make_probability()

	var prob_map [2][2][2][2]float64 = diff.Map_probagbility(seq)

	fmt.Println("Seedsetf")
	fmt.Println(SeedSet_F)

	fmt.Println("prob_map")
	fmt.Println(prob_map)

	SeedSet_F_strong2 := make([]int, len(adj))
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
	for j:=0;j<len(adj);j++{
		for k:=0;k<len(adj);k++{
			if(adj[j][k] != 0){
				infler_num += 1
				break
			}
		}
	}
	greedy_ans, _,greedy_ans_v := opt.Greedy(0,100,adj,SeedSet_F_strong2, prob_map,pop_list,interest_list,assum_list,infler_num,true,1000)

	fmt.Println(greedy_ans)
	fmt.Println(greedy_ans_v)

	os.Exit(0)


	kurikaesi := 100 //it mean loop num nearly sample_size

	fmt.Println("infler_num:",infler_num)
	ans3 := make([][]float64,0)//values?
	ans5 := make([][]int,0)//nodes

	for i:=1;i<infler_num;i++{
		// kurikaesi = i * len(adj)
		_,ans2,ans4 := opt.RandomSuppression(adj, i, SeedSet_F_strong2,  prob_map, pop_list, interest_list, assum_list, kurikaesi,OnlyInfler)

		ans3 = append(ans3,ans2...)
		ans5 = append(ans5,ans4...)
	//

	file1, err := os.Create(result_Path+strconv.Itoa(i)+"kurikasi"+strconv.Itoa(kurikaesi)+".csv")
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

	file1, err := os.Create(result_Path+"allkurikasi"+strconv.Itoa(kurikaesi)+".csv")
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


 file1, err = os.Create(result_Path+"Nodeallkurikasi"+strconv.Itoa(kurikaesi)+".csv")
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
	for j:=0;j<50;j++{

		selected_list := opt.CallKumiawase_Impression(adj , j, j+4, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list)

		opt.Selected_Suppression(adj, selected_list, SeedSet_F_strong2,  prob_map , pop_list, interest_list, assum_list)
	}

	node_num:=100
	for i:=0;i<20;i++{
		adjFilePath = "adj_json"+strconv.Itoa(node_num)+"node"+strconv.Itoa(i)+"seed.txt"
		bytes, err = ioutil.ReadFile(adjFilePath)
		if err != nil {
			panic(err)
		}

		// fmt.Println(string(bytes))

		dataJson = string(bytes)

		arr = make(map[int]map[int]int)
		// var arr []string
		_ = json.Unmarshal([]byte(dataJson), &arr)
		// fmt.Println(arr)

		// fmt.Println(arr[0][1])

		adj = make([][]int, n)

		for i := 0; i < n; i++ {
			adj[i] = make([]int, n)
			for j := 0; j < n; j++ {
				adj[i][j] = arr[j][i]
			}
		}



		for j:=0;j<node_num-1;j++{

			selected_list := opt.CallKumiawase_Impression(adj , j, j+4, SeedSet_F_strong2, prob_map, pop_list, interest_list, assum_list)
			fmt.Print(j,":")

			opt.Selected_Suppression(adj, selected_list, SeedSet_F_strong2,  prob_map , pop_list, interest_list, assum_list)
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
	pram_data := Parameter{GraphPath: adjFilePath, Node_n: n, Random_seed: seed, K_F: K_F, K_T: K_T, Mont_sample_size: sample_size, Pop_list: pop_list, SeedSet_F: seedsetf_pram, Interest_list: interest_list, Assum_list: assum_list, Seq: seq, Prob_map: prob_map}

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
	str := strings.Replace(time.Now().Format(layout2), ":", "-", -1)
	folder_path := "result/" + strconv.Itoa(n) + "node" + str
	err = os.Mkdir(folder_path, os.ModePerm)
	if err != nil {
		fmt.Println("error create" + folder_path)
		log.Fatal(err)
	}

	file, _ := json.MarshalIndent(pram_data, "", " ")
	_ = ioutil.WriteFile(folder_path+"/param_json.txt", file, 0644)
	//test part
	var S []int
	var hist [][]float64
	//

	SeedSet_F_strong := make([]int, len(adj))
	SeedSet_F_strong[0] = 1 //here
	// SeedSet_Greedy[1] = 1
	//偽情報の発信源を色々と
	// greedy_ans1, greedy_value1, greedy_value21 := opt.Greedy(1, 100, adj, SeedSet_F_strong, prob_map, pop_list, interest_list, assum_list, 3, true, 1000)
	// fmt.Println(greedy_ans1, greedy_value1, greedy_value21)
	// os.Exit(0)

	sample_size = 1000000
	S, hist = sim_submod(adj, sample_size, pop_list, interest_list, assum_list, SeedSet_F_strong, K_T, prob_map, folder_path)

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

func sim_submod(adj [][]int, sample_size int, pop_list [2]int, interest_list [][]int, assum_list [][]int, SeedSet_F []int, K_T int, prob_map [2][2][2][2]float64, folder_path string) ([]int, [][]float64) {
	var S []int
	var hist [][]float64
	S, hist = opt.Check_submod(1, K_T, sample_size, adj, SeedSet_F, prob_map, pop_list, interest_list, assum_list, folder_path)

	return S, hist
}

func main() {
	var withFalseInfo = [][][]int{
        {
            {29, 120, 166, 217},
            {29, 34, 118},
            {29, 34, 67, 69, 92, 115, 297},
            {29, 34, 69, 92, 182, 472},
        },
        {
            {34, 158, 219, 310},
            {29, 34, 360, 395},
            {29, 34, 258, 269, 395, 414, 473},
            {29, 34, 269, 280, 285, 297, 375},
        },
        {
            {174},
            {34, 327, 395},
            {190, 230, 278},
        },
        {
            {34, 230, 395, 404, 413, 473},
            {177},
            {67, 92, 108, 146, 168},
            {34, 92, 220},
        },
        {
            {34, 158, 321, 407, 474},
            {34, 220, 258},
            {29, 34, 152, 166, 208, 258},
            {19, 29, 208, 419, 438},
        },
        {
            {29, 34, 158, 196, 239, 366},
            {143},
            {29, 34, 168, 192, 196, 264, 395},
            {29, 34, 192, 196, 230, 258, 326, 361},
        },
        {
            {19, 34, 124, 273},
            {29, 34, 258, 309},
            {19, 34, 322, 419},
            {29, 34, 258, 322, 352, 366, 463},
        },
        {
            {6, 34, 98, 100, 120, 196, 202},
            {6, 34, 69, 157},
            {6, 34, 147, 258, 333, 404},
            {6, 34, 147, 389},
        },
        {
            {34, 258, 395, 407},
            {34, 322},
            {29, 34, 67, 258, 269, 280, 321, 361, 364},
            {34, 276, 393},
        },
    }

    // 虚偽情報なしの解
    var withoutFalseInfo = [][][]int{
        {
            {29, 34, 106, 206},
            {29, 34, 218},
            {29, 34, 83, 92, 115, 166},
            {29, 39, 83, 92, 98, 168, 189, 316, 329},
        },
        {
            {285, 288},
            {29, 188},
            {6, 29, 134, 152, 189, 297},
            {29, 69, 134, 146, 166, 173, 189, 235, 252, 395, 405},
        },
        {
            {6, 67, 217},
            {137},
            {29, 190, 235, 284, 364},
        },
        {
            {395, 404, 420},
            {34, 278},
            {92, 98, 182, 267},
            {34, 82, 92, 98, 103, 196},
        },
        {
            {29, 144},
            {258, 290},
            {158, 190, 196, 258, 280, 284},
            {34, 152, 187, 190, 196, 273, 342},
        },
        {
            {29, 42, 274, 285},
            {122},
            {29, 192, 196, 202, 230, 264},
            {29, 42, 192, 196, 202, 221, 230, 390},
        },
        {
            {34, 174, 258, 363},
            {34, 317},
            {34, 52, 190, 258},
            {34, 67, 190, 258, 395, 404, 412},
        },
        {
            {6, 67, 115, 413},
            {347},
            {69, 111, 182, 285},
            {6, 29, 111, 173, 230, 258, 308, 430},
        },
        {
            {34, 43, 162},
            {34, 258, 472},
            {29, 67, 269, 321, 342},
            {29, 43, 67, 269, 282, 285, 321, 345, 353, 395, 474},
        },
    }

	for i:=0;i<9;i++{
		fmt.Println()
		fmt.Println()
		// user_weight := 0.1*float64(i)
		user_weight := 0.0
		// fmt.Println("user_weight",user_weight)
		fmt.Println("seed",i)
		rand.Seed(int64(i))


		var seed int64 = int64(i)
		adjFilePath := "adj_jsonTwitterInteractionUCongress.txt"
		adj,interest_list,assum_list := Make_adj_interest_assum(adjFilePath,seed)
		use_user := false
		use_infl := true
		use_kaiki := false
		S_f_type := 1
		// num2 := 0


		// adjFilePath = "Graphs/adj_json50node.txt"
		adjFilePath = "adj_jsonTwitterInteractionUCongress.txt"
		adj,interest_list,assum_list = Make_adj_interest_assum(adjFilePath,seed)
		fmt.Println("len adj",len(adj))
		// os.Exit(0)
		if i == 0{

			cal_max_users(adj,7)
		}
		capacity := 302.0
		//コスト=拡散量用

		// os.Exit(0)

		use_user = false
		use_infl = true
		use_kaiki = false
		S_f_type = 2
		for j:=1.0;j<5.0;j++{
			if(i==2 && j == 4.0){
				continue
			}
			if use_infl{
				capacity = j*100
			}else{
				capacity = j
			}
			DP_ans_d := withFalseInfo[i][int(j)-1]
			DP_ans_s := withoutFalseInfo[i][int(j)-1]
			use_DP(adj,interest_list,assum_list,user_weight,capacity,use_kaiki,use_user,use_infl,1,S_f_type,false,DP_ans_d,DP_ans_s)
		}
		//
		// fmt.Println()
		// use_strict(adj,interest_list,assum_list,user_weight)
	}
}
