// Library for inserting ANSI escapes into prompt strings.

package prompt

import "fmt"

type StyleMarker struct {
	// A string containing an ANSI terminal escape code.
	escapeCode string
	// The position in the prompt string where this escape code should be
	// inserted.
	pos int
}

// A prompt string, with some formatting markers.
type Prompt struct {
	text string
	// Sorted by ascending pos.
	styleMarkers []StyleMarker
}

const (
	Black   = iota
	Red     = iota
	Green   = iota
	Yellow  = iota
	Blue    = iota
	Magenta = iota
	Cyan    = iota
	White   = iota
)

func (prompt *Prompt) Len() int {
	return len(prompt.text)
}

// Appends some text to this Prompt.
func (prompt *Prompt) Write(text string) {
	prompt.text += text
}

// Applies a new style at the end of this prompt, using the given foreground
// color and boldness.
func (prompt *Prompt) Style(color int, bold bool) {
	var boldness int
	if bold {
		boldness = 1
	} else {
		boldness = 0
	}
	escape := fmt.Sprintf("\033[%d;%dm", boldness, color+30)
	prompt.styleMarkers =
		append(prompt.styleMarkers, StyleMarker{escape, len(prompt.text)})
}

// Resets the style at the end of this Prompt to have no special styling.
func (prompt *Prompt) ClearStyle() {
	prompt.styleMarkers =
		append(prompt.styleMarkers, StyleMarker{"\033[0m", len(prompt.text)})
}

// Concatenates 'other' onto this Prompt.
func (prompt *Prompt) Append(other *Prompt) {
	offset := len(prompt.text)
	prompt.text += other.text
	for _, marker := range other.styleMarkers {
		prompt.styleMarkers = append(prompt.styleMarkers,
			StyleMarker{marker.escapeCode,
				marker.pos + offset})
	}
}

// Serializes this Prompt to a string with embedded ANSI escape sequences.
func (prompt *Prompt) Dump() string {
	buffer := "%{\033[0m%}" // Start with a clean format.
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
