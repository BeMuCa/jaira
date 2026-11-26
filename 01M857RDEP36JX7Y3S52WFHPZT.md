---
id: 01M857RDEP36JX7Y3S52WFHPZT
title: "Der Lane-Verlauf eines Boards ist als Graph zu sehen und einzustellen: Lanes als Knoten, Uebergaenge und Loops als Kanten, je Lane Agent, Modell und Anzahl Agenten"
status: critique
ready: true
creator: BeMuCa
goal: "Eine Person sieht den Ablauf eines Boards als Graph - jede Lane ein Knoten, jeder Uebergang und jeder Loop (critique/optimize/testing zurueck nach in-progress) eine Kante mit ihrer Grenze - und stellt je Knoten ein, welcher Agent bzw. welches Modell die Lane laeuft und wie viele Agenten sie gleichzeitig bearbeiten (Fan-out), ohne Lane-Dateien von Hand zu editieren."
context: |-
  Wunsch von Berk 2026-11-23: 'lane-verlauf als graph mit nodes, loops und lanes, wo man einstellen kann wieviele agenten die tickets bearbeiten sollen (fan out) oder loops und welcher agent die lane durchlaufen soll'. Vergleich: LangGraph-artige Ansicht.
  Heute: der Ablauf steht verteilt in .jaira/lanes/order (Reihenfolge), den Lane-Dateien (agentic, model-tier, Prompt) und den Loop-Regeln (max 3 Rueckgaben). Sehen kann man ihn nur als Text: 'jaira lanes', der Lanes-Block in CLAUDE.md, die Lane-Vorschau in den Settings (2B4KC9).
  Was es dazu schon gibt, nur jeweils ein Teil: ATDG74 (ideas) - N Agenten je Lane (Schwarm), parallel wenn RAM reicht; CXH98W (signoff) - Modell fuer Orchestrator und Worker in den Settings; 403F8T (ideas) - Schwarm-Debatte; ENCFVZ (signoff) - Lane-Arten mit skill/script/receipt; mehrere Boards je Repo (Fan-out ueber Boards, Milestone jaira-board-fan-out).
  Keines zeigt den Ablauf als Bild, und keines stellt Agent, Anzahl und Loop-Grenze an einer Stelle ein.
  Offen fuer den Brainstorm: TUI oder jaira web oder beides; was davon nur Ansicht ist und was in die Lane-Dateien zurueckschreibt (Lane-Dateien bleiben die Quelle, die Datei ist die API); wie N Ergebnisse zusammengefuehrt werden (ATDG74).
definition-of-done: "Brainstorm: Ort (TUI/web), was angezeigt und was eingestellt wird, wie es in die Lane-Dateien zurueckschreibt; als Notiz am Ticket"
tags:
  - tui
blocked-by: []
related:
  - 01M67K72RDJY2N2R2KQGATDG74
  - 01M658C7SP4V29E95574CXH98W
  - 01M6M72GP39ZAGPEKY6F403F8T
  - 01M67K7GEZF694985EHTENCFVZ
  - 01M58G28438KJJD3YDGM2B4KC9
commits: []
created-at: 2026-11-23T22:17:22Z
updated-at: 2026-11-26T10:54:18Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: worker-wfhpzt
claimed-at: 2026-11-26T10:46:36Z
outcome-what: "The Compact View (v) sets a lane: t steps the focused lane's tier, + and - its agents (1-9). Each press rewrites one frontmatter line of the lane file (core/lane/tune.go SetTier, SetAgents), parsed before it is written; the step shows the planned count in front of its tier (3×cheap), a line under the graph names the file written."
outcome-why: "The flow was visible (R68J0X) and the key existed (FSPJP9), but changing who runs a lane still meant opening the lane file by hand."
outcome-resolves: "Concurrency: atomic write; a hand edit between read and write loses (noted). Error paths: a refused or failed write opens the REFUSED box and leaves the file alone; a reload error is shown, not dropped. Unusual input: tier taken from other lane files goes through SetScalar quoting; agents clamps at 1 and 9 and writes nothing there; lanes without agent step or without a file refuse. Permissions: writes only l.Source and, for a byte-equal custom lane, its catalogue twin as E does. Open: a shipped lane set this way is skipped by 'jaira update --lanes'."
review-summary: "internal/tui/pipeline.go tuneLane: the callbacks hand back a 'what' string and stepAgents returns early at the ends, so - on a lane without agent step does nothing while + on it is refused, and the line under the graph says what was asked for, not what the file says. Instead: callbacks return only the error, tuneLane reloads and writes the line from the reloaded lane (tier and agents), and says nothing when neither changed; stepAgents always calls lane.SetAgents with the clamped number."
round: 2
back-from: critique
back-by: lane
back-count: 1
---

# Der Lane-Verlauf eines Boards ist als Graph zu sehen und einzustellen: Lanes als Knoten, Uebergaenge und Loops als Kanten, je Lane Agent, Modell und Anzahl Agenten

## Definition of Done

- [x] Brainstorm: Ort (TUI/web), was angezeigt und was eingestellt wird, wie es in die Lane-Dateien zurueckschreibt; als Notiz am Ticket
  proof: Notiz 'Brainstorm (Teamlead, 2026-11-23)' am Ticket; Folgetickets R68J0X, FSPJP9, 3ADC5S
- [x] Der Graph zeigt jede Lane als Knoten und jeden Uebergang und Loop als Kante, Loops mit ihrer Grenze
  proof: internal/tui/pipeline.go renderPipeline (arrows between steps, loop legend with limit, from R68J0X) and renderStep (planned agents in front of the tier); tests TestPipelineShowsEveryStepInOrder, TestPipelineSetsTierAndAgentsInTheLaneFile ('3×cheap')
- [x] Je Lane lassen sich Agent bzw. Modell und Anzahl Agenten einstellen; die Einstellung steht danach in der Lane-Datei
  proof: core/lane/tune.go SetTier/SetAgents, internal/tui/pipeline.go pipelineKey (t, +, -) and tuneLane; tests TestTuneWritesOnlyItsOwnLine, TestTuneRefusesAndLeavesTheFileAlone, TestPipelineSetsTierAndAgentsInTheLaneFile, TestPipelineRefusedSettingSaysWhy, TestPipelineSettingReachesTheCatalogueTwin
- [x] Release-Notiz in core/release/notes/<ID>.md
  proof: core/release/notes/WFHPZT.md

## Options

- [x] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 22:39 · BeMuCa** — Berk 2026-11-23: die heutige Compact View (Taste v) wird zu dieser Board-Visualisierung umgebaut statt eine neue Ansicht zu bauen. Sie soll zeigen: Fan-out (N Agenten in einer Lane, ATDG74), wie zwei Boards in einer Join-Lane Infos austauschen oder diskutieren (YTE87T), mehrere Worker, die dasselbe Board parallel durchlaufen, mit ihren Worktrees gebuendelt - damit man alle Staende auf einem Board sieht. Milestone lane-graph; Fable als Teamlead, Haiku-Agenten setzen um.
- **2026-11-23 22:42 · BeMuCa** — Brainstorm (Teamlead, 2026-11-23). Was der Code heute zeigt: die Compact View (internal/tui/pipeline.go, Taste v) zeichnet schon jede Lane als Kasten mit Pfeil zur naechsten, in einer Schlange, mit Punkten je Ticket und Zahl der Agenten. Es fehlen: (1) die Loops - rejects-to und max-rounds stehen in der Lane-Datei (core/lane/lane.go RejectsTo, MaxRounds), die Ansicht zeigt sie nicht; (2) wer die Lane laeuft - model-tier, Person, fresh context stehen in der Lane-Datei, die Ansicht zeigt sie nicht; (3) Einstellen - die Ansicht schreibt nichts; (4) Anzahl Agenten je Lane - dafuer gibt es noch keinen Schluessel in der Lane-Datei (ATDG74); (5) Worker ueber Worktrees - core/session laedt nur die Sessions des eigenen Worktrees (ticket.Store.SessionsDir ist je Worktree), die Ansicht sieht fremde Worker nicht; (6) Join-Lane - haengt an YTE87T, dessen Form Berk erst freigeben muss.
Wege: A) nur TUI, Compact View umbauen (Berks Vorgabe 22:39): klein, nutzt was da ist, gibt jaira web auf. B) TUI und jaira web: doppelter Aufwand, zwei Darstellungen in Schritt halten - gegen 'kleiner als paca'. C) eigener Graph-Screen neben der Compact View: von Berk schon abgelehnt.
Entschieden: A. Die Lane-Datei bleibt die Quelle: die Ansicht liest lane.Set und schreibt nur Frontmatter-Schluessel zurueck (model-tier, agents), nie den Prompt. Welches Modell hinter strong/cheap steht, stellt CXH98W in den Settings ein (noch nicht gemergt) - der Graph stellt je Lane nur den Tier um, kein Modellname in der Lane-Datei, sonst gaebe es zwei Stellen fuer dasselbe.
Schnitt in baubare Teile: 1. Folgeticket 'Loops und Besetzung anzeigen' (nur Ansicht, sofort baubar auf master). 2. ATDG74-Folgeticket 'Lane-Schluessel agents: N' (sofort baubar). 3. WFHPZT selbst: Tier und agents im Graph einstellen, schreibt in die Lane-Datei - nach 1 und 2. 4. Folgeticket 'Worker ueber Worktrees gebuendelt' (Backlog, braucht erst einen Weg, Sessions anderer Worktrees zu lesen). 5. Join-Lane im Graph: erst wenn YTE87T entschieden ist, kein Ticket vorher.
- **2026-11-23 22:44 · BeMuCa** — Folgetickets angelegt: R68J0X (Loops und Besetzung anzeigen, zuerst), FSPJP9 (Lane-Schluessel agents: N, folgt auf ATDG74), 3ADC5S (Worker ueber Worktrees gebuendelt, Backlog). WFHPZT selbst baut danach das Einstellen im Graph (Tier und agents je Lane, schreibt in die Lane-Datei) und wartet auf R68J0X und FSPJP9. Der Brainstorm zu ATDG74 steht im Kontext von FSPJP9: ATDG74 gehoert berkc, eine BeMuCa-Session kann es ohne --steal nicht ziehen und deshalb keine Notiz daran schreiben.
- **2026-11-26 10:48 · BeMuCa** — Worker 2026-11-26, open decisions before the plan: none that need a person. Settled from the code, each open to disagreement: (1) Keys in the Compact View: t steps the focused lane's tier, + (or =) and - change agents 1..9; a/d/w/s, enter, v, 1-9 stay. (2) Tier has no fixed vocabulary in the code (core/lane/template.go calls it a local alias), so t cycles through cheap, strong and every tier a lane file of this board already names; a new alias is still typed by hand. (3) agents back at 1 removes the key instead of writing 'agents: 1' - absent means 1 (FSPJP9), and + then - leaves the file byte-identical. (4) The write goes through lane.Parse before it lands: what the loader would refuse (agents above 1 on a lane that writes outcome-what or plan) is refused with that message and nothing is written. A lane with no agent step refuses both keys. (5) A custom lane whose catalogue copy (~/.jaira/lanes) was byte-equal gets the same write-through the E edit in the lane settings does (internal/tui/lanes.go writeThrough), otherwise the lane settings would call it drifted. (6) Known cost, not solved here: a shipped lane changed this way is no longer a version jaira shipped, so 'jaira update --lanes' skips it (core/lane/refresh.go RefreshSkipped) - same as any hand edit of a lane file. (7) DoD 2 edges: R68J0X drew transitions as arrows and loops as the mark in the node plus the legend under the graph, and decided against a drawn back-arrow across the snake; this ticket keeps that and adds the configured agent count to the node.
- **2026-11-26 10:53 · BeMuCa** — in-progress: seen in a real terminal (tmux, scratch board without remote): t and + on Review wrote 'model-tier: cheap' and 'agents: 3' into .jaira/lanes/review.md, the step read '3×cheap'; + on Todo and HITL opened the REFUSED box. Measured: 'jaira update --lanes --dry-run' afterwards reports review.md as skipped (edited here). Had to find out: Doc.SetScalar writes a number quoted (agents: "3"), so the count goes through SetRaw. The help screen was exactly 60 rows in TestHelpExplainsEveryCardMark; one more help line pushed '@name' below the fold, so that test's terminal is 61 rows now (the help scrolls). Not done: no key for a loop's max-rounds - the brainstorm limits the view to model-tier and agents. Race left as it is: a hand edit of the same lane file between the view's read and its write is overwritten (the write itself is atomic).
- **2026-11-26 10:54 · BeMuCa** — critique: internal/tui/pipeline.go tuneLane: the callbacks hand back a 'what' string and stepAgents returns early at the ends, so - on a lane without agent step does nothing while + on it is refused, and the line under the graph says what was asked for, not what the file says. Instead: callbacks return only the error, tuneLane reloads and writes the line from the reloaded lane (tier and agents), and says nothing when neither changed; stepAgents always calls lane.SetAgents with the clamped number.
