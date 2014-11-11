package main

import "code.google.com/p/sbp-go-utils/prompt"
import "flag"
import "fmt"
import "os"

var exitCode = flag.Int("exitcode", 0,
  "Exit code of previous command. If absent, 0 is assumed.")
var format = flag.String("format", "",
  "Format to output. Possible values are \"prompt\" and \"title\".")
var width = flag.Int("width", -1,
  "Maximum number of characters which the output may occupy. If absent, the " +
  "value of $COLUMNS is used.")

func main() {
  flag.Parse()

  var info = ""
  var flag = ""

  var env = prompt.DefaultPromptEnv()
  if *width >= 0 {
    env.Width = *width
  }

  gitInfo, err := prompt.GetGitInfo(env.Pwd)
  if err == nil {
    info = gitInfo.String()
    flag = "git"
  }

  // Send results to stdout.
  // Interpret the "format" flag.
  switch *format {
    case "prompt": fmt.Print(prompt.MakePrompt(env, info, *exitCode, flag))
    case "title": fmt.Print(prompt.MakeTitle(env, info))
    default:
      fmt.Fprintf(os.Stderr, "Unrecognized value for --format: %s\n", *format)
      os.Exit(1)
  }
}
