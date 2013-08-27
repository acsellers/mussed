package mussed

import (
	"strings"
	"text/template/parse"
)

var (
	LeftDelim        = "{{"
	RightDelim       = "}}"
	LeftEscapeDelim  = "{{{"
	RightEscapeDelim = "}}}"
)

func Parse(templateName, templateContent string) (map[string]*parse.Tree, error) {
	i := strings.Index(templateName, ".mustache")
	name := templateName[:i] + templateName[i+len(".mustache"):]

	proto := &protoTree{
		source:     templateContent,
		localRight: RightDelim,
		localLeft:  LeftDelim,
		tree: &parse.Tree{
			Name:      name,
			ParseName: templateName,
			Root: &parse.ListNode{
				NodeType: parse.NodeList,
			},
		},
	}
	proto.list = proto.tree.Root
	proto.parse()

	return proto.templates(), proto.err
}
