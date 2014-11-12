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

func TestEvalCommand(t *testing.T) {
  output, err := EvalCommand("/", "echo", "hi")
  if err != nil {
    t.Errorf("Got an error: %v", err)
  }
  if output != "hi" {
    t.Errorf("Expected \"hi\", got \"%s\"", output)
  }

  output, err = EvalCommand("/", "not-a-valid-command")
  if err == nil {
    t.Errorf("Expected an error")
  }
}

func TestEvalCommandAsync(t *testing.T) {
  var outputChan = make(chan string, 1)
  var errChan = make(chan error, 1)
  var output string
  var err error

  EvalCommandAsync(outputChan, errChan, "/", "echo", "hi")
  select {
    case output = <-outputChan:
      if output != "hi" {
        t.Errorf("Expected \"hi\", got \"%s\"", output)
      }
    case err = <-errChan:
      t.Errorf("Got an error: %v", err)
  }

  EvalCommandAsync(outputChan, errChan, "/", "not-a-valid-command")
  select {
    case output = <-outputChan:
      t.Errorf("Didn't get an error")
    case err = <-errChan: // OK.
  }
}
