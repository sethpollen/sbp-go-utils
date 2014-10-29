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

const escapeIntro = "\033["
const {
  black = iota
  red = iota
  green = iota
  yellow = iota
  blue = iota
  magenta = iota
  cyan = iota
  white = iota
}

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
  escape := escapeIntro
  if bold {
    escape += "1;"
  }
  escape += fmt.Sprintf("%dm", color + 30)
  append(prompt.styleMarkers, StyleMarker{escape, len(prompt.text)})
}

// Resets the style at the end of this Prompt to have no special styling.
func (prompt *Prompt) ClearStyle() {
  append(prompt.styleMarkers, StyleMarker{"\033[0m", len(prompt.text)})
}

// Concatenates 'other' onto this Prompt.
func (prompt *Prompt) Append(other *Prompt) {
	offset := len(prompt.text)
	prompt.text += other.text
	for marker := range other.styleMarkers {
		append(prompt.styleMarkers, StyleMarker{marker.escapeCode, marker.pos + offset})
	}
}

// Serializes this Prompt to a string with embedded ANSI escape sequences.
func (prompt *Prompt) Dump() string {
  buffer := ""
  nextPos := 0

  for marker := range prompt.styleMarkers {
    // Take the characters up to this marker.
    buffer += prompt.text[nextPos:marker.pos]
    nextPos = marker.pos
    // Append the marker, wrapping it in prompt-safe %{ %} delimiters.
    buffer += "%{" + marker.escapeCode + "%}"
  }

  // Take the remaining characters after the final style marker.
  buffer += prompt.text[nextPos:]

  return buffer
}
