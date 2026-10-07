---
id: 01M4ATQ46NZ08TFVQARZ9WSS49
title: "Die README nennt jede Rolle in einer Zeile, und das Tutorial liegt im Repo"
status: testing
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
updated-at: 2026-10-07T09:39:09Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-3654451
claimed-at: 2026-10-07T09:38:39Z
outcome-what: "README: Tabelle mit einer Zeile pro Skill und Link auf die Tutorial-Seite; Tutorial-Seite unter docs/, quirks-fest, ohne toten CSS"
outcome-why: "Berk wusste nicht, welche Skills es gibt und wie man sie anspricht; critique fand die Seite unverlinkt, optimize zwei wirkungslose CSS-Zeilen"
outcome-resolves: "DoD 1: Tabelle vor den bestehenden Absaetzen, die bleiben. DoD 2: committete Datei = Artifact-Version 5"
executed-by: opus
review-summary: none
review-gaps: "Looked at README.md:623-638,662 and all of docs/how-to-work-with-jaira.html; removed nothing. Duplication: no other file lists every role (grep README, docs/, AGENTS.md); README table vs the four detailed bullets and the German page vs the README table left as critique ruled. Dead code: every CSS class and custom property in the page is used. Left two inert CSS lines on purpose: how-to-work-with-jaira.html:177 prefers-reduced-motion rule (the page declares no transition or animation; may still act on whatever the artifact host injects) and :95 ol.chain gap: 0 (grid default). Removing either breaks DoD 2 - the file must equal artifact version 4 - and costs a republish this lane cannot do, for no visible change. Fluff/cost: none; copy-button JS fallback covers a real case (clipboard missing or refused), fonts load once."
---

# Die README nennt jede Rolle in einer Zeile, und das Tutorial liegt im Repo

## Definition of Done

- [x] README.md listet im Rollen-Abschnitt jeden Skill (jaira plus die neun aus 'jaira roles list') mit genau einer Zeile, was er tut; die bestehenden ausfuehrlichen Absaetze bleiben
  proof: README.md:623-634 Tabelle Skill|What it does, 10 Zeilen (jaira + die neun aus 'jaira roles list' von 0.3.5); Absaetze darunter unveraendert
- [x] docs/how-to-work-with-jaira.html ist committet und deckt sich mit der veroeffentlichten Seite
  proof: docs/how-to-work-with-jaira.html; dieselbe Datei als Version 5 nach https://claude.ai/artifact/FAnMJiZz6dG2hb3it2GM2c veroeffentlicht

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
- **2026-10-07 09:36 · BeMuCa** — critique pass 2 (read only round 1's three findings and a02ca51): all three closed, nothing new. L1 README 'jaira roles list shows the roles the binary you run carries' and L2 link line after the table (README.md:636-638) as proposed. L3 closed by the different fix in the 09:32 note: both harms the finding named are gone - table { font: inherit } undoes the quirks-mode table reset of font/line-height (font-size set right after, text-align set on th/td), and lang=de on .wrap covers every visible element incl. the JS copy button. Keeping the file in the publisher's input format (the Artifact host wraps the page in its own skeleton) is a defensible reason not to add a doctype; not reopened. DoD 2 republish taken from the 09:32 note, not re-checked here.
- **2026-10-07 09:38 · BeMuCa** — Nach optimize: die zwei wirkungslosen CSS-Zeilen, die optimize nur wegen des Neu-Veroeffentlichens stehen liess (gap: 0 auf ol.chain, reduced-motion-Regel ohne Transitions auf der Seite), entfernt und als Version 5 veroeffentlicht. Ein Republish kostet einen Aufruf; kein Grund, toten Code zu behalten.
