---
id: 01M67K6X23DJ7TDVQ6DVDVPAPE
title: "Jede Person hat pro Projekt einen festen jaira-Namen, einmal abgefragt und lokal gespeichert; er ersetzt Git-Name und Aliasse"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Ein Mensch erscheint auf dem Board unter genau einem Namen: abgefragt beim ersten Start in einem Projekt, gespeichert als lokale, gitignorierte Projekteinstellung, benutzt fuer Zuweisung, Tracking und die Benennung seiner Refs."
context: |-
  Wunsch von berkc am 31.10.2026 (Punkt 2).
  Problem heute: das Board fuehrt dieselben Menschen unter mehreren Namen - BeMuCa (121) und berkc (88) fuer berkc, 'Alexander Sacharov' (48), 'alex', 'AlSa' fuer seinen Kollegen.
  Heute kommt der Name aus git user.name (internal/cli/root.go:353 identity -> core/identity); Aliasse stehen global in ~/.jaira/identity, nicht pro Projekt.
  Projekteinstellungen gibt es nur als committete Dateien (.jaira/chat-model, .jaira/release, bald .jaira/active-boards aus KVMG0K); eine lokale, gitignorierte Projekteinstellung gibt es nicht - die braucht es dafuer (z.B. .jaira/local/settings.json, gitignoriert).
  Web-Nutzer ohne Claude/CLI (PMs) werden im Web nach ihrem Namen gefragt.
  Offen: bestehende Tickets mit alten Namen umschreiben (Migration) oder nur als Alias lesen?
definition-of-done: "Erster Start in einem Projekt fragt den jaira-Namen ab und speichert ihn in einer lokalen, gitignorierten Projekteinstellung; Test"
tags:
  - cli
blocked-by: []
related:
  - 01M5MNWJF240306ZVDA41K3RA8
commits: []
created-at: 2026-10-30T23:44:36Z
updated-at: 2026-10-30T23:44:36Z
---

# Jede Person hat pro Projekt einen festen jaira-Namen, einmal abgefragt und lokal gespeichert; er ersetzt Git-Name und Aliasse

## Definition of Done

- [ ] Erster Start in einem Projekt fragt den jaira-Namen ab und speichert ihn in einer lokalen, gitignorierten Projekteinstellung; Test
- [ ] creator, assignee, claimed-by und neue Refs nutzen diesen Namen; alte Namen werden als Alias derselben Person gelesen
- [ ] jaira web fragt einen Nutzer ohne Namen im Browser danach
- [ ] Release-Note-Zeile

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

