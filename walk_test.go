package hns

import "testing"

func TestWalk(t *testing.T) {
	root, err := ParseString(`
<div>
	<p>foo</p>
	<p>bar</p>
	<p>baz</p>
	<a>foo</a>
	<div>
		<a>bar</a>
	</div>
</div>
	`)
	if err != nil {
		t.Fatal(err)
	}

	var res []*Node
	root.Walk(
		DescendantTagEq("div",
			DescendantTagEq("p",
				Do(func(node *Node) {
					res = append(res, node)
				}))))
	if len(res) != 3 {
		t.Fatalf("p result not match")
	}
	if res[0].Text != "foo" || res[1].Text != "bar" || res[2].Text != "baz" {
		t.Fatalf("p result error")
	}

	res = res[:0]
	root.Walk(
		ChildrenTagEq("div",
			ChildrenTagEq("a",
				Do(func(n *Node) {
					res = append(res, n)
				}))))
	if len(res) != 1 {
		t.Fatalf("a result not match")
	}
	if res[0].Text != "foo" {
		t.Fatalf("a result error")
	}

	res = res[:0]
	root.Walk(Descendant(func(_ *WalkCtx, node *Node) bool {
		return node.Tag == "a" && node.Text == "bar"
	}, Do(func(node *Node) {
		res = append(res, node)
	})))
	if len(res) != 1 {
		t.Fatalf("abar result not match")
	}
	if res[0].Parent.Tag != "div" {
		t.Fatalf("abar result error")
	}
}
