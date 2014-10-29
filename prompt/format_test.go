package prompt

import "strconv"
import "testing"

func TestEmptyPrompt(t *testing.T) {
	var p Prompt
	if p.Dump() != "%{\033[0m%}" {
		t.Error("Dump ==", strconv.Quote(p.Dump()))
	}
	if p.Len() != 0 {
		t.Error("Len ==", p.Len())
	}
}

func TestNoFormatting(t *testing.T) {
	var p Prompt
	p.Write("ABC")
	if p.Dump() != "%{\033[0m%}ABC" {
		t.Error("Dump ==", strconv.Quote(p.Dump()))
	}
	if p.Len() != 3 {
		t.Error("Len ==", p.Len())
	}
}

func TestFormatting(t *testing.T) {
	var p Prompt
	p.Write("A")
	p.Style(yellow, true)
	p.Write("B")
	p.Style(cyan, false)
	p.Write("C")
	p.ClearStyle()
	p.Write("D")
	if p.Dump() != "%{\033[0m%}A%{\033[1;33m%}B%{\033[0;36m%}C%{\033[0m%}D" {
		t.Error("Dump ==", strconv.Quote(p.Dump()))
	}
	if p.Len() != 4 {
		t.Error("Len ==", p.Len())
	}
}

func TestAppend(t *testing.T) {
	var p, q Prompt
	p.Style(yellow, true)
	p.Write("This is p.")
	q.Write("This ")
	q.Style(white, true)
	q.Write("is q.")
	p.Append(&q)
	if p.Dump() != "%{\033[0m%}%{\033[1;33m%}This is p.This %{\033[1;37m%}is q." {
		t.Error("Dump ==", strconv.Quote(p.Dump()))
	}
	if p.Len() != 20 {
		t.Error("Len ==", p.Len())
	}
}
