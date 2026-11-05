---
id: 01M6NHNZV41ZW2MYHS3YE2HSK1
title: "panda-bobr/jaira-plugin-template: lauffaehige Vorlage mit Beispielbefehl und Release-Pipeline"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer 'Use this template' drueckt und einen Tag setzt, hat ein Release, das die Katalog-CI ohne Aenderung annimmt."
context: |-
  Das Repo https://github.com/panda-bobr/jaira-plugin-template ist seit 05.11.2026 leer und als Template markiert.
  Inhalt: Go-Modul, Beispielbefehl 'jaira-hello' (liest 'jaira list --json', schreibt eine note), jaira-plugin.json, goreleaser-Konfig (Assets jaira-<id>_<os>_<arch>, checksums.txt), Action fuer Test und Release auf Tag, README mit den drei Schritten bis zum Katalog-PR.
  Format und Regeln: R6X6FH (docs/PLUGINS.md).
definition-of-done: "Ein Repo aus der Vorlage baut mit 'git tag v0.1.0 && git push --tags' ein Release mit allen Assets und checksums.txt; einmal vorgefuehrt"
tags:
  - ci
blocked-by:
  - R6X6FH
related: []
commits: []
created-at: 2026-11-05T09:47:15Z
updated-at: 2026-11-05T10:32:02Z
updated-by: Alexander Sacharov
---

# panda-bobr/jaira-plugin-template: lauffaehige Vorlage mit Beispielbefehl und Release-Pipeline

## Definition of Done

- [ ] Ein Repo aus der Vorlage baut mit 'git tag v0.1.0 && git push --tags' ein Release mit allen Assets und checksums.txt; einmal vorgefuehrt
- [ ] Das Release besteht die CI von panda-bobr/jaira-releases

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [~] Go-Modul + Befehl jaira-hello (greet: list --json lesen, note schreiben; manifest --json aus eingebettetem jaira-plugin.json; Versionspruefung JAIRA_VERSION) + Tests
- [x] goreleaser-Konfig (Assets jaira-<id>_<os>_<arch>, checksums.txt, jaira-plugin.json als Asset) + Actions: test auf push/PR, release auf Tag mit Pruefung Tag == version
- [x] README: drei Schritte von Use this template bis Katalog-PR
- [ ] Push nach main von panda-bobr/jaira-plugin-template, kein Tag
- [ ] panda-bobr/jaira-hello-demo aus der Vorlage anlegen, Tag v0.1.0, Release-Assets und checksums.txt pruefen -> DoD 1
- [ ] E2HSK1 nach blocked auf D44NVR (DoD 2 wartet auf die Katalog-CI)

## Progress
- **2026-11-05 10:23 · Alexander Sacharov** — Teamlead (05.11.2026): der Standard (R6X6FH, abgenommen) liegt noch nicht auf master, sondern auf origin/feat/r6x6fh: docs/PLUGINS.md und docs/schema/*.json - lesen mit 'git show origin/feat/r6x6fh:docs/PLUGINS.md'. Die Arbeit selbst gehoert in das Repo panda-bobr/e2hsk1; Alex ist dort admin. Ein Agent pusht seinen Branch und oeffnet keinen PR.
- **2026-11-05 10:24 · Alexander Sacharov** — Korrektur zur vorigen Teamlead-Notiz: das Ziel-Repo ist panda-bobr/jaira-plugin-template (nicht panda-bobr/e2hsk1). Das Release des Beispiels muss die CI aus D44NVR (panda-bobr/jaira-releases) bestehen; D44NVR laeuft parallel in einem eigenen Tab.
- **2026-11-05 10:24 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/1770.sock
- **2026-11-05 10:24 · Alexander Sacharov** — Worker-Befund vor dem Plan (05.11.2026): (a) panda-bobr/jaira-releases ist leer, keine CI - die kommt erst mit D44NVR (backlog). DoD 2 ist bis dahin nicht erfuellbar. (b) Die ganze Arbeit liegt im fremden Repo panda-bobr/jaira-plugin-template (oeffentlich, leer): Commit, Push nach main und ein Tag v0.1.0 erzeugen ein oeffentliches Release - ausserhalb dieses Worktrees, nach aussen sichtbar.
- **2026-11-05 10:24 · Alexander Sacharov** — Frage 1 (DoD 2 / D44NVR): 1) E2HSK1 auf D44NVR blockieren und D44NVR zuerst bauen (Recommended) - DoD bleibt ehrlich, eine Runde Wartezeit. 2) Vorlage jetzt bauen, DoD 1 abhaken, E2HSK1 danach in blocked bis D44NVR - Vorlage steht frueher, Ticket haengt offen. 3) DoD 2 ersetzen durch lokale Pruefung gegen das Format aus docs/PLUGINS.md - sofort fertig, aber keine echte Katalog-CI.
- **2026-11-05 10:24 · Alexander Sacharov** — Frage 2 (fremdes Repo): Darf der Worker in panda-bobr/jaira-plugin-template direkt nach main pushen und zum Vorfuehren v0.1.0 taggen (oeffentliches Release)? 1) Ja, push nach main + Tag v0.1.0 im Template selbst (Recommended) - einfachste Vorfuehrung, Release liegt danach im Template (spaeter loeschbar). 2) Push nach main, aber Vorfuehrung in einem Wegwerf-Repo aus der Vorlage (panda-bobr/jaira-hello-demo) - Template bleibt ohne Release, ein Repo mehr. 3) Nur Branch pushen, Alex merged und taggt selbst - kein Push nach main durch den Agenten.
- **2026-11-05 10:24 · Alexander Sacharov** — docs/PLUGINS.md (R6X6FH, done) steht nur auf feat/r6x6fh, nicht auf master. Die Vorlage richtet sich nach dem Stand dc278130 dort.
- **2026-11-05 10:27 · Alexander Sacharov** — Alex (05.11.2026, via Teamlead): Frage 2 -> Push in main von panda-bobr/jaira-plugin-template erlaubt, KEIN Tag/Release im Template selbst; die Demonstration laeuft in einem Einweg-Repo panda-bobr/jaira-hello-demo, angelegt aus dem Template (gh repo create --template), dort Tag v0.1.0. Frage 1 entscheidet der Teamlead: D44NVR laeuft bereits parallel in eigenem Tab, also Template jetzt bauen, DoD 1 abhaken, DoD 2 pruefen, sobald die CI von D44NVR auf main von jaira-releases liegt; bis dahin E2HSK1 nach blocked mit Verweis auf D44NVR.
- **2026-11-05 10:28 · Alexander Sacharov** — Worker: Mode bleibt autonom (zwei beantwortete Rahmenfragen, keine Designfrage offen). Weiter: todo -> pre-process.
- **2026-11-05 10:28 · Alexander Sacharov** — decision (worker): Demo-id = hello, Programm jaira-hello, ein Befehl greet. manifest --json gibt das per go:embed eingebettete jaira-plugin.json aus - so koennen Datei und Ausgabe nie auseinanderlaufen. Release laedt jaira-plugin.json zusaetzlich als Asset hoch, weil D44NVR die Datei am Default-Branch liest, ein spaeteres jaira plugin install sie aber am Release braucht. Release-Workflow bricht ab, wenn Tag != v<version> im Manifest. Arbeitskopie der Vorlage: Scratchpad dieser Sitzung; nach jedem Schritt gepusht, damit ein Neustart nichts verliert.
- **2026-11-05 10:31 · Alexander Sacharov** — Vorlage Schritt 1: go test gruen; gebautes jaira-hello gegen dieses Board: greet -> 198 Tickets exit 0, unbekannte id exit 5, JAIRA_VERSION=0.4.0 exit 3. Den note-Schreibweg nur im Unit-Test (gefaelschtes runJaira) geprueft - absichtlich keine Hello-Notiz auf das echte Board geschrieben. goreleaser lokal laufen zu lassen wurde vom Auto-Mode verweigert (Binary laden+ausfuehren); die Release-Konfig wird erst im Demo-Repo jaira-hello-demo echt geprueft.
