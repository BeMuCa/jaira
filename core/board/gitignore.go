package board

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BeMuCa/jaira/core/ticket"
)

// IgnoreLine is the entry that keeps a board private.
const IgnoreLine = "/.jaira/"

// LanesIgnoreLine keeps a project's own lane directory out of git even once
// the rest of the board is shared. .jaira/lanes/ is scoped to this machine
// (see D-03) and was never meant to reach a teammate. This is an ordinary
// ignore rule on a path outside any ignored tree, not a negation pattern: by
// the time this line matters, RemoveIgnore has already stopped ignoring
// .jaira/ itself, so there is no ignored parent to negate out of.
const LanesIgnoreLine = "/.jaira/lanes/"

// MilestonesIgnoreLine keeps the milestone files out of git on a shared board.
// A milestone travels on its ref, refs/jaira/milestones/<ID>, and every
// command keeps the file in step with it; a committed copy only rides along
// with a branch, falls behind the ref, and comes back in every merge as a
// conflict nobody needs.
const MilestonesIgnoreLine = "/.jaira/milestones/"

// privateComment heads the entry; RemoveIgnore takes it out with the line,
// recognising it by privateMarker.
const (
	privateComment = "# jaira board — private to this machine. Run 'jaira share' to publish it."
	privateMarker  = "jaira board — private"
	lanesComment   = "# this project's lane files are this machine's, not the team's"
	lanesMarker    = "lane files are this machine's"
	// milestonesComment heads MilestonesIgnoreLine; RemoveMilestonesIgnore
	// knows it by milestonesMarker.
	milestonesComment = "# jaira milestones travel on their refs; the files are working copies, never committed"
	milestonesMarker  = "jaira milestones travel on their refs"
)

// A board starts private. Committing tickets is what makes a board visible to
// everyone who can read the repository, and that should be a decision rather than
// a default — the tool has no way to know whether these notes are ready to be
// seen, or whether this repository is even yours to publish into.
//
// Private and shared differ only in whether `.jaira/` is ignored. Nothing about
// the tickets themselves changes, so going public later is not a migration.
//
// A private board is excluded in <git-common-dir>/info/exclude, not in the
// .gitignore of the working tree. The .gitignore is tracked: writing to it
// leaves every new worktree with a modified file that belongs to no change,
// and a repository that guards its pushes against stray changes refuses
// every push. info/exclude lives in the git directory all worktrees share,
// applies to each of them, and never appears in a diff. Only outside a git
// repository, where there is no info/exclude, does the entry go into
// .gitignore — untracked there, and still in force if 'git init' comes later.

// gitPlace is where git keeps the exclude file for root, and the pattern that
// names root's .jaira/ in it: exclude patterns are relative to the top of the
// working tree, so a board in a subdirectory is "/sub/.jaira/". Asked of git,
// like commonGitDir in core/ticket, rather than resolving a worktree's .git
// file by hand. Empty outside a repository or without git.
func gitPlace(root string) (exclude, line string) {
	_, exclude, line = place(root)
	return exclude, line
}

// place is gitPlace and whether root lies inside a working tree, from one
// git process.
func place(root string) (inRepo bool, exclude, line string) {
	inRepo, exclude, prefix := rawPlace(root)
	if exclude == "" {
		return inRepo, "", ""
	}
	return inRepo, exclude, entryLine(prefix)
}

// rawPlace is place with root's prefix below the top of the working tree as
// git prints it, not yet made an exclude pattern. --is-inside-work-tree goes
// last: git before 2.31 prints the --path-format it does not know back as an
// extra first line, and the last line stays the answer either way.
func rawPlace(root string) (inRepo bool, exclude, prefix string) {
	if _, err := exec.LookPath("git"); err != nil {
		return false, "", ""
	}
	out, err := exec.Command("git", "-C", root, "rev-parse", "--path-format=absolute", "--git-common-dir", "--show-prefix", "--is-inside-work-tree").Output()
	if err != nil {
		return false, "", ""
	}
	return parseProbe(string(out))
}

// parseProbe reads rawPlace's git output: parsePlace's lines, then "true"
// or "false" for --is-inside-work-tree. A git directory that is no working
// tree still has its exclude file.
func parseProbe(out string) (inRepo bool, exclude, prefix string) {
	lines := strings.Split(strings.TrimRight(out, "\r\n"), "\n")
	inRepo = strings.TrimSpace(lines[len(lines)-1]) == "true"
	exclude, prefix = parsePlace(strings.Join(lines[:len(lines)-1], "\n"))
	return inRepo, exclude, prefix
}

// parsePlace reads two lines of git output: the common git dir —
// "C:/repo/.git" on Windows, which filepath makes a path of — and the prefix
// of root below the top of the working tree, "" or "sub/dir/", always with
// forward slashes, as exclude patterns are written. Only the line ending
// comes off a line: a space at either end belongs to a directory's name
// (" sub/"), and an entry without it excludes nothing.
func parsePlace(out string) (exclude, prefix string) {
	parts := strings.Split(strings.TrimRight(out, "\r\n"), "\n")
	common := strings.TrimSuffix(parts[0], "\r")
	if common == "" {
		return "", ""
	}
	if len(parts) > 1 {
		prefix = strings.TrimSuffix(parts[1], "\r")
	}
	// git before 2.31 does not know --path-format and prints it back as the
	// first line, exit 0: anything that is not an absolute path is not the
	// answer, and the board falls back to .gitignore rather than to a
	// directory of that name inside the repository.
	common = filepath.FromSlash(common)
	if !filepath.IsAbs(common) {
		return "", ""
	}
	return filepath.Join(common, "info", "exclude"), prefix
}

// entryLine is the exclude pattern for the board below prefix.
func entryLine(prefix string) string {
	return "/" + globEscape.Replace(prefix) + ".jaira/"
}

// globEscape keeps a directory name such as "a[1]" literal in an exclude
// pattern, where git would read [1] as a character class.
var globEscape = strings.NewReplacer(`\`, `\\`, "*", `\*`, "?", `\?`, "[", `\[`)

// IgnoreFileName is the file a private board's entry goes into, as a person
// reads it: info/exclude in a git repository — absolute where it lies outside
// root, as a worktree's shared git dir does — and .gitignore outside one.
func IgnoreFileName(root string) string {
	path, _ := gitPlace(root)
	if path == "" {
		return ".gitignore"
	}
	// git prints the path resolved (/private/var on macOS, the long name on
	// Windows) while root may come through a symlink or an 8.3 short name;
	// resolve both, or the relative name falls back to the absolute path.
	if rel, err := filepath.Rel(resolved(root), resolved(path)); err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return filepath.ToSlash(rel)
	}
	return path
}

// resolved is p with its symlinks evaluated, or p itself where that fails.
func resolved(p string) string {
	if r, err := filepath.EvalSymlinks(p); err == nil {
		return r
	}
	return p
}

// boardLine reports whether line excludes the board in a file whose patterns
// are relative to the board's root (.gitignore) — or, given the anchored
// pattern for a subdirectory, in info/exclude.
func boardLine(line, anchored string) bool {
	switch t := strings.TrimSpace(line); t {
	case IgnoreLine, ".jaira/", ".jaira":
		return anchored == "" || anchored == IgnoreLine
	default:
		return anchored != "" && t == anchored
	}
}

func lists(path, anchored string) bool {
	b, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(b), "\n") {
		if boardLine(line, anchored) {
			return true
		}
	}
	return false
}

// Ignored reports whether jaira's entry for the board is written: in
// info/exclude, or as a line in the root .gitignore — an older jaira wrote it
// there, and a person may have written it there by hand. It is not whether
// git ignores the board; GitIgnores answers that.
func Ignored(root string) bool {
	if InGitignore(root) {
		return true
	}
	exclude, line := gitPlace(root)
	return exclude != "" && lists(exclude, line)
}

// Probe answers what every command on a board with a merge driver needs
// before it binds that driver — is root inside a working tree, and does
// jaira's entry exclude the board (Ignored) — with one git process at most.
// Asking IsRepo and Ignored one after the other costs two, and runs on every
// command of a shared board; the entry in the root .gitignore, read from
// disk, settles "ignored" without git at all — inRepo is then false, as
// nothing about an ignored board needs it.
//
// inRepo is therefore only an answer where ignored is false: true then means
// root lies inside a working tree, false that it does not. Where ignored is
// true, inRepo says nothing about the repository — a caller that needs it for
// an ignored board asks gitrepo.Repo.IsRepo.
func Probe(root string) (inRepo, ignored bool) {
	if InGitignore(root) {
		return false, true
	}
	inRepo, exclude, line := place(root)
	return inRepo, exclude != "" && lists(exclude, line)
}

// GitIgnores is the one answer to "does git ignore this board": it asks git
// (check-ignore) for the board's .jaira/ and gives back the rule that does it,
// as "<file>:<line>:<pattern>", or "" when nothing ignores it. Every rule git
// applies counts — jaira's entry, a person's line in info/exclude, a nested
// .gitignore, core.excludesFile, a "!/.jaira/" negation — on any git version.
// A git that cannot answer (exit other than 0 or 1: a symlinked .jaira, no
// repository) is an error, never "not ignored". Ignored, by contrast, only
// says whether jaira's own entry is written.
func GitIgnores(root string) (by string, err error) {
	cmd := exec.Command("git", "-C", root, "check-ignore", "-v", ".jaira/")
	var stderr strings.Builder
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	var exit *exec.ExitError
	switch {
	case err == nil:
		by, _, _ = strings.Cut(strings.TrimSpace(string(out)), "\t")
		return by, nil
	case errors.As(err, &exit) && exit.ExitCode() == 1:
		return "", nil
	}
	if msg := strings.TrimSpace(stderr.String()); msg != "" {
		return "", fmt.Errorf("git check-ignore: %s", msg)
	}
	return "", fmt.Errorf("git check-ignore: %w", err)
}

// AddIgnore keeps the board out of git, in the file IgnoreFileName names.
func AddIgnore(root string) (changed bool, err error) {
	if InGitignore(root) {
		return false, nil
	}
	exclude, line := gitPlace(root)
	if exclude == "" {
		return true, appendEntry(filepath.Join(root, ".gitignore"), privateComment, IgnoreLine)
	}
	err = locked(exclude, func() error {
		if lists(exclude, line) {
			return nil
		}
		changed = true
		return appendEntry(exclude, privateComment, line)
	})
	return changed, err
}

// RemoveIgnore stops excluding the board, wherever it is excluded.
func RemoveIgnore(root string) (changed bool, err error) {
	changed, err = dropEntry(filepath.Join(root, ".gitignore"), privateMarker, func(t string) bool { return boardLine(t, "") })
	if err != nil {
		return changed, err
	}
	exclude, line := gitPlace(root)
	if exclude == "" || !lists(exclude, line) {
		return changed, nil
	}
	// Only jaira's own, anchored entry leaves info/exclude. An unanchored
	// ".jaira/" there is a person's line and excludes every board in the
	// repository, one in a subdirectory too: taking it out to share the top
	// board would publish that one as well. It stays, and Ignored still says so.
	err = locked(exclude, func() error {
		ex, err := dropEntry(exclude, privateMarker, func(t string) bool { return t == line })
		changed = changed || ex
		return err
	})
	return changed, err
}

// StaleEntry is the entry a jaira before 689E0T wrote into info/exclude for
// a board whose path below the top of the working tree starts with a space:
// it cut the space off, so the board in " sub" got "/sub/.jaira/", a line
// that excludes nothing of it. It is "" wherever that line is not provably
// this board's own: the path starts with no space (the old entry was the
// right one), the line is not there, a copy of it lacks jaira's comment
// directly above (a person wrote it), a board stands in the directory the
// line names in any worktree of the repository (it is that board's entry),
// or git cannot list those worktrees or names one that is not there.
func StaleEntry(root string) (exclude, line string) {
	_, exclude, prefix := rawPlace(root)
	trimmed := strings.TrimSpace(prefix)
	if exclude == "" || trimmed == prefix {
		return "", ""
	}
	line = entryLine(trimmed)
	if !ownEntry(exclude, line) {
		return "", ""
	}
	// info/exclude belongs to the git directory every worktree shares, so
	// the line may be the entry of a board in sub of any of them.
	tops, err := worktrees(root)
	if err != nil {
		return "", ""
	}
	// A worktree git lists at a path that is not there — moved by hand, on
	// a drive not mounted — may hold a board in sub nobody can look at.
	for _, top := range tops {
		if _, err := os.Stat(top); err != nil {
			return "", ""
		}
		if _, err := os.Lstat(filepath.Join(top, filepath.FromSlash(trimmed), ticket.DirName)); !errors.Is(err, fs.ErrNotExist) {
			return "", ""
		}
	}
	return exclude, line
}

// worktrees is the top of every working tree of root's repository, the
// main one and each linked one, as 'git worktree list' names them.
func worktrees(root string) ([]string, error) {
	out, err := exec.Command("git", "-C", root, "worktree", "list", "--porcelain").Output()
	if err != nil {
		return nil, err
	}
	var tops []string
	for _, l := range strings.Split(string(out), "\n") {
		if top, ok := strings.CutPrefix(strings.TrimSuffix(l, "\r"), "worktree "); ok {
			tops = append(tops, filepath.FromSlash(top))
		}
	}
	if len(tops) == 0 {
		return nil, errors.New("git worktree list named no working tree")
	}
	return tops, nil
}

// ownEntry reports whether line stands in path, and every time with jaira's
// comment directly above it.
func ownEntry(path, line string) bool {
	b, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	lines := strings.Split(string(b), "\n")
	found := false
	for i, l := range lines {
		if strings.TrimSpace(l) != line {
			continue
		}
		if i == 0 || !strings.Contains(lines[i-1], privateMarker) {
			return false
		}
		found = true
	}
	return found
}

// DropStaleEntry takes StaleEntry's line out of info/exclude, with its
// comment, and gives it back; "" where there is none. The line excludes
// nothing of this board, so taking it out never publishes it.
func DropStaleEntry(root string) (dropped string, err error) {
	exclude, line := StaleEntry(root)
	if line == "" {
		return "", nil
	}
	// Every check again under the lock AddIgnore takes: a 'jaira init' in
	// sub between the first check and the drop finds the line already
	// written, adds none, and must not lose it here.
	err = locked(exclude, func() error {
		if _, again := StaleEntry(root); again != line {
			return nil
		}
		changed, err := dropEntry(exclude, privateMarker, func(t string) bool { return t == line })
		if changed {
			dropped = line
		}
		return err
	})
	return dropped, err
}

// InGitignore reports whether the root .gitignore excludes the board — the
// place an older jaira wrote the entry, and the one MoveIgnore empties.
func InGitignore(root string) bool {
	return lists(filepath.Join(root, ".gitignore"), "")
}

// MoveIgnore carries an entry an older jaira wrote into the tracked .gitignore
// over to info/exclude: the line leaves .gitignore, which the person then
// commits once, and the board stays exactly as private as it was. Nothing
// happens outside a git repository, where .gitignore is the place, and on a
// .gitignore without the entry. The exclude entry is written first: a failure
// between the two leaves the board excluded twice, never not at all.
func MoveIgnore(root string) (moved bool, err error) {
	exclude, line := gitPlace(root)
	if exclude == "" || !InGitignore(root) {
		return false, nil
	}
	err = locked(exclude, func() error {
		if lists(exclude, line) {
			return nil
		}
		return appendEntry(exclude, privateComment, line)
	})
	if err != nil {
		return false, err
	}
	return dropEntry(filepath.Join(root, ".gitignore"), privateMarker, func(t string) bool { return boardLine(t, "") })
}

// locked runs fn holding a lock beside exclude. Every worktree of a repository
// shares that one file, so the store's own locks, which live per working
// tree, would not keep two jaira processes from reading it before either has
// written and both adding the entry.
func locked(exclude string, fn func() error) error {
	if err := os.MkdirAll(filepath.Dir(exclude), 0o755); err != nil {
		return err
	}
	unlock, err := ticket.LockFile(exclude + ".jaira.lock")
	if err != nil {
		return err
	}
	defer unlock()
	return fn()
}

// appendEntry adds a comment and one pattern to the end of path, after a blank
// line, creating the file where it does not exist, in the line endings the
// file already uses. Written whole and renamed into place (ticket.WriteAtomic),
// so git never reads it half written. A file that exists but cannot be read
// is an error, never an empty file: writing "nothing + the entry" over it
// would lose every line a person or git put there.
func appendEntry(path, comment, line string) error {
	existing, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	body := string(existing)
	nl := "\n"
	if strings.Contains(body, "\r\n") {
		nl = "\r\n"
	}
	if body != "" && !strings.HasSuffix(body, "\n") {
		body += nl
	}
	body += nl + comment + nl + line + nl
	return ticket.WriteAtomic(path, []byte(body))
}

// dropEntry takes every line match accepts out of path, with the comment
// carrying marker directly above it and the blank line that follows. A
// comment above a line match does not accept stays: info/exclude can hold
// the entries of two boards, the top one and one in a subdirectory.
func dropEntry(path, marker string, match func(string) bool) (changed bool, err error) {
	existing, err := os.ReadFile(path)
	if err != nil {
		return false, nil
	}
	lines := strings.Split(string(existing), "\n")
	var out []string
	// droppedComment: the line before was a marker comment taken out with
	// its entry. blankFollows: the line before was such an entry, so a blank
	// line here is the one that belonged to it — and only here: a blank line
	// further down belongs to whatever stands above it.
	var droppedComment, blankFollows bool
	for i, line := range lines {
		t := strings.TrimSpace(line)
		wasComment, wasEntry := droppedComment, blankFollows
		droppedComment, blankFollows = false, false
		if match(t) {
			changed = true
			blankFollows = wasComment
			continue
		}
		if strings.Contains(t, marker) && i+1 < len(lines) && match(strings.TrimSpace(lines[i+1])) {
			droppedComment = true
			continue
		}
		if wasEntry && t == "" {
			continue
		}
		out = append(out, line)
	}
	if !changed {
		return false, nil
	}
	body := strings.Join(out, "\n")
	// The blank line before the entry is what ends the file once the entry is
	// gone: "" in an LF file, and in a CRLF file a "\r" that would be left
	// standing alone after the last "\r\n".
	body = strings.TrimSuffix(body, "\r")
	// An entry with no blank line before it took the file's last newline
	// with it: "a\n# jaira…\n/.jaira/\n" must end "a\n", not "a".
	if body != "" && !strings.HasSuffix(body, "\n") && strings.HasSuffix(string(existing), "\n") {
		if strings.HasSuffix(string(existing), "\r\n") {
			body += "\r"
		}
		body += "\n"
	}
	return true, ticket.WriteAtomic(path, []byte(body))
}

// AddLanesIgnore keeps this project's lane directory private, independent of
// whether the rest of the board is shared. It goes into the tracked .gitignore:
// it only matters once the board is shared, and a shared board's .gitignore is
// committed with it. Idempotent: calling it again once the line is already
// present reports changed=false and writes nothing.
func AddLanesIgnore(root string) (changed bool, err error) {
	return addGitignoreLine(root, lanesComment, LanesIgnoreLine, false)
}

// AddMilestonesIgnore puts MilestonesIgnoreLine into the tracked .gitignore,
// where a shared board's teammates get it with the next pull. It is what
// 'jaira share', and 'jaira init' and 'jaira update' on a shared board, call;
// a private board needs no line, its whole .jaira/ is excluded already.
// Idempotent like AddLanesIgnore. dry reports whether it would write and
// writes nothing.
func AddMilestonesIgnore(root string, dry bool) (changed bool, err error) {
	return addGitignoreLine(root, milestonesComment, MilestonesIgnoreLine, dry)
}

// RemoveMilestonesIgnore takes the milestones line out again, for 'share
// --undo': the private board's /.jaira/ covers it.
func RemoveMilestonesIgnore(root string) (changed bool, err error) {
	return dropEntry(filepath.Join(root, ".gitignore"), milestonesMarker, func(t string) bool { return t == MilestonesIgnoreLine })
}

// addGitignoreLine appends line with its comment to the root .gitignore
// unless a line there already says it.
func addGitignoreLine(root, comment, line string, dry bool) (changed bool, err error) {
	path := filepath.Join(root, ".gitignore")
	existing, _ := os.ReadFile(path)
	for _, l := range strings.Split(string(existing), "\n") {
		if strings.TrimSpace(l) == line {
			return false, nil
		}
	}
	if dry {
		return true, nil
	}
	return true, appendEntry(path, comment, line)
}

// RemoveLanesIgnore removes the lanes-only ignore line. Used by 'share
// --undo': once the whole board is private again, /.jaira/ already covers
// .jaira/lanes/, and a leftover line would be a puzzle for the next reader.
func RemoveLanesIgnore(root string) (changed bool, err error) {
	return dropEntry(filepath.Join(root, ".gitignore"), lanesMarker, func(t string) bool { return t == LanesIgnoreLine })
}
