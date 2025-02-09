package texts

func RuneSetFromString(s string) map[rune]struct{} {
	set := make(map[rune]struct{}, len(s))
	for _, r := range s {
		set[r] = struct{}{}
	}

	return set
}
