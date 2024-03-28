package commands

import (
	"fmt"

	"github.com/prethora/ittf/internal/aliases"
)

func execAliases() {
	fmt.Println()

	maxKeyLength := 0
	// Find the length of the longest key
	for key := range aliases.Aliases {
		if len(key) > maxKeyLength {
			maxKeyLength = len(key)
		}
	}
	for key, value := range aliases.Aliases {
		fmt.Printf("%-*s  =>  %s\n", maxKeyLength, key, value)
	}
	fmt.Println()
}
