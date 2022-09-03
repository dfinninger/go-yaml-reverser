package reverser

import (
	"bufio"
	"strings"
	"testing"
)

// TODO, struct of tests to test Yaml documents

type docTest struct {
	input    string
	expected string
}

func TestDocument(t *testing.T) {
	tests := []docTest{
		// basic test
		{
			`foo: "bar"`,
			`foo: "bar"`,
		},
		// document start
		{
			`---
foo: "bar"`,
			`---
foo: "bar"`,
		},
		// yaml version spec
		{
			`%YAML 1.2
---
foo: "bar"`,
			`%YAML 1.2
---
foo: "bar"`,
		},
		// document end
		{
			`%YAML 1.2
---
foo: "bar"
...`,
			`%YAML 1.2
---
foo: "bar"
...`,
		},
		// document end and bare doc
		{
			`%YAML 1.2
---
foo: "bar"
...
a: "b"`,
			`a: "b"
%YAML 1.2
---
foo: "bar"
...`,
		},
		// two doc start, no doc end
		{
			`---
foo: "bar"
---
baz: "quux"`,
			`---
baz: "quux"
---
foo: "bar"`,
		},
	}

	runTests(t, tests)
}

func runTests(t *testing.T, tests []docTest) {
	t.Helper()

	for _, tt := range tests {
		in := bufio.NewScanner(strings.NewReader(tt.input))
		doc, err := reverseStream(in)
		if err != nil {
			t.Fatalf("decode error:\ninput:\n%s\n\noutput:\n%v", tt.input, doc)
		}

		out := doc.String()
		if out != tt.expected {
			t.Errorf("unexpected output\n=== actual ===\n%s\n=== expected ===\n%s", out, tt.expected)
		}
	}
}
