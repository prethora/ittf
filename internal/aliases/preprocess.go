package aliases

import "regexp"

func PreprocessAlias(value string) string {
	matches := regexp.MustCompile(`^\s*\$(\S+(?:\s+\S+)*)\s*$`).FindStringSubmatch(value)
	if len(matches) == 2 {
		mappedValue, exists := Aliases[matches[1]]
		if exists {
			return mappedValue
		} else {
			return value
		}
	} else {
		return value
	}
}
