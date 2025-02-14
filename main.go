package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alifoo/pokedexcli/internal"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config.c = internal.NewCache(5 * time.Second)
	caughtPokemons = make(map[string]Pokemon)
	for {
		exec := false
		fmt.Print("Pokedex > ")
		scanner.Scan()
		var command string
		var areaName string
		if len(scanner.Text()) > 0 {
			txt := strings.ToLower(scanner.Text())
			args := strings.Fields(txt)
			command = args[0]
			if len(args) > 1 {
				areaName = args[1]
			}
		} else {
			command = "Unknown"
		}

		for key, value := range commands {
			if command == key {
				exec = true
				value.callback(&config, areaName)
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

func commandExit(config *Config, areaName string) {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
}

func commandHelp(config *Config, areaName string) {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	// for _, value := range commands {
	// 	fmt.Printf("%s: %s\n", value.name, value.description)
	// }
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
}

func commandMap(config *Config, areaName string) {
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

func commandMapB(config *Config, areaName string) {
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
		defer res.Body.Close()
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

func commandExplore(config *Config, areaName string) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", areaName)
	d, exists := config.c.Get(url)

	if !exists {
		fmt.Println("Data not in cache! Requesting...")
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
		}

		res, err := http.DefaultClient.Do(req)
		defer res.Body.Close()
		if err != nil {
			fmt.Println(err)
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}

		var locationAreaResponse LocationArea
		err = json.Unmarshal(data, &locationAreaResponse)
		if err != nil {
			fmt.Println(err)
		}

		for _, poke := range locationAreaResponse.PokemonEncounters {
			fmt.Println(poke.Pokemon.Name)
		}
	} else {
		fmt.Println("Data still in cache! Using existing data.")
		var locationAreaResponse LocationArea
		err := json.Unmarshal(d, &locationAreaResponse)
		if err != nil {
			fmt.Println(err)
		}

		for _, poke := range locationAreaResponse.PokemonEncounters {
			fmt.Println(poke.Pokemon.Name)
		}
	}
}

func commandCatch(config *Config, pokemonName string) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokemonName)
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

	var pokemon Pokemon
	err = json.Unmarshal(data, &pokemon)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	chance := rand.IntN(pokemon.BaseExperience * 2)
	if chance >= pokemon.BaseExperience {
		fmt.Printf("Congratulations! You catched %s!\n", pokemonName)
		caughtPokemons[pokemonName] = pokemon
	} else {
		fmt.Printf("You were unable to catch it. Your score was %v. Try throwing another pokeball!\n", chance)
		fmt.Printf("%s base exp: %v\n", pokemonName, pokemon.BaseExperience)
	}
}

func commandInspect(config *Config, pokemonName string) {
	pokemonInfo := caughtPokemons[pokemonName]
	if len(pokemonInfo.Name) == 0 {
		fmt.Println("Pokemon not found!")
		return
	}

	fmt.Printf("Name: %s\n", pokemonInfo.Name)
	fmt.Printf("Height: %v\n", pokemonInfo.Height)
	fmt.Printf("Weight: %v\n", pokemonInfo.Weight)
	fmt.Println("Stats:")
	for _, s := range pokemonInfo.Stats {
		fmt.Printf("  -%v: %v\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemonInfo.Types {
		fmt.Printf("  - %v\n", t.Type.Name)
	}
}

type cliCommand struct {
	name string
	description string
	callback func(*Config, string)
}

type Config struct {
	Next *string
	Previous *string
	c *internal.Cache
}

var config Config
var caughtPokemons map[string]Pokemon

var commands = map[string]cliCommand{
	"exit": {
		name: "exit",
		description: "Exit the Pokedex",
		callback: func(config *Config, areaName string) { commandExit(config, areaName) },
	},
	"help": {
		name: "help",
		description: "Displays a help message",
		callback: func(config *Config, areaName string) { commandHelp(config, areaName ) },
	},
	"map": {
		name: "map",
		description: "Open the map",
		callback: func(config *Config, areaName string) { commandMap(config, areaName) },
	},
	"mapb": {
		name: "mapb",
		description: "Goes back in the map",
		callback: func(config *Config, areaName string)  {
			commandMapB(config, areaName)
		},
	},
	"explore": {
		name: "explore",
		description: "Explore a specific area",
		callback: func(config *Config, areaName string) {
			commandExplore(config, areaName)
		},
	},
	"catch": {
		name: "catch",
		description: "Catch a pokemon",
		callback: func(config *Config, pokemonName string) {
			commandCatch(config, pokemonName)
		},
	},
	"inspect": {
		name: "inspect",
		description: "Inspect a caught pokemon",
		callback: func(config *Config, pokemonName string) {
			commandInspect(config, pokemonName)
		},
	},
}