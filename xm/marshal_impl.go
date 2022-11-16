package xm

import (
	"encoding"
	"reflect"
	"strconv"
)

func marshal_attr(val reflect.Value) (RawAttr, bool) {
	// handle nil pointers
	for val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, false
		}
		val = val.Elem()
	}

	// handle ContMarshaler values
	typ := val.Type()
	if val.CanInterface() && typ.Implements(attrMarshalerType) {
		v := val.Interface().(AttrMarshaler)
		return v.MarshalXAttr()
	}
	if val.CanAddr() {
		pv := val.Addr()
		if pv.CanInterface() && pv.Type().Implements(attrMarshalerType) {
			v := pv.Interface().(AttrMarshaler)
			return v.MarshalXAttr()
		}
	}

	// handle encoding.TextMarshaler values
	if val.CanInterface() && typ.Implements(textMarshalerType) {
		return textMarshalerToAttr(val.Interface().(encoding.TextMarshaler))
	}
	if val.CanAddr() {
		pv := val.Addr()
		if pv.CanInterface() && pv.Type().Implements(textMarshalerType) {
			return textMarshalerToAttr(pv.Interface().(encoding.TextMarshaler))
		}
	}

	// handle booleans, integer, and floating point values
	if s, ok := reflectCoreToStr(val); ok {
		return RawAttr(s), true
	}

	panic(&ErrUnsupportedType{val.Type()})
}

func marshal_content(p Printer, val reflect.Value) {
	// handle nil pointers
	for val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}

	// handle ContMarshaler values
	typ := val.Type()
	if val.CanInterface() && typ.Implements(contMarshalerType) {
		v := val.Interface().(ContMarshaler)
		v.MarshalXCont(p)
		return
	}
	if val.CanAddr() {
		pv := val.Addr()
		if pv.CanInterface() && pv.Type().Implements(contMarshalerType) {
			v := pv.Interface().(ContMarshaler)
			v.MarshalXCont(p)
			return
		}
	}

	// handle encoding.TextMarshaler values
	if val.CanInterface() && typ.Implements(textMarshalerType) {
		textMarshalerToCont(p, val.Interface().(encoding.TextMarshaler))
		return
	}
	if val.CanAddr() {
		pv := val.Addr()
		if pv.CanInterface() && pv.Type().Implements(textMarshalerType) {
			textMarshalerToCont(p, pv.Interface().(encoding.TextMarshaler))
			return
		}
	}

	// handle booleans, integer, and floating point values
	if s, ok := reflectCoreToStr(val); ok {
		p.Content(RawCont(s))
		return
	}

	panic(&ErrUnsupportedType{val.Type()})
}

var (
	attrMarshalerType = reflect.TypeOf((*AttrMarshaler)(nil)).Elem()
	contMarshalerType = reflect.TypeOf((*ContMarshaler)(nil)).Elem()
	textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

func textMarshalerToAttr(v encoding.TextMarshaler) ([]byte, bool) {
	b, e := v.MarshalText()
	if e == nil {
		r := ScrambleAttr(string(b))
		return r, len(b) > 0
	} else {
		panic(e)
	}
}

func textMarshalerToCont(p Printer, v encoding.TextMarshaler) {
	b, e := v.MarshalText()
	if e == nil {
		p.Content(ScrambleCont(string(b)))
	} else {
		panic(e)
	}
}

// UnsupportedTypeError is returned when Marshal encounters a type
// that cannot be converted into XML.
type ErrUnsupportedType struct {
	reflect.Type
}

func (e ErrUnsupportedType) Error() string {
	return "xml: unsupported type: " + e.String()
}

// coreToStr converts supported core types to string without using reflect.
func coreToStr(val any) (string, bool) {
	switch v := val.(type) {
	case bool:
		return strconv.FormatBool(v), true
	case int:
		return strconv.FormatInt(int64(v), 10), true
	case int8:
		return strconv.FormatInt(int64(v), 10), true
	case int16:
		return strconv.FormatInt(int64(v), 10), true
	case int32:
		return strconv.FormatInt(int64(v), 10), true
	case int64:
		return strconv.FormatInt(v, 10), true
	case uint:
		return strconv.FormatUint(uint64(v), 10), true
	case uint8:
		return strconv.FormatUint(uint64(v), 10), true
	case uint16:
		return strconv.FormatUint(uint64(v), 10), true
	case uint32:
		return strconv.FormatUint(uint64(v), 10), true
	case uint64:
		return strconv.FormatUint(v, 10), true
	case uintptr:
		return strconv.FormatUint(uint64(v), 10), true
	case float32:
		return strconv.FormatFloat(float64(v), 'g', -1, 32), true
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64), true
	default:
		return "", false
	}
}

func reflectCoreToStr(val reflect.Value) (string, bool) {
	switch val.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(val.Bool()), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(val.Uint(), 10), true
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'g', -1, val.Type().Bits()), true
	default:
		return "", false
	}
}
