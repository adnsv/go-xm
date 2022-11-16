# go-xm 

[![GoDoc](https://godoc.org/github.com/adnsv/go-xm?status.svg)](https://godoc.org/github.com/adnsv/go-xm)

go-xm is a GO library for writing XML documents

While GO standard library already provides a package for marshaling custom data to XML, it has certain limitations that go-xm is trying to address.

Features:

- Supports declarative and functional writing
- Full control over block/inline indentation
- Fast and simple to use
- Type safety
- Custom type support
- No external dependencies

Not implemented yet (but planned):

- marshaling custom structs with field tags
- simplified high level support for block and inline comments
- simplified high level support for other elements (cdata, pi, etc.)

Notice that comments, dtd, pi, cdata, etc... can be injected in current
implementation with the functional `func(Printer)` call (see below).


Warning: unstable API, WIP

## Example

```go
func ExampleWriter() {
	buf := strings.Builder{}
	p := NewPrinter(Indent2Spaces,
		func(s []byte) { buf.Write(s) },
		func(n string) TagKind {
			if n == "em" || n == "strong" {
				return Inline
			} else {
				return Block
			}
		})
	w := NewWriter(p)
	p.XmlDecl()

	attrs := map[string]any{
		"k1": "v1",
		"k2": 42,
		"k3": 3.14,
	}

	notanymap := map[string]string{
		"s1": "v1",
		"s2": "v2",
		"s3": "v3",
	}

	w.Tag("root",
		Attr("key", "value"), Attr("bool", true), Attr("int", 42),

		// block tag indenting
		Tag("div"),
		Tag("div", Tag("subdiv", Tag("subsubdiv"))),

		// indentation is handled differently for block and inline tags
		Tag("em"),
		Tag("em", Tag("em", Tag("em"))),

		// simple content
		Tag("h1", "String Content"),

		// plain content between tags
		"Plain text content with automatic\ncharacter reference scrambling\nthat also supports aligned wrapping",

		// inline child tags within content
		Tag("p", "String content with ", Tag("strong", "inline formatting"), " is handled as expected"),

		// block child tags within content
		Tag("p", "String content with", Tag("div", "block subtags"), "is handled differently"),

		// declarative attributes
		Tag("div", Attr("key", "val"), "automatically sorts between attributes and content", attrs),
		Tag("div", Attrs(notanymap), "maps that are not `map[string]any` must be wrapped with Attrs(mymap)"),

		// functional attributes
		Tag("div", "functional and subfunctional attribute writing", func(attrs AttrWriter) {
			attrs.Attr("k", "v")
			attrs.Attr("subfunc", func(sub AttrWriter) { sub.Attr("k2", "subfunc") })
		}),

		// funcional content
		func(sub Writer) {
			sub.Tag("p", "functional content writing")
			sub.Tag("p", func(subsub Writer) { subsub.Content("can be nested") })
		},

		// low level printing
		Tag("div", func(p Printer) {
			p.Content(nil) // start new line
			p.Linebreak()
			p.Content(RawCont("direct raw writing with higher performance"))
			p.OTag("p")
			p.Content(RawCont("<!--and and flexibility-->"))
			p.CTag()
			p.Content(RawCont("<![CDATA[...]]>"))
			p.Linebreak()
			p.Content(ScrambleCont("make sure you pair OTag/CTag calls\nand avoid writing <things> that do not comply with XML syntax"))
			p.StopInline() // make sure the following block level closing tag is indented and aligned nicely
		}),
	)

	fmt.Println(buf.String())
}
```

produces

```xml
<?xml version='1.0' encoding='UTF-8'?>
<root key='value' bool='true' int='42'>
  <div/>
  <div>
    <subdiv>
      <subsubdiv/>
    </subdiv>
  </div>
  <em/><em><em><em/></em></em>
  <h1>String Content</h1>
  Plain text content with automatic
  character reference scrambling
  that also supports aligned wrapping
  <p>String content with <strong>inline formatting</strong> is handled as expected</p>
  <p>String content with
    <div>block subtags</div>
    is handled differently</p>
  <div key='val' k1='v1' k2='42' k3='3.14'>automatically sorts between attributes and content</div>
  <div s1='v1' s2='v2' s3='v3'>maps that are not `map[string]any` must be wrapped with Attrs(mymap)</div>
  <div k='v' k2='subfunc'>functional and subfunctional attribute writing</div>
  <p>functional content writing</p>
  <p>can be nested</p>
  <div>
    direct raw writing with higher performance
    <p><!--and and flexibility--></p>
    <![CDATA[...]]>
    make sure you pair OTag/CTag calls
    and avoid writing &lt;things&gt; that do not comply with XML syntax
  </div>
</root>
```


## Documentation

Automatically generated documentation for the package can be viewed online here:
http://pkg.go.dev/github.com/adnsv/go-xm