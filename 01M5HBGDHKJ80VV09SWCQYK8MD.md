---
id: 01M5HBGDHKJ80VV09SWCQYK8MD
title: Die Session-ID traegt keinen Rechnernamen mehr
status: review
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
updated-at: 2026-10-26T19:43:36Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: qyk8md-driver
claimed-at: 2026-10-26T19:11:29Z
outcome-what: "Session ids are now a random six-character token kept once per machine in machine-id under JAIRA_HOME or ~/.jaira, plus the parent pid; the hostname is only compared, never written. An old-format claim that names this machine counts as abandoned at once."
outcome-why: "claimed-by is committed with the ticket, so a hostname in it published the machine name, and on the maintainer's machine his employer, in ticket files and refs."
outcome-resolves: "No new claimed-by carries a hostname; old claims from other machines still block until ClaimTTL; this machine's own old claims need no --steal after the update."
review-summary: |-
  Session ids no longer contain the machine's name. When none of CLAUDE_SESSION_ID, JAIRA_SESSION_ID or TERM_SESSION_ID is set, DefaultID (core/session/session.go) now builds the id as <token>-<parent pid> instead of <hostname>-<parent pid>. The token is six random characters (the random tail of a fresh ULID), made on first use and kept in the file machine-id under JAIRA_HOME, else ~/.jaira; if that file cannot be written, the word 'local' stands in. That id is what 'jaira claim' writes as claimed-by, so a committed ticket now shows e.g. 'claimed-by: KZ78NQ-616260'.
  The hostname is still read, but only to recognise the old form: LegacyLocal (session.go) is true for '<this machine's hostname, sanitized>-<digits>'. ClaimActive and StaleClaim (internal/cli/claim.go) treat such a claim as abandoned at once, because nothing on this machine writes that form any more: 'jaira claim' takes it over without --steal and prints the usual 'carries an abandoned claim by ...' line, and 'jaira next' stops skipping the ticket. An old-format claim from another machine is compared as a plain string, as before, and blocks until the 30-minute ClaimTTL.
  Tests: TestDefaultIDCarriesNoHostname and TestLegacyLocal (core/session/machine_test.go), TestLegacyHostnameClaims (internal/cli/staleclaim_test.go). Release note: core/release/notes/QYK8MD.md, one line. Old claimed-by values already on master are not touched (history rewrite, out of scope).
round: 3
review-gaps: |-
  No blocker. DoD 3 is open on purpose: the full go test ./... -race runs once before the PR (house rule); its release-note half is done.
  1. DoD 2, read literally, says an old hostname claim keeps counting as another session's claim. That holds for other machines. For THIS machine's old claims the code goes further: they count as abandoned at once. That is what the 23.10 note on this ticket asked for, so it is a decision, not a defect - but the DoD text and the code differ here, and signoff should know.
  2. LegacyLocal matches on the hostname alone. Another machine with the same hostname (default names like 'localhost' or 'ubuntu', cloned VMs) that still runs an old build: its live claims look abandoned here. Only while old builds are still around; claims are advisory. (From the code, not reproduced.)
  3. The token is per home directory, not per machine, despite the file name machine-id. A ~/.jaira shared by several machines (NFS home, dotfiles sync) gives them one token; two sessions then collide only if the parent pid is equal too. The old hostname form kept them apart. (From the code.)
  4. The 'local' fallback (no home directory) is the same word on every such machine, so 'local-<pid>' can collide across machines. Rare: the ticket store itself falls back into the repo in that case. (From the code.)
  5. First-use race, let stand by critique: two jaira processes creating machine-id in the same instant can each use a different token once. A claim taken in that call then looks foreign to its own shell until ClaimTTL or --steal - the critique note's 'self-heals on the next read' is true for the token, not for that claim. Very rare; not worth code now. (From the code, not reproduced.)
test-verdict: "pass: gofmt/vet/build clean, race tests core/session + cli (LegacyHostname|Claim|Stale) ok, release note is one line; token claims, foreign/own-old/stale old-format claims behave on a scratch board; full suite pending pre-PR run"
review-verdict: "Ready for signoff once DoD 3's full race suite has run green; by house rule that run happens once before the PR, so DoD 3 is still open and is the only thing missing. The diff does what DoD 1 and 2 ask: no new session id carries the hostname (TestDefaultIDCarriesNoHostname fails when the old DefaultID is put back, checked), foreign old-format claims still block until ClaimTTL, and this machine's own old claims are taken over without --steal, as the 23.10 note asked (TestLegacyHostnameClaims fails without the LegacyLocal check, checked). Gaps 1-5 are edge cases or a decision to confirm; none needs a change before merge. gofmt -l empty, go vet ./... and go build ./... green, go test -race ./core/session and ./internal/cli -run 'LegacyHostname|Claim|Stale' green, scratch-board run matches the review-check below."
review-round: 3
review-check: |-
  Run in one terminal, in order; every step uses the same shell.
  1. S=$(mktemp -d)   -> nothing printed; a scratch folder
  2. cd /home/berk/git/jAIra-QYK8MD && go build -o "$S/jaira" ./cmd/jaira   -> nothing printed
  3. mkdir -p "$S/home" "$S/jhome" "$S/board" && cd "$S/board" && git init -q && git commit -q --allow-empty -m init   -> nothing printed; a git repo with no remote
  4. j() { env -u CLAUDE_SESSION_ID -u JAIRA_SESSION_ID -u TERM_SESSION_ID HOME="$S/home" JAIRA_HOME="$S/jhome" JAIRA_NO_FETCH=1 "$S/jaira" "$@"; }   -> nothing printed; 'j' is now the new binary on the scratch home
  5. j init   -> last line: Run 'jaira share' when you want your team to have it.
  6. j create probe --goal g --context c --dod d   -> first line: Created <HANDLE>  probe
  7. ID=<the HANDLE from step 6>
  8. j claim $ID   -> Claimed <HANDLE> until HH:MM
  9. grep claimed-by .jaira/tickets/*.md   -> claimed-by: <6 characters>-<digits>, for example KZ78NQ-616260
  10. cat "$S/jhome/machine-id"   -> the same 6 characters as in step 9
  11. grep -ci "$(hostname)" .jaira/tickets/*.md   -> 0
  12. j claim $ID --session host-a-4242 --steal   -> Claimed ... (pretends an old build on ANOTHER machine holds the ticket)
  13. j claim $ID; echo $?   -> jaira: <HANDLE> is claimed by host-a-4242 (moments ago); pass --steal to take it anyway, then 3
  14. j claim $ID --session "$(hostname)-4242" --steal   -> Claimed ... (pretends an old build on THIS machine holds it; if 'hostname' prints a dot, use "$(hostname | tr . -)-4242")
  15. j claim $ID; echo $?   -> Claimed <HANDLE> until HH:MM, then: <HANDLE> carries an abandoned claim by <your hostname>-4242, last renewed moments ago, then 0 - no --steal needed
  16. grep claimed-by .jaira/tickets/*.md   -> the token form again, same 6 characters as step 10
  17. cd /home/berk/git/jAIra-QYK8MD && go test ./core/session ./internal/cli -run 'TestDefaultIDCarriesNoHostname|TestLegacyLocal|TestLegacyHostnameClaims' -count=1   -> two lines starting with ok
  18. rm -rf "$S"   -> scratch gone
  DoD 3 is not checked by this: it needs the full go test ./... -race run before the PR.
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
- **2026-10-26 19:35 · BeMuCa** — testing: gofmt -l empty, go vet and go build clean; go test -race core/session ok, internal/cli -run LegacyHostname|Claim|Stale ok; notes/QYK8MD.md is 1 unwrapped '- ' line. Scratch board (no remote, built -race binary): (a) claim wrote claimed-by 6chars-digits, machine-id holds the same 6 chars, second claim accepted RC=0; (b) fresh foreign host-a-4242: claim refused RC=3 'pass --steal', --release refused RC=3 'claim is still live'; (c) this machine's old form fresh: claimed RC=0 without --steal, 'abandoned claim', new claimed-by is token form; (d) foreign old form 2h old: claimed RC=0 without --steal. Full go test ./... pending pre-PR run (DoD 3 open).
- **2026-10-26 19:43 · BeMuCa** — review: no blocker; DoD 3 waits on the pre-PR race suite. Five minor points in review-gaps: (1) DoD 2 text vs code - this machine's old claims count as abandoned at once, as the 23.10 note asked; (2) LegacyLocal matches on hostname alone, so a same-named machine on an old build gets walked over; (3) machine-id is per home directory, a shared ~/.jaira gives several machines one token; (4) the 'local' fallback is the same on every machine without a home; (5) the first-use race can leave one claim looking foreign to its own shell until ClaimTTL. None needs a change before merge.
