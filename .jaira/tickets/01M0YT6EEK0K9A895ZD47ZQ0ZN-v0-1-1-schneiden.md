---
id: 01M0YT6EEK0K9A895ZD47ZQ0ZN
title: v0.1.1 schneiden
status: human
ready: true
creator: BeMuCa
goal: Alles nach 62989f1 ist als Release veroeffentlicht
context: "Seit dem letzten Release ist sehr viel gelandet: die Liste von sieben, drei fremde PRs, der board-bewusste Block. Schneiden heisst: ein Block '## 0.1.1' in core/release/NOTES.md, taggen, Tag pushen. In die Notizen gehoert die Verhaltensaenderung, nicht nur die Funktionsliste: ein handgeschriebenes '[-]' in einer Checkliste galt frueher als offen und gilt jetzt als zurueckgezogen, blockiert also den Abschluss nicht mehr."
definition-of-done: "core/release/NOTES.md hat einen 0.1.1-Block, der die [-]-Aenderung nennt; der Tag ist gepusht; jaira self upgrade --check findet das Release"
blocked-by: []
commits: []
created-at: 2026-08-26T10:35:02Z
updated-at: 2026-08-31T19:06:35Z
claimed-by: EE-3NX6GL3-2976569
claimed-at: 2026-08-31T19:05:43Z
updated-by: BeMuCa
assignee: BeMuCa
question: "Tag setzen? Zwei Befehle: git tag v0.1.1 && git push origin v0.1.1 - goreleaser schneidet dann das Release, danach findet jaira self upgrade --check es und der letzte DoD-Punkt ist erfuellt. Der NOTES-Block ist schon auf master."
outcome-what: "core/release/NOTES.md hat den 0.1.1-Block: 13 Einzeiler in Anweisungs-Stimme, neueste zuerst, die [-]-Verhaltensaenderung als erste Zeile (DoD-Pflicht); committet und gepusht"
outcome-why: "Release-Schnitt war faellig - 130+ Commits seit v0.1.0 inkl. zweier Renames unreleased Surface; der Block ist die Haelfte des DoD, die ein Agent liefern kann, der Tag-Push veroeffentlicht und gehoert dir"
outcome-resolves: go test ./core/release -count=1 gruen (der NOTES-Parser liest den Block); Zeilenregel eingehalten (jede Aenderung genau eine Zeile)
---

# v0.1.1 schneiden

## Definition of Done

- [ ] core/release/NOTES.md hat einen 0.1.1-Block, der die [-]-Aenderung nennt; der Tag ist gepusht; jaira self upgrade --check findet das Release

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

