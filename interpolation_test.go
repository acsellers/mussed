package mussed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"testing"
)

/*
Interpolation tags are used to integrate dynamic content into the template.

The tag's content MUST be a non-whitespace character sequence NOT containing
the current closing delimiter.

This tag's content names the data to replace the tag.  A single period (`.`)
indicates that the item currently sitting atop the context stack should be
used; otherwise, name resolution is as follows:
  1) Split the name on periods; the first part is the name to resolve, any
  remaining parts should be retained.
  2) Walk the context stack from top to bottom, finding the first context
  that is a) a hash containing the name as a key OR b) an object responding
  to a method with the given name.
  3) If the context is a hash, the data is the value associated with the
  name.
  4) If the context is an object, the data is the value returned by the
  method with the given name.
  5) If any name parts were retained in step 1, each should be resolved
  against a context stack containing only the result from the former
  resolution.  If any part fails resolution, the result should be considered
  falsey, and should interpolate as the empty string.
Data should be coerced into a string (and escaped, if appropriate) before
interpolation.

The Interpolation tags MUST NOT be treated as standalone.

*/

func TestINTERPOLATION0(t *testing.T) {
	// No Interpolation

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `Hello from {Mustache}!
`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`Hello from {Mustache}!
`, b.String())
	})
}

func TestINTERPOLATION1(t *testing.T) {
	// Basic Interpolation

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `Hello, {{subject}}!
`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"subject":"world"}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`Hello, world!
`, b.String())
	})
}

func TestINTERPOLATION2(t *testing.T) {
	// HTML Escaping

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `These characters should be HTML escaped: {{forbidden}}
`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"forbidden":"& \" \u003c \u003e"}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`These characters should be HTML escaped: &amp; &#34; &lt; &gt;
`, b.String())
	})
}

func TestINTERPOLATION3(t *testing.T) {
	// Triple Mustache

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `These characters should not be HTML escaped: {{{forbidden}}}
`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"forbidden":"& \" < >"}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`These characters should not be HTML escaped: & " < >
`, b.String())
	})
}

func TestINTERPOLATION4(t *testing.T) {
	// Ampersand

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `These characters should not be HTML escaped: {{&forbidden}}
`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"forbidden":"& \" \u003c \u003e"}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`These characters should not be HTML escaped: & " < >
`, b.String())
	})
}

func TestINTERPOLATION5(t *testing.T) {
	// Basic Integer Interpolation

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `"{{mph}} miles an hour!"`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"mph":85}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`"85 miles an hour!"`, b.String())
	})
}

func TestINTERPOLATION6(t *testing.T) {
	// Triple Mustache Integer Interpolation

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `"{{{mph}}} miles an hour!"`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"mph":85}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "tree", data))
		test.AreEqual(`"85 miles an hour!"`, b.String())
	})
}

func TestINTERPOLATION7(t *testing.T) {
	// Ampersand Integer Interpolation

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `"{{&mph}} miles an hour!"`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"mph":85}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`"85 miles an hour!"`, b.String())
	})
}

func TestINTERPOLATION8(t *testing.T) {
	// Basic Decimal Interpolation

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `"{{power}} jiggawatts!"`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"power":1.21}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`"1.21 jiggawatts!"`, b.String())
	})
}

func TestINTERPOLATION9(t *testing.T) {
	// Triple Mustache Decimal Interpolation

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `"{{{power}}} jiggawatts!"`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"power":1.21}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`"1.21 jiggawatts!"`, b.String())
	})
}

func TestINTERPOLATION10(t *testing.T) {
	// Ampersand Decimal Interpolation

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `"{{&power}} jiggawatts!"`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"power":1.21}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`"1.21 jiggawatts!"`, b.String())
	})
}

func TestINTERPOLATION11(t *testing.T) {
	// Basic Context Miss Interpolation

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `I ({{cannot}}) be seen!`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`I () be seen!`, b.String())
	})
}

func TestINTERPOLATION12(t *testing.T) {
	// Triple Mustache Context Miss Interpolation

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `I ({{{cannot}}}) be seen!`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "tree", data))
		test.AreEqual(`I () be seen!`, b.String())
	})
}

func TestINTERPOLATION13(t *testing.T) {
	// Ampersand Context Miss Interpolation

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `I ({{&cannot}}) be seen!`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "tree", data))
		test.AreEqual(`I () be seen!`, b.String())
	})
}

func TestINTERPOLATION14(t *testing.T) {
	// Dotted Names - Basic Interpolation

	fmt.Println("here")
	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `"{{person.name}}" == "{{#person}}{{name}}{{/person}}"`)
		fmt.Println("wat")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"person":{"name":"Joe"}}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`"Joe" == "Joe"`, b.String())
	})

	fmt.Println("here")
}

func TestINTERPOLATION15(t *testing.T) {
	// Dotted Names - Triple Mustache Interpolation

	fmt.Println("here")
	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `"{{{person.name}}}" == "{{#person}}{{{name}}}{{/person}}"`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"person":{"name":"Joe"}}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`"Joe" == "Joe"`, b.String())
	})
}

func TestINTERPOLATION16(t *testing.T) {
	// Dotted Names - Ampersand Interpolation

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `"{{&person.name}}" == "{{#person}}{{&name}}{{/person}}"`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"person":{"name":"Joe"}}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`"Joe" == "Joe"`, b.String())
	})
}

func TestINTERPOLATION17(t *testing.T) {
	// Dotted Names - Arbitrary Depth

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `"{{a.b.c.d.e.name}}" == "Phil"`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"a":{"b":{"c":{"d":{"e":{"name":"Phil"}}}}}}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`"Phil" == "Phil"`, b.String())
	})
}

func TestINTERPOLATION18(t *testing.T) {
	// Dotted Names - Broken Chains

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `"{{a.b.c}}" == ""`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"a":{}}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`"" == ""`, b.String())
	})
}

func TestINTERPOLATION19(t *testing.T) {
	// Dotted Names - Broken Chain Resolution

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `"{{a.b.c.name}}" == ""`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"a":{"b":{}},"c":{"name":"Jim"}}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`"" == ""`, b.String())
	})
}

func TestINTERPOLATION20(t *testing.T) {
	// Dotted Names - Initial Resolution

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `"{{#a}}{{b.c.d.e.name}}{{/a}}" == "Phil"`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"a":{"b":{"c":{"d":{"e":{"name":"Phil"}}}}},"b":{"c":{"d":{"e":{"name":"Wrong"}}}}}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`"Phil" == "Phil"`, b.String())
	})
}

func TestINTERPOLATION21(t *testing.T) {
	// Interpolation - Surrounding Whitespace

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `| {{string}} |`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"string":"---"}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`| --- |`, b.String())
	})
}

func TestINTERPOLATION22(t *testing.T) {
	// Triple Mustache - Surrounding Whitespace

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `| {{{string}}} |`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"string":"---"}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`| --- |`, b.String())
	})
}

func TestINTERPOLATION23(t *testing.T) {
	// Ampersand - Surrounding Whitespace

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `| {{&string}} |`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"string":"---"}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`| --- |`, b.String())
	})
}

func TestINTERPOLATION24(t *testing.T) {
	// Interpolation - Standalone

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `  {{string}}
`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"string":"---"}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`  ---
`, b.String())
	})
}

func TestINTERPOLATION25(t *testing.T) {
	// Triple Mustache - Standalone

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `  {{{string}}}
`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"string":"---"}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`  ---
`, b.String())
	})
}

func TestINTERPOLATION26(t *testing.T) {
	// Ampersand - Standalone

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `  {{&string}}
`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"string":"---"}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`  ---
`, b.String())
	})
}

func TestINTERPOLATION27(t *testing.T) {
	// Interpolation With Padding

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `|{{ string }}|`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"string":"---"}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`|---|`, b.String())
	})
}

func TestINTERPOLATION28(t *testing.T) {
	// Triple Mustache With Padding

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `|{{{ string }}}|`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"string":"---"}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`|---|`, b.String())
	})
}

func TestINTERPOLATION29(t *testing.T) {
	// Ampersand With Padding

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", `|{{& string }}|`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"string":"---"}`), &data))
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", data)
		test.AreEqual(`|---|`, b.String())
	})
}
