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

	"github.com/breeze/pokedexcli/internal/pokecache"
)

func main() {
	startRepl(commands)
}

var commands map[string]cliCommand
var appCache *pokecache.Cache

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	trimmed := strings.TrimSpace(lowered)
	final := strings.Fields(trimmed)
	return final
}

func parseFirstWord(text string) string {
	cleaned := cleanInput(text)
	if len(cleaned) == 0 {
		return ""
	}
	return cleaned[0]
}

func commandExit(config *Config, _ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, _ []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for name, cmd := range commands {
		fmt.Printf("%s: %s\n", name, cmd.description)
	}
	return nil
}

func commandMap(config *Config, _ []string) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if config.Next != "" {
		url = config.Next
	}
	cachedData, found := appCache.Get(url)
	var body []byte
	var err error
	if found {
		fmt.Println("Using cached data!")
		body = cachedData
	} else {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		appCache.Add(url, body)
	}
	var locationResp LocationAreaResponse
	err = json.Unmarshal(body, &locationResp)
	if err != nil {
		return err
	}
	for _, area := range locationResp.Results {
		fmt.Println(area.Name)
	}
	if locationResp.Next != nil {
		config.Next = *locationResp.Next
	} else {
		config.Next = ""
	}
	if locationResp.Previous != nil {
		config.Previous = *locationResp.Previous
	} else {
		config.Previous = ""
	}
	return nil
}

func commandMapb(config *Config, _ []string) error {
	if config.Previous == "" {
		if config.Next == "" {
			fmt.Println("Showing the first page of locations:")
			return commandMap(config, []string{})
		}
		fmt.Println("you're on the first page")
		return nil
	}
	url := config.Previous
	cachedData, found := appCache.Get(url)
	var body []byte
	var err error
	if found {
		fmt.Println("Using cached data!")
		body = cachedData
	} else {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		appCache.Add(url, body)
	}
	var locationResp LocationAreaResponse
	err = json.Unmarshal(body, &locationResp)
	if err != nil {
		return err
	}
	for _, area := range locationResp.Results {
		fmt.Println(area.Name)
	}
	if locationResp.Next != nil {
		config.Next = *locationResp.Next
	} else {
		config.Next = ""
	}
	if locationResp.Previous != nil {
		config.Previous = *locationResp.Previous
	} else {
		config.Previous = ""
	}
	return nil
}

func commandExplore(config *Config, params []string) error {
	if len(params) == 0 {
		return fmt.Errorf("please provide a location area name")
	}
	locationAreaName := strings.Join(params, "-")
	fmt.Printf("Exploring %s...\n", locationAreaName)
	url := "https://pokeapi.co/api/v2/location-area/" + locationAreaName
	var data map[string]interface{}
	cachedData, found := appCache.Get(url)
	var body []byte
	if found {
		body = cachedData
	} else {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		appCache.Add(url, body)
	}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	fmt.Println("Found Pokemon:")
	pokemonEncounters, ok := data["pokemon_encounters"].([]interface{})
	if !ok {
		return fmt.Errorf("error parsing pokemon encounters")
	}
	for _, encounter := range pokemonEncounters {
		encounterMap, ok := encounter.(map[string]interface{})
		if !ok {
			continue
		}
		pokemon, ok := encounterMap["pokemon"].(map[string]interface{})
		if !ok {
			continue
		}
		name, ok := pokemon["name"].(string)
		if !ok {
			continue
		}
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, []string) error
}

func init() {
	commands = map[string]cliCommand{
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
			description: "Displays Next 20 Locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays Previous 20 Locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays Pokemon in an area",
			callback:    commandExplore,
		},
	}
}

func startRepl(commands map[string]cliCommand) {
	scanner := bufio.NewScanner(os.Stdin)
	config := &Config{}
	appCache = pokecache.NewCache(time.Minute * 5)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		parts := strings.SplitN(input, " ", 2)
		command := parts[0]
		var params []string
		if len(parts) > 1 {
			params = []string{parts[1]}
		}
		if cmd, exists := commands[command]; exists {
			err := cmd.callback(config, params)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

type LocationAreaResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Config struct {
	Next     string
	Previous string
}
