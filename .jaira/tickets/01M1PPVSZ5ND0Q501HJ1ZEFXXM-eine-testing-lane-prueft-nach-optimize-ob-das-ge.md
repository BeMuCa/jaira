---
id: 01M1PPVSZ5ND0Q501HJ1ZEFXXM
title: "Eine Testing-Lane prueft nach optimize, ob das Geforderte existiert und funktioniert"
status: critique
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Zwischen optimize und human sitzt eine agentische Testing-Lane: sie baut, laesst die Suite laufen, prueft die DoD Punkt fuer Punkt am Baum und exerziert das geaenderte Verhalten; Findings gehen mit Befund und Loesungsvorschlag (test-verdict-Feld + Note) zurueck nach Implementing, der Loop laeuft danach erneut durch critique/optimize/testing"
context: "Berk am 04.09.: 'man sollte ein Testing lane nach optimize setzen, die testet ob die geforderte Sache im Ticket umgesetzt wurde, und ob sie funktional ist' - Rueckweg nach in-progress mit Vermerk was falsch ist und wie man es loesen koennte. Ort: Katalog-Lane lanes/testing.md (wie critique/optimize, Invariante 13), auf diesem Board adoptiert; Reihenfolge ueber after: optimize + order-Datei. Findings-Uebergabe: test-verdict-Feld (Lane-declared, rendert als Lane fields) plus jaira note fuer den Implementing-Agenten - bewusst OHNE die builtin in-progress-Lane anzufassen (input-requires-Erweiterung waere eine Zeile, kostet aber Dauer-Drift-Warnung; Berks Zuruf). Tier cheap: die Lane fuehrt aus und vergleicht, die Urteils-Lanes davor/danach sind strong."
definition-of-done: lanes/testing.md existiert und parst; das Board zeigt Testing zwischen Optimize und HITL; ein Ticket kann testing nur mit test-verdict verlassen; rejects-to nennt in-progress; die Lane-Doku (lanes show testing) traegt Prompt und Kontrakt
tags: []
blocked-by: []
commits: []
created-at: 2026-09-04T17:18:34Z
updated-at: 2026-09-04T17:33:20Z
claimed-by: EE-3NX6GL3-34378
claimed-at: 2026-09-04T17:20:10Z
updated-by: BeMuCa
outcome-what: "Katalog-Lane lanes/testing.md (agentic, tier cheap, rejects-to in-progress, Output test-verdict): drei Paesse - Gates bauen/Suite mit -race, DoD Punkt fuer Punkt am Baum verifizieren, das geaenderte Verhalten selbst exerzieren; Findings als test-verdict fail + Note mit Befund/Beleg/Loesungsvorschlag zurueck nach Implementing; auf diesem Board adoptiert und zwischen optimize und human eingeordnet"
outcome-why: "Berk am 04.09.: nach optimize soll getestet werden, ob das Geforderte umgesetzt und funktional ist, mit Rueckweg samt Vermerk - der Loop laeuft danach erneut durch critique/optimize/testing"
outcome-resolves: Kontrakt sichtbar in jaira lanes show testing; Reihenfolge in jaira lanes; shipped-Parsing gruen; das Verlassen-ohne-verdict-Gate beweist dieses Ticket auf seinem eigenen Weg durch testing gleich selbst
executed-by: fable
---

# Eine Testing-Lane prueft nach optimize, ob das Geforderte existiert und funktioniert

## Definition of Done

- [x] lanes/testing.md existiert und parst; das Board zeigt Testing zwischen Optimize und HITL; ein Ticket kann testing nur mit test-verdict verlassen; rejects-to nennt in-progress; die Lane-Doku (lanes show testing) traegt Prompt und Kontrakt
  proof: lanes/testing.md parst (core/lane shipped-Tests gruen); adoptiert + added, order-Datei: optimize->testing->human ('jaira lanes' zeigt die Reihe); rejects-to in-progress und output test-verdict stehen im Kontrakt; Gate-Beweis folgt am eigenen Weg dieses Tickets durch testing

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-04 17:33 · BeMuCa** — Entscheidungen: tier cheap (die Lane fuehrt aus und vergleicht; Urteil liegt in critique/optimize/review) - eine Zeile im Lane-File, falls Berk strong will. Findings-Uebergabe: test-verdict-Feld + jaira note an den Implementing-Agenten; die builtin in-progress-Lane wurde bewusst NICHT angefasst (input-requires um test-verdict zu erweitern waere eine Zeile in .jaira/lanes/in-progress.md, kostet aber eine Dauer-Drift-Warnung gegen das builtin - Berks Zuruf). lanes add haengte testing ans Ende der order-Datei (bekannter Anker-Fall, backlog): von Hand hinter optimize gestellt.
