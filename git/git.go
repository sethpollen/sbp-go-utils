// Library for querying info from a local Git repository.
package git

import "bufio"
import "path"
import "regexp"
import "strings"
import . "github.com/sethpollen/sbp-go-utils/format"
import "github.com/sethpollen/sbp-go-utils/prompt"
import "github.com/sethpollen/sbp-go-utils/util"

type GitInfo struct {
	// Name of this Git repo.
	RepoName string
	// Pwd, relative to the root repo path.
	RelativePwd string
	// The name of the current branch, or a short hash if we are in a detached
	// head.
	Branch string
	// True iff there are uncommitted local changes.
	Dirty bool
	// True iff there are unpushed local commits.
	Ahead bool
}

// Regex to match the "branch" line from git status --branch --porcelain. If
// this matches, the local branch is ahead of the remote branch.
var statusBranchAheadRegex = regexp.MustCompile("^\\#\\# .* \\[ahead [0-9]+\\]$")

// Queries a GitInfo for the repository that parents 'pwd'. If 'pwd' is not in
// a Git repository, returns an error.
func GetGitInfo(pwd string) (*GitInfo, error) {
	repoPath, err :=
		util.EvalCommandSync(pwd, "git", "rev-parse", "--show-toplevel")
	if err != nil {
		return nil, err
	}

	branch, err := util.EvalCommandSync(pwd, "git", "symbolic-ref", "HEAD")
	if err == nil {
		var branchParts = strings.Split(branch, "/")
		branch = branchParts[len(branchParts)-1]
	} else {
		// We may be in a detached head. In that case, find the hash of the detached
		// head revision.
		branch, err =
			util.EvalCommandSync(pwd, "git", "rev-parse", "--short", "HEAD")
		if err != nil {
			return nil, err
		}
	}

	status, err :=
		util.EvalCommandSync(pwd, "git", "status", "--branch", "--porcelain")
	if err != nil {
		return nil, err
	}

	var info = new(GitInfo)
	info.RepoName = path.Base(repoPath)
	info.RelativePwd = util.RelativePath(pwd, repoPath)
	info.Branch = branch

	info.Dirty = false
	info.Ahead = false

	// Parse the git status result.
	var scanner = bufio.NewScanner(strings.NewReader(status))
	// Stop looping of we set both Ahead and Dirty to true.
	for scanner.Scan() && !(info.Ahead && info.Dirty) {
		var line = scanner.Text()
		if strings.HasPrefix(line, "## ") {
			// This is the "branch" line.
			if statusBranchAheadRegex.FindStringIndex(line) != nil {
				info.Ahead = true
			}
		} else {
			// This is not the "branch" line, so it must indicate that a file is
			// dirty.
			info.Dirty = true
		}
	}

	return info, nil
}

// Formats a GitInfo as a string, suitable for use as an 'info' string in a
// prompt.
func (info *GitInfo) String() string {
	var str = info.RepoName
	if info.Branch != "master" {
		str += ": " + info.Branch
	}
	if info.Ahead || info.Dirty {
		str += " "
		if info.Ahead {
			str += "^"
		}
		if info.Dirty {
			str += "*"
		}
	}
	return str
}

// A prompt.Modlue that matches any directory inside a Git repo.
type module struct{}

func (self module) Prepare(env *prompt.PromptEnv) {}

func (self module) Match(env *prompt.PromptEnv, updateCache bool) bool {
	gitInfo, err := GetGitInfo(env.Pwd)
	if err != nil {
		return false
	}
	env.Info = gitInfo.String()
	env.Flag = append(env.Flag, Stylize("git", Red, Intense)...)
	env.Pwd = gitInfo.RelativePwd
	return true
}

func (self module) Description() string {
	return "git"
}

func Module() module {
	return module{}
}
