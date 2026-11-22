package pr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
)

type glUser struct {
	Username string `json:"username"`
}

type glMR struct {
	IID             int    `json:"iid"`
	Title           string `json:"title"`
	WebURL          string `json:"web_url"`
	State           string `json:"state"`
	Draft           bool   `json:"draft"`
	SourceBranch    string `json:"source_branch"`
	SourceProjectID int    `json:"source_project_id"`
	SHA             string `json:"sha"`
	HeadPipeline    *struct {
		Status string `json:"status"`
	} `json:"head_pipeline"`
}

type glDiscussion struct {
	IndividualNote bool `json:"individual_note"`
	Notes          []struct {
		Author     *glUser `json:"author"`
		CreatedAt  string  `json:"created_at"`
		System     bool    `json:"system"`
		Resolvable bool    `json:"resolvable"`
		Resolved   bool    `json:"resolved"`
	} `json:"notes"`
}

type glReviewer struct {
	User  *glUser `json:"user"`
	State string  `json:"state"`
}

// readGitLab reads merge requests through 'glab api': the list, then per
// merge request its pipeline, discussions, reviewers and commits — the list
// carries none of them. Every list is read to its last page (--paginate):
// on GitLab each comment outside a thread is a discussion of its own, and a
// thread left open on page two is what settled must not miss.
func readGitLab(call func(...string) ([]byte, error), origin string, q Query) (*Status, error) {
	repo, own, err := gitlabTarget(call, origin)
	if err != nil {
		return nil, err
	}
	project := "projects/" + url.PathEscape(repo)
	get := func(path string, into any) error {
		out, err := call("api", path)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(out, into); err != nil {
			return fmt.Errorf("glab api %s: unreadable answer: %v", path, err)
		}
		return nil
	}
	var me glUser
	if err := get("user", &me); err != nil {
		return nil, err
	}
	st := &Status{Forge: "gitlab", Repo: repo, Viewer: me.Username, Branch: q.Branch, PRs: []PR{}}
	var list []glMR
	if q.Number > 0 {
		list = make([]glMR, 1)
		if err := get(project+"/merge_requests/"+strconv.Itoa(q.Number), &list[0]); err != nil {
			return nil, err
		}
	} else {
		path := project + "/merge_requests?state=opened&per_page=100&order_by=updated_at"
		if q.Branch != "" {
			path += "&source_branch=" + url.QueryEscape(q.Branch)
		}
		found, err := pages[glMR](call, path)
		if err != nil {
			return nil, err
		}
		for _, m := range found {
			// source_branch names a branch in any project: a fork's branch
			// of the same name is somebody else's. The one asked for is
			// origin's.
			if q.Branch != "" && m.SourceProjectID != own {
				continue
			}
			list = append(list, m)
		}
	}
	for _, m := range list {
		mr := project + "/merge_requests/" + strconv.Itoa(m.IID)
		// The list leaves head_pipeline out; the single merge request has it.
		if q.Number == 0 {
			if err := get(mr, &m); err != nil {
				return nil, err
			}
		}
		discussions, err := pages[glDiscussion](call, mr+"/discussions?per_page=100")
		if err != nil {
			return nil, err
		}
		reviewers, err := pages[glReviewer](call, mr+"/reviewers?per_page=100")
		if err != nil {
			return nil, err
		}
		commits, err := pages[struct {
			Title string `json:"title"`
		}](call, mr+"/commits?per_page=100")
		if err != nil {
			return nil, err
		}
		var subjects []string
		for _, c := range commits {
			subjects = append(subjects, c.Title)
		}
		st.PRs = append(st.PRs, fromGitLab(m, discussions, reviewers, subjects, st.Viewer))
	}
	return st, nil
}

// pages reads every page of a list (--paginate) — glab prints one JSON
// array per page, one after the other — into one slice.
func pages[T any](call func(...string) ([]byte, error), path string) ([]T, error) {
	out, err := call("api", "--paginate", path)
	if err != nil {
		return nil, err
	}
	var list []T
	dec := json.NewDecoder(bytes.NewReader(out))
	for n := 0; ; n++ {
		var page []T
		if err := dec.Decode(&page); err == io.EOF && n > 0 {
			return list, nil
		} else if err != nil {
			return nil, fmt.Errorf("glab api %s: unreadable answer: %v", path, err)
		}
		list = append(list, page...)
	}
}

func gitlabTarget(call func(...string) ([]byte, error), origin string) (target string, own int, err error) {
	out, err := call("repo", "view", origin, "-F", "json")
	if err != nil {
		return "", 0, err
	}
	var v struct {
		ID                int    `json:"id"`
		PathWithNamespace string `json:"path_with_namespace"`
		ForkedFrom        *struct {
			PathWithNamespace string `json:"path_with_namespace"`
		} `json:"forked_from_project"`
	}
	if err := json.Unmarshal(out, &v); err != nil {
		return "", 0, fmt.Errorf("glab repo view: unreadable answer: %v", err)
	}
	if v.PathWithNamespace == "" {
		return "", 0, fmt.Errorf("glab repo view: no project named for %s", origin)
	}
	if v.ForkedFrom != nil && v.ForkedFrom.PathWithNamespace != "" {
		return v.ForkedFrom.PathWithNamespace, v.ID, nil
	}
	return v.PathWithNamespace, v.ID, nil
}

func username(u *glUser) string {
	if u == nil {
		return "ghost"
	}
	return u.Username
}

func fromGitLab(m glMR, discussions []glDiscussion, reviewers []glReviewer, subjects []string, viewer string) PR {
	state := m.State
	if state == "opened" {
		state = "open"
	}
	p := PR{Number: m.IID, Title: m.Title, URL: m.WebURL, State: state, Branch: m.SourceBranch, Draft: m.Draft}
	p.Tickets = ticketsOf(subjects, m.SourceBranch)
	for _, d := range discussions {
		thread, resolved, last := false, true, ""
		for _, n := range d.Notes {
			if n.System {
				continue
			}
			a := username(n.Author)
			last = a
			if n.Resolvable {
				thread = true
				resolved = resolved && n.Resolved
			}
			if a != viewer {
				p.Comments++
				kind := "comment"
				if !d.IndividualNote {
					kind = "thread"
				}
				seen(&p.LastForeign, Activity{At: n.CreatedAt, Author: a, Kind: kind})
			}
		}
		switch {
		case !thread:
		case resolved:
			p.Threads.Resolved++
		case last == viewer:
			p.Threads.Answered++
		default:
			p.Threads.Waiting++
		}
	}
	for _, r := range reviewers {
		switch name := username(r.User); r.State {
		case "requested_changes":
			p.ChangesRequested = append(p.ChangesRequested, name)
		case "approved":
			p.Approved = append(p.Approved, name)
		case "unreviewed", "unapproved", "review_started":
			p.ReviewRequested = append(p.ReviewRequested, name)
		}
	}
	p.CI = "none"
	p.HeadSHA = m.SHA
	if m.HeadPipeline != nil {
		p.CI = gitlabCI(m.HeadPipeline.Status)
	}
	p.ReviewRequested = orEmpty(p.ReviewRequested)
	p.ChangesRequested = orEmpty(p.ChangesRequested)
	p.Approved = orEmpty(p.Approved)
	settle(&p)
	return p
}

// gitlabCI sorts a pipeline status as githubCI sorts a rollup state. A
// canceled pipeline is red, as a cancelled check is on GitHub, so it holds
// settled; skipped ran nothing and counts as no checks.
func gitlabCI(status string) string {
	switch status {
	case "success":
		return "success"
	case "failed", "canceled":
		return "failure"
	case "created", "waiting_for_resource", "preparing", "pending", "running", "scheduled", "manual":
		return "pending"
	}
	return "none"
}
