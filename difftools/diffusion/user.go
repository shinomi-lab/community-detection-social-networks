package diffusion

import (
	"math"
	"math/rand"

	// "fmt"
	// "os"
	"sync"
)

const (
	Interest_low  int = 0
	Interest_high int = 1
	Interests_n   int = 2
)

const (
	Assum_F  int = 0
	Assum_T  int = 1
	Assums_n int = 2
)

func MakeInterestList(n int, r *rand.Rand) [][]int {
	var interestList = make([][]int, n)
	for i := range interestList {
		interestList[i] = make([]int, InfoTypes_n)
		interestList[i][InfoType_F] = r.Intn(2)
		interestList[i][InfoType_T] = r.Intn(2)
	}
	return interestList
}

func MakeAssumList(n int, r *rand.Rand) [][]int {
	var assumList = make([][]int, n)
	for i := range assumList {
		assumList[i] = make([]int, InfoTypes_n)
		assumList[i][Pop_low] = r.Intn(2)
		assumList[i][Pop_high] = r.Intn(2)

	}

	return assumList
}

func makeProbabilities() [16]float64 {
	var x [16]float64
	x[1] = 1

	for i := 1; i < 17; i++ {
		x[i-1] = math.Pow(10.0, float64(-i)/16.0) / 8
		// x[i-1] = math.Pow(10.0, float64(-i)/16.0)/32
	}

	return x
}

type UserProbTable [2][2][2][2]float64

func initProbabilityTable() UserProbTable {
	prob := makeProbabilities()

	prob_1011 := prob[0]
	prob_1111 := prob[1]
	prob_0011 := prob[2]
	prob_0111 := prob[3]
	prob_1001 := prob[4]
	prob_1101 := prob[5]
	prob_0001 := prob[6]
	prob_0101 := prob[7]

	prob_1000 := prob[8]
	prob_1100 := prob[9]
	prob_0000 := prob[10]
	prob_0100 := prob[11]

	prob_1010 := prob[12]
	prob_1110 := prob[13]
	prob_0010 := prob[14]
	prob_0110 := prob[15]

	var t UserProbTable
	t[Pop_low][InfoType_F][Interest_low][Assum_F] = prob_0000
	t[Pop_low][InfoType_F][Interest_low][Assum_T] = prob_0001
	t[Pop_low][InfoType_F][Interest_high][Assum_F] = prob_0010
	t[Pop_low][InfoType_F][Interest_high][Assum_T] = prob_0011
	t[Pop_low][InfoType_T][Interest_low][Assum_F] = prob_0100
	t[Pop_low][InfoType_T][Interest_low][Assum_T] = prob_0101
	t[Pop_low][InfoType_T][Interest_high][Assum_F] = prob_0110
	t[Pop_low][InfoType_T][Interest_high][Assum_T] = prob_0111
	t[Pop_high][InfoType_F][Interest_low][Assum_F] = prob_1000
	t[Pop_high][InfoType_F][Interest_low][Assum_T] = prob_1001
	t[Pop_high][InfoType_F][Interest_high][Assum_F] = prob_1010
	t[Pop_high][InfoType_F][Interest_high][Assum_T] = prob_1011
	t[Pop_high][InfoType_T][Interest_low][Assum_F] = prob_1100
	t[Pop_high][InfoType_T][Interest_low][Assum_T] = prob_1101
	t[Pop_high][InfoType_T][Interest_high][Assum_F] = prob_1110
	t[Pop_high][InfoType_T][Interest_high][Assum_T] = prob_1111

	return t
}

// 一度作ったテーブルをを再利用する
var (
	table UserProbTable
	once  sync.Once
)

func GetUserProbTable() UserProbTable {
	once.Do(func() {
		table = initProbabilityTable()
	})
	return table
}
