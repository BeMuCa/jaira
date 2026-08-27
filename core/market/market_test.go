package market

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
