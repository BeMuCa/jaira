package identity

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

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

// TestSlugOfPathTraversalIsSafe asserts a name built to escape its parent
// directory cannot survive Slug — the result becomes a path component under
// .jaira/shared/, so ".." and "/" must never appear in it (T-5-03).
func TestSlugOfPathTraversalIsSafe(t *testing.T) {
	got := Slug("../../etc/passwd")
	if strings.Contains(got, "..") || strings.Contains(got, "/") {
		t.Errorf("Slug(%q) = %q must not contain .. or /", "../../etc/passwd", got)
	}
	if got == "" {
		t.Error("Slug() of a traversal string returned empty")
	}
}

// TestSlugOfLeadingDotIsSafe asserts a leading dot does not survive into the
// slug, since a leading-dot directory name has its own special meaning on
// most filesystems.
func TestSlugOfLeadingDotIsSafe(t *testing.T) {
	got := Slug(".hidden")
	if got == "" {
		t.Fatal("Slug(\".hidden\") returned empty")
	}
	if strings.HasPrefix(got, ".") {
		t.Errorf("Slug(%q) = %q must not start with a dot", ".hidden", got)
	}
}

// TestSlugOfEmptyStringFallsBackToNonEmpty asserts an empty name still
// yields a non-empty slug, since an empty path component would silently
// write to the parent directory.
func TestSlugOfEmptyStringFallsBackToNonEmpty(t *testing.T) {
	if got := Slug(""); got == "" {
		t.Error("Slug(\"\") returned empty, which would silently write to the parent directory")
	}
}

// A person is not one string. jaira knew an identity as exactly one, so a ticket
// whose assignee carried a work email read as somebody else's while git carried
// a user.name — and every move on your own ticket needed --force, which trains
// the override rather than protecting anything.
func TestAliasesIncludeTheGitEmail(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("JAIRA_USER", "")
	gitInit(t, dir, "BeMuCa", "berk@example.test")

	got := Aliases(dir)

	if len(got) == 0 || got[0] != "BeMuCa" {
		t.Fatalf("aliases = %v, want the canonical name first", got)
	}
	if !IsMe(dir, "berk@example.test") {
		t.Errorf("the git email is not recognised as me: %v", got)
	}
	if !IsMe(dir, "BeMuCa") {
		t.Errorf("the git name is not recognised as me: %v", got)
	}
}

// The alias file is how a name jaira cannot discover — a work address used by a
// teammate's tooling — becomes yours.
func TestAliasesReadTheAliasFile(t *testing.T) {
	dir := t.TempDir()
	home := t.TempDir()
	t.Setenv("JAIRA_HOME", home)
	t.Setenv("JAIRA_USER", "")
	gitInit(t, dir, "BeMuCa", "berk@example.test")
	if err := os.WriteFile(filepath.Join(home, "identity"),
		[]byte("# other names\n\nberk.calabakan@partner.example\n  BeMuCa  \n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if !IsMe(dir, "berk.calabakan@partner.example") {
		t.Errorf("a configured alias is not recognised: %v", Aliases(dir))
	}
	// Comments, blanks and a duplicate of a name already known add nothing.
	if got := Aliases(dir); len(got) != 3 {
		t.Errorf("aliases = %v, want three distinct names", got)
	}
}

func TestSomeoneElseIsNotMe(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("JAIRA_USER", "")
	gitInit(t, dir, "BeMuCa", "berk@example.test")

	if IsMe(dir, "alexander@example.test") {
		t.Error("a teammate reads as me")
	}
	// An unassigned ticket belongs to nobody, not to everybody.
	if Same("", "") || IsMe(dir, "") {
		t.Error("the empty string reads as a person")
	}
}

func gitInit(t *testing.T, dir, name, email string) {
	t.Helper()
	for _, args := range [][]string{
		{"init", "-q"},
		{"config", "user.name", name},
		{"config", "user.email", email},
	} {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
}

// TestInitials is table-driven: a sync folder is stamped with whoever took the
// ticket off the board, and a folder name should tell a teammate who that was
// at a glance.
func TestInitials(t *testing.T) {
	cases := []struct {
		name, want string
	}{
		{"Alexander Sacharov", "as"},
		{"  anna   marie  ross ", "amr"},
		// A single word has no initials to take, so its slug is used instead —
		// a folder named "b-20260823" tells a teammate nothing, but one named
		// "berk-20260823" does.
		{"berk", "berk"},
		{"", "unnamed"},
	}
	for _, c := range cases {
		if got := Initials(c.name); got != c.want {
			t.Errorf("Initials(%q) = %q, want %q", c.name, got, c.want)
		}
	}
}
