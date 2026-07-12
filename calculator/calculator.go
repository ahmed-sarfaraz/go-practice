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

func readInt(prompt string) (int, error) {
	fmt.Print(prompt)

	var input string
	if _, err := fmt.Scan(&input); err != nil {
		return 0, err
	}

	// strconv.Atoi explicitly forces base-10, bypassing the octal/leading-zero trap
	num, err := strconv.Atoi(input)
	if err != nil {
		return 0, errors.New("Invalid number format")
	}

	return num, nil
}

func Calculator() {
	num1, err := readInt("Enter first number: ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	num2, err := readInt("Enter second number: ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("\nChoose an operation:")
	fmt.Println("1. Addition\n2. Subtraction\n3. Multiplication\n4. Division")
	op, err := readInt("Your choice (1-4): ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	res, err := runOperation(num1, num2, op)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Fixed the variable order here so the math matches the display
	fmt.Printf("%d %s %d = %d\n", num1, getOperation(op), num2, res)

}

func main() {
	Calculator()

}
