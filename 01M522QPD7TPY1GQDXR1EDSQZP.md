---
id: 01M522QPD7TPY1GQDXR1EDSQZP
title: "Cockpit-Band: Update-Hinweis mit Knopf 'Update' und Knopf 'Verwerfen'"
status: review
ready: true
creator: BeMuCa
goal: "Wenn eine neue jaira-Version da ist, zeigt das Cockpit-Band einen Hinweis mit zwei Knoepfen: Update (fuehrt den Ja-Pfad aus) und Verwerfen (blendet den Hinweis fuer diese Version aus)"
context: "Berk wollte: im Cockpit einen Update-Hinweis mit zwei Knoepfen. Das Band zeigt heute keinen Update-Hinweis. T8A8KM hat ihn verworfen, weil das Cockpit nur mit installiertem Cockpit laeuft und das TS nicht getestet war. Das ist kein Grund gegen den Knopf, nur einer fuer den Weg. Die stderr-Zeile in next/resume/claim/move erreicht nur Sitzungen, die gerade einen jaira-Befehl fahren. Der Knopf erreicht jede Sitzung, die das Cockpit offen hat. Ja-Pfad: erst 'jaira self upgrade --check', dann eine Vorschau 'jaira update --roles --lanes --dry-run', erst nach Ja das echte Update. Verwerfen merkt sich die Version, damit die Zeile nicht wiederkommt."
definition-of-done: Band zeigt bei neuerer Version Hinweis mit Knopf Update und Knopf Verwerfen; Update fuehrt Vorschau und danach Ja-Pfad aus; Verwerfen unterdrueckt den Hinweis fuer genau diese Version; Test mit Node-Lauf im Cockpit
tags:
  - release
  - tui
blocked-by: []
related: []
follows: 01M21PJTWPK4DZ89PRY6T8A8KM
commits: []
created-at: 2026-10-16T10:04:46Z
updated-at: 2026-10-16T13:00:01Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-682638
claimed-at: 2026-10-16T12:59:08Z
outcome-what: "Runde 1: jaira-cockpit-Band zeigt bei neuerer Version (Release-Cache $JAIRA_HOME|~/.jaira/update-check.json gegen 'jaira --version', jede Minute gelesen) 'jaira X is released (this is Y)' mit [ Update ] und Dismiss. Update: 'jaira self upgrade --check --json' + 'jaira update --roles --lanes --dry-run', Vorschau im Pane cockpit-upgrade (Cancel hat den Fokus), erst 'Update to X' fuehrt 'jaira self upgrade' und 'jaira update --roles --lanes' aus. Dismiss merkt die Version im Store. Runde 2: Ergebnis-Pane zeigt auch stderr von self upgrade (Host-Override-Warnung); ein Dry-run mit Exit ausser 0/3 oeffnet das Fehler-Pane statt der Vorschau und bietet kein Update an. Agenten-Weg: stderr-Zeile (internal/cli/newerrelease.go) und Agenten-Block (core/board/announce.go, CLAUDE.md, AGENTS.md) sagen jetzt: erst --check und --dry-run laufen lassen, dem User zeigen, fragen, erst auf Ja die echten Befehle. README, docs/COMMANDS.md, zwei Release-Notiz-Zeilen, shipped.sha256 neu."
outcome-why: "Die stderr-Zeile erreicht nur Sitzungen, die gerade jaira-Befehle fahren; das Band erreicht jede Sitzung mit offenem Cockpit. Berk wollte den Knopf mit Vorschau vor jedem Schreiben und einem Verwerfen, das sich die Version merkt - und dieselbe Vorschau auch auf dem Agenten-Weg ('kein Dry run noetig?')."
outcome-resolves: "DoD 1: Hinweis mit Update und Verwerfen (Dismiss), Update = Vorschau dann Ja-Pfad, Verwerfen nur fuer genau diese Version; claude plugin test 157 pass. DoD 2: stderr-Zeile und Agenten-Block verlangen die Vorschau vor der Frage; Go-Tests auf die eine Zeile und den Block."
review-summary: "Two commits. 9d2cfd51 (cockpit mod): once a minute the band reads jaira's own release cache (update-check.json under JAIRA_HOME, else ~/.jaira) and compares its 'latest' with 'jaira --version'. When the cache names a newer release that was not dismissed, the band shows 'jaira X is released (this is Y)' with [ Update ] and Dismiss. Update runs 'jaira self upgrade --check --json' and 'jaira update --roles --lanes --dry-run' and shows both in a pane where Cancel has the focus. Only 'Update to X' in that pane runs 'jaira self upgrade' and then 'jaira update --roles --lanes', and shows their output. Dismiss keeps the version in the mod's store, so exactly that release stays hidden across sessions and the next one shows again. With no board or on a dev build there is no notice. efa8e9cd: the result pane now shows self upgrade's stderr too. A dry run that exits with anything other than 0 or 3 opens the error pane and offers no update (exit 3 means edited files are left alone, and it still counts as a preview). The stderr line after next/resume/claim/move and the agent block in CLAUDE.md/AGENTS.md now tell an agent: run --check and --dry-run first, show the user, ask, and run the real two commands only on a yes. Also changed: README, docs/COMMANDS.md, two release-note lines, and shipped.sha256 regenerated."
review-gaps: "No blocking gap. Minor: (1) After the yes, the result pane is titled 'jaira updated to X' whatever 'jaira update --roles --lanes' exited with (register.tsx upgradeJaira). Its output is in the pane, but a failed setup step only shows if you read that output. (2) The agent block and the stderr line say the two previews write nothing, but 'self upgrade --check' does rewrite update-check.json under ~/.jaira (or JAIRA_HOME). Nothing in the repo or the binary changes, so this is harmless wording. (3) Already on the ticket as follow-up candidates, not new: the 20s press lock (BUSY_MS) against a 300s upgrade lets a second press during a slow download start a second preview and a second upgrade; a brew or go-install build gets offered 'Update to X', which can only fail. Neither one is filed as a ticket yet. (4) Round 1's review-check said 155 pass and 12 tests; the real numbers are 157 and 14, and review-check is rewritten."
test-verdict: "pass (round 2, efa8e9cd): go build ok, go test -race -count=1 ./... RC=0 (32 ok, 0 FAIL); claude plugin validate passed, claude plugin test 157 pass 0 fail; DoD 1+2 proof lines checked in the tree; scratch with real binaries (0.5.1 -> stub release 0.6.0, scratch HOME/JAIRA_HOME): jaira next prints the new one-line instruction, the agent's two previews (self upgrade --check, update --roles --lanes --dry-run) left the snapshot unchanged; round 1 yes-path run stands; not seen: the band live in a Claude Code terminal"
review-verdict: "The diff meets both definition-of-done items, and I found no defect that should block it. Checked by running, not only reading: (a) claude plugin test: 157 pass, 0 fail, 14 tests in 'a newer jaira release'. (b) This branch built as 0.5.1 into scratch, with a cache naming 9.9.9: 'jaira next --json' prints one stderr line naming check, dry-run, show, ask and only-on-yes, and stdout is clean. (c) The committed CLAUDE.md and AGENTS.md agent blocks are byte-identical to what that build writes. (d) 'update --roles --lanes --dry-run' wrote nothing (before/after snapshot equal). (e) shipped.sha256 regenerated in scratch is identical to the committed file. (f) gofmt and vet are clean. (g) node 24 on board.ts: offeredRelease (newer, dismissed, next after a dismissed one, same, dev, 0.10.0 vs 0.9.9) and upgradeCheckOf behave as claimed. Not checked by me: the cockpit in a live Claude Code session, and the Go test suite (not run, by the owner's rule). The gaps are small: (1) and (2) are cosmetic, and (3) is already known. Merge after T8A8KM and PWAF22: the branch stacks on them."
review-round: 2
review-check: |-
  1. cd /home/berk/git/jAIra-EDSQZP/core/mod/builtin/jaira-cockpit && claude plugin test  -> it ends with '157 pass' and '0 fail', and 14 lines start with '(pass) a newer jaira release >'.
  2. Agent line, in scratch only, nothing installed is replaced: cd /home/berk/git/jAIra-EDSQZP && go build -ldflags '-X main.version=0.5.1' -o /tmp/jh/bin/jaira ./cmd/jaira   (always -o into /tmp, never ~/.local/bin)
  3. mkdir -p /tmp/jh/home/jh /tmp/jh/repo && echo "{\"checked_at\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",\"latest\":\"9.9.9\"}" > /tmp/jh/home/jh/update-check.json   (a fresh checked_at stops a real release check from overwriting 9.9.9)
  4. cd /tmp/jh/repo && git init -q && git commit -q --allow-empty -m init
  5. export HOME=/tmp/jh/home JAIRA_HOME=/tmp/jh/home/jh JAIRA_RELEASE_API=http://127.0.0.1:9/x && /tmp/jh/bin/jaira init && /tmp/jh/bin/jaira next --json >/dev/null  -> exactly one line on screen: "jaira: jaira 9.9.9 is released (this is 0.5.1) - run 'jaira self upgrade --check' and 'jaira update --roles --lanes --dry-run', show the user what they report and ask the user whether to update; only on yes run 'jaira self upgrade', then 'jaira update --roles --lanes'"
  6. grep -n 'newer release' -A5 /tmp/jh/repo/CLAUDE.md  -> the same six lines as grep -n 'newer release' -A5 /home/berk/git/jAIra-EDSQZP/CLAUDE.md (check, dry-run, show, ask once, only on a yes).
  7. Cockpit, live (not run by me in this review): in a NEW terminal (without the exports from step 5): mkdir -p /tmp/jc && echo '{"checked_at":"2026-10-16T12:00:00Z","latest":"9.9.9"}' > /tmp/jc/update-check.json && cd /home/berk/git/jAIra-EDSQZP && JAIRA_HOME=/tmp/jc claude --plugin-dir core/mod/builtin/jaira-cockpit  -> the jaira-cockpit frame shows 'jaira 9.9.9 is released (this is 0.5.1) [ Update ] Dismiss'. If your installed cockpit also loads, use the frame that has this line.
  8. Press Dismiss -> the line goes, and the notice says 'jaira 9.9.9 dismissed: the band offers the next release again'. Change 9.9.9 to 9.9.10 in /tmp/jc/update-check.json and wait up to one minute -> the line comes back with 9.9.10.
  9. Press Update -> while 0.5.1 is still the newest real release, no pane opens, the notice says 'jaira 0.5.1 is already the latest release', and the line goes. Run sha256sum ~/.local/bin/jaira before and after: the same hash both times.
  10. The preview pane, the dry-run error pane and 'Update to X' need a real release newer than yours. Step 1 covers them (tests 'Update path with preview', 'the yes in the pane', 'a dry run that gives no preview', 'an upgrade jaira refuses').
round: 2
---

# Cockpit-Band: Update-Hinweis mit Knopf 'Update' und Knopf 'Verwerfen'

## Definition of Done

- [x] Band zeigt bei neuerer Version Hinweis mit Knopf Update und Knopf Verwerfen; Update fuehrt Vorschau und danach Ja-Pfad aus; Verwerfen unterdrueckt den Hinweis fuer genau diese Version; Test mit Node-Lauf im Cockpit
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1180 band row Update/Dismiss, :356 previewUpgrade (self upgrade --check --json + update --roles --lanes --dry-run, pane; :369 no preview -> error pane, no 'Update to X'), :378 upgradeJaira (self upgrade + update --roles --lanes, only from the pane's yes, stdout and stderr shown), :346 dismissRelease (store dismissed-release), board.ts:520 offeredRelease; tests/board.test.ts:1572 describe 'a newer jaira release' (14 tests) - claude plugin test: 157 pass 0 fail; node 24 run of board.ts on real jaira output in scratch
- [x] Agenten-Weg mit Vorschau: stderr-Zeile (internal/cli/newerrelease.go) und Agenten-Block (core/board/announce.go) sagen: erst 'jaira self upgrade --check' und 'jaira update --roles --lanes --dry-run', dem User zeigen was sich aendert, fragen, erst auf Ja 'jaira self upgrade' und 'jaira update --roles --lanes'; Zeile bleibt eine Zeile; CLAUDE.md/AGENTS.md neu erzeugt
  proof: internal/cli/newerrelease.go:47 one stderr line: run 'jaira self upgrade --check' and 'jaira update --roles --lanes --dry-run', show the user, ask, only on yes the real commands; core/board/announce.go:112 agent block the same; CLAUDE.md:108 and AGENTS.md:81 regenerated with 'jaira update' of this branch; tests TestNextAndResumeNameANewerReleaseOnStderr / TestClaimAndMoveAskAboutANewerReleaseOnStderr (internal/cli/newerrelease_test.go, wants on the one release line) and TestAgentNoteTellsAgentsToAskAboutANewerRelease (core/board/announce_test.go:191)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] board.ts: releaseOf(cache json) und offeredRelease(installed, latest, dismissed); texts.ts: Woerter fuer Hinweis, Vorschau, Ergebnis
- [x] register.tsx: Cache jede Minute lesen (jairaHome geteilt mit settingsJson), Hinweiszeile mit Update und Dismiss im Band, Dismiss merkt Version im Store
- [x] register.tsx: Update = self upgrade --check --json + update --roles --lanes --dry-run, Pane mit Vorschau; Bestaetigen = self upgrade + update --roles --lanes; types/index.d.ts
- [x] tests/board.test.ts: Hinweis, Update mit Vorschau, nichts vor Bestaetigung, Ja-Pfad, Dismiss fuer genau diese Version; claude plugin test + validate
- [x] release note core/release/notes/EDSQZP.md, shipped.sha256 neu, Ja-Pfad mit echtem Binary gegen Stub-Release in Scratch-JAIRA_HOME pruefen
- [x] Cockpit: stderr von self upgrade ins Ergebnis-Pane (gap 1); Dry-run mit Exit ausser 0/3 -> Fehler-Pane, kein 'Update to X'; Tests
- [x] Agenten-Weg: newerrelease.go-Zeile und announce.go-Block auf check + dry-run + zeigen + fragen; Tests-Wants ergaenzen; CLAUDE.md/AGENTS.md mit jaira update dieses Branches neu; README/COMMANDS.md nachziehen
- [x] Release-Notiz zweite Zeile, shipped.sha256 neu, gofmt/vet/build, Commit

## Progress
- **2026-10-16 11:20 · BeMuCa** — Entscheidungen vor der Arbeit (keine offen fuer Berk, alle aus Code/Ticket/Default gesetzt):
1. Versionsdaten: Cockpit liest dieselbe Cache-Datei wie TUI und stderr-Zeile ($JAIRA_HOME bzw. ~/.jaira/update-check.json, Feld latest) und vergleicht mit 'jaira --version'; kein zweiter Check. Aufgefrischt wird der Cache weiter von TUI und next/resume/claim/move (PollCache, stuendlich); das Cockpit liest ihn jede Minute, wie das TUI.
2. Knopf-Texte englisch: 'Update' und 'Dismiss' statt 'Verwerfen' - texts.ts: jeder Cockpit-Text ist englisch, egal welche Sprache die Person spricht.
3. Verwerfen merkt sich die Version im Store des Cockpits ($.store, ueber Sitzungen). Gilt nur fuer das Band; die stderr-Zeile der CLI bleibt. Eine neuere Version zeigt den Hinweis wieder.
4. Update: 'jaira self upgrade --check --json', dann 'jaira update --roles --lanes --dry-run', beides im Board-Verzeichnis; Vorschau in einem Pane mit Knopf zum Bestaetigen und Abbrechen. Erst die Bestaetigung startet 'jaira self upgrade' und danach 'jaira update --roles --lanes'. --check schreibt nur den Cache, sonst nichts.
5. Ohne Board und in dev-Builds kein Hinweis (wie die TUI-Pille und die Regel des Cockpits: was ein Board braucht, bleibt ohne Board unsichtbar).
6. Mod-Dateien geaendert -> core/shipped/shipped.sha256 mit scripts/shipped-history.sh neu erzeugen (Test erzwingt es).
- **2026-10-16 11:32 · BeMuCa** — Gepruefte Wege (fuer critique/testing/review):
- claude plugin test im Mod-Ordner: 155 pass, 0 fail; 7 der 12 neuen Tests fallen mit dem alten register.tsx (die 5 'kein Hinweis'-Faelle sind Gegenproben und bestehen ohne Feature).
- Ja-Pfad mit echten Binaries in Scratch (HOME, JAIRA_HOME, PATH auf scratchpad/yes): Build dieses Branches als 0.5.1 und 0.6.0 (-X main.version, nicht release.Current - das setzt cli.Execute), Stub-Release per python http.server auf 127.0.0.1 ueber JAIRA_RELEASE_API/JAIRA_RELEASE_DOWNLOADS. jaira next schrieb den Cache (PollCache), --check JSON wie vom Cockpit gelesen, dry-run aenderte nichts (Snapshot gleich), self upgrade ersetzte nur das Scratch-Binary und zog die alte register.tsx in beiden Skills-Ordnern nach. ~/.local/bin/jaira vorher/nachher gleicher sha256.
- node 24 laedt hooks/board.ts direkt (Type-Stripping): offeredRelease/upgradeCheckOf auf der echten Ausgabe geprueft.
- tsc (npx typescript@5.6.3, tsconfig ausserhalb des Mod-Ordners): hooks/ und types/ fehlerfrei; tests/board.test.ts hat 5 alte Fehler (doppeltes toasts im World-Typ Z.42/69/142, Mounted.input Z.780/782) - nicht von hier, nicht angefasst.
Bekannte Grenze: die Vorschau ist der dry-run des HEUTE installierten jaira; was erst das neue Binary mitbringt, steht nicht darin (self upgrade laesst danach ohnehin das neue 'update --roles' laufen). Ein Upgrade im Terminal waehrend der Sitzung: das Cockpit merkt sich die Version je Sitzung, der Hinweis bleibt, bis Update gedrueckt wird - dann sagt --check 'already the latest' und der Hinweis geht.
- **2026-10-16 11:34 · BeMuCa** — critique (erster Pass, ganzer Diff 9d2cfd51): kein Befund. Geprueft und stehen gelassen: (1) eigener Pane cockpit-upgrade statt $.ui.ask wie beim Cache-Guard - ask ist der AskUserQuestion-Dialog fuer eine kurze Frage mit 'Other'-Freitext; die Vorschau ist mehrzeilig und folgt dem Muster des Explanation-Panes (Wahl in einem Pane mit Fokus). (2) run+output statt des jaira()-Helfers, weil die Vorschau stdout braucht und jaira() nur ok/why liefert. (3) 'update --roles --lanes' nach 'self upgrade' wiederholt --roles: noetig, weil ein Upgrade von 0.5.1 vom alten Binary gemacht wird, das update --roles noch nicht startet (Release-Notiz PWAF22). (4) Dismiss im Store statt in ~/.jaira: Entscheidung 3 oben.
- **2026-10-16 11:43 · BeMuCa** — review (fresh context, diff 9d2cfd51): checked by running, not only reading.
- claude plugin test in core/mod/builtin/jaira-cockpit: 155 pass, 0 fail.
- Counter-check in a scratch copy with the pre-change register.tsx: 7 of the 12 new tests fail (the 5 'no notice' cases pass either way, as the implementer said).
- scripts/shipped-history.sh regenerated into scratch: identical to core/shipped/shipped.sha256.
- node 24 on hooks/board.ts: offeredRelease/upgradeCheckOf as claimed, 0.10.0 counts as newer than 0.9.9.
- Go side matches: core/selfupdate/cache.go Path() (JAIRA_HOME or ~/.jaira/update-check.json, field latest); internal/cli/self.go --check --json fields current/latest/up_to_date; --check is not refused for Homebrew/go install, only the real upgrade is.
- The worktree_diff of the payload was only untracked .jaira/milestones/*.md - not this ticket.
Kept here because review-gaps replaces it: optimize's review-gaps was "nothing removed: the one duplicate (JAIRA_HOME-or-~/.jaira in settingsJson) is already folded into jairaHome(); every new text key and helper has a caller; considered one helper for 'stdout+stderr without warnings' in previewUpgrade/upgradeJaira - two callers with different fallbacks, left; pre-existing and left alone: metrics.csv is read and written under HOME/.jaira, not JAIRA_HOME". critique's review-summary was "none".
- **2026-10-16 12:45 · BeMuCa** — decision (Berk 'Finish alles', Hauptsitzung mit den empfohlenen Optionen): (1) gap 1 - stderr von self upgrade kommt ins Ergebnis-Pane. (2) Dry-run-Exit pruefen: schlaegt der Dry-run fehl, zeigt das Pane den Fehler und bietet kein 'Update to X'. (3) Dismiss bleibt nur im Band. (4) Neu: auch der Agenten-Weg zeigt vorher eine Vorschau - stderr-Zeile (internal/cli/newerrelease.go) und Agenten-Block (core/board/announce.go) sagen: erst 'jaira self upgrade --check' und 'jaira update --roles --lanes --dry-run', dem User zeigen, fragen, erst auf Ja die echten Befehle.
- **2026-10-16 12:45 · BeMuCa** — follow-up candidates (kein Ticket, Berk): (a) once() gibt die Update-Sperre nach BUSY_MS=20s frei, 'jaira self upgrade' darf 300s laufen - ein zweiter Druck waehrend eines langsamen Downloads kann eine zweite Vorschau und ein zweites Upgrade starten; (b) bei brew/go-install-Builds zeigt das Pane die Vorschau und 'Update to X', das nur scheitern kann - install_method steht schon im --check-JSON und koennte das vorher sagen; (c) vorher schon da: metrics.csv wird unter HOME/.jaira gelesen und geschrieben, nicht unter JAIRA_HOME (core/mod/builtin/jaira-cockpit/hooks/register.tsx, writeMetrics und der Ticket-Pane).
- **2026-10-16 12:52 · BeMuCa** — critique (Runde 2, nur efa8e9cd - die neue Arbeit; Runde 1 hatte keinen Befund): kein Befund. Geprueft: stderr-Zeile und Agenten-Block sagen dasselbe in derselben Reihenfolge (Vorschau, zeigen, fragen, erst auf Ja schreiben); die Zeile bleibt ein Fprintf mit einem \n. Exit 3 als Vorschau gelten lassen: update gibt editedError (ExitValidation=3) zurueck, nachdem es die ganze Vorschau gedruckt hat (internal/cli/update.go:115), das ist keine gescheiterte Vorschau. Die beiden Randfaelle 20s-Sperre und brew/go-install-Vorschau sind laut Berk Folge-Kandidaten, nicht Befunde.
- **2026-10-16 13:00 · BeMuCa** — review round 2 (fresh context, diff 9d2cfd51+efa8e9cd), checked by running, not only reading:
- claude plugin test in core/mod/builtin/jaira-cockpit: 157 pass, 0 fail, 14 tests in 'a newer jaira release'; the yes-path test pins the self-upgrade stderr in the result pane, and the dry-run-failure test pins 'no Update to X'.
- This branch built as 0.5.1 into scratch (HOME/JAIRA_HOME in scratch, release API on a dead local port), cache latest 9.9.9: 'jaira next --json' gives one stderr line with check, dry-run, show, ask and only-on-yes; stdout has no 9.9.9.
- Agent block in the committed CLAUDE.md and AGENTS.md is byte-identical to what that build writes on init.
- 'update --roles --lanes --dry-run' in the scratch repo: exit 0, before/after snapshot of repo and HOME equal.
- scripts/shipped-history.sh into scratch: identical to core/shipped/shipped.sha256. gofmt and vet clean on internal/cli and core/board.
- node 24 on hooks/board.ts: offeredRelease newer/dismissed/next-after-dismissed/same/dev/0.10.0>0.9.9 and upgradeCheckOf past the older-board warning, all as claimed.
- ~/.local/bin/jaira sha256 8b2014bd... before and after.
- worktree_diff of the payload was only the untracked .jaira/milestones/*.md, which is not this ticket.
Minor gaps only (see review-gaps): the result title says updated even if the setup step failed; the 'write nothing' wording ignores the cache file that --check rewrites. The two follow-up candidates (20s lock, brew/go-install) are still unfiled.
