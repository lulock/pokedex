package main

import "testing"

func TestCleanInput(t *testing.T) {
	// start by creating a slice of test case structs
	cases := []struct{
		input string
		expected []string
	}{
		{
			input: "   hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "    ",
			expected: []string{},
		},
		{
			input: "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		if len(actual) != len(c.expected){
			t.Errorf("Expected: %v, but got %v.", c.expected, actual)
			return
		}
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("Expected: %v, but got: %v.", expectedWord, word )
				return
			}
		}
	}
}
