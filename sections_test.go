package mussed

import (
	"bytes"
	"encoding/json"
	"html/template"
	"testing"
)

/*
Section tags and End Section tags are used in combination to wrap a section
of the template for iteration

These tags' content MUST be a non-whitespace character sequence NOT
containing the current closing delimiter; each Section tag MUST be followed
by an End Section tag with the same content within the same section.

This tag's content names the data to replace the tag.  Name resolution is as
follows:
  1) Split the name on periods; the first part is the name to resolve, any
  remaining parts should be retained.
  2) Walk the context stack from top to bottom, finding the first context
  that is a) a hash containing the name as a key OR b) an object responding
  to a method with the given name.
  3) If the context is a hash, the data is the value associated with the
  name.
  4) If the context is an object and the method with the given name has an
  arity of 1, the method SHOULD be called with a String containing the
  unprocessed contents of the sections; the data is the value returned.
  5) Otherwise, the data is the value returned by calling the method with
  the given name.
  6) If any name parts were retained in step 1, each should be resolved
  against a context stack containing only the result from the former
  resolution.  If any part fails resolution, the result should be considered
  falsey, and should interpolate as the empty string.
If the data is not of a list type, it is coerced into a list as follows: if
the data is truthy (e.g. `!!data == true`), use a single-element list
containing the data, otherwise use an empty list.

For each element in the data list, the element MUST be pushed onto the
context stack, the section MUST be rendered, and the element MUST be popped
off the context stack.

Section and End Section tags SHOULD be treated as standalone when
appropriate.

*/

func TestSECTIONS0(t *testing.T) {
	// Truthy

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `"{{#boolean}}This should be rendered.{{/boolean}}"`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"This should be rendered."`, b.String())
	})
}

func TestSECTIONS1(t *testing.T) {
	// Falsey

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `"{{#boolean}}This should not be rendered.{{/boolean}}"`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":false}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`""`, b.String())
	})
}

func TestSECTIONS2(t *testing.T) {
	// Context

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `"{{#context}}Hi {{name}}.{{/context}}"`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"context":{"name":"Joe"}}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"Hi Joe."`, b.String())
	})
}

func TestSECTIONS3(t *testing.T) {
	// Deeply Nested Contexts

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `{{#a}}
{{one}}
{{#b}}
{{one}}{{two}}{{one}}
{{#c}}
{{one}}{{two}}{{three}}{{two}}{{one}}
{{#d}}
{{one}}{{two}}{{three}}{{four}}{{three}}{{two}}{{one}}
{{#e}}
{{one}}{{two}}{{three}}{{four}}{{five}}{{four}}{{three}}{{two}}{{one}}
{{/e}}
{{one}}{{two}}{{three}}{{four}}{{three}}{{two}}{{one}}
{{/d}}
{{one}}{{two}}{{three}}{{two}}{{one}}
{{/c}}
{{one}}{{two}}{{one}}
{{/b}}
{{one}}
{{/a}}
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"a":{"one":1},"b":{"two":2},"c":{"three":3},"d":{"four":4},"e":{"five":5}}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`1
121
12321
1234321
123454321
1234321
12321
121
1
`, b.String())
	})
}

func TestSECTIONS4(t *testing.T) {
	// List

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `"{{#list}}{{item}}{{/list}}"`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"list":[{"item":1},{"item":2},{"item":3}]}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"123"`, b.String())
	})
}

func TestSECTIONS5(t *testing.T) {
	// Empty List

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `"{{#list}}Yay lists!{{/list}}"`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"list":[]}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`""`, b.String())
	})
}

func TestSECTIONS6(t *testing.T) {
	// Doubled

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `{{#bool}}
* first
{{/bool}}
* {{two}}
{{#bool}}
* third
{{/bool}}
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"bool":true,"two":"second"}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`* first
* second
* third
`, b.String())
	})
}

func TestSECTIONS7(t *testing.T) {
	// Nested (Truthy)

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `| A {{#bool}}B {{#bool}}C{{/bool}} D{{/bool}} E |`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"bool":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`| A B C D E |`, b.String())
	})
}

func TestSECTIONS8(t *testing.T) {
	// Nested (Falsey)

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `| A {{#bool}}B {{#bool}}C{{/bool}} D{{/bool}} E |`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"bool":false}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`| A  E |`, b.String())
	})
}

func TestSECTIONS9(t *testing.T) {
	// Context Misses

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `[{{#missing}}Found key 'missing'!{{/missing}}]`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`[]`, b.String())
	})
}

func TestSECTIONS10(t *testing.T) {
	// Implicit Iterator - String

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `"{{#list}}({{.}}){{/list}}"`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"list":["a","b","c","d","e"]}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"(a)(b)(c)(d)(e)"`, b.String())
	})
}

func TestSECTIONS11(t *testing.T) {
	// Implicit Iterator - Integer

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `"{{#list}}({{.}}){{/list}}"`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"list":[1,2,3,4,5]}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"(1)(2)(3)(4)(5)"`, b.String())
	})
}

func TestSECTIONS12(t *testing.T) {
	// Implicit Iterator - Decimal

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `"{{#list}}({{.}}){{/list}}"`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"list":[1.1,2.2,3.3,4.4,5.5]}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"(1.1)(2.2)(3.3)(4.4)(5.5)"`, b.String())
	})
}

func TestSECTIONS13(t *testing.T) {
	// Dotted Names - Truthy

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `"{{#a.b.c}}Here{{/a.b.c}}" == "Here"`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"a":{"b":{"c":true}}}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"Here" == "Here"`, b.String())
	})
}

func TestSECTIONS14(t *testing.T) {
	// Dotted Names - Falsey

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `"{{#a.b.c}}Here{{/a.b.c}}" == ""`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"a":{"b":{"c":false}}}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"" == ""`, b.String())
	})
}

func TestSECTIONS15(t *testing.T) {
	// Dotted Names - Broken Chains

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `"{{#a.b.c}}Here{{/a.b.c}}" == ""`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"a":{}}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"" == ""`, b.String())
	})
}

func TestSECTIONS16(t *testing.T) {
	// Surrounding Whitespace

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", ` | {{#boolean}}	|	{{/boolean}} | 
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(` | 	|	 | 
`, b.String())
	})
}

func TestSECTIONS17(t *testing.T) {
	// Internal Whitespace

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", ` | {{#boolean}} {{! Important Whitespace }}
 {{/boolean}} | 
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(` |  
  | 
`, b.String())
	})
}

func TestSECTIONS18(t *testing.T) {
	// Indented Inline Sections

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", ` {{#boolean}}YES{{/boolean}}
 {{#boolean}}GOOD{{/boolean}}
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(` YES
 GOOD
`, b.String())
	})
}

func TestSECTIONS19(t *testing.T) {
	// Standalone Lines

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `| This Is
{{#boolean}}
|
{{/boolean}}
| A Line
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`| This Is
|
| A Line
`, b.String())
	})
}

func TestSECTIONS20(t *testing.T) {
	// Indented Standalone Lines

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `| This Is
  {{#boolean}}
|
  {{/boolean}}
| A Line
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`| This Is
|
| A Line
`, b.String())
	})
}

func TestSECTIONS21(t *testing.T) {
	// Standalone Line Endings

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `|
{{#boolean}}
{{/boolean}}
|`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`|
|`, b.String())
	})
}

func TestSECTIONS22(t *testing.T) {
	// Standalone Without Previous Line

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `  {{#boolean}}
#{{/boolean}}
/`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`#
/`, b.String())
	})
}

func TestSECTIONS23(t *testing.T) {
	// Standalone Without Newline

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `#{{#boolean}}
/
  {{/boolean}}`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`#
/
`, b.String())
	})
}

func TestSECTIONS24(t *testing.T) {
	// Padding

	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache", `|{{# boolean }}={{/ boolean }}|`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`|=|`, b.String())
	})
}
