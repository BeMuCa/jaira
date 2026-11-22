---
id: 01M81VFQRXRFE870G9010GMYQB
title: "Das Verdict zeigt jede Notiz mit Ampel-Marker und in der Sprache der Person, verständlich umformuliert"
status: optimize
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Alex sieht im Verdict jede Notiz als eine kurze Zeile auf Russisch mit Marker und weiß ohne Nachlesen, welche blockiert."
context: |-
  Problem: im Verdict des Cockpits stehen die Notizen (review-gaps) als englischer Fließtext ohne Einordnung. Alex liest langsam und auf Russisch; er muss jede Notiz lesen, um zu sehen, ob sie blockiert.
  - Wunsch Alex 22.11.2026: "können wir direkt die Notes mit einer Ampel markieren? und ich möchte sie in meiner Sprache bekommen, auch wenn sie im Ticket englisch stehen; wenn wir sie eh durch ein LLM schicken, dann gleich richtig, schön und verständlich".
  - Vorbild: die Signoff-Berichte des Teamleads im Chat. Jede Notiz eine Zeile mit Marker: 🔴 Blocker, 🟡 kleiner Rand (eine Runde ~10 min), ⚪ bekannte Grenze, 🔵 nicht dieses Ticket. Alex findet das sofort klar.
  - AJ41SC bringt die Ampel für das ganze Ticket (erste Zeile im Verdict); dieses Ticket macht dasselbe pro Notiz.
  - Heute: texts.ts:1 sagt, das Cockpit zeigt alles auf Englisch. Explain/Discuss-Chats antworten schon in der Sprache der Person (texts.ts:138), das Verdict-Pane nicht.
  - Das Ticket bleibt englisch/deutsch; übersetzt wird nur die Anzeige.
  - Modell: .jaira/chat-model (heute sonnet) oder billiger; das Ergebnis pro Ticket und Notizen-Stand cachen, damit nicht jeder Blick ein LLM-Aufruf ist.
  - Offen für pre-process: woher die Sprache der Person kommt (Einstellung in ~/.jaira/settings.json oder die Sprache der Session).
definition-of-done: Jede Notiz im Verdict beginnt mit einem Marker 🔴/🟡/⚪/🔵 nach ihrer Klasse (Test pro Klasse)
tags:
  - acceptance
  - tui
blocked-by: []
related: []
follows: 01M81G681S9YWPBZWQ7HAJ41SC
commits: []
created-at: 2026-11-22T14:45:11Z
updated-at: 2026-11-22T15:09:25Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3863384
claimed-at: 2026-11-22T14:47:40Z
outcome-what: "Verdict- und Give-back-Pane zeigen jede Review-Notiz mit 🔴/🟡/⚪/🔵 (Klasse aus gapItems); mit \"language\" in ~/.jaira/settings.json (auch im S-Fenster: Notes language) schreibt haiku die Notizen kurz in diese Sprache um, gecacht pro Ticket+Notizen in $.store; Ticket bleibt im Original."
outcome-why: "Alex liest langsam und auf Russisch; er musste jede englische Notiz lesen, um zu sehen, ob sie blockiert."
outcome-resolves: "DoD 1-4, 6-8; DoD 5 (Alex sieht es) offen. Geprüft: Nebenläufigkeit (ein Aufruf pro Ticket+Notizen, Set rephrasing), Fehlerpfade (Modell wirft/antwortet nicht/falsche Form -> Original), ungewöhnliche Eingabe (leere/kaputte settings, falsche Anzahl, leere Strings), Rechte (Sprache nur aus eigener settings.json, Notiztext nur zur Anzeige)."
review-summary: none
review-gaps: "round 2: one error return for the two unmarshals in settings.Update. Left: settings.Save is now called only by tests for setup - kept as the plain writer; earlier round's choices stand (stepAfter, keep() vs rephrase store pattern)."
test-verdict: "pass: go vet + go build + go test ./... RC=0, -race tui/settings RC=0, cockpit 316 pass, check-ticket-commits and shipped-history --check RC=0; DoD 1-4,6-8 verified in tree; real haiku answered the prompt in Russian, parsed to 3 notes, paths kept; Save failure (JAIRA_HOME=/dev/null/nope) shows mkdir error, row stays none"
hunt-verdict: "fail: bughunt done at ddd89fd6a, 2 change (1 blocker: settings step overwrites an unparseable settings.json; stale rephrased text after language set to none or failed call)"
round: 2
back-from: testing
back-by: lane
back-count: 1
---

# Das Verdict zeigt jede Notiz mit Ampel-Marker und in der Sprache der Person, verständlich umformuliert

## Definition of Done

- [x] Jede Notiz im Verdict beginnt mit einem Marker 🔴/🟡/⚪/🔵 nach ihrer Klasse (Test pro Klasse)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'each class has its mark, in the Verdict pane and in the Give back pane' (NOTE_MARK in hooks/board.ts)
- [x] Die Notizen erscheinen in der eingestellten Sprache der Person, kurz und einfach umformuliert; das Ticket selbst bleibt unverändert
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'a person who set Русский/Deutsch reads the notes in it, the ticket keeps its words' (rephrase in hooks/register.tsx, rephrasePrompt in hooks/board.ts)
- [x] Die Umformulierung wird pro Ticket und Notizen-Stand gecacht; ohne LLM (Fehler, offline) zeigt das Verdict den Originaltext mit Marker
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'kept per ticket and notes...', 'an answer that does not fit...', 'a model that does not answer...', 'a rewrite shown before goes once the language is taken away...'
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/0GMYQB.md; core/shipped/shipped.sha256 (core/role TestEveryInstalledFileIsAShippedVersion green)
- [ ] seen by a person: Alex liest das Verdict eines Tickets in signoff auf Russisch mit Markern
- [x] Das Give-back-Fenster zeigt die Notizen genauso: Marker pro Notiz und in der Sprache der Person (Alex liest dort am häufigsten, was zu tun ist) - Alex 22.11.
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'Give back shows the notes in the language and writes the originals onto the ticket'
- [x] Die Sprache gilt pro Person: Alex sieht Russisch, ein Kollege mit Deutsch sieht Deutsch; nichts davon steht im Board (Test mit zwei Spracheinstellungen) - Alex 22.11.
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: two tests 'a person who set <language>...'; internal/tui/settings_test.go TestSettingsEnterStepsTheLanguage (two JAIRA_HOMEs, board dir unchanged)
- [x] Neue Einstellung für die Sprache in ~/.jaira/settings.json und im Settings-Fenster (S); ohne Einstellung bleibt der Originaltext - Alex 22.11.
  proof: core/settings/settings.go Settings.Language + Update; internal/tui/settings.go stepLanguage; TestSettingsEnterStepsTheLanguage, TestSettingsWritesNeverWipeABrokenFile, TestUpdateNeverWipesTheFile; cockpit test 'without a language nothing is sent'

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 14:45 · Alexander Sacharov** — decision (Alex 22.11.): Give back gehört in dieses Ticket - dort liest Alex am häufigsten, was getan wird.
- **2026-11-22 14:46 · Alexander Sacharov** — decision (Alex 22.11.): Sprache pro Person, nicht fest Russisch - wer auf Deutsch arbeitet, sieht Deutsch. Die Sprache kommt von der Person bzw. ihrer Maschine (z. B. ~/.jaira/settings.json oder Sprache der Session), nicht vom Board.
- **2026-11-22 14:46 · Alexander Sacharov** — decision (Alex 22.11.): die Sprache ist eine Einstellung in den jaira-Settings (~/.jaira/settings.json, auch im Settings-Fenster S), pro Maschine/Person.
- **2026-11-22 14:47 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 14:49 · Alexander Sacharov** — decisions (worker, settled from code, no open question left): (1) Klasse pro Notiz kommt aus gapItems (board.ts), wie im Give-back-Fenster - Verdict und Give back zeigen dieselbe Klasse; Marker 🔴 blocker, 🟡 edge, ⚪ limit, 🔵 elsewhere ersetzen die Text-Labels [blocker]/[not this ticket]/[known limit]. (2) Einstellung heisst "language" in ~/.jaira/settings.json, frei als Sprachname (an das Modell gegeben); S-Fenster steppt durch none, Deutsch, English, Русский, andere von Hand. (3) Modell: haiku über $.model.complete wie der Thread-Classifier - billiger als chat-model, eine Übersetzung braucht kein sonnet. (4) Cache in $.store je Ticket (Schlüssel rephrased:<id>), gültig solange Sprache und Notizentexte gleich sind; höchstens 50 Tickets. (5) Geschrieben wird immer der Originaltext (decisionLines), nur die Anzeige wird umformuliert.
- **2026-11-22 14:59 · Alexander Sacharov** — known limits (in-progress): (a) 40+ lange Notizen können maxTokens 4000 sprengen -> Antwort passt nicht, Originale werden gezeigt; (b) der Notiztext geht ans Modell - eine Notiz, die das Modell anweist, kann nur die Anzeige verfälschen, nie das Ticket (Anzahl wird geprüft); (c) zwei gleichlautende Notizen teilen eine Übersetzung (Map nach Text). Verdict-Pane listet Notizen jetzt nach Klasse sortiert (blocker zuerst), Give back in Reihenfolge der gaps (Index = Toggle).
- **2026-11-22 15:00 · Alexander Sacharov** — critique (pass 1): none. Left standing: Explain pane keeps its coloured [label] marks - its text comes from the Explain chat, not this ticket's DoD; haiku instead of chat-model per decision note (3); cache keyed by handle only, language and source compared, so a handle collision costs one more call, never a wrong text.
- **2026-11-22 15:04 · Alexander Sacharov** — bughunt: change: blocker: internal/tui/settings.go:179-186 (stepLanguage) - class: a settings-screen step does settings.Load() then settings.Save(), and Load returns an empty Settings for any file it cannot read whole, so the Save overwrites the person's settings.json with defaults plus the stepped field; keys the Go struct does not know are dropped the same way. The change adds two ways in: the new Notes language row, and the new typed Language field - a non-string "language" (was an ignored unknown key before) now makes Load fail, so every jaira command ignores remote/hook too, and one press on any settings row wipes them. Same pattern, older (60155dd, 881cb39): switchDispatch settings.go:194-203, editKey settings.go:229-231. example: settings.json {"remote":"upstream","hook":"/bin/notify",} (trailing comma) or {...,"language":5}, one enter on Notes language -> file is {"language":"Deutsch"}; {"remote":"upstream","future-key":"x",...} loses future-key. ran: scratch go test calling newSettingsScreen + stepLanguage under JAIRA_HOME (deleted afterwards)
- **2026-11-22 15:04 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1078-1110 (rephrase) - class: the rephrasedAtom entry for a ticket is only ever replaced, never cleared, so whenever the current rephrase produces nothing the panes go on showing the last rewrite; register.tsx:1773 and :2002 read said[text] ?? text. example: language set to Deutsch, Verdict opened (German shown), then Notes language stepped to none -> languageOf is undefined, rephrase returns at line 1081 before touching the atom, and Verdict/Give back still show German for the rest of the session, against the settings row's 'none shows them as the ticket has them'; same when the language is switched to Русский and the haiku call fails or is offline -> German stays. read
- **2026-11-22 15:04 · Alexander Sacharov** — bughunt: done at ddd89fd6a - 2 findings (2 change, 1 of them blocker; 0 pre-existing) in 8 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/types/index.d.ts core/release/notes/0GMYQB.md core/settings/settings.go core/shipped/shipped.sha256 internal/tui/settings.go internal/tui/settings_test.go
- **2026-11-22 15:05 · Alexander Sacharov** — testing: class: a settings write from a copy Load could not read whole (and a typed field that makes a whole file unreadable); example: hunter's settings.json with trailing comma or language:5, enter on Notes language -> file reduced to language; fix: one Update path in core/settings that refuses to write over a file that does not parse and keeps unknown keys, used by all three screen writes; Language tolerant of a non-string. Second: rephrase clears the ticket's shown texts whenever it produces none.
- **2026-11-22 15:08 · Alexander Sacharov** — round 2 fix: settings.Update (core/settings) is now the one write path of the S screen - logbook days, dispatch, language: refuses a file that does not parse (error names the path), keeps unknown keys, drops a field the change emptied. Language is its own type that reads a non-string as none, so a typo no longer makes Load drop remote/hook. Save stays for tests and no screen calls it. rephrase clears the ticket's shown texts on no language/failed call.
- **2026-11-22 15:08 · Alexander Sacharov** — critique (pass 2, repair only): none. settings.Update sits in core/settings beside Load/Save, where the file is owned; all three screen writes use it. Save left for test setup.
