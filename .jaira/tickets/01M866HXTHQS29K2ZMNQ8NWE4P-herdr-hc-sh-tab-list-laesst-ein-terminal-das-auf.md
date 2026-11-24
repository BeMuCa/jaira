---
id: 01M866HXTHQS29K2ZMNQ8NWE4P
title: "herdr-hc.sh 'tab list' laesst ein Terminal, das auf 'hc get' nicht antwortet, still weg und endet mit Exit 0"
status: testing
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
updated-at: 2026-11-24T07:26:23Z
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-1283952
claimed-at: 2026-11-24T07:16:12Z
outcome-what: "herdr-hc.sh 'tab list': Terminal ohne Antwort -> Liste mit \"partial\":true, stderr nennt es, Exit 0; not_found -> weggelassen. Give back (register.tsx:1791) liest partial wie eine gescheiterte Liste. Tests Go (7 Faelle) und Cockpit."
outcome-why: "Ein stummes Terminal fiel still aus der Liste und Give back startete einen zweiten Worker; Exit 3 (Runde 1) haette reap.sh und die Tab-Anzeigen fuer jedes fremde stumme Terminal lahmgelegt."
outcome-resolves: "Nur der Aufrufer, der Vollstaendigkeit braucht, reagiert; reap.sh/agentsOf/workers wie vorher. Geprueft: Rennen ls/get (not_found), hc ls scheitert (Exit 1), leere Liste."
review-summary: none
review-gaps: "Kommentar herdr-hc.sh:85 umformuliert (unklarer Satz); sonst none - tabsPartial dupliziert das try/catch von tabsOf bewusst (eigener Rueckgabetyp, ein Aufrufer); kein toter Code (Exit-3-Zweig ganz entfernt)."
test-verdict: "pass: go build/vet + go test -count=1 ./... green (rc=0), -race on role green, claude plugin test 544/544; DoD 1-3 in the tree verified (Proof-Zeilen herdr-hc.sh:85-101, register.tsx:1790, notes/8NWE4P.md); herdr-hc.sh mit hc-Stub: alle antworten -> Exit 0 + 2 Tabs; t7 stumm -> Exit 3, Liste ohne t7, stderr nennt t7; hc ls scheitert -> Exit 1, nur stderr. Schreibt nichts (keine Datei, kein Ticket) - Zwei-Schreiber-Test entfaellt."
hunt-verdict: "fail: bughunt done at eb0a645dc, 2 change: Exit 3 sperrt reap.sh/agentsOf/workers bei einem fremden stummen Terminal; not_found (Terminal zwischen ls und get entfernt) zaehlt als stumm"
round: 2
back-from: testing
back-by: lane
back-count: 1
---

# herdr-hc.sh 'tab list' laesst ein Terminal, das auf 'hc get' nicht antwortet, still weg und endet mit Exit 0

## Definition of Done

- [x] herdr-hc.sh 'tab list' meldet ein nicht antwortendes Terminal (Exit ungleich 0 oder Feld partial:true); Test mit einem hc-Stub, dessen get fuer ein Terminal scheitert
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:92-108 (partial:true, not_found weggelassen); TestHerdrHcTabListSaysWhichTerminalDidNotAnswer (core/role/role_test.go)
- [x] Give back behandelt eine unvollstaendige Liste wie eine gescheiterte (W80HNH): kein neuer Worker, Meldung an die Person; Test
  proof: register.tsx:1791 + sessions.ts tabsPartial; board.test.ts 'a tab list missing a terminal that gave no answer starts no worker either'
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
