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

// Runs 'name' in 'pwd' with 'args'.  Will either send one string containing the
// command's stdout to 'outputChan' or send one error to 'errorChan'.
func EvalCommand(outputChan chan<-string, errorChan chan<-error, pwd string,
    name string, args ...string) {
  var cmd = exec.Command(name, args...)
  cmd.Dir = pwd
  text, err := cmd.Output()
  if err != nil {
    errorChan <- err
  } else {
    outputChan <- strings.TrimSpace(string(text))
  }
}

