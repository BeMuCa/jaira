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
updated-at: 2026-10-11T18:32:20Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-522427
claimed-at: 2026-10-11T17:39:04Z
outcome-what: "metrics.sh summiert den Verbrauch pro Agent (entdoppelt nach message id) in Spalten input/cache_write/cache_read/output, cost_eq und peak_context; judge-prep.sh legt die Richter-Ordner als mktemp-Ordner ausserhalb OUT an; Arm- und Richter-Prompts nennen nur Klon bzw. Richter-Ordner und beginnen mit cd; docs/LANE-EVAL.md nachgezogen"
outcome-why: "Signoff-Rueckgabe: tokens war die Spitze des letzten Zuges und beguenstigte den worker-Arm; Blindheit und Isolation hingen nur am Prompt"
outcome-resolves: "DoD 5 und 6; A/B-Kontrolle ergibt cost_eq A 2,32 Mio, B 1,50 Mio"
review-summary: "Adds a lane-eval harness: .claude/workflows/lane-eval.js (Workflow script) runs each ticket N times in two arms (a fresh agent per lane vs one worker up to review), then a fresh review agent per arm, then two blind judges per ticket/run with X/Y swapped (disagreement = split), and has an agent write results.json. scripts/lane-eval/clone.sh makes one git clone per arm/ticket/run with origin removed at a fixed base, re-creating a logbook (golden) ticket fresh in todo and writing the accepted diff as reference; judge-prep.sh writes p1/p2 judge folders (X.diff, Y.diff, ticket.md, reference.diff) and the arm mapping only to OUT/KEY.txt; metrics.sh reads agent-*.jsonl transcripts whose first prompt line carries the [lane-eval id=...] marker and writes metrics.csv + summary.csv. docs/LANE-EVAL.md explains start, golden tickets and evaluation."
review-gaps: "1) tokens in metrics.sh is the LAST turn's context (input+cache+output of the final message), not tokens consumed: a synthetic transcript with turns of 122 and 228 reports 228, not 350. Summed per arm this compares many small peaks (lane-per-agent) with one large peak (worker) and understates the long single session's real cost - the cost column the eval exists for is biased toward the worker arm. 2) Blindness is prompt-only: KEY.txt sits at OUT/KEY.txt and the arm-named clone dirs at OUT/<arm>/..., both readable by the judge agent, which is only told not to look; moving KEY.txt and clones outside OUT/judge's parent would make it structural. 3) Clone isolation is likewise instruction-only: agents start in the session cwd (the real repo) and use the real ~/.jaira (roles, notify hook); a forgotten cd writes to the real board. The clones themselves are clean (0 remotes, verified). 4) Unverified until the dry run: that Workflow agents write agent-*.jsonl with the prompt as first user message (true for Agent-tool subagents here) - metrics.sh fails loudly if not."
test-verdict: "pass: go build ./... rc=0, go test -count=1 ./... mit leerem absolutem HOME rc=0; DoD 1-4 im Baum nachgeprueft: clone.sh legt fuer BG5QJ6 (logbook) zwei Klone ohne Remote an (neue Handles, todo), judge-prep.sh schreibt p1/p2 vertauscht + KEY.txt, X/Y-Diffs nennen keinen Arm, metrics.sh schreibt metrics.csv/summary.csv aus markierten Transkripten; node --check auf dem gewrappten Skript rc=0; Mock-Lauf des Workflows (8 Klone, 57 Agenten) ordnet Sieger richtig zu. Nicht pruefbar hier: ein echter Workflow-Lauf - der gehoert der Hauptsitzung."
review-verdict: "The diff meets all four DoD items and the mechanics work (clone, judge swap, CSV verified by hand); the script matches the Workflow API (pure-literal meta, pipeline stage signature, null handling, no Date/random). Approve with reservations: the token metric measures peak context rather than spend, which can skew the very comparison this is for - worth deciding before the first real run."
review-check: "1. cd into the s4m0tx worktree; S=$(mktemp -d)  2. bash scripts/lane-eval/clone.sh $PWD $S/ev worker BG5QJ6 1 - prints JSON with handle, lane todo; git -C $S/ev/worker/BG5QJ6/r1 remote prints nothing  3. bash scripts/lane-eval/clone.sh $PWD $S/ev lane-per-agent BG5QJ6 1, then echo x >> $S/ev/worker/BG5QJ6/r1/README.md  4. bash scripts/lane-eval/judge-prep.sh $S/ev BG5QJ6 1 lane-per-agent worker - p1/Y.diff and p2/X.diff are the README change, the other two empty; cat $S/ev/KEY.txt shows the two mappings swapped  5. The real check is the dry run in the main session: Workflow lane-eval with args {repo, out, tickets:['BG5QJ6'], n:1}, then bash scripts/lane-eval/metrics.sh <out> <transcriptDir> - metrics.csv must have rows for both arms."
---

# Ein Lane-Eval vergleicht zwei Arbeitsweisen am selben Ticket und schreibt die Metriken

## Definition of Done

- [x] Workflow .claude/workflows/lane-eval.js nimmt Ticket-Ids, Basis-Commit, Arme (lane-per-agent, worker) und N; legt pro Arm/Ticket/Lauf einen Klon ohne Remote an und raeumt ihn nicht weg
  proof: .claude/workflows/lane-eval.js (args tickets/base/arms/n, runUnit -> setup); scripts/lane-eval/clone.sh: git clone + remote remove, never removed; tried on BG5QJ6 and GDB5YJ
- [x] Pro Ticket bewertet ein blinder Richter die Diffs zweimal mit vertauschtem X/Y; die Zuordnung steht nur in einer KEY-Datei
  proof: .claude/workflows/lane-eval.js judgePair (p1/p2 swapped, split on disagreement); scripts/lane-eval/judge-prep.sh writes judge/<t>-r<n>-p1|p2 and KEY.txt; tried on BG5QJ6
- [x] Ergebnis ist eine CSV (arm, ticket, run, lane, Token-Spalten, duration_ms, tool_uses, loops, verdict) plus eine Zusammenfassung pro Arm
  proof: scripts/lane-eval/metrics.sh -> metrics.csv (arm,ticket,run,lane,input/cache_write/cache_read/output_tokens,cost_eq,peak_context,duration_ms,tool_uses,loops,verdict) + summary.csv; run on marked copies of the A/B transcripts
- [x] Eine kurze Anleitung in docs/ sagt, wie man es startet und wie man ein Ticket aus dem logbook als goldenes Ticket nimmt; das Skript ist syntaktisch geprueft (node --check) und die Klon-/CSV-Schritte sind als Bash-Teil einzeln an einem Ticket ausprobiert
  proof: docs/LANE-EVAL.md; node --check on the wrapped script rc=0; mock run of the script with fake agent() (8 clones, verdict mapping right); clone/judge-prep/metrics run on BG5QJ6
- [x] metrics.sh misst Verbrauch statt Spitze: pro Agent Nachrichten nach message id entdoppelt, Summe ueber alle Zuege von input_tokens, cache_creation_input_tokens, cache_read_input_tokens als Spalten plus cost_eq = input + cache_write*1.25 + cache_read*0.1; die alte Zahl bleibt als Spalte peak_context; summary.csv und docs/LANE-EVAL.md nachgezogen; an den alten A/B-Transkripten ergibt sich A 2,32 Mio, B 1,50 Mio cost_eq
  proof: scripts/lane-eval/metrics.sh per_agent (dedup by message.id, sums, cost_eq floor, peak_context); on marked copies of the A/B transcripts (usage.py agent lists) summary.csv: A cost_eq 2322803, B 1502474, input/cache columns equal usage.py exactly; docs/LANE-EVAL.md 'Read the result'
- [x] Blindheit und Klon-Isolation stehen nicht nur im Prompt: die Eingaben des Richters liegen in einem Ordner mit neutralen Namen ohne KEY.txt, und Prompt (und cwd, soweit agent() das erlaubt) jedes Arm-Agenten zeigen auf seinen Klon; was bleibt, ist im Ticket notiert
  proof: scripts/lane-eval/judge-prep.sh: folders via mktemp outside <out>, dir only in KEY.txt (tried on BG5QJ6, rerun replaces KEY lines); .claude/workflows/lane-eval.js rules()/judge(): prompts name only clone resp. judge folder, cd first (mock run: no repo/out path in arm or judge prompts); residuals in docs/LANE-EVAL.md 'What keeps the judge blind'

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
- [x] metrics.sh: Verbrauch summieren, peak_context behalten, summary.csv, an alten A/B-Transkripten pruefen
- [x] judge-prep.sh/lane-eval.js: Richter-Ordner ausserhalb OUT, Arm-Prompt beginnt mit cd in den Klon
- [x] docs/LANE-EVAL.md nachziehen

## Progress
- **2026-10-11 17:37 · Alexander Sacharov** — Den echten Workflow-Lauf kann nur die Hauptsitzung starten (Workflow-Tool braucht Opt-in des Menschen). Der Worker schreibt und prueft das Skript; der Trockenlauf mit einem Ticket und N=1 passiert danach in der Hauptsitzung.
- **2026-10-11 17:40 · Alexander Sacharov** — Entscheidungen vor dem Plan (keine braucht den Menschen): (1) agent() im Workflow liefert KEINE Tokens/Dauer/Tool-Zahl, nur Text bzw. Schema-Objekt; das Skript hat auch kein Dateisystem. Darum beginnt jeder Agent-Prompt mit einem Marker [lane-eval arm=.. ticket=.. run=.. step=.. lane=..], und scripts/lane-eval/metrics.sh liest danach die agent-*.jsonl im Transkript-Ordner (usage der letzten Antwort = tokens, erster/letzter timestamp = duration_ms, eindeutige tool_use-ids = tool_uses). Gegen den Handversuch geprueft: duration und tool_uses stimmen exakt, tokens ~1% darunter. (2) Goldenes Ticket aus dem logbook wird im Klon per 'jaira create' frisch angelegt (Titel/Ziel/Kontext/DoD aus der logbook-Datei), damit Plan, Notizen, Outcome und Review nicht durchsickern; es bekommt dadurch eine neue Id, die in .lane-eval.json steht. (3) loops: lane-per-agent zaehlt Rueckspruenge nach in-progress aus den Agent-Ergebnissen, worker meldet sie selbst. (4) Richter: Lauf i von Arm A gegen Lauf i von Arm B, zweimal mit vertauschtem X/Y; verdict pro Zeile = win/lose/tie/split dieses Laufs.
- **2026-10-11 17:46 · Alexander Sacharov** — Gefunden: (a) jaira create nimmt den Titel nur positional, --title gibt es nicht als Flag-Ersatz bei fehlendem Argument (Exit 2). (b) Ein Ticket aus der Basis kann in backlog stehen (GDB5YJ an 43080ed) - clone.sh schiebt es bis vor die erste Agent-Lane (todo), sonst stoppt lane-per-agent sofort. (c) node --check lehnt das top-level return eines Workflows ab; Pruefung ueber eine Wrapper-Funktion (docs/LANE-EVAL.md). (d) Ungeprueft: wo die Workflow-Laufzeit ihre agent-*.jsonl ablegt - metrics.sh sucht darum rekursiv unter jedem angegebenen Ordner und filtert per Marker+eval-id. Das zeigt erst der Trockenlauf in der Hauptsitzung.
- **2026-10-11 17:46 · Alexander Sacharov** — critique 1: (1) clone.sh verschluckt jaira-Fehler (2>/dev/null + set -e = stiller Abbruch) -> Helfer j(), der nur die 'older version'-Zeile filtert. (2) docs/LANE-EVAL.md muss sagen, dass repo die Helfer unter scripts/lane-eval/ tragen muss. Stehen gelassen: Record-Agent schreibt results.json per Heredoc - der Workflow hat kein Dateisystem, ein anderer Weg braucht mehr Teile.
- **2026-10-11 17:47 · Alexander Sacharov** — critique 2: beide Befunde aus Runde 1 behoben (clone.sh j(), Fehlerfall NOSUCH meldet sich mit rc=2; docs nennt repo-Voraussetzung). Nichts Neues.
- **2026-10-11 17:49 · Alexander Sacharov** — testing: alles gruen. Offen fuer den echten Trockenlauf (Hauptsitzung, ein Ticket, n=1): ob die Workflow-Agenten ihre agent-*.jsonl im selben Format ablegen wie Agent-Tool-Subagenten und in welchem Ordner - metrics.sh sucht rekursiv; falls 'no agent of eval ... found', den Ordner aus dem Workflow-Ergebnis (transcriptDir) bzw. ~/.claude/projects/<projekt>/<session>/ uebergeben.
- **2026-10-11 17:51 · Alexander Sacharov** — review: approved with reservations. Most serious: metrics.sh tokens = last-turn context (peak), not consumption - checked with a synthetic transcript (228 reported, 350 summed). Blindness and clone isolation hold only by prompt (KEY.txt and arm dirs readable under OUT; agents start in the real repo cwd). Verified: clone has no remote, judge-prep swap and KEY.txt, metrics.csv/summary.csv, script parses.
- **2026-10-11 18:03 · Alexander Sacharov** — Review-Befund, Rueckgabe empfohlen: metrics.sh liest tokens als Kontextgroesse des letzten Zuges (Spitze), nicht als Verbrauch. Das beguenstigt den worker-Arm. Fix: pro Agent die Summe ueber alle Zuege von cache_creation_input_tokens, cache_read_input_tokens, input_tokens schreiben, dazu eine Spalte cost_eq = cache_write*1.25 + cache_read*0.1 + input; die Spitze als eigene Spalte peak_context behalten. Nachgerechnet am A/B-Versuch: A 2,32 Mio, B 1,50 Mio cost_eq.
- **2026-10-11 18:32 · Alexander Sacharov** — Runde 2: agent() im Workflow hat KEINE cwd-Option (nur label/phase/schema/model/effort/isolation/agentType; isolation:'worktree' waere ein Worktree des echten Repos, nicht der Klon). Darum nur strukturell, was billig geht: Arm-Prompt nennt keinen Pfad ausser dem Klon und beginnt mit cd; Richter-Ordner sind mktemp-Ordner ausserhalb OUT, Pfad nur in KEY.txt. Bleibt offen (in docs notiert): vergessenes cd schreibt ins echte Board; echtes ~/.jaira; alte Richter-Ordner bleiben in TMPDIR. Pruefung der Summe: Kopien der A/B-Transkripte mit eingesetztem Marker (Skript scratchpad/mkab.py) -> metrics.sh ergibt exakt die usage.py-Spalten, cost_eq A 2322803 / B 1502474 (pro Agent abgerundet, daher 7 unter der Gesamtrechnung). Summierte Spitzen waren 1,10 Mio / 0,45 Mio - die alte Spalte uebertrieb den Abstand.
