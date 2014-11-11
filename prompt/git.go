// Library for querying info from a local Git repository.

package prompt

import "os/exec"
import "path"
import "strings"

type GitInfo struct {
	// Name of this Git repo.
	RepoName string
  // Full path to the root of this Git repo.
  RepoPath string
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

	status, err := runCommand(pwd, "git", "status", "--porcelain")
	if err != nil {
		return nil, err
	}

	var info = new(GitInfo)
	info.RepoPath = repoPath
	info.RepoName = path.Base(repoPath)

  if strings.HasPrefix(pwd, repoPath) {
    info.RelativePwd = pwd[len(repoPath):]
    // If the relative PWD is more than just "/", remove the leading slash.
    if strings.HasPrefix(info.RelativePwd, "/") && len(info.RelativePwd) >= 2 {
      info.RelativePwd = info.RelativePwd[1:]
    }
  } else {
    // We can't seem to remove the repo prefix, so just preserve the PWD.
    info.RelativePwd = pwd
  }

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

func runCommand(pwd string, name string, arg ...string) (string, error) {
	var cmd = exec.Command(name, arg...)
	cmd.Dir = pwd
	text, err := cmd.Output()
	return strings.TrimSpace(string(text)), err
}
