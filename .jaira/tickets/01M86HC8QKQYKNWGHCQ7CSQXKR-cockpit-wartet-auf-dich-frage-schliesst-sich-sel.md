---
id: 01M86HC8QKQYKNWGHCQ7CSQXKR
title: "Cockpit 'Wartet auf dich': Frage schliesst sich selbst oder per x, zeigt ihre Zeit und uebersetzt nur auf Knopfdruck"
status: critique
ready: true
creator: Sascha
assignee: Sascha
goal: "Im Pane 'Wartet auf dich' stehen nur Fragen, die wirklich noch offen sind; jede zeigt, wann sie gestellt wurde; sie erscheint in der Originalsprache und wird nur auf Knopfdruck uebersetzt."
context: |-
  Was falsch ist (Pane aus 1M3X77, core/mod/builtin/jaira-cockpit/hooks/threads.ts + register.tsx, Texte texts.ts:267-274):
  1) Eine Frage verschwindet nur, wenn die Person im Pane eine Option waehlt oder der Klassifizierer eine Chat-Antwort erkennt. Wer im Cockpit per Accept/Give back entschieden hat, behaelt die Frage fuer immer in der Liste. Es gibt keinen Weg, sie zu schliessen.
  2) Es fehlt die Zeit: die Person liest eine Frage und weiss nicht, zu welcher Nachricht sie gehoert.
  3) Das Pane uebersetzt immer (settings language). Sascha (Alex) 24.11.2026: 'Knopf zum Uebersetzen, nicht immer; ich brauche das nicht immer'. Vorbild: FMWBA8 (Abnahme-Pane, Taste t / Translate, ein Aufruf fuer alles, nochmal druecken = Original).
  Ausloeser: Sascha 24.11.2026 im Teamlead-Chat.
definition-of-done: "Eine Frage zu einem Ticket verschwindet von selbst, wenn das Ticket die Lane verlassen hat, in der sie gestellt wurde (z.B. signoff -> done oder zurueck); Test"
tags:
  - cockpit
blocked-by: []
related: []
follows: 01M8537VMKP5N8TJ2YC91M3X77
commits: []
created-at: 2026-11-24T10:24:44Z
updated-at: 2026-11-24T10:44:08Z
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-69066
claimed-at: 2026-11-24T10:25:16Z
outcome-what: "Cockpit 'Waiting for you': x settles a question; a question settles itself once its ticket leaves the lane it was asked in; settled ones stay 24h dimmed under 'Answered'; asked/answered times; questions kept in the original, t translates all at once"
outcome-why: "Questions decided in the Acceptance pane stayed forever, had no time, and were always shown rephrased (Sascha 24.11.2026)"
outcome-resolves: CSQXKR
review-summary: "core/mod/builtin/jaira-cockpit/hooks/threads.ts settled() keeps the options of a settled question, which the Answered section never shows, so Translate sends them to the model and they bust its cache - drop them there (options: []) | core/mod/builtin/jaira-cockpit/hooks/register.tsx ASKS pane casts e.props for bodyColumns; the Overview pane in the same file reads e.props.bodyColumns directly - do the same"
round: 2
back-from: critique
back-by: lane
back-count: 1
---

# Cockpit 'Wartet auf dich': Frage schliesst sich selbst oder per x, zeigt ihre Zeit und uebersetzt nur auf Knopfdruck

## Definition of Done

- [x] Eine Frage zu einem Ticket gilt von selbst als erledigt, wenn das Ticket die Lane verlassen hat, in der sie gestellt wurde (z.B. signoff -> done oder zurueck); Test
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts 'settling questions, their times and Translate (CSQXKR) > a question on a ticket goes to Answered once the ticket leaves the lane it was asked in' + 'lane after: *' (threads.ts moved, register.tsx settleMoved)
- [x] Jede offene Zeile hat x bzw. ein Klick-Ziel 'erledigt', das sie ohne Antwort in den Chat als erledigt markiert; Test
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts 'settling questions, their times and Translate (CSQXKR) > x settles a question without a word into the session; ...' (register.tsx dismissAsk, key ask-N-done)
- [x] Jede Frage zeigt Datum und Uhrzeit, zu der sie gestellt wurde; eine erledigte zusaetzlich, wann sie erledigt wurde (z.B. 'gefragt 07:12 · beantwortet 07:40'); Test
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts 'settling questions, their times and Translate (CSQXKR) > times: today by the clock, another day with its date' + x test matches 'asked today HH:MM' + lane test 'asked today .. · answered today ..'
- [x] Das Pane zeigt die Fragen im Original; Taste t bzw. ein Knopf Translate uebersetzt alle auf einmal, nochmal druecken zeigt wieder das Original; Test
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts 'settling questions, their times and Translate (CSQXKR) > the pane opens on the original; t translates every question at once, t again shows the original' + 'without a language set there is no Translate'
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/CSQXKR.md
- [x] Erledigte Fragen (Antwort im Pane, im Chat, per x oder automatisch) verschwinden nicht sofort: sie stehen unter einer Trennlinie in einem eigenen Abschnitt 'Beantwortet' unter den offenen, gedimmt, und fallen nach 24 Stunden heraus; offene stehen oben; Test
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts 'settling questions, their times and Translate (CSQXKR) > what a change settled carries both times; a day later it is gone' + 'settled questions are kept across a restart, and one settled a day ago is no longer shown' + 1M3X77 'the band opens the pane...' (Answered after a pick)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-24 10:25 · Sascha** — parent session: uds:/run/user/1000/cc-socks/3417.sock
- **2026-11-24 10:26 · Sascha** — decision: Sascha 24.11.2026: beide Wege (automatisch nach Lane-Wechsel UND manuell per x). Nicht sofort entfernen: erledigte Fragen einen Tag lang unter einer Trennlinie als 'Beantwortet' zeigen, damit klar ist, was neu und was schon beantwortet ist. Zeit der Frage und Zeit der Erledigung anzeigen.
- **2026-11-24 10:26 · Sascha** — decisions (worker, no open question): 1) lane at asking time = status from the same 'jaira show' that looks up the title, only for the ticket the question names (about); the thread's own ticket is not looked up (as today). Every board refresh (60s) compares with 'jaira list' of all boards; a ticket in none of them is asked once with 'jaira show'; unknown or other lane -> done. 2) closed = any open question that is gone after a change (pick, chat filing, x, lane left, a new question on the same thread replacing it) - goes to 'Answered' with answeredAt, dropped after 24h. 3) times in English, cockpit is English only (NZZDRH): 'asked today 07:12 · answered today 07:40', other day '23.11. 21:40'. 4) question is now kept in the assistant's own wording (classifier prompt), Translate = one haiku call for every question+option shown, into settings language; no language set = no button (as FMWBA8); a pick sends the original option text.
- **2026-11-24 10:34 · Sascha** — waiting: cockpit test run (background b0ym9m1cy) since 11:34; code for DoD 1-4,6 written, release note written, shipped.sha256 still to regenerate (go test ./core/role fails until then)
- **2026-11-24 10:40 · Sascha** — in-progress: the pane no longer closes after the last pick (it would hide 'Answered'; the band button only appears while a question is open, so a closed pane could not show it). shipped.sha256 must be regenerated (scripts/shipped-history.sh --write) after every change to the cockpit files, or go test ./core/role fails. Known limit: a question whose ticket is only the thread's own (no 'about') does not settle by lane - its lane is never read (as its title is not).
- **2026-11-24 10:41 · Sascha** — critique: settled() drops options (never shown under Answered, only cost Translate tokens); ASKS pane reads e.props.bodyColumns without the cast, as the Overview pane does. Rest stands: separate asksShown atom (questions are no ticket, the rephrased/translating atoms are keyed by ticket), settling by replacement (decision note above).
