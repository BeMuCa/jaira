---
id: 01M851YE2ACMD981Z6V1Q030WZ
title: Unter JAIRA_USER zaehlt jaira die per E-Mail zugewiesenen Tickets des Checkout-Besitzers als die der anderen Person
status: testing
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ist JAIRA_USER gesetzt, gilt nur dessen Name samt seinen Aliassen als 'ich' - nicht zusaetzlich die git user.email des Checkouts."
context: |-
  Was falsch ist: core/identity/identity.go:112 (func Aliases) nimmt git user.email immer dazu, auch wenn JAIRA_USER eine andere Person nennt.
  Folge: unter JAIRA_USER=grace zaehlt IsMe die Tickets, die per E-Mail dem Besitzer des Checkouts gehoeren, als grace' Tickets (list --mine, Claim-Pruefungen).
  Gefunden im Review zu DVPAPE am 23.11.2026; gab es schon vorher auf master.
  Erst nach DVPAPE starten: DVPAPE baut core/identity gerade um (Personenliste auf refs/jaira/people), beide Tickets aendern dieselbe Datei.
definition-of-done: "Unter JAIRA_USER liefert Aliases nur JAIRA_USER und dessen Aliasse aus der Personenliste, nicht die git user.email des Checkouts; Test"
tags:
  - cli
blocked-by: []
related: []
follows: 01M67K6X23DJ7TDVQ6DVDVPAPE
commits: []
created-at: 2026-11-23T20:35:47Z
updated-at: 2026-11-23T21:10:50Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-37226
claimed-at: 2026-11-23T21:01:32Z
outcome-what: Aliases returns only JAIRA_USER and its people-list names when JAIRA_USER is set; whoami marks the alias file as unread then
outcome-why: "git user.email and ~/.jaira/identity belong to the checkout owner; under JAIRA_USER another person acts and must not own the owner's tickets"
outcome-resolves: "list --mine / claim / move checks under JAIRA_USER counted the owner's email-assigned tickets as the other person's"
review-summary: none
review-gaps: "folded the four strings.TrimSpace(os.Getenv(JAIRA_USER)) reads (currentIn, Aliases, whoami) into identity.ByEnv(); raw() keeps its own read because it needs the value; nothing dead, no cost on hot paths"
---

# Unter JAIRA_USER zaehlt jaira die per E-Mail zugewiesenen Tickets des Checkout-Besitzers als die der anderen Person

## Definition of Done

- [x] Unter JAIRA_USER liefert Aliases nur JAIRA_USER und dessen Aliasse aus der Personenliste, nicht die git user.email des Checkouts; Test
  proof: core/identity/identity.go Aliases (byEnv return); TestAliasesUnderJairaUserLeaveOutTheOwnersNames
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/Q030WZ.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 21:00 · Alexander Sacharov** — teamlead 23.11.2026: auf Alex' Wunsch sofort gestartet, gestapelt auf feat/dvpape (DVPAPE noch nicht auf master; beide aendern core/identity). Gemergt wird nach DVPAPE. Aliases soll unter JAIRA_USER die Personenliste aus DVPAPE (core/people) nutzen.
- **2026-11-23 21:01 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/48662.sock
- **2026-11-23 21:02 · Alexander Sacharov** — decision (worker): open decisions none. Under JAIRA_USER, Aliases also skips the machine's alias file (~/.jaira/identity, JAIRA_HOME): it belongs to the checkout's owner just like git user.email, and the DoD says 'nur JAIRA_USER und dessen Aliasse aus der Personenliste'. Without JAIRA_USER nothing changes.
- **2026-11-23 21:10 · Alexander Sacharov** — critique+optimize in one pass (49 lines, no lock/parser). critique: none — fix sits in the callee (Aliases), every caller (gate, list --mine, tui, release) picks it up. optimize: four inline JAIRA_USER reads folded into identity.ByEnv().
