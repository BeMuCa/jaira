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
updated-at: 2026-10-04T12:16:06Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-23999
claimed-at: 2026-10-04T11:51:17Z
outcome-what: "Getestet"
outcome-why: "Review-Befunde eingearbeitet"
outcome-resolves: testing
review-summary: |-
  Reiner Prompt-Text. jaira-dispatcher/SKILL.md bekommt den Abschnitt 'How you ask' (Z.164-181): jede Frage an den Menschen per AskUserQuestion, eine Frage pro Entscheidung, 2-4 Optionen mit einzeiligem Trade-off, Empfehlung zuerst mit '(Recommended)', keine 'other'-Option, Ja/Nein als Zwei-Optionen-Wahl; ohne AskUserQuestion (Subagent/headless) gehen die Optionen nummeriert in den Bericht und der Dispatcher stoppt. Beide Fragestellen verweisen darauf (Z.53 Schritt 3, Z.393-394 dritte Runde). jaira-teamlead/SKILL.md Z.33-37: Fragen kommen als AskUserQuestion, eine nummerierte Liste vom Dispatcher wird unverändert gestellt. NOTES.md: eine Zeile unter '## Unreleased', als Anweisung formuliert, nennt 'jaira roles install --global --force'. Querverweise ('below'/'above'/'step 4 above') stimmen; go test ./core/role/ ./core/release/ grün.
  Runde 2 nach Review: Ergänzungen widersprechen nichts; still.
review-gaps: "1) Rückfallweg verliert die Antwort: dispatcher/SKILL.md:180-181 verlangt 'jaira note' vor dem Handeln, aber im Rückfall hat der Dispatcher bereits gestoppt - die Antwort bekommt der Teamlead, und teamlead/SKILL.md:33-37 sagt ihm weder, die Antwort per 'jaira note <id>' auf das Ticket zu schreiben, noch, den Dispatcher danach weiterzuführen (SendMessage) bzw. neu zu starten. Folge: Entscheidung steht nirgends auf dem Board (Core Value), und ein neuer Dispatcher zählt sie in Schritt 1 (Z.44-49) erneut als offen und fragt dieselbe Frage nochmal. 2) Mehrere offene Entscheidungen im Rückfall: Z.175-177 schickt 'die Optionen' als eine Liste hoch, Schritt 3 (Z.52-53) verlangt eine Frage nach der anderen, weil Antworten voneinander abhängen können - offen, ob der Dispatcher alle Listen auf einmal liefert oder nur die erste. Klein. 3) Klein: der Rückfall-Bericht mit nummerierter Liste sprengt 'Three lines at the end' (dispatcher/SKILL.md:414-422); ein Satz, dass die Liste zur dritten Zeile gehört, würde den Widerspruch schließen."
test-verdict: |-
  go test ./core/role/ ./core/release/ grün. Verhalten (fragt der Dispatcher wirklich mit Auswahl?) nur in einem echten Lauf prüfbar.
  Runde 2: go test ./core/role/ ./core/release/ grün.
question: "Nach 'jaira roles install --global --force' einen Dispatcher auf ein Ticket mit offener Entscheidung ansetzen: kommt die Frage als Auswahl?"
review-verdict: "Zurück an in-progress. DoD formal erfüllt, keine Widersprüche im Hauptweg, Querverweise korrekt, NOTES-Zeile regelkonform, Tests grün. Aber der neue Rückfallweg (Dispatcher als Subagent) hat ein Loch: niemand ist angewiesen, die Antwort des Menschen auf das Ticket zu schreiben und den Dispatcher fortzusetzen - genau der Fall, den Schritt 4 verhindern soll. Fix ist ein Satz in teamlead/SKILL.md:33-37. Bin mir sicher beim Befund, unsicher nur, ob Alex ihn für diese Runde als blockierend sieht."
review-check: "1. cd /home/alex/projects/jaira && go test ./core/role/ ./core/release/ - beide 'ok'. 2. sed -n '164,181p' core/role/builtin/jaira-dispatcher/SKILL.md - Abschnitt 'How you ask' mit AskUserQuestion, 2-4 Optionen, '(Recommended)', Rückfall nummerierte Liste. 3. sed -n '50,56p;390,397p' derselben Datei - beide Stellen verweisen auf 'How you ask'. 4. sed -n '33,37p' core/role/builtin/jaira-teamlead/SKILL.md - AskUserQuestion-Satz; prüfen: steht dort, dass der Teamlead die Antwort mit 'jaira note' aufs Ticket schreibt? (heute: nein - das ist die Lücke). 5. sed -n '16,18p' core/release/NOTES.md - genau eine '- '-Zeile unter '## Unreleased'. 6. Echter Lauf: 'jaira roles install --global --force', dann Teamlead ohne Herdr einen Dispatcher (Subagent) auf ein Ticket mit offener Entscheidung ansetzen - Erwartung: Teamlead stellt die Frage als Auswahl mit Empfehlung zuerst; danach 'jaira show <id> --json' - steht die Antwort als Note drauf?"
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
