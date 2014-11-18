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
	p.Style(Yellow, Bold)
	p.Write("B")
	p.Style(Cyan, Dim)
	p.Write("C")
	p.Style(Red, Intense)
	p.Write("D")
	p.ClearStyle()
	p.Write("E")
	if p.String() !=
		"%{\033[0m%}A%{\033[1;93m%}B%{\033[0;36m%}C%{\033[0;91m%}D%{\033[0m%}E" +
    "%{\033[0m%}" {
		t.Error("String ==", strconv.Quote(p.String()))
	}
	if p.Len() != 5 {
		t.Error("Len ==", p.Len())
	}
}

func TestAppend(t *testing.T) {
	var p = NewPrompt()
	var q = NewPrompt()
	p.Style(Yellow, Bold)
	p.Write("This is p.")
	q.Write("This ")
	q.Style(White, Bold)
	q.Write("is q.")
	p.Append(q)
	if p.String() !=
		"%{\033[1;93m%}This is p.%{\033[0m%}This %{\033[1;97m%}is q.%{\033[0m%}" {
		t.Error("String ==", strconv.Quote(p.String()))
	}
	if p.Len() != 20 {
		t.Error("Len ==", p.Len())
	}
}
