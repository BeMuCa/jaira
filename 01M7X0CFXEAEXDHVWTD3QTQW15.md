---
id: 01M7X0CFXEAEXDHVWTD3QTQW15
title: "Give back im Cockpit startet den neuen Worker mit Teamlead, im Worktree mit dem Code und nie doppelt"
status: backlog
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
updated-at: 2026-11-20T20:01:18Z
updated-by: Alexander Sacharov
---

# Give back im Cockpit startet den neuen Worker mit Teamlead, im Worktree mit dem Code und nie doppelt

## Definition of Done

- [ ] Der Start-Prompt eines Give back nennt die parent session des Tickets und verbietet AskUserQuestion; danach geht eine Zeile '<ID> given back: <Entscheidung>, Worker <name>' an den Teamlead (Test im Cockpit)
- [ ] Give back findet den Branch mit den Commits des Tickets (Commits mit der ID, PR-Branch) und dessen Worktree und startet dort; nur ohne solchen Branch entsteht feat/<id> (Test: Branch mit anderem Namen)
- [ ] Laeuft fuer das Ticket schon ein Worker, startet Give back keinen zweiten, sondern schickt ihm die Entscheidung (Test: zweimal Give back hintereinander = ein Worker)
- [ ] Eingabeklassen: Ticket ohne parent-Note, parent-Session tot, zwei Branches mit Commits der ID, Worktree geloescht aber Branch da, zwei Give back gleichzeitig

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
