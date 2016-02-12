package prompt

// TODO: These tests currently fail if you are running tmux at the same time.
// We need to make them more hermetic.

import "strconv"
import "testing"
import "time"
import . "github.com/sethpollen/sbp-go-utils/format"
import "github.com/sethpollen/sbp-go-utils/shell"

var env = PromptEnv{time.Unix(0, 0), "/home/me", "", "myhost.example.com", "",
	"", make(StyledString, 0), 0, 100, *shell.NewEnvironMod(), nil}

func assertMakePrompt(t *testing.T, expected string, width int, info string,
	info2 string, pwd string, exitCode int, flag string) {
	var myEnv PromptEnv = env
	myEnv.Pwd = pwd
	myEnv.Info = info
	myEnv.Info2 = info2
	myEnv.Flag = Stylize(flag, Red, Intense)
	myEnv.Width = width
	myEnv.ExitCode = exitCode
	var p = myEnv.makePrompt(nil)
	if p.PlainString() != expected {
		t.Errorf("\nExpected %s"+
			"\nGot      %s",
			strconv.Quote(expected), strconv.Quote(p.PlainString()))
	}
}

func assertMakeTitle(t *testing.T, expected string, info string, pwd string) {
	var myEnv = env
	myEnv.Pwd = pwd
	myEnv.Info = info
	var actual = myEnv.makeTitle(nil)
	if actual != expected {
		t.Errorf("Expected %s\nGot %s",
			strconv.Quote(expected), strconv.Quote(actual))
	}
}

func TestMakePromptSimple(t *testing.T) {
	assertMakePrompt(t,
		"12/31 18:00 myhost /pw/d\nflag$ ",	100, "", "", "/pw/d", 0, "flag")
}

func TestMakePromptHomeCollapsing(t *testing.T) {
	assertMakePrompt(t,
		"12/31 18:00 myhost ~/place\nflag$ ", 100, "", "", "/home/me/place", 0,
    "flag")
}

func TestMakePromptWithInfoAndExitCode(t *testing.T) {
	assertMakePrompt(t,
		"12/31 18:00 myhost [info] /pw/d [15]\nflag$ ", 100, "info", "info2",
    "/pw/d", 15, "flag")
}

func TestMakePromptTruncatedPwd(t *testing.T) {
	assertMakePrompt(t,
		"12/31 18:00 myhost [info] …6789012345678901234567890\nflag$ ",
		52, "info", "info2", "1234567890123456789012345678901234567890", 0, "flag")
}

func TestMakePromptPwdOnItsOwnLine(t *testing.T) {
	assertMakePrompt(t,
		"12/31 18:00 myhost [info] \n…23456789012345678901234567890\nflag$ ",
		30, "info", "info2", "1234567890123456789012345678901234567890", 0, "flag")
}

func TestMakeTitle(t *testing.T) {
	assertMakeTitle(t, "[info]pwd", "info", "pwd")
}
