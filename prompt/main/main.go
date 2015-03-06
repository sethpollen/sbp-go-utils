package main

import "log"
import "github.com/sethpollen/sbp-go-utils/git"
import "github.com/sethpollen/sbp-go-utils/hg"
import "github.com/sethpollen/sbp-go-utils/prompt"

func main() {
	err := prompt.DoMain([]prompt.Module{git.Module(), hg.Module()}, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
