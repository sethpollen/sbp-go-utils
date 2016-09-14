// Library for constructing prompt strings of the specific form that I like.
package prompt

import "fmt"
import "os"
import "os/exec"
import "os/user"
import "regexp"
import "strings"
import "time"
import "unicode/utf8"
import "github.com/bradfitz/gomemcache/memcache"
import . "github.com/sethpollen/sbp-go-utils/format"
import "github.com/sethpollen/sbp-go-utils/shell"

// Collects information during construction of a prompt string.
type PromptEnv struct {
	Now      time.Time
	Home     string
	Pwd      string
	Hostname string
	// Text to include in the prompt, along with the PWD.
	Info string
	// A secondary info string. Displayed using $RPROMPT.
	Info2 string
	// A short string to place before the final $ in the prompt.
	Flag StyledString
	// Exit code of the last process run in the shell.
	ExitCode int
	// Maximum number of characters which prompt may occupy horizontally.
	Width int
	// Environment variables which should be emitted to the shell which uses this
	// prompt.
	EnvironMod shell.EnvironMod
	// Handle to the local memcache instance.
	Memcache *memcache.Client
}

// Generates a PromptEnv based on current environment variables. The maximum
// number of characters which the prompt may occupy must be passed as 'width'.
func NewPromptEnv(width int, exitCode int, mc *memcache.Client) *PromptEnv {
	var self = new(PromptEnv)
	self.Now = time.Now()
	self.Memcache = mc

	user, err := user.Current()
	if err != nil {
		self.Home = ""
	} else {
		self.Home = user.HomeDir
	}

  // If possible, get the pwd from $PWD, as this usually does the right thing
  // with symlinks (i.e. it shows the path you used to get here, not the
  // actual physical path). If $PWD fails, fall back on os.Getwd().
  self.Pwd = os.Getenv("PWD")
  if len(self.Pwd) == 0 {
    self.Pwd, _ = os.Getwd()
  }

	self.Hostname, _ = os.Hostname()
	self.Info = ""
	self.Info2 = ""
	self.ExitCode = exitCode
	self.Width = width
	self.EnvironMod = *shell.NewEnvironMod()

	return self
}

// Generates a shell prompt string.
func (self *PromptEnv) makePrompt(
	pwdMod func(in StyledString) StyledString) StyledString {
	// If the hostname is a full domain name, remove all but the first domain
	// component.
	// TODO: This info should really be part of the PromptEnv, so we don't
	// have to compute it both here and in MakeTitle.
	var shortHostname = strings.SplitN(self.Hostname, ".", 2)[0]
	var runningOverSsh = (os.Getenv("SSH_TTY") != "")
	var tmuxStatus = getTmuxStatus("ssh")

	// Format the date and time.
	var dateTime = self.Now.Format("01/02 15:04")

	// Construct the prompt text which must precede the PWD.
	var promptBeforePwd StyledString

	// Date and time.
	promptBeforePwd = Stylize(dateTime+" ", Cyan, Bold)

	// Hostname.
	if runningOverSsh {
		promptBeforePwd = append(promptBeforePwd, Stylize("(", Yellow, Dim)...)
	}
	promptBeforePwd = append(promptBeforePwd,
		Stylize(shortHostname, Magenta, Bold)...)
	if runningOverSsh {
		promptBeforePwd = append(promptBeforePwd, Stylize(")", Yellow, Dim)...)
	}

	switch tmuxStatus {
	case TmuxNone:
		// Do nothing.
	case TmuxRunning:
		if os.Getenv("TMUX") != "" {
			// Do nothing; we are already inside tmux.
		} else {
			// Show a subtle % to indicate "running".
			promptBeforePwd = append(promptBeforePwd, Stylize("%%", Yellow, Dim)...)
		}
	case TmuxBell:
		// Show a bold ! to indicate "bell".
		promptBeforePwd = append(promptBeforePwd, Stylize("!", Yellow, Bold)...)
	}
	promptBeforePwd = append(promptBeforePwd, Unstyled(" ")...)

	// Info (if we got one).
	if self.Info != "" {
		promptBeforePwd = append(promptBeforePwd, Stylize("[", White, Dim)...)
		promptBeforePwd = append(promptBeforePwd,
			Stylize(self.Info, White, Bold)...)
		promptBeforePwd = append(promptBeforePwd, Stylize("] ", White, Dim)...)
	}

	// Construct the prompt text which must follow the PWD.
	var promptAfterPwd StyledString

	// Exit code.
	if self.ExitCode != 0 {
		promptAfterPwd = Stylize(fmt.Sprintf(" [%d]", self.ExitCode), Red, Bold)
	}

	// Determine how much space is left for the PWD.
	var pwdWidth = self.Width - len(promptBeforePwd) - len(promptAfterPwd)
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

	var pwdPrompt = self.formatPwd(pwdMod, pwdWidth)

	// Build the complete prompt string.
	var fullPrompt StyledString = promptBeforePwd
	if pwdOnItsOwnLine {
		fullPrompt = append(fullPrompt, promptAfterPwd...)
		fullPrompt = append(fullPrompt, Unstyled("\n")...)
		fullPrompt = append(fullPrompt, pwdPrompt...)
	} else {
		fullPrompt = append(fullPrompt, pwdPrompt...)
		fullPrompt = append(fullPrompt, promptAfterPwd...)
	}
	fullPrompt = append(fullPrompt, Unstyled("\n")...)
	fullPrompt = append(fullPrompt, self.Flag...)
	fullPrompt = append(fullPrompt, Stylize("$ ", Yellow, Bold)...)

	return fullPrompt
}

// Generates a shell RPROMPT string. This will be printed on the right-hand
// side of the second line of the prompt. It will disappear if the user types
// a long command, so it should not be super important. self.Info2 will be the
// content displayed.
func (self *PromptEnv) makeRPrompt() StyledString {
	var rPrompt StyledString
	if self.Info2 != "" {
		rPrompt = Stylize(self.Info2, White, Dim)
	}
	return rPrompt
}

// Generates a terminal emulator title bar string. Similar to a shell prompt
// string, but lacks formatting escapes.
func (self *PromptEnv) makeTitle(
	pwdMod func(in StyledString) StyledString) string {
	var runningOverSsh = (os.Getenv("SSH_TTY") != "")

	var host = ""
	if runningOverSsh {
		var shortHostname = strings.SplitN(self.Hostname, ".", 2)[0]
		host = "(" + shortHostname + ")"
	}

	var info = ""
	if self.Info != "" {
		info = fmt.Sprintf("[%s]", self.Info)
		// Just return the plain self.Info if we detect we are running inside
		// a tmux session (i.e. the $TMUX environment variable is set). When running
		// in tmux, these title strings don't get displayed on the xterm window;
		// they get shown in the tmux tab bar. Space there is constrained, so we
		// don't want to see lengthy PWDs.
		var runningUnderTmux = (os.Getenv("TMUX") != "")
		if runningUnderTmux {
			return info
		}
	}

	var pwdWidth = self.Width - utf8.RuneCountInString(info)
	return host + info + self.formatPwd(pwdMod, pwdWidth).PlainString()
}

// Formats the PWD for use in a prompt. 'mod' is an arbitrary transformation
// to apply to the full PWD before it is (potentially) truncated.
func (self *PromptEnv) formatPwd(
	mod func(in StyledString) StyledString, width int) StyledString {
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

	var styledPwd StyledString = Stylize(pwd, Cyan, Bold)

	if mod != nil {
		styledPwd = mod(styledPwd)
	}

	// Dim slashes in the PWD.
	for i := range styledPwd {
		if styledPwd[i].Text == '/' {
			styledPwd[i].Style.Modifier = Dim
		}
	}

	// Subtract 1 in case we have to include the ellipsis character.
	var pwdRunes = utf8.RuneCountInString(styledPwd.PlainString())
	var start = pwdRunes - (width - 1)
	if start > 0 {
		// Truncate the PWD.
		if start >= pwdRunes {
			// There is no room for the PWD at all.
			styledPwd = make(StyledString, 0)
		} else {
			styledPwd = styledPwd[start:]
			var withEllipsis StyledString = Stylize("…", Cyan, Dim)
			withEllipsis = append(withEllipsis, styledPwd...)
			styledPwd = withEllipsis
		}
	}

	return styledPwd
}

// Renders all the information from this PromptEnv into a shell script which
// may be sourced. The following variables will be set:
//   PROMPT
//   RPROMPT
//   TERM_TITLE
//   ... plus any other variables set in self.EnvironMod.
func (self *PromptEnv) ToScript(
	pwdMod func(in StyledString) StyledString) string {
	// Start by making a copy of the custom EnvironMod.
	var mod = self.EnvironMod
	// Now add our variables to it.
	mod.SetVar("PROMPT", self.makePrompt(pwdMod).String())
	mod.SetVar("RPROMPT", self.makeRPrompt().String())
	mod.SetVar("TERM_TITLE", self.makeTitle(pwdMod))
	// Include the Info string separately, since it is sometimes useful
	// on its own (i.e. as the name of the current repo).
	mod.SetVar("INFO", self.Info)
	return mod.ToScript()
}

// Tmux statuses.
const (
	// The tmux session is not running.
	TmuxNone = iota
	// The tmux session is running but has no bell.
	TmuxRunning
	// The tmux session is running and has an unviewed bell.
	TmuxBell
)

// Returns TmuxNone, TmuxRunning, or TmuxBell based on the status of the given
// tmux session.
func getTmuxStatus(session string) int {
	output, err := exec.Command("tmux",
		"list-windows",
		"-F", "#{session_name} #{window_flags}").Output()
	if err != nil {
		return TmuxNone
	}
	output_str := string(output)

	matched, err := regexp.MatchString(fmt.Sprintf("%s ", session), output_str)
	if err != nil || !matched {
		return TmuxNone
	}

	// The "!" flag indicates a bell.
	matched, err = regexp.MatchString(fmt.Sprintf("%s .*\\!", session),
		output_str)
	if err != nil {
		return TmuxNone
	}

	if matched {
		return TmuxBell
	}
	return TmuxRunning
}
