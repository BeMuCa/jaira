package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A lane's --for-lane payload names the model its tier runs on, from settings:
// review is strong and runs on opus at the default, in-progress is cheap and
// runs on haiku at high effort. A settings.json entry replaces the tier's model.
func TestForLaneCarriesTheTierModel(t *testing.T) {
	lanesTestCatalogue(t)
	t.Setenv("JAIRA_DEFAULT_BOARD", filepath.Join(t.TempDir(), "none.md"))
	home := t.TempDir()
	t.Setenv("JAIRA_HOME", home)
	dir, h := newModeTicket(t)

	review := jsonCLI(t, dir, "show", h, "--for-lane", "review")
	if review["model_tier"] != "strong" || review["model"] != "opus" || review["effort"] != "" {
		t.Errorf("review --json: tier=%v model=%v effort=%v, want strong/opus/empty", review["model_tier"], review["model"], review["effort"])
	}
	progress := jsonCLI(t, dir, "show", h, "--for-lane", "in-progress")
	if progress["model"] != "haiku" || progress["effort"] != "high" {
		t.Errorf("in-progress --json: model=%v effort=%v, want haiku/high", progress["model"], progress["effort"])
	}

	if err := os.WriteFile(filepath.Join(home, "settings.json"), []byte(`{"models":{"strong":{"model":"sonnet","effort":"high"}}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	review = jsonCLI(t, dir, "show", h, "--for-lane", "review")
	if review["model"] != "sonnet" || review["effort"] != "high" {
		t.Errorf("review with settings models.strong: model=%v effort=%v, want sonnet/high", review["model"], review["effort"])
	}

	out, err := runCLI(t, dir, "show", h, "--for-lane", "review")
	if err != nil {
		t.Fatalf("show --for-lane review: %v\n%s", err, out)
	}
	if !strings.Contains(out, "tier: strong (sonnet high)") {
		t.Errorf("plain-text head does not name the model:\n%s", out)
	}
}
