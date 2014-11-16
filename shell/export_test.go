package shell

import "testing"

func TestToScript(t *testing.T) {
	var mod = NewEnvironMod()
	mod.SetVar("A", "B")
	mod.SetVar("B", "~!@#$%^&*()_+ :;<>,.?/\"'\t\r\n日本")
	mod.UnsetVar("A")
	var actual = mod.ToScript()
	var expected = "unset A\nexport B='~!@#$%^&*()_+ :;<>,.?/\"\\'\t\r\n日本'\n"
	if actual != expected {
		// Find the point where the two strings diverge.
		var actualRunes = toRunes(actual)
		var expectedRunes = toRunes(expected)
		if len(actualRunes) != len(expectedRunes) {
			t.Errorf("Expected %d runes, got %d runes",
				len(expectedRunes), len(actualRunes))
		} else {
			for i := 0; i < len(actualRunes); i++ {
				var actualRune = actualRunes[i]
				var expectedRune = expectedRunes[i]
				if actualRune != expectedRune {
					t.Errorf("At position %d, expected rune 0x%X, got rune 0x%X",
						i, expectedRune, actualRune)
				}
			}
		}
	}
}

func toRunes(text string) []rune {
	var runes []rune
	for _, r := range text {
		runes = append(runes, r)
	}
	return runes
}
