// main in sample_check_submodu_sub.go
package main

import (
	diff "difftools/diffusion"
	exp "difftools/experiment"
	"difftools/network"
	"fmt"
	"math/rand"
)

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

	for i := 0; i < 9; i++ {
		fmt.Println()
		fmt.Println()
		// user_weight := 0.1*float64(i)
		user_weight := 0.0
		// fmt.Println("user_weight",user_weight)
		fmt.Println("seed", i)
		r := rand.New(rand.NewSource(int64(i)))

		adjFilePath := "adj_jsonTwitterInteractionUCongress.txt"
		net := network.ReadAdjJson(adjFilePath)
		interestList := diff.MakeInterestList(net.N, r)
		assumList := diff.MakeAssumList(net.N, r)
		probTable := diff.GetUserProbTable()
		use_user := false
		use_infl := true
		use_kaiki := false
		S_f_type := 1
		// num2 := 0

		// adjFilePath = "Graphs/adj_json50node.txt"
		// adjFilePath = "adj_jsonTwitterInteractionUCongress.txt"
		// adj, interest_list, assum_list = Make_adj_interest_assum(adjFilePath, seed)
		fmt.Println("len adj", net.N)
		// os.Exit(0)
		if i == 0 {
			exp.CalMaxUsers(net, 7)
		}
		capacity := 302.0
		//コスト=拡散量用

		// os.Exit(0)

		use_user = false
		use_infl = true
		use_kaiki = false
		S_f_type = 2
		for j := 1.0; j < 5.0; j++ {
			if i == 2 && j == 4.0 {
				continue
			}
			if use_infl {
				capacity = j * 100
			} else {
				capacity = j
			}
			DP_ans_d := withFalseInfo[i][int(j)-1]
			DP_ans_s := withoutFalseInfo[i][int(j)-1]
			exp.ComputeDP(
				net,
				interestList,
				assumList,
				probTable,
				user_weight,
				capacity,
				use_kaiki,
				use_user,
				use_infl,
				1,
				S_f_type,
				false,
				DP_ans_d,
				DP_ans_s,
				r,
			)
		}
		//
		// fmt.Println()
		// use_strict(adj,interest_list,assum_list,user_weight)
	}
}
