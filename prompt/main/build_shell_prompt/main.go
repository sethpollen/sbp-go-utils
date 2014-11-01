// TODO: there should just be one main file for "prompt". We may want to rename
// the "prompt" package to sbp_shprompt and corp_shprompt. Use flags to indicate
// whether a prompt stirng or shell title is desired.

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
