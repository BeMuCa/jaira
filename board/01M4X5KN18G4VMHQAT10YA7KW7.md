---
id: 01M4X5KN18G4VMHQAT10YA7KW7
title: "Eine Bitte ueber zwei Boards wird zu zwei verlinkten Tickets mit festem Beruehrungspunkt, und ihre Agenten tauschen sich aus"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Betrifft eine Bitte zwei Boards, teilt der Teamlead sie in ein Ticket pro Board, verlinkt sie, legt den Beruehrungspunkt fest (z.B. den API-Endpunkt), und ein Agent kann dem anderen mitteilen, was der braucht"
context: |-
  Was heute fehlt: Ein Ticket lebt auf einem Board. Arbeit, die Frontend und Backend zugleich betrifft, ist ein Ticket oder zwei lose Tickets ohne vereinbarte Schnittstelle.
  Berks Beispiel vom 14.10.: 'im Frontend soll noch ein Import-Button rein, der eine Datei annimmt und dann in einem S3-Bucket lagert'. Der Teamlead teilt das in 2 Tickets mit festem Kontaktpunkt, wo sich die beiden beruehren. Eines laeuft im ui-Board, eines im backend-Board, beide sind verlinkt.
  Dazu: 'es besteht die Moeglichkeit, dass ein Agent Informationen an den anderen uebertraegt, wie z.B. backend sagt dem frontend-Agenten den API-Endpunkt, den er fuer den Button bereitstellt.'
  Heute gibt es schon Links zwischen Tickets (related, follows, blocked-by) und Notizen. Es gibt aber keinen festen Beruehrungspunkt und keinen Weg, wie ein Worker einem anderen etwas uebergibt.
  Verwandt: F914QH (Epics: eine Funktion sammelt ihre Tickets), S0EG0M (mehrere Boards), der Teamlead-Routing-Ticket.
definition-of-done: "Der Teamlead legt fuer eine Bitte ueber zwei Boards zwei Tickets an, verlinkt sie und schreibt den Beruehrungspunkt (z.B. Endpunkt, Datenformat) auf beide"
tags: []
blocked-by: []
related:
  - 01M4X5JKRN4AV3KX56SES0EG0M
  - 01M4WS20HSNGR8XKP3RCF914QH
commits: []
created-at: 2026-10-14T12:18:47Z
updated-at: 2026-10-14T12:18:47Z
---

# Eine Bitte ueber zwei Boards wird zu zwei verlinkten Tickets mit festem Beruehrungspunkt, und ihre Agenten tauschen sich aus

## Definition of Done

- [ ] Der Teamlead legt fuer eine Bitte ueber zwei Boards zwei Tickets an, verlinkt sie und schreibt den Beruehrungspunkt (z.B. Endpunkt, Datenformat) auf beide
- [ ] Ein Worker kann dem verlinkten Ticket eine Uebergabe schreiben (z.B. 'POST /api/import, multipart, Antwort {key}'), die der andere Worker beim naechsten Lane-Input sieht; Test
- [ ] Aendert sich der Beruehrungspunkt, sieht der andere Worker das, bevor er weiterarbeitet; Test

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

