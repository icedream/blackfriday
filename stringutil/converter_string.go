package stringutil

import "strings"

type stringConverter struct {
	old         string
	replacement string
}

func NewStringConverter(old string, replacementStr string) *stringConverter {
	return &stringConverter{
		old:         old,
		replacement: replacementStr,
	}
}

func (conv *stringConverter) Process(text string) string {
	return strings.Replace(text, conv.old, conv.replacement, -1)
}
