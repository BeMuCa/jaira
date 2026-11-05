---
id: 01M6NHNZV41ZW2MYHS3YE2HSK1
title: "panda-bobr/jaira-plugin-template: lauffaehige Vorlage mit Beispielbefehl und Release-Pipeline"
status: review
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
updated-at: 2026-11-05T17:05:43Z
updated-by: Alexander Sacharov
blocked-reason: "DoD 2 wartet auf die Katalog-CI aus D44NVR (panda-bobr/jaira-releases, noch nicht auf main)"
outcome-what: "panda-bobr/jaira-plugin-template@fc5d7e1: Go-Plugin jaira-hello (greet liest jaira list --json und schreibt jaira note; manifest --json aus eingebettetem jaira-plugin.json; Exit 3 bei zu altem JAIRA_VERSION), Tests inkl. Manifest-Regeln, goreleaser (6 Binaries, checksums.txt, jaira-plugin.json), Actions test und release-auf-Tag mit Tag==version-Pruefung, README mit drei Schritten bis Katalog-PR. Demo panda-bobr/jaira-hello-demo v0.1.0; Katalog-PR panda-bobr/jaira-releases#4."
outcome-why: "Plugin-Autoren sollen ohne eigene Release-Technik ein Release bekommen, das die Katalog-CI aus D44NVR ohne Aenderung annimmt (Standard R6X6FH, docs/PLUGINS.md)."
outcome-resolves: "DoD 1: Tag v0.1.0 im unveraenderten Demo-Repo erzeugte alle 6 Assets + checksums.txt (Run 42011801567). DoD 2: check-catalog auf PR #4 gruen (Run 42026722958)."
review-summary: "Im Repo panda-bobr/jaira-plugin-template (main = 66de84a) liegt eine Go-Vorlage: jaira-hello mit einem Befehl greet (liest 'jaira list --json', zaehlt Tickets; mit Ticket-Handle schreibt es eine Note ueber 'jaira note'), 'manifest --json' gibt das per go:embed eingebaute jaira-plugin.json aus, Exit 3 bei zu altem JAIRA_VERSION, jaira-Exitcodes werden durchgereicht. Ein Test prueft das Manifest gegen die Regeln aus docs/PLUGINS.md. goreleaser baut 6 Binaries jaira-<id>_<os>_<arch>(.exe), checksums.txt und laedt jaira-plugin.json mit hoch; die id kommt nur aus dem Manifest. Der release-Workflow bricht ab, wenn Tag != v<version>. README fuehrt in drei Schritten bis zum Katalog-PR. Bewiesen: das unveraenderte Demo-Repo panda-bobr/jaira-hello-demo (nur 'Initial commit' aus der Vorlage) hat Release v0.1.0 mit allen 6 Binaries, checksums.txt und jaira-plugin.json (Run 42011801567 gruen); der Katalog-PR panda-bobr/jaira-releases#4 (nur plugins.json) hat den Check 'check' gruen. Im jaira-Repo selbst aendert sich nur die Ticket-Datei."
round: 2
review-gaps: "none - main.go hat keinen Code ohne Aufrufer mehr (version in 66de84a entfernt); .goreleaser.yaml project_name/changelog waeren entbehrlich, bleiben aber: die Konfig ist in jaira-hello-demo v0.1.0 so bewiesen, lokal laesst sich goreleaser hier nicht laufen"
test-verdict: "pass - template main 66de84a: test-Workflow Run 42070699388 success (go vet + go test); Demo jaira-hello-demo v0.1.0 Release-Run 42011801567 success, 6 Binaries + checksums.txt, sha256sum -c alle OK; Katalog-PR jaira-releases#4 check-catalog Run 42026722958 pass. Lokal gegen dieses Board: greet exit 0, unbekannte id exit 5, zu altes JAIRA_VERSION exit 3."
claimed-by: DESKTOP-RFTCH11-9733
claimed-at: 2026-11-05T17:04:34Z
---

# panda-bobr/jaira-plugin-template: lauffaehige Vorlage mit Beispielbefehl und Release-Pipeline

## Definition of Done

- [x] Ein Repo aus der Vorlage baut mit 'git tag v0.1.0 && git push --tags' ein Release mit allen Assets und checksums.txt; einmal vorgefuehrt
  proof: https://github.com/panda-bobr/jaira-hello-demo/releases/tag/v0.1.0 - Tag v0.1.0 -> Run 42011801567 gruen; 6 Binaries jaira-hello_<os>_<arch>(.exe), checksums.txt (sha256sum -c alle OK), jaira-plugin.json; Vorlage = panda-bobr/jaira-plugin-template@fc5d7e1
- [x] Das Release besteht die CI von panda-bobr/jaira-releases
  proof: https://github.com/panda-bobr/jaira-releases/pull/4 - check-catalog run 42026722958 pass, Kommentar 'Catalogue check passed' fuer hello -> panda-bobr/jaira-hello-demo v0.1.0

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [ ] Go-Modul + Befehl jaira-hello (greet: list --json lesen, note schreiben; manifest --json aus eingebettetem jaira-plugin.json; Versionspruefung JAIRA_VERSION) + Tests
- [x] goreleaser-Konfig (Assets jaira-<id>_<os>_<arch>, checksums.txt, jaira-plugin.json als Asset) + Actions: test auf push/PR, release auf Tag mit Pruefung Tag == version
- [x] README: drei Schritte von Use this template bis Katalog-PR
- [x] Push nach main von panda-bobr/jaira-plugin-template, kein Tag
- [x] panda-bobr/jaira-hello-demo aus der Vorlage anlegen, Tag v0.1.0, Release-Assets und checksums.txt pruefen -> DoD 1
- [x] E2HSK1 nach blocked auf D44NVR (DoD 2 wartet auf die Katalog-CI)
- [x] Katalog-PR vorbereiten: Branch feat/e2hsk1 in panda-bobr/jaira-releases mit hello-Eintrag; Alex pusht und oeffnet den PR

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
- **2026-11-05 10:32 · Alexander Sacharov** — Vorlage fertig als Commit fc5d7e1 (feat(E2HSK1)) in der Arbeitskopie /tmp/claude-1000/-home-alex-projects--worktrees-jaira-e2hsk1/e870470a-cbe4-4763-84f0-083ada1efa9f/scratchpad/jaira-plugin-template - NOCH NICHT GEPUSHT. Der Push nach main wurde vom Auto-Mode-Klassifikator verweigert ('Create Public Surface'); Alex' Erlaubnis kam nur ueber den Teamlead, das zaehlt fuer die Sitzung nicht. Alex muss selbst pushen und jaira-hello-demo anlegen (Zeilen gingen an den Teamlead). Die Arbeitskopie lebt im Scratchpad dieser Sitzung: stirbt die Sitzung vor dem Push, ist sie aus diesem Ticket nicht wiederherstellbar.
- **2026-11-05 10:35 · Alexander Sacharov** — Gefunden im Demo-Lauf: goreleaser haengt .exe bei formats: [binary] selbst an (bestaetigt: jaira-hello_windows_amd64.exe). jaira-plugin.json ist Release-Asset, steht aber nicht in checksums.txt (extra_files werden nicht gehasht) - fuer D44NVR egal, die prueft nur Binaries + checksums.txt. Weiter, sobald D44NVR auf main von jaira-releases liegt: ein Katalog-PR (oder dessen Check) fuer hello -> https://github.com/panda-bobr/jaira-hello-demo muss gruen sein -> DoD 2. Danach Ticket-Datei mit einem Commit chore(E2HSK1) hier mitnehmen; Code liegt im Template-Repo.
- **2026-11-05 12:06 · Alexander Sacharov** — DoD 2 vorbereitet: Branch feat/e2hsk1 (8fdfee6) in Klon /tmp/claude-1000/-home-alex-projects--worktrees-jaira-e2hsk1/e870470a-cbe4-4763-84f0-083ada1efa9f/scratchpad/jaira-releases, nicht gepusht (Teamlead: Alex pusht und oeffnet). Lokal mit check_catalog.py von main 44d293d und --author sashasoft90 geprueft: einziger Fehler ist die Mitgliedschaft (sashasoft90 nicht oeffentliches Mitglied von panda-bobr); Schema, id, jaira-plugin.json, Release v0.1.0 und checksums.txt bestehen. DoD 2 gilt erst bei gruenem Check im echten PR. Merge listet hello wirklich im Katalog - ob der Test-PR gemergt oder geschlossen wird, ist Alex' Wahl.
- **2026-11-05 12:24 · Alexander Sacharov** — critique: (1) test.yml push nur auf main, pull_request bleibt - sonst Doppellauf je PR. (2) version/--version aus main.go und usage raus - nicht im Standard, manifest --json reicht. README-Link auf BeMuCa/jaira bleibt bewusst: der Umzug nach panda-bobr ist entschieden, aber noch nicht passiert. jaira-hello-demo bleibt auf dem alten Stand fc5d7e1 - es beweist DoD 1, kein Grund es nachzuziehen.
- **2026-11-05 12:25 · Alexander Sacharov** — Kritik-Fixes als 66de84a in der Template-Arbeitskopie (Scratchpad), noch nicht gepusht - Alex pusht. go vet/test gruen.
- **2026-11-05 12:25 · Alexander Sacharov** — critique Runde 2: beide Befunde aus Runde 1 in 66de84a behoben (test.yml:5, main.go ohne version); sonst nichts. goreleaser-Konfig bewusst nicht angefasst.
- **2026-11-05 12:26 · Alexander Sacharov** — testing wartet: 66de84a muss auf main von jaira-plugin-template (Alex pusht), dann test-Workflow dort gruen -> test-verdict. Lokal: go vet + go test gruen (exit 0).
