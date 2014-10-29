// Library for constructing prompt strings of the specific form that I like.

package prompt

import "fmt"
import "os"
import "os/exec"
import "strings"
import "time"

// Main entry point for this module.
//   maxWidth - Maximum width that the prompt string may occupy, in characters.
//   info - An "info" string, which appears next to the PWD.
//   pwd - The current working directory to display.
//   errorcode - The result code of the previous shell command.
//   flag - A short "flag" string, which appears before the final $.
func MakePrompt(maxWidth int, info string, pwd string, errorCode int,
	flag string) *Prompt {
	// Perform tilde collapsing on the PWD.
	var home = os.Getenv("HOME")
	if strings.HasPrefix(pwd, home) {
		pwd = "~" + strings.TrimPrefix(pwd, home)
	}

	// If the hostname is a full domain name, remove all but the first domain
	// component.
	var hostname, _ = os.Hostname()
	var shortHostname = strings.SplitN(hostname, ".", 2)[0]
	var runningOverSsh = (len(os.Getenv("SSH_TTY")) > 0)

	// Format the date and time.
	var dateTime = time.Now().Format("01/02 15:04")

	// Now we must do the negotiation to decide how many characters to give to
	// each part of the prompt line. Start with non-negotiables.

	// Date and time.
	var dateTimePrompt Prompt
	dateTimePrompt.Style(Cyan, true)
	dateTimePrompt.Write(dateTime + " ")

	// Hostname.
	var hostnamePrompt Prompt
	if runningOverSsh {
		hostnamePrompt.Style(Yellow, false)
		hostnamePrompt.Write("(")
	}
	hostnamePrompt.Style(Magenta, true)
	hostnamePrompt.Write(shortHostname)
	if runningOverSsh {
		hostnamePrompt.Style(Yellow, false)
		hostnamePrompt.Write(") ")
	}

	// Info (if we got one).
	var infoPrompt Prompt
	if len(info) > 0 {
		infoPrompt.Style(White, false)
		infoPrompt.Write("[")
		infoPrompt.Style(White, true)
		infoPrompt.Write(info)
		infoPrompt.Style(White, false)
		infoPrompt.Write("] ")
	}

	var p = new(Prompt)
	return p
}
