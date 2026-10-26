---
id: 01M5HBGDHKJ80VV09SWCQYK8MD
title: Die Session-ID traegt keinen Rechnernamen mehr
status: testing
ready: true
creator: BeMuCa
goal: "claimed-by und jede andere Stelle, die eine Session benennt, enthalten keinen Hostnamen: ein Teammitglied kann aus einem Ticket nicht ablesen, auf welcher Maschine es bearbeitet wurde."
context: |-
  Am 2026-10-22 fiel auf: core/session/session.go:129 bildet die Session-ID aus os.Hostname() plus Zahl, und jaira claim schreibt sie als claimed-by ins Ticket (internal/cli/claim.go:151).
  Der Hostname von Berks Arbeitsrechner traegt das Kuerzel seines Arbeitgebers als Praefix. Damit steht der Arbeitgeber in 45 Ticketdateien auf master und in 10 Ticket-Refs - die Erwaehnung, die T22Z2C gerade verbieten soll.
  Ein Hostname ist auch sonst nichts, was ein Ticket braucht: er identifiziert die Session nur zufaellig und verraet Infrastruktur.
  Vorschlag: eine zufaellige, pro Rechner einmal in ~/.jaira erzeugte Kennung (oder ein Hash des Hostnamens) statt des Klartexts; die Zahl dahinter bleibt. Alte claimed-by-Werte muessen weiter als 'fremde Session' erkannt werden, damit ein Update keinen Claim verliert (siehe die Regel: Update verliert keine Daten).
  Die Altdaten auf master und in den Refs werden beim History-Rewrite am 2026-10-26 ersetzt; das ist nicht Teil dieses Tickets.
definition-of-done: "Eine neue Session-ID enthaelt keinen Hostnamen; ein Test prueft, dass os.Hostname() nicht mehr in der ID vorkommt"
tags:
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T08:26:42Z
updated-at: 2026-10-26T19:35:25Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: qyk8md-driver
claimed-at: 2026-10-26T19:11:29Z
outcome-what: "Session ids are now a random six-character token kept once per machine in machine-id under JAIRA_HOME or ~/.jaira, plus the parent pid; the hostname is only compared, never written. An old-format claim that names this machine counts as abandoned at once."
outcome-why: "claimed-by is committed with the ticket, so a hostname in it published the machine name, and on the maintainer's machine his employer, in ticket files and refs."
outcome-resolves: "No new claimed-by carries a hostname; old claims from other machines still block until ClaimTTL; this machine's own old claims need no --steal after the update."
review-summary: none
round: 3
review-gaps: |-
  pass 2 (4529c17b only): nothing new to remove. Both pass-1 fixes landed: core/session/session.go:170-171 LegacyLocal now reads 'return id == prefix' (same results, sanitize caps at 48 ASCII bytes); core/session/machine_test.go:43-46 adds the two host-built negatives only when the sanitized hostname is at most 43 characters (44 would also be safe, critique pass 3; a conservative bound, not a cleanup)
  let stand (pass 1, not reopened): sanitizeForTest (internal/cli/staleclaim_test.go:159) is a deliberate frozen copy of the format old builds wrote; the JAIRA_HOME-else-~/.jaira lookup in machineToken is the seventh private copy with no shared helper, merging them is its own ticket; LegacyLocal sits in both ClaimActive and StaleClaim (internal/cli/claim.go:29 and :54) because those two were already separate; internal/cli/flow.go:398 reads machine-id once per actionable ticket, hoisting it is optional; Session.Host (core/session/session.go:28) is pre-existing dead code, left alone
test-verdict: "pass: gofmt/vet/build clean, race tests core/session + cli (LegacyHostname|Claim|Stale) ok, release note is one line; token claims, foreign/own-old/stale old-format claims behave on a scratch board; full suite pending pre-PR run"
---

# Die Session-ID traegt keinen Rechnernamen mehr

## Definition of Done

- [x] Eine neue Session-ID enthaelt keinen Hostnamen; ein Test prueft, dass os.Hostname() nicht mehr in der ID vorkommt
  proof: core/session/session.go:126 DefaultID uses machine token; core/session/machine_test.go:10 TestDefaultIDCarriesNoHostname
- [x] Ein Ticket mit einer alten, hostnamenhaltigen claimed-by gilt weiter als von einer anderen Session beansprucht; ein Test deckt den Fall
  proof: internal/cli/staleclaim_test.go:124 TestLegacyHostnameClaims (foreign old-format live claim stays active; this machine's old-format claim is stale)
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-23 15:26 · BeMuCa** — Feedback 23.10. (Punkt 12): auf dem Board stehen 7 verwaiste Claims mit Rechnernamen drin. Beim Fix mitdenken: ein verwaister Claim mit altem Format soll beim naechsten claim desselben Tickets als abgelaufen gelten (stale), damit niemand --steal braucht.
- **2026-10-26 19:14 · BeMuCa** — Driver: pre-process skipped (option planning not ticked). Design handed to in-progress: random 6-char token once per machine in JAIRA_HOME or ~/.jaira, DefaultID = token-ppid, env overrides unchanged. Orphaned old-format = a claim in the old hostname-digits form that names THIS machine: no session here can renew it any more, so it counts as stale even within ClaimTTL. Old-format claims from other machines stay foreign until ClaimTTL, as before (plain string compare). The old claims on master are all older than ClaimTTL and were already stale by time. Claimed with JAIRA_SESSION_ID set to a neutral name so this ticket carries no hostname.
- **2026-10-26 19:20 · BeMuCa** — critique (pass 1, whole diff): three findings, all clear fixes, see review-summary. (1) internal/cli/claim.go StaleClaim: drop the separate LegacyLocal branch and add '|| session.LegacyLocal(t.ClaimedBy)' to the existing 'd > ClaimTTL' test - same results, the zero-timestamp branch above already returns 0. (2) core/session/session.go machineToken: use ticket.Handle(ticket.NewID(time.Now())) instead of the [20:] slice and its comment. (3) core/session/machine_test.go: look for sanitize(host), not the raw hostname - the old build wrote the sanitized form, so a dotted hostname (name.local, FQDN) passes the current test on the old code too. Let stand, so later passes need not reopen them: LegacyLocal treating any this-machine '<host>-<digits>' claim as abandoned at once, standing on BeMuCa's note of 23.10 that asked for exactly that (an old build still running on the same machine would be walked over, but self upgrade replaces the binary in place, so that overlap is near zero); the token written with plain os.WriteFile rather than O_EXCL create - two jaira processes starting in the same microseconds on a machine with no machine-id yet could each keep a different token for one call, which self-heals on the next read and was judged not worth the extra code; sanitizeForTest in internal/cli/staleclaim_test.go copying session.sanitize is left to the optimize lane.
- **2026-10-26 19:23 · BeMuCa** — critique (pass 2, only 434df5e8 and the three pass-1 findings): all three landed - (1) internal/cli/claim.go StaleClaim folds LegacyLocal into the 'd > ClaimTTL' test, zero-timestamp branch still returns first, same results; (2) core/session/session.go machineToken uses ticket.Handle(ticket.NewID(...)), Handle is the last six characters (core/ticket/id.go:88); (3) core/session/machine_test.go compares against sanitize(hostname). Nothing new; the three let-stand items of pass 1 stay closed. gofmt/build/vet green.
- **2026-10-26 19:29 · BeMuCa** — optimize: four passes over ded0919b+434df5e8. Duplication: searched for per-user-dir lookups (JAIRA_HOME/UserHomeDir), random-token makers (core/milestone/id.go NewID makes six random characters too, but it would pull the milestone package into session for a different idea; ticket.Handle(ticket.NewID) reuses a dependency session already has), digit-only helpers (none; ticket.isNumericLike also allows signs and dots), and sanitizer copies (core/tag Normalize refuses characters instead of replacing them, so it is a different rule). Kept sanitizeForTest: it is a frozen copy of what old builds wrote, and folding it would need session.sanitize exported for a test or a DefaultID detour that must also clear CLAUDE_SESSION_ID. Dead code: the change orphaned nothing; Session.Host (core/session/session.go:28) is pre-existing dead, set and read nowhere, left alone, and it writes no hostname anywhere. Fluff: only the prefix[:48] slice in LegacyLocal. Cost: LegacyLocal runs only for live foreign claims, after the TTL test short-circuits; the machine-id read per iteration at internal/cli/flow.go:398 is negligible, left. Found on the way, not a cleanup: TestLegacyLocal breaks on hostnames of 46+ characters, see review-gaps. Lane prompt says this lane edits; the driver said report only, so no code changed.
- **2026-10-26 19:31 · BeMuCa** — critique (pass 3, only 4529c17b and the two optimize review-gaps fixes): both landed, nothing new. (1) core/session/machine_test.go TestLegacyLocal adds the two host-built negatives only when len(sanitize(host)) <= 43; a scratch copy of sanitize/LegacyLocal over host lengths 1-52 shows the test passes at every length (44 would also be safe; 45 is the first length where sanitize(host+'-12x4') is cut to an accepted id, so 43 is one conservative step, not a defect). (2) core/session/session.go LegacyLocal 'return id == prefix': sanitize caps at 48 ASCII bytes, so the old prefix[:48] was prefix, same results. Standing on critique pass 1/2 notes and the optimize let-stand list: nothing reopened. go test ./core/session -run 'TestLegacyLocal|TestDefaultID' green, gofmt clean.
- **2026-10-26 19:33 · BeMuCa** — optimize (pass 2, only 4529c17b): both pass-1 fixes landed (LegacyLocal 'return id == prefix'; TestLegacyLocal host-built negatives only for sanitized hostnames <= 43 chars); the commit adds no duplication, dead code, fluff or cost; pass-1 let-stand items not reopened. go test ./core/session -run 'TestLegacyLocal|TestDefaultID' green, gofmt and vet clean. No code changed (driver said report only).
