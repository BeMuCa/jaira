package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const marketCritique = "---\nid: critique\nname: Critique\nafter: in-progress\nprecedence: 45\nagentic: true\nmodel-tier: strong\ndescription: Judges the approach.\n---\n# Prompt\n\nCriticise.\n"

// marketServer stands in for GitHub's contents API and raw files.
func marketServer(t *testing.T) {
	t.Helper()
	mux := http.NewServeMux()
	var srv *httptest.Server
	mux.HandleFunc("/contents/lanes", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode([]map[string]any{
			{"name": "README.md", "path": "lanes/README.md", "type": "file", "download_url": srv.URL + "/raw/README.md"},
			{"name": "critique.md", "path": "lanes/critique.md", "type": "file", "download_url": srv.URL + "/raw/critique.md"},
		})
	})
	mux.HandleFunc("/raw/critique.md", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(marketCritique)) })
	mux.HandleFunc("/raw/README.md", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("# Catalogue\n")) })
	srv = httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	t.Setenv("JAIRA_MARKET_API", srv.URL+"/contents/lanes")
}

// TestLanesMarketListsWhatTheRepositoryOffers: id, name and description from
// the parsed files; the README is not a lane; the override is named.
func TestLanesMarketListsWhatTheRepositoryOffers(t *testing.T) {
	lanesTestCatalogue(t)
	root := lanesTestProject(t)
	marketServer(t)

	out, err := runLanes(t, root, "market")
	if err != nil {
		t.Fatalf("lanes market: %v\n%s", err, out)
	}
	for _, want := range []string{"critique", "Judges the approach.", "JAIRA_MARKET_API", "market adopt"} {
		if !strings.Contains(out, want) {
			t.Errorf("listing lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "README") {
		t.Errorf("the catalogue README is not a lane:\n%s", out)
	}

	jout, err := runLanes(t, root, "market", "--json")
	if err != nil {
		t.Fatalf("lanes market --json: %v\n%s", err, jout)
	}
	var payload struct {
		Lanes []map[string]any `json:"lanes"`
	}
	if err := json.Unmarshal([]byte(jout), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, jout)
	}
	if len(payload.Lanes) != 1 || payload.Lanes[0]["id"] != "critique" || payload.Lanes[0]["model_tier"] != "strong" {
		t.Errorf("json lanes = %v", payload.Lanes)
	}
}

// TestLanesMarketAdoptLandsInTheCatalogueAndAddFindsIt is the whole point:
// two commands from "saw it on GitHub" to "on my board".
func TestLanesMarketAdoptLandsInTheCatalogueAndAddFindsIt(t *testing.T) {
	catalogue := lanesTestCatalogue(t)
	root := lanesTestProject(t)
	marketServer(t)

	out, err := runLanes(t, root, "market", "adopt", "critique")
	if err != nil {
		t.Fatalf("market adopt: %v\n%s", err, out)
	}
	got, err := os.ReadFile(filepath.Join(catalogue, "critique.md"))
	if err != nil {
		t.Fatalf("expected the lane in the catalogue: %v", err)
	}
	if string(got) != marketCritique {
		t.Errorf("catalogue copy differs from what was served:\n%s", got)
	}
	if !strings.Contains(out, "lanes add critique") {
		t.Errorf("adopt output does not say how to put it on the board:\n%s", out)
	}

	if _, err := runLanes(t, root, "market", "adopt", "critique"); err == nil {
		t.Error("a second adopt without --force must refuse")
	}
	if out, err := runLanes(t, root, "add", "critique"); err != nil {
		t.Fatalf("lanes add after market adopt: %v\n%s", err, out)
	}
}

// TestLanesMarketUnknownIdNamesWhatExists: a typo gets the list, exit 2.
func TestLanesMarketUnknownIdNamesWhatExists(t *testing.T) {
	lanesTestCatalogue(t)
	root := lanesTestProject(t)
	marketServer(t)
	out, err := runLanes(t, root, "market", "adopt", "critiqe")
	if err == nil {
		t.Fatalf("expected a refusal, got:\n%s", out)
	}
	if !strings.Contains(err.Error(), "critique") {
		t.Errorf("refusal does not name what exists: %v", err)
	}
}

// TestLanesMarketWithoutNetworkFailsPlainly: nothing listening is an error
// message with exit 1, not a panic and not a partial write.
func TestLanesMarketWithoutNetworkFailsPlainly(t *testing.T) {
	catalogue := lanesTestCatalogue(t)
	root := lanesTestProject(t)
	t.Setenv("JAIRA_MARKET_API", "http://127.0.0.1:1/contents/lanes")
	if _, err := runLanes(t, root, "market"); err == nil {
		t.Fatal("expected an error with nothing listening")
	}
	if _, err := runLanes(t, root, "market", "adopt", "critique"); err == nil {
		t.Fatal("expected adopt to fail with nothing listening")
	}
	if entries, _ := os.ReadDir(catalogue); len(entries) != 0 {
		t.Errorf("a failed adopt must write nothing, catalogue has %d entries", len(entries))
	}
}
