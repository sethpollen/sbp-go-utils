package main

import "fmt"
import "os"
import "code.google.com/p/sbp-go-utils/prompt"

func main() {
  err := prompt.DoMain([]prompt.PwdMatcher{prompt.GitMatcher})
  if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
