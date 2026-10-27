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

// ErrNoSuchBoardType is returned for a type name jaira does not ship.
var ErrNoSuchBoardType = errors.New("no such board type")

// ErrBoardExists is returned when the board to create already has lanes.
var ErrBoardExists = errors.New("board already exists")

// BoardType is one shipped type, as 'jaira boards types' lists it.
type BoardType struct {
	Name        string
	Description string
	Lanes       []string
}

// BoardTypes lists the shipped types, sorted by name.
func BoardTypes() ([]BoardType, error) {
	entries, err := fs.ReadDir(boardTypesFS, boardTypesDir)
	if err != nil {
		return nil, err
	}
	var out []BoardType
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		bt, err := boardType(e.Name())
		if err != nil {
			return nil, err
		}
		out = append(out, bt)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func boardType(name string) (BoardType, error) {
	desc, err := fs.ReadFile(boardTypesFS, path.Join(boardTypesDir, name, descriptionFileName))
	if err != nil {
		return BoardType{}, fmt.Errorf("board type %s: %w", name, err)
	}
	order, err := fs.ReadFile(boardTypesFS, path.Join(boardTypesDir, name, orderFileName))
	if err != nil {
		return BoardType{}, fmt.Errorf("board type %s: %w", name, err)
	}
	return BoardType{
		Name:        name,
		Description: strings.Join(strings.Fields(string(desc)), " "),
		Lanes:       strings.Fields(string(order)),
	}, nil
}

// typeLanes resolves the lane ids of type typ, in its order, with each
// lane's file content: the type's own file for an id it has one for, the
// built-in lane's otherwise.
func typeLanes(typ string) ([]string, [][]byte, error) {
	bt, err := boardType(typ)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil, fmt.Errorf("%w %q; 'jaira boards types' lists them", ErrNoSuchBoardType, typ)
		}
		return nil, nil, err
	}
	own := map[string][]byte{}
	files, _ := fs.Glob(boardTypesFS, path.Join(boardTypesDir, typ, "*.md"))
	for _, f := range files {
		raw, err := fs.ReadFile(boardTypesFS, f)
		if err != nil {
			return nil, nil, err
		}
		l, err := Parse(raw, "boardtype:"+f)
		if err != nil {
			return nil, nil, err
		}
		own[l.ID] = raw
	}
	builtins, err := Builtins()
	if err != nil {
		return nil, nil, err
	}
	builtinByID := indexByID(builtins)

	var raws [][]byte
	for _, id := range bt.Lanes {
		raw, ok := own[id]
		if !ok {
			b, found := builtinByID[id]
			if !found {
				return nil, nil, fmt.Errorf("board type %s names lane %q, which it has no file for and jaira does not ship", typ, id)
			}
			if raw, err = Bytes(b); err != nil {
				return nil, nil, err
			}
		}
		raws = append(raws, raw)
	}
	return bt.Lanes, raws, nil
}

// CreateBoard writes board name of type typ into
// <root>/.jaira/boards/<name>/lanes/: the type's lane files, its order file,
// and description as the board's description. It refuses today's board and
// a board that already has lane files, with ErrBoardExists, and writes
// nothing then. It returns the lanes directory it wrote.
func CreateBoard(root, name, typ, description string) (string, []string, error) {
	norm, err := NormalizeBoard(name)
	if err != nil {
		return "", nil, err
	}
	if norm == "" {
		return "", nil, fmt.Errorf("%w: %q is today's board, in %s", ErrBoardExists, DefaultBoardName, ProjectLanesDir(root))
	}
	dir := boardLanesDir(root, norm)
	if hasLaneFiles(dir) {
		return "", nil, fmt.Errorf("%w: %s already has lanes in %s", ErrBoardExists, norm, dir)
	}
	ids, raws, err := typeLanes(typ)
	if err != nil {
		return "", nil, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", nil, err
	}
	for i, id := range ids {
		if err := os.WriteFile(filepath.Join(dir, id+".md"), raws[i], 0o644); err != nil {
			return "", nil, err
		}
	}
	if err := os.WriteFile(filepath.Join(dir, orderFileName), []byte(strings.Join(ids, "\n")+"\n"), 0o644); err != nil {
		return "", nil, err
	}
	if err := os.WriteFile(filepath.Join(dir, descriptionFileName), []byte(strings.TrimSpace(description)+"\n"), 0o644); err != nil {
		return "", nil, err
	}
	return dir, ids, nil
}
