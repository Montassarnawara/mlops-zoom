package main

import (
	"fmt"
)

// Define a struct to represent a person
type Person struct {
	Name string
	Age  int
}

func main() {
	// Create a slice of Person
	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}

	// Print the slice of people
	fmt.Println("People:", people)

	// Access and modify elements in the slice
	for i, person := range people {
		fmt.Printf("Person %d: Name=%s, Age=%d\n", i+1, person.Name, person.Age)
		if person.Age < 30 {
			person.Age += 5 // Increment age by 5 if less than 30
			fmt.Printf("Updated Person %d: Name=%s, Age=%d\n", i+1, person.Name, person.Age)
		}
	}

	// Create a map to store people by name
	peopleMap := make(map[string]Person)
	for _, person := range people {
		peopleMap[person.Name] = person
	}
	// Print the map of people
	fmt.Println("People Map:")
	for name, person := range peopleMap {
		fmt.Printf("Name: %s, Age: %d\n", name, person.Age)
	}

	// Create an array of integers
	numbers := [5]int{1, 2, 3, 4, 5}
	// Print the array of numbers
	fmt.Println("Numbers Array:", numbers)
	// Access and modify elements in the array
	for i, num := range numbers {
		fmt.Printf("Number %d: %d\n", i+1, num)
		if num%2 == 0 {
			numbers[i] *= 2 // Double the even numbers
			fmt.Printf("Updated Number %d: %d\n", i+1, numbers[i])
		}
	}
}
