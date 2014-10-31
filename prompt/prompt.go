// Library for constructing prompt strings of the specific form that I like.

package prompt

import "fmt"
import "os"
import "os/user"
import "strconv"
import "strings"
import "time"
import "unicode/utf8"

// Injectable data for testing MakePrompt.
type PromptEnv struct {
	Now      time.Time
	Home     string
	Pwd      string
	Hostname string
	// Maximum number of characters which prompt may occupy horizontally.
	Width int
}

// Generates a PromptEnv based on current environment variables.
func DefaultPromptEnv() *PromptEnv {
	var env = new(PromptEnv)
	env.Now = time.Now()

  user, err := user.Current()
  if err != nil {
    env.Home = ""
  } else {
	  env.Home = user.HomeDir
  }

  env.Pwd, _ = os.Getwd()
	env.Hostname, _ = os.Hostname()

  var widthStr = os.Getenv("COLUMNS")
	width, err := strconv.ParseInt(widthStr, 10, 32)
	if err != nil {
		// Pick a reasonable default.
		width = 100
	}
	env.Width = int(width)

  return env
}

// Generates a shell prompt string.
//   info - An "info" string, which appears next to the PWD.
//   exitCode - The result code of the previous shell command.
//   flag - A short "flag" string, which appears before the final $.
func MakePrompt(env *PromptEnv, info string, exitCode int,
	flag string) *Prompt {
	// If the hostname is a full domain name, remove all but the first domain
	// component.
	var shortHostname = strings.SplitN(env.Hostname, ".", 2)[0]
	var runningOverSsh = (os.Getenv("SSH_TTY") != "")

	// Format the date and time.
	var dateTime = env.Now.Format("01/02 15:04")

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
	if info != "" {
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
	var pwdWidth = env.Width - promptBeforePwd.Len() - promptAfterPwd.Len()
	if pwdWidth < 0 {
		pwdWidth = 0
	}
	var pwdOnItsOwnLine = false
	if env.Width >= 25 && pwdWidth < 25 {
		// Don't cram the PWD into a tiny space; put it on its own line.
		pwdWidth = env.Width
		pwdOnItsOwnLine = true
	}

	var pwd = formatPwd(env, pwdWidth)

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

// Generates a terminal emulator title bar string. Similar to a shell prompt
// string, but lacks formatting escapes.
func MakeTitle(env *PromptEnv, info string) string {
	if info != "" {
		info = fmt.Sprintf("[%s]", info)
	}
	var pwdWidth = env.Width - utf8.RuneCountInString(info)
	return info + formatPwd(env, pwdWidth)
}

// Formats the PWD for use in a prompt.
func formatPwd(env *PromptEnv, width int) string {
	// Perform tilde collapsing on the PWD.
	var home = env.Home
	if strings.HasSuffix(home, "/") {
		home = home[:len(home)-1]
	}
	var pwd = env.Pwd
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
