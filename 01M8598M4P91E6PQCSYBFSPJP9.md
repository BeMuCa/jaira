---
id: 01M8598M4P91E6PQCSYBFSPJP9
title: "Eine Lane-Datei kann mit agents: N sagen, dass N Agenten die Lane laufen; show --for-lane reicht es weiter"
status: human
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Eine lesende Lane laesst sich in ihrer Lane-Datei mit agents: N besetzen; jaira prueft den Wert, reicht ihn an den Agenten weiter, und die Rollen-Skills sagen, wie N Ergebnisse zu einem werden."
context: |-
  Heute laeuft jede Lane mit genau einem Agenten; es gibt keinen Schluessel dafuer (ATDG74, Idee von Berk: 5x Haiku in critique).
  jaira startet selbst keine Agenten. Der Schwarm ist deshalb: ein Schluessel in der Lane-Datei, den 'jaira show <id> --for-lane <lane> --json' ausgibt, und ein Absatz in den Skills jaira-worker und jaira-role-lane (core/role/builtin/).
  Entschieden im Brainstorm zu ATDG74 (Teamlead, 2026-11-23):
  - Schluessel 'agents: N' im Frontmatter, ganze Zahl 1 bis 9, fehlt er, gilt 1 und nichts aendert sich.
  - Erlaubt nur auf lesenden Lanes. agents > 1 auf einer Lane, deren output-produces outcome-what oder plan enthaelt, ist ein Fehler beim Laden - N Agenten im selben Worktree kollidieren.
  - Zusammenfuehren: Vereinigung der Befunde; der tragende Agent prueft jeden Befund am Code, streicht Doppelte und falsche, und schreibt EIN Ergebnis und EIN Verdict. Keine Mehrheit (verliert den Befund, den nur einer fand), kein eigener Zusammenfasser (prueft nichts).
  - Speicher: vor jedem Start free -h; unter 2 GiB verfuegbar nacheinander statt parallel; nie mehr als 4 zugleich. Das steht im Skill-Text, das Binary startet nichts.
  Code: core/lane/lane.go (Lane-Struct, Parser um Zeile 333-440, max-rounds ist das Muster fuer einen Zahlen-Schluessel; Gleichheit um Zeile 960), core/lane/facts.go, die Ausgabe von show --for-lane in internal/cli (dort wo model_tier gesetzt wird), 'jaira lanes show'.
  Achtung: der Zweig origin/feat/ENCFVZ-skill-lane (nicht gemergt) aendert dieselben Stellen in lane.go; klein halten, damit der Konflikt klein bleibt. Aendert sich eine Datei unter core/role/, verlangt ein Test eine neue core/shipped/shipped.sha256 - nachlesen, wie die erzeugt wird.
definition-of-done: "agents: N in einer Lane-Datei wird gelesen; 0, negative Werte, mehr als 9 und Nicht-Zahlen werden beim Laden mit einer Meldung abgelehnt, die Datei und Schluessel nennt; Test in core/lane"
tags:
  - concurrency
  - cli
blocked-by: []
related: []
follows: 01M67K72RDJY2N2R2KQGATDG74
commits:
  - 71f812811bd997015ae675670c35cb4dbceefa15
created-at: 2026-11-23T22:43:41Z
updated-at: 2026-11-23T23:01:07Z
updated-by: BeMuCa
claimed-by: worker-fspjp9
claimed-at: 2026-11-23T22:56:02Z
outcome-what: "agents: N is read from lane files (1..9, default 1), refused above 1 on outcome-what or plan lanes, and carried in show --for-lane and lanes show; both role skills say how to unite N agents' findings"
outcome-why: "the lanes that judge a ticket can then run several readers at once, without N agents writing the same worktree"
outcome-resolves: "class checked: bad values (0, -1, 10, x, 2.5, empty, huge) refused naming file and key; writing lanes refused; concurrency: jaira starts nothing, the skill caps at 4 and checks free -h; error path: load refuses before any lane is used; permissions: not applicable, no process started"
review-summary: "core/lane/lane.go agents refusal checks output-produces for outcome-what or plan only; core/lane/builtin/26-optimize.md has output-produces [review-gaps] but says 'Make the changes. This lane edits', so agents: 3 on a copied optimize lane passes the load and three agents edit one worktree. The ticket's rule (read-only lanes only) is not enforced for it. Needs a person's choice, see the question."
question: "agents above 1 passes the load on a lane that edits code. optimize produces only review-gaps, so the check (outcome-what or plan) cannot see that it writes. Options: (a) keep the check and say in the docs and the role skills that agents above 1 is for read-only lanes only; nothing stops a person who copies optimize and adds the key. Cheapest, no new key. (b) add a lane key that says the lane writes, and refuse agents above 1 on it; one more key to validate and document. (c) refuse agents above 1 on the lane with id optimize by name; a special case in code. Recommendation: (a) now, because the shipped optimize lane sets no agents key and a person has to add it on purpose; (b) if a copied lane is expected to be common. Which?"
---

# Eine Lane-Datei kann mit agents: N sagen, dass N Agenten die Lane laufen; show --for-lane reicht es weiter

## Definition of Done

- [x] agents: N in einer Lane-Datei wird gelesen; 0, negative Werte, mehr als 9 und Nicht-Zahlen werden beim Laden mit einer Meldung abgelehnt, die Datei und Schluessel nennt; Test in core/lane
  proof: core/lane/agents_test.go TestParseAgents; core/lane/lane.go Agents parser
- [x] agents > 1 auf einer Lane, die outcome-what oder plan liefert, wird beim Laden abgelehnt; Test
  proof: core/lane/agents_test.go TestParseAgentsRefusedOnWritingLane; core/lane/lane.go refusal
- [x] jaira show <id> --for-lane <lane> --json traegt agents (1 ohne Schluessel); die Textausgabe und jaira lanes show nennen es nur, wenn es groesser als 1 ist; Test in internal/cli
  proof: internal/cli/agents_test.go TestShowForLaneCarriesAgents, TestLanesShowCarriesAgents; internal/cli/flow.go, tickets.go
- [x] jaira-worker und jaira-role-lane sagen, was bei agents > 1 zu tun ist: N Agenten, Speicherregel, Befunde vereinigen und am Code pruefen, ein Verdict
  proof: core/role/builtin/jaira-worker/SKILL.md 'A lane that runs several agents'; core/role/builtin/jaira-role-lane/SKILL.md 'A lane with several agents'
- [x] docs und Release-Notiz in core/release/notes/<ID>.md
  proof: README.md 'A lane can run several agents'; core/release/notes/FSPJP9.md; core/lane/template.go

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] read the max-rounds path end to end: lane.go parser 422-431, flow.go 750 and 815, tickets.go lanes show 1523 and 1536
- [x] failing tests in core/lane: default 1, 1 to 9 accepted, 0 -1 10 x 2.5 refused with file and key named, agents>1 refused on an outcome-what or plan lane
- [x] failing test in internal/cli: show --for-lane json agents 1 without key and 3 with key; text line only above 1; lanes show
- [x] implement agents in core/lane/lane.go: field next to MaxRounds, parser, equality; no change to facts.go (the board note never shows lane agents)
- [x] implement show --for-lane payload and text in internal/cli/flow.go; lanes show in internal/cli/tickets.go
- [x] skill paragraphs in core/role/builtin/jaira-worker/SKILL.md and jaira-role-lane/SKILL.md, then add the new shipped.sha256 lines with scripts/shipped-history.sh --write
- [x] template.go comment line, README paragraph beside max-rounds, release note core/release/notes/FSPJP9.md
- [x] gofmt, go vet -p 2 on touched packages, targeted tests only, then move to in-progress

## Progress
- **2026-11-23 22:57 · BeMuCa** — pre-process: agents has no runtime effect in jaira; facts.go left alone because the board note never shows lane agents. Default 1 is set in the parser so show prints 1 for an absent key. Lanes that write (outcome-what, plan) reject agents>1 at load, because N agents in one worktree collide.
- **2026-11-23 23:01 · BeMuCa** — critique: agents above 1 is only refused where output-produces names outcome-what or plan. optimize (core/lane/builtin/26-optimize.md) produces review-gaps and edits code ('Make the changes. This lane edits'), so the rule lets agents: N through on it. review and critique produce review-gaps/review-summary and only read, so output-produces cannot tell optimize apart from them. Left alone: the choice is a person's (see question). Everything else in the diff read and found sound: parser, equality, show JSON and text, lanes show, shipped manifest only appended.
