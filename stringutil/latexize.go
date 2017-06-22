package stringutil

var additionalTexReplacements = []Converter{
	NewStringConverter("{", "\\{"),
	NewStringConverter("}", "\\}"),
	NewStringConverter("\\", "\\textbackslash{}"),

	NewStringConverter("&", "\\&"),
	NewStringConverter("%", "\\%"),
	NewStringConverter("$", "\\$"),
	NewStringConverter("#", "\\#"),
	NewStringConverter("_", "\\_"),
	NewStringConverter("~", "\\textasciitilde{}"),
	NewStringConverter("^", "\\textasciicircum{}"),
	NewStringConverter("ß", "\\ss{}"),
}

/*
Latexize takes an input text as parsed from the value of any field in a project
file and turns it into LaTeX code.
*/
func Latexize(input string) string {
	renderedString := CustomEscape(input, additionalTexReplacements...)
	return renderedString
}
