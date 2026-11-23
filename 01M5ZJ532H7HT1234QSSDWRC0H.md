---
id: 01M5ZJ532H7HT1234QSSDWRC0H
title: "Der Agent-Block sagt, wie man Goal, Kontext und DoD nach dem Anlegen aendert"
status: done
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
updated-at: 2026-11-23T22:38:29Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: worker-dwrc0h
claimed-at: 2026-11-23T01:25:43Z
outcome-what: "Agent-Block (announce.go, CLAUDE.md, AGENTS.md) nennt jaira set goal=/context= und jaira dod --add/--text fuer Aenderungen nach dem Anlegen; create --help verweist darauf; Test fuer dod --add mit stabilen Nummern und Proofs; Release-Note."
outcome-why: "Agenten hielten Goal/DoD nach dem Anlegen fuer unveraenderlich und behalfen sich mit Notizen, weil der Agent-Block die vorhandenen Befehle nicht nannte."
outcome-resolves: "Alle drei DoD-Punkte; kein neuer Befehl noetig. Vier Kategorien: keine neue Logik, nur Text und ein Test; Nebenlaeufigkeit/Fehlerpfade/Rechte unveraendert, unuebliche Eingabe (Leerzeichen im set-Wert) auf Scratch-Board geprueft."
review-summary: "Docs-only change: one new bullet in the agent block (core/board/announce.go agentNote, mirrored verbatim in CLAUDE.md and AGENTS.md) says a ticket is changed after capture through the CLI, not a note - 'jaira set <id> goal=... context=...' for goal and context, 'jaira dod <id> --add' to append a DoD item, 'jaira dod <id> <n> --text' to reword one, with existing numbers and proofs unchanged. 'jaira create --help' gains a closing paragraph pointing at set and dod --add. A release note core/release/notes/DWRC0H.md tells users to run jaira update. No new command and no logic change; a CLI test added in the first commit was removed again in the second as a duplicate of core/ticket/proof_test.go TestAddItemAfterProofCarryingLastItem, which is the proof for DoD 2."
review-gaps: |-
  note: DoD 2 rests on a unit test of ticket.AddItem (core/ticket/proof_test.go:115); internal/cli/checklist.go:240 does call AddItem, but no test drives 'jaira dod --add' end to end since the CLI test was dropped in 60387d0e0 - re-checked by hand on a scratch board, item 2 kept [x] and its proof after --add
  note: the create --help paragraph (internal/cli/tickets.go:219-221) names set and dod --add but not dod <n> --text, which the agent block does name; one more clause would make the two match
test-verdict: "pass: gofmt clean, go vet (core/board, internal/cli, core/ticket) RC=0, build RC=0; go test not run by house rule (main session runs it). DoD 1 and 3 read in the tree (announce.go:57-61, create --help last paragraph); DoD 2 proof test core/ticket/proof_test.go:115 read, not run. Scratch board: set goal/context and dod --add/--text behave as documented; set on a missing id exits 5, dod --add with empty text exits 2, nothing written; two parallel dod --add both landed (3 items). Agent-file command part identical to jaira init output."
hunt-verdict: "pass: bughunt done at 60387d0e0, 0 findings (docs-only change)"
review-verdict: "accept: all three DoD items hold - the agent block names jaira set for goal and context with the literal syntax, dod --add keeps existing numbers and proofs (unit test plus a hand try), and create --help points at set; gofmt, go vet and go build are clean; no defect found in the touched text or the commands it documents"
review-round: 1
review-verdict-gaps: |-
  note: DoD 2 rests on a unit test of ticket.AddItem (core/ticket/proof_test.go:115); internal/cli/checklist.go:240 does call AddItem, but no test drives 'jaira dod --add' end to end since the CLI test was dropped in 60387d0e0 - re-checked by hand on a scratch board, item 2 kept [x] and its proof after --add
  note: the create --help paragraph (internal/cli/tickets.go:219-221) names set and dod --add but not dod <n> --text, which the agent block does name; one more clause would make the two match
review-check: |-
  1. In the worktree run: go build -o /tmp/jr ./cmd/jaira
  2. Make a scratch repo: mkdir /tmp/jt && cd /tmp/jt && git init && git commit --allow-empty -m init && /tmp/jr init
  3. /tmp/jr create probe --goal g --context c --dod one --dod two   (note the 6-letter id it prints)
  4. /tmp/jr dod <id> 2 --done --proof P2
  5. /tmp/jr set <id> goal="new goal" context="new context"
  6. /tmp/jr dod <id> --add three
  7. /tmp/jr show <id>  -> goal and context are the new ones; item 1 one [ ], item 2 two [x] with proof P2, item 3 three [ ] - nothing renumbered
  8. grep -n "changes after capture" CLAUDE.md  -> the new bullet with jaira set and dod --add
  9. /tmp/jr create --help | tail -15  -> paragraph "What create sets can be changed later"
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
- **2026-11-23 01:29 · BeMuCa** — testing: 'jaira validate' im Repo warnt, der jaira-Block in CLAUDE.md/AGENTS.md passe nicht zu Lanes/.jaira/pr. Der Befehlsteil (bis '## This board's lanes') ist mit der Generator-Ausgabe identisch (diff gegen 'jaira init' auf Scratch-Board); die Abweichung liegt im Lane-/PR-Teil und kommt nicht aus diesem Ticket. 'jaira update' im Worktree hat der Auto-Mode-Classifier abgelehnt, daher nicht weiter eingegrenzt.
