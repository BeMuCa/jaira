---
id: 01M2HS9CJZVKSH7352DT179VBX
title: "Die PR-Rolle kann nur GitHub, obwohl ein Board schon auf GitLab liegt"
status: backlog
ready: true
creator: Alexander Sacharov
goal: "Die Rolle, die fertige Arbeit zur Abgabe bringt, funktioniert auf einem GitLab-Repository genauso wie auf GitHub - ohne dass jemand den ausgelieferten Prompt von Hand umschreibt."
context: |-
  core/role/builtin/jaira-role-pr/SKILL.md ruft an fuenf Stellen 'gh' auf: :22 'gh pr list --head ... --state open', :72 'gh pr create --title ... --body-file ...', und :91-92 die drei Verbote 'gh pr create', 'gh pr merge', 'gh pr review --approve'.

  Auf GitLab gibt es keinen dieser Befehle. Das Gegenstueck heisst 'glab mr list --source-branch' und 'glab mr create'. Auch das Wort stimmt nicht: dort heisst es Merge Request, und der Prompt sagt durchgehend pull request.

  Das ist kein hypothetischer Fall. Auf diesem Rechner liegt unter /home/alex/projects/requirementsgenie ein Board mit .jaira/, und sein einziger Remote ist git.esprit-engineering.de - GitLab. glab ist installiert, sowohl als Windows-Datei als auch als Linux-Binary in ~/.local/bin/glab. Wer dort die PR-Rolle startet, bekommt einen Prompt, dessen erste Handlung ins Leere greift.

  jaira selbst ist davon nicht betroffen: es reist auf git-Refs und kennt keinen Forge. Betroffen ist allein der Rollen-Prompt.

  Was NICHT aufgeweicht werden darf, egal wie die Loesung aussieht: die Regel aus 13VMA8 gilt auf beiden Seiten gleich - der Agent schiebt den Zweig und haelt an, das Aufmachen gibt der Mensch frei, und Mergen und Freigeben bleiben verboten. Ein zweiter Weg durch den Prompt ist eine zweite Stelle, an der diese Regel stehen und stimmen muss.

  Offen und Sache der Plan-Lane: woran der Forge erkannt wird. Die Remote-URL ist der naheliegende Weg und braucht keine Einstellung; eine Einstellung je Board waere der andere; zwei getrennte Rollen der dritte und vermutlich der schlechteste, weil die Regel dann doppelt gepflegt werden muss.
definition-of-done: "Auf einem Repository, dessen Remote auf GitLab zeigt, listet die Rolle die offenen Merge Requests des aktuellen Zweigs und schreibt dem Menschen eine lauffaehige 'glab mr create'-Zeile aus. Nachgestellt auf dem requirementsgenie-Board oder einem Fixture mit einem GitLab-Remote."
tags:
  - cli
blocked-by: []
related:
  - 01M28FMQGC4CNQT8Z9WY13VMA8
commits: []
created-at: 2026-09-15T05:40:26Z
updated-at: 2026-09-15T05:45:28Z
assignee: "Alexander Sacharov"
updated-by: Alexander Sacharov
---

# Die PR-Rolle kann nur GitHub, obwohl ein Board schon auf GitLab liegt

## Definition of Done

- [-] Auf einem Repository, dessen Remote auf GitLab zeigt, listet die Rolle die offenen Merge Requests des aktuellen Zweigs und schreibt dem Menschen eine lauffaehige 'glab mr create'-Zeile aus. Nachgestellt auf dem requirementsgenie-Board oder einem Fixture mit einem GitLab-Remote.
  proof: als Punkte 2-4 nach 13VMA8 gewandert, 2026-09-15
- [ ] Auf GitHub aendert sich nichts: derselbe Prompt tut dort weiterhin genau das, was er heute tut. Nachgestellt auf diesem Repository.
- [ ] Der Forge wird erkannt, ohne dass jemand pro Board etwas einstellen muss, solange der Remote es hergibt - und wenn er es nicht hergibt, sagt die Rolle das, statt den falschen Befehl zu raten.
- [ ] Die Regel aus 13VMA8 steht auf beiden Wegen und stimmt: der Agent schreibt die Zeile aus und fuehrt sie nie aus, merged nie und gibt nie frei. Nachgestellt, indem beide Wege gelesen werden - kein 'Never run' darf auf einem Weg fehlen.
- [ ] Eine Zeile in core/release/NOTES.md unter ## Unreleased.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-15 05:45 · Alexander Sacharov** — Am 2026-09-15 von Alex in 13VMA8 zusammengelegt und hier stillgelegt. Die Kriterien stehen jetzt dort als Punkte 2 bis 4, in seiner Formulierung: welches Werkzeug laeuft, soll waehlbar sein und nicht nur geraten.

Grund fuer das Zusammenlegen: es ist dieselbe Datei. core/role/builtin/jaira-role-pr/SKILL.md hat gerade vier Runden Ueberarbeitung hinter sich; ein eigenes Ticket haette spaeter um dieselben Zeilen gekaempft.
