package nw

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
	<ul class="list" id="uu"></ul>
</div>
<ul id="u" bar="BAR" class="list"></ul>
	`)
	if err != nil {
		t.Fatal(err)
	}

	var res []*Node

	// DescendantTagEq
	res = res[:0]
	root.Walk(
		Descendant(TagEq("div"),
			Descendant(TagEq("p"),
				func(node *Node) {
					res = append(res, node)
				})))
	if len(res) != 3 {
		t.Fatalf("DescendantTagEq not match")
	}
	if res[0].Text != "foo" || res[1].Text != "bar" || res[2].Text != "baz" {
		t.Fatalf("DescendantTagEq result error")
	}

	// DescendantIdEq
	res = res[:0]
	root.Walk(Descendant(IdEq("foo"), func(node *Node) {
		res = append(res, node)
	}))
	if len(res) != 1 {
		t.Fatalf("DescendantIdEq result not match")
	}
	if res[0].Tag != "div" || len(res[0].Children) != 1 {
		t.Fatalf("DescendantIdEq result error")
	}

	// ChildrenTagEq
	res = res[:0]
	root.Walk(
		Children(TagEq("div"),
			Children(TagEq("a"),
				func(n *Node) {
					res = append(res, n)
				})))
	if len(res) != 1 {
		t.Fatalf("ChildrenTagEq result not match")
	}
	if res[0].Text != "foo" {
		t.Fatalf("ChildrenTagEq result error")
	}

	// ChildrenIdEq
	res = res[:0]
	root.Walk(Children(IdEq("u"), Append(&res)))
	if len(res) != 1 || res[0].Tag != "ul" {
		t.Fatalf("ChildrenIdEq result error")
	}

	// AllDescendantIdEq
	res = res[:0]
	root.Walk(AllDescendant(IdEq("1c"), Append(&res)))
	if len(res) != 1 || res[0].Tag != "p" {
		t.Fatalf("AllDescendantIdEq result error")
	}

	// Descendant
	res = res[:0]
	root.Walk(Descendant(func(node *Node) bool {
		return node.Tag == "a" && node.Text == "bar"
	}, func(node *Node) {
		res = append(res, node)
	}))
	if len(res) != 1 {
		t.Fatalf("Descendant result not match")
	}
	if res[0].Parent.Tag != "div" {
		t.Fatalf("Descendant result error")
	}

	// AllDescendantTagEq
	res = res[:0]
	root.Walk(AllDescendant(TagEq("div"), Append(&res)))
	if len(res) != 2 {
		t.Fatalf("reuslt not match")
	}
	if len(res[0].Children) != 6 || len(res[1].Children) != 1 {
		t.Fatalf("result error")
	}

	// TagMatch
	res = res[:0]
	root.Walk(Descendant(TagMatch("p|ul"), Append(&res)))
	if len(res) != 5 || res[0].Tag != "p" || res[1].Tag != "p" || res[2].Tag != "p" || res[3].Tag != "ul" || res[4].Tag != "ul" {
		t.Fatalf("DescendantTagMatch result error")
	}
	res = res[:0]
	root.Walk(AllDescendant(TagMatch("div|a"), Append(&res)))
	if len(res) != 4 || res[0].Tag != "div" || res[1].Tag != "a" || res[2].Tag != "div" || res[3].Tag != "a" {
		p("%v\n", res)
		t.Fatalf("AllDescendantTagMatch result error")
	}
	res = res[:0]
	root.Walk(Children(TagMatch("div|ul"), Append(&res)))
	if len(res) != 2 || res[0].Attr["id"] != "1" || res[1].Attr["id"] != "u" {
		t.Fatalf("ChildrenTagMatch result error")
	}

	// IdMatch
	res = res[:0]
	root.Walk(Descendant(IdMatch("1[a-z]+"), Append(&res)))
	if len(res) != 4 || res[0].Attr["id"] != "1a" {
		t.Fatalf("DescendantIdMatch result error")
	}
	res = res[:0]
	root.Walk(AllDescendant(IdMatch("1.*"), Append(&res)))
	if len(res) != 5 || res[3].Attr["id"] != "1c" {
		t.Fatalf("AllDescendantIdMatch result error")
	}
	res = res[:0]
	root.Walk(Children(IdMatch(".*"), Append(&res)))
	if len(res) != 2 || res[1].Attr["id"] != "u" {
		t.Fatalf("ChildrenIdMatch result error")
	}

	// AttrEq
	res = res[:0]
	root.Walk(Descendant(AttrEq("foo", "foo"), Append(&res)))
	if len(res) != 2 || res[0].Attr["id"] != "1a" || res[1].Attr["id"] != "1b" {
		t.Fatalf("DescendantAttrEq result error")
	}
	res = res[:0]
	root.Walk(AllDescendant(AttrEq("bar", "bar"), Append(&res)))
	if len(res) != 3 || res[0].Id != "1" || res[1].Id != "foo" || res[2].Id != "div-a" {
		t.Fatalf("AllDescendantAttrEq result error")
	}
	res = res[:0]
	root.Walk(Children(AttrEq("bar", "bar"), Append(&res)))
	if len(res) != 1 || res[0].Id != "1" {
		t.Fatalf("ChildrenAttrEq result error")
	}

	// AttrMatch
	res = res[:0]
	root.Walk(Descendant(AttrMatch("foo", "foo"), Append(&res)))
	if len(res) != 2 || res[0].Attr["id"] != "1a" || res[1].Attr["id"] != "1b" {
		t.Fatalf("DescendantAttrMatch result error")
	}
	res = res[:0]
	root.Walk(AllDescendant(AttrMatch("bar", "bar"), Append(&res)))
	if len(res) != 3 || res[0].Id != "1" || res[1].Id != "foo" || res[2].Id != "div-a" {
		t.Fatalf("AllDescendantAttrMatch result error")
	}
	res = res[:0]
	root.Walk(Children(AttrMatch("bar", "bar"), Append(&res)))
	if len(res) != 1 || res[0].Id != "1" {
		t.Fatalf("ChildrenAttrMatch result error")
	}

	// ClassEq
	res = res[:0]
	root.Walk(Descendant(ClassEq("foo"), Append(&res)))
	if len(res) != 1 || res[0].Id != "1" {
		t.Fatalf("DescendantClassEq result error")
	}
	res = res[:0]
	root.Walk(AllDescendant(ClassEq("div"), Append(&res)))
	if len(res) != 2 || res[0].Id != "1" || res[1].Id != "foo" {
		t.Fatalf("AllDescendantClassEq result error")
	}
	res = res[:0]
	root.Walk(Children(ClassEq("list"), Append(&res)))
	if len(res) != 1 || res[0].Id != "u" {
		t.Fatalf("ChildrenClassEq result error")
	}

	// ClassMatch
	res = res[:0]
	root.Walk(Descendant(ClassMatch("bar|list"), Append(&res)))
	if len(res) != 3 || res[0].Id != "foo" || res[1].Id != "uu" || res[2].Id != "u" {
		t.Fatalf("DescendantClassMatch result error")
	}
	res = res[:0]
	root.Walk(AllDescendant(ClassMatch("foo|bar|list"), Append(&res)))
	if len(res) != 4 || res[0].Id != "1" || res[1].Id != "foo" || res[2].Id != "uu" || res[3].Id != "u" {
		t.Fatalf("AllDescendantClassMatch result error")
	}
	res = res[:0]
	root.Walk(Children(ClassMatch("list"), Append(&res)))
	if len(res) != 1 || res[0].Id != "u" {
		t.Fatalf("ChildrenClassMatch result error")
	}

	// Multi and Current
	res = res[:0]
	root.Walk(Descendant(TagEq("div"), Multi(
		Children(TagEq("p"), Append(&res)),
		Children(TagEq("a"), Append(&res)),
		Children(TagEq("ul"), Append(&res)),
		Children(TagEq("div"), Current(AttrEq("bar", "bar"), Append(&res))),
	)))
	if len(res) != 6 || res[0].Id != "1a" || res[3].Id != "1d" || res[5].Id != "foo" {
		t.Fatal()
	}
}
