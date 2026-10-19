---
id: 01M5B0PCJW8KCF3XFCTSWVYG9H
title: Tickets und Logbook nennen keine anderen Projekte mehr
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Kein Text unter .jaira/ (Tickets, Logbook, Archiv, Ticket-Refs) nennt mehr den Namen eines anderen Projekts des Maintainers."
context: |-
  Was falsch ist: Ticket-Texte und Notes nennen andere, private Projekte des Maintainers beim Namen. jaira ist oeffentlich auf GitHub, das ist ein Datenleck.
  Die Namen stehen NICHT hier, damit dieses Ticket das Leck nicht wiederholt. Liste: ~/.jaira/foreign-names (eine Zeile je Name, nur auf Alex' Rechner).
  Wo sie stehen (Stand 19.10.2026): upstream/master: gut 20 Dateien unter .jaira/tickets/, .jaira/logbook/ und .jaira/archive/ (git grep -ciE -f ~/.jaira/foreign-names upstream/master). Refs: 9 Tickets unter refs/jaira/tickets/* (aktueller Inhalt jedes Refs mit git show <ref>:<datei> pruefen).
  Anlass: Alex 19.10.2026 - 'nie Projektnamen verwenden'. Auf CRDH6B wurde eine Note schon per sed anonymisiert (Commit 54583386 auf feat/crdh6b).
  Entschieden von Alex: nur der heutige Text wird bereinigt, mit normalen Commits und einem PR. Die git-Historie (master und Refs) wird NICHT umgeschrieben, kein force-push.
  Ausnahme: die CLI kann Notes nicht aendern. Deshalb ist hier ausdruecklich erlaubt, Dateien unter .jaira/ von Hand zu aendern - nur das Ersetzen der Namen, sonst nichts.
  Ersatz: neutral beschreiben, was gemeint war ('ein anderes Projekt', 'ein Epic in einem anderen Projekt', 'das Video-Board' usw.), damit der Satz verstaendlich bleibt.
  Offen und zuerst klaeren: wie eine von Hand geaenderte Ticket-Datei auf ihren Ref kommt (refsync, z.B. core/refsync/refsync.go Record), ohne dass der naechste jaira-Aufruf den alten Text vom Ref zurueckschreibt. Pruefen mit git show auf den Ref nach dem Schreiben.
  Im Commit und PR selbst die Namen auch nicht nennen.
definition-of-done: git grep -iE -f ~/.jaira/foreign-names auf dem Branch findet unter .jaira/ nichts mehr
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-10-19T21:22:16Z
updated-at: 2026-10-19T21:22:45Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-10599
claimed-at: 2026-10-19T21:22:42Z
---

# Tickets und Logbook nennen keine anderen Projekte mehr

## Definition of Done

- [ ] git grep -iE -f ~/.jaira/foreign-names auf dem Branch findet unter .jaira/ nichts mehr
- [ ] Der aktuelle Inhalt jedes Refs unter refs/jaira/ enthaelt keinen Namen aus der Liste, auf upstream geprueft nach dem Push
- [ ] Jeder ersetzte Satz bleibt verstaendlich; kein Feld ausser dem Namen ist veraendert (git diff --word-diff zeigt nur Namensersetzungen)
- [ ] Commit-Nachrichten, Ticket-Notes und PR-Text dieses Tickets nennen keinen der Namen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-19 21:22 · Alexander Sacharov** — Teamlead 19.10.2026: der Worktree startet auf dem lokalen master 7f008dec, upstream/master ist weiter (36bde2df, release 0.5.4). Vor dem ersten Commit auf upstream/master aufsetzen und dort greppen - sonst fehlen die neuesten Logbook-Dateien.
- **2026-10-19 21:22 · Alexander Sacharov** — parent session: jaira-20
