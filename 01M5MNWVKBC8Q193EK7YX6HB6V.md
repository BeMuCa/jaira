---
id: 01M5MNWVKBC8Q193EK7YX6HB6V
title: "Die Einzeiler des Cockpits nennen die Lane, in der das Ticket wirklich steht"
status: signoff
ready: true
creator: BeMuCa
goal: "Die Kurzfassung eines Agentenberichts im Cockpit stimmt mit dem Board ueberein: sie nennt den tatsaechlichen Lanewechsel, oder keinen."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkt 10): die Einzeiler waren zweimal falsch ('GH0GKD · list → human', 'done → in-progress') und nannten einmal eine Lane, die das Ticket nie betreten hat. Der volle Bericht musste jedes Mal gelesen werden - dann ist die Kurzfassung wertlos.
  Ursache vermutlich: core/mod/builtin/jaira-cockpit/hooks/digest.ts liest den Lanewechsel aus dem Text des Berichts (Regex auf 'a -> b'), nicht aus dem Board.
  Vorschlag: den Wechsel aus jaira list --json vor und nach dem Bericht ableiten (der Cockpit liest das Board ohnehin, board.ts:3) und den Text nur fuer die Frage nutzen.
definition-of-done: "Der Einzeiler leitet den Lanewechsel aus dem Board ab, nicht aus dem Berichtstext; nennt der Text eine Lane, die das Ticket nicht hat, erscheint kein Wechsel"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:25:58Z
updated-at: 2026-10-26T20:21:17Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: x6hb6v-driver
claimed-at: 2026-10-26T20:10:25Z
outcome-what: "Der Cockpit-Einzeiler zeigt 'vorher → jetzt' nur noch aus dem Board: shorten() in core/mod/builtin/jaira-cockpit/hooks/register.tsx liest die Lane des einen Tickets per `jaira show <id> --json`, erst nachdem digest() ein Ticket und keinen Fehlschlag gefunden hat (digest ist jetzt async). Die Lane, die der Berichtstext nennt, wird nie gelesen; laneChangeOf liegt in digest.ts."
outcome-why: "Der Berichtstext nannte beliebige Lanes ('list -> human', 'done -> in-progress', 'moved to review'), die das Ticket nie hatte; das Board ist die einzige Quelle, die stimmt. lastSeen wird nur beim ersten Auftauchen und bei Berichten gesetzt, weil refreshBoard jede Minute liest und ein minutenlang arbeitender Subagent sonst nie einen Wechsel zeigte."
outcome-resolves: "Der Einzeiler stimmt mit dem Board ueberein: der echte Wechsel oder keiner. Vier Tests in core/mod/builtin/jaira-cockpit/tests/cockpit.test.ts ('the lane in the line comes from the board') decken die drei beobachteten Faelle und den positiven Fall ab; Release-Notiz core/release/notes/X6HB6V.md."
review-summary: "Der Cockpit-Einzeiler liest keine Lane mehr aus dem Berichtstext. digest.ts: die Lane-Liste LANES, die Regexe ARROW ('a -> b') und INTO ('moved to X', 'nach X') und laneOf sind gestrichen; digest() ist jetzt async und fragt den Lanewechsel ueber einen Rueckruf ab, und zwar erst, nachdem es ein Ticket und keinen Fehlschlag gefunden hat. register.tsx: shorten() liest fuer genau dieses eine Ticket 'jaira show <id> --json' und nimmt dessen status (show findet das Ticket auf jedem Board und auf seinem Ref). Verglichen wird mit der Map lastSeen (Handle → Lane): refreshBoard traegt jedes Ticket beim ersten Sehen ein und ueberschreibt nie; jeder Bericht schreibt danach die aktuelle Lane hinein. Gezeigt wird 'vorher → jetzt' nur, wenn beide Lanes bekannt sind und sich unterscheiden, sonst gar keine Lane. Eine Lane, die nur im Text steht, kann deshalb nie mehr erscheinen. In den alten Tests faellt nur die aus dem Text gelesene Lane weg; Verdict, Frage, Fehlschlag und Marker bleiben gleich. Vier neue Tests ('the lane in the line comes from the board') decken 'list -> human', 'done -> in-progress', 'moved to review' und einen echten Wechsel ab. Release-Notiz: core/release/notes/X6HB6V.md (eine Zeile)."
round: 2
review-gaps: |-
  1. Release-Notiz ungenau (klein): sie sagt 'before → now since the cockpit last saw that ticket'. Das Cockpit liest das Board aber jede Minute; die Basis ist die Lane beim letzten Bericht zu diesem Ticket, oder beim ersten Sehen in dieser Sitzung. Vorschlag: '... since the last report about that ticket, or since the session first saw it'.
  2. Bedeutung des Wechsels (Frage an die Person, kein Fehler): er umfasst alles seit dem letzten Bericht, auch Moves, die der Treiber oder die Hauptsitzung selbst gemacht hat. Wo der Treiber nach jedem Bericht schiebt, zeigt z. B. der Bericht des Review-Agenten 'testing → review' - wahr auf dem Board, aber nicht das, was dieser Agent getan hat.
  3. Kein Repo-Test fuer das erste Sehen in refreshBoard (klein): die neuen Tests setzen die Basis mit einem ersten Bericht. Ich habe den Pfad in einer Scratch-Kopie mit einem Wegwerf-Test geprueft (Sitzungsstart sieht in-progress, Move nach critique, Turn-Ende, Bericht zeigt 'in-progress → critique'; zweiter Bericht ohne Lane) - gruen. Ebenso dort: ein Fehlschlag-Bericht startet kein 'jaira show'.
  Sonst nichts gefunden: alle drei DoD-Punkte erfuellt, die Outcome-Angaben stimmen mit dem Diff.
test-verdict: "pass: cockpit 164/164, validate ok, gofmt/vet/build clean, 4 new tests cover the 3 observed cases + positive; jaira show --json status verified on default and second board"
review-verdict: "Freigabe mit kleinem Vorbehalt. Der Diff erfuellt alle drei DoD-Punkte; Fehler habe ich keine gefunden. Geprueft: 164 Cockpit-Tests gruen; die vier neuen Tests schlagen gegen den alten Code von master fehl, sie schuetzen die Aenderung also wirklich; gofmt, go vet, go build sauber; 'jaira show --json' kostet ca. 0,1 s pro Bericht. Vorbehalt: Wortlaut der Release-Notiz (Luecke 1) und die Bedeutung 'seit dem letzten Bericht' (Luecke 2) - beides sollte die Person kennen, beides blockiert nicht. Nicht geprueft: das Verhalten in einer echten Claude-Code-Sitzung (nur im Test-Harness)."
review-round: 2
review-check: |-
  Dauer ca. 10 Minuten. Du tauschst kurz dein globales Cockpit gegen das aus diesem Build und stellst es am Ende zurueck; starte in der Zeit keine anderen Claude-Sitzungen.
  1. cd /home/berk/git/jAIra-X6HB6V && go build -o /tmp/x6/jaira ./cmd/jaira
  2. /tmp/x6/jaira roles install --into /tmp/x6/skills  -> letzte Zeile 'roles in /tmp/x6/skills'; der Ordner /tmp/x6/skills/jaira-cockpit existiert.
  3. mv ~/.claude/skills/jaira-cockpit /tmp/x6/cockpit-alt && cp -r /tmp/x6/skills/jaira-cockpit ~/.claude/skills/
  4. grep -c laneChangeOf ~/.claude/skills/jaira-cockpit/hooks/digest.ts  -> gibt 1 aus (das neue Cockpit ist drin).
  5. Board-Kopie ohne Remote: mkdir -p /tmp/x6/repo && cd /tmp/x6/repo && git init -q && jaira init
  6. jaira create Probe --goal Probe --context Probe --dod Probe, dann jaira list  -> merk dir das 6-stellige Handle, unten H.
  7. jaira move H --to todo  -> 'H → todo'.
  8. cd /tmp/x6/repo && claude  -> dem Ordner vertrauen.
  9. Tippe: 'Starte einen Hintergrund-Agenten. Er fuehrt jaira move H --to in-progress aus und antwortet danach nur mit diesem Satz: Ticket H: jaira list -> human, done -> in-progress, moved to review.'
  10. Wenn er fertig ist, siehst du genau eine Zeile: '[shortened by jaira-cockpit] H · todo → in-progress — full report: /cockpit report H'. Es darf NICHT 'list → human', 'done → in-progress' oder 'review' darin stehen.
  11. Tippe: 'Starte noch einen Hintergrund-Agenten. Er verschiebt nichts und antwortet nur: Ticket H moved to review.'
  12. Du siehst: '[shortened by jaira-cockpit] H — full report: /cockpit report H' - ohne Lane, weil das Ticket noch in in-progress steht (jaira show H zeigt 'lane in-progress').
  13. /cockpit report H  -> zeigt beide Berichte im vollen Wortlaut.
  14. Zurueckstellen: claude beenden, dann rm -rf ~/.claude/skills/jaira-cockpit && mv /tmp/x6/cockpit-alt ~/.claude/skills/jaira-cockpit
  Nur per Test pruefbar: cd /home/berk/git/jAIra-X6HB6V/core/mod/builtin/jaira-cockpit && claude plugin test . 2>&1 | tail -4  -> '164 pass, 0 fail'; der Block 'the lane in the line comes from the board' (4 Tests) deckt die drei beobachteten Faelle ab. Das erste Sehen beim Sitzungsstart deckt nur Schritt 10 ab, kein Repo-Test.
---

# Die Einzeiler des Cockpits nennen die Lane, in der das Ticket wirklich steht

## Definition of Done

- [x] Der Einzeiler leitet den Lanewechsel aus dem Board ab, nicht aus dem Berichtstext; nennt der Text eine Lane, die das Ticket nicht hat, erscheint kein Wechsel
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx laneOf()/shorten(), digest.ts async digest(text, laneChange)
- [x] Ein Test deckt die drei beobachteten Fehlfaelle ab
  proof: tests/cockpit.test.ts 'the lane in the line comes from the board' (4 tests, fake board answers jaira show)
- [x] Release-Notiz unter core/release/notes/
  proof: core/release/notes/X6HB6V.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-26 19:48 · BeMuCa** — driver: pre-process is not ticked for this ticket (Options), so the design is written here. Why not 'previous poll vs now': refreshBoard reads jaira list --json every 60 s and after every turn; a subagent works for minutes, so the poll sees its moves long before its report arrives - 'last poll vs now' would show no change almost always. Design instead: register.tsx keeps one map handle -> lane 'as of the last report about this ticket'. refreshBoard only ADDS tickets it has not seen yet (first sight), never overwrites. On a report, shorten() reads jaira list --json fresh (default board plus the other boards refreshBoard reads), emits '<before> -> <now>' when both are known and differ, then stores <now>. The report text is no longer read for lanes at all (laneOf/ARROW/INTO/LANES go), so a lane only the text names can never show up. Verdict stays text-derived: not in the DoD, not one of the observed failures.
- **2026-10-26 19:51 · BeMuCa** — lanesNow() liest bei jedem Bericht frisch; in shorten wird board nicht gecacht (sonst entfiele checkVersion in refreshBoard). Bei Lesefehler einer Nebenboard-Liste: gar kein Lanewechsel.
- **2026-10-26 19:57 · BeMuCa** — critique (first pass): five findings, all in review-summary. (1+2, one fix) read the lane only once digest has found a ticket and no failure — make digest's laneChange callback async — and read it with one 'jaira show <ticket> --json' (status) instead of lanesNow: show resolves the id on every board (core/ticket/store.go:1174; checked on a scratch board: an 'ops' ticket is missing from 'jaira list --json' but 'jaira show <id> --json' gives status and board), so the copy of refreshBoard's board loop and its separate failure rule go, and failure reports (which shorten sees twice, session.receive then prompt.submit) spawn no jaira at all. (3) laneChangeOf to digest.ts — it formats the line, board.ts parses jaira. (4) release note as an instruction, and not 'as of the last report' (the first report compares with first sight). (5) cockpit.test.ts first(): drop the no-op lanes/to parameters. Let stand: lastSeen seeded at first sight and updated only on reports (driver note 19:48) — found no case where it shows a lane the ticket did not take; no board cache in shorten (note 19:51); verdict stays text-derived (19:48); the edits to the existing tests drop only the text-derived lane, every verdict/question/failure/usage/marker assertion is unchanged.
- **2026-10-26 20:01 · BeMuCa** — critique (second pass): none left - checked f04f549e against the five first-pass findings, all addressed: (1) digest is async and awaits laneChange only after isFailure/ticketOf (digest.ts:85-92); (2) lanesNow and its board loop are gone, laneOf reads one 'jaira show <id> --json' status (register.tsx:218-223; checked: show --json carries top-level status); (3) laneChangeOf now in digest.ts:76; (4) release note is one instruction line; (5) first($, ticket) without the no-op parameters. Repair introduced nothing that breaks the DoD; first-pass let-stands untouched.
- **2026-10-26 20:06 · BeMuCa** — testing: pass. Gates green. Verified on scratch board (no remote): show --json status backlog->todo on default and --board second ticket; plain 'jaira show <id>' finds second-board ticket. Not exercised: live Claude Code band. Failure path: digest() returns before calling laneChange when isFailure, so no jaira show starts; existing tests assert text stays whole but do not assert no show call (minor gap).
- **2026-10-26 20:20 · BeMuCa** — review: Freigabe mit kleinem Vorbehalt, nichts zurueckgeschickt. (1) Release-Notiz core/release/notes/X6HB6V.md sagt 'since the cockpit last saw that ticket' - die Basis ist aber die Lane beim letzten Bericht (oder beim ersten Sehen in der Sitzung), nicht der letzte Board-Read (jede Minute); Vorschlag: 'since the last report about that ticket, or since the session first saw it'. (2) Der Wechsel umfasst auch Moves des Treibers/der Hauptsitzung seit dem letzten Bericht (z. B. Review-Bericht zeigt 'testing → review'). (3) Kein Repo-Test fuer firstSight in refreshBoard; in Scratch-Kopie mit Wegwerf-Test geprueft: gruen, ebenso 'Fehlschlag startet kein jaira show'.
