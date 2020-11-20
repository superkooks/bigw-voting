package bgw

import (
	"sort"
)

// SortByAddr returns a list of ints ordered by IP address
func SortByAddr(m map[string]int) []int {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var out []int
	for _, k := range keys {
		out = append(out, m[k])
	}

	return out
}
