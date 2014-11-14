// Helper library for implementers of main functions which use build prompts.
package prompt

import "errors"
import "flag"
import "fmt"
import "io/ioutil"
import "log"
import "os"
import "time"

// Required flags.
var width = flag.Int("width", -1,
  "Maximum number of characters which the output may occupy.")

// Optional flags.
var exitCode = flag.Int("exitcode", 0,
  "Exit code of previous command. If absent, 0 is assumed.")
var promptFile = flag.String("prompt_file", "",
  "File to write prompt string to.")
var rPromptFile = flag.String("rprompt_file", "",
  "File to write RPROMPT string to.")
var titleFile = flag.String("title_file", "",
  "File to write title string to.")
var varFile = flag.String("var_file", "",
  "File to write additional environment variables to. This file will be in " +
  "a format appropriate for sourcing in your shell.")

var printTiming = flag.Bool("print_timing", false,
  "True to log diagnostics about how long each part of the program takes.")

var processStart = time.Now()

// Type for a function which may match a PWD and produce an info string.
type PwdMatcher interface {
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

  var env = MakePromptEnv(*width)
  for _, matcher := range matchers {
    LogTime(fmt.Sprintf("Begin matcher \"%s\"", matcher.Description()))
    var done bool = matcher.Match(env)
    LogTime(fmt.Sprintf("End matcher \"%s\"", matcher.Description()))

    if done {
      break
    }
  }

  // Write results.
  if *promptFile != "" {
    var prompt = MakePrompt(env, *exitCode).String()
    err := ioutil.WriteFile(*promptFile, []byte(prompt), 0770)
    if err != nil {
      return err
    }
  }
  if *rPromptFile != "" {
    var rPrompt = MakeRPrompt(env).String()
    err := ioutil.WriteFile(*rPromptFile, []byte(rPrompt), 0770)
    if err != nil {
      return err
    }
  }
  if *titleFile != "" {
    var title = MakeTitle(env)
    err := ioutil.WriteFile(*titleFile, []byte(title), 0770)
    if err != nil {
      return err
    }
  }
  if *varFile != "" {
    file, err := os.Create(*varFile)
    if err != nil {
      return err
    }
    for name, value := range env.Vars {
      fmt.Fprintf(file, "export %s=\"%s\"\n", name, value)
    }
    err = file.Close()
    if err != nil {
      return err
    }
  }

  return nil
}

func LogTime(message string) {
  if !*printTiming {
    return
  }
  var elapsed = time.Now().Sub(processStart)
  log.Printf("(%v) %s\n", elapsed, message)
}
