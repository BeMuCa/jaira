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
updated-at: 2026-10-04T15:15:35Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-23999
claimed-at: 2026-10-04T11:51:17Z
outcome-what: "Getestet"
outcome-why: "Befund Runde 4 eingearbeitet"
outcome-resolves: testing
review-summary: "Reiner Prompt-Text. dispatcher/SKILL.md: Abschnitt 'How you ask' (Z.164-185) - jede Frage per AskUserQuestion, 2-4 Optionen, Empfehlung zuerst mit '(Recommended)', Ja/Nein als Zwei-Optionen-Wahl; ohne AskUserQuestion gehen alle offenen Entscheidungen als nummerierte Listen in Fragereihenfolge zusätzlich zu den drei Berichtszeilen hoch, und der Dispatcher stoppt. Schritt 3 (Z.53) und dritte Runde (Z.398-399) verweisen darauf. teamlead/SKILL.md Z.33-42: Teamlead stellt die Listen unverändert als Auswahl, schreibt jede Antwort per 'jaira note' und setzt dann selbst 'mode=conversational', bevor ein frischer Dispatcher startet; Subagent-Prompt Z.70-75 verlangt die Optionen neben den drei Zeilen. Runde 4 (225eb82): die Ausnahme in Dispatcher-Schritt 2, die den Modus aus Notes ableitete, ist entfernt (Z.50-51 wieder wie vor 8b0bb42). NOTES.md: eine '- '-Zeile unter '## Unreleased'. go test ./core/role/ ./core/release/ grün."
review-gaps: "1) Runde-3-Lücke ist zu (Dispatcher leitet den Modus nicht mehr aus Notes ab), aber dieselbe Fehlwirkung sitzt jetzt beim Teamlead: teamlead/SKILL.md:37-39 setzt mode=conversational nach JEDER weitergereichten Antwort. Der Dispatcher reicht über 'How you ask' (dispatcher/SKILL.md:175-182, Verweis Z.398-399) auch die Ja/Nein-Frage der dritten Runde als Liste hoch, und der Subagent-Prompt (teamlead/SKILL.md:72-74) fordert 'any question you could not ask'. Realistischer Lauf ohne Herdr: Testing schickt dreimal zurück -> Teamlead fragt 'vierte Runde?' -> Note + mode=conversational -> der neue Dispatcher fährt den Rest mit --no-worktree und Commit-Zeilen statt Commits, ohne dass jemand den Ticket-Zuschnitt entschieden hat. Fix: Modus nur setzen, wenn die Antworten offene Entscheidungen vor der Plan-Lane schließen; bei der Dritte-Runde-Frage nur die Note. Nicht blockierend: a) teamlead/SKILL.md:42 'nobody else sets it' stimmt nicht ganz - der Dispatcher setzt ihn selbst in Schritt 5 (dispatcher/SKILL.md:57-61), wenn er selbst gefragt hat. b) dispatcher/SKILL.md:46-49 (vorbestehend) beschreibt den Fall 'gestorben zwischen Schritt 4 und 5, Modus noch nicht gesetzt', aber nach Wegfall der Ausnahme setzt ihn dort niemand mehr - das Ticket läuft autonom weiter; galt schon vor diesem Ticket."
test-verdict: |-
  go test ./core/role/ ./core/release/ grün. Verhalten (fragt der Dispatcher wirklich mit Auswahl?) nur in einem echten Lauf prüfbar.
  Runde 2: go test ./core/role/ ./core/release/ grün.
  Runde 3: go test ./core/role/ grün.
  Runde 4: go test ./core/role/ ./core/release/ grün.
  Runde 5: go test ./core/role/ ./core/release/ grün.
question: "Nach 'jaira roles install --global --force' einen Dispatcher auf ein Ticket mit offener Entscheidung ansetzen: kommt die Frage als Auswahl?"
review-verdict: "Zurück an in-progress. DoD erfüllt, Runde-3-Befund geschlossen, Tests grün, NOTES-Zeile regelkonform. Aber die Vereinfachung verschiebt das Problem: die Regel 'Teamlead setzt mode=conversational' (teamlead/SKILL.md:37-39) ist unbedingt und greift damit auch bei der Dritte-Runde-Ja/Nein-Frage, die über denselben Rückfallweg kommt. Fix ist ein Halbsatz Einschränkung. Sicher beim Befund; ob Alex den Dritte-Runde-ohne-Herdr-Fall für häufig genug hält, um zu blockieren, ist seine Entscheidung."
review-check: "1. cd /home/alex/projects/jaira && go test ./core/role/ ./core/release/ - beide Zeilen 'ok'. 2. sed -n '44,61p' core/role/builtin/jaira-dispatcher/SKILL.md - Schritt 2 endet bei 'needs nobody.', keine Ausnahme mehr über Notes. 3. sed -n '33,42p' core/role/builtin/jaira-teamlead/SKILL.md - Teamlead schreibt 'jaira note' und 'jaira set <id> mode=conversational'; prüfen: steht dort eine Einschränkung auf offene Entscheidungen vor der Plan-Lane? (heute: nein - die Lücke). 4. sed -n '395,400p' core/role/builtin/jaira-dispatcher/SKILL.md - die Dritte-Runde-Frage verweist auf 'How you ask', geht also im Subagent-Fall über denselben Weg zum Teamlead. 5. sed -n '16,18p' core/release/NOTES.md - genau eine '- '-Zeile unter '## Unreleased'. 6. Echter Lauf (optional): 'jaira roles install --global --force', Teamlead ohne Herdr, Dispatcher als Subagent auf ein Ticket mit einer offenen Entscheidung - Teamlead fragt als Auswahl, danach 'jaira show <id> --json': Note vorhanden und mode=conversational."
---

# Der Dispatcher fragt in Prosa statt mit Auswahl

## Definition of Done

- [x] jaira-dispatcher/SKILL.md: ein Abschnitt, wie gefragt wird (AskUserQuestion, 2–4 Optionen, Empfehlung zuerst, Rückfall nummerierte Liste); beide Fragestellen verweisen darauf
  proof: core/role/builtin/jaira-dispatcher/SKILL.md:164 Abschnitt 'How you ask'; Verweise core/role/builtin/jaira-dispatcher/SKILL.md:53 und core/role/builtin/jaira-dispatcher/SKILL.md:394
- [x] jaira-teamlead/SKILL.md: Fragen eines Dispatchers werden als AskUserQuestion weitergegeben, nicht als Text
  proof: core/role/builtin/jaira-teamlead/SKILL.md:33-45
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
- **2026-10-04 14:01 · Alexander Sacharov** — Review Runde 4 schickt zurück: teamlead/SKILL.md:37-39 setzt mode=conversational nach jeder weitergereichten Antwort - auch nach der Ja/Nein-Frage der dritten Runde (dispatcher/SKILL.md:398-399 -> 'How you ask' Rückfallweg). Ohne Herdr schaltet ein vierter Testing-Durchlauf so den Modus ein (--no-worktree, Commit-Zeilen), ohne dass jemand den Zuschnitt entschieden hat. Fix: Modus nur, wenn die Antworten offene Entscheidungen vor der Plan-Lane schließen; sonst nur die Note. Nebenbei: Z.42 'nobody else sets it' übersieht Dispatcher-Schritt 5.
- **2026-10-04 15:15 · Alexander Sacharov** — Alex hat entschieden (2026-10-04): Fix + fünfte Review-Runde. Der Teamlead setzt mode nur, wenn die Antworten offene Entscheidungen vor der Plan-Lane schließen; jede andere Antwort (z. B. Stopp-Frage nach drei Runden) bekommt nur die Note.
