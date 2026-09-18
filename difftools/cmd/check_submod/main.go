// main in sample_check_submodu.go
package main

import (
	diff "difftools/diffusion"
	exp "difftools/experiment"
	"difftools/network"
	"difftools/optimization"
	"fmt"
	"math/rand"
)

func main() {
	var num_pick_users int = 7
	for i := 0; i < 9; i++ {
		fmt.Println()
		fmt.Println()
		// user_weight := 0.1*float64(i)
		user_weight := 0.0
		// fmt.Println("user_weight",user_weight)
		fmt.Println("seed", i)
		r := rand.New(rand.NewSource(int64(i)))
		// var seed int64 = int64(i)
		// adjFilePath := "adj_jsonTwitterInteractionUCongress.txt"
		// adjFilePath := "community_31.txt"
		// adj, interest_list, assum_list := Make_adj_interest_assum(adjFilePath, seed)
		// use_user := true     //コスト：ユーザ
		// use_infl := false    //コスト：拡散量
		// use_kaiki := false   //コスト：予想拡散量　使ってない
		// use_follower := true //コスト：総フォロワー数　使ってない
		// S_f_type := 1
		// num2 := 0

		// adjFilePath = "Graphs/adj_json50node.txt"
		adjFilePath := "community_21_adjmat.txt"
		use_congress := true
		// use_congress := true
		// adjFilePath = "adj_json_egoTwitter_kirinuki.txt"
		// adj, interest_list, assum_list := exp.Make_adj_interest_assum(adjFilePath, seed)

		net := network.ReadJson(adjFilePath)
		interestList := diff.MakeInterestList(net.N, r)
		assumList := diff.MakeAssumList(net.N, r)
		probTable := diff.GetUserProbTable()

		fmt.Println("len adj", net.N)
		// os.Exit(0)
		if i == 0 {
			exp.CalMaxUsers(net, num_pick_users)
			// exp.CalMaxUsersFixed(net, num_pick_users)
		}
		capacity := 302.0
		//コスト=拡散量用

		// use_user := false     //コスト：ユーザ
		// use_infl := true      //コスト：拡散量
		// use_kaiki := false    //コスト：予想拡散量　使ってない
		// use_follower := false //コスト：総フォロワー数　使ってない
		costFunc := optimization.CostInfluence

		S_f_type := 2 //1:単独　2:複数
		for j := 1.0; j < 5.0; j++ {
			switch costFunc {
			case optimization.CostInfluence:
				if use_congress {
					capacity = j * 100
				} else {
					capacity = j * 5
				}
			case optimization.CostFollower:
				capacity = j * 30
			default:
				capacity = j
			}
			// if costFunc == optimization.CostInfluence {
			// }
			// if use_infl && use_congress {
			// 	capacity = j * 100
			// 	// fmt.Println("okokokok")
			// } else if use_infl && !use_congress {
			// 	capacity = j * 5
			// } else if use_follower {
			// 	capacity = j * 30
			// } else {
			// 	capacity = j
			// }
			// fmt.Println(capacity)
			//
			exp.ComputeMaximizationDP(
				net,
				interestList,
				assumList,
				probTable,
				user_weight,
				capacity,
				// use_kaiki,
				// use_user,
				// use_infl,
				// use_follower,
				costFunc,
				1,
				S_f_type,
				false,
				r,
			)
		}
		//
		// fmt.Println()
		// use_strict(adj,interest_list,assum_list,user_weight)
	}
}
