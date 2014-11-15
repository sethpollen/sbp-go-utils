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

// An invoker of this helper must assemble a list of "modules" to be executed
// for each command prompt.
type Module interface {
  // Always invoked on every Module before trying to match any of them.
  Prepare(env *PromptEnv)

  // If the match succeeds, modifies 'env' in-place and returns true. Otherwise,
  // returns false.
  Match(env *PromptEnv) bool

  // Returns a short string describing this Module.
  Description() string
}

// Entry point. Executes 'modules' against the current PWD, stopping once one
// of them returns true.
func DoMain(modules []Module) error {
  flag.Parse()

  LogTime("Begin DoMain")
  defer LogTime("End DoMain")

  // Check flags.
  if *width < 0 {
    return errors.New("--width must be specified")
  }

  var env = NewPromptEnv(*width, *exitCode)
  for _, module := range modules {
    module.Prepare(env)
  }
  for _, module := range modules {
    LogTime(fmt.Sprintf("Begin matching module \"%s\"", module.Description()))
    var done bool = module.Match(env)
    LogTime(fmt.Sprintf("End matching module \"%s\"", module.Description()))

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
