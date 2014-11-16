package prompt

import "strconv"
import "testing"

func TestEmptyPrompt(t *testing.T) {
	var p = NewPrompt()
	if p.String() != "%{\033[0m%}%{\033[0m%}" {
		t.Error("String ==", strconv.Quote(p.String()))
	}
	if p.Len() != 0 {
		t.Error("Len ==", p.Len())
	}
}

func TestNoFormatting(t *testing.T) {
	var p = NewPrompt()
	p.Write("ABC")
	if p.String() != "%{\033[0m%}ABC%{\033[0m%}" {
		t.Error("String ==", strconv.Quote(p.String()))
	}
	if p.Len() != 3 {
		t.Error("Len ==", p.Len())
	}
}

func TestFormatting(t *testing.T) {
	var p = NewPrompt()
	p.Write("A")
	p.Style(Yellow, true)
	p.Write("B")
	p.Style(Cyan, false)
	p.Write("C")
	p.ClearStyle()
	p.Write("D")
	if p.String() !=
		"%{\033[0m%}A%{\033[1;33m%}B%{\033[0;36m%}C%{\033[0m%}D%{\033[0m%}" {
		t.Error("String ==", strconv.Quote(p.String()))
	}
	if p.Len() != 4 {
		t.Error("Len ==", p.Len())
	}
}

func TestAppend(t *testing.T) {
	var p = NewPrompt()
	var q = NewPrompt()
	p.Style(Yellow, true)
	p.Write("This is p.")
	q.Write("This ")
	q.Style(White, true)
	q.Write("is q.")
	p.Append(q)
	if p.String() !=
		"%{\033[1;33m%}This is p.%{\033[0m%}This %{\033[1;37m%}is q.%{\033[0m%}" {
		t.Error("String ==", strconv.Quote(p.String()))
	}
	if p.Len() != 20 {
		t.Error("Len ==", p.Len())
	}
}
