// Library for querying info from a local Git repository.

package prompt

import "os/exec"
import "path"
import "strings"

type GitInfo struct {
  // Name of this Git repo.
  Repo string
  // The name of the current branch, or a short hash if we are in a detached
  // head.
  Branch string
}

// Queries GitInfo for the repository that parents 'pwd'.
func GetGitInfo(pwd string) (*GitInfo, error) {
  repoPath, err := runCommand(pwd, "git", "rev-parse", "--show-toplevel")
  if err != nil {
    return nil, err
  }

  branch, err := runCommand(pwd, "git", "symbolic-ref", "HEAD")
  if err == nil {
    var branchParts = strings.Split(branch, "/")
    branch = branchParts[len(branchParts)-1]
  } else {
    // We may be in a detached head. In that case, find the hash of the detached
    // head revision.
    branch, err = runCommand(pwd, "git", "rev-parse", "--short", "HEAD")
    if err != nil {
      return nil, err
    }
  }

  var info = new (GitInfo)
  info.Repo = strings.TrimSpace(path.Base(string(repoPath)))
  info.Branch = strings.TrimSpace(branch)
  return info, nil
}

func runCommand(pwd string, name string, arg ...string) (string, error) {
  var cmd = exec.Command(name, arg...)
  cmd.Dir = pwd
  text, err := cmd.Output()
  return string(text), err
}
