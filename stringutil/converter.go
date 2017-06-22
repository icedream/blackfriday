package stringutil

type Converter interface {
	Process(text string) string
}
