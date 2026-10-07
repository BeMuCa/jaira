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
updated-at: 2026-10-07T09:32:39Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-3634832
claimed-at: 2026-10-07T09:31:48Z
outcome-what: "README: Tabelle mit einer Zeile pro Skill, Link auf die Tutorial-Seite, 'names the rest' umformuliert; Tutorial-Seite unter docs/, quirks-fest"
outcome-why: "Berk wusste nicht, welche Skills es gibt und wie man sie anspricht; critique Runde 1 fand die Seite nicht verlinkt und den alten Hinweis veraltet"
outcome-resolves: "DoD 1: Tabelle vor den bestehenden Absaetzen, die bleiben. DoD 2: committete Datei = Artifact-Version 4. Die drei Befunde aus Runde 1 sind in der Note vom Implementierer beantwortet"
executed-by: opus
review-summary: |-
  README.md:658 "`jaira roles list` names the rest: single-lane workers, a tester, ..." is obsolete now that the table above names every role (it is a pointer, not one of the detailed paragraphs the DoD keeps) - reword it to "`jaira roles list` shows the roles the binary you run carries."
  README.md:635 nothing in the repo links to docs/how-to-work-with-jaira.html (git grep finds no reference), so the newcomer the goal names never finds it - add one line after the table in the pattern of README.md:680, e.g. "A walk-through in German, which route fits which task: **[docs/how-to-work-with-jaira.html](docs/how-to-work-with-jaira.html)** (open it in a browser; GitHub shows the source)."
  docs/how-to-work-with-jaira.html:1 is the artifact host's input format: it starts at <meta charset> with no doctype and no <html lang>, which the claude.ai host adds but a file opened from a clone does not get - it renders in quirks mode (table line-height resets) and is announced as English. Put <!doctype html> and <html lang="de"> above line 1; a host that wraps it ignores the stray doctype and merges lang onto its own <html>, so republish the same file to keep DoD 2.
---

# Die README nennt jede Rolle in einer Zeile, und das Tutorial liegt im Repo

## Definition of Done

- [x] README.md listet im Rollen-Abschnitt jeden Skill (jaira plus die neun aus 'jaira roles list') mit genau einer Zeile, was er tut; die bestehenden ausfuehrlichen Absaetze bleiben
  proof: README.md:623-634 Tabelle Skill|What it does, 10 Zeilen (jaira + die neun aus 'jaira roles list' von 0.3.5); Absaetze darunter unveraendert
- [x] docs/how-to-work-with-jaira.html ist committet und deckt sich mit der veroeffentlichten Seite
  proof: docs/how-to-work-with-jaira.html; dieselbe Datei als Version 4 nach https://claude.ai/artifact/FAnMJiZz6dG2hb3it2GM2c veroeffentlicht

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-07 09:31 · BeMuCa** — critique (pass 1, whole diff of 6ab980c read): three findings, fix is clear for each.
1) README.md:658 bullet "`jaira roles list` names the rest: ..." is obsolete - the new table names every role. It is a pointer, not a detailed paragraph the DoD protects. Reword to "`jaira roles list` shows the roles the binary you run carries."
2) Nothing links to docs/how-to-work-with-jaira.html (git grep: no reference outside .jaira). Goal says a newcomer finds it. Add one English line after the table (README.md:635), in the style of README.md:680, saying the page is German and must be opened in a browser (GitHub shows .html as source).
3) docs/how-to-work-with-jaira.html:1 starts at <meta charset>: no <!doctype html>, no <html lang="de">. claude.ai adds both when publishing; a file opened from a clone gets neither -> quirks mode + wrong language. Add both above line 1, then republish the same file to the artifact (FAnMJiZz6dG2hb3it2GM2c) so DoD 2 still holds. HTML parser ignores a doctype inside an already-open document and merges lang onto the existing <html>, so the published copy is unaffected.
Let stand on purpose: README table repeating the four detailed bullets (DoD keeps the bullets); German page duplicating the README table (person asked for both); the data-theme CSS that only the artifact host uses (needed for the published copy, which must match); hand-maintained page drifting by version (footer names the refresh steps). Content correctness (e.g. jaira pull, cp of SKILL.md) not judged here - that is review/testing.
- **2026-10-07 09:32 · BeMuCa** — Antwort auf critique Runde 1: L1 (README 'names the rest' umformuliert) und L2 (Link auf die Seite hinter der Tabelle) wie vorgeschlagen. L3 anders geloest: kein <!doctype>/<html lang> in die Datei - die Artifact-Vorgabe des Publishers verbietet eigene doctype/html/head/body-Tags ausdruecklich, er wrappt selbst. Stattdessen die sichtbare Quirks-Folge behoben: table { font: inherit } (Quirks setzt Tabellen-Schrift auf initial) und lang=de auf dem Inhalts-Wrapper. Datei als Version 4 veroeffentlicht, DoD 2 haelt.
