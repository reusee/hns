package hns

import "testing"

func TestWalk(t *testing.T) {
	root, err := ParseString(`
<div id="1">
	<p id="1a">foo</p>
	<p id="1b">bar</p>
	<p id="1c">baz</p>
	<a id="1d">foo</a>
	<div id="foo">
		<a id="div-a">bar</a>
	</div>
</div>
<ul id="u"></ul>
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

	// ChildrenIdEq
	res = root.Walk(ChildrenIdEq("u", Return)).Return
	if len(res) != 1 || res[0].Tag != "ul" {
		t.Fatalf("ChildrenIdEq result error")
	}

	// AllDescendantIdEq
	res = root.Walk(AllDescendantIdEq("1c", Return)).Return
	if len(res) != 1 || res[0].Tag != "p" {
		t.Fatalf("AllDescendantIdEq result error")
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

	// Return and AllDescendantTagEq
	res = root.Walk(AllDescendantTagEq("div", Return)).Return
	if len(res) != 2 {
		t.Fatalf("Return reuslt not match")
	}
	if len(res[0].Children) != 5 || len(res[1].Children) != 1 {
		t.Fatalf("Return result error")
	}

	// TagMatch
	res = root.Walk(DescendantTagMatch("p|ul", Return)).Return
	if len(res) != 4 || res[0].Tag != "p" || res[1].Tag != "p" || res[2].Tag != "p" || res[3].Tag != "ul" {
		t.Fatalf("DescendantTagMatch result error")
	}
	res = root.Walk(AllDescendantTagMatch("div|a", Return)).Return
	if len(res) != 4 || res[0].Tag != "div" || res[1].Tag != "a" || res[2].Tag != "div" || res[3].Tag != "a" {
		p("%v\n", res)
		t.Fatalf("AllDescendantTagMatch result error")
	}
	res = root.Walk(ChildrenTagMatch("div|ul", Return)).Return
	if len(res) != 2 || res[0].Attr["id"] != "1" || res[1].Attr["id"] != "u" {
		t.Fatalf("ChildrenTagMatch result error")
	}

	// IdMatch
	res = root.Walk(DescendantIdMatch("1[a-z]+", Return)).Return
	if len(res) != 4 || res[0].Attr["id"] != "1a" {
		t.Fatalf("DescendantIdMatch result error")
	}
	res = root.Walk(AllDescendantIdMatch("1.*", Return)).Return
	if len(res) != 5 || res[3].Attr["id"] != "1c" {
		t.Fatalf("AllDescendantIdMatch result error")
	}
	res = root.Walk(ChildrenIdMatch(".*", Return)).Return
	if len(res) != 2 || res[1].Attr["id"] != "u" {
		t.Fatalf("ChildrenIdMatch result error")
	}
}
