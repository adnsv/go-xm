package xm

import (
	"reflect"
	"sort"
)

type writer_impl struct {
	p Printer
}

func (w *writer_impl) attrEx(key string, val any, optional bool) {
	var raw RawAttr
	var ok bool

	switch v := val.(type) {
	case RawAttr:
		raw, ok = v, len(v) > 0

	case string:
		raw = ScrambleAttr(v)
		ok = len(raw) > 0

	case func(AttrWriter):
		// disregard optional flag in subfuncs
		v(w)
		return

	default:
		var s string
		if s, ok = coreToStr(val); ok {
			raw = RawAttr(s)
		} else {
			raw, ok = marshal_attr(reflect.ValueOf(val))
		}
	}

	if optional && !ok {
		return
	}
	w.p.Attr(key, raw)
}

// Attr implements AttrWriter.Attr().
func (w *writer_impl) Attr(key string, val any) {
	w.attrEx(key, val, false)
}

// OptAttr implements AttrWriter.OptAttr().
func (w *writer_impl) OptAttr(key string, val any) {
	w.attrEx(key, val, true)
}

// OptAttr implements AttrWriter.Attrs().
func (w *writer_impl) Attrs(aa map[string]any) {
	kk := make([]string, 0, len(aa))
	for a := range aa {
		kk = append(kk, a)
	}
	sort.Strings(kk)
	for _, k := range kk {
		w.Attr(k, aa[k])
	}
}

// Content implements ContentWriter.Content().
func (w *writer_impl) Content(args ...any) {
	for _, arg := range args {
		switch a := arg.(type) {
		case RawCont:
			w.p.Content(a)
		case string:
			w.p.Content(ScrambleCont(a))
		case func(ContWriter):
			a(w)
		case func(TagWriter):
			a(w)
		case func(Writer):
			a(w)
		case func(Printer):
			a(w.p)
		default:
			if r, ok := coreToStr(a); ok {
				w.p.Content(RawCont(r))
			} else {
				marshal_content(w.p, reflect.ValueOf(a))
			}
		}
	}
}

// Tag implements TagWriter.Tag().
func (w *writer_impl) Tag(name string, args ...any) {
	w.p.OTag(name)
	defer w.p.CTag()

	// attributes
	for _, arg := range args {
		switch a := arg.(type) {
		case map[string]any:
			w.Attrs(a)
		case func(AttrWriter):
			a(w)
		}
	}

	// content
	for _, arg := range args {
		switch a := arg.(type) {
		case map[string]any, func(AttrWriter):
			// skip attrs
			continue
		default:
			w.Content(a)
		}
	}
}
