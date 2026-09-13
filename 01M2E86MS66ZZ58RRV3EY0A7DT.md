---
id: 01M2E86MS66ZZ58RRV3EY0A7DT
title: "Ein Beispiel-Hook liegt bei, damit Lane-Wechsel jemanden erreichen"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "jaira bringt ein lauffaehiges Beispiel-Hookskript im Binary mit und druckt es auf Zuruf, damit der erste Lane-Wechsel ohne eigene Erfindung bei einem Menschen ankommt - und nur dann, wenn sich ohne ihn nichts bewegt."
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
updated-at: 2026-09-13T21:50:07Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-60947
claimed-at: 2026-09-13T21:46:35Z
outcome-what: "Plan fuer das mitgelieferte Hook-Beispiel festgelegt: neun Schritte von Skript ueber go:embed und 'jaira hook example' bis README, NOTES.md und Tests"
outcome-why: "Die Brainstorm-Lane hatte die drei offenen Fragen entschieden, aber nicht gesagt, wo die Datei liegt, was die Zustellzeile ist und wie das Hoerbare pruefbar wird"
outcome-resolves: "Ticket hat eine Plan-Checkliste, die Schritt fuer Schritt abgearbeitet werden kann"
---

# Ein Beispiel-Hook liegt bei, damit Lane-Wechsel jemanden erreichen

## Definition of Done

- [ ] Ein Befehl legt ein lauffaehiges Hook-Beispiel ab und sagt, wie es scharfgeschaltet wird; das Beispiel laeuft ohne Herdr fehlerfrei durch und tut dann nichts; ein Lane-Wechsel in eine Lane, die einem Menschen gehoert, ist hoerbar von einem gewoehnlichen unterscheidbar; ein Hinweis auf das Beispiel steht dort, wo 'hook' dokumentiert ist; eine Zeile in core/release/NOTES.md unter ## Unreleased; go test ./... -race gruen

## Options

- [x] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Beispielskript core/hook/example/notify.sh schreiben: Zustellzeile ist die Terminal-Glocke, ntfy- und notify-send-Einzeiler auskommentiert direkt darunter; Lane-Tabelle human|signoff -> Bitte-Ton, done -> Abschluss-Ton, blocked und alle Agenten-Lanes stumm; Kopfkommentar englisch mit den JAIRA_*-Variablen, der 5s-Kappung und dem Satz, welche Zeile zu ersetzen ist
- [x] Skript einbetten: neue Datei core/hook/example.go mit go:embed example/notify.sh und einer Funktion ExampleScript() string; core/hook/hook.go bleibt unberuehrt, der Vertrag aendert sich nicht
- [~] Befehl 'jaira hook example' in internal/cli/hook.go ergaenzen: druckt das Skript auf stdout, schreibt keine Datei und keine Einstellung; Long-Text nennt 'jaira hook example > ~/.jaira/notify.sh && chmod +x ~/.jaira/notify.sh' und die fehlende Zeile "hook": "/home/<du>/.jaira/notify.sh"
- [ ] Test internal/cli/hook_example_test.go nach dem Muster von runStopHook: Ausgabe von 'hook example' in eine Datei schreiben und mit sh und gesetzten JAIRA_*-Variablen ausfuehren; human, signoff und done geben hoerbare Ausgabe, in-progress, todo und blocked geben nichts, jeder Fall endet mit exit 0
- [ ] Test fuer den leeren Rechner: ohne jedes Zustellwerkzeug auf PATH laeuft das Skript fehlerfrei durch und tut nichts (das ist die DoD-Zeile 'laeuft ohne Herdr fehlerfrei durch')
- [ ] README.md: im Hook-Absatz bei Zeile 344-350 zwei Saetze auf 'jaira hook example'; bei Zeile 561 einen Halbsatz, dass 'hook print' und 'hook example' zwei unverwandte Dinge sind
- [ ] core/release/NOTES.md: eine englische Zeile unter ## Unreleased, als Anweisung formuliert
- [ ] go test ./... -race gruen, danach die DoD-Zeile mit --proof abhaken
- [ ] Folgeticket erfassen: 'jaira lanes --json' gibt requires-question und requires-human-exit nicht aus, deshalb stehen die Lane-Namen im Beispiel hart und es ist auf einem Board mit eigenen Lanes stumm

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
- **2026-09-13 21:43 · Alexander Sacharov** — Die drei offenen Fragen, entschieden.

Frage 1 - eingebettet und per Befehl, oder nur im Repository: eingebettet, per Befehl auf stdout. Begruendung oben (Weg B). Der Befehl heisst nicht "jaira hook print" - das ist vergeben - sondern muss den move/claim-Hook benennen, Vorschlag "jaira hook example". Umbenennen von "hook print" ist kein Teil dieses Tickets.

Frage 2 - wie viele Beispiele: genau eines, und zwar keines, das auf diesem Rechner haengt. ~/.jaira/notify.sh braucht HERDR_ENV und HERDR_BIN_PATH und beendet sich bei allen anderen mit exit 0 - als mitgeliefertes Beispiel waere es ein Skript, das fuer fast jeden nichts tut, und das ist schlimmer als keins, weil man den Fehler nicht sieht. Was ohne Konto, ohne URL und ohne Installation laeuft, ist die Terminal-Glocke; ntfy.sh braucht ein Topic und Slack eine Webhook-URL, beides ist Einrichtung und kein Beispiel. Also: ein Skript, dessen Zustellzeile die eine offensichtlich zu ersetzende Zeile ist, mit ntfy und einem Desktop-Toast als auskommentierte Einzeiler direkt darunter. Eine Datei, drei Zustellwege, kein Menue aus Befehlen - und damit auch die Antwort auf "ist das kleiner als paca".

Frage 3 - schreibt der Befehl settings.json: nein. Er druckt und nennt die fehlende Zeile ("hook": "/home/<du>/.jaira/notify.sh"), so wie "jaira hook print" es vormacht. Kein --write in der ersten Fassung.

Die Regel, die das Beispiel tragen muss, und was sie konkret heisst: einen Ton verdient nur der Zustand, in dem sich ohne den Menschen nichts bewegt. Das sind die Lanes mit requires-question oder requires-human-exit (human, signoff) - dort steht die Arbeit bis jemand entscheidet. done bekommt einen anderen, abschliessenden Ton, weil es das Ende ist und nicht die Bitte. blocked ist der Grenzfall: die Lane wartet laut Definition auf eine externe Abhaengigkeit, nicht auf den Menschen - das handgeschriebene notify.sh gibt ihr trotzdem request. Fuer das mitgelieferte Beispiel bleibt blocked stumm, weil sonst die Regel, die es lehren soll, im Beispiel selbst schon gebrochen ist. Alles andere - jeder Wechsel zwischen Agenten-Lanes - ist stumm.

Offen und bewusst nicht hier geloest: dass die Lane-Namen im Beispiel hart stehen muessen, weil "jaira lanes --json" die Gates nicht ausgibt. Auf einem Board mit eigenen Lanes ist das Beispiel dann stumm statt falsch, was die ertraeglichere Richtung ist. Gehoert als eigenes Ticket erfasst, nicht hier.
- **2026-09-13 21:47 · Alexander Sacharov** — Warum der Plan so aussieht (Lane pre-process; gelesen: internal/cli/hook.go, internal/cli/hook_test.go, core/hook/hook.go, core/release/NOTES.md, README.md 338-350 und 550-567, ~/.jaira/notify.sh).

Drei Entscheidungen, die die Brainstorm-Notizen offen gelassen haben und die der Plan jetzt festlegt:

1. Wo das Skript liegt: core/hook/example/notify.sh, eingebettet von einer neuen Datei core/hook/example.go. Nicht in internal/cli, weil jaira seine Vorlagen in core/ einbettet (core/lane/builtin/*.md, core/release/NOTES.md) und internal/cli heute kein go:embed hat. core/hook/hook.go selbst wird nicht angefasst - das Ticket verbietet Aenderungen am Vertrag, eine zusaetzliche Datei im selben Paket ist keine.

2. Warum die Terminal-Glocke die Zustellzeile ist und nicht nur ein Kommentar: sie ist das Einzige, was ohne Konto, URL und Installation wirklich etwas tut - und damit das Einzige, was ein Test beobachten kann. Ein Beispiel, dessen aktive Zeile auskommentiert ist, laeuft zwar fehlerfrei durch, aber die DoD-Zeile "ein Lane-Wechsel in eine Lane, die einem Menschen gehoert, ist hoerbar von einem gewoehnlichen unterscheidbar" waere dann nicht pruefbar. Mit der Glocke ist sie es: das Byte ist da oder nicht.

3. Wie getestet wird: exakt nach dem Muster, das internal/cli/hook_test.go fuer 'hook print' schon aufgebaut hat - die Ausgabe des Befehls wird als das ausgefuehrt, was der Nutzer einfuegt, nicht eine Go-Umschreibung davon. Der Unterschied: hier gibt es keinen Stub auf PATH, sondern gesetzte JAIRA_*-Variablen und die Frage, ob Ausgabe entsteht.

Was der Plan bewusst nicht tut: 'jaira hook print' umbenennen. Die zwei Befehle unter 'hook' meinen weiterhin Unterschiedliches (Claude-Code-Stop-Hook vs. move/claim-Hook); der Plan loest das mit einem Halbsatz in der README statt mit einer Umbenennung, die kein Teil dieses Tickets ist.

Reihenfolge: Skript zuerst, dann Einbettung, dann Befehl, dann Tests. Umgekehrt haette man einen Befehl, der eine Datei druckt, deren Inhalt noch nicht entschieden ist - und die Zustellzeile ist die einzige echte Entscheidung im Skript.
