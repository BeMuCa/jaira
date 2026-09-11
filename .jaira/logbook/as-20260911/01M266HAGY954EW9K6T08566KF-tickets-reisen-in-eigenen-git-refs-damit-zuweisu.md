---
id: 01M266HAGY954EW9K6T08566KF
title: "Tickets reisen in eigenen Git-Refs, damit Zuweisungen ohne gemeinsamen Branch ankommen"
status: done
ready: true
creator: Alexander Sacharov
goal: "Ein Ticket und sein Besitzer erreichen den Kollegen ueber refs/jaira/tickets/<id>, ohne dass ein Branch geteilt oder gemergt werden muss, und eine Zuweisung loest bei ihm eine Benachrichtigung aus"
context: "Berk will laut Issue #7 benachrichtigt werden, wenn ihm jemand ein Ticket zuweist. Heute geht das nicht: die Ticketdatei liegt im Branch des Schreibers. Ist der Branch ungepusht, veraltet oder force-gepusht, sieht der Empfaenger nichts - weder Titel noch Id. Alle Branches abzuscannen ist teuer und trotzdem falsch. Geprueft und belegt: ein Push auf refs/jaira/... nimmt GitHub an (exit 0), ein zweiter Push aufs selbe Ref wird als non-fast-forward abgelehnt (exit 1) - das ist ein Compare-and-Swap ohne Server. Der Tree des Refs kann die Ticketdatei selbst tragen, ein Kollege liest sie per git show ohne Branch und ohne Checkout. Der Default-Refspec zieht diese Refs nicht mit, wer jaira nicht nutzt merkt nichts. Ausgeschlossen: getrennte claims-/log-Refs (ein Ref pro Ticket reicht, --force-with-lease auf den gelesenen Sha ist das CAS) und ein Release per Commit mit Parent (laesst das Ref bestehen, das Ticket klebt dann fuer immer am ersten Besitzer). Die Ticketdatei bleibt zusaetzlich im Code-Commit, damit der Reviewer Diff und Grund an einer Stelle sieht. Details, Schritte und alle Messergebnisse stehen im Body."
definition-of-done: "core/gitref liest, schreibt und loescht refs/jaira/tickets/<id> mit CAS und unterscheidet 'Rennen verloren' von 'Netz weg'; create/claim/move/note/dod pushen das Ticket zusaetzlich ins Ref und melden bei Ablehnung wer schneller war; ein Fetch holt die Refs und das Board zeigt Ref-Tickets samt Besitzer, auch ohne den Branch des Schreibers; eine Zuweisung an mich erzeugt eine Desktop-Benachrichtigung, abschaltbar; ein optionaler Hook bei move/claim wird aufgerufen; ein ins Logbuch gelegtes Ticket loescht sein Ref; Tests mit zwei Klonen eines Bare-Repos belegen Rennen, Uebernahme und den Merge zweier Branches am selben Ticket"
tags:
  - concurrency
  - cli
blocked-by: []
commits:
  - 03961b0ca45f8d14ac927537376b562381590dee
  - d3bc99ef64908478ecd5764b1e0b48003b5d9051
  - d66fcad478207589fea9ab085aa86f4355e98b7c
  - 98009909c4a329efecdda38ce1126a37f40690f9
  - 9f6e629553243faafeb128b5ea138d40da7ce4bf
  - 68aae28834ff19a76ff59b0010771567258ae872
  - 7cd9cd046d6488f72821e060081bebf3a0982197
  - cc94119280855fc2230411228dcbdc7aaae894d6
  - 3dbf02839cb83e92d9cb8f0e38eb225b454782ea
  - 3cee82285f08edadb47aab486c45c6d3f25e7236
  - 7811cb3c37c239f596c20e73dc34a0ea2088de55
  - ec8c7e2b0e293869e0b84b3019a6b4d0fd4bc0aa
created-at: 2026-09-10T17:41:04Z
updated-at: 2026-09-11T13:31:35Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-133823
claimed-at: 2026-09-11T11:19:07Z
question: "Zwei Fragen, bevor das weiterlaeuft: (1) Reicht dir, dass Tickets, die nur auf einem Ref liegen, als Zaehler in der Hinweiszeile stehen und 'jaira fetch' sie vollstaendig listet - oder sollen sie als eigene read-only Karten auf dem Board erscheinen? Letzteres ist ein Folgeticket, weil es einen Durchgangspfad in der Eingabebehandlung braucht. (2) Der Push geht erst nach dem Kommando raus und nur in einem Prozess, der selbst geschrieben hat. Ein reines Lesekommando fetcht damit nie von sich aus - willst du das so, oder soll 'jaira list' gelegentlich auch nachziehen?"
outcome-what: "second-model review of the git-refs feature (core/gitref, outbox, refsync, notify, hook, settings, fetch, board markers, mergebranches_test)"
outcome-why: "diff verified against build + full test suite + source, all DoD lines confirmed; two non-blocking notes and two open implementer questions left for the person signing off"
outcome-resolves: "review-summary, review-gaps, review-verdict, review-check written; approve with notes"
review-summary: "Tickets now travel on refs/jaira/tickets/<id> instead of only living in a branch. core/gitref reads/writes/deletes that ref with compare-and-swap (--force-with-lease), separating a lost race (ErrRaceLost) from an unreachable network (ErrOffline). Every write goes through an outbox (core/outbox) so claim/create/move/note/dod stay offline-capable and flush opportunistically through the two real write gates (Store.Mutate/Create) via a WriteRecorder interface (core/refsync.Syncer). 'jaira fetch' pulls the refs and the board shows ref-only/unsent markers merged field-by-field with the local file (core/merge). core/notify sends a desktop notification (notify-send/osascript/PowerShell) only when a ticket is both mine and changed since last seen, gated by a new ~/.jaira/settings.json (remote, notify-off, hook). core/hook shells out to an optional script on move/claim with env-var contract, silent on failure. Logbook/archive intentionally leave the ref standing until the ticket lands in a configured landing branch, at which point core/snapshot reaps it - avoiding a window where a filed ticket is invisible to everyone. internal/cli/mergebranches_test.go proves two branches editing the same ticket merge field-aware via a real git merge driver."
review-gaps: "One thing worth a look, not blocking: the DoD line 'a logbooked ticket deletes its ref' is satisfied indirectly, not immediately - RecordFiled keeps the ref alive until core/snapshot reaps it after the ticket lands in a landing branch (refsync.go:150-163, snapshot.go reap/Drop). That's a deliberate, documented tradeoff from an earlier ticket (avoids a window where a filed ticket is invisible to everyone) but reads differently from the DoD's literal wording; worth a one-line note in the DoD or README so a future reader doesn't file it as a bug. Process note, not a code gap: the ticket's own 'commits' field/diff handed to this review lane covers only ec8c7e2 (the delta since the last move, out of human), not the full feature - the bulk of the implementation (core/gitref, outbox, refsync, notify, hook, settings, fetch, mergebranches_test) landed in earlier commits already covered by this ticket's critique/optimize/testing loop passes. I additionally re-verified by building, running the full test suite for every listed package, and checking each DoD line against source rather than relying on the handed diff alone. Two open questions from the implementer (read-only ref-only cards on the board; whether 'jaira list' should also opportunistically fetch) are still unanswered in the ticket's question field and belong to whoever signs off next, not to this review. Otherwise: build and full test suite green, every DoD line checked against source/tests."
review-verdict: "Matches the DoD: CAS with typed race/offline errors, configurable board remote, all five write commands recorded through the two store gates, fetch + board display with ref-only/unsent markers, opt-out desktop notification, optional hook, tests with two clones of a bare repo covering race/takeover/merge, and README covering the feature and the fork boundary. Approve, with the two gaps above noted for the record - the ref-deletion timing deserves a documentation tweak, not a rework."
review-check: |-
  1. cd into the repo and run: go build ./... - should exit 0 with no output.
  2. Run: go test ./core/gitref/... ./core/outbox/... ./core/refsync/... ./core/notify/... ./core/hook/... ./core/settings/... ./internal/cli/... -run 'TestTwoBranches|TestAWriteOnTheStore|TestARejectedWrite|TestSecondWriterLosesTheRace' -v - every listed test should print PASS.
  3. Make a throwaway bare repo: mkdir /tmp/board.git && git init --bare /tmp/board.git, then clone it twice into /tmp/a and /tmp/b.
  4. In /tmp/a: jaira init, jaira create 'ping' --goal g --dod d --context c, then check git ls-remote /tmp/board.git 'refs/jaira/*' shows one ref.
  5. In /tmp/b: jaira fetch - should list the ticket with its id and lane, without /tmp/b ever having the branch from /tmp/a (git branch -a in /tmp/b stays empty).
  6. In /tmp/b, claim the ticket, then in /tmp/a try to move it - jaira should refuse and print who claimed it, read from the ref, not silently overwrite.
---

## Warum Refs und nicht Branches

Ein Alert braucht einen Kanal. jaira hat genau einen: das Git-Remote. Sitzt jeder
in seinem eigenen Branch, hilft das nicht — den Branch des anderen sieht man
nicht, er kann ungepusht, veraltet oder force-gepusht sein. Alle Branches
abzuscannen ist teuer und trotzdem falsch.

Ein Ref ausserhalb von `refs/heads` loest das: es gehoert zu keinem Branch, der
Default-Refspec zieht es nicht mit (wer jaira nicht nutzt, merkt nichts), und
der Push darauf ist ein Compare-and-Swap, das die Ref-Sperre des Hosters nutzt.

## Das Ref-Layout

    refs/jaira/tickets/<id>   # ein Ref pro Ticket, Tree enthaelt die Ticketdatei

Ein Ref pro Ticket, sein Tree enthaelt `<id>.md` — nicht nur die Id, sondern den
vollen Inhalt. Damit liest ein Kollege das Ticket ohne Branch und ohne Checkout:

    git show refs/jaira/tickets/<id>:<id>.md

Lane und Besitzer stehen im Frontmatter, die Historie ist `git log` auf dem Ref.
Getrennte claims-/log-Refs sind dadurch ueberfluessig.

## Compare-and-Swap

Jeder Schreibvorgang pusht mit `--force-with-lease=refs/jaira/tickets/<id>:<sha
den man gelesen hat>`. Hat jemand dazwischen geschrieben, lehnt Git ab (exit 1),
und jaira sagt: neu einlesen, das Ticket hat sich bewegt. Ein Rennen wird sofort
erkannt, nicht drei Tage spaeter beim Merge.

## Verhaeltnis zur Ticketdatei im Commit

Beides bleibt, mit klarer Zustaendigkeit:

| | Datei im Commit | Ref |
|---|---|---|
| Antwortet auf | was war das fuer eine Aenderung, warum | wo steht das Ticket jetzt, bei wem |
| Liest | Reviewer des Diffs | Board, Alerts, andere Maschine |
| Aendert sich | einmal, mit dem Code | bei jedem move/claim/note |
| Lebt | dauerhaft in der Master-Historie | solange das Ticket auf dem Board ist |

Der Ablauf aus CLAUDE.md bleibt unveraendert: Ticket bewegen, Ticketdatei
zusammen mit dem Code committen. Zusaetzlich pusht dieselbe Operation das Ref
mit denselben Bytes. Keine zweite Wahrheit, derselbe Blob auf einem Kanal, der
ohne den Branch sichtbar ist.

Solange der Branch nicht gemergt ist, ist das Ref massgeblich fuer Lane und
Besitzer. Die Datei im Commit ist eine Momentaufnahme: auf diesem Commit stand
das Ticket hier. Genau das braucht der Reviewer.

## Zwei Branches am selben Ticket

Auf Ref-Ebene entsteht der Konflikt gar nicht erst — CAS lehnt den Zweiten ab.
Auf Dateiebene loest der bestehende Merge-Driver (`core/merge`): `status` nach
Fortschritt in der Lane-Kette, Listen (`commits`, `tags`, `blocked-by`) als
Union, Prosafelder als einziger echter Konflikt. Gegen die reale Registrierung
geprueft, siehe unten.

## Offen: Claim ohne Netz

Ohne Verbindung gibt es kein CAS. Vorschlag: `claim` verlangt Netz, alles andere
(`move`, `note`, `dod`) bleibt vollstaendig offline. Eine einzige
Online-Operation ist leicht zu erklaeren und unmoeglich misszuverstehen. Die
Alternative — lokal provisorisch claimen und beim ersten Push nach Timestamp
abgleichen — ist gebaut komplizierter und im Konfliktfall unangenehm.

## Was schon geprueft ist (gegen echtes GitHub, nicht theoretisch)

- Push eines neuen `refs/jaira/...` auf GitHub: angenommen, exit 0.
- Zweiter Push ohne Parent auf dasselbe Ref: abgelehnt, non-fast-forward, exit 1.
- `--force-with-lease` auf den gelesenen Sha: der Zweite wird abgelehnt, der
  Erste gewinnt, der Serverzustand bleibt konsistent.
- `git ls-remote origin 'refs/jaira/*'`: alle Refs in einem Roundtrip sichtbar.
- Loeschen per `:refs/jaira/...`: funktioniert.
- Branches, PRs und die GitHub-Weboberflaeche bleiben unberuehrt — die Refs sind
  dort nirgends sichtbar. Deshalb muss `jaira board` die Besitzer selbst zeigen.
- Ticketdatei im Ref-Tree, ohne jeden Branch beim Empfaenger lesbar.
- Merge-Driver: `in-progress` vs `review` ergibt `review`, obwohl `in-progress`
  den neueren Timestamp trug; `done` vs `backlog` ergibt `done`; `tags` und
  `commits` beider Seiten bleiben erhalten.

## Umsetzungsschritte

1. `core/gitref`: Refs lesen/schreiben/loeschen ueber `os/exec`, CAS mit
   `--force-with-lease`, Retry bei Netzfehler, klarer Fehlertyp fuer "Rennen
   verloren" (unterscheidbar von "Netz weg").
2. Schreibpfad: `create`, `claim`, `move`, `note`, `dod` pushen dasselbe Ticket
   zusaetzlich ins Ref. Bei Ablehnung: neu einlesen und dem Aufrufer sagen, wer
   schneller war.
3. Lesepfad: `jaira fetch` (bzw. Hintergrund-Fetch im TUI) holt
   `+refs/jaira/tickets/*:refs/jaira/tickets/*` und mischt Ref-Tickets in die
   Boardansicht — mit Kennzeichnung, dass sie noch in keinem Branch liegen.
4. Alert: Aenderung an einem Ticket mit einem selbst als Assignee loest eine
   Desktop-Benachrichtigung aus (`notify-send` / `osascript`), Intervall
   30–60 s, per Konfiguration abschaltbar.
5. Hook-Punkt fuer sofortige Zustellung: ein optionales Skript, das bei
   `move`/`claim` aufgerufen wird. jaira ruft nur auf und bringt keine
   Abhaengigkeit mit (Slack-Webhook, ntfy.sh, Telegram sind Sache des Nutzers).
6. Aufraeumen: gemergte bzw. ins Logbuch gelegte Tickets loeschen ihr Ref.
7. Tests: zwei Klone eines Bare-Repos in `t.TempDir()`, echtes Rennen um ein
   Ticket, Release und Uebernahme, Merge zweier Branches am selben Ticket.

## Was ausdruecklich nicht dazugehoert

Kein Server, kein Daemon, keine zentrale Koordination. Alles laeuft ueber das
Remote, das ohnehin da ist. Wer jaira nicht benutzt, sieht von alldem nichts.

## Options

- [x] planning

## Plan

- [x] core/gitref: Ref-Inhalt ohne Checkout bauen (hash-object -w, mktree, commit-tree mit dem geleasten Sha als Parent) und lesen (git show <sha>:<id>.md)
  proof: core/gitref/gitref.go:Write/Read, TestTicketArrivesWithoutASharedBranch
- [x] core/gitref: Push mit --force-with-lease=<ref>:<gelesener sha>, Fehler typisieren: ErrRaceLost (non-fast-forward/stale info) vs ErrOffline (Netz/Remote weg)
  proof: core/gitref/gitref.go:classify + push, TestSecondWriterLosesTheRace / TestUnreachableRemoteIsOfflineNotARace
- [x] Entschieden: asynchron. Mutate/Create legen den Schreibvorgang in eine Outbox, gepusht wird ausserhalb - Voraussetzung fuer offline claim, loest zugleich TUI und Merge-Driver
  proof: core/refsync/refsync.go - Record/RecordDelete legen nur ab, Flush sendet ausserhalb des Schreibpfads
- [x] Schreibpfad anbinden: der Push haengt an Store.Mutate + Store.Create, nicht an create/claim/move/note/dod einzeln (alle 20+ Aufrufstellen laufen durch diese zwei)
  proof: core/ticket/store.go: Mutate/Create rufen record(), Archive/Delete/Logbook rufen recordDelete(); internal/cli/root.go openStore + internal/tui/model.go haengen den Recorder an
- [x] Bei Ablehnung: Ticket aus dem Ref neu einlesen und dem Aufrufer melden, wer schneller war (Assignee + updated-at aus dem fremden Ref)
  proof: core/refsync/refsync.go:winner + Winner.Describe, internal/cli/refs.go flushRefs; TestARejectedWriteNamesWhoWasQuicker
- [x] Lesepfad: fetch von +refs/jaira/tickets/*:refs/jaira/tickets/* plus git ls-remote als billiger Ueberblick
  proof: core/gitref Fetch(--prune)/ListRemote + core/refsync Incoming(); internal/cli/fetch.go 'jaira fetch'
- [x] Board: eine Geschichte pro Ticket - lokale Datei als Basis, Ref via core/merge.Merge dazu (base = Blob des Parent-Commits), plus die Marker ref-only und unsent an der Karte
  proof: core/refsync Reconcile (merge.Merge, base = Blob des Parent-Commits) + TUI-Marker: unsent an der Karte (internal/tui/view.go), ref-only als Zaehler in der Hinweiszeile; Tests TestReconcile*
- [x] Settings-Datei in ~/.jaira/ anlegen (existiert heute nicht, nur projects.json und state/) - traegt den Abschalter fuer die Benachrichtigung
  proof: core/settings/settings.go (~/.jaira/settings.json), 4 Tests
- [x] Desktop-Benachrichtigung bei Zuweisung an mich: notify-send / osascript per os/exec, kein neues Modul, faellt still aus wenn nichts da ist
  proof: core/notify/notify.go (notify-send/osascript/PowerShell, WSL-Fallback), announceArrivals in internal/cli/fetch.go
- [x] Optionaler Hook bei move/claim: jaira ruft ein Skript auf und bringt keine Abhaengigkeit mit
  proof: core/hook/hook.go + fireHook in internal/cli/flow.go (move) und claim.go (claim); Smoke-Test: 'move 01M26C7P... todo' im Hook-Log
- [x] logbook und archive loeschen das Ref (git push origin :refs/jaira/tickets/<id>)
  proof: core/ticket/store.go Archive/Delete/Logbook -> recordDelete -> outbox OpDelete; TestRecordDeleteTakesTheRefDown
- [x] Tests: Bare-Repo plus zwei Klone in t.TempDir(), echtes Rennen, Uebernahme, Merge zweier Branches am selben Ticket
  proof: core/gitref/gitref_test.go (Rennen, Uebernahme), core/refsync/refsync_test.go, internal/cli/mergebranches_test.go (echter git merge zweier Branches am selben Ticket)
- [x] Outbox unter ~/.jaira/state/<worktree>/: ein noch nicht gepushter Schreibvorgang pro Ticket, mit dem gelesenen Sha als Lease
  proof: core/outbox/outbox.go (Queue/Pending/List/Drop/Flush), 8 Tests in core/outbox/outbox_test.go, darunter TestAnOfflineWriteArrivesWhenTheNetworkIsBack gegen echtes git
- [x] Outbox abarbeiten beim naechsten Kommando, das Netz hat; abgelehnter nachgeholter Push meldet dem Benutzer, wer schneller war
  proof: internal/cli/root.go Execute ruft flushRefs nach dem Kommando; internal/cli/refs.go meldet rejected/unsent/failed auf stderr
- [x] Board-Remote konfigurierbar (Default origin) - die Settings-Datei traegt es neben dem Abschalter der Benachrichtigung
  proof: core/settings Settings.RemoteName(), von attachRefs an refsync uebergeben
- [x] README: wozu die Refs da sind, welche Befehle sie schreiben und lesen, wie das Board-Remote konfiguriert wird, und die Fork-Grenze ausdruecklich (kein Push-Recht = Ticketdatei im Branch plus PR)
  proof: README.md Abschnitt 'Tickets travel on their own git refs', inkl. Fork-Grenze und settings.json
- [x] TUI: Fetch nur im Hintergrund (tea.Cmd + Program.Send), gezeichnet wird aus lokalen Dateien und bereits gefetchten Refs - nie Netz in View/Update
  proof: internal/tui/refs.go: fetchRefs als tea.Cmd, refTick alle 60s, refMarks liest nur Platte

## Progress

- **2026-09-10 18:19 · Alexander Sacharov** — Warum der Plan so aussieht, drei Funde im Code (gemessen, nicht vermutet):

1. Es gibt genau zwei Schreibtore: core/ticket/store.go:644 Mutate und :620 Create. Alle Aufrufstellen (create, claim, note, dod, tags, move, TUI, Merge-Driver - 20+ per grep) laufen durch sie. Der Ref-Push gehoert deshalb dorthin, nicht in die fuenf CLI-Kommandos einzeln - genau das Muster von JJ32B4 (eine Move-Funktion, vier delegierende Aufrufer).

2. Genau daraus folgt aber ein Problem, das die Umsetzungsschritte im Ticket nicht nennen: Mutate wird auch aus dem TUI (internal/tui/edit.go, external.go, followup.go) und aus dem Merge-Driver (internal/cli/mergedriver.go:221) gerufen. Ein synchroner Push in Mutate haengt das Board am Netz und laesst den Merge-Driver mitten im git-merge pushen. Deshalb ist Schritt 3 eine Entscheidung und keine Implementierung.

3. core/gitrepo ist heute reines Lesen (git.go, 141 Zeilen: Commits, Diff, Stat, HeadSHA) und kennt kein Remote. Ein Schreib-Paket core/gitref daneben ist additiv, nichts Bestehendes muss umgebaut werden.

4. Fuer den Abschalter der Benachrichtigung gibt es noch keinen Ort: in ~/.jaira/ liegen nur projects.json und state/<worktree>/. Die Settings-Datei ist damit Teil dieses Tickets, nicht Voraussetzung.

Nicht geprueft, absichtlich: ob GitHub die Refs annimmt - das steht im Ticketbody als bereits gemessen (exit 0 / exit 1 / force-with-lease / ls-remote / Loeschen), ich habe es nicht wiederholt.
- **2026-09-10 18:44 · Alexander Sacharov** — Geprueft, ob der Plan ueber einen Fork funktioniert. Er tut es nicht.

Test: ein Ref nach origin (mein Fork sashasoft90/jaira) gepusht, dann beide Remotes abgefragt.
  git ls-remote origin   'refs/jaira/*'  -> refs/jaira/probe sichtbar
  git ls-remote upstream 'refs/jaira/*'  -> leer

Forks teilen auf GitHub den Objektspeicher, aber jedes Repository hat eigene Refs. Branches sieht man beim anderen nur, weil ein PR ein eigener Kopiermechanismus ist. Custom Refs fahren da nicht mit.

Folge fuer die Umsetzung: es braucht ein Board-Remote - ein Repository, in das alle Beteiligten die jaira-Refs pushen. Nicht 'jeder in seinen Fork'. Default origin, per Konfiguration ueberschreibbar. Ohne das schreibe ich in meinen Fork und Berk schaut in upstream, und beide sehen nichts.

Auch geprueft: ich habe push auf BeMuCa/jaira (permissions.push = true), master ist nicht geschuetzt. Upstream als Board-Remote ist also ohne neues Repository moeglich.

Nicht loesbar und deshalb festhalten: GitHub kann Rechte nicht auf ein Ref-Namespace einschraenken. push auf refs/jaira/* ohne push auf refs/heads/* gibt es nicht - Schreibrecht gilt immer fuers ganze Repository. Teilnehmer am Board = Person mit Push-Recht auf das Board-Repo. Fuer ein Team in Ordnung, fuer einen fremden Contributor nicht: der bleibt beim heutigen Weg, Ticketdatei im Branch plus PR.

Naechster Schritt: Konfigurationsfeld fuers Board-Remote in Schritt 1 der Umsetzung aufnehmen, bevor core/gitref geschrieben wird.
- **2026-09-10 18:47 · Alexander Sacharov** — Entschieden (Alexander, 10.09.2026): claim bleibt offline moeglich. Offline ist selten genug, dass es die Bedienung nicht diktieren darf.

Der Ablauf: claim schreibt lokal mit Timestamp, der Push wird nachgeholt, sobald Netz da ist, und wird er dann abgelehnt, sagt jaira es dem Benutzer (wer schneller war), statt still zu verlieren.

Damit ist die Alternative aus dem Ticketbody ('claim verlangt Netz') verworfen - der Vorschlag dort ist nicht mehr gueltig.

Folge fuer den Plan: Schritt 3 ist damit entschieden, und zwar auf asynchron. Ein synchroner Push in Mutate kann es gar nicht sein, denn ein lokal moeglicher claim braucht ohnehin einen Ort, an dem ein noch nicht gepushter Schreibvorgang liegt (Outbox unter ~/.jaira/state/<worktree>/). Das loest gleichzeitig das TUI- und Merge-Driver-Problem aus Note 1: Mutate legt nur ab, gepusht wird ausserhalb.

Damit gehoert zum Ticket zusaetzlich: die Outbox selbst, ihr Abarbeiten beim naechsten Kommando mit Netz, und die Meldung bei Ablehnung eines nachgeholten Pushs.
- **2026-09-10 18:51 · Alexander Sacharov** — Abweichung vom Plan, Schritt 1: der Ref-Commit bekommt den geleasten Sha als Parent, nicht 'ohne Parent' wie im Ticketbody vorgesehen.

Grund: mit Parent ist ein veralteter Schreibvorgang von sich aus ein non-fast-forward, den git ohnehin ablehnt - das CAS haengt dann nicht allein an --force-with-lease. Ausserdem wird 'git log refs/jaira/tickets/<id>' damit zur Bewegungshistorie des Tickets, was der Ticketbody selbst voraussetzt ('die Historie ist git log auf dem Ref') - ohne Parent gaebe es genau einen Commit und keine Historie.

--force-with-lease bleibt trotzdem gesetzt: beim allerersten Schreiben hat kein Commit einen Parent, mit dem git vergleichen koennte, und nur der leere Lease sagt 'ich erwarte, dass es das Ref noch nicht gibt'. Belegt in TestFirstWriteRefusesAnExistingRef.

Zwei Dinge, die beim Bauen aufgefallen sind und nicht im Plan standen:

- fetch braucht --prune. Ohne das bleibt ein geloeschtes Ref (Logbuch, Archiv) auf jedem anderen Klon fuer immer stehen: fuer den Besitzer weg, fuer alle anderen unsterblich. Steht jetzt als Vertrag im Doc-Kommentar von Fetch, nicht als Aufraeumen.
- Der Push setzt GIT_TERMINAL_PROMPT=0 und gitref schreibt Author/Committer selbst in die Env. Sonst haengt ein Push, der im Hintergrund eines fremden Kommandos laeuft, an einem Credential-Prompt, und auf einer Maschine ohne user.name (CI, frischer Container) scheitert commit-tree mit 'please tell me who you are'.

Lokales Ref wird erst nach dem akzeptierten Push bewegt (update-ref danach), damit der naechste Lease immer ein Sha ist, das das Remote wirklich gesehen hat. Belegt am Ende von TestUnreachableRemoteIsOfflineNotARace.
- **2026-09-10 19:00 · Alexander Sacharov** — Entschieden (Alexander + ich, 10.09.2026): das Board zeigt eine Geschichte pro Ticket, nicht zwei. Zwei Quellen ja, aber sie werden unter der Oberflaeche zusammengefuehrt.

Warum nicht nur das Ref: es faellt in vier Faellen aus, in denen das Board trotzdem stimmen muss - git ist gar nicht Voraussetzung (jaira laeuft in einem Verzeichnis, das kein Repo ist), offline liegt ein noch nicht gesendeter Schreibvorgang in der Outbox und muss sofort sichtbar sein, handverlesene Aenderungen passieren an der Datei im Arbeitsbaum und nicht am Ref, und auf einem noch nicht geteilten Board (init gitignored .jaira/) gibt es kein einziges Ref. Die lokale Datei ist damit die Basis, das Ref die zweite Quelle.

Warum nicht zwei Spalten im TUI: das laesst den Menschen bei jedem Blick mit den Augen mergen - genau die Arbeit, die core/merge/merge.go:80 schon macht (status nach Fortschritt in der Lane-Kette, Listen als Union, Prosa als einziger echter Konflikt).

Der Fund, der es billig macht: weil der Ref-Commit den geleasten Sha als Parent hat, liefert das Ref die Merge-Basis gratis. Der Blob im Parent-Commit ist genau 'der Stand, von dem ich ausgegangen bin'. Damit ist es ein echter Drei-Wege-Merge und keine Zwei-Wege-Raterei nach Timestamp:

  base   = Blob aus dem Parent-Commit des Refs
  ours   = .jaira/tickets/<id>.md
  theirs = Blob aus dem Ref

Kein neues Vorrangregelwerk. Genau das ist auch der zweite Grund, den Parent zu behalten (siehe die Abweichungs-Note oben).

Was der Mensch stattdessen sieht: zwei Zustandsmarker an der Karte, keine zwei Quellen - 'ref-only' (das Ticket liegt in keinem Branch, jemand hat es mir zugewiesen; der Fall, fuer den dieses Ticket existiert) und 'unsent' (die Outbox haelt eine Schreibung). Ein Prosakonflikt geht den bestehenden Weg: conflict-theirs-<field> plus jaira resolve, kein dritter Mechanismus.

Harte Anforderung ans TUI: Fetch ist Netz und hat in View/Update nichts zu suchen. Hintergrund per tea.Cmd und Program.Send; gezeichnet wird aus den lokalen Dateien und den bereits gefetchten Refs, die danach lokal und damit offline billig lesbar sind.
- **2026-09-10 19:05 · Alexander Sacharov** — Outbox gebaut, drei Regeln, die beim Bauen entschieden wurden und nicht offensichtlich sind:

1. Ein zweiter lokaler Schreibvorgang, waehrend der erste noch liegt, ersetzt den Inhalt, behaelt aber den Lease des ersten. Die Frage, die das Remote beantwortet, ist 'hat seit meinem letzten Lesen jemand geschrieben' - der eigene ungesendete Schreibvorgang ist keine Antwort darauf. Nimmt man den neueren Lease, vergleicht das Remote gegen ein Sha, das es nie gesehen hat, und der Push wird aus dem falschen Grund abgelehnt. Belegt in TestSupersedingAWriteKeepsTheConfirmedLease. QueuedAt bleibt dabei ebenfalls das des ersten - das Ticket wartet seit dann, nicht seit der letzten Aenderung.

2. Der Eintrag traegt die ganze Ticketdatei, keinen Patch. Ein wiedergespielter Patch koennte sauber auf einen Stand passen, den nie jemand geprueft hat. Deshalb ist Ersetzen auch nicht verlustbehaftet: die neueren Bytes enthalten schon alles, was die aelteren sagten.

3. Flush bricht beim ersten unerreichbaren Remote ab, laeuft aber bei einem verlorenen Rennen weiter. Netz ist eine Eigenschaft der Maschine, nicht des Tickets: hat ein Push bewiesen, dass es keine Route gibt, wartet jeder weitere in denselben Timeout, und ein Kommando, das der Benutzer aus einem voellig anderen Grund gestartet hat, haengt so lange wie die Queue ist. Ein verlorenes Rennen betrifft dagegen ein Ticket, das naechste kann durchgehen.

Ein verlorenes Rennen loescht seinen Eintrag: erneut senden wuerde nur erneut verlieren, und die lokale Datei liegt weiter da - das Board fuehrt das Ref per core/merge dazu. Ein unklassifizierter git-Fehler behaelt den Eintrag, denn Wegwerfen bei einem Fehler, den niemand eingeordnet hat, ist der einzige Ausgang, der wirklich Arbeit verliert.

Was die Outbox nicht macht: sie sperrt nicht. Zwei Sessions, die gleichzeitig flushen, pushen beide, eine verliert das Rennen und raeumt ihren Eintrag - das korrigiert sich selbst. Eine Sperre waere hier teurer als das Problem.
- **2026-09-10 19:15 · Alexander Sacharov** — Schreibpfad angebunden, drei Dinge, die beim Bauen entschieden oder gefunden wurden:

1. core/ticket bekommt eine Recorder-Schnittstelle (WriteRecorder mit Record und RecordDelete), core/refsync erfuellt sie. Umgekehrt waere es ein Importzyklus: outbox und refsync lesen Tickets, also darf core/ticket sie nicht kennen. Die Schnittstelle steht damit aus demselben Grund da wie das Feld Actor - der Aufrufer weiss, wer er ist, und core/ticket muss nichts von git wissen.

2. Ein Fehler des Recorders macht die Mutation nicht rueckgaengig. Zu dem Zeitpunkt liegt die Datei schon auf der Platte, also ist es kein Veto. Dafuer gibt es ticket.ErrNotRecorded: der Aufrufer kann melden 'lokal geschrieben, nicht weitergegeben', ohne dass die Mutation als gescheitert gilt.

3. Geflusht wird nur, wenn dieser Prozess selbst etwas eingereiht hat (refsync.Dirty). Ein Lesekommando bleibt damit vollstaendig vom Netz weg - 'jaira list' darf nie auf ein Remote warten. Gesendet wird in cli.Execute nach dem Kommando, auch wenn das Kommando gescheitert ist: ein Ticket, das geschrieben wurde, soll das Team sehen.

Gefunden und behoben, war im Plan nicht vorgesehen: 'kein Repository' und 'Remote nicht erreichbar' sind zwei verschiedene Zustaende, und mein erster Wurf hat beide als offline behandelt. Ein Verzeichnis ohne .git ist ein unterstuetzter Weg, jaira zu benutzen (README: git ist optional) - dort darf nichts eingereiht werden, sonst zeigt die Karte fuer immer 'unsent'. Also gibt es jetzt gitref.ErrNoRepo und Repo.Usable(), einmal pro Prozess gefragt und gecacht, statt es aus einem gescheiterten Push zu erraten. Belegt in TestABoardWithoutARemoteRecordsNothing und TestUsableRejectsAMissingRemote.

Nebenbei aufgefallen: meine ersten Offline-Tests haben Offline mit einem Remote-Pfad simuliert, der kein Repository ist - das ist gerade nicht offline. Sie zeigen jetzt auf 127.0.0.1:1, was sofort refused liefert, ohne echten Timeout im Test.

Smoke-Test mit dem gebauten Binary gegen ein Bare-Repo und zwei Klone: create legt das Ref an, der zweite Klon liest das ganze Ticket per git show ohne einen einzigen Branch (git branch -a ist leer), und nach einem fremden Push meldet 'jaira note' auf stderr: 'was not sent: someone else wrote it first, it is now in review'.
- **2026-09-10 19:21 · Alexander Sacharov** — Lesepfad, Einstellungen, Benachrichtigung und Hook stehen. Was dabei entschieden wurde:

1. 'jaira fetch' ist der Lesepfad: ein Roundtrip auf refs/jaira/tickets/*, kein Branch, kein Checkout, kein Merge. Es markiert pro Karte '@you', 'new' und 'ref-only' - letzteres ist die Eigenschaft, die den Fall benennt, um den es geht: das Ticket liegt in keinem Branch, den du hast.

2. Benachrichtigt wird nur bei Mine UND Changed. Ein Ticket, das seit letzter Woche auf mich zeigt, ist keine Nachricht, und bei jedem Fetch zu poppen erzieht den Benutzer dazu, die Benachrichtigung zu ignorieren. 'Changed' braucht Gedaechtnis: refs-seen.json unter dem State-Verzeichnis haelt pro Ticket das zuletzt gemeldete Ref-Sha. Ein Schreibfehler dort ist absichtlich still - der Preis ist eine doppelte Benachrichtigung, nicht ein gescheiterter Fetch.

3. core/notify shellt aus: notify-send, osascript, PowerShell-Balloon. Keine Bibliothek, weil das cgo oder eine Windows-API-Bindung bedeuten wuerde und dieses Werkzeug ein statisches Binary ohne Laufzeitabhaengigkeit ist. Jeder Fehlschlag ist still: im Container, ueber ssh und auf CI gibt es keinen Notifier, und ein Board, das das jedes Mal meldet, ist schlimmer als eines, das leise nichts tut. Auf dieser Maschine (WSL) fehlt notify-send, powershell.exe ist erreichbar - deshalb der WSL-Zweig.

4. Einstellungen: ~/.jaira/settings.json neben projects.json, mit remote, notify-off und hook. Wichtig ist die Formulierung 'notify-off' und nicht 'notify': der Nullwert - keine Datei, oder eine, die das Feld nicht nennt - laesst die Benachrichtigung an. Die Funktion wurde angefragt, und ein Default, der sie still zurueckhaelt, antwortet auf die Anfrage mit nichts. Eine kaputte Datei faellt auf die Defaults zurueck statt das Board zu verweigern.

5. Der Hook ist der Ausweg fuer 'sag es jetzt': jaira ruft ein Skript, der Vertrag sind Umgebungsvariablen (JAIRA_EVENT, JAIRA_TICKET, JAIRA_TITLE, JAIRA_STATUS, JAIRA_ASSIGNEE, JAIRA_ACTOR, JAIRA_ROOT), keine Argumente - so liest ein Skript nur, was es braucht, und ein spaeter ergaenztes Feld bricht nichts. Timeout 5s, Ausgabe verworfen (ein geschwaetziger Hook darf sich nicht in das stdout mischen, das ein Agent als JSON liest), jeder Fehlschlag still.

Smoke-Test mit dem gebauten Binary, zwei Klone: ada legt ein Ticket an und weist es berk zu, berk laeuft 'jaira fetch' und sieht 'J2QRYV backlog cookie dropped on 302 @you new ref-only' - ohne einen einzigen Branch von ada. Der zweite Fetch meldet changed=false, also genau einmal. Ein legaler Move schreibt 'move <id> todo' ins Hook-Log.

Nebenbei belegt: der bestehende Assignee-Gate hat meinen ersten Move-Versuch abgelehnt ('belongs to berk'), also greifen die alten Gates unveraendert weiter.
- **2026-09-10 19:28 · Alexander Sacharov** — Board und Doku fertig, damit ist alles aus der Definition of Done abgehakt. Zwei Entscheidungen am Board, die ich bewusst so und nicht groesser gemacht habe:

1. Ref-Tickets, die in keinem Branch dieses Klons liegen, erscheinen NICHT als eigene Karten, sondern als Zaehler in der Hinweiszeile ('⇢ N on refs'), und 'jaira fetch' listet sie vollstaendig mit Titel, Lane und Besitzer. Grund: eine Karte, auf der das Board keine Aktion zulassen darf, braucht einen eigenen read-only-Durchgangspfad in der Eingabebehandlung - sonst schreibt der erste Tastendruck auf eine Datei, die es hier nicht gibt. Das ist eine eigene Aenderung, kein Nebeneffekt dieser. Verschwiegen wird nichts: der Zaehler ist da, und der Befehl zeigt alles.

2. Der Hintergrund-Fetch laeuft alle 60s als tea.Cmd, nicht im 2s-Takt des lokalen Rescans. Ein Netz-Roundtrip in dieser Frequenz waere falsch, und niemand bekommt zweimal pro Minute ein Ticket zugewiesen. refMarks (unsent, ref-only) liest ausschliesslich Platte - Outbox-Dateien und lokale Refs - und darf deshalb im reload laufen.

Der Zwei-Branch-Merge ist jetzt gegen echtes git belegt (internal/cli/mergebranches_test.go): der Test baut das Binary, registriert damit den echten Driver, laesst ada taggen und berk nach todo bewegen und merged die Branches. Ergebnis: status: todo (die weitere Lane gewinnt, obwohl adas Schreibvorgang der spaetere war), beide Tags erhalten, keine Konfliktmarker.

Beim Schreiben des Tests aufgefallen und beachtet: 'move ... --to pre-process' scheitert mit exit 3, weil diese Lane die Option 'planning' verlangt. Der Test benutzt deshalb todo. Das ist kein Fehler, sondern der bestehende Options-Gate - er greift unveraendert weiter.

Offen und absichtlich nicht in diesem Ticket: Ref-Tickets als eigene, read-only Karten auf dem Board (siehe 1). Das gehoert in ein Folgeticket.
- **2026-09-11 11:23 · Alexander Sacharov** — Review durch eine zweite Sitzung auf einem anderen Modell (Sonnet, eigener Kontext, ohne Schreibrecht am Code): Vertrag erfuellt, 'approve with notes'. Zwei Anmerkungen, beide bearbeitet:

1. Die DoD-Zeile ueber das Loeschen des Refs beim Logbuch las sich anders als der Code arbeitet - das Loeschen ist seit RA7PFE bis zum Landen aufgeschoben. Zeile umformuliert statt sie stehenzulassen: sonst legt in einem halben Jahr jemand genau das als Fehler an.

2. Prozessanmerkung des Reviewers, die stimmt und festgehalten gehoert: der Lane wurde nur der Diff seit dem letzten Move gereicht (ein Commit), nicht das ganze Feature - die Implementierung liegt in frueheren Commits. Er hat deshalb selbst gebaut, die Testsuite jedes genannten Pakets laufen lassen und jede DoD-Zeile am Quelltext geprueft, statt sich auf den gereichten Diff zu verlassen. Wer den naechsten Review-Lauf an einem lange laufenden Ticket macht, sollte dasselbe tun.

Die beiden offenen Fragen im question-Feld sind inzwischen von Alexander beantwortet: 'jaira list' fetcht nicht selbst (dafuer gibt es PP5SCQ, den Hintergrund-Fetch), und Ref-Tickets sind inzwischen echte Karten auf dem Board.

## Definition of Done

- [x] core/gitref liest, schreibt und loescht refs/jaira/tickets/<id> mit CAS und unterscheidet 'Rennen verloren' von 'Netz weg'
  proof: core/gitref/gitref.go (Write/Read/Delete/SHA, ErrRaceLost vs ErrOffline), 7 Tests in core/gitref/gitref_test.go gegen ein Bare-Repo mit zwei Klonen
- [x] Ein konfigurierbares Board-Remote entscheidet, wohin die Refs gehen; Default origin
  proof: core/settings/settings.go: remote in ~/.jaira/settings.json, Default origin (TestDefaultsWithNoFile, TestBlankRemoteIsStillOrigin)
- [x] create, claim, move, note und dod pushen das Ticket zusaetzlich ins Ref und melden bei Ablehnung, wer schneller war
  proof: internal/cli/root.go openStore haengt den Recorder an, Execute flusht danach; die Meldung bei Ablehnung ist im Smoke-Test gegen echtes git erschienen
- [x] Ein Fetch holt die Refs und das Board zeigt Ref-Tickets samt Besitzer, auch ohne den Branch des Schreibers
  proof: jaira fetch (internal/cli/fetch.go) zeigt Ref-Tickets samt assignee ohne Branch; TUI zaehlt sie in der Hinweiszeile und markiert unsent
- [x] Eine Zuweisung an mich erzeugt eine Desktop-Benachrichtigung, per Konfiguration abschaltbar
  proof: core/notify + announceArrivals: nur Mine und nur Changed, abschaltbar per notify-off; refs-seen.json macht 'einmal statt bei jedem Fetch'
- [x] Ein optionaler Hook bei move und claim wird aufgerufen, ohne dass jaira eine Abhaengigkeit mitbringt
  proof: core/hook/hook.go, 5 Tests; fireHook bei move und claim; Smoke-Test mit echtem Skript
- [x] Ein ins Logbuch gelegtes Ticket gibt sein Ref frei: der Endzustand wird daraufgeschrieben, und entfernt wird es, sobald das Ticket in einer Landebranch angekommen ist (RA7PFE) - sofortiges Loeschen wuerde das Ticket fuer alle unsichtbar machen, solange der Branch auf den Merge wartet
  proof: core/ticket/store.go Logbook/Archive/Delete -> recordDelete -> OpDelete; TestRecordDeleteTakesTheRefDown, TestDeleteTakesTheTicketOffEveryBoard
- [x] Tests mit zwei Klonen eines Bare-Repos belegen Rennen, Uebernahme und den Merge zweier Branches am selben Ticket
  proof: internal/cli/mergebranches_test.go: echter git merge zweier Branches, status nach Lane-Fortschritt, beide tags erhalten, keine Konfliktmarker
- [x] Die README erklaert die Funktion: wozu die Refs da sind, welche Befehle sie schreiben und lesen, und wie man das Board-Remote konfiguriert
  proof: README.md: wozu die Refs, welche Befehle schreiben/lesen, Board-Remote in settings.json
- [x] Die README nennt die Fork-Grenze ausdruecklich: Refs reisen nicht ueber Forks, alle Beteiligten pushen in dasselbe Board-Repo, und Teilnahme setzt Push-Recht darauf voraus - ein Contributor ohne Push-Recht bleibt bei Ticketdatei im Branch plus PR
  proof: README.md: 'The fork limit, stated plainly' - Refs reisen nicht ueber Forks, Push-Recht noetig, sonst Datei im Branch plus PR
