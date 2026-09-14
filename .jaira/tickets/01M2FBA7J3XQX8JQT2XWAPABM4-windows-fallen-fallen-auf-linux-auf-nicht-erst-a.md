---
id: 01M2FBA7J3XQX8JQT2XWAPABM4
title: "Windows-Fallen fallen auf Linux auf, nicht erst acht Minuten spaeter in CI"
status: optimize
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
updated-at: 2026-09-14T18:14:16Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-463488
claimed-at: 2026-09-14T18:14:12Z
outcome-what: "Regel 4 auf 'go build -o' verengt: isGoBuildOutput (internal/wintrap/wintrap_scan_test.go:443) akzeptiert nur noch das Literal \"build\", die Zweige \"install\" und \"test\" sind raus."
outcome-why: "Der Fundtext behauptet 'a binary is built with go build -o'. 'go install' kennt kein -o, der Zweig konnte nie feuern; 'go test -o' baut ein Test-Binary, das der Text nicht beschreibt."
outcome-resolves: "go test ./... gruen, go vet ./... und GOOS=windows GOARCH=amd64 go vet ./... gruen. TestEachPatternFires und TestRepositoryIsClean gruen ohne aufgeweichte Regel."
review-summary: |-
  internal/wintrap/wintrap_scan_test.go:235 binaryExt is a 14-entry extension list guarding a case this module does not have: every //go:embed in it (core/role/role.go:33, core/lane/lane.go:34, core/release/release.go:17, core/hook/example.go:9) pulls only .md and .sh. Git already answers this with -text/binary, which gitattributes_scan_test.go:55 parses and then scores as NOT pinned. Delete binaryExt and let -text/binary count as covered in eolLF — one mechanism, and the one git owns.
  internal/wintrap/wintrap_scan_test.go:185 selName2 is a seven-line adapter with exactly one caller (hasGOOS, line 177), existing only to widen selName's ast.Expr parameter to ast.Node; the numbered name is what you call a function you did not want to name. Inline it: in hasGOOS type-assert n.(*ast.SelectorExpr) and test X == runtime / Sel == GOOS, then delete selName2.
  internal/wintrap/wintrap_scan_test.go:424 the isGoBuildOutput doc says "Anything narrower than this and the check fires on every unrelated tool that happens to take a -o flag". Narrower fires on less; it is looser that fires on sort -o and tar -o — the round-1 finding this function was written to answer. Write "looser".
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
- **2026-09-14 16:17 · Alexander Sacharov** — critique Runde 2: ein Fund, eine Zeile.

wintrap_scan_test.go:442 - isGoBuildOutput nimmt neben "build" auch "install" und "test", der Fundtext auf :413 behauptet aber "a binary is built with go build -o". Nachgemessen: "go install -o" gibt es nicht (flag provided but not defined: -o), der Zweig kann also nie greifen; "go test -o" baut ein Testbinary, das die Meldung nicht beschreibt. Das ist der Fund aus Runde 1 in kleiner Form - die Regel prueft wieder etwas anderes, als sie meldet. Abhilfe: im switch auf :442 nur "build" stehen lassen.

Geprueft und absichtlich NICHT gemeldet:
- binaryExt (:235) sieht nach totem Code aus, verhindert aber echten Schaden: ohne die Liste bekaeme ein eingebettetes .png den Rat, eol=lf zu setzen. Bleibt.
- Das Paket besteht jetzt nur noch aus _test.go-Dateien; go build ./..., go vet ./... und GOOS=windows go vet ./... laufen damit alle gruen. Die Umbenennung aus Runde 1 hat nichts kaputt gemacht.
- mentionsExe schaltet weiter eine ganze Funktion stumm, sobald ein ".exe"-Literal darin steht. Das war in Runde 1 die bewusste Entscheidung und wird hier nicht neu aufgemacht.
- README wiederholt die fuenf Muster aus dem Paketkommentar. Da der Paketkommentar jetzt in einer _test.go steht und go doc ihn nicht mehr zeigt, ist die README-Kopie die einzige auffindbare - kein Fund.
- **2026-09-14 16:18 · Alexander Sacharov** — critique Runde 3: ein Fund, eine Zeile. isGoBuildOutput nimmt nur noch "build", nicht mehr "install"/"test". Nachgeprueft, warum der Fund stimmt: 'go install' kennt kein -o (flag provided but not defined: -o), der Zweig konnte also nie feuern; 'go test -o' baut ein Test-Binary, und das beschreibt der Fundtext nicht. TestEachPatternFires bleibt gruen - das Fixture testdata/rule4/build.go:10 ist exec.Command("go","build","-o",...), also weiter im verengten Muster.
- **2026-09-14 16:21 · Alexander Sacharov** — critique Runde 3: drei Funde, alle klein und alle mit klarer Abhilfe, kein Punkt fuer die Nutzerin zu entscheiden. Runde-1- und Runde-2-Funde sind abgearbeitet und werden nicht wiederholt.

1. wintrap_scan_test.go:235 - binaryExt zaehlt 14 Endungen auf, fuer einen Fall, den dieses Modul nicht hat: nachgesehen, alle vier //go:embed (core/role/role.go:33 all:builtin, core/lane/lane.go:34, core/release/release.go:17, core/hook/example.go:9) ziehen ausschliesslich .md und .sh; unter core/role/builtin und core/lane/builtin liegen nur diese beiden Endungen. Git hat fuer eine Binaerdatei schon die Antwort - -text bzw. binary -, und gitattributes_scan_test.go:55 liest das Attribut bereits, wertet es dann aber ausdruecklich als 'nicht gepinnt'. binaryExt loeschen und -text/binary in eolLF als abgedeckt zaehlen: ein Mechanismus statt zweier, und der, den git selbst besitzt.

2. wintrap_scan_test.go:185 - selName2 ist ein siebenzeiliger Adapter mit genau einem Aufrufer (hasGOOS, Zeile 177) und existiert nur, um selName von ast.Expr auf ast.Node zu weiten. Der durchnummerierte Name ist das, was man einer Funktion gibt, die man nicht benennen wollte. In hasGOOS direkt n.(*ast.SelectorExpr) pruefen und X == runtime / Sel == GOOS testen, selName2 raus.

3. wintrap_scan_test.go:424 - der Doc-Kommentar zu isGoBuildOutput sagt 'Anything narrower than this and the check fires on every unrelated tool that happens to take a -o flag'. Enger trifft weniger; es ist weiter, was auf sort -o und tar -o anschlaegt - genau der Fund aus Runde 1, fuer den diese Funktion geschrieben wurde. 'looser' hinschreiben.

Nicht gemeldet, absichtlich: (a) fsPackages enthaelt 'ioutil', das im ganzen Repository nirgends vorkommt, und skipDirs kennt 'vendor' und 'node_modules', die es in diesem reinen Go-Modul nicht gibt - drei Map-Eintraege, deren Entfernung nichts aendert und deren Verbleib nichts kostet. (b) TestMatchPattern steht in wintrap_test.go, waehrend matchPattern in gitattributes_scan_test.go liegt; beide sind _test.go-Dateien im selben Paket, das laeuft, und eine Runde dafuer lohnt nicht. (c) core/role/role.go:93 und core/settings/settings.go:213 tragen den Grund in einem Kommentar ueber dem //wintrap:ok statt auf der Markerzeile wie internal/cli/mergebranches_test.go:34 - der Grund steht direkt daneben und ist lesbar.
- **2026-09-14 16:23 · Alexander Sacharov** — Dispatcher stopped here under the three-round rule: critique has sent the work back three times (round 1: rule 4 too wide plus the :// escape hatch; round 2: install/test arms of rule 4; round 3: binaryExt duplicates git's own -text/binary, and selName2 is a one-caller adapter). Findings shrink every round and none were re-raised, but the rule does not allow a fourth round on that reasoning. Round 3's two findings are dead-code removals - that is the optimize lane's job, not critique's, which may be the real reason the loop will not terminate. A person decides: run optimize next and let it take them, or accept and move on. Branch feat/APABM4-windows-traps is pushed to origin; no pull request opened, by instruction. go test ./... green, GOOS=windows vet and build green.
