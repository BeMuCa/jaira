---
id: 01M263PXWYJ9F6TG6EQSZWK0PF
title: "Ein Ticket mit offener Frage faellt auf, egal in welcher Lane es liegt"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Eine gestellte, unbeantwortete Frage ist auf dem Board sichtbar, auch wenn das Ticket nicht in der human-Lane sitzt"
context: |-
  Auf dem requirementsgenie-Board tragen 9GEGTB, NEEY9A und SXDBX6 ein gefuelltes question-Feld (931, 1059 und 19 Zeichen) und liegen in review, nicht in human.
  9GEGTB fragt woertlich 'Do I delete the stray {clarify_response} from prompts/v1_powerpoint.txt?' - eine Frage, die vor der Abnahme beantwortet gehoert.
  In review ist sie unsichtbar: die Lane heisst 'Review', nicht 'wartet auf deine Antwort', und die Karte zeigt sie nicht.
  jaira kann das schon anzeigen - internal/tui/view.go:1026 rendert bei einem question-Feld das Kennzeichen 'waiting on your answer'. Vor dem Bauen pruefen, ob das nur in der Detailansicht haengt oder auch auf der Karte, und warum es hier nicht greift.
definition-of-done: "Ein Ticket mit gefuelltem question-Feld ist auf der Board-Karte als solches erkennbar, unabhaengig von seiner Lane; ein Test deckt eine Frage ausserhalb der human-Lane ab; go test ./... -race gruen"
tags:
  - cli
blocked-by: []
commits: []
created-at: 2026-09-10T16:51:42Z
updated-at: 2026-10-13T14:21:12Z
updated-by: BeMuCa
---

# Ein Ticket mit offener Frage faellt auf, egal in welcher Lane es liegt

## Definition of Done

- [ ] Ein Ticket mit gefuelltem question-Feld ist auf der Board-Karte als solches erkennbar, unabhaengig von seiner Lane; ein Test deckt eine Frage ausserhalb der human-Lane ab; go test ./... -race gruen
- [ ] Verlaesst ein Ticket die human-Lane, leert jaira das question-Feld; Frage und Antwort bleiben als Notiz im Ticket lesbar (wie die Antwort hineinkommt, z.B. --answer beim move, entscheidet brainstorm); damit heisst ein gefuelltes question-Feld wieder 'unbeantwortet'; ein Test deckt den Weg aus human ab

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-13 13:49 · BeMuCa** — 2026-10-13, Berk und Claude: Der Fragetext steht im Detail weit unten (internal/tui/view.go:1354, nach Ziel, Kontext und DoD); oben steht nur das Kennzeichen '▲ waiting on your answer' an der Lane-Zeile (view.go:1255), und das nur in einer requires-question-Lane. Auf der Karte steht nichts. Berk sah PNC854 in human und fragte 'da ist doch gar keine Frage?'. Mit diesem Ticket mitnehmen: der Fragetext gehoert oben ins Detail. Achtung fuer die DoD hier: PNC854 liegt jetzt in signoff mit einer laengst beantworteten Frage im question-Feld - 'gefuelltes question-Feld' ist also nicht dasselbe wie 'unbeantwortet'.
- **2026-10-13 14:21 · BeMuCa** — Berk 2026-10-13: Variante A gewaehlt. Beim Verlassen von human wird das question-Feld geleert, Frage und Antwort bleiben als Notiz. Grund: heute bleibt eine beantwortete Frage im Feld stehen (PNC854 liegt in signoff mit seiner laengst beantworteten Frage), und dann wuerde dieses Ticket eine beantwortete Frage als offen auf die Karte malen. Verworfen: Frage nur in human zeigen (widerspricht dem Ziel hier, Fragen ausserhalb von human zu sehen). Naehe zu B4MGTP (Frage und Antwort als Paar) beachten, nicht vorwegnehmen.
