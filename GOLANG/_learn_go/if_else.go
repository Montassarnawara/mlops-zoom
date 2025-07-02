package main

import (
	"fmt"
)

func main() {
	// Declare and initialize variables
	var a int = 10
	var b int = 20
	fmt.Println("a:", a)
	fmt.Println("b:", b)

	// If-else statement
	if a > b {
		fmt.Println("a is greater than b")
	} else if a < b {
		fmt.Println("a is less than b")
	} else {
		fmt.Println("a is equal to b")
	}

	// Switch statement
	switch {
	case a > b:
		fmt.Println("In switch: a is greater than b")
	case a < b:
		fmt.Println("In switch: a is less than b")
	default:
		fmt.Println("In switch: a is equal to b")
	}
}
