package rules

import (
	"math"
)

func MatchRuleRegExp(str string, ruleRegExp *RuleRegExp, useSubgroupMatch bool) (string, bool) {
	maxSearch := ruleRegExp.Index
	requiredCount := int(math.Abs(float64(ruleRegExp.Index)))

	if ruleRegExp.Index < 0 {
		maxSearch = -1 // If the index is negative, do not limit the number of occurences, as we are counting backwards.
	}

	matchesSet := ruleRegExp.RegEx.FindAllStringSubmatch(str, maxSearch)
	if requiredCount <= len(matchesSet) { // check if enough matches were found to include the provided index.
		var requiredIndex int

		if ruleRegExp.Index > 0 {
			requiredIndex = ruleRegExp.Index - 1
		} else {
			requiredIndex = len(matchesSet) + ruleRegExp.Index
		}

		matches := matchesSet[requiredIndex]

		var result string

		// if the --subgroup-match flag was set, attempt to extract a single subgroup match.
		if len(matches) == 2 && useSubgroupMatch {
			result = matches[1]
		} else {
			result = matches[0]
		}

		return result, true
	} else {
		return "", false
	}
}
