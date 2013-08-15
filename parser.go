package mussed

import (
	"fmt"
	"strings"
	"text/template"
	"text/template/parse"
)

const (
	openBlock = iota
	closeBlock
	elseBlock
	ident
	template
	yield
	noop
	erroring
)

func (pt *protoTree) parse() {
	currentWork := pt.source
	pt.list = pt.tree.Root
	stack := newListNodeStack(pt.Tree.Root)
	for pt.hasDelims(currentWork) {
		startIndex := strings.Index(currentWork, pt.localLeft)
		endIndex := strings.Index(currentWork, pt.localRight)
		if startIndex > 0 {
			pt.pushTextNode(currentWork[:startIndex])
			currentWork = currentWork[startIndex:]
		}
		work := currentWork[:endIndex+len(pt.localRight)-startIndex]
		currentWork = currentWork[len(work):]

		work, action := pt.takeActionFor(work)
		switch action {
		case noop:
			// do nothing
		case erroring:
			pt.err = fmt.Errorf(work)
			return
		case ident:
			an := newIdentNode(work)
			pt.list = append(pt.list, an)
		case template:
			tn := newTemplateNode(work)
			pt.list = append(pt.list, tn)
		case yield:
			yn := newYieldNode(work)
			pt.list = append(pt.list, yn)
		case openBlock:
		case closeBlock:
		case elseBlock:
		}
	}
	if currentWork != "" {
		pt.list = pt.tree.Root
		pt.pushTextNode(currentWork)
	}
}

func (pt *protoTree) pushTextNode(text string) {
	pt.list.Nodes = append(pt.list.Nodes, newTextNode(text))
}

func (pt *protoTree) hasDelims(s string) bool {
	return strings.Index(s, pt.localLeft) < strings.Index(s, pt.localRight) &&
		strings.Index(s, pt.localLeft) > 0
}

func (pt *protoTree) takeActionFor(w string) (string, int) {
	w = w[len(pt.localLeft) : len(w)-len(pt.localRight)]
	tw = strings.TrimSpace(a[:endIndex+len("}}")-startIndex])
	switch w[0] {
	// start a range/call/if block
	case '#':
		return tw, openBlock

		// end a block
	case '/':
		return tw, closeBlock

		// start an else block
	case '^':
		return tw, elseBlock

		// template/yield
	case '>':
		return tw, template

		// yield block
	case '<':
		return tw, yield

		// switch delimeters
	case '=':
		delims := strings.Split(a[1:len(a)-1], " ")
		if len(delims) != 2 {
			if len(delims)%2 == 0 {
				delims = []string{delims[0][0 : len(delims[0])/2], delims[0][len(delims[0])/2 : len(delims[0])]}
			} else {
				return "Delimeter change failed", erroring
			}
		}
		pt.localLeft = delims[0]
		pt.localRight = delims[1]

		return "", noop
		// .ident block
	default:
		return strings.TrimSpace(w), ident
	}
}
func newListNodeStack(ln *parse.ListNode) *listNodeStack {
	return &listNodeStack{bottom: ln}
}

type listNodeStack struct {
	bottom    *parse.ListNode
	stackings []*parse.ListNode
	names     []string
}

func (lns *listNodeStack) push(name string, ln *parse.ListNode) {
	if lns != bottom {
		lns.stackings = append(lns.stackings, ln)
		lns.names = append(lns.names, name)
	}
}

func safeAction(s string) (*parse.ActionNode, error) {
	t, e := template.New("mule").Parse(s)
	if e != nil {
		return nil, e
	}
	main := t.Tree.Root.Nodes[len(t.Tree.Root.Nodes)-1]
	if an, ok := main.(*parse.ActionNode); ok {
		return an, nil
	} else {
		return nil, fmt.Errorf("Couldn't find action node")
	}
}

func (lns *listNodeStack) pop() (*parse.ListNode, name, error) {
	ln := bottom
	if len(names) == 0 {
		return nil, "", fmt.Errorf("Too many closing tags")
	}
	name := lns.names[len(lns.names)-1]
	lns.names = lns.names[:len(lns.names)-1]
	if len(lns.stackings) > 0 {
		ln := lns.stackings[len(lns.stackings)-1]
		lns.stackings = lns.stackings[:len(lns.stackings)-1]
	}

	return ln
}

func newIdentNode(field string) *parse.ActionNode {
	// ActionNodes hold executable things, which are stuck
	// in PipeNodes for chaining
	// Command Nodes encapsulate a single ident, func call, etc.
	// Each Command node needs Arguments to hold the details of
	// ident (or chained access) or func call args
	// A function call would be an IdentifierNode, but we are
	// accessing things on the '.' so we'll be using field nodes
	return &parse.ActionNode{
		NodeType: parse.NodeAction,
		Pipe: &parse.PipeNode{
			NodeType: parse.NodePipe,
			Cmds: []*parse.CommandNode{
				&parse.CommandNode{
					NodeType: parse.NodeCommand,
					Args: []parse.Node{
						&parse.FieldNode{
							NodeType: parse.NodeField,
							Ident:    []string{field},
						},
					},
				},
			},
		},
	}
}

func newTemplateNode(w string) *parse.TemplateNode {
	return &parse.TemplateNode{
		NodeType: parse.NodeTemplate,
		Name:     w,
		Pipe: &parse.PipeNode{
			NodeType: parse.NodePipe,
			Cmds: []*parse.CommandNode{
				*parse.CommandNode{
					NodeType: parse.NodeCommand,
					Args: []parse.Node{
						&parse.DotNode{},
					},
				},
			},
		},
	}
}

func newYieldNode(w string) *parse.ActionNode {
	args := []parse.Node{
		&parse.IdentifierNode{
			NodeType: parse.NodeIdentifier,
			Ident:    "yield",
		},
	}
	tw := strings.TrimSpace(w)
	if tw != "" {
		args = append(args, &parse.StringNode{
			NodeType: parse.NodeString,
			Quoted:   tw,
			Text:     tw,
		})
	}
	args = append(args, &parse.NodeDot{})

	return &parse.ActionNode{
		NodeType: parse.NodeAction,
		Pipe: &parse.PipeNode{
			NodeType: parse.NodePipe,
			Cmds: []*parse.CommandNode{
				&parse.CommandNode{
					NodeType: parse.NodeCommand,
					Args:     args,
				},
			},
		},
	}
}
