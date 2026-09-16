---
id: 01M2N12WVABVRYM7C923WATJDQ
title: Der GitLab-Arm der pr-Rolle benutzt ein abgekuendigtes Flag und laesst --head weg
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Auf einem GitLab-Board oeffnet /jaira-role-pr den Merge Request im richtigen Projekt und aus dem richtigen Fork-Branch, mit Flags, die glab noch kennt."
context: |-
  In core/role/builtin/jaira-role-pr/SKILL.md steht der GitLab-Zweig der Zielrepository-Leiter. Die review-Lane von A3R6YC hat ihn am echten glab geprueft und zwei Fehler gefunden:

  1. 'glab mr create --target-project <pfad>' ist in glab 1.114 abgekuendigt. Nachstellen: 'glab mr create --target-project foo/bar --title x --description y' - erste Ausgabezeile ist 'Flag --target-project has been deprecated, Use --repo instead.' Der Befehl legt nichts an, er bricht danach ab, weil hier kein GitLab-Remote konfiguriert ist.
  2. Im Fork-Fall fehlt '--head <fork>'. 'glab mr create --help' zeigt im Beispielblock 'glab mr create --repo upstream/project --head your-namespace/project ...' - zwei Flags, wo SKILL.md nur eines nennt. Ohne --head sucht glab den Quell-Branch im Zielprojekt und findet ihn nicht.

  Der GitHub-Arm ist nachgewiesen richtig und wurde auf dem Fork sashasoft90/jaira durchgespielt; nur GitLab ist ungeprueft. Dieses Board liegt auf GitHub, deshalb wurde A3R6YC mit diesem Rest angenommen statt zurueckgeschickt - der Fehler faellt hier niemandem auf und trifft erst das erste Board auf GitLab.

  Dazu offen und NICHT nachgestellt: die review-Lane vermutet, dass 'gh pr create' ohne '--head' aus einem Fork-Clone haengen bleibt. Wer dieses Ticket nimmt, kann das gleich mitpruefen - es ist dieselbe Stelle.

  Die Aenderung ist reiner Text in einer Datei.
definition-of-done: "core/role/builtin/jaira-role-pr/SKILL.md nennt im GitLab-Zweig '--repo' statt '--target-project' und ergaenzt '--head <namespace>/<projekt>' fuer den Fork-Fall; 'glab mr create --help' belegt beide Flags; eine Zeile unter ## Unreleased in core/release/NOTES.md; go test ./core/role/... gruen"
tags:
  - docs
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-16T11:54:25Z
updated-at: 2026-09-16T11:54:25Z
---

# Der GitLab-Arm der pr-Rolle benutzt ein abgekuendigtes Flag und laesst --head weg

## Definition of Done

- [ ] core/role/builtin/jaira-role-pr/SKILL.md nennt im GitLab-Zweig '--repo' statt '--target-project' und ergaenzt '--head <namespace>/<projekt>' fuer den Fork-Fall; 'glab mr create --help' belegt beide Flags; eine Zeile unter ## Unreleased in core/release/NOTES.md; go test ./core/role/... gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

