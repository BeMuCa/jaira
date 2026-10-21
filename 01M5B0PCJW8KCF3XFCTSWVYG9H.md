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
updated-at: 2026-10-21T07:48:46Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-145108
claimed-at: 2026-10-21T07:48:28Z
question: ""
mode: conversational
---

# Tickets und Logbook nennen keine anderen Projekte mehr

## Definition of Done

- [ ] git grep -iE -f ~/.jaira/foreign-names auf dem Branch findet unter .jaira/ nichts mehr
- [ ] Der aktuelle Inhalt jedes Refs unter refs/jaira/ enthaelt keinen Namen aus der Liste, auf upstream geprueft nach dem Push
- [ ] Jeder ersetzte Satz bleibt verstaendlich; kein Feld ausser dem Namen ist veraendert (git diff --word-diff zeigt nur Namensersetzungen)
- [ ] Commit-Nachrichten, Ticket-Notes und PR-Text dieses Tickets nennen keinen der Namen
- [ ] Die drei Dateien ausserhalb .jaira/ (.planning/HANDOFF.md, internal/cli/staleclaim_test.go, scripts/migrate-tickets.py) nennen keinen Namen aus der Liste mehr; go test ./internal/cli/ bleibt gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-19 21:22 · Alexander Sacharov** — Teamlead 19.10.2026: der Worktree startet auf dem lokalen master 7f008dec, upstream/master ist weiter (36bde2df, release 0.5.4). Vor dem ersten Commit auf upstream/master aufsetzen und dort greppen - sonst fehlen die neuesten Logbook-Dateien.
- **2026-10-19 21:22 · Alexander Sacharov** — parent session: jaira-20
- **2026-10-19 21:25 · Alexander Sacharov** — decision (Ref-Weg, aus dem Code geklaert): ein Ticket-Ref ist eine Commit-Kette, jeder Commit hat genau eine Datei <ULID>.md als ganzen Baum (core/gitref/gitref.go refWrite/commitFile). jaira schreibt einen Ref nur aus dem Outbox, und nur fuer Tickets, die hier geschrieben wurden - ein Ref, fuer den nichts in der Queue liegt, wird von keinem jaira-Aufruf ueberschrieben. Weg also ohne CLI: hash-object + mktree + commit-tree -p <alter Ref> + git push upstream <neu>:<ref> --force-with-lease=<ref>:<alt>, danach update-ref lokal. Das ist ein normaler Commit oben drauf (fast-forward), keine Historie umgeschrieben. jaira pull scheidet aus: es claimt und aendert damit assignee auf dem Ref.
- **2026-10-19 21:25 · Alexander Sacharov** — decision: updated_at bleibt unveraendert (DoD 3: nur Namen ersetzt). Restrisiko: ein Klon, der eine alte Kopie derselben Datei auf der Platte hat und nach einem fetch darauf schreibt, schickt den alten Text mit gueltigem Lease wieder hoch (CheckFresh vergleicht nur updated_at). Betrifft nur Tickets, die auch als Datei auf master liegen; mit dem Merge dieses PRs ist auch die master-Kopie sauber. Pruefung nach dem Push: git show auf jeden Ref.
- **2026-10-19 21:25 · Alexander Sacharov** — question (vor dem Plan, an Alex): Umfang. Gezaehlt auf master+Refs: 1) drei weitere Namen eigener Projekte, die NICHT in ~/.jaira/foreign-names stehen (je 1 Fundstelle, eine davon auch auf einem Ref), dazu ein Compose-Praefix und ein Herdr-Workspace-Label, die vom Hauptnamen abgeleitet sind. 2) der GitLab-Host der Firma, 11 Stellen. 3) drei Dateien ausserhalb .jaira/ nennen einen Listen-Namen (.planning/HANDOFF.md, internal/cli/staleclaim_test.go, scripts/migrate-tickets.py). Optionen je Punkt: F1 a) mitnehmen und die drei Namen in ~/.jaira/foreign-names ergaenzen [empfohlen] b) nur die Liste. F2 a) Host neutral ersetzen ('ein firmeneigenes GitLab') [empfohlen] b) Host bleibt. F3 a) mitnehmen, DoD-Punkt dazu [empfohlen] b) eigenes Ticket c) bleibt.
- **2026-10-21 07:47 · Alexander Sacharov** — decision (Alex, 2026-10-21, ueber Teamlead) zum Umfang: F1 a) die drei weiteren Projektnamen samt abgeleitetem Compose-Praefix und Workspace-Label mitnehmen und die drei Namen in ~/.jaira/foreign-names ergaenzen. F2 a) den Firmen-GitLab-Host neutral ersetzen ('ein firmeneigenes GitLab'). F3 a) die drei Dateien ausserhalb .jaira/ (.planning/HANDOFF.md, internal/cli/staleclaim_test.go, scripts/migrate-tickets.py) hier mitnehmen, DoD-Punkt dazu. Frage damit geschlossen.
- **2026-10-21 07:48 · Alexander Sacharov** — Teamlead jaira-46, 2026-10-21: alter Worker-Tab (w3:tJ3) geschlossen - Session nicht mehr erreichbar (ListAgents kennt sie nicht, workers-Tool: Kontext unbekannt). Neuer Worker in derselben Worktree, Parent jaira-46. Der Ref dieses Tickets war lokal kaputt (leeres git-Objekt nach WSL-Absturz), von upstream neu geholt.
- **2026-10-21 07:48 · Alexander Sacharov** — parent session: jaira-46
