package mussed

import (
	"strings"
	"text/template/parse"
)

var mangleNum int

type protoTree struct {
	source     string
	tree       *parse.Tree
	childTrees []*parse.Tree
	list       *parse.ListNode
	stack      []*parse.ListNode
	err        error
	localLeft  string
	localRight string
}

func (pt *protoTree) pop() *parse.ListNode {
	if len(pt.stack) == 0 {
		return pt.tree.Root
	}
	ln := pt.stack[len(pt.stack)-1]
	pt.stack = pt.stack[:len(pt.stack)-2]
	return ln
}
func (pt *protoTree) push(ln *parse.ListNode) {
	pt.stack = append(pt.stack, ln)
}

type stash struct {
	tree    *protoTree
	content string
}

func (s *stash) needsMoreText() bool {
	normalOpen := strings.Index(s.content, s.tree.localLeft)
	normalUnescape := strings.Index(s.content, LeftEscapeDelim)

	if normalUnescape >= 0 && normalUnescape < normalOpen {
		closeIndex := strings.Index(s.content, RightEscapeDelim)
		return !(closeIndex >= 0 && closeIndex > normalUnescape)
	}
	if normalOpen >= 0 {
		closeIndex := strings.Index(s.content, s.tree.localRight)
		return !(closeIndex >= 0 && closeIndex > normalOpen)
	}

	return false
}

func (s *stash) Append(t string) {
	s.content += t
}
func (s *stash) hasAction() bool {
	return strings.Contains(s.content, s.tree.localLeft) ||
		strings.Contains(s.content, LeftEscapeDelim)
}

func (s *stash) pullToAction() (string, string) {
	var text, action string
	loc, abnormal := s.nextActionLocation()
	text = s.content[:loc]
	s.content = s.content[loc:]
	if abnormal {
		action = s.content[:len(LeftEscapeDelim)]
		closeLocation := strings.Index(s.content, RightEscapeDelim)
		action += s.content[:closeLocation]
	} else {
		action = s.content[:len(s.tree.localLeft)]
		closeLocation := strings.Index(s.content, s.tree.localRight)
		action += s.content[:closeLocation]
	}
	return text, action
}

func (s *stash) nextActionLocation() (int, bool) {
	normalOpen := strings.Index(s.content, s.tree.localLeft)
	normalUnescape := strings.Index(s.content, LeftEscapeDelim)

	if normalUnescape >= 0 && normalUnescape < normalOpen {
		return normalUnescape, true
	}
	if normalOpen >= 0 {
		return normalOpen, false
	}
	return -1, false
}
