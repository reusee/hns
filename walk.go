package hns

import "regexp"

type WalkFunc func(*WalkCtx, *Node)

type WalkCtx struct {
	Return []*Node
}

func (n *Node) Walk(fn WalkFunc) []*Node {
	ctx := &WalkCtx{}
	fn(ctx, n)
	return ctx.Return
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

// member: tag, class, id, attr

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

// scope X member

func DescendantTagEq(tag string, cont WalkFunc) WalkFunc {
	return TagEq(Descendant, tag, cont)
}

func DescendantIdEq(id string, cont WalkFunc) WalkFunc {
	return IdEq(Descendant, id, cont)
}

func AllDescendantTagEq(tag string, cont WalkFunc) WalkFunc {
	return TagEq(AllDescendant, tag, cont)
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

// no need for AllDescendantIdEq

func ChildrenTagEq(tag string, cont WalkFunc) WalkFunc {
	return TagEq(Children, tag, cont)
}

func ChildrenIdEq(id string, cont WalkFunc) WalkFunc {
	return IdEq(Children, id, cont)
}
