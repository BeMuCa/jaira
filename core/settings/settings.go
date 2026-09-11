// Package settings holds the few choices that belong to a person rather than
// to a board.
//
// It exists because two of them arrived with the ticket refs: which remote the
// refs go to, and whether an assignment pops up a desktop notification. Neither
// belongs in the repository — one teammate wanting to be notified says nothing
// about another, and a remote name is a property of a clone, not of the project.
//
// The file is ~/.jaira/settings.json, beside projects.json, and every field has
// a working default: a missing file, an unreadable one and an empty one all mean
// "the defaults", because settings a user never opened must never be the reason
// a board does not start.
package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/snapshot"
)

// Settings is the whole file.
type Settings struct {
	// Remote is where ticket refs are pushed and fetched. Empty means origin.
	Remote string `json:"remote,omitempty"`

	// Notify turns the desktop notification for an assignment off. It is
	// phrased as the negative so that the zero value — a file that does not
	// mention it, or no file at all — leaves notifications on: the feature was
	// asked for, and a default that silently withholds it would answer the
	// request with nothing.
	NotifyOff bool `json:"notify-off,omitempty"`

	// Hook is a script called when a ticket moves or is claimed, for immediate
	// delivery through whatever the user already uses (a Slack webhook,
	// ntfy.sh, Telegram). jaira only calls it and brings no dependency of its
	// own; an empty value calls nothing.
	Hook string `json:"hook,omitempty"`

	// LandingBranches are the branches whose contents mean a ticket has
	// arrived: filed away in one of them, it is finished for everybody and its
	// ref has done its job.
	//
	// A list with globs rather than one "main branch", because there is no
	// answer to what the important branch is called — main, master, develop, a
	// release line — and a list turns a guess into a setting. Empty falls back
	// to the remote's own HEAD, and if that cannot be resolved either, nothing
	// is ever removed.
	LandingBranches []string `json:"landing-branches,omitempty"`

	// SnapshotBranch holds the board as files, as a backup for whoever clones
	// for the first time and has no refs yet. Empty means jaira/board.
	SnapshotBranch string `json:"snapshot-branch,omitempty"`

	// SnapshotEvery is how old the backup may get, as a duration ("72h").
	// Empty means three days: every participant's clone holds the refs, so the
	// board survives any one machine, and this is for the slow cases only.
	SnapshotEvery string `json:"snapshot-every,omitempty"`

	// FetchEvery is how often the ticket refs are brought up to date in the
	// background, as a duration ("10m"). Empty means ten minutes.
	//
	// One setting for the CLI and the board alike. They used to disagree — the
	// board had its own hard-wired minute — and somebody who changed this would
	// have found the board carrying on at its own pace.
	FetchEvery string `json:"fetch-every,omitempty"`

	// LandingGrace is how long a finished ticket may go without arriving in a
	// landing branch before jaira mentions it, as a duration ("72h"). Empty
	// means three days.
	LandingGrace string `json:"landing-grace,omitempty"`
}

// Path is the settings file, honouring JAIRA_HOME so tests and a sandboxed run
// do not touch the real one.
func Path() string {
	if v := os.Getenv("JAIRA_HOME"); v != "" {
		return filepath.Join(v, "settings.json")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".jaira", "settings.json")
}

// Load returns the settings, or the defaults for anything missing.
//
// An unreadable or malformed file is deliberately not an error: it is one
// person's preferences, and refusing to open a board over it would trade a
// small annoyance for a total one.
func Load() Settings {
	var s Settings
	path := Path()
	if path == "" {
		return s
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return s
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return Settings{}
	}
	return s
}

// Save writes the settings, creating ~/.jaira if it is not there yet.
func Save(s Settings) error {
	path := Path()
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

// RemoteOr returns the configured remote, or the default.
func (s Settings) RemoteName() string {
	if r := strings.TrimSpace(s.Remote); r != "" {
		return r
	}
	return gitref.DefaultRemote
}

// NotifyEnabled reports whether an assignment should raise a desktop
// notification.
func (s Settings) NotifyEnabled() bool { return !s.NotifyOff }

// SnapshotBranchName returns the configured snapshot branch, or the default.
func (s Settings) SnapshotBranchName() string {
	if b := strings.TrimSpace(s.SnapshotBranch); b != "" {
		return b
	}
	return snapshot.DefaultBranch
}

// every parses one of the interval settings, falling back to the default.
//
// A value nobody can parse is not worth refusing to start over: it is one
// person's preference file, the cost of ignoring it is that a background job
// keeps its usual pace, and the cost of failing on it is a board that will not
// open. Zero and negative are treated the same way, since "every 0s" is not an
// interval anybody means.
func every(value string, fallback time.Duration) time.Duration {
	v := strings.TrimSpace(value)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		return fallback
	}
	return d
}

// SnapshotInterval returns how old a snapshot may get.
func (s Settings) SnapshotInterval() time.Duration {
	return every(s.SnapshotEvery, snapshot.DefaultEvery)
}

// FetchInterval returns how often the refs are refreshed in the background,
// for the CLI and the board alike.
func (s Settings) FetchInterval() time.Duration {
	return every(s.FetchEvery, refsync.DefaultFetchEvery)
}

// DefaultLandingGrace is how long a finished ticket may take to arrive before
// it is worth mentioning.
//
// Three days, chosen against how long a review actually takes here rather than
// as a round number: a week is long enough that a forgotten branch stops being
// news by the time anybody hears about it.
const DefaultLandingGrace = 3 * 24 * time.Hour

// LandingGraceInterval returns that, or whatever the settings say instead.
func (s Settings) LandingGraceInterval() time.Duration {
	return every(s.LandingGrace, DefaultLandingGrace)
}

// Landing returns the branches to check for a landed ticket, as revisions on
// the remote.
//
// The configured list wins over the remote's HEAD, and that order matters: the
// HEAD is the hosting provider's default, not necessarily the branch a team
// cares about. Both are prefixed with the remote, since what counts is what
// has arrived where everybody can see it — not what has landed in somebody's
// local copy of a branch.
func (s Settings) Landing(remote string, remoteHead func() string) []string {
	if remote == "" {
		remote = gitref.DefaultRemote
	}
	if len(s.LandingBranches) > 0 {
		out := make([]string, 0, len(s.LandingBranches))
		for _, b := range s.LandingBranches {
			b = strings.TrimSpace(b)
			if b == "" {
				continue
			}
			if strings.HasPrefix(b, remote+"/") || strings.HasPrefix(b, "refs/") {
				out = append(out, b)
				continue
			}
			out = append(out, remote+"/"+b)
		}
		return out
	}
	if remoteHead != nil {
		if head := strings.TrimSpace(remoteHead()); head != "" {
			return []string{head}
		}
	}
	// Nothing resolved: no branch, and therefore nothing ever removed.
	return nil
}
