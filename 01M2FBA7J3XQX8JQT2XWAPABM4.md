---
id: 01M2FBA7J3XQX8JQT2XWAPABM4
title: "Windows-Fallen fallen auf Linux auf, nicht erst acht Minuten spaeter in CI"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die fuenf bekannten Windows-Fallen scheitern auf einem Linux-Rechner in Sekunden: ein Go-Test benennt jede mit der Abhilfe im Fehlertext, und ein GOOS=windows-Lauf von vet und build auf dem ubuntu-Job faengt ab, was schon beim Uebersetzen bricht."
context: |-
  Der windows-latest-Job faellt regelmaessig um, und wir erfahren es erst, wenn die Arbeit schon fertiggemeldet ist.

  Der Job braucht 8 Minuten. ubuntu-latest braucht 2:43, macos-latest 3:47. internal/tui allein 269 Sekunden davon.

  Sieben Windows-Fixes stehen in der Historie, und sie sind fuenf Muster, kein einziges neues Problem:
  - t.Setenv("HOME") ohne USERPROFILE: os.UserHomeDir liest auf Windows USERPROFILE, der Test lief am temp-Home vorbei (f085553, internal/cli/roles_test.go).
  - Ein go:embed-Verzeichnis ohne eol=lf-Zeile in .gitattributes: ein Windows-Checkout mit core.autocrlf gibt go:embed CRLF-Bytes (1b3bf08 fuer core/role/builtin, 40dca0d fuer core/lane/builtin).
  - Annahmen ueber das Execute-Bit und chmod, die es auf Windows nicht gibt (1b3bf08, 17912ba).
  - Ein Binaername ohne .exe (36f7055).
  - Pfade mit festem Vorwaertsschraegstrich (4d3cace).

  Keine dieser fuenf Ursachen braucht Windows, um gefunden zu werden. Alle fuenf sind auf Linux mit einem Griff in den Quelltext sichtbar.

  Wine ist geprueft und verworfen, das muss niemand nochmal durchrechnen: 'GOOS=windows GOARCH=amd64 go test -exec=wine ./...' laeuft grundsaetzlich, aber -race faellt weg (der Detektor braucht CGO und einen mingw-Kreuztoolchain), 9 der 30 Pakete rufen git und brauchen dafuer Git-for-Windows im wine-Praefix, und die CRLF-Klasse ist unter wine gar nicht sichtbar, weil sie eine Eigenschaft des Windows-Checkouts ist und nicht der Laufzeit. Das sind drei der sieben Faelle fuer einen halben Tag Aufbau und dauerhafte Pflege.

  Windows bleibt eine ausgelieferte Plattform: .goreleaser.yaml baut sechs Binaries, eines davon fuer Windows. Den Job abschalten oder auf continue-on-error setzen ist deshalb keine Loesung, das waere blind ausliefern.

  Es gibt in diesem Repository kein Taskfile.yml und kein Makefile. Der GOOS=windows-Lauf ist also ein Schritt in .github/workflows/ci.yaml plus eine dokumentierte Zeile, kein task-Target.
definition-of-done: |-
  Ein Test, der auf Linux laeuft, faellt bei jedem der fuenf Muster um und nennt im Fehlertext die Abhilfe, nicht nur den Fund: t.Setenv("HOME") ohne USERPROFILE danebenl; ein //go:embed-Verzeichnis ohne passende eol=lf-Zeile in .gitattributes; eine Berechtigungs- oder chmod-Behauptung ohne runtime.GOOS-Abzweig; ein erwarteter Binaername ohne .exe-Behandlung; ein zusammengebauter Pfad mit festem / statt filepath.Join. Jedes Muster hat einen eigenen Unterfall, der zeigt, dass er tatsaechlich anschlaegt.

  Der Test ist auf dem heutigen Stand des Repositories gruen, ohne dass dafuer eine der fuenf Regeln aufgeweicht wurde.

  Der ubuntu-latest-Job in .github/workflows/ci.yaml fuehrt zusaetzlich 'GOOS=windows GOARCH=amd64 go vet ./...' und 'GOOS=windows GOARCH=amd64 go build ./cmd/jaira' aus, und ein absichtlich eingebauter Windows-Uebersetzungsfehler laesst diesen Schritt fallen.

  Die beiden GOOS=windows-Zeilen stehen so in der Entwickler-Dokumentation, dass jemand sie vor dem Push von Hand laufen lassen kann.

  Keine Zeile in core/release/NOTES.md: Tests und CI sind von aussen am Binary nicht zu beobachten.
tags:
  - ci
blocked-by: []
related: []
commits: []
created-at: 2026-09-14T06:57:45Z
updated-at: 2026-09-14T06:58:00Z
updated-by: Alexander Sacharov
---

# Windows-Fallen fallen auf Linux auf, nicht erst acht Minuten spaeter in CI

## Definition of Done

- [ ] Ein Test, der auf Linux laeuft, faellt bei jedem der fuenf Muster um und nennt im Fehlertext die Abhilfe, nicht nur den Fund: t.Setenv("HOME") ohne USERPROFILE danebenl; ein //go:embed-Verzeichnis ohne passende eol=lf-Zeile in .gitattributes; eine Berechtigungs- oder chmod-Behauptung ohne runtime.GOOS-Abzweig; ein erwarteter Binaername ohne .exe-Behandlung; ein zusammengebauter Pfad mit festem / statt filepath.Join. Jedes Muster hat einen eigenen Unterfall, der zeigt, dass er tatsaechlich anschlaegt.

Der Test ist auf dem heutigen Stand des Repositories gruen, ohne dass dafuer eine der fuenf Regeln aufgeweicht wurde.

Der ubuntu-latest-Job in .github/workflows/ci.yaml fuehrt zusaetzlich 'GOOS=windows GOARCH=amd64 go vet ./...' und 'GOOS=windows GOARCH=amd64 go build ./cmd/jaira' aus, und ein absichtlich eingebauter Windows-Uebersetzungsfehler laesst diesen Schritt fallen.

Die beiden GOOS=windows-Zeilen stehen so in der Entwickler-Dokumentation, dass jemand sie vor dem Push von Hand laufen lassen kann.

Keine Zeile in core/release/NOTES.md: Tests und CI sind von aussen am Binary nicht zu beobachten.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

