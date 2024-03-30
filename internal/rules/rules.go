package rules

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type rawRule struct {
	Vendor   []interface{} `yaml:"Matches"`
	Index    int           `yaml:"Index"`
	After    interface{}   `yaml:"After"`
	Before   interface{}   `yaml:"Before"`
	Date     string        `yaml:"Date"`
	BaseName string        `yaml:"Basename"`
	FileName string        `yaml:"Output"`
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
	BaseName   string
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
