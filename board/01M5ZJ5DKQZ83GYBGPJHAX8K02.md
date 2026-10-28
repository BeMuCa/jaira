---
id: 01M5ZJ5DKQZ83GYBGPJHAX8K02
title: "Die PR-Rolle hat einen festen Schritt vor dem Push: Scan, master, Suite, DoD, Logbook"
status: backlog
ready: false
creator: BeMuCa
goal: "jaira-role-pr fuehrt vor jedem Push einen festen Ablauf aus: Branch auf E-Mail, Denylist-Strings und Rechnernamen pruefen, master mergen, die volle Race-Suite einmal laufen lassen, die Suite-DoD-Punkte ticken, je Ticket logbook plus chore-Commit, dann push."
context: |-
  Was schief ist: fuer zwei Batches (19 Tickets, 22.-27.10.2026) wurde dieser Schritt jedes Mal von Hand aus dem Kopf gemacht (O5/O17). Die Teile stehen verstreut: Suite nur vor dem PR (Hausregel, 4N1B66), keine Firmen-Strings und kein Hostname im Repo (T22Z2C, QYK8MD), Logbook je Ticket (3SJMV8).
  Heute: core/role/builtin/jaira-role-pr/SKILL.md kennt logbook, aber weder Scan noch Suite noch DoD-Tick (grep race|denylist|hostname|email findet nichts).
  Der Scan liest die Denylist von ausserhalb des Repos (~/.config/jaira/mention-denylist, T22Z2C); die Strings selbst duerfen nie in den Skill.
  Nicht gewollt: PR oeffnen oder mergen - das bleibt beim Maintainer.
definition-of-done: SKILL.md hat den Ablauf als nummerierte Liste; jeder Schritt nennt den Befehl und das Abbruchkriterium
tags:
  - docs
  - release
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:52:23Z
updated-at: 2026-10-27T20:52:23Z
---

# Die PR-Rolle hat einen festen Schritt vor dem Push: Scan, master, Suite, DoD, Logbook

## Definition of Done

- [ ] SKILL.md hat den Ablauf als nummerierte Liste; jeder Schritt nennt den Befehl und das Abbruchkriterium
- [ ] Ein Treffer im Scan bricht vor dem Push ab und nennt Datei und Commit, nie den String selbst im Report
- [ ] Die Suite-DoD-Punkte werden mit --proof auf den Lauf getickt, danach braucht der done-Move kein --force
- [ ] Auf einem Scratch-Branch nachgestellt: der Ablauf laeuft durch und pusht

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

