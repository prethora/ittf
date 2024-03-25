package rules

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type rawRule struct {
	Vendor     []interface{} `yaml:"vendor"`
	Date       interface{}   `yaml:"date"`
	DateFormat string        `yaml:"dateFormat"`
	FileName   string        `yaml:"fileName"`
}

type RuleRegExp struct {
	Index  int
	RegEx  *regexp.Regexp
	After  *RuleRegExp
	Before *RuleRegExp
}

type Rule struct {
	Vendor     []*RuleRegExp
	Date       *RuleRegExp
	DateFormat string
	FileName   string
}

type Rules []*Rule

func ReadRulesFromFile(filePath string) (Rules, error) {
	var rawRules []rawRule

	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read rules file: %v", err)
	}

	// Unmarshal the YAML into the rules slice
	err = yaml.Unmarshal(content, &rawRules)
	if err != nil {
		return nil, fmt.Errorf("could not parse rules file: %v", err)
	}

	rules := Rules{}

	// Parse the raw rules into a slice of Rule structs, ready to be applied to the input file for testing
	for index, raw := range rawRules {
		rule, err := parseRawRule(raw)
		if err != nil {
			return nil, fmt.Errorf("could not parse rule [%d]: %v", index+1, err)
		}
		rules = append(rules, rule)
	}

	return rules, nil
}
