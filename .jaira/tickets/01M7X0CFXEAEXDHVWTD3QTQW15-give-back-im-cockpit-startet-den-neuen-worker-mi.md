---
id: 01M7X0CFXEAEXDHVWTD3QTQW15
title: "Give back im Cockpit startet den neuen Worker mit Teamlead, im Worktree mit dem Code und nie doppelt"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Give back aus dem Cockpit gibt dem neuen Worker die parent session des Tickets mit, setzt ihn in den vorhandenen Worktree/Branch mit den Commits des Tickets und startet keinen, wenn fuer das Ticket schon einer laeuft - dann bekommt der laufende die Entscheidung."
context: |-
  Gemeldet vom Teamlead eines anderen Projekts am 20.11.2026, an einem Tag mehrfach beobachtet:
  1. Der neue Worker kennt den Teamlead nicht und fragt die Person per AskUserQuestion im Chat, statt per SendMessage den Teamlead.
  2. Er startet in einem neuen leeren Worktree auf feat/<id> = main, obwohl die Commits auf einem anderen Branch liegen (z.B. dem PR-Branch mit anderem Namen). Drei Tickets betroffen.
  3. Fuer ein Ticket liefen drei Worker-Sessions gleichzeitig.
  Quelle der parent session heute: Note "parent session: uds:/..." am Ticket (spawn.sh schreibt sie), kein Feld.
  Der laufende Worker ist ueber den workers-Tool des Cockpits / ListAgents auffindbar.
  Verwandt: 76ERHH (Decide scheitert bei Ticket nur auf dem Ref - Punkt 4 der Meldung), R80Z4G (spawn.sh nimmt vorhandenen Branch wieder auf), X2TGXG (Worker findet seinen Teamlead nach Neustart).
definition-of-done: "Der Start-Prompt eines Give back nennt die parent session des Tickets und verbietet AskUserQuestion; danach geht eine Zeile '<ID> given back: <Entscheidung>, Worker <name>' an den Teamlead (Test im Cockpit)"
tags: []
blocked-by:
  - 01M7WX00JXRK922S029Q76ERHH
related:
  - 01M7WX00JXRK922S029Q76ERHH
  - 01M6M878R42NEHQ7HDP0R80Z4G
  - 01M5MNWDKYJGC78B3EKMX2TGXG
commits: []
created-at: 2026-11-20T17:34:35Z
updated-at: 2026-11-22T19:12:54Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1167205
claimed-at: 2026-11-22T18:52:27Z
outcome-what: "Give back im Cockpit: neuer Worker berichtet an die parent session des Tickets (sonst an die Cockpit-Session, sonst eigener Tab), startet auf dem Branch mit dem neuesten Commit der ID (vorhandener Worktree oder spawn.sh --branch), ein arbeitender Worker bekommt die Entscheidung statt ersetzt zu werden; Accept/Give back/Decide schicken '<ID>: ... (Person, Cockpit)' in die Session, Accept auch als Note. spawn.sh: --branch, vorhandener <prefix>/<slug> wird wieder ausgecheckt, Setup-Fehler loescht nur selbst angelegte Branches, Prompt mit parent nennt SendMessage statt AskUserQuestion."
outcome-why: "Ein Teamlead eines anderen Projekts sah am 20.11. Worker, die per AskUserQuestion fragten, in leeren feat/<id>-Worktrees starteten und dreifach liefen."
outcome-resolves: QTQW15 DoD 1-6
review-summary: |-
  core/mod/builtin/jaira-cockpit/hooks/register.tsx workPlace repeats the branch search decideStep already does inline (worktree branch, */<slug> listing, branch --contains the newest commit naming the id) - move the git queries into one helper both call, with carrierFor deciding
  core/mod/builtin/jaira-cockpit/hooks/register.tsx givingBack: a worker at work with no record is reported to the teamlead as 'worker none', which invites it to start the very second worker this press refuses to - the line must name the working tab and say the decision is to be typed there
round: 2
back-from: critique
back-by: lane
back-count: 1
---

# Give back im Cockpit startet den neuen Worker mit Teamlead, im Worktree mit dem Code und nie doppelt

## Definition of Done

- [x] Der Start-Prompt eines Give back nennt die parent session des Tickets und verbietet AskUserQuestion; danach geht eine Zeile '<ID> given back: <Entscheidung>, Worker <name>' an den Teamlead (Test im Cockpit)
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1424 + spawn.sh asks=: test 'the new worker reports to the parent session the ticket names, and the teamlead hears who took it'; Go TestSpawnTypesTheRoleForItsLaneName
- [x] Give back findet den Branch mit den Commits des Tickets (Commits mit der ID, PR-Branch) und dessen Worktree und startet dort; nur ohne solchen Branch entsteht feat/<id> (Test: Branch mit anderem Namen)
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1385 ticketBranch (shared with Accept's decideStep), workPlace; tests 'the branch the ticket\'s commits ride on wins over an empty feat/<id>...', 'a branch whose worktree was removed is checked out again...'; Go TestSpawnTakesUpTheBranchTheWorkRidesOn
- [x] Laeuft fuer das Ticket schon ein Worker, startet Give back keinen zweiten, sondern schickt ihm die Entscheidung (Test: zweimal Give back hintereinander = ein Worker)
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1344: tests 'a worker at work gets the decision...', 'a worker at work with nothing recorded...', 'Give back twice in a row: one worker...'
- [x] Eingabeklassen: Ticket ohne parent-Note, parent-Session tot, zwei Branches mit Commits der ID, Worktree geloescht aber Branch da, zwei Give back gleichzeitig
  proof: tests in core/mod/builtin/jaira-cockpit/tests/board.test.ts 'give back': no parent note -> 'no parent anywhere...'/'reports to this session (no parent noted)'; parent dead -> 'a parent session that is gone...'; two branches -> 'two branches hold the commit...'; worktree deleted -> 'a branch whose worktree was removed...' + Go TestSpawnTakesUpTheBranchTheWorkRidesOn; two at once -> 'Give back twice in a row' + existing 'a move another session made first'
- [x] Jede Entscheidung der Person im Cockpit (Accept, Give back, Decide) schickt dem Teamlead des Tickets sofort eine Zeile '<ID>: <Entscheidung> (Person, Cockpit)'; heute erfaehrt der Teamlead davon nichts und listet laengst entschiedene Tickets weiter als offen (Alex 20.11.2026: 'du bekommst die Info nicht schnell und wiederholst dich')
  proof: core/mod/builtin/jaira-cockpit/hooks/texts.ts toSessionAccepted/toSessionGivenBack/decisionFor: '<ID>: ... (Person, Cockpit)'; tests 'the line a press hands this session'
- [x] Entscheidung Teamlead 22.11.: DoD5 nur fuer die eigene Session - Accept, Give back und Decide geben ihre Zeile per prompt.submit in die Session des Cockpits und schreiben sie zusaetzlich als Note ans Ticket; ein Teamlead in einer anderen Session ist eine bekannte Grenze (Note)
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:679 Accept-Note; Give back/Decide notes vorhanden; test 'an option marked [accept] accepts the ticket as the band\'s Accept does'

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-20 17:34 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-22 18:52 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 18:56 · Alexander Sacharov** — decision (worker, settled): parent tot (DoD4) -> Give back nimmt die Session des Cockpits selbst als parent (CLAUDE_CODE_MESSAGING_SOCKET), wie Discuss heute; ohne beides --no-parent + Note. Zwei Branches mit Commits der ID -> der mit dem neuesten Commit, Note nennt den anderen. Branch ohne Worktree -> spawn.sh bekommt --branch <name> und legt den Worktree auf dem vorhandenen Branch an (statt -b feat/<id>). Laufender Worker = Tab <id>/worker|dispatch mit agent_status working -> bekommt die Entscheidung per send-text, egal was chooseWorker sagt; idle folgt weiter chooseWorker (der alte Tab wird vor dem Neustart geschlossen). Zwei Presses gleichzeitig -> Sperre je Ticket im Store (prozessuebergreifend), der zweite bekommt 'laeuft schon'.
- **2026-11-22 18:56 · Alexander Sacharov** — question: DoD5 - der Teamlead des Tickets ist eine andere Session als die des Cockpits; ein Mod kann kein SendMessage und prompt.submit erreicht nur die eigene Session. Wie kommt die Zeile hin?
1) Herdr send-text in den Pane des Teamleads: jede Session mit Cockpit legt im Store socket->pane ab, Give back/Accept/Decide tippt die Zeile mit [jaira-cockpit] dort ein (Recommended; sofort, wie heute beim alten Worker; ohne Herdr nur prompt.submit hier)
2) Briefkasten im Store: das Cockpit des Teamleads liest ihn beim naechsten Turn/Render und reicht die Zeile per prompt.submit ein (kein Herdr noetig, kommt aber erst beim naechsten Ereignis dort an)
3) Nur dieselbe Session (prompt.submit wie heute) + Note am Ticket; fremde Teamleads ausser Umfang (kleinstes, DoD5 nur teilweise)
- **2026-11-22 18:56 · Alexander Sacharov** — decision (teamlead 22.11.): F DoD5 = Option 3. Cockpit und Teamlead laufen bei Alex in derselben Session (Cockpit-Meldungen kommen heute schon per prompt.submit beim Teamlead an). prompt.submit in die eigene Session + note auf dem Ticket reicht; ein Teamlead in einer anderen Session ist ausserhalb des Umfangs (bekannte Grenze, als note).
- **2026-11-22 18:58 · Alexander Sacharov** — plan (worker): 1) spawn.sh: --branch <name> legt den Worktree auf vorhandenem Branch an; ein schon vorhandener <prefix>/<slug> wird wieder aufgenommen statt -b (deckt R80Z4G mit ab); scheitert die Setup, wird der Branch nur geloescht, wenn spawn.sh ihn selbst angelegt hat (sonst Datenverlust); bei worker/dispatch mit parent haengt der Prompt 'Fragen per SendMessage an <parent>, nie AskUserQuestion' an. 2) Cockpit givingBack: Branch = der mit dem neuesten Commit, der die ID nennt (eigener Worktree-Branch bevorzugt, wenn er ihn enthaelt), sonst Worktree-Branch, sonst */<slug>; Worktree da -> --no-worktree dort, sonst --branch. 3) parent = neueste Note 'parent session:', spawn exit 2 -> eigener Socket -> --no-parent. 4) Tab working -> Entscheidung per send-text, nie schliessen. 5) Zeilen an die Session im DoD-Format + Accept-Note. Zwei gleichzeitige Give back: jaira move antwortet dem zweiten 'already in in-progress' (movedBefore) - keine eigene Sperre noetig, nur Test.
- **2026-11-22 19:07 · Alexander Sacharov** — known limit (worker): Give back sucht Branches nur lokal (git log --branches). Ein PR-Branch, der nur noch auf origin liegt (lokal geloescht), wird nicht gefunden -> neuer feat/<id>. Ein Worktree der Ticket-ID mitten im Rebase (kein Branch) und Commits auf anderem Branch: spawn.sh verweigert mit 'already stands on no branch', die Note sagt es. Ein arbeitender Worker ohne Session-Record: kein Zweiter, die Note nennt seinen Tab - die Entscheidung muss dort eingetippt werden.
- **2026-11-22 19:09 · Alexander Sacharov** — critique: (1) workPlace und decideStep suchen den Branch des Tickets zweimal mit denselben drei git-Abfragen - ein Helfer fuer beide, carrierFor entscheidet. Folge fuer Accept: ein leerer feat/<id>-Worktree verliert gegen den Branch mit den Commits der ID, auch fuer den PR - gewollt, derselbe Fehler. (2) Arbeitender Worker ohne Record: die Zeile an den Teamlead sagt 'worker none' und laedt zum zweiten Worker ein - sie muss den Tab nennen.
- **2026-11-22 19:12 · Alexander Sacharov** — in-progress (Runde 2): beide critique-Funde behoben - ticketBranch ist der eine Helfer fuer Accept (decideStep) und Give back (workPlace); Accept waehlt damit auch den Branch mit dem neuesten Commit der ID statt eines leeren feat/<id> (Release-Note-Zeile dazu). Arbeitender Worker ohne Record: Zeile an den Teamlead nennt den Tab und 'start no other'. 370/370 Mod-Tests, go test core/role ok.
