package regtree

import (
	"regexp"

	"gopkg.in/baa.v1"
)

// Tree provlider router store for baa with regexp and radix tree
type Tree struct {
	static    bool
	alpha     byte
	pattern   string
	handlers  []baa.HandlerFunc
	params    []string
	format    []byte
	re        *regexp.Regexp
	schildren []*Tree
	rchildren []*Tree
	parent    *Tree
}

// NewTree create a new tree route node
func NewTree(pattern string, handlers []baa.HandlerFunc) *Tree {
	if pattern == "" {
		panic("tree.new: pattern can be empty")
	}
	return &Tree{
		static:   true,
		alpha:    pattern[0],
		pattern:  pattern,
		format:   []byte(pattern),
		handlers: handlers,
	}
}

// Get returns matched node and param values for key
func (t *Tree) Get(pattern string, c *baa.Context) []baa.HandlerFunc {
	// regexp rule
	if !t.static {
		matches := t.re.FindStringSubmatchIndex(pattern)
		if len(matches) == (len(t.params)+1)*2 && matches[0] == 0 {
			for j := range t.params {
				c.SetParam(t.params[j], pattern[matches[(j+1)*2]:matches[(j+1)*2+1]])
			}
			return t.handlers
		}
	}
	// static rule
	matched := 0
	for ; matched < len(pattern) && matched < len(t.pattern) && pattern[matched] == t.pattern[matched]; matched++ {
	}
	// no prefix
	if matched != len(t.pattern) {
		return nil
	}
	// found
	if matched == len(pattern) {
		if t.handlers != nil {
			return t.handlers
		}
	}
	// node is prefix
	pattern = pattern[matched:]
	// first, static rule
	if len(pattern) > 0 {
		if snode := t.findChild(pattern[0]); snode != nil {
			if h := snode.Get(pattern, c); h != nil {
				return h
			}
		}
	}

	// then, regexp rule
	for i := range t.rchildren {
		matches := t.rchildren[i].re.FindStringSubmatchIndex(pattern)
		if len(matches) != (len(t.rchildren[i].params)+1)*2 || matches[0] > 0 {
			continue
		}
		for j := range t.rchildren[i].params {
			c.SetParam(t.rchildren[i].params[j], pattern[matches[(j+1)*2]:matches[(j+1)*2+1]])
		}
		return t.rchildren[i].handlers
	}

	return nil
}

// Add return new node with key and val
func (t *Tree) Add(pattern string, handlers []baa.HandlerFunc) *Tree {
	// find the common prefix
	matched := 0
	for ; matched < len(pattern) && matched < len(t.pattern) && pattern[matched] == t.pattern[matched]; matched++ {
	}

	// no prefix
	if matched == 0 {
		return nil
	}

	if matched == len(t.pattern) {
		// the node pattern is the same as the pattern: make the current node as data node
		if matched == len(pattern) {
			if handlers != nil {
				if t.handlers != nil {
					panic("the route is be exists: " + t.String())
				}
				t.handlers = handlers
			}
			return t
		}

		// the node pattern is a prefix of the pattern: create a child node
		pattern = pattern[matched:]
		for _, child := range t.schildren {
			if node := child.Add(pattern, handlers); node != nil {
				return node
			}
		}

		// no child match, to be a new child
		return t.addChild(pattern, handlers)
	}

	// the pattern is a prefix of node pattern: create a new node instead of child
	if matched == len(pattern) {
		node := NewTree(t.pattern[matched:], t.handlers)
		node.schildren = t.schildren
		node.rchildren = t.rchildren
		node.parent = t
		t.pattern = pattern
		t.format = []byte(t.pattern)
		t.handlers = handlers
		t.schildren = []*Tree{node}
		t.rchildren = nil
		return t
	}

	// the node pattern shares a partial prefix with the key: split the node pattern
	node := NewTree(t.pattern[matched:], t.handlers)
	node.schildren = t.schildren
	node.rchildren = t.rchildren
	node.parent = t
	t.pattern = pattern[:matched]
	t.format = []byte(t.pattern)
	t.handlers = nil
	t.schildren = nil
	t.rchildren = nil
	t.schildren = append(t.schildren, node)
	return t.addChild(pattern[matched:], handlers)
}

func (t *Tree) addChild(pattern string, handlers []baa.HandlerFunc) *Tree {
	// check it is a static route child or not
	var staticPattern, param, rule string
	var params []string
	var newPattern, format []byte
	var i, j, k int
	for i = 0; i < len(pattern); i++ {
		if pattern[i] == '*' {
			// set static prefix
			if len(staticPattern) == 0 && len(params) == 0 && i > 0 {
				staticPattern = pattern[:i]
			}
			rule = "(^.*)"
			param = ""
			newPattern = append(newPattern, rule...)
			format = append(format, "%v"...)
			params = append(params, param)
			continue
		}
		if pattern[i] == ':' {
			for j = i + 1; j < len(pattern) && baa.IsParamChar(pattern[j]); j++ {
			}
			// set static prefix
			if len(staticPattern) == 0 && len(params) == 0 && i > 0 {
				staticPattern = pattern[:i]
			}
			param = pattern[i+1 : j]
			i = j - 1
			// check regexp rule
			rule = ""
			if j < len(pattern) && pattern[j] == '(' {
				for k = j + 1; k < len(pattern) && pattern[k] != ')'; k++ {
				}
				rule = pattern[j+1 : k]
				i = k
			}
			if rule == "" {
				rule = "([^\\/]+)"
			} else if rule == "int" {
				rule = "([\\d]+)"
			} else if rule == "string" {
				rule = "([\\w]+)"
			} else {
				rule = "(" + rule + ")"
			}
			newPattern = append(newPattern, rule...)
			format = append(format, "%v"...)
			params = append(params, param)
			continue
		}
		newPattern = append(newPattern, pattern[i])
		format = append(format, pattern[i])
	}

	var reNode, staticNode *Tree
	var err error
	if len(params) > 0 {
		// key has regexp rule, new regexp rule
		reNode = NewTree(string(newPattern[len(staticPattern):]), handlers)
		reNode.static = false
		reNode.params = params
		reNode.format = format[len(staticPattern):]
		reNode.re, err = regexp.Compile(reNode.pattern + "$")
		if err != nil {
			panic("tree.addChild: " + err.Error())
		}
		// set pattern with static prefix
		pattern = staticPattern
	}

	if len(pattern) > 0 {
		// key has static rule
		staticNode = NewTree(pattern, nil)
		staticNode.parent = t
		if reNode != nil {
			reNode.parent = staticNode
			staticNode.rchildren = append(staticNode.rchildren, reNode)
			t.schildren = append(t.schildren, staticNode)
			return reNode
		}
		staticNode.handlers = handlers
		t.schildren = append(t.schildren, staticNode)
		return staticNode
	}

	// key has regexp rule without static rule
	reNode.parent = t
	for _, child := range t.rchildren {
		if child.pattern == reNode.pattern {
			panic("the route is be exists: " + child.String())
		}
	}
	t.rchildren = append(t.rchildren, reNode)
	return reNode
}

// findChild find a match node from tree
func (t *Tree) findChild(b byte) *Tree {
	var i int
	var j = len(t.schildren)
	for ; i < j; i++ {
		if t.schildren[i].alpha == b {
			return t.schildren[i]
		}
	}
	return nil
}

// String return full key
func (t *Tree) String() string {
	s := t.pattern
	if t.parent != nil {
		s = t.parent.String() + s
	}
	return s
}

// formatStr return parsed format string
func (t *Tree) formatStr() []byte {
	var s []byte
	s = append(s, t.format...)
	if t.parent != nil {
		t := t.parent.formatStr()
		t = append(t, s...)
		s = t
	}
	return s
}
