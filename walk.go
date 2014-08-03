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

// utils

func Do(fn func(*Node)) WalkFunc {
	return func(ctx *WalkCtx, node *Node) {
		fn(node)
	}
}

func Return(ctx *WalkCtx, node *Node) {
	ctx.Return = append(ctx.Return, node)
}

// scope

type WalkPredict func(*WalkCtx, *Node) bool
type ScopeCombinator func(WalkPredict, WalkFunc) WalkFunc

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

// member: tag, id, attr, class

func TagEq(scope ScopeCombinator, tag string, cont WalkFunc) WalkFunc {
	return scope(func(_ *WalkCtx, node *Node) bool {
		return node.Tag == tag
	}, cont)
}

func TagMatch(scope ScopeCombinator, pattern string, cont WalkFunc) WalkFunc {
	p := regexp.MustCompile(pattern)
	return scope(func(_ *WalkCtx, node *Node) bool {
		return p.MatchString(node.Tag)
	}, cont)
}

func IdEq(scope ScopeCombinator, id string, cont WalkFunc) WalkFunc {
	return scope(func(_ *WalkCtx, node *Node) bool {
		return node.Attr["id"] == id
	}, cont)
}

func IdMatch(scope ScopeCombinator, pattern string, cont WalkFunc) WalkFunc {
	p := regexp.MustCompile(pattern)
	return scope(func(_ *WalkCtx, node *Node) bool {
		return p.MatchString(node.Attr["id"])
	}, cont)
}

func AttrEq(scope ScopeCombinator, key, value string, cont WalkFunc) WalkFunc {
	return scope(func(_ *WalkCtx, node *Node) bool {
		return node.Attr[key] == value
	}, cont)
}

func AttrMatch(scope ScopeCombinator, key, value string, cont WalkFunc) WalkFunc {
	p := regexp.MustCompile(value)
	return scope(func(_ *WalkCtx, node *Node) bool {
		return p.MatchString(node.Attr[key])
	}, cont)
}

func ClassEq(scope ScopeCombinator, class string, cont WalkFunc) WalkFunc {
	return scope(func(_ *WalkCtx, node *Node) bool {
		for _, c := range node.Class {
			if c == class {
				return true
			}
		}
		return false
	}, cont)
}

// scope X member

func DescendantTagEq(tag string, cont WalkFunc) WalkFunc {
	return TagEq(Descendant, tag, cont)
}

func AllDescendantTagEq(tag string, cont WalkFunc) WalkFunc {
	return TagEq(AllDescendant, tag, cont)
}

func ChildrenTagEq(tag string, cont WalkFunc) WalkFunc {
	return TagEq(Children, tag, cont)
}

func DescendantTagMatch(pattern string, cont WalkFunc) WalkFunc {
	return TagMatch(Descendant, pattern, cont)
}

func AllDescendantTagMatch(pattern string, cont WalkFunc) WalkFunc {
	return TagMatch(AllDescendant, pattern, cont)
}

func ChildrenTagMatch(pattern string, cont WalkFunc) WalkFunc {
	return TagMatch(Children, pattern, cont)
}

func DescendantIdEq(id string, cont WalkFunc) WalkFunc {
	return IdEq(Descendant, id, cont)
}

func AllDescendantIdEq(id string, cont WalkFunc) WalkFunc {
	return IdEq(AllDescendant, id, cont)
}

func ChildrenIdEq(id string, cont WalkFunc) WalkFunc {
	return IdEq(Children, id, cont)
}

func DescendantIdMatch(pattern string, cont WalkFunc) WalkFunc {
	return IdMatch(Descendant, pattern, cont)
}

func AllDescendantIdMatch(pattern string, cont WalkFunc) WalkFunc {
	return IdMatch(AllDescendant, pattern, cont)
}

func ChildrenIdMatch(pattern string, cont WalkFunc) WalkFunc {
	return IdMatch(Children, pattern, cont)
}

func DescendantAttrEq(key, value string, cont WalkFunc) WalkFunc {
	return AttrEq(Descendant, key, value, cont)
}

func AllDescendantAttrEq(key, value string, cont WalkFunc) WalkFunc {
	return AttrEq(AllDescendant, key, value, cont)
}

func ChildrenAttrEq(key, value string, cont WalkFunc) WalkFunc {
	return AttrEq(Children, key, value, cont)
}

func DescendantAttrMatch(key, value string, cont WalkFunc) WalkFunc {
	return AttrMatch(Descendant, key, value, cont)
}

func AllDescendantAttrMatch(key, value string, cont WalkFunc) WalkFunc {
	return AttrMatch(AllDescendant, key, value, cont)
}

func ChildrenAttrMatch(key, value string, cont WalkFunc) WalkFunc {
	return AttrMatch(Children, key, value, cont)
}

func DescendantClassEq(class string, cont WalkFunc) WalkFunc {
	return ClassEq(Descendant, class, cont)
}

func AllDescendantClassEq(class string, cont WalkFunc) WalkFunc {
	return ClassEq(AllDescendant, class, cont)
}

func ChildrenClassEq(class string, cont WalkFunc) WalkFunc {
	return ClassEq(Children, class, cont)
}
