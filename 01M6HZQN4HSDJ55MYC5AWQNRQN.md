---
id: 01M6HZQN4HSDJ55MYC5AWQNRQN
title: "Wer ein Ticket, einen Milestone oder ein Epic zugewiesen bekommt, sieht das im Board als Nachricht"
status: brainstorm
ready: true
creator: berkc
assignee: berkc
goal: "Bekommt jemand von einem Kollegen oder einem PM (jaira web) etwas zugewiesen, erscheint das im Board als Nachricht, bis er es gesehen hat"
context: |-
  Wunsch berkc 31.10.2026, Punkt 5, im Original abgebrochen: 'after you get assigned a ticket or milestone (or epic) from a colleague or the web user (PM) you see it in the board as a ...'. Am 04.11. ergaenzt: 'as a message or something - weiss aber nicht wie man das umsetzen kann, muss man brainstormen'.
  Heute: Zuweisung landet per Ref beim naechsten jaira fetch; es gibt eine Desktop-Benachrichtigung (8566KF), im Board selbst nichts.
  Offen fuer brainstorm: Form (Banner oben, Markierung auf der Karte, Posteingang-Fenster), wann sie verschwindet (erstes Oeffnen, Quittieren), wo 'gesehen' gespeichert wird (lokal, nicht im Ticket - sonst Commit-Laerm), Gegenstueck im Web. Haengt an DVPAPE (fester Name pro Person).
definition-of-done: "Brainstorm hat Form, Verschwinden und Speicherort festgelegt"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-11-04T00:35:52Z
updated-at: 2026-11-04T21:27:41Z
updated-by: berkc
claimed-by: berk-ca-1280500
claimed-at: 2026-11-04T21:24:22Z
question: "Welche Form: 1 lokaler Posteingang+Banner+Kartenmarker (empfohlen), 2 nur Kartenmarker, 3 Posteingang als Git-Ref pro Person? Und Milestone/Epic: B1 gebuendelte Ticketzuweisung (empfohlen) oder B2 neues owner-Feld?"
---

# Wer ein Ticket, einen Milestone oder ein Epic zugewiesen bekommt, sieht das im Board als Nachricht

## Definition of Done

- [ ] Brainstorm hat Form, Verschwinden und Speicherort festgelegt
- [ ] Nach einer Zuweisung durch jemand anderen zeigt das Board die Nachricht bis zum Gesehen; Test
- [ ] Release-Notiz unter core/release/notes/

## Options

- [x] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-04 21:27 · berkc** — BRAINSTORM (Code-Befund, dann Optionen, Empfehlung zuerst)

Was heute wirklich passiert (gelesen, nicht vermutet):
- Zuweisung kommt per refs/jaira/tickets/* an; refsync.Incoming (core/refsync/refsync.go:947) meldet Arrival{Mine, Changed}.
- TUI zeigt dann 30 s lang eine Zeile in der Fusszeile: "● CW2X0F is yours: <Titel>" (internal/tui/refs.go:150, view.go:1104, flashFor model.go:1001). Danach weg.
- "Gesehen" wird beim FETCH gespeichert, nicht beim Hinschauen: ~/.jaira/state/<tree>/refs-seen.json (refsync.go:123, saveSeen). Holt ein Agent oder 'jaira fetch' die Refs, ist die Nachricht verbraucht, bevor das Board sie je zeigt. Das ist der eigentliche Fehler.
- Changed = Ref-SHA hat sich bewegt, nicht "Assignee hat gewechselt": jede Notiz eines Kollegen auf einem Ticket, das schon mir gehoert, meldet erneut "is yours". Absender wird nicht gezeigt, obwohl updated-by im Ticket steht (core/ticket/schema.go:36).
- Milestones und Epics haben kein Assignee-/Owner-Feld (.jaira/milestones/*.md: nur color, created-at, Ticketliste). "Milestone zugewiesen" gibt es im Datenmodell nicht.
- Wer bin ich: JAIRA_USER > git user.name > $USER, Aliasse ~/.jaira/identity (core/identity/identity.go:16). DVPAPE (fester Name) ist noch backlog.
- jaira web ist read-only, kein Begriff "aktueller Nutzer".

Optionen:

1. (EMPFOHLEN) Bleibender Posteingang lokal + Banner + Kartenmarker
   Speicher: ~/.jaira/state/<repo>/inbox.json (RepoStateDir, alle Worktrees eines Clones teilen ihn), nie im Ticket.
   Eintrag entsteht, wenn der Assignee im Ref von "nicht ich" auf "ich" wechselt UND updated-by != ich. Verschwindet erst, wenn ich das Ticket oeffne oder 'x' druecke. 'jaira inbox [--json] [--ack]' fuer Agenten.
     jaira · default                                          berkc
     ✉ 2 neu fuer dich: CW2X0F von alex · pm-web (5 Tickets) von anna   i oeffnen  x gesehen
     backlog            todo               in-progress
     ┌──────────────┐   ┌──────────────┐
     │✉ CW2X0F alex │   │ ...          │
     │jaira web legt│
     └──────────────┘
     [i] ┌ Fuer dich ─────────────────────────────────┐
         │▸ CW2X0F  jaira web legt Tickets an  alex 2h │
         │  pm-web  Milestone, 5 Tickets        anna 1d│
         │ enter oeffnen · x gesehen · esc             │
         └─────────────────────────────────────────────┘
   Kostet: neue Datei + Diff alter/neuer Assignee in Incoming, ein Panel im TUI. Gibt auf: gesehen gilt pro Rechner, nicht pro Person.

2. Nur Markierung auf der Karte, gesehen = einmal geoeffnet
     todo
     ┌──────────────────┐
     │● NEU CW2X0F alex→│
     │jaira web legt ...│
     └──────────────────┘
   Kostet am wenigsten (ein Set von IDs lokal, ein Marker im Card-Render). Gibt auf: Karte in eingeklappter/gefilterter Lane oder auf anderem Board sieht man nicht; Milestone/Epic hat keine Karte; Agenten bekommen nichts.

3. Posteingang als Git-Ref pro Person (refs/jaira/inbox/<name>), quittieren pusht
   Gleiches Bild wie 1, zusaetzlich im Web und auf jedem Rechner derselbe Stand; ein PM sieht "noch nicht gesehen".
   Kostet: neue Ref-Art, Push bei jedem Quittieren, CAS-Konflikte, braucht DVPAPE zwingend (Name als Ref-Schluessel). Gibt auf: "lokal, kein Commit-Laerm" aus dem Kontext.

Warum 1: behebt genau den Fehler (gesehen beim Fetch statt beim Menschen), bleibt git-only ohne neuen Ref und ohne Laerm, funktioniert vor DVPAPE ueber identity.IsMe, und ein Web-Gegenstueck kann spaeter dieselbe Datei lesen (jaira web laeuft lokal). 3 bleibt als Ausbau, falls PMs "gesehen" sehen wollen.

Zweite Entscheidung (Milestone/Epic):
 B1 (empfohlen) Kein neues Feld: werden mehrere Tickets eines Milestones/Epics in einem Fetch mir zugewiesen, wird daraus EIN Eintrag "pm-web (5 Tickets) von anna".
 B2 Neues Feld owner: in Milestone- und Epic-Datei, eigene Nachricht "pm-web gehoert jetzt dir". Formataenderung, mehr Arbeit, eigenes Ticket wert.

Vorschlag geschaerftes Ziel: Weist mir jemand anderes ein Ticket zu (einzeln oder gebuendelt aus einem Milestone/Epic), zeigt das TUI-Board eine bleibende Nachricht mit Absender, bis ich das Ticket oeffne oder quittiere; "gesehen" liegt lokal unter ~/.jaira/state, nie im Ticket.
Vorschlag DoD:
 1 Eintrag nur bei Assignee-Wechsel auf mich durch jemand anderen; Notizen/eigene Zuweisung erzeugen keinen; Test core/refsync
 2 Eintrag ueberlebt 'jaira fetch' aus einer anderen Session und Board-Neustart, bis gesehen; Test
 3 TUI: Banner + Kartenmarker + Liste (i), Oeffnen oder x quittiert; TUI-Test
 4 'jaira inbox [--json] [--ack <id>|--all]'
 5 Mehrere Zuweisungen aus einem Milestone/Epic in einem Fetch = ein Eintrag
 6 Release-Notiz core/release/notes/WQNRQN.md
Nicht hier: Web-Gegenstueck (nach DVPAPE/CW2X0F eigenes Ticket).
