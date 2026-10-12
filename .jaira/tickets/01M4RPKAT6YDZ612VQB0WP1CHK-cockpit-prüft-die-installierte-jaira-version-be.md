---
id: 01M4RPKAT6YDZ612VQB0WP1CHK
title: "Cockpit prüft die installierte jaira-Version, bevor es sich auf sie verlässt"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Das jaira-cockpit erkennt eine zu alte jaira-Binary und sagt es, statt still falsche Spuren auf dem Ticket zu hinterlassen"
context: |-
  Give back im Cockpit schreibt eine Notiz aufs Ticket, die nicht stimmen kann, wenn jaira zu alt ist.
  - Seit QRHB35 schiebt giveBack (core/mod/builtin/jaira-cockpit/hooks/register.tsx, ~Zeile 443-446) ein Ticket aus signoff ohne --force nach in-progress zurück.
  - Das geht erst ab jaira 0.3.7 (Q484YS). Ältere jaira verweigert den move.
  - Annahme in QRHB35: das Cockpit ist per go:embed in der Binary, also ist jaira immer >= 0.3.7.
  - Die Annahme ist falsch: .claude-plugin/marketplace.json installiert jaira-cockpit ohne jaira. Dann läuft es neben einer beliebig alten jaira.
  - Folge: giveBack schreibt zuerst 'Given back from signoff: <grund>' per jaira note, dann scheitert der move. Toast meldet den Fehler, Ticket bleibt in signoff, aber die Notiz behauptet das Gegenteil.
  - Gefunden als review-gap 1 von QRHB35 (siehe dessen review-Feld).
  - Plugin-Version steht in core/mod/builtin/jaira-cockpit/.claude-plugin/plugin.json (0.2.0); eine Mindestversion von jaira gibt es nirgends.
  - Offen für brainstorm: nur prüfen und warnen (z. B. 'jaira --version' beim Start), oder Funktionen je nach Version abschalten.
definition-of-done: Das Cockpit liest die installierte jaira-Version und kennt eine Mindestversion an genau einer Stelle; Test
tags: []
blocked-by: []
related: []
follows: 01M4RBYV2C8F1BPQ70XCQRHB35
commits: []
created-at: 2026-10-12T18:39:31Z
updated-at: 2026-10-12T18:52:45Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-47494
claimed-at: 2026-10-12T18:49:34Z
outcome-what: "jaira-cockpit reads 'jaira --version' once when it finds the board and warns (toast + band) when it is older than MIN_JAIRA 0.3.7; Give back moves first and notes only after the move succeeded; README names the minimum for the marketplace install; plugin 0.2.1"
outcome-why: "Installed from the marketplace the cockpit runs beside any jaira; an older one refuses Give back's move without --force, and the note written before it claimed the ticket went back"
outcome-resolves: review gap 1 of QRHB35
---

# Cockpit prüft die installierte jaira-Version, bevor es sich auf sie verlässt

## Definition of Done

- [x] Das Cockpit liest die installierte jaira-Version und kennt eine Mindestversion an genau einer Stelle; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts:237 MIN_JAIRA (only place), versionOf/isOlder :242/:247; core/mod/builtin/jaira-cockpit/hooks/register.tsx:285 checkVersion runs 'jaira --version' once when the board is found; tests 'the jaira it runs beside' in core/mod/builtin/jaira-cockpit/tests/board.test.ts (claude plugin test: 87 pass, exit 0)
- [x] Ist jaira älter als die Mindestversion, sagt das Cockpit es sichtbar (Toast/Band) mit installierter und benötigter Version; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:288 say(t.tooOld(installed, MIN_JAIRA)) = toast + band notice; text core/mod/builtin/jaira-cockpit/hooks/texts.ts:19; test 'one older than the cockpit needs is said once, with the version installed and the one needed', and 'jaira 0.3.7/0.10.0/dev is not warned about'
- [x] giveBack schreibt keine Notiz, wenn der move nicht möglich ist: erst Version prüfen bzw. move, Notiz nur bei Erfolg; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:460-462 move first, note only after it succeeded; tests 'a refused move writes no note: the ticket stays where it is and the refusal is shown' and the move-before-note order in 'notes why and moves the ticket back to in-progress…'
- [x] README/Install-Hinweis des Marketplace-Wegs nennt die Mindestversion von jaira
  proof: README.md:739-743 under '/plugin install jaira-cockpit --marketplace BeMuCa/jaira': needs jaira 0.3.7 or later
- [x] Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/WP1CHK.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-12 18:49 · Alexander Sacharov** — parent session: jaira-f5
- **2026-10-12 18:49 · Alexander Sacharov** — Teamlead 2026-10-12: Branch feat/wp1chk (Worktree /home/alex/projects/.worktrees/jaira-wp1chk) ist von feat/qrhb35 (82aa5654) abgezweigt, nicht von master: der Cockpit-Code mit workers()/giveBack ohne --force existiert nur dort. feat/qrhb35 ist noch nicht gepusht/gemergt; der PR dieses Tickets stapelt sich darauf. Parallel laufen WP1CHK (giveBack/Version), GPFA7B (accept), TWB0R6 (Band-Render) - alle in hooks/register.tsx; nur den eigenen Bereich anfassen.
- **2026-10-12 18:50 · Alexander Sacharov** — Keine offenen Entscheidungen; selbst entschieden: (1) nur prüfen und warnen, keine Funktionen abschalten — die DoD verlangt nur Warnung + sichere Reihenfolge. (2) Mindestversion als Konstante MIN_JAIRA in hooks/board.ts, gelesen per 'jaira --version' (Ausgabe 'jaira version X.Y.Z'). (3) Eine Version, die keine X.Y.Z ist (Dev-Build 'dev'), gilt als unbekannt und wird nicht gewarnt — sonst warnt jeder go-run-Build. (4) giveBack: erst move, dann Notiz — die Reihenfolge allein macht die Notiz wahr, unabhängig von der Version. (5) Warnung einmal pro Session beim Finden des Boards, per say() (Toast + notice im Band). (6) README: Abschnitt '/plugin install jaira-cockpit --marketplace' nennt jaira >= 0.3.7.
- **2026-10-12 18:52 · Alexander Sacharov** — plugin.json 0.2.0 -> 0.2.1: marketplace installs only see an update when the plugin version changes, and this fix is for exactly those installs.
