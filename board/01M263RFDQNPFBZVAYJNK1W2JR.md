---
id: 01M263RFDQNPFBZVAYJNK1W2JR
title: Ein verwaister Claim laesst sich aufloesen
status: in-progress
ready: true
creator: BeMuCa
assignee: Alexander Sacharov
goal: "Ein Claim, dessen Session gestorben ist, laesst sich mit einem Befehl loesen und blockiert das Ticket nicht auf Dauer"
context: |-
  Auf dem requirementsgenie-Board: 'jaira next --per-lane' meldet '5M11NC carries an abandoned claim by EE-3NX6GL3-690539, last renewed 11h ago'. Zwei Tage vorher meldete dasselbe Board '2CQ613 ... last renewed 436h ago' - also 18 Tage.
  jaira erkennt den verwaisten Claim also und benennt ihn praezise. Was fehlt, ist der Befehl daneben: 'jaira --help' listet claim, aber nichts, was einen Claim wieder loest, und er laeuft auch nicht von selbst ab.
  NICHT geprueft und vor dem Bauen zu klaeren: ob es intern doch eine Ablaufzeit gibt, die nur nicht greift - ich habe das aus der Befehlsliste geschlossen, nicht im Code nachgelesen.
  Warum es zaehlt: mehrere parallele Sessions schreiben dieses Board, und eine getoetete Session hat nie einen Zug, ihren Claim zurueckzugeben. Genau der Fall, fuer den das Board gebaut ist.
definition-of-done: "Ein verwaister Claim laesst sich per Befehl loesen; die Meldung in 'jaira next' nennt diesen Befehl; ein Test deckt Loesen und die Meldung ab; go test ./... -race gruen"
tags:
  - cli
blocked-by: []
commits: []
created-at: 2026-09-10T16:52:33Z
updated-at: 2026-10-11T20:51:37Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-20501
claimed-at: 2026-10-11T20:45:29Z
outcome-what: "jaira next nennt neben der Warnung zum verwaisten Claim den Befehl 'jaira claim <id> --release', der ihn loest"
outcome-why: "Ablauf (30 min) und --release gab es schon; die Warnung wiederholte sich bei jedem next ohne Hinweis, was zu tun ist"
outcome-resolves: "Ein verwaister Claim laesst sich mit dem Befehl loesen, den die Meldung nennt"
review-summary: "internal/cli/claimrelease_test.go:16 setzt HOME/USERPROFILE in claimFixture von Hand; upstream/master hat dafuer TestMain + internal/testhome (AK3Z6D, 860b776) - master in die Branch mergen und die vier Zeilen samt filepath-Import streichen; dann ist go test ./... -race auch ohne HOME-Trick gruen"
---

# Ein verwaister Claim laesst sich aufloesen

## Definition of Done

- [x] Ein verwaister Claim laesst sich per Befehl loesen; die Meldung in 'jaira next' nennt diesen Befehl; ein Test deckt Loesen und die Meldung ab; go test ./... -race gruen
  proof: internal/cli/flow.go:510 (warnStaleClaims nennt 'jaira claim <id> --release'); internal/cli/staleclaim_test.go:79 TestNextNamesTheCommandThatClearsAnAbandonedClaim, :120 TestClaimTakeoverDoesNotSuggestReleasing; go test ./... -race exit 0 mit leerem HOME

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-11 19:45 · Alexander Sacharov** — Code aus dem A/B-Versuch vom 11.10.2026 uebernommen (Arm A; blinder Richter knapp fuer diese Variante; Review: approve). Befund: Claim-Ablauf (ClaimTTL 30 min) und 'jaira claim <id> --release' gab es schon, der Kontext irrte; die Aenderung nennt den Befehl in der Warnung von jaira next. Offen laut Review: der DoD-Proof nannte releaseHint, das inzwischen in warnStaleClaims eingefaltet ist; '--release' ist in 'jaira --help' nicht sichtbar (nicht DoD). Die Fixture-Reparatur in claimrelease_test.go ueberschneidet sich mit AK3Z6D (HOME in TestMain). Alex hat das Ticket von BeMuCa uebernommen.
- **2026-10-11 20:45 · Alexander Sacharov** — parent session: jaira-61
- **2026-10-11 20:45 · Alexander Sacharov** — Worker: keine offenen Entscheidungen. Code liegt schon in 9081d17. '--release' steht in 'jaira claim --help' (Flag + Absatz); dass 'jaira --help' Flags von Unterbefehlen nicht listet, ist cobra-Standard und nicht DoD.
- **2026-10-11 20:50 · Alexander Sacharov** — go test ./... -race faellt auf dieser Ticket-Branch mit dem echten HOME durch (15 Tests, u.a. TestArchiveFromTheBoard): die Tests lesen ~/.jaira/default-board.md der Maschine - das ist AK3Z6D, gefixt in upstream/master (860b776), aber nicht in dieser Branch. Mit leerem HOME (HOME=/tmp/... USERPROFILE=...) exit 0. Nicht hier gemergt: gehoert nicht zum Ticket; beim Merge der PR auf master verschwindet es.
- **2026-10-11 20:51 · Alexander Sacharov** — critique: internal/cli/claimrelease_test.go:16 setzt HOME/USERPROFILE in claimFixture von Hand; upstream/master hat dafuer TestMain + internal/testhome (AK3Z6D, 860b776) - master in die Branch mergen und die vier Zeilen samt filepath-Import streichen; dann ist go test ./... -race auch ohne HOME-Trick gruen. Sonst nichts: Hinweis im Aufrufer warnStaleClaims statt in staleClaimLine ist richtig (claim druckt dieselbe Zeile nach der Uebernahme, Kommentar flow.go:503ff), Format eine Zeile, keine neue Abstraktion.
