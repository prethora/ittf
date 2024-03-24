package commands

import (
	"fmt"
	"os"

	"github.com/prethora/ittf/internal/rules"
)

func execValidate() {
	// ANSI escape codes
	green := "\033[32m"
	lightRed := "\033[91m"
	reset := "\033[0m"

	// Read and parse the rules from the provided rules file path.
	rulesList, err := rules.ReadRulesFromFile(rulesFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%serror: %v%s\n", lightRed, err, reset)
		os.Exit(1)
	}

	// Prepare rule(s) suffix based on the number of rules found
	suffix := "s"
	if len(rulesList) == 1 {
		suffix = ""
	}

	// Output success status.
	fmt.Printf("%s✓ Rules file is OK\n\n%s%d rule%s found.\n", green, reset, len(rulesList), suffix)
}
