package diffusion

import (
	"math"
	"math/rand"
	// "fmt"
	// "os"
)

var Interest_low int = 0
var Interest_high int = 1
var Interests_n int = 2

var Assum_F int = 0
var Assum_T int = 1
var Assums_n int = 2

func Make_interest_list(n int, seed int64) [][]int {
	_ = seed
	var interest_list = make([][]int, n)
	for i := range interest_list {
		interest_list[i] = make([]int, InfoTypes_n)
		interest_list[i][InfoType_F] = rand.Intn(2)
		interest_list[i][InfoType_T] = rand.Intn(2)
	}
	return interest_list
}

func Make_assum_list(n int, seed int64) [][]int {
	_ = seed
	var assum_list = make([][]int, n)
	for i := range assum_list {
		assum_list[i] = make([]int, InfoTypes_n)
		assum_list[i][Pop_low] = rand.Intn(2)
		assum_list[i][Pop_high] = rand.Intn(2)

	}

	return assum_list
}

func Make_probability() [16]float64 {
	var x [16]float64
	x[1] = 1

	for i := 1; i < 17; i++ {
		x[i-1] = math.Pow(10.0, float64(-i)/16.0)/8
		// x[i-1] = math.Pow(10.0, float64(-i)/16.0)/32
	}

	return x
}

func Map_probagbility(prob [16]float64) [2][2][2][2]float64 {

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

	var a [2][2][2][2]float64

	a[Pop_low][InfoType_F][Interest_low][Assum_F] = prob_0000
	a[Pop_low][InfoType_F][Interest_low][Assum_T] = prob_0001
	a[Pop_low][InfoType_F][Interest_high][Assum_F] = prob_0010
	a[Pop_low][InfoType_F][Interest_high][Assum_T] = prob_0011
	a[Pop_low][InfoType_T][Interest_low][Assum_F] = prob_0100
	a[Pop_low][InfoType_T][Interest_low][Assum_T] = prob_0101
	a[Pop_low][InfoType_T][Interest_high][Assum_F] = prob_0110
	a[Pop_low][InfoType_T][Interest_high][Assum_T] = prob_0111
	a[Pop_high][InfoType_F][Interest_low][Assum_F] = prob_1000
	a[Pop_high][InfoType_F][Interest_low][Assum_T] = prob_1001
	a[Pop_high][InfoType_F][Interest_high][Assum_F] = prob_1010
	a[Pop_high][InfoType_F][Interest_high][Assum_T] = prob_1011
	a[Pop_high][InfoType_T][Interest_low][Assum_F] = prob_1100
	a[Pop_high][InfoType_T][Interest_low][Assum_T] = prob_1101
	a[Pop_high][InfoType_T][Interest_high][Assum_F] = prob_1110
	a[Pop_high][InfoType_T][Interest_high][Assum_T] = prob_1111

	return a
}
