---
id: 01M6NHNZV41ZW2MYHS3YE2HSK1
title: "panda-bobr/jaira-plugin-template: lauffaehige Vorlage mit Beispielbefehl und Release-Pipeline"
status: done
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
blocked-by: []
related: []
commits:
  - d06a7cf3e092f9dd3dea08ffb60f3ddca4991cc1
created-at: 2026-11-05T09:47:15Z
updated-at: 2026-11-05T21:29:51Z
updated-by: Alexander Sacharov
blocked-reason: ""
outcome-what: "panda-bobr/jaira-plugin-template@fc5d7e1: Go-Plugin jaira-hello (greet liest jaira list --json und schreibt jaira note; manifest --json aus eingebettetem jaira-plugin.json; Exit 3 bei zu altem JAIRA_VERSION), Tests inkl. Manifest-Regeln, goreleaser (6 Binaries, checksums.txt, jaira-plugin.json), Actions test und release-auf-Tag mit Tag==version-Pruefung, README mit drei Schritten bis Katalog-PR. Demo panda-bobr/jaira-hello-demo v0.1.0; Katalog-PR panda-bobr/jaira-releases#4."
outcome-why: "Plugin-Autoren sollen ohne eigene Release-Technik ein Release bekommen, das die Katalog-CI aus D44NVR ohne Aenderung annimmt (Standard R6X6FH, docs/PLUGINS.md)."
outcome-resolves: "DoD 1: Tag v0.1.0 im unveraenderten Demo-Repo erzeugte alle 6 Assets + checksums.txt (Run 42011801567). DoD 2: check-catalog auf PR #4 gruen (Run 42026722958)."
review-summary: "Im Repo panda-bobr/jaira-plugin-template (main = 7678fd2) liegt eine Go-Plugin-Vorlage: main.go mit jaira-hello (greet zaehlt per 'jaira list --json' die Tickets bzw. schreibt per 'jaira note --json' eine Notiz; manifest --json gibt das per go:embed eingebettete jaira-plugin.json aus; Exit 3 bei JAIRA_VERSION < minJairaVersion; Exit-Code von jaira wird durchgereicht), 8 Tests inkl. Manifest-Regeln, .goreleaser.yaml (6 Binaries jaira-<id>_<os>_<arch>, Id aus dem Manifest via PLUGIN_ID, jaira-plugin.json als Asset und seit 7678fd2 auch in checksums.txt), release.yml (Tag muss v+version sein, dann go test und goreleaser) und test.yml (push main + PR). README fuehrt in drei Schritten bis zum Katalog-PR. Bewiesen wurde das an panda-bobr/jaira-hello-demo (1 Commit, Stand fc5d7e1): Release v0.1.0 mit 6 Binaries, checksums.txt, jaira-plugin.json; Katalog-PR jaira-releases#4 aendert nur plugins.json, Check 'check' gruen. In diesem Repo traegt der Branch nur die Ticketdatei (d06a7cf3)."
round: 3
review-gaps: |-
  note: checksum.extra_files (7678fd2) ist durch keinen Release-Lauf bewiesen - die Demo steht auf fc5d7e1, ihre checksums.txt enthaelt jaira-plugin.json nicht; erster echter Beweis ist das naechste Release aus der Vorlage
  note: Katalog-PR jaira-releases#4 ist noch offen; mergen listet hello wirklich im Katalog, schliessen nicht - Alex' Wahl
  note: greet <ticket> ruft jaira note ohne --board; bei gesetztem JAIRA_BOARD sucht list auf dem Board, note verlaesst sich auf die Handle-Aufloesung
  note: Plan-Punkt 1 ist ungehakt, obwohl fc5d7e1 ihn erfuellt - nur Kosmetik, der DoD ist vollstaendig gehakt
  note: go test lokal hier nicht gelaufen (Auto-Mode verweigert), gestuetzt auf test-Run 42074852946 gruen fuer 7678fd2
test-verdict: "pass - template main 7678fd2: test-Workflow Run 42074852946 success (go vet + go test). DoD-Belege unveraendert: jaira-hello-demo v0.1.0 (Run 42011801567, 6 Binaries + checksums.txt, sha256sum -c OK), jaira-releases#4 check-catalog Run 42026722958 pass. checksum.extra_files fuer jaira-plugin.json ist per goreleaser-Quelle v2.0.0 belegt, aber noch durch keinen Release-Lauf - erstes Release aus 7678fd2 zeigt es."
claimed-by: DESKTOP-RFTCH11-9733
claimed-at: 2026-11-05T17:04:34Z
review-verdict: "Beide DoD-Punkte sind mit oeffentlichen Belegen erfuellt und von mir nachgeprueft (Release-Assets der Demo, 1 Commit, PR #4 Check gruen, nur plugins.json). Kein Blocker. Einzige echte Unsicherheit: der Fix 'jaira-plugin.json in checksums.txt' ist Konfig ohne Release-Lauf - die Schluesselnamen stimmen mit goreleaser v2 ueberein, bewiesen ist er aber nicht."
review-round: 3
review-check: "1. Oeffne https://github.com/panda-bobr/jaira-hello-demo/releases/tag/v0.1.0 - du siehst 6 Dateien jaira-hello_<os>_<arch> (zwei mit .exe), checksums.txt und jaira-plugin.json  2. Oeffne https://github.com/panda-bobr/jaira-hello-demo/commits - genau ein Commit, also nichts an der Vorlage geaendert  3. Oeffne https://github.com/panda-bobr/jaira-releases/pull/4 - Check 'check' gruen, geaendert nur plugins.json  4. Oeffne https://github.com/panda-bobr/jaira-plugin-template/actions - test-Lauf fuer 7678fd2 gruen  5. Oeffne https://github.com/panda-bobr/jaira-plugin-template/blob/main/.goreleaser.yaml - unter checksum: steht extra_files mit ./jaira-plugin.json"
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
- **2026-11-05 17:20 · Alexander Sacharov** — Alex (via Teamlead, 05.11.2026): Review-Notes vor Abnahme beheben: (1) jaira-plugin.json in checksums.txt via checksum.extra_files, (2) README:19 Aussage zur id korrigieren und Fundstellen auflisten, (3) blocked_by R6X6FH und blocked_reason leeren. jaira-hello-demo bleibt wie es ist.
- **2026-11-05 17:21 · Alexander Sacharov** — Review-Notes behoben in Template 7678fd2 (nicht gepusht): .goreleaser.yaml checksum.extra_files ./jaira-plugin.json - geprueft: goreleaser Checksum-Struct hat ExtraFiles schon in pkg/config/config.go@v2.0.0, Doku www/content/customization/package/checksum.md (OSS, nicht Pro); release.yml pinnt ~> v2, aktuell v2.20.0. NICHT durch einen Release-Lauf bewiesen: goreleaser lokal ist verweigert, jaira-hello-demo bleibt auf Wunsch unveraendert - erster echter Beweis ist das naechste Release aus der Vorlage. README: id-Fundstellen (README Intro, Katalog-Eintrag, Abschnitt jaira-hello; main.go Kopfkommentar) aufgelistet; Kopfkommentar .goreleaser.yaml angepasst. blocked_by/blocked_reason geleert.
- **2026-11-05 17:21 · Alexander Sacharov** — critique Runde 3 (Diff 66de84a..7678fd2): nichts. checksum.extra_files neben release.extra_files ist kein Doppel - das eine hasht, das andere laedt hoch.
- **2026-11-05 17:29 · Alexander Sacharov** — review Runde 3 (frischer Kontext): Template 7678fd2 + Demo + PR #4 nachgeprueft, kein Blocker. go test lokal vom Auto-Mode verweigert - auf CI-Run 42074852946 gestuetzt.
- **2026-11-05 21:29 · Alexander Sacharov** — Alex (05.11.2026): abgenommen. Offen ausserhalb des Tickets: PR jaira-releases#4 (hello) mergen oder schliessen ist Alex' Wahl; checksum.extra_files beweist erst das naechste Release aus der Vorlage.
