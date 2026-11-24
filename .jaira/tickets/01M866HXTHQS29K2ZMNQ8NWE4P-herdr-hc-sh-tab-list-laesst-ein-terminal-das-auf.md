---
id: 01M866HXTHQS29K2ZMNQ8NWE4P
title: "herdr-hc.sh 'tab list' laesst ein Terminal, das auf 'hc get' nicht antwortet, still weg und endet mit Exit 0"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Kann herdr-hc.sh ein Terminal nicht abfragen, ist 'tab list' als unvollstaendig erkennbar (Exit ungleich 0 oder ein Markierungsfeld), und das Cockpit startet daraufhin keinen zweiten Worker."
context: |-
  Was falsch ist: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:88 - in 'tab list' wird jedes Terminal mit 'hc get' gefragt; antwortet eines nicht, faellt es per '|| continue' aus der Liste, und das Skript endet trotzdem mit Exit 0.
  Folge: das Cockpit (sessions.ts:97, Give back) sieht den Tab eines lebenden Workers nicht und startet einen zweiten Worker auf demselben Ticket und Worktree.
  W80HNH hat den Fall 'tab list scheitert ganz' behoben (kein Worker, Meldung an die Person); eine Teil-Liste mit Exit 0 erkennt es nicht. Gefunden von Bughunt und Review zu W80HNH am 24.11.2026; gab es schon vorher.
  Gilt nur unter Herdr App (hc); unter dem alten herdr liefert 'tab list' alles auf einmal.
definition-of-done: "herdr-hc.sh 'tab list' meldet ein nicht antwortendes Terminal (Exit ungleich 0 oder Feld partial:true); Test mit einem hc-Stub, dessen get fuer ein Terminal scheitert"
tags:
  - cockpit
blocked-by: []
related: []
follows: 01M8543JH7M7VXA9THHZW80HNH
commits: []
created-at: 2026-11-24T07:15:35Z
updated-at: 2026-11-24T07:31:10Z
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-1283952
claimed-at: 2026-11-24T07:16:12Z
outcome-what: "herdr-hc.sh 'tab list': Terminal ohne Antwort -> Liste mit \"partial\":true, stderr nennt es, Exit 0; not_found -> weggelassen. Give back (register.tsx:1791) liest partial wie eine gescheiterte Liste. Tests Go (7 Faelle) und Cockpit."
outcome-why: "Ein stummes Terminal fiel still aus der Liste und Give back startete einen zweiten Worker; Exit 3 (Runde 1) haette reap.sh und die Tab-Anzeigen fuer jedes fremde stumme Terminal lahmgelegt."
outcome-resolves: "Nur der Aufrufer, der Vollstaendigkeit braucht, reagiert; reap.sh/agentsOf/workers wie vorher. Geprueft: Rennen ls/get (not_found), hc ls scheitert (Exit 1), leere Liste."
review-summary: none
review-gaps: "Kommentar herdr-hc.sh:85 umformuliert (unklarer Satz); sonst none - tabsPartial dupliziert das try/catch von tabsOf bewusst (eigener Rueckgabetyp, ein Aufrufer); kein toter Code (Exit-3-Zweig ganz entfernt)."
test-verdict: "pass: Runde 2 (seit eb0a645dc): go test ./role green, claude plugin test 544/544 green, shipped-history --check ok; DoD 1-3 im Baum geprueft; herdr-hc.sh mit hc-Stub: alle antworten -> Exit 0 ohne partial; t7 Timeout -> Exit 0, partial:true, stderr nennt t7; t9 not_found -> weggelassen, kein partial; hc ls scheitert -> Exit 1. Schreibt nur eine mktemp-Datei, per trap entfernt - kein geteilter Schreibort."
hunt-verdict: "fail: bughunt done at fd7a26ff1, 1 change blocker: partial:true liest nur Give back; workers/agentsOf/Overview/Chat-Cleanup lesen die Teil-Liste als vollstaendig; 1 pre-existing (reap.sh:418 pane get-Fehler = Tab zu)"
round: 3
back-from: testing
back-by: lane
back-count: 2
---

# herdr-hc.sh 'tab list' laesst ein Terminal, das auf 'hc get' nicht antwortet, still weg und endet mit Exit 0

## Definition of Done

- [x] herdr-hc.sh 'tab list' meldet ein nicht antwortendes Terminal (Exit ungleich 0 oder Feld partial:true); Test mit einem hc-Stub, dessen get fuer ein Terminal scheitert
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:96-109 (not_found weggelassen, sonst Exit 3 + stderr); TestHerdrHcTabListSaysWhichTerminalDidNotAnswer (core/role/role_test.go)
- [x] Give back behandelt eine unvollstaendige Liste wie eine gescheiterte (W80HNH): kein neuer Worker, Meldung an die Person; Test
  proof: register.tsx:1790 (W80HNH) liest Exit 3 als gescheitert; board.test.ts 'a tab list missing a terminal that gave no answer starts no worker either'
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/8NWE4P.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-24 07:15 · Alexander Sacharov** — teamlead 24.11.2026: Alex: starten. Gestapelt auf feat/w80hnh (W80HNH wird gerade gelandet); vor dem PR upstream/master einmergen.
- **2026-11-24 07:16 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/3417.sock
- **2026-11-24 07:17 · Alexander Sacharov** — decision: offene Entscheidung (Exit!=0 oder partial:true) selbst entschieden: Exit 3 + stderr-Zeile je stummem Terminal, die Liste bleibt auf stdout. Grund: alle drei Cockpit-Aufrufer (register.tsx:348 agentsOf, :1785 Give back, :1949 workers) und reap.sh:133 pruefen schon den Exit-Code; Give back lehnt seit W80HNH ab und nennt stderr. partial:true haette drei TS-Stellen gebraucht, und reap.sh haette eine Teil-Liste als 'herdr up' gelesen. Kosten: ein dauerhaft stummes fremdes Terminal sperrt Give back/reap, bis es weg ist - sicherer Fehler, Meldung nennt das Terminal.
- **2026-11-24 07:20 · Alexander Sacharov** — Geprueft: Rennen (Terminal zwischen hc ls und hc get entfernt) -> einmal Exit 3, Give back lehnt ab, naechster Versuch geht - sicher. hc ls scheitert -> set -e, Exit !=0 wie vorher. reap.sh:133 liest Exit 3 als 'herdr antwortet nicht' und raeumt nichts ab - gewollt. workers-Tool (register.tsx:1949) zeigt bei Exit 3 alle Records ohne Tab-Filter - wie bei W80HNH.
- **2026-11-24 07:20 · Alexander Sacharov** — critique: nichts. Einfachste Form: 9 Zeilen im Aufgerufenen (herdr-hc.sh), kein Cockpit-Code - W80HNH-Zweig register.tsx:1790 deckt Exit!=0 schon ab. Exit 3 statt 1 unterscheidet kein Aufrufer; bleibt, weil er im Log 'Liste unvollstaendig' von 'hc kaputt' trennt. Der zweite Cockpit-Test sieht wie der W80HNH-Test aus, prueft aber, dass eine Teil-Liste auf stdout bei Exit!=0 nicht gelesen wird.
- **2026-11-24 07:21 · Alexander Sacharov** — waiting: bughunt (Hintergrund-Subagent) since 08:21
- **2026-11-24 07:22 · Sascha** — bughunt: change: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:91 - class: every failure of 'hc get' counts as 'gave no answer', including a definite not_found for a terminal that was removed between 'hc ls' and 'hc get' (a concurrent 'tab close'/'hc rm' from close-tab.sh or reap.sh, or any terminal of the person exiting); example: ls lists t2, t2 is removed, 'hc get t2' answers {"kind":"not_found"} exit 3, tab list exits 3 and Give back refuses to start a worker although the list was complete; ran: hc get t99999 --json (exit 3, kind not_found) + fake HC_BIN where get t2 fails (tab list exit=3). Only a timeout/refusal should mark a terminal mute; not_found means the tab is closed and can be left out
- **2026-11-24 07:22 · Sascha** — bughunt: change: core/role/builtin/jaira-dispatcher/scripts/reap.sh:133 - class: callers of 'tab list' that read any non-zero exit as 'Herdr down' now throw away the whole list because of one mute terminal, which may not even be a jaira tab; places: reap.sh:133 (herdr_up stays empty, so parked tabs are never closed and every worktree is kept with 'herdr not reachable' for as long as any terminal is mute), cockpit register.tsx:353 agentsOf (board shows no Herdr tabs at all), register.tsx:1953 workers (tab filter dropped, closed-tab records listed as live), register.tsx:1182 (listed=false after a failed spawn.sh removes the chat's ticket copy from under a chat whose tab did come up); example: fake HC_BIN with t1 answering and t2 failing -> 'tab list >/dev/null 2>&1 && herdr_up=1' leaves herdr_up empty; ran: that line against the fake hc. Callers are older (blame e0bbc68e8, feed258f4, 05331f861) but get the new exit 3 only from this change; only Give back (the case the ticket targets) wants exit 3 to mean 'unknown'
- **2026-11-24 07:22 · Sascha** — hunt report: 2 findings, both change, neither a blocker, as notes on 8NWE4P | read whole: herdr-hc.sh, the give-back/agentsOf/workers/openChat paths of cockpit register.tsx and sessions.ts tab helpers, reap.sh herdr_up path; tests (board.test.ts, role_test.go) and shipped.sha256 checked, hashes match and go test ./role passes | not run against the real Herdr App: whether 'hc ls' can list a terminal that 'hc get' then reports not_found depends on hc timing; reproduced with a fake hc and a real 'hc get' on a missing address
- **2026-11-24 07:22 · Sascha** — bughunt: done at eb0a645dc - 2 findings (2 change, 0 pre-existing) in 5 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/8NWE4P.md core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-24 07:23 · Sascha** — testing: class: ein Teil-Ergebnis von 'tab list' als Exit!=0 meldet allen Aufrufern 'Herdr down', obwohl nur Give back Vollstaendigkeit braucht; example: reap.sh:133 herdr_up leer -> nichts wird mehr abgeraeumt, solange irgendein Terminal stumm ist (Bughunt-Note); fix: Exit 0 lassen, Feld partial:true in der Liste, nur Give back (register.tsx:1790) liest es als gescheitert; dazu not_found aus 'hc get' als geschlossen weglassen, nicht als stumm. Test ueber Klassen: antwortet / not_found / Timeout.
- **2026-11-24 07:25 · Sascha** — Runde 2: Exit 3 verworfen (Bughunt: reap.sh, agentsOf, workers lasen ein fremdes stummes Terminal als 'Herdr down'). Jetzt Exit 0 + "partial":true + stderr; nur Give back liest das Feld (sessions.ts tabsPartial). workers/agentsOf zeigen eine Teil-Liste weiter an - eine Anzeige, kein Start, bewusst gelassen. not_found aus 'hc get' (stderr, Exit 3, real geprueft mit hc get t99999) = geschlossen, weggelassen.
- **2026-11-24 07:26 · Sascha** — critique Runde 2: nichts. partial-Feld im Aufgerufenen, ein Leser (tabsPartial neben tabsOf in sessions.ts, gleiches try/catch-Muster); Exit-Code-Zweig von W80HNH unveraendert mitbenutzt.
- **2026-11-24 07:26 · Sascha** — waiting: bughunt Runde 2 (Hintergrund-Subagent) since 08:26
- **2026-11-24 07:28 · Sascha** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1952-1956 + sessions.ts:152 - class: every reader of 'herdr tab list' except givingBack (register.tsx:1791) reads a list marked "partial":true as complete, because 34dc412bf moved the 'incomplete' signal from the exit code (which every reader honoured) into a JSON field only Give back reads; places: workers() register.tsx:1952 -> workerLines sessions.ts:152 drops the carrier whose tab is the mute one, so the teamlead's 'workers' tool says no worker runs (sessions.ts:72-76 says Give back and 'workers' share one rule - they no longer do); agentsOf register.tsx:353 (listed=true) feeds the band, the Overview workers pane (readWorkers :1281) and the Discuss/Explain cleanup at :1182, which deletes the copy of a chat whose tab is the mute terminal; example: hc ls names t1,t7, hc get t7 times out -> tab list exit 0 with partial:true and only t1 -> workerLines returns [] for MG677R's worker, while without tabs it lists it; ran: bun scratch on sessions.ts (partial list -> [], no list -> 1 line) and herdr-hc.sh with a fake hc (exit 0, partial:true); unclear provenance: on master a mute terminal was also dropped with exit 0, so vs master behaviour is the same - marked change because the path runs through 34dc412bf's lines and eb0a645dc had it covered
- **2026-11-24 07:28 · Sascha** — bughunt: pre-existing: core/role/builtin/jaira-dispatcher/scripts/reap.sh:418 - class: reap.sh reads any failed 'herdr pane get' (timeout, transport error) as 'the person closed the tab', not only an unknown pane; example: a done+merged+clean ticket whose worker tab is a Herdr App terminal that times out on hc get -> herdr_up=1 (tab list exits 0, reap.sh:133 ignores partial) -> pane get exits 1 -> tab empty -> record counted gone -> teardown, 'git worktree remove --force', 'branch -D' and the records deleted while claude still sits in that tab, which no later round closes because its record is gone; herdr-hc.sh:87-89 justifies exit 0 by 'reap.sh only asks whether Herdr answers', which is not so; eb0a645dc's exit 3 had kept reap off this path, 34dc412bf restores master's behaviour; blame: 5c5207ba67 (2026-10-17, before the ticket); ran: herdr-hc.sh with a fake hc - tab list exit 0, pane get t7 exit 1
- **2026-11-24 07:28 · Sascha** — hunt report: 2 findings (1 change blocker, 1 pre-existing), notes on 8NWE4P | read whole: core/mod/builtin/jaira-cockpit/hooks/sessions.ts, register.tsx at every tab-list reader (givingBack, agentsOf and its 4 callers, workers()) not all 2847 lines; read as diff since eb0a645dc: herdr-hc.sh (plus whole file), role_test.go, board.test.ts, release note, shipped.sha256; also read reap.sh/close-tab.sh as callers | not checked: real Herdr App timeout behaviour of hc get (only not_found confirmed against the real hc: exit 3, kind not_found); close-tab.sh treats a failed pane get as gone too but only keeps, so no harm there
- **2026-11-24 07:28 · Sascha** — bughunt: done at fd7a26ff1 - 2 findings (1 change, 1 pre-existing) in 7 modules, follow-up hunt since eb0a645dc; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/sessions.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/8NWE4P.md core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-24 07:29 · Sascha** — testing Runde 2: class: das Signal 'Liste unvollstaendig' muss jeder Leser von 'tab list' sehen, nicht nur Give back; example: workers() register.tsx:1952 verliert den Worker im stummen Terminal (Bughunt-Note). decision (Worker): zurueck zu Exit 3, not_found bleibt weggelassen. Grund: Exit!=0 lesen alle Leser schon als 'unbekannt' (Cockpit-Zweige, reap.sh:133); ein partial-Feld braeuchte dieselbe Pruefung an 5 Stellen plus reap.sh = Exit 3 mit mehr Code. Die Kosten aus Bughunt Runde 1 (reap.sh raeumt nichts ab, solange ein Terminal stumm ist) sind die sichere Seite: mit Exit 0 laeuft reap.sh bis reap.sh:418 und baut den Worktree eines lebenden Workers ab (pre-existing Note Runde 2). Dieser pre-existing Fund ist mit Exit 3 fuer hc-Timeouts zu, fuer andere pane-get-Fehler bleibt er - Note fuer die Person.
