// Helper library for implementers of main functions which use build prompts.
// Prints to stdout a shell script which should then be sourced to set up the
// shell.
package prompt

import "flag"
import "fmt"
import "log"
import "time"
import . "github.com/sethpollen/sbp-go-utils/format"
import "github.com/sethpollen/sbp-go-utils/util"

var width = flag.Int("width", 100,
	"Maximum number of characters which the output may occupy.")
var updateCache = flag.Bool("update_cache", false,
	"True to perform expensive operations and update the cache.")
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
	// returns false. If 'updateCache' is true, this call should do expensive
	// operations and write their results to the cache.
	Match(env *PromptEnv, updateCache bool) bool

	// Returns a short string describing this Module.
	Description() string
}

// Entry point. Executes 'modules' against the current PWD, stopping once one
// of them returns true. 'pwdMod' is an optional function to apply additional
// formatting to the PWD before it is printed.
func DoMain(modules []Module,
	pwdMod func(in StyledString) StyledString) error {
	flag.Parse()

	LogTime("Begin DoMain")

	var env = NewPromptEnv(*width, *exitCode, util.LocalMemcache())
	for _, module := range modules {
		LogTime(fmt.Sprintf("Begin Prepare(\"%s\")", module.Description()))
		module.Prepare(env)
		LogTime(fmt.Sprintf("Begin Prepare(\"%s\")", module.Description()))
	}
	for _, module := range modules {
		LogTime(fmt.Sprintf("Begin Match(\"%s\")", module.Description()))
		var done bool = module.Match(env, *updateCache)
		LogTime(fmt.Sprintf("End Match(\"%s\")", module.Description()))

		if done {
			break
		}
	}

	// Report the amount of time we spent generating the prompt.
	var elapsed = time.Now().Sub(processStart)
	env.EnvironMod.SetVar("PROMPT_GENERATION_SECONDS",
		fmt.Sprintf("%f", elapsed.Seconds()))

	// Write results.
	fmt.Println(env.ToScript(pwdMod))

	LogTime("End DoMain")
	return nil
}

func LogTime(message string) {
	if !*printTiming {
		return
	}
	var elapsed = time.Now().Sub(processStart)
	log.Printf("(%v) %s\n", elapsed, message)
}
