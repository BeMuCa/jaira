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
updated-at: 2026-11-22T16:56:29Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1768841
claimed-at: 2026-11-22T16:50:44Z
outcome-what: Slug schneidet an einer Zeichengrenze; jaira update nennt Dateien unter .jaira/ mit ungültigem UTF-8-Namen samt git mv
outcome-why: "Ein Schnitt mitten in ü/ß/Kyrillisch erzeugte Dateinamen, die macOS nicht auschecken kann (PR #184)"
outcome-resolves: Ticket- und Logbook-Dateien sind immer gültiges UTF-8; Altlasten werden gemeldet
---

# Ticket-Dateinamen schneiden Umlaute nicht mehr mitten im Zeichen ab

## Definition of Done

- [x] Slug kürzt auf höchstens 48 Bytes an einer Zeichengrenze; Tabellentest mit Umlaut, ß, Kyrillisch und Emoji genau an der Grenze, jedes Ergebnis utf8.ValidString
  proof: core/ticket/id.go:59 (Schnitt an utf8.RuneStart); core/ticket/id_test.go TestSlugCutsAtRuneBoundary + TestSlugAlwaysValidUTF8
- [x] Vorhandene Dateien mit ungültigem Namen: jaira doctor oder update meldet sie (oder benennt sie um) - Entscheidung im Plan, Test
  proof: internal/cli/update.go invalidNames/printInvalidNames; internal/cli/updatenames_test.go TestUpdateNamesFilesWithInvalidUTF8 (führt das gedruckte git mv per bash aus)
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
