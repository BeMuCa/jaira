// Package wintrap finds the five ways this repository has broken on Windows
// before, by reading the source on whatever machine you are sitting at.
//
// Every one of the seven Windows-only CI failures in this repository's history
// was one of five patterns, and every one of them is visible in the source
// without running Windows:
//
//  1. t.Setenv("HOME") without USERPROFILE beside it — os.UserHomeDir reads
//     USERPROFILE on Windows, so the test walks straight past its temp home.
//  2. A //go:embed target with no eol=lf line in .gitattributes — a Windows
//     checkout with core.autocrlf hands go:embed CRLF bytes.
//  3. A permission claim with no runtime.GOOS branch — Windows carries no
//     POSIX permission bits and os.Chmod only flips the read-only attribute.
//  4. An expected binary name with no .exe handling.
//  5. A path glued together with a literal "/" instead of filepath.Join.
//
// A site that looks like one of these but is provably not — a URL, a git
// refname, an io/fs path, which is always slash-separated on every platform —
// is exempted with a //wintrap:ok comment on the line or the line above it,
// carrying the reason. The rule still fires; the exemption is visible in the
// diff and in review, which a loosened rule would not be.
package wintrap

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Finding is one site that will break on Windows. Problem says what was found,
// Fix says what to do about it — a finding that only names the problem leaves
// the reader exactly where the red CI job left them.
type Finding struct {
	Rule    int
	File    string // relative to the scanned root, slash-separated
	Line    int
	Problem string
	Fix     string
}

func (f Finding) String() string {
	return fmt.Sprintf("%s:%d: windows trap %d: %s\n\tfix: %s", f.File, f.Line, f.Rule, f.Problem, f.Fix)
}

// skipDirs are not ours to police: .git is not source, and testdata holds this
// package's own deliberate violations.
var skipDirs = map[string]bool{".git": true, "testdata": true, "vendor": true, "node_modules": true}

// Scan walks root and reports every site matching one of the five patterns.
func Scan(root string) ([]Finding, error) {
	attrs, err := loadAttributes(root)
	if err != nil {
		return nil, err
	}
	var out []Finding
	fset := token.NewFileSet()
	err = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if p != root && skipDirs[d.Name()] {
				return fs.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".go") {
			return nil
		}
		file, perr := parser.ParseFile(fset, p, nil, parser.ParseComments)
		if perr != nil {
			// Unparseable Go is the compiler's complaint, not ours.
			return nil
		}
		rel := relSlash(root, p)
		exempt := exemptLines(fset, file)
		var found []Finding
		found = append(found, checkSetenv(fset, file, rel)...)
		found = append(found, checkEmbed(fset, file, rel, root, filepath.Dir(p), attrs)...)
		found = append(found, checkPerm(fset, file, rel)...)
		found = append(found, checkExe(fset, file, rel)...)
		found = append(found, checkSlash(fset, file, rel)...)
		for _, f := range found {
			if !exempt[f.Line] {
				out = append(out, f)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].File != out[j].File {
			return out[i].File < out[j].File
		}
		if out[i].Line != out[j].Line {
			return out[i].Line < out[j].Line
		}
		return out[i].Rule < out[j].Rule
	})
	return out, nil
}

func relSlash(root, p string) string {
	r, err := filepath.Rel(root, p)
	if err != nil {
		return filepath.ToSlash(p)
	}
	return filepath.ToSlash(r)
}

// exemptLines collects the lines silenced by a //wintrap:ok comment, which
// applies to its own line and to the line after it.
func exemptLines(fset *token.FileSet, f *ast.File) map[int]bool {
	out := map[int]bool{}
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			if !strings.Contains(c.Text, "wintrap:ok") {
				continue
			}
			l := fset.Position(c.Pos()).Line
			out[l] = true
			out[l+1] = true
		}
	}
	return out
}

// funcBodies yields one body per top-level function, which is the unit every
// "in the same function" rule below is measured against.
func funcBodies(f *ast.File) []ast.Node {
	var out []ast.Node
	for _, d := range f.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok && fd.Body != nil {
			out = append(out, fd.Body)
		}
	}
	return out
}

func line(fset *token.FileSet, n ast.Node) int { return fset.Position(n.Pos()).Line }

func selName(e ast.Expr) (pkg, sel string, ok bool) {
	s, ok := e.(*ast.SelectorExpr)
	if !ok {
		return "", "", false
	}
	if id, ok := s.X.(*ast.Ident); ok {
		return id.Name, s.Sel.Name, true
	}
	return "", s.Sel.Name, true
}

func strLit(e ast.Expr) (string, bool) {
	b, ok := e.(*ast.BasicLit)
	if !ok || b.Kind != token.STRING {
		return "", false
	}
	v, err := strconv.Unquote(b.Value)
	if err != nil {
		return "", false
	}
	return v, true
}

// anyNode reports whether pred holds for n or for any node below it, which is
// the shape every "does this function mention X" question below takes.
func anyNode(n ast.Node, pred func(ast.Node) bool) bool {
	found := false
	ast.Inspect(n, func(x ast.Node) bool {
		if !found && pred(x) {
			found = true
		}
		return !found
	})
	return found
}

// litContains reports whether n holds a string literal containing sub, compared
// case-insensitively.
func litContains(n ast.Node, sub string) bool {
	return anyNode(n, func(x ast.Node) bool {
		lit, ok := x.(*ast.BasicLit)
		if !ok {
			return false
		}
		v, ok := strLit(lit)
		return ok && strings.Contains(strings.ToLower(v), sub)
	})
}

func hasGOOS(body ast.Node) bool {
	return anyNode(body, func(n ast.Node) bool {
		s, ok := n.(*ast.SelectorExpr)
		if !ok || s.Sel.Name != "GOOS" {
			return false
		}
		id, ok := s.X.(*ast.Ident)
		return ok && id.Name == "runtime"
	})
}

// ---- rule 1: HOME without USERPROFILE -------------------------------------

func checkSetenv(fset *token.FileSet, f *ast.File, rel string) []Finding {
	var out []Finding
	for _, body := range funcBodies(f) {
		set := map[string]bool{}
		var homeLines []int
		ast.Inspect(body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) == 0 {
				return true
			}
			_, sel, ok := selName(call.Fun)
			if !ok || sel != "Setenv" {
				return true
			}
			key, ok := strLit(call.Args[0])
			if !ok {
				return true
			}
			set[key] = true
			if key == "HOME" {
				homeLines = append(homeLines, line(fset, call))
			}
			return true
		})
		if len(homeLines) == 0 || set["USERPROFILE"] {
			continue
		}
		for _, l := range homeLines {
			out = append(out, Finding{
				Rule: 1, File: rel, Line: l,
				Problem: `Setenv("HOME", ...) with no USERPROFILE beside it: os.UserHomeDir reads USERPROFILE on Windows, so this redirection has no effect there and the test reads the real home directory`,
				Fix:     `add t.Setenv("USERPROFILE", <the same directory>) next to it`,
			})
		}
	}
	return out
}

// ---- rule 2: go:embed with nothing pinning it in .gitattributes -----------

func checkEmbed(fset *token.FileSet, f *ast.File, rel, root, dir string, attrs *attributes) []Finding {
	var out []Finding
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			text := strings.TrimPrefix(c.Text, "//")
			if !strings.HasPrefix(text, "go:embed ") {
				continue
			}
			var uncovered []string
			for _, pat := range embedPatterns(strings.TrimPrefix(text, "go:embed ")) {
				for _, file := range embedFiles(dir, pat) {
					r := relSlash(root, file)
					if attrs.pinned(r) {
						continue
					}
					uncovered = append(uncovered, r)
				}
			}
			if len(uncovered) == 0 {
				continue
			}
			sort.Strings(uncovered)
			out = append(out, Finding{
				Rule: 2, File: rel, Line: fset.Position(c.Pos()).Line,
				Problem: fmt.Sprintf("//go:embed pulls in %d file(s) that no line in .gitattributes pins, starting at %s: a Windows checkout with core.autocrlf embeds them with CRLF, and the embedded bytes stop matching what the code expects", len(uncovered), uncovered[0]),
				Fix:     fmt.Sprintf("add a line to .gitattributes pinning them, e.g. %q — or %q if they are not text", coverPattern(uncovered[0])+" text eol=lf", coverPattern(uncovered[0])+" binary"),
			})
		}
	}
	return out
}

// embedPatterns splits a go:embed argument list, which is space-separated with
// optionally quoted entries and an optional all: prefix per pattern.
func embedPatterns(s string) []string {
	var out []string
	for _, fld := range strings.Fields(s) {
		if unq, err := strconv.Unquote(fld); err == nil {
			fld = unq
		}
		fld = strings.TrimPrefix(fld, "all:")
		if fld != "" {
			out = append(out, fld)
		}
	}
	return out
}

// embedFiles resolves one go:embed pattern to the regular files it pulls in.
func embedFiles(dir, pattern string) []string {
	matches, err := filepath.Glob(filepath.Join(dir, filepath.FromSlash(pattern)))
	if err != nil {
		return nil
	}
	var out []string
	for _, m := range matches {
		_ = filepath.WalkDir(m, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if !d.IsDir() {
				out = append(out, p)
			}
			return nil
		})
	}
	return out
}

// coverPattern suggests the .gitattributes line for an uncovered file: the
// whole extension inside its directory, which is what the existing lines do.
func coverPattern(rel string) string {
	ext := path.Ext(rel)
	if ext == "" {
		return rel
	}
	return path.Dir(rel) + "/*" + ext
}

// ---- rule 3: permission claims without a GOOS branch ----------------------

func checkPerm(fset *token.FileSet, f *ast.File, rel string) []Finding {
	var out []Finding
	for _, body := range funcBodies(f) {
		if hasGOOS(body) {
			continue
		}
		var lines []int
		ast.Inspect(body, func(n ast.Node) bool {
			switch e := n.(type) {
			case *ast.CallExpr:
				if _, sel, ok := selName(e.Fun); ok {
					if sel == "Perm" && isModeCall(e.Fun) {
						lines = append(lines, line(fset, e))
					}
					// Chmod to 0 exists only to make an operation fail; on
					// Windows it does not, so the test it guards never fires.
					if sel == "Chmod" && len(e.Args) > 0 && isZeroMode(e.Args[len(e.Args)-1]) {
						lines = append(lines, line(fset, e))
					}
				}
			case *ast.BinaryExpr:
				// Only an octal mask is a permission claim. Masking a Mode()
				// against a named constant asks what kind of file it is —
				// os.ModeCharDevice and friends mean the same on every platform.
				if e.Op == token.AND && (isModeExpr(e.X) && isOctalLit(e.Y) || isModeExpr(e.Y) && isOctalLit(e.X)) {
					lines = append(lines, line(fset, e))
				}
			}
			return true
		})
		for _, l := range lines {
			out = append(out, Finding{
				Rule: 3, File: rel, Line: l,
				Problem: "a claim about POSIX permission bits with no runtime.GOOS branch in this function: Windows carries no execute bit and os.Chmod there only flips the read-only attribute, so this claim is false on Windows however correct the code is",
				Fix:     `guard it with if runtime.GOOS != "windows" { ... }, or skip the test there, and say in a comment why`,
			})
		}
	}
	return out
}

func isModeCall(fun ast.Expr) bool {
	s, ok := fun.(*ast.SelectorExpr)
	return ok && isModeExpr(s.X)
}

func isModeExpr(e ast.Expr) bool {
	call, ok := e.(*ast.CallExpr)
	if !ok {
		return false
	}
	_, sel, ok := selName(call.Fun)
	return ok && sel == "Mode"
}

func isOctalLit(e ast.Expr) bool {
	b, ok := e.(*ast.BasicLit)
	return ok && b.Kind == token.INT
}

func isZeroMode(e ast.Expr) bool {
	b, ok := e.(*ast.BasicLit)
	if !ok || b.Kind != token.INT {
		return false
	}
	v, err := strconv.ParseInt(strings.NewReplacer("0o", "", "0O", "", "_", "").Replace(b.Value), 8, 64)
	return err == nil && v == 0
}

// ---- rule 4: an expected binary name without .exe -------------------------

func checkExe(fset *token.FileSet, f *ast.File, rel string) []Finding {
	var out []Finding
	for _, body := range funcBodies(f) {
		// A function that names ".exe" itself already accounts for the suffix. A
		// site that delegates it to a helper carries //wintrap:ok with its reason
		// instead — a silencer nobody can see at the site it silences is worse
		// than the finding.
		if hasGOOS(body) || litContains(body, ".exe") {
			continue
		}
		var lines []int
		ast.Inspect(body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if isGoBuildOutput(call) {
				lines = append(lines, line(fset, call))
			}
			return true
		})
		for _, l := range lines {
			out = append(out, Finding{
				Rule: 4, File: rel, Line: l,
				Problem: `a binary is built with "go build -o" and nothing in this function accounts for the .exe suffix: on Windows the file lands next to the name you asked for, and every later reference to it misses`,
				Fix:     `append ".exe" to the output name when runtime.GOOS == "windows"`,
			})
		}
	}
	return out
}

// isGoBuildOutput reports whether call is exec.Command("go", "build", ..., "-o", ...) —
// the only shape the rule 4 message describes. "install" takes no -o at all and
// "test -o" writes a test binary, which is not what the finding says. Anything
// looser than this and the check fires on every unrelated tool that happens to
// take a "-o" flag, which is what a finding is not allowed to do: claim
// something it did not see.
func isGoBuildOutput(call *ast.CallExpr) bool {
	pkg, sel, ok := selName(call.Fun)
	if !ok || pkg != "exec" || (sel != "Command" && sel != "CommandContext") {
		return false
	}
	var lits []string
	for _, a := range call.Args {
		if v, ok := strLit(a); ok {
			lits = append(lits, v)
		}
	}
	if len(lits) == 0 || lits[0] != "go" {
		return false
	}
	build, out := false, false
	for _, v := range lits[1:] {
		switch v {
		case "build":
			build = true
		case "-o":
			out = true
		}
	}
	return build && out
}

// ---- rule 5: a path glued together with a literal "/" ---------------------

var fsPackages = map[string]bool{"os": true, "filepath": true, "ioutil": true}

var prefixFuncs = map[string]bool{
	"TrimPrefix": true, "HasPrefix": true, "CutPrefix": true,
	"TrimSuffix": true, "HasSuffix": true, "CutSuffix": true,
}

func checkSlash(fset *token.FileSet, f *ast.File, rel string) []Finding {
	var out []Finding
	report := func(n ast.Node, where string) {
		out = append(out, Finding{
			Rule: 5, File: rel, Line: line(fset, n),
			Problem: "a path is glued together with a literal \"/\" and handed to " + where + ": Windows separates with \\, so the result never matches the paths the filesystem hands back",
			Fix:     `build it with filepath.Join, and compare against filepath.ToSlash(...) if you need slashes — or mark the site //wintrap:ok if it is a URL, a git refname or an io/fs path`,
		})
	}
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		pkg, sel, ok := selName(call.Fun)
		if !ok {
			return true
		}
		switch {
		case fsPackages[pkg]:
			for _, a := range call.Args {
				if sepConcat(a) {
					report(call, pkg+"."+sel)
					break
				}
			}
		case pkg == "strings" && prefixFuncs[sel] && len(call.Args) == 2:
			if sepConcat(call.Args[1]) {
				report(call, "strings."+sel)
			}
		}
		return true
	})
	return out
}

// sepConcat reports whether e concatenates a literal path separator into a
// value. It does not try to spot URLs: this only runs on values already headed
// into os.*, filepath.* or a prefix test against a filesystem path, where a URL
// has no business being. A site that is genuinely one carries //wintrap:ok.
func sepConcat(e ast.Expr) bool {
	bin, ok := e.(*ast.BinaryExpr)
	if !ok || bin.Op != token.ADD {
		return false
	}
	return litContains(e, "/")
}
