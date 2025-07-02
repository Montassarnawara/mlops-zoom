package main

import (
	"fmt"
)

// Define a struct to represent a person
type user struct {
	Name string
	Age  int
}

func main() {
	// Create an instance of Person
	person1 := user{Name: "Alice", Age: 30}

	// Access fields of the struct
	fmt.Println("Name:", person1.Name)
	fmt.Println("Age:", person1.Age)

	// Modify fields of the struct
	person1.Age = 31
	fmt.Println("Updated Age:", person1.Age)

	// Create another instance of Person
	person2 := user{Name: "Bob", Age: 25}
	fmt.Println("Person 2:", person2)
}
