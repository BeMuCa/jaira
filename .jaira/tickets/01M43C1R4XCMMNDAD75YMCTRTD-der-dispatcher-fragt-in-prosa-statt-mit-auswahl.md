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
updated-at: 2026-10-04T14:00:09Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-23999
claimed-at: 2026-10-04T11:51:17Z
outcome-what: "Getestet"
outcome-why: "Vereinfachung umgesetzt"
outcome-resolves: testing
review-summary: "Reiner Prompt-Text, drei Commits. a5b445c: jaira-dispatcher/SKILL.md bekommt 'How you ask' (Z.168-189): jede Frage per AskUserQuestion, eine pro Entscheidung, 2-4 Optionen, Empfehlung zuerst mit '(Recommended)', kein 'other', Ja/Nein als Zwei-Optionen-Wahl; ohne AskUserQuestion gehen die Optionen nummeriert in den Bericht und der Dispatcher stoppt. Beide Fragestellen verweisen darauf (Z.57, Z.402-403). teamlead/SKILL.md Z.33-37: Fragen als Auswahl, Dispatcher-Listen unverändert stellen. NOTES.md: eine Zeile unter Unreleased. 976af20: Rückfall bündelt alle offenen Entscheidungen in Frage-Reihenfolge, Listen zusätzlich zu den drei Berichtszeilen (dispatcher Z.182-186); Teamlead schreibt jede Antwort per jaira note aufs Ticket, bevor ein frischer Dispatcher startet (teamlead Z.37-40). cab71ad (Runde-2-Lücken): dispatcher Schritt 2 (Z.50-55) — sind Entscheidungen durch Notes geschlossen und mode leer, setzt der Dispatcher mode=conversational (Schritt 5) und läuft weiter; teamlead Subagent-Prompt (Z.68-72) verlangt die nummerierten Optionen zusätzlich zu den drei Zeilen. Beide Runde-2-Lücken sind damit geschlossen; go test ./core/role/ ./core/release/ grün."
review-gaps: "1) dispatcher/SKILL.md:51-55 — die neue Ausnahme greift bei JEDER Note, die eine Entscheidung schließt, nicht nur bei der Antwort eines Menschen. Die brainstorm-Lane schreibt aber per Prompt 'what you would do, and why' als Note, und jaira-role-lane/SKILL.md:71 verlangt 'why this and not that'-Notes. Ein Dispatcher, der vor der Plan-Lane solche Agenten-Notes liest, zählt die Entscheidung als 'closed in step 1' und schaltet mode=conversational ein — obwohl kein Mensch entschieden hat. Folge: ein autonomer Ticket-Lauf startet Worker mit --no-worktree, Worker committen nicht mehr und warten auf einen Menschen, der nicht mitliest. Die Begründung im selben Satz ('a person decided the shape') stimmt dann nicht. Fix: Bedingung auf Notes einschränken, die die Antwort einer Person festhalten (z.B. Teamlead/Dispatcher schreibt die Antwort als 'Person entschied: ...'), und das in teamlead Z.38 und dispatcher Schritt 4 als Form der Note festlegen. Sonst keine Widersprüche gefunden; Querverweise stimmen."
test-verdict: |-
  go test ./core/role/ ./core/release/ grün. Verhalten (fragt der Dispatcher wirklich mit Auswahl?) nur in einem echten Lauf prüfbar.
  Runde 2: go test ./core/role/ ./core/release/ grün.
  Runde 3: go test ./core/role/ grün.
  Runde 4: go test ./core/role/ ./core/release/ grün.
question: "Nach 'jaira roles install --global --force' einen Dispatcher auf ein Ticket mit offener Entscheidung ansetzen: kommt die Frage als Auswahl?"
review-verdict: "Zurück an in-progress. Die beiden Runde-2-Lücken sind geschlossen (Modus im Rückfallweg, Subagent-Prompt), Tests grün, NOTES-Zeile regelkonform. Aber die neue Ausnahme in Schritt 2 (dispatcher Z.51-55) unterscheidet nicht zwischen der Antwort eines Menschen und einer Agenten-Note (brainstorm-Empfehlung, 'why this and not that') und kann so autonome Tickets ungewollt in den conversational-Modus kippen. Ein Satz behebt es. Sicher beim Befund; unsicher, wie oft er in der Praxis auftritt — Alex mag ihn als nicht-blockierend werten."
review-check: "1. cd /home/alex/projects/jaira && go test ./core/role/ ./core/release/ — zweimal 'ok'. 2. sed -n '44,65p' core/role/builtin/jaira-dispatcher/SKILL.md — Schritt 2 enthält die Ausnahme 'if notes answered decisions ... and mode is still empty ... Do step 5 now'. Prüfen: steht dort, dass die Note die Antwort einer PERSON sein muss? (heute: nein — das ist die Lücke). 3. sed -n '168,189p' derselben Datei — 'How you ask': AskUserQuestion, 2-4 Optionen, '(Recommended)', Rückfall nummerierte Listen in Frage-Reihenfolge, zusätzlich zu den drei Zeilen. 4. sed -n '33,40p;65,73p' core/role/builtin/jaira-teamlead/SKILL.md — Teamlead stellt die Listen als Auswahl, schreibt jede Antwort per 'jaira note <id>' vor dem neuen Dispatcher; Subagent-Prompt verlangt 'plus the numbered options'. 5. sed -n '16,18p' core/release/NOTES.md — genau eine '- '-Zeile unter '## Unreleased'. 6. Echter Lauf: 'jaira roles install --global --force', Teamlead ohne Herdr, Dispatcher als Subagent auf ein Ticket mit offener Entscheidung — erwarten: Frage kommt als Auswahl, danach 'jaira show <id> --json' zeigt die Antwort als Note und nach Start des frischen Dispatchers 'mode: conversational'. 7. Gegenprobe: Ticket, dessen brainstorm-Lane nur eine Empfehlungs-Note schrieb, per Dispatcher starten — erwarten: mode bleibt leer (heute vermutlich nicht)."
---

# Der Dispatcher fragt in Prosa statt mit Auswahl

## Definition of Done

- [x] jaira-dispatcher/SKILL.md: ein Abschnitt, wie gefragt wird (AskUserQuestion, 2–4 Optionen, Empfehlung zuerst, Rückfall nummerierte Liste); beide Fragestellen verweisen darauf
  proof: core/role/builtin/jaira-dispatcher/SKILL.md:164 Abschnitt 'How you ask'; Verweise core/role/builtin/jaira-dispatcher/SKILL.md:53 und core/role/builtin/jaira-dispatcher/SKILL.md:394
- [x] jaira-teamlead/SKILL.md: Fragen eines Dispatchers werden als AskUserQuestion weitergegeben, nicht als Text
  proof: core/role/builtin/jaira-teamlead/SKILL.md:33-42
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
- **2026-10-04 12:21 · Alexander Sacharov** — Review Runde 3 schickt zurück: Runde-2-Lücken geschlossen. Neu: dispatcher/SKILL.md:51-55 setzt mode=conversational, sobald irgendeine Note eine Entscheidung schließt — auch eine Agenten-Note (brainstorm-Prompt verlangt 'what you would do, and why'; jaira-role-lane/SKILL.md:71 'why this and not that'). Dann kippt ein autonomer Lauf in den conversational-Modus, ohne dass ein Mensch entschieden hat. Fix: Ausnahme nur für Notes, die die Antwort einer Person festhalten, und diese Form in teamlead Z.38 / dispatcher Schritt 4 festlegen.
- **2026-10-04 13:59 · Alexander Sacharov** — Alex hat entschieden (2026-10-04, Auswahl nach drei Review-Runden): vereinfachen. Die Ausnahme in Schritt 2 des Dispatchers entfällt; der Teamlead, der die Antwort bekommt, schreibt die Note UND setzt mode=conversational selbst. Kein Raten aus Notizen.
