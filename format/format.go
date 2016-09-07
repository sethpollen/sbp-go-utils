// Library for inserting ANSI escapes into prompt strings.
package prompt

import "bytes"
import "fmt"
import "unicode"
import "unicode/utf8"

// A string of text, with some formatting markers.
type StyledString []StyledRune

type StyledRune struct {
	Style Style
	Text  rune
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

type Style struct {
	Color    int // Black, Red, etc.
	Modifier int // Dim, Intense, or Bold.
}

const resetStyleEscape = "\033[0m"

// Constructs a StyledString containing the given 'text' with the given
// 'color' and style 'modifier'.
func Stylize(text string, color int, modifier int) StyledString {
	var runes = utf8.RuneCountInString(text)
	var result StyledString = make([]StyledRune, runes)
	for i, r := range text {
		result[i] = StyledRune{Style{color, modifier}, r}
	}
	return result
}

// Constructs a StyledString containing the given 'text' and a "don't care"
// style. Good for use with whitespace.
func Unstyled(text string) StyledString {
	return Stylize(text, Black, Dim)
}

// Formats a Style as an ANSI escape sequence and returns the escape sequence.
func (self Style) toAnsi() string {
	var boldness int = 0
	var colorOffset int = 30
	switch self.Modifier {
	case Dim: // Nothing.
	case Intense:
		colorOffset = 90
	case Bold:
		boldness = 1
		colorOffset = 90
	}
	// Always precede the new style escape with a reset to avoid leakage of any
	// style elements.
	return fmt.Sprintf("%s\033[%d;%dm", resetStyleEscape, boldness,
		self.Color+colorOffset)
}

// Serializes this StyledString to a string with embedded ANSI escape
// sequences.
func (self StyledString) String() string {
	var buffer bytes.Buffer
	var first = true
	var lastStyle Style

	for _, r := range self {
		if unicode.IsSpace(r.Text) {
			// Don't bother applying style.
			buffer.WriteRune(r.Text)
			continue
		}
		if first || lastStyle != r.Style {
			// The style is changing, so insert a new style escape.
			buffer.WriteString("%{")
			buffer.WriteString(r.Style.toAnsi())
			buffer.WriteString("%}")

			first = false
			lastStyle = r.Style
		}
		buffer.WriteRune(r.Text)
	}

	// Clear style before ending.
	if !first {
		buffer.WriteString("%{")
		buffer.WriteString(resetStyleEscape)
		buffer.WriteString("%}")
	}
	return buffer.String()
}

// Returns just the text fro this StyledString, without any formatting.
func (self StyledString) PlainString() string {
	var buffer = bytes.NewBuffer(make([]byte, 0, len(self)))
	for _, r := range self {
		buffer.WriteRune(r.Text)
	}
	return buffer.String()
}
