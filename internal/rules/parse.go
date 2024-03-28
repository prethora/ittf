package rules

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/prethora/ittf/internal/aliases"
)

func parseRegExp(value interface{}) (*regexp.Regexp, error) {
	switch valueValue := value.(type) {
	case string:
		regex, err := regexp.Compile(aliases.PreprocessAlias(strings.TrimSpace(valueValue)))
		if err != nil {
			return nil, err
		}
		return regex, nil
	case int:
		regex, err := regexp.Compile(strconv.Itoa(valueValue))
		if err != nil {
			return nil, err
		}
		return regex, nil
	case float64:
		regex, err := regexp.Compile(strconv.FormatFloat(valueValue, 'f', -1, 64))
		if err != nil {
			return nil, err
		}
		return regex, nil
	default:
		return nil, nil
	}
}

func parseRuleRegExpIndex(value map[string]interface{}) (int, error) {
	if index, exists := value["Index"]; exists {
		if indexInt, ok := index.(int); ok {
			if indexInt != 0 {
				return indexInt, nil
			} else {
				return 0, errors.New("'Index' if set must be a non-zero value")
			}
		} else {
			return 0, errors.New("'Index' if set must be an integer")
		}
	} else {
		return 1, nil
	}
}

func parseRuleRegExpMatch(value map[string]interface{}) (*regexp.Regexp, error) {
	if match, exists := value["Match"]; exists {
		regexp, err := parseRegExp(match)
		if err != nil {
			return nil, err
		}
		if regexp == nil {
			return nil, errors.New("could not recognize match field type")
		}
		return regexp, nil
	} else {
		return nil, errors.New("the match field is required")
	}
}

func parseRuleRegExpSubRuleRegExp(value map[string]interface{}, fieldName string) (*RuleRegExp, error) {
	if raw, exists := value[fieldName]; exists {
		regexpRule, err := parseRuleRegExp(raw, fmt.Sprintf(".%s", fieldName))
		if err != nil {
			return nil, err
		}
		return regexpRule, nil
	} else {
		return nil, nil
	}
}

// parseRuleRegExp takes a string input and attempts to parse it into a RuleRegExp struct, which includes the one-based index
// of the expected occurrence and the compiled regex.
func parseRuleRegExp(raw interface{}, errorPrefix string) (*RuleRegExp, error) {
	ruleRegexp := &RuleRegExp{
		Index:  1, // The default value, if the prefix is not specified
		RegEx:  nil,
		After:  nil,
		Before: nil,
	}

	regexp, err := parseRegExp(raw)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", errorPrefix, err)
	}
	if regexp != nil {
		ruleRegexp.RegEx = regexp
	} else {
		if rawValue, ok := raw.(map[string]interface{}); ok {
			index, err := parseRuleRegExpIndex(rawValue)
			if err != nil {
				return nil, fmt.Errorf("%s: %v", errorPrefix, err)
			}
			ruleRegexp.Index = index

			regex, err := parseRuleRegExpMatch(rawValue)
			if err != nil {
				return nil, fmt.Errorf("%s: %v", errorPrefix, err)
			}
			ruleRegexp.RegEx = regex

			after, err := parseRuleRegExpSubRuleRegExp(rawValue, "After")
			if err != nil {
				return nil, fmt.Errorf("%s%v", errorPrefix, err)
			}
			if after != nil {
				ruleRegexp.After = after
			}

			before, err := parseRuleRegExpSubRuleRegExp(rawValue, "Before")
			if err != nil {
				return nil, fmt.Errorf("%s%v", errorPrefix, err)
			}
			if before != nil {
				ruleRegexp.Before = before
			}

		} else {
			return nil, fmt.Errorf("%s: unrecognized type", errorPrefix)
		}
	}

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
			return nil, fmt.Errorf("vendor regexp [%d] is empty", index+1)
		}

		// Parse each vendor regex string into a RuleRegExp struct.
		vendorRegex, err := parseRuleRegExp(value, fmt.Sprintf("vendor regexp[%d]", index+1))
		if err != nil {
			return nil, fmt.Errorf("could not parse %v", err)
		}
		// Append the parsed RuleRegExp to the Vendor slice of the Rule.
		rule.Vendor = append(rule.Vendor, vendorRegex)
	}

	// Validate that the date field was provided in the rules file and parse it.
	if raw.Date == "" {
		return nil, errors.New("date regexp is empty or missing")
	}

	dateRegex, err := parseRuleRegExp(raw.Date, "date regexp")
	if err != nil {
		return nil, fmt.Errorf("could not compile %v", err)
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

	rule.DateFormat = aliases.PreprocessAlias(strings.TrimSpace(raw.DateFormat))
	rule.BaseName = strings.TrimSpace(raw.BaseName)
	rule.FileName = strings.TrimSpace(raw.FileName)

	if raw.FileName != "" && !strings.Contains(rule.FileName, "(date)") {
		return nil, fmt.Errorf("invalid fileName format (%s): it does not contain the (date) placeholder", rule.FileName)
	}

	// Return the populated Rule struct.
	return rule, nil
}
