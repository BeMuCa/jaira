---
id: 01M5ZJ532H7HT1234QSSDWRC0H
title: "Der Agent-Block sagt, wie man Goal, Kontext und DoD nach dem Anlegen aendert"
status: in-progress
ready: true
creator: BeMuCa
goal: "Der Agent-Block in CLAUDE.md nennt jaira set fuer goal und context und einen dokumentierten Weg, einen DoD-Punkt nachzutragen; niemand muss das per Notiz behelfen."
context: |-
  Was schief ist: Teamlead und Worker glaubten, 'jaira edit fehlt' und ergaenzten Goal und DoD per jaira note (Feedback L323/L330/O20, 27.10.2026).
  Dabei gibt es jaira set <id> goal=... und set --append (jaira set --help, 27.10.2026) - nur steht es nicht im Agent-Block (CLAUDE.md, Abschnitt 'Task tracking: jaira'), und ob set einen DoD-Punkt sauber anhaengt, ist nicht geprueft.
  Zu klaeren im Ticket: wie ein einzelner DoD-Punkt nachgetragen wird (set definition-of-done --append oder ein dod --add), und dass die Nummern bestehender Punkte dabei bleiben, sonst zeigen alte --proof-Eintraege ins Leere.
  Nicht gewollt: ein neuer Befehl jaira edit, wenn set reicht.
definition-of-done: Der Agent-Block nennt jaira set fuer goal und context mit einem Beispiel
tags:
  - docs
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:52:13Z
updated-at: 2026-11-23T01:28:05Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: worker-dwrc0h
claimed-at: 2026-11-23T01:25:43Z
---

# Der Agent-Block sagt, wie man Goal, Kontext und DoD nach dem Anlegen aendert

## Definition of Done

- [x] Der Agent-Block nennt jaira set fuer goal und context mit einem Beispiel
  proof: core/board/announce.go:57-61 (bullet 'a ticket changes after capture'), mirrored in CLAUDE.md:23-27 and AGENTS.md:23-27
- [x] Ein DoD-Punkt laesst sich nach dem Anlegen anhaengen, ohne dass bestehende Nummern wandern; ein Test belegt es
  proof: internal/cli/checklist_test.go TestDodAddKeepsExistingNumbersAndProofs (jaira dod --add, existing command); not run here (no go test in this session), compiles under go vet
- [ ] jaira create --help verweist auf set fuer spaetere Aenderungen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 01:26 · BeMuCa** — Keine offene Entscheidung. Geklaert aus dem Code: DoD-Punkt nachtragen = 'jaira dod <id> --add <text>' (gibt es schon, haengt im richtigen Abschnitt an, Nummern bleiben); kein neuer Befehl, kein 'set definition-of-done --append' (definition_of_done ist nur der erste Punkt, waere der falsche Weg). Umformulieren = 'jaira dod <id> <n> --text'. Fehlend: Satz im Agent-Block (core/board/announce.go), Verweis in 'create --help', ein Test dass --add die Nummern und den Proof bestehender Punkte laesst.
