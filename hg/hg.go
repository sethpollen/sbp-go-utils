// Utilities for dealing with Mercurial repositories. Mercurial is written in
// Python, and as such it is pretty slow. It takes about 50 ms to run
// "hg status" on a repo. This is too long to run inline in my command prompt.
// So this file contains some hacks to query Mercurial repo information without
// invoking any part of the Mercurial codebase.
package hg

import "errors"
import "os"
import "path"
import "code.google.com/p/sbp-go-utils/prompt"
import "code.google.com/p/sbp-go-utils/util"

// Encapsulates information about an Hg repo.
type HgInfo struct {
  // Name of this Hg repo.
  RepoName string
  // Pwd, relative to the root repo path.
  RelativePwd string
}

func GetHgInfo(pwd string) (*HgInfo, error) {
  repoPath, err := util.SearchParents(pwd, isHgRepo)
  if err != nil {
    return nil, errors.New("Not in an Hg repo")
  }
  var info = new(HgInfo)
  info.RepoName = path.Base(repoPath)
  info.RelativePwd = util.RelativePath(pwd, repoPath)
  return info, nil
}

func isHgRepo(pwd string) bool {
  var metaDir = path.Join(pwd, ".hg")
  fileInfo, err := os.Stat(metaDir)
  return err == nil && fileInfo.IsDir()
}

// A prompt.Module that matches any directory inside an Hg repo.
type module struct {}

func (self module) Prepare(env *prompt.PromptEnv) {}

func (self module) Match(env *prompt.PromptEnv) bool {
  hgInfo, err := GetHgInfo(env.Pwd)
  if err != nil {
    return false
  }
  env.Info = hgInfo.RepoName
  env.Flag.Style(prompt.Magenta, true)
  env.Flag.Write("hg")
  env.Pwd = hgInfo.RelativePwd
  return true
}

func (self module) Description() string {
  return "hg"
}

func Module() module {
  return module{}
}
