// Utilities for dealing with Mercurial repositories. Mercurial is written in
// Python, and as such it is pretty slow. It takes about 50 ms to run
// "hg status" on a repo. This is too long to run inline in my command prompt.
// So this file contains some hacks to query Mercurial repo information without
// invoking any part of the Mercurial codebase.
package hg

import "encoding/json"
import "errors"
import "os"
import "path"
import "strings"
import "code.google.com/p/sbp-go-utils/prompt"
import "code.google.com/p/sbp-go-utils/util"
import "github.com/bradfitz/gomemcache/memcache"

// Encapsulates information about an Hg repo.
type HgInfo struct {
	// Name of this Hg repo.
	RepoName string
  // Full path to the root of this repo.
  RepoPath string
	// Pwd, relative to the root repo path.
	RelativePwd string
  // True if there are uncommitted local changes.
  Dirty bool
}

// Encapsulates information about an Hg repo which is expensive to compute.
type ExpensiveHgInfo struct {
  // True if there are unpushed local commits.
  Unpushed bool
}

func (self *HgInfo) key() string {
  return self.RepoPath
}

func GetHgInfo(pwd string) (*HgInfo, error) {
	repoPath, err := getHgRepoRoot(pwd)
  if err != nil {
    return nil, err
  }

  // Now that we know we are in an Hg repo, it's worth paying the cost to run
  // hg status.
  status, err := util.EvalCommandSync(pwd, "hg", "status")
  if err != nil {
    return nil, err
  }

	var info = new(HgInfo)
	info.RepoName = path.Base(repoPath)
  info.RepoPath = repoPath
	info.RelativePwd = util.RelativePath(pwd, repoPath)
  info.Dirty = (status != "")
	return info, nil
}

func getHgRepoRoot(pwd string) (string, error) {
	repoPath, err := util.SearchParents(pwd, isHgRepoRoot)
	if err != nil {
		return "", errors.New("Not in an Hg repo")
	}
  return repoPath, nil
}

func isHgRepoRoot(pwd string) bool {
	var metaDir = path.Join(pwd, ".hg")
	fileInfo, err := os.Stat(metaDir)
	return err == nil && fileInfo.IsDir()
}

func GetExpensiveHgInfo(pwd string) (*ExpensiveHgInfo, error) {
  // This command takes a while because it actually contacts the remote server.
  outgoing, err := util.EvalCommandSync(pwd, "hg", "outgoing", "--limit=1")
  if err != nil {
    return nil, err
  }

  var info = new(ExpensiveHgInfo)
  info.Unpushed = !strings.Contains(outgoing, "no changes found")
  return info, nil
}

// A prompt.Module that matches any directory inside an Hg repo.
type module struct{}

func (self module) Prepare(env *prompt.PromptEnv) {}

func (self module) Match(env *prompt.PromptEnv, updateCache bool) bool {
	hgInfo, err := GetHgInfo(env.Pwd)
	if err != nil {
		return false
	}

  var expensiveInfo *ExpensiveHgInfo = nil
  if updateCache {
    expensiveInfo, err = GetExpensiveHgInfo(env.Pwd)
    if err == nil {
      value, err := json.Marshal(expensiveInfo)
      if err == nil {
        var item memcache.Item
        item.Key = hgInfo.key()
        item.Value = value
        env.Memcache.Set(&item)
      }
    }
  } else {
    // Just try to read the cache.
    expensiveInfo, _ = readCachedInfo(hgInfo, env.Memcache)
  }

	env.Info = hgInfo.RepoName
  if hgInfo.Dirty || expensiveInfo.Unpushed {
    env.Info += " "
    if expensiveInfo.Unpushed {
      env.Info += "^"
    }
    if hgInfo.Dirty {
      env.Info += "*"
    }
  }
	env.Flag.Style(prompt.Magenta, true)
	env.Flag.Write("hg")
	env.Pwd = hgInfo.RelativePwd
	return true
}

func readCachedInfo(hgInfo *HgInfo, mc *memcache.Client) (
    *ExpensiveHgInfo, error) {
  item, err := mc.Get(hgInfo.key())
  if err != nil {
    return nil, err
  }
  var expensiveInfo = new(ExpensiveHgInfo)
  err = json.Unmarshal(item.Value, expensiveInfo)
  if err != nil {
    return nil, err
  }
  return expensiveInfo, nil
}

func (self module) Description() string {
	return "hg"
}

func Module() module {
	return module{}
}
