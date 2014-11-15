package prompt

import "strconv"
import "testing"
import "time"

var env = PromptEnv{time.Unix(0, 0), "/home/me", "", "myhost.example.com", "",
                    "", Prompt{}, 100, make(map[string]*string)}

func assertMakePrompt(t *testing.T, expected string, width int, info string,
	info2 string, pwd string, exitCode int, flag string) {
	var myEnv = env
	myEnv.Pwd = pwd
  myEnv.Info = info
  myEnv.Info2 = info2
  myEnv.Flag.Write(flag)
	myEnv.Width = width
	var p = MakePrompt(&myEnv, exitCode)
	if p.String() != expected {
		t.Errorf("Expected %s\nGot %s",
			strconv.Quote(expected), strconv.Quote(p.String()))
	}
}

func assertMakeTitle(t *testing.T, expected string, info string, pwd string) {
  var myEnv = env
  myEnv.Pwd = pwd
  myEnv.Info = info
  var actual = MakeTitle(&myEnv)
  if actual != expected {
    t.Errorf("Expected %s\nGot %s", 
             strconv.Quote(expected), strconv.Quote(actual))
  }
}

func TestMakePromptSimple(t *testing.T) {
	assertMakePrompt(t,
		"%{\033[0m%}%{\033[1;36m%}12/31 18:00 "+
			"%{\033[1;35m%}myhost "+
			"%{\033[1;36m%}/pw/d "+
			"\nflag%{\033[1;33m%}$ %{\033[0m%}",
		100, "", "", "/pw/d", 0, "flag")
}

func TestMakePromptHomeCollapsing(t *testing.T) {
	assertMakePrompt(t,
		"%{\033[0m%}%{\033[1;36m%}12/31 18:00 "+
			"%{\033[1;35m%}myhost "+
			"%{\033[1;36m%}~/place "+
			"\nflag%{\033[1;33m%}$ %{\033[0m%}",
		100, "", "", "/home/me/place", 0, "flag")
}

func TestMakePromptWithInfoAndExitCode(t *testing.T) {
	assertMakePrompt(t,
		"%{\033[0m%}%{\033[1;36m%}12/31 18:00 "+
			"%{\033[1;35m%}myhost "+
			"%{\033[0;37m%}[%{\033[1;37m%}info%{\033[0;37m%}] "+
			"%{\033[1;36m%}/pw/d "+
			"%{\033[1;31m%}[15]"+
			"\nflag%{\033[1;33m%}$ %{\033[0m%}",
		100, "info", "info2", "/pw/d", 15, "flag")
}

func TestMakePromptTruncatedPwd(t *testing.T) {
	assertMakePrompt(t,
		"%{\033[0m%}%{\033[1;36m%}12/31 18:00 "+
			"%{\033[1;35m%}myhost "+
			"%{\033[0;37m%}[%{\033[1;37m%}info%{\033[0;37m%}] "+
			"%{\033[1;36m%}..789012345678901234567890 "+
			"\nflag%{\033[1;33m%}$ %{\033[0m%}",
		52, "info", "info2", "1234567890123456789012345678901234567890", 0, "flag")
}

func TestMakePromptPwdOnItsOwnLine(t *testing.T) {
  assertMakePrompt(t,
    "%{\033[0m%}%{\033[1;36m%}12/31 18:00 "+
      "%{\033[1;35m%}myhost "+
      "%{\033[0;37m%}[%{\033[1;37m%}info%{\033[0;37m%}] \n"+
      "%{\033[1;36m%}..3456789012345678901234567890"+
      "\nflag%{\033[1;33m%}$ %{\033[0m%}",
    30, "info", "info2", "1234567890123456789012345678901234567890", 0, "flag")
}

func TestMakeTitle(t *testing.T) {
  assertMakeTitle(t, "[info]pwd", "info", "pwd")
}

func TestMakeVarScript(t *testing.T) {
  var myEnv = env
  myEnv.UnsetVar("A")
  myEnv.SetVar("B", "hello, world")
  var varText = MakeVarScript(&myEnv)
  
  // Map iteration order is not deterministic.
  const lineA = "unset A\n"
  const lineB = "export B=\"hello, world\"\n"
  
  if varText != lineA + lineB && varText != lineB + lineA {
    t.Errorf("Unexpected var script: %s", varText)
  }
}
