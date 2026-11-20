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
updated-at: 2026-11-20T16:59:54Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-59469
claimed-at: 2026-11-20T16:43:41Z
question: ""
outcome-what: "Board-Dateien .jaira/pr (person|agent, Default person) und .jaira/pr-reviewer; settings.PRFor/PRReviewerFor; 'jaira roles list' zeigt pr/pr_source/pr_reviewer; jaira-role-pr: Abschnitt 'Whether this board lets an agent open it', create nur bei person-Aufruf oder agent+alle Tickets done, Reviewer anfordern, eine Kommentarrunde ueber neues scripts/pr-round.sh (threads, comments, reply=antworten+resolve, comment, request); teamlead/dispatcher nennen die Ausnahme; Release-Note."
outcome-why: "Ein Teamlead machte jede PR-Runde von Hand (push, PR-Text, Copilot anfordern, antworten, resolve); Alex erlaubt Agenten das inzwischen, aber nur pro Repo."
outcome-resolves: "Mit .jaira/pr=agent oeffnet jaira-role-pr den PR nach der Abnahme aller Tickets selbst, fordert den Reviewer an und faehrt eine Kommentarrunde; ohne Datei bleibt alles wie heute; Merge/Approve bleibt in jedem Fall beim Menschen."
review-summary: none
round: 2
review-gaps: "none - geprueft: PRFor folgt bewusst ReleaseFor (gleiches Muster, eigene Werte), Zeichenpruefung seit critique 1 eine Funktion word(); pr-round.sh hat kein Gegenstueck im Repo (grep requested_reviewers/resolveReviewThread: nur hier); kein toter Code, alle neuen Funktionen haben Aufrufer; roles list liest zwei kleine Dateien mehr, kein heisser Pfad"
---

# Ein Agent darf auf Wunsch den PR oeffnen und eine Review-Runde (Copilot) selbst fahren

## Definition of Done

- [x] Ohne Einstellung oeffnet kein Agent einen PR (Test oder Skill-Text belegt)
  proof: core/settings/settings_test.go TestPRDefaultsToThePerson; core/role/builtin/jaira-role-pr/SKILL.md:37 (person default -> hand the line back); core/role/pr_round_test.go TestPROpensOnlyWhereTheBoardSaysAndNeverAccepts
- [x] Mit Einstellung oeffnet jaira-role-pr den PR selbst und fordert den konfigurierten Reviewer (z.B. Copilot) an
  proof: core/role/builtin/jaira-role-pr/SKILL.md:37 (agent + every ticket done), core/role/builtin/jaira-role-pr/SKILL.md:323 Ask for the board's reviewer; core/role/pr_round_test.go TestPRRoundRequestsTheReviewer, TestPROpensOnlyWhereTheBoardSaysAndNeverAccepts
- [x] Eine Kommentarrunde: jeder Kommentar bekommt 'Behoben in <sha>: ...' oder eine Begruendung, wird resolved, das Review neu angefordert
  proof: core/role/builtin/jaira-role-pr/scripts/pr-round.sh reply/request; core/role/pr_round_test.go TestPRRoundRepliesThenResolves; core/role/builtin/jaira-role-pr/SKILL.md:350 Answering review comments
- [x] Merge und Approve macht kein Agent, auch mit Einstellung
  proof: core/role/pr_round_test.go TestPRRoundNeverMergesOrApproves, TestPROpensOnlyWhereTheBoardSaysAndNeverAccepts ('.jaira/pr never reaches merge or approve')
- [x] Eingabeklasse PR existiert schon: kein zweiter PR, direkt zur Kommentarrunde
  proof: core/role/builtin/jaira-role-pr/SKILL.md:270 One listed -> straight to Answering review comments; the setting does not change that step
- [x] Eingabeklasse Fork vs. gleiches Repo: PR geht an das Ziel aus 'Which repository it goes to', die Kommentarrunde fragt dasselbe Repo
  proof: core/role/builtin/jaira-role-pr/SKILL.md:356 the round runs against <owner/repo> as settled in 'Which repository it goes to', never origin; pr-round.sh takes owner/repo explicitly
- [x] Eingabeklasse gh/glab nicht angemeldet oder ohne Scope: das Skript bricht mit klarer Meldung und Exit != 0 ab, nichts halb erledigt verschwiegen
  proof: core/role/pr_round_test.go TestPRRoundRefusesBeforeItWrites (not signed in -> exit 3, no write)
- [x] Eingabeklasse Kommentar ohne Thread (PR-Konversation, Review-Body): wird beantwortet, aber nicht als resolved gemeldet
  proof: core/role/pr_round_test.go TestPRRoundListsOnlyWhatIsOpen (comments listed apart); pr-round.sh comment answers without resolving
- [x] Eingabeklasse Copilot nicht verfuegbar (Reviewer-Anfrage abgelehnt): wird gemeldet, PR bleibt offen, kein Abbruch der Runde
  proof: core/role/pr_round_test.go TestPRRoundRequestsTheReviewer (GH_POST=fail -> exit 4 with gh's reason); core/role/builtin/jaira-role-pr/SKILL.md:340
- [x] Eingabeklasse zweite Runde auf denselben Kommentaren: bereits resolved Threads werden nicht erneut beantwortet
  proof: core/role/pr_round_test.go TestPRRoundListsOnlyWhatIsOpen (resolved PRRT_done not listed)
- [x] Eingabeklasse .jaira/pr/.jaira/pr-reviewer leer, CRLF, Kommentar, Tippfehler, Shell-Sonderzeichen: Tabellentest, Default person / Reviewer leer
  proof: core/settings/settings_test.go TestPRDefaultsToThePerson, TestPRReviewerIsALoginOrNothing (empty, CRLF, comment, typo, shell chars)

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
