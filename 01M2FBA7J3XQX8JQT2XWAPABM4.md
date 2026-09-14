---
id: 01M2FBA7J3XQX8JQT2XWAPABM4
title: "Windows-Fallen fallen auf Linux auf, nicht erst acht Minuten spaeter in CI"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: "Alexander Sacharov"
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
updated-at: 2026-09-14T16:05:23Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-273044
claimed-at: 2026-09-14T15:48:01Z
---

# Windows-Fallen fallen auf Linux auf, nicht erst acht Minuten spaeter in CI

## Definition of Done

- [x] Ein Test, der auf Linux laeuft, faellt bei jedem der fuenf Muster um und nennt im Fehlertext die Abhilfe, nicht nur den Fund: ein t.Setenv("HOME") ohne ein USERPROFILE daneben; ein //go:embed-Verzeichnis ohne passende eol=lf-Zeile in .gitattributes; eine Berechtigungs- oder chmod-Behauptung ohne runtime.GOOS-Abzweig; ein erwarteter Binaername ohne .exe-Behandlung; ein zusammengebauter Pfad mit festem / statt filepath.Join. Jedes Muster hat einen eigenen Unterfall, der zeigt, dass er tatsaechlich anschlaegt.
  proof: internal/wintrap/wintrap_test.go:TestEachPatternFires (fixtures internal/wintrap/testdata/rule1..rule5)
- [x] Der Test ist auf dem heutigen Stand des Repositories gruen, ohne dass dafuer eine der fuenf Regeln aufgeweicht wurde.
  proof: internal/wintrap/wintrap_test.go:TestRepositoryIsClean
- [x] Der ubuntu-latest-Job in .github/workflows/ci.yaml fuehrt zusaetzlich 'GOOS=windows GOARCH=amd64 go vet ./...' und 'GOOS=windows GOARCH=amd64 go build ./cmd/jaira' aus, und ein absichtlich eingebauter Windows-Uebersetzungsfehler laesst diesen Schritt fallen.
  proof: .github/workflows/ci.yaml:29 'Cross-check the Windows build'
- [ ] Die beiden GOOS=windows-Zeilen stehen so in der Entwickler-Dokumentation, dass jemand sie vor dem Push von Hand laufen lassen kann.
- [ ] Keine Zeile in core/release/NOTES.md: Tests und CI sind von aussen am Binary nicht zu beobachten.

Der Test ist auf dem heutigen Stand des Repositories gruen, ohne dass dafuer eine der fuenf Regeln aufgeweicht wurde.

Der ubuntu-latest-Job in .github/workflows/ci.yaml fuehrt zusaetzlich 'GOOS=windows GOARCH=amd64 go vet ./...' und 'GOOS=windows GOARCH=amd64 go build ./cmd/jaira' aus, und ein absichtlich eingebauter Windows-Uebersetzungsfehler laesst diesen Schritt fallen.

Die beiden GOOS=windows-Zeilen stehen so in der Entwickler-Dokumentation, dass jemand sie vor dem Push von Hand laufen lassen kann.

Keine Zeile in core/release/NOTES.md: Tests und CI sind von aussen am Binary nicht zu beobachten.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] internal/wintrap anlegen: Finding-Typ mit Datei, Zeile und Abhilfetext, Scan(root) laeuft ueber den Baum und ueberspringt .git/ und testdata/
- [x] Testgeruest: ein Lauf gegen die Repository-Wurzel (per Aufstieg zur go.mod gefunden) erwartet null Findings, ein Lauf je Muster gegen ein testdata-Fixture erwartet genau dessen Fund
- [x] Muster 1 bauen: t.Setenv("HOME") ohne t.Setenv("USERPROFILE") in derselben Funktion, per AST, plus Fixture
- [x] Muster 2 bauen: //go:embed-Muster aufloesen und gegen einen minimalen .gitattributes-Matcher (*, ** ) auf eol=lf pruefen, plus Fixture
- [x] Muster 3 bauen: Chmod-Aufruf oder Mode().Perm()-Behauptung ohne runtime.GOOS in derselben Funktion, plus Fixture
- [x] Muster 4 bauen: erwarteter Binaername aus 'go build -o' oder exec.Command ohne .exe-Behandlung in derselben Funktion, plus Fixture
- [x] Muster 5 bauen: Verkettung mit festem / deren Ergebnis in einem os.*/filepath.*-Argument oder einem TrimPrefix gegen einen filepath-Wert landet, plus Fixture
- [x] Fehlertexte durchgehen: jedes Finding nennt die Abhilfe (USERPROFILE danebensetzen, Zeile in .gitattributes, runtime.GOOS-Abzweig, .exe-Suffix, filepath.Join), nicht nur den Fund
- [x] Repository-Lauf gruen machen, ohne eine Regel aufzuweichen: identity_test.go USERPROFILE, .gitattributes fuer core/release/NOTES.md und core/hook/example/notify.sh, role.go:90 und browse.go:318 auf filepath umstellen
- [x] ci.yaml: im ubuntu-latest-Job 'GOOS=windows GOARCH=amd64 go vet ./...' und 'GOOS=windows GOARCH=amd64 go build ./cmd/jaira' ergaenzen
- [x] Beweisen, dass der neue Schritt greift: einen Windows-Uebersetzungsfehler absichtlich einbauen, den Fehlschlag zeigen, Ergebnis als note festhalten, Fehler zuruecknehmen
- [x] README '## Development': die beiden GOOS=windows-Zeilen so aufnehmen, dass sie vor dem Push von Hand laufen
- [~] Abschluss: go test ./... und go vet ./... gruen, kein Eintrag in core/release/NOTES.md

## Progress
- **2026-09-14 15:51 · Alexander Sacharov** — Entwurfsentscheidung Ablage: die Pruefer kommen als eigenes Paket internal/wintrap, nicht als Test-Helfer in einem bestehenden Paket. Grund: die Definition of Done verlangt fuer jedes Muster einen Unterfall, der zeigt, dass der Pruefer anschlaegt. Dafuer muss der Pruefer eine beliebige Wurzel scannen koennen (ein synthetischer Baum unter testdata/), nicht fest das Repository. Also: Scan(root string) []Finding in einer normalen .go-Datei, der Test ruft es zweimal - einmal gegen testdata-Verstoesse (erwartet: schlaegt an), einmal gegen die Repository-Wurzel (erwartet: nichts). Die Repository-Wurzel wird von der Testdatei aus per Aufstieg zur go.mod gesucht. Der Scan muss .git/ und testdata/ ueberspringen, sonst schlagen die eigenen Fixtures im Repository-Lauf an. internal/ und nicht core/, weil das ein Entwicklungs-Waechter ist und core/ nichts aus internal/ importieren darf - kein Pfad zum Binary, also keine Auswirkung auf die Auslieferung.
- **2026-09-14 15:51 · Alexander Sacharov** — Muster 2 (go:embed ohne eol=lf) - zwei Wege geprueft. (a) 'git check-attr text eol -- <datei>' aufrufen: exakt die echte git-Semantik, braucht aber git und fuer jedes testdata-Fixture ein 'git init'. (b) Einen minimalen .gitattributes-Matcher selbst schreiben. Empfehlung: (b). Im Repository stehen heute genau drei Musterformen (dir/*.md, dir/**/*.md, dir/**/*.sh), die path.Match plus eine **-Aufloesung abdeckt, und (b) ist hermetisch - die Fixtures sind dann einfache Verzeichnisse ohne Repository. Wichtige Richtung bei Zweifel: im Zweifel melden, nicht durchwinken. Ein falscher Alarm kostet eine Zeile in .gitattributes, ein uebersehener Fall kostet wieder einen roten windows-latest-Job.
- **2026-09-14 15:51 · Alexander Sacharov** — Muster 5 (fester Vorwaertsschraegstrich) laesst sich nicht als grep bauen. Zahlen aus dem heutigen Stand: 12 Treffer fuer '+ "/' ausserhalb von Tests, 26 in Tests. Die allermeisten sind richtig so: URLs (selfupdate.go:208), git-Refnamen (gitref.go:467, settings.go:215), die .gitignore-Zeile (board/gitignore.go:10) und Vergleiche, die vorher bewusst durch filepath.ToSlash gehen (link.go:197, tickets.go:1332). Ein naiver Textfilter meldet alle und ist damit unbrauchbar. Empfehlung: per AST nur dann melden, wenn das Ergebnis der Verkettung in derselben Funktion in einem Argument einer os.*- oder filepath.*-Funktion landet, oder in strings.TrimPrefix/HasPrefix gegen einen Wert, der aus einem filepath.*-Aufruf stammt. Das faengt genau die beiden echten Faelle (core/role/role.go:90, internal/tui/browse.go:318) und laesst URLs und Refnamen in Ruhe.
- **2026-09-14 15:51 · Alexander Sacharov** — Was auf dem heutigen Stand tatsaechlich anschlagen wird - das ist die Arbeit hinter 'der Test ist gruen, ohne eine Regel aufzuweichen'. Muster 1: core/identity/identity_test.go:27 setzt t.Setenv("HOME", t.TempDir()) ohne USERPROFILE daneben (die beiden anderen Stellen, cache_test.go:26 und roles_test.go:130, sind bereits versorgt). Muster 2: core/release/NOTES.md und core/hook/example/notify.sh werden per go:embed eingebunden, stehen aber nicht in .gitattributes - NOTES.md wird zeilenweise gescannt, notify.sh ist mit CRLF nicht ausfuehrbar, beide brauchen die Zeile wirklich. Muster 5: core/role/role.go:90 und internal/tui/browse.go:318 schneiden mit strings.TrimPrefix(p, root+"/") ab, was auf Windows nie trifft. Muster 3 und 4 haben nach heutigem Stand vermutlich keinen echten Treffer mehr - das ist erst nach dem Bau des Pruefers sicher.
- **2026-09-14 16:00 · Alexander Sacharov** — Muster 5, Abweichung vom Plan: core/role/role.go:90 ist KEIN echter Fehler. Der Plan (Schritt 9) wollte ihn auf filepath umstellen - das waere falsch. Die Zeile laeuft ueber fs.WalkDir auf einer embed.FS, und io/fs-Pfade sind auf jeder Plattform mit / getrennt; filepath.Join wuerde sie auf Windows kaputt machen. Gleiches gilt fuer core/settings/settings.go:211, das ist ein git-Refname. Beide sind deshalb mit //wintrap:ok plus Begruendung markiert statt umgeschrieben. Die Ausnahme-Markierung ist bewusst so gebaut, dass die Regel weiter anschlaegt und die Ausnahme im Diff steht - eine aufgeweichte Regel waere unsichtbar. Echt und umgestellt sind nur internal/tui/browse.go:318 (Pfade aus einem Dateisystem-Scan) sowie core/tag/tag_test.go und core/gate/gate_test.go.
- **2026-09-14 16:02 · Alexander Sacharov** — Der neue CI-Schritt greift nachweislich. Probe: eine Datei internal/cli/probe_windowsbreak.go mit "syscall.Umask(0)" - existiert auf unix, nicht auf windows. Ergebnis: "go build ./..." auf Linux bleibt gruen, "GOOS=windows GOARCH=amd64 go vet ./..." bricht mit "internal/cli/probe_windowsbreak.go:6:40: undefined: syscall.Umask", ebenso "GOOS=windows GOARCH=amd64 go build ./cmd/jaira". Nach Entfernen der Probe wieder gruen. Die Probe ist bewusst nicht im Repository geblieben - der Beweis steht hier, nicht als dauerhaft kaputte Datei. go vet deckt dabei auch die _test.go-Dateien ab, und genau dort lagen die meisten bisherigen Windows-Fehler.
