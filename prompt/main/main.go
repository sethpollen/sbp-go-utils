package main

import "code.google.com/p/sbp-go-utils/prompt"
import "flag"
import "fmt"
import "os"
import "strings"

// Optional flags.
var exitCode = flag.Int("exitcode", 0,
  "Exit code of previous command. If absent, 0 is assumed.")

// Required flags.
var format = flag.String("format", "",
  "Format to output. Possible values are \"prompt\" and \"title\".")
var width = flag.Int("width", -1,
  "Maximum number of characters which the output may occupy.")

func main() {
  flag.Parse()

  if *width < 0 {
    fmt.Fprintln(os.Stderr, "--width must be specified")
    os.Exit(1)
    return
  }
  if *format == "" {
    fmt.Fprintln(os.Stderr, "--format must be specified")
    os.Exit(1)
    return
  }

  var env = prompt.DefaultPromptEnv(*width)
  var info = ""
  var flag = ""

  gitInfo, err := prompt.GetGitInfo(env.Pwd)
  if err == nil {
    info = gitInfo.String()
    flag = "git"
    if strings.HasPrefix(env.Pwd, gitInfo.RepoPath) {
      env.Pwd = env.Pwd[len(gitInfo.RepoPath):]
    }
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
