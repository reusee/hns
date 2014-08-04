package hns

import "testing"

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ParseString(`
<div>
	<p>
		<a>
		</a>
	</p>
</div>
	`)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCss(b *testing.B) {
	root, err := ParseString(`
<div>
	<p>
		<a>
			<img />
		</a>
	</p>
</div>
	`)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root.Walk(Css("div p a img", Return))
	}
}
