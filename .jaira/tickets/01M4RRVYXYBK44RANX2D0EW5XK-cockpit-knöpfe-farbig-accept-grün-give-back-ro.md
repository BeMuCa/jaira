---
id: 01M4RRVYXYBK44RANX2D0EW5XK
title: "Cockpit-Knöpfe farbig: Accept grün, Give back rot, Explain gelb"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die Knöpfe im Cockpit-Band sind auf einen Blick unterscheidbar: Accept grün, Give back rot, Explain gelb"
context: |-
  Im Cockpit-Band sehen Accept, Give back und Explain gleich aus; man muss lesen, um den richtigen zu treffen.
  - Wunsch von Alex 2026-10-12: Accept grün, Give back rot, Explain gelb.
  - Code: Knöpfe in core/mod/builtin/jaira-cockpit/hooks/register.tsx (~Zeile 856-868, Branch feat/13pqb3). Accept hat heute variant="primary" (blau hinterlegt, siehe Screenshots), die anderen keine variant.
  - Unbekannt: ob <Button> der Mod-Engine eine Farbe annimmt (color-Prop oder weitere variants wie success/danger/warning). Erst die Engine-Typen prüfen; gibt es keine Farbe, die Beschriftung in einem farbigen <Text> rendern oder melden, was geht.
  - Die Farben sind Theme-Farben (success/error/warning wie bei 'threads:' in gelb), keine festen Hex-Werte, damit helles und dunkles Terminal-Theme funktionieren.
  - Read explanation und Verdict bleiben neutral.
definition-of-done: "Accept ist grün, Give back rot, Explain gelb im Band, über Theme-Farben der Engine; Test"
tags: []
blocked-by: []
related: []
commits: []
created-at: 2026-10-12T19:19:11Z
updated-at: 2026-10-12T19:25:07Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-41108
claimed-at: 2026-10-12T19:19:37Z
outcome-what: "Accept, Give back und Explain im Cockpit-Band zeichnen ihre Beschriftung als <Text> in den Theme-Farben success, error, warning; Read explanation und Verdict bleiben neutral"
outcome-why: "Die Knöpfe sahen gleich aus; man musste lesen, um den richtigen zu treffen"
outcome-resolves: "Accept grün, Give back rot, Explain gelb über Theme-Farben; Fokus bleibt durch die Inversion der Engine erkennbar; Release-Note 0EW5XK.md"
---

# Cockpit-Knöpfe farbig: Accept grün, Give back rot, Explain gelb

## Definition of Done

- [x] Accept ist grün, Give back rot, Explain gelb im Band, über Theme-Farben der Engine; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:858-864; test 'waiting for you > Accept, Give back and Explain are told apart by the theme's colours; Read explanation and Verdict stay plain' (tests/board.test.ts:239)
- [x] Der fokussierte Knopf bleibt erkennbar (Fokus-Ring/Hervorhebung geht durch die Farbe nicht verloren)
  proof: register.tsx:858-864 set no plain/inverse; the engine marks focus and pointer with inverse whatever the Text colour (ButtonProps.plain/variant docs and ui.focus 'the engine's inverse marks it', claude-code.d.ts 2.1.296)
- [x] Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/0EW5XK.md:1

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-12 19:19 · Alexander Sacharov** — Teamlead 2026-10-12: Branch feat/0ew5xk ist von feat/13pqb3 abgezweigt (Kette QRHB35 -> TWB0R6 -> 13PQB3): 13PQB3 hat die Knopfreihe (Verdict-Knopf, Explain-Regel) umgebaut. Nur die Knöpfe einfärben; GPFA7B arbeitet parallel an accept()/afterAccept in derselben Datei.
- **2026-10-12 19:19 · Alexander Sacharov** — parent session: jaira-f5
- **2026-10-12 19:23 · Alexander Sacharov** — Keine offenen Entscheidungen. Engine 2.1.296: ButtonProps hat keine color-Prop, variant nur 'primary'|'secondary'. Weg: Button behält label (Adresse/Name) und bekommt ein <Text color=success|error|warning> als Kind — erlaubt laut ButtonProps.label-Doku. Fokus: die Engine markiert Fokus und Pointer per inverse, unabhängig von der Textfarbe; wir setzen kein plain/inverse. Accept behält variant=primary (Desktop zeichnet so den Hauptknopf; im Terminal bleiben die Klammern in Akzentfarbe, die Beschriftung grün).
- **2026-10-12 19:25 · Alexander Sacharov** — DoD 2 ist aus der Engine-Doku belegt, nicht live gesehen: die Test-API zeigt keinen Fokus-Zustand. Im Terminal prüfen: ctrl+x tab ins Band, Fokus auf Accept — der Knopf muss invertiert erscheinen.
