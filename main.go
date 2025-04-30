package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"net/http"
	"encoding/json"
	"io"
)

// this is a registry of commands
// it provides abstraciton for managing commands added
type cliCommand struct {
	name string
	description string
	callback func(*config) error
}

type config struct {
	Next string
	Previous string
}

type LocationAreas struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}


// cleans input by removing whitespace and returning slice of all words in lowercase
func cleanInput(text string) []string {
	text_trimmed := strings.TrimSpace(text)
	text_lowercase := strings.ToLower(text_trimmed)
	result := strings.Fields(text_lowercase)
	return result
}

// exits the programme
func commandExit(conf *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// displays the names of 20 location areas in the Pokemon world

func commandMap(conf *config) error {
	resp, err := http.Get(conf.Next)
	if err != nil {
		fmt.Println(fmt.Sprintf("Could not get locations:%v",err))
		return err
	}

	locAreas := LocationAreas{}

	data, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = json.Unmarshal(data, &locAreas)
	
	if err != nil {
		fmt.Println(err)
		return err
	}
	
	for _, loc := range locAreas.Results {
		fmt.Println(loc.Name)
	}
	conf.Next = locAreas.Next
	conf.Previous = locAreas.Previous
	return nil
}

func commandMapb(conf *config) error {
	resp, err := http.Get(conf.Previous)
	if err != nil {
		fmt.Println(fmt.Sprintf("Could not get locations:%v",err))
		return err
	}

	locAreas := LocationAreas{}

	data, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = json.Unmarshal(data, &locAreas)
	
	if err != nil {
		fmt.Println(err)
		return err
	}
	
	for _, loc := range locAreas.Results {
		fmt.Println(loc.Name)
	}
	conf.Next = locAreas.Next
	conf.Previous = locAreas.Previous
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin) // wait for user input using bufio.NewScanner which blocks the code and waits for input, once the user types something and presses enter, the code continues and the input is available in the returned bufio.Scanner
	// map the supported commands:
	conf := config{Next: "https://pokeapi.co/api/v2/location-area/"}
	validCommands := map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"map": {
			name: "map",
			description: "Displays the names of the next 20 location areas in the Pokemon world",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "Displayes the names of the previous 20 locations in the Pokemon world",
			callback: commandMapb,
		},
	}
	validCommands["help"] = cliCommand{
			name: "help",
			description: "Displays a help message",
			callback: func(conf *config) error {
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
					if err := cmd.callback(&conf); err != nil {
					fmt.Println(err)
					}
				}
			}
		}
	}

}
