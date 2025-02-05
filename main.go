package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		exec := false
		fmt.Print("Pokedex > ")
		scanner.Scan()
		var command string
		if len(scanner.Text()) > 0 {
			txt := strings.ToLower(scanner.Text())
			command = strings.Fields(txt)[0]
		} else {
			command = "Unknown"
		}

		for key, value := range commands {
			if command == key {
				exec = true
				value.callback()
			}
		}

		if !exec {
			fmt.Println("Unknown command")
		}

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

func commandExit() {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
}

func commandHelp() {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n")
	// for _, value := range commands {
	// 	fmt.Printf("%s: %s\n", value.name, value.description)
	// }
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
}

type cliCommand struct {
	name string
	description string
	callback func()
}

var commands = map[string]cliCommand{
	"exit": {
		name: "exit",
		description: "Exit the Pokedex",
		callback: commandExit,
	},
	"help": {
		name: "help",
		description: "Displays a help message",
		callback: commandHelp,
	},
}