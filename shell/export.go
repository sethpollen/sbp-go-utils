// Library for exporting a set of environment variables to a shell.
package shell

import "bytes"
import "fmt"
import "sort"

// Represents a set of variables to set or unset in the environment.
type EnvironMod struct {
	// Keys are variable names. Values are nil for variables to unset. Otherwise,
	// values are pointers to variable values to set.
	mods map[string]*string
}

func NewEnvironMod() *EnvironMod {
	var mod = new(EnvironMod)
	mod.mods = make(map[string]*string)
	return mod
}

func (self *EnvironMod) SetVar(key, value string) {
	self.mods[key] = &value
}

func (self *EnvironMod) UnsetVar(key string) {
	self.mods[key] = nil
}

// Generates a shell script which can be sourced in a shell to apply this
// EnvironMod.
func (self *EnvironMod) ToScript() string {
	// We want to output the keys in sorted order. We have to do this sorting
	// ourselves.
	var keys []string
	for key := range self.mods {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var buf = bytes.NewBufferString("")
	for _, key := range keys {
		var value = self.mods[key]
		if value == nil {
			fmt.Fprintf(buf, "unset %s\n", key)
		} else {
			fmt.Fprintf(buf, "export %s=%s\n", key, quote(*value))
		}
	}
	return buf.String()
}

// Escapes and quotes 'text' so it may safely be embedded into a shell script.
func quote(text string) string {
	var buf = bytes.NewBuffer(make([]byte, 0, 2+2*len(text)))
	// Use single quote to avoid variable substitution.
	fmt.Fprint(buf, "'")
	for _, c := range text {
		if c == '\'' {
			fmt.Fprint(buf, "\\'")
		} else {
			// In a POSIX shell, this works even for newlines!
			fmt.Fprintf(buf, "%c", c)
		}
	}
	fmt.Fprint(buf, "'")
	return buf.String()
}
