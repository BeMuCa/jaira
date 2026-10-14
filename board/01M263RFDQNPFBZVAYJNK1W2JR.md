---
id: 01M263RFDQNPFBZVAYJNK1W2JR
title: Ein verwaister Claim laesst sich aufloesen
status: done
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
commits:
  - 6bd21c73399d119c3edd146759e47737c2828ef0
  - 9081d17ead79671ba2cbe70ffbd9b087766f0d60
  - 9ba73f5f3a4956c8ddd68036a8fa78c61418fac6
  - 6d4566ebc9c5a9b0d27837d42187b8e543c47dce
created-at: 2026-09-10T16:52:33Z
updated-at: 2026-10-12T06:39:56Z
updated-by: Alexander Sacharov
outcome-what: "Review: Diff gegen DoD geprueft, Tests selbst gelaufen"
outcome-why: "fresh-context Review vor der Abnahme durch eine Person"
outcome-resolves: "Approve, keine blockierenden Luecken"
review-summary: "'jaira next' (auch --per-lane) haengt an seine Warnung ueber einen verwaisten Claim jetzt '; clear it with: jaira claim <id> --release' an (internal/cli/flow.go:510, warnStaleClaims). Den Befehl gab es schon: --release loest einen fremden Claim, sobald er aelter als ClaimTTL (30 min) ist, und einen live erneuerten verweigert er - der Kontext des Tickets irrte, es fehlte nur der Hinweis. Der Hinweis sitzt bewusst im Aufrufer und nicht in staleClaimLine, weil 'jaira claim' dieselbe Zeile nach einer Uebernahme druckt, wo Loesen falsch waere. Zwei Tests in internal/cli/staleclaim_test.go: next nennt den Befehl, genau dieser Befehl loest den Claim, danach verschwindet die Warnung; und die Uebernahme-Meldung von claim enthaelt kein --release. Eine Unreleased-Zeile in core/release/NOTES.md."
review-gaps: "Klein, nicht blockierend: (1) Der Test laeuft nur ueber 'next --all'; der Pfad 'next --per-lane' aus dem Ticket-Kontext (flow.go:466) ruft dieselbe warnStaleClaims auf, ist aber nicht eigens getestet. (2) Ein Claim ohne Zeitstempel bekommt denselben Hinweis - korrekt, denn ClaimActive behandelt ihn als abgelaufen, also gelingt --release. (3) Die Ticket-Notiz zu 15 roten Tests mit echtem HOME ist durch den Merge von master (9ba73f5) erledigt: go test ./... -race hier mit echtem HOME exit 0, 0 FAIL. Sonst keine Defekte; die DoD-Proofs (flow.go:510, staleclaim_test.go:84 und :75) stimmen mit dem Diff ueberein."
test-verdict: "pass: go test ./... -race RC=0 (echtes HOME, nach Merge master); DoD 1 im Baum: flow.go:510, staleclaim_test.go:84 + :74; auf Scratch-Board: claim mit Session 'dead', claimed-at auf -12h -> 'jaira next --all' meldet '... last renewed 12h ago; clear it with: jaira claim 47V14F --release', genau dieser Befehl -> 'Released', rc 0, naechstes next ohne Warnung"
review-verdict: "Approve. Der Diff erfuellt die DoD: Befehl zum Loesen existiert, die Meldung in 'jaira next' nennt ihn, ein Test deckt Meldung und Loesen ab, go test ./... -race gruen (selbst gelaufen). Aenderung ist minimal (eine Formatzeile + Kommentar), keine neue Abstraktion."
review-check: "1. go build -o /tmp/jaira ./cmd/jaira  2. In einem Wegwerf-Repo: /tmp/jaira init und /tmp/jaira create test --goal g --context c --dod d  3. In der Ticketdatei unter .jaira/tickets/ claimed_by: tote-session und claimed_at auf eine Zeit vor 2 Stunden setzen (nur fuer den Test von Hand)  4. /tmp/jaira next --all - auf stderr steht '... carries an abandoned claim by tote-session, last renewed 2h ago; clear it with: jaira claim <id> --release'  5. Genau diesen Befehl ausfuehren - Ausgabe 'Released <id>'  6. /tmp/jaira next --all erneut - die Warnung ist weg. Ohne Handarbeit: go test ./internal/cli -run 'TestNextNamesTheCommandThatClearsAnAbandonedClaim|TestClaimReportsTakingOverAnAbandonedClaim' -v ist gruen."
---

# Ein verwaister Claim laesst sich aufloesen

## Definition of Done

- [x] Ein verwaister Claim laesst sich per Befehl loesen; die Meldung in 'jaira next' nennt diesen Befehl; ein Test deckt Loesen und die Meldung ab; go test ./... -race gruen
  proof: internal/cli/flow.go:510 (warnStaleClaims nennt 'jaira claim <id> --release'); internal/cli/staleclaim_test.go:84 TestNextNamesTheCommandThatClearsAnAbandonedClaim (Meldung + Loesen), :74 in TestClaimReportsTakingOverAnAbandonedClaim (kein Hinweis bei Uebernahme); go test ./... -race exit 0

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
- **2026-10-11 20:55 · Alexander Sacharov** — in-progress (Runde 2): upstream/master gemergt (Konflikt nur NOTES.md, beide Unreleased-Zeilen behalten), HOME-Zeilen aus claimFixture gestrichen - claimrelease_test.go ist jetzt identisch mit master. go test ./... -race mit echtem HOME exit 0; die Notiz zum HOME-Trick ist damit erledigt.
- **2026-10-11 20:55 · Alexander Sacharov** — critique Runde 2: nichts mehr. Diff gegen master: flow.go eine Formatzeile + Kommentar, zwei Tests, eine NOTES-Zeile. HOME-Befund aus Runde 1 erledigt (9ba73f5).
- **2026-10-11 20:56 · Alexander Sacharov** — optimize: internal/cli/staleclaim_test.go: TestClaimTakeoverDoesNotSuggestReleasing wiederholte den Aufbau von TestClaimReportsTakingOverAnAbandonedClaim; die --release-Pruefung in diesen Test gefaltet, eigener Test geloescht (Mutation geprueft: Hinweis in staleClaimLine laesst ihn rot werden). Sonst nichts entfernt.
- **2026-10-11 20:57 · Alexander Sacharov** — critique/optimize Runde 3: nichts mehr nach 6d4566e.
- **2026-10-11 20:58 · Alexander Sacharov** — Worker: review an frischen Agenten per run-lane.sh --no-worktree uebergeben (Herdr-Tab). Teamlead: Abschaltung ~00:50, Alex schlaeft - Fragen gehen an jaira-61.
