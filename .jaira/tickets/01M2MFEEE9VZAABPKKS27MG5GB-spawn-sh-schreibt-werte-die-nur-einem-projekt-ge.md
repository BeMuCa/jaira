---
id: 01M2MFEEE9VZAABPKKS27MG5GB
title: "spawn.sh schreibt Werte, die nur einem Projekt gehoeren"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "spawn.sh legt den worktree an, oeffnet die Kachel und startet den Arbeiter - mehr nicht. Alles, was ein Projekt zum Loslaufen braucht, macht ein Skript im Repository, das spawn.sh aufruft und mit worktree-Pfad, Slug, Port-Versatz und Wurzel versorgt. Ohne dieses Skript legt spawn.sh nur den worktree an."
context: |-
  spawn.sh legt fuer jeden Arbeiter einen worktree an und schreibt ihm eine eigene .env. Darin stehen heute HTTP_PORT und DB_PORT_HOST - Namen, die sich ein einzelnes Projekt ausgedacht hat. Ein anderes Projekt nennt seine Ports anders oder hat keine, und bekommt die Zeilen trotzdem.

  Die Folge sieht man an einem Board, das jaira benutzt: dessen spawn.sh ist von 100 auf 159 Zeilen gewachsen, weil dort IMAGE_NS, COMPOSE_PROFILES, VITE_PORT_HOST und BACKEND_PORT_HOST nachgetragen wurden. Alles vier gehoert diesem einen Projekt.

  Das Schlimme daran ist nicht das Ueberschreiben - 'roles install' laesst eine bearbeitete Datei stehen und meldet sie. Es ist das Gegenteil: die Datei ist eingefroren. Jede Verbesserung an spawn.sh wird bei diesem Anwender still uebersprungen, und er merkt es erst, wenn er die Ausgabe von 'roles install' liest. Der Fork ist auch nicht sortenrein: --no-worktree und der Lane-Name 'dispatch' stehen darin und sind allgemein, sie gehoeren hierher.

  Am 16.09.2026 aufgekommen, als in einem worktree die COMPOSE_PROFILES des Hauptverzeichnisses mitkopiert wurden und der Zweig acht Container zu viel hochfuhr.

  Noch offen und Teil der Loesung: welche Angaben der Hook braucht. Vorschlag: Slug und Port-Versatz als Argumente, damit die Datei ihre Ports genauso verteilen kann wie spawn.sh heute.
definition-of-done: "spawn.sh ruft .jaira/worktree-env auf, wenn die Datei existiert und ausfuehrbar ist, und haengt ihre Ausgabe an die .env des worktree"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-16T06:46:09Z
updated-at: 2026-09-16T07:50:16Z
updated-by: Alexander Sacharov
---

# spawn.sh schreibt Werte, die nur einem Projekt gehoeren

## Definition of Done

- [ ] spawn.sh ruft ein Einrichtungsskript des Repositorys auf und uebergibt ihm worktree-Pfad, Slug, Port-Versatz und Repository-Wurzel
- [ ] Ohne dieses Skript legt spawn.sh den worktree an und richtet sonst nichts ein
- [ ] Weder .env-Kopie noch HTTP_PORT oder DB_PORT_HOST stehen noch in spawn.sh; die Dokumentation zeigt ein Beispielskript, das beides macht
- [ ] Die Hilfe von 'roles install' sagt, dass eine bearbeitete Datei nicht nur unveraendert bleibt, sondern auch keine Verbesserungen mehr bekommt
- [ ] Ein Board, das bisher ohne Einrichtungsskript lief, bekommt beim naechsten Arbeiter eine verstaendliche Meldung statt einer stillen Verhaltensaenderung

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-16 06:46 · Alexander Sacharov** — Praezisiert am 16.09.2026 auf Ansage: jaira soll den worktree anlegen und sonst nichts. Die erste Fassung dieses Tickets hat noch von einer Datei gesprochen, die Zeilen fuer die .env ausgibt - das ist zu eng gedacht.

Warum zu eng: das Kopieren der .env ist selbst schon projektspezifisch. Ein Projekt ohne .env braucht es nicht, ein anderes will daneben noch 'npm install' oder eine Datenbank hochfahren. Wenn jaira die .env kopiert und ein Hook nur Zeilen anhaengen darf, bleibt die eine Entscheidung, die am haeufigsten falsch ist, weiter in jaira.

Also: ein Einrichtungsskript, kein Zeilenlieferant. jaira ruft es auf, gibt ihm worktree-Pfad, Slug, Port-Versatz und Repository-Wurzel, und was danach in dem Verzeichnis steht, ist Sache des Projekts. Rueckwaerts vertraeglich bleibt es dadurch, dass ein Board ohne dieses Skript einen worktree ohne .env bekommt - wer die alte Bequemlichkeit will, schreibt sich das Skript einmal hin, und jaira zeigt eins als Vorlage.
- **2026-09-16 06:48 · Alexander Sacharov** — Reihenfolge, entschieden von der Teamlead-Sitzung am 2026-09-16, damit niemand hier anfaengt, bevor 7KX89C durch ist.

Drei Tickets fassen dieselbe Datei an: 7KX89C (spawn.sh faehrt im Repository mit), dieses hier (Projektwerte raus aus spawn.sh), 3YRPXJ (spawn.sh kann nur eine Rolle starten). Das ist eine Arbeit an einem File, keine drei nebeneinander.

7KX89C geht zuerst, und der Grund steht im Kontext DIESES Tickets: spawn.sh liegt heute nur in ~/.claude/skills/jaira-dispatcher/ auf einer Maschine. Wer hier zuerst baut, aendert eine Datei, die in keinem Diff steht, in keinem PR auftaucht und bei Berk nicht ankommt - und muss dieselbe Aenderung ein zweites Mal machen, sobald 7KX89C die Datei ins Repository holt. Genau das Einfrieren, das dieses Ticket beklagt, wuerde die Behebung dieses Tickets selbst treffen.

Danach dieses Ticket, dann 3YRPXJ.

Was aus der Sitzung noch dazugehoert: der Teamlead hat dreimal eine gepatchte Kopie von spawn.sh im Scratchpad gefahren, weil die Prompt-Zeile fest verdrahtet ist. Das ist 3YRPXJ und dort notiert - aber es ist dasselbe Muster wie hier: wer spawn.sh nicht erweitern kann, forkt es. Eine Loesung fuer dieses Ticket, die nur die .env herausnimmt, laesst den Fork-Grund von 3YRPXJ stehen. Die Plan-Lane sollte beide Haken - Einrichtungs-Hook und Prompt-Argument - als eine Erweiterbarkeitsfrage ansehen, auch wenn sie in zwei Tickets gebaut werden.
- **2026-09-16 07:50 · Alexander Sacharov** — Die Schnittstelle ist unten schon in Betrieb, im Board requirementsgenie - sie muss nicht erfunden, nur uebernommen werden. Aufruf, wie er dort laeuft: setup="$root/.jaira/worktree-setup"; if [ -x "$setup" ]; then off=$(( ( $(printf '%s' "$slug" | cksum | cut -d' ' -f1) % 40 ) + 1 )); "$setup" "$wt" "$slug" "$off" "$root"; fi

Vier Dinge, die der Umzug nicht von selbst mitbringt: (1) Der Port-Versatz bleibt bei jaira - die einzige Zusicherung, die ein Projekt nicht geben kann: zwei Arbeiter duerfen nie einen Stapel teilen. jaira rechnet ihn aus dem Slug und uebergibt ihn; das Skript des Projekts verteilt ihn nur auf Namen, die es selbst kennt. (2) Bei --no-worktree darf das Skript nicht laufen: dort ist wt == root, und es wuerde in die lebende .env des Repositorys schreiben. Heute schuetzt das nur zufaellig das 'if [ ! -d "$wt" ]' darum herum. (3) Ein fehlschlagendes Skript darf den Arbeiter nicht toeten: spawn.sh hat 'set -euo pipefail', ein Exit ungleich 0 aus dem Hook bricht ab, bevor die Kachel ueberhaupt aufgeht. Abfangen und melden, nicht abbrechen. (4) 'fehlt' und 'nicht ausfuehrbar' getrennt melden, sonst sieht ein vergessenes chmod +x aus wie ein Projekt ohne Einrichtung.

Was der Umzug loest, gemessen: das Exemplar von spawn.sh in requirementsgenie war auf 159 Zeilen gewachsen - IMAGE_NS, COMPOSE_PROFILES, AZIMUTT_PORT, VITE_PORT_HOST, BACKEND_PORT_HOST, dazu ein pg_dump der Hauptdatenbank. Alles davon traegt Namen, die nur dieses eine Projekt kennt. 'roles install' ueberschreibt so eine Datei nicht, sie friert nur ein: jede Verbesserung an spawn.sh wird bei dem Anwender still uebersprungen.

Nicht projektspezifisch und daher in spawn.sh gehoerend, nicht in den Hook: --no-worktree und der Lane-Name 'dispatch'.
