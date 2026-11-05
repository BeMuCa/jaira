---
id: 01M11YYQZQ4HW090F77EFCMP17
title: Lanes geben Artefakte aus und reichen sie an die naechste Lane weiter
status: backlog
ready: false
creator: BeMuCa
goal: "Eine Lane kann benannte Dateien erzeugen, die gespeichert werden und der naechsten Lane als Eingabe uebergeben werden"
context: |-
  Berk am 27.08.: Lanes sollen Artefakte ausgeben koennen. Ablage unter .jaira/artifacts/ mit zwei Ordnern: einer gitignoriert (local/), einer geteilt und nach dem GitHub-Handle benannt (bemuca/). Die Lane deklariert die Namen ihrer Artefakte; die Folge-Lane deklariert sie als Eingabe und 'jaira show --for-lane' reicht sie mit.

  Was heute da ist: nichts davon. input-requires und output-produces auf einer Lane meinen Ticket-Felder (plan, outcome-what, review-summary), keine Dateien. Der einzige Ort fuer Text, der nicht ins Frontmatter passt, ist der Ticket-Body (Plan, Progress).

  Offen vor dem Bauen: (1) Groesse und Format - alles Markdown/Text, oder auch Binaerdateien? (2) Benennung: pro Ticket (artifacts/<handle>/<ticket>-<name>.md) damit zwei Tickets nicht kollidieren. (3) Was 'weitergeben' heisst: Pfad im --for-lane JSON reicht fuer eine Datei; Inhalt inline sprengt die begrenzte Eingabe, die die Lane absichtlich hat. (4) Ob das geteilte Verzeichnis dieselbe Sicherheitsfrage aufwirft wie shared/ Lanes: ein committetes Artefakt ist fremder Text, der als Eingabe in ein Modell laeuft.
definition-of-done: "Eine Lane deklariert output-artifacts; nach dem Move liegt die Datei unter .jaira/artifacts/; die Folge-Lane deklariert sie als input-artifacts und 'jaira show --for-lane' nennt den Pfad; local/ ist gitignoriert, der Handle-Ordner nicht; ein Test deckt Ablage und Weitergabe ab"
blocked-by: []
commits: []
created-at: 2026-08-27T15:55:56Z
updated-at: 2026-11-04T21:29:21Z
updated-by: berkc
---

# Lanes geben Artefakte aus und reichen sie an die naechste Lane weiter

## Definition of Done

- [ ] Eine Lane deklariert output-artifacts; nach dem Move liegt die Datei unter .jaira/artifacts/; die Folge-Lane deklariert sie als input-artifacts und 'jaira show --for-lane' nennt den Pfad; local/ ist gitignoriert, der Handle-Ordner nicht; ein Test deckt Ablage und Weitergabe ab

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-04 21:29 · berkc** — Follow-up of ENCFVZ (skill/script/receipt lanes): the maintainer listed output-artifacts / input-artifacts as lane attributes on 04.11.2026 and left them to this ticket. Refinements from ENCFVZ's brainstorm note of 30.10.2026: path per ticket under .jaira/artifacts/<handle>/<ticket>/<lane>-<name>.<ext>, a gitignored local/ for large or private files, size limit for shared ones (proposed 1 MB), --for-lane hands paths not contents, the gate checks output-artifacts exist on leaving.
