package pr

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type glUser struct {
	Username string `json:"username"`
}

type glMR struct {
	IID          int    `json:"iid"`
	Title        string `json:"title"`
	WebURL       string `json:"web_url"`
	State        string `json:"state"`
	Draft        bool   `json:"draft"`
	SourceBranch string `json:"source_branch"`
	HeadPipeline *struct {
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
// carries none of them.
func readGitLab(call func(...string) ([]byte, error), origin string, q Query) (*Status, error) {
	repo, err := gitlabTarget(call, origin)
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
		path := project + "/merge_requests?state=opened&per_page=50&order_by=updated_at"
		if q.Branch != "" {
			path += "&source_branch=" + url.QueryEscape(q.Branch)
		}
		if err := get(path, &list); err != nil {
			return nil, err
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
		var discussions []glDiscussion
		if err := get(mr+"/discussions?per_page=100", &discussions); err != nil {
			return nil, err
		}
		var reviewers []glReviewer
		if err := get(mr+"/reviewers", &reviewers); err != nil {
			return nil, err
		}
		var commits []struct {
			Title string `json:"title"`
		}
		if err := get(mr+"/commits?per_page=100", &commits); err != nil {
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

func gitlabTarget(call func(...string) ([]byte, error), origin string) (string, error) {
	out, err := call("repo", "view", origin, "-F", "json")
	if err != nil {
		return "", err
	}
	var v struct {
		PathWithNamespace string `json:"path_with_namespace"`
		ForkedFrom        *struct {
			PathWithNamespace string `json:"path_with_namespace"`
		} `json:"forked_from_project"`
	}
	if err := json.Unmarshal(out, &v); err != nil {
		return "", fmt.Errorf("glab repo view: unreadable answer: %v", err)
	}
	if v.ForkedFrom != nil && v.ForkedFrom.PathWithNamespace != "" {
		return v.ForkedFrom.PathWithNamespace, nil
	}
	if v.PathWithNamespace == "" {
		return "", fmt.Errorf("glab repo view: no project named for %s", origin)
	}
	return v.PathWithNamespace, nil
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
	if m.HeadPipeline != nil {
		p.CI = gitlabCI(m.HeadPipeline.Status)
	}
	p.ReviewRequested = orEmpty(p.ReviewRequested)
	p.ChangesRequested = orEmpty(p.ChangesRequested)
	p.Approved = orEmpty(p.Approved)
	settle(&p)
	return p
}

func gitlabCI(status string) string {
	switch status {
	case "success":
		return "success"
	case "failed":
		return "failure"
	case "created", "waiting_for_resource", "preparing", "pending", "running", "scheduled", "manual":
		return "pending"
	}
	return "none"
}
