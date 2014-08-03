package hns

import (
	"strings"
	"testing"
)

const testHtml = `
<!-foobar->
<div class="foo" id="bar">
	<p>FOOBARBAZ</p>
	<img src="foobar" width="300" />
</div>
`

func TestParse(t *testing.T) {
	root, err := Parse(strings.NewReader(testHtml))
	if err != nil {
		t.Fatal(err)
	}
	if err := root.Compare(&Node{
		Tag: "ROOT",
		Children: []*Node{
			{
				Tag: "div",
				Attr: map[string]string{
					"class": "foo",
					"id":    "bar",
				},
				Raw: `<div class="foo" id="bar">
	<p>FOOBARBAZ</p>
	<img src="foobar" width="300" />
</div>`,
				Children: []*Node{
					{
						Tag:  "p",
						Text: "FOOBARBAZ",
						TextParts: []string{
							"FOOBARBAZ",
						},
						Raw: `<p>FOOBARBAZ</p>`,
					},
					{
						Tag: "img",
						Attr: map[string]string{
							"src":   "foobar",
							"width": "300",
						},
						Raw: `<img src="foobar" width="300" />`,
					},
				},
			},
		},
		Raw: testHtml,
	}); err != nil {
		t.Fatal(err)
	}
}

func TestParseError(t *testing.T) {
	root, err := Parse(strings.NewReader(testHtml))
	if err != nil {
		t.Fatal(err)
	}

	if err := root.Compare(&Node{}); err == nil || !strings.HasPrefix(err.Error(), "number of children") {
		t.Fatalf("allowing number of children not matched or error mismatched: %v", err)
	}

	if err := root.Compare(&Node{
		Children: []*Node{
			{
				Children: []*Node{
					{}, {},
				},
			},
		},
	}); err == nil || !strings.HasPrefix(err.Error(), "tag") {
		t.Fatalf("allowing tag mismatched or error mismatched: %v", err)
	}

	if err := root.Compare(&Node{
		Tag: "ROOT",
		Children: []*Node{
			{
				Tag: "div",
				Children: []*Node{
					{
						Tag: "p",
					},
					{
						Tag: "img",
					},
				},
			},
		},
	}); err == nil || !strings.HasPrefix(err.Error(), "text") {
		t.Fatalf("allowing text mismatched or error mismatched: %v", err)
	}

	if err := root.Compare(&Node{
		Tag: "ROOT",
		Children: []*Node{
			{
				Tag: "div",
				Children: []*Node{
					{
						Tag:  "p",
						Text: "FOOBARBAZ",
					},
					{
						Tag: "img",
					},
				},
			},
		},
	}); err == nil || !strings.HasPrefix(err.Error(), "textparts length") {
		t.Fatalf("allowing textparts length mismatched or error mismatched: %v", err)
	}

	if err := root.Compare(&Node{
		Tag: "ROOT",
		Children: []*Node{
			{
				Tag: "div",
				Children: []*Node{
					{
						Tag:  "p",
						Text: "FOOBARBAZ",
						TextParts: []string{
							"",
						},
					},
					{
						Tag: "img",
					},
				},
			},
		},
	}); err == nil || !strings.HasPrefix(err.Error(), "textparts") {
		t.Fatalf("allowing textparts mismatched or error mismatched: %v", err)
	}

	if err := root.Compare(&Node{
		Tag: "ROOT",
		Children: []*Node{
			{
				Tag: "div",
				Children: []*Node{
					{
						Tag:  "p",
						Text: "FOOBARBAZ",
						TextParts: []string{
							"FOOBARBAZ",
						},
					},
					{
						Tag: "img",
					},
				},
			},
		},
	}); err == nil || !strings.HasPrefix(err.Error(), "raw") {
		t.Fatalf("allowing raw mismatched or error mismatched: %v", err)
	}

	if err := root.Compare(&Node{
		Tag: "ROOT",
		Children: []*Node{
			{
				Tag: "div",
				Children: []*Node{
					{
						Tag:  "p",
						Text: "FOOBARBAZ",
						TextParts: []string{
							"FOOBARBAZ",
						},
						Raw: `<p>FOOBARBAZ</p>`,
					},
					{
						Tag: "img",
					},
				},
			},
		},
	}); err == nil || !strings.HasPrefix(err.Error(), "number of attr") {
		t.Fatalf("allowing number of attr mismatched or error mismatched: %v", err)
	}

	if err := root.Compare(&Node{
		Tag: "ROOT",
		Children: []*Node{
			{
				Tag: "div",
				Children: []*Node{
					{
						Tag:  "p",
						Text: "FOOBARBAZ",
						TextParts: []string{
							"FOOBARBAZ",
						},
						Raw: `<p>FOOBARBAZ</p>`,
					},
					{
						Tag: "img",
						Attr: map[string]string{
							"src":   "",
							"width": "",
						},
					},
				},
			},
		},
	}); err == nil || !strings.HasPrefix(err.Error(), "attr") {
		t.Fatalf("allowing attr mismatched or error mismatched: %v", err)
	}

}

func TestParseTagMismatched(t *testing.T) {
	_, err := Parse(strings.NewReader(`<div><p></div>`))
	if err == nil || err.Error() != "end tag mismatched, expected p, got div" {
		t.Fatalf("allowing tag mismatched")
	}
}
