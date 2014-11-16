// Library for constructing prompt strings of the specific form that I like.
// TODO: unit tests are failing
package prompt

import "fmt"
import "os"
import "os/user"
import "strings"
import "time"
import "unicode/utf8"
import "code.google.com/p/sbp-go-utils/shell"

// Collects information during construction of a prompt string.
type PromptEnv struct {
	Now        time.Time
	Home       string
	Pwd        string
	Hostname   string
  // Text to include in the prompt, along with the PWD.
  Info       string
  // A secondary info string. Displayed using $RPROMPT.
  Info2      string
  // A short string to place before the final $ in the prompt.
  Flag       Prompt
  // Exit code of the last process run in the shell.
  ExitCode   int
	// Maximum number of characters which prompt may occupy horizontally.
	Width      int
  // Environment variables which should be emitted to the shell which uses this
  // prompt.
  EnvironMod shell.EnvironMod
}

// Generates a PromptEnv based on current environment variables. The maximum
// number of characters which the prompt may occupy must be passed as 'width'.
func NewPromptEnv(width int, exitCode int) *PromptEnv {
	var self = new(PromptEnv)
	self.Now = time.Now()

  user, err := user.Current()
  if err != nil {
    self.Home = ""
  } else {
	  self.Home = user.HomeDir
  }

  self.Pwd, _ = os.Getwd()
	self.Hostname, _ = os.Hostname()
  self.Info = ""
  self.Info2 = ""
  self.ExitCode = exitCode
	self.Width = width
  self.Flag = *NewPrompt()
  self.EnvironMod = *shell.NewEnvironMod()

  return self
}

// Generates a shell prompt string.
func (self *PromptEnv) makePrompt() *Prompt {
	// If the hostname is a full domain name, remove all but the first domain
	// component.
	var shortHostname = strings.SplitN(self.Hostname, ".", 2)[0]
	var runningOverSsh = (os.Getenv("SSH_TTY") != "")

	// Format the date and time.
	var dateTime = self.Now.Format("01/02 15:04")

	// Construct the prompt text which must precede the PWD.
	var promptBeforePwd = NewPrompt()

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
	if self.Info != "" {
		promptBeforePwd.Style(White, false)
		promptBeforePwd.Write("[")
		promptBeforePwd.Style(White, true)
		promptBeforePwd.Write(self.Info)
		promptBeforePwd.Style(White, false)
		promptBeforePwd.Write("] ")
	}

	// Construct the prompt text which must follow the PWD.
	var promptAfterPwd = NewPrompt()

	// Exit code.
	if self.ExitCode != 0 {
		promptAfterPwd.Style(Red, true)
		promptAfterPwd.Write(fmt.Sprintf("[%d]", self.ExitCode))
	}

	// Determine how much space is left for the PWD.
	var pwdWidth = self.Width - promptBeforePwd.Len() - promptAfterPwd.Len()
	if pwdWidth < 0 {
		pwdWidth = 0
	}
	var pwdOnItsOwnLine = false
	if pwdWidth < 20 && utf8.RuneCountInString(self.Pwd) >= 20 &&
     self.Width >= 20 {
		// Don't cram the PWD into a tiny space; put it on its own line.
		pwdWidth = self.Width
		pwdOnItsOwnLine = true
	}

	var pwd = self.formatPwd(pwdWidth)

	// Build the complete prompt string.
	var fullPrompt = NewPrompt()
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
  fullPrompt.Write("\n")
  fullPrompt.Append(&self.Flag)
	fullPrompt.Style(Yellow, true)
	fullPrompt.Write("$ ")

	return fullPrompt
}

// Generates a shell RPROMPT string. This will be printed on the right-hand
// side of the second line of the prompt. It will disappear if the user types
// a long command, so it should not be super important. self.Info2 will be the
// content displayed.
// TODO: unit test
func (self *PromptEnv) makeRPrompt() *Prompt {
  var rPrompt = NewPrompt()
  if self.Info2 != "" {
    rPrompt.Style(White, false)
    rPrompt.Write(self.Info2)
  }
  return rPrompt
}

// Generates a terminal emulator title bar string. Similar to a shell prompt
// string, but lacks formatting escapes.
func (self *PromptEnv) makeTitle() string {
  var info = ""
	if self.Info != "" {
		info = fmt.Sprintf("[%s]", self.Info)
	}
	var pwdWidth = self.Width - utf8.RuneCountInString(info)
	return info + self.formatPwd(pwdWidth)
}

// Formats the PWD for use in a prompt.
func (self *PromptEnv) formatPwd(width int) string {
	// Perform tilde collapsing on the PWD.
	var home = self.Home
	if strings.HasSuffix(home, "/") {
		home = home[:len(home)-1]
	}
	var pwd = self.Pwd
	if strings.HasPrefix(pwd, home) {
		pwd = "~" + pwd[len(home):]
	}
	if pwd == "" {
		pwd = "/"
	}

	// Subtract 2 in case we have to include the ".." characters.
	var pwdRunes = utf8.RuneCountInString(pwd)
	var start = pwdRunes - (width - 2)
	if start > 0 {
		// Truncate the PWD.
		if start >= pwdRunes {
			// There is no room for the PWD at all.
			pwd = ""
		} else {
			pwd = ".." + pwd[start:]
		}
	}
	return pwd
}

// Renders all the information from this PromptEnv into a shell script which
// may be sourced. The following variables will be set:
//   PROMPT
//   RPROMPT
//   TERM_TITLE
//   ... plus any other variables set in self.EnvironMod.
func (self *PromptEnv) ToScript() string {
  // Start by making a copy of the custom EnvironMod.
  var mod = self.EnvironMod
  // Now add our variables to it.
  mod.SetVar("PROMPT", self.makePrompt().String())
  mod.SetVar("RPROMPT", self.makeRPrompt().String())
  mod.SetVar("TERM_TITLE", self.makeTitle())
  return mod.ToScript()
}
