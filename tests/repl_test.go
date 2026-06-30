package tests

import (
	"testing"

	hel "github.com/mochatai/pokedex/helpFunc"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "a s d f g",
			expected: []string{"a", "s", "d", "f", "g"},
		},
		{
			input:    "HELLO",
			expected: []string{"hello"},
		},
	}

	for _, c := range cases {
		actual := hel.CleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Error("not having the same length")
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("not the same word")
			}
		}

	}

}
