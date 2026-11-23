---
id: 01M5ZJ532H7HT1234QSSDWRC0H
title: "Der Agent-Block sagt, wie man Goal, Kontext und DoD nach dem Anlegen aendert"
status: testing
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
updated-at: 2026-11-23T01:29:46Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: worker-dwrc0h
claimed-at: 2026-11-23T01:25:43Z
outcome-what: "Agent-Block (announce.go, CLAUDE.md, AGENTS.md) nennt jaira set goal=/context= und jaira dod --add/--text fuer Aenderungen nach dem Anlegen; create --help verweist darauf; Test fuer dod --add mit stabilen Nummern und Proofs; Release-Note."
outcome-why: "Agenten hielten Goal/DoD nach dem Anlegen fuer unveraenderlich und behalfen sich mit Notizen, weil der Agent-Block die vorhandenen Befehle nicht nannte."
outcome-resolves: "Alle drei DoD-Punkte; kein neuer Befehl noetig. Vier Kategorien: keine neue Logik, nur Text und ein Test; Nebenlaeufigkeit/Fehlerpfade/Rechte unveraendert, unuebliche Eingabe (Leerzeichen im set-Wert) auf Scratch-Board geprueft."
review-summary: none
review-gaps: "removed internal/cli/checklist_test.go TestDodAddKeepsExistingNumbersAndProofs: core/ticket/proof_test.go:115 TestAddItemAfterProofCarryingLastItem already proves the same (append lands after the proof-carrying last item, earlier numbers and proofs stay); DoD 2 proof repointed there. Nothing else: text only."
---

# Der Agent-Block sagt, wie man Goal, Kontext und DoD nach dem Anlegen aendert

## Definition of Done

- [x] Der Agent-Block nennt jaira set fuer goal und context mit einem Beispiel
  proof: core/board/announce.go:57-61 (bullet 'a ticket changes after capture'), mirrored in CLAUDE.md:23-27 and AGENTS.md:23-27
- [x] Ein DoD-Punkt laesst sich nach dem Anlegen anhaengen, ohne dass bestehende Nummern wandern; ein Test belegt es
  proof: core/ticket/proof_test.go:115 TestAddItemAfterProofCarryingLastItem (existing): AddItem, the function behind 'jaira dod --add', appends below the last item and its proof, item 1 keeps number and proof; also tried by hand on a scratch board (item 2 [x] with proof, --add three -> 1,2 unchanged, 3 new)
- [x] jaira create --help verweist auf set fuer spaetere Aenderungen
  proof: internal/cli/tickets.go:219-221, last paragraph of 'jaira create --help'

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 01:26 · BeMuCa** — Keine offene Entscheidung. Geklaert aus dem Code: DoD-Punkt nachtragen = 'jaira dod <id> --add <text>' (gibt es schon, haengt im richtigen Abschnitt an, Nummern bleiben); kein neuer Befehl, kein 'set definition-of-done --append' (definition_of_done ist nur der erste Punkt, waere der falsche Weg). Umformulieren = 'jaira dod <id> <n> --text'. Fehlend: Satz im Agent-Block (core/board/announce.go), Verweis in 'create --help', ein Test dass --add die Nummern und den Proof bestehender Punkte laesst.
- **2026-11-23 01:28 · BeMuCa** — Bullet in den Abschnitt 'Capturing' gesetzt, nicht unter 'Working a ticket': feat/hpfd7d aendert die Milestone-Zeilen weiter unten im selben Block, so bleibt der Merge konfliktfrei. CLAUDE.md/AGENTS.md nur um diese 5 Zeilen ergaenzt, gegen 'jaira init' auf einem Scratch-Board verglichen (identisch). Beispiel set goal=... context=... und dod --add auf dem Scratch-Board von Hand ausprobiert. Test nicht gelaufen (go test in dieser Session verboten).
- **2026-11-23 01:28 · BeMuCa** — critique: nichts. Nur Text und ein Test; die Befehle gab es schon (dod --add/--text, set), kein neues Muster. Bullet steht im Capturing-Abschnitt neben create, wo ein Agent nach dem Anlegen sucht.
- **2026-11-23 01:29 · BeMuCa** — optimize: eigener CLI-Test entfernt, weil core/ticket/proof_test.go:115 genau diese Eigenschaft schon testet; DoD 2 zeigt jetzt darauf.
- **2026-11-23 01:29 · BeMuCa** — bughunt: done at 60387d0e0 - 0 findings (0 change, 0 pre-existing) in 6 files, first hunt; concurrency, error paths, input, permissions checked; the documented set goal/context, dod --add and dod --text were run on a scratch board and keep numbers and proofs; files: AGENTS.md CLAUDE.md core/board/announce.go core/release/notes/DWRC0H.md internal/cli/checklist_test.go internal/cli/tickets.go
