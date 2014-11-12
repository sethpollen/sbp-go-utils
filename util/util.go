package util

import "os/exec"
import "strings"

// Tries to make a version of 'path' which is relative to 'prefix'. If that
// fails, returns 'path' unchanged.
func RelativePath(path string, prefix string) string {
  if !strings.HasPrefix(path, prefix) {
    return path
  }

  // Remove the prefix.
  path = path[len(prefix):]

  if (path == "") || (path == "/") {
    return "/"
  }

  // Remove a leading slash.
  if strings.HasPrefix(path, "/") {
    return path[1:]
  }

  return path
}

// Runs 'name' in 'pwd' with 'args', returning its stdout.
func EvalCommand(pwd string, name string, args ...string) (string, error) {
  var cmd = exec.Command(name, args...)
  cmd.Dir = pwd
  text, err := cmd.Output()
  return strings.TrimSpace(string(text)), err
}

// Async version of EvalCommand. Will either send one string containing the
// command's stdout to 'output' or ssend one error to 'err'.
func EvalCommandAsync(output chan<- string, err chan<- error,
  pwd string, name string, args ...string) {
  myOutput, myErr := EvalCommand(pwd, name, args...)
  if myErr != nil {
    err <- myErr
  } else {
    output <- myOutput
  }
}

