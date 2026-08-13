package lane

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BeMuCa/jaira/core/ticket"
)

// DefaultBoard is a per-user setting: which lanes a freshly initialised board
// gets, and which ticket Options start ticked. An absent file, an absent
// lanes: field, or an empty list all mean the same thing — the built-ins —
// which is easy to get backwards, so it is stated here rather than only in
// the file format.
type DefaultBoard struct {
	Lanes    []string
	Options  []string
	Path     string
	Warnings []string
	// Body is the prose below the frontmatter, preserved across a save so
	// SaveDefaultBoard never erases a note the user left themselves.
	Body string
}

// DefaultBoardPath is the per-user file that decides a new board's lanes and
// which ticket Options start ticked. $JAIRA_DEFAULT_BOARD overrides it — an
// explicit override rather than a sibling of UserLanesDir, because a test
// needs to point at a temp file and deriving a sibling path from an
// already-overridable directory breaks the day someone points
// JAIRA_LANES_DIR at something that is not <something>/lanes.
func DefaultBoardPath() string {
	if v := os.Getenv("JAIRA_DEFAULT_BOARD"); v != "" {
		return v
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".jaira", "default-board.md")
}

// LoadDefaultBoard reads DefaultBoardPath(). No file at all is the normal
// state, not a missing-configuration error, and returns a zero board. A file
// that exists but does not parse is warned about and treated the same as
// absent, rather than taking 'jaira init' or 'jaira lanes' down over one bad
// file — a teammate's shared default board naming their own lane, or a typo,
// must not break your board.
func LoadDefaultBoard() (*DefaultBoard, error) {
	path := DefaultBoardPath()
	b := &DefaultBoard{Path: path}
	raw, err := os.ReadFile(path)
	if err != nil {
		return b, nil // absent or unreadable: treated as no default board
	}
	d, err := ticket.ParseDoc(raw)
	if err != nil {
		b.Warnings = append(b.Warnings, fmt.Sprintf(
			"default board %s did not parse and was treated as absent: %v", path, err))
		return b, nil
	}
	if lanes, err := d.List("lanes"); err == nil {
		b.Lanes = lanes
	}
	if opts, err := d.List("options"); err == nil {
		b.Options = opts
	}
	b.Body = strings.TrimSpace(d.Body())
	return b, nil
}

// defaultBoardBody is the prose a freshly saved default board gets when it
// has none of its own yet — a file with an empty body reads as broken rather
// than as "nothing to say here".
const defaultBoardBody = `# Default board

Decides which lanes a freshly initialised board gets, and which ticket
Options start ticked. An absent file, an absent lanes: field, or an empty
list all mean the built-ins. Edit this file directly, or from jaira's home
screen with 'd'.`

// SaveDefaultBoard writes b to b.Path with lanes: and options: in a fixed
// order, so a saved file diffs predictably. The prose body is preserved
// verbatim if the file already had one.
func SaveDefaultBoard(b *DefaultBoard) error {
	if err := os.MkdirAll(filepath.Dir(b.Path), 0o755); err != nil {
		return err
	}
	body := strings.TrimSpace(b.Body)
	if body == "" {
		body = defaultBoardBody
	}
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString("lanes: " + yamlFlowList(b.Lanes) + "\n")
	sb.WriteString("options: " + yamlFlowList(b.Options) + "\n")
	sb.WriteString("---\n\n")
	sb.WriteString(body)
	sb.WriteString("\n")
	return os.WriteFile(b.Path, []byte(sb.String()), 0o644)
}

// yamlFlowList renders a flow-style YAML list, [] for an empty one. Lane ids
// and option names are both simple, already-constrained identifiers
// ([a-z0-9-] for ids; option names come from the same well-behaved
// frontmatter fields), so no quoting is needed.
func yamlFlowList(items []string) string {
	if len(items) == 0 {
		return "[]"
	}
	return "[" + strings.Join(items, ", ") + "]"
}

// Differs reports whether board's selection differs from the built-in set:
// a different set of ids, or a selected lane that resolves — through set, the
// catalogue-aware load already in hand — to something other than the
// built-in of that id. Compared as a set, not a list: reordering alone must
// not trigger materialisation, since order is precedence's business, not the
// default board's.
func Differs(board *DefaultBoard, set *Set) bool {
	if len(board.Lanes) == 0 {
		return false // absent or empty selection means the built-ins
	}
	builtins, err := Builtins()
	if err != nil {
		return true // cannot prove sameness; safer to materialise than to hide
	}
	if len(board.Lanes) != len(builtins) {
		return true
	}
	builtinIDs := make(map[string]bool, len(builtins))
	for _, b := range builtins {
		builtinIDs[b.ID] = true
	}
	for _, id := range board.Lanes {
		if !builtinIDs[id] {
			return true
		}
		l, ok := set.Get(id)
		if !ok || !l.Builtin {
			return true // missing, or resolves to an override rather than the built-in
		}
	}
	return false
}

// Materialise writes <root>/.jaira/lanes/ only when board's selection differs
// from the built-in set — the point of the criterion is that a repo whose
// owner changed nothing carries no lane files at all. set is used to resolve
// each selected id to its lane (built-in or catalogue override) and to read
// its bytes; an id the set does not have is warned about, on b, rather than
// failing the whole write.
func Materialise(root string, set *Set, b *DefaultBoard) ([]string, error) {
	if !Differs(b, set) {
		return nil, nil
	}
	dir := ProjectLanesDir(root)
	var written []string
	for _, id := range b.Lanes {
		l, ok := set.Get(id)
		if !ok {
			b.Warnings = append(b.Warnings, fmt.Sprintf(
				"default board names lane %q, which is not installed; skipped", id))
			continue
		}
		raw, err := Bytes(l)
		if err != nil {
			return written, err
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return written, err
		}
		// Named from the resolved lane's own, already-validID-constrained ID —
		// never from the string in the default board — so a default board
		// naming something odd cannot escape dir (T-5-07).
		dst := filepath.Join(dir, l.ID+".md")
		if err := os.WriteFile(dst, raw, 0o644); err != nil {
			return written, err
		}
		written = append(written, dst)
	}
	return written, nil
}

// Validate checks a default board against the currently installed lanes,
// returning warnings in the same shape core/lane's own loader already uses:
// an unknown lane id, an unknown option, or (carried over from
// LoadDefaultBoard) an unparseable file. There is one warning channel,
// already surfaced by 'jaira lanes', --json, the TUI warnings block and
// every command through loadEnv — a second reporting path for the same
// class of problem would be the mistake here.
//
// board.Warnings is included as-is rather than re-derived: LoadDefaultBoard
// is the only place that knows whether the file failed to parse.
func Validate(board *DefaultBoard, set *Set) []string {
	warnings := append([]string{}, board.Warnings...)
	for _, id := range board.Lanes {
		if _, ok := set.Get(id); !ok {
			warnings = append(warnings, fmt.Sprintf(
				"default board names lane %q, which is not installed", id))
		}
	}
	known := make(map[string]bool, len(set.Options()))
	for _, o := range set.Options() {
		known[o] = true
	}
	for _, o := range board.Options {
		if !known[o] {
			warnings = append(warnings, fmt.Sprintf(
				"default board names option %q, which no installed lane requires", o))
		}
	}
	return warnings
}

// ResolveOptions turns a default board's option choices into the ticked/
// unticked list a new ticket's Options checklist needs, against the options
// this set's installed lanes actually declare. An option the board names
// that no lane requires is silently dropped here — 'jaira lanes' reporting
// that typo is Validate's job, not this function's.
func ResolveOptions(set *Set, board *DefaultBoard) []ticket.BodyOption {
	ticked := make(map[string]bool, len(board.Options))
	for _, o := range board.Options {
		ticked[o] = true
	}
	names := set.Options()
	out := make([]ticket.BodyOption, 0, len(names))
	for _, n := range names {
		out = append(out, ticket.BodyOption{Name: n, Ticked: ticked[n]})
	}
	return out
}
