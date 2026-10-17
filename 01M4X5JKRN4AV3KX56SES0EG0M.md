---
id: 01M4X5JKRN4AV3KX56SES0EG0M
title: "Ein Repo kann mehrere Boards haben, jedes mit eigenen Lanes fuer sein Thema"
status: in-progress
ready: true
creator: BeMuCa
assignee: Alexander Sacharov
goal: "In einem Repo laufen mehrere Boards nebeneinander (z.B. backend, ui, video), jedes mit eigenen Lanes; ein Ticket gehoert genau einem Board und laeuft durch dessen Lanes"
context: |-
  Was heute fehlt: Ein Repo hat genau ein Board. Ein Board ist sein Lane-Ordner (.jaira/lanes/, siehe README 'A board is its lane directory'). Alle Tickets laufen durch dieselben Lanes.
  Unsere Lanes passen fuer Software-Entwicklung (SWE). Fuer andere Arbeit passen sie nicht.
  Berks Idee vom 14.10., seine Beispiele:
  - UI-Board: Lanes speziell fuer UI, z.B. verlangt testing dort immer Playwright-Tests.
  - Video-Board: Lanes fuer Ideen-Generierung, Skript, Voice-over und Schnitt.
  - Backend-Board: z.B. fuer 'in der Datenbank fuer user eine Spalte hinzufuegen'.
  Jedes ist ein eigenes Board. Der Teamlead verteilt Tickets darauf (eigenes Ticket), und ein jaira-worker bekommt Ticket plus Board und treibt es durch dessen Lanes.
  Offen fuer brainstorm:
  - wo die Boards liegen (z.B. .jaira/boards/<name>/lanes/) und wie ein Ticket sein Board traegt (Feld board:)
  - wie TUI und Launcher mehrere Boards eines Repos zeigen
  - wie es mit dem heutigen Einzel-Board vertraeglich bleibt (Berks Regel: ein Update darf nichts verschwinden lassen)
  Haengt am Ticket-Schema (BGZ0SW, NFJCTK): eine Lane eines solchen Boards muss eigene Felder erklaeren koennen.
definition-of-done: "Design als Spec unter docs/superpowers/specs/, von Berk freigegeben: Ablage der Boards, Board-Feld am Ticket, Anzeige in TUI und Launcher, Migration des heutigen Boards"
tags: []
blocked-by: []
related:
  - 01M4X5J1QRR70AM2B2JPBGZ0SW
  - 01M11YYR0SMVCXNQPZ68NFJCTK
commits: []
created-at: 2026-10-14T12:18:13Z
updated-at: 2026-10-17T12:19:47Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-73711
claimed-at: 2026-10-17T12:17:54Z
---

# Ein Repo kann mehrere Boards haben, jedes mit eigenen Lanes fuer sein Thema

## Definition of Done

- [x] Entscheidungen stehen als decision-Notes am Ticket (Ablage .jaira/boards/<name>/lanes/, Feld board:, Boardwechsel, Migration); keine eigene Spec, Alexander hat die Richtung am 17.10. vorgegeben
  proof: decision-Notes vom 17.10. am Ticket (Ablage, Name/Existenz, Feld board:/Migration, list/next-Default, Boardwechsel)
- [ ] TUI: Shift+F waehlt ein Board; das Board zeigt dann dessen Lanes und nur dessen Tickets; ohne Wahl das heutige Board; Test. Mehr TUI-Darstellung nicht
- [ ] jaira create, list, next und show nehmen ein Board (z.B. --board ui); ohne Angabe gilt das heutige Board, nichts aendert sich still; Test
- [ ] Ein bestehendes Einzel-Board laeuft nach dem Update unveraendert weiter; Test
- [ ] Eine Zeile in core/release/notes/<ID>.md
- [ ] Ein Ticket wechselt das Board mit jaira set <id> board=<name>; existiert seine Lane dort nicht, landet es in der ersten Lane des neuen Boards und jaira sagt das; Test

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-17 12:17 · Alexander Sacharov** — Alexander Sacharov took this ticket over from BeMuCa
- **2026-10-17 12:17 · Alexander Sacharov** — decision 17.10. (Alexander, Teamlead-Session): Alexander uebernimmt das Ticket von Berk und grenzt den Umfang ein: kein von Berk freigegebener Spec, TUI nur Filter per Shift+F. Tickets bekommen KEINE Ordner pro Board — sie leben bis zum Logbook auf refs/jaira/tickets/<ID>, der Pfad traegt keine Information. Zugehoerigkeit ist das Frontmatter-Feld board:; fehlt es, gilt das heutige Board (.jaira/lanes/), so laeuft jedes bestehende Board unveraendert. Ordner nur fuer die Lanes eines weiteren Boards: .jaira/boards/<name>/lanes/ mit eigener order. Boardwechsel = eine Zeile im Frontmatter statt Rename (kein rename/modify-Konflikt). Ausloeser: ~/projects/demo-video braucht ein Board 'video' (Skript -> human -> Bau, Szenen parallel -> Probe-Render -> human -> Stimme -> Finalrender), damit eine Designerin ohne Programmierkenntnis Videos nach Vorlage anlegt, waehrend Code und Animationen (shared/, scripts/) weiter durch die SWE-Lanes laufen.
- **2026-10-17 12:18 · Alexander Sacharov** — parent session: jaira-79
- **2026-10-17 12:19 · Alexander Sacharov** — decision: Ablage. Ein weiteres Board ist .jaira/boards/<name>/lanes/ mit Lane-Dateien und eigener order-Datei. Das heutige Board (.jaira/lanes/) hat keinen Namen; auf der Kommandozeile steht 'default' fuer es (--board default, set board=default oder board= leer). Ein Board-Name ist klein, Ziffern und Bindestriche, 'default' ist reserviert.
- **2026-10-17 12:19 · Alexander Sacharov** — decision: Ein benanntes Board existiert genau dann, wenn sein lanes/-Ordner Lane-Dateien enthaelt. Ein unbekannter Name ist ein Fehler (Exit 2) und nennt die vorhandenen Boards. Kein automatisches Anlegen wie beim heutigen Board: ein Tippfehler in --board darf kein neues Board erzeugen. Ein Board legt man an, indem man den Ordner mit Lane-Dateien fuellt (z.B. aus .jaira/lanes/ kopiert); kein eigenes Kommando - Scope.
- **2026-10-17 12:19 · Alexander Sacharov** — decision: Feld board: im Frontmatter. Fehlt es, gehoert das Ticket zum heutigen Board - das ist die ganze Migration: bestehende Tickets und .jaira/lanes/ bleiben unberuehrt, nichts wird umgeschrieben.
- **2026-10-17 12:19 · Alexander Sacharov** — decision: jaira list und next ohne --board zeigen nur Tickets des heutigen Boards, wie die TUI ohne Wahl. Kommandos auf ein einzelnes Ticket (show, move, dod, claim, set, note) brauchen kein --board: sie lesen das Board vom Ticket und pruefen es gegen dessen Lanes.
- **2026-10-17 12:19 · Alexander Sacharov** — decision: Boardwechsel ist jaira set <id> board=<name>, eine Zeile im Frontmatter. Gibt es die Lane des Tickets im neuen Board nicht, setzt jaira status auf die erste Lane in der order des neuen Boards und sagt das auf stdout.
