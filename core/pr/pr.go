// Package pr reads where the open pull requests of a repository stand — review
// threads waiting for an answer, comments from others, CI, reviews still
// requested — from the forge's own command line tool, gh or glab. jaira keeps
// no copy: every call asks the forge, and nothing is written.
//
// It is what 'jaira pr' prints and what the cockpit polls to tell a teamlead
// session that a pull request needs its next round. 'settled' is the same
// question jaira-role-pr's 'pr-round.sh settled' answers: no thread
// unresolved, no review still requested, no reviewer's latest verdict asking
// for changes.
package pr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/BeMuCa/jaira/core/gitrepo"
)

// Errors a caller tells apart: each is a reason nothing could be read, not
// a pull request that is in some state.
var (
	// ErrNoForge: origin is missing, or names a host that is neither GitHub
	// nor GitLab and git config jaira.forge says nothing.
	ErrNoForge = errors.New("no forge")
	// ErrNoTool: gh or glab is not installed.
	ErrNoTool = errors.New("forge tool missing")
	// ErrAuth: the tool is not signed in, or its sign-in ran out.
	ErrAuth = errors.New("not signed in")
)

// Run runs one command in dir and returns its standard output. A failure
// carries the exit code — an *exec.ExitError, or any error with an
// ExitCode() method — and the command's standard error.
type Run func(dir, name string, args ...string) ([]byte, error)

type exitCoder interface{ ExitCode() int }

// Status is every open pull request the query found, newest activity first.
type Status struct {
	Forge  string `json:"forge"`
	Repo   string `json:"repo"`
	Viewer string `json:"viewer"`
	Branch string `json:"branch,omitempty"`
	PRs    []PR   `json:"prs"`
}

// Threads counts the review threads of one pull request: Waiting are not
// resolved and their last comment is somebody else's — a point nobody
// answered yet; Answered are not resolved but the viewer wrote last; Resolved
// are done.
type Threads struct {
	Waiting  int `json:"waiting"`
	Answered int `json:"answered"`
	Resolved int `json:"resolved"`
}

// Activity is the newest comment or review someone other than the viewer
// left. Kind is comment, thread or review; State the review's verdict.
type Activity struct {
	At     string `json:"at"`
	Author string `json:"author"`
	Kind   string `json:"kind"`
	State  string `json:"state,omitempty"`
}

// PR is one open pull request (merge request on GitLab).
type PR struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	// State is open, merged or closed.
	State   string   `json:"state"`
	URL     string   `json:"url"`
	Branch  string   `json:"branch"`
	Draft   bool     `json:"draft"`
	Tickets []string `json:"tickets"`
	Threads Threads  `json:"threads"`
	// Comments counts what others wrote, in threads and on the conversation.
	Comments    int       `json:"comments"`
	LastForeign *Activity `json:"last_foreign"`
	// CI is success, failure, pending or none, for the head commit HeadSHA.
	CI               string   `json:"ci"`
	HeadSHA          string   `json:"head_sha"`
	ReviewRequested  []string `json:"review_requested"`
	ChangesRequested []string `json:"changes_requested"`
	Approved         []string `json:"approved"`
	Settled          bool     `json:"settled"`
	Unsettled        []string `json:"unsettled"`
	// partial names what could not be read whole and settled depends on:
	// it keeps the pull request unsettled rather than settle it on half.
	partial []string
}

// Query says what to read: pull request Number in whatever state it is in,
// else the open ones whose head is Branch, else every open one.
type Query struct {
	Dir    string
	Branch string
	Number int
}

// Read asks the forge of the repository in q.Dir for its open pull requests.
func Read(run Run, look func(string) (string, error), q Query) (*Status, error) {
	forge, origin, err := forgeOf(run, q.Dir)
	if err != nil {
		return nil, err
	}
	tool := map[string]string{"github": "gh", "gitlab": "glab"}[forge]
	if _, err := look(tool); err != nil {
		return nil, fmt.Errorf("%w: %s is not installed", ErrNoTool, tool)
	}
	call := func(args ...string) ([]byte, error) {
		out, err := run(q.Dir, tool, args...)
		if err != nil {
			return nil, classify(tool, err)
		}
		return out, nil
	}
	if forge == "github" {
		return readGitHub(call, origin, q)
	}
	return readGitLab(call, origin, q)
}

// forgeOf decides the forge as the cockpit's forgeOf does: git config
// jaira.forge first, else the host of origin.
func forgeOf(run Run, dir string) (forge, origin string, err error) {
	out, err := run(dir, "git", "remote", "get-url", "origin")
	if err != nil {
		return "", "", fmt.Errorf("%w: this repository has no remote origin", ErrNoForge)
	}
	origin = strings.TrimSpace(string(out))
	set, _ := run(dir, "git", "config", "jaira.forge")
	forge = strings.ToLower(strings.TrimSpace(string(set)))
	if forge == "" {
		switch lower := strings.ToLower(origin); {
		case strings.Contains(lower, "github.com"):
			forge = "github"
		case strings.Contains(lower, "gitlab"):
			forge = "gitlab"
		}
	}
	if forge != "github" && forge != "gitlab" {
		return "", "", fmt.Errorf("%w: origin %s is neither GitHub nor GitLab — set git config jaira.forge github|gitlab", ErrNoForge, origin)
	}
	return forge, origin, nil
}

// signedOut is how gh and glab say the sign-in is missing or gone: an HTTP
// 401 from the API, GitHub's "Bad credentials", the hint to log in. A bare
// 401 is not one — "no pull request with the number of 401" is not about
// the sign-in.
var signedOut = regexp.MustCompile(`(?i)\bHTTP 401\b|\b401 Unauthorized\b|Bad credentials|auth login`)

// classify turns a failed gh/glab call into ErrAuth where it is a sign-in
// problem — gh ends with 4 then, both say so in signedOut's words — and keeps
// the tool's own words otherwise.
func classify(tool string, err error) error {
	code := -1
	var ec exitCoder
	if errors.As(err, &ec) {
		code = ec.ExitCode()
	}
	msg := err.Error()
	var ee *exec.ExitError
	if errors.As(err, &ee) && len(ee.Stderr) > 0 {
		msg = strings.TrimSpace(string(ee.Stderr))
	}
	if (tool == "gh" && code == 4) || signedOut.MatchString(msg) {
		return fmt.Errorf("%w: %s — run: %s auth login", ErrAuth, firstLine(msg), tool)
	}
	// gh and glab often start their own message with their name already.
	return fmt.Errorf("%s: %s", tool, strings.TrimPrefix(firstLine(msg), tool+": "))
}

func firstLine(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

// ticketsOf names the tickets a pull request carries: every handle its
// commit subjects name the way this board's commits do, then the branch's
// last part where it is a handle — spawn.sh names a ticket's branch
// <prefix>/<slug>.
func ticketsOf(subjects []string, branch string) []string {
	out := gitrepo.HandlesInSubjects(strings.Join(subjects, "\n"))
	if out == nil {
		out = []string{}
	}
	last := branch[strings.LastIndexByte(branch, '/')+1:]
	if h := strings.ToUpper(last); handle.MatchString(h) && !slices.Contains(out, h) {
		out = append(out, h)
	}
	return out
}

var handle = regexp.MustCompile(`^[0-9A-Z]{6}$`)

// values decodes what a --paginate call printed: one JSON value per page
// (gh and glab write the pages one after another), each handed to each.
func values(out []byte, what string, each func(json.RawMessage) error) error {
	dec := json.NewDecoder(bytes.NewReader(out))
	n := 0
	for {
		var v json.RawMessage
		if err := dec.Decode(&v); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("%s: unreadable answer: %v", what, err)
		}
		n++
		if err := each(v); err != nil {
			return err
		}
	}
	if n == 0 {
		return fmt.Errorf("%s: unreadable answer: empty", what)
	}
	return nil
}

// newer reports whether RFC 3339 time a is after b; an unreadable a never is.
func newer(a, b string) bool {
	ta, err := time.Parse(time.RFC3339, a)
	if err != nil {
		return false
	}
	tb, err := time.Parse(time.RFC3339, b)
	return err != nil || ta.After(tb)
}

// seen keeps the newest foreign activity.
func seen(last **Activity, a Activity) {
	if a.At == "" {
		return
	}
	if *last == nil || newer(a.At, (*last).At) {
		c := a
		*last = &c
	}
}

// settle fills Settled and Unsettled from the counts, in the words
// pr-round.sh settled uses.
func settle(p *PR) {
	p.Unsettled = []string{}
	switch open := p.Threads.Waiting + p.Threads.Answered; open {
	case 0:
	case 1:
		p.Unsettled = append(p.Unsettled, "1 open thread")
	default:
		p.Unsettled = append(p.Unsettled, fmt.Sprintf("%d open threads", open))
	}
	for _, r := range p.ReviewRequested {
		p.Unsettled = append(p.Unsettled, "review pending: "+r)
	}
	for _, r := range p.ChangesRequested {
		p.Unsettled = append(p.Unsettled, "changes requested: "+r)
	}
	p.Unsettled = append(p.Unsettled, p.partial...)
	p.Settled = len(p.Unsettled) == 0
}

func orEmpty(s []string) []string {
	if s == nil {
		return []string{}
	}
	sort.Strings(s)
	return s
}

// ownerName splits owner/name; GitLab's nested groups keep every part but
// the last in owner.
func ownerName(repo string) (string, string, bool) {
	i := strings.LastIndexByte(repo, '/')
	if i <= 0 || i == len(repo)-1 {
		return "", "", false
	}
	return repo[:i], repo[i+1:], true
}
