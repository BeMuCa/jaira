---
id: 01M11YYR1T040049KA8YB4MGTP
title: "Ein Ticket zeigt seine Geschichte faltbar: Fragen, Antworten, Kritik und Loesung pro Runde"
status: backlog
ready: false
creator: BeMuCa
goal: "Auf den ersten Blick zeigt ein Ticket nur seine Grundform; was es in HITL, critique und optimize durchlaufen hat, ist aufklappbar und vollstaendig"
context: |-
  Berk am 27.08.: definieren, wie ein Ticket aussieht - Grundform zu Beginn, welche Felder jede Stufe hinzufuegt, was eingefaltet und was immer sichtbar ist. Beispiel: ein Ticket, das HITL durchlief, bekommt einen faltbaren Block mit Frage und Antwort; eine critique/in-progress-Schleife zeigt pro Runde die Kritik und wie sie geloest wurde. Nichts soll den ersten Blick verstopfen, nichts verloren gehen.

  Was heute da ist: review-summary, review-gaps, review-verdict, review-check und question sind EINWERTIG - jede Runde ueberschreibt die letzte. Der critique-Prompt sagt deshalb 'Findings zusaetzlich als Notiz schreiben, damit sie den naechsten Ueberschreib ueberleben'. Die Geschichte steckt also in ## Progress als Notizen, ohne Struktur. Die Detailansicht (internal/tui/view.go renderDetail) zeigt Felder flach, nichts ist faltbar. TEMPLATE.md pro Board bestimmt die Grundform.

  Zwei Teile: (1) Datenformat - Runden im Body als Abschnitte (## Critique 1, ## Answer ...), hand-editierbar und diff-lesbar, per CLI geschrieben statt einwertige Felder zu ueberschreiben. (2) Anzeige - Abschnitte in der Detailansicht falten (Taste), Grundform immer offen. Erst (1), (2) folgt daraus.
definition-of-done: "Eine zweite critique-Runde ueberschreibt die erste nicht mehr, beide sind im Ticket lesbar; eine Frage an HITL und ihre Antwort stehen als Paar im Ticket; die Detailansicht zeigt die Grundform offen und die Runden eingefaltet, eine Taste klappt auf; das Format ist im Ticket-Markdown von Hand lesbar"
blocked-by: []
commits: []
created-at: 2026-08-27T15:55:56Z
updated-at: 2026-09-10T16:53:34Z
updated-by: BeMuCa
---

# Ein Ticket zeigt seine Geschichte faltbar: Fragen, Antworten, Kritik und Loesung pro Runde

## Definition of Done

- [ ] Eine zweite critique-Runde ueberschreibt die erste nicht mehr, beide sind im Ticket lesbar; eine Frage an HITL und ihre Antwort stehen als Paar im Ticket; die Detailansicht zeigt die Grundform offen und die Runden eingefaltet, eine Taste klappt auf; das Format ist im Ticket-Markdown von Hand lesbar

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-10 16:53 · BeMuCa** — Aus dem Board-Audit vom 09./10.09. auf einem anderen Projekt: dieses Ticket ist nicht mehr nur Komfort, es ist die Voraussetzung, um ueberhaupt pruefen zu koennen, ob ein Ticket seine Lanes durchlaufen hat.
Konkret: 11 Tickets standen in review, 8 ohne review-verdict. Ich konnte NICHT rekonstruieren, welchen Weg sie genommen haben - das Ticket traegt nur updated-at und updated-by, keine Historie. Dass critique und optimize uebersprungen wurden, weiss ich einzig, weil ein Agent es zufaellig in review-summary von 9GEGTB hineingeschrieben hat.
Ohne Historie ist 'ist das Board sauber gelaufen' keine Frage, die man beantworten kann - nur eine, die man raet.
Haengt an D28H7V (Sprung ueber eine Lane faellt auf): dessen DoD verlangt, dass der Sprung am Ticket ablesbar ist, und genau dafuer braucht es hier ein Feld.
