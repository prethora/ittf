package rules

import (
	"math"
)

func matchRuleRegExpIndex(str string, ruleRegExp *RuleRegExp, useSubgroupMatch bool) ([]int, bool) {
	beg := 0
	end := len(str)

	if ruleRegExp.After != nil {
		if match, matched := matchRuleRegExpIndex(str, ruleRegExp.After, false); matched {
			beg = match[1]
		} else {
			return nil, false
		}
	}

	if ruleRegExp.Before != nil {
		if match, matched := matchRuleRegExpIndex(str, ruleRegExp.Before, false); matched {
			end = match[0]
		} else {
			return nil, false
		}
	}

	if end < beg {
		return nil, false
	}

	str = str[beg:end]

	maxSearch := ruleRegExp.Index
	requiredCount := int(math.Abs(float64(ruleRegExp.Index)))

	if ruleRegExp.Index < 0 {
		maxSearch = -1 // If the index is negative, do not limit the number of occurences, as we are counting backwards.
	}

	matchesSet := ruleRegExp.RegEx.FindAllStringSubmatchIndex(str, maxSearch)
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

func MatchRuleRegExp(str string, ruleRegExp *RuleRegExp, useSubgroupMatch bool) (string, bool) {
	if match, matched := matchRuleRegExpIndex(str, ruleRegExp, useSubgroupMatch); matched {
		return str[match[0]:match[1]], true
	} else {
		return "", false
	}
}
