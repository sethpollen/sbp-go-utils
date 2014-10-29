// Library for constructing prompt strings of the specific form that I like.

package prompt

import "fmt"
import "math"
import "os"
import "os/exec"
import "strings"
import "time"

// Main entry point for this module.
//   maxWidth - Maximum width that the prompt string may occupy, in characters.
//   info - An "info" string, which appears next to the PWD.
//   pwd - The current working directory to display.
//   exitCode - The result code of the previous shell command.
//   flag - A short "flag" string, which appears before the final $.
func MakePrompt(maxWidth int, info string, pwd string, exitCode int,
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
		promptBeforePwd.Write(") ")
	}

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
    promptAfterPwd.Write(fmt.Sprintf("[%d]", exitCode)
  }

  // Determine how much space is left for the PWD.
  var pwdWidth =
    math.Max(0, maxWidth - promptBeforePwd.Len() - promptAfterPwd.Len())
  var pwdOnItsOwnLine = false
  if maxWidth >= 25 && pwdWidth < 25 {
    // Don't cram the PWD into a tiny space; put it on its own line.
    pwdWidth = maxWidth
    pwdOnItsOwnLine = true
  }

  // Subtract 2 in case we have to include the ".." characters.
  var pwdStart = len(pwd) - (pwdWidth - 2)
  if start > 0 {
    // Truncate the pwd.
    if start >= len(pwd) {
      // There is no room for the PWD at all.
      pwd = ""
    } else {
      pwd = ".." + pwd[start:]
    }
  }

  // Build the complete prompt string.
  var fullPrompt Prompt
  fullPrompt.Append(promptBeforePwd)
  if pwdOnItsOwnLine {
    fullPrompt.Append(promptAfterPwd)
    fullPrompt.Write("\n")
    fullPrompt.Style(Cyan, true)
    fullPrompt.Write(pwd)
  } else {
    fullPrompt.Style(Cyan, true)
    fullPrompt.Write(pwd + " ")
    fullPrompt.Append(promptAfterPwd)
  }
  fullPrompt.Style(Yellow, true)
  fullPrompt.Write("\n" + flag + "$ ")

  return fullPrompt
}
