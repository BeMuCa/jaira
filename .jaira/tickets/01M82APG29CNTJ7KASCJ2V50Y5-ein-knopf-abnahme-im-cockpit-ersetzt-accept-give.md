---
id: 01M82APG29CNTJ7KASCJ2V50Y5
title: "Ein Knopf 'Abnahme' im Cockpit ersetzt Accept, Give back und Verdict: pro Notiz ok / beheben / Ticket wählen, dann ausführen"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die Person nimmt ein Ticket in einem Fenster ab: sie sieht Ampel und Notizen, wählt pro Notiz ok, beheben oder Ticket, und ein Knopf 'Ausführen' macht daraus done, Give back und neue Tickets."
context: |-
  Was heute falsch ist (Alex 22.11.2026): Accept, Give back und Verdict zeigen dieselben Notizen in drei Fenstern. Alex liest meist im Give back, weil er dort pro Notiz umschalten kann, und muss dann trotzdem einen anderen Knopf drücken.
  Vorschlag Alex: ein Knopf 'Abnahme' (RU: 'Приём'), darin wie heute im Give back pro Notiz umschalten:
  - alle ok -> Ticket nach done (wie Accept heute)
  - mindestens eine 'beheben' -> Give back mit genau diesen Notizen
  - 'Ticket' -> neues Ticket mit follows auf dieses, Text aus der Notiz
  Bestehendes: Ampel oben (AJ41SC), Marker pro Notiz und Übersetzung (0GMYQB), Give back startet Worker (QTQW15), Accept-Gates (DoD, Bughunt, pr-when, CRASZS).
  Code: core/mod/builtin/jaira-cockpit/hooks/register.tsx (Verdict-, Give-back- und Accept-Pfade), board.ts (gapItems).
  Erst starten, wenn QTQW15, DF46B9 und T5QN9P gemergt sind - gleiche Dateien.
definition-of-done: "Ein Knopf 'Abnahme' pro Ticket in signoff ersetzt Accept, Give back und Verdict in der Zeile; das Fenster zeigt Ampel, DoD/Bughunt-Stand und jede Notiz mit Marker in der Sprache der Person (Test)"
tags:
  - cockpit
blocked-by: []
related: []
commits: []
created-at: 2026-11-22T19:11:01Z
updated-at: 2026-11-22T21:08:08Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-6414
claimed-at: 2026-11-22T19:38:15Z
outcome-what: "Explain all back under the row's one lock; Execute pressed while the row is busy says so (t.busyRow) and is drawn dim instead of being dropped."
outcome-why: "bughunt round 2: a lock of its own let an Explain chat run beside a press that pulled and then removed the ticket copy"
outcome-resolves: "bughunt note of round 2"
review-summary: "none"
round: 4
back-from: testing
back-by: lane
back-count: 2
review-gaps: "none: round 3 fix diff adds no duplicate or dead code; earlier: removed nextChoice, Waiting.explainable, the Verdict pane and seven unused texts, folded store-and-trim into storeNewest"
test-verdict: "fail: suite green (410/410 on 2.1.329, go test shipped/release RC=0); bughunt round 2 found one change finding: Explain all's own lock lets a chat run beside a press that pulled the ticket copy"
hunt-verdict: "fail: 1 change finding - register.tsx Explain all lock <id>:chat races a press holding the row lock"
---

# Ein Knopf 'Abnahme' im Cockpit ersetzt Accept, Give back und Verdict: pro Notiz ok / beheben / Ticket wählen, dann ausführen

## Definition of Done

- [x] Ein Knopf 'Abnahme' pro Ticket in signoff ersetzt Accept, Give back und Verdict in der Zeile; das Fenster zeigt Ampel, DoD/Bughunt-Stand und jede Notiz mit Marker in der Sprache der Person (Test)
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx on('ui.render', requestId: ACCEPTANCE) and the row's acceptance-<id> button; tests/board.test.ts 'Acceptance is the one coloured button of the row', 'the head says what the ticket is for', 'the head comes in the person's language'
- [x] Pro Notiz drei Zustände ok / beheben / Ticket; 'Ausführen' macht: alle ok -> done über dieselben Gates wie Accept heute; eine 'beheben' -> Give back mit genau diesen Notizen; 'Ticket' -> jaira create mit follows, auch zusammen mit done oder Give back (Tabellentest über die Kombinationen)
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts executionOf; tests/board.test.ts 'the Acceptance pane, read alone > Execute: *' (17-row table), 'Give back writes exactly the notes to fix', 'a note that is no blocker ... decisions go onto the ticket once it is accepted'
- [x] Ein Gate, das done verweigert (DoD offen, Bughunt fehlt, PR nicht settled), sagt es im Fenster und ändert nichts; Entscheidungen der Person gehen nicht verloren (Test)
  proof: tests/board.test.ts 'a gate refusing done is said in the Acceptance pane: nothing moves, nothing is written, the choices stay', 'a refusal is shown whole in the Acceptance pane'
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/2V50Y5.md; core/shipped/shipped.sha256 via scripts/shipped-history.sh --write (go test ./core/shipped ok, --check ok)
- [ ] seen by a person: Alex nimmt ein Ticket über 'Abnahme' ab
- [x] Kopf des Fensters, schön gesetzt und in der Sprache der Person: Titel, Ziel (goal), was gemacht wurde und warum (outcome-what/why), Ampel mit Grund - Alex 22.11.
  proof: tests/board.test.ts 'the head says what the ticket is for, what was done and why', 'the head comes in the person's language' (headPrompt, board.ts)
- [x] Verlauf auf einen Blick aus jaira-Daten: Runden (round, Rückgaben je Lane, back-count/back-from), Commits und geänderte Dateien, DoD erledigt/superseded/gesamt, Bughunt auf HEAD, Review-Urteil, Zeit von created bis jetzt, PR-Stand wenn vorhanden (jaira pr), Release-Note-Zeile(n) - Alex 22.11.
  proof: tests/board.test.ts 'the head says what the ticket is for...': round/back/age, DoD with superseded, bughunt on HEAD, review, commits+files, pull request, release note; board.ts storyOf/commitsOf/releaseLines
- [x] Pro Notiz ein Knopf 'erklären': ein kurzer Satz in Alltagssprache, warum sie wichtig ist oder nicht (ersetzt das heutige Explain, das nur bei Blocker/Vorbehalt erschien) - Alex 22.11.
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx explainNote; tests/board.test.ts 'Explain on one note asks the cheap model once for a sentence', 'a model that gives no sentence says so'
- [x] Ampel ist Pflicht: jede Notiz trägt ihren Marker 🔴/🟡/⚪/🔵, und die Gesamtampel oben folgt aus den Markern (ein 🔴 -> rot, sonst ein 🟡 -> gelb, sonst grün) zusammen mit den Gates; Explain-Knopf pro Notiz und für das ganze Ticket (Test) - Alex 22.11.
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts lightOf(w, hunt, classes); tests/board.test.ts 'a known limit or a note for another ticket colours nothing', 'Explain all opens the explaining chat from the pane'
- [x] Eine Erklärung (explain:-Notiz) aus einer früheren Runde gilt nach Give back nicht mehr: 'Read explanation' verschwindet, sobald das Ticket in eine neue Runde geht, und erscheint erst mit einer neuen Erklärung wieder (heute bleibt der Knopf stehen, gesehen bei T5QN9P 22.11.) (Test) - Alex 22.11.
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts GIVEN_BACK in explanationNotes; tests/board.test.ts 'Read explanation goes once the ticket was given back after it', 'an explanation a give-back came after is none'
- [x] Fenster folgt dem Layout-Entwurf vom 22.11. (Option 1): Kopf mit Ampel, Block 'Ход', Notizen mit [ок][исправить][тикет][?], Fuß mit Итог + [Объяснить всё][Выполнить] - Alex 22.11.
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx Pane ACCEPTANCE drawn as the 22.11. sketch: light+head, History, notes with [ok][fix][ticket][?], summary line, Execute + Explain all
- [x] Edges: leere review-gaps (kein Review) -> Fenster zeigt nur Kopf/Verlauf, Ausführen = Accept; sehr viele Notizen scrollen; Notiz mit '-' am Anfang oder Klammern wird Folge-Ticket ohne Flag-Verwechslung (bestehendes --); zweiter Druck auf Ausführen während der erste läuft tut nichts (once); Ticket inzwischen woanders (stillIn) -> Fenster sagt es, nichts geschrieben (Test)
  proof: tests/board.test.ts 'decided meanwhile by someone else', 'very many notes are all rows', 'every value of the follow-up goes with =', 'a pressed row locks ... a second press does nothing', 'no notes: the pane asks what is missing'
- [ ] seen by a person: das Fenster 'Приём' im Terminal - Kopf, Ход, Marker und Knöpfe stehen sauber, kein abgeschnittener Text

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Pure Funktionen in board.ts: executionOf (Tabelle über alle Kombinationen), acceptanceOf (Kopf+Verlauf aus show-JSON), explanationNotes mit Rundenregel, notePrompt/noteReasonOf - Tests zuerst
- [x] texts.ts: Wörter für Fenster, Knopf Приём/Acceptance, Verlaufszeilen, Gate-Zeile (de/en/ru wie bestehend)
- [x] register.tsx: acceptanceAtom + openAcceptance (show-JSON, git log, release note, PR-Stand, rephrase von goal/what/why/notes)
- [x] register.tsx: Pane ACCEPTANCE zeichnen nach Layout Option 1 (Kopf, Ход, Заметки mit 4 Knöpfen, Fuß mit Input/Button)
- [x] register.tsx: Ausführen über executionOf -> giveBackWith/acceptWith, die ein Ergebnis zurückgeben; Verweigerung bleibt im Fenster
- [x] Band-Zeile umbauen: Accept/Give back/Verdict raus, Приём rein; alte Panes GIVE_BACK/VERDICT entfernen, wenn nichts sie mehr öffnet
- [x] Cockpit-Tests (cockpit.test.ts) für Knopf, Gate-Verweigerung, Rundenregel; release note + shipped.sha256; claude plugin test + go test ./core/shipped

## Progress
- **2026-11-22 19:12 · Alexander Sacharov** — decision (Alex 22.11.): das Abnahme-Fenster soll der Person das Ticket verständlich machen - mehr Information, schön gesetzt: Titel, Ziel, was/warum, Runden und Verlauf, alles was jaira liefert; Explain pro Notiz statt nur bei Blocker. Layout-Entwurf zuerst als ASCII an den Teamlead (Frage), damit Alex ihn sieht, bevor gebaut wird.
- **2026-11-22 19:38 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 19:38 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 19:39 · Alexander Sacharov** — question: Layout-Entwurf 'Abnahme' (Приём) - so bauen?
┌ Приём · 2V50Y5 ─────────────────────────────── 🟡 жёлтый ┐
│ Ein Knopf 'Abnahme' ersetzt Accept, Give back, Verdict   │
│ Цель:   <goal, übersetzt>                                │
│ Сделано: <outcome-what>                                  │
│ Зачем:   <outcome-why>                                   │
│ Светофор: жёлтый - 1 мелочь 🟡, гейты ок                 │
├ Ход ─────────────────────────────────────────────────────┤
│ Раунд 2 · возвратов: critique 1, testing 0 · 3 дня       │
│ DoD 9/10 (1 superseded) · Bughunt на HEAD ✓ · Review: ok │
│ 7 коммитов · 12 файлов · PR #195 open/settled · Note: 1  │
├ Заметки ─────────────────────────────────────────────────┤
│ 🟡 <Notiz 1 übersetzt>        [ок] [исправить] [тикет] [?]│
│ ⚪ <Notiz 2>                  [ок] [исправить] [тикет] [?]│
│ 🔵 <Notiz 3>                  [ок] [исправить] [тикет] [?]│
├──────────────────────────────────────────────────────────┤
│ Итог: done + 1 новый тикет      [Объяснить всё] [Выполнить]│
│ (Gate verweigert -> rote Zeile hier, nichts geändert)     │
└──────────────────────────────────────────────────────────┘
[?] = erklären: ein Satz Alltagssprache, inline unter der Notiz.
Optionen: 1) so bauen (Recommended) 2) Verlauf einklappbar, Notizen zuerst 3) anders - bitte sagen
Zusatz: Bau startet erst nach Merge von QTQW15 (signoff), DF46B9 und T5QN9P (done, nicht in master) - gleiche Dateien.
- **2026-11-22 19:42 · Alexander Sacharov** — decision (Alex 22.11.): Layout Option 1 - so bauen wie im Entwurf (Kopf, Verlauf, Notizen mit [ok][beheben][Ticket][?], Ergebnis unten). Erst nach dem Merge von QTQW15/DF46B9/T5QN9P mit dem Code beginnen.
- **2026-11-22 20:07 · Alexander Sacharov** — decision (pre-process, settled from the code): (1) die drei Zustände sind die bestehenden Choice-Werte accept/now/later (ok/beheben/Ticket), Vorbelegung wie heute (Blocker -> beheben, sonst ok). (2) 'Ausführen' ist eine reine Funktion executionOf(items, text, acceptTo) in board.ts: eine 'beheben' oder getippter Text -> Give back (giveBackWith, Folge-Tickets für 'Ticket' inklusive); sonst Accept über acceptWith (dieselben Gates: dry-run, otherRefusal); ein behaltener Blocker braucht wie heute einen Grund; acceptTo null (human-Lane) und alles ok -> Fenster sagt es, nichts passiert. (3) Ein Textfeld: Label wechselt - 'warum trotz Blocker annehmen' (Pflicht) wenn ein Blocker nicht beheben ist und nichts beheben, sonst 'was noch zu beheben ist' (optional, nicht leer -> Give back); submitLabel 'Выполнить'. Ohne Input (mobil) ein Button 'Выполнить'. (4) Band-Zeile: [Приём] statt Accept/Give back/Verdict; Explain-Chat wandert als [Объяснить всё] ins Fenster; Read explanation und Discuss bleiben in der Zeile. (5) Fenster bleibt bei Gate-Verweigerung offen, Verweigerung als rote Zeile im Fuß, Entscheidungen bleiben im Atom - heute schließt close() den Give-back-Pane vor acceptWith und die Wahl ist weg. (6) DoD 10: eine explain:-Notiz gilt nicht mehr, sobald nach ihr eine 'Given back from'-Notiz steht (explanationNotes); jaira stempelt Rundenwechsel nicht, und das Cockpit schreibt diese Notiz bei jedem Give back. Ein Rückschieben per CLI ohne Notiz bleibt unerkannt (known limit). (7) Verlauf aus einem jaira show --json (round, back, dod_items state, created_at, review) + git log --grep=<ID> im Ticket-Worktree (Commits, Dateien) + core/release/notes/<ID>.md falls vorhanden; PR-Stand über die vorhandenen pullTarget/pullRequestsCommand/hasPullRequest - ein 'jaira pr' gibt es nicht. (8) Erklären pro Notiz: ein haiku-Aufruf beim Druck, Antwort im Store je Notiztext+Sprache, ohne gesetzte Sprache in der Sprache der Notiz.
- **2026-11-22 20:12 · Alexander Sacharov** — in-progress: board.ts hat executionOf/storyOf/releaseLines/commitsOf/reasonPrompt/reasonOf und die Rundenregel in explanationNotes, mit Tabellentest (403 grün). Entscheidung: die festen Wörter des Fensters bleiben englisch wie alle Cockpit-Texte (texts.ts: 'one set of words'); in der Sprache der Person sind Titel, Ziel, was/warum, Ampel-Grund und Notizen (ein rephrase-Aufruf), wie DoD 6 sie aufzählt. Der Knopf heißt 'Acceptance'. Pro Notiz vier Knöpfe [ok][fix][ticket][?] wie im Layout statt eines Umschalt-Knopfs. Execute ist ein eigener Knopf; das Feld meldet seinen Text über onInput (Engine-Input kennt change/submit), Enter im Feld führt ebenfalls aus.
- **2026-11-22 20:40 · Alexander Sacharov** — in-progress: Fenster gebaut, alte Tests migriert (399 grün). Gefunden und behoben: (a) das Öffnen hielt die Zeilen-Sperre, bis git log/PR-Abfrage fertig waren - ein Execute in der Zeit ging verloren; jetzt füllt readStory die Historie nach dem Öffnen außerhalb der Sperre. (b) pullTarget macht git fetch - für die reine Anzeige zu schwer; openPullRequest fragt nur gh/glab (repo view + pr list), ohne fetch. (c) Enter im Feld nimmt nie an (ein Enter zu viel brächte das Ticket nach done) - nur der Knopf Execute; mit Text oder 'fix' führt Enter aus. (d) Kopf-Übersetzung ist ein eigener Modellaufruf (headPrompt, Store-Key <id>:head), damit die Notiz-Übersetzung unverändert bleibt. Tests: Hilfen acceptancePane/pressAccept/pressExplain/cycle in board.test.ts; quietly() filtert, was das Öffnen ausführt, aus w.ran.
- **2026-11-22 20:49 · Alexander Sacharov** — critique: openPullRequest (register.tsx) kopiert die erste Hälfte von pullTarget (forge, origin, repo view, targetRepoOf) - in einen gemeinsamen Helfer forgeRepo ziehen. Bewusst stehen gelassen: Give back schließt das Fenster vor dem Move wie bisher (DoD 3 betrifft die done-Gates); der Release-Note-Pfad core/release/notes/<ID>.md ist jaira-spezifisch und fehlt in anderen Repos einfach (DoD 7 verlangt ihn); canType-Rückfall stammt aus dem bisherigen hasInput.
- **2026-11-22 20:50 · Alexander Sacharov** — critique round 2: forgeRepo ist jetzt der eine Ort für forge+Ziel-Repo (pullTarget und openPullRequest), Fehlertexte unverändert; nichts weiter.
- **2026-11-22 21:00 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/board.ts:144 (executionOf) with hooks/texts.ts:269,272 - class: a lane that does not end the work (acceptTo null, e.g. human) while a blocker is not marked fix: typed words are read as the reason for accepting (asksWhy), and accepting is impossible there, so Execute refuses and the refusal tells the person to 'type what is missing' - which they did; the field is also labelled 'Why accept it with the blocker?' in a lane where nothing can be accepted, and on mobile (no field) the same text asks for typing that cannot happen. Only marking the blocker fix gets the ticket back. example: human lane, gaps 'blocker: x' set to ok (or ticket), typed 'the docs are missing', Execute -> step notTerminal, nothing moves. ran: node --experimental-strip-types on executionOf([{cls:'blocker',choice:'accept'}],'the docs are missing',null) -> {step:'notTerminal',asksWhy:true}; same for choice 'later'
- **2026-11-22 21:00 · Alexander Sacharov** — bughunt: change: core/release/notes/2V50Y5.md:4 against hooks/register.tsx:2069-2076 - class: the release note promises a safety property the code does not keep; 'Enter in the pane's field never accepts - press Execute' is false whenever a blocker is not marked fix and a reason is typed: onSubmit lets step 'accept' with asksWhy through to execute(), which moves the ticket to done. example: signoff, 'blocker: the header' set to ok, type a reason, press Enter -> jaira move --to done --force. read; the shipped test 'a blocker not marked fix now is accepted only with a reason...' (tests/board.test.ts:2675) proves exactly that Enter accepts. Either the note line or the onSubmit guard is wrong
- **2026-11-22 21:00 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1255 (explainNote) - class: a failed model call is cached as the answer for the rest of the session; put(t.noReason) writes 'no explanation came back' into reasonsAtom, and the guard at :1256 returns early on any value but undefined, so the ? beside that note can never ask again until a restart; reasonsAtom is also keyed by ticket and note text only, so after the notes language changes the old-language sentence stays. example: haiku times out once -> every later press of ? on that note shows 'no explanation came back'. read (store is skipped on failure, atom is not cleared)
- **2026-11-22 21:00 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2084-2085 with :895 (once) and :983 (openChat returns 'hold') - class: Explain all and Execute now sit side by side under the same row lock; Explain all without Herdr returns 'hold', which keeps the lock for BUSY_MS (20 s), and with Herdr spawn.sh holds it until it ends or 20 s pass; an Execute pressed in that time is dropped by once() without a word, and the pane's Execute is not drawn dim the way the band's button is, so the person sees nothing happen and the ticket stays. example: press Explain all, then Execute within 20 s -> no move, no message. read
- **2026-11-22 21:00 · Alexander Sacharov** — bughunt: done at 592e88dc6 - 4 findings (4 change, 0 blocker, 0 pre-existing) in 6 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/types/index.d.ts core/release/notes/2V50Y5.md core/shipped/shipped.sha256
- **2026-11-22 21:01 · Alexander Sacharov** — testing: zurück an in-progress für die 4 bughunt-Funde: (1) asksWhy nur, wo angenommen werden kann (acceptTo != null) - sonst ist Text 'was fehlt' und gibt zurück; (2) Release-Note: Enter nimmt nur mit Grund für einen behaltenen Blocker an; (3) explainNote: Fehler nicht als Antwort festhalten, beim Öffnen des Fensters die Sätze zurücksetzen, damit ein Sprachwechsel neu fragt; (4) Explain all unter eigener Sperre <id>:chat, damit Execute nicht 20 s verloren geht.
- **2026-11-22 21:03 · Alexander Sacharov** — critique round 3 (fix diff 592e88dc6..087af7b44): die vier Fixes sitzen an der Stelle des Fehlers (executionOf, explainNote/openAcceptance, eigene Sperre für Explain all); nichts weiter.
- **2026-11-22 21:05 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2091 - class: Explain all now takes its own lock (`${w.id}:chat`), so openChat runs beside any press on the row's lock w.id that pulled the ticket's copy into the main checkout itself - Execute (accept at :664-672, giveBack at :1477-1482) and Discuss (:1932/:1973); the guard at :966-968 sees only chats that registered a copy, and ticketHere (:582-583) answers 'current' with no copy while that press holds it; example: ticket only on its ref, no own worktree - Execute pulls the copy and runs dry-run/move/follow-ups for seconds, Explain all pressed meanwhile finds the file current, opens its chat on the main checkout with no copy of its own, then accept's here.leave() runs git clean on that file and the chat's 'jaira note' meets a ticket on its ref and not on disk (refused). Same if Discuss and Explain all are pressed within the ticketHere window: two chats, each believing it owns or need not own the copy, the first to end removes it from under the other. Before 087af7b44 both buttons shared w.id and the second press did nothing. Not a blocker: narrow timing, the chat sees the refusal, the accept or give-back itself is unaffected; read (call path followed, not reproduced in the harness)
- **2026-11-22 21:05 · Alexander Sacharov** — bughunt: done at 087af7b44 - 1 finding (1 change, 0 blocker, 0 pre-existing) in 6 modules, follow-up hunt (diff since 592e88dc6: board.ts executionOf, register.tsx explainNote/openAcceptance/Explain all lock, notes, tests); concurrency, error paths, input, permissions checked; claude plugin test: 410 pass; files: core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/types/index.d.ts core/release/notes/2V50Y5.md core/shipped/shipped.sha256
- **2026-11-22 21:07 · Alexander Sacharov** — testing round 2 -> in-progress. decision (kleineres Design statt neuer Regel): die eigene Sperre für Explain all wieder entfernt - eine Sperre pro Ticket, zwei Schreiber auf einem Ticket nie gleichzeitig. Execute während einer laufenden Sperre wird jetzt gesagt (t.busyRow, Knopf gedimmt) statt still verworfen; das war der eigentliche Fund aus Runde 1.
