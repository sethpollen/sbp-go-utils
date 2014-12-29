package prompt

import "strconv"
import "testing"

func TestEmptyStyledString(t *testing.T) {
	var p = NewStyledString()
	if p.String() != "%{\033[0m%}%{\033[0m%}" {
		t.Error("String ==", strconv.Quote(p.String()))
	}
	if p.Len() != 0 {
		t.Error("Len ==", p.Len())
	}
}

func TestNoFormatting(t *testing.T) {
	var p = NewStyledString()
	p.Write("ABC")
	if p.String() != "%{\033[0m%}ABC%{\033[0m%}" {
		t.Error("String ==", strconv.Quote(p.String()))
	}
	if p.Len() != 3 {
		t.Error("Len ==", p.Len())
	}
}

func TestFormatting(t *testing.T) {
	var p = NewStyledString()
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
	var p = NewStyledString()
	var q = NewStyledString()
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

func TestTrimFirst(t *testing.T) {
  var p = NewStyledString()
  p.Write("ZZZ")
  p.Style(Yellow, Bold)
  p.Write("ABC")
  p.Style(Green, Dim)
  p.Write("DEF")
  p.Style(Cyan, Intense)
  p.Write("GHI")

  // Issue a trim that doesn't touch any style markers.
  p.TrimFirst(3)

  var expected1 = NewStyledString()
  expected1.Style(Yellow, Bold)
  expected1.Write("ABC")
  expected1.Style(Green, Dim)
  expected1.Write("DEF")
  expected1.Style(Cyan, Intense)
  expected1.Write("GHI")
  if p.String() != expected1.String() {
    t.Error("String ==", strconv.Quote(p.String()))
  }

  // Trim so that the first style marker must be advanced.
  p.TrimFirst(1)

  var expected2 = NewStyledString()
  expected2.Style(Yellow, Bold)
  expected2.Write("BC")
  expected2.Style(Green, Dim)
  expected2.Write("DEF")
  expected2.Style(Cyan, Intense)
  expected2.Write("GHI")
  if p.String() != expected2.String() {
    t.Error("String ==", strconv.Quote(p.String()))
  }

  // Trim away the remaining two style markers.
  p.TrimFirst(6)

  var expected3 = NewStyledString()
  expected3.Style(Cyan, Intense)
  expected3.Write("HI")
  if p.String() != expected3.String() {
    t.Error("String ==", strconv.Quote(p.String()))
  }

  // Trim away everything.
  p.TrimFirst(2)

  var expected4 = NewStyledString()
  expected4.Style(Cyan, Intense)
  if p.String() != expected4.String() {
    t.Error("String ==", strconv.Quote(p.String()))
  }
}

func TestTrimLast(t *testing.T) {
  var p = NewStyledString()
  p.Style(Yellow, Bold)
  p.Write("ABC")
  p.Style(Green, Dim)
  p.Write("DEF")
  p.Style(Cyan, Intense)
  p.Write("GHI")

  // Issue a trim that doesn't touch any style markers.
  p.TrimLast(1)

  var expected1 = NewStyledString()
  expected1.Style(Yellow, Bold)
  expected1.Write("ABC")
  expected1.Style(Green, Dim)
  expected1.Write("DEF")
  expected1.Style(Cyan, Intense)
  expected1.Write("GH")
  if p.String() != expected1.String() {
    t.Error("String ==", strconv.Quote(p.String()))
  }

  // Issue a trim that removes two style markers.
  p.TrimLast(6)

  var expected2 = NewStyledString()
  expected2.Style(Yellow, Bold)
  expected2.Write("AB")
  if p.String() != expected2.String() {
    t.Error("String ==", strconv.Quote(p.String()))
  }

  // Trim away everything.
  p.TrimLast(2)

  var expected3 = NewStyledString()
  expected3.Style(Yellow, Bold)
  if p.String() != expected3.String() {
    t.Error("String ==", strconv.Quote(p.String()))
  }
}

