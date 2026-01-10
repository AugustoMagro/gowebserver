package tests

import (
	"testing"

	"github.com/AugustoMagro/gowebserver/internal/auth"
)

func TestHashPassword(t *testing.T) {

	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "Password#",
			expected: "Password#",
		},
		{
			input:    "Augusto@v$magro!@",
			expected: "Augusto@v$magro!@",
		},
		{
			input:    "Bat#123@test",
			expected: "Bat#123@test",
		},
	}

	for _, c := range cases {
		actual, err := auth.HashPassword(c.input)
		if err != nil {
			t.Errorf("Error, test failed")
		}

		valid, err := auth.CheckPasswordHash(c.expected, actual)
		if err != nil {
			t.Errorf("Error, test failed")
		}

		if !valid {
			t.Errorf("Error, test failed")
		}
	}

}
