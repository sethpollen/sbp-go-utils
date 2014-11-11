// My utilities for dealing with paths.
package prompt

import "strings"

// Tries to make a version of 'path' which is relative to 'prefix'. If that
// fails, returns 'path' unchanged.
func RelativePath(path string, prefix string) string {
  if !strings.HasPrefix(path, prefix) {
    return path
  }

  // Remove the prefix.
  path = path[len(prefix):]

  // Remove a leading slash, unless the path is just a slash.
  if strings.HasPrefix(path, "/") && path != "/" {
    path = path[1:]
  }

  return path
}
