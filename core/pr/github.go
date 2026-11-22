package pr

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// One query for everything a pull request's stand needs: threads with their
// newest comments, the conversation, reviews, requested reviewers, each
// reviewer's latest verdict, the commit subjects (for the tickets) and the
// head commit's check rollup. Bounded so a repository with many pull
// requests stays inside GitHub's node limit: 50 pull requests, 100 threads of
// 20 comments each.
const prFields = `number title url state isDraft headRefName
 reviewRequests(first:20){nodes{requestedReviewer{__typename ... on User{login} ... on Bot{login} ... on Mannequin{login} ... on Team{slug}}}}
 latestOpinionatedReviews(first:20){nodes{state author{login}}}
 reviews(last:30){nodes{state submittedAt author{login}}}
 comments(last:50){nodes{createdAt author{login}}}
 reviewThreads(first:100){nodes{isResolved comments(last:20){nodes{createdAt author{login}}}}}
 subjects:commits(first:100){nodes{commit{messageHeadline}}}
 head:commits(last:1){nodes{commit{statusCheckRollup{state}}}}`

const listQuery = `query($owner:String!,$name:String!,$head:String){viewer{login}
 repository(owner:$owner,name:$name){pullRequests(states:OPEN,headRefName:$head,first:50,orderBy:{field:UPDATED_AT,direction:DESC}){nodes{` + prFields + `}}}}`

const oneQuery = `query($owner:String!,$name:String!,$number:Int!){viewer{login}
 repository(owner:$owner,name:$name){pullRequest(number:$number){` + prFields + `}}}`

type ghLogin struct {
	Login string `json:"login"`
}

type ghPR struct {
	Number         int    `json:"number"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	State          string `json:"state"`
	IsDraft        bool   `json:"isDraft"`
	HeadRefName    string `json:"headRefName"`
	ReviewRequests struct {
		Nodes []struct {
			RequestedReviewer *struct {
				Login string `json:"login"`
				Slug  string `json:"slug"`
			} `json:"requestedReviewer"`
		} `json:"nodes"`
	} `json:"reviewRequests"`
	LatestOpinionatedReviews struct {
		Nodes []struct {
			State  string   `json:"state"`
			Author *ghLogin `json:"author"`
		} `json:"nodes"`
	} `json:"latestOpinionatedReviews"`
	Reviews struct {
		Nodes []struct {
			State       string   `json:"state"`
			SubmittedAt string   `json:"submittedAt"`
			Author      *ghLogin `json:"author"`
		} `json:"nodes"`
	} `json:"reviews"`
	Comments      ghComments `json:"comments"`
	ReviewThreads struct {
		Nodes []struct {
			IsResolved bool       `json:"isResolved"`
			Comments   ghComments `json:"comments"`
		} `json:"nodes"`
	} `json:"reviewThreads"`
	Subjects struct {
		Nodes []struct {
			Commit struct {
				MessageHeadline string `json:"messageHeadline"`
			} `json:"commit"`
		} `json:"nodes"`
	} `json:"subjects"`
	Head struct {
		Nodes []struct {
			Commit struct {
				StatusCheckRollup *struct {
					State string `json:"state"`
				} `json:"statusCheckRollup"`
			} `json:"commit"`
		} `json:"nodes"`
	} `json:"head"`
}

type ghComments struct {
	Nodes []struct {
		CreatedAt string   `json:"createdAt"`
		Author    *ghLogin `json:"author"`
	} `json:"nodes"`
}

type ghAnswer struct {
	Data *struct {
		Viewer     ghLogin `json:"viewer"`
		Repository *struct {
			PullRequests *struct {
				Nodes []ghPR `json:"nodes"`
			} `json:"pullRequests"`
			PullRequest *ghPR `json:"pullRequest"`
		} `json:"repository"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func readGitHub(call func(...string) ([]byte, error), origin string, q Query) (*Status, error) {
	repo, err := githubTarget(call, origin)
	if err != nil {
		return nil, err
	}
	owner, name, ok := ownerName(repo)
	if !ok {
		return nil, fmt.Errorf("gh repo view: %q is not owner/name", repo)
	}
	args := []string{"api", "graphql", "-f", "owner=" + owner, "-f", "name=" + name}
	switch {
	case q.Number > 0:
		args = append(args, "-F", "number="+strconv.Itoa(q.Number), "-f", "query="+oneQuery)
	case q.Branch != "":
		args = append(args, "-f", "head="+q.Branch, "-f", "query="+listQuery)
	default:
		args = append(args, "-f", "query="+listQuery)
	}
	out, err := call(args...)
	if err != nil {
		return nil, err
	}
	var a ghAnswer
	if err := json.Unmarshal(out, &a); err != nil {
		return nil, fmt.Errorf("gh api graphql: unreadable answer: %v", err)
	}
	if len(a.Errors) > 0 {
		return nil, fmt.Errorf("gh api graphql: %s", a.Errors[0].Message)
	}
	if a.Data == nil || a.Data.Repository == nil {
		return nil, fmt.Errorf("gh api graphql: no repository %s in the answer", repo)
	}
	st := &Status{Forge: "github", Repo: repo, Viewer: a.Data.Viewer.Login, Branch: q.Branch, PRs: []PR{}}
	var nodes []ghPR
	if q.Number > 0 {
		if a.Data.Repository.PullRequest == nil {
			return nil, fmt.Errorf("gh api graphql: %s has no pull request #%d", repo, q.Number)
		}
		nodes = []ghPR{*a.Data.Repository.PullRequest}
	} else if a.Data.Repository.PullRequests != nil {
		nodes = a.Data.Repository.PullRequests.Nodes
	}
	for _, n := range nodes {
		st.PRs = append(st.PRs, fromGitHub(n, st.Viewer))
	}
	return st, nil
}

// githubTarget is the repository pull requests of this checkout go to: the
// parent where origin is a fork, as jaira-role-pr and the cockpit settle it.
func githubTarget(call func(...string) ([]byte, error), origin string) (string, error) {
	out, err := call("repo", "view", origin, "--json", "nameWithOwner,parent")
	if err != nil {
		return "", err
	}
	var v struct {
		NameWithOwner string `json:"nameWithOwner"`
		Parent        *struct {
			Name  string  `json:"name"`
			Owner ghLogin `json:"owner"`
		} `json:"parent"`
	}
	if err := json.Unmarshal(out, &v); err != nil {
		return "", fmt.Errorf("gh repo view: unreadable answer: %v", err)
	}
	if v.Parent != nil && v.Parent.Owner.Login != "" && v.Parent.Name != "" {
		return v.Parent.Owner.Login + "/" + v.Parent.Name, nil
	}
	if v.NameWithOwner == "" {
		return "", fmt.Errorf("gh repo view: no repository named for %s", origin)
	}
	return v.NameWithOwner, nil
}

func login(l *ghLogin) string {
	if l == nil {
		return "ghost"
	}
	return l.Login
}

func fromGitHub(n ghPR, viewer string) PR {
	p := PR{
		Number: n.Number, Title: n.Title, URL: n.URL, State: strings.ToLower(n.State),
		Branch: n.HeadRefName, Draft: n.IsDraft,
	}
	var subjects []string
	for _, c := range n.Subjects.Nodes {
		subjects = append(subjects, c.Commit.MessageHeadline)
	}
	p.Tickets = ticketsOf(subjects, n.HeadRefName)
	for _, t := range n.ReviewThreads.Nodes {
		cs := t.Comments.Nodes
		last := ""
		if len(cs) > 0 {
			last = login(cs[len(cs)-1].Author)
		}
		switch {
		case t.IsResolved:
			p.Threads.Resolved++
		case last == viewer:
			p.Threads.Answered++
		default:
			p.Threads.Waiting++
		}
		for _, c := range cs {
			if a := login(c.Author); a != viewer {
				p.Comments++
				seen(&p.LastForeign, Activity{At: c.CreatedAt, Author: a, Kind: "thread"})
			}
		}
	}
	for _, c := range n.Comments.Nodes {
		if a := login(c.Author); a != viewer {
			p.Comments++
			seen(&p.LastForeign, Activity{At: c.CreatedAt, Author: a, Kind: "comment"})
		}
	}
	for _, r := range n.Reviews.Nodes {
		// A review still PENDING is a draft nobody else sees yet.
		if a := login(r.Author); a != viewer && r.State != "PENDING" {
			seen(&p.LastForeign, Activity{At: r.SubmittedAt, Author: a, Kind: "review", State: strings.ToLower(r.State)})
		}
	}
	for _, r := range n.ReviewRequests.Nodes {
		if r.RequestedReviewer == nil {
			continue
		}
		name := r.RequestedReviewer.Login
		if name == "" {
			name = r.RequestedReviewer.Slug
		}
		p.ReviewRequested = append(p.ReviewRequested, name)
	}
	for _, r := range n.LatestOpinionatedReviews.Nodes {
		switch r.State {
		case "CHANGES_REQUESTED":
			p.ChangesRequested = append(p.ChangesRequested, login(r.Author))
		case "APPROVED":
			p.Approved = append(p.Approved, login(r.Author))
		}
	}
	p.CI = "none"
	if hs := n.Head.Nodes; len(hs) > 0 && hs[0].Commit.StatusCheckRollup != nil {
		p.CI = githubCI(hs[0].Commit.StatusCheckRollup.State)
	}
	p.ReviewRequested = orEmpty(p.ReviewRequested)
	p.ChangesRequested = orEmpty(p.ChangesRequested)
	p.Approved = orEmpty(p.Approved)
	settle(&p)
	return p
}

func githubCI(state string) string {
	switch state {
	case "SUCCESS":
		return "success"
	case "FAILURE", "ERROR":
		return "failure"
	case "PENDING", "EXPECTED":
		return "pending"
	}
	return "none"
}
