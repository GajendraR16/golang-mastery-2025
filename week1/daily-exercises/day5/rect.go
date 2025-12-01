package main

import (
	"errors"
	"fmt"
)

type Rectangle struct {
	Length, Width float64
}

func NewRectangle(length, width float64) (Rectangle, error) {
	if length >= 0 && width >= 0 {
		return Rectangle{Length: length, Width: width}, nil
	}
	return Rectangle{}, errors.New("length and width must be positive")
}

func (r Rectangle) IsSquare() bool {
	return r.Length == r.Width
}

func (r *Rectangle) Scale(factor float64) {
	r.Length *= factor
	r.Width *= factor
}

func (r Rectangle) String() string {
	// TODO: Return formatted string like "Rectangle(10x5)"
	return fmt.Sprintf("Rectangle(%.2fx%.2f)", r.Length, r.Width)
}
func (r Rectangle) Area() float64 {
	return r.Length * r.Width
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Length + r.Width)
}

func main() {
	// Test NewRectangle with validation
	rect, err := NewRectangle(10, 5)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(rect)
	fmt.Printf("Area: %.2f\n", rect.Area())
	fmt.Printf("Perimeter: %.2f\n", rect.Perimeter())
	fmt.Printf("Is square? %v\n", rect.IsSquare())

	// Test Scale
	fmt.Println("\nScaling by 2...")
	rect.Scale(2)
	fmt.Println(rect)
	fmt.Printf("Area: %.2f\n", rect.Area())

	// Test square
	fmt.Println("\nCreating a square...")
	square, _ := NewRectangle(7, 7)
	fmt.Println(square)
	fmt.Printf("Is square? %v\n", square.IsSquare())

	// Test invalid input
	fmt.Println("\nTrying invalid dimensions...")
	invalid, err := NewRectangle(-5, 10)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(invalid)
	}
}
