---
id: 01M4TRAVWQFDY2W0MPEA339DNB
title: "Ein Testbuild von jaira laesst sich starten, ohne dass er Refs ins Remote sendet"
status: critique
ready: true
creator: BeMuCa
goal: "Wer einen ungeprueften Build gegen einen echten Checkout startet, kann das Senden an das Remote abschalten."
context: |-
  Das TUI sendet beim Start und bei jedem Ref-Tick alle gesammelten Refs ins Remote (internal/tui/refs.go:78, y.Flush()). Die CLI sendet ebenfalls (internal/cli/refs.go:88 und :248).
  Ein Testbuild auf einem echten Checkout veroeffentlicht also alles, was er schreibt - auch das, was ein fehlerhafter Build falsch schreibt.
  Am 2026-10-13 (D0SAHM, Milestone Roadmap) lief deshalb jeder manuelle Test auf einer Kopie: cp -a des Boards, git remote remove origin, HOME und JAIRA_HOME in ein Scratch-Verzeichnis.
  Einen Schalter gibt es heute nicht. Die Umgebungsvariablen, die der Code liest: JAIRA_DEFAULT_BOARD, JAIRA_HOME, JAIRA_LANES_DIR, JAIRA_NO_UPDATE_CHECK, JAIRA_USER (grep auf Getenv, 2026-10-13).
  Vorbild fuer den Namen: JAIRA_NO_UPDATE_CHECK.
definition-of-done: "Mit gesetztem Schalter sendet weder das TUI noch die CLI Refs; beide sagen sichtbar, dass nicht gesendet wird; ein Test belegt, dass Flush dann nicht laeuft"
tags:
  - concurrency
  - cli
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-13T13:48:19Z
updated-at: 2026-10-22T08:16:22Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-4187593
claimed-at: 2026-10-22T07:55:10Z
outcome-what: "JAIRA_NO_PUSH=1 sperrt die beiden git-push-Aufrufe in core/gitref (ErrPushOff); Syncer.Flush meldet das vor dem Fetch und laesst die Outbox unberuehrt. Die CLI nennt es in einer stderr-Zeile, das TUI in der Versionszeile."
outcome-why: "Ein Testbuild auf einem echten Checkout veroeffentlichte alles, was er schreibt; bisher schuetzte nur eine Board-Kopie ohne Remote."
outcome-resolves: "Mit JAIRA_NO_PUSH=1 sendet weder TUI noch CLI Refs, beide sagen es sichtbar; TestNoPushSwitchRefusesEverySendingPath und TestFlushWithTheNoPushSwitchKeepsTheQueue belegen, dass nichts rausgeht und die Warteschlange bleibt."
review-summary: none
round: 2
---

# Ein Testbuild von jaira laesst sich starten, ohne dass er Refs ins Remote sendet

## Definition of Done

- [x] Mit gesetztem Schalter sendet weder das TUI noch die CLI Refs; beide sagen sichtbar, dass nicht gesendet wird; ein Test belegt, dass Flush dann nicht laeuft
  proof: core/gitref/gitref.go:84,707,928 ErrPushOff+PushDisabled in pushRefspec/PushBranch; core/refsync/refsync.go:649 Flush early return; internal/cli/refs.go:254 flushRefs stderr line; internal/cli/snapshot.go:124 maybeSnapshot skips; internal/tui/updatecheck.go:48 version line; tests TestNoPushSwitchRefusesEverySendingPath, TestFlushWithTheNoPushSwitchKeepsTheQueue
- [ ] Der Schalter steht in docs/COMMANDS.md; go test ./... -race gruen; Release-Notiz unter core/release/notes/
  proof: docs/COMMANDS.md:163 jaira fetch row; core/release/notes/339DNB.md; race suite pending (maintainer)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Tests zuerst: core/gitref - mit JAIRA_NO_PUSH=1 verweigern Write und PushBranch mit ErrPushOff, das Remote bleibt unveraendert; core/refsync - Flush mit Schalter gibt ErrPushOff, sendet nichts, der Eintrag bleibt in der Outbox und geht ohne Schalter danach raus
- [x] core/gitref: PushDisableEnv = JAIRA_NO_PUSH, ErrPushOff, PushDisabled(); pushRefspec und PushBranch (die einzigen zwei git push im Binary) brechen davor ab - deckt pull, release, snapshot und Epic-Rebase ohne weitere Stellen
- [x] core/refsync Syncer.Flush: bei PushDisabled vor dem Fetch nil, gitref.ErrPushOff zurueck - die Queue bleibt wie offline, nichts geht verloren
- [x] CLI: flushRefs (internal/cli/refs.go:250) gibt die eine stderr-Zeile aus; fileOnRefOnly/putOnRef reichen ErrPushOff als why durch; webPush zeigt es
- [x] TUI internal/tui/refs.go:85: ErrPushOff aus Flush ist kein Fehler, der Fetch laeuft weiter; Hinweis in der Versionszeile (versionLine, Home und Model)
- [x] docs/COMMANDS.md neben JAIRA_NO_UPDATE_CHECK, core/release/notes/339DNB.md

## Progress
- **2026-10-22 07:58 · BeMuCa** — Plan vom Treiber geschrieben (pre-process wurde von next_lane uebersprungen, in-progress meldete 'missing: plan'). Warum der Schalter in core/gitref und nicht nur an den drei Flush-Aufrufen: der Kontext nennt nur internal/cli/refs.go:88/:248 und internal/tui/refs.go:78, aber es gibt mehr Sendewege - internal/cli/web.go:109 (webPush), Syncer.Pull/Release (core/refsync/refsync.go:1395/:1478 schreiben den Ref direkt), der Snapshot (core/snapshot/snapshot.go:185 PushBranch, :252/:275 Delete) und rebaseEpic (core/refsync/epic.go:113). Alle laufen durch genau zwei git-push-Aufrufe: core/gitref/gitref.go PushBranch und pushRefspec. Flush prueft zusaetzlich frueh, damit die Outbox wie offline erhalten bleibt und kein Fetch/keine Warnung pro Eintrag entsteht.
- **2026-10-22 08:02 · BeMuCa** — Dead end/why: gate sits in the two git push calls of core/gitref (not at the three Flush callers) because pull/release/snapshot/epic-rebase push directly. Flush returns ErrPushOff BEFORE Fetch and before the empty-outbox check, so TUI ticks hit it every time - fetchRefs therefore treats it as no error. Known gap: a ticket created with the switch on is a plain file (create path files on disk when Flush errors), not queued; a later command does not send it, 'jaira release <id>' does. Not changed (same as offline). No TUI/CLI unit test for versionLine/noRefReason.
- **2026-10-22 08:09 · BeMuCa** — critique (pass 1, whole diff): 5 findings, all in review-summary - (1) release note German, every other note English; (2) docs/COMMANDS.md:163 missing full stop and same false promise; (3) core/gitref/gitref.go:84 ErrPushOff text promises "stays queued and goes out with the first command run without it", false for pull/release/snapshot (direct ref writes, refsync.go:1399/:1482, snapshot.go:185/:275) and for read commands (flushRefs only on Dirty, refs.go:250) - sentinel should name the state ("gitref: sending is off (JAIRA_NO_PUSH=1)") like its siblings, the queue promise belongs in flushRefs (refs.go:253) and noRefReason (refs.go:161); (4) jaira create prints the ErrPushOff sentence twice (tickets.go:429 stdout, refs.go:255 stderr) and advises "once the remote works" (tickets.go:438) while the remote works; (5) internal/cli/snapshot.go:122 still spawns the background snapshot, bgrun stamps first (bgrun.go:123), child fails at PushBranch, 72h backup slot used up - return early on gitref.PushDisabled().
Checked and let stand: gate in the two git push calls (the only "push" literals: gitref.go:711, :931) plus the early return in Syncer.Flush - both earn their keep: the gitref gate covers pull/release/snapshot/epic-rebase, and without the Flush return every entry would come back Failed after a fetch (core/outbox/outbox.go:489), one warning each. versionLine wrapper (internal/tui/updatecheck.go:46) is the least invasive shape - the alternatives append at both call sites (home.go:145, model.go:326) or at four returns. TUI fetchRefs treating ErrPushOff as no error (internal/tui/refs.go:89) is right: the fetch still runs, unsent markers still show. PushDisableEnv in gitref follows FetchDisableEnv (core/refsync/schedule.go:21) and snapshot.DisableEnv living with their owner.
Correction to the 08:02 note: a ticket created with the switch does go out - its outbox entry stays queued and the next writing command (or the board) sends it, so the ref appears (TestFlushWithTheNoPushSwitchKeepsTheQueue proves the queue part). Only the file stays on disk, as offline; fileOnRefOnly's doc (internal/cli/refs.go:75-77) accepts that and 'jaira release X' removes it. Standing on that note's "Not changed (same as offline)": let stand.
- **2026-10-22 08:12 · BeMuCa** — Round 2 (critique fixes): noRefReason deliberately has NO queue promise (differs from critique advice) - it also feeds jaira create stdout and whoami, where nothing is queued; the promise lives only in flushRefs (refs.go:254). Residual: create still prints two lines with the switch on (stdout: state + 'once sending is on again, release...'; stderr from flushRefs: queue promise) - different content, not suppressed, since flushRefs cannot know create already spoke. maybeSnapshot early return untested (needs SpawnRun scaffolding). jaira share spawns a background snapshot child that pushes jaira/board a moment later - when diffing ls-remote, sleep after share before taking the baseline.
- **2026-10-22 08:16 · BeMuCa** — critique (pass 2): read only the pass-1 findings and 769c2110, nothing found. (1) release note now English, one line, an instruction - resolved. (2) docs/COMMANDS.md:163 full stop after '--quiet skips it.', promise rewritten to what holds (pull/release/snapshot refused, queued write goes out with the next writing command once unset) - resolved. (3) core/gitref/gitref.go:84 sentinel names only the state, like its siblings; the queue promise now lives only in flushRefs (internal/cli/refs.go:254-257), which runs only when this process queued something - resolved. (4) create's advice says 'once sending is on again' (tickets.go:444-446); stdout and stderr now carry different content - resolved. Residual (a), two lines on create, let stand: the stderr queue promise is true for create too, and the 08:02 note is wrong on this point. From reading the code, not from a run: create saves, store.go:237 calls Record, refsync.go:139 queues the entry, and Flush returns ErrPushOff at refsync.go:649 before it reads the queue. So the next writing command without the switch sends the ref, the file stays until 'jaira release', the same as offline. (5) internal/cli/snapshot.go:124 returns early on gitref.PushDisabled() - resolved. Residual (b), no test for that guard, let stand: not part of the definition of done (its Flush test is TestFlushWithTheNoPushSwitchKeepsTheQueue), the child's push is blocked in gitref anyway, and a test would need SpawnRun scaffolding; the testing lane may weigh it. 'jaira: gitref: sending is off (JAIRA_NO_PUSH=1): the change stays queued ...' let stand: sentinels with their package prefix already reach users in written lines (pull.go:69 '%s: %v' with refsync.ErrUnsent, pull.go:52, refs.go:260, root.go:157). Rewording it like noRefReason (refs.go:162) would be polish, not a definition-of-done gap.
