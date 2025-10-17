package optimization

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	diff "m/difftools/diffusion"
	"math/rand"
	"sort"
	// "os"
)

func Strict(seed int64, sample_size int, adj [][]int, Seed_set []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int, ans_len int, Count_true bool, sample_size2 int) ([]int, float64, float64) {
	var n int = len(diff.Set)
	var max float64 = 0
	var result float64
	var ans []int

	ans = make([]int, 0, ans_len)
	S := make([]int, len(Seed_set))
	_ = copy(S, Seed_set)
	S_test := make([]int, len(Seed_set))
	_ = copy(S_test, Seed_set)

	var info_num int

	if Count_true {
		info_num = 2
	} else {
		info_num = 1
	}

	for i := 0; i < n; i++ {
		if S_test[diff.Set[i]] != 0 { //すでに発信源のユーザだったら
			continue
		}
		for j := i; j < n; j++ {
			if S_test[diff.Set[j]] != 0 { //すでに発信源のユーザだったら
				continue
			}
			for k := j; k < n; k++ {
				if S_test[diff.Set[k]] != 0 { //すでに発信源のユーザだったら
					continue
				}
				//main loop
				_ = copy(S_test, S)
				S_test[diff.Set[i]] = info_num
				S_test[diff.Set[j]] = info_num
				S_test[diff.Set[k]] = info_num
				//complete set Seedsets

				dist := Infl_prop_exp(seed, sample_size, adj, S_test, prob_map, pop, interest_list, assum_list)
				if Count_true {
					result = dist[diff.InfoType_T]
				} else {
					result = dist[diff.InfoType_F]
				}
				if result > max {
					max = result
					ans = []int{diff.Set[i], diff.Set[j], diff.Set[k]} //いける？
				}

			}
		}
	} //mainloop end
	_ = copy(S_test, S)
	for i := 0; i < len(ans); i++ {
		S_test[ans[i]] = info_num
	}

	//complete set Seedsets

	var max2 float64
	dist2 := Infl_prop_exp(seed, sample_size2, adj, S_test, prob_map, pop, interest_list, assum_list)
	if Count_true {
		max2 = dist2[diff.InfoType_T]
	} else {
		max2 = dist2[diff.InfoType_F]
	}

	return ans, max, max2
}

func Strict2(seed int64, sample_size int, adj [][]int, Seed_set []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int, ans_len int, Count_true bool, sample_size2 int, under int, upper int) ([]int, float64, float64) {
	var n int = len(diff.Set)
	var max float64 = 0
	var result float64
	var ans []int

	ans = make([]int, 0, ans_len)
	S := make([]int, len(Seed_set))
	_ = copy(S, Seed_set)
	S_test := make([]int, len(Seed_set))
	_ = copy(S_test, Seed_set)

	var info_num int

	if Count_true {
		info_num = 2
	} else {
		info_num = 1
	}

	for i := 0; i < n; i++ {
		if S_test[diff.Set[i]] != 0 { //すでに発信源のユーザだったら
			continue
		}
		for j := i; j < n; j++ {
			if S_test[diff.Set[j]] != 0 { //すでに発信源のユーザだったら
				continue
			}
			for k := j; k < n; k++ {
				if S_test[diff.Set[k]] != 0 { //すでに発信源のユーザだったら
					continue
				}
				//main loop
				_ = copy(S_test, S)
				S_test[diff.Set[i]] = info_num
				S_test[diff.Set[j]] = info_num
				S_test[diff.Set[k]] = info_num
				//complete set Seedsets

				dist := Infl_prop_exp(seed, sample_size, adj, S_test, prob_map, pop, interest_list, assum_list)
				if Count_true {
					result = dist[diff.InfoType_T]
				} else {
					result = dist[diff.InfoType_F]
				}
				if result > max {
					max = result
					ans = []int{diff.Set[i], diff.Set[j], diff.Set[k]} //いける？
				}

			}
		}
	} //mainloop end
	_ = copy(S_test, S)
	for i := 0; i < len(ans); i++ {
		S_test[ans[i]] = info_num
	}

	//complete set Seedsets

	var max2 float64
	dist2 := Infl_prop_exp(seed, sample_size2, adj, S_test, prob_map, pop, interest_list, assum_list)
	if Count_true {
		max2 = dist2[diff.InfoType_T]
	} else {
		max2 = dist2[diff.InfoType_F]
	}

	return ans, max, max2
}

// public values
var aaa [][]int
var saiki int
var ketteizumi int
var counter int
var infler_cost_list_copy []int
var prob_map_copy [2][2][2][2]float64
var pop_copy [2]int
var interest_list_copy [][]int
var assum_list_copy [][]int

// [0,1,1,0,1]→[1,2,4]に型変換して代入
func printCombination(pattern []int, elems []int, n int) {
	oneOf := make([]int, 0)
	for i := 0; i < n; i++ {
		//同じのを複数個選べるようにfor文になっている(pattern[i] = 1だったら追加)
		for j := 0; j < pattern[i]; j++ {
			oneOf = append(oneOf, elems[i])
		}
	}
	aaa = append(aaa, oneOf)
}

/* n個の要素からr個の要素を選ぶ場合の全パターンを列挙する */
func combination(adj [][]int, pattern []int, elems []int, n int, undder int, upper int, num_decided int, OnlyInfler bool, use_cost_infl bool) {
	var use_get_num_selected func([][]int, []int, int, []int) int

	if use_cost_infl {
		use_get_num_selected = getNumSelected_infl
	} else {
		use_get_num_selected = getNumSelected
	}
	// fmt.Println("num_decided:",num_decided)
	num_selected := use_get_num_selected(adj, pattern, num_decided, elems)

	// if(num_decided != counter){
	//   counter = num_decided
	//   fmt.Println("plus one",counter,n)
	// }

	if num_decided == n {
		ketteizumi = ketteizumi + 1
		if ketteizumi%1000 == 0 {
			fmt.Println("決定", ketteizumi)
		}

		/* n個全ての要素に対して"選ぶ"or"選ばない"が決定ずみ */
		if num_selected <= upper && num_selected > undder {
			/* r個だけ選ばれている場合のみ、選ばれた要素を表示 */
			printCombination(pattern, elems, n)
		}
		return
	}

	/* num_decided個目の要素を"選ばない"場合のパターンを作成 */
	pattern[num_decided] = 0
	combination(adj, pattern, elems, n, undder, upper, num_decided+1, OnlyInfler, use_cost_infl)
	if num_selected <= upper {
		/* num_decided個目の要素を"選ぶ"場合のパターンを作成 */
		if OnlyInfler {
			if FolowerSize(adj, num_decided) != 0 {
				pattern[num_decided] = 1
				combination(adj, pattern, elems, n, undder, upper, num_decided+1, OnlyInfler, use_cost_infl)
			} else {
				// fmt.Println("除外している")
			}
		} else {
			pattern[num_decided] = 1
			combination(adj, pattern, elems, n, undder, upper, num_decided+1, OnlyInfler, use_cost_infl)

		}

	} else {
		// fmt.Println("除外している")
	}
}

/* n個の要素からr個の要素を選ぶ場合の全パターンを列挙する */
func combination2(adj [][]int, pattern []int, elems []int, n int, undder float64, upper float64, num_decided int, OnlyInfler bool, max_user int, user_weight float64) {

	num_selected := getNumSelected2(adj, pattern, num_decided, elems, max_user, user_weight)

	if num_decided == n {
		/* n個全ての要素に対して"選ぶ"or"選ばない"が決定ずみ */
		if num_selected <= upper && num_selected > undder {
			/* r個だけ選ばれている場合のみ、選ばれた要素を表示 */
			printCombination(pattern, elems, n)
		}
		return
	}

	/* num_decided個目の要素を"選ばない"場合のパターンを作成 */
	pattern[num_decided] = 0
	combination2(adj, pattern, elems, n, undder, upper, num_decided+1, OnlyInfler, max_user, user_weight)

	/* num_decided個目の要素を"選ぶ"場合のパターンを作成 */
	if OnlyInfler {
		if FolowerSize(adj, num_decided) != 0 {
			pattern[num_decided] = 1
			combination2(adj, pattern, elems, n, undder, upper, num_decided+1, OnlyInfler, max_user, user_weight)
		}
	} else {
		pattern[num_decided] = 1
		combination2(adj, pattern, elems, n, undder, upper, num_decided+1, OnlyInfler, max_user, user_weight)

	}
}

func SameImporession(adj [][]int, pattern []int, elems []int, n int, undder int, upper int, num_decided int, SeedSet []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int) {

	num_selected := getNumSelected_SameImpression(adj, pattern, num_decided, elems, SeedSet, prob_map, pop, interest_list, assum_list)

	if num_decided == n {
		/* n個全ての要素に対して"選ぶ"or"選ばない"が決定ずみ */
		if num_selected <= float64(upper) && num_selected > float64(undder) {
			/* r個だけ選ばれている場合のみ、選ばれた要素を表示 */
			printCombination(pattern, elems, n)
		}
		return
	}

	/* num_decided個目の要素を"選ばない"場合のパターンを作成 */
	pattern[num_decided] = 0
	SameImporession(adj, pattern, elems, n, undder, upper, num_decided+1, SeedSet, prob_map, pop, interest_list, assum_list)

	/* num_decided個目の要素を"選ぶ"場合のパターンを作成 */
	infler_num := 0
	for i := 0; i < len(adj); i++ {
		if adj[num_decided][i] != 0 {
			infler_num = 1
			break
		}
	}
	if infler_num != 0 {
		pattern[num_decided] = 1
		SameImporession(adj, pattern, elems, n, undder, upper, num_decided+1, SeedSet, prob_map, pop, interest_list, assum_list)

	}
}

func getNumSelected(adj [][]int, pattern []int, n int, elems []int) int {
	/* "選ぶ"と決定された要素の数を計算 */
	// printf("pattern\t");
	num_selected := 0
	for i := 0; i < n; i++ {
		if pattern[i] == 1 {
			num_selected += FolowerSize(adj, elems[i])
		}
	}
	return num_selected
}

func getNumSelected_infl(adj [][]int, pattern []int, n int, elems []int) int {
	/* "選ぶ"と決定された要素の数を計算 */
	// printf("pattern\t");
	num_selected := 0
	for i := 0; i < n; i++ {
		if pattern[i] == 1 {
			num_selected += Cal_cost_infl_int(adj, elems[i], prob_map_copy, pop_copy, interest_list_copy, assum_list_copy)
		}
	}
	return num_selected
}

func getNumSelected2(adj [][]int, pattern []int, n int, elems []int, max_user int, user_weight float64) float64 {
	/* "選ぶ"と決定された要素の数を計算 nownow*/
	// printf("pattern\t");
	num_selected := 0.0
	for i := 0; i < n; i++ {
		if pattern[i] == 1 {
			num_selected += Cal_cost(user_weight, 1-user_weight, adj, elems[i], max_user)
		}
	}
	// fmt.Println("pattern:",pattern,"elems:",elems)
	// os.Exit(0)
	return num_selected
}

func getNumSelected_SameImpression(adj [][]int, pattern []int, n int, elems []int, SeedSet []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int) float64 {
	/* "選ぶ"と決定された要素の数を計算 */
	// printf("pattern\t");

	S_test := make([]int, len(SeedSet))
	_ = copy(S_test, SeedSet)

	for i := 0; i < len(SeedSet); i++ {
		if pattern[i] == 1 {
			S_test[i] = 2 //true = 2
		}
	}

	dist2 := Infl_prop_exp(0, 100, adj, S_test, prob_map, pop, interest_list, assum_list) //後で直す

	ans := dist2[diff.InfoType_T]

	return ans
}

func FolowerSize(adj [][]int, node int) int {
	ans := 0
	for _, isEdge := range adj[node] {
		ans += isEdge
	}

	return ans
}

// func ImpressionSize(adj [][]int,node int, SeedSet []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int)int{
//   ans := 0
//
//   S_test := make([]int ,len(Seed_set))
//   _ = copy(S_test, Seed_set)
//
//   dist2 := Infl_prop_exp(0, 1000, adj, S_F, prob_map, pop, interest_list, assum_list)
//
//   ans = dist2[diff.InfoType_T]
//
//   return ans
// }

func CallKumiawase(adj [][]int, under int, upper int, SeedSet []int, OnlyInfler bool, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int, use_cost_infl bool) [][]int {
	ketteizumi = 0
	counter = 0
	aaa = make([][]int, 0)
	fmt.Println("calling CallKumiawase")
	//dest := make([]int, len(src))
	//public関数にコピー(毎回引数にするのがめんどくさいから(値の変更は行わないから宣言の仕方とか変えて中身を変更できないようにした方がよさそう))
	for i := range prob_map {
		for j := range prob_map[i] {
			for k := range prob_map[i][j] {
				for l := range prob_map[i][j][k] {
					prob_map_copy[i][j][k][l] = prob_map[i][j][k][l]
				}
			}
		}
	}
	pop_copy[0] = pop[0]
	pop_copy[1] = pop[1]

	for i := range interest_list {
		for j := range interest_list[i] {
			interest_list_copy[i][j] = interest_list[i][j]

		}
	}
	for i := range assum_list {
		for j := range assum_list[i] {
			assum_list_copy[i][j] = assum_list[i][j]

		}
	}
	//nを指定することで選べるユーザ数の上限を決めれる
	n := len(adj)
	var a int
	k := 0
	// n = 5
	//情報の発信源となりうる(出次数0を消す)ユーザをまとめる　a
	//型は[0,2,4,6,7]って感じ
	elems := make([]int, 0, n)
	for i := 0; i < n; i++ {
		a = 0
		if SeedSet[i] == 1 {
			continue
		}
		for j := 0; j < len(adj); j++ {
			if adj[i][j] != 0 {
				a = a + 1
			}
		}
		if a != 0 {
			elems = append(elems, i)
			k = k + 1
		}
	}
	//a End
	pattern := make([]int, n)
	saiki = 0
	fmt.Println("calling combination")
	combination(adj, pattern, elems, k, under, upper, 0, OnlyInfler, use_cost_infl)
	// fmt.Println("most important", aaa)
	return aaa
}

func CallKumiawase2(adj [][]int, under float64, upper float64, SeedSet []int, OnlyInfler bool, max_user int, user_weight float64) [][]int {
	aaa = make([][]int, 0)
	//nを指定することで選べるユーザ数の上限を決めれる
	n := len(adj)
	var a int
	k := 0
	// n = 5
	//情報の発信源となりうる(出次数0を消す)ユーザをまとめる　a
	//型は[0,2,4,6,7]って感じ
	elems := make([]int, 0, n)
	for i := 0; i < n; i++ {
		a = 0
		if SeedSet[i] == 1 {
			continue
		}
		for j := 0; j < len(adj); j++ {
			if adj[i][j] != 0 {
				a = a + 1
			}
		}
		if a != 0 {
			elems = append(elems, i)
			k = k + 1
		}
	}
	//a End
	pattern := make([]int, n)
	saiki = 0
	combination2(adj, pattern, elems, k, under, upper, 0, OnlyInfler, max_user, user_weight)
	// fmt.Println("most important", aaa)
	return aaa
}

// func CallKumiawase_infl(adj [][]int,under int, upper int, SeedSet []int, OnlyInfler bool)[][]int {
//     ketteizumi = 0
//     counter = 0
//     aaa = make([][]int,0)
//     fmt.Println("calling CallKumiawase")
//     //nを指定することで選べるユーザ数の上限を決めれる
//     n := len(adj)
//     var a int
//     k := 0
//     // n = 5
//     //情報の発信源となりうる(出次数0を消す)ユーザをまとめる　a
//     //型は[0,2,4,6,7]って感じ
//     elems := make([]int,0,n)
//     for i:=0;i<n;i++{
//       a = 0
//       if SeedSet[i] == 1{
//         continue
//       }
//       for j:=0;j<len(adj);j++{
//         if (adj[i][j] != 0){
//           a = a + 1
//         }
//       }
//       if(a != 0){
//         elems = append(elems,i)
//         k = k + 1
//       }
//     }
//     //a End
//     pattern := make([]int,n)
//     saiki = 0
//     fmt.Println("calling combination")
//     combination(adj, pattern, elems, k, under, upper, 0, OnlyInfler);
//     // fmt.Println("most important", aaa)
//     return aaa
// }

func CallKumiawase_Impression(adj [][]int, under int, upper int, SeedSet []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int) [][]int {
	aaa = make([][]int, 0)
	// fmt.Println("calling CallKumiawase")
	n := len(adj)
	var a int
	k := 0
	// n = 5
	elems := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if SeedSet[i] == 1 {
			continue
		}
		a = 0
		for j := 0; j < len(adj); j++ {
			if adj[i][j] != 0 {
				a = a + 1
			}
		}
		if a != 0 {
			elems = append(elems, i)
			// elems[k] = i
			k = k + 1
		}
	}

	pattern := make([]int, n)
	saiki = 0

	SameImporession(adj, pattern, elems, k, under, upper, 0, SeedSet, prob_map, pop, interest_list, assum_list)
	// fmt.Println("most important", aaa)
	return aaa
}

func RandomSuppression(adj [][]int, node_num int, SeedSet []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int, kurikaesi int, OnlyInfler bool) (float64, [][]float64, [][]int) {
	var ans float64
	var ans2 [][]float64
	var ans3 [][]int

	ans2 = make([][]float64, 0)
	ans3 = make([][]int, 0)
	ans = 0
	var num int

	S_test := make([]int, len(SeedSet))
	_ = copy(S_test, SeedSet)

	for z := 0; z < kurikaesi; z++ {
		ans2_v := make([]float64, 4)
		ans3_v := make([]int, 0)
		S_test := make([]int, len(SeedSet))
		_ = copy(S_test, SeedSet)
		for i := 0; i < node_num; i++ {
			num = rand.Intn(len(SeedSet))
			if S_test[num] != 0 {
				i = i - 1

			} else if OnlyInfler {
				adj_len := len(adj)
				f_num := 0
				for i := 0; i < adj_len; i++ {
					f_num += adj[num][i]
				}
				if f_num > 0 {
					S_test[num] = 2
					ans3_v = append(ans3_v, num)
					// fmt.Println("selected is ",num)
				} else {
					i = i - 1
				}

			} else {
				S_test[num] = 2
				ans3_v = append(ans3_v, num)
			}
		}
		sum := 0
		for _, x := range S_test {
			sum += x
		}
		sum = (sum - 1) / 2
		// fmt.Println("start")
		dist := Infl_prop_exp(0, 100, adj, S_test, prob_map, pop, interest_list, assum_list)
		// fmt.Println("end")
		ans2_v[0] = float64(node_num)
		ans2_v[1] = float64(CalFolower(adj, S_test))
		ans2_v[2] = dist[diff.InfoType_T]
		ans2_v[3] = dist[diff.InfoType_F]

		// fmt.Println("ans2_v",ans2_v[0],ans2_v[1],ans2_v[2])
		ans += dist[diff.InfoType_T]

		ans2 = append(ans2, ans2_v)
		ans3 = append(ans3, ans3_v)
		// ans += dist[diff.InfoType_T] - sum   //発信者自信を含まない拡散を出力
	}

	// fmt.Println("RandomSuppression return:",ans / float64(kurikaesi))
	fmt.Println(node_num, ":", ans/float64(kurikaesi))
	return ans / float64(kurikaesi), ans2, ans3

}

func PythonSuppression(adj [][]int, SeedSet []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int, OnlyInfler bool) ([]float64, [][]int) {
	var ans float64
	var ans2 []float64
	var ans3 [][]int

	ans2 = make([]float64, 0)
	ans3 = make([][]int, 0)
	ans = 0
	var num int

	S_test := make([]int, len(SeedSet))
	_ = copy(S_test, SeedSet)

	node_list_path := "Python_random_nodelists/node_list.txt"
	bytes, err := ioutil.ReadFile(node_list_path)
	if err != nil {
		panic(err)
	}

	var dataJson string = string(bytes)

	arr := make(map[int]map[int]int)
	// var arr []string
	_ = json.Unmarshal([]byte(dataJson), &arr)
	// fmt.Println(arr)

	// fmt.Println(arr[2][0])
	// fmt.Println(arr)
	// n := len(arr[0])
	var node_lists [][]int = make([][]int, len(arr))

	//make node_lists
	for i := 0; i < len(arr); i++ {
		node_lists[i] = make([]int, len(arr[i]))
		for j := 0; j < len(arr[i]); j++ {
			node_lists[i][j] = arr[i][j]
			// print(node_lists[i][j])
		}
		// print("---------------")
	}
	// fmt.Println(node_lists)

	for z := 0; z < len(node_lists); z++ {
		ans2_v := make([]float64, 4)
		ans3_v := make([]int, 0)
		S_test := make([]int, len(SeedSet))
		_ = copy(S_test, SeedSet)
		for i := 0; i < len(node_lists[z]); i++ {
			//
			num = node_lists[z][i]
			S_test[num] = 2
			// copy(ans3_v, node_list[z])
			//   if(S_test[num] != 0){
			//     i = i - 1
			//
			//   }else if(OnlyInfler){
			//     adj_len := len(adj)
			//     f_num := 0
			//     for i:=0;i<adj_len;i++{
			//       f_num += adj[num][i]
			//     }
			//     if(f_num > 0){
			//       S_test[num] = 2
			//       ans3_v = append(ans3_v,num)
			//       // fmt.Println("selected is ",num)
			//     }else{
			//       i = i - 1
			//     }
			//
			//   }else{
			//     S_test[num] = 2
			//     ans3_v = append(ans3_v,num)
			//   }
		}
		sum := 0
		for _, x := range S_test {
			sum += x
		}
		sum = (sum - 1) / 2
		// fmt.Println("start")
		dist := Infl_prop_exp(-1, 1000, adj, S_test, prob_map, pop, interest_list, assum_list)

		// print("あ\t")
		// fmt.Println("end")
		// ans2_v[0] = float64(node_num)
		ans2_v[1] = float64(CalFolower(adj, S_test))
		ans2_v[2] = dist[diff.InfoType_T]
		ans2_v[3] = dist[diff.InfoType_F]

		// fmt.Println("ans2_v",ans2_v[0],ans2_v[1],ans2_v[2])
		ans += dist[diff.InfoType_T]

		ans2 = append(ans2, ans2_v[2])
		ans3 = append(ans3, ans3_v)
		// ans += dist[diff.InfoType_T] - sum   //発信者自信を含まない拡散を出力
	}

	// fmt.Println("RandomSuppression return:",ans / float64(len(arr)))
	// fmt.Println(node_num,":",ans/float64(len(arr)))
	fmt.Println()
	return ans2, node_lists

}

func Selected_Suppression(adj [][]int, selected_list [][]int, SeedSet []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int) float64 {
	var ans float64
	ans = 0

	for i := 0; i < len(selected_list); i++ {
		S_test := make([]int, len(SeedSet))
		_ = copy(S_test, SeedSet)
		selected := selected_list[i]
		for j := 0; j < len(selected); j++ {
			node := selected[j]
			if S_test[node] != 0 {
				fmt.Println("ERROR in Selected_Suppression")
			}
			S_test[node] = 2
		}
		dist := Infl_prop_exp(1, 100, adj, S_test, prob_map, pop, interest_list, assum_list)

		ans += dist[diff.InfoType_T]
	}
	// for selected := range selected_list{
	//   for node := range selected{
	//     if S_test[node] != 0{
	//       fmt.Println("ERROR in Selected_Suppression")
	//     }
	//     S_test[num] = 2
	//   }
	//   dist := Infl_prop_exp(0, 1000, adj, S_test, prob_map, pop, interest_list, assum_list)
	//
	//   ans  += dist[diff.InfoType_T]
	// }
	if len(selected_list) == 0 {
		fmt.Println("selected_list is zero")
		return 0.0
	} else {
		// fmt.Println("Selected Suppression return:",ans / float64(len(selected_list)))
		fmt.Println(":", ans/float64(len(selected_list)))
		return ans / float64(len(selected_list))
	}
}

func Selected_SuppressionReturnList(adj [][]int, selected_list [][]int, SeedSet []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int) []float64 {
	// var ans float64
	ans := make([]float64, len(selected_list))

	for i := 0; i < len(selected_list); i++ {
		S_test := make([]int, len(SeedSet))
		_ = copy(S_test, SeedSet)
		selected := selected_list[i]
		for j := 0; j < len(selected); j++ {
			node := selected[j]
			if S_test[node] != 0 {
				fmt.Println("ERROR in Selected_Suppression")
			}
			S_test[node] = 2
		}
		dist := Infl_prop_exp(1, 100, adj, S_test, prob_map, pop, interest_list, assum_list) //here

		ans[i] = dist[diff.InfoType_T]
	}

	return ans
}

func Selected_Suppression_Maximum(adj [][]int, selected_list [][]int, SeedSet []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int) ([]int, float64, float64) {
	var ans float64
	var max_users []int
	ans = 0
	max := 0.0
	var diff_f float64

	for i := 0; i < len(selected_list); i++ {
		S_test := make([]int, len(SeedSet))
		_ = copy(S_test, SeedSet)
		selected := selected_list[i]
		sort.Slice(selected, func(i, j int) bool { return selected[i] < selected[j] })
		// rand.Seed(0)
		ans = 0
		for j := 0; j < len(selected); j++ {
			node := selected[j]
			if S_test[node] != 0 {
				fmt.Println("ERROR in Selected_Suppression")
			}
			S_test[node] = 2
		}
		// fmt.Println(S_test)
		rand.Seed(100)
		dist := Infl_prop_exp(-1, 1000, adj, S_test, prob_map, pop, interest_list, assum_list)

		ans = dist[diff.InfoType_T]
		diff_f = dist[diff.InfoType_F]
		// fmt.Println(ans)
		if ans > max {
			max = ans
			max_users = make([]int, len(selected))
			diff_f = dist[diff.InfoType_F]
			_ = copy(max_users, selected)
		}
	}

	if len(selected_list) == 0 {
		fmt.Println("selected_list is zero")
		slice := make([]int, 0)
		return slice, 0.0, 0.0
	} else {
		// fmt.Println("Selected Suppression return:",ans / float64(len(selected_list)))
		// fmt.Println(max)
		return max_users, max, diff_f
	}
}

func CalFolower(adj [][]int, nodes []int) int {
	ans := 0
	for node, v := range nodes {
		if v == 2 {
			ans += FolowerSize(adj, node)
		}
	}

	return ans
}

// func deepCopy(original [][]int) [][]int {
//     // コピー先のスライスを作成
//     copy := make([][]int, len(original))
//
//     // 各スライスをループして新しいメモリにコピー
//     for i := range original {
//         copy[i] = make([]int, len(original[i]))
//         copy(copy[i], original[i])
//     }
//
//     return copy
// }
