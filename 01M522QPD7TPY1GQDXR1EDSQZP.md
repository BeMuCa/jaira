---
id: 01M522QPD7TPY1GQDXR1EDSQZP
title: "Cockpit-Band: Update-Hinweis mit Knopf 'Update' und Knopf 'Verwerfen'"
status: testing
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
updated-at: 2026-10-16T12:52:22Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-482020
claimed-at: 2026-10-16T11:16:45Z
outcome-what: "Runde 1: jaira-cockpit-Band zeigt bei neuerer Version (Release-Cache $JAIRA_HOME|~/.jaira/update-check.json gegen 'jaira --version', jede Minute gelesen) 'jaira X is released (this is Y)' mit [ Update ] und Dismiss. Update: 'jaira self upgrade --check --json' + 'jaira update --roles --lanes --dry-run', Vorschau im Pane cockpit-upgrade (Cancel hat den Fokus), erst 'Update to X' fuehrt 'jaira self upgrade' und 'jaira update --roles --lanes' aus. Dismiss merkt die Version im Store. Runde 2: Ergebnis-Pane zeigt auch stderr von self upgrade (Host-Override-Warnung); ein Dry-run mit Exit ausser 0/3 oeffnet das Fehler-Pane statt der Vorschau und bietet kein Update an. Agenten-Weg: stderr-Zeile (internal/cli/newerrelease.go) und Agenten-Block (core/board/announce.go, CLAUDE.md, AGENTS.md) sagen jetzt: erst --check und --dry-run laufen lassen, dem User zeigen, fragen, erst auf Ja die echten Befehle. README, docs/COMMANDS.md, zwei Release-Notiz-Zeilen, shipped.sha256 neu."
outcome-why: "Die stderr-Zeile erreicht nur Sitzungen, die gerade jaira-Befehle fahren; das Band erreicht jede Sitzung mit offenem Cockpit. Berk wollte den Knopf mit Vorschau vor jedem Schreiben und einem Verwerfen, das sich die Version merkt - und dieselbe Vorschau auch auf dem Agenten-Weg ('kein Dry run noetig?')."
outcome-resolves: "DoD 1: Hinweis mit Update und Verwerfen (Dismiss), Update = Vorschau dann Ja-Pfad, Verwerfen nur fuer genau diese Version; claude plugin test 157 pass. DoD 2: stderr-Zeile und Agenten-Block verlangen die Vorschau vor der Frage; Go-Tests auf die eine Zeile und den Block."
review-summary: "none"
review-gaps: "nothing removed in round 2: the release line has one producer (internal/cli/newerrelease.go:47) and the block one (core/board/announce.go:112), CLAUDE.md/AGENTS.md are generated from it; releaseLineOf in newerrelease_test.go has two callers; dryRunFailed has one caller. Round 1 still stands: jairaHome() is the one shared JAIRA_HOME lookup; pre-existing metrics.csv under HOME/.jaira left alone (noted as follow-up candidate)"
test-verdict: "pass: go build ./... ok, go test -race -count=1 ./... RC=0 (32 ok, 0 FAIL); claude plugin validate passed, claude plugin test 155 pass 0 fail (12 in 'a newer jaira release'); DoD 1 proof lines checked in the tree; yes-path run with real binaries built as 0.5.1/0.6.0 against a local stub release in a scratch HOME/JAIRA_HOME: check + dry-run wrote nothing (snapshot equal), self upgrade replaced only the scratch binary, update --roles --lanes ran; ~/.local/bin/jaira unchanged; not seen: the band live in a Claude Code terminal"
review-verdict: "The diff meets the definition of done: a newer release gets a line with Update and Dismiss, Update shows a preview and writes only after 'Update to X', Dismiss hides exactly that version, and 12 cockpit tests cover it (7 of them fail without the change; 155 pass, 0 fail). One small defect - gap 1, the release-host override warning is dropped from the result pane - is worth a one-line fix but I would not hold the merge for it; gaps 2-4 are edges. Merge after T8A8KM and PWAF22."
review-round: 1
review-check: |-
  1. cd /home/berk/git/jAIra-EDSQZP/core/mod/builtin/jaira-cockpit && claude plugin test  -> ends with '155 pass' and '0 fail'; the block 'a newer jaira release' has 12 (pass) lines (notice, preview + Cancel, yes path, refused upgrade, Dismiss, next session).
  2. Hand check, nothing is replaced: mkdir -p /tmp/jh && echo '{"checked_at":"2026-10-16T09:00:00Z","latest":"9.9.9"}' > /tmp/jh/update-check.json
  3. cd /home/berk/git/jAIra-EDSQZP && JAIRA_HOME=/tmp/jh claude --plugin-dir core/mod/builtin/jaira-cockpit  -> the jaira-cockpit frame shows 'jaira 9.9.9 is released (this is 0.5.1) [ Update ] Dismiss'. (Not verified by me: whether your installed cockpit loads beside it; if you see two frames, use the one with this line.)
  4. Press Dismiss -> the line is gone and the notice says 'jaira 9.9.9 dismissed: the band offers the next release again'.
  5. Change 9.9.9 to 9.9.10 in /tmp/jh/update-check.json and wait up to one minute -> the line is back, now with 9.9.10.
  6. Press Update -> while 0.5.1 is still the newest real release: no pane opens, the notice says 'jaira 0.5.1 is already the latest release', the line is gone, and /tmp/jh/update-check.json now says latest 0.5.1. sha256sum ~/.local/bin/jaira is the same before and after.
  7. The preview pane and 'Update to X' need a real release newer than yours; step 1 covers them (tests 'Update path with preview' and 'the yes in the pane').
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
