package gitrepo

import (
	"regexp"
	"strings"

	"github.com/BeMuCa/jaira/core/ticket"
)

// Hunt is the newest finished bug hunt: the commit it ran at, and what it
// says about the code now. A chat handed four notes cannot work that out: the
// hunt's note has scrolled out of them by the time a ticket waits for a
// person, and a merge from master moves HEAD past the hunt's commit without
// changing anything it read. State is one of:
//
//   - ok: nothing the ticket committed since touched the files it read
//   - fail: a commit of the ticket since then did, or this checkout does not
//     contain the hunt's commit
//   - unknown: no such commit, but a file it read is changed and not
//     committed — git cannot tell whether the hunt saw that change; also
//     when git cannot answer at all. Never drawn red: no clock is read
//   - skipped: the ticket sets stop-after-fix, so a fix round runs its tests
//     and no hunt — a hunt that no longer covers the code, or none at all,
//     is the choice the ticket made, not a fault
type Hunt struct {
	At    string `json:"at"`
	State string `json:"state"`
}

// The states a Hunt reports.
const (
	HuntOK      = "ok"
	HuntFail    = "fail"
	HuntUnsure  = "unknown"
	HuntSkipped = "skipped"
)

// huntDone is the closing line a bug hunt writes: "bughunt: done at <commit>
// ...; files: <path> <path>". The file list is the hunt's own account of what
// it read; without one the hunt is taken to have read all code.
var huntDone = regexp.MustCompile(`^bughunt: done at ([0-9a-f]{7,40})\b`)

// LastHunt reads the newest hunt off the notes, nil when none finished and
// the ticket does not set stop-after-fix.
func (r *Repo) LastHunt(t *ticket.Ticket) *Hunt {
	h := r.newestHunt(t)
	if t.StopAfterFix && (h == nil || h.State == HuntFail) {
		if h == nil {
			h = &Hunt{}
		}
		h.State = HuntSkipped
	}
	return h
}

func (r *Repo) newestHunt(t *ticket.Ticket) *Hunt {
	notes := ticket.Notes(t.Body)
	for i := len(notes) - 1; i >= 0; i-- {
		text := strings.TrimSpace(notes[i].Text)
		m := huntDone.FindStringSubmatch(text)
		if m == nil {
			continue
		}
		var files []string
		if at := strings.LastIndex(text, "files: "); at >= 0 {
			files = strings.Fields(text[at+len("files: "):])
		}
		state, err := r.HuntCoverage(m[1], t.ID, files)
		switch {
		case err != nil || state == HuntUnknown:
			return &Hunt{At: m[1], State: HuntUnsure}
		case state == HuntStale:
			return &Hunt{At: m[1], State: HuntFail}
		default:
			return &Hunt{At: m[1], State: HuntOK}
		}
	}
	return nil
}

// TouchesCode reports whether any of these commits changed a path that is
// not board data (isBoardData): a ticket whose commits carry only its own
// file, tags or milestones changed no code, and has nothing to hunt.
func (r *Repo) TouchesCode(shas []string) (bool, error) {
	if len(shas) == 0 {
		return false, nil
	}
	out, err := r.run(append([]string{"-c", "core.quotePath=false", "show", "--no-renames", "--format=", "--name-only", "-z"}, shas...)...)
	if err != nil {
		return false, err
	}
	return anyCode(out), nil
}

// HuntReport has the fields of gate.HuntReport, so a caller converts one to
// the other. gitrepo cannot import gate: core/board's tests import gitrepo,
// and gate reaches core/board through lane.
type HuntReport struct {
	CodeChanged bool
	State       string
	At          string
	Head        string
}

// HuntCheck is what the gate asks before a ticket reaches the lane a person
// releases (gate.Env.Hunt): whether its commits changed code, its newest
// hunt and HEAD. Git failing to list the commits reads as no code changed —
// the gate refuses on what git shows, never on what it could not read.
func (r *Repo) HuntCheck(t *ticket.Ticket) HuntReport {
	shas, err := r.CommitsForTicket(t.Path, t.ID)
	if err != nil {
		return HuntReport{}
	}
	if code, err := r.TouchesCode(shas); err != nil || !code {
		return HuntReport{}
	}
	rep := HuntReport{CodeChanged: true}
	rep.Head, _ = r.HeadSHA()
	if h := r.LastHunt(t); h != nil {
		rep.State, rep.At = h.State, h.At
	}
	return rep
}
