package funcs

func Set_Has(setA []int, num int)bool{
  for _, a := range setA{
    if a == num{
      return true
    }
  }
  return false
}

func Set_Sum(setA []int, setB []int)[]int{
  ans := make([]int,len(setA), len(setA)+len(setB))
  _ = copy(ans,setA)

  for _, b := range setB{
    if Set_Has(ans,b) == false{
      ans = append(ans,b)
    }
  }
  return ans
}
