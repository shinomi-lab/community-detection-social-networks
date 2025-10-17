package optimization

import (
	"os"
	"fmt"
	// diff "m/difftools/diffusion"
  "bufio"
)

func SameImpressionCostFollower(sample_size int, adj [][]int, SeedSet_F []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int,under int, upper int, exit_f bool,use_cost_infl bool){
  fmt.Println("calling SameImpressionCostFollower")
  zerolist := make([]int,len(adj))
  sameFollowerList := CallKumiawase(adj, under,upper, SeedSet_F,true,prob_map,pop,interest_list,assum_list,use_cost_infl)
	

  fmt.Println("start selected suprresion return list",sameFollowerList)
  if exit_f{
    zerolist = SeedSet_F
  }
  suppList := Selected_SuppressionReturnList(adj, sameFollowerList, zerolist,  prob_map , pop, interest_list, assum_list)
  file, err := os.Create("SameImporessionCostFollower.csv")
  if err != nil {
      fmt.Println("Error creating file:", err)
      return
  }
  // 関数の終了時にファイルを閉じる
  defer file.Close()
  writer := bufio.NewWriter(file)

	for j := 0; j < len(sameFollowerList); j++ {

    // ファイルに書き込む
    _, err = writer.WriteString(fmt.Sprintf("%d",j)+","+fmt.Sprintf("%d",len(sameFollowerList[j]))+","+fmt.Sprintf("%f",suppList[j])+"\n")
    // _, err = file.WriteString("aaaaa\n")
    // fmt.Println("aa")
    if err != nil {
        fmt.Println("Error writing to file:", err)
        return
    }

	}
  err = writer.Flush()
    if err != nil {
        fmt.Println("Error flushing buffer:", err)
        return
    }
}
