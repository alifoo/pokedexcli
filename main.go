package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
}

func cleanInput(text string) []string {
	if len(text) == 0 {
		var input []string
		return input
	}
	input := strings.Fields(text)
	return input
}