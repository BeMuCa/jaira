---
id: 01M0YT5W77TGZQW0M75GDNAEPN
title: "Die Versionszeile schweigt, wenn sie nichts zu sagen hat"
status: done
ready: true
creator: BeMuCa
goal: Auf einem Selbstbau steht keine Versionszeile in der Fusszeile
context: "Berk fragt, wofuer 'jaira dev' in der Fusszeile steht. Auf einem aus dem Quellcode gebauten Binary ist release.Current gleich 'dev', und versionLine (internal/tui/updatecheck.go) gibt dann genau 'jaira dev' zurueck - eine Zeile, die nichts beantwortet. Auf einem echten Release steht dort 'jaira 0.1.0 - up to date' oder der Upgrade-Hinweis, und der ist die Zeile wert. Gezeichnet wird sie an zwei Stellen: home.go:369 und view.go:703."
definition-of-done: Auf einem dev-Build erscheint die Zeile weder im Startbildschirm noch in der Board-Fusszeile; auf einem Release unveraendert; ein Test deckt beide Faelle
blocked-by: []
commits:
  - 73b73a023850e048df38f484f1f3353ba8124c98
  - 377c451a3057edaaf53e4769bd5dcc82564c004f
  - 6c5d81b28e18049c24bafd21a1d639c263474a32
created-at: 2026-08-26T10:34:43Z
updated-at: 2026-09-03T15:48:33Z
claimed-by: EE-3NX6GL3-2626872
claimed-at: 2026-08-31T16:29:26Z
updated-by: BeMuCa
assignee: BeMuCa
outcome-what: "Positiv-Anker in TestHomeFooterOmitsVersionLineOnDevBuild ('q quit') und TestBoardStatusBarOmitsVersionLineOnDevBuild ('? help') ergaenzt, in derselben Assertion wie die bestehende Negativ-Pruefung auf 'jaira dev'."
outcome-why: "Die reine Negativ-Pruefung waere auch gruen geblieben, wenn der ganze Footer aus einem fremden Grund verschwindet - der Anker beweist, dass tatsaechlich nur die Versionszeile fehlt und nicht der gesamte Footer."
outcome-resolves: "Anker-Strings gegen aktuellen Code verifiziert: home.go:391 fuehrt 'q quit' in der Hint-Liste, view.go:728 fuehrt '? help' in den Statusbar-keys. go test ./internal/tui -count=1 gruen, go test ./... -race -count=1 gruen, gofmt -l core internal listet weiterhin nur die vorbestehende internal/cli/tickets.go."
executed-by: sonnet
review-summary: none
review-gaps: "Runde 2: nichts entfernt; der Review-Fund (Negativ-Tests ohne Positiv-Anker) ist behoben - 'q quit' (home.go:391) und '? help' (view.go:728) ankern dieselbe Render-Ausgabe im selben Assertion-Pfad."
review-check: "go build -o /tmp/claude-1000/jaira-check ./cmd/jaira && /tmp/claude-1000/jaira-check im Repo starten: Startbildschirm und Board-Fusszeile ohne 'jaira dev', die Tastenzeilen (q quit / ? help) stehen da. go test ./internal/tui -run 'OmitsVersionLine|VersionLineSaysNothing|CarriesTheVersionIndicator' -count=1 -v: alles gruen."
review-verdict: "accept (Opus-Re-Review): beide Anker existieren in den aktuellen Hint-Listen und haengen am selben Render-String; hints_test pinnt die Strings unabhaengig, ein Rename bricht laut"
---

# Die Versionszeile schweigt, wenn sie nichts zu sagen hat

## Definition of Done

- [ ] Auf einem dev-Build erscheint die Zeile weder im Startbildschirm noch in der Board-Fusszeile; auf einem Release unveraendert; ein Test deckt beide Faelle

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-08-31 16:32 · BeMuCa** — home.go:394 und view.go:744 pruefen bereits auf versionLine != "" und blenden die Zeile dann aus — die Behebung brauchte nur eine Aenderung in updatecheck.go (return "" statt 'jaira dev'), keine Aenderung an den beiden Renderstellen. Der existierende Test devfooter_test.go pruefte genau das alte (jetzt falsche) Verhalten und musste umgeschrieben werden: statt auf 'if release.Current != "dev": t.Skipf' zu bauen, nutzt er jetzt setReleaseCurrent(t, "dev") aus updatecheck_test.go, damit der Test unabhaengig vom Build-Typ deterministisch laeuft.
- **2026-08-31 17:18 · BeMuCa** — review (Opus, batch): accept mit einem Testfund - die zwei neuen Footer-Tests sind reine Negativ-Substring-Pruefungen und bestehen auch, wenn der ganze Footer aus fremdem Grund verschwindet. Fix: Positiv-Anker in dieselbe Assertion ('q quit' im Home-Footer, '? help' in der Statusbar).
- **2026-08-31 17:46 · BeMuCa** — Review-Fund war berechtigt: die beiden neuen Footer-Tests pruefften nur Abwesenheit von 'jaira dev', nicht Anwesenheit des restlichen Footers - waeren also auch bei komplett leerem Rendering gruen geblieben. Anker verifiziert gegen den aktuellen Code: home.go:391 hat 'q quit' in der Hint-Liste, view.go:728 hat '? help' in den keys der Statusbar. Zwischenzeitlich kamen sechs fremde Commits vor meinem eigenen dazu (u.a. 7387693, clampBlock jetzt am Ende von Home.render) - Tests liefen auf dem aktuellen HEAD unveraendert gruen, keine Anpassung noetig.
