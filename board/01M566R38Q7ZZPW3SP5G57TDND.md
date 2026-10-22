---
id: 01M566R38Q7ZZPW3SP5G57TDND
title: "In den Einstellungen waehlt man, welche Boards der Teamlead fuer neue Tickets nehmen darf"
status: backlog
ready: true
creator: berkc
assignee: BeMuCa
goal: "Der Mensch schaltet in den Einstellungen Boards fuer den Teamlead an oder aus; der Teamlead ordnet neue Tickets nur aktiven Boards zu, nach Name und Beschreibung."
context: |-
  Berk am 18.10.: 'Der Teamlead soll entscheiden, welches Board fuer welche Aufgabe genommen wird. Die Boards brauchen eine Beschreibung, die dem Agenten sagt, ob ein Ticket dorthin gehoert. Und der User kann in den Einstellungen auswaehlen, welche Boards fuer den Teamlead aktiv sind. Wir koennen z.B. zwei Boards fuer UI-Tickets haben, dann waehlt man eins davon aus.'
  Was es schon gibt (seit 0.5.3, S0EG0M/8NHD8Q done): jedes Board hat eine Datei 'description' in seinem Lanes-Ordner; 'jaira boards' listet Name, offene Tickets und Beschreibung; core/role/builtin/jaira-teamlead/SKILL.md (ab Zeile ~107) liest 'jaira boards' und legt ein Ticket auf das Board, dessen Beschreibung passt - passen zwei oder fehlt eine Beschreibung, fragt er den Menschen.
  Was fehlt: es gibt keinen Schalter 'aktiv fuer den Teamlead'. Zwei passende Boards fuehren heute immer zur Rueckfrage.
  Offen fuer brainstorm: wo der Schalter liegt - im Repo und geteilt (das Team sieht dieselbe Auswahl) oder pro Person in ~/.jaira/settings.json; was 'jaira boards' und 'jaira create' fuer ein inaktives Board tun (nur der Teamlead filtert, oder alle); Anzeige im Settings-Screen (S) der TUI.
  Verwandt: Y2C9QN (Boards anlegen aus Typen, dort muss eine Beschreibung Pflicht werden).
definition-of-done: "Im Settings-Screen (S) gibt es eine Liste der Boards mit Name, Beschreibung und einem Schalter aktiv/inaktiv"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-18T00:31:52Z
updated-at: 2026-10-18T00:31:52Z
---

# In den Einstellungen waehlt man, welche Boards der Teamlead fuer neue Tickets nehmen darf

## Definition of Done

- [ ] Im Settings-Screen (S) gibt es eine Liste der Boards mit Name, Beschreibung und einem Schalter aktiv/inaktiv
- [ ] 'jaira boards' zeigt den Schalter und '--json' gibt ihn aus
- [ ] jaira-teamlead nimmt nur aktive Boards in Betracht; ein Test oder Skill-Text pinnt das
- [ ] Tests; Release-Note

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

