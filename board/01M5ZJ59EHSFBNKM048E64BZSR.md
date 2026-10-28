---
id: 01M5ZJ59EHSFBNKM048E64BZSR
title: "jaira show nennt die optionalen Lanes, an denen next_lane vorbeigeht"
status: backlog
ready: false
creator: BeMuCa
goal: "jaira show --json fuehrt neben next_lane die optionalen Lanes (brainstorm, pre-process) auf, die uebersprungen werden, und wie man sie einschaltet (jaira dod <id> --option <lane>); ein Driver kann so Planung waehlen, statt sie unbemerkt zu verpassen."
context: |-
  Was schief ist: Driver uebersprangen pre-process bei Tickets, die einen Plan gebraucht haetten, bis jemand sagte, dass man --option vorher ticken muss (O11, 27.10.2026).
  Heute: jaira dod --option <lane> schaltet eine Schritt-Option ein (jaira dod --help); jaira show --json zeigt nur next_lane, nicht, dass eine optionale Lane uebersprungen wird; der Agent-Block erwaehnt --option nicht.
  Nicht gewollt: optionale Lanes standardmaessig einschalten - dann sind sie nicht optional.
definition-of-done: "jaira show <id> --json hat ein Feld, das die uebersprungenen optionalen Lanes und den Befehl zum Einschalten nennt"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:52:19Z
updated-at: 2026-10-27T20:52:19Z
---

# jaira show nennt die optionalen Lanes, an denen next_lane vorbeigeht

## Definition of Done

- [ ] jaira show <id> --json hat ein Feld, das die uebersprungenen optionalen Lanes und den Befehl zum Einschalten nennt
- [ ] Der Agent-Block nennt --option in einem Satz
- [ ] Release-Notiz in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

