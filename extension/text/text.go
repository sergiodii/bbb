package text

import (
	"regexp"
)

func GetPartOfString(s string, regexString string) string {
	re := regexp.MustCompile(regexString)
	matches := re.FindStringSubmatch(s)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
