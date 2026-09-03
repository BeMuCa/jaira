---
id: 01M1CACHXX142XS20N7SBDV0HM
title: Doku und Gate widersprechen sich beim Verlassen der Human-Review-Lane
status: human
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Es gibt genau eine Wahrheit dazu, ob ein Agent eine Abnahme auf Zuruf des Nutzers per --force ausfuehren darf - Doku und Gate sagen dasselbe"
context: "Gefunden am 31.08., als Berk sagte 'alle bei Human review auf done setzen': jaira move <id> --to done --force wurde fuer alle 8 Tickets verweigert. docs/AGENTS.md:155-160 und .claude/skills/jaira/SKILL.md:275-282 dokumentieren --force als den sanktionierten Weg, wenn der Nutzer die Abnahme im Gespraech ausgesprochen hat ('the decision stays human, only the keystroke is delegated'). core/gate/gate.go:377 blockt aber jede nicht-interaktive Bewegung aus einer requires-human-exit-Lane, unabhaengig von --force (CodeNeedsHuman). Eine Seite muss der anderen folgen; welche, ist eine Design-Entscheidung."
definition-of-done: Gate und beide Doku-Stellen sagen dasselbe; ein Test pinnt das gewaehlte Verhalten
blocked-by: []
commits: []
created-at: 2026-08-31T16:28:07Z
updated-at: 2026-09-01T09:40:38Z
question: "Soll das Gate --force auf deinen Zuruf wieder durchlassen (Doku bleibt, Verhalten wie dokumentiert) - oder bleibt das Gate hart und die zwei Doku-Stellen werden auf 'nur im Board' umgeschrieben? Kostenabwaegung: durchlassen = bequem, aber ein Agent kann 'der Nutzer hat es gesagt' halluzinieren; hart = jede Abnahme kostet dich den Gang ins Board."
updated-by: BeMuCa
---

# Doku und Gate widersprechen sich beim Verlassen der Human-Review-Lane

## Definition of Done

- [ ] Gate und beide Doku-Stellen sagen dasselbe; ein Test pinnt das gewaehlte Verhalten

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-01 09:40 · BeMuCa** — ZURUECKGEZOGEN am 01.09.: die Praemisse war falsch. Doku und Gate widersprechen sich NICHT - flow.go:198 laesst --force die needs_human-Violation ueberschreiben, genau wie SKILL.md/AGENTS.md es beschreiben; per --dry-run --force belegt ('would be allowed; It would override 1 gate refusal(s)'). Der Fehlalarm entstand, weil 'move --force | tail -1' die Override-Bullet der ERFOLGS-Ausgabe zeigt, die wie eine Verweigerung aussieht. Die 8 Zuruf-Tickets waren seit dem 31.08. in done.
