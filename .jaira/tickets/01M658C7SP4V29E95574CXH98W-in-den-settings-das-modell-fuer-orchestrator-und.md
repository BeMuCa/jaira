---
id: 01M658C7SP4V29E95574CXH98W
title: "In den Settings das Modell fuer Orchestrator und Worker waehlen; Default Opus 5.5 fuer den Orchestrator, Haiku (high) fuer Worker"
status: in-progress
ready: true
creator: berkc
assignee: berkc
goal: "Der Teamlead (Orchestrator) und die Worker starten mit den Modellen, die in den Settings stehen; ohne Einstellung Opus 5.5 fuer den Teamlead und Haiku mit hohem Effort fuer Worker."
context: |-
  Wunsch von berkc am 30.10.2026: 'add that you can decide the model of orchestrator and worker agents in Settings. For example i want opus 5.5 for my orchestrator and then haiku on high for my workers. (make that default)'
  Heute: Lanes tragen eine Modellstufe (tier strong/cheap); Settings kennt nur 'Chat model' fuer die Cockpit-Chats (internal/tui/settings.go, .jaira/chat-model). spawn.sh startet Worker - pruefen, ob es ein Modell/Effort uebergeben kann.
  Offen: Verhaeltnis zu Lane-Tiers - ueberschreibt die Worker-Einstellung den Tier, oder ist sie nur der Default fuer 'cheap'?
definition-of-done: "Settings hat Zeilen 'Orchestrator model' und 'Worker model' (Modell + Effort), gespeichert in settings.json; Test"
tags:
  - tui
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-30T01:56:48Z
updated-at: 2026-10-31T02:40:36Z
updated-by: berkc
claimed-by: berk-ca-3129445
claimed-at: 2026-10-31T01:27:42Z
---

# In den Settings das Modell fuer Orchestrator und Worker waehlen; Default Opus 5.5 fuer den Orchestrator, Haiku (high) fuer Worker

## Definition of Done

- [ ] Settings hat Zeilen 'Orchestrator model' und 'Worker model' (Modell + Effort), gespeichert in settings.json; Test
- [ ] spawn.sh / Teamlead startet Worker mit dem eingestellten Modell und Effort; Teamlead-Skill nennt sein Modell
- [ ] Ohne Einstellung: Opus 5.5 Orchestrator, Haiku high Worker
- [ ] Release-Note-Zeile

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/settings/settings.go: add type AgentModel struct{Model string json:model,omitempty; Effort string json:effort,omitempty} and Settings field Models map[string]AgentModel json:"models,omitempty" (keys orchestrator, worker, strong, cheap). Add var DefaultModels = {orchestrator:{opus,""}, worker:{haiku,high}, strong:{opus,""}, cheap:{haiku,high}} and func (s Settings) ModelFor(key string) (AgentModel, source string): entry present in s.Models wins whole (model via existing modelName(), effort kept only if in low|medium|high|xhigh|max), else DefaultModels[key] with source 'the default'; unknown tier key -> worker entry. Add var ModelChoices = [{opus,''},{opus,high},{sonnet,''},{sonnet,high},{haiku,''},{haiku,high}] and func (m AgentModel) String() 'opus' / 'haiku high'. Aliases only (opus/sonnet/haiku): the Agent tool takes no full ids
- [x] core/settings/settings_test.go: TestModelsDefault (no file -> orchestrator opus, worker haiku high, strong opus, cheap haiku high), TestModelsFromFileWinWhole (worker {sonnet} -> effort empty, not high), TestModelsRejectBadValues (model 'x;rm' -> default; effort 'huge' -> ''), TestSaveKeepsModels round trip. go test ./core/settings -count=1
- [x] internal/cli/flow.go --for-lane: next to model_tier add "model" and "effort" from settings.Load().ModelFor(l.ModelTier), only when l.ModelTier != "". Plain-text head: 'tier: strong (opus)' / 'tier: cheap (haiku high)'. Per-ticket model-tier field is NOT read (unchanged). Test in internal/cli: a strong lane's --for-lane --json carries model opus, effort ''; with JAIRA_HOME settings.json models.strong={sonnet,high} it carries sonnet/high
- [x] internal/cli/roles.go 'roles list': --json adds orchestrator_model, orchestrator_effort, worker_model, worker_effort (from ModelFor); text output one line 'models: orchestrator opus, worker haiku high, strong opus, cheap haiku high (models in <settings path>)'; extend the command's Long help with one paragraph. Test the json keys
- [x] internal/tui/settings.go: four rows after 'Chat model': 'Orchestrator model', 'Worker model', 'Strong lanes model', 'Cheap lanes model' (action settingsActionModel + key). Value shows ModelFor(key).String() plus ' (default)' when not in the file. enter steps: absent -> ModelChoices[0..5] -> absent (delete key; drop map when empty), saved at once via Load/modify/Save like switchDispatch. desc says 'model + effort for <...>; any other by hand in settings.json'. Test in internal/tui/settings_test.go like TestSettingsEnterStepsTheChatModel (JAIRA_HOME temp, step through, check settings.Load().Models)
- [ ] core/role/builtin/jaira-dispatcher/scripts/spawn.sh: compute model/effort before building $claude: lane dispatch|worker -> worker_model/worker_effort from 'jaira -C "$wt" roles list --json | jaira json get worker_model --default ""'; explain|discuss -> chat_model as today (effort none); pr -> none; any other lane -> 'jaira -C "$wt" show "$ticket" --for-lane "$lane" --json | jaira json get model --default ""' and effort likewise. Same regex guard for model, effort must match ^(low|medium|high|xhigh|max)?$. Append --effort to all three start lines (claude var and the windows line). Update usage text one sentence
- [ ] core/role/builtin/jaira-worker/SKILL.md: loop step 2 -> before working a lane compare payload model/effort with worker_model/worker_effort from 'jaira roles list --json'; differs -> the lane is not yours: hand it to a new agent exactly as 'The fresh lane' (run-lane.sh under Herdr, spawn.sh picks the model itself; else Agent(subagent_type: claude, model: <model>, effort: <effort>, prompt: jaira-role-lane ...)). Fresh-lane subagent line also passes model/effort
- [ ] core/role/builtin/jaira-dispatcher/SKILL.md: Subagents path (around line 364) -> Agent per lane with model/effort from the lane's --for-lane payload; one sentence that run-lane/spawn.sh already pick it. core/role/builtin/jaira-teamlead/SKILL.md: section naming its own model: orchestrator_model from roles list; if the session runs another model say once 'restart with claude --model <m> [--effort <e>]', never switch silently; no-Herdr Agent(...) call for dispatcher/worker gets model: worker_model, effort: worker_effort
- [ ] shipped hashes: run go test ./core/role ./core/shipped ./core/lane -count=1 after the SKILL.md/spawn.sh edits; TestEveryInstalledFileIsAShippedVersion prints the missing lines -> APPEND them to core/shipped/shipped.sha256 (never remove old hashes). bash -n spawn.sh
- [ ] core/release/notes/CXH98W.md, one line each, never wrapped: '- Pick the orchestrator, worker, strong-lane and cheap-lane model (and effort) under Settings, or in "models" in ~/.jaira/settings.json; without it the teamlead runs on opus, workers and cheap lanes on haiku high, strong lanes on opus.' and '- ! Run jaira update --roles so workers start each lane on the model of its tier; jaira show --for-lane --json now carries model and effort.'
- [ ] verify: go test ./core/settings ./internal/cli ./internal/tui ./core/role ./core/shipped -count=1 (skip the 3 hangers); tick DoD with proofs; commit feat(CXH98W) with the ticket file

## Progress
- **2026-10-30 21:39 · berkc** — berkc 30.10.2026: Effort muss mit einstellbar sein (Modell + Effort). Pro Lane ein festes Modell: Lanes tragen heute schon 'tier' (strong/cheap, siehe CLAUDE.md-Lane-Liste). Vorschlag statt neuem Ticket-Attribut: Settings legt fest, welches Modell+Effort zu jedem Tier gehoert (z.B. strong = Opus 5.5 high, cheap = Haiku high) plus Orchestrator-Modell; eine Lane darf zusaetzlich ein festes 'model:' tragen. Weicht das Modell einer Lane vom Modell des laufenden Workers ab, uebergibt der Worker diese Lane an einen neuen Agenten mit frischem Kontext (wie heute bei context: fresh). Offen fuer berkc: zusaetzlich ein model-Feld pro Ticket, oder reicht Lane + Settings?
- **2026-10-30 23:45 · berkc** — berkc 31.10.2026 (Punkt 8): Die review-Lane soll immer ein starkes Modell nutzen - Opus 5.5. Heute hat review 'tier strong' (jaira lanes), aber welches Modell 'strong' bedeutet, legt nichts fest; genau diese Zuordnung Tier -> Modell+Effort soll dieses Ticket in den Settings liefern (Default strong = Opus 5.5).
- **2026-10-31 02:34 · berkc** — Pre-process decisions (planner, 31.10.2026): (1) Worker model and tier mapping are separate entries in one settings.json map 'models' {orchestrator, worker, strong, cheap}, each {model, effort}; an entry wins whole, missing -> default (opus / haiku high / opus / haiku high). Answers the context's open question: the worker setting does NOT override tiers; it is the model the jaira-worker session is started on, and a lane whose tier model differs is handed to a new agent on that model (same procedure as a fresh lane). With defaults: worker=cheap=haiku high, so in-progress/testing run in the worker, brainstorm/pre-process/critique/optimize/review go to opus subagents. (2) Aliases only (opus/sonnet/haiku), not full ids: the Agent tool's model param takes aliases; assumption, unverified: alias 'opus' resolves to Opus 5.5 today. (3) Machine-level only (~/.jaira/settings.json), no board file like chat-model: maintainer asked for settings.json; jaira update does not write settings.json (only Load/Save in core/settings), and adding the field to the struct keeps it across every Save. (4) No per-ticket model override (YAGNI). Found: tickets already carry a model-tier field (jaira create --tier, core/ticket/schema.go FieldModelTier) that --for-lane ignores (flow.go reports l.ModelTier only) - left as is. Open for berkc: should a ticket's own model-tier override the lane's tier? (5) Teamlead cannot change its own model; its skill names orchestrator_model and tells the human the restart line if it differs. (6) The dispatcher session itself starts on worker model; claude CLI has --effort low|medium|high|xhigh|max (checked claude --help).
