package main

import(
  // "encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	// "log"
	diff "m/difftools/diffusion"
	opt "m/difftools/optimization"
	// "os"
	// "strconv"
	// "strings"
	// "time"
	"math/rand"
	// "reflect"
)

func main(){
  var seed int64
  seed = 1
  adjFilePath := "Graphs/adj_json50node.txt"
  adj,interest_list,assum_list := Make_adj_interest_assum(adjFilePath,seed)
  var pop_list [2]int
  pop_list[0] = diff.Pop_high
  pop_list[1] = diff.Pop_high
  var seq [16]float64 = diff.Make_probability()

	var prob_map [2][2][2][2]float64 = diff.Map_probagbility(seq)

  InflTest(adj,prob_map,pop_list,interest_list,assum_list)
  Selected_Suppression_MaximumTest(adj,prob_map,pop_list,interest_list,assum_list)
}

func InflTest(adj [][]int,prob_map [2][2][2][2]float64,pop_list [2]int,interest_list [][]int,assum_list [][]int){

  SeedSetF := make([]int,len(adj))
  SeedSetF[2] = 1
  SeedSetF[0] = 2
  rand.Seed(100)
  hist := opt.Infl_prop_exp(0,1000,adj,SeedSetF,prob_map,pop_list,interest_list,assum_list)

  fmt.Println(hist)
}

func Selected_Suppression_MaximumTest(adj [][]int,prob_map [2][2][2][2]float64,pop_list [2]int,interest_list [][]int,assum_list [][]int){
  SeedSetF := make([]int,len(adj))
  SeedSetF[2] = 1
  greedy_ans := []int{0}
  greedy_ans2 := make([][]int,0)
  greedy_ans2 = append(greedy_ans2,greedy_ans)

  _,test_greedy_ans_v := opt.Selected_Suppression_Maximum(adj,greedy_ans2, SeedSetF,  prob_map , pop_list, interest_list, assum_list)

  fmt.Println(test_greedy_ans_v)
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
