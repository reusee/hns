package hns

import "testing"

func TestWalk(t *testing.T) {
	root, err := ParseString(`
<div>
	<p>foo</p>
	<p>bar</p>
	<p>baz</p>
	<a>foo</a>
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
		t.Fatalf("result not match")
	}
	if res[0].Text != "foo" || res[1].Text != "bar" || res[2].Text != "baz" {
		t.Fatalf("result error")
	}
}
