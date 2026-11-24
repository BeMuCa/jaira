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
updated-at: 2026-11-24T07:20:18Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1283952
claimed-at: 2026-11-24T07:16:12Z
---

# herdr-hc.sh 'tab list' laesst ein Terminal, das auf 'hc get' nicht antwortet, still weg und endet mit Exit 0

## Definition of Done

- [x] herdr-hc.sh 'tab list' meldet ein nicht antwortendes Terminal (Exit ungleich 0 oder Feld partial:true); Test mit einem hc-Stub, dessen get fuer ein Terminal scheitert
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:85-101 (exit 3 + stderr); TestHerdrHcTabListSaysWhichTerminalDidNotAnswer (core/role/role_test.go)
- [x] Give back behandelt eine unvollstaendige Liste wie eine gescheiterte (W80HNH): kein neuer Worker, Meldung an die Person; Test
  proof: board.test.ts 'a tab list missing a terminal that gave no answer starts no worker either' - register.tsx:1790 (W80HNH) reads exit 3 as failed
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
