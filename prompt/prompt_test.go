package prompt

import "strconv"
import "testing"
import "time"
import "code.google.com/p/sbp-go-utils/shell"

var env = PromptEnv{time.Unix(0, 0), "/home/me", "", "myhost.example.com", "",
	"", *NewStyledString(), 0, 100, *shell.NewEnvironMod(), nil}

func assertMakePrompt(t *testing.T, expected string, width int, info string,
	info2 string, pwd string, exitCode int, flag string) {
	var myEnv PromptEnv = env
	myEnv.Pwd = pwd
	myEnv.Info = info
	myEnv.Info2 = info2
	myEnv.Flag.Write(flag)
	myEnv.Width = width
	myEnv.ExitCode = exitCode
	var p = myEnv.makePrompt()
	if p.String() != expected {
		t.Errorf("\nExpected %s"+
			"\nGot      %s",
			strconv.Quote(expected), strconv.Quote(p.String()))
	}
}

func assertMakeTitle(t *testing.T, expected string, info string, pwd string) {
	var myEnv = env
	myEnv.Pwd = pwd
	myEnv.Info = info
	var actual = myEnv.makeTitle()
	if actual != expected {
		t.Errorf("Expected %s\nGot %s",
			strconv.Quote(expected), strconv.Quote(actual))
	}
}

func TestMakePromptSimple(t *testing.T) {
	assertMakePrompt(t,
		"%{\033[1;96m%}12/31 18:00 "+
			"%{\033[1;95m%}myhost "+
			"%{\033[1;96m%}/pw/d"+
			"%{\033[0m%}\nflag%{\033[1;93m%}$ %{\033[0m%}",
		100, "", "", "/pw/d", 0, "flag")
}

func TestMakePromptHomeCollapsing(t *testing.T) {
	assertMakePrompt(t,
		"%{\033[1;96m%}12/31 18:00 "+
			"%{\033[1;95m%}myhost "+
			"%{\033[1;96m%}~/place"+
			"%{\033[0m%}\nflag%{\033[1;93m%}$ %{\033[0m%}",
		100, "", "", "/home/me/place", 0, "flag")
}

func TestMakePromptWithInfoAndExitCode(t *testing.T) {
	assertMakePrompt(t,
		"%{\033[1;96m%}12/31 18:00 "+
			"%{\033[1;95m%}myhost "+
			"%{\033[0;37m%}[%{\033[1;97m%}info%{\033[0;37m%}] "+
			"%{\033[1;96m%}/pw/d"+
			"%{\033[1;91m%} [15]"+
			"\n%{\033[0m%}flag%{\033[1;93m%}$ %{\033[0m%}",
		100, "info", "info2", "/pw/d", 15, "flag")
}

func TestMakePromptTruncatedPwd(t *testing.T) {
	assertMakePrompt(t,
		"%{\033[1;96m%}12/31 18:00 "+
			"%{\033[1;95m%}myhost "+
			"%{\033[0;37m%}[%{\033[1;97m%}info%{\033[0;37m%}] "+
			"%{\033[1;96m%}…6789012345678901234567890"+
			"%{\033[0m%}\nflag%{\033[1;93m%}$ %{\033[0m%}",
		52, "info", "info2", "1234567890123456789012345678901234567890", 0, "flag")
}

func TestMakePromptPwdOnItsOwnLine(t *testing.T) {
	assertMakePrompt(t,
		"%{\033[1;96m%}12/31 18:00 "+
			"%{\033[1;95m%}myhost "+
			"%{\033[0;37m%}[%{\033[1;97m%}info%{\033[0;37m%}] %{\033[0m%}\n"+
			"%{\033[1;96m%}…23456789012345678901234567890"+
			"\n%{\033[0m%}flag%{\033[1;93m%}$ %{\033[0m%}",
		30, "info", "info2", "1234567890123456789012345678901234567890", 0, "flag")
}

func TestMakeTitle(t *testing.T) {
	assertMakeTitle(t, "[info]pwd", "info", "pwd")
}
