package main

import "log"
import "code.google.com/p/sbp-go-utils/git"
import "code.google.com/p/sbp-go-utils/prompt"

func main() {
  err := prompt.DoMain([]prompt.PwdMatcher{git.GitMatcher})
  if err != nil {
    log.Fatalln(err)
  }
}
