package main

import (
	"fmt"
)

func add(a int, b int) int {
	return a + b
}
func addAndSubtract(a int, b int) (int, int) {
	return a + b, a - b
}
func multiplyAndDivide(a int, b int) (mul int, div float64) {
	mul = a * b
	if b != 0 {
		div = float64(a) / float64(b)
	} else {
		div = 0 // Handle division by zero
	}
	return
}
func main() {
	// Function call
	result := add(10, 20)
	fmt.Println("Sum:", result)

	// Function with multiple return values
	sum, diff := addAndSubtract(30, 10)
	fmt.Println("Sum:", sum, "Difference:", diff)

	// Function with named return values
	mul, div := multiplyAndDivide(20, 4)
	fmt.Println("Multiplication:", mul, "Division:", div)
}
