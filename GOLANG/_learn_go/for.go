package main

import (
	"fmt"
)

func main() {
	// For loop
	for i := 0; i < 5; i++ {
		fmt.Println("For loop iteration:", i)
	}

	// While-like loop using for
	j := 0
	for j < 5 {
		fmt.Println("While-like loop iteration:", j)
		j++
	}
}
