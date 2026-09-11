---
id: 01M27JJ75T5PBP44D4DCPP5SCQ
title: "Auch wer nur die CLI benutzt, bekommt die Refs ohne danach zu fragen"
status: done
ready: true
creator: Alexander Sacharov
goal: "Ein Ticket, das mir jemand zuweist, erreicht mich auch dann, wenn ich das Board nie oeffne und nur lesende Kommandos benutze - ohne dass ein einziges Kommando auf das Netz wartet"
context: |-
  Gefunden beim Abnehmen von 8566KF (Frage von Alexander): 'wenn wir jaira nicht als Tool starten, kommt der fetch dann ueberhaupt nie?'

  Heute holt die Refs: das TUI alle 60 Sekunden (internal/tui/refs.go) und jede SCHREIBENDE Kommandozeile, weil refsync.Flush vor dem Senden fetcht. Wer nur in der CLI lebt und nur liest - jaira list, next, show - fetcht nie. Fuer den existiert eine Zuweisung nicht, bis er von Hand 'jaira fetch' tippt. Das ist genau die Person, fuer die das Feature gebaut wurde: ein Agent an einer bash-Zeile.

  Ausdruecklich NICHT die Loesung: den Fetch in 'jaira list' haengen. Ein Lesekommando darf nie auf ein Remote warten - das ist die Regel, die das Board schnell haelt, und Alexander hat sie beim Beantworten ausdruecklich bestaetigt ('надо подправить но не в list').

  Das Muster dafuer liegt schon zweimal im Haus und ist beide Male dasselbe: die Update-Pruefung (core/selfupdate/cache.go SpawnRefresh) und der Snapshot (core/snapshot/schedule.go). Stempel VOR dem Lauf, abgekoppeltes Kind, Rekursionsschutz per Umgebungsvariable, kein Wait, stdio auf DevNull, nie aus einem Testbinary. Ein drittes handgeschriebenes Exemplar davon waere die Stelle, an der eine dieser Feinheiten irgendwann fehlt.
definition-of-done: "nach jedem Kommando - auch einem rein lesenden - wird ein abgekoppelter 'jaira fetch' gestartet, wenn der letzte laenger her ist als das eingestellte Intervall (Default zehn Minuten, 'fetch-every-minutes' in settings.json), und kein Kommando wartet darauf; das Muster aus Stempel-vor-dem-Lauf, Rekursionsschutz, kein Wait, DevNull und Testbinary-Schutz steht an EINER Stelle, die der Snapshot-Lauf und dieser Fetch beide benutzen, statt ein drittes Mal abgeschrieben zu werden; JAIRA_NO_FETCH=1 schaltet es ab und ist zugleich der Rekursionsschutz des Kindes; Tests decken ab, dass faellig/nicht faellig richtig entschieden wird und dass ein abgeschaltetes Board nie faellig ist"
tags:
  - cli
  - concurrency
blocked-by: []
commits:
  - d879a7a2408c86429b376bc4a11906a9ebf91547
  - 3dbf02839cb83e92d9cb8f0e38eb225b454782ea
created-at: 2026-09-11T06:30:31Z
updated-at: 2026-09-11T13:31:55Z
claimed-by: DESKTOP-RFTCH11-13032
claimed-at: 2026-09-11T06:30:57Z
updated-by: Alexander Sacharov
assignee: Alexander Sacharov
question: "Fertig: nach jedem Kommando - auch einem lesenden - wird ein abgekoppelter Fetch gestartet, wenn der letzte aelter als zehn Minuten ist. list bleibt gemessen bei 0,02 s. Passt das Intervall, oder lieber fuenf Minuten?"
outcome-what: "core/bgrun als gemeinsame Stelle fuer Hintergrundlaeufe, core/refsync FetchTask, Aufruf nach jedem Kommando"
outcome-why: "Wer nur die CLI benutzt und nur liest, hat nie gefetcht - fuer den existierte eine Zuweisung nicht, bis er von Hand 'jaira fetch' tippte"
outcome-resolves: "Alle vier Punkte der DoD: der Fetch laeuft nach jedem Kommando abgekoppelt und nur wenn faellig, das Muster steht an einer Stelle die Snapshot und Fetch teilen, JAIRA_NO_FETCH=1 schaltet ab und ist zugleich der Rekursionsschutz, und fuenf Tests decken faellig/nicht faellig, den Ausschalter, einen kaputten Stempel und 'Stempel vor dem Lauf' ab"
---

# Auch wer nur die CLI benutzt, bekommt die Refs ohne danach zu fragen

## Definition of Done

- [x] nach jedem Kommando - auch einem rein lesenden - wird ein abgekoppelter 'jaira fetch' gestartet, wenn der letzte laenger her ist als das eingestellte Intervall (Default zehn Minuten, 'fetch-every-minutes' in settings.json), und kein Kommando wartet darauf; das Muster aus Stempel-vor-dem-Lauf, Rekursionsschutz, kein Wait, DevNull und Testbinary-Schutz steht an EINER Stelle, die der Snapshot-Lauf und dieser Fetch beide benutzen, statt ein drittes Mal abgeschrieben zu werden; JAIRA_NO_FETCH=1 schaltet es ab und ist zugleich der Rekursionsschutz des Kindes; Tests decken ab, dass faellig/nicht faellig richtig entschieden wird und dass ein abgeschaltetes Board nie faellig ist
  proof: core/bgrun (Stamp-vor-Lauf, Rekursionsschutz, kein Wait, DevNull, Testbinary-Schutz, Dir) - von Snapshot und Fetch gemeinsam benutzt; core/refsync/schedule.go FetchTask; internal/cli/refs.go afterCommand; 5 Tests in core/bgrun; Smoke: 'jaira list' bleibt bei 0,02 s und das zugewiesene Ticket steht beim naechsten list da, ohne je 'jaira fetch' getippt zu haben

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-11 06:34 · Alexander Sacharov** — Gebaut. Zwei Dinge:

1. Das Muster steht jetzt an einer Stelle (core/bgrun) und wird von Snapshot und Fetch benutzt, statt ein drittes Mal abgeschrieben zu werden. Genau dabei kam heraus, warum das wichtig war - siehe 2.

2. Beim Smoke-Test gefunden, und es betraf auch den schon fertigen Snapshot: ein abgekoppeltes Kind erbt das ARBEITSVERZEICHNIS des Elternprozesses, nicht dessen -C-Flag. 'jaira -C /pfad/zum/board list' hat also einen Fetch gestartet, der auf dem Verzeichnis lief, in dem die Shell zufaellig stand. Sichtbar wurde es nur, weil der Test das Ergebnis geprueft hat ('steht das Ticket jetzt da?') und nicht bloss, ob der Stempel geschrieben wurde. bgrun.Task hat dafuer jetzt ein Dir-Feld, das beide Aufrufer setzen muessen.

Nicht in 'jaira list' gehaengt, wie Alexander es verlangt hat: list bleibt bei 0,02 s, der Fetch laeuft daneben in einem anderen Prozess, und sein Ergebnis ist beim naechsten Kommando da.
