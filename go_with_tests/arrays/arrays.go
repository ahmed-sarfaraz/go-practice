package arrays

func VariadicAdd(val ...[]int) []int {
	res := make([]int, len(val))
	for k, i := range val {
		sum := 0
		for _, j := range i {
			sum += j
		}
		res[k] = sum
	}
	return res
}
