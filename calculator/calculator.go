package main

import (
	"errors"
	"fmt"
	"strconv"
)

func runOperation(a, b, c int) (int, error) {
	switch c {
	case 1:
		return a + b, nil
	case 2:
		return a - b, nil
	case 3:
		return a * b, nil
	case 4:
		if b == 0 {
			return 0, errors.New("Divide by zero error.")
		}
		return a / b, nil
	default:
		return 0, errors.New("Invalid operation. Please try again.")
	}
}

func getOperation(op int) string {
	switch op {
	case 1:
		return "+"
	case 2:
		return "-"
	case 3:
		return "*"
	case 4:
		return "/"
	}
	return ""
}

func Calculator() {
	fmt.Println("Operator")
	var a string
	_, err := fmt.Scan(&a)
	Operator, err := strconv.Atoi(a)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Operand")
	var b string
	_, err = fmt.Scan(&b)
	Operand, err := strconv.Atoi(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Operation")
	fmt.Println("1. Addition")
	fmt.Println("2. Subtraction")
	fmt.Println("3. Multiplication")
	fmt.Println("4. Division")
	var c string
	_, err = fmt.Scan(&c)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	Operation, err := strconv.Atoi(c)
	res, err := runOperation(Operator, Operand, Operation)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%d %s %d = %d\n", Operand, getOperation(Operation), Operator, res)

}

func main() {
	Calculator()

}
