package main

import "log"
import "code.google.com/p/sbp-go-utils/git"
import "code.google.com/p/sbp-go-utils/hg"
import "code.google.com/p/sbp-go-utils/prompt"

func main() {
  err := prompt.DoMain([]prompt.Module{git.Module(), hg.Module()})
  if err != nil {
    log.Fatalln(err)
  }
}
