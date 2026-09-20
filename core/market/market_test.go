package market

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/release"
)

// serve stands in for GitHub: a contents listing of lanes/ and the files.
func serve(t *testing.T, files map[string]string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	var srv *httptest.Server
	mux.HandleFunc("/contents/lanes", func(w http.ResponseWriter, _ *http.Request) {
		var items []map[string]any
		for name := range files {
			items = append(items, map[string]any{
				"name": name, "path": "lanes/" + name, "type": "file",
				"download_url": srv.URL + "/raw/" + name,
			})
		}
		items = append(items, map[string]any{"name": "sub", "path": "lanes/sub", "type": "dir"})
		_ = json.NewEncoder(w).Encode(items)
	})
	mux.HandleFunc("/raw/", func(w http.ResponseWriter, r *http.Request) {
		body, ok := files[strings.TrimPrefix(r.URL.Path, "/raw/")]
		if !ok {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(body))
	})
	srv = httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	t.Setenv("JAIRA_MARKET_API", srv.URL+"/contents/lanes")
	return srv
}

const critique = "---\nid: critique\nname: Critique\nafter: in-progress\nprecedence: 45\nagentic: true\nmodel-tier: strong\ndescription: Judges the approach.\n---\n# Prompt\n\nCriticise.\n"
const broken = "---\nid: \nname: Broken\n---\n"

// TestListParsesEveryLaneAndSkipsReadmeAndBroken: the README is not a lane,
// a file that does not parse is a warning, and what remains is sorted by id.
func TestListParsesEveryLaneAndSkipsReadmeAndBroken(t *testing.T) {
	serve(t, map[string]string{
		"README.md":   "# Catalogue lanes\n",
		"critique.md": critique,
		"broken.md":   broken,
		"aaa.md":      strings.Replace(critique, "id: critique", "id: aaa", 1),
	})
	entries, warnings, err := New().List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, e := range entries {
		ids = append(ids, e.Lane.ID)
	}
	if got := strings.Join(ids, ","); got != "aaa,critique" {
		t.Errorf("ids = %s, want aaa,critique (sorted, README and broken left out)", got)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "broken.md") {
		t.Errorf("warnings = %v, want one naming broken.md", warnings)
	}
	if entries[1].Lane.Description != "Judges the approach." || entries[1].Path != "lanes/critique.md" {
		t.Errorf("entry = %+v", entries[1])
	}
}

// TestListWithoutAHostIsAnErrorNotAPanic: no server answers, the command gets
// an error it can print.
func TestListWithoutAHostIsAnErrorNotAPanic(t *testing.T) {
	t.Setenv("JAIRA_MARKET_API", "http://127.0.0.1:1/contents/lanes")
	if _, _, err := New().List(context.Background()); err == nil {
		t.Fatal("expected an error with nothing listening")
	}
}

// TestOverrideRefusesCleartextOffLoopback: the same rule as the release
// host — https, or http on loopback for tests, nothing else.
func TestOverrideRefusesCleartextOffLoopback(t *testing.T) {
	t.Setenv("JAIRA_MARKET_API", "http://example.com/contents/lanes")
	if got := apiBase(); !strings.HasPrefix(got, "https://api.github.com/") {
		t.Errorf("apiBase = %q, want the default when the override is cleartext off loopback", got)
	}
	if Overridden() != "" {
		t.Error("a refused override must not be reported as in effect")
	}
}

// TestFetchReturnsTheFileAsServed: what adopt copies is byte-for-byte what
// the repository holds.
func TestFetchReturnsTheFileAsServed(t *testing.T) {
	serve(t, map[string]string{"critique.md": critique})
	c := New()
	entries, _, err := c.List(context.Background())
	if err != nil || len(entries) != 1 {
		t.Fatalf("List = %v, %v", entries, err)
	}
	raw, err := c.Fetch(context.Background(), entries[0])
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != critique {
		t.Errorf("Fetch = %q, want the served file", raw)
	}
}

// TestListPinsTheCatalogueToTheRunningTag is DoD 6: a released binary must be
// served the catalogue of its own tag, not the default branch's HEAD. The
// server records the ref it was asked for.
func TestListPinsTheCatalogueToTheRunningTag(t *testing.T) {
	var gotRef string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRef = r.URL.Query().Get("ref")
		_, _ = w.Write([]byte("[]"))
	}))
	t.Cleanup(srv.Close)
	t.Setenv("JAIRA_MARKET_API", srv.URL+"/contents/lanes")

	old := release.Current
	t.Cleanup(func() { release.Current = old })
	release.Current = "0.3.0"

	if _, _, err := New().List(context.Background()); err != nil {
		t.Fatalf("List: %v", err)
	}
	if gotRef != "v0.3.0" {
		t.Errorf("ref = %q, want %q — an old binary served today's lanes is the bug this pins shut", gotRef, "v0.3.0")
	}
	if Unpinned() != "" {
		t.Errorf("a released build must not announce itself as unpinned, got %q", Unpinned())
	}
}

// TestDevBuildSendsNoRefAndSaysSo asserts the other half of DoD 6: a source
// build has no tag to ask for, so it fetches the development branch — and
// says so rather than letting the user believe otherwise. The sentence is
// returned, not printed: apiBase runs on every request and --json has no room
// for prose.
func TestDevBuildSendsNoRefAndSaysSo(t *testing.T) {
	old := release.Current
	t.Cleanup(func() { release.Current = old })
	release.Current = "dev"

	var asked bool
	var gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		asked, gotQuery = true, r.URL.RawQuery
		_, _ = w.Write([]byte("[]"))
	}))
	t.Cleanup(srv.Close)
	t.Setenv("JAIRA_MARKET_API", srv.URL+"/contents/lanes")

	if _, _, err := New().List(context.Background()); err != nil {
		t.Fatalf("List: %v", err)
	}
	if !asked {
		t.Fatal("the catalogue was never asked")
	}
	if strings.Contains(gotQuery, "ref=") {
		t.Errorf("a dev build must send no ref, got query %q", gotQuery)
	}
	u := Unpinned()
	if u == "" {
		t.Fatal("a dev build must say the catalogue is the development branch")
	}
	if !strings.Contains(u, "dev") {
		t.Errorf("the announcement must name the version it reports: %q", u)
	}
}

// TestRefIsSetOnAnAddressThatAlreadyHasAQuery is the concatenation bug this
// must not have: gluing "?ref=" onto an override that already carries one
// produces "?x=1?ref=v0.3.0", a single nonsense parameter. The test server
// sets a query of its own, so this is the ordinary case, not a corner.
func TestRefIsSetOnAnAddressThatAlreadyHasAQuery(t *testing.T) {
	old := release.Current
	t.Cleanup(func() { release.Current = old })
	release.Current = "0.3.0"

	t.Setenv("JAIRA_MARKET_API", "https://example.test/contents/lanes?token=abc")
	got := apiBase()
	u, err := url.Parse(got)
	if err != nil {
		t.Fatalf("apiBase produced an unparseable address %q: %v", got, err)
	}
	if u.Query().Get("ref") != "v0.3.0" {
		t.Errorf("ref = %q in %q", u.Query().Get("ref"), got)
	}
	if u.Query().Get("token") != "abc" {
		t.Errorf("the address' own query was lost: %q", got)
	}
}

// TestBlankVersionStillReadsAsASentence: pinnedRef already treats a blanked
// release.Current as unpinned; Unpinned used to paste the empty string
// straight into its sentence and produce "reports version , which is no
// released tag".
func TestBlankVersionStillReadsAsASentence(t *testing.T) {
	old := release.Current
	t.Cleanup(func() { release.Current = old })
	release.Current = ""

	u := Unpinned()
	if u == "" {
		t.Fatal("a build with no version must still say the catalogue is the development branch")
	}
	if strings.Contains(u, "version ,") {
		t.Errorf("the sentence has a hole where the version would be: %q", u)
	}
	if !strings.Contains(u, "no version") {
		t.Errorf("the sentence must say there is no version to name: %q", u)
	}
}
