---
id: 01M85APMVRFHSKJB3XW63WYNV4
title: Eine statische Seite zeigt die Lanes und Boards des Katalogs und den Befehl zum Uebernehmen
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Wer den Katalog im Browser oeffnet, sieht jede Lane mit Vertrag und ganzem Prompt, jedes Board als Abfolge seiner Lanes, und kopiert den Befehl oder laedt die Datei - ohne dass irgendwo ein Server laeuft."
context: |-
  Heute sieht man den Katalog nur im Terminal ('jaira lanes market') oder als Dateiliste auf GitHub. Wer jaira noch nicht hat, sieht nicht, was es gibt.
  Loesung: eine Seite, die ein Workflow im Katalog-Repo nach jedem Merge aus index.json und den Dateien erzeugt und ueber GitHub Pages ausliefert. Kein Backend, keine Datenbank, keine Anmeldung.
  - Startseite: Liste mit Suche und Filter nach Art (Lane, Board).
  - Lane-Seite: Name, Beschreibung, Vertrag (Tier, liest, schreibt, sitzt nach, schickt zurueck nach), der ganze Prompt als Text, sha256, Knopf 'Befehl kopieren' (jaira lanes market adopt <id>), Link auf die Rohdatei.
  - Board-Seite: die Lanes als Reihe von links nach rechts mit den Rueckwaerts-Kanten, darunter je Lane ein Link.
  - 'Einreichen': ein Formular, das nichts hochlaedt, sondern einen Link auf GitHubs eigene Seite 'neue Datei' im Katalog-Repo baut (.../new/main/lanes?filename=<id>.md&value=<inhalt>). GitHub legt Fork und Pull Request selbst an. Braucht ein GitHub-Konto, aber kein git.
  Bewusst nicht: Hochladen ueber die Seite ohne GitHub. Das braucht einen Dienst, der ein Token haelt, Spam abwehrt und Konten verwaltet - genau der Server, den jaira nicht hat. Obsidian ist diesen Weg gegangen (community.obsidian.md, Anmeldung mit GitHub-Konto) und betreibt dafuer ein Portal.
  Der Prompt wird als Text ausgegeben, nie als HTML gerendert - ein eingereichter Prompt darf auf der Seite kein Skript ausfuehren.
  Der Generator liegt im Katalog-Repo, nicht im jaira-Binary.
  Vorgabe: docs/MARKET.md.
definition-of-done: Die Seite ist unter der Pages-Adresse erreichbar und listet jeden Eintrag aus index.json
tags:
  - web
blocked-by:
  - TGYQGJ
related: []
commits: []
created-at: 2026-11-23T23:08:49Z
updated-at: 2026-11-23T23:08:49Z
---

# Eine statische Seite zeigt die Lanes und Boards des Katalogs und den Befehl zum Uebernehmen

## Definition of Done

- [ ] Die Seite ist unter der Pages-Adresse erreichbar und listet jeden Eintrag aus index.json
- [ ] Eine Lane-Seite zeigt Vertrag, ganzen Prompt, sha256 und den Befehl zum Kopieren
- [ ] Eine Board-Seite zeigt die Lane-Abfolge mit Rueckwaerts-Kanten
- [ ] Ein Prompt mit <script>-Text wird als Text angezeigt und nicht ausgefuehrt (Test im Generator)
- [ ] Der Einreichen-Link oeffnet GitHubs Neue-Datei-Seite mit vorbelegtem Inhalt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

