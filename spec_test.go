package mussed

import "testing"

// These tests are modified from the mustache/spec repo
func TestDelimters(t *testing.T) {
	within(t, func(test *aTest) {
	})
}

var interpolationTests = `[
    {"name":"No Interpolation","data":{},"expected":"Hello from {Mustache}!\n","template":"Hello from {Mustache}!\n","desc":"Mustache-free templates should render as-is."},
    {"name":"Basic Interpolation","data":{"subject":"world"},"expected":"Hello, world!\n","template":"Hello, {{subject}}!\n","desc":"Unadorned tags should interpolate content into the template."},
    {"name":"HTML Escaping","data":{"forbidden":"& \" < >"},"expected":"These characters should be HTML escaped: &amp; &quot; &lt; &gt;\n","template":"These characters should be HTML escaped: {{forbidden}}\n","desc":"Basic interpolation should be HTML escaped."},
    {"name":"Triple Mustache","data":{"forbidden":"& \" < >"},"expected":"These characters should not be HTML escaped: & \" < >\n","template":"These characters should not be HTML escaped: {{{forbidden}}}\n","desc":"Triple mustaches should interpolate without HTML escaping."},
    {"name":"Ampersand","data":{"forbidden":"& \" < >"},"expected":"These characters should not be HTML escaped: & \" < >\n","template":"These characters should not be HTML escaped: {{&forbidden}}\n","desc":"Ampersand should interpolate without HTML escaping."},
    {"name":"Basic Integer Interpolation","data":{"mph":85},"expected":"\"85 miles an hour!\"","template":"\"{{mph}} miles an hour!\"","desc":"Integers should interpolate seamlessly."},
    {"name":"Triple Mustache Integer Interpolation","data":{"mph":85},"expected":"\"85 miles an hour!\"","template":"\"{{{mph}}} miles an hour!\"","desc":"Integers should interpolate seamlessly."},
    {"name":"Ampersand Integer Interpolation","data":{"mph":85},"expected":"\"85 miles an hour!\"","template":"\"{{&mph}} miles an hour!\"","desc":"Integers should interpolate seamlessly."},
    {"name":"Basic Decimal Interpolation","data":{"power":1.21},"expected":"\"1.21 jiggawatts!\"","template":"\"{{power}} jiggawatts!\"","desc":"Decimals should interpolate seamlessly with proper significance."},
    {"name":"Triple Mustache Decimal Interpolation","data":{"power":1.21},"expected":"\"1.21 jiggawatts!\"","template":"\"{{{power}}} jiggawatts!\"","desc":"Decimals should interpolate seamlessly with proper significance."},
    {"name":"Ampersand Decimal Interpolation","data":{"power":1.21},"expected":"\"1.21 jiggawatts!\"","template":"\"{{&power}} jiggawatts!\"","desc":"Decimals should interpolate seamlessly with proper significance."},
    {"name":"Basic Context Miss Interpolation","data":{},"expected":"I () be seen!","template":"I ({{cannot}}) be seen!","desc":"Failed context lookups should default to empty strings."},
    {"name":"Triple Mustache Context Miss Interpolation","data":{},"expected":"I () be seen!","template":"I ({{{cannot}}}) be seen!","desc":"Failed context lookups should default to empty strings."},
    {"name":"Ampersand Context Miss Interpolation","data":{},"expected":"I () be seen!","template":"I ({{&cannot}}) be seen!","desc":"Failed context lookups should default to empty strings."},
    {"name":"Dotted Names - Basic Interpolation","data":{"person":{"name":"Joe"}},"expected":"\"Joe\" == \"Joe\"","template":"\"{{person.name}}\" == \"{{#person}}{{name}}{{/person}}\"","desc":"Dotted names should be considered a form of shorthand for sections."},
    {"name":"Dotted Names - Triple Mustache Interpolation","data":{"person":{"name":"Joe"}},"expected":"\"Joe\" == \"Joe\"","template":"\"{{{person.name}}}\" == \"{{#person}}{{{name}}}{{/person}}\"","desc":"Dotted names should be considered a form of shorthand for sections."},
    {"name":"Dotted Names - Ampersand Interpolation","data":{"person":{"name":"Joe"}},"expected":"\"Joe\" == \"Joe\"","template":"\"{{&person.name}}\" == \"{{#person}}{{&name}}{{/person}}\"","desc":"Dotted names should be considered a form of shorthand for sections."},
    {"name":"Dotted Names - Arbitrary Depth","data":{"a":{"b":{"c":{"d":{"e":{"name":"Phil"}}}}}},"expected":"\"Phil\" == \"Phil\"","template":"\"{{a.b.c.d.e.name}}\" == \"Phil\"","desc":"Dotted names should be functional to any level of nesting."},
    {"name":"Dotted Names - Broken Chains","data":{"a":{}},"expected":"\"\" == \"\"","template":"\"{{a.b.c}}\" == \"\"","desc":"Any falsey value prior to the last part of the name should yield ''."},
    {"name":"Dotted Names - Broken Chain Resolution","data":{"a":{"b":{}},"c":{"name":"Jim"}},"expected":"\"\" == \"\"","template":"\"{{a.b.c.name}}\" == \"\"","desc":"Each part of a dotted name should resolve only against its parent."},
    {"name":"Dotted Names - Initial Resolution","data":{"a":{"b":{"c":{"d":{"e":{"name":"Phil"}}}}},"b":{"c":{"d":{"e":{"name":"Wrong"}}}}},"expected":"\"Phil\" == \"Phil\"","template":"\"{{#a}}{{b.c.d.e.name}}{{/a}}\" == \"Phil\"","desc":"The first part of a dotted name should resolve as any other name."},
    {"name":"Interpolation - Surrounding Whitespace","data":{"string":"---"},"expected":"| --- |","template":"| {{string}} |","desc":"Interpolation should not alter surrounding whitespace."},{"name":"Triple Mustache - Surrounding Whitespace","data":{"string":"---"},"expected":"| --- |","template":"| {{{string}}} |","desc":"Interpolation should not alter surrounding whitespace."},
    {"name":"Ampersand - Surrounding Whitespace","data":{"string":"---"},"expected":"| --- |","template":"| {{&string}} |","desc":"Interpolation should not alter surrounding whitespace."},
    {"name":"Interpolation - Standalone","data":{"string":"---"},"expected":"  ---\n","template":"  {{string}}\n","desc":"Standalone interpolation should not alter surrounding whitespace."},
    {"name":"Triple Mustache - Standalone","data":{"string":"---"},"expected":"  ---\n","template":"  {{{string}}}\n","desc":"Standalone interpolation should not alter surrounding whitespace."},
    {"name":"Ampersand - Standalone","data":{"string":"---"},"expected":"  ---\n","template":"  {{&string}}\n","desc":"Standalone interpolation should not alter surrounding whitespace."},
    {"name":"Interpolation With Padding","data":{"string":"---"},"expected":"|---|","template":"|{{ string }}|","desc":"Superfluous in-tag whitespace should be ignored."},
    {"name":"Triple Mustache With Padding","data":{"string":"---"},"expected":"|---|","template":"|{{{ string }}}|","desc":"Superfluous in-tag whitespace should be ignored."},
    {"name":"Ampersand With Padding","data":{"string":"---"},"expected":"|---|","template":"|{{& string }}|","desc":"Superfluous in-tag whitespace should be ignored."}]`
