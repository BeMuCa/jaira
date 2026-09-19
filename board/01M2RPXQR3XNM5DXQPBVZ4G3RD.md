---
id: 01M2RPXQR3XNM5DXQPBVZ4G3RD
title: Der Lane-Payload zeigt den Diff des Tickets und nicht den des Branches
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer 'jaira show <id> --for-lane <lane> --json' liest, bekommt entweder den Diff der ganzen Ticket-Arbeit oder eine Angabe, dass er nur einen Ausschnitt sieht - nicht stillschweigend beides verwechselbar"
context: |-
  Gefunden in der review-Lane von 01M2NCH8ZK9524JZ00J4GTQHNH (GTQHNH) am 2026-09-17, dort ausdruecklich als vorbestehend und ausserhalb jenes Tickets stehengelassen.

  Was schiefgeht: showForLane in internal/cli/flow.go nimmt 'shas := t.Commits' und fragt git nur dann nach den Commits des Tickets, wenn dieses Feld leer ist. Auf GTQHNH stehen in 'commits:' drei SHAs, auf dem Branch liegen einundzwanzig - alle mit dem Handle im Betreff, alle die Ticket-Datei anfassend. Der review-Worker bekam also den Diff aus drei von einundzwanzig Commits, und er kam mit 'complete: true'. Achtzehn Commits, darunter der gesamte validate-, resume- und announce-Teil, waren unsichtbar, und nichts im Payload sagte das.

  Warum das Feld veraltet: 'commits:' wird beim 'jaira move' aus der Implementierungs-Lane einmal geschrieben. Jede spaetere Runde durch critique/in-progress schreibt es nicht nach. Bei einem Ticket, das sechzehn Kritik-Runden laeuft, ist der Eintrag nach der ersten Runde falsch.

  Die Folge trifft die zwei Lanes, die genau dafuer da sind, den Diff zu lesen: critique und review. Ein Reviewer, der dem Payload glaubt, unterschreibt einen Bruchteil.

  Was schon entschieden ist: die Prompts von jaira-role-lane und jaira-dispatcher sagen seit GTQHNH, wie man es merkt (Zaehlprobe 'jq .commits | length' gegen 'git log master..HEAD') und was man stattdessen liest ('git diff master...HEAD'). Das ist ein Pflaster, keine Reparatur.

  Zu entscheiden ist, was 'commits:' bedeuten soll: die Aufzeichnung dessen, was eine Lane abgegeben hat (dann muss der Payload sie ergaenzen statt ersetzen), oder ein Zwischenspeicher der Ableitung (dann muss er nachgefuehrt oder fallengelassen werden). Das ist der Grund, warum es nicht als Einzeiler in GTQHNH mitlief.
definition-of-done: "Ein Ticket, dessen 'commits:' aelter ist als sein Branch, bekommt von 'jaira show <id> --for-lane critique --json' entweder den Diff aller seiner Commits oder ein Feld, das sagt, wie viele fehlen - nachgestellt an einem Ticket mit gefuelltem 'commits:' und weiteren Commits danach"
tags:
  - cli
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-09-17T22:13:48Z
updated-at: 2026-09-17T22:13:48Z
---

# Der Lane-Payload zeigt den Diff des Tickets und nicht den des Branches

## Definition of Done

- [ ] Ein Ticket, dessen 'commits:' aelter ist als sein Branch, bekommt von 'jaira show <id> --for-lane critique --json' entweder den Diff aller seiner Commits oder ein Feld, das sagt, wie viele fehlen - nachgestellt an einem Ticket mit gefuelltem 'commits:' und weiteren Commits danach
- [ ] Die Bedeutung von 'commits:' ist in docs/ einmal aufgeschrieben, damit die naechste Aenderung nicht wieder raet
- [ ] Die Zaehlprobe in core/role/builtin/jaira-role-lane/SKILL.md und core/role/builtin/jaira-dispatcher/SKILL.md ist entfernt oder auf den neuen Zustand umgeschrieben - das Pflaster bleibt nicht neben der Reparatur stehen
- [ ] Eine Zeile in core/release/NOTES.md unter ## Unreleased

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

