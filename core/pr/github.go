package pr

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// One query for everything a pull request's stand needs: threads with their
// newest comments, the conversation, reviews, requested reviewers, each
// reviewer's latest verdict, the commit subjects (for the tickets), the head
// repository (to tell a fork's branch from the parent's of the same name) and
// the head commit with its check rollup. Pull requests come 25 to a page and
// threads 100 to a page, both read to the end (GitHub's node limit allows no
// more in one query); what settled depends on and cannot be paged — more
// than 100 reviewers requested or with a verdict — makes the pull request
// unsettled rather than read short.
const prFields = `number title url state isDraft headRefName headRepository{nameWithOwner}
 reviewRequests(first:100){pageInfo{hasNextPage} nodes{requestedReviewer{__typename ... on User{login} ... on Bot{login} ... on Mannequin{login} ... on Team{slug}}}}
 latestOpinionatedReviews(first:100){pageInfo{hasNextPage} nodes{state author{login}}}
 reviews(last:30){nodes{state submittedAt author{login}}}
 comments(last:50){nodes{createdAt author{login}}}
 reviewThreads(first:100){pageInfo{hasNextPage endCursor} nodes{isResolved comments(last:20){nodes{createdAt author{login}}}}}
 subjects:commits(first:100){nodes{commit{messageHeadline}}}
 head:commits(last:1){nodes{commit{oid statusCheckRollup{state}}}}`

const listQuery = `query($owner:String!,$name:String!,$head:String,$after:String){viewer{login}
 repository(owner:$owner,name:$name){pullRequests(states:OPEN,headRefName:$head,first:25,after:$after,orderBy:{field:UPDATED_AT,direction:DESC}){pageInfo{hasNextPage endCursor} nodes{` + prFields + `}}}}`

const oneQuery = `query($owner:String!,$name:String!,$number:Int!){viewer{login}
 repository(owner:$owner,name:$name){pullRequest(number:$number){` + prFields + `}}}`

const threadsQuery = `query($owner:String!,$name:String!,$number:Int!,$after:String){
 repository(owner:$owner,name:$name){pullRequest(number:$number){
 reviewThreads(first:100,after:$after){pageInfo{hasNextPage endCursor} nodes{isResolved comments(last:20){nodes{createdAt author{login}}}}}}}}`

type ghPage struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
}

type ghThreads struct {
	PageInfo ghPage `json:"pageInfo"`
	Nodes    []struct {
		IsResolved bool       `json:"isResolved"`
		Comments   ghComments `json:"comments"`
	} `json:"nodes"`
}

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
	HeadRepository *struct {
		NameWithOwner string `json:"nameWithOwner"`
	} `json:"headRepository"`
	ReviewRequests struct {
		PageInfo ghPage `json:"pageInfo"`
		Nodes    []struct {
			RequestedReviewer *struct {
				Login string `json:"login"`
				Slug  string `json:"slug"`
			} `json:"requestedReviewer"`
		} `json:"nodes"`
	} `json:"reviewRequests"`
	LatestOpinionatedReviews struct {
		PageInfo ghPage `json:"pageInfo"`
		Nodes    []struct {
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
	ReviewThreads ghThreads  `json:"reviewThreads"`
	Subjects      struct {
		Nodes []struct {
			Commit struct {
				MessageHeadline string `json:"messageHeadline"`
			} `json:"commit"`
		} `json:"nodes"`
	} `json:"subjects"`
	Head struct {
		Nodes []struct {
			Commit struct {
				OID               string `json:"oid"`
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
				PageInfo ghPage `json:"pageInfo"`
				Nodes    []ghPR `json:"nodes"`
			} `json:"pullRequests"`
			PullRequest *ghPR `json:"pullRequest"`
		} `json:"repository"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func readGitHub(call func(...string) ([]byte, error), origin string, q Query) (*Status, error) {
	repo, own, err := githubTarget(call, origin)
	if err != nil {
		return nil, err
	}
	owner, name, ok := ownerName(repo)
	if !ok {
		return nil, fmt.Errorf("gh repo view: %q is not owner/name", repo)
	}
	ask := func(extra ...string) (*ghAnswer, error) {
		out, err := call(append([]string{"api", "graphql", "-f", "owner=" + owner, "-f", "name=" + name}, extra...)...)
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
		return &a, nil
	}
	st := &Status{Forge: "github", Repo: repo, Branch: q.Branch, PRs: []PR{}}
	var nodes []ghPR
	if q.Number > 0 {
		a, err := ask("-F", "number="+strconv.Itoa(q.Number), "-f", "query="+oneQuery)
		if err != nil {
			return nil, err
		}
		if a.Data.Repository.PullRequest == nil {
			return nil, fmt.Errorf("gh api graphql: %s has no pull request #%d", repo, q.Number)
		}
		st.Viewer = a.Data.Viewer.Login
		nodes = []ghPR{*a.Data.Repository.PullRequest}
	} else {
		after := ""
		for {
			args := []string{"-f", "query=" + listQuery}
			if q.Branch != "" {
				args = append(args, "-f", "head="+q.Branch)
			}
			if after != "" {
				args = append(args, "-f", "after="+after)
			}
			a, err := ask(args...)
			if err != nil {
				return nil, err
			}
			st.Viewer = a.Data.Viewer.Login
			list := a.Data.Repository.PullRequests
			if list == nil {
				break
			}
			for _, n := range list.Nodes {
				// headRefName names a branch in any repository: on a fork's
				// parent, master or feat/x of somebody else's fork matches
				// too. The branch asked for is the one in origin.
				if q.Branch != "" && (n.HeadRepository == nil || !strings.EqualFold(n.HeadRepository.NameWithOwner, own)) {
					continue
				}
				nodes = append(nodes, n)
			}
			if !list.PageInfo.HasNextPage || list.PageInfo.EndCursor == "" {
				break
			}
			after = list.PageInfo.EndCursor
		}
	}
	for i := range nodes {
		// Threads past the first hundred: a thread left open on page two
		// is exactly what settled must not miss.
		for page := nodes[i].ReviewThreads.PageInfo; page.HasNextPage && page.EndCursor != ""; {
			a, err := ask("-F", "number="+strconv.Itoa(nodes[i].Number), "-f", "after="+page.EndCursor, "-f", "query="+threadsQuery)
			if err != nil {
				return nil, err
			}
			pr := a.Data.Repository.PullRequest
			if pr == nil {
				return nil, fmt.Errorf("gh api graphql: pull request #%d went away while its threads were read", nodes[i].Number)
			}
			nodes[i].ReviewThreads.Nodes = append(nodes[i].ReviewThreads.Nodes, pr.ReviewThreads.Nodes...)
			page = pr.ReviewThreads.PageInfo
		}
		st.PRs = append(st.PRs, fromGitHub(nodes[i], st.Viewer))
	}
	return st, nil
}

// githubTarget is the repository pull requests of this checkout go to: the
// parent where origin is a fork, as jaira-role-pr and the cockpit settle it.
func githubTarget(call func(...string) ([]byte, error), origin string) (target, own string, err error) {
	out, err := call("repo", "view", origin, "--json", "nameWithOwner,parent")
	if err != nil {
		return "", "", err
	}
	var v struct {
		NameWithOwner string `json:"nameWithOwner"`
		Parent        *struct {
			Name  string  `json:"name"`
			Owner ghLogin `json:"owner"`
		} `json:"parent"`
	}
	if err := json.Unmarshal(out, &v); err != nil {
		return "", "", fmt.Errorf("gh repo view: unreadable answer: %v", err)
	}
	if v.NameWithOwner == "" {
		return "", "", fmt.Errorf("gh repo view: no repository named for %s", origin)
	}
	if v.Parent != nil && v.Parent.Owner.Login != "" && v.Parent.Name != "" {
		return v.Parent.Owner.Login + "/" + v.Parent.Name, v.NameWithOwner, nil
	}
	return v.NameWithOwner, v.NameWithOwner, nil
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
	if hs := n.Head.Nodes; len(hs) > 0 {
		p.HeadSHA = hs[0].Commit.OID
		if hs[0].Commit.StatusCheckRollup != nil {
			p.CI = githubCI(hs[0].Commit.StatusCheckRollup.State)
		}
	}
	if n.ReviewRequests.PageInfo.HasNextPage {
		p.partial = append(p.partial, "more than 100 reviews requested: not all read")
	}
	if n.LatestOpinionatedReviews.PageInfo.HasNextPage {
		p.partial = append(p.partial, "more than 100 reviewers with a verdict: not all read")
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
