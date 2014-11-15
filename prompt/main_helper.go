// Helper library for implementers of main functions which use build prompts.
// Prints to stdout a shell script which should then be sourced to set up the
// shell.
package prompt

import "errors"
import "flag"
import "fmt"
import "log"
import "time"

// Required flags.
var width = flag.Int("width", -1,
  "Maximum number of characters which the output may occupy.")

// Optional flags.
var exitCode = flag.Int("exitcode", 0,
  "Exit code of previous command. If absent, 0 is assumed.")

var printTiming = flag.Bool("print_timing", false,
  "True to log diagnostics about how long each part of the program takes.")

var processStart = time.Now()

// Type for a function which may match a PWD and produce an info string.
type PwdMatcher interface {
  // Always invoked on every PwdMatcher before trying to match any of them.
  Prepare(env *PromptEnv)

  // If the match succeeds, modifies 'env' in-place and returns true. Otherwise,
  // returns false.
  Match(env *PromptEnv) bool

  // Returns a short string describing this PwdMatcher.
  Description() string
}

// Entry point. Executes 'matchers' against the current PWD, stopping once one
// of them returns true.
func DoMain(matchers []PwdMatcher) error {
  flag.Parse()

  LogTime("Begin DoMain")
  defer LogTime("End DoMain")

  // Check flags.
  if *width < 0 {
    return errors.New("--width must be specified")
  }

  var env = NewPromptEnv(*width, *exitCode)
  for _, matcher := range matchers {
    matcher.Prepare(env)
  }
  for _, matcher := range matchers {
    LogTime(fmt.Sprintf("Begin matcher \"%s\"", matcher.Description()))
    var done bool = matcher.Match(env)
    LogTime(fmt.Sprintf("End matcher \"%s\"", matcher.Description()))

    if done {
      break
    }
  }

  // Write results.
  fmt.Println(env.ToScript())
  return nil
}

func LogTime(message string) {
  if !*printTiming {
    return
  }
  var elapsed = time.Now().Sub(processStart)
  log.Printf("(%v) %s\n", elapsed, message)
}
