package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/role"
	"github.com/BeMuCa/jaira/core/settings"
)

// nudgeSkillsIfStale prints one line on stderr when the jaira skills in a personal skills directory
// are older than this binary: a file 'jaira update --roles' would write or replace. Edited files are
// the person's and never count. It checks once per binary version — the stamp is the version in
// <jaira home>/skills-check/version — so every later command pays one small file read.
//
// Like nudgeIfStale it writes to os.Stderr directly: stdout is reserved for the payload an agent
// parses, and the line must reach the terminal under --json too.
func nudgeSkillsIfStale() {
	if release.Current == "dev" || updateInProgress {
		return
	}
	// No home means no jaira home: filepath.Dir("") is "." and would stamp into the cwd.
	p := settings.Path()
	if p == "" {
		return
	}
	stamp := filepath.Join(filepath.Dir(p), "skills-check")
	if release.Stamped(stamp) == release.Current {
		return
	}
	runs, err := refreshRoles("", true)
	if err != nil {
		return
	}
	var behind []string
	for _, r := range runs {
		if role.Behind(r.results) > 0 {
			behind = append(behind, r.target.Dir)
		}
	}
	// Stamped even when nothing is behind: one check per version, not one per command.
	_ = release.Stamp(stamp)
	if len(behind) == 0 {
		return
	}
	fmt.Fprintf(os.Stderr, "jaira: the jaira skills in %s are older than this jaira (%s) — run 'jaira update --roles' to bring them up to date\n",
		strings.Join(behind, " and "), release.Current)
}
