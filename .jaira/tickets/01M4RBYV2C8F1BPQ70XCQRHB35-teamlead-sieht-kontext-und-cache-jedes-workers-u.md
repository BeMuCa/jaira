---
id: 01M4RBYV2C8F1BPQ70XCQRHB35
title: Teamlead sieht Kontext und Cache jedes Workers und schliesst kalte Worker
status: testing
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Der Teamlead weiss zu jedem laufenden Worker, wie gross sein Kontext ist und wie lange sein Prompt-Cache noch warm ist, und schliesst den Herdr-Tab eines Workers, dessen Ticket auf einen Menschen wartet und dessen Cache kalt ist oder dessen Kontext ueber ~150k liegt."
context: |-
  Was heute falsch ist: der Teamlead sieht nicht, wie voll oder wie kalt ein Worker ist. Alex musste am 2026-10-12 selbst bemerken, dass der NZZDRH-Worker bei ~300k Kontext lag.
  Folge: Worker-Tabs bleiben nach signoff offen und warten Stunden auf Alex. Ihr Cache (Abo 1 h, API-Key 5 min) ist dann kalt, und Give back an diesen Worker zahlt den ganzen Kontext neu.
  Was schon da ist: das Cockpit (core/mod/builtin/jaira-cockpit, hooks/register.tsx writeRecord) schreibt pro Sitzung einen Record mit tokens, ttl (Cache-Restzeit) und lastAnswerAt; Give back liest ihn (chooseWorker: alter Worker nur bei warmem Cache und <150k).
  Was fehlt: 1) der Teamlead bekommt diese Zahlen nicht (weder in den Drei-Zeilen-Berichten der Worker noch per Befehl); 2) es gibt keine Regel, wann der Teamlead einen wartenden Worker schliesst.
  Auslöser: Alex am 2026-10-12: 'teamlead soll Info bekommen, wo der Cache ist; wenn er lange nicht geantwortet hat, Worker schliessen'.
  Erst nach NZZDRH anfangen: beide aendern register.tsx.
definition-of-done: "Ein Befehl (z. B. '/cockpit workers' oder 'jaira'-naher Weg) listet jeden laufenden Worker dieses Boards mit Ticket, Lane, Kontext-Tokens und Cache-Restzeit; Test"
tags: []
blocked-by:
  - NZZDRH
related:
  - 01M4P5FK818EAFZ68XFYNZZDRH
commits: []
created-at: 2026-10-12T15:33:33Z
updated-at: 2026-10-12T18:26:55Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-68234
claimed-at: 2026-10-12T18:17:28Z
outcome-what: "jaira-cockpit lists the board's workers (/cockpit workers, tool workers) with lane, context, cache and the next-round verdict; chooseWorker limits context by a share of the model window (workerContextShare, default 50%); Give back moves without --force; worker/dispatcher reports carry own context and cache; teamlead skill keeps a worker's tab open until its ticket comes back and then follows the workers verdict"
outcome-why: "the teamlead could not see how full or how cold a worker was, and closed or reused workers by a rule of its own while a fixed 150k limit hit large tasks"
outcome-resolves: "DoD 1-5: workers() register.tsx + workerLines sessions.ts with tests; skill report lines; teamlead section; release note QRHB35; Give back without --force"
review-summary: "none"
review-gaps: "folded carrierTab(json) into one carrierTab(tabs) - the wrapper only forwarded; unexported contextLimit and UNKNOWN_WINDOW (no caller outside sessions.ts). Left: workerLines computes cache minutes itself rather than cache.ts cacheView, because cache.ts imports sessions.ts (TYPED_BY_COCKPIT) and the reverse import would be circular; the size() wording beside texts.ts t.context stays, as chooseWorker's why-strings already live in sessions.ts"
---

# Teamlead sieht Kontext und Cache jedes Workers und schliesst kalte Worker

## Definition of Done

- [x] Ein Befehl (z. B. '/cockpit workers' oder 'jaira'-naher Weg) listet jeden laufenden Worker dieses Boards mit Ticket, Lane, Kontext-Tokens und Cache-Restzeit; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:487 workers() behind '/cockpit workers' and tool mcp__jaira-cockpit__workers, rows from hooks/sessions.ts:119 workerLines(); tests 'the workers of the board' in tests/board.test.ts (claude plugin test: 81 pass, exit 0)
- [x] Der Drei-Zeilen-Bericht eines Workers (jaira-worker / jaira-dispatcher Skill) nennt seinen Kontext und seine Cache-Restzeit
  proof: core/role/builtin/jaira-worker/SKILL.md:315 and core/role/builtin/jaira-dispatcher/SKILL.md:480 second report line carries own context and cache from the first line of mcp__jaira-cockpit__workers (hooks/sessions.ts selfLine); test 'a worker writes the board it belongs to and its window into its record'
- [x] Teamlead-Skill: ein Worker wird nie angehalten oder geschlossen, waehrend er arbeitet oder auf einen Menschen wartet. Entschieden wird erst, wenn ein Ticket aus signoff/human in eine neue Runde zurueckgeht (Teamlead per SendMessage wie Cockpit-Give-back): Cache warm und Kontext unter der Grenze -> derselbe Worker macht weiter; Cache kalt oder Kontext ueber der Grenze -> Tab schliessen, neuer Worker ueber spawn.sh. Die Grenze ist einstellbar, als Anteil am Kontextfenster des Modells (Default 50%), nicht fest 150k. Dieselbe Regel wie chooseWorker im Cockpit, nicht eine zweite
  proof: core/role/builtin/jaira-teamlead/SKILL.md:158 'A tab stays open until its ticket comes back': never closed while working or waiting, decided on a new round by the verdict of the workers tool = chooseWorker (hooks/sessions.ts:73) with share = userConfig workerContextShare (default 50, .claude-plugin/plugin.json); tests 'a context under half of its window goes back to the same worker…', 'past half of its window a new worker takes the round', 'the share is a setting'
- [x] Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/QRHB35.md
- [x] Give back (register.tsx giveBack) schickt ohne --force zurueck, sobald jaira 0.3.7 (Q484YS) installiert ist; --force hebt sonst still alle anderen Weigerungen auf (Review-Gap aus NZZDRH)
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:446 move to in-progress without --force; test 'notes why and moves the ticket back to in-progress, without --force, so any other refusal still stops it' (tests/board.test.ts)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-12 15:34 · Alexander Sacharov** — Alex 2026-10-12: a fixed 150k limit would hit huge tasks. Decision: the trigger is a cold cache while the ticket waits for a person; size is only a configurable quality switch as a share of the model's window (default 50%), never applied to a worker that is still working. DoD 3 reworded.
- **2026-10-12 15:35 · Alexander Sacharov** — Alex 2026-10-12: the worker must never be stopped. Close/replace only when a ticket comes back from signoff into a new round - that is where to try a fresh worker. DoD 3 reworded again (no closing while waiting).
- **2026-10-12 18:17 · Alexander Sacharov** — parent session: jaira-4c
- **2026-10-12 18:18 · Alexander Sacharov** — Open decisions counted before the plan: none left for a person. Settled from the code: (1) DoD1 = one function behind both '/cockpit workers' (person) and the cockpit MCP tool 'workers' (mcp__jaira-cockpit__workers) - the teamlead is an agent and cannot type a slash command; it lists carrier records (worker/dispatch) of this project (records gain 'repo' = $.session.repo().root, shared by all worktrees) with a live Herdr tab, plus ticket lane from 'jaira list --json', context k/window, cache min left, and chooseWorker's verdict for a next round. (2) DoD2: the tool's first line is the calling session itself from live $.session.usage(); worker/dispatcher skills put it into their report. (3) DoD3: limit = userConfig workerContextShare (percent, default 50) x usage.context.window, stored in the record as 'window'; a record without window (written by an older build, <24h) counts as a 200k window - the conservative side. The teamlead skill reads the verdict from the tool instead of restating the rule, so there is one rule. (4) DoD5: drop --force without a version check - the mod is go:embedded in the jaira binary, so any build carrying this mod is >= 0.3.7.
- **2026-10-12 18:21 · Alexander Sacharov** — Learned: $.session.usage().context.window exists (tests mock it at 1_000_000) - that is where the record's window comes from. Without Herdr in the asking session workerLines cannot see tabs, so it lists every record of this repo and judges as if the tab were alive; records are only written with HERDR_PANE_ID, so that case is rare.
- **2026-10-12 18:23 · Alexander Sacharov** — Teamlead skill: the old sections 'Close what you started' (close at push) and 'Bringing one back' (always a new session) contradicted DoD3 and were rewritten into one section; the 'resume to ask it something' paragraph was kept. core/role/role_test.go TestCrossReferencesCarryThePrefix refuses a literal '/worker' in a skill, so the tab label is described in words.
- **2026-10-12 18:24 · Alexander Sacharov** — critique round 1: (1) board.ts lanesOf landed under the whoOf comment - move; (2) workerLines without a tab list judged the tab alive, Give back without Herdr gives to nobody - use tab !== undefined; (3) share as module state does not fit guarded(threshold, usdPerMTok) - pass it as an argument. Left alone: the separate selfLine wording (DoD2 wants the worker's own figures, which only the live usage has).
- **2026-10-12 18:25 · Alexander Sacharov** — critique round 2: code clean after round 1; teamlead skill step 2 must make the kept worker re-read the ticket (stale-board warning of the removed 'Bringing one back' section), step 3 reflow.
- **2026-10-12 18:26 · Alexander Sacharov** — critique round 3: nothing left - round 1 and 2 findings repaired (notes above).
