---
id: 01M846QY4DXDWFZV58WSMK8HD2
title: "Cockpit-Abnahme: Abschnitte optisch trennen, damit der Text schnell lesbar ist"
status: critique
ready: true
creator: Alexander Sacharov
goal: "Die Abnahme-Ansicht (Приём) liest sich auf einen Blick: Kopf (Titel, Goal, Done, Why), History, Release note und Review-Notizen sind sichtbar getrennte Blöcke, und innerhalb des Kopfs springt jedes Feld ins Auge."
context: |-
  Was heute falsch ist: in der Abnahme-Ansicht des jaira-cockpit steht fast alles als weißer Fließtext untereinander. Goal, Done und Why im Kopf haben nur ein fettes Wort davor und verschwimmen zu einem Block; History und Release note gehen nahtlos ineinander über.
  Unten bei den Review-Notizen klappt es schon: jede Notiz hat ihre Farbmarke (🔴🟡⚪🔵) und eine Knopfzeile darunter - das liest der Maintainer gut.
  Wunsch des Maintainers (23.11.2026, mit Screenshot): Text so teilen, dass er schnell gelesen wird; 'Streifen' zwischen den Teilen; weniger Weiß, das sich vermischt; oben soll klar sein, wo Goal ist und wo was.
  Wo: core/mod/builtin/jaira-cockpit/hooks/register.tsx (Abnahme-Ansicht, Story/readStory, Kopf und History), Texte in hooks/texts.ts.
  Leser: hat ADHS, liest langsam - Ziel ist Überfliegen, nicht schöner Satz.
  Offen für den Worker: welche Mittel die Engine bietet (Trennlinien, Hintergrund, Farbe, Einrückung) - das stärkste davon nehmen, das im Terminal sicher rendert.
definition-of-done: "Zwischen Kopf, History, Release note und Review-Notizen steht je eine sichtbare Trennung (Linie oder farbiger Abschnittstitel)"
tags:
  - cockpit
  - acceptance
blocked-by: []
related: []
commits: []
created-at: 2026-11-23T12:40:23Z
updated-at: 2026-11-23T12:47:34Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-91230
claimed-at: 2026-11-23T12:40:52Z
outcome-what: "Acceptance pane: History, Release note and Review notes open with a coloured title and a dim rule; Goal/Done/Why each have a coloured bar and a blank line; history side info (round, age, commits) dim, the DoD/hunt/review line bright; release note lines one per line"
outcome-why: the pane was one white block; the maintainer skims and could not find Goal or where History ended
outcome-resolves: "MK8HD2: DoD 1-4, 6; DoD 5 is the maintainer's look at acceptance"
---

# Cockpit-Abnahme: Abschnitte optisch trennen, damit der Text schnell lesbar ist

## Definition of Done

- [x] Zwischen Kopf, History, Release note und Review-Notizen steht je eine sichtbare Trennung (Linie oder farbiger Abschnittstitel)
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx section() + tests/board.test.ts 'the head says what the ticket is for'
- [x] Im Kopf hat jedes Feld (Goal, Done, Why) eine eigene Farbe oder Marke am Zeilenanfang und eine Leerzeile Abstand, sodass man ein Feld findet, ohne zu lesen
  proof: register.tsx field(): BAR + colour per field, marginTop 1; test checks suggestion/success/warning
- [x] Langer Fließtext ist nicht mehr gleichfarbig weiß: Beschriftungen und Nebeninfos (Runde, Zeit, Commits) sind gedimmt, der Inhalt hell
  proof: register.tsx history {dim}; test: round/commits dimColor, DoD line not
- [x] Hell- und Dunkelmodus des Terminals bleiben lesbar
  proof: only theme keys (claude/suggestion/success/warning) and dimColor, no hex colours: register.tsx section()/field()
- [ ] seen by a person: der Maintainer findet in der Ansicht Goal, Done, Why und Release note mit einem Blick
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/MK8HD2.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] section header helper: coloured title + dim rule
- [x] head: mark + colour per field, blank line between
- [x] history: side info dimmed, content bright
- [x] tests + release note

## Progress
- **2026-11-23 12:40 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/6958.sock
- **2026-11-23 12:42 · Alexander Sacharov** — decision: engine has no per-side Box borders (BoxProps: borderStyle only, all four sides). Separator = section title in colour + dim '─' run cut at pane width (Text wrap=truncate-end). Head fields: coloured '▌' mark + coloured bold label per field (Goal suggestion, Done success, Why warning), marginTop 1. Theme keys only, so light/dark follow the theme. No decision left for a person.
- **2026-11-23 12:47 · Alexander Sacharov** — the release note lines are now one per line under their title, no longer joined with ' / '. Section titles lost their trailing colon (texts.ts). shipped.sha256 regenerated (gate shipped-manifest).
