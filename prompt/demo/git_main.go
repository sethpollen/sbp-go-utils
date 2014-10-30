// A small main program to test git.go.

package main

import "fmt"
import "code.google.com/p/sbp-go-utils/prompt"

func main() {
	info, err := prompt.GetGitInfo(".")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(info)
	}
}
