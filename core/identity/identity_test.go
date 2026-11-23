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
	// USERPROFILE is the same redirection on Windows: os.UserHomeDir reads it
	// there and ignores HOME entirely.
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
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

// The board's people list decides the acting name, found by any of git's two
// names or by JAIRA_USER; a person the list does not know acts as before.
func TestCurrentIsTheJairaNameThePeopleListGives(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	dir := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	git("init", "--quiet")
	git("config", "user.name", "BeMuCa")
	git("config", "user.email", "berk@example.test")
	t.Setenv("JAIRA_USER", "")
	if got := Current(dir); got != "BeMuCa" {
		t.Fatalf("without a list: %q", got)
	}
	list := filepath.Join(dir, "people")
	if err := os.WriteFile(list, []byte("berkc: bemuca\nother: someone\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git("update-ref", "refs/jaira/people", strings.TrimSpace(gitOut(t, dir, "commit-tree", strings.TrimSpace(gitOut(t, dir, "write-tree-for", list)), "-m", "people")))
	cases := []struct{ name, email, user, want string }{
		{"BeMuCa", "x@example.test", "", "berkc"},                     // by user.name, any case
		{"Berk on laptop", "berk@example.test", "", "Berk on laptop"}, // email unlisted
		{"nobody", "x@example.test", "someone", "other"},              // JAIRA_USER is looked up too
		{"nobody", "x@example.test", "", "nobody"},
		{"nobody", "someone", "", "other"},      // by user.email
		{"nobody", "someone", "grace", "grace"}, // JAIRA_USER is not overruled by git's email                    // unknown stays itself
	}
	for _, c := range cases {
		git("config", "user.name", c.name)
		git("config", "user.email", c.email)
		t.Setenv("JAIRA_USER", c.user)
		if got := Current(dir); got != c.want {
			t.Errorf("%+v: Current = %q", c, got)
		}
	}
	git("config", "user.name", "BeMuCa")
	git("config", "user.email", "berk@example.test")
	t.Setenv("JAIRA_USER", "")
	if a := Aliases(dir); a[0] != "berkc" || !IsMe(dir, "BEMUCA") || !IsMe(dir, "berkc") || IsMe(dir, "someone") {
		t.Errorf("aliases = %v", a)
	}
}

// Under JAIRA_USER only that name and the people list's names for it mean me:
// git's email and the alias file are the checkout owner's, and grace must not
// see the owner's tickets as hers.
func TestAliasesUnderJairaUserLeaveOutTheOwnersNames(t *testing.T) {
	dir := t.TempDir()
	home := t.TempDir()
	t.Setenv("JAIRA_HOME", home)
	gitInit(t, dir, "BeMuCa", "berk@example.test")
	if err := os.WriteFile(filepath.Join(home, "identity"), []byte("berk.c@partner.example\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	list := filepath.Join(dir, "people")
	if err := os.WriteFile(list, []byte("grace: grace@example.test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	commit := strings.TrimSpace(gitOut(t, dir, "commit-tree", strings.TrimSpace(gitOut(t, dir, "write-tree-for", list)), "-m", "people"))
	gitOut(t, dir, "update-ref", "refs/jaira/people", commit)
	t.Setenv("JAIRA_USER", "grace")

	for _, owner := range []string{"berk@example.test", "BeMuCa", "berk.c@partner.example"} {
		if IsMe(dir, owner) {
			t.Errorf("%q reads as grace: %v", owner, Aliases(dir))
		}
	}
	for _, mine := range []string{"grace", "GRACE@example.test"} {
		if !IsMe(dir, mine) {
			t.Errorf("%q is not grace: %v", mine, Aliases(dir))
		}
	}
}

// gitOut runs git in dir; "write-tree-for <file>" is shorthand for a tree that
// holds that one file under the name "people".
func gitOut(t *testing.T, dir string, args ...string) string {
	t.Helper()
	if args[0] == "write-tree-for" {
		blob := strings.TrimSpace(gitOut(t, dir, "hash-object", "-w", args[1]))
		cmd := exec.Command("git", "-C", dir, "mktree")
		cmd.Stdin = strings.NewReader("100644 blob " + blob + "\tpeople\n")
		out, err := cmd.Output()
		if err != nil {
			t.Fatal(err)
		}
		return string(out)
	}
	out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).Output()
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
	return string(out)
}
