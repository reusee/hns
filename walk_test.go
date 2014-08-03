package hns

import "testing"

func TestWalk(t *testing.T) {
	root, err := ParseString(`
<div>
	<p>foo</p>
	<p>bar</p>
	<p>baz</p>
	<a>foo</a>
	<div id="foo">
		<a>bar</a>
	</div>
</div>
	`)
	if err != nil {
		t.Fatal(err)
	}

	// DescendantTagEq
	var res []*Node
	root.Walk(
		DescendantTagEq("div",
			DescendantTagEq("p",
				Do(func(node *Node) {
					res = append(res, node)
				}))))
	if len(res) != 3 {
		t.Fatalf("DescendantTagEq not match")
	}
	if res[0].Text != "foo" || res[1].Text != "bar" || res[2].Text != "baz" {
		t.Fatalf("DescendantTagEq result error")
	}

	// ChildrenTagEq
	res = res[:0]
	root.Walk(
		ChildrenTagEq("div",
			ChildrenTagEq("a",
				Do(func(n *Node) {
					res = append(res, n)
				}))))
	if len(res) != 1 {
		t.Fatalf("ChildrenTagEq result not match")
	}
	if res[0].Text != "foo" {
		t.Fatalf("ChildrenTagEq result error")
	}

	// Descendant
	res = res[:0]
	root.Walk(Descendant(func(_ *WalkCtx, node *Node) bool {
		return node.Tag == "a" && node.Text == "bar"
	}, Do(func(node *Node) {
		res = append(res, node)
	})))
	if len(res) != 1 {
		t.Fatalf("Descendant result not match")
	}
	if res[0].Parent.Tag != "div" {
		t.Fatalf("Descendant result error")
	}

	// DescendantIdEq
	res = res[:0]
	root.Walk(DescendantIdEq("foo", Do(func(node *Node) {
		res = append(res, node)
	})))
	if len(res) != 1 {
		t.Fatalf("DescendantIdEq result not match")
	}
	if res[0].Tag != "div" || len(res[0].Children) != 1 {
		t.Fatalf("DescendantIdEq result error")
	}
}
