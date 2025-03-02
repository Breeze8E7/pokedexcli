package main

import (
	"bufio"
	"fmt"
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

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n")
	for name, cmd := range commands {
		fmt.Printf("%s: %s\n", name, cmd.description)
	}
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
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
	}
}

func startRepl(commands map[string]cliCommand) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		command := parseFirstWord(input)
		if cmd, exists := commands[command]; exists {
			err := cmd.callback()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
