package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	startRepl(commands)
}

var commands map[string]cliCommand

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

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n")
	for name, cmd := range commands {
		fmt.Printf("%s: %s\n", name, cmd.description)
	}
	return nil
}

func commandMap(config *Config) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if config.Next != "" {
		url = config.Next
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
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

func commandMapb(config *Config) error {
	if config.Previous == "" {
		if config.Next == "" {
			fmt.Println("Showing the first page of locations:")
			return commandMap(config)
		}
		fmt.Println("you're on the first page")
		return nil
	}
	resp, err := http.Get(config.Previous)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
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

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
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
	}
}

func startRepl(commands map[string]cliCommand) {
	scanner := bufio.NewScanner(os.Stdin)
	config := &Config{}
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		command := parseFirstWord(input)
		if cmd, exists := commands[command]; exists {
			err := cmd.callback(config)
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
