package prompt

import "testing"

func TestEmpty(t *testing.T) {
  var r = RelativePath("", "abc")
  if r != "" {
		t.Errorf("Expected \"\", got \"%s\"", r)
  }
  r = RelativePath("abc", "")
  if r != "abc" {
    t.Errorf("Expected \"abc\", got \"%s\"", r)
  }
}

func TestNotAPrefix(t *testing.T) {
  var r = RelativePath("/a/b/c", "a/b")
  if r != "/a/b/c" {
    t.Errorf("Expected \"/a/b/c\", got \"%s\"", r)
  }
}

func TestPrefix(t *testing.T) {
  var r = RelativePath("/a/b/c", "/a/b")
  if r != "c" {
    t.Errorf("Expected \"c\", got \"%s\"", r)
  }
}

func TestLoneSlash(t *testing.T) {
  var r = RelativePath("/a/b/c", "/a/b/c")
  if r != "/" {
    t.Errorf("Expected \"/\", got \"%s\"", r)
  }
  r = RelativePath("/a/b/c/", "/a/b/c")
  if r != "/" {
    t.Errorf("Expected \"/\", got \"%s\"", r)
  }
}
