---
id: 01M7X2WM5Q4Z6VPF1RAK0K8Y6S
title: "shipped.sha256 kennt jede Version, die in einem Release-Tag ausgeliefert wurde, und CI haelt das fest"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "core/shipped/shipped.sha256 enthaelt den Hash jeder Datei-Version, die je in einem v*-Tag stand; ein CI-Schritt schlaegt fehl, sobald eine ausgelieferte Version fehlt, sodass jaira update --roles eine unveraenderte alte Installation nie mehr als Handaenderung ueberspringt."
context: |-
  Gefunden vom bughunt auf TW5FE2 (20.11.2026, Note 16:50): 8 Datei-Versionen aus Release-Tags fehlen im Manifest, z.B. core/mod/builtin/jaira-cockpit/hooks/board.ts aus v0.5.2/v0.5.3, dazu Versionen von register.tsx, texts.ts, types/index.d.ts und jaira-role-acceptance/template.html.
  Folge: jaira update --roles (core/role/install.go:147 writeFile) haelt eine unberuehrte Installation dieser Version fuer von Hand geaendert und ueberspringt sie (Skipped, nur mit --force). Der Nutzer bekommt still keine neue Rolle/kein neues Cockpit.
  Nachweis: scripts/shipped-history.sh erzeugt 59 Zeilen, die im committeten Manifest fehlen; 8 davon standen in v*-Tags (git ls-tree je Tag); git log -S<hash> auf dem Manifest ist leer, sie waren nie drin.
  Ursache vermutlich: das Manifest wird bei Merges von Hand zusammengefuehrt (Konflikte append-only aufgeloest), Zwischenstaende gehen verloren.
  Alex 20.11.2026: eigenes Ticket, ja.
definition-of-done: "Alle Datei-Versionen aus allen v*-Tags stehen in core/shipped/shipped.sha256 (Test oder Skript-Pruefung)"
tags:
  - release
blocked-by: []
related: []
commits: []
created-at: 2026-11-20T18:18:21Z
updated-at: 2026-11-20T18:18:33Z
updated-by: Alexander Sacharov
---

# shipped.sha256 kennt jede Version, die in einem Release-Tag ausgeliefert wurde, und CI haelt das fest

## Definition of Done

- [ ] Alle Datei-Versionen aus allen v*-Tags stehen in core/shipped/shipped.sha256 (Test oder Skript-Pruefung)
- [ ] Ein CI-Schritt prueft gegen die Tags und schlaegt fehl, wenn eine ausgelieferte Version fehlt; laeuft auch im release-notes-Workflow vor dem Release-PR
- [ ] jaira update --roles auf einer unveraenderten v0.5.2-Installation von board.ts ueberschreibt sie statt Skipped (Test)
- [ ] Eingabeklassen: Datei in einem Tag geloescht und spaeter wieder angelegt, CRLF-Version derselben Datei, Tag ohne die Datei, Merge-Konflikt im Manifest

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

