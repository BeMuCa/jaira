// Package settings holds the few choices that belong to a person rather than
// to a board.
//
// It exists because two of them arrived with the ticket refs: which remote the
// refs go to, and whether an assignment pops up a desktop notification. Neither
// belongs in the repository — one teammate wanting to be notified says nothing
// about another, and a remote name is a property of a clone, not of the project.
//
// That last point is why the remote here is only a default. This file is read
// by every board opened on the machine, so "remote": "upstream" — right for the
// one checkout that has a fork as its origin — reached every other board too
// and took down every ref command on repositories that have no upstream at all.
// The name of a remote is a property of one clone, so that is where the answer
// for one board lives: git config jaira.remote <name>, in the clone's own
// config. Not a file in .jaira/, which would be committed and therefore wrong
// for the next person, whose clone may know the same repository by another
// name; and not a path-keyed section here, which git worktrees would split into
// several entries for one clone.
//
// RemoteFor is the only way to ask: it resolves the two against each other.
// See its comment for the order and for why a per-board name never falls back.
// RemoteSourceFor is the same ladder when the caller also has to say which step
// answered.
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
	"github.com/BeMuCa/jaira/core/ticket"
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

	// LogbookDays is how many days back the board shows what went into the
	// logbook, under the terminal lane, and how far 'jaira logbook' lists
	// without --since. A pointer because 0 is a real answer — "show none" —
	// and has to be told apart from a file that never mentions it, which
	// means DefaultLogbookDays. Read it through LogbookWindow.
	LogbookDays *int `json:"logbook-days,omitempty"`

	// Dispatch picks how a ticket is carried through its lanes when an agent
	// is handed the whole of it: "worker", one agent through every lane with
	// only a lane marked context: fresh handed to a new one, through
	// jaira-worker, or "lanes", one new agent per lane through
	// jaira-dispatcher. Empty or anything else means "worker"; only an
	// explicit "lanes" keeps the dispatcher. Read it through DispatchMode.
	Dispatch string `json:"dispatch,omitempty"`

	// ChatModel is the model the cockpit's chats start on — the explain and
	// discuss tabs, which write a few lines for a person and need no more than
	// a cheaper model. It is handed to 'claude --model' as it stands, so an
	// alias (sonnet) and a full id both work. Empty means no --model at all:
	// the account's own default, as before there was a setting. No lane
	// reads it; their model comes from the lane's tier.
	//
	// This is only the machine's fallback: a board that names its own in
	// .jaira/chat-model wins, because what a chat about the board's tickets
	// costs is the team's choice. Read it through ChatModelFor.
	ChatModel string `json:"chat-model,omitempty"`

	// Python is the command a role runs its Python scripts with — the
	// acceptance page, the logbook summary. A command, not a path, so
	// "uv run python" works as well as "py -3": on Windows python3 is often
	// missing or the Store's stub, which prints "Python" and exits 49
	// (issue #79). Empty means python3. JAIRA_PYTHON in the environment wins
	// over it. Read it through PythonCommand.
	Python string `json:"python,omitempty"`
}

// PythonEnv overrides the python setting for one shell or one machine.
const PythonEnv = "JAIRA_PYTHON"

// PythonCommand returns the command roles run Python with, and where it came
// from in words a person can read: JAIRA_PYTHON, then settings.json, else
// python3.
func (s Settings) PythonCommand() (cmd, source string) {
	if v := strings.TrimSpace(os.Getenv(PythonEnv)); v != "" {
		return v, "from " + PythonEnv
	}
	if v := strings.TrimSpace(s.Python); v != "" {
		return v, "from settings.json on this machine"
	}
	return "python3", "the default"
}

// ChatModelChoices are the values the settings screen steps through, the
// empty account default first. Anything else is written by hand.
var ChatModelChoices = []string{"", "sonnet", "haiku", "opus"}

// chatModelFile is the board's own chat model under .jaira: one line, the
// model, shared with everybody who works the board. Its own file, like
// lanes/order and tags, because a line is the smallest thing git merges.
const chatModelFile = "chat-model"

// BoardChatModelPath is where the board at root names its chat model.
func BoardChatModelPath(root string) string {
	return filepath.Join(root, ticket.DirName, chatModelFile)
}

// chatModelHeader is written above the model, for whoever opens the file.
const chatModelHeader = `# The model jaira's cockpit starts its explain and discuss chats on, for
# everybody working this board: one line, as 'claude --model' takes it
# (sonnet, haiku, opus, or a full model id). Without this file each machine's
# "chat-model" in ~/.jaira/settings.json applies, and without that the account
# default.
`

// BoardChatModel returns the model the board at root names, or "" when it
// names none, or names something that is not a model name.
func BoardChatModel(root string) string {
	return modelName(boardLine(BoardChatModelPath(root)))
}

// boardLine returns the first line of a one-line board file that is neither
// blank nor a # comment, trimmed, or "" when there is none or no file.
func boardLine(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			return line
		}
	}
	return ""
}

// SetBoardChatModel writes the board's chat model; "" removes the file, so
// the machine's fallback applies again.
func SetBoardChatModel(root, model string) error {
	path := BoardChatModelPath(root)
	if model == "" {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	return os.WriteFile(path, []byte(chatModelHeader+model+"\n"), 0o644)
}

// ChatModelFor returns the model the chats on the board at root start on,
// and where it came from in words a person can read: the board's file first,
// then this machine's settings.json, else "" — no --model, the account's
// own default. root may be "" outside any board.
func (s Settings) ChatModelFor(root string) (name, source string) {
	if root != "" {
		if m := BoardChatModel(root); m != "" {
			return m, "from .jaira/chat-model, shared with the board"
		}
	}
	if m := modelName(s.ChatModel); m != "" {
		return m, "from settings.json on this machine"
	}
	return "", "nothing configured"
}

// modelName returns v trimmed, or "" when it is not a model name. spawn.sh
// puts it into the command line a tab runs, so a value with anything but
// letters, digits and . _ - [ ] is treated as unset rather than carried into
// a shell.
func modelName(v string) string {
	v = strings.TrimSpace(v)
	for _, r := range v {
		ok := r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' ||
			strings.ContainsRune("._-[]", r)
		if !ok {
			return ""
		}
	}
	return v
}

// The two values of the board's release file.
const (
	ReleaseWorker   = "worker"
	ReleaseTeamlead = "teamlead"
)

// releaseFile is the board's answer to who touches the repository's version
// and changelog: "teamlead" when they are set once, at release, by whoever
// cuts it; "worker" when every ticket brings its own. A board file and not a
// field in settings.json, because it is a property of the repository — one
// machine's answer would reach every other repository on it, which is how
// "remote" once broke boards that had no upstream.
const releaseFile = "release"

// BoardReleasePath is where the board at root says who owns the release.
func BoardReleasePath(root string) string {
	return filepath.Join(root, ticket.DirName, releaseFile)
}

// ReleaseFor returns ReleaseTeamlead when the board at root names it, and
// ReleaseWorker otherwise, with where it came from in words a person can read.
// Worker is the default because it is what ran before there was a choice:
// every ticket may bump the version and write the changelog itself, as the
// repository's own instructions say. A typo falls back to it too, like
// DispatchMode, rather than stopping a worker from starting.
func ReleaseFor(root string) (owner, source string) {
	if root != "" && boardLine(BoardReleasePath(root)) == ReleaseTeamlead {
		return ReleaseTeamlead, "from .jaira/release, shared with the board"
	}
	return ReleaseWorker, "the default"
}

// The two values of Dispatch.
const (
	DispatchLanes  = "lanes"
	DispatchWorker = "worker"
)

// DispatchMode returns DispatchLanes only when the settings ask for it and
// DispatchWorker otherwise — a typo falls back to the default rather than
// stopping a teamlead from starting anything.
func (s Settings) DispatchMode() string {
	if strings.TrimSpace(s.Dispatch) == DispatchLanes {
		return DispatchLanes
	}
	return DispatchWorker
}

// DispatchRole is the role a teamlead starts for one ticket under that mode.
func (s Settings) DispatchRole() string {
	if s.DispatchMode() == DispatchWorker {
		return "jaira-worker"
	}
	return "jaira-dispatcher"
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

// RemoteFor returns the remote this board's ticket refs travel on.
//
// The order, and why each step is where it is:
//
//  1. git config jaira.remote in the clone. Somebody decided this for this
//     repository, so it is used as given and never falls back. If that remote
//     is gone, the ref command stops and says so — see Repo.Usable. A silent
//     fallback here would be the expensive kind of wrong: in a fork, origin is
//     the fork, and quietly pushing a ticket ref there loses it exactly where
//     nobody looks.
//  2. The machine-wide "remote" in settings.json, if this repository has a
//     remote by that name. It is a default for every board, so it applies here
//     only when it is true here.
//  3. The repository's only remote, when it has exactly one. There is nothing
//     to get wrong: one remote is where anything can go. This is what makes a
//     plain single-remote checkout work without anybody configuring anything,
//     even while settings.json names a remote it has never heard of.
//  4. Otherwise the configured name is returned unchanged, so the failure is
//     loud and names it. Several remotes and none of them the one asked for is
//     ambiguous, and guessing among them is the fork case again.
func (s Settings) RemoteFor(dir string) string {
	name, _ := s.RemoteSourceFor(dir)
	return name
}

// RemoteSourceFor is RemoteFor with the step of the ladder that answered, in
// words a person can read.
//
// The two are one function because a command that explains which remote is used
// must not compute the name a second time: the ladder has four steps, it is
// changed as one thing, and a second copy of it means the command whose only job
// is to tell the truth about the remote is the command that names a different
// one than the code that failed.
func (s Settings) RemoteSourceFor(dir string) (name, source string) {
	if board := strings.TrimSpace(gitref.BoardRemote(dir)); board != "" {
		return board, "from git config jaira.remote, set for this clone"
	}
	want := strings.TrimSpace(s.Remote)
	configured := want != ""
	if !configured {
		want = gitref.DefaultRemote
	}
	have := gitref.Remotes(dir)
	for _, n := range have {
		if n == want {
			if configured {
				return want, "from settings.json on this machine"
			}
			return want, "the default, nothing configured"
		}
	}
	if len(have) == 1 {
		return have[0], "the only remote this repository has"
	}
	if configured {
		return want, "from settings.json on this machine, and not a remote here"
	}
	return want, "the default, nothing configured"
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

// DefaultLogbookDays is the logbook window nobody chose: four weeks, what
// 'jaira logbook' listed before the window became a setting.
const DefaultLogbookDays = 28

// LogbookWindow returns how many days of the logbook the board shows and
// 'jaira logbook' lists: the setting, or DefaultLogbookDays when it is absent
// or negative. Zero means none on the board — and everything in the listing,
// which is what --since 0 already means there.
func (s Settings) LogbookWindow() int {
	if s.LogbookDays == nil || *s.LogbookDays < 0 {
		return DefaultLogbookDays
	}
	return *s.LogbookDays
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
			// A git refname is slash-separated on every platform; it is not
			// a filesystem path.
			//wintrap:ok
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
