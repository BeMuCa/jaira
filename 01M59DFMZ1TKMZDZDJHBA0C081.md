---
id: 01M59DFMZ1TKMZDZDJHBA0C081
title: "Ein Ticket laesst sich mit Grund ablehnen: jaira reject, und das Archiv sagt warum"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer ein Ticket nicht macht, sagt mit einem Befehl warum - und Monate spaeter steht der Grund am Ticket, getrennt von bloss weggeraeumten."
context: |-
  Was heute fehlt: ein abgelehntes Ticket sieht genauso aus wie ein weggeraeumtes.
  - 'jaira archive <id>' (core/ticket/store.go:301, Taste x auf dem Board) verschiebt nach .jaira/archive/ und schreibt KEINEN Grund.
  - 'jaira delete <id>' loescht endgueltig - nur fuer versehentlich angelegte Tickets, keine Ablehnung.
  - Folge: 'jaira archive' (Liste) und das Web-Portal koennen nicht sagen, ob ein Ticket ein Duplikat, verworfen oder ueberholt war. Das verletzt den Kernwert 'was war das Ticket und wo steht es'.
  Anlass: Alex am 19.10.2026 - 'wir haben keine richtige Moeglichkeit, ein Ticket zu loeschen oder abzulehnen'.
  Entschieden von Alex (19.10.2026):
  - Neuer Befehl 'jaira reject <id> --why <text> [--duplicate-of <id>]'. Art: wontfix | duplicate | obsolete (--duplicate-of setzt duplicate).
  - Schreibt ins Frontmatter einen Block 'rejected:' mit kind, why, by, at; danach archiviert er wie 'jaira archive' (gleicher Code, kein zweiter Weg).
  - Nur ein Mensch lehnt ab. Ein Agent schiebt das Ticket stattdessen nach human mit --question 'ablehnen? Grund: ...'. Wie ein Agent erkannt wird: denselben Mechanismus nehmen, der heute den Ausgang aus human sperrt (requires-human-exit, siehe VW1WB1) - nichts Neues erfinden.
  Verworfen: Grund nur an archive haengen (Aufraeumen und Ablehnen blieben ununterscheidbar); eigene Lane 'rejected' (aendert die Lanes jedes Boards).
  Offen fuer die Plan-Lane: ob x auf dem Board kuenftig nach dem Grund fragt oder eine eigene Taste ablehnt (x bleibt archivieren, siehe model.go:1532-Kommentar).
definition-of-done: "'jaira reject <id> --why <text>' schreibt rejected.kind/why/by/at ins Frontmatter und verschiebt das Ticket nach .jaira/archive/; ohne --why Exit 2"
tags:
  - cli
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-19T06:27:18Z
updated-at: 2026-10-19T06:36:40Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-16821
claimed-at: 2026-10-19T06:27:46Z
mode: conversational
---

# Ein Ticket laesst sich mit Grund ablehnen: jaira reject, und das Archiv sagt warum

## Definition of Done

- [x] 'jaira reject <id> --why <text>' schreibt rejected.kind/why/by/at ins Frontmatter und verschiebt das Ticket nach .jaira/archive/; ohne --why Exit 2
  proof: internal/cli/reject.go newRejectCmd; core/ticket/store.go Store.Reject; tests TestRejectWritesTheReasonAndArchives, TestRejectWithoutWhyIsAUsageError
- [x] --duplicate-of <id> setzt kind=duplicate und verlinkt das Original; eine unbekannte id wird mit Exit 5 abgelehnt
  proof: core/ticket/store.go Store.Reject (LoadAnywhere, related); test TestRejectDuplicateOfLinksTheOriginal, TestRejectDuplicateOfAnUnknownTicketIsNotFound
- [x] Ein Agent wird mit Exit 3 abgewiesen und die Meldung nennt den Weg ueber human --question
  proof: internal/cli/reject.go !force -> ExitValidation needs_human, questionLane; test TestRejectWithoutForceSendsAnAgentToTheHumanLane
- [ ] Auf dem Board laesst sich ein Ticket mit einer einzeiligen Begruendung ablehnen
- [ ] 'jaira archive' (Liste) und 'jaira show' zeigen rejected-kind und rejected-why; 'jaira restore' nimmt die rejected-* Felder wieder ab
- [ ] Release-Note core/release/notes/<ID>.md, Agent-Block/jaira-Skill nennen reject; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Schema: Felder rejected-kind/why/by/at in core/ticket/schema.go, Ticket.Rejected, Store.Reject (schreibt Felder, ruft Archive) und Store.Restore entfernt sie
- [ ] CLI: internal/cli/reject.go - --why Pflicht (Exit 2), --kind wontfix|obsolete, --duplicate-of (Exit 5 bei unbekannter id, related setzen), ohne --force Exit 3 mit Weg ueber human --question
- [ ] TUI: Taste R fragt eine Zeile Grund und ruft Store.Reject; x bleibt archivieren
- [ ] Anzeige: jaira archive (Liste) und jaira show zeigen kind und why (Web-Portal: TPCKZB)
- [ ] Release-Note, Agent-Block, jaira-Skill; go test ./... -race

## Progress
- **2026-10-19 06:27 · Alexander Sacharov** — parent session: jaira-a7
- **2026-10-19 06:28 · Alexander Sacharov** — decision (worker, settled): Art waehlt --kind wontfix|obsolete, Default wontfix; --duplicate-of setzt duplicate und verbietet --kind. Auf dem Board wird mit Default wontfix abgelehnt - eine Zeile Grund reicht (DoD 4).
- **2026-10-19 06:28 · Alexander Sacharov** — question 1 (an jaira-a7): Der Mechanismus aus VW1WB1 ist 'Interactive' (core/gate/gate.go:395) - nur das Board setzt ihn, jeder CLI-Aufruf gilt als Agent. Damit widersprechen sich DoD 1 (CLI reject klappt) und DoD 3 (Agent Exit 3). Optionen: 1) (Recommended) CLI reject ohne --force -> Exit 3 mit Hinweis auf human --question bzw. das Board; ein Mensch in der Shell nimmt --force, genau wie heute bei move aus human. 2) CLI reject nur ueber das Board, kein --force-Weg - DoD 1 gilt dann fuer den Store-Befehl, den das Board aufruft. 3) CLI reject immer erlaubt, DoD 3 faellt weg - Agenten werden nur per Skill-Regel abgehalten.
- **2026-10-19 06:28 · Alexander Sacharov** — question 2 (an jaira-a7): Welche Taste lehnt auf dem Board ab? 1) (Recommended) eigene Taste R (in der Detailansicht und auf dem Board frei), fragt nach einer Zeile Grund; x bleibt reines Archivieren - Aufraeumen und Ablehnen bleiben getrennt, wie der Kommentar model.go:1532 es fuer x verlangt. 2) x fragt nach einem Grund, leer = nur archivieren - eine Taste weniger, aber jedes Archivieren bekommt eine Frage dazwischen.
- **2026-10-19 06:31 · Alexander Sacharov** — decision (Alex via jaira-a7, 19.10.2026): Q1 -> Option 1: CLI reject ohne --force -> Exit 3, Meldung nennt human --question und das Board; ein Mensch in der Shell nimmt --force, wie bei move aus human. Q2 -> Option 1: eigene Taste R fragt nach einer Zeile Grund; x bleibt reines Archivieren.
- **2026-10-19 06:32 · Alexander Sacharov** — question 3 (an jaira-a7): Kontext/DoD verlangen einen Block 'rejected:' mit kind/why/by/at. Das Format kennt aber nur flache Top-Level-Keys: core/ticket/schema.go:69-72 begruendet, warum outcome-what/why/resolves flach sind (Merge-Driver loest pro Top-Level-Key, ein Feld = eine Zeile), und der Writer (core/ticket/frontmatter.go SetScalar/SetList) kann keine verschachtelten Mappings schreiben. Optionen: 1) (Recommended) flache Keys rejected-kind, rejected-why, rejected-by, rejected-at - wie outcome-*, kein neuer Writer, 'jaira set' und Merge-Driver funktionieren sofort. 2) verschachtelter Block 'rejected:' wie im Ticket - braucht einen neuen Block-Writer im Frontmatter-Code und weicht von der Regel in schema.go ab.
- **2026-10-19 06:33 · Alexander Sacharov** — decision (Alex via jaira-a7): Q3 -> Option 1: flache Keys rejected-kind, rejected-why, rejected-by, rejected-at wie outcome-*. DoD 1/5 'rejected.kind/why/by/at' und 'rejected-Block' sind damit diese vier Keys.
- **2026-10-19 06:35 · Alexander Sacharov** — scope (Teamlead via jaira-a7): Web-Portal-Teil von DoD 5 wandert nach TPCKZB (Web-Filter, parallel auf feat/tpckzb) - zwei Branches an core/web gleichzeitig wuerden kollidieren. TPCKZB zeigt rejected-* an, wenn vorhanden. Hier wird core/web nicht angefasst.
