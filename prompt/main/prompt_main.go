package main

import "code.google.com/p/sbp-go-utils/prompt"
import "flag"
import "fmt"

var exitCode = flag.Int("exitcode", 0, "Exit code of previous command.")

func main() {
  flag.Parse()

  var env = prompt.DefaultPromptEnv()
  gitInfo, err := prompt.GetGitInfo(env.Pwd)

  var info = ""
  var flag = ""
  if err == nil {
    info = gitInfo.String()
    flag = "git"
  }

  fmt.Print(prompt.MakePrompt(env, info, *exitCode, flag))
}
