package texts

import "strings"

func JoinListLinesWithPrefix(list []string, prefix string) string {
	if len(list) == 0 {
		return ""
	}

	result := make([]string, 0, len(list))
	for _, line := range list {
		result = append(result, prefix+line)
	}

	return strings.Join(result, "\n")
}
