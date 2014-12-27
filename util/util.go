package util

import "errors"
import "io/ioutil"
import "os/exec"
import "path"
import "strings"
import "github.com/bradfitz/gomemcache/memcache"

// Gets a connection to the local memcached.
func LocalMemcache() *memcache.Client {
  return memcache.New("localhost:11211")
}

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
func EvalCommand(outputChan chan<- string, errorChan chan<- error, pwd string,
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

// Synchronous wrapper around EvalCommand.
func EvalCommandSync(pwd string, name string, args ...string) (string, error) {
  var outputChan = make(chan string)
  var errorChan = make(chan error)
  go EvalCommand(outputChan, errorChan, pwd, name, args...)
  select {
  case err := <-errorChan:
    return "", err
  case output := <-outputChan:
    return output, nil
  }
}

// Returns the shortest prefix of 'p' for which 'test' returns true. Returns
// an error if no prefix matched.
func SearchParents(p string, test func(p string) bool) (string, error) {
	// Build a list of prefixes, beginning with the longest.
	var prefixes []string
	for {
		prefixes = append(prefixes, p)
		var oldP = p
		p = path.Dir(p)
		if p == oldP {
			break
		}
	}

	// Search through the list backwards to find the shortest matching prefix.
	for i := len(prefixes) - 1; i >= 0; i-- {
		var prefix = prefixes[i]
		if test(prefix) {
			return prefix, nil
		}
	}

	return "", errors.New("No prefix matched")
}

// Returns the longest common prefix of strings in 'p'.
func GetLongestCommonPrefix(p []string) string {
  // TODO: make this rune-aware
  if len(p) == 0 {
    return ""
  }
  var longestPrefix = p[0]
  for _, s := range p {
    var minLen = min(len(s), len(longestPrefix))
    var i = 0
    for ; i < minLen; i++ {
      if s[i] != longestPrefix[i] {
        break
      }
    }
    // 'i' now contains the index of the first character that did not match
    // between 's' and 'longestPrefix'.
    if i <= len(longestPrefix) {
      longestPrefix = longestPrefix[0:i]
    }
  }
  return longestPrefix
}

// Splits a path into all of its individual components.
func SplitPath(p string) []string {
  p = path.Clean(p)

  var components []string
  for {
    if p == "" {
      break
    }
    if p == "/" {
      components = append(components, "/")
      break
    }
    parent, child := path.Split(p)
    p = parent
    components = append(components, child)
  }

  // We must now reverse the components slice.
  for i := 0; i < len(components) / 2; i++ {
    var j = (len(components) - 1) - i
    var temp = components[i]
    components[i] = components[j]
    components[j] = temp
  }

  return components
}

// Gets a prefix of 'individual' which uniquely identifies it among all the
// 'population'. Returns the longest such prefix up to 'stopLen', after which
// returns the shortest such prefix.
func getShortestDistinctPrefix(
    individual string, population []string, stopLen int) string {
  var individualRunes = []rune(individual)
  if len(individualRunes) <= stopLen {
    // No shortening is required.
    return individual
  }
  var runesToKeep = stopLen
  // Look over the population, possibly increasing runesToKeep in order to
  // preserve distinctness.
  for _, member := range population {
    for i, memberRune := range member {
      if i >= len(individualRunes) {
        // There are no more runes we can keep from 'individual'.
        break;
      }
      // We should keep this rune.
      runesToKeep = max(runesToKeep, i)
      if individualRunes[i] != memberRune {
        // We found a dissimilarity, so we can stop keeping more runes.
        // Move on to the next member.
        break
      }
    }
  }
  return string(individualRunes[0:runesToKeep])
}

// Takes in a file path and compresses its components to make the
// path occupy fewer characters, but don't compress smaller than 'stopLen'.
// Preserves a long enough prefix for each component to disambiguate it from its
// sibling path elements.
func CompressPath(p string, stopLen int) (string, error) {
  var original = []rune(p)
  var compressed = make([]rune, 0, len(original))

  for pos := 0; pos < len(original); {
    var ch = original[pos]

    if ch == '/' {
      // Preserve path separators.
      compressed = append(compressed, ch)
      pos++
      continue
    }

    // Find the next slash.
    pos++
    for pos < len(original) && original[pos] != '/' {
      pos++
    }
    if pos >= len(original) {
      // This is the last path component.
      // TODO:
    }

    // Grab this path component.
    // var pathUpToThisComponent = original[0:pos]
    // TODO:
  }

  // TODO:
  return string(compressed), nil
}

// Compresses the final component of the path 'p'. Returns the length (in runes)
// of the compressed prefix of the final component.
func compressPathComponent(p string) (int, error) {
  dir, file := path.Split(p)
  var fileRunes = []rune(file)

  entries, err := ioutil.ReadDir(dir)
  if err != nil {
    return -1, err
  }

  // The length of the prefix of 'file' which is required to disambiguate it
  // from all of its siblings in the filesystem.
  var prefixLen = 0
  for _, entry := range entries {
    var sibling = entry.Name()
    if sibling == file {
      // This is the file itself; no disambiguation is necessary.
      continue
    }
    var siblingRunes = []rune(sibling)

    for i := 0; i < min(len(fileRunes), len(siblingRunes)); i++ {
      // Include this character in the compressed prefix.
      prefixLen = max(prefixLen, i)
      if fileRunes[i] != siblingRunes[i] {
        // This character disambiguates.
        break
      }
    }
  }
  return prefixLen, nil
}

func min(a, b int) int {
  if a < b {
    return a
  } else {
    return b
  }
}

func max(a, b int) int {
  if a > b {
    return a
  } else {
    return b
  }
}

