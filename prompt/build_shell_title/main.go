package main

import "code.google.com/p/sbp-go-utils/prompt"
import "fmt"

func main() {
  var env = prompt.DefaultPromptEnv()
  gitInfo, err := prompt.GetGitInfo(env.Pwd)

  var info = ""
  if err == nil {
    info = gitInfo.String()
  }

  fmt.Print(prompt.MakeTitle(env, info))
}
