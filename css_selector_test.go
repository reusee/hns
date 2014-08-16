package nw

import "testing"

var cssTestHtml = `
<div class="div" id="main">
	<div id="foo">FOO</div>
	<div id="bar">BAR</div>
</div>
`

func TestCssSelector(t *testing.T) {
	root, err := ParseString(cssTestHtml)
	if err != nil {
		t.Fatal(err)
	}

	var ret []*Node

	ret = ret[:0]
	root.Walk(Css("  div.div", Append(&ret)))
	if len(ret) != 1 || ret[0].Id != "main" {
		t.Fatal()
	}

	ret = ret[:0]
	root.Walk(Css("div", Append(&ret)))
	if len(ret) != 3 {
		t.Fatal()
	}

	ret = ret[:0]
	root.Walk(Css("#bar", Append(&ret)))
	if len(ret) != 1 || ret[0].Tag != "div" {
		t.Fatal()
	}

	ret = ret[:0]
	root.Walk(Css("div.div #bar", Append(&ret)))
	if len(ret) != 1 || ret[0].Text != "BAR" {
		t.Fatal()
	}

	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Fatal()
			}
		}()
		root.Walk(Css("<>", Append(&ret)))
	}()

	ret = ret[:0]
	root.Walk(Css("div div div", Append(&ret)))
	if len(ret) != 0 {
		t.Fatal()
	}
}
