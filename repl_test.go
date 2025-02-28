package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "aLL your Base    are   belong TO  us",
			expected: []string{"all", "your", "base", "are", "belong", "to", "us"},
		},
		{
			input:    "Charmander Bulbasaur Squirtle",
			expected: []string{"charmander", "bulbasaur", "squirtle"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q) returned %d words, expected %d", c.input, len(actual), len(c.expected))
			continue
		}
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("cleanInput(%q) word %d: got %q, expected %q", c.input, i, actual[i], c.expected[i])
			}
		}
	}
}
