package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		result := parseFirstWord(input)
		fmt.Println("Your command was:", result)
	}
}

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
