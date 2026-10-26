---
id: 01M5HBGDHKJ80VV09SWCQYK8MD
title: Die Session-ID traegt keinen Rechnernamen mehr
status: critique
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
updated-at: 2026-10-26T19:23:44Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: qyk8md-driver
claimed-at: 2026-10-26T19:11:29Z
outcome-what: "Session ids are now a random six-character token kept once per machine in machine-id under JAIRA_HOME or ~/.jaira, plus the parent pid; the hostname is only compared, never written. An old-format claim that names this machine counts as abandoned at once."
outcome-why: "claimed-by is committed with the ticket, so a hostname in it published the machine name, and on the maintainer's machine his employer, in ticket files and refs."
outcome-resolves: "No new claimed-by carries a hostname; old claims from other machines still block until ClaimTTL; this machine's own old claims need no --steal after the update."
review-summary: none
round: 2
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
