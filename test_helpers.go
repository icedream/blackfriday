//
// Blackfriday Markdown Processor
// Available at http://github.com/russross/blackfriday
//
// Copyright © 2011 Russ Ross <russ@russross.com>.
// Distributed under the Simplified BSD License.
// See README.md for details.
//

//
// Helper functions for unit testing
//

package blackfriday

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"testing"
)

func runMarkdownBlockWithRenderer(input string, extensions Extensions, renderer Renderer) string {
	return string(Markdown([]byte(input), renderer, extensions))
}

func runMarkdownBlock(input string, extensions Extensions) string {
	renderer := HTMLRenderer(UseXHTML, extensions, "", "")
	return runMarkdownBlockWithRenderer(input, extensions, renderer)
}

func runnerWithRendererParameters(parameters HTMLRendererParameters) func(string, Extensions) string {
	return func(input string, extensions Extensions) string {
		renderer := HTMLRendererWithParameters(UseXHTML, extensions, "", "", parameters)
		return runMarkdownBlockWithRenderer(input, extensions, renderer)
	}
}

func doTestsBlock(t *testing.T, tests []string, extensions Extensions) {
	doTestsBlockWithRunner(t, tests, extensions, runMarkdownBlock)
}

func doTestsBlockWithRunner(t *testing.T, tests []string, extensions Extensions, runner func(string, Extensions) string) {
	// catch and report panics
	var candidate string
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("\npanic while processing [%#v]: %s\n", candidate, err)
		}
	}()

	for i := 0; i+1 < len(tests); i += 2 {
		input := tests[i]
		candidate = input
		expected := tests[i+1]
		actual := runner(candidate, extensions)
		if actual != expected {
			t.Errorf("\nInput   [%#v]\nExpected[%#v]\nActual  [%#v]",
				candidate, expected, actual)
		}

		// now test every substring to stress test bounds checking
		if !testing.Short() {
			for start := 0; start < len(input); start++ {
				for end := start + 1; end <= len(input); end++ {
					candidate = input[start:end]
					_ = runMarkdownBlock(candidate, extensions)
				}
			}
		}
	}
}

func runMarkdownInline(input string, opts Options, htmlFlags HTMLFlags, params HTMLRendererParameters) string {
	opts.Extensions |= Autolink
	opts.Extensions |= Strikethrough

	htmlFlags |= UseXHTML

	renderer := HTMLRendererWithParameters(htmlFlags, opts.Extensions, "", "", params)

	return string(MarkdownOptions([]byte(input), renderer, opts))
}

func doTestsInline(t *testing.T, tests []string) {
	doTestsInlineParam(t, tests, Options{}, 0, HTMLRendererParameters{})
}

func doLinkTestsInline(t *testing.T, tests []string) {
	doTestsInline(t, tests)

	prefix := "http://localhost"
	params := HTMLRendererParameters{AbsolutePrefix: prefix}
	transformTests := transformLinks(tests, prefix)
	doTestsInlineParam(t, transformTests, Options{}, 0, params)
	doTestsInlineParam(t, transformTests, Options{}, CommonHtmlFlags, params)
}

func doSafeTestsInline(t *testing.T, tests []string) {
	doTestsInlineParam(t, tests, Options{}, Safelink, HTMLRendererParameters{})

	// All the links in this test should not have the prefix appended, so
	// just rerun it with different parameters and the same expectations.
	prefix := "http://localhost"
	params := HTMLRendererParameters{AbsolutePrefix: prefix}
	transformTests := transformLinks(tests, prefix)
	doTestsInlineParam(t, transformTests, Options{}, Safelink, params)
}

func doTestsInlineParam(t *testing.T, tests []string, opts Options, htmlFlags HTMLFlags, params HTMLRendererParameters) {
	// catch and report panics
	var candidate string
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("\npanic while processing [%#v] (%v)\n", candidate, err)
		}
	}()

	for i := 0; i+1 < len(tests); i += 2 {
		input := tests[i]
		candidate = input
		expected := tests[i+1]
		actual := runMarkdownInline(candidate, opts, htmlFlags, params)
		if actual != expected {
			t.Errorf("\nInput   [%#v]\nExpected[%#v]\nActual  [%#v]",
				candidate, expected, actual)
		}

		// now test every substring to stress test bounds checking
		if !testing.Short() {
			for start := 0; start < len(input); start++ {
				for end := start + 1; end <= len(input); end++ {
					candidate = input[start:end]
					_ = runMarkdownInline(candidate, opts, htmlFlags, params)
				}
			}
		}
	}
}

func transformLinks(tests []string, prefix string) []string {
	newTests := make([]string, len(tests))
	anchorRe := regexp.MustCompile(`<a href="/(.*?)"`)
	imgRe := regexp.MustCompile(`<img src="/(.*?)"`)
	for i, test := range tests {
		if i%2 == 1 {
			test = anchorRe.ReplaceAllString(test, `<a href="`+prefix+`/$1"`)
			test = imgRe.ReplaceAllString(test, `<img src="`+prefix+`/$1"`)
		}
		newTests[i] = test
	}
	return newTests
}

func runMarkdownReference(input string, flag Extensions) string {
	renderer := HTMLRenderer(0, flag, "", "")
	return string(Markdown([]byte(input), renderer, flag))
}

func doTestsReference(t *testing.T, files []string, flag Extensions) {
	// catch and report panics
	var candidate string
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("\npanic while processing [%#v]\n", candidate)
		}
	}()

	for _, basename := range files {
		filename := filepath.Join("testdata", basename+".text")
		inputBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Errorf("Couldn't open '%s', error: %v\n", filename, err)
			continue
		}
		input := string(inputBytes)

		filename = filepath.Join("testdata", basename+".html")
		expectedBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Errorf("Couldn't open '%s', error: %v\n", filename, err)
			continue
		}
		expected := string(expectedBytes)

		// fmt.Fprintf(os.Stderr, "processing %s ...", filename)
		actual := string(runMarkdownReference(input, flag))
		if actual != expected {
			t.Errorf("\n    [%#v]\nExpected[%#v]\nActual  [%#v]",
				basename+".text", expected, actual)
		}
		// fmt.Fprintf(os.Stderr, " ok\n")

		// now test every prefix of every input to check for
		// bounds checking
		if !testing.Short() {
			start, max := 0, len(input)
			for end := start + 1; end <= max; end++ {
				candidate = input[start:end]
				// fmt.Fprintf(os.Stderr, "  %s %d:%d/%d\n", filename, start, end, max)
				_ = runMarkdownReference(candidate, flag)
			}
		}
	}
}