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
updated-at: 2026-09-13T21:43:15Z
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
- **2026-09-13 21:42 · Alexander Sacharov** — Befund im Code (Lane brainstorm; gelesen: core/hook/hook.go, core/settings/settings.go, internal/cli/hook.go, core/lane/lane.go, README.md 338-350 und 561-567, hooks/sync-tasks.sh).

Der Vertrag steht wirklich: hook.Run nimmt einen Skriptpfad, setzt sieben JAIRA_*-Variablen, kappt bei 5s, verschluckt Ausgabe und Exit-Code. Daran fehlt nichts.

Fehlend ist nicht Dokumentation, sondern ein erster Lauf. Ueber den Hook existiert heute: das Feld "hook" in core/settings/settings.go:43 und vier Zeilen README. Beides beschreibt die Schnittstelle und liefert keine einzige Zeile, die man starten koennte. Wer ein leeres Feld und eine Variablenliste sieht, muss das Skript selbst erfinden - und schreibt es deshalb nicht.

Zwei Dinge aus dem Code, die die Ticket-Notiz nicht kennt:

1. Der Name "hook" ist schon belegt. internal/cli/hook.go hat "jaira hook print", und das meint etwas ganz anderes: ein Claude-Code-Stop-Hook-Snippet fuer ~/.claude/settings.json. Ein Beispiel fuer den move/claim-Hook muss unter diesem Elternbefehl unterscheidbar heissen, sonst stehen zwei unverwandte Dinge unter einem Wort.

2. Ein Skript kann heute nicht herausfinden, welche Lane einem Menschen gehoert. "jaira lanes --json" liefert agentic, terminal, precedence - aber nicht die Gates requires-question / requires-human-exit, die human und signoff erst zu Menschen-Lanes machen (core/lane/builtin/30-human.md, 45-signoff.md). agentic:false allein reicht nicht: backlog, todo, done und blocked sind ebenfalls agentic:false. Jedes Beispiel muss die Lane-Namen also hart hinschreiben und ist auf einem Board mit eigenen Lanes falsch.
- **2026-09-13 21:43 · Alexander Sacharov** — Drei Wege, wie das Beispiel ausgeliefert werden kann.

A) Nur im Repository, unter examples/ oder hooks/, in der README genannt.
Kostet fast nichts und aendert kein Byte am Binary. Gibt auf, wen es erreicht: jaira wird als einzelne Datei verteilt, "jaira self update" holt ein Binary und kein Arbeitsverzeichnis. Wer nicht geklont hat, sieht die Datei nie - und das ist genau die Person, die das leere Feld anstarrt. Praezedenzfall hooks/sync-tasks.sh zeigt den Defekt: das Skript existiert seit Langem und muss in der README erklaert werden, weil es nirgends greifbar ist.

B) Eingebettet (go:embed) und von einem Befehl auf stdout gedruckt.
Kostet rund eine Seite Text im Binary und einen Befehl. Gibt: das Beispiel ist da, wo das Programm ist, und die Zeile "jaira <befehl> > ~/.jaira/notify.sh && chmod +x" ist der erste Lauf. Genau das Muster, das "jaira hook print" schon benutzt, und genau das Muster von core/lane/builtin/*.md - jaira bettet seine Vorlagen ein, das ist hier Hausstil und kein Sonderfall.

C) Eingebettet und der Befehl schreibt Datei plus ~/.jaira/settings.json.
Gibt: null Handgriffe. Kostet: der Befehl fasst eine Datei an, die schon einen funktionierenden Hook eintragen kann, und ueberschreibt ihn im schlimmsten Fall. internal/cli/hook.go begruendet fuer den Claude-Code-Fall ausdruecklich das Gegenteil ("settings.json is the user's file"), und ~/.jaira/settings.json ist zwar jairas Datei, aber der eingetragene Hook ist es nicht.

Empfehlung: B. Erreichbarkeit ist der ganze Zweck des Tickets, also faellt A. Der Unterschied zwischen B und C ist ein einziger cp-Handgriff, C bezahlt ihn mit einem Schreibzugriff auf fremde Konfiguration - das ist kein Tausch, den ein Beispiel wert ist.
