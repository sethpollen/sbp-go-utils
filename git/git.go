// Library for querying info from a local Git repository.
package git

import "path"
import "strings"
import "code.google.com/p/sbp-go-utils/prompt"
import "code.google.com/p/sbp-go-utils/util"

type GitInfo struct {
	// Name of this Git repo.
	RepoName string
  // Pwd, relative to the root repo path.
  RelativePwd string
	// The name of the current branch, or a short hash if we are in a detached
	// head.
	Branch string
	// True if there are uncommitted local changes.
	Dirty bool
}

// Queries a GitInfo for the repository that parents 'pwd'. If 'pwd' is not in
// a Git repository, returns an error.
func GetGitInfo(pwd string) (*GitInfo, error) {
	repoPath, err := util.EvalCommand(pwd, "git", "rev-parse", "--show-toplevel")
	if err != nil {
		return nil, err
	}

	branch, err := util.EvalCommand(pwd, "git", "symbolic-ref", "HEAD")
	if err == nil {
		var branchParts = strings.Split(branch, "/")
		branch = branchParts[len(branchParts)-1]
	} else {
		// We may be in a detached head. In that case, find the hash of the detached
		// head revision.
		branch, err = util.EvalCommand(pwd, "git", "rev-parse", "--short", "HEAD")
		if err != nil {
			return nil, err
		}
	}

	status, err := util.EvalCommand(pwd, "git", "status", "--porcelain")
	if err != nil {
		return nil, err
	}

	var info = new(GitInfo)
	info.RepoName = path.Base(repoPath)
  info.RelativePwd = util.RelativePath(pwd, repoPath)
	info.Branch = branch
	info.Dirty = (status != "")
	return info, nil
}

// Formats a GitInfo as a string, suitable for use as an 'info' string in a
// prompt.
func (info *GitInfo) String() string {
	var str = info.RepoName
	if info.RepoName != info.Branch {
		str += ": " + info.Branch
	}
	if info.Dirty {
		str += " *"
	}
	return str
}

// A PwdMatcher that matches any directory inside a Git repo.
var GitMatcher prompt.PwdMatcher = func(env *prompt.PromptEnv) bool {
  gitInfo, err := GetGitInfo(env.Pwd)
  if err != nil {
    return false
  }
  env.Info = gitInfo.String()
  env.Flag = "git"
  env.Pwd = gitInfo.RelativePwd
  return true
}