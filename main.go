package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

// this is a registry of commands
// it provides abstraciton for managing commands added
type cliCommand struct {
	name string
	description string
	callback func() error
}


// cleans input by removing whitespace and returning slice of all words in lowercase
func cleanInput(text string) []string {
	text_trimmed := strings.TrimSpace(text)
	text_lowercase := strings.ToLower(text_trimmed)
	result := strings.Fields(text_lowercase)
	return result
}

// exits the programme
func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin) // wait for user input using bufio.NewScanner which blocks the code and waits for input, once the user types something and presses enter, the code continues and the input is available in the returned bufio.Scanner
	// map the supported commands:
	validCommands := map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
	}
	validCommands["help"] = cliCommand{
			name: "help",
			description: "Displays a help message",
			callback: func() error {
				fmt.Println("Welcome to the Pokedex!")
				fmt.Println("Usage:")
				fmt.Println()
				for _, v := range (validCommands) {
					fmt.Println(fmt.Sprintf("%v: %v", v.name, v.description))
				}
				return nil
			},
		}
	for i := 0; ; i++ {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			usrinput := scanner.Text()
			cleanedUsrinput := cleanInput(usrinput)
			if len(cleanedUsrinput) > 0 {
				command := cleanedUsrinput[0]
				cmd, ok := validCommands[command]
				if !ok {
					fmt.Println("Unknown command")
				} else {
					if err := cmd.callback(); err != nil {
					fmt.Println(err)
					}
				}
			}
		}
	}

}
