---
id: 01M2FBA7J3XQX8JQT2XWAPABM4
title: "Windows-Fallen fallen auf Linux auf, nicht erst acht Minuten spaeter in CI"
status: critique
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
commits:
  - 3f0c8bf38ea2f302ccdfe7ebc462df927636eff9
created-at: 2026-09-14T06:57:45Z
updated-at: 2026-09-14T16:16:47Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-382690
claimed-at: 2026-09-14T16:16:39Z
outcome-what: "Vier critique-Funde abgearbeitet. Regel 4 verengt: checkExe meldet nur noch, wenn der Aufruf exec.Command/CommandContext mit erstem Literal \"go\" und \"build\"/\"install\"/\"test\" plus \"-o\" ist (neues isGoBuildOutput in internal/wintrap/wintrap_scan_test.go) - die Gegenproben sort -o und tar -c -o melden nicht mehr. Der CallExpr-Zweig aus mentionsExe ist geloescht, nur das \".exe\"-Literal zaehlt noch; die zwei Stellen, die davon lebten (internal/cli/mergebranches_test.go:34 und :194, beide rufen exeSuffix()), tragen jetzt //wintrap:ok mit Begruendung. Der \"://\"-Zweig in sepConcat ist geloescht. wintrap.go und gitattributes.go sind zu wintrap_scan_test.go und gitattributes_scan_test.go umbenannt, das Modul traegt kein exportiertes Scan mehr ausserhalb des Tests."
outcome-why: "Jeder der drei Regelfunde war eine Meldung, die etwas anderes prueft als sie behauptet, oder ein Stummschalter, den man an der stummgeschalteten Stelle nicht sieht - beides macht den Waechter unglaubwuerdig, und ein Waechter, dem man nicht glaubt, wird abgeschaltet. Der vierte Fund nahm 640 Zeilen Entwickler-Werkzeug aus dem ausgelieferten Modul, ohne dass eine Zeile davon anders arbeitet."
outcome-resolves: "go test ./... und go vet ./... gruen, GOOS=windows go vet ./... und go build ./cmd/jaira gruen. TestEachPatternFires laeuft unveraendert weiter - alle fuenf Fixtures schlagen an, Regel 4 also trotz der Verengung nicht vakuum. TestRepositoryIsClean gruen ohne aufgeweichte Regel: die zwei neuen Stellen sind mit //wintrap:ok plus Grund ausgenommen, nicht durch eine Lockerung. Kein Eintrag in core/release/NOTES.md, von aussen am Binary ist nichts davon zu beobachten."
review-summary: "internal/wintrap/wintrap_scan_test.go:442 isGoBuildOutput accepts \"install\" and \"test\" beside \"build\", but the finding text at wintrap_scan_test.go:413 asserts that \"a binary is built with go build -o\". Verified: \"go install -o\" does not exist (flag provided but not defined: -o), so that arm can never fire, and \"go test -o\" builds a test binary the message does not describe. Keep only \"build\" in the switch at :442, so the rule again verifies exactly what it claims."
---

# Windows-Fallen fallen auf Linux auf, nicht erst acht Minuten spaeter in CI

## Definition of Done

- [x] Ein Test, der auf Linux laeuft, faellt bei jedem der fuenf Muster um und nennt im Fehlertext die Abhilfe, nicht nur den Fund: ein t.Setenv("HOME") ohne ein USERPROFILE daneben; ein //go:embed-Verzeichnis ohne passende eol=lf-Zeile in .gitattributes; eine Berechtigungs- oder chmod-Behauptung ohne runtime.GOOS-Abzweig; ein erwarteter Binaername ohne .exe-Behandlung; ein zusammengebauter Pfad mit festem / statt filepath.Join. Jedes Muster hat einen eigenen Unterfall, der zeigt, dass er tatsaechlich anschlaegt.
  proof: internal/wintrap/wintrap_test.go:TestEachPatternFires (fixtures internal/wintrap/testdata/rule1..rule5)
- [x] Der Test ist auf dem heutigen Stand des Repositories gruen, ohne dass dafuer eine der fuenf Regeln aufgeweicht wurde.
  proof: internal/wintrap/wintrap_test.go:TestRepositoryIsClean
- [x] Der ubuntu-latest-Job in .github/workflows/ci.yaml fuehrt zusaetzlich 'GOOS=windows GOARCH=amd64 go vet ./...' und 'GOOS=windows GOARCH=amd64 go build ./cmd/jaira' aus, und ein absichtlich eingebauter Windows-Uebersetzungsfehler laesst diesen Schritt fallen.
  proof: .github/workflows/ci.yaml:29 'Cross-check the Windows build'
- [x] Die beiden GOOS=windows-Zeilen stehen so in der Entwickler-Dokumentation, dass jemand sie vor dem Push von Hand laufen lassen kann.
  proof: README.md:796 '## Development'
- [x] Keine Zeile in core/release/NOTES.md: Tests und CI sind von aussen am Binary nicht zu beobachten.
  proof: core/release/NOTES.md unchanged (git status clean for that path)

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
- [x] Abschluss: go test ./... und go vet ./... gruen, kein Eintrag in core/release/NOTES.md

## Progress
- **2026-09-14 15:51 · Alexander Sacharov** — Entwurfsentscheidung Ablage: die Pruefer kommen als eigenes Paket internal/wintrap, nicht als Test-Helfer in einem bestehenden Paket. Grund: die Definition of Done verlangt fuer jedes Muster einen Unterfall, der zeigt, dass der Pruefer anschlaegt. Dafuer muss der Pruefer eine beliebige Wurzel scannen koennen (ein synthetischer Baum unter testdata/), nicht fest das Repository. Also: Scan(root string) []Finding in einer normalen .go-Datei, der Test ruft es zweimal - einmal gegen testdata-Verstoesse (erwartet: schlaegt an), einmal gegen die Repository-Wurzel (erwartet: nichts). Die Repository-Wurzel wird von der Testdatei aus per Aufstieg zur go.mod gesucht. Der Scan muss .git/ und testdata/ ueberspringen, sonst schlagen die eigenen Fixtures im Repository-Lauf an. internal/ und nicht core/, weil das ein Entwicklungs-Waechter ist und core/ nichts aus internal/ importieren darf - kein Pfad zum Binary, also keine Auswirkung auf die Auslieferung.
- **2026-09-14 15:51 · Alexander Sacharov** — Muster 2 (go:embed ohne eol=lf) - zwei Wege geprueft. (a) 'git check-attr text eol -- <datei>' aufrufen: exakt die echte git-Semantik, braucht aber git und fuer jedes testdata-Fixture ein 'git init'. (b) Einen minimalen .gitattributes-Matcher selbst schreiben. Empfehlung: (b). Im Repository stehen heute genau drei Musterformen (dir/*.md, dir/**/*.md, dir/**/*.sh), die path.Match plus eine **-Aufloesung abdeckt, und (b) ist hermetisch - die Fixtures sind dann einfache Verzeichnisse ohne Repository. Wichtige Richtung bei Zweifel: im Zweifel melden, nicht durchwinken. Ein falscher Alarm kostet eine Zeile in .gitattributes, ein uebersehener Fall kostet wieder einen roten windows-latest-Job.
- **2026-09-14 15:51 · Alexander Sacharov** — Muster 5 (fester Vorwaertsschraegstrich) laesst sich nicht als grep bauen. Zahlen aus dem heutigen Stand: 12 Treffer fuer '+ "/' ausserhalb von Tests, 26 in Tests. Die allermeisten sind richtig so: URLs (selfupdate.go:208), git-Refnamen (gitref.go:467, settings.go:215), die .gitignore-Zeile (board/gitignore.go:10) und Vergleiche, die vorher bewusst durch filepath.ToSlash gehen (link.go:197, tickets.go:1332). Ein naiver Textfilter meldet alle und ist damit unbrauchbar. Empfehlung: per AST nur dann melden, wenn das Ergebnis der Verkettung in derselben Funktion in einem Argument einer os.*- oder filepath.*-Funktion landet, oder in strings.TrimPrefix/HasPrefix gegen einen Wert, der aus einem filepath.*-Aufruf stammt. Das faengt genau die beiden echten Faelle (core/role/role.go:90, internal/tui/browse.go:318) und laesst URLs und Refnamen in Ruhe.
- **2026-09-14 15:51 · Alexander Sacharov** — Was auf dem heutigen Stand tatsaechlich anschlagen wird - das ist die Arbeit hinter 'der Test ist gruen, ohne eine Regel aufzuweichen'. Muster 1: core/identity/identity_test.go:27 setzt t.Setenv("HOME", t.TempDir()) ohne USERPROFILE daneben (die beiden anderen Stellen, cache_test.go:26 und roles_test.go:130, sind bereits versorgt). Muster 2: core/release/NOTES.md und core/hook/example/notify.sh werden per go:embed eingebunden, stehen aber nicht in .gitattributes - NOTES.md wird zeilenweise gescannt, notify.sh ist mit CRLF nicht ausfuehrbar, beide brauchen die Zeile wirklich. Muster 5: core/role/role.go:90 und internal/tui/browse.go:318 schneiden mit strings.TrimPrefix(p, root+"/") ab, was auf Windows nie trifft. Muster 3 und 4 haben nach heutigem Stand vermutlich keinen echten Treffer mehr - das ist erst nach dem Bau des Pruefers sicher.
- **2026-09-14 16:00 · Alexander Sacharov** — Muster 5, Abweichung vom Plan: core/role/role.go:90 ist KEIN echter Fehler. Der Plan (Schritt 9) wollte ihn auf filepath umstellen - das waere falsch. Die Zeile laeuft ueber fs.WalkDir auf einer embed.FS, und io/fs-Pfade sind auf jeder Plattform mit / getrennt; filepath.Join wuerde sie auf Windows kaputt machen. Gleiches gilt fuer core/settings/settings.go:211, das ist ein git-Refname. Beide sind deshalb mit //wintrap:ok plus Begruendung markiert statt umgeschrieben. Die Ausnahme-Markierung ist bewusst so gebaut, dass die Regel weiter anschlaegt und die Ausnahme im Diff steht - eine aufgeweichte Regel waere unsichtbar. Echt und umgestellt sind nur internal/tui/browse.go:318 (Pfade aus einem Dateisystem-Scan) sowie core/tag/tag_test.go und core/gate/gate_test.go.
- **2026-09-14 16:02 · Alexander Sacharov** — Der neue CI-Schritt greift nachweislich. Probe: eine Datei internal/cli/probe_windowsbreak.go mit "syscall.Umask(0)" - existiert auf unix, nicht auf windows. Ergebnis: "go build ./..." auf Linux bleibt gruen, "GOOS=windows GOARCH=amd64 go vet ./..." bricht mit "internal/cli/probe_windowsbreak.go:6:40: undefined: syscall.Umask", ebenso "GOOS=windows GOARCH=amd64 go build ./cmd/jaira". Nach Entfernen der Probe wieder gruen. Die Probe ist bewusst nicht im Repository geblieben - der Beweis steht hier, nicht als dauerhaft kaputte Datei. go vet deckt dabei auch die _test.go-Dateien ab, und genau dort lagen die meisten bisherigen Windows-Fehler.
- **2026-09-14 16:09 · Alexander Sacharov** — critique: vier Funde, alle mit klarer Abhilfe, kein Punkt fuer die Nutzerin zu entscheiden.

1. wintrap.go:406 - Regel 4 trifft jeden Aufruf mit einem Argument "-o", behauptet im Fundtext aber "go build -o". Nachgemessen mit einem eigenen Fixture: exec.Command("sort","-o",out,in) und exec.Command("tar","-c","-o","x.tar",dir) melden beide Falle 4. Abhilfe: nur melden, wenn der Aufruf exec.Command/exec.CommandContext ist, erstes Literal "go" und "build" unter den Argumenten. Sonst prueft die Regel nicht, was ihre eigene Meldung behauptet.

2. wintrap.go:517 - der "://"-Zweig in sepConcat schuetzt einen Zustand, der nicht eintreten kann: eine URL, die in os.*/filepath.* oder in TrimPrefix gegen einen Pfad landet. Die eine echte URL-Verkettung (core/selfupdate/selfupdate.go:208) ist ein nacktes return, das checkSlash nie ansieht. Loeschen - der Abhilfetext in wintrap.go:468 schickt URL-Stellen schon zu //wintrap:ok, und das ist der Mechanismus, der wirkt.

3. wintrap.go:435 - mentionsExe schweigt fuer eine ganze Funktion, sobald irgendein gerufener Name "exe" enthaelt. Das ist ein dritter Stummschalter neben //wintrap:ok und dem hasGOOS-Funktionssprung, und der einzige, den man an der stummgeschalteten Stelle nicht sieht. Den CallExpr-Zweig raus, das ".exe"-Literal behalten, Hilfsfunktionen per //wintrap:ok mit Begruendung ausnehmen.

4. Einordnung: wintrap.go und gitattributes.go sind 640 Zeilen Nicht-Test-Code mit exportiertem Scan, dessen einziger Aufrufer der eigene Test ist. Das Repository hat den Waechtertest schon als Muster - core/lane/lane_test.go:698 TestModelTierNeverComparedToModelName liest den Quellbaum aus einer _test.go im bewachten Paket. Beide Dateien in _test.go umbenennen: testdata, Scan und alle vier Tests laufen unveraendert weiter, und das Modul traegt keine exportierte Entwickler-API mehr, die niemand importiert. Die note vom 15:51 begruendet das eigene Paket mit "Scan muss eine beliebige Wurzel scannen koennen" - das kann eine _test.go genauso.

Nicht gemeldet, absichtlich: gitattributes.go:30 liest nur root/.gitattributes, nicht die in Unterverzeichnissen. Im Repository gibt es heute keine, und der Fehlalarm kostet eine Zeile in der Wurzeldatei. Nicht wert, dafuer git-Semantik nachzubauen.
- **2026-09-14 16:13 · Alexander Sacharov** — critique-Runde abgearbeitet, vier Funde. Was dabei herauskam und nicht im Diff steht:

- Regel 4 auf exec.Command/CommandContext mit erstem Literal "go" plus "build"/"install"/"test" und "-o" verengt. Nachgemessen mit den beiden Gegenproben aus der critique (sort -o, tar -c -o) in testdata/clean: melden nicht mehr.
- Der CallExpr-Zweig in mentionsExe ist raus. Folge, die vorher niemand sah: internal/cli/mergebranches_test.go:34 und :194 waren nur deshalb still, weil sie exeSuffix() rufen. Beide tragen jetzt //wintrap:ok mit Begruendung - die Stelle ist korrekt, aber die Stummschaltung steht jetzt dort, wo man sie liest.
- Der "://"-Zweig in sepConcat ist geloescht. Kein neuer Treffer im Repository-Lauf, damit ist die Analyse der critique bestaetigt: eine URL erreicht os.*/filepath.*/TrimPrefix hier nirgends.
- wintrap.go -> wintrap_scan_test.go, gitattributes.go -> gitattributes_scan_test.go. Die note vom 15:51 hat das eigene Nicht-Test-Paket mit "Scan muss eine beliebige Wurzel scannen koennen" begruendet - das war kein Argument, eine _test.go kann das genauso. Damit traegt das Modul kein exportiertes Scan mehr, das niemand importiert. README ## Development bleibt woertlich richtig, dort steht "internal/wintrap, das go test ./... schon laeuft" und nicht der Dateiname.
