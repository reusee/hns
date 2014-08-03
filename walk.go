package hns

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

type ScopeCombinator func(func(*WalkCtx, *Node) bool, WalkFunc) WalkFunc

func Descendant(predict func(*WalkCtx, *Node) bool, cont WalkFunc) WalkFunc {
	var f func(*WalkCtx, *Node)
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

func AllDescendant(predict func(*WalkCtx, *Node) bool, cont WalkFunc) WalkFunc {
	var f func(*WalkCtx, *Node)
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

func Children(predict func(*WalkCtx, *Node) bool, cont WalkFunc) WalkFunc {
	return func(ctx *WalkCtx, node *Node) {
		for _, child := range node.Children {
			if predict(ctx, child) {
				cont(ctx, child)
			}
		}
	}
}

// tag, class, id, attr

func TagEq(scope ScopeCombinator, tag string, cont WalkFunc) WalkFunc {
	return scope(func(_ *WalkCtx, node *Node) bool {
		return node.Tag == tag
	}, cont)
}

func IdEq(scope ScopeCombinator, id string, cont WalkFunc) WalkFunc {
	return scope(func(_ *WalkCtx, node *Node) bool {
		return node.Attr["id"] == id
	}, cont)
}

// helpers

func DescendantTagEq(tag string, cont WalkFunc) WalkFunc {
	return TagEq(Descendant, tag, cont)
}

func DescendantIdEq(id string, cont WalkFunc) WalkFunc {
	return IdEq(Descendant, id, cont)
}

func AllDescendantTagEq(tag string, cont WalkFunc) WalkFunc {
	return TagEq(AllDescendant, tag, cont)
}

// no need for AllDescendantIdEq

func ChildrenTagEq(tag string, cont WalkFunc) WalkFunc {
	return TagEq(Children, tag, cont)
}

func ChildrenIdEq(id string, cont WalkFunc) WalkFunc {
	return IdEq(Children, id, cont)
}
