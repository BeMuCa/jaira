---
id: 01M85AKKJ5S50P9BXTYY5TED33
title: "jaira market lint prueft ein Katalog-Verzeichnis, bevor ein Eintrag hineinkommt"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Ein Befehl sagt fuer ein Katalog-Verzeichnis, ob jede Datei hineindarf, mit Exit-Code und einer Zeile je Verstoss - CI und ein Firmen-Spiegel nutzen dieselbe Pruefung."
context: |-
  Heute prueft nur ein Go-Test, ob eine Katalog-Lane parst: core/lane/shipped_test.go laedt lanes/*.md. Mehr wird nicht geprueft.
  Nicht geprueft: welche Dateien ein Pull Request anlegt, wie gross sie sind, ob ein Prompt versteckten Text traegt.
  Der Test laeuft nur im Repo BeMuCa/jaira. Ein eigenes Katalog-Repo (Milestone extension-marketplace) kann ihn nicht aufrufen.
  Loesung: die Pruefung zieht ins Binary, als 'jaira market lint <verzeichnis>'. Das Katalog-Repo laedt ein Release-Binary und ruft sie auf. Ein Firmen-Spiegel hat damit dieselbe Pruefung ohne eigenes Werkzeug.
  Regeln, je mit eigener Fehlerzeile:
  - nur erlaubte Pfade: lanes/<id>.md, boards/<id>/description, boards/<id>/order, boards/<id>/lanes/<id>.md, README.md. Alles andere ist ein Verstoss (kein Skript, kein Binary, kein Symlink).
  - jede Lane parst und laedt ohne Warnung (lane.Parse / lane.Load, wie shipped_test.go), id gleich Dateiname, kebab-case.
  - Groesse: hoechstens 32 KB je Datei.
  - keine versteckten Zeichen: Zero-Width- und Bidi-Steuerzeichen, HTML-Kommentare im Prompt. Dort versteckt man Anweisungen, die ein Reviewer im Diff nicht sieht.
  - description ist nicht leer; bei einem Board nennt order nur Lanes, die im Board-Ordner liegen oder im Binary stecken.
  Die Regel 'kein script:' ist ein eigenes Ticket, weil das Feld erst mit ENCFVZ existiert.
  Exit 0 ohne Verstoss, sonst Exit 1; --json gibt die Verstoesse als Liste.
  Vorgabe: docs/MARKET.md.
definition-of-done: "jaira market lint <dir> endet mit 0 fuer ein gueltiges Verzeichnis und mit 1 und einer Zeile je Verstoss sonst"
tags:
  - cli
  - ci
blocked-by:
  - 01M8BQFBMS778FE9WK6YDAW7WP
related: []
commits: []
created-at: 2026-11-23T23:07:10Z
updated-at: 2026-11-26T10:50:51Z
updated-by: Berk
---

# jaira market lint prueft ein Katalog-Verzeichnis, bevor ein Eintrag hineinkommt

## Definition of Done

- [ ] jaira market lint <dir> endet mit 0 fuer ein gueltiges Verzeichnis und mit 1 und einer Zeile je Verstoss sonst
- [ ] Je Regel ein Test: fremder Pfad, Symlink, zu grosse Datei, Zero-Width-Zeichen, HTML-Kommentar, id ungleich Dateiname, leere description
- [ ] --json liefert die Verstoesse mit path und rule
- [ ] core/lane/shipped_test.go ruft dieselbe Pruefung auf statt eigener Logik
- [ ] Release-Note core/release/notes/<TICKET>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

