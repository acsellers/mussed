package mussed

import (
	"text/template/parse"
)

var (
	LeftDelim  = "{{"
	RightDelim = "}}"
)

func Parse(templateName, templateContent string) (map[string]*parse.Tree, error) {
	proto := &protoTree{source: text, localRight: RightDelim, localLeft: LeftDelim}
	proto.parse()

	i := strings.Index(name, ".mustache")

	return map[string]*parse.Tree{
		name[:i] + name[i+5:]: proto.tree,
	}, proto.err
}

type protoTree struct {
	source     string
	tokenList  []token
	tree       *parse.Tree
	list       *parse.ListNode
	err        error
	localLeft  string
	localRight string
}
