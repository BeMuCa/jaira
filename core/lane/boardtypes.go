package lane

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// A board type is what 'jaira boards create --type' builds a board from: a
// directory under boardtypes/ shaped like a board's own lanes directory.
//
//	description  one line saying what the type is for; 'jaira boards types' prints it
//	order        the board's lane ids, one per line, in column order
//	<id>.md      a lane of the type's own, optional
//
// An id in order with no file of the type's own is the built-in lane of that
// id, so a type that only rearranges the shipped lanes ships no lane file,
// and one with lanes of its own — an ideas board, say — ships just those.
//
//go:embed boardtypes
var boardTypesEmbed embed.FS

// boardTypesFS is the shipped types, one directory per type under
// boardtypes/. A variable so a test can lay other types in.
var boardTypesFS fs.FS = boardTypesEmbed

const boardTypesDir = "boardtypes"

// DefaultBoardType is the type a board created without one gets: the
// built-in lanes for code.
const DefaultBoardType = "code"

// UserBoardsDir is the personal board catalogue, ~/.jaira/boards (or under
// $JAIRA_HOME): one directory per board, shaped like a shipped type. Every
// repository reads it, so a board saved or adopted there once is a type
// 'jaira boards create --type' offers everywhere, while the boards a
// repository already has stay its own copies and are never rewritten from it.
func UserBoardsDir() string {
	return inUserHome("boards")
}

// ErrNoSuchBoardType is returned for a type name that is neither shipped nor
// in the catalogue.
var ErrNoSuchBoardType = errors.New("no such board type")

// ErrBoardExists is returned when the board to create already has lanes.
var ErrBoardExists = errors.New("board already exists")

// BoardType is one type, as 'jaira boards types' lists it: shipped, or a
// board in the personal catalogue (Catalogue).
type BoardType struct {
	Name        string
	Description string
	Lanes       []string
	Catalogue   bool
}

// BoardTypes lists the shipped types and the catalogue's boards, sorted by
// name. A catalogue board named like a shipped type is left out: the shipped
// one wins, so a new binary's version of it is never hidden behind an old
// copy. A catalogue directory that does not read as a type is skipped rather
// than failing the list.
func BoardTypes() ([]BoardType, error) {
	entries, err := fs.ReadDir(boardTypesFS, boardTypesDir)
	if err != nil {
		return nil, err
	}
	var out []BoardType
	shipped := map[string]bool{}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		bt, err := boardType(boardTypesFS, path.Join(boardTypesDir, e.Name()), e.Name())
		if err != nil {
			return nil, err
		}
		shipped[bt.Name] = true
		out = append(out, bt)
	}
	if dir := UserBoardsDir(); dir != "" {
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			if !e.IsDir() || !boardNameRe.MatchString(e.Name()) || shipped[e.Name()] {
				continue
			}
			bt, err := boardType(os.DirFS(dir), e.Name(), e.Name())
			if err != nil {
				continue
			}
			bt.Catalogue = true
			out = append(out, bt)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// boardType reads the type in directory dir of fsys: its description and
// its order.
func boardType(fsys fs.FS, dir, name string) (BoardType, error) {
	desc, err := fs.ReadFile(fsys, path.Join(dir, descriptionFileName))
	if err != nil {
		return BoardType{}, fmt.Errorf("board type %s: %w", name, err)
	}
	order, err := fs.ReadFile(fsys, path.Join(dir, orderFileName))
	if err != nil {
		return BoardType{}, fmt.Errorf("board type %s: %w", name, err)
	}
	return BoardType{
		Name:        name,
		Description: strings.Join(strings.Fields(string(desc)), " "),
		Lanes:       strings.Fields(string(order)),
	}, nil
}

// isShippedType reports whether jaira ships a type called name.
func isShippedType(name string) bool {
	_, err := fs.Stat(boardTypesFS, path.Join(boardTypesDir, name, orderFileName))
	return err == nil
}

// typeLanes resolves the lane ids of type typ, in its order, with each
// lane's file content: the type's own file for an id it has one for, the
// built-in lane's otherwise. A shipped type is looked up first, then the
// catalogue.
func typeLanes(typ string) ([]string, [][]byte, error) {
	fsys, dir := fs.FS(boardTypesFS), path.Join(boardTypesDir, typ)
	if !isShippedType(typ) && boardNameRe.MatchString(typ) && UserBoardsDir() != "" {
		fsys, dir = os.DirFS(UserBoardsDir()), typ
	}
	bt, err := boardType(fsys, dir, typ)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil, fmt.Errorf("%w %q; 'jaira boards types' lists them", ErrNoSuchBoardType, typ)
		}
		return nil, nil, err
	}
	raws, err := dirLanes(fsys, dir, bt)
	return bt.Lanes, raws, err
}

// dirLanes is the file content of each lane of bt, whose files are in dir
// of fsys: its own file for an id it has one for, the built-in lane's
// otherwise. Every own file must parse.
func dirLanes(fsys fs.FS, dir string, bt BoardType) ([][]byte, error) {
	own := map[string][]byte{}
	files, _ := fs.Glob(fsys, path.Join(dir, "*.md"))
	for _, f := range files {
		raw, err := fs.ReadFile(fsys, f)
		if err != nil {
			return nil, err
		}
		l, err := Parse(raw, "boardtype:"+f)
		if err != nil {
			return nil, err
		}
		own[l.ID] = raw
	}
	builtins, err := Builtins()
	if err != nil {
		return nil, err
	}
	builtinByID := indexByID(builtins)

	var raws [][]byte
	for _, id := range bt.Lanes {
		raw, ok := own[id]
		if !ok {
			b, found := builtinByID[id]
			if !found {
				return nil, fmt.Errorf("board type %s names lane %q, which it has no file for and jaira does not ship", bt.Name, id)
			}
			if raw, err = Bytes(b); err != nil {
				return nil, err
			}
		}
		raws = append(raws, raw)
	}
	return raws, nil
}

// CreateBoard writes board name of type typ into
// <root>/.jaira/boards/<name>/lanes/: the type's lane files, its order file,
// and description as the board's description. It refuses today's board and
// a board that already has lane files, with ErrBoardExists, and writes
// nothing then. It returns the lanes directory it wrote.
func CreateBoard(root, name, typ, description string) (string, []string, error) {
	dir, err := newBoardDir(root, name)
	if err != nil {
		return "", nil, err
	}
	ids, raws, err := typeLanes(typ)
	if err != nil {
		return "", nil, err
	}
	return dir, ids, writeBoardDir(dir, ids, raws, description)
}

// CreateBoardWithLanes is CreateBoard with the lanes picked one by one
// instead of a type: ids, in that order, each one a lane of set (today's
// board) or one Installable offers — the same lanes the lane settings
// screen shows. An id that is neither is refused before anything is written.
func CreateBoardWithLanes(root, name string, set *Set, ids []string, description string) (string, error) {
	dir, err := newBoardDir(root, name)
	if err != nil {
		return "", err
	}
	if len(ids) == 0 {
		return "", errors.New("a board needs at least one lane")
	}
	pickable, err := PickableLanes(set)
	if err != nil {
		return "", err
	}
	byID := indexByID(pickable)
	seen := map[string]bool{}
	raws := make([][]byte, 0, len(ids))
	for _, id := range ids {
		l, ok := byID[id]
		if !ok {
			var known []string
			for _, p := range pickable {
				known = append(known, p.ID)
			}
			return "", fmt.Errorf("no lane %q to put on the board; lanes there are: %s", id, strings.Join(known, ", "))
		}
		if seen[id] {
			return "", fmt.Errorf("lane %q is named twice", id)
		}
		seen[id] = true
		raw, err := Bytes(l)
		if err != nil {
			return "", err
		}
		raws = append(raws, raw)
	}
	return dir, writeBoardDir(dir, ids, raws, description)
}

// PickableLanes is every lane a new board can be built from: set's own,
// then each built-in and catalogue lane set does not have (Installable).
func PickableLanes(set *Set) ([]*Lane, error) {
	more, err := Installable(set)
	if err != nil {
		return nil, err
	}
	return append(append([]*Lane{}, set.Lanes...), more...), nil
}

// newBoardDir checks that board name can be created in root and returns its
// lanes directory.
func newBoardDir(root, name string) (string, error) {
	norm, err := NormalizeBoard(name)
	if err != nil {
		return "", err
	}
	if norm == "" {
		return "", fmt.Errorf("%w: %q is today's board, in %s", ErrBoardExists, DefaultBoardName, ProjectLanesDir(root))
	}
	dir := boardLanesDir(root, norm)
	if hasLaneFiles(dir) {
		return "", fmt.Errorf("%w: %s already has lanes in %s", ErrBoardExists, norm, dir)
	}
	return dir, nil
}

// writeBoardDir writes a board-type-shaped directory: one file per lane,
// the order file and the description.
func writeBoardDir(dir string, ids []string, raws [][]byte, description string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for i, id := range ids {
		if err := os.WriteFile(filepath.Join(dir, id+".md"), raws[i], 0o644); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filepath.Join(dir, orderFileName), []byte(strings.Join(ids, "\n")+"\n"), 0o644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, descriptionFileName), []byte(strings.TrimSpace(description)+"\n"), 0o644)
}
