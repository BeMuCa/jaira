---
id: 01M7WWYG213BWB06A7D1CRASZS
title: "Ein offener PR wird begleitet: jaira sieht neue Kommentare, CI-Stand und wartende Re-Requests und meldet sie"
status: optimize
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Pro Ticket mit PR zeigt jaira den Stand (Kommentare offen/beantwortet, CI, Review angefordert) und meldet einer Session, wenn sich daran etwas aendert, damit die naechste Runde startet."
context: |-
  Heute: den Stand eines PR (z.B. "3 offen, 2 beantwortet, Re-Request ausstehend") haelt der Teamlead im Kopf; jaira weiss von PRs nichts ausser in jaira-role-pr.
  Alex 20.11.2026: Automatik, um einen PR zu begleiten, und Moeglichkeiten, ihn zu verfolgen.
  Kein Daemon (Grundsatz von jaira): Abfrage beim Board-Refresh, im Cockpit oder per Befehl; gh/glab als Quelle.
  Verwandt: QFWRXH (PR-Spalte auf dem Board), das Ticket mit der optionalen PR-Runde.
definition-of-done: "Ein Befehl oder eine Cockpit-Ansicht zeigt je Ticket mit PR: offene/beantwortete Kommentare, CI-Stand, ob ein Review aussteht"
tags: []
blocked-by:
  - 01M5ZJ5DKQZ83GYBGPJHAX8K02
related:
  - 01M7WWYB5XZMD01DZY3MJG3REW
  - 01M6HZQPZZFAF3MXY3AWQFWRXH
  - 01M7WX040JC9Z202EH852EF48S
commits: []
created-at: 2026-11-20T16:34:31Z
updated-at: 2026-11-22T16:50:16Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3910889
claimed-at: 2026-11-22T14:48:02Z
outcome-what: "Neuer Befehl 'jaira pr [<branch>] [--all|--number n] [--json] [--watch [--every d]]' (core/pr: gh-GraphQL bzw. glab-API; Threads wartend/beantwortet/erledigt, fremde Kommentare, CI, Reviews, Tickets, settled). --watch: Ausgangsstand als erste Zeile, dann eine Zeile je Aenderung, Ende bei merged/closed. Cockpit: reicht einen neuen PR einmal an die Teamlead-Session weiter (nach Accept begrenzte Nachschau 20x60s; Agentenbericht mit PR-URL), Store-Schluessel gegen Doppeltmelden; Accept bei pr-when=signoff oeffnet erst den PR (Ticket bleibt in signoff) und laesst done erst bei settled zu."
outcome-why: "Den Stand eines PR hielt der Teamlead im Kopf; jaira wusste nichts davon. Alex 22.11.: Befehl im Binary, kein Dauer-Polling im Cockpit, Teamlead beobachtet selbst per Monitor; JG3REW-Rest (Accept in signoff) gehoert hierher."
outcome-resolves: "DoD 1-7. Vier Kategorien geprueft: Nebenlaeufigkeit - zwei Cockpit-Sessions koennen denselben PR im selben Moment beide melden (Store get/set nicht atomar), schlimmstenfalls eine Zeile doppelt; Fehlerpfade - halber Lauf ist ein Fehler, nie ein PR ohne Threads; Watch sagt einen Lesefehler einmal und liest weiter, verlorener Login beendet ihn (Exit 3); ungewoehnliche Eingabe - Branchnamen gehen als ein Argument an gh/glab bzw. URL-escaped an die GitLab-API, Fenster: 50 PRs, 100 Threads, 20 Kommentare je Thread; Rechte - nur Lesen mit dem gh/glab-Login der Person."
review-summary: "none"
round: 3
back-from: testing
back-by: lane
back-count: 1
review-gaps: "round 2: values() (one caller) folded into gitlab.go pages[T]; pr.go imports trimmed. Left as before: handle regex vs gitrepo handleShape, Go/TS forge detection, pr-round.sh settled (critique note)."
test-verdict: "pass: go vet ./... + GOOS=windows vet clean; go test ./... RC=0, -race on core/pr core/role core/shipped internal/cli RC=0; claude plugin test (cockpit) 313 pass. DoD 1-7 verified in the tree (proof files/tests opened and run). Live: jaira pr on BeMuCa/jaira lists #184/#183 with CI and threads, --number 91 shows 1 answered thread; --watch --number 179 prints the baseline line and 'merged', exit 0. Failure paths: GH_TOKEN=invalid -> exit 3 'not signed in' (also under --watch), no gh on PATH -> exit 3, repo without origin -> exit 3; no PR -> exit 0 with empty prs. jaira pr writes nothing; the cockpit's one write is a store key per pull request (pr-followed:<project>:<n>), two sessions at once give at most one duplicate line."
hunt-verdict: "fail: 7 change findings (2 blockers: first-page-only reads make settled wrong; CI green again never fires across pending)"
---

# Ein offener PR wird begleitet: jaira sieht neue Kommentare, CI-Stand und wartende Re-Requests und meldet sie

## Definition of Done

- [x] Ein Befehl oder eine Cockpit-Ansicht zeigt je Ticket mit PR: offene/beantwortete Kommentare, CI-Stand, ob ein Review aussteht
  proof: jaira pr / --all / --json: internal/cli/pr.go printPRs + core/pr/watch.go Summary (threads waiting/answered/resolved, CI, review pending); live: jaira pr --number 91 -> 1 answered thread
- [x] Kommt ein neuer Kommentar oder wird CI rot, erfaehrt es die Session des Teamleads ohne Nachfragen
  proof: jaira pr --watch (internal/cli/pr.go watchPR, core/pr/watch.go Changes+Carry) under Monitor, started on the cockpit's followPr line; tests internal/cli/pr_test.go TestPRWatchReportsEachChangeOnceAndEnds, core/pr/pr_test.go TestChanges, TestChangesCIAcrossPending
- [x] 'jaira pr [--json]' zeigt den Stand des PR der Branch (neue Kommentare/Threads, CI, Reviews) über gh und glab (Test) - Alex 22.11.
  proof: internal/cli/pr.go newPRCmd; core/pr/github.go readGitHub, core/pr/gitlab.go readGitLab; tests core/pr/pr_test.go TestReadInputClasses+TestReadGitLab, internal/cli/pr_test.go TestPRShowsTheCheckedOutBranch
- [x] Kein Dauer-Polling im Cockpit: wird ein PR geöffnet (Cockpit-Accept oder Agent), schreibt das Cockpit eine Zeile 'PR #N für <ID> offen – begleiten' in die Teamlead-Session; der Teamlead startet selbst eine Beobachtung über 'jaira pr' (eine Zeile je Änderung, endet bei merged/closed). Nach Accept bekommt der Teamlead ebenfalls eine Zeile (Test) - Alex 22.11.
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx startPrLane->followPullRequest/lookForPullRequest (60 s, 20 looks), shorten()+followLine (report with PR URL), toSessionAccepted unchanged; tests core/mod/builtin/jaira-cockpit/tests/board.test.ts 'following a pull request (CRASZS)', core/mod/builtin/jaira-cockpit/tests/cockpit.test.ts 'a report naming a pull request it opened...'
- [x] Cockpit-Accept öffnet bei pr-when=signoff den PR schon in signoff und wartet vor done auf 'settled' (JG3REW-Rest, Test) - Alex 22.11.
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx heldInSignoff (called first in accept); tests core/mod/builtin/jaira-cockpit/tests/board.test.ts '.jaira/pr-when says signoff' (no PR -> pr lane, stays; not settled -> pane; settled/merged -> done; pending ticket; forge unchecked)
- [x] Eingabeklassen fuer 'jaira pr': kein PR, kein gh/glab, nicht angemeldet, Fork, Branch mit Sonderzeichen, mehrere PRs, leere/kaputte Antwort, Fehler auf halbem Weg - jede mit eigenem Ausgang (Test)
  proof: core/pr/pr_test.go TestReadInputClasses, TestReadGitLab, TestReadPagesAndOrigin (pages, fork namesake, empty/broken page), TestClassify (401 only as sign-in); internal/cli/pr_test.go TestPRExitCodes
- [x] Kein Doppeltmelden: derselbe Stand meldet nichts erneut, auch nach Neustart der Session; ein neuer PR ist Ausgangsstand, keine Meldung (Test)
  proof: watch: first read is the baseline (internal/cli/pr.go watchPR 'watching' line), Changes compares with the previous read only (core/pr/pr_test.go TestChanges 'same activity read again'); cockpit: store key pr-followed:<project>:<n> (board.test.ts 'handed over before ... also by a later session', cockpit.test.ts 'once')

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/pr: Forge aus jaira.forge/origin, Ziel-Repo (Fork->Parent), offene PRs ueber gh api graphql bzw. glab api lesen und in ein Modell giessen (Threads offen/beantwortet/erledigt, fremde Kommentare + letzter, CI, Reviews ausstehend/Aenderung verlangt, Tickets, settled)
- [x] internal/cli: 'jaira pr [<branch>] [--all] [--json] [--watch [--every d]]' - Text und JSON; --watch: Ausgangsstand ohne Meldung, dann eine Zeile je Aenderung (fremder Kommentar/Review, CI rot/gruen, settled), Ende bei merged/closed; Exit 3 bei fehlendem gh/glab oder Login; Release-Note
- [x] Tests Go: Tabellentest ueber Fake-gh/glab (kein PR, mehrere, Fork, CI-Zustaende, leere Antwort, kaputtes JSON, Fehler halb im Lauf)
- [x] Cockpit: kein Dauer-Polling. Nach Accept, das die pr-Lane startete: begrenzte Nachschau (je 60 s, hoechstens 20 Mal) per 'jaira pr <branch> --json', sobald der PR da ist einmal 'PR #N fuer <ID> offen - begleiten: jaira pr --watch <branch>' per prompt.submit; ein Agentenbericht mit PR-URL bekommt dieselbe Zeile; einmal je PR (Store)
- [x] Cockpit-Accept bei pr-when=signoff: kein PR -> pr-Lane starten, Ticket bleibt in signoff; PR offen und nicht settled -> Liste, bleibt; settled -> done
- [x] Gates lokal: go vet, go test (cli/pr), claude plugin test, shipped-manifest falls Rollen beruehrt

## Progress
- **2026-11-20 16:49 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-20 17:58 · Alexander Sacharov** — von JG3REW (teamlead jaira-64, 2026-11-20): zwei Punkte gehoeren hierher. (1) Auf Boards mit .jaira/pr-when=signoff darf ein Ticket erst nach done, wenn der PR keine offenen Threads und kein ausstehendes Review hat. JG3REW prueft das nur in der Rolle jaira-role-pr (scripts/pr-round.sh settled <owner/repo> <number>, Exit 3 mit Liste). Ein Gate in 'jaira move ... done' im Binary fehlt bewusst: jeder move haengte sonst an Netz und gh-Login. (2) Der Cockpit-Accept (core/mod/builtin/jaira-cockpit/hooks/board.ts stepAfterAccept) oeffnet den PR erst, wenn alle Tickets in done sind. Bei pr-when=signoff muss er den PR schon in signoff oeffnen und vor done 'settled' abwarten.
- **2026-11-22 14:48 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 14:49 · Alexander Sacharov** — Offene Entscheidungen vor dem Plan (worker CRASZS):
F1 Wo lebt der PR-Stand?
 1. (Empfohlen) Befehl 'jaira pr [--json]' im Binary (gh/glab); Cockpit liest ihn, QFWRXH (PR-Spalte im TUI) kann ihn wiederverwenden - eine Quelle, mehr Go-Code.
 2. Nur im Cockpit (TS ruft gh) - weniger Code, aber ohne Cockpit sieht niemand etwas.
 3. Neues Kommando in pr-round.sh (status) - am billigsten, nur GitHub, schwer zu finden.
F2 Wie erfaehrt der Teamlead eine Aenderung?
 1. (Empfohlen) Cockpit fragt alle 5 min in der Session der Person ab und schickt je Aenderung eine Zeile per prompt.submit in die Session - der Teamlead reagiert sofort; nur solange das Cockpit laeuft.
 2. Nur Toast/Band-Zeile - die Person sieht es, das Modell reagiert nicht.
 3. 'jaira pr --watch' fuer Monitor/loop - ohne Cockpit, muss aber jede Session neu gestartet werden.
F3 Die zwei Punkte aus JG3REW (Notiz 20.11.: Cockpit-Accept oeffnet den PR bei pr-when=signoff schon in signoff und wartet vor done auf 'settled')
 1. (Empfohlen) In diesem Ticket, als DoD-Punkte - gehoert zum Begleiten eines PR.
 2. Eigenes Ticket mit follows CRASZS - dieses bleibt kleiner.
Selbst entschieden (widersprechen moeglich): Ausloeser = neuer Kommentar/Thread, CI wird rot, Review abgegeben; Zuordnung PR->Ticket ueber Branchname und Commit-Subjects wie im Cockpit heute.
- **2026-11-22 16:15 · Alexander Sacharov** — decision (Alex 22.11.): F1=1 Befehl 'jaira pr [--json]' im Binary (gh/glab), Cockpit liest ihn; F2=1 Cockpit fragt alle 5 min ab und schreibt je Änderung eine Zeile per prompt.submit in die Session; F3=1 die zwei JG3REW-Punkte als DoD in diesem Ticket.
- **2026-11-22 16:18 · Alexander Sacharov** — decision (Alex 22.11., ersetzt F2=1): kein 5-min-Polling im Cockpit. Signal nur wenn nötig: Cockpit meldet 'PR offen – begleiten' bzw. 'angenommen' einmal in die Teamlead-Session; das Beobachten startet der Teamlead selbst (Monitor über 'jaira pr', endet bei merged/closed). DoD 4 umformuliert.
- **2026-11-22 16:19 · Alexander Sacharov** — Plan nach F2-Wechsel: Benachrichtigung des Cockpits nach Accept braucht eine begrenzte Nachschau, weil die pr-Lane (spawn.sh pr) ohne --parent laeuft - ihr Bericht erreicht die Teamlead-Session nicht. Nicht gewaehlt: spawn.sh/jaira-role-pr um --parent erweitern (Rollentext + shipped-Manifest, groesser).
- **2026-11-22 16:24 · Alexander Sacharov** — jaira pr steht (696778fc8). Gelernt: GitHub-GraphQL pullRequests nimmt headRefName als Filter, also ein Aufruf fuer Branch und --all; GitLab braucht je MR vier Aufrufe (MR fuer head_pipeline - die Liste hat es nicht -, discussions, reviewers, commits). 'offen/beantwortet' = ungeloester Thread, letzter Kommentar fremd vs. eigener (viewer). Watch meldet nur gegen den vorigen Lauf; Ausgangsstand = erste Zeile 'watching', keine Meldung -> DoD 7 fuer den Watch. Rest: Cockpit (Plan 4, 5).
- **2026-11-22 16:34 · Alexander Sacharov** — critique 1: maxArgs (internal/cli/pr.go) gehoert zu exactArgs/minArgs nach internal/cli/root.go. Stehen gelassen, mit Grund: (a) core/pr settle wiederholt die Regel von pr-round.sh settled (GitHub-only Rollenskript); die Rolle auf 'jaira pr --number N --json' umzustellen waere ein Eingriff in jaira-role-pr + Manifest, nicht Teil dieses Tickets. (b) Forge-/Fork-Erkennung gibt es jetzt in Go (core/pr) und TS (board.ts forgeOf/targetRepoOf) - zwei Sprachen, nicht zusammenlegbar. (c) followPr-Zeile auf Englisch statt 'PR #N fuer <ID> offen - begleiten': texts.ts schreibt jeden Cockpit-Text englisch vor. (d) --all bleibt: DoD 1 verlangt den Stand je Ticket mit PR, ueber alle PRs.
- **2026-11-22 16:34 · Alexander Sacharov** — critique 2: maxArgs liegt jetzt in internal/cli/root.go neben exactArgs/minArgs (5239634bd) - erledigt, nichts weiter.
- **2026-11-22 16:40 · Alexander Sacharov** — bughunt: change: blocker: core/pr/gitlab.go:90 - class: only the first page of a paginated forge list is read, so the newest items are missing and Settled/LastForeign are computed from a truncated list; also gitlab.go:100 (commits, tickets), github.go:21 reviewThreads(first:100) (first = the OLDEST 100). example: a GitLab MR with more than 100 discussions (GitLab counts every plain comment and every system note such as 'added 1 commit' as its own discussion, oldest first) - an unresolved thread on page 2 is not counted, settled=true, heldInSignoff lets Accept move the ticket to done, and --watch never reports the new comments; read (GitLab discussions API pages ascending, per_page max 100, no page loop here)
- **2026-11-22 16:40 · Alexander Sacharov** — bughunt: change: blocker: core/pr/watch.go:62 - class: a CI transition is only reported against the immediately previous read, so any intermediate state hides it; example: red -> fix pushed -> pending -> success (the normal recovery, a push always passes through pending at a 1-minute interval) prints no 'CI green again', and red -> pending -> red after a rerun prints no second 'CI red' either; the documented 'CI ... green again' line effectively never fires; ran: scratch test in core/pr calling Changes(red,pending) and Changes(pending,green) - both returned []
- **2026-11-22 16:40 · Alexander Sacharov** — bughunt: change: core/pr/pr.go:166 - class: auth detection by substring '401' anywhere in the tool's error text, so any number containing 401 is read as a lost sign-in; example: 'jaira pr --number 401' on a missing PR (gh exits 1: 'Could not resolve to a PullRequest with the number of 401') exits 3 not_signed_in instead of 1, a project path like group/api-4012 in a glab 404 does the same, and inside --watch such a transient error ENDS the watch instead of retrying; ran: scratch test classify(gh/glab, ...) - errors.Is(err, ErrAuth) true for both
- **2026-11-22 16:40 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:655 - class: a non-idempotent step (spawn the pr lane) behind a button the person is told to press again, with no guard for one already running; example: Accept on a pr-when=signoff board starts spawn.sh pr, the ticket stays in signoff, and a second Accept before the lane has created the PR (has===false still) spawns a second pr agent in the same worktree - two agents pushing and answering the same review comments; spawn.sh (core/role/builtin/jaira-dispatcher/scripts/spawn.sh) has no check for an existing tab of the ticket/lane; read
- **2026-11-22 16:40 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:640 - class: the signoff hold fails open where its own lookups fail, although its comment says a forge that cannot be asked holds the ticket; also :643 (jaira lanes failing -> arrived=[] -> not held); example: 'jaira roles list --json' timing out (30s) or exiting non-zero on a pr-when=signoff board makes heldInSignoff return false and Accept moves the ticket to done with no PR or an unsettled round; read
- **2026-11-22 16:40 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:849 - class: a transient failure treated as permanent; lookForPullRequest stops looking on any non-zero exit, but only exit 3 (no forge/tool/sign-in) is permanent - exit 1 (forge_failed: network blip, 5xx, rate limit) also ends the 20 looks, so the teamlead is never handed the new PR; read
- **2026-11-22 16:40 · Alexander Sacharov** — bughunt: change: core/pr/github.go:26 and register.tsx:248 - class: a pull request matched by something that does not identify it in this repository; (a) GitHub's headRefName filter ignores the head repository, so on the parent of a fork 'jaira pr' on branch master (or any name another contributor's fork also uses) lists other people's PRs and --watch refuses with 'several pull requests'; (b) pullRequestIn takes the first /pull/<n> URL in an agent report whatever repository it points at, so a report quoting another project's PR hands the teamlead 'jaira pr --watch --number <n>' for an unrelated number here and burns that number's once-only store key; read
- **2026-11-22 16:40 · Alexander Sacharov** — bughunt: done at 5fc0b0ba5 - 7 findings (7 change, 2 blocker, 0 pre-existing) in 4 modules (core/pr, internal/cli pr, cockpit hooks, cockpit texts), first hunt; concurrency, error paths, input, permissions checked; files: README.md core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/tests/cockpit.test.ts core/pr/github.go core/pr/gitlab.go core/pr/pr.go core/pr/pr_test.go core/pr/watch.go core/release/notes/CRASZS.md core/shipped/shipped.sha256 docs/COMMANDS.md internal/cli/pr.go internal/cli/pr_test.go internal/cli/root.go
- **2026-11-22 16:41 · Alexander Sacharov** — testing 1 -> in-progress: bughunt 7x change. Fix-Plan je Klasse: (1) Fenster -> paginieren (GitHub PRs per --paginate, Threads je PR nachladen; GitLab glab api --paginate als JSON-Strom), wo nicht nachladbar fail-closed als unsettled; (2) CI-Vergleich gegen das letzte endgueltige Ergebnis + head-SHA statt gegen den letzten Lauf; (3) Login-Erkennung nur an gh-Exit 4 / 'HTTP 401' / 'Bad credentials' / 'auth login', nicht an jeder 401; (4) zweites Accept startet keine zweite pr-Lane (Store-Marker 20 min); (5) roles/lanes-Fehler halten in signoff fest (fail-closed); (6) Nachschau bricht nur bei Exit 2/3 ab, nicht bei Exit 1; (7) PR nur aus dem eigenen Head-Repo (headRepository bzw. source_project_id) und PR-URL im Bericht nur, wenn ihr Repo in git remote -v steht.
- **2026-11-22 16:49 · Alexander Sacharov** — Runde 2 (bughunt-Fixes): GitHub-Liste selbst paginiert statt gh --paginate - gh sucht das erste pageInfo in der Antwort, und die Abfrage hat verschachtelte pageInfo (Threads), das waere die falsche Seite. GitLab dagegen per glab api --paginate (REST-Link-Header, keine Verschachtelung), Seiten als JSON-Strom dekodiert (core/pr/pr.go values). 'more than 100 reviewers' nicht nachgeladen, sondern fail-closed als unsettled. CI-Vergleich: Carry() traegt das letzte endgueltige Ergebnis + head-SHA ueber pending hinweg.
- **2026-11-22 16:49 · Alexander Sacharov** — critique 3: die letzte Kritik-Runde endete mit none; diese Runde beantwortete den bughunt, nicht die Kritik. Die Reparatur (Seiten, Herkunft, CI ueber pending, fail-closed) aendert die Form nicht - nichts neu.
