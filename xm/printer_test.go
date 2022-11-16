package xm

import (
	"bytes"
	"fmt"
)

func ExamplePrinter() {
	buf := bytes.Buffer{}

	p := NewPrinter(Indent2Spaces,
		func(s []byte) { buf.Write(s) },
		func(n string) TagKind {
			if n == "em" || n == "strong" {
				return Inline
			} else {
				return Block
			}
		})

	p.OTag("root")
	p.Attr("key", RawAttr("val"))

	p.OTag("div")
	p.CTag()

	p.OTag("h1")
	p.Content(RawCont("heading"))
	p.CTag() // p
	p.OTag("p")
	p.Content(RawCont("paragraph"))
	p.CTag() // p
	p.OTag("p")
	p.CTag() // p

	// content block
	p.Content(RawCont("Hello "))
	p.OTag("em")
	p.Content(RawCont("World!"))
	p.CTag()

	// content block
	p.OTag("strong")
	p.Content(RawCont("Hello"))
	p.CTag()
	p.Content(RawCont(" World!"))

	p.OTag("p")
	p.OTag("strong")
	p.Content(RawCont("bold"))
	p.CTag()
	p.CTag() // p

	p.OTag("div")
	p.OTag("div")
	p.CTag()
	p.CTag()

	p.OTag("p")
	// note: this can be turned off with the PreserveInlineWhitespace flag
	p.Content(RawCont("line breaks must\nalign nicely with additional\nindentation that matches parent's\nblock level"))
	p.CTag()

	p.CTag() // root

	fmt.Println(buf.String())

	// Output:
	// <root key='val'>
	//   <div/>
	//   <h1>heading</h1>
	//   <p>paragraph</p>
	//   <p/>
	//   Hello <em>World!</em><strong>Hello</strong> World!
	//   <p><strong>bold</strong></p>
	//   <div>
	//     <div/>
	//   </div>
	//   <p>line breaks must
	//     align nicely with additional
	//     indentation that matches parent's
	//     block level</p>
	// </root>

}
