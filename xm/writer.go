package xm

import (
	"sort"
)

// AttrWriter is an interface for writing XML attributes.
type AttrWriter interface {
	// Attr writes key='val' pair into tag's attributes. Accepted val types are:
	//
	//   - RawAttr, a special version of []byte that is written as-is:
	//   - string, gets scrambled with the ScrambleAttr() function
	//   - nils, or pointer types resolving to nil empty attribute
	//   - types supporting AttrMarshaler interface are resolved as t.MarshalXAttr()
	//   - types supporting encoding.TextMarshaler are marshaled into text, then scrambled with ScrambleAttr()
	//   - boolean types are resolved to 'true' or 'false'
	//   - integer types are resolved to their decimal representation
	//   - floating point types are converted to strings with strconv.FormatFloat using fmt='g' and prec=-1
	//   - all other types will panic with ErrUnsupportedType
	Attr(string, any)

	// OptAttr works similar to Attr, but it will skip writing the whole key='val'
	// pair if val is empty:
	//
	//   - RawAttr is considered empty if its length is zero
	//   - strings are considered empty if their length is zero
	//   - nils and pointer types resolving to nil are considered empty
	//   - types that support AttrMarshaler are considered empty if the bool part returned by t.MarshalXAttr() is false
	//   - types that support encoding.TextMarshaler are considered empty if v.MarshalText() returns empty byte slice
	//   - floating/integer/boolean types are never considered empty
	//   - all other types will panic with ErrUnsupportedType
	OptAttr(string, any)

	// Attrs writes map[string]any as attributes. The keys in the map are treated as
	// attribute keys. These are written raw without any scrambling or validation,
	// make sure you don't pass maps with keys that don't conform to xml attribute
	// key syntax.
	Attrs(map[string]any)
}

// ContWriter is an interface for writing content between tags.
type ContWriter interface {
	Content(...any)
}

// TagWriter is an interface for writing XML tags with optional attributes and content.
type TagWriter interface {
	// Tag writes an XML tag with attributes and content from args, where accepted arguments are:
	//
	//   - map[string]any - is written into tag attributes
	//   - func(AttrWriter) - is written into tag attribute
	//   - Attrs(map[string]T) - is processed into attributes, wrapped as func(AttrWriter)
	//   - all types supported by ContWriter - written into tag content
	Tag(string, ...any)
}

// Writer combines AttrWriter and TagWriter
type Writer interface {
	ContWriter
	TagWriter
}

// NewWriter wraps Printer p providing TagWriter API. Notice, that for a valid
// XML document, you will need to write exactly one tag into it that becomes the
// root.
func NewWriter(p Printer) Writer {
	return &writer_impl{p: p}
}

// Attrs takes a generic map[string]T and turns it into a functor for writing
// attributes that can be passed to TagWriter.
func Attrs[M ~map[string]T, T any](m M) func(AttrWriter) {
	return func(w AttrWriter) {
		kk := make([]string, 0, len(m))
		for a := range m {
			kk = append(kk, a)
		}
		sort.Strings(kk)
		for _, k := range kk {
			w.Attr(k, m[k])
		}
	}
}

func Attr[T any](key string, val T) func(AttrWriter) {
	return func(w AttrWriter) {
		w.Attr(key, val)
	}
}

func Tag(name string, args ...any) func(TagWriter) {
	return func(w TagWriter) {
		w.Tag(name, args...)
	}
}
