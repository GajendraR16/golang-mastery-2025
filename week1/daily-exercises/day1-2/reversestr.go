package main

import "fmt"

func main() {
	s := "🙂Go!"
	rev := []rune(s)
	n := len(rev) - 1
	i, j := 0, n
	for i < j {
		rev[i], rev[j] = rev[j], rev[i]
		i++
		j--
	}
	fmt.Printf("%s", string(rev))
}
