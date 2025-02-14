package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alifoo/pokedexcli/internal"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config.c = internal.NewCache(5 * time.Second)
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
				value.callback(&config)
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

func commandExit(config *Config) {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
}

func commandHelp(config *Config) {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	// for _, value := range commands {
	// 	fmt.Printf("%s: %s\n", value.name, value.description)
	// }
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
}

func commandMap(config *Config) {
	url := "https://pokeapi.co/api/v2/location-area/"
	d, exists := config.c.Get(url)
	if exists {
		fmt.Println("Data still in cache! Using existing data.")
		var locationResponse Location
		err := json.Unmarshal(d, &locationResponse)
		if err != nil {
			fmt.Println(err)
		}

		for _, loc := range locationResponse.Results {
			fmt.Println(loc.Name)
		}
	} else {
		fmt.Println("Data not in cache! Requesting...")
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()
		code := res.StatusCode
		fmt.Printf("Code: %v\n", code)


		data, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}

		config.c.Add(url, data)
		var locationResponse Location
		err = json.Unmarshal(data, &locationResponse)
		if err != nil {
			fmt.Println(err)
		}

		for _, loc := range locationResponse.Results {
			fmt.Println(loc.Name)
		}
	}
}

func commandMapB(config *Config) {
	if config.Previous == nil {
		fmt.Println("you're on the first page")
	} else {
		fmt.Println("cool")
		url := *config.Previous
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}

		var locationResponse Location
		err = json.Unmarshal(data, &locationResponse)
		if err != nil {
			fmt.Println(err)
		}

		for _, loc := range locationResponse.Results {
			fmt.Println(loc.Name)
		}

	}

}

type cliCommand struct {
	name string
	description string
	callback func(*Config)
}

type Config struct {
	Next *string
	Previous *string
	c *internal.Cache
}

var config Config

var commands = map[string]cliCommand{
	"exit": {
		name: "exit",
		description: "Exit the Pokedex",
		callback: func(config *Config) { commandExit(config) },
	},
	"help": {
		name: "help",
		description: "Displays a help message",
		callback: func(config *Config) { commandHelp(config) },
	},
	"map": {
		name: "map",
		description: "Open the map",
		callback: func(config *Config) { commandMap(config) },
	},
	"mapb": {
		name: "mapb",
		description: "Goes back in the map",
		callback: func(config *Config)  {
			commandMapB(config)
		},
	},
}