package prompt

import "testing"

func TestRelativePathEmpty(t *testing.T) {
  var r = RelativePath("", "abc")
  if r != "" {
		t.Errorf("Expected \"\", got \"%s\"", r)
  }
  r = RelativePath("abc", "")
  if r != "abc" {
    t.Errorf("Expected \"abc\", got \"%s\"", r)
  }
}

func TestRelativePathNotAPrefix(t *testing.T) {
  var r = RelativePath("/a/b/c", "a/b")
  if r != "/a/b/c" {
    t.Errorf("Expected \"/a/b/c\", got \"%s\"", r)
  }
}

func TestRelativePathPrefix(t *testing.T) {
  var r = RelativePath("/a/b/c", "/a/b")
  if r != "c" {
    t.Errorf("Expected \"c\", got \"%s\"", r)
  }
}

func TestRelativePathLoneSlash(t *testing.T) {
  var r = RelativePath("/a/b/c", "/a/b/c")
  if r != "/" {
    t.Errorf("Expected \"/\", got \"%s\"", r)
  }
  r = RelativePath("/a/b/c/", "/a/b/c")
  if r != "/" {
    t.Errorf("Expected \"/\", got \"%s\"", r)
  }
}

func TestRunCommand(t *testing.T) {
  r, err := RunCommand("/", "echo", "hi")
  if err != nil {
    t.Errorf("Got an error: %v", err)
  }
  if r != "hi" {
    t.Errorf("Expected \"hi\", got \"%s\"", r)
  }

  r, err = RunCommand("/", "not-a-valid-command")
  if err == nil {
    t.Errorf("Expected an error")
  }
}
