package stringutil

// BackslashEscape modifies a string so any occurrences of given substrings
// are prefixed with a backslash.
func BackslashEscape(s string, substrings ...string) string {
	replacements := []Converter{}

	for _, substring := range substrings {
		replacements = append(replacements, NewStringConverter(substring, "\\"+substring))
	}

	return CustomEscape(s, replacements...)
}

// CustomEscape modifies a string according to the given conversion rules in
// the exact given order.
func CustomEscape(s string, replacements ...Converter) string {
	for _, replacer := range replacements {
		s = replacer.Process(s)
	}

	return s
}
