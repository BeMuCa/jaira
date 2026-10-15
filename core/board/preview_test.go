package board

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func writeAgentFile(t *testing.T, root, name, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// snapshot lists every file under root with its bytes.
func snapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		b, err := os.ReadFile(filepath.Join(root, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		out[e.Name()] = string(b)
	}
	return out
}

// The preview says what the write does, word for word, and touches nothing.
func TestPreviewAgentFilesSaysWhatTheWriteDoesAndWritesNothing(t *testing.T) {
	stale := "# mine\n\n" + jairaMarkerStart + "\nstale\n" + jairaMarkerEnd + "\n"
	blockOnly := jairaMarkerStart + "\nstale\n" + jairaMarkerEnd + "\n"
	cases := []struct {
		name, choice string
		files        map[string]string
	}{
		{"stale block", "", map[string]string{"CLAUDE.md": stale}},
		{"fresh root", "", nil},
		{"agents, block removed and import added", AgentFileAgents, map[string]string{"CLAUDE.md": stale}},
		{"agents, claude file deleted", AgentFileAgents, map[string]string{"CLAUDE.md": blockOnly}},
		{"claude", AgentFileClaude, map[string]string{"AGENTS.md": stale}},
		{"both", AgentFileBoth, map[string]string{"AGENTS.md": stale}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root := t.TempDir()
			for name, body := range c.files {
				writeAgentFile(t, root, name, body)
			}
			before := snapshot(t, root)

			preview, err := PreviewAgentFiles(root, c.choice, nil)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(snapshot(t, root), before) {
				t.Fatal("the preview wrote")
			}

			var written []string
			if c.choice == "" {
				written, err = AnnounceInAgentFiles(root, nil)
			} else {
				written, err = ChooseAgentFiles(root, c.choice, nil)
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(written) == 0 {
				t.Fatal("the case changes nothing; it proves nothing")
			}
			if !reflect.DeepEqual(preview, written) {
				t.Errorf("preview %q, write %q", preview, written)
			}
		})
	}
}
