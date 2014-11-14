// Library for constructing prompt strings of the specific form that I like.
package prompt

import "bytes"
import "fmt"
import "os"
import "os/user"
import "strings"
import "time"
import "unicode/utf8"

// Collects information during construction of a prompt string.
type PromptEnv struct {
	Now      time.Time
	Home     string
	Pwd      string
	Hostname string
  // Text to include in the prompt, along with the PWD.
  Info     string
  // A secondary info string. Displayed using $RPROMPT.
  Info2    string
  // A short string to place before the final $ in the prompt.
  Flag     string
	// Maximum number of characters which prompt may occupy horizontally.
	Width    int
  // Environment variables which should be emitted to the shell which uses this
  // prompt. Values will not be escaped, so don't put any weird characters in
  // here. Values will be quoted. Entries with nil values will be unset in the
  // shell.
  Vars     map[string]*string
}

func (self *PromptEnv) SetVar(name string, value string) {
  self.Vars[name] = &value
}

func (self *PromptEnv) UnsetVar(name string) {
  self.Vars[name] = nil
}

// Generates a PromptEnv based on current environment variables. The maximum
// number of characters which the prompt may occupy must be passed as 'width'.
func MakePromptEnv(width int) *PromptEnv {
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
  env.Info = ""
  env.Info2 = ""
  env.Flag = ""
	env.Width = width
  env.Vars = make(map[string]*string)

  return env
}

// Generates a shell prompt string.
//   exitCode - The result code of the previous shell command.
func MakePrompt(env *PromptEnv, exitCode int) *Prompt {
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
	if env.Info != "" {
		promptBeforePwd.Style(White, false)
		promptBeforePwd.Write("[")
		promptBeforePwd.Style(White, true)
		promptBeforePwd.Write(env.Info)
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
	if pwdWidth < 20 && utf8.RuneCountInString(env.Pwd) >= 20 && env.Width >= 20 {
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
	fullPrompt.Write("\n" + env.Flag + "$ ")

	return fullPrompt
}

// Generates a shell RPROMPT string. This will be printed on the right-hand
// side of the second line of the prompt. It will disappear if the user types
// a long command, so it should not be super important. env.Info2 will be the
// content displayed.
// TODO: unit test
func MakeRPrompt(env *PromptEnv) *Prompt {
  var rPrompt = new(Prompt)
  if env.Info2 != "" {
    rPrompt.Style(White, false)
    rPrompt.Write(env.Info2)
  }
  return rPrompt
}

// Generates a terminal emulator title bar string. Similar to a shell prompt
// string, but lacks formatting escapes.
func MakeTitle(env *PromptEnv) string {
  var info = ""
	if env.Info != "" {
		info = fmt.Sprintf("[%s]", env.Info)
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

func MakeVarScript(env *PromptEnv) string {
  var buf = bytes.NewBufferString("")
  for name, value := range env.Vars {
    if value == nil {
      fmt.Fprintf(buf, "unset %s\n", name)
    } else {
      fmt.Fprintf(buf, "export %s=\"%s\"\n", name, *value)
    }
  }
  return buf.String()
}
