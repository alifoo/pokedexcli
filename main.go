package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Print("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		txt := strings.ToLower(scanner.Text())
		command := strings.Fields(txt)[0]
		fmt.Println("Your command was:", command)
	}
}

func cleanInput(text string) []string {
	if len(text) == 0 {
		var input []string
		return input
	}
	input := strings.Fields(text)
	return input
}