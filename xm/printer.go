package xm

// RawAttr is used when writing XML attribute values to indicate that a value
// can be written without any further processing (bypasses ScrambleAttr)
type RawAttr []byte

// RawCont is used when writing XML content to indicate that it does not need
// any further processing (bypasses ScrambleCont).
type RawCont []byte

// DeclPrinter handles generation of top matter in the XML document.
type DeclPrinter interface {
	// BOM writes UTF-8 byte order mask.
	BOM()

	// XmlDecl writes XmlDecl at the top of the file.
	XmlDecl()
}

// AttrPrinter is an interface for writing XML tag attributes.
type AttrPrinter interface {
	// Attr adds key='val' pairs to a previously opened tag. Notice that Attr works
	// only immediately after opening the tag, once the opening tag is finalized,
	// calling Attr will panic. See description of OTag() for more details.
	Attr(key string, val RawAttr)
}

// ContPrinter is an interface for writing content between XML tags.
type ContPrinter interface {
	// Content places non-tagged content into the output. For indentation purposes,
	// this content is considered as inline. Typically, this would be some text
	// between or inside other tags, but you can also use this for emiting commends,
	// cdata, and processing instructions.
	Content(RawCont)

	// Linebreak can be used for inserting '\n' linebreaks after tags that have
	// line break semantics, like <br/> tags in html.
	Linebreak()

	StopInline()
}

type TagPrinter interface {
	// OTag starts a new open tag:
	//
	//  <name
	//
	// The name parameter is written verbatim, make sure all symbols in it
	// conform to XML standards.
	//
	// Once a tag is open, you can add attributes to it with the Attr command:
	//
	//    <name key='val' key2='val2'
	//
	// After writing the attributes, you can call CTag, which immediately closes
	// the tag:
	//
	//    <name key='val' key2='val2'/>
	//
	// Or you can start writing child content with the Content(...) call:
	//
	//    <name key='val' key2='val2'>...
	//
	// Or you can start writing subtags with the another OTag()/CTag() call:
	//
	//    <name key='val' key2='val2'><subtag/> ...
	//
	// All Content() and child OTag()/CTag() calls may be repeated and
	// interleaved. Once done, call CTag() to finalize the tag:
	//
	//    <name key='val' key2='val2'> ... child content and subtags ... </name>
	//
	// Once done with the tag, a matching CTag() call must be envoked.
	//
	OTag(name string)

	// CTag closes the tag that was previously opened with the OTag() call.
	//
	// If there was no content written after opening the tag, you will get
	//
	//    <tag optional='attributes' />
	//
	// If some content was written, you will get it wrapped in an open/close tag
	// pair:
	//
	//    <tag optional='attributes'> ... content ... </tag>
	//
	// If you want to have open/close pair with empty content, then make a dummy
	// empty content call after opening the tag: Content(nil):
	//
	//    <tag optional='attributes'></tag>
	//
	CTag()
}

// Printer combines DeclPrinter, AttrPrinter, ContPrinter, and TagPrinter
// interfaces into one providing the complete support for XML syntax.
type Printer interface {
	DeclPrinter
	AttrPrinter
	ContPrinter
	TagPrinter
}

// TagKind is used to customize the behavior of tags when styling the XML
// document's appearance with automatic indentation.
type TagKind int

// Accepted values for TagKind:
const (
	Block  = TagKind(iota) // block level indentation (default)
	Inline                 // inline tag
)

// IndentStyle specifies the indentation in the XML document.
type IndentStyle int

const (
	IndentTabs    = IndentStyle(0)  // each block level starts on a new line, indented with one '\t' per level
	Indent2Spaces = IndentStyle(2)  // each block level starts on a new line, indented with 2 spaces per level
	Indent4Spaces = IndentStyle(4)  // each block level starts on a new line, indented with 4 spaces per level
	IndentNone    = IndentStyle(-1) // no new lines, no indentation
)

type PrinterFlags uint

const (
	PreserveInlineWhitespace = PrinterFlags(1 << iota)
)

// NewPrinter creates a new Printer for writing XML files.
//
// The tagger parameter is a callback that allows to customize indentation for
// certain tags. If tagger is nil, then all the tags will be treated as block
// level tags.
//
func NewPrinter(indenter IndentStyle, putter func([]byte), tagger func(string) TagKind) Printer {
	return &printer_impl{
		put:         putter,
		indent:      indenter,
		on_tag_kind: tagger,
	}
}
