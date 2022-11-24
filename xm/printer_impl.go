package xm

import (
	"bytes"
)

type printer_impl struct {
	put          func([]byte)
	names        []string // stack of tag names, used for closing tags
	block_level  int
	inline_level int
	inline_mode  bool
	in_tag       bool
	indent       IndentStyle
	flags        PrinterFlags
	on_tag_kind  func(n string) TagKind
}

func (p *printer_impl) BOM() {
	p.put([]byte("\uFEFF")) // writes \xef\xbb\xbf
}

func (p *printer_impl) XmlDecl() {
	if len(p.names) > 0 {
		panic("xml writer: invalid XmlDecl placement")
	}
	p.put([]byte("<?xml version='1.0' encoding='UTF-8'?>\n"))
}

func (p *printer_impl) Content(s RawCont) {
	if p.in_tag {
		p.in_tag = false
		p.put([]byte(">"))
	} else if !p.inline_mode {
		p.putIndent()
	}
	p.inline_mode = true

	if p.flags&PreserveInlineWhitespace != 0 || p.indent == IndentNone {
		p.put(s)
	} else {
		// re-indent after linebreaks
		i := bytes.IndexByte(s, '\n')
		if i < 0 {
			p.put(s)
			return
		}
		for {
			line := s[:i]
			p.put(line)
			s = s[i+1:]
			i = bytes.IndexByte(s, '\n')
			if i < 0 {
				p.putIndent()
				break
			} else if i == 0 {
				// a special handler for '\n\n' sequences to avoid generating
				// empty lines that only have spaces or tabs before the next '\n'
				p.put([]byte{'\n'})
			} else {
				p.putIndent()
			}
		}
		p.put(s)
	}
}

func (p *printer_impl) Linebreak() {
	if p.flags&PreserveInlineWhitespace == 0 {
		p.putIndent()
	} else {
		p.put([]byte{'\n'})
	}
}

func (p *printer_impl) StopInline() {
	if p.inline_level == 0 {
		p.inline_mode = false
	}
}

func (p *printer_impl) Attr(key string, val RawAttr) {
	if !p.in_tag {
		panic("xml writer: invalid xml printer.Attr call")
	}
	p.put([]byte{' '})
	p.put([]byte(key))
	p.put([]byte("='"))
	p.put(val)
	p.put([]byte{'\''})
}

func (p *printer_impl) OTag(name string) {
	if len(name) == 0 {
		panic("xml writer: trying to write a tag with empty name")
	}
	k := p.kindOf(name)
	p.names = append(p.names, name)

	was_in_tag := p.in_tag
	if p.in_tag {
		p.in_tag = false
		p.put([]byte{'>'})
	}

	if p.inline_level > 0 || k == Inline {
		if p.inline_mode || was_in_tag {
			p.inline_level++
		} else {
			p.putIndent()
			p.inline_mode = true
			p.inline_level++
		}
	} else { // block tag
		if p.inline_mode {
			p.inline_mode = false
		}
		p.putIndent()
		p.block_level++
	}

	p.put([]byte{'<'})
	p.put([]byte(name))
	p.in_tag = true
}

func (p *printer_impl) CTag() {
	stack_len := len(p.names)
	if stack_len == 0 {
		panic("xml writer: tag stack underflow, unpaired CTag call")
	}
	name := p.names[stack_len-1]

	was_inline := p.inline_mode

	if p.inline_mode {
		if p.inline_level > 0 {
			p.inline_level--
		} else {
			p.inline_mode = false
			p.block_level--
		}
	} else {
		p.block_level--
	}

	if p.in_tag {
		p.in_tag = false
		p.put([]byte("/>"))
	} else {
		if !was_inline {
			p.putIndent()
		}
		p.put([]byte("</"))
		p.put([]byte(name))
		p.put([]byte{'>'})
	}

	p.names = p.names[:stack_len-1]
}

const (
	eols_8    = "\n\n\n\n\n\n\n\n"
	tabs_8    = "\t\t\t\t\t\t\t\t"
	spaces_16 = "                "
)

func (p *printer_impl) putIndent() {
	if p.indent == IndentNone {
		return
	}
	if len(p.names) > 0 {
		p.put([]byte{'\n'})
	}
	if p.indent == IndentTabs {
		n := p.block_level
		for n > 8 {
			p.put([]byte(tabs_8))
			n -= 8
		}
		p.put([]byte(tabs_8[:n]))
	} else {
		n := p.block_level * int(p.indent)
		for n > 16 {
			p.put([]byte(spaces_16))
			n -= 16
		}
		p.put([]byte(spaces_16[:n]))
	}
}

func (p *printer_impl) kindOf(n string) TagKind {
	if p.on_tag_kind != nil {
		return p.on_tag_kind(n)
	}
	return Block
}
