---
id: 01M1CJ5TX2F4J6BPH29CJW7R34
title: jaira lanes add setzt die neue Lane neben ihren Anker statt ans Ende der order-Datei
status: backlog
ready: false
creator: BeMuCa
goal: "Eine per lanes add installierte Lane erscheint dort, wo ihr after:-Anker sie hinstellt - oder die Doku verspricht es nirgends"
context: "Fund des End-Reviews zu YBC0MT (31.08.): after: wird nur konsultiert, wenn keine order-Datei existiert; jaira init schreibt aber eine. Verifiziert: auf frischem Board haengt 'jaira lanes add secrets-scan' die Lane als LETZTE Zeile an .jaira/lanes/order - die Spalte landet rechts von blocked, nie hinter in-progress. Ursache core/lane/order.go:240 (ids = append(ids, l.ID)). Moeglicher Fix: am Anker einfuegen (slices.Insert nach indexOfID). ABER: order()-Verhalten ist Invarianten-Territorium (Anker-Reihenfolge routete schon einmal critique vor human; 'fix' der Sibling-Reihenfolge wurde bewusst zurueckgezogen) - deshalb Entscheidung statt Cleanup. Uebergangsweise wurde das README-Versprechen entschaerft (YBC0MT-Fix)."
definition-of-done: "Entweder fuegt lanes add am Anker ein (mit Test) oder die Doku sagt ehrlich, dass die Spalte ans Ende kommt und die order-Datei dem Nutzer gehoert"
blocked-by: []
commits: []
created-at: 2026-08-31T18:44:15Z
updated-at: 2026-08-31T18:44:15Z
---

# jaira lanes add setzt die neue Lane neben ihren Anker statt ans Ende der order-Datei

## Definition of Done

- [ ] Entweder fuegt lanes add am Anker ein (mit Test) oder die Doku sagt ehrlich, dass die Spalte ans Ende kommt und die order-Datei dem Nutzer gehoert

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

