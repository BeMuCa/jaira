---
id: 01M2E86MS66ZZ58RRV3EY0A7DT
title: "Ein Beispiel-Hook liegt bei, damit Lane-Wechsel jemanden erreichen"
status: brainstorm
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "jaira bringt ein Hook-Beispiel mit, das man mit einem Befehl einsetzt; danach meldet sich jeder Lane-Wechsel dort, wo der Mensch hinsieht"
context: |-
  core/hook/hook.go ruft bei jedem move und jedem claim ein Skript auf und uebergibt JAIRA_EVENT, JAIRA_TICKET, JAIRA_TITLE, JAIRA_STATUS, JAIRA_ASSIGNEE, JAIRA_ACTOR und JAIRA_ROOT als Umgebungsvariablen. Der Mechanismus ist fertig und gut: jaira bringt bewusst keine eigene Abhaengigkeit mit, die Anbindung gehoert dem Nutzer.
  Was fehlt, ist der erste Schritt. Wer 'hook' in den Einstellungen sieht, hat ein leeres Feld und keinen Anhaltspunkt, was hineingehoert. Praktisch schreibt es deshalb niemand.
  Am 2026-09-13 ist hier eines von Hand entstanden, ~/.jaira/notify.sh, rund 25 Zeilen: es prueft HERDR_ENV, verwandelt den Lane-Wechsel in 'herdr notification show', und waehlt den Ton nach Ziel-Lane - ein Lane, der einem Menschen gehoert, bekommt request, done bekommt done, alles andere bleibt stumm. Die Regel dahinter, die das Beispiel transportieren soll: einen Ton verdient nur der Zustand, in dem sich ohne den Menschen nichts bewegt.
  Zu klaeren, absichtlich nicht entschieden:
  - ob das Beispiel eingebettet wird (wie die Lanes, go:embed) und per Befehl geschrieben, oder ob es nur im Repository unter einem Beispielordner liegt und in der Dokumentation genannt wird.
  - ob mehr als ein Beispiel beiliegt. Herdr ist das, was hier laeuft; ntfy.sh und eine Terminal-Glocke sind die naechstliegenden, weil sie keine Konten brauchen. Slack braucht eine Webhook-URL und ist damit kein Beispiel, sondern eine Einrichtung.
  - ob der Einsetz-Befehl auch settings.json schreibt oder nur die Datei hinlegt und sagt, welche Zeile fehlt. Fremde Einstellungen ungefragt zu aendern ist die unangenehmere Variante.
  Nicht Teil dieses Tickets: an core/hook selbst etwas aendern. Der Vertrag steht.
definition-of-done: "Ein Befehl legt ein lauffaehiges Hook-Beispiel ab und sagt, wie es scharfgeschaltet wird; das Beispiel laeuft ohne Herdr fehlerfrei durch und tut dann nichts; ein Lane-Wechsel in eine Lane, die einem Menschen gehoert, ist hoerbar von einem gewoehnlichen unterscheidbar; ein Hinweis auf das Beispiel steht dort, wo 'hook' dokumentiert ist; eine Zeile in core/release/NOTES.md unter ## Unreleased; go test ./... -race gruen"
tags:
  - cli
blocked-by: []
parent: 01M2E248SM9X1JRZBNTHC9V7ZV
related: []
commits: []
created-at: 2026-09-13T20:44:07Z
updated-at: 2026-09-13T21:42:41Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-55522
claimed-at: 2026-09-13T21:42:41Z
---

# Ein Beispiel-Hook liegt bei, damit Lane-Wechsel jemanden erreichen

## Definition of Done

- [ ] Ein Befehl legt ein lauffaehiges Hook-Beispiel ab und sagt, wie es scharfgeschaltet wird; das Beispiel laeuft ohne Herdr fehlerfrei durch und tut dann nichts; ein Lane-Wechsel in eine Lane, die einem Menschen gehoert, ist hoerbar von einem gewoehnlichen unterscheidbar; ein Hinweis auf das Beispiel steht dort, wo 'hook' dokumentiert ist; eine Zeile in core/release/NOTES.md unter ## Unreleased; go test ./... -race gruen

## Options

- [x] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

