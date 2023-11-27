package main

import (
	"fmt"
)

// main will be the thing that actually runs the command

// hint below
// the below block of code creates a new struct called cliCommand which 
// we will we will map to string value that we read in from the buffer
type cliCommand struct {
	name        string
	description string
	callback    func() error
}
// the below is an example structure of a map that maps strings to cliCommands
// the callbacks are the functions that I will want to create so that if a 
// command matches that key, do whatever that function says to do
// e.g. for "help" we might want the callback "commandHelp" to print a guide on 
// how to use the pokedex
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
}


func main() {
	for {
		// block the for loop to wait for input
	}
	// use infinite for loop to keep the REPL running
	// at start of loop, block it and wait for some input
	// when info is recieved, parse it and then execute a command
	// Once command is finished, print output then go to the next iteration of the loop                 
	fmt.Println("Hello, world!") 
}