// Helper library for implementers of main functions which use build prompts.
package prompt

import "flag"
import "fmt"
import "os"

// Required flags.
var format = flag.String("format", "",
  "Format to output. Possible values are \"prompt\" and \"title\".")
var width = flag.Int("width", -1,
  "Maximum number of characters which the output may occupy.")

// Optional flags.
var exitCode = flag.Int("exitcode", 0,
  "Exit code of previous command. If absent, 0 is assumed.")

// Type for a function which may match a PWD and produce an info string.
// If the match succeeds, modifies 'env' in-place and returns true. Otherwise,
// returns false.
type PwdMatcher func(env *PromptEnv) bool

// Entry point. Executes 'matchers' against the current PWD, stopping once one
// of them returns true.
func DoMain(matchers []PwdMatcher) {
  flag.Parse()

  // Check flags.
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

  var env = MakePromptEnv(*width)
  for _, matcher := range matchers {
    if matcher(env) {
      break
    }
  }

  // Send results to stdout.
  // Interpret the "format" flag.
  switch *format {
    case "prompt": fmt.Print(MakePrompt(env, *exitCode))
    case "title": fmt.Print(MakeTitle(env))
    default:
      fmt.Fprintf(os.Stderr, "Unrecognized value for --format: %s\n", *format)
      os.Exit(1)
  }
}

// A PwdMatcher that matches any directory inside a Git repo.
var GitMatcher PwdMatcher = func(env *PromptEnv) bool {
  gitInfo, err := GetGitInfo(env.Pwd)
  if err == nil {
    env.Info = gitInfo.String()
    env.Flag = "git"
    env.Pwd = gitInfo.RelativePwd
    return true
  }
  return false
}
