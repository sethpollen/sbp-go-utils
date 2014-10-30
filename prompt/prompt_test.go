package prompt

import "strconv"
import "testing"
import "time"

var env = PromptEnv{time.Unix(0, 0), "/home/me", "", "myhost.example.com"}

func assertMakePrompt(t *testing.T, expected string, maxWidth int, info string,
	pwd string, exitCode int, flag string) {
	var myEnv = env
	myEnv.pwd = pwd
	var p = MakePrompt(&myEnv, maxWidth, info, exitCode, flag)
	if p.String() != expected {
		t.Errorf("Expected %s\nGot %s",
			strconv.Quote(expected), strconv.Quote(p.String()))
	}
}

func TestMakePromptSimple(t *testing.T) {
	assertMakePrompt(t,
		"%{\033[0m%}%{\033[1;36m%}12/31 18:00 "+
			"%{\033[1;35m%}myhost "+
			"%{\033[1;36m%}/pw/d "+
			"%{\033[1;33m%}\nflag$ %{\033[0m%}",
		100, "", "/pw/d", 0, "flag")
}

func TestMakePromptHomeCollapsing(t *testing.T) {
	assertMakePrompt(t,
		"%{\033[0m%}%{\033[1;36m%}12/31 18:00 "+
			"%{\033[1;35m%}myhost "+
			"%{\033[1;36m%}~/place "+
			"%{\033[1;33m%}\nflag$ %{\033[0m%}",
		100, "", "/home/me/place", 0, "flag")
}

func TestMakePromptWithInfoAndExitCode(t *testing.T) {
	assertMakePrompt(t,
		"%{\033[0m%}%{\033[1;36m%}12/31 18:00 "+
			"%{\033[1;35m%}myhost "+
			"%{\033[0;37m%}[%{\033[1;37m%}info%{\033[0;37m%}] "+
			"%{\033[1;36m%}/pw/d "+
			"%{\033[1;31m%}[15]"+
			"%{\033[1;33m%}\nflag$ %{\033[0m%}",
		100, "info", "/pw/d", 15, "flag")
}

func TestMakePromptTruncatedPwd(t *testing.T) {
	assertMakePrompt(t,
		"%{\033[0m%}%{\033[1;36m%}12/31 18:00 "+
			"%{\033[1;35m%}myhost "+
			"%{\033[0;37m%}[%{\033[1;37m%}info%{\033[0;37m%}] "+
			"%{\033[1;36m%}..789012345678901234567890 "+
			"%{\033[1;33m%}\nflag$ %{\033[0m%}",
		52, "info", "1234567890123456789012345678901234567890", 0, "flag")
}

func TestMakePromptPwdOnItsOwnLine(t *testing.T) {

}
