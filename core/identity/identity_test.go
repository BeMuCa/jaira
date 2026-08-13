package identity

import "testing"

func TestCurrentHonoursJairaUserAboveEverything(t *testing.T) {
	t.Setenv("JAIRA_USER", "berk")
	if got := Current("."); got != "berk" {
		t.Errorf("Current() = %q, want %q", got, "berk")
	}
}

func TestCurrentFallsThroughWhenJairaUserIsEmpty(t *testing.T) {
	t.Setenv("JAIRA_USER", "")
	t.Setenv("USER", "fallback-user")
	t.Setenv("USERNAME", "")
	t.Setenv("LOGNAME", "")
	// HOME is redirected to an empty directory so this test does not pick up
	// the real developer's ~/.gitconfig, which would otherwise make "git
	// config user.name" succeed even outside a repo and mask the fallback
	// this test is checking.
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	// A directory with no git repo (or one with no configured user.name) falls
	// through to the USER/USERNAME/LOGNAME environment variables.
	if got := Current(t.TempDir()); got != "fallback-user" {
		t.Errorf("Current() = %q, want %q", got, "fallback-user")
	}
}

func TestSlugLowercasesSimpleName(t *testing.T) {
	if got := Slug("BeMuCa"); got != "bemuca" {
		t.Errorf("Slug(%q) = %q, want %q", "BeMuCa", got, "bemuca")
	}
}

func TestSlugOfNameWithSpacesAndAccentsIsPathSafe(t *testing.T) {
	got := Slug("Anna Müller")
	if got == "" {
		t.Fatal("Slug() returned empty for a non-empty name")
	}
	if got[0] == '-' || got[len(got)-1] == '-' {
		t.Errorf("Slug(%q) = %q has a leading or trailing separator", "Anna Müller", got)
	}
	for _, r := range got {
		isLower := r >= 'a' && r <= 'z'
		isDigit := r >= '0' && r <= '9'
		if !isLower && !isDigit && r != '-' {
			t.Errorf("Slug(%q) = %q contains a character outside [a-z0-9-]: %q", "Anna Müller", got, r)
		}
	}
}

func TestSlugOfOnlyPunctuationFallsBackToNonEmpty(t *testing.T) {
	got := Slug("!!!")
	if got == "" {
		t.Error("Slug() of only punctuation returned empty, which would silently write to the parent directory")
	}
}
