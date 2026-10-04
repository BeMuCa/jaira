---
id: 01M43C1R4XCMMNDAD75YMCTRTD
title: Der Dispatcher fragt in Prosa statt mit Auswahl
status: review
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Jede Frage des Dispatchers an den Menschen kommt als Auswahl aus 2–4 Optionen, Empfehlung zuerst — beantwortbar mit einem Klick statt mit einem Absatz."
context: |-
  Was falsch ist: Der Dispatcher stellt seine Fragen als Fließtext. Alex muss lesen, abwägen und frei antworten.
  Auslöser: Alex am 2026-10-04: 'der Dispatcher soll bei seinen Fragen immer Varianten zur Auswahl geben, interaktiv'.
  Wo gefragt wird: core/role/builtin/jaira-dispatcher/SKILL.md, Abschnitt 'Before the plan lane' Schritt 3 (offene Entscheidungen) und 'When to stop' (dritte Runde, ja/nein).
  Werkzeug: AskUserQuestion in Claude Code, 2–4 Optionen, Empfehlung als erste mit '(Recommended)'.
  Haken: Läuft der Dispatcher als Subagent, hat er AskUserQuestion evtl. nicht. Dann gehen die Optionen nummeriert in den Bericht, und der Teamlead stellt sie mit AskUserQuestion (core/role/builtin/jaira-teamlead/SKILL.md).
definition-of-done: "jaira-dispatcher/SKILL.md: ein Abschnitt, wie gefragt wird (AskUserQuestion, 2–4 Optionen, Empfehlung zuerst, Rückfall nummerierte Liste); beide Fragestellen verweisen darauf"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-04T11:51:06Z
updated-at: 2026-10-04T12:20:06Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-23999
claimed-at: 2026-10-04T11:51:17Z
outcome-what: "Getestet"
outcome-why: "Befunde eingearbeitet"
outcome-resolves: testing
review-summary: "Reiner Prompt-Text, zwei Commits. a5b445c: jaira-dispatcher/SKILL.md bekommt 'How you ask' (Z.164-185): jede Frage an den Menschen per AskUserQuestion, eine Frage pro Entscheidung, 2-4 Optionen mit einzeiligem Trade-off, Empfehlung zuerst mit '(Recommended)', keine 'other'-Option, Ja/Nein als Zwei-Optionen-Wahl. Ohne AskUserQuestion (Subagent/headless) gehen die Optionen als nummerierte Liste in den Bericht, der Dispatcher stoppt. Beide Fragestellen verweisen darauf (Z.53, Z.398-399). teamlead/SKILL.md Z.33-40: Fragen als AskUserQuestion, Dispatcher-Listen unverändert stellen. NOTES.md: eine Zeile unter '## Unreleased'. 976af20 schließt die drei Runde-1-Lücken: (1) teamlead Z.37-40: jede Antwort per 'jaira note <id>' aufs Ticket, bevor ein frischer Dispatcher startet; (2) dispatcher Z.178-181: im Rückfall gehen alle offenen Entscheidungen zusammen hoch, je eine Liste, in Frage-Reihenfolge; der Teamlead fragt einzeln und lässt durch frühere Antworten erledigte weg; (3) Z.182: Listen kommen zusätzlich zu den drei Berichtszeilen. go test ./core/role/ ./core/release/ grün."
review-gaps: "1) Rückfallweg verliert den Modus: Bisher setzte der fragende Dispatcher nach den Antworten mode=conversational (dispatcher/SKILL.md:57-61, Schritt 5). Im Rückfall stoppt er vorher; der Teamlead schreibt nur die Notes (teamlead/SKILL.md:37-40) und setzt den Modus nicht. Der frische Dispatcher zählt in Schritt 1 alle Entscheidungen als per Note geschlossen (Z.46-49), landet bei Schritt 2 'None open? ... run on as usual' (Z.50-51) und läuft autonom - kein --no-worktree, keine Diffs nach jedem DoD-Punkt, keine begleitende Critique. Z.47-48 deutet ('the mode is not yet set') nur an, dass er Schritt 5 nachholen müsste, sagt es aber nicht; vorher war das ein seltener Absturzfall, mit dem Rückfall ist es der Normalweg jedes Subagent-Dispatchers. Fix: ein Satz in teamlead Z.37-40 ('... dann jaira set <id> mode=conversational') oder in dispatcher Schritt 1/2 ('per Note beantwortete Entscheidungen zählen für Schritt 5 mit'). 2) Klein: der Subagent-Startprompt in teamlead/SKILL.md:69-70 sagt 'Report in three lines.' - widerspricht dispatcher Z.182 ('lists come on top of the three report lines'); ein Subagent nimmt die Zeile im eigenen Prompt im Zweifel wörtlicher als die Skill-Datei. Ergänzung '... plus any numbered option lists' schließt das."
test-verdict: |-
  go test ./core/role/ ./core/release/ grün. Verhalten (fragt der Dispatcher wirklich mit Auswahl?) nur in einem echten Lauf prüfbar.
  Runde 2: go test ./core/role/ ./core/release/ grün.
  Runde 3: go test ./core/role/ grün.
question: "Nach 'jaira roles install --global --force' einen Dispatcher auf ein Ticket mit offener Entscheidung ansetzen: kommt die Frage als Auswahl?"
review-verdict: "Zurück an in-progress. Die drei Runde-1-Lücken sind geschlossen (teamlead:37-40, dispatcher:178-182), neue Sätze widersprechen dem Rest nicht direkt, NOTES-Zeile ist eine Zeile unter '## Unreleased' und als Anweisung formuliert, Tests grün. Aber der jetzt geschlossene Rückfallweg hat eine Folgelücke: niemand setzt mode=conversational, und der neu gestartete Dispatcher fährt das Ticket mit beantworteten offenen Entscheidungen autonom (dispatcher:50-51 vs. 57-61). Sicher beim Befund; unsicher nur, ob Alex ihn in diesem Ticket statt in einem eigenen will - die Ursache ist aber der hier eingeführte Rückfall. Befund 2 (teamlead:69-70) ist klein und allein kein Rücksendegrund."
review-check: "1. cd /home/alex/projects/jaira && go test ./core/role/ ./core/release/ - beide Zeilen 'ok'. 2. sed -n '164,185p' core/role/builtin/jaira-dispatcher/SKILL.md - Abschnitt 'How you ask': AskUserQuestion, 2-4 Optionen, '(Recommended)', Rückfall mit nummerierten Listen, mehrere Entscheidungen zusammen, Listen zusätzlich zu den drei Zeilen. 3. sed -n '33,40p' core/role/builtin/jaira-teamlead/SKILL.md - Satz 'Each answer goes onto the ticket with jaira note'. Prüfen: steht dort auch 'jaira set <id> mode=conversational'? (heute: nein - Lücke 1). 4. sed -n '44,61p' core/role/builtin/jaira-dispatcher/SKILL.md - lesen als frischer Dispatcher, dessen offene Entscheidungen alle per Note beantwortet sind: Schritt 2 sagt 'run on as usual', Schritt 5 (Modus setzen) wird nie erreicht. 5. sed -n '67,71p' core/role/builtin/jaira-teamlead/SKILL.md - Subagent-Prompt sagt 'Report in three lines.' ohne die Listen (Lücke 2). 6. sed -n '16,18p' core/release/NOTES.md - genau eine '- '-Zeile unter '## Unreleased'. 7. Echter Lauf nach 'jaira roles install --global --force': Teamlead ohne Herdr, Dispatcher als Subagent auf ein Ticket mit offener Entscheidung - Erwartung: Frage kommt als Auswahl, Antwort steht danach als Note (jaira show <id> --json), und 'mode' ist 'conversational' (heute vermutlich leer)."
---

# Der Dispatcher fragt in Prosa statt mit Auswahl

## Definition of Done

- [x] jaira-dispatcher/SKILL.md: ein Abschnitt, wie gefragt wird (AskUserQuestion, 2–4 Optionen, Empfehlung zuerst, Rückfall nummerierte Liste); beide Fragestellen verweisen darauf
  proof: core/role/builtin/jaira-dispatcher/SKILL.md:164 Abschnitt 'How you ask'; Verweise core/role/builtin/jaira-dispatcher/SKILL.md:53 und core/role/builtin/jaira-dispatcher/SKILL.md:394
- [x] jaira-teamlead/SKILL.md: Fragen eines Dispatchers werden als AskUserQuestion weitergegeben, nicht als Text
  proof: core/role/builtin/jaira-teamlead/SKILL.md:33-40
- [x] go test ./core/role/ grün
  proof: go test ./core/role/ ok
- [x] Zeile in core/release/NOTES.md unter Unreleased
  proof: core/release/NOTES.md:17

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-04 11:52 · Alexander Sacharov** — Mehrere offene Entscheidungen werden NICHT gebündelt in einem AskUserQuestion-Aufruf gestellt: Schritt 3 verlangt eine Frage nach der anderen, weil spätere Antworten von früheren abhängen können. Eine Option 'Other' wird nicht angelegt, Claude Code bietet freie Eingabe ohnehin an.
- **2026-10-04 12:14 · Alexander Sacharov** — Review schickt zurück: im Rückfallweg (Dispatcher ohne AskUserQuestion, z.B. als Subagent) stoppt der Dispatcher mit nummerierter Liste, der Teamlead fragt - aber teamlead/SKILL.md:33-37 sagt ihm nicht, die Antwort per 'jaira note <id>' aufs Ticket zu schreiben und den Dispatcher fortzusetzen/neu zu starten. Ohne Note fragt ein neuer Dispatcher in Schritt 1 (dispatcher/SKILL.md:44-49) dieselbe Frage erneut, und die Entscheidung fehlt auf dem Board. Nebenbei: bei mehreren Entscheidungen im Rückfall klären, ob eine oder alle Listen hochgehen (Schritt 3: eine nach der anderen); Rückfall-Liste vs. 'Three lines' im Report (Z.414).
- **2026-10-04 12:15 · Alexander Sacharov** — Review-Lücken behoben: (1) teamlead SKILL.md:37 — Antwort per jaira note aufs Ticket, bevor ein frischer Dispatcher startet; (2) dispatcher SKILL.md:178 — im Rückfallweg gehen alle offenen Entscheidungen zusammen hoch, der Teamlead fragt sie einzeln in Reihenfolge; (3) dieselbe Stelle — die Listen kommen zusätzlich zu den drei Berichtszeilen.
- **2026-10-04 12:19 · Alexander Sacharov** — Review Runde 2 schickt zurück: Runde-1-Lücken geschlossen, aber im Rückfallweg setzt niemand mode=conversational. Der Teamlead schreibt nur die Notes (teamlead/SKILL.md:37-40); der frische Dispatcher zählt alles als geschlossen, nimmt Schritt 2 'run on as usual' (dispatcher/SKILL.md:50-51) und überspringt Schritt 5 (Z.57-61) - das Ticket läuft autonom, obwohl es offene Entscheidungen hatte. Fix: Teamlead setzt den Modus nach den Notes, oder Schritt 1/2 sagt, dass per Note beantwortete Entscheidungen für Schritt 5 zählen. Nebenbei: Subagent-Prompt teamlead/SKILL.md:69-70 'Report in three lines.' um die Optionslisten ergänzen.
- **2026-10-04 12:19 · Alexander Sacharov** — Review Runde 2 behoben: dispatcher SKILL.md:50-55 — Schritt 2 setzt mode=conversational, wenn Notizen Entscheidungen beantwortet haben und mode leer ist (Antwort kam über den Teamlead oder vor dem Tod eines Dispatchers). teamlead SKILL.md:69-71 — der Subagent-Prompt nennt die Optionslisten neben den drei Zeilen.
