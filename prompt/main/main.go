package main

import "code.google.com/p/sbp-go-utils/prompt"

func main() {
  prompt.DoMain([]prompt.PwdMatcher{prompt.GitMatcher})
}
