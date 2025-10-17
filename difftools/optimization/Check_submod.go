package optimization

import (
	"encoding/csv"
	"fmt"
	"log"
	diff "m/difftools/diffusion"
	"math/rand"
	"os"
	"strconv"
	"strings"

	// "time"
	"sort"
)

func Check_submod(seed int64, k int, sample_size int, adj [][]int, SeedSet_F []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int, folder_path string) ([]int, [][]float64) {

	var n int = len(adj)
	var S []int = make([]int, n)

	for i, f := range SeedSet_F {
		if f > 0 {
			S[i] = 1
		}
	}

	var hist [][]float64 = make([][]float64, k)

	for i := 0; i < k; i++ {
		hist[i] = make([]float64, n)
	}

	sizes := []int{3, 4, 5}
	loop := 400
	// sets_len := len(sizes)*loop

	mont_loop := 1
	// mont_num := sample_size*mont_loop

	temp := make([][13]float64, len(sizes)*loop)
	temp[1][1] = 1.0

	var s_dist []float64 = make([]float64, n)

	//create file
	new_folder_path := folder_path + "/Check_submod"
	err := os.Mkdir(new_folder_path, os.ModePerm)
	if err != nil {
		fmt.Println("error create Check_submod")
		log.Fatal(err)
	}
	filename := new_folder_path + "/random" + strconv.Itoa(sample_size) + ".csv"
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	w := csv.NewWriter(f)

	colmns := []string{"K_T", "SetA", "SetB", "SetA_r", "SetB_r", "AandB_r", "AorB_r", "IsSubmodularity", "menos"}
	w.Write(colmns)
	//main loop
	for j := 0; j < len(sizes); j++ {
		fmt.Println(("percent"))
		fmt.Println(j, len(sizes))
		for v := 0; v < loop; v++ {

			var SetA []int
			var SetB []int

			Sets := make([][]int, 2)
			Sets[0] = make([]int, n)
			Sets[1] = make([]int, n)

			SetA = make([]int, n)
			_ = copy(SetA, S)
			// setA_list := Make_SeedSet_T_Random(SetA, sizes[j], adj)
			setA_list := Make_SeedSet_T_Strong(SetA, sizes[j], adj, 10)
			SetB = make([]int, n)
			_ = copy(SetB, S)
			// setB_list := Make_SeedSet_T_Random(SetB, sizes[j], adj)
			setB_list := Make_SeedSet_T_Strong(SetB, sizes[j], adj, 10)

			Sets[0] = setA_list
			Sets[1] = setB_list

			result := make([]float64, 4)
			// conf := make([]float64, 4)

			var SetAandB []int
			SetAandB = make([]int, n)
			_ = copy(SetAandB, S)
			append_seedset_T(SetAandB, Set_Multi(setA_list, setB_list))

			var SetAorB []int
			SetAorB = make([]int, n)
			_ = copy(SetAorB, S)
			append_seedset_T(SetAorB, Set_Sum(setA_list, setB_list))

			Set_use := make([][]int, 4)
			Set_use[0] = SetA
			Set_use[1] = SetB
			Set_use[2] = SetAandB
			Set_use[3] = SetAorB

			for i, set := range Set_use {
				dist := Infl_prop_exp(-1, sample_size*mont_loop, adj, set, prob_map, pop, interest_list, assum_list)
				result[i] = dist[diff.InfoType_T]
				//here
			}

			Sets_string := make([][]string, 2)
			Sets_string[0] = Int_to_String(Sets[0])
			Sets_string[1] = Int_to_String(Sets[1])

			part0 := []string{strings.Join(Sets_string[0], "-"), strings.Join(Sets_string[1], "-")} //here

			a := []float64{result[0], result[1], result[2], result[3], BoolToInt(result[0]+result[1] >= result[2]+result[3]), (result[0] + result[1]) - (result[2] + result[3])}

			part1 := Float_to_String(a)

			retu := append(part0, part1...)

			w.Write(retu)

		}

		//毎回初期化するけど宣言は外でできる？

		// for i:=0;i<n;i++{
		//   if
		s_dist[j] = 0
		// }
	}

	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	return S, hist

}

func FocusLoop(loop_n int, list1 []int, list2 []int, SeedSet_F []int, seed int64, sample_size int, adj [][]int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int, folder_path string) {

	_ = seed

	// now := time.Now()

	n := len(adj)
	SetA := make([]int, n)
	_ = copy(SetA, SeedSet_F)

	for _, n := range list1 { //多分appendSeedsetTでやれる
		if SetA[n] == 1 {
			fmt.Println("ここはS_fのところです")
		} else {
			SetA[n] = 2
		}
	}

	SetB := make([]int, n)
	_ = copy(SetB, SeedSet_F)

	for _, n := range list2 {
		if SetB[n] == 1 {
			fmt.Println("ここはS_fのところです")
		} else {
			SetB[n] = 2
		}
	}

	// var SetAandB []int
	SetAandB := make([]int, n)
	_ = copy(SetAandB, SeedSet_F)
	append_seedset_T(SetAandB, Set_Multi(list1, list2))

	// var SetAorB []int
	SetAorB := make([]int, n)
	_ = copy(SetAorB, SeedSet_F)
	append_seedset_T(SetAorB, Set_Sum(list1, list2))

	Set_use := make([][]int, 4)
	Set_use[0] = SetA
	Set_use[1] = SetB
	Set_use[2] = SetAandB
	Set_use[3] = SetAorB

	result := make([]float64, 4)

	filename := folder_path + "/focus" + strings.Join(Int_to_String(list1), "-") + strings.Join(Int_to_String(list2), "-") + strconv.Itoa(sample_size) + "-" + strconv.Itoa(loop_n) + ".csv"

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	w := csv.NewWriter(f)

	colmns := []string{"K_T", "SetA", "SetB", "SetA_r", "SetB_r", "AandB_r", "AorB_r", "IsSubmodularity", "menos"}
	w.Write(colmns)

	for j := 0; j < loop_n; j++ {
		if j%(loop_n/10) == 0 {
			fmt.Println("abc")
		}

		for i, set := range Set_use {
			dist := Infl_prop_exp(-1, sample_size, adj, set, prob_map, pop, interest_list, assum_list)
			result[i] = dist[diff.InfoType_T]
		}

		Sets_string := make([][]string, 2)
		Sets_string[0] = Int_to_String(list1)
		Sets_string[1] = Int_to_String(list2)

		part0 := []string{strings.Join(Sets_string[0], "-"), strings.Join(Sets_string[1], "-")}

		a := []float64{result[0], result[1], result[2], result[3], BoolToInt(result[0]+result[1] >= result[2]+result[3]), (result[0] + result[1]) - (result[2] + result[3])}

		part1 := Float_to_String(a)

		retu := append(part0, part1...)

		w.Write(retu)

	}

	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func Make_SeedSet_T_Random(Su []int, k int, adj [][]int) []int {
	n := len(Su)
	var sets []int
	Set := make([]int, 0, len(Su)) //選ばれる可能性があるノードたち(出次数が1以上)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if adj[i][j] > 0 {
				if Su[i] == 0 {
					Set = append(Set, i)
					break
				}
			}
		}
	}
	if len(Set) < k {
		k = len(Set)
		fmt.Println("十分な数の候補がありません")
	}

	for i := 0; i < k; {
		result := Set[rand.Intn(len(Set))]

		if Su[result] == 0 {
			Su[result] = 2
			sets = append(sets, result)
			i++
		}
	}
	return sets
}

type Node_power struct {
	NodeName int
	power    int //次数
}

func Make_SeedSet_T_Strong(Su []int, k int, adj [][]int, size int) []int {
	n := len(Su)
	var sets []int
	Set := make([]Node_power, size)
	var count int

	for i := 0; i < n; i++ {
		count = 0
		for j := 0; j < n; j++ {
			if adj[i][j] > 0 {
				if Su[i] == 0 {
					count++
				}
			}
		}
		node := Node_power{i, count}

		if Set[size-1].power < node.power {
			Set[size-1] = node
			sort.Slice(Set, func(i, j int) bool { return Set[i].power > Set[j].power })
		}

	}
	sets = append(sets, Set[0].NodeName) //一番強い奴は固定
	Su[Set[0].NodeName] = 2
	for i := 0; i < k-1; {

		result := Set[rand.Intn(len(Set))].NodeName

		if Su[result] == 0 {
			Su[result] = 2
			sets = append(sets, result)
			i++
		}
	}
	return sets
}

func append_seedset_T(Su []int, seedset_t_list []int) {
	for _, i := range seedset_t_list {
		Su[i] = 2
	}
}

func Set_Has(setA []int, num int) bool {
	for _, a := range setA {
		if a == num {
			return true
		}
	}
	return false
}

func Set_Sum(setA []int, setB []int) []int {
	ans := make([]int, len(setA), len(setA)+len(setB))
	_ = copy(ans, setA)

	for _, b := range setB {
		if Set_Has(ans, b) == false {
			ans = append(ans, b)
		}
	}
	return ans
}

func Set_Multi(setA []int, setB []int) []int {
	ans := make([]int, 0, len(setA))
	for _, b := range setB {
		if Set_Has(setA, b) {
			ans = append(ans, b)
		}
	}
	return ans
}

func Slice_Sum(slice []float64) float64 {
	var ans float64
	for _, n := range slice {
		ans = ans + n
	}
	return ans
}

func Int_to_String(slice []int) []string {
	l1 := len(slice)
	ans := make([]string, l1)
	for i := 0; i < l1; i++ {
		ans[i] = strconv.Itoa(slice[i])
	}
	return ans
}

func Float_to_String(slice []float64) []string {
	l1 := len(slice)
	ans := make([]string, l1)
	for i := 0; i < l1; i++ {
		ans[i] = strconv.FormatFloat(slice[i], 'f', 5, 64)
	}
	return ans
}

// func IntSlice_to_csv(slice [][]int, filename string, colmns []string){
//   f, err := os.Create(filename)
//     if err != nil {
//         log.Fatal(err)
//     }
//
//     w := csv.NewWriter(f)
//
//     w.Write(colmns)
//     strings := Int_to_String(slice)
//     fmt.Println(strings)
//     w.WriteAll(strings)
//
//     w.Flush()
//
//     if err := w.Error(); err != nil {
//         log.Fatal(err)
//     }
//
//
// }

func BoolToInt(b bool) float64 {
	if b {
		return 1.0
	} else {
		return 0.0
	}
}
