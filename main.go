package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"net/http"
	"encoding/json"
	"io"
	"github.com/lulock/pokedex/internal/pokecache"
	"time"
	"math/rand"
)

// this is a registry of commands
// it provides abstraciton for managing commands added
type cliCommand struct {
	name string
	description string
	callback func(*config, ...string) error
}

type config struct {
	Next string
	Previous string
	Cache *pokecache.Cache
	Pokedex map[string]pokemon
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

type pokemonInArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}
type pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order        any `json:"order"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []struct {
		Abilities []struct {
			Ability  any  `json:"ability"`
			IsHidden bool `json:"is_hidden"`
			Slot     int  `json:"slot"`
		} `json:"abilities"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"past_abilities"`
	PastTypes []any `json:"past_types"`
	Species   struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       string `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  string `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      string `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale string `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       string `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       string `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  string `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      string `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale string `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  string `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

// cleans input by removing whitespace and returning slice of all words in lowercase
func cleanInput(text string) []string {
	text_trimmed := strings.TrimSpace(text)
	text_lowercase := strings.ToLower(text_trimmed)
	result := strings.Fields(text_lowercase)
	return result
}

// exits the programme
func commandExit(conf *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// displays the names of 20 location areas in the Pokemon world
func commandMap(conf *config, args ...string) error {
	// first check cache?
	data, exists := conf.Cache.Get(conf.Next)
	locAreas := LocationAreas{}
	// if it does not exist,get
	if !exists {
		resp, err := http.Get(conf.Next)
		if err != nil {
			fmt.Println(fmt.Sprintf("Could not get locations: %v", err))
			return err
		}
		data, err = io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println(err)
			return err
		}
		// now add to cache
		conf.Cache.Add(conf.Next, data)

	} else {
		fmt.Println("using cache to go forwards!")
	}

	err := json.Unmarshal(data, &locAreas)
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

// displays the names of 20 previous locations 
func commandMapb(conf *config, args ...string) error {
	// first check cache?
	data, exists := conf.Cache.Get(conf.Previous)
	// if it does not exist,get
	if !exists {
		resp, err := http.Get(conf.Previous)
		if err != nil {
			fmt.Println(fmt.Sprintf("Could not get locations: %v", err))
			return err
		}
		data, err = io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println(err)
			return err
		}
		// now add to cache
		conf.Cache.Add(conf.Previous, data)
	} else {
		fmt.Println("using cache 2 go back")
	}


	locAreas := LocationAreas{}

	err := json.Unmarshal(data, &locAreas)
	
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

// explore command takes the name of a location area and lists 
// all the Pokemon located there.
func commandExplore(conf *config, args ...string) error {
	fmt.Println(fmt.Sprintf("Looking around %v for pokemon 🧐", args[0]))
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v", args[0])
	// first check cache?
	data, exists := conf.Cache.Get(url)
	// if it does not exist,get
	if !exists {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		data, err = io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println(err)
			return err
		}
		// now add to cache
		conf.Cache.Add(url, data)
	} else {
		fmt.Println("using cache 2 list pokemon")
	}

	pokemon := pokemonInArea{}
	err := json.Unmarshal(data, &pokemon)
	if err != nil {
		return err
	}
	fmt.Println("Found these fellas:")
	for _, poke := range pokemon.PokemonEncounters {
		fmt.Println(fmt.Sprintf(". %v", poke.Pokemon.Name))
	}
	return nil
}

// catch command takes the name of a pokemon and tries to catch them
func commandCatch(conf *config, args ...string) error {
	pokename := args[0]
	fmt.Println(fmt.Sprintf("Throwing a Pokeball at %v...", pokename))
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v", pokename)
	// first check cache?
	data, exists := conf.Cache.Get(url)
	// if it does not exist,get
	if !exists {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		data, err = io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println(err)
			return err
		}
		// now add to cache
		conf.Cache.Add(url, data)
	} else {
		fmt.Println("using cache 2 try catchin pokemon again")
	}

	pokemon := pokemon{}
	err := json.Unmarshal(data, &pokemon)
	if err != nil {
		return err
	}
	randomInt := rand.Intn(100)
	chance := (340 - pokemon.BaseExperience)/340
	if chance < 30 {
		chance = 30
	}
	if chance > 90 {
		chance = 90
	}
	
	isCaught := randomInt < chance
	
	if isCaught {
		fmt.Println(fmt.Sprintf("%v was caught!", pokemon.Name))
		conf.Pokedex[pokemon.Name] = pokemon
	} else {	
		fmt.Println(fmt.Sprintf("%v escaped!", pokemon.Name))
	}

	return nil
}

func commandInspect(conf *config, args ...string) error {
	pokename := args[0]
	pokemon, ok := conf.Pokedex[pokename]
	if !ok {
		fmt.Println(fmt.Sprintf("you have not caught that pokemon"))
	} else {
		fmt.Println(fmt.Sprintf("Name: %v", pokemon.Name))
		fmt.Println(fmt.Sprintf("Height: %v", pokemon.Height))
		fmt.Println(fmt.Sprintf("Weight: %v", pokemon.Weight))
		fmt.Println(fmt.Sprintf("Stats:"))
		fmt.Println(fmt.Sprintf("  . hp: %v", pokemon.Stats[0].BaseStat))
		fmt.Println(fmt.Sprintf("  . attack: %v", pokemon.Stats[1].BaseStat))
		fmt.Println(fmt.Sprintf("  . defense: %v", pokemon.Stats[2].BaseStat))
		fmt.Println(fmt.Sprintf("  . special-attack: %v", pokemon.Stats[3].BaseStat))
		fmt.Println(fmt.Sprintf("  . special-defense: %v", pokemon.Stats[4].BaseStat))
		fmt.Println(fmt.Sprintf("  . speed: %v", pokemon.Stats[5].BaseStat))
		fmt.Println(fmt.Sprintf("Types:"))
		for _, poketype := range pokemon.Types {
			fmt.Println(fmt.Sprintf("  . %v", poketype.Type.Name))
		}
	}
	
	return nil
}

func commandPokedex(conf *config, args ...string) error {
	if len(conf.Pokedex) == 0 {
		fmt.Println("You haven't caught any Pokemon yet! Use the Catch command and try to catch 'em all.")
	} else {
		fmt.Println("Your Pokedex:")
		for k, _ := range (conf.Pokedex) {
			fmt.Println(fmt.Sprintf(" . %v", k))
		}
	}
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin) // wait for user input using bufio.NewScanner which blocks the code and waits for input, once the user types something and presses enter, the code continues and the input is available in the returned bufio.Scanner
	// make a cache
	// const duration := 5 * time.Millisecond
	// cache := NewCache(duration)
	// map the supported commands:
	conf := config{
		Next: "https://pokeapi.co/api/v2/location-area/",
		Cache: pokecache.NewCache(5 * time.Second),
		Pokedex: make(map[string]pokemon),
	}

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
			description: "Displays the names of the previous 20 locations in the Pokemon world",
			callback: commandMapb,
		},
		"explore" : {
			name: "explore",
			description: "Displays a list of all Pokemon located in the area passed as input",
			callback: commandExplore,
		},
		"catch" : {
			name: "catch",
			description: "Tries to catch a Pokemon",
			callback: commandCatch,
		},
		"inspect" : {
			name: "inspect",
			description: "Inspects Pokemon",
			callback: commandInspect,
		},
		"pokedex" : {
			name: "pokedex",
			description: "Lists all caught Pokemon",
			callback: commandPokedex,
		},



	}
	validCommands["help"] = cliCommand{
			name: "help",
			description: "Displays a help message",
			callback: func(conf *config, args ...string) error {
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
				args := ""

				if len(cleanedUsrinput) > 1 {
					args = cleanedUsrinput[1]
				}
				cmd, ok := validCommands[command]
				if !ok {
					fmt.Println("Unknown command")
				} else {
					if err := cmd.callback(&conf, args); err != nil {
					fmt.Println(err)
					}
				}
			}
		}
	}

}
