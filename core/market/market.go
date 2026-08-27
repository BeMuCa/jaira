// Package market lists and fetches the lanes published in the jaira
// repository's own lanes/ directory on GitHub — the catalogue anyone can add
// to with a pull request, checked by CI before it lands.
//
// The directory is the marketplace. There is no registry, no index file and no
// format of its own: what GitHub's contents API says is in lanes/ is what is on
// offer, and a lane arrives here the same way a teammate's shared lane does —
// parsed first, then adopted into the catalogue on purpose.
package market

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/selfupdate"
)

// repo and dir name where the catalogue lives.
const (
	repo = "BeMuCa/jaira"
	dir  = "lanes"
)

// maxBody bounds every response read: a lane file is a few kilobytes, a
// directory listing a few more, and an unbounded read of a remote body is a
// memory bomb waiting on a misbehaving server.
const maxBody = 4 << 20

// apiBase is the contents-API address of the catalogue directory.
// JAIRA_MARKET_API overrides it, with the same rule as the release host
// overrides (https, or plain http on loopback for the test suite only) — see
// selfupdate.OverrideOf for why that rule and why not a flag.
func apiBase() string {
	if v := selfupdate.OverrideOf("JAIRA_MARKET_API"); v != "" {
		return v
	}
	return "https://api.github.com/repos/" + repo + "/contents/" + dir
}

// Overridden reports the override in effect, so a command can name it before
// installing a prompt that host served.
func Overridden() string {
	if v := selfupdate.OverrideOf("JAIRA_MARKET_API"); v != "" {
		return "JAIRA_MARKET_API=" + v
	}
	return ""
}

// Entry is one lane on offer: the parsed lane, and where its file lives.
type Entry struct {
	Lane *lane.Lane
	Path string // path inside the repository, e.g. lanes/critique.md
	URL  string // where the raw file is fetched from
	Raw  []byte // the file as served, parsed into Lane
}

// Client talks to the catalogue host.
type Client struct {
	HTTP *http.Client
}

// New returns a Client with an explicit timeout — the zero-value http.Client
// has none, and a hung TLS handshake would hang the command forever.
func New() *Client { return &Client{HTTP: &http.Client{Timeout: 30 * time.Second}} }

func (c *Client) client() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return New().HTTP
}

// List fetches the catalogue: every *.md in the directory except its README,
// each downloaded and parsed so the listing can show id, name and description
// rather than filenames. A file that does not parse is a warning, not a
// failure — one broken contribution must not hide the rest.
func (c *Client) List(ctx context.Context) ([]Entry, []string, error) {
	body, err := c.fetch(ctx, apiBase())
	if err != nil {
		return nil, nil, fmt.Errorf("listing the lane marketplace: %w", err)
	}
	// GitHub's contents payload has many fields; only these three are read,
	// so anything else can change shape without breaking this.
	var items []struct {
		Name        string `json:"name"`
		Path        string `json:"path"`
		Type        string `json:"type"`
		DownloadURL string `json:"download_url"`
	}
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, nil, fmt.Errorf("decoding the lane marketplace listing: %w", err)
	}
	var entries []Entry
	var warnings []string
	for _, it := range items {
		if it.Type != "file" || !strings.HasSuffix(it.Name, ".md") || strings.EqualFold(it.Name, "README.md") {
			continue
		}
		raw, err := c.fetch(ctx, it.DownloadURL)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("%s: could not fetch: %v", it.Path, err))
			continue
		}
		l, err := lane.Parse(raw, it.Path)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("%s did not parse and was skipped: %v", it.Path, err))
			continue
		}
		entries = append(entries, Entry{Lane: l, Path: it.Path, URL: it.DownloadURL, Raw: raw})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Lane.ID < entries[j].Lane.ID })
	return entries, warnings, nil
}

// errNotFound distinguishes a missing file from any other failure.
var errNotFound = errors.New("not found")

func (c *Client) fetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json, text/plain;q=0.9, */*;q=0.8")
	req.Header.Set("User-Agent", "jaira")
	resp, err := c.client().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, errNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: HTTP %d", url, resp.StatusCode)
	}
	return io.ReadAll(io.LimitReader(resp.Body, maxBody))
}
