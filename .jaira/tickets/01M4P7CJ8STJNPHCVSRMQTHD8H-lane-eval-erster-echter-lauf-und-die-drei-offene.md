---
id: 01M4P7CJ8STJNPHCVSRMQTHD8H
title: "lane-eval: erster echter Lauf und die drei offenen Luecken aus dem Review"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "lane-eval ist einmal echt gelaufen (ein Ticket, N=1) und die bekannten Luecken in Blindheit und Kostenzaehlung sind geschlossen oder ehrlich dokumentiert"
context: |-
  lane-eval kam mit S4M0TX (PR BeMuCa/jaira#46, gemergt). Es wurde NIE mit dem echten Workflow-Tool gestartet - nur Mock-Lauf und Einzelschritte. Unbekannt ist daher: wo Workflow-Agenten ihre Transkripte ablegen und ob die Markierungszeile im ersten Prompt so ankommt, wie scripts/lane-eval/metrics.sh sie sucht.
  Offene Punkte aus dem Review (Notes auf S4M0TX im logbook .jaira/logbook/as-20261011/):
  1. docs/LANE-EVAL.md legt im Beispiel <out> nach /tmp, also neben die Richter-Ordner (mktemp in TMPDIR) - KEY.txt ist dann einen 'ls /tmp' vom Richter entfernt. Beispiel unter $HOME legen oder TMPDIR fuer judge-prep setzen.
  2. Subagents, die ein Arm-Agent selbst startet (Agent-Tool), tragen keine Markierung und werden nicht gezaehlt. Das kann den worker-Arm zu billig aussehen lassen, weil jaira-worker fuer review einen frischen Agenten startet.
  3. cost_eq laesst Output-Tokens weg; die stehen in den Transkripten nur unvollstaendig.
  Zusaetzlich offen (Frage von Alex am 11.10.): soll der Richter auch die Ticket-Datei sehen? Im A/B-Versuch sassen die typischen Fehler des Lane-Pipeline-Arms dort (veralteter --proof, zwei SHAs in einem commits-Eintrag), und der Richter sah sie nicht.
  Den Workflow startet nur die Hauptsitzung des Menschen (Workflow-Tool braucht sein Opt-in).
definition-of-done: Ein echter Lauf mit einem goldenen Ticket aus dem logbook und N=1 ist gelaufen; metrics.csv hat fuer jeden Agenten Zahlen; Ergebnis als Note
tags: []
blocked-by: []
related: []
commits: []
created-at: 2026-10-11T19:35:11Z
updated-at: 2026-10-11T20:19:13Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-789618
claimed-at: 2026-10-11T19:45:23Z
mode: ""
outcome-what: "lane-eval: clone.sh/judge-prep.sh verweigern ein <out> unter dem Temp-Root; metrics.sh zaehlt verschachtelte Agenten ueber parentAgentId mit (Spalte subagents) und findet die Markierung dort, wo der Workflow-Harness sie hinlegt; der Richter bekommt X/Y.ticket.md (DoD mit proof, outcome, commits); ein echter Lauf BG5QJ6 N=1 lief, Ergebnis als Note"
outcome-why: "lane-eval war nie echt gelaufen; KEY.txt lag einen ls vom Richter entfernt, verschachtelte Agenten fehlten in den Kosten, und der Richter sah die typischen Fehler auf dem Ticket nicht"
outcome-resolves: "Alle drei DoD-Punkte: Lauf mit Zahlen fuer jeden Agenten, Punkte 1 und 2 behoben, Richter sieht die Ticket-Datei (Alex' Entscheidung)"
review-summary: "scripts/lane-eval/judge-prep.sh:26-30 wiederholt die Temp-Root-Pruefung aus clone.sh:28-32; clone.sh laeuft immer zuerst auf demselben <out> - die Kopie in judge-prep.sh streichen, Kommentar dort auf 'clone.sh refuses' aendern"
---

# lane-eval: erster echter Lauf und die drei offenen Luecken aus dem Review

## Definition of Done

- [x] Ein echter Lauf mit einem goldenen Ticket aus dem logbook und N=1 ist gelaufen; metrics.csv hat fuer jeden Agenten Zahlen; Ergebnis als Note
  proof: Lauf wf_9ea1a9f1-c81: /home/alex/lane-eval/2026-10-11-qthd8h/metrics.csv hat 9 Zeilen mit Zahlen fuer alle Arm-Agenten (7 lane-per-agent, 2 worker); Ergebnis als Note
- [x] Punkte 1 und 2 sind behoben (Beispiel-Pfad; Subagents eines Arm-Agenten werden mitgezaehlt) oder in docs/LANE-EVAL.md mit ihrer Auswirkung genannt
  proof: scripts/lane-eval/clone.sh:31 verweigert ein <out> unter dem Temp-Root (Beispiel unter $HOME, docs/LANE-EVAL.md:35); scripts/lane-eval/metrics.sh:79 rechnet unmarkierte Agenten ueber parentAgentId dem markierten Vorfahren zu, Spalte subagents - Fixture mit 2 Ebenen + Richter-Helfer geprueft
- [x] Entschieden und umgesetzt, ob der Richter die Ticket-Datei sieht
  proof: Entscheidung Alex: ja, nur DoD/Outcome/commits. scripts/lane-eval/judge-prep.sh:61 schreibt X.ticket.md/Y.ticket.md (armticket), Richter-Prompt .claude/workflows/lane-eval.js:189, im Lauf wf_9ea1a9f1-c81 angekommen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] judge-prep.sh: verweigern, wenn <out> unter dem Temp-Wurzelordner der Richter liegt; Doku-Beispiel nach $HOME
- [x] metrics.sh: Agent ohne Markierung ueber meta.json/parentAgentId dem markierten Vorfahren zuschlagen
- [x] judge-prep.sh: X.ticket.md/Y.ticket.md aus der Arm-Ticket-Datei, nur DoD, Outcome, commits; Richter-Prompt und Doku
- [x] Echter Lauf BG5QJ6 N=1, metrics.sh, Ergebnis als Note

## Progress
- **2026-10-11 19:46 · Alexander Sacharov** — Worker-Entscheidungen vor dem Plan (selbst entschieden, widersprechbar):
- Punkt 2 ist behebbar statt nur zu dokumentieren: Claude Code legt neben jedes agent-<id>.jsonl ein agent-<id>.meta.json mit parentAgentId; verschachtelte Agenten liegen flach im selben subagents/-Ordner (gesehen in ~/.claude/projects/-home-alex-projects-jaira/a1fb63e1-.../subagents). metrics.sh ordnet einen Agenten ohne Markierung dem naechsten markierten Vorfahren zu.
- Punkt 1 wird im Code verhindert, nicht nur im Beispiel: judge-prep.sh verweigert, wenn <out> unter demselben Temp-Wurzelordner liegt wie die Richter-Ordner; Beispiel-Pfad nach $HOME.
- Punkt 3 (cost_eq ohne Output) ist nicht im DoD und bleibt dokumentiert wie bisher.
Offen fuer Alex: Richter sieht Ticket-Datei ja/nein; echter Lauf jetzt aus dieser Sitzung (Workflow-Opt-in).
- **2026-10-11 19:54 · Alexander Sacharov** — parent session: jaira-61 (Teamlead; Worker wurde ohne --parent gestartet)
Antworten von Alex (in meinem Tab, 11.10.):
1. Richter sieht die Ticket-Datei jeder Arm: ja, aber nur Definition of Done (mit --proof), Outcome und commits - keine Notes, kein Autor, damit der Arm nicht erkennbar wird.
2. Echter Lauf: BG5QJ6, N=1, beide Arme, nach den Fixes. Alex waehlte 'aus dieser Sitzung jetzt'; Teamlead jaira-61 sagt, die Hauptsitzung startet den Lauf nach dem Review. Vor dem Start mit jaira-61 abstimmen, nicht doppelt laufen lassen.
- **2026-10-11 19:55 · Alexander Sacharov** — Abgestimmt mit jaira-61: Alex' Wahl im Worker-Tab gilt als Opt-in - der Worker startet den echten Lauf selbst nach den Fixes und schickt jaira-61 drei Zeilen (cost_eq je Arm, Richterurteile, Brueche in metrics.sh).
- **2026-10-11 19:57 · Alexander Sacharov** — metrics.sh: Falle gefunden - capture(...)? liefert bei fehlendem Treffer 'empty', nicht null; frueher fiel ein unmarkierter Agent dadurch still raus. Jetzt [capture(...)?] | first. Herdr-Tabs (spawn.sh) sind eigene Sitzungen ohne parentAgentId und bleiben ungezaehlt - in LANE-EVAL.md genannt; die Workflow-Prompts starten keine.
- **2026-10-11 19:58 · Alexander Sacharov** — Echter Lauf gestartet: Workflow lane-eval, repo=/home/alex/projects/.worktrees/jaira-qthd8h (hat die Fixes), out=/home/alex/lane-eval/2026-10-11-qthd8h, tickets=[BG5QJ6], n=1. Stirbt diese Sitzung: Clones und r1.json liegen unter out, results.json fehlt dann.
- **2026-10-11 19:58 · Alexander Sacharov** — Gefunden: Workflow-Agenten legen ihre Transkripte unter ~/.claude/projects/<projekt>/<session>/subagents/workflows/<runId>/ ab (dieser Lauf: .../-home-alex-projects--worktrees-jaira-qthd8h/0c43ac98-.../subagents/workflows/wf_9ea1a9f1-c81). Run-ID wf_9ea1a9f1-c81, resumable aus derselben Sitzung.
- **2026-10-11 20:00 · Alexander Sacharov** — Echter Lauf zeigt: Workflow-Transkripte beginnen NICHT mit dem Prompt. 1. user-Nachricht = '[Workflow harness — user request]' mit der Anfrage, die den Workflow startete (hier '/jaira-worker QTHD8H' - die Arm-Agenten lesen das als Nutzerwillen, der bei Konflikt gewinnt!), 2. user-Nachricht = Prompt, jede Zeile 2 Leerzeichen eingerueckt. metrics.sh fand deshalb nichts; jetzt sucht es die Markierung am Zeilenanfang der ersten zwei user-Nachrichten. Workflow-Agenten haben meta.json ohne parentAgentId, agentType workflow-subagent, description = Label.
- **2026-10-11 20:18 · Alexander Sacharov** — ERGEBNIS echter Lauf (wf_9ea1a9f1-c81, BG5QJ6, N=1, out=/home/alex/lane-eval/2026-10-11-qthd8h, 15 Agenten, 19 min Wandzeit):
- cost_eq: lane-per-agent 1.19 M (7 Agenten), worker 1.00 M (2 Agenten) - worker 16% billiger. Beim Hand-A/B waren es 2.32 M gegen 1.50 M.
- Zeit (Summe der Agenten): 1060 s gegen 917 s; Tool-Aufrufe 123 gegen 85; beide 1 Schleife (critique schickte einmal zurueck); subagents 0 in beiden Armen.
- Richter: split. Beide Durchgaenge waehlten die Diff auf Position Y (p1: Y=lane-per-agent, p2: Y=worker) - die Positionsneigung, gegen die der Tausch gebaut ist, hat hier entschieden. Begruendungen: lane-per-agent loest den Abnahmeweg (acceptance/teamlead) vollstaendiger; worker erneuert AGENTS.md/CLAUDE.md, korrigiert den review-Prompt und testet end-to-end, fuehrt aber ein spekulatives Lane.OnPath ein.
- Beide Richter zitierten Inhalte aus X/Y.ticket.md nicht ausdruecklich; die Dateien kamen aber an.
- output_tokens wirken zu klein (710 bei 49 Tool-Aufrufen) - Punkt 3 aus dem Review bestaetigt, bleibt dokumentiert.
- Nebenfund im worker-Arm: 'jaira move --what/--why/--resolves' in review ueberschrieb die outcome-Felder des Implementierers; der Review-Agent stellte sie mit jaira set wieder her. Kein Teil dieses Tickets.
- **2026-10-11 20:18 · Alexander Sacharov** — critique: doppelte Temp-Root-Pruefung - nur clone.sh behalten (frueher Abbruch ist der eigentliche Wert), judge-prep.sh-Kopie und Doku-Satz anpassen. Sonst nichts: armticket-Schleife braucht git, deshalb nicht in jq; metrics-Aenderung ist minimal fuer das Gefundene.
