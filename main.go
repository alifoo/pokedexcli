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


var commands map[string]cliCommand
var config Config
var caughtPokemons map[string]Pokemon

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

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

func main() {
	commands = map[string]cliCommand{
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
			description: "Explore a specific area, seeing all pokemons found there",
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
		"pokedex": {
			name: "pokedex",
			description: "Inspect a all caught pokemon",
			callback: func(config *Config, s string) {
				commandPokedex(config, s)
			},
		},
	}


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
	fmt.Println(Red + "Closing the Pokedex... Goodbye!" + Reset)
	os.Exit(0)
}

func commandHelp(config *Config, areaName string) {
	fmt.Print(Green + "Welcome to the Pokedex!\nUsage:\n" + Reset)
	for _, value := range commands {
		fmt.Printf(Blue + "%s: %s\n" + Reset, value.name, value.description)
	}
}

func commandMap(config *Config, areaName string) {
	url := "https://pokeapi.co/api/v2/location-area/"
	d, exists := config.c.Get(url)
	if exists {
		var locationResponse Location
		err := json.Unmarshal(d, &locationResponse)
		if err != nil {
			fmt.Println(err)
		}

		for _, loc := range locationResponse.Results {
			fmt.Println(loc.Name)
		}
	} else {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()

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
		fmt.Println("You're on the first page")
	} else {
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
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Error in creating request: %v", err)
			return
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("Error in making request: %v", err)
			return
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			fmt.Printf("Error: Area '%s' not found\n", areaName)
			return
		}
		data, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error in reading response body: %v", err)
			return
		}

		var locationAreaResponse LocationArea
		err = json.Unmarshal(data, &locationAreaResponse)
		if err != nil {
			fmt.Printf("Error in Unmarshal of data: %v", err)
			return
		}

		for _, poke := range locationAreaResponse.PokemonEncounters {
			fmt.Println(Purple + poke.Pokemon.Name + Reset)
		}
	} else {
		var locationAreaResponse LocationArea
		err := json.Unmarshal(d, &locationAreaResponse)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, poke := range locationAreaResponse.PokemonEncounters {
			fmt.Println(Purple + poke.Pokemon.Name + Reset)
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
		fmt.Printf(Green + "Congratulations! You catched %s!\n" + Reset, pokemonName)
		fmt.Println("You may now inspect it with the inspect command.")
		caughtPokemons[pokemonName] = pokemon
	} else {
		fmt.Printf(Yellow + "%v escaped the PokeBall! Try throwing another PokeBall.\n" + Reset, pokemonName)
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

func commandPokedex(config *Config, s string) {
	if len(caughtPokemons) == 0 {
		fmt.Println(Yellow + "You still have 0 pokemons. Try using the command 'catch' and a pokemon name!" + Reset)
		return
	}

	fmt.Println(Green + "Your pokedex:" + Reset)
	for _, p := range caughtPokemons {
		fmt.Printf(Cyan + "  - %v\n" + Reset, p.Name)
	}
}

