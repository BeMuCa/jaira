---
id: 01M822MVB17D08MDDTKKBWT53T
title: Ticket-Dateinamen schneiden Umlaute nicht mehr mitten im Zeichen ab
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Slug kürzt auf 48 Bytes nur an einer Zeichengrenze, damit jeder Dateiname gültiges UTF-8 ist und auf macOS ausgecheckt werden kann."
context: |-
  Was falsch ist: core/ticket/id.go Slug() kürzt mit s[:48] auf Bytes. Fällt Byte 48 mitten in ein Mehrbyte-Zeichen (ü, ö, ß, Kyrillisch), endet der Dateiname mit einem halben Zeichen.
  Ausloeser 22.11.2026: PR #184 (AJ41SC) - CI macos-latest scheitert schon beim Checkout: 'unable to create file .jaira/logbook/.../...-ampel-gr\303.md'. APFS lehnt ungültiges UTF-8 ab; Linux und Windows-CI merkten nichts.
  Behelf: Datei im Branch per git mv umbenannt. Betroffen ist jede Ticket- und Logbook-Datei, deren Titel um Byte 48 ein Umlaut hat.
  Bekannt: master hat heute keinen ungültigen Dateinamen (geprüft mit git ls-tree).
definition-of-done: "Slug kürzt auf höchstens 48 Bytes an einer Zeichengrenze; Tabellentest mit Umlaut, ß, Kyrillisch und Emoji genau an der Grenze, jedes Ergebnis utf8.ValidString"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-11-22T16:50:19Z
updated-at: 2026-11-22T17:02:44Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1768841
claimed-at: 2026-11-22T16:50:44Z
outcome-what: Slug schneidet an einer Zeichengrenze; jaira validate meldet Dateien unter .jaira/ mit ungültigem UTF-8-Namen samt mv (exit 3)
outcome-why: "Ein Schnitt mitten in ü/ß/Kyrillisch erzeugte Dateinamen, die macOS nicht auschecken kann (PR #184)"
outcome-resolves: Ticket- und Logbook-Dateien sind immer gültiges UTF-8; Altlasten meldet validate
review-summary: "none"
round: 3
back-from: testing
back-by: lane
back-count: 1
review-gaps: "none: kein vorhandener Shell-Quoter (grep shellQuote/shQuote); truncate in internal/cli/tickets.go:1782 und internal/tui/view.go:1930 kürzen nach Runen mit Ellipse, nicht nach Bytes - nicht wiederverwendbar für das 48-Byte-Limit; kein toter Code, update.go ist gleich master"
test-verdict: "fail: suite green (go test -race core/ticket internal/cli core/validate RC=0, go test ./... RC=0), DoD 1-3 verified in the tree; scratch board: create with umlaut at byte 48 gives valid names, planted \\xc3 name -> validate exit 3 with git mv; unreadable dir no crash; validate writes nothing, no concurrency case. Finding: invalid_name Path and git mv are relative to the repo root while every other Problem path is absolute (internal/cli/validate.go invalidNames) - the printed git mv fails when run from a subdirectory; use the absolute path"
hunt-verdict: "change: internal/cli/validate.go invalidNames - printed git mv only works from the repo root and only for a tracked file (unshared board / untracked file: exit 128)"
---

# Ticket-Dateinamen schneiden Umlaute nicht mehr mitten im Zeichen ab

## Definition of Done

- [x] Slug kürzt auf höchstens 48 Bytes an einer Zeichengrenze; Tabellentest mit Umlaut, ß, Kyrillisch und Emoji genau an der Grenze, jedes Ergebnis utf8.ValidString
  proof: core/ticket/id.go:59 (Schnitt an utf8.RuneStart); core/ticket/id_test.go TestSlugCutsAtRuneBoundary + TestSlugAlwaysValidUTF8
- [x] Vorhandene Dateien mit ungültigem Namen: jaira doctor oder update meldet sie (oder benennt sie um) - Entscheidung im Plan, Test
  proof: internal/cli/validate.go invalidNames (Problem invalid_name, exit 3, mv mit absoluten Pfaden); internal/cli/validatenames_test.go TestValidateNamesFilesWithInvalidUTF8 (führt das gedruckte mv aus einem Unterordner aus)
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/BWT53T.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 16:50 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 16:51 · Alexander Sacharov** — decision DoD2: jaira update meldet ungültige UTF-8-Dateinamen unter .jaira/ (tickets, logbook, alles) mit git-mv-Hinweis, benennt nicht um. Grund: Umbenennen in update erzeugt Konflikte mit anderen Branches, die dieselbe Datei anfassen; master hat heute keine solche Datei; Dateiname ist nur Hinweis (Slug-Kommentar), der Frontmatter-Titel zählt. Kein doctor-Befehl im Code - kein neuer Befehl dafür.
- **2026-11-22 16:56 · Alexander Sacharov** — Klasse geprüft: einziger Byte-Schnitt in einem Dateinamen ist ticket.Slug. identity.Slug (Logbook-Ordner) und session.sanitize sind nur ASCII. Emoji ist kein Buchstabe und wird schon zum Bindestrich. Ungültiges UTF-8 im Titel wird von range zu U+FFFD -> Bindestrich, Slug bleibt gültig. Die Ausgabe von update quotet mit bash $'...', weil "..." die \\xNN-Escapes nicht liest; PowerShell-Nutzer müssen umschreiben.
- **2026-11-22 16:56 · Alexander Sacharov** — critique: Namensprüfung von update nach jaira validate verlegen (validate.Problem code invalid_name, error, Path, Message mit git mv $'...'). validate ist der Doktor-Befehl, den das DoD meint; exit 3 und JSON-Form sind dort schon da. update bleibt unberührt.
- **2026-11-22 16:58 · Alexander Sacharov** — critique-Befund umgesetzt: Prüfung liegt jetzt in jaira validate, update ist wieder unverändert gegenüber master. Path im Problem ist der Name mit U+FFFD statt der kaputten Bytes (JSON würde ihn ohnehin so ersetzen); der exakte Name steht im git-mv der Message.
- **2026-11-22 17:00 · Alexander Sacharov** — testing: invalid_name nutzt relative Pfade, alle anderen validate-Problems absolute; das gedruckte git mv geht nur aus dem Repo-Root. Fix: absolute Pfade in Path und git mv.
- **2026-11-22 17:00 · Alexander Sacharov** — bughunt: change: internal/cli/validate.go:163-172 (invalidNames) - class: the printed rename assumes the file is tracked by git and the command runs from the repository root, neither of which validate checks; examples: (1) on a board not shared yet (.jaira/ in .git/info/exclude since jaira init) or any untracked file under .jaira/, the printed git mv fails 'fatal: not under version control' (exit 128), and validate still exits 3 for a file macOS teammates never receive; (2) run from a subdirectory (validate works there), the root-relative paths fail 'fatal: bad source' (exit 128). Nothing is damaged, the rename just does not run as printed. Not a blocker. ran: scratch repo, jaira init, planted .jaira/logbook/t-2026/01ABC-gr\xc3.md, eval of the printed line from sub/ (exit 128) and from root while unshared (exit 128); after git add -f it runs (exit 0) and validate passes
- **2026-11-22 17:00 · Alexander Sacharov** — bughunt: done at 9b1a3d418 - 1 finding (1 change, 0 pre-existing) in 3 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/release/notes/BWT53T.md core/ticket/id.go core/ticket/id_test.go internal/cli/update.go internal/cli/updatenames_test.go internal/cli/validate.go internal/cli/validatenames_test.go
- **2026-11-22 17:02 · Alexander Sacharov** — Fix testing/bughunt: validate druckt jetzt mv statt git mv, mit absoluten Pfaden wie die anderen Problems. mv geht auf ungeteilten Boards und ungetrackten Dateien; git erkennt bei getrackten Dateien die Umbenennung trotzdem.
