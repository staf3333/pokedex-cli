package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/staf3333/pokedexcli/internal/pokeapi"
	"github.com/staf3333/pokedexcli/internal/pokecache"
)

// main will be the thing that actually runs the command

// hint below
// the below block of code creates a new struct called cliCommand which
// we will we will map to string value that we read in from the buffer
type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

// the below is an example structure of a map that maps strings to cliCommands
// the callbacks are the functions that I will want to create so that if a
// command matches that key, do whatever that function says to do
// e.g. for "help" we might want the callback "commandHelp" to print a guide on
// how to use the pokedex
func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display names of 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name: "explore <area_name>",
			description: "Display pokemon in given area",
			callback: commandExplore,
		},
		"catch": {
			name: "catch <pokemon_name>",
			description: "Capture a pokemon",
			callback: commandCatch,
		},
		"inspect": {
			name: "inspect <pokemon_name>",
			description: "Inspect a pokemon in pokedex",
			callback: commandInspect,
		},
	}
}

// need some way to keep track of what is the current page, next page, and prev page
// to do this, define struct to hold the next and prev page. Then pass refs (pointer) to
// these structs when you call the respective commands
// Commands need to accept a pointer to a config struct as a param!
// config struct contains the next and previous urls
type config struct {
	previous *string
	next     string
	cache pokecache.Cache
	pokedex map[string]pokeapi.PokeAPIPokemonResponse
}

// displays the names of 20 location areas in the Pokemon world
// each subsequent call to map should display the next 20 locations
func commandMap(config *config, args ...string) error {
	nextURL := config.next
	locationResponse := pokeapi.PokeAPILocationResponse{}
	body := pokeapi.GetData(nextURL, &config.cache)
	err := json.Unmarshal(body, &locationResponse)
	if err != nil {			
		fmt.Println("error with API request")
		return err
	}
	for _, location := range locationResponse.Results {
		fmt.Println(location.Name)
	}
	previous, next := locationResponse.Previous, locationResponse.Next
	config.previous = previous
	config.next = next
	return nil
}

// similar to map command, displays the previous 20 locations
// suggests, need a way to keep track of the page that you're currently on
func commandMapb(config *config, args ...string) error {
	if config.previous == nil {
		return errors.New("you're on the first page")
	}

	// Derefence pointer to get the string value
	previousURL := *config.previous
	locationResponse := pokeapi.PokeAPILocationResponse{}
	body := pokeapi.GetData(previousURL, &config.cache)
	err := json.Unmarshal(body, &locationResponse)
	if err != nil {			
		fmt.Println("Invalid URL")
		return err
	}
	for _, location := range locationResponse.Results {
		fmt.Println(location.Name)
	}
	previous, next := locationResponse.Previous, locationResponse.Next
	config.previous = previous
	config.next = next
	return nil
}

func commandExplore(config *config, args ...string) error {
	if len(args) < 1 {
		fmt.Println("You need to enter a location area")
		return errors.New("not enough arguments")
	}
	if len(args) > 1 {
		fmt.Println("Only one location may be accepted")
		return errors.New("too many arguments")
	}
	areaName := args[0]
	fmt.Printf("Exploring %v \n", areaName)
	locationAreaResponse := pokeapi.PokeAPILocationAreaResponse{}
	areaURL := "https://pokeapi.co/api/v2/location-area/" + areaName
	body := pokeapi.GetData(areaURL, &config.cache)
	err := json.Unmarshal(body, &locationAreaResponse)
	if err != nil {
		fmt.Println("Invalid LocationID")
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, encounter := range locationAreaResponse.PokemonEncounters {
		fmt.Printf("- %v \n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *config, args ...string) error {
	if len(args) < 1 {
		fmt.Println("You need to enter a pokemon to capture")
		return errors.New("not enough arguments")
	}
	if len(args) > 1 {
		fmt.Println("Can only capture one pokemon at a time")
		return errors.New("too many arguments")
	}
	pokemonName := strings.ToLower(args[0])
	pokemonUrl := "https://pokeapi.co/api/v2/pokemon/" + pokemonName
	pokemonResponse := pokeapi.PokeAPIPokemonResponse{}
	body := pokeapi.GetData(pokemonUrl, &config.cache)
	err := json.Unmarshal(body, &pokemonResponse)
	if err != nil {
		fmt.Println("Invalid Pokemon Name")
		return err
	}
	fmt.Printf("Throwing a pokeball at %s... \n", pokemonName)
	catchChance := pokemonResponse.BaseExperience
	catchRoll := rand.Intn(620)
	if catchRoll > catchChance {
		fmt.Printf("%s was caught! \n", pokemonName)
		// add pokemon to pokedex
		config.pokedex[pokemonName] = pokemonResponse
	} else {
		fmt.Printf("%s escaped! \n", pokemonName)
	}
	
	return nil
}

func commandInspect(config *config, args ...string) error {
	if len(args) < 1 {
		fmt.Println("You need to choose a pokemon to inspect")
		return errors.New("not enough arguments")
	}
	if len(args) > 1 {
		fmt.Println("You can only inspect one pokemon at a time")
		return errors.New("too many arguments")
	}

	pokemonName := strings.ToLower(args[0])
	pokemon, ok := config.pokedex[pokemonName]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return errors.New("pokemon not registered")
	}
	// print the name, height, weight, stats and type(s) of the Pokemon
	fmt.Printf("Name: %s \n",pokemon.Name)
	fmt.Printf("Height: %v \n", pokemon.Height)
	fmt.Printf("Weight: %v \n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, statList := range pokemon.Stats {
		fmt.Printf(" -%s: %v \n", statList.Stat.Name, statList.BaseStat)
	}

	fmt.Println("Types:")
	for _, typeList := range pokemon.Types {
		fmt.Printf(" -%s \n", typeList.Type.Name)
	}

	return nil
}

func commandHelp(config *config, args ...string) error {
	// do what criteria says when help command is called
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	fmt.Println()
	return nil
}

func commandExit(config *config, args ...string) error {
	fmt.Println("Exiting Pokedex")
	os.Exit(0)
	// if no errors, return nil
	return nil
}

func parseInput(input string) (string, []string) {
	parts := strings.Split(input, " ")
	if len(parts) == 2 {
		return parts[0], parts[1:]
	}
	return parts[0], nil
}

func main() {
	config := config{
		next:     "https://pokeapi.co/api/v2/location/?limit=20",
		previous: nil,
		cache: *pokecache.NewCache(100 * time.Second),
		pokedex: map[string]pokeapi.PokeAPIPokemonResponse{},
	}
	for {
		// Create new scanner to read from stdin
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("Pokedex > ")

		// read input line by line
		for scanner.Scan() {
			input := scanner.Text()
			// destructure command name and params from input
			commandName, args := parseInput(input)
			command, exists := getCommands()[commandName]
			if exists {
				err := command.callback(&config, args...)
				if err != nil {
					//handle error some type of way
					fmt.Println("Error: ", err)
				}
				break
			} else {
				fmt.Println("Command does not exist")
			}
		}
	}
}
