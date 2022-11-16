package xm

import "errors"

type AttrMarshaler interface {
	MarshalXAttr() (RawAttr, bool)
}

type ContMarshaler interface {
	MarshalXCont(w Printer)
}

var ErrEmptyAttribute = errors.New("xml: empty sttribute")
