// A small main program to test git.go.

package main

import "fmt"
import "code.google.com/p/sbp-go-utils/prompt"

func main() {
  fmt.Println(prompt.GetGitInfo("."))
}
