---
id: 01M4P0M3CJRGBSDH1X62S4M0TX
title: Ein Lane-Eval vergleicht zwei Arbeitsweisen am selben Ticket und schreibt die Metriken
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein wiederholbarer Workflow fuehrt dasselbe Ticket in zwei Arbeitsweisen aus, laesst die Diffs blind bewerten und schreibt Agenten, Tokens, Zeit, Schleifen und Urteil in eine CSV"
context: |-
  Der A/B-Versuch vom 11.10.2026 (Report: https://claude.ai/artifact/7hu5TjpnChgLF6PkcuugfS) lief von Hand: ich habe sechs Klone angelegt, jede Lane einzeln gestartet und Tokens aus den Agent-Meldungen in eine CSV getippt. Das ist nicht wiederholbar und nicht vergleichbar mit dem naechsten Lauf.
  Was fehlt, damit es ein Standard wird:
  - Ein Skript statt Handarbeit: .claude/workflows/lane-eval.js (Claude Code Workflow-Tool, Skill 'workflow-authoring'). NICHT im jaira-Binary - es ist Entwicklerwerkzeug, kein Board-Feature (Scope: kleiner als paca).
  - Isolation wie im Versuch: je Arm und Ticket ein 'git clone' ohne Remote in einem Temp-Verzeichnis, Basis-Commit fest, Tickets im Klon auf den Ausfuehrenden umgehaengt. Nie pushen, nie das echte Board anfassen.
  - Blinder Richter pro Ticket, ZWEIMAL mit vertauschtem X/Y (LLM-Richter bevorzugen die erste Position).
  - Wiederholungen: N Laeufe pro Arm und Ticket (Default 3), weil ein Lauf Zufall sein kann.
  - Goldenes Ticketset: abgeschlossene Tickets aus .jaira/logbook/ - Basis = Commit vor dem ersten Ticket-Commit, der akzeptierte Diff dient dem Richter als Referenz.
  Bekannte Stolperfalle: auf dieser Maschine fallen ~16 Tests schon auf master durch, weil Tests ~/.jaira/default-board.md lesen. Im Versuch galt: DoD 'Suite gruen' ist erfuellt, wenn dieselben Tests fehlschlagen wie auf der Basis. Mit HOME auf ein ABSOLUTES Temp-Verzeichnis ist die Suite gruen (relativ bricht go mit GOMODCACHE ab).
definition-of-done: "Workflow .claude/workflows/lane-eval.js nimmt Ticket-Ids, Basis-Commit, Arme (lane-per-agent, worker) und N; legt pro Arm/Ticket/Lauf einen Klon ohne Remote an und raeumt ihn nicht weg"
tags: []
blocked-by: []
related: []
commits: []
created-at: 2026-10-11T17:36:58Z
updated-at: 2026-10-11T17:47:28Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-522427
claimed-at: 2026-10-11T17:39:04Z
outcome-what: "clone.sh: jaira-Aufrufe ueber j(), das nur die 'older version'-Zeile filtert; docs/LANE-EVAL.md: repo muss scripts/lane-eval/ tragen."
outcome-why: "critique 1: stille Abbrueche in clone.sh, fehlende Voraussetzung in der Anleitung."
outcome-resolves: "Ein scheiternder jaira-Schritt im Setup nennt jetzt seinen Fehler; die Anleitung sagt, welches repo gemeint ist."
review-summary: |-
  scripts/lane-eval/clone.sh: jeder jaira-Aufruf laeuft mit 2>/dev/null unter set -e - scheitert restore/delete/set, bricht das Skript ohne ein Wort ab; einen Helfer j() einfuehren, der nur die 'older version'-Warnung filtert und echte Fehler durchlaesst
  docs/LANE-EVAL.md 'Start it': sagt nicht, dass args.repo scripts/lane-eval/ enthalten muss (der Workflow ruft ${repo}/scripts/lane-eval/*.sh) - vor dem Merge also den Worktree des Branches als repo angeben
---

# Ein Lane-Eval vergleicht zwei Arbeitsweisen am selben Ticket und schreibt die Metriken

## Definition of Done

- [x] Workflow .claude/workflows/lane-eval.js nimmt Ticket-Ids, Basis-Commit, Arme (lane-per-agent, worker) und N; legt pro Arm/Ticket/Lauf einen Klon ohne Remote an und raeumt ihn nicht weg
  proof: .claude/workflows/lane-eval.js (args tickets/base/arms/n, runUnit -> setup); scripts/lane-eval/clone.sh: git clone + remote remove, never removed; tried on BG5QJ6 and GDB5YJ
- [x] Pro Ticket bewertet ein blinder Richter die Diffs zweimal mit vertauschtem X/Y; die Zuordnung steht nur in einer KEY-Datei
  proof: .claude/workflows/lane-eval.js judgePair (p1/p2 swapped, split on disagreement); scripts/lane-eval/judge-prep.sh writes judge/<t>-r<n>-p1|p2 and KEY.txt; tried on BG5QJ6
- [x] Ergebnis ist eine CSV (arm, ticket, run, lane, tokens, duration_ms, tool_uses, loops, verdict) plus eine Zusammenfassung pro Arm
  proof: scripts/lane-eval/metrics.sh -> metrics.csv (arm,ticket,run,lane,tokens,duration_ms,tool_uses,loops,verdict) + summary.csv; checked against hand metrics.csv row A,8WFNQY,6 (duration, tool_uses exact)
- [x] Eine kurze Anleitung in docs/ sagt, wie man es startet und wie man ein Ticket aus dem logbook als goldenes Ticket nimmt; das Skript ist syntaktisch geprueft (node --check) und die Klon-/CSV-Schritte sind als Bash-Teil einzeln an einem Ticket ausprobiert
  proof: docs/LANE-EVAL.md; node --check on the wrapped script rc=0; mock run of the script with fake agent() (8 clones, verdict mapping right); clone/judge-prep/metrics run on BG5QJ6

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Helfer scripts/lane-eval/clone.sh: Klon ohne Remote pro Arm/Ticket/Lauf, Basis-Commit, goldenes Ticket aus logbook frisch anlegen, Referenz-Diff, .lane-eval.json
- [x] Helfer scripts/lane-eval/judge-prep.sh: Diffs zweier Laeufe als X/Y in zwei Richter-Ordnern (vertauscht), KEY.txt nur ausserhalb
- [x] Helfer scripts/lane-eval/metrics.sh: Agent-Transkripte (Marker im Prompt) + results.json -> metrics.csv und summary.txt pro Arm
- [x] .claude/workflows/lane-eval.js: Setup, Arme (lane-per-agent Schleife, worker), frischer review, Richter x2, results.json; node --check
- [x] docs/LANE-EVAL.md: Start, goldenes Ticket aus logbook, Auswertung
- [x] Bash-Teile an einem logbook-Ticket einzeln ausprobieren (Klon, judge-prep, metrics mit Fake-Transkript)

## Progress
- **2026-10-11 17:37 · Alexander Sacharov** — Den echten Workflow-Lauf kann nur die Hauptsitzung starten (Workflow-Tool braucht Opt-in des Menschen). Der Worker schreibt und prueft das Skript; der Trockenlauf mit einem Ticket und N=1 passiert danach in der Hauptsitzung.
- **2026-10-11 17:40 · Alexander Sacharov** — Entscheidungen vor dem Plan (keine braucht den Menschen): (1) agent() im Workflow liefert KEINE Tokens/Dauer/Tool-Zahl, nur Text bzw. Schema-Objekt; das Skript hat auch kein Dateisystem. Darum beginnt jeder Agent-Prompt mit einem Marker [lane-eval arm=.. ticket=.. run=.. step=.. lane=..], und scripts/lane-eval/metrics.sh liest danach die agent-*.jsonl im Transkript-Ordner (usage der letzten Antwort = tokens, erster/letzter timestamp = duration_ms, eindeutige tool_use-ids = tool_uses). Gegen den Handversuch geprueft: duration und tool_uses stimmen exakt, tokens ~1% darunter. (2) Goldenes Ticket aus dem logbook wird im Klon per 'jaira create' frisch angelegt (Titel/Ziel/Kontext/DoD aus der logbook-Datei), damit Plan, Notizen, Outcome und Review nicht durchsickern; es bekommt dadurch eine neue Id, die in .lane-eval.json steht. (3) loops: lane-per-agent zaehlt Rueckspruenge nach in-progress aus den Agent-Ergebnissen, worker meldet sie selbst. (4) Richter: Lauf i von Arm A gegen Lauf i von Arm B, zweimal mit vertauschtem X/Y; verdict pro Zeile = win/lose/tie/split dieses Laufs.
- **2026-10-11 17:46 · Alexander Sacharov** — Gefunden: (a) jaira create nimmt den Titel nur positional, --title gibt es nicht als Flag-Ersatz bei fehlendem Argument (Exit 2). (b) Ein Ticket aus der Basis kann in backlog stehen (GDB5YJ an 43080ed) - clone.sh schiebt es bis vor die erste Agent-Lane (todo), sonst stoppt lane-per-agent sofort. (c) node --check lehnt das top-level return eines Workflows ab; Pruefung ueber eine Wrapper-Funktion (docs/LANE-EVAL.md). (d) Ungeprueft: wo die Workflow-Laufzeit ihre agent-*.jsonl ablegt - metrics.sh sucht darum rekursiv unter jedem angegebenen Ordner und filtert per Marker+eval-id. Das zeigt erst der Trockenlauf in der Hauptsitzung.
- **2026-10-11 17:46 · Alexander Sacharov** — critique 1: (1) clone.sh verschluckt jaira-Fehler (2>/dev/null + set -e = stiller Abbruch) -> Helfer j(), der nur die 'older version'-Zeile filtert. (2) docs/LANE-EVAL.md muss sagen, dass repo die Helfer unter scripts/lane-eval/ tragen muss. Stehen gelassen: Record-Agent schreibt results.json per Heredoc - der Workflow hat kein Dateisystem, ein anderer Weg braucht mehr Teile.
