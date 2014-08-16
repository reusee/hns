package nw

import "regexp"

type WalkFunc func(*Node)

func (n *Node) Walk(fn WalkFunc) {
	fn(n)
}

// util combinators

func Multi(funcs ...WalkFunc) WalkFunc {
	return func(node *Node) {
		for _, f := range funcs {
			f(node)
		}
	}
}

func Assign(p **Node) WalkFunc {
	return func(n *Node) {
		*p = n
	}
}

func Append(p *[]*Node) WalkFunc {
	return func(n *Node) {
		*p = append(*p, n)
	}
}

// scope combinators

func Descendant(predict WalkPredict, cont WalkFunc) WalkFunc {
	var f WalkFunc
	f = func(node *Node) {
		for _, child := range node.Children {
			if predict(child) {
				cont(child)
			} else {
				f(child)
			}
		}
	}
	return f
}

func AllDescendant(predict WalkPredict, cont WalkFunc) WalkFunc {
	var f WalkFunc
	f = func(node *Node) {
		for _, child := range node.Children {
			if predict(child) {
				cont(child)
			}
			f(child)
		}
	}
	return f
}

func Children(predict WalkPredict, cont WalkFunc) WalkFunc {
	return func(node *Node) {
		for _, child := range node.Children {
			if predict(child) {
				cont(child)
			}
		}
	}
}

func Current(predict WalkPredict, cont WalkFunc) WalkFunc {
	return func(node *Node) {
		if predict(node) {
			cont(node)
		}
	}
}

//TODO Siblings SiblingsBefore SiblingsAfter

// predicts

type WalkPredict func(*Node) bool

func TagEq(tag string) WalkPredict {
	return func(node *Node) bool {
		return node.Tag == tag
	}
}

func TagMatch(pattern string) WalkPredict {
	p := regexp.MustCompile(pattern)
	return func(node *Node) bool {
		return p.MatchString(node.Tag)
	}
}

func IdEq(id string) WalkPredict {
	return func(node *Node) bool {
		return node.Attr["id"] == id
	}
}

func IdMatch(pattern string) WalkPredict {
	p := regexp.MustCompile(pattern)
	return func(node *Node) bool {
		return p.MatchString(node.Attr["id"])
	}
}

func AttrEq(key, value string) WalkPredict {
	return func(node *Node) bool {
		return node.Attr[key] == value
	}
}

func AttrMatch(key, pattern string) WalkPredict {
	p := regexp.MustCompile(pattern)
	return func(node *Node) bool {
		return p.MatchString(node.Attr[key])
	}
}

func ClassEq(class string) WalkPredict {
	return func(node *Node) bool {
		for _, c := range node.Class {
			if c == class {
				return true
			}
		}
		return false
	}
}

func ClassMatch(pattern string) WalkPredict {
	p := regexp.MustCompile(pattern)
	return func(node *Node) bool {
		for _, c := range node.Class {
			if p.MatchString(c) {
				return true
			}
		}
		return false
	}
}
