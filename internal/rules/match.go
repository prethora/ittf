package rules

import (
	"math"
	"regexp"
	"strings"
	"time"
)

func replaceUsingMap(s string, replacements map[string]string) string {
	for old, new := range replacements {
		s = strings.Replace(s, old, new, -1) // -1 means replace all occurrences
	}
	return s
}

func PreProcessDate(dateStr string) string {
	regex := regexp.MustCompile(`(?i)(\d+)(st|nd|rd|th)`)
	replacements := map[string]string{
		"Janu ": "Jan ",
		"janu ": "jan ",
		"Febr ": "Feb ",
		"febr ": "feb ",
		"Marc ": "Mar ",
		"marc ": "mar ",
		"Apri ": "Apr ",
		"apri ": "apr ",
		"Augu ": "Aug ",
		"augu ": "aug ",
		"Sept ": "Sep ",
		"sept ": "sep ",
		"Octo ": "Oct ",
		"octo ": "oct ",
		"Nove ": "Nov ",
		"nove ": "nov ",
		"Dece ": "Dec ",
		"dece ": "dec ",
		// Versions with "/"
		"Janu/": "Jan/",
		"janu/": "jan/",
		"Febr/": "Feb/",
		"febr/": "feb/",
		"Marc/": "Mar/",
		"marc/": "mar/",
		"Apri/": "Apr/",
		"apri/": "apr/",
		"Augu/": "Aug/",
		"augu/": "aug/",
		"Sept/": "Sep/",
		"sept/": "sep/",
		"Octo/": "Oct/",
		"octo/": "oct/",
		"Nove/": "Nov/",
		"nove/": "nov/",
		"Dece/": "Dec/",
		"dece/": "dec/",
	}
	return replaceUsingMap(regex.ReplaceAllString(dateStr, `$1`), replacements)
}

func matchRuleRegExpIndex(str string, ruleRegExp *RuleRegExp, useSubgroupMatch bool, mustMatchDateFormat string) ([]int, bool) {
	beg := 0
	end := len(str)

	if ruleRegExp.After != nil {
		if match, matched := matchRuleRegExpIndex(str, ruleRegExp.After, false, ""); matched {
			beg = match[1]
		} else {
			return nil, false
		}
	}

	if ruleRegExp.Before != nil {
		if match, matched := matchRuleRegExpIndex(str, ruleRegExp.Before, false, ""); matched {
			end = match[0]
		} else {
			return nil, false
		}
	}

	if end < beg {
		return nil, false
	}

	str = str[beg:end]

	requiredCount := int(math.Abs(float64(ruleRegExp.Index)))

	matchesSet := ruleRegExp.RegEx.FindAllStringSubmatchIndex(str)

	if mustMatchDateFormat != "" {
		_matchesSet := matchesSet
		matchesSet = [][]int{}

		for _, matches := range _matchesSet {
			var result []int
			if len(matches) == 4 && useSubgroupMatch {
				result = matches[2:]
			} else {
				result = matches[0:2]
			}
			matchStr := str[result[0]:result[1]]
			_, err := time.Parse(mustMatchDateFormat, PreProcessDate(matchStr))
			if err == nil {
				matchesSet = append(matchesSet, matches)
			}
		}
	}

	if requiredCount <= len(matchesSet) { // check if enough matches were found to include the provided index.
		var requiredIndex int

		if ruleRegExp.Index > 0 {
			requiredIndex = ruleRegExp.Index - 1
		} else {
			requiredIndex = len(matchesSet) + ruleRegExp.Index
		}

		matches := matchesSet[requiredIndex]

		var result []int

		// if the --subgroup-match flag was set, attempt to extract a single subgroup match.
		if len(matches) == 4 && useSubgroupMatch {
			result = matches[2:]
		} else {
			result = matches[0:2]
		}

		result[0] += beg
		result[1] += beg

		return result, true
	} else {
		return nil, false
	}
}

func MatchRuleRegExp(str string, ruleRegExp *RuleRegExp, useSubgroupMatch bool, mustMatchDateFormat string) (string, bool) {
	if match, matched := matchRuleRegExpIndex(str, ruleRegExp, useSubgroupMatch, mustMatchDateFormat); matched {
		return str[match[0]:match[1]], true
	} else {
		return "", false
	}
}
