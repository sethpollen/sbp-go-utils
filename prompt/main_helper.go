// Helper library for implementers of main functions which use build prompts.
package prompt

import "errors"
import "flag"
import "io/ioutil"

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

// Type for a function which may match a PWD and produce an info string.
// If the match succeeds, modifies 'env' in-place and returns true. Otherwise,
// returns false.
type PwdMatcher func(env *PromptEnv) bool

// Entry point. Executes 'matchers' against the current PWD, stopping once one
// of them returns true.
func DoMain(matchers []PwdMatcher) error {
  flag.Parse()

  // Check flags.
  if *width < 0 {
    return errors.New("--width must be specified")
  }

  var env = MakePromptEnv(*width)
  for _, matcher := range matchers {
    if matcher(env) {
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

  return nil
}
