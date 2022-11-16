package xm

import (
	"fmt"
	"testing"
)

func TestScramble(t *testing.T) {
	tests := []struct {
		s         string
		want_attr string
		want_cont string
	}{
		{"", "", ""},
		{"abc", "abc", "abc"},
		{"\x00", "&#00;", "&#00;"},
		{"\t", "&#09;", "\t"},
		{"\n", "&#0a;", "\n"},
		{"<", "&lt;", "&lt;"},
		{">", "&gt;", "&gt;"},
		{"'", "&apos;", "'"},
		{"\"", "\"", "\""},
		{"&", "&amp;", "&amp;"},
		{"x 'y' z", "x &apos;y&apos; z", "x 'y' z"},
		{"x \"y\" z", "x \"y\" z", "x \"y\" z"},
		{"a&b", "a&amp;b", "a&amp;b"},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%q", tt.s)

		t.Run(name, func(t *testing.T) {
			got_attr := string(ScrambleAttr(tt.s))
			got_cont := string(ScrambleCont(tt.s))

			if got_attr != tt.want_attr || got_cont != tt.want_cont {
				t.Errorf("ScrambleAttr(), ScrambleCont() = %q, %q; want %q, %q",
					got_attr, got_cont, tt.want_attr, tt.want_cont)
			}
		})
	}
}
