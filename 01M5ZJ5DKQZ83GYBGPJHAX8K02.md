---
id: 01M5ZJ5DKQZ83GYBGPJHAX8K02
title: "Die PR-Rolle hat einen festen Schritt vor dem Push: Scan, master, Suite, DoD, Logbook"
status: backlog
ready: true
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
blocked-by:
  - 01M2HZKCQ48BPQYN1YPTXB97NJ
related:
  - 01M7WWYB5XZMD01DZY3MJG3REW
commits: []
created-at: 2026-10-27T20:52:23Z
updated-at: 2026-11-20T16:49:23Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
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
- **2026-11-20 16:49 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
