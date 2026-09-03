---
id: 01M0YT73TG5GVY6D0XRQ4DQPMS
title: lanes remove schreibt eine Lane-Datei ohne Anker
status: done
ready: true
creator: BeMuCa
goal: Nach jaira lanes remove warnt das Board nicht mehr ueber einen leeren after-Wert
context: "Gefunden am 26.08. beim Einrichten des Boards im jaira-Repo selbst. 'jaira lanes remove my-lane' materialisiert das Projekt-Lane-Verzeichnis und schreibt dabei alle eingebauten Lanes als Dateien heraus - backlog.md bekommt dabei kein 'after:'. Seither meldet jeder Befehl: 'lane .jaira/lanes/backlog.md: anchor \"\" is not installed; placed before the terminal lane'. Die Lane funktioniert, aber die Warnung steht bei jedem Aufruf da. Vermutlich schreibt der Export die eingebaute Lane heraus, ohne dass die eingebaute ueberhaupt ein after hat - dann muss der Export das Feld weglassen statt es leer zu schreiben."
definition-of-done: Nach lanes remove auf einem frischen Board kommt keine Anker-Warnung; eine schon geschriebene backlog.md ohne after loest ebenfalls keine Warnung aus
blocked-by: []
commits: []
created-at: 2026-08-26T10:35:24Z
updated-at: 2026-08-31T16:21:34Z
claimed-by: EE-3NX6GL3-2087828
claimed-at: 2026-08-26T15:48:43Z
updated-by: BeMuCa
assignee: BeMuCa
outcome-what: "Kritik umgesetzt: insertAt entfernt, stattdessen ein 'else if' vor der Ankerpruefung. Netto bleibt eine geaenderte Bedingung plus Kommentar in core/lane/order(); die Einfuegezeile ist unberuehrt."
outcome-why: "Ein Helfer, dessen beide Klemmungen von keiner Aufrufstelle erreichbar sind, ist Fehlerbehandlung fuer einen Zustand der nicht auftreten kann - und die obere Klemmung war schon vorher toter Schutz, den das Verschieben zu meinem gemacht haette."
outcome-resolves: "Die DoD ist unveraendert erfuellt: derselbe Test laeuft, gofmt unveraendert, ~/.local/bin/jaira neu gebaut, keine Warnung auf diesem Board. Der Diff ist jetzt kleiner als der, den die Kritik gelesen hat."
review-summary: "core/lane/order() hat zwei Faelle verwechselt: 'diese Lane hat absichtlich keinen Anker' und 'der Anker dieser Lane existiert nicht'. Beide endeten bei idx < 0. Jetzt steht 'if l.After == \"\"' als else-if vor der Ankerpruefung, setzt idx auf terminalIndex(out)-1 und laeuft in dieselbe Einfuegezeile. Netto: eine Bedingung, ein Kommentar, ein Test. Zwei Commits: 9b8b981 mit einem insertAt-Helfer, 0f1142e nimmt ihn nach der Kritik wieder weg."
review-gaps: "Nichts entfernt, und das ist geprueft statt behauptet. Doppelung: 'After == \"\"' und terminalIndex kommen im ganzen Repo nur in core/lane/lane.go:643-651 und :686 vor, es gibt keine zweite Stelle, die eine Lane ohne Anker platziert. Toter Code: insertAt ist in der Kritik-Runde schon weggefallen, nichts weiter verwaist. Stehen gelassen und benannt: die Klemmung 'if at > len(out)' bei der Einfuegezeile ist unerreichbar, war das aber schon vorher - fremder toter Code, nicht meiner. Auch stehen gelassen: gofmt -l listet weiterhin internal/cli/tickets.go, unberuehrt. Kosten: terminalIndex laeuft in der Schleife, ueber maximal ein Dutzend Lanes und nur beim Laden - nicht angefasst. Der neue Test doppelt TestOrderNoAnchorLandsBeforeTerminal nicht: der deckt 'kein Anker, Builtins vorhanden' ab, der neue 'kein Anker, Liste leer'."
review-check: |-
  1. Im jaira-Repo bauen: go build -o ~/.local/bin/jaira ./cmd/jaira - laeuft ohne Ausgabe durch.
  2. Im jaira-Repo 'jaira list' aufrufen. Die erste Zeile muss 'Backlog' sein. Keine Zeile darf mit 'jaira: warning:' anfangen. Vorher stand dort: 'lane .jaira/lanes/backlog.md: anchor "" is not installed'.
  3. Frisches Board pruefen: 'cd $(mktemp -d) && git init -q . && jaira init'. Danach 'jaira lanes remove brainstorm' - es antwortet 'removed brainstorm from this project'.
  4. Im selben Ordner 'jaira list' aufrufen. Es kommt nur 'No tickets match.' und keine Warnung.
  5. Gegenprobe, dass Schritt 4 ueberhaupt etwas prueft: dieselben Schritte mit einem Binary von 7a861b3 (git worktree add, go build) - dort erscheint die Anker-Warnung. Gemessen: sie erscheint.
  6. go clean -testcache && go test ./... -race - alle Pakete ok.
review-verdict: "Angenommen, mit einer Einschraenkung, die der Leser kennen soll: dieselbe Sitzung hat implementiert, kritisiert, optimiert und diesen Diff beurteilt - es war kein zweites Modell im Spiel, weil in dieser Sitzung keine Subagenten erlaubt waren. Was tatsaechlich geprueft ist: go test ./... -race mit geleertem Cache gruen, gofmt -l unveraendert, und der Test faellt auf dem alten Code. Die Gegenprobe aus review-check Schritt 5 wurde ausgefuehrt, nicht nur aufgeschrieben: ein Binary von 7a861b3 zeigt die Warnung, das neue nicht. Verhalten aendert sich nur an einer Stelle: eine Lane ohne Anker warnt nicht mehr. Ihre Position ist dieselbe wie vorher. Zwei Dinge bewusst nicht gemacht: die Umsortierung, die dabei gefunden wurde (QA3GN1), und der schon vorher unerreichbare Klemmungs-Zweig."
---

# lanes remove schreibt eine Lane-Datei ohne Anker

## Definition of Done

- [x] Nach lanes remove auf einem frischen Board kommt keine Anker-Warnung; eine schon geschriebene backlog.md ohne after loest ebenfalls keine Warnung aus
  proof: core/lane/lane.go:626 insertAt; TestMaterialisedWorkingSetLoadsWithoutAnchorWarning; live: jaira list auf diesem Board ohne Warnung

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-08-26 15:54 · BeMuCa** — Die Vermutung im Ticket war falsch: der Export laesst 'after:' korrekt weg, backlog.md hat das Feld gar nicht. Der Fehler sass im Laden, nicht im Schreiben.

core/lane/order() nahm an, dass die eingebauten Lanes schon platziert sind. Nach MaterialiseWorkingSet ist keine Lane mehr Builtin (replaceLane markiert eine Ersetzung absichtlich nicht als Builtin), also startet die Liste leer, terminalIndex einer leeren Liste ist 0, idx wird -1 und der Fall 'kein Anker' fiel in den Zweig 'Anker nicht installiert', weil present[""] false ist.

Beim Verifizieren fiel ein zweiter Fehler auf, jetzt QA3GN1: jaira init schreibt keine order-Datei, und ohne sie ordnet order() Lanes mit gleichem Anker um. Absichtlich NICHT hier mitgefixt - 'direkt hinter den Anker' ist gewollt (so landet critique vor human), die Reparatur gehoert in Materialise.
- **2026-08-26 15:56 · BeMuCa** — critique: insertAt ersetzen durch ein 'else if' vor der Ankerpruefung. Die beiden Klemmungen in insertAt sind unerreichbar (at liegt immer in [0, len(out)]), und ein Helfer, dessen Schutz nie greift, ist Fehlerbehandlung fuer einen Zustand der nicht auftreten kann. Die alte Einfuegezeile bleibt dann unangetastet - auch die vorher schon tote obere Klemmung, die ich sonst durch das Verschieben zu meiner mache.
