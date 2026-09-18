package funcs

import "slices"

// func IsInSet(setA []int, num int) bool {
// 	for _, a := range setA {
// 		if a == num {
// 			return true
// 		}
// 	}
// 	return false
// }

func UnionSets(setA []int, setB []int) []int {
	ans := make([]int, len(setA), len(setA)+len(setB))
	_ = copy(ans, setA)

	for _, b := range setB {
		if slices.Contains(ans, b) == false {
			ans = append(ans, b)
		}
	}
	return ans
}
