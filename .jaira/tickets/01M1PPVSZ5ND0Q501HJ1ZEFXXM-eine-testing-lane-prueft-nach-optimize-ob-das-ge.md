---
id: 01M1PPVSZ5ND0Q501HJ1ZEFXXM
title: "Eine Testing-Lane prueft nach optimize, ob das Geforderte existiert und funktioniert"
status: review
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
updated-at: 2026-09-04T17:34:41Z
claimed-by: EE-3NX6GL3-34378
claimed-at: 2026-09-04T17:20:10Z
updated-by: BeMuCa
outcome-what: "Katalog-Lane lanes/testing.md (agentic, tier cheap, rejects-to in-progress, Output test-verdict): drei Paesse - Gates bauen/Suite mit -race, DoD Punkt fuer Punkt am Baum verifizieren, das geaenderte Verhalten selbst exerzieren; Findings als test-verdict fail + Note mit Befund/Beleg/Loesungsvorschlag zurueck nach Implementing; auf diesem Board adoptiert und zwischen optimize und human eingeordnet"
outcome-why: "Berk am 04.09.: nach optimize soll getestet werden, ob das Geforderte umgesetzt und funktional ist, mit Rueckweg samt Vermerk - der Loop laeuft danach erneut durch critique/optimize/testing"
outcome-resolves: Kontrakt sichtbar in jaira lanes show testing; Reihenfolge in jaira lanes; shipped-Parsing gruen; das Verlassen-ohne-verdict-Gate beweist dieses Ticket auf seinem eigenen Weg durch testing gleich selbst
executed-by: fable
review-summary: "Kritik: exakt das critique/optimize-Muster (Katalog-Datei, adopt+add, Invariante 13) - kein neuer Mechanismus, nur ein neuer Schritt; tier cheap begruendet (ausfuehren+vergleichen, nicht urteilen); der Prompt verlangt woertliche Fehlerausgaben statt Zusammenfassungen und verbietet Zertifizieren ohne Ausfuehren. Bewusst offen gelassen: builtin in-progress unangetastet (Findings via Feld+Note statt bounded input) - dokumentiert in der Note, Berks Zuruf."
review-gaps: "Nichts entfernt. Gelassen: test-verdict als EIN Feld (pass/fail + Kurzbefund) statt Feld-Paar - die Detailtiefe traegt die Note; kein eigener Farb-/UI-Support - Lane-Felder rendern generisch (KA9CFA haelt Zeilen)."
test-verdict: "pass: Gates gruen (core/lane -race nach Katalog-Erweiterung, Build ok), DoD am Baum verifiziert (Datei, Reihenfolge, Kontrakt), Verhalten exerziert - der Exit-Gate-Beweis lief an diesem Ticket selbst"
review-verdict: "accept (koordinator-verifiziert, offengelegt: Content-only - eine Katalog-Lane-Datei plus Board-Adoption, kein Go-Code; Beweise am lebenden Board: Reihenfolge in 'jaira lanes', Exit-Gate verweigerte den Move ohne test-verdict woertlich, shipped-Parsing -race gruen). Das erste echte Ticket durch die Lane ist die eigentliche Feuertaufe - sein Weg zeigt, ob der Prompt traegt."
---

# Eine Testing-Lane prueft nach optimize, ob das Geforderte existiert und funktioniert

## Definition of Done

- [x] lanes/testing.md existiert und parst; das Board zeigt Testing zwischen Optimize und HITL; ein Ticket kann testing nur mit test-verdict verlassen; rejects-to nennt in-progress; die Lane-Doku (lanes show testing) traegt Prompt und Kontrakt
  proof: lanes/testing.md parst (core/lane -race gruen); 'jaira lanes': optimize->testing->human; Gate BEWIESEN: 'move --to review' ohne verdict verweigert mit 'testing declares it produces test-verdict, which is still empty'; Kontrakt in lanes show testing

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-04 17:33 · BeMuCa** — Entscheidungen: tier cheap (die Lane fuehrt aus und vergleicht; Urteil liegt in critique/optimize/review) - eine Zeile im Lane-File, falls Berk strong will. Findings-Uebergabe: test-verdict-Feld + jaira note an den Implementing-Agenten; die builtin in-progress-Lane wurde bewusst NICHT angefasst (input-requires um test-verdict zu erweitern waere eine Zeile in .jaira/lanes/in-progress.md, kostet aber eine Dauer-Drift-Warnung gegen das builtin - Berks Zuruf). lanes add haengte testing ans Ende der order-Datei (bekannter Anker-Fall, backlog): von Hand hinter optimize gestellt.
