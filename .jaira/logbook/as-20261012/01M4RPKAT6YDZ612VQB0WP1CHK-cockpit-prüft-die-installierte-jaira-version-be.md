---
id: 01M4RPKAT6YDZ612VQB0WP1CHK
title: "Cockpit prüft die installierte jaira-Version, bevor es sich auf sie verlässt"
status: done
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
commits:
  - 328f00e002da3c09b6c57ceecb0e48ddf8185307
  - 38e03dfc89d15653f32d569fc51a07381d64b454
  - efda32015253401fd1dc881f9d31f04e98ca5185
created-at: 2026-10-12T18:39:31Z
updated-at: 2026-10-12T19:39:08Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-92968
claimed-at: 2026-10-12T19:31:57Z
outcome-what: "jaira-cockpit reads 'jaira --version' once when it finds the board and warns (toast + band) when it is older than MIN_JAIRA 0.3.7; Give back moves first and notes only after the move succeeded, and a note that fails after the move opens a pane with jaira's reason and the lost text as a ready 'jaira note' line while the handover goes on; README names the minimum for the marketplace install; plugin 0.2.1"
outcome-why: "Installed from the marketplace the cockpit runs beside any jaira; an older one refuses Give back's move without --force, and the note written before it claimed the ticket went back. After the reorder a failing note was silent, losing the one thing the person typed"
outcome-resolves: review gap 1 of QRHB35; review gap 1 of WP1CHK
review-summary: "When the cockpit finds the board it runs 'jaira --version' once per session. If the version is older than MIN_JAIRA 0.3.7 (hooks/board.ts, the only place it is written down), it shows a toast and a band notice naming both versions. A 'dev' build and any version that is not X.Y.Z get no warning. Give back now moves the ticket first and writes the 'Given back from <lane>: <reason>' note only after the move succeeded, so a refused move leaves no note. If that note fails after the move, a 'Note on <id> failed' pane opens with jaira's reason and a ready, shell-quoted 'jaira note' line; the handover to the worker still goes ahead. Also in the change: a README paragraph for the marketplace install, plugin.json 0.2.1, two release-note lines and 9 engine tests."
review-gaps: "Two small gaps, neither blocks acceptance. (1) The too-old warning is said once, through say(). That puts it in noticeAtom, and the next say() or fail() overwrites it: the first Accept or Give back clears the band. After that only the refused Give back itself explains the problem, and its pane does show jaira's own refusal. (2) versionOf reads '0.3.7-rc.1' as 0.3.7, so a pre-release of the minimum version gets no warning. The DoD is fully covered by the diff: MIN_JAIRA in one place, the warning with both versions, move before note, a pane for the failed note, the README and the release notes. Every claim in the outcome is backed by the diff. claude plugin test: 88 pass, 0 fail."
test-verdict: "pass: claude plugin test 88 pass, exit 0; the new test fails (87/1) with the fail() line disabled; the pasted 'jaira note' line keeps a ' and $ intact through bash; DoD 1-6 proofs point at existing lines. Not run live with a failing 'jaira note' in a Claude Code session"
review-verdict: "Ready to accept. The diff satisfies all 6 DoD items, with no defects found. The two gaps are cosmetic: a band notice that can be overwritten, and pre-release versions."
review-check: "1. cd core/mod/builtin/jaira-cockpit && claude plugin test → '88 pass, 0 fail'. 2. The tests named 'the jaira it runs beside' check the warning for 0.3.6, that 0.3.7, 0.10.0 and dev get no warning, and that jaira is not asked without a board. 3. 'a refused move writes no note' and 'a note that fails after the move is shown…' cover Give back. 4. By hand, with an old binary: put jaira 0.3.6 first on PATH, open Claude Code in this repo → a toast 'jaira 0.3.6 is too old for this cockpit: it needs 0.3.7 or later…'. 5. With that binary press Give back on a signoff ticket → a 'Give back <id> failed' pane, and 'jaira show <id>' has no 'Given back from' note. 6. Read README.md:740-743 → it names 'jaira 0.3.7 or later'."
---

# Cockpit prüft die installierte jaira-Version, bevor es sich auf sie verlässt

## Definition of Done

- [x] Das Cockpit liest die installierte jaira-Version und kennt eine Mindestversion an genau einer Stelle; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts:237 MIN_JAIRA (only place), versionOf/isOlder :242/:247; core/mod/builtin/jaira-cockpit/hooks/register.tsx:285 checkVersion runs 'jaira --version' once when the board is found; tests 'the jaira it runs beside' in core/mod/builtin/jaira-cockpit/tests/board.test.ts (claude plugin test: 87 pass, exit 0)
- [x] Ist jaira älter als die Mindestversion, sagt das Cockpit es sichtbar (Toast/Band) mit installierter und benötigter Version; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:288 say(t.tooOld(installed, MIN_JAIRA)) = toast + band notice; text core/mod/builtin/jaira-cockpit/hooks/texts.ts:19; test 'one older than the cockpit needs is said once, with the version installed and the one needed', and 'jaira 0.3.7/0.10.0/dev is not warned about'
- [x] giveBack schreibt keine Notiz, wenn der move nicht möglich ist: erst Version prüfen bzw. move, Notiz nur bei Erfolg; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:460-467 move first, note only after it succeeded; tests 'a refused move writes no note: the ticket stays where it is and the refusal is shown' and the move-before-note order in 'notes why and moves the ticket back to in-progress…'
- [x] README/Install-Hinweis des Marketplace-Wegs nennt die Mindestversion von jaira
  proof: README.md:739-743 under '/plugin install jaira-cockpit --marketplace BeMuCa/jaira': needs jaira 0.3.7 or later
- [x] Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/WP1CHK.md (two lines)
- [x] Scheitert in giveBack die Note nach einem erfolgreichen move, zeigt das Cockpit es sichtbar (fail-Pane mit Grund und dem verlorenen Text); Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:467 fail(t.noteOn, t.noteLost) when the note after the move fails, text hooks/texts.ts:19-20; test 'a note that fails after the move is shown with the reason typed, and the worker still gets the ticket' tests/board.test.ts:605 (claude plugin test: 88 pass, exit 0)

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
- **2026-10-12 18:53 · Alexander Sacharov** — critique round 1: core/mod/builtin/jaira-cockpit/hooks/board.ts:243 versionOf guards states that cannot occur: goreleaser sets main.version={{.Version}} without a leading v (.goreleaser.yaml:22), and 'jaira --version' prints no board warning (checked: plain 'jaira version 0.3.7' on this outdated board) — drop 'v?' and the withoutWarnings() call, keep /\bversion\s+(\d+\.\d+\.\d+)\b/. Left alone: checkVersion in refreshBoard rather than findBoard (findBoard only answers whether a board is there; the warning is a side effect that belongs with the other say() calls); the test mock's /^0\.3\.[0-6]$/ stands in for an old jaira and need not reuse isOlder.
- **2026-10-12 18:53 · Alexander Sacharov** — critique round 2: round 1 finding repaired (board.ts:243, no v?, no withoutWarnings); nothing left.
- **2026-10-12 18:54 · Alexander Sacharov** — testing: not exercised live in a Claude Code session with an old jaira (no 0.3.6 binary at hand); the engine tests drive the same hooks with 'jaira version 0.3.6' and a move refused without --force.
- **2026-10-12 18:56 · Alexander Sacharov** — review: accept-worthy. One small gap: register.tsx:462 ignores the result of the 'Given back from…' note after a successful move, so a failing note is silent; fix by fail() on !ok after the move, or leave knowingly.
- **2026-10-12 18:57 · Alexander Sacharov** — worker: stopped at signoff. Put to parent jaira-f5: review gap 1 (register.tsx:462 ignores a failed note after the move) — 1) fix now: back to in-progress, fail() on a failed note plus a test (recommended); 2) leave it knowingly and accept with 'jaira move WP1CHK --to done --force'.
- **2026-10-12 19:28 · Alexander Sacharov** — explain: Befund: Give back verschiebt das Ticket jetzt zuerst und schreibt danach die Notiz mit dem Grund; schlägt diese Notiz fehl, merkt es niemand (register.tsx:462 ignoriert das Ergebnis von 'jaira note').
Warum: Der Grund ist das Einzige, was die Person beim Zurückgeben selbst eintippt; geht er verloren, steht das Ticket in in-progress und keiner weiß, was zu tun ist. Der Fall ist selten, weil der move im selben Verzeichnis gerade geklappt hat.
a) Jetzt reparieren: zurück nach in-progress, bei fehlgeschlagener Notiz fail() mit Begründung zeigen, dazu ein Test (ca. 15 Minuten).
b) Bewusst so lassen und annehmen: jaira move WP1CHK --to done --force.
Empfehlung: a - die Reparatur ist klein und schützt genau den Text, den ein Mensch geschrieben hat.
- **2026-10-12 19:31 · Alexander Sacharov** — decision: a) Jetzt reparieren: zurück nach in-progress, bei fehlgeschlagener Notiz fail() mit Begründung zeigen, dazu ein Test (ca. 15 Minuten).
- **2026-10-12 19:31 · Alexander Sacharov** — Given back from signoff: repair
- **2026-10-12 19:31 · Alexander Sacharov** — Alex 2026-10-12: Review-Gap 1 jetzt beheben (register.tsx ~462, fail() + Test). Dann wieder signoff.
- **2026-10-12 19:31 · Alexander Sacharov** — cockpit: given back to a new worker - old worker's context 125k >= 100k (50% of its window)
- **2026-10-12 19:32 · Alexander Sacharov** — worker (new, after give back): DoD 6 settled — a failed note after a successful move opens the fail pane (title 'Note on <id> failed') with jaira's reason and the lost text as a ready 'jaira note' line, and stops there: no worker is handed the ticket, because a new worker would read a ticket without the reason. The person re-adds the note and presses nothing else; the ticket already sits in in-progress.
- **2026-10-12 19:32 · Alexander Sacharov** — worker: revised the DoD 6 decision above — the handover goes on after the fail pane. Stopping would leave the ticket in in-progress with no worker, and the cockpit has no button to hand it out from there (Give back is offered only on signoff/human). The pane carries the lost text as a ready 'jaira note' line; the old worker gets the reason in its typed message anyway.
- **2026-10-12 19:34 · Alexander Sacharov** — critique round 3 (repair for review gap 1, efda3201 only): nothing. Let stand: backMessage (hooks/sessions.ts:161) still tells the old worker the reason is 'also a note on the ticket' when that note failed — the worker has the reason in the message itself, and the pane tells the person; t.noteOn/t.noteLost follow texts.ts holding every person-facing string.
- **2026-10-12 19:39 · Alexander Sacharov** — Alex 2026-10-12: angenommen; die zwei kosmetischen Review-Gaps (Warnung vom nächsten say() überschrieben, 0.3.7-rc.1 = 0.3.7) bleiben bewusst stehen.
