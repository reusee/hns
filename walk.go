package hns

import "regexp"

type WalkFunc func(*WalkCtx, *Node)

type WalkCtx struct {
	Return []*Node
}

func (n *Node) Walk(fn WalkFunc) *WalkCtx {
	ctx := &WalkCtx{}
	fn(ctx, n)
	return ctx
}

// util combinators

func Do(fn func(*Node)) WalkFunc {
	return func(ctx *WalkCtx, node *Node) {
		fn(node)
	}
}

func Return(ctx *WalkCtx, node *Node) {
	ctx.Return = append(ctx.Return, node)
}

// scope combinators

func Descendant(predict WalkPredict, cont WalkFunc) WalkFunc {
	var f WalkFunc
	f = func(ctx *WalkCtx, node *Node) {
		for _, child := range node.Children {
			if predict(ctx, child) {
				cont(ctx, child)
			} else {
				f(ctx, child)
			}
		}
	}
	return f
}

func AllDescendant(predict WalkPredict, cont WalkFunc) WalkFunc {
	var f WalkFunc
	f = func(ctx *WalkCtx, node *Node) {
		for _, child := range node.Children {
			if predict(ctx, child) {
				cont(ctx, child)
			}
			f(ctx, child)
		}
	}
	return f
}

func Children(predict WalkPredict, cont WalkFunc) WalkFunc {
	return func(ctx *WalkCtx, node *Node) {
		for _, child := range node.Children {
			if predict(ctx, child) {
				cont(ctx, child)
			}
		}
	}
}

// predicts

type WalkPredict func(*WalkCtx, *Node) bool

func TagEq(tag string) WalkPredict {
	return func(_ *WalkCtx, node *Node) bool {
		return node.Tag == tag
	}
}

func TagMatch(pattern string) WalkPredict {
	p := regexp.MustCompile(pattern)
	return func(_ *WalkCtx, node *Node) bool {
		return p.MatchString(node.Tag)
	}
}

func IdEq(id string) WalkPredict {
	return func(_ *WalkCtx, node *Node) bool {
		return node.Attr["id"] == id
	}
}

func IdMatch(pattern string) WalkPredict {
	p := regexp.MustCompile(pattern)
	return func(_ *WalkCtx, node *Node) bool {
		return p.MatchString(node.Attr["id"])
	}
}

func AttrEq(key, value string) WalkPredict {
	return func(_ *WalkCtx, node *Node) bool {
		return node.Attr[key] == value
	}
}

func AttrMatch(key, value string) WalkPredict {
	p := regexp.MustCompile(value)
	return func(_ *WalkCtx, node *Node) bool {
		return p.MatchString(node.Attr[key])
	}
}

func ClassEq(class string) WalkPredict {
	return func(_ *WalkCtx, node *Node) bool {
		for _, c := range node.Class {
			if c == class {
				return true
			}
		}
		return false
	}
}
