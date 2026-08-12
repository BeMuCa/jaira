// Package ticket implements the on-disk ticket format and store.
//
// The file format IS the API: tickets are hand-edited and read in git diffs, so
// writing one field must leave every other byte of the file untouched. That
// constraint rules out the obvious implementation (unmarshal into a struct,
// mutate, marshal back), and it also rules out re-printing a YAML AST — both
// silently renormalize quoting, key order, and comment alignment on fields the
// caller never asked to change.
//
// Instead, writes use the parser purely as a locator: it yields the line and
// column of the value being replaced, and the replacement is spliced into the
// original bytes. Everything outside the spliced span is preserved by
// construction rather than by the serializer's good behaviour.
package ticket

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

const delim = "---"

var (
	// ErrNoFrontmatter means the file does not open with a YAML frontmatter block.
	ErrNoFrontmatter = errors.New("ticket: file does not begin with '---'")
	// ErrUnterminated means the frontmatter block is never closed.
	ErrUnterminated = errors.New("ticket: frontmatter block is not terminated by '---'")
	// ErrUnsafeYAML means the frontmatter uses constructs this tool cannot
	// rewrite without risking corruption. Reported rather than mangled.
	ErrUnsafeYAML = errors.New("ticket: frontmatter uses YAML features jaira cannot safely rewrite")
)

// Doc is a ticket file split into the three zones that get different treatment:
// the opening delimiter, the frontmatter body (the only part YAML ever sees),
// and the markdown body (never parsed, only carried).
type Doc struct {
	open string // "---\n", plus any leading BOM/whitespace exactly as found
	fm   string // frontmatter body, always newline-terminated
	tail string // closing delimiter line and everything after it
}

// ParseDoc splits a ticket file into its zones without interpreting the YAML.
func ParseDoc(src []byte) (*Doc, error) {
	s := string(src)

	// Tolerate a UTF-8 BOM, which some editors add, without letting it break
	// the delimiter check or reappear in the wrong place on write.
	var bom string
	if strings.HasPrefix(s, "\ufeff") {
		bom, s = "\ufeff", strings.TrimPrefix(s, "\ufeff")
	}

	first, rest, ok := cutLine(s)
	if !ok || strings.TrimRight(first, " \t\r") != delim {
		return nil, ErrNoFrontmatter
	}
	// Preserve the opening line's exact ending. cutLine strips a trailing CR so
	// the delimiter compares equal on CRLF files, but dropping it on write would
	// be a byte-fidelity violation on the one line the whole format depends on.
	eol := "\n"
	if idx := strings.IndexByte(s, '\n'); idx > 0 && s[idx-1] == '\r' {
		eol = "\r\n"
	}

	// Find the closing delimiter: the first line that is exactly "---".
	offset := 0
	for {
		line, after, ok := cutLine(rest[offset:])
		if !ok && line == "" {
			return nil, ErrUnterminated
		}
		if strings.TrimRight(line, " \t\r") == delim {
			return &Doc{
				open: bom + first + eol,
				fm:   rest[:offset],
				tail: rest[offset:],
			}, nil
		}
		if !ok {
			return nil, ErrUnterminated
		}
		offset += len(rest[offset:]) - len(after)
	}
}

// cutLine splits off the first line, returning it without its newline. ok is
// false when there was no newline (i.e. the final, unterminated line).
func cutLine(s string) (line, rest string, ok bool) {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return strings.TrimSuffix(s[:i], "\r"), s[i+1:], true
	}
	return s, "", false
}

// String reassembles the file exactly.
func (d *Doc) String() string { return d.open + d.fm + d.tail }

// Bytes reassembles the file exactly.
func (d *Doc) Bytes() []byte { return []byte(d.String()) }

// Frontmatter returns the raw frontmatter zone, for callers that want to
// unmarshal it normally. Reads do not need byte fidelity, so unmarshalling here
// is safe; only writes must go through the splice path.
func (d *Doc) Frontmatter() string { return d.fm }

// Body returns the markdown body with the closing delimiter removed.
func (d *Doc) Body() string {
	_, rest, _ := cutLine(d.tail)
	return rest
}

// SetBody replaces the markdown body, leaving frontmatter untouched.
func (d *Doc) SetBody(body string) {
	line, _, _ := cutLine(d.tail)
	if body != "" && !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	d.tail = line + "\n" + body
}

// Validate rejects frontmatter containing constructs the splice writer cannot
// safely round-trip. Anchors and aliases are the realistic hazard: a hand-edited
// ticket using them could be silently restructured by a naive rewrite, so the
// tool refuses to touch such a file rather than risking corruption (STORE-10).
func (d *Doc) Validate() error {
	f, err := parser.ParseBytes([]byte(d.fm), parser.ParseComments)
	if err != nil {
		return fmt.Errorf("ticket: frontmatter is not valid YAML: %w", err)
	}
	if len(f.Docs) == 0 || f.Docs[0].Body == nil {
		return nil // empty frontmatter is structurally fine
	}
	if _, ok := f.Docs[0].Body.(*ast.MappingNode); !ok {
		// A single key/value parses as MappingValueNode rather than MappingNode.
		if _, ok := f.Docs[0].Body.(*ast.MappingValueNode); !ok {
			return fmt.Errorf("%w: root is %T, expected a mapping", ErrUnsafeYAML, f.Docs[0].Body)
		}
	}
	var bad string
	ast.Walk(visitorFunc(func(n ast.Node) ast.Visitor {
		switch n.(type) {
		case *ast.AnchorNode:
			bad = "anchor (&name)"
		case *ast.AliasNode:
			bad = "alias (*name)"
		case *ast.MergeKeyNode:
			bad = "merge key (<<)"
		}
		return nil
	}), f.Docs[0].Body)
	if bad != "" {
		return fmt.Errorf("%w: found a YAML %s; edit the file by hand or remove it", ErrUnsafeYAML, bad)
	}
	return nil
}

type visitorFunc func(ast.Node) ast.Visitor

func (v visitorFunc) Visit(n ast.Node) ast.Visitor {
	if r := v(n); r != nil {
		return r
	}
	return v
}

// topLevel returns the frontmatter's top-level key/value pairs in source order.
func (d *Doc) topLevel() ([]*ast.MappingValueNode, error) {
	f, err := parser.ParseBytes([]byte(d.fm), parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("ticket: frontmatter is not valid YAML: %w", err)
	}
	if len(f.Docs) == 0 || f.Docs[0].Body == nil {
		return nil, nil
	}
	switch b := f.Docs[0].Body.(type) {
	case *ast.MappingNode:
		return b.Values, nil
	case *ast.MappingValueNode:
		return []*ast.MappingValueNode{b}, nil
	default:
		return nil, fmt.Errorf("%w: root is %T", ErrUnsafeYAML, b)
	}
}

func (d *Doc) find(key string) (*ast.MappingValueNode, error) {
	vals, err := d.topLevel()
	if err != nil {
		return nil, err
	}
	for _, v := range vals {
		if v.Key.GetToken().Value == key {
			return v, nil
		}
	}
	return nil, nil
}

// Has reports whether a top-level key is present.
func (d *Doc) Has(key string) bool {
	v, err := d.find(key)
	return err == nil && v != nil
}

// Keys lists the top-level keys in source order.
func (d *Doc) Keys() []string {
	vals, err := d.topLevel()
	if err != nil {
		return nil
	}
	out := make([]string, 0, len(vals))
	for _, v := range vals {
		out = append(out, v.Key.GetToken().Value)
	}
	return out
}

// SetScalar sets a top-level scalar field. An existing value is spliced in
// place, preserving the line's indentation, quoting style and trailing comment.
// A missing key is appended at the end of the frontmatter.
func (d *Doc) SetScalar(key, val string) error {
	node, err := d.find(key)
	if err != nil {
		return err
	}
	if node == nil {
		d.appendLine(fmt.Sprintf("%s: %s", key, encodeScalar(val)))
		return nil
	}

	switch node.Value.(type) {
	case *ast.SequenceNode, *ast.MappingNode, *ast.MappingValueNode:
		return fmt.Errorf("ticket: %q holds a collection, not a scalar", key)
	}

	tok := node.Value.GetToken()
	if tok == nil || tok.Position == nil {
		return fmt.Errorf("ticket: no source position for %q", key)
	}

	lines := strings.Split(d.fm, "\n")
	li := tok.Position.Line - 1
	if li < 0 || li >= len(lines) {
		return fmt.Errorf("ticket: %q reports line %d, outside the frontmatter", key, tok.Position.Line)
	}
	line := lines[li]
	col := tok.Position.Column - 1
	if col < 0 || col > len(line) {
		return fmt.Errorf("ticket: %q reports column %d, outside line %d", key, tok.Position.Column, tok.Position.Line)
	}

	oldLen, err := scalarSourceLen(line[col:], tok.Value)
	if err != nil {
		return fmt.Errorf("ticket: %q: %w", key, err)
	}

	// Multi-line scalars (block literals, folded) span lines the locator does
	// not describe; refuse rather than truncate the remainder.
	if strings.ContainsAny(line[col:col+oldLen], "|>") && oldLen <= 2 {
		return fmt.Errorf("ticket: %q uses a multi-line block scalar; not supported by SetScalar", key)
	}

	lines[li] = line[:col] + encodeScalarLike(line[col:col+oldLen], val) + line[col+oldLen:]
	d.fm = strings.Join(lines, "\n")
	return nil
}

// scalarSourceLen measures how many bytes of src the scalar occupies, so the
// splice replaces exactly the value and not any trailing comment.
func scalarSourceLen(src, value string) (int, error) {
	if src == "" {
		return 0, errors.New("value position is at end of line")
	}
	switch src[0] {
	case '"', '\'':
		q := src[0]
		for i := 1; i < len(src); i++ {
			if src[i] == '\\' && q == '"' {
				i++
				continue
			}
			if src[i] == q {
				// Doubled single-quote is an escaped quote, not the end.
				if q == '\'' && i+1 < len(src) && src[i+1] == '\'' {
					i++
					continue
				}
				return i + 1, nil
			}
		}
		return 0, errors.New("unterminated quoted scalar")
	}
	if strings.HasPrefix(src, value) && value != "" {
		return len(value), nil
	}
	// Plain scalar: runs to a comment or end of line. A '#' only starts a
	// comment when preceded by whitespace.
	end := len(src)
	for i := 0; i < len(src); i++ {
		if src[i] == '#' && i > 0 && (src[i-1] == ' ' || src[i-1] == '\t') {
			end = i
			break
		}
	}
	return len(strings.TrimRight(src[:end], " \t")), nil
}

// encodeScalarLike keeps the original quoting style where possible, so a field
// written as "quoted" stays quoted and a bare field stays bare.
func encodeScalarLike(old, val string) string {
	if len(old) >= 2 {
		switch old[0] {
		case '"':
			return `"` + escapeDouble(val) + `"`
		case '\'':
			return `'` + strings.ReplaceAll(val, "'", "''") + `'`
		}
	}
	return encodeScalar(val)
}

// encodeScalar quotes only when a bare scalar would be ambiguous or invalid.
func encodeScalar(v string) string {
	if v == "" {
		return `""`
	}
	if needsQuoting(v) {
		return `"` + escapeDouble(v) + `"`
	}
	return v
}

func needsQuoting(v string) bool {
	if strings.TrimSpace(v) != v {
		return true
	}
	if strings.ContainsAny(v, ":#\n\"'{}[]&*!|>%@`,") {
		return true
	}
	switch strings.ToLower(v) {
	case "true", "false", "null", "~", "yes", "no", "on", "off":
		return true
	}
	// Anything that would parse as a number must be quoted to stay a string.
	if isNumericLike(v) {
		return true
	}
	return false
}

func isNumericLike(v string) bool {
	seenDigit := false
	for i, r := range v {
		switch {
		case r >= '0' && r <= '9':
			seenDigit = true
		case r == '-' || r == '+':
			if i != 0 {
				return false
			}
		case r == '.' || r == 'e' || r == 'E':
			// allowed inside a number
		default:
			return false
		}
	}
	return seenDigit
}

func escapeDouble(v string) string {
	var b strings.Builder
	for _, r := range v {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\t':
			b.WriteString(`\t`)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// SetRaw sets a field to a literal YAML scalar without quoting it. Used for
// values that are unambiguous by construction — booleans and RFC3339 timestamps
// — so tickets read idiomatically instead of carrying `ready: "false"`. Callers
// are responsible for passing something valid; anything user-supplied should go
// through SetScalar, which quotes defensively.
func (d *Doc) SetRaw(key, literal string) error {
	node, err := d.find(key)
	if err != nil {
		return err
	}
	if node == nil {
		d.appendLine(key + ": " + literal)
		return nil
	}
	switch node.Value.(type) {
	case *ast.SequenceNode, *ast.MappingNode, *ast.MappingValueNode:
		return fmt.Errorf("ticket: %q holds a collection, not a scalar", key)
	}
	tok := node.Value.GetToken()
	if tok == nil || tok.Position == nil {
		return fmt.Errorf("ticket: no source position for %q", key)
	}
	lines := strings.Split(d.fm, "\n")
	li := tok.Position.Line - 1
	if li < 0 || li >= len(lines) {
		return fmt.Errorf("ticket: %q reports line %d, outside the frontmatter", key, tok.Position.Line)
	}
	line := lines[li]
	col := tok.Position.Column - 1
	if col < 0 || col > len(line) {
		return fmt.Errorf("ticket: %q reports column %d, outside line %d", key, tok.Position.Column, tok.Position.Line)
	}
	oldLen, err := scalarSourceLen(line[col:], tok.Value)
	if err != nil {
		return fmt.Errorf("ticket: %q: %w", key, err)
	}
	lines[li] = line[:col] + literal + line[col+oldLen:]
	d.fm = strings.Join(lines, "\n")
	return nil
}

// SetList sets a top-level list field. The whole block for that key is replaced
// (it necessarily spans multiple lines), but no other key is disturbed. A
// trailing comment on the key line is carried over.
func (d *Doc) SetList(key string, items []string) error {
	node, err := d.find(key)
	if err != nil {
		return err
	}
	if node == nil {
		d.appendBlock(renderList(key, items, ""))
		return nil
	}

	keyTok := node.Key.GetToken()
	if keyTok == nil || keyTok.Position == nil {
		return fmt.Errorf("ticket: no source position for %q", key)
	}
	start := keyTok.Position.Line - 1
	end := lastLine(node.Value) - 1
	if end < start {
		end = start
	}

	lines := strings.Split(d.fm, "\n")
	if start < 0 || end >= len(lines) {
		return fmt.Errorf("ticket: %q spans lines %d-%d, outside the frontmatter", key, start+1, end+1)
	}

	comment := trailingComment(lines[start])
	replacement := strings.Split(renderList(key, items, comment), "\n")

	out := append([]string{}, lines[:start]...)
	out = append(out, replacement...)
	out = append(out, lines[end+1:]...)
	d.fm = strings.Join(out, "\n")
	return nil
}

// List reads a top-level list field. Absent or empty yields nil.
func (d *Doc) List(key string) ([]string, error) {
	node, err := d.find(key)
	if err != nil || node == nil {
		return nil, err
	}
	seq, ok := node.Value.(*ast.SequenceNode)
	if !ok {
		return nil, fmt.Errorf("ticket: %q is %T, not a list", key, node.Value)
	}
	out := make([]string, 0, len(seq.Values))
	for _, v := range seq.Values {
		out = append(out, v.GetToken().Value)
	}
	return out, nil
}

// Scalar reads a top-level scalar field.
func (d *Doc) Scalar(key string) (string, bool, error) {
	node, err := d.find(key)
	if err != nil || node == nil {
		return "", false, err
	}
	switch node.Value.(type) {
	case *ast.SequenceNode, *ast.MappingNode, *ast.MappingValueNode:
		return "", false, fmt.Errorf("ticket: %q holds a collection, not a scalar", key)
	case *ast.NullNode:
		return "", true, nil
	}
	return node.Value.GetToken().Value, true, nil
}

func renderList(key string, items []string, comment string) string {
	suffix := ""
	if comment != "" {
		suffix = " " + comment
	}
	if len(items) == 0 {
		return key + ": []" + suffix
	}
	var b strings.Builder
	b.WriteString(key + ":" + suffix)
	for _, it := range items {
		b.WriteString("\n  - " + encodeScalar(it))
	}
	return b.String()
}

// trailingComment extracts a '#' comment from a line, if any.
func trailingComment(line string) string {
	inS, inD := false, false
	for i := 0; i < len(line); i++ {
		switch line[i] {
		case '\'':
			if !inD {
				inS = !inS
			}
		case '"':
			if !inS {
				inD = !inD
			}
		case '#':
			if !inS && !inD && i > 0 && (line[i-1] == ' ' || line[i-1] == '\t') {
				return strings.TrimRight(line[i:], " \t")
			}
		}
	}
	return ""
}

// lastLine finds the final source line a node occupies.
func lastLine(n ast.Node) int {
	max := 0
	if t := n.GetToken(); t != nil && t.Position != nil {
		max = t.Position.Line
	}
	if seq, ok := n.(*ast.SequenceNode); ok {
		for _, v := range seq.Values {
			if l := lastLine(v); l > max {
				max = l
			}
		}
	}
	return max
}

// appendLine adds a new key at the end of the frontmatter zone.
func (d *Doc) appendLine(line string) { d.appendBlock(line) }

func (d *Doc) appendBlock(block string) {
	if d.fm != "" && !strings.HasSuffix(d.fm, "\n") {
		d.fm += "\n"
	}
	// Keep a trailing blank line at the end of frontmatter as the last thing,
	// so appended keys land above it rather than after it.
	if strings.HasSuffix(d.fm, "\n\n") {
		d.fm = strings.TrimSuffix(d.fm, "\n") + block + "\n\n"
		return
	}
	d.fm += block + "\n"
}
