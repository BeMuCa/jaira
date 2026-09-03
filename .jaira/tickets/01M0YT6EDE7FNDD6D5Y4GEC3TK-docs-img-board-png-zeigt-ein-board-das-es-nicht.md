---
id: 01M0YT6EDE7FNDD6D5Y4GEC3TK
title: "docs/img/board.png zeigt ein Board, das es nicht mehr gibt"
status: done
ready: true
creator: BeMuCa
goal: "Das Board-Bild in der Doku zeigt den heutigen Stand, und es gibt ein Skript dafuer"
context: "Das Bild zeigt eine Meldung 'N lane(s) off-screen', die es seit c02e49c nicht mehr gibt, und kein 'v compact'. Seither hat sich mehr geaendert: die fokussierte Lane sitzt mittig (3067c44), z versteckt leere Lanes (9f01288), Karten tragen neue Marken. Die anderen drei Bilder wurden geprueft und stimmen. Zum Neubauen fehlt eine Demo-Welt - die wurde nie geskriptet, das Rezept steht im Handoff bei 9cf71cc. Diesmal skripten, sonst steht dasselbe Problem beim naechsten Board-Umbau wieder da."
definition-of-done: "Das Bild zeigt die heutige Oberflaeche; ein eingechecktes Skript baut die Demo-Welt, aus der es entsteht"
blocked-by: []
commits: []
created-at: 2026-08-26T10:35:02Z
updated-at: 2026-09-03T10:28:51Z
claimed-by: EE-3NX6GL3-2856189
claimed-at: 2026-08-31T18:14:30Z
updated-by: BeMuCa
assignee: BeMuCa
outcome-what: "Zwei Skript-Funde behoben: der gedruckte Renderbefehl schreibt nach $dir/board.png plus separater cp-Zeile statt direkt nach docs/img/board.png, und mk() prueft create-Status und Handle selbst statt sich auf eine Pipeline unter set -eu zu verlassen"
outcome-why: "so verschmutzt kein Iterieren an der Demo-Welt mehr ein eingechecktes Asset, und ein abgelehntes create bricht mit eigener Meldung ab statt Exit 0 mit leerem Handle zu liefern und ein Kommando spaeter falsch zu sterben"
outcome-resolves: "sh -n sauber; beide Fehlerzweige einzeln getestet (jaira=false und jaira=true, je Exit 1 vor dem naechsten Kommando); voller Lauf in einem Temp-Verzeichnis gruen und docs/img/board.png unveraendert (md5 6e41babb92f6dc022341939fb446deff vor und nach dem Lauf)"
review-summary: none
review-gaps: "Nichts entfernt. Gelassen: die drei intakten Bilder (home/pipeline/signoff) bewusst nicht neu gerendert - Drift-Vermeidung steht im Skript-Header; 150x20 statt 150x27 (toter Leerraum unter den Karten weg); projects.json mit festen Timestamps (last_open ist sekundengenau, drei Boards in einem Lauf wuerden zufaellig sortieren)."
review-check: "./scripts/demoworld.sh ausfuehren (nutzt mktemp, beruehrt dein Board nicht; baut jaira aus dem Checkout), dann die zwei export/go-run-Zeilen aus der Skript-Ausgabe: docs/img/board.png entsteht neu und git diff zeigt keine Aenderung (deterministisch). Bild oeffnen: Tastenzeile mit 'z thin empty', Karten mit 'o spec' und 'unworked', keine 'N lane(s) off-screen'-Zeile."
review-verdict: "accept (Opus-End-Review auf 9cfcf58: PNG pixel-identisch reproduziert, Isolation vollstaendig verifiziert; Delta f873556 vom Koordinator gelesen: Rezept rendert neben der Welt + cp-Zeile, Fehler-Guards in mk() beidseitig getestet)"
---

# docs/img/board.png zeigt ein Board, das es nicht mehr gibt

## Definition of Done

- [x] Das Bild zeigt die heutige Oberflaeche; ein eingechecktes Skript baut die Demo-Welt, aus der es entsteht
  proof: docs/img/board.png neu gerendert aus scripts/demoworld.sh; Rezept im Skript-Header

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Rezept aus Handoff 9cf71cc bestaetigen: shotgen + scripts/termshot.py, PIL und DejaVu-Font sind vorhanden
  proof: scripts/termshot.py eingecheckt; PIL 12.2.0 und /usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf vorhanden; board.png 1571x630 = 150 cols x 27 Zeilen
- [x] Demo-Welt als scripts/demoworld.sh skripten: drei Boards (checkout-service, mobile-app, infra-tooling) in einem Temp-Verzeichnis, JAIRA_HOME/JAIRA_LANES_DIR/JAIRA_USER gescoped, Tickets ueber die CLI in Backlog/Todo/Implementing/Human Review/Blocked
  proof: scripts/demoworld.sh: drei Boards, 5 Tickets in Backlog/Todo/Implementing/Human Review/Blocked, Lauf in /tmp .../scratchpad/final erfolgreich
- [x] Skript laufen lassen, Board-Ansicht mit shotgen 150x28 rendern und mit termshot.py nach docs/img/board.png schreiben
  proof: go run ./scripts/shotgen <world>/checkout-service board 150 20 | python3 scripts/termshot.py /dev/stdin docs/img/board.png --cols 150 -> 1571x410, 17 Zeilen
- [x] Bild pruefen: kein 'lane(s) off-screen', neuer Footer, zentrierte Lane; die anderen drei Bilder nicht anfassen
  proof: Bild zeigt Footer 'v compact / z thin empty', Marken '○ spec' und '◇ unworked'; keine off-screen-Meldung; home/pipeline/signoff.png unberuehrt (git status)
- [x] gofmt / go test nur falls Go-Code beruehrt wurde; Commit docs(GEC3TK)
  proof: kein Go-Code angefasst; gofmt -l core internal => nur internal/cli/tickets.go (vorbestehend)

## Progress
- **2026-08-31 18:20 · BeMuCa** — Das Rezept stand komplett im Handoff 9cf71cc und war schon eingecheckt: scripts/shotgen (rendert durch die echten Models) plus scripts/termshot.py (ANSI nach PNG via Pillow). Kein freeze, kein Screenshot per Hand, nichts aus dem Netz - PIL 12.2.0 und DejaVuSansMono liegen hier. Fehlte nur die Demo-Welt.

Was die Bilder verraten haben: board.png ist 1571x630, das sind bei DejaVuSansMono 17pt genau 150 Spalten (cw=10.234375, width=int(cw*cols)+2*PAD). So kam die Spaltenzahl heraus, sie stand nirgends.

Die Demo-Welt braucht keine eigenen Lanes: die zehn Namen in den alten Bildern (Backlog..Blocked, 'Implementing', 'HITL', 'Human Review') sind genau die Built-ins, die 'jaira init' in .jaira/lanes/ schreibt. Assignee 'Demo' kommt aus JAIRA_USER, nicht aus git config.

Die Projektzeile oben wird nach last_open sortiert, und FormatTime hat Sekundenauflösung. Drei Boards in einem Skriptlauf landen alle in derselben Sekunde, also war die Reihenfolge zufällig. Deshalb schreibt das Skript projects.json am Ende selbst mit festen Zeitstempeln.

Handles sind generiert, nicht wählbar. Das Skript liest sie aus 'create --json' zurück; 'list --json | head -1' waere falsch, weil die Reihenfolge nicht die Anlagereihenfolge ist.

Das Backlog-Ticket hat absichtlich kein --goal. Nur so ist gate.Ready false und die Karte zeigt '○ spec'. Mit Goal verschwand die Marke, die im alten Bild noch zu sehen war.

Nur board.png neu. Der Arbeitsbaum hat parallel fremde Aenderungen (README.md, core/validate, internal/cli/*) von einer anderen Session; der Commit fasst nur scripts/demoworld.sh und docs/img/board.png an.
- **2026-08-31 18:44 · BeMuCa** — End-Review accept mit 2 Skript-Funden: das gedruckte Rezept schreibt direkt in das getrackte docs/img/board.png (Iterieren verschmutzt ein committetes Asset) - Ziel $dir/board.png plus separate cp-Zeile; mk() ist eine Pipeline unter set -eu ohne pipefail - leerer Handle stirbt eine Zeile spaeter irrefuehrend, [ -n ]-Guard nach jedem Capture.
- **2026-08-31 18:46 · BeMuCa** — Zwei Review-Funde aus dem Skript nachgezogen.

1) Der gedruckte Renderbefehl zeigte direkt auf docs/img/board.png. Wer die Welt zweimal baut, macht damit jedes Mal ein eingechecktes Asset dreckig. Jetzt rendert er nach $dir/board.png, und ein eigener 'cp'-Schritt kopiert es erst, wenn es passt.

2) mk() war eine Pipeline. Unter set -eu ist der Status der von sed, nicht der von jaira: ein abgelehntes create lieferte Exit 0 und einen leeren Handle, und das Skript starb erst beim naechsten Kommando mit einer Fehlermeldung, die auf das falsche Ding zeigte. sh hat kein pipefail, also wird die Ausgabe erst in eine Variable geholt und dann geprueft.

Abweichung vom Review-Vorschlag: der Check steht in mk() statt hinter jeder der sechs Handle-Zuweisungen. Gleiche Abdeckung, eine Stelle statt sechs, und der mobile-app-Aufruf ist mit abgedeckt, obwohl er gar keinen Handle einsammelt (er laeuft ohne Kommandosubstitution). Beide Zweige einzeln getestet: mit jaira=false greift der create-Zweig, mit jaira=true (Exit 0, keine Ausgabe) der Leer-Handle-Zweig, beide mit Exit 1 vor dem naechsten Kommando.

Das eingecheckte PNG wurde nicht neu gerendert: md5 vor und nach dem Lauf identisch (6e41babb92f6dc022341939fb446deff). Das nebenher gerenderte Bild unterscheidet sich nur in den ULIDs.
