package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/staf3333/pokedexcli/internal/pokeapi"
)

// main will be the thing that actually runs the command

// hint below
// the below block of code creates a new struct called cliCommand which
// we will we will map to string value that we read in from the buffer
type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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
	}
}




// need some way to keep track of what is the current page, next page, and prev page
// to do this, define struct to hold the next and prev page. Then pass refs (pointer) to 
// these structs when you call the respective commands
// Commands need to accept a pointer to a config struct as a param!
// config struct contains the next and previous urls
type config struct {
	previous *string
	next string
}

// displays the names of 20 location areas in the Pokemon world
// each subsequent call to map should display the next 20 locations
func commandMap(config *config) error {
	nextURL := config.next
	previous, next := pokeapi.GetFromPokeAPI(nextURL)
	config.previous = previous
	config.next = next
	return nil
}

// similar to map command, displays the previous 20 locations
// suggests, need a way to keep track of the page that you're currently on
func commandMapb(config *config) error {
	if config.previous == nil {
		fmt.Println("No previous page available")
		return nil
	}
	
	// Derefence pointer to get the string value
	previousURL := *config.previous
	previous, next := pokeapi.GetFromPokeAPI(previousURL)
	config.previous = previous
	config.next = next
	return nil
}

func commandHelp(config *config) error {
	// do what criteria says when help command is called
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	fmt.Println()
	// getFromPokeAPI()
	// if no errors, return nil
	return nil
}

func commandExit(config *config) error {
	fmt.Println("Exiting Pokedex")
	os.Exit(0)
	// if no errors, return nil
	return nil
}

func main() {
	pageTracker := config{ 
		next: "https://pokeapi.co/api/v2/location/?limit=20",
		previous: nil,
	}
	for {
		// Create new scanner to read from stdin
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("Pokedex > ")

		// read input line by line
		for scanner.Scan() {
			input := scanner.Text()
			command, exists := getCommands()[input]
			if exists {
				err := command.callback(&pageTracker)
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
