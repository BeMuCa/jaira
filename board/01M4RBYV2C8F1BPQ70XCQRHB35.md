---
id: 01M4RBYV2C8F1BPQ70XCQRHB35
title: Teamlead sieht Kontext und Cache jedes Workers und schliesst kalte Worker
status: done
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
commits:
  - 98718866a2108af5cc179fce1c326c45b09936c2
  - 174bbd1c2eba494d40309d1c5e31e033596083f6
  - e2e93bfa782ca404df82d69ba8859141c75242e4
  - 82aa5654a3f692d2950c9bb4592d0d88e137344a
created-at: 2026-10-12T15:33:33Z
updated-at: 2026-10-12T19:09:08Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-68234
claimed-at: 2026-10-12T18:17:28Z
outcome-what: "jaira-cockpit lists the board's workers (/cockpit workers, tool workers) with lane, context, cache and the next-round verdict; chooseWorker limits context by a share of the model window (workerContextShare, default 50%); Give back moves without --force; worker/dispatcher reports carry own context and cache; teamlead skill keeps a worker's tab open until its ticket comes back and then follows the workers verdict"
outcome-why: "the teamlead could not see how full or how cold a worker was, and closed or reused workers by a rule of its own while a fixed 150k limit hit large tasks"
outcome-resolves: "DoD 1-5: workers() register.tsx + workerLines sessions.ts with tests; skill report lines; teamlead section; release note QRHB35; Give back without --force"
review-summary: "The jaira-cockpit gains one function, workers() (register.tsx), behind both '/cockpit workers' for a person and a new MCP tool mcp__jaira-cockpit__workers for the teamlead agent. Its first line is the asking session's own context and cache from live usage (selfLine); then one line per ticket of this board carried by a worker/dispatcher with an open Herdr tab: ticket, lane (from jaira list --json), context k of window k, cache minutes left or 'cache warm (working)', and the verdict of chooseWorker for a next round ('same worker'/'new worker' and why), skipped for a tab whose agent is working. To make that possible each session record now stores repo (main working tree, shared by all worktrees - verified against the store: worktree sessions file reports under /home/alex/projects/jaira) and window. chooseWorker's fixed 150k limit became a share of the record's window (userConfig workerContextShare, default 50%; a record without window counts as 200k), passed in as an argument like the guard threshold; Give back uses the same share. Give back now moves back to in-progress without --force. Skills: worker and dispatcher put their own context/cache (first line of the workers tool) into their second report line; the teamlead skill replaces 'close at push' and 'always a new session' with one section: a tab stays open while working or waiting, and on a new round the teamlead reads the workers verdict - same worker via SendMessage plus a re-read of jaira show, or close tab and spawn.sh. Release note core/release/notes/QRHB35.md with 4 one-line entries."
review-gaps: "1) DoD 5 said 'once jaira 0.3.7 is installed'; the change drops --force unconditionally on the argument that the mod is go:embedded so any build carrying it is >=0.3.7. That is not the only route: .claude-plugin/marketplace.json installs jaira-cockpit without jaira, so it can run beside an older jaira binary, where 'move --to in-progress' out of signoff is refused. The failure is loud (fail toast, ticket stays), but giveBack has already written 'Given back from signoff: <reason>' onto the ticket before the refused move, so the ticket says it was given back when it was not. Not a blocker; a version check or a README line on the minimum jaira version would close it. 2) Teamlead skill: tabs now stay open until done, so idle sessions accumulate on a busy board - decided by Alex in the notes, flagged only so the person sees the cost. 3) Live check only showed 'nothing recorded' rows (running sessions had the older cockpit without 'repo'), so the full row with context/cache/verdict is proven by tests (tests/board.test.ts 'the workers of the board'), not yet seen live. No defects found in the code itself; claude plugin test: 81 pass; go test ./core/role ./core/mod ./core/release ok."
test-verdict: "pass: go build/vet, go test -race -count=1 ./... RC=0, claude plugin test 81 pass RC=0, plugin validate passed; DoD 1-5 verified in the tree; /cockpit workers exercised live with --plugin-dir"
review-verdict: "The diff satisfies DoD 1-5 and the outcome claims match it; no code defects found. One reservation: Give back without --force assumes jaira >=0.3.7, which the marketplace install route does not guarantee (gap 1) - acceptable to ship, worth a README line or follow-up."
review-check: "1. Run 'jaira roles install' from this branch's build (go build -o /tmp/jaira ./cmd/jaira && /tmp/jaira roles install), so running sessions get this cockpit. 2. Start a worker on any ticket via the teamlead (it opens a Herdr tab labelled <ID>/worker) and let it answer once. 3. In the teamlead session type '/cockpit workers'. You should see a first line 'this session: context Nk of Mk, cache warm while it works and N min after its last answer', then 'Agents carrying a ticket of this board:' and a line '<ID> · <lane> · worker · context Nk of 1000k · cache N min left · next round: same worker (...)', or 'working: nothing to decide' while the worker is busy. 4. Open the cockpit settings and set workerContextShare to 1; run '/cockpit workers' again: the verdict now reads 'next round: new worker (old worker's context Nk >= 10k (1% of its window))'. 5. On a ticket in signoff press Give back in the band and give a reason: the ticket goes to in-progress and 'jaira show <ID>' notes whom it went to; no --force is used. Automated: cd core/mod/builtin/jaira-cockpit && claude plugin test . (81 pass)."
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
- **2026-10-12 18:30 · Alexander Sacharov** — testing, live: claude -p --plugin-dir core/mod/builtin/jaira-cockpit '/cockpit workers' in this worktree printed the self line and two worker tabs of this board (776466 signoff, QRHB35 testing), both 'nothing recorded': the running sessions still have the installed older cockpit, whose records carry no 'repo' and are filtered out. Expected for the transition; each worker rewrites its record on its next turn once 'jaira roles install' put this build in place.
- **2026-10-12 18:35 · Alexander Sacharov** — review: verdict pass with one reservation - Give back drops --force unconditionally, but the marketplace route (.claude-plugin/marketplace.json) can run the cockpit beside a jaira older than 0.3.7, where the move out of signoff is refused after the 'Given back' note is already written. Not sent back; for the person to decide (README line or version check).
- **2026-10-12 18:37 · Alexander Sacharov** — explain: Review sagt bestanden – DoD 1–5 erfüllt, kein Code-Fehler. Ein Vorbehalt: Give back schickt jetzt ohne --force zurück und setzt jaira >= 0.3.7 voraus. Über marketplace.json lässt sich das Cockpit aber neben einem älteren jaira installieren; dort wird der Zug aus signoff verweigert – NACHDEM die Notiz 'Given back from signoff' schon auf dem Ticket steht. Das Ticket behauptet dann eine Rückgabe, die nicht passiert ist (Fehler-Toast kommt, Ticket bleibt stehen). Optionen: (a) so ausliefern, (b) eine README-Zeile 'Cockpit braucht jaira >= 0.3.7', (c) Versionsprüfung oder Notiz erst nach erfolgreichem move schreiben. Reviewer empfiehlt: ausliefern, Lücke per README-Zeile oder Folgearbeit schließen. Nur zur Kenntnis: offene Worker-Tabs sammeln sich jetzt an (deine Entscheidung), und die volle Worker-Zeile ist nur per Test belegt, live noch nicht gesehen.
