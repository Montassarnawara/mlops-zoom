package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
	var ch = "hello bro"
	fmt.Println(ch)

	ch = "hello bro 2"
	fmt.Println(ch)
	// Declare and initialize variables
	var a int = 10
	var b int = 20
	println("a:", a)
	println("b:", b)
	var c int = a + b
	fmt.Printf("%d + %d = %d", a, b, c)

	// Declare bool  variables
	//var test bool = true
	//fmt.Println("\nIs it true?", isTrue)
	// declare en line
	var a1, b1, c1 int = 10, 20, 30
	fmt.Println("\nValues of a, b, c:", a1, b1, c1)
}
