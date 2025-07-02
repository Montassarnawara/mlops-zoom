package main

import (
	"fmt"
)

// Rename or remove this main function to avoid redeclaration error.
// For example, rename to mainPointers for demonstration:

func main() {
	var a = 1
	var b = a // b is a copy of a

	fmt.Println("a:", a, "b:", b) // a: 1 b: 1

	b += 1                        // b is now 2, a is still 1
	fmt.Println("a:", a, "b:", b) // a: 1 b:

	var c = 1
	var d *int                     // d is a pointer to an int, it is nil by default
	d = &c                         // d now points to c
	fmt.Println("c:", c, "d:", *d) // c: 1 d: 1
	*d += 1                        // dereference d to change c
	fmt.Println("c:", c, "d:", *d) // c: 2 d: 2
	//
	var x = &a // x now points to a
	// dereference x to change a
	fmt.Println("a:", a, "b:", b, "x:", x, "c:", c, "d:", *d) // a: 1 b: 2 c: 2 d: 2
}
