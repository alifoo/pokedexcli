package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	} {
		{
			input: "   hello   world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "p ikachu charmander       venosaur",
			expected: []string{"p", "ikachu", "charmander", "venosaur"},
		},
		{
			input: "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Unmatching length between expected output and cleanInput function")
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Output different than the actual expected string list in word: %v", word)
			}
		}
	}

}