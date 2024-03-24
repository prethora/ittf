package commands

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/prethora/ittf/internal/rules"
)

// execRoot is the main function that takes the input file path and the rules file path, processes the input content based on
// the rules, and generates an output filename - or exits with status 1 if no match is found.
func execRoot() {
	// Read the input content from the file specified by inputFilePath.

	_, err := time.Parse(dateOutputFormat, dateOutputFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: provided date output format is invalid (%s)\n", dateOutputFormat)
		os.Exit(2)
	}

	inputContent, err := os.ReadFile(inputFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not read input file: %v\n", err)
		os.Exit(2)
	}

	// Convert the input content to a string for processing.
	inputContentStr := string(inputContent)

	// Read and parse the rules from the provided rules file path.
	rulesList, err := rules.ReadRulesFromFile(rulesFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	// Iterate through each rule in the rules list.
	for index, rule := range rulesList {
		var vendorMatches []string // To hold matches for vendor patterns.
		var dateMatch string       // To hold the date match.

		// Process each vendor rule regex pattern.
		for _, ruleRegExp := range rule.Vendor {
			matchesSet := ruleRegExp.RegEx.FindAllStringSubmatch(inputContentStr, ruleRegExp.Index+1) // match up to the provided index (defaults to 0)
			if ruleRegExp.Index < len(matchesSet) {                                                   // check if enough matches were found to include the provided index
				matches := matchesSet[ruleRegExp.Index]
				// Append the first subgroup match, if available, otherwise the entire match.
				if len(matches) == 2 {
					vendorMatches = append(vendorMatches, matches[1])
				} else {
					vendorMatches = append(vendorMatches, matches[0])
				}
			} else {
				break // Exit the loop if no match is found for the current index.
			}
		}

		// Continue to the next rule if the number of vendor matches is less than expected (meaning at least one did not match).
		if len(vendorMatches) < len(rule.Vendor) {
			continue
		}

		// Process the date rule regex pattern.
		matchesSet := rule.Date.RegEx.FindAllStringSubmatch(inputContentStr, rule.Date.Index+1)
		if rule.Date.Index >= len(matchesSet) {
			continue // Skip to the next rule if no date match was found.
		}

		matches := matchesSet[rule.Date.Index]
		// Extract the date match, preferring the first subgroup match if available.
		if len(matches) == 2 {
			dateMatch = matches[1]
		} else {
			dateMatch = matches[0]
		}

		// Parse the matched date string using the rule's date format.
		date, err := time.Parse(rule.DateFormat, dateMatch)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: matched rule [%d] but failed to parse date (%s): %v\n", index, dateMatch, err)
			os.Exit(2)
		}

		// Format the parsed date according to the required output format.
		dateOutput := date.Format(dateOutputFormat)

		// Initialize a map to hold the placeholders for file path template and their corresponding values.
		filePathTemplateMap := map[string]string{}

		// Populate the map with vendor matches.
		for index, vendorMatch := range vendorMatches {
			filePathTemplateMap[fmt.Sprintf("(%d)", index)] = vendorMatch
		}

		// Add the formatted date to the map.
		filePathTemplateMap["(date)"] = dateOutput

		// Start constructing the output file name using a builder for efficient string concatenation.
		fileName := rule.FileName
		var builder strings.Builder

		// Compile a regex to match placeholders in the file name template.
		sepRegex := regexp.MustCompile(`\([date0-9]+\)`)

		// Loop through the file name, replacing placeholders with their corresponding values from the map.
		for {
			match := sepRegex.FindStringIndex(fileName)
			if match != nil {
				// Append the literal part of the file name.
				builder.WriteString(fileName[0:match[0]])
				// Extract the placeholder.
				placeholder := fileName[match[0]:match[1]]
				// Replace the placeholder with its value if it exists in the map.
				if value, exists := filePathTemplateMap[placeholder]; exists {
					builder.WriteString(value)
				} else {
					// If the placeholder doesn't have a corresponding value, keep it as is.
					builder.WriteString(placeholder)
				}
				// Move to the next segment of the file name.
				fileName = fileName[match[1]:]
			} else {
				// Append the remaining part of the file name and break out of the loop.
				builder.WriteString(fileName)
				break
			}
		}

		// Finalize the constructed file name.
		result := builder.String()

		// Output the result and exit the program successfully.
		fmt.Println(result)
		os.Exit(0)
	}

	// If no rule matched, exit the program with status 1 indicating failure.
	os.Exit(1)
}