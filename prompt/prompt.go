// Library for constructing prompt strings of the specific form that I like.

package prompt

import "fmt"
import "os"
import "strings"
import "time"

// Injectable data for testing MakePrompt.
type PromptEnv struct {
	now      time.Time
	home     string
	hostname string
}

func DefaultPromptEnv() *PromptEnv {
	var env = new(PromptEnv)
	env.now = time.Now()
	env.home = os.Getenv("HOME")
	env.hostname, _ = os.Hostname()
	return env
}

// Main entry point for this module.
//   maxWidth - Maximum width that the prompt string may occupy, in characters.
//   info - An "info" string, which appears next to the PWD.
//   pwd - The current working directory to display.
//   exitCode - The result code of the previous shell command.
//   flag - A short "flag" string, which appears before the final $.
func MakePrompt(env *PromptEnv, maxWidth int, info string, pwd string,
	exitCode int, flag string) *Prompt {
	// Perform tilde collapsing on the PWD.
	var home = env.home
	if strings.HasSuffix(home, "/") {
		home = home[:len(home)-1]
	}
	if strings.HasPrefix(pwd, home) {
		pwd = "~" + pwd[len(home):]
	}

	// If the hostname is a full domain name, remove all but the first domain
	// component.
	var shortHostname = strings.SplitN(env.hostname, ".", 2)[0]
	var runningOverSsh = (len(os.Getenv("SSH_TTY")) > 0)

	// Format the date and time.
	var dateTime = env.now.Format("01/02 15:04")

	// Construct the prompt text which must precede the PWD.
	var promptBeforePwd Prompt

	// Date and time.
	promptBeforePwd.Style(Cyan, true)
	promptBeforePwd.Write(dateTime + " ")

	// Hostname.
	if runningOverSsh {
		promptBeforePwd.Style(Yellow, false)
		promptBeforePwd.Write("(")
	}
	promptBeforePwd.Style(Magenta, true)
	promptBeforePwd.Write(shortHostname)
	if runningOverSsh {
		promptBeforePwd.Style(Yellow, false)
		promptBeforePwd.Write(")")
	}
	promptBeforePwd.Write(" ")

	// Info (if we got one).
	if len(info) > 0 {
		promptBeforePwd.Style(White, false)
		promptBeforePwd.Write("[")
		promptBeforePwd.Style(White, true)
		promptBeforePwd.Write(info)
		promptBeforePwd.Style(White, false)
		promptBeforePwd.Write("] ")
	}

	// Construct the prompt text which must follow the PWD.
	var promptAfterPwd Prompt

	// Exit code.
	if exitCode != 0 {
		promptAfterPwd.Style(Red, true)
		promptAfterPwd.Write(fmt.Sprintf("[%d]", exitCode))
	}

	// Determine how much space is left for the PWD.
	var pwdWidth = maxWidth - promptBeforePwd.Len() - promptAfterPwd.Len()
	if pwdWidth < 0 {
		pwdWidth = 0
	}
	var pwdOnItsOwnLine = false
	if maxWidth >= 25 && pwdWidth < 25 {
		// Don't cram the PWD into a tiny space; put it on its own line.
		pwdWidth = maxWidth
		pwdOnItsOwnLine = true
	}

	// Subtract 2 in case we have to include the ".." characters.
	var pwdStart = len(pwd) - (pwdWidth - 2)
	if pwdStart > 0 {
		// Truncate the pwd.
		if pwdStart >= len(pwd) {
			// There is no room for the PWD at all.
			pwd = ""
		} else {
			pwd = ".." + pwd[pwdStart:]
		}
	}

	// Build the complete prompt string.
	var fullPrompt = new(Prompt)
	fullPrompt.Append(&promptBeforePwd)
	if pwdOnItsOwnLine {
		fullPrompt.Append(&promptAfterPwd)
		fullPrompt.Write("\n")
		fullPrompt.Style(Cyan, true)
		fullPrompt.Write(pwd)
	} else {
		fullPrompt.Style(Cyan, true)
		fullPrompt.Write(pwd + " ")
		fullPrompt.Append(&promptAfterPwd)
	}
	fullPrompt.Style(Yellow, true)
	fullPrompt.Write("\n" + flag + "$ ")

	return fullPrompt
}
