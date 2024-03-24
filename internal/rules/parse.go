package rules

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// parseRuleRegExp takes a string input and attempts to parse it into a RuleRegExp struct, which includes the zero-based index
// of the expected occurrence and the compiled regex.
func parseRuleRegExp(raw string) (*RuleRegExp, error) {
	// Initialize a new RuleRegExp with default values.
	ruleRegexp := &RuleRegExp{
		Index: 0,
		RegEx: nil,
	}

	// Trim leading and trailing whitespace from the input string.
	raw = strings.TrimSpace(raw)

	// Attempt to match and extract the index part from the beginning of the string.
	matches := regexp.MustCompile(`^(\d+)\/`).FindStringSubmatch(raw)
	if len(matches) > 0 {
		// If a match is found, parse the index as an integer.
		index, err := strconv.Atoi(matches[1])
		if err == nil {
			ruleRegexp.Index = index
			// Remove the index part from the raw string to leave only the regex pattern.
			raw = raw[len(matches[0]):]
		}
	}

	// Compile the remaining part of the string as a regex pattern.
	regex, err := regexp.Compile(raw)
	if err != nil {
		return nil, err
	}

	// Assign the compiled regex to the RuleRegExp struct.
	ruleRegexp.RegEx = regex

	// Return the populated RuleRegExp struct.
	return ruleRegexp, nil
}

// parseRawRule takes a rawRule struct and attempts to parse and validate it into a Rule struct, which can be directly applied
// to test the input file
func parseRawRule(raw rawRule) (*Rule, error) {
	// Initialize a new Rule with default values.
	rule := &Rule{
		Vendor:     []*RuleRegExp{},
		Date:       nil,
		DateFormat: "",
		FileName:   "",
	}

	// Validate that the vendor array was present in the rules file.
	if raw.Vendor == nil {
		return nil, errors.New("vendor array is missing")
	}

	// Validate that at least one vendor regexp was provided
	if len(raw.Vendor) == 0 {
		return nil, errors.New("vendor array is empty")
	}

	// Iterate through each vendor regex string.
	for index, value := range raw.Vendor {
		if value == "" {
			return nil, fmt.Errorf("vendor regexp [%d] is empty", index)
		}

		// Parse each vendor regex string into a RuleRegExp struct.
		vendorRegex, err := parseRuleRegExp(value)
		if err != nil {
			return nil, fmt.Errorf("could not compile vendor regexp [%d] (%s): %v", index, value, err)
		}
		// Append the parsed RuleRegExp to the Vendor slice of the Rule.
		rule.Vendor = append(rule.Vendor, vendorRegex)
	}

	// Validate that the date field was provided in the rules file and parse it.
	if raw.Date == "" {
		return nil, errors.New("date regexp is empty or missing")
	}

	dateRegex, err := parseRuleRegExp(raw.Date)
	if err != nil {
		return nil, fmt.Errorf("could not compile date regexp (%s): %v", raw.Date, err)
	}

	rule.Date = dateRegex

	// Trim leading and trailing whitespace from the dateFormat field and validate it.
	raw.DateFormat = strings.TrimSpace(raw.DateFormat)
	if raw.DateFormat == "" {
		return nil, errors.New("dateFormat is empty or missing")
	}

	// Validate the DateFormat by trying to parse it with itself (ensuring it's a valid date format string).
	_, err = time.Parse(raw.DateFormat, raw.DateFormat)
	if err != nil {
		return nil, fmt.Errorf("invalid dateFormat (%s)", raw.DateFormat)
	}

	rule.DateFormat = raw.DateFormat

	// Trim leading and trailing whitespace from the FileName field and set a default value if it's empty.
	raw.FileName = strings.TrimSpace(raw.FileName)
	if raw.FileName == "" {
		// The fileName defaults to {date} - {last vendor regexp match}.pdf
		rule.FileName = fmt.Sprintf("(date) - (%d).pdf", len(raw.Vendor)-1)
	} else {
		rule.FileName = raw.FileName
	}

	if !strings.Contains(rule.FileName, "(date)") {
		return nil, fmt.Errorf("invalid fileName format (%s): it does not contain the (date) placeholder", rule.FileName)
	}

	// Return the populated Rule struct.
	return rule, nil
}
