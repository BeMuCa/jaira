---
id: 01M43C1R4XCMMNDAD75YMCTRTD
title: Der Dispatcher fragt in Prosa statt mit Auswahl
status: done
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
commits:
  - 2be7f7557c302785f8534ca19258ceed9c2f75f3
  - e8cde06c368d6dceac55e33af066881aa23e0f02
  - e69bed3b60e494ee4f432ab37acc5fef8f42d262
  - 11987efb28e40a644e46efdb0ece01f1a82d38ca
  - ee0f4ee76904b4150d325549f07252ee6e55a1da
created-at: 2026-10-04T11:51:06Z
updated-at: 2026-10-04T15:22:58Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-23999
claimed-at: 2026-10-04T11:51:17Z
outcome-what: "Review bestanden"
outcome-why: "Runde 5 ohne blockierenden Befund"
outcome-resolves: review
review-summary: "Reiner Prompt-Text. dispatcher/SKILL.md: Abschnitt 'How you ask' (Z.164-185) - jede Frage per AskUserQuestion, 2-4 Optionen, Empfehlung zuerst mit '(Recommended)', Ja/Nein als Zwei-Optionen-Wahl, keine 'other'-Option; ohne AskUserQuestion (Subagent/headless) gehen alle offenen Entscheidungen als nummerierte Listen in Fragereihenfolge zusätzlich zu den drei Berichtszeilen hoch, und der Dispatcher stoppt. Schritt 3 (Z.53) und die Dritte-Runde-Frage (Z.398-399) verweisen darauf. teamlead/SKILL.md Z.33-45: Teamlead stellt die Listen unverändert als Auswahl und schreibt jede Antwort per 'jaira note', bevor ein frischer Dispatcher startet; nur wenn die Antworten die offenen Entscheidungen vor der Plan-Lane schließen, setzt er zusätzlich 'mode=conversational' (Ersatz für Dispatcher-Schritt 5); jede andere Antwort, z. B. Ja/Nein nach drei Runden, bekommt nur die Note. Subagent-Prompt Z.72-77 verlangt die Optionslisten neben den drei Zeilen. Runde 5 (43e1de7): Modus-Regel eingeschränkt, ungenaues 'nobody else sets it' entfernt. NOTES.md: eine '- '-Zeile unter '## Unreleased'. go test ./core/role/ ./core/release/ grün."
review-gaps: "none (blockierend). Runde-4-Lücke geschlossen: teamlead/SKILL.md:40-45 setzt den Modus nur noch für Antworten auf offene Entscheidungen vor der Plan-Lane, die Dritte-Runde-Frage bekommt nur die Note; 'nobody else sets it' ist weg. Nicht blockierend: a) Die Optionslisten im Rückfallweg (dispatcher/SKILL.md:175-182) tragen keine Kennung, ob sie Vor-Plan-Entscheidungen oder die Stopp-Frage sind - der Teamlead muss das aus Lane und Berichtszeilen ablesen (Stopp-Frage kommt mit den Befunden der dritten Runde, Vor-Plan-Fragen vor pre-process); im realistischen Lauf unterscheidbar. b) vorbestehend (galt vor 3e95568): dispatcher/SKILL.md:46-49 beschreibt 'gestorben zwischen Schritt 4 und 5, Modus noch nicht gesetzt', aber kein frischer Dispatcher setzt ihn dann nach."
test-verdict: |-
  go test ./core/role/ ./core/release/ grün. Verhalten (fragt der Dispatcher wirklich mit Auswahl?) nur in einem echten Lauf prüfbar.
  Runde 2: go test ./core/role/ ./core/release/ grün.
  Runde 3: go test ./core/role/ grün.
  Runde 4: go test ./core/role/ ./core/release/ grün.
  Runde 5: go test ./core/role/ ./core/release/ grün.
question: "Nach 'jaira roles install --global --force' einen Dispatcher auf ein Ticket mit offener Entscheidung ansetzen: kommt die Frage als Auswahl?"
review-verdict: "Bestanden. DoD erfüllt und belegt, Runde-4-Befund geschlossen, keine neue Fehlwirkung im Haupt- oder Rückfallweg gefunden, Querverweise stimmen, NOTES-Zeile regelkonform (eine Zeile, Anweisung, unter Unreleased), Tests grün. Unsicher bleibt nur, ob sich der Teamlead im echten Lauf an die Unterscheidung Vor-Plan-Frage vs. Stopp-Frage hält - das zeigt erst ein Lauf ohne Herdr."
review-check: "1. cd /home/alex/projects/jaira && go test ./core/role/ ./core/release/ - beide Zeilen 'ok'. 2. sed -n '164,185p' core/role/builtin/jaira-dispatcher/SKILL.md - Abschnitt 'How you ask': AskUserQuestion, 2-4 Optionen, '(Recommended)', Rückfall nummerierte Listen. 3. sed -n '52,53p;398,399p' derselben Datei - beide Fragestellen verweisen auf 'How you ask'. 4. sed -n '33,45p' core/role/builtin/jaira-teamlead/SKILL.md - 'jaira note' für jede Antwort; 'mode=conversational' nur 'When the answers closed the open decisions a dispatcher found before the plan lane'; Satz 'Any other answer, such as the yes or no after three rounds ... leaves the mode alone' vorhanden; kein 'nobody else sets it' mehr. 5. sed -n '16,18p' core/release/NOTES.md - genau eine '- '-Zeile unter '## Unreleased'. 6. Echter Lauf (optional): 'jaira roles install --global --force', Teamlead ohne Herdr, Dispatcher als Subagent auf ein Ticket mit einer offenen Entscheidung - Teamlead fragt als Auswahl mit Empfehlung zuerst; danach 'jaira show <id> --json': Note vorhanden und mode=conversational."
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
- **2026-10-04 15:16 · Alexander Sacharov** — Review Runde 5 bestanden. Offen, nicht blockierend: (a) die Optionslisten im Rückfallweg (dispatcher SKILL.md:175-182) sagen nicht, ob sie Entscheidungen vor der Plan-Lane oder die Stopp-Frage sind — der Teamlead liest das aus Lane und Berichtszeilen; (b) älter als dieses Ticket: stirbt ein Dispatcher zwischen Schritt 4 und 5, setzt der frische den Modus nicht nach (dispatcher SKILL.md:46-49).
