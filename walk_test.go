package hns

import "testing"

func TestWalk(t *testing.T) {
	root, err := ParseString(`
<div id="1" bar="bar" class="div foo">
	<p id="1a" foo="foo">foo</p>
	<p id="1b" foo="foo">bar</p>
	<p id="1c">baz</p>
	<a id="1d">foo</a>
	<div id="foo" bar="bar" class="bar div ">
		<a id="div-a" bar="bar">bar</a>
	</div>
	<ul class="list"></ul>
</div>
<ul id="u" bar="BAR" class="list"></ul>
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
	if len(res[0].Children) != 6 || len(res[1].Children) != 1 {
		t.Fatalf("Return result error")
	}

	// TagMatch
	res = root.Walk(DescendantTagMatch("p|ul", Return)).Return
	if len(res) != 5 || res[0].Tag != "p" || res[1].Tag != "p" || res[2].Tag != "p" || res[3].Tag != "ul" || res[4].Tag != "ul" {
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

	// AttrEq
	res = root.Walk(DescendantAttrEq("foo", "foo", Return)).Return
	if len(res) != 2 || res[0].Attr["id"] != "1a" || res[1].Attr["id"] != "1b" {
		t.Fatalf("DescendantAttrEq result error")
	}
	res = root.Walk(AllDescendantAttrEq("bar", "bar", Return)).Return
	if len(res) != 3 || res[0].Id != "1" || res[1].Id != "foo" || res[2].Id != "div-a" {
		t.Fatalf("AllDescendantAttrEq result error")
	}
	res = root.Walk(ChildrenAttrEq("bar", "bar", Return)).Return
	if len(res) != 1 || res[0].Id != "1" {
		t.Fatalf("ChildrenAttrEq result error")
	}

	// AttrMatch
	res = root.Walk(DescendantAttrMatch("foo", "foo", Return)).Return
	if len(res) != 2 || res[0].Attr["id"] != "1a" || res[1].Attr["id"] != "1b" {
		t.Fatalf("DescendantAttrMatch result error")
	}
	res = root.Walk(AllDescendantAttrMatch("bar", "bar", Return)).Return
	if len(res) != 3 || res[0].Id != "1" || res[1].Id != "foo" || res[2].Id != "div-a" {
		t.Fatalf("AllDescendantAttrMatch result error")
	}
	res = root.Walk(ChildrenAttrMatch("bar", "bar", Return)).Return
	if len(res) != 1 || res[0].Id != "1" {
		t.Fatalf("ChildrenAttrMatch result error")
	}

	// ClassEq
	res = root.Walk(DescendantClassEq("foo", Return)).Return
	if len(res) != 1 || res[0].Id != "1" {
		t.Fatalf("DescendantClassEq result error")
	}
	res = root.Walk(AllDescendantClassEq("div", Return)).Return
	if len(res) != 2 || res[0].Id != "1" || res[1].Id != "foo" {
		t.Fatalf("AllDescendantClassEq result error")
	}
	res = root.Walk(ChildrenClassEq("list", Return)).Return
	if len(res) != 1 || res[0].Id != "u" {
		t.Fatalf("ChildrenClassEq result error")
	}
}
