---
id: 01M86H5QJWEKFHXJXA710TREV1
title: "reap.sh:418 reads any failed 'herdr pane get' as a closed tab, not only…"
status: in-progress
ready: true
creator: Sascha
context: |-
  Was falsch ist: core/role/builtin/jaira-dispatcher/scripts/reap.sh:418 - schlaegt 'herdr pane get <pane>' fehl, bleibt tab leer, und reap.sh behandelt den Tab als schon geschlossen.
  Folge: hat herdr nur einen Moment nicht geantwortet (Timeout, Neustart), entfernt reap.sh Worktree und Branch, obwohl claude im Tab noch darin arbeitet.
  Der Check herdr_up (reap.sh:133, 'tab list') laeuft einmal am Anfang der Runde; ein Ausfall danach faellt nicht auf.
  8NWE4P (24.11.) hat nur den Fall 'hc get antwortet bei tab list nicht' abgefangen (Exit 3 -> reap behaelt alles).
  Fehlerformate: Herdr App (hc) schreibt {"error":{"kind":"not_found"}} auf stderr, Exit 3; altes herdr.exe kennt den Code pane_not_found (aus strings herdr.exe).
  Gleiches Muster, nicht dieses Ticket: close-tab.sh:159 (dort parkt es nur aus, entfernt nichts).
tags:
  - cockpit
blocked-by: []
related: []
follows: 01M866HXTHQS29K2ZMNQ8NWE4P
commits: []
created-at: 2026-11-24T10:21:10Z
updated-at: 2026-11-26T08:07:02Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-35248
claimed-at: 2026-11-26T07:54:32Z
goal: "reap.sh liest ein gescheitertes 'herdr pane get' nur dann als geschlossenen Tab, wenn herdr den Pane nicht kennt (not_found); jede andere Fehlschlag (Timeout, herdr weg) behaelt Tab, Worktree und Record fuer die naechste Runde."
outcome-what: "reap.sh: pane get/tab get ueber hget(); nur not_found (hc) bzw. pane_not_found (altes herdr) gilt als geschlossener Tab, jeder andere Fehlschlag -> 'kept <t>: herdr gave no answer ...', Runde fuer das Ticket haelt an. close-tab.sh: gleiche Regel, ein stummes herdr hebt den Park nicht auf. Tests: reap 6 Faelle, close-tab 5 Faelle."
outcome-why: "reap.sh las jeden Fehlschlag von 'herdr pane get' als geschlossenen Tab und entfernte Worktree und Branch, auch wenn herdr nur kurz nicht antwortete und claude noch im Tab sass."
outcome-resolves: Ein Timeout oder Neustart von herdr kostet keinen Worktree mehr; nur herdr selbst kann einen Tab fuer geschlossen erklaeren.
review-summary: none
round: 3
back-from: testing
back-by: lane
back-count: 1
review-gaps: "reap.sh hget(): fails with 1 instead of returning 2, callers use 'if ! hget' (rc bookkeeping gone); left the twin helper close-tab.sh ask() - no shared library between the scripts; herdr-hc.sh 'tab list' keeps its own grep for \"not_found\" (different loop, stderr file already open); mktemp per call is per record, not hot"
test-verdict: "pass: go vet ./... and go test ./... green (RC=0), shipped-history --check ok; DoD 1-4 verified in the tree (TestReapReadsOnlyNotFoundAsAClosedTab 7 cases, close-tab cases); hget against real hc: unknown terminal -> treated as closed, own terminal -> tab id, timing-out hc stub -> rc 1 with its message; failure path: worktree, branch and records stay (test); the change adds no writer - reap.lock and park's mv unchanged"
hunt-verdict: "fail: bughunt done at 4899c2218, 1 change (close-tab.sh ask(): mktemp/tr fail outside the ERR trap -> exit 1 without kept), 0 pre-existing"
---

# reap.sh:418 reads any failed 'herdr pane get' as a closed tab, not only…

## Definition of Done

- [x] reap.sh:418 schliesst nur bei not_found/pane_not_found auf einen geschlossenen Tab; jeder andere Fehlschlag von 'pane get' gibt 'kept <ticket>: ...' aus und laesst Worktree, Branch und Record stehen; Test in core/role/reap_test.go mit einem herdr-Stub, dessen pane get ohne not_found scheitert
  proof: reap.sh hget() + pane/tab get in the close loop; TestReapReadsOnlyNotFoundAsAClosedTab (timeout, silent, tab get timeout: worktree, branch, records kept)
- [x] Ein Pane, den herdr nicht kennt (not_found bzw. pane_not_found), wird wie bisher als geschlossen behandelt; Test
  proof: TestReapReadsOnlyNotFoundAsAClosedTab cases 'pane unknown to hc', 'pane unknown to the old herdr', 'tab unknown to hc'
- [x] Release-Note in core/release/notes/0TREV1.md
  proof: core/release/notes/0TREV1.md
- [x] Gleicher Fehler in close-tab.sh:159 (Geschwister): ein stummes herdr hebt den Park nicht mehr auf; Test
  proof: close-tab.sh ask(); TestReapClosesAParkedTabOnceItsCacheIsCold/'herdr gives no answer', TestCloseTabKeepsAWorkerThatIsNotFinished timeout/not_found/'no temp directory' cases

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-26 07:54 · Sascha** — parent session: uds:/run/user/1000/cc-socks/62036.sock
- **2026-11-26 07:56 · Sascha** — decision: keine offenen Entscheidungen. Unbekannter Fehlschlag von pane get -> 'kept <t>: herdr gave no answer for pane <p>' und die ganze Runde fuer das Ticket bleibt stehen (wie 'closing failed'). Erkennung per grep auf not_found im kombinierten stdout+stderr: deckt hc (kind not_found) und altes herdr (pane_not_found) ab, ohne den Stream des alten herdr zu kennen. close-tab.sh:159 bleibt draussen (entfernt nichts).
- **2026-11-26 07:59 · Sascha** — gefunden: altes herdr.exe-Format nur aus strings bekannt (pane_not_found), Stream und Exit-Code unbekannt - darum grep auf stdout+stderr. Geschwister ohne Schaden gelassen: run-lane.sh:153 sagt 'already gone' auch bei Timeout (nur Text, Tab bleibt offen), spawn.sh:475 meldet ohnehin einen Fehler. Bestehender Test 'herdr not reachable' erwartet jetzt 'herdr gave no answer' statt 'does not know pane' - richtiger. Vier Kategorien: Nebenlaeufigkeit - reap.lock unveraendert; Fehlerpfad - jeder Nicht-not_found-Fehler haelt die Runde an (failed=1); Eingabe - leere Fehlermeldung -> 'exit N, no message'; Rechte - nichts Neues von aussen.
- **2026-11-26 08:00 · Sascha** — critique: (1) *not_found* zu weit - herdr kennt command_not_found/plugin_not_found/agent_not_found; nur die gequoteten Codes "not_found" und "pane_not_found" zaehlen, Test mit command_not_found. (2) hget Rueckgabe 1 ungenutzt - wie ask(): 0 mit leerem hout.
- **2026-11-26 08:02 · Sascha** — critique Runde 2: nichts mehr. Zwei Kopien des Helfers (reap.sh hget, close-tab.sh ask) bewusst stehen gelassen: die Skripte teilen keine Bibliothek, und eine gesourcte Datei braeuchte Install- und Manifest-Aenderungen fuer 10 Zeilen.
- **2026-11-26 08:03 · Sascha** — testing: hget gegen echtes hc geprueft - t99999 (not_found) -> leer, rc 0; eigenes Terminal -> tab_id; Stub mit Timeout -> rc 1, herr 'hc: request timed out'. waiting: bughunt (Hintergrund-Agent) und go test ./... seit 09:03
- **2026-11-26 08:05 · Sascha** — bughunt: change: core/role/builtin/jaira-dispatcher/scripts/close-tab.sh:165 - class: a failing command inside a function in close-tab.sh ends the script with exit 1 and no 'kept' line, because the ERR trap (line 95) is not inherited by functions without set -E; new in ask() (mktemp at :165, the tr read at :167), and already there in park() (:153). The usage text promises exit 0 with 'kept <ticket>: ...'. example: TMPDIR not writable -> mktemp exits 1 -> close-tab.sh exits 1 with only mktemp's message. Not a blocker: reap.sh calls close-tab with '|| true', and a worker's tab just stays open. ran: bash model of the same set -e + ERR trap + function shape (rc=1, trap silent); TMPDIR=/nonexistent mktemp (rc=1)
- **2026-11-26 08:05 · Sascha** — hunt report: 1 finding, 1 change (an edge, not a blocker), as a note on 0TREV1 | read whole: reap.sh, close-tab.sh (plus herdr-hc.sh, which both call), reap_test.go and closetab_test.go fakes, release note, shipped.sha256 (both script hashes are present) | could not check how old (non-hc) herdr answers 'tab get' for a gone tab (tab_not_found?): no herdr binary on this machine. If it answers that way, the pane-to-tab race is read as 'no answer' and kept for one round only. hc's not_found (exit 3) and unavailable (exit 5) were confirmed live
- **2026-11-26 08:05 · Sascha** — bughunt: done at 4899c2218 - 1 finding (1 change, 0 pre-existing) in 6 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/release/notes/0TREV1.md core/role/builtin/jaira-dispatcher/scripts/close-tab.sh core/role/builtin/jaira-dispatcher/scripts/reap.sh core/role/closetab_test.go core/role/reap_test.go core/shipped/shipped.sha256
- **2026-11-26 08:06 · Sascha** — bughunt-Fix: ask() prueft mktemp und tr selbst (kept-Zeile). set -E verworfen: dann feuert die ERR-Falle auch in jeder $(...)-Substitution, und ein ungeschuetztes x="$(cmd)" bekaeme die kept-Zeile als Wert statt die Runde zu beenden. park() (close-tab.sh:153) hat dieselbe Luecke seit vorher - nicht angefasst, Notiz fuer die Person. Altes herdr: strings herdr.exe kennt kein tab_not_found, also kein weiterer Code noetig.
