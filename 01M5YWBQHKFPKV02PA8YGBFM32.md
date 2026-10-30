---
id: 01M5YWBQHKFPKV02PA8YGBFM32
title: "jaira update bringt standardmaessig alles auf Stand - Abschnitt, Rollen, Cockpit, Skill - und 'jaira update diff' zeigt, was zurueckliegt"
status: in-progress
ready: true
creator: berkc
assignee: berkc
goal: "Ein einziges 'jaira update' bringt Board-Abschnitt, Rollen-Skills, Cockpit-Mod und jaira-Skill auf den Stand des installierten Binaries; 'jaira update diff' sagt vorher, was zurueckliegt."
context: |-
  Wunsch von berkc am 27.10.2026: 'Jaira update should do all of it by default, from skill to mod to roles. jaira update diff can be used to see which lanes or roles are behind.'
  Heute (0.5.6): 'jaira update' schreibt nur den jaira-Abschnitt in CLAUDE.md/AGENTS.md und die Release-Notes. Rollen und Cockpit kommen nur ueber 'jaira roles install --global|--project'.
  PR #91 (T8A8KM) fuegt 'update --roles', '--lanes' und '--dry-run' hinzu - Rollen sind dort noch opt-in.
  Offen: Lanes auch per Default? Lanes sind geteilte Board-Konfiguration und brauchen einen Commit - Vorschlag: Rollen/Mod/Skill per Default, Lanes nur mit --lanes, 'diff' zeigt beides.
  Gefunden am 23.10.: 'roles install' hielt zwei alte, unveraenderte Cockpit-Dateien (board.ts aus 5bc22fb1, register.tsx aus 2c4c6790) fuer 'von dir bearbeitet', weil ihr Hash in core/shipped/shipped.sha256 fehlt - auch in #91. Vermutung (ungeprueft): scripts/shipped-history.sh geht nicht durch zweite Merge-Eltern. Ohne den Fix wuerde ein Default-Update genau diese Dateien stehen lassen.
definition-of-done: "'jaira update' ohne Flags aktualisiert Abschnitt, Rollen, Cockpit-Mod und jaira-Skill ueberall, wo sie schon installiert sind; bearbeitete Dateien bleiben (exit 3)"
tags:
  - cli
blocked-by: []
related:
  - 01M21PJTWPK4DZ89PRY6T8A8KM
commits: []
created-at: 2026-10-27T14:31:21Z
updated-at: 2026-10-30T22:13:02Z
updated-by: berkc
claimed-by: berk-ca-2129445
claimed-at: 2026-10-30T21:41:03Z
---

# jaira update bringt standardmaessig alles auf Stand - Abschnitt, Rollen, Cockpit, Skill - und 'jaira update diff' zeigt, was zurueckliegt

## Definition of Done

- [ ] 'jaira update' ohne Flags aktualisiert Abschnitt, Rollen, Cockpit-Mod und jaira-Skill ueberall, wo sie schon installiert sind; bearbeitete Dateien bleiben (exit 3)
- [ ] 'jaira update diff' listet je Rolle, Mod-Datei und Lane: aktuell, zurueck (welche Version) oder bearbeitet; schreibt nichts
- [ ] Alte ausgelieferte Versionen aus Merge-Zweigen (board.ts 5bc22fb1, register.tsx 2c4c6790) gelten als ausgeliefert, nicht als bearbeitet
- [ ] Tests; Release-Note

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] find the real cause of the false 'edited' cockpit files in core/shipped (done: list is stale, script is correct; see notes)
  proof: core/shipped/shipped.sha256 regenerated; scripts/shipped-history.sh
- [x] scripts/shipped-history.sh: add the release column (first vX.Y.Z tag that holds the blob at that path, '-' if none)
- [x] regenerate core/shipped/shipped.sha256 with the script; check append-only (no committed line lost) and every tagged embedded file present
- [x] core/shipped/shipped.go: Known keeps sha+path matching; add Release(path, b) to read the column; update shipped_test.go regex
- [x] core/role/install.go and core/lane/refresh.go: Result / Refreshed carry the release of an Updated file
- [x] internal/cli/update.go + updatesetup.go: roles default on, --roles accepted as no-op, personal-only path kept, hintLines keeps lanes only
- [x] internal/cli/update.go: 'jaira update diff' subcommand, one line per role, mod file and lane: current | behind (release) | edited; writes nothing
- [x] tests: adjust update_test.go / updatedry_test.go / updatesetup_test.go for new defaults; add diff tests incl. a no-write check
- [ ] release note core/release/notes/GBFM32.md (one '- ' line per change, instruction form)
- [ ] go vet ./... and go test ./... (exit codes checked, output verbatim on failure)

## Progress
- **2026-10-30 21:41 · berkc** — open decisions at start: none that only the person can take. Settled by the person (ticket context, 27.10.): roles, cockpit mod and jaira skill update by default; lanes opt-in via --lanes; 'update diff' writes nothing. Settled here by default: diff output is one line per role/mod file/lane with state current|behind (version)|edited, the same for all three; 'update' touches only what is already installed. Parent session: none (started directly as worker subagent, no AskUserQuestion).
- **2026-10-30 21:44 · berkc** — finding 1 (shipped list): scripts/shipped-history.sh does walk second parents: it uses 'git log --full-history -m' since its first commit 299c65aa, and every commit reachable from HEAD that touches the embedded paths is in that walk (303 vs 375 commits, none missing). Its output today is 558 lines; core/shipped/shipped.sha256 has 503 and is a strict subset of it (0 lines committed-only, 55 lines script-only). The list is append-only and was patched by hand (736c80ee) and never re-run after the 44 later commits that changed embedded files. Hence board.ts 9406af (added by 5bc22fb1, F80QZT) and register.tsx are reported as edited. Fix: regenerate the list with the script. Checked: every file of tags v0.4.1..v0.5.7 under core/role/builtin, core/mod/builtin, core/lane/builtin, core/role/jaira has its sha in the regenerated list (0 missing); the committed list misses 22 of them (v0.5.2..v0.5.6 cockpit files).
- **2026-10-30 21:46 · berkc** — finding 2 (cause of the cockpit false 'edited'): not the script's walk. The committed list is an append-only file: script output at 299c65aa (438 lines) was extended by hand (736c80ee) and never re-run after 44 later commits to embedded paths. Committed list = strict subset of script output today (0 lines lost, 55 missing). Fix = regenerate, then keep append-only. Also a manifest that keeps 2 columns cannot say which release a behind file came from, so DoD 2 'welche Version' needs a release column.
- **2026-10-30 21:46 · berkc** — decisions settled here (reversible, not person-only): (a) manifest line becomes 'sha  path  release' where release is the first vX.Y.Z tag whose tree holds that blob at that path, '-' if none; shipped.Known keeps matching on sha+path. (b) update: roles, cockpit mod and jaira skill default on; --roles stays accepted as a no-op; lanes stay --lanes. (c) 'update diff' is a subcommand, one line per role/mod file and per lane: current | behind (release, or 'missing' for files not yet written) | edited; it reuses the dry-run of refreshRoles and lane.RefreshShipped and writes nothing. Outside a board, roles only, as before.
- **2026-10-30 21:47 · berkc** — plan written as Plan items (see jaira show). Options 'planning' was left unticked by the person, so the pre-process lane is skipped; the plan above is the pre-process output, written by the worker.
- **2026-10-30 21:56 · berkc** — built and checked by hand on a scratch clone (HOME=scratch, stale v0.5.4 cockpit files): update diff lists board.ts behind (shipped in v0.5.4), texts.ts/types behind (v0.5.2); bare update rewrites them to current; diff afterwards shows only lanes (critique edited, optimize/testing behind: those files in the clone are old). update diff writes nothing (tests: TestUpdateDiffWritesNothingAndNamesWhatIsBehind, tree snapshot). Lane Load writes lane files when a board has none: diff skips lanes there on purpose (ProjectLanesActive guard).
