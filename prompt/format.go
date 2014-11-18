// Library for inserting ANSI escapes into prompt strings.
package prompt

import "fmt"
import "unicode/utf8"

type StyleMarker struct {
	// A string containing an ANSI terminal escape code.
	escapeCode string
	// The position in the prompt string where this escape code should be
	// inserted. This is a byte offset; it is agnostic to UTF-8 rune encoding.
	pos int
}

// A prompt string, with some formatting markers.
type Prompt struct {
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

func NewPrompt() *Prompt {
	var p = new(Prompt)
	p.ClearStyle()
	return p
}

func (prompt *Prompt) Len() int {
	return utf8.RuneCountInString(prompt.text)
}

// Appends some text to this Prompt.
func (prompt *Prompt) Write(text string) {
	prompt.text += text
}

// Appends a new style, starting at the end of the current text.
func (self *Prompt) appendMarker(escapeCode string, pos int) {
	var newMarker = StyleMarker{escapeCode, pos}
	var lastMarker = self.lastMarker()
	if lastMarker != nil && lastMarker.pos == pos {
		// Replace lastMarker with newMarker.
		self.styleMarkers[len(self.styleMarkers)-1] = newMarker
	} else if lastMarker != nil && lastMarker.escapeCode == escapeCode {
		// The new marker is the same as the existing style, so don't add anything.
	} else {
		// Append newMarker to the list.
		self.styleMarkers = append(self.styleMarkers, newMarker)
	}
}

// Gets the last StyleMarker in this Prompt, or nil if there are no StyleMarkers
// in this Prompt.
func (self *Prompt) lastMarker() *StyleMarker {
	var lastMarkerIndex = len(self.styleMarkers) - 1
	if lastMarkerIndex >= 0 {
		return &self.styleMarkers[lastMarkerIndex]
	} else {
		return nil
	}
}

// Applies a new style at the end of this prompt, using the given foreground
// color and modifier.
func (prompt *Prompt) Style(color int, modifier int) {
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
	prompt.appendMarker(escape, len(prompt.text))
}

// Resets the style at the end of this Prompt to have no special styling.
func (prompt *Prompt) ClearStyle() {
	prompt.appendMarker("\033[0m", len(prompt.text))
}

// Concatenates 'other' onto this Prompt.
func (prompt *Prompt) Append(other *Prompt) {
	offset := len(prompt.text)
	prompt.text += other.text
	for _, marker := range other.styleMarkers {
		prompt.appendMarker(marker.escapeCode, marker.pos+offset)
	}
}

// Serializes this Prompt to a string with embedded ANSI escape sequences.
func (prompt *Prompt) String() string {
	buffer := ""
	nextPos := 0

	for _, marker := range prompt.styleMarkers {
		// Take the characters up to this marker.
		buffer += prompt.text[nextPos:marker.pos]
		nextPos = marker.pos
		// Append the marker, wrapping it in prompt-safe %{ %} delimiters.
		buffer += "%{" + marker.escapeCode + "%}"
	}

	// Take the remaining characters after the final style marker.
	buffer += prompt.text[nextPos:]

	// End with a clean format.
	buffer += "%{\033[0m%}"
	return buffer
}
