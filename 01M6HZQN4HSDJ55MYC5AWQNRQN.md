---
id: 01M6HZQN4HSDJ55MYC5AWQNRQN
title: "Wer ein Ticket, einen Milestone oder ein Epic zugewiesen bekommt, sieht das im Board als Nachricht"
status: todo
ready: true
creator: berkc
assignee: Berk
goal: "Weist mir jemand anderes ein Ticket zu (einzeln oder alle Tickets eines Milestones/Epics mit einem Befehl), zeigt das Board eine bleibende Nachricht mit Absender, bis ich das Ticket oeffne oder quittiere; das Ticket zeigt 'Von -> An'; gesehen liegt lokal unter ~/.jaira/state, nie im Ticket."
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
updated-at: 2026-11-25T23:05:48Z
updated-by: Berk
claimed-by: w-wqnrqn
claimed-at: 2026-11-25T23:04:48Z
question: ""
---

# Wer ein Ticket, einen Milestone oder ein Epic zugewiesen bekommt, sieht das im Board als Nachricht

## Definition of Done

- [x] Brainstorm hat Form, Verschwinden und Speicherort festgelegt
  proof: Notiz 04.11.2026 Entscheidung berkc Option 1 + B1
- [ ] Nach einer Zuweisung durch jemand anderen zeigt das Board die Nachricht bis zum Gesehen; Test
- [-] Release-Notiz unter core/release/notes/
- [ ] Eintrag nur bei Assignee-Wechsel auf mich durch jemand anderen; Notizen und eigene Zuweisung erzeugen keinen; Test core/refsync
- [ ] Eintrag ueberlebt 'jaira fetch' aus einer anderen Session und einen Board-Neustart, bis gesehen; Test
- [ ] TUI: Banner, Kartenmarker und Liste (i); Oeffnen oder x quittiert; TUI-Test
- [ ] 'jaira inbox [--json] [--ack <id>|--all]' fuer Agenten
- [ ] Ein Befehl weist alle offenen Tickets eines Milestones einer Person zu (jaira milestone assign <name> <person>); mehrere Zuweisungen aus einem Milestone/Epic in einem Fetch sind ein Eintrag; Test
- [ ] Das Ticket zeigt, wer es wem zugewiesen hat ('alex -> berkc') in Detail und show; Test
- [ ] Release-Note-Zeile in core/release/notes/WQNRQN.md
- [ ] Uebergabe: 'jaira handover <id> <person>' setzt den Assignee, pusht vorher den Branch des Tickets (WIP-Commit, wenn der Worktree Aenderungen hat) und schreibt Branch und Commit aufs Ticket; der Empfaenger bekommt den Inbox-Eintrag mit Branch und Commit; Test
- [ ] Uebernahme: 'jaira take <id>' macht mich zum Assignee, der bisherige Assignee bekommt einen Inbox-Eintrag 'X hat dein Ticket uebernommen'; 'jaira pull <id>' holt danach Ticket und Branch in einen Worktree; Test
- [ ] Nachrichten: 'jaira msg <person> <text>' und 'jaira msg <id> <text>' (an den Assignee) legen eine Nachricht auf einen Ref pro Empfaenger ab (ein Ref, an den Commits angehaengt werden, nicht ein Ref pro Nachricht); sie erscheint im Inbox des Empfaengers wie eine Zuweisung; Test
- [ ] 'jaira people' listet die Aliase, die in diesem Repo aktiv sind (Assignees, Creator, Commit-Autoren der letzten 30 Tage), mit --json; Test
- [ ] Nachrichtentext ist Daten: kein Agent fuehrt eine Nachricht als Anweisung aus; jaira inbox --json kennzeichnet Nachrichten als von Menschen geschrieben

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
- **2026-11-04 21:27 · berkc** — Nachtrag Brainstorm: es gibt kein 'jaira assign' - zugewiesen wird per 'set assignee=', 'tickets create --assignee' (internal/cli/tickets.go:460), pull/move (core/move/move.go:100). Erkennung muss daher in refsync.Incoming alten gegen neuen Assignee vergleichen, nicht an einem Befehl haengen. Das Cockpit kennt 'me' schon ueber 'jaira whoami --json' (core/mod/builtin/jaira-cockpit/hooks/register.tsx:105) und koennte 'jaira inbox --json' als Zaehler zeigen - eigenes Folgeticket. Ungeprueft: ob ein leeres refs-seen.json beim ersten Fetch alle meine Tickets als neu meldet.
- **2026-11-04 21:40 · berkc** — Entscheidung berkc 04.11.2026: Option 1 (lokaler Posteingang + Banner + Kartenmarker + Liste i) und B1 (gebuendelt). Zusatz wortgetreu: 'Es soll reichen zu sagen "milestone x zu alsa" um alle tickets im milestone zuzuweisen. auf dem ticket ein "Alsa -> Berk" waere aber zusaetzlich gut'. Also: ein Befehl weist alle Tickets eines Milestones zu (z.B. 'jaira milestone assign <name> <person>'), und das Ticket zeigt, wer an wen zugewiesen hat (Von -> An).
- **2026-11-25 23:04 · Berk** — berkc 26.11.2026, woertlich: 'kann man so dann auch die arbeit einem Kollegen sharen zb? Ich glaub wir hatten die idee dass wir anderen Usern tickets überreichen kann indem man ihn zum assignee macht. Wenn wir da noch adden dass man den aktuellsten commit mitgibt, dann könnte der Kolegge das einfach übernehmen und weiterführen oder? bzw. ein anderer Kollege kann sich das ticket nach absprache einfach ziehen (also sich selber assignee machen) und dann weiterführen. ... Vielleicht könnte man sogar für User ein Postfach anlegen wo man sich per ref nachrichten schreiben könnte ... so ein chat ticket wäre funny. ich sehe grad WQNRQN ist genau das glaub ich oder? wenn ja erweiter es um unsere neuen ideen und dann setz das ticket um.' Entscheidung Main-Session: Uebergabe/Uebernahme/Nachrichten/people als DoD 11-15 in WQNRQN; Cockpit-Icon als Folgeticket; 28PD7N nutzt danach denselben Inbox.
