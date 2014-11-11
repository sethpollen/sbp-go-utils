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
