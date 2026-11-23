package pr

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const stateQuery = `query($owner:String!,$name:String!,$number:Int!){
 repository(owner:$owner,name:$name){pullRequest(number:$number){state}}}`

// State asks the forge of the repository in dir whether pull request number
// is open, merged or closed — nothing else, so a blocked-by entry pr#<n>
// costs two small calls, not the whole stand Read collects. The forge and
// the repository are settled as Read settles them: git config jaira.forge or
// origin's host, and a fork's parent.
func State(run Run, look func(string) (string, error), dir string, number int) (string, error) {
	forge, origin, call, err := connect(run, look, dir)
	if err != nil {
		return "", err
	}
	if forge == "github" {
		return githubState(call, origin, number)
	}
	return gitlabState(call, origin, number)
}

func githubState(call func(...string) ([]byte, error), origin string, number int) (string, error) {
	repo, _, err := githubTarget(call, origin)
	if err != nil {
		return "", err
	}
	owner, name, ok := ownerName(repo)
	if !ok {
		return "", fmt.Errorf("gh repo view: %q is not owner/name", repo)
	}
	out, err := call("api", "graphql", "-f", "owner="+owner, "-f", "name="+name,
		"-F", "number="+strconv.Itoa(number), "-f", "query="+stateQuery)
	if err != nil {
		return "", err
	}
	var a struct {
		Data *struct {
			Repository *struct {
				PullRequest *struct {
					State string `json:"state"`
				} `json:"pullRequest"`
			} `json:"repository"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(out, &a); err != nil {
		return "", fmt.Errorf("gh api graphql: unreadable answer: %v", err)
	}
	if len(a.Errors) > 0 {
		return "", fmt.Errorf("gh api graphql: %s", a.Errors[0].Message)
	}
	if a.Data == nil || a.Data.Repository == nil || a.Data.Repository.PullRequest == nil {
		return "", fmt.Errorf("gh api graphql: %s has no pull request #%d", repo, number)
	}
	return knownState(strings.ToLower(a.Data.Repository.PullRequest.State), "gh")
}

func gitlabState(call func(...string) ([]byte, error), origin string, number int) (string, error) {
	repo, _, err := gitlabTarget(call, origin)
	if err != nil {
		return "", err
	}
	path := "projects/" + url.PathEscape(repo) + "/merge_requests/" + strconv.Itoa(number)
	out, err := call("api", path)
	if err != nil {
		return "", err
	}
	var m struct {
		State string `json:"state"`
	}
	if err := json.Unmarshal(out, &m); err != nil {
		return "", fmt.Errorf("glab api %s: unreadable answer: %v", path, err)
	}
	state := m.State
	switch state {
	case "opened", "locked":
		state = "open"
	}
	return knownState(state, "glab")
}

// knownState keeps only the three states a caller acts on; anything else
// is an answer jaira cannot read, not a fourth state.
func knownState(s, tool string) (string, error) {
	switch s {
	case "open", "merged", "closed":
		return s, nil
	}
	return "", fmt.Errorf("%s: pull request state %q is not open, merged or closed", tool, s)
}
