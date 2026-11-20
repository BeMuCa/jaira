---
id: 01M7WWYB5XZMD01DZY3MJG3REW
title: Ein Agent darf auf Wunsch den PR oeffnen und eine Review-Runde (Copilot) selbst fahren
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Eine Einstellung erlaubt Agenten, den Branch zu pushen, den PR zu oeffnen und eine Kommentarrunde abzuarbeiten (Guard, Push, Antwort je Kommentar mit Commit, resolve, Review neu anfordern); ohne die Einstellung bleibt alles wie heute."
context: |-
  Heute: Regel in CLAUDE.md und jaira-role-pr: ein Agent pusht und stoppt, PR oeffnen/mergen ist Sache des Menschen.
  Folge: ein Teamlead macht bei jedem Ticket und jeder Fix-Runde von Hand checkout, Guard, Push, PR-Text, Copilot anfordern, antworten, resolve.
  Alex 20.11.2026: er erlaubt Agenten inzwischen oefter PRs zu erstellen, auch mit Copilot als Reviewer, aber nicht ueberall -> optional, pro Repo/Board.
  Gewuenscht vom Teamlead: eine Rolle oder ein Skript pr-round <ID>; die Person gibt nur die Freigabe.
  Mergen und Approven bleiben in jedem Fall beim Menschen.
  Verwandt: AX8K02 (fester Schritt vor dem Push), XB97NJ (Board sagt, ob Arbeit per PR kommt).
definition-of-done: Ohne Einstellung oeffnet kein Agent einen PR (Test oder Skill-Text belegt)
tags: []
blocked-by: []
related:
  - 01M2HZKCQ48BPQYN1YPTXB97NJ
  - 01M5ZJ5DKQZ83GYBGPJHAX8K02
  - 01M2N12WVABVRYM7C923WATJDQ
  - 01M7WWYG213BWB06A7D1CRASZS
commits: []
created-at: 2026-11-20T16:34:26Z
updated-at: 2026-11-20T17:54:58Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-87225
claimed-at: 2026-11-20T17:38:56Z
question: ""
outcome-what: "Board-Dateien .jaira/pr (person|agent, Default person) und .jaira/pr-reviewer; settings.PRFor/PRReviewerFor; 'jaira roles list' zeigt pr/pr_source/pr_reviewer; jaira-role-pr: Abschnitt 'Whether this board lets an agent open it', create nur bei person-Aufruf oder agent+alle Tickets done, Reviewer anfordern, eine Kommentarrunde ueber neues scripts/pr-round.sh (threads, comments, reply=antworten+resolve, comment, request); teamlead/dispatcher nennen die Ausnahme; Release-Note."
outcome-why: "Ein Teamlead machte jede PR-Runde von Hand (push, PR-Text, Copilot anfordern, antworten, resolve); Alex erlaubt Agenten das inzwischen, aber nur pro Repo."
outcome-resolves: "Mit .jaira/pr=agent oeffnet jaira-role-pr den PR nach der Abnahme aller Tickets selbst, fordert den Reviewer an und faehrt eine Kommentarrunde; ohne Datei bleibt alles wie heute; Merge/Approve bleibt in jedem Fall beim Menschen."
review-summary: "none"
round: 5
review-gaps: |-
  note: pr-round.sh comments lists every conversation comment and review body each round, the agent's own earlier answers included - not answering twice there rests on the agent reading the list (DoD 10 covers threads only)
  note: the Copilot request uses the REST login copilot-pull-request-reviewer[bot]; tested only against the fake gh in core/role/pr_round_test.go, never against live GitHub
  note: reply checks then answers in two calls - two sessions on the same thread at once could each post the answer
  note: the gate 'every ticket on the branch is in done' lives only in jaira-role-pr/SKILL.md; nothing in the binary checks it
  note: this repository's own CLAUDE.md still says agents never open PRs - correct while this board has no .jaira/pr, but the two texts must be kept apart if one is ever written
test-verdict: "pass: go build + go test -count=1 -race ./... RC=0 auf 5db2eee4; DoD 1-11 im Baum geprueft; roles list auf Scratch-Board (default, agent+copilot aus CRLF, Shell-Zeichen -> leer, unlesbare Datei -> person); pr-round.sh threads/comments live gegen BeMuCa/jaira rc=0, falsche Nummer rc=1 ohne Schreiben; reply/request nur gegen gh-Stub (live wuerde auf echten PR schreiben); roles list schreibt nichts"
hunt-verdict: "pass: bughunt done at 5db2eee4, 0 change, 0 pre-existing (Runde 3; Runden 1-2: 4 change, alle behoben)"
review-verdict: "All 11 DoD items are met by the diff; settings, roles list, pr-round.sh and the role texts read whole show no defect. Related tests pass (go test ./core/settings ./core/role ./internal/cli). Only uncertainty: the Copilot reviewer request has not been tried against real GitHub."
review-round: 4
review-check: "1. cd into a scratch clone of a GitHub repo you own, with an open test PR (number N)  2. run: jaira roles list  -> last line reads 'pr: person - who opens a pull request (the default); reviewer: none'  3. echo agent > .jaira/pr; echo copilot > .jaira/pr-reviewer; jaira roles list  -> 'pr: agent ... (from .jaira/pr, shared with the board); reviewer: copilot'  4. echo Agent > .jaira/pr; jaira roles list  -> back to 'pr: person' (a typo never enables it)  5. bash core/role/builtin/jaira-role-pr/scripts/pr-round.sh request <owner/repo> N copilot  -> prints 'requested: copilot' and Copilot appears as reviewer on the PR (or exit 4 with GitHub's reason if Copilot is not enabled)  6. leave one review comment on the PR, then: pr-round.sh threads <owner/repo> N  -> one JSON line with its id  7. printf 'Fixed in abc123: test' > /tmp/a.txt; pr-round.sh reply <owner/repo> N <id> /tmp/a.txt  -> the thread shows the answer and is resolved; run the same line again -> 'resolved: already', no second answer"
---

# Ein Agent darf auf Wunsch den PR oeffnen und eine Review-Runde (Copilot) selbst fahren

## Definition of Done

- [x] Ohne Einstellung oeffnet kein Agent einen PR (Test oder Skill-Text belegt)
  proof: core/settings/settings_test.go TestPRDefaultsToThePerson; core/role/builtin/jaira-role-pr/SKILL.md:37 (person default -> hand the line back); core/role/pr_round_test.go TestPROpensOnlyWhereTheBoardSaysAndNeverAccepts
- [x] Mit Einstellung oeffnet jaira-role-pr den PR selbst und fordert den konfigurierten Reviewer (z.B. Copilot) an
  proof: core/role/builtin/jaira-role-pr/SKILL.md:37 (agent + every ticket done), core/role/builtin/jaira-role-pr/SKILL.md:343 Ask for the board's reviewer, core/role/builtin/jaira-role-pr/SKILL.md:325 pr-round.sh open; core/role/pr_round_test.go TestPRRoundRequestsTheReviewer, TestPROpensOnlyWhereTheBoardSaysAndNeverAccepts, TestPROpenRunsTheGateFirst
- [x] Eine Kommentarrunde: jeder Kommentar bekommt 'Behoben in <sha>: ...' oder eine Begruendung, wird resolved, das Review neu angefordert
  proof: core/role/builtin/jaira-role-pr/SKILL.md:367 Answering review comments; core/role/builtin/jaira-role-pr/scripts/pr-round.sh reply/comment/request; core/role/pr_round_test.go TestPRRoundRepliesThenResolves
- [x] Merge und Approve macht kein Agent, auch mit Einstellung
  proof: core/role/pr_round_test.go TestPRRoundNeverMergesOrApproves, TestPROpensOnlyWhereTheBoardSaysAndNeverAccepts ('.jaira/pr never reaches merge or approve')
- [x] Eingabeklasse PR existiert schon: kein zweiter PR, direkt zur Kommentarrunde
  proof: core/role/builtin/jaira-role-pr/SKILL.md:270 One listed -> straight to Answering review comments; the setting does not change that step
- [x] Eingabeklasse Fork vs. gleiches Repo: PR geht an das Ziel aus 'Which repository it goes to', die Kommentarrunde fragt dasselbe Repo
  proof: core/role/builtin/jaira-role-pr/SKILL.md:376 the round runs against <owner/repo> as settled in 'Which repository it goes to', never origin; pr-round.sh takes owner/repo explicitly, open names --repo and --base; core/role/pr_round_test.go TestPROpenRunsTheGateFirst (pr create --repo up/r --base master)
- [x] Eingabeklasse gh/glab nicht angemeldet oder ohne Scope: das Skript bricht mit klarer Meldung und Exit != 0 ab, nichts halb erledigt verschwiegen
  proof: core/role/pr_round_test.go TestPRRoundRefusesBeforeItWrites (not signed in -> exit 3, no write)
- [x] Eingabeklasse Kommentar ohne Thread (PR-Konversation, Review-Body): wird beantwortet, aber nicht als resolved gemeldet
  proof: core/role/pr_round_test.go TestPRRoundListsOnlyWhatIsOpen (comments listed apart), TestPRRoundCommentSaysItOnce (comment answers, never resolves, never twice)
- [x] Eingabeklasse Copilot nicht verfuegbar (Reviewer-Anfrage abgelehnt): wird gemeldet, PR bleibt offen, kein Abbruch der Runde
  proof: core/role/pr_round_test.go TestPRRoundRequestsTheReviewer (GH_POST=fail -> exit 4 with gh's reason); core/role/builtin/jaira-role-pr/SKILL.md:360
- [x] Eingabeklasse zweite Runde auf denselben Kommentaren: bereits resolved Threads werden nicht erneut beantwortet
  proof: core/role/pr_round_test.go TestPRRoundListsOnlyWhatIsOpen (resolved thread not listed), TestPRRoundReplyIsSafeToRunAgain, TestPRRoundCommentSaysItOnce (said in an earlier round -> nothing posted)
- [x] Eingabeklasse .jaira/pr/.jaira/pr-reviewer leer, CRLF, Kommentar, Tippfehler, Shell-Sonderzeichen: Tabellentest, Default person / Reviewer leer
  proof: core/settings/settings_test.go TestPRDefaultsToThePerson, TestPRReviewerIsALoginOrNothing (empty, CRLF, comment, typo, shell chars)
- [x] Das Gate 'alle Tickets des Branches in done' prueft das Binary: jaira roles pr-gate <base> mit Exit 3 und Liste der offenen Tickets; pr-round.sh open oeffnet nur nach bestandenem Gate
  proof: internal/cli/roles.go newRolesPRGateCmd, core/gitrepo/named.go TicketsNamed; internal/cli/prgate_test.go TestPRGateOpensOnlyWhenEveryTicketOfTheBranchIsDone, TestPRGateRefusesWhatItCannotJudge; core/role/pr_round_test.go TestPROpenRunsTheGateFirst
- [x] Zwei Sessions, die denselben Thread gleichzeitig beantworten, hinterlassen genau eine Antwort
  proof: core/role/builtin/jaira-role-pr/scripts/pr-round.sh dedupe (threads) and comment read-back; core/role/pr_round_test.go TestPRRoundReplyLeavesOneAnswerWhenTwoSessionsRace, TestPRRoundCommentSaysItOnce

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/settings: PRFor/PRReviewerFor aus .jaira/pr und .jaira/pr-reviewer, Default person; Tests
- [x] jaira roles list: pr, pr_source, pr_reviewer in Text und --json; Test
- [x] scripts/pr-round.sh in jaira-role-pr: threads, reply (antworten+resolve), rerequest; GitHub ueber gh api
- [x] jaira-role-pr SKILL.md: Abschnitt zur Einstellung, done-Gate, Reviewer anfordern, Kommentarrunde; teamlead/dispatcher/worker-Saetze angleichen; role_test
- [x] Release-Note core/release/notes/JG3REW.md, Suite, Commit

## Progress
- **2026-11-20 16:43 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/35112.sock
- **2026-11-20 16:44 · Alexander Sacharov** — decision (settled, worker): die Einstellung ist eine Board-Datei wie .jaira/release - .jaira/pr, eine Zeile 'person' (Default, wie heute) oder 'agent'; der Reviewer steht in .jaira/pr-reviewer (z.B. copilot). Grund: Vorbild ReleaseFor in core/settings/settings.go, reist mit dem Repo (XB97NJ: nicht settings.json, nicht git config); sichtbar ueber 'jaira roles list --json' wie release.
- **2026-11-20 16:44 · Alexander Sacharov** — decision (settled, worker): die Kommentarrunde laeuft ueber Rollentext plus ein Skript in jaira-role-pr/scripts (offene Threads listen, antworten, resolve per gh api graphql, Review neu anfordern) - der Teamlead wollte 'pr-round <ID>'; ein Skript macht den GraphQL-Teil wiederholbar. Merge/Approve bleibt in jedem Fall verboten.
- **2026-11-20 16:44 · Alexander Sacharov** — question: Wann darf ein Agent mit .jaira/pr=agent den PR oeffnen? 1) (Recommended) Sobald alle Tickets des Branches in done sind - wie der Cockpit-Accept heute; die Abnahme der Person ist die Freigabe, Copilot reviewt danach. 2) Schon in signoff - Copilot reviewt, bevor die Person abnimmt; die Person sieht Copilot-Ergebnis beim Accept, aber ein PR existiert fuer noch nicht abgenommene Arbeit. 3) Nur auf ausdrueckliche Freigabe pro PR (Teamlead gibt --by-person weiter, nachdem die Person ja gesagt hat) - Einstellung spart nur die Handarbeit, nicht die Frage.
- **2026-11-20 16:47 · Alexander Sacharov** — answer (Alex via teamlead, 2026-11-20): 1) - ein Agent oeffnet den PR mit .jaira/pr=agent erst, wenn alle Tickets des Branches in done sind; die Abnahme der Person ist die Freigabe, Copilot reviewt danach. Die selbst getroffenen Entscheidungen bleiben stehen. Mode bleibt autonom (eine beantwortete Frage).
- **2026-11-20 16:48 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026 via Teamlead, gelten bis TW5FE2): 1) DoD nennt Fehlerklassen (leer, gross, Sonderzeichen/CRLF, Leerzeichen-Pfade, Windows, zwei Prozesse, Fehler auf halbem Weg; Heuristik -> Annahmeregel). 2) in-progress prueft vor dem Move Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte; Tabellentests ueber die Klasse; Geschwister-Stellen per grep mitziehen; Fix an der Ursache. 3) Fund nennt Klasse+Beispiel, Rueckweg verlangt Fix der Klasse. 4) Blocker nur falsches Ergebnis/Datenverlust/Absturz/Sicherheitsluecke; nach Fix-Durchgang Test+Bughunt nur ueber Diff seit letztem Hunt, dann review; nach 3 Rueckgaben derselben Lane Frage an Teamlead. 5) kleine review-gaps in derselben Runde schliessen; Vorbestehendes als 'note: pre-existing: ...'.
- **2026-11-20 16:51 · Alexander Sacharov** — pr-round.sh request: gh pr edit --add-reviewer kennt keinen Spezialwert fuer Copilot (gh 2.83.1 --help sagt das ausdruecklich); darum REST POST requested_reviewers mit 'copilot-pull-request-reviewer[bot]'. Gegen ein echtes Repo nicht ausprobiert - nur gegen einen gh-Stub; ein abgelehnter Request endet mit Exit 4 und die Rolle meldet ihn. Der Stub wendet den --jq-Filter mit echtem jq an, damit der Filter selbst getestet ist.
- **2026-11-20 16:57 · Alexander Sacharov** — critique 1: (a) PRReviewerFor dupliziert modelName-Zeichenpruefung -> gemeinsamer Helfer. (b) GitLab-Runde: ':id' im glab api ist das Projekt des aktuellen Verzeichnisses (glab api --help: 'populated with values from the repository of the current directory'), im Fork also origin statt Ziel; glab mr note antwortet nicht im Thread. Fix: URL-kodiertes <owner/repo> und discussions-API.
- **2026-11-20 16:59 · Alexander Sacharov** — critique 2: beide Funde aus critique 1 behoben in d527d5a9 (word() in settings.go; GitLab ueber discussions-API mit <owner%2Frepo>). Nichts Neues.
- **2026-11-20 17:01 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-role-pr/scripts/pr-round.sh:71 - 'threads' passes owner and name with -F, which gh converts to a JSON integer when the value is all digits (a repo named 2048, a numeric login); the query declares them String!, which a GraphQL server refuses for an Int, so 'threads' fails on such a repository; use -f for owner and name; ran: GH_TOKEN=dummy GH_DEBUG=api gh api graphql -F owner=someone -F name=2048 ... (gh 2.83.1 sent GraphQL variables {"name":2048,...})
- **2026-11-20 17:01 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-role-pr/scripts/pr-round.sh:88-94 - 'reply' posts the answer and then resolves in a second call; when the resolve fails (exit 1) the thread is answered but stays open, and SKILL.md:404 says the next round picks it up because 'threads' still lists it, so the next round posts a second answer to the same thread; read
- **2026-11-20 17:01 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-role-pr/scripts/pr-round.sh:109 - a .jaira/pr-reviewer of '@alice' passes PRReviewerFor (core/settings/settings.go:323 allows @) but only '@copilot' loses its @; 'request' sends reviewers[]=@alice, which is not a login, so every review request fails with exit 4; read
- **2026-11-20 17:01 · Alexander Sacharov** — bughunt: done at d527d5a9 - 3 findings (3 change, 0 pre-existing) in 6 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/release/notes/JG3REW.md core/role/builtin/jaira-dispatcher/SKILL.md core/role/builtin/jaira-role-pr/SKILL.md core/role/builtin/jaira-role-pr/scripts/pr-round.sh core/role/builtin/jaira-teamlead/SKILL.md core/role/pr_round_test.go core/role/role_test.go core/settings/settings.go core/settings/settings_test.go core/shipped/shipped.sha256 internal/cli/roles.go internal/cli/roles_test.go
- **2026-11-20 17:02 · Alexander Sacharov** — testing 1: Klasse statt Beispiel - (1) jede GraphQL-String-Variable mit -f statt -F (auch thread id); (2) reply wird idempotent: steht als letzter Kommentar im Thread schon einer des angemeldeten Nutzers, wird nur noch resolved; (3) ein fuehrendes @ wird fuer jeden Reviewer entfernt, nicht nur copilot.
- **2026-11-20 17:05 · Alexander Sacharov** — critique 3 (nur Fix 1dd8e6d7): -f fuer Strings, reply prueft vor dem Schreiben, @ wird entfernt - passt, nichts Neues.
- **2026-11-20 17:06 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-role-pr/scripts/pr-round.sh:99 - reply treats a thread as already answered whenever its last comment is by the gh login itself, but agents run under the maintainer's own gh account: a line comment the maintainer left on the PR (or a follow-up he wrote last in a thread) makes reply skip the answer, print 'answered: already, resolving', resolve the thread and exit 0 - the request is closed unanswered; written by 1dd8e6d7 (JG3REW); ran: fake gh with viewer=alex and last comment author=alex, pr-round.sh reply o/r 7 PRRT_x answer.md -> exit 0, 0 addPullRequestReviewThreadReply calls, 1 resolveReviewThread call
- **2026-11-20 17:06 · Alexander Sacharov** — bughunt: done at 1dd8e6d7 - 1 finding (1 change, 0 pre-existing) in 4 modules, follow-up hunt (diff since d527d5a9; .jaira/milestones/* new to the list, read whole: board data, no code); concurrency, error paths, input, permissions checked; files: .jaira/milestones/abnahme-2026-11-A4140M.md .jaira/milestones/hackathon-vorbereitung-XD97YV.md .jaira/milestones/jaira-board-fan-out-2H8BYF.md .jaira/milestones/jaira-erweiterungen-abschauen-DM1Y32.md .jaira/milestones/jaira-improvement-findings-1-QXJ7YA.md .jaira/milestones/jaira-improvement-findings-2-J89MM5.md .jaira/milestones/jaira-improvement-findings-3-KGN25M.md .jaira/milestones/pm-web-A79Q9V.md .jaira/milestones/readme-videos-ZDEJZG.md .jaira/milestones/teamlead-intake-2S92X9.md .jaira/milestones/tickets-lanes-improvement-4E3ZYF.md core/release/notes/JG3REW.md core/role/builtin/jaira-dispatcher/SKILL.md core/role/builtin/jaira-role-pr/SKILL.md core/role/builtin/jaira-role-pr/scripts/pr-round.sh core/role/builtin/jaira-teamlead/SKILL.md core/role/pr_round_test.go core/role/role_test.go core/settings/settings.go core/settings/settings_test.go core/shipped/shipped.sha256 internal/cli/roles.go internal/cli/roles_test.go
- **2026-11-20 17:07 · Alexander Sacharov** — testing 2: Klasse = 'beantwortet' wird am Autor statt an der Antwort erkannt. Fix: als beantwortet gilt nur, wenn der letzte Kommentar vom selben Login stammt UND sein Text gleich dem Antworttext ist (CR und Rand-Leerraum ignoriert). Sonst wird gepostet.
- **2026-11-20 17:09 · Alexander Sacharov** — critique 4 (nur Fix 5db2eee4): beantwortet = gleicher Login UND gleicher Text; deckt Wiederholung ab, laesst Maintainer-Bemerkungen offen. Nichts Neues.
- **2026-11-20 17:10 · Alexander Sacharov** — bughunt: done at 5db2eee4 - 0 findings in 1 module, follow-up hunt (diff since 1dd8e6d7: pr-round.sh reply, pr_round_test.go, shipped.sha256; no new files); concurrency, error paths, input, permissions checked; files: .jaira/milestones/abnahme-2026-11-A4140M.md .jaira/milestones/hackathon-vorbereitung-XD97YV.md .jaira/milestones/jaira-board-fan-out-2H8BYF.md .jaira/milestones/jaira-erweiterungen-abschauen-DM1Y32.md .jaira/milestones/jaira-improvement-findings-1-QXJ7YA.md .jaira/milestones/jaira-improvement-findings-2-J89MM5.md .jaira/milestones/jaira-improvement-findings-3-KGN25M.md .jaira/milestones/pm-web-A79Q9V.md .jaira/milestones/readme-videos-ZDEJZG.md .jaira/milestones/teamlead-intake-2S92X9.md .jaira/milestones/tickets-lanes-improvement-4E3ZYF.md core/release/notes/JG3REW.md core/role/builtin/jaira-dispatcher/SKILL.md core/role/builtin/jaira-role-pr/SKILL.md core/role/builtin/jaira-role-pr/scripts/pr-round.sh core/role/builtin/jaira-teamlead/SKILL.md core/role/pr_round_test.go core/role/role_test.go core/settings/settings.go core/settings/settings_test.go core/shipped/shipped.sha256 internal/cli/roles.go internal/cli/roles_test.go
- **2026-11-20 17:25 · Alexander Sacharov** — explain: Alle 11 Punkte der Definition of Done sind erfüllt, Tests laufen grün. Offen ist nur: die Copilot-Anfrage (Login copilot-pull-request-reviewer[bot]) lief bisher nur gegen ein gefälschtes gh, nie gegen echtes GitHub.
Warum: Lehnt GitHub den Login ab, scheitert nur der Schritt 'Review neu anfordern' (Exit 4 mit Begründung); Push, Antworten und resolve sind davon nicht betroffen. Zwei kleine Lücken bleiben: das Gate 'alle Tickets im Branch sind done' steht nur im SKILL.md, und zwei Sessions könnten dieselbe Antwort doppelt posten.
1) Zurückgeben: das Gate im Binary prüfen und die Doppelantwort sperren, bevor es angenommen wird [give back]
2) Annehmen und Schritte 5-7 der Prüfliste bei der Abnahme gegen ein echtes Test-PR laufen lassen [accept]
Empfehlung: 2 - die Lücken sind dokumentiert und ungefährlich, der Live-Test ist die einzige Unsicherheit und gehört in die Abnahme.
- **2026-11-20 17:38 · Alexander Sacharov** — decision: 1) Zurückgeben: das Gate im Binary prüfen und die Doppelantwort sperren, bevor es angenommen wird
- **2026-11-20 17:38 · Alexander Sacharov** — Given back from signoff: 1) Zurückgeben: das Gate im Binary prüfen und die Doppelantwort sperren, bevor es angenommen wird
- **2026-11-20 17:38 · Alexander Sacharov** — cockpit: given back to a new worker - nothing recorded about the old worker session (cache and context unknown)
- **2026-11-20 17:40 · Alexander Sacharov** — runde 5 (zurueck aus signoff, Alex 17:38): decision (worker, settled): das Gate wird ein Kommando 'jaira roles pr-gate <base>' - liest .jaira/pr, die Handles aus 'git log --format=%s <base>..HEAD' (Scope type(ID[,ID]): wie ticketsNamed im Cockpit, board.ts), und jeden Handle auf dem Board (Store.List, ref-Stand eingeschlossen): nicht im Terminal-Lane seines Boards -> Exit 3 mit Liste; nicht auf dem Board (logbook) haelt nichts auf; kein Handle -> Exit 3. Unter roles, weil 'roles list' schon pr zeigt - kein neues Top-Level-Kommando. pr-round.sh bekommt 'open', das das Gate vor gh pr create laufen laesst; GitLab: SKILL ruft das Gate vor glab mr create.
- **2026-11-20 17:40 · Alexander Sacharov** — decision (worker, settled): Doppelantwort - nach dem Posten liest reply die Kommentare des Threads; der frueheste Kommentar desselben Logins mit demselben Text bleibt, ist es nicht der eigene, loescht reply den eigenen (deletePullRequestReviewComment). Gilt fuer jede Reihenfolge zweier Sessions, auch auf zwei Rechnern; ein flock wuerde nur eine Maschine abdecken und fehlt auf macOS.
- **2026-11-20 17:48 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/35112.sock (teamlead jaira-64) - Fragen per SendMessage, Schlussbericht dorthin. Interim-Regel 5: die review-gaps-notes werden in dieser Runde geschlossen: (1) comment prueft vor dem Posten, ob derselbe Text desselben Logins schon im Gespraech steht; (2) Copilot live bleibt Abnahmeschritt (review-check 5); (3) Doppelantwort und (4) Gate im Binary siehe oben; (5) CLAUDE.md nennt die Ausnahme. Dazu Bughunt-Rest: Vergleich ignoriert Leerraum am Zeilenende, falls GitHub ihn kappt.
- **2026-11-20 17:54 · Alexander Sacharov** — critique 5 (nur c471d660): nichts. Gelesen: pr-gate (roles.go, gitrepo/named.go), pr-round.sh open/said/first_of/dedupe/comment. Absichtlich so gelassen: reply erkennt 'schon beantwortet' am LETZTEN Kommentar, comment an IRGENDEINEM - im Thread muss nach einem Widerspruch des Reviewers erneut geantwortet werden, das Gespraech hat nichts zu resolven. Die Handle-Lesart ist board.ts ticketsNamed nachgebaut; das Cockpit hat einen anderen Ausloeser (Person drueckt Accept) und bleibt bei sich.
