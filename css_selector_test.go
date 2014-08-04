package hns

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

	ret := root.Walk(Css("  div.div", Return)).Return
	if len(ret) != 1 || ret[0].Id != "main" {
		t.Fatal()
	}

	ret = root.Walk(Css("div", Return)).Return
	if len(ret) != 3 {
		t.Fatal()
	}

	ret = root.Walk(Css("#bar", Return)).Return
	if len(ret) != 1 || ret[0].Tag != "div" {
		t.Fatal()
	}

	ret = root.Walk(Css("div.div #bar", Return)).Return
	if len(ret) != 1 || ret[0].Text != "BAR" {
		t.Fatal()
	}

	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Fatal()
			}
		}()
		root.Walk(Css("<>", Return))
	}()

	ret = root.Walk(Css("div div div", Return)).Return
	if len(ret) != 0 {
		t.Fatal()
	}
}
