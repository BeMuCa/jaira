---
id: 01M4ATQ46NZ08TFVQARZ9WSS49
title: "Die README nennt jede Rolle in einer Zeile, und das Tutorial liegt im Repo"
status: critique
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Wer jaira neu benutzt, findet in der README in einer Zeile pro Skill, was er tut, und im Repo eine Seite, die zeigt, welcher Weg fuer welche Aufgabe passt"
context: |-
  Heute beschreibt README.md (Abschnitt 'Roles: a teamlead, a dispatcher, and the workers', ca. Zeile 615) nur vier Rollen ausfuehrlich; die anderen fuenf stehen nur hinter 'jaira roles list names the rest'.
  Berk am 06.10.: hat nicht gewusst, welche Skills es gibt und ob man mit Claude oder mit den Skills redet. Daraus entstand im Chat eine Liste mit einem Einzeiler pro Skill und eine HTML-Seite (Workflow-Baum, alle Rollen, Abnahmeseite, Einrichten).
  Berk am 07.10.: die Einzeiler sollen so in die README; die Seite soll ins Repo unter docs/.
  Die Seite ist auf Deutsch (so gewuenscht), README und docs/*.md sind Englisch - die Einzeiler kommen deshalb auf Englisch in die README.
  Die Seite wurde ausserhalb des Repos gebaut und liegt bisher unversioniert in /home/berk/git/jAIra/docs/how-to-work-with-jaira.html (auf Branch feat/PAP369-link-lines, falscher Ort).
  Kein Release-Notes-Eintrag: README und docs/ sind nicht im Binary.
definition-of-done: "README.md listet im Rollen-Abschnitt jeden Skill (jaira plus die neun aus 'jaira roles list') mit genau einer Zeile, was er tut; die bestehenden ausfuehrlichen Absaetze bleiben"
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-10-07T09:22:07Z
updated-at: 2026-10-07T09:25:46Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-3605060
claimed-at: 2026-10-07T09:23:32Z
outcome-what: README-Rollenabschnitt bekommt eine Tabelle mit einer Zeile pro Skill; die Tutorial-Seite liegt jetzt unter docs/ und ist auf 0.3.5 gezogen
outcome-why: "Berk wusste nicht, welche Skills es gibt und wie man sie anspricht; die README nannte fuenf der neun Rollen nur als 'the rest'"
outcome-resolves: "DoD 1: die Tabelle steht vor den bestehenden Absaetzen, die bleiben unangetastet. DoD 2: die committete Datei ist genau die, aus der Artifact-Version 3 veroeffentlicht wurde"
executed-by: opus
---

# Die README nennt jede Rolle in einer Zeile, und das Tutorial liegt im Repo

## Definition of Done

- [x] README.md listet im Rollen-Abschnitt jeden Skill (jaira plus die neun aus 'jaira roles list') mit genau einer Zeile, was er tut; die bestehenden ausfuehrlichen Absaetze bleiben
  proof: README.md:623-634 Tabelle Skill|What it does, 10 Zeilen (jaira + die neun aus 'jaira roles list' von 0.3.5); Absaetze darunter unveraendert
- [x] docs/how-to-work-with-jaira.html ist committet und deckt sich mit der veroeffentlichten Seite
  proof: docs/how-to-work-with-jaira.html; dieselbe Datei als Version 3 nach https://claude.ai/artifact/FAnMJiZz6dG2hb3it2GM2c veroeffentlicht

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

