package prompt

import "strconv"
import "testing"

func TestEmpty(t *testing.T) {
	var p StyledString
	if p.String() != "" {
		t.Error("String ==", strconv.Quote(p.String()))
	}
	if p.PlainString() != "" {
		t.Error("PlainString ==", strconv.Quote(p.PlainString()))
	}
}

func TestStyle(t *testing.T) {
	var p StyledString = Stylize("abc", Red, Intense)
	if p.String() != "%{\033[0m\033[0;91m%}abc%{\033[0m%}" {
		t.Error("String ==", strconv.Quote(p.String()))
	}
	if p.PlainString() != "abc" {
		t.Error("PlainString ==", strconv.Quote(p.PlainString()))
	}
}

func TestAppend(t *testing.T) {
	var p StyledString = Stylize("ab", Red, Intense)
	p = append(p, Stylize("cd", Red, Intense)...)
	p = append(p, Stylize(" ef", Black, Intense)...)
	if p.String() !=
		"%{\033[0m\033[0;91m%}abcd %{\033[0m\033[0;90m%}ef%{\033[0m%}" {
		t.Error("String ==", strconv.Quote(p.String()))
	}
	if p.PlainString() != "abcd ef" {
		t.Error("PlainString ==", strconv.Quote(p.PlainString()))
	}
}

