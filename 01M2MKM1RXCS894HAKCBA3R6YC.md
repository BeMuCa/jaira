---
id: 01M2MKM1RXCS894HAKCBA3R6YC
title: "Die pr-Rolle oeffnet den Pull Request selbst, wenn ein Mensch sie aufruft"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer /jaira-role-pr tippt, bekommt einen offenen Pull Request und keine Befehlszeile zum Einfuegen; ruft ein Agent die Rolle auf, bleibt es beim Pushen."
context: |-
  Heute pusht die Rolle jaira-role-pr den Branch und schreibt eine 'gh pr create'-Zeile aus, die der Mensch danach selbst einfuegt und ausfuehrt.

  Alex am 2026-09-16: 'если я стартую /jaira-role-pr то она должна сразу и создавать, просто без команды'. Sein Aufruf IST die Anweisung, den Pull Request zu oeffnen - der Kopierschritt danach kauft nichts und kostet eine Runde.

  Die Regel dahinter bleibt: wer die Aenderung geschrieben hat, entscheidet nicht, dass sie ankommt. Sie wird nur anders aufgehaengt - nicht mehr an 'kein Agent oeffnet je einen PR', sondern an 'nur der Aufruf eines Menschen oeffnet einen'. Ein Dispatcher oder Teamlead, der die Rolle aufruft, pusht weiter nur und reicht die Zeile zurueck. Merge und Approve bleiben in jedem Fall verboten.

  Der Text ist bereits geschrieben und liegt auf dem lokalen Branch chore/pr-role-opens-the-pr (Commit 5b849ca): core/role/builtin/jaira-role-pr/SKILL.md umgeschrieben, eine Zeile in core/release/NOTES.md. Er ist noch nicht durch critique, optimize oder review gelaufen und traegt keine Ticket-Id im Commit.

  Dabei gefunden und offen: in einem Fork zielt 'gh pr create' nicht selbstverstaendlich auf das richtige Repository - origin ist der Fork, die Pull Requests dieses Boards gehen aber alle nach BeMuCa/jaira. Die Rolle muss das pruefen, bevor sie das Kommando ausfuehrt, sonst steht der Pull Request im falschen Repository und ein Mensch muss ihn von Hand schliessen.
definition-of-done: "core/role/builtin/jaira-role-pr/SKILL.md sagt: Aufruf durch einen Menschen -> pushen und oeffnen; Aufruf durch einen Agenten -> pushen und die Zeile zurueckgeben; merge und approve bleiben in beiden Faellen verboten; die Rolle prueft vor dem Oeffnen, in welches Repository der Pull Request geht; eine Zeile unter ## Unreleased in core/release/NOTES.md; go test ./core/role/... gruen"
tags:
  - docs
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-16T07:59:07Z
updated-at: 2026-09-16T08:02:06Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-11544
claimed-at: 2026-09-16T07:59:59Z
---

# Die pr-Rolle oeffnet den Pull Request selbst, wenn ein Mensch sie aufruft

## Definition of Done

- [ ] core/role/builtin/jaira-role-pr/SKILL.md sagt: Aufruf durch einen Menschen -> pushen und oeffnen; Aufruf durch einen Agenten -> pushen und die Zeile zurueckgeben; merge und approve bleiben in beiden Faellen verboten; die Rolle prueft vor dem Oeffnen, in welches Repository der Pull Request geht; eine Zeile unter ## Unreleased in core/release/NOTES.md; go test ./core/role/... gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Den Zielrepository-Check in core/role/builtin/jaira-role-pr/SKILL.md vor den create-Befehl ziehen und konkret machen: wie man das Ziel ermittelt (gh repo view / git config jaira.remote), was bei einem Fork gilt, und dass ohne klares Ziel nicht geoeffnet wird
- [ ] NOTES.md-Zeile unter ## Unreleased so ergaenzen, dass der Zielrepository-Check darin steht
- [ ] go test ./core/role/... laufen lassen, Ticket verschieben und zusammen mit dem Code committen

## Progress

