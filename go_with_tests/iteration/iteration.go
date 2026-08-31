package iteration

import "strings"

const numberOfIterations = 5

func Iterate(s string) string {
	res := ""
	for i := 0; i < numberOfIterations; i++ {
		res += s
	}
	return res
}

func IterateWithStringPackage(s string) string {
	var res strings.Builder

	for i := 0; i < numberOfIterations; i++ {
		res.WriteString(s)
	}

	return res.String()

}
