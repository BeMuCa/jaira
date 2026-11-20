---
id: 01M7WWYB5XZMD01DZY3MJG3REW
title: Ein Agent darf auf Wunsch den PR oeffnen und eine Review-Runde (Copilot) selbst fahren
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Eine Einstellung erlaubt Agenten, den Branch zu pushen, den PR zu oeffnen und eine Kommentarrunde abzuarbeiten (Guard, Push, Antwort je Kommentar mit Commit, resolve, Review neu anfordern); ohne die Einstellung bleibt alles wie heute."
context: |-
  Heute: Regel in CLAUDE.md und jaira-role-pr: ein Agent pusht und stoppt, PR oeffnen/mergen ist Sache des Menschen.
  Folge: ein Teamlead macht bei jedem Ticket und jeder Fix-Runde von Hand checkout, Guard, Push, PR-Text, Copilot anfordern, antworten, resolve.
  Alex 20.11.2026: er erlaubt Agenten inzwischen oefter PRs zu erstellen, auch mit Copilot als Reviewer, aber nicht ueberall -> optional, pro Repo/Board.
  Gewuenscht vom Teamlead: eine Rolle oder ein Skript pr-round <ID>; die Person gibt nur die Freigabe.
  Mergen und Approven bleiben in jedem Fall beim Menschen.
  Verwandt: AX8K02 (fester Schritt vor dem Push), XB97NJ (Board sagt, ob Arbeit per PR kommt).
definition-of-done: Ohne Einstellung oeffnet kein Agent einen PR (Test oder Skill-Text belegt)
tags: []
blocked-by: []
related:
  - 01M2HZKCQ48BPQYN1YPTXB97NJ
  - 01M5ZJ5DKQZ83GYBGPJHAX8K02
  - 01M2N12WVABVRYM7C923WATJDQ
  - 01M7WWYG213BWB06A7D1CRASZS
commits: []
created-at: 2026-11-20T16:34:26Z
updated-at: 2026-11-20T16:47:34Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-59469
claimed-at: 2026-11-20T16:43:41Z
question: ""
---

# Ein Agent darf auf Wunsch den PR oeffnen und eine Review-Runde (Copilot) selbst fahren

## Definition of Done

- [ ] Ohne Einstellung oeffnet kein Agent einen PR (Test oder Skill-Text belegt)
- [ ] Mit Einstellung oeffnet jaira-role-pr den PR selbst und fordert den konfigurierten Reviewer (z.B. Copilot) an
- [ ] Eine Kommentarrunde: jeder Kommentar bekommt 'Behoben in <sha>: ...' oder eine Begruendung, wird resolved, das Review neu angefordert
- [ ] Merge und Approve macht kein Agent, auch mit Einstellung

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-20 16:43 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/35112.sock
- **2026-11-20 16:44 · Alexander Sacharov** — decision (settled, worker): die Einstellung ist eine Board-Datei wie .jaira/release - .jaira/pr, eine Zeile 'person' (Default, wie heute) oder 'agent'; der Reviewer steht in .jaira/pr-reviewer (z.B. copilot). Grund: Vorbild ReleaseFor in core/settings/settings.go, reist mit dem Repo (XB97NJ: nicht settings.json, nicht git config); sichtbar ueber 'jaira roles list --json' wie release.
- **2026-11-20 16:44 · Alexander Sacharov** — decision (settled, worker): die Kommentarrunde laeuft ueber Rollentext plus ein Skript in jaira-role-pr/scripts (offene Threads listen, antworten, resolve per gh api graphql, Review neu anfordern) - der Teamlead wollte 'pr-round <ID>'; ein Skript macht den GraphQL-Teil wiederholbar. Merge/Approve bleibt in jedem Fall verboten.
- **2026-11-20 16:44 · Alexander Sacharov** — question: Wann darf ein Agent mit .jaira/pr=agent den PR oeffnen? 1) (Recommended) Sobald alle Tickets des Branches in done sind - wie der Cockpit-Accept heute; die Abnahme der Person ist die Freigabe, Copilot reviewt danach. 2) Schon in signoff - Copilot reviewt, bevor die Person abnimmt; die Person sieht Copilot-Ergebnis beim Accept, aber ein PR existiert fuer noch nicht abgenommene Arbeit. 3) Nur auf ausdrueckliche Freigabe pro PR (Teamlead gibt --by-person weiter, nachdem die Person ja gesagt hat) - Einstellung spart nur die Handarbeit, nicht die Frage.
- **2026-11-20 16:47 · Alexander Sacharov** — answer (Alex via teamlead, 2026-11-20): 1) - ein Agent oeffnet den PR mit .jaira/pr=agent erst, wenn alle Tickets des Branches in done sind; die Abnahme der Person ist die Freigabe, Copilot reviewt danach. Die selbst getroffenen Entscheidungen bleiben stehen. Mode bleibt autonom (eine beantwortete Frage).
