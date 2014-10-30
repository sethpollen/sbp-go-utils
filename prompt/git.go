// Library for querying info from a local Git repository.

package prompt

import "os/exec"
import "path"

type GitInfo struct {
  RepoName string
}

// Queries GitInfo for the repository that parents 'pwd'.
func GetGitInfo(pwd string) (*GitInfo, error) {
  var cmd = exec.Command("git", "rev-parse", "--show-toplevel")
  cmd.Dir = pwd
  var repoPath, err = cmd.Output()
  if err != nil {
    return nil, err
  }
  var info = new (GitInfo)
  info.RepoName = path.Base(string(repoPath))
  return info, nil
}
