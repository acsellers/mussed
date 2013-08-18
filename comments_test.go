package mussed

import (
	"bytes"
	"encoding/json"
	"testing"
	"text/template"
)

/*
Comment tags represent content that should never appear in the resulting
output.

The tag's content may contain any substring (including newlines) EXCEPT the
closing delimiter.

Comment tags SHOULD be treated as standalone when appropriate.

*/

func TestCOMMENTS0(t *testing.T) {
	// Inline

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.mustache", `12345{{! Comment Block! }}67890`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte("{}"), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`1234567890`, b.String())
	})
}

func TestCOMMENTS1(t *testing.T) {
	// Multiline

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.mustache", `12345{{!
  This is a
  multi-line comment...
}}67890
`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte("{}"), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`1234567890
`, b.String())
	})
}

func TestCOMMENTS2(t *testing.T) {
	// Standalone

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.mustache", `Begin.
{{! Comment Block! }}
End.
`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte("{}"), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`Begin.

End.
`, b.String())
	})
}

func TestCOMMENTS3(t *testing.T) {
	// Indented Standalone

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.mustache", `Begin.
  {{! Indented Comment Block! }}
End.
`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte("{}"), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`Begin.
  
End.
`, b.String())
	})
}

func TestCOMMENTS4(t *testing.T) {
	// Standalone Line Endings

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.mustache", `|
{{! Standalone Comment }}
|`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte("{}"), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`|

|`, b.String())
	})
}

func TestCOMMENTS5(t *testing.T) {
	// Standalone Without Previous Line

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.mustache", `  {{! I'm Still Standalone }}
!`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte("{}"), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`  
!`, b.String())
	})
}

func TestCOMMENTS6(t *testing.T) {
	// Standalone Without Newline

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.mustache", `!
  {{! I'm Still Standalone }}`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte("{}"), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`!
  `, b.String())
	})
}

func TestCOMMENTS7(t *testing.T) {
	// Multiline Standalone

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.mustache", `Begin.
{{!
Something's going on here...
}}
End.
`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte("{}"), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`Begin.

End.
`, b.String())
	})
}

func TestCOMMENTS8(t *testing.T) {
	// Indented Multiline Standalone

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.mustache", `Begin.
  {{!
    Something's going on here...
  }}
End.
`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte("{}"), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`Begin.
  
End.
`, b.String())
	})
}

func TestCOMMENTS9(t *testing.T) {
	// Indented Inline

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.mustache", `  12 {{! 34 }}
`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte("{}"), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`  12 
`, b.String())
	})
}

func TestCOMMENTS10(t *testing.T) {
	// Surrounding Whitespace

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.mustache", `12345 {{! Comment Block! }} 67890`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte("{}"), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`12345  67890`, b.String())
	})
}
