package utils

import "strings"

func mergeEnv(base []string, overrides []string) []string {
	overrideKeys := make(map[string]bool)
	for _, e := range overrides {
		overrideKeys[strings.SplitN(e, "=", 2)[0]] = true
	}

	result := make([]string, 0, len(base)+len(overrides))
	for _, e := range base {
		if !overrideKeys[strings.SplitN(e, "=", 2)[0]] {
			result = append(result, e)
		}
	}

	return append(result, overrides...)
}