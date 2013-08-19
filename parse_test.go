package mussed

import (
	"bytes"
	"html/template"
	"testing"
)

func TestTextParse(t *testing.T) {
	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.mustache", "<html></html>")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", nil)
		test.AreEqual("<html></html>", b.String())
	})
}

func TestParse(t *testing.T) {
	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.mustache", "<html>{{name}}</html>")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", map[string]interface{}{"name": "Andrew"})
		test.AreEqual("<html>Andrew</html>", b.String())
	})
}

func TestParseTwo(t *testing.T) {
	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.mustache", "<html>{{welcome}}{{name}}</html>")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree", map[string]interface{}{"name": "Andrew", "welcome": "Bonjour, "})
		test.AreEqual("<html>Bonjour, Andrew</html>", b.String())
	})
}

func TestSwitchDelim(t *testing.T) {
	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.mustache", "<html>{{name}}{{=<% %>=}}<title><%title%></title><%={{ }}=%><body>{{body}}</body></html>")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree",
			map[string]interface{}{"name": "Andrew", "title": "Wat", "body": "do you lift?"})
		test.AreEqual("<html>Andrew<title>Wat</title><body>do you lift?</body></html>", b.String())
	})
}

func TestIfSingle(t *testing.T) {
	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", "<html>{{#name}}<title>Title</title>{{/name}}</html>")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "tree",
			map[string]interface{}{"name": "Andrew"}))
		test.AreEqual("<html><title>Title</title></html>", b.String())
	})
}

func TestRanging(t *testing.T) {
	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", "<html>{{#name}}<title>Title</title>{{/name}}</html>")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "tree",
			map[string]interface{}{"name": "Andrew"})
		test.AreEqual("<html><title>Title</title></html>", b.String())

		b.Reset()
		t.ExecuteTemplate(b, "tree",
			map[string]interface{}{"name": []string{"Andrew", "John"}})
		test.AreEqual("<html><title>Title</title><title>Title</title></html>", b.String())

	})

}

func TestNot(t *testing.T) {
	within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		tree, err := Parse("test.mustache", "<html>{{^name}}<title>Title</title>{{/name}}</html>")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "tree", nil))
		test.AreEqual("<html><title>Title</title></html>", b.String())
	})
}
