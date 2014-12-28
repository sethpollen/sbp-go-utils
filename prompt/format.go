// Library for inserting ANSI escapes into prompt strings.
package prompt

import "fmt"
import "unicode/utf8"

type StyleMarker struct {
	// A string containing an ANSI terminal escape code.
	escapeCode string
	// The position in the string where this escape code should be
	// inserted. This is a byte offset; it is agnostic to UTF-8 rune encoding.
	pos int
}

// A string of text, with some formatting markers.
type StyledString struct {
	text string
	// Sorted by ascending pos.
	styleMarkers []StyleMarker
}

// Colors.
const (
	Black = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Font/color modifiers.
const (
  Dim = iota
  Intense
  Bold
)

const resetStyleEscape = "\033[0m"

func NewStyledString() *StyledString {
	var p = new(StyledString)
	p.ClearStyle()
	return p
}

func (self *StyledString) Len() int {
	return utf8.RuneCountInString(self.text)
}

// Appends some text to this StyledString.
func (self *StyledString) Write(text string) {
	self.text += text
}

// Appends a new style, starting at the end of the current text.
func (self *StyledString) appendMarker(escapeCode string, pos int) {
	var newMarker = StyleMarker{escapeCode, pos}
	var lastMarker = self.lastMarker()
	if lastMarker != nil && lastMarker.pos == pos {
		// Replace lastMarker with newMarker.
		self.styleMarkers[len(self.styleMarkers)-1] = newMarker
	} else if lastMarker != nil && lastMarker.escapeCode == escapeCode {
		// The new marker is the same as the existing style, so don't
		// add anything.
	} else {
		// Append newMarker to the list.
		self.styleMarkers = append(self.styleMarkers, newMarker)
	}
}

// Gets the last StyleMarker in this StyledString, or nil if there are no
// StyleMarkers in this StyledString.
func (self *StyledString) lastMarker() *StyleMarker {
	var lastMarkerIndex = len(self.styleMarkers) - 1
	if lastMarkerIndex >= 0 {
		return &self.styleMarkers[lastMarkerIndex]
	} else {
		return nil
	}
}

// Applies a new style at the end of this string, using the given foreground
// color and modifier.
func (self *StyledString) Style(color int, modifier int) {
	var boldness int = 0
  var colorOffset int = 30
  switch modifier {
    case Dim:
    case Intense:
      colorOffset = 90
    case Bold:
      boldness = 1
      colorOffset = 90
  }
	var escape = fmt.Sprintf("\033[%d;%dm", boldness, color + colorOffset)
	self.appendMarker(escape, len(self.text))
}

// Resets the style at the end of this StyledString to have no special
// styling.
func (self *StyledString) ClearStyle() {
	self.appendMarker("\033[0m", len(self.text))
}

// Concatenates 'other' onto this StyledString.
func (self *StyledString) Append(other *StyledString) {
	offset := len(self.text)
	self.text += other.text
	for _, marker := range other.styleMarkers {
		self.appendMarker(marker.escapeCode, marker.pos+offset)
	}
}

// Removes the first 'trim' bytes of text.
// TODO: test
func (self *StyledString) TrimLeft(trim int) {
  // TODO:
}

// Removes the last 'trim' bytes of text.
// TODO: test
func (self *StyledString) TrimRight(trim int) {
  self.text = self.text[0:len(self.text)-trim]
  self.trimMarkers()
}

// Removes markers which are off the end of the string.
func (self *StyledString) trimMarkers() {
  for {
    var last = self.lastMarker()
    if last == nil {
      break
    }
    if last.pos > len(self.text) {
      // Drop this marker, since it won't even apply to the next character
      // added to the string.
      self.styleMarkers = self.styleMarkers[0:len(self.styleMarkers)-1]
      continue
    }
    // Don't drop any more markers.
    break
  }
}

// Serializes this StyledString to a string with embedded ANSI escape
// sequences.
func (self *StyledString) String() string {
	buffer := ""
	nextPos := 0

	for _, marker := range self.styleMarkers {
		// Take the characters up to this marker.
		buffer += self.text[nextPos:marker.pos]
		nextPos = marker.pos
		// Append the marker, wrapping it in prompt-safe %{ %}
		// delimiters.
		buffer += "%{" + marker.escapeCode + "%}"
	}

	// Take the remaining characters after the final style marker.
	buffer += self.text[nextPos:]

	// End with a clean format.
	buffer += "%{\033[0m%}"
	return buffer
}

// Returns just the text fro this StyledString, without any formatting.
func (self* StyledString) PlainString() string {
  return self.text
}
