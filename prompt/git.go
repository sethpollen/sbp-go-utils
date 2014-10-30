// Library for querying info from a local Git repository.

package prompt

import "os/exec"
import "path"

type GitInfo struct {
  RepoName string
}

// Queries GitInfo for the repository that parents 'pwd'.
func GetGitInfo(pwd string) (*GitInfo, error) {
  var repoPath, err = runCommand(pwd, "git", "rev-parse", "--show-toplevel")
  if err != nil {
    repoPath = ""
  }

  var branch, err = runCommand(pwd, "git", "symbolic-ref", "HEAD")
  if err != nil {
    // We may be in a detached head. In that case, find the hash of the detached
    // head revision.
    branch, err = runCommand(pwd, "git", "rev-parse", "HEAD", "--abbrev-ref")
    if err != nil {
      branch = ""
    }
  }

  var info = new (GitInfo)
  info.RepoName = path.Base(string(repoPath))
  return info, nil
}

func runCommand(pwd string, name string, arg ...string) (string, error) {
  var cmd = exec.Command(name, arg)
  cmd.Dir = pwd
  return cmd.Output()
}
