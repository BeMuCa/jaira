// Package link answers one question: what else is connected to this ticket,
// and where does it live now.
//
// It exists as its own package because four callers need the same answer —
// the gate, validate, the CLI and the TUI — and three copies of it would be
// three truths. The gate in particular must not learn to read the filesystem
// (see core/gate's doc comment), so it takes the one predicate it needs as an
// injected function rather than importing this package's index wholesale.
//
// The connections are of five kinds and every one of them is read in both
// directions:
//
//	blocked-by / blocks   a hard dependency, the only kind a gate enforces
//	parent / children     containment, stored only on the child
//	related               a loose association, stored on whichever side was edited
//	follows / followed-by a follow-up produced by a review
//
// Children are never stored. A children list on the parent could not be
// written at all in the common case — the parent may be in the logbook, in
// the archive, or on a ref with no file here — and two sessions adding two
// children would collide on one line. The child names its parent, one writer
// per file, and the tree is rebuilt by reading.
package link

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// Place is where a ticket was found. A link whose other end has left the
// board is still a link, and saying only "not found" is what made every
// connection rot the moment the work it pointed at finished.
type Place string

const (
	// PlaceBoard means the ticket has a file under .jaira/tickets/.
	PlaceBoard Place = "board"
	// PlaceRef means the board can see it but nobody here has pulled it in.
	PlaceRef Place = "ref"
	// PlaceLogbook means it is finished and filed under .jaira/logbook/.
	PlaceLogbook Place = "logbook"
	// PlaceArchive means it was taken off the board into .jaira/archive/.
	PlaceArchive Place = "archive"
	// PlaceUnknown means the id resolves to nothing here at all.
	PlaceUnknown Place = "unknown"
)

// Kind names one direction of one relation. Both directions are separate
// kinds: "blocks" and "blocked by" read as opposite facts to whoever is
// looking, and collapsing them would hide which way the arrow points.
type Kind string

const (
	KindBlockedBy  Kind = "blocked-by"
	KindBlocks     Kind = "blocks"
	KindParent     Kind = "parent"
	KindChild      Kind = "child"
	KindRelated    Kind = "related"
	KindFollows    Kind = "follows"
	KindFollowedBy Kind = "followed-by"
)

// Title is how a kind is written out for a reader.
func (k Kind) Title() string {
	switch k {
	case KindBlockedBy:
		return "waiting on"
	case KindBlocks:
		return "blocking"
	case KindParent:
		return "part of"
	case KindChild:
		return "contains"
	case KindRelated:
		return "related to"
	case KindFollows:
		return "follows"
	case KindFollowedBy:
		return "followed by"
	}
	return string(k)
}

// Ref is one end of a link: enough to show it and to find it again, without
// the caller having to know which of the four places it came out of.
type Ref struct {
	ID    string
	Title string
	Lane  string
	Place Place
	// Path is the file, empty for a ticket that only exists on a ref.
	Path string
	// Done means the ticket sits in a terminal lane. A blocker that is done
	// no longer blocks, wherever it is filed.
	Done bool
	// Depth is how far below the subject a child sits: 0 for a direct child,
	// 1 for its child, and so on. Zero for every other kind.
	Depth int
}

// Entry pairs a ref with the relation that produced it.
type Entry struct {
	Kind Kind
	Ref  Ref
}

// Index holds every ticket this repository can see, whichever of the four
// places it lives in.
//
// Building it reads the board eagerly (it is already in memory for every
// caller) and the filed-away directories lazily: a logbook grows without
// bound, and paying for it on every board render would break the promise
// that the board opens instantly. Nothing here reads a filed file until a
// question is asked that only the file can answer.
type Index struct {
	lanes *lane.Set

	// board is every ticket the store listed — including the ones that exist
	// only on a ref, which the store already folds in.
	board map[string]*ticket.Ticket
	// filed maps id to file for the logbook and the archive. It comes from a
	// directory listing, so having it costs nothing.
	filed map[string]string

	// loaded caches filed tickets that have actually been read.
	loaded map[string]*ticket.Ticket
	// allFiled records that every filed ticket has been read, which is what
	// reverse links need and nothing else does.
	allFiled bool
}

// Build indexes what the store can see. The tickets are passed in rather than
// listed here because every caller has already listed them and a second walk
// would double the cost of opening the board.
func Build(s *ticket.Store, lanes *lane.Set, board []*ticket.Ticket) *Index {
	ix := &Index{
		lanes:  lanes,
		board:  make(map[string]*ticket.Ticket, len(board)),
		filed:  map[string]string{},
		loaded: map[string]*ticket.Ticket{},
	}
	for _, t := range board {
		ix.board[t.ID] = t
	}
	if s != nil {
		for id, path := range s.FiledAwayIDs() {
			// A ticket that is both on the board and filed away is reported
			// as a problem elsewhere; here the board copy is the live one.
			if _, onBoard := ix.board[id]; !onBoard {
				ix.filed[id] = path
			}
		}
	}
	return ix
}

// terminal reports whether a lane id is a terminal one. An unknown lane is
// not terminal: this installation cannot vouch for a lane it does not have.
func (ix *Index) terminal(laneID string) bool {
	if ix == nil || ix.lanes == nil {
		return false
	}
	l, ok := ix.lanes.Get(laneID)
	return ok && l.Terminal
}

// placeOf reads the place out of a filed ticket's path.
func placeOf(path string) Place {
	if strings.Contains(filepath.ToSlash(path), "/"+ticket.ArchiveSubdir+"/") {
		return PlaceArchive
	}
	return PlaceLogbook
}

// loadFiled reads one filed ticket, caching it.
func (ix *Index) loadFiled(id string) *ticket.Ticket {
	if t, ok := ix.loaded[id]; ok {
		return t
	}
	path, ok := ix.filed[id]
	if !ok {
		return nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		ix.loaded[id] = nil
		return nil
	}
	d, err := ticket.ParseDoc(raw)
	if err != nil {
		ix.loaded[id] = nil
		return nil
	}
	t, err := ticket.Decode(d, path)
	if err != nil {
		ix.loaded[id] = nil
		return nil
	}
	ix.loaded[id] = t
	return t
}

// Lookup resolves one id to everything a reader needs to see it, wherever it
// lives. The second result is false when the id is not known here at all —
// which, unlike before, now genuinely means "nowhere", not merely "not on the
// board".
func (ix *Index) Lookup(id string) (Ref, bool) {
	if ix == nil || id == "" {
		return Ref{ID: id, Place: PlaceUnknown}, false
	}
	if t, ok := ix.board[id]; ok {
		place := PlaceBoard
		if t.ReadOnly {
			place = PlaceRef
		}
		return Ref{
			ID:    t.ID,
			Title: t.Title,
			Lane:  t.Status,
			Place: place,
			Path:  t.Path,
			Done:  ix.terminal(t.Status),
		}, true
	}
	path, filed := ix.filed[id]
	if !filed {
		return Ref{ID: id, Place: PlaceUnknown}, false
	}
	ref := Ref{ID: id, Place: placeOf(path), Path: path}
	// The logbook only ever receives finished work, so a ticket found there
	// is done even before its file is read. The archive is a drawer, not a
	// verdict, so that one has to be opened.
	if ref.Place == PlaceLogbook {
		ref.Done = true
	}
	if t := ix.loadFiled(id); t != nil {
		ref.Title = t.Title
		ref.Lane = t.Status
		if ref.Place == PlaceArchive {
			ref.Done = ix.terminal(t.Status)
		}
	}
	return ref, true
}

// Satisfied reports whether a ticket may be treated as a cleared dependency:
// it exists somewhere and it is finished.
//
// This is the whole of what the gate needs, which is why it is a single
// boolean: core/gate stays free of the filesystem and takes this as an
// injected function.
func (ix *Index) Satisfied(id string) bool {
	ref, ok := ix.Lookup(id)
	return ok && ref.Done
}

// Known reports whether the id resolves to a ticket anywhere.
func (ix *Index) Known(id string) bool {
	_, ok := ix.Lookup(id)
	return ok
}

// readAll makes sure every filed ticket has been read. Reverse links — who
// names me as their parent, who blocks on me — cannot be answered any other
// way, and they are only ever asked when somebody opens the link view.
func (ix *Index) readAll() []*ticket.Ticket {
	out := make([]*ticket.Ticket, 0, len(ix.board)+len(ix.filed))
	for _, t := range ix.board {
		out = append(out, t)
	}
	if !ix.allFiled {
		for id := range ix.filed {
			ix.loadFiled(id)
		}
		ix.allFiled = true
	}
	for _, t := range ix.loaded {
		if t != nil {
			out = append(out, t)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Relations lists every connection of one ticket, in both directions,
// grouped kind by kind in a fixed order so the view does not reshuffle
// between openings.
func (ix *Index) Relations(id string) []Entry {
	if ix == nil || id == "" {
		return nil
	}
	all := ix.readAll()
	subject := ix.ticket(id)

	var out []Entry
	add := func(k Kind, otherID string, depth int) {
		if otherID == "" || otherID == id {
			return
		}
		ref, ok := ix.Lookup(otherID)
		if !ok {
			// An id nothing here knows is still worth showing: a dangling
			// link is a fact about this ticket, and hiding it is how a typo
			// survives for months.
			ref = Ref{ID: otherID, Place: PlaceUnknown}
		}
		ref.Depth = depth
		out = append(out, Entry{Kind: k, Ref: ref})
	}

	if subject != nil {
		for _, dep := range subject.BlockedBy {
			add(KindBlockedBy, dep, 0)
		}
		add(KindParent, subject.Parent, 0)
	}
	for _, other := range all {
		if other.ID == id {
			continue
		}
		for _, dep := range other.BlockedBy {
			if dep == id {
				add(KindBlocks, other.ID, 0)
				break
			}
		}
	}
	// Children come as a tree, deepest branch included, so an epic reads as
	// one thing rather than as a list that has to be re-opened per level.
	for _, c := range ix.childTree(id, all, map[string]bool{id: true}, 0) {
		out = append(out, c)
	}
	for _, otherID := range ix.relatedBoth(id, subject, all) {
		add(KindRelated, otherID, 0)
	}
	if subject != nil {
		add(KindFollows, subject.Follows, 0)
	}
	for _, other := range all {
		if other.Follows == id {
			add(KindFollowedBy, other.ID, 0)
		}
	}
	return out
}

// childTree walks the children of id depth-first. seen carries every id on
// the path from the root, so a parent ring — which a hand-edited file can
// always produce — terminates instead of hanging the caller.
func (ix *Index) childTree(id string, all []*ticket.Ticket, seen map[string]bool, depth int) []Entry {
	var out []Entry
	for _, t := range all {
		if t.Parent != id || seen[t.ID] {
			continue
		}
		seen[t.ID] = true
		ref, ok := ix.Lookup(t.ID)
		if !ok {
			ref = Ref{ID: t.ID, Place: PlaceUnknown}
		}
		ref.Depth = depth
		out = append(out, Entry{Kind: KindChild, Ref: ref})
		out = append(out, ix.childTree(t.ID, all, seen, depth+1)...)
	}
	return out
}

// relatedBoth collects the ids related to this ticket from both sides: the
// ones it names, and the ones that name it. The relation is symmetric in
// meaning but stored on one side only, because a link that had to be written
// twice would be half-written most of the time.
func (ix *Index) relatedBoth(id string, subject *ticket.Ticket, all []*ticket.Ticket) []string {
	seen := map[string]bool{id: true}
	var out []string
	take := func(other string) {
		if seen[other] {
			return
		}
		seen[other] = true
		out = append(out, other)
	}
	if subject != nil {
		for _, r := range subject.Related {
			take(r)
		}
	}
	for _, t := range all {
		for _, r := range t.Related {
			if r == id {
				take(t.ID)
			}
		}
	}
	sort.Strings(out)
	return out
}

// ticket returns the subject itself from whichever place holds it.
func (ix *Index) ticket(id string) *ticket.Ticket {
	if t, ok := ix.board[id]; ok {
		return t
	}
	return ix.loadFiled(id)
}

// ParentCycle reports the ring a ticket's parent chain runs into, if it runs
// into one, as the ids on that ring starting from this ticket. Nil means the
// chain terminates.
//
// Worth its own function because the chain is followed recursively wherever
// it is rendered: a ring that nobody checked for is a hung board, not a bad
// diagram.
func (ix *Index) ParentCycle(id string) []string {
	if ix == nil {
		return nil
	}
	seen := map[string]bool{}
	var chain []string
	for cur := id; cur != ""; {
		if seen[cur] {
			return append(chain, cur)
		}
		seen[cur] = true
		chain = append(chain, cur)
		t := ix.ticket(cur)
		if t == nil {
			return nil
		}
		cur = t.Parent
	}
	return nil
}

// Whereabouts says where a linked ticket lives now, in the words a reader
// needs. It is here rather than in each view because the CLI and the TUI
// must not disagree about it: "logbook · done" is the difference between a
// dependency that cleared and one that is still owed, and two spellings of
// that would be two answers.
func (r Ref) Whereabouts() string {
	switch r.Place {
	case PlaceBoard:
		return "on the board · " + r.Lane
	case PlaceRef:
		return "on its ref · jaira pull " + ticket.Handle(r.ID)
	case PlaceLogbook:
		return "logbook · done"
	case PlaceArchive:
		if r.Done {
			return "archive · done"
		}
		return "archive · unfinished"
	}
	return "nowhere here · dangling link"
}

// Label is the ticket as one line: handle and title, with a stand-in when
// the title is only readable where the ticket actually lives.
func (r Ref) Label() string {
	title := r.Title
	if title == "" {
		title = "(title unknown here)"
	}
	return ticket.Handle(r.ID) + "  " + title
}
