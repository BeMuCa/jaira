---
id: 01M67K70X843TBF0TYRPMXMDFH
title: "Die Release-Notes erscheinen als Popup statt Seite, getrennt in Features zuerst und 'Andere' darunter, luftig und leicht lesbar"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Nach einem Upgrade zeigt ein Popup die Aenderungen: oben Features und Feature-Aenderungen, darunter unter der Ueberschrift 'Andere' der Rest, mit Abstaenden, freundlich und schnell zu erfassen."
context: |-
  Wunsch von berkc am 31.10.2026 (Punkt 6): 'I wanted the release notes to be a pop up. You made them a page. check what i wrote in the ticket again and redo. i wanted pop up only and a split between features (or feature changes) and others. Also make it more fun and easily readable. You can use space and gaps ... Features and feature changes are the most relevant, so they come first and then a header with others.'
  Heute (0.5.7, GK1QGT): eine Version pro Seite, links/rechts blaettert; 'w' im Launcher. Lies GK1QGT und seine Notizen, was berkc urspruenglich verlangt hat.
  Die Notizdateien core/release/notes/<ID>.md haben heute nur '- ' und '- ! ' (Handlungsbedarf); eine Kennzeichnung Feature vs. anderes fehlt - braucht ein Format (z.B. '- + ' fuer Feature) und die Regeln in CLAUDE.md.
definition-of-done: "Popup (Overlay ueber dem Board), keine eigene Seite; links/rechts blaettert weiter zwischen Versionen"
tags:
  - tui
blocked-by: []
related:
  - 01M5659BR7RYZ6RVC84HGK1QGT
commits: []
created-at: 2026-10-30T23:44:40Z
updated-at: 2026-10-30T23:48:37Z
updated-by: berkc
claimed-by: berk-ca-2691736
claimed-at: 2026-10-30T23:47:40Z
---

# Die Release-Notes erscheinen als Popup statt Seite, getrennt in Features zuerst und 'Andere' darunter, luftig und leicht lesbar

## Definition of Done

- [ ] Popup (Overlay ueber dem Board), keine eigene Seite; links/rechts blaettert weiter zwischen Versionen
- [ ] Je Version zuerst Features/Feature-Aenderungen, dann Ueberschrift 'Andere'; Leerzeilen und Einrueckung fuer Lesbarkeit; Test
- [ ] Notizformat kennzeichnet Features; CLAUDE.md-Abschnitt zu Release-Notes beschreibt es
- [ ] Release-Note-Zeile

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-30 23:48 · berkc** — Open decisions, counted once (before the plan). 1 feature marker: '- + text' for a feature or feature change, combinable with the must-act marker as '- + ! text'. Settled: mirrors '- ! ' (same parser style, ASCII, a plain line to old binaries which show '+ text'). Rejected: a word prefix such as 'Feature:' (a sentence starting with that word would be misread); a separate file per kind (breaks one file per ticket); a unicode star (old binaries and terminals render it unevenly).
- **2026-10-30 23:48 · berkc** — 2 section header labels: 'Features' above the features, 'Others' above the rest. English, because the popup title and the hint are English ('What is new since', '←/→ older/newer'). The ticket's 'Andere' in its DoD is superseded by the person's own wording 'others'. 3 the 'Others' header is printed only when the version has at least one feature line: a version with only plain notes (every version before this change) looks exactly as it did. 4 a must-act feature line sits under Features and keeps its warning colour; a must-act plain line sits under Others as before. 5 paging stays as GK1QGT decided: the popup pages only unread versions, 'w' pages all, newest first, left/h older, right/l newer. Settled, not reopened.
