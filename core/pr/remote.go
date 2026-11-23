package pr

import "strings"

// remote is a git remote URL taken apart: the host it lives on, without a
// port or a user, and the repository's path there, owner/repo on GitHub,
// group/…/project on GitLab, without .git.
type remote struct {
	Host string
	Repo string
}

// parseRemote takes apart the URL forms git accepts for a remote:
// scheme://[user@]host[:port]/owner/repo(.git) for https, ssh, git and
// git+ssh, and the scp form [user@]host:owner/repo(.git), whatever the user
// before the @ is called. A GitHub Enterprise remote often reads
// alice@ghe.example.com:team/app.git, which gh cannot resolve as it stands
// (XSGWPX). false where the URL names no host or no owner/repo — a local
// path, for one.
func parseRemote(raw string) (remote, bool) {
	s := strings.TrimSpace(raw)
	var host, path string
	if i := strings.Index(s, "://"); i > 0 {
		rest := s[i+3:]
		j := strings.IndexByte(rest, '/')
		if j < 0 {
			return remote{}, false
		}
		host, path = rest[:j], rest[j+1:]
		if k := strings.LastIndexByte(host, '@'); k >= 0 {
			host = host[k+1:]
		}
		if k := strings.LastIndexByte(host, ':'); k >= 0 && !strings.HasSuffix(host, "]") {
			host = host[:k]
		}
	} else {
		// scp form: no slash before the first colon, which a Windows drive
		// (C:\…) or a local path never satisfies with a host in front.
		colon := strings.IndexByte(s, ':')
		if colon <= 0 || strings.ContainsAny(s[:colon], `/\`) {
			return remote{}, false
		}
		host, path = s[:colon], s[colon+1:]
		if k := strings.LastIndexByte(host, '@'); k >= 0 {
			host = host[k+1:]
		}
		if len(host) == 1 {
			return remote{}, false // a drive letter, not a host
		}
	}
	path = strings.Trim(path, "/")
	path = strings.TrimSuffix(path, ".git")
	if host == "" || !strings.Contains(path, "/") {
		return remote{}, false
	}
	return remote{Host: strings.ToLower(host), Repo: path}, true
}

// ghRepo is how gh is told the repository of origin: on a host other than
// github.com as host/owner/repo — gh's own [HOST/]OWNER/REPO, which -R and
// 'gh repo view' both take — with the host for --hostname. github.com, and
// a URL that cannot be taken apart, go as they are: gh reads those itself.
func ghRepo(origin string) (repo, host string) {
	r, ok := parseRemote(origin)
	if !ok || r.Host == "github.com" || r.Host == "www.github.com" {
		return origin, ""
	}
	return r.Host + "/" + r.Repo, r.Host
}
