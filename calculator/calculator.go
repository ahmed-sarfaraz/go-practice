package main

import (
	"fmt"
)

func runCalculator(a, b, c int) int {
	switch c {
	case 1:
		return a + b
	case 2:
		return a - b
	case 3:
		return a * b
	case 4:
		if b == 0 {
			fmt.Println("Divide by zero error.")
			return 0
		}
		return a / b
	default:
		fmt.Println("Invalid operation. Please try again.")
		return 0
	}

}

func main() {
	fmt.Println("Operator")
	var a int
	_, err := fmt.Scan(&a)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Operand")
	var b int
	_, err = fmt.Scan(&b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Operation")
	fmt.Println("1. Addition")
	fmt.Println("2. Subtraction")
	fmt.Println("3. Multiplication")
	fmt.Println("4. Division")
	var c int
	_, err = fmt.Scan(&c)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Print(runCalculator(a, b, c))
}
