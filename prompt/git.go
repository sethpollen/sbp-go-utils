// Library for querying info from a local Git repository.

package prompt

import "fmt"
import "os/exec"

type GitInfo struct {
  repoName string
}

// Queries GitInfo for the repository that parents 'pwd'.
func GetGitInfo(pwd string) *GitInfo {
  var repoPath, err =
    exec.Command(['git', 'rev-parse', '--show_toplevel']).Output()
  if err != nil {
    return nil, err
  }
  var info = new (GitInfo)
  return info, nil
}
