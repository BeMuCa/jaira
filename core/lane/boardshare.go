package lane

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/BeMuCa/jaira/core/ticket"
)

// Boards travel the way lanes do (share.go): a board of a repository is
// saved into the personal catalogue (UserBoardsDir), published into
// .jaira/shared/<you>/boards/<name>/ for teammates, and adopted from there
// into the catalogue. Each of those is a directory shaped like a shipped
// board type, holding every lane's file, so it reads the same wherever it
// lands and needs nothing else to be built from.

// sharedBoardsSubdir is the folder under .jaira/shared/<who>/ that holds
// published boards. One level deeper than a published lane, so Shared's
// shared/*/*.md glob never reads a board's lane as a lane of its own.
const sharedBoardsSubdir = "boards"

// SaveBoard copies board name of root, as it loads, into dst as a
// board-type directory: every lane's file, the order file and the
// description. An existing dst is refused with ErrBoardExists unless
// overwrite, which replaces it whole.
func SaveBoard(root, name, dst string, overwrite bool) ([]string, error) {
	set, err := LoadBoard(root, name)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(set.Lanes))
	raws := make([][]byte, 0, len(set.Lanes))
	for _, l := range set.Lanes {
		raw, err := Bytes(l)
		if err != nil {
			return nil, err
		}
		ids = append(ids, l.ID)
		raws = append(raws, raw)
	}
	if err := clearDst(dst, overwrite); err != nil {
		return nil, err
	}
	return ids, writeBoardDir(dst, ids, raws, Description(root, name))
}

// clearDst refuses an existing dst unless overwrite, and removes it then.
func clearDst(dst string, overwrite bool) error {
	if _, err := os.Stat(dst); err == nil {
		if !overwrite {
			return fmt.Errorf("%w: %s already exists; refusing to overwrite", ErrBoardExists, dst)
		}
		return os.RemoveAll(dst)
	} else if !os.IsNotExist(err) {
		return err
	}
	return nil
}

// SaveToCatalogue saves board name of root into the personal catalogue
// under its own name, when the catalogue has no board of that name yet and
// jaira ships no type of that name. It reports the catalogue path and
// whether it wrote it; a catalogue board already there is never replaced.
func SaveToCatalogue(root, name string) (string, bool, error) {
	norm, err := NormalizeBoard(name)
	if err != nil || norm == "" || isShippedType(norm) || UserBoardsDir() == "" {
		return "", false, err
	}
	dst := filepath.Join(UserBoardsDir(), norm)
	if _, err := os.Stat(dst); err == nil {
		return dst, false, nil
	}
	if _, err := SaveBoard(root, norm, dst, false); err != nil {
		return "", false, err
	}
	return dst, true, nil
}

// PublishBoard copies board name of root into
// <sharedDir>/<who>/boards/<name>/ for teammates to adopt; today's board is
// published as "default". Refused when it is there already, unless
// overwrite.
func PublishBoard(root, name, sharedDir, who string, overwrite bool) (string, error) {
	norm, err := NormalizeBoard(name)
	if err != nil {
		return "", err
	}
	label := norm
	if label == "" {
		label = DefaultBoardName
	}
	dst := filepath.Join(sharedDir, who, sharedBoardsSubdir, label)
	if _, err := SaveBoard(root, norm, dst, overwrite); err != nil {
		return "", err
	}
	return dst, nil
}

// SharedBoard is one board found under .jaira/shared/<folder>/boards/.
type SharedBoard struct {
	BoardType
	Folder string
	Path   string
}

// SharedBoards lists every board published under root's .jaira/shared/. A
// directory that does not read as a board type, or whose lanes do not
// parse, is skipped and named in warnings. Like Shared, nothing here is
// ever loaded onto a board: a shared board is used only once adopted.
func SharedBoards(root string) (boards []SharedBoard, warnings []string, err error) {
	if root == "" {
		return nil, nil, nil
	}
	matches, err := filepath.Glob(filepath.Join(root, ticket.DirName, ticket.SharedSubdir, "*", sharedBoardsSubdir, "*"))
	if err != nil {
		return nil, nil, err
	}
	sort.Strings(matches)
	for _, p := range matches {
		if fi, err := os.Stat(p); err != nil || !fi.IsDir() {
			continue
		}
		bt, _, err := readBoardDir(p)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("shared board %s was skipped: %v", p, err))
			continue
		}
		boards = append(boards, SharedBoard{BoardType: bt, Folder: filepath.Base(filepath.Dir(filepath.Dir(p))), Path: p})
	}
	return boards, warnings, nil
}

// readBoardDir reads the board-type directory dir, with every lane it names
// resolved and parsed.
func readBoardDir(dir string) (BoardType, [][]byte, error) {
	fsys := os.DirFS(dir)
	bt, err := boardType(fsys, ".", filepath.Base(dir))
	if err != nil {
		return BoardType{}, nil, err
	}
	raws, err := dirLanes(fsys, ".", bt)
	return bt, raws, err
}

// AdoptBoard copies the published board in srcDir — a path 'jaira boards
// shared' lists — into the personal catalogue under its directory's name,
// after checking it reads as a board type. It refuses a name jaira ships a
// type of, and an existing catalogue board unless overwrite.
func AdoptBoard(srcDir string, overwrite bool) (BoardType, string, error) {
	bt, raws, err := readBoardDir(srcDir)
	if err != nil {
		return BoardType{}, "", err
	}
	if !boardNameRe.MatchString(bt.Name) {
		return BoardType{}, "", fmt.Errorf("board name %q: use lower-case letters, digits and dashes", bt.Name)
	}
	if isShippedType(bt.Name) {
		return BoardType{}, "", fmt.Errorf("%w: jaira ships a board type called %s", ErrBoardExists, bt.Name)
	}
	if UserBoardsDir() == "" {
		return BoardType{}, "", fmt.Errorf("no home directory for the board catalogue")
	}
	dst := filepath.Join(UserBoardsDir(), bt.Name)
	if err := clearDst(dst, overwrite); err != nil {
		return BoardType{}, "", err
	}
	bt.Catalogue = true
	return bt, dst, writeBoardDir(dst, bt.Lanes, raws, bt.Description)
}
