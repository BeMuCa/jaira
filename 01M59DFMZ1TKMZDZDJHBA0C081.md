---
id: 01M59DFMZ1TKMZDZDJHBA0C081
title: "Ein Ticket laesst sich mit Grund ablehnen: jaira reject, und das Archiv sagt warum"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer ein Ticket nicht macht, sagt mit einem Befehl warum - und Monate spaeter steht der Grund am Ticket, getrennt von bloss weggeraeumten."
context: |-
  Was heute fehlt: ein abgelehntes Ticket sieht genauso aus wie ein weggeraeumtes.
  - 'jaira archive <id>' (core/ticket/store.go:301, Taste x auf dem Board) verschiebt nach .jaira/archive/ und schreibt KEINEN Grund.
  - 'jaira delete <id>' loescht endgueltig - nur fuer versehentlich angelegte Tickets, keine Ablehnung.
  - Folge: 'jaira archive' (Liste) und das Web-Portal koennen nicht sagen, ob ein Ticket ein Duplikat, verworfen oder ueberholt war. Das verletzt den Kernwert 'was war das Ticket und wo steht es'.
  Anlass: Alex am 19.10.2026 - 'wir haben keine richtige Moeglichkeit, ein Ticket zu loeschen oder abzulehnen'.
  Entschieden von Alex (19.10.2026):
  - Neuer Befehl 'jaira reject <id> --why <text> [--duplicate-of <id>]'. Art: wontfix | duplicate | obsolete (--duplicate-of setzt duplicate).
  - Schreibt ins Frontmatter einen Block 'rejected:' mit kind, why, by, at; danach archiviert er wie 'jaira archive' (gleicher Code, kein zweiter Weg).
  - Nur ein Mensch lehnt ab. Ein Agent schiebt das Ticket stattdessen nach human mit --question 'ablehnen? Grund: ...'. Wie ein Agent erkannt wird: denselben Mechanismus nehmen, der heute den Ausgang aus human sperrt (requires-human-exit, siehe VW1WB1) - nichts Neues erfinden.
  Verworfen: Grund nur an archive haengen (Aufraeumen und Ablehnen blieben ununterscheidbar); eigene Lane 'rejected' (aendert die Lanes jedes Boards).
  Offen fuer die Plan-Lane: ob x auf dem Board kuenftig nach dem Grund fragt oder eine eigene Taste ablehnt (x bleibt archivieren, siehe model.go:1532-Kommentar).
definition-of-done: "'jaira reject <id> --why <text>' schreibt rejected.kind/why/by/at ins Frontmatter und verschiebt das Ticket nach .jaira/archive/; ohne --why Exit 2"
tags:
  - cli
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-19T06:27:18Z
updated-at: 2026-10-19T07:32:19Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-9569
claimed-at: 2026-10-19T07:04:10Z
mode: conversational
outcome-what: "Neuer Befehl jaira reject <id> --why [--kind wontfix|obsolete] [--duplicate-of] [--force] schreibt rejected-kind/why/by/at und archiviert ueber Store.Archive; Taste R auf dem Board; archive-Liste und show zeigen den Grund; restore nimmt die Felder ab (Doc.Remove)."
outcome-why: "Ein abgelehntes Ticket sah im Archiv genauso aus wie ein weggeraeumtes; Monate spaeter liess sich nicht sagen, ob es Duplikat, verworfen oder ueberholt war."
outcome-resolves: "DoD 1-3: internal/cli/reject.go + Store.Reject, Tests in internal/cli/reject_test.go; DoD 4: modeReject in internal/tui/model.go, internal/tui/reject_test.go; DoD 5: archivedRejections/rejectionLine, unreject; DoD 6: release note, agent block, Skill, go test ./... -race gruen."
review-summary: "Neuer Befehl 'jaira reject <id> --why <text> --force' (internal/cli/reject.go): schreibt vier flache Keys rejected-kind/why/by/at per Store.Reject (core/ticket/store.go) und archiviert dann ueber Store.Archive - derselbe Weg wie 'jaira archive'. --kind wontfix|obsolete (Default wontfix); --duplicate-of loest das Original ueber LoadAnywhere auf, setzt kind=duplicate und traegt es in related ein (unbekannt -> Exit 5, sich selbst -> Exit 1). Ohne --why Exit 2, ohne --force Exit 3 mit der fertigen 'jaira move ... --to human --question'-Zeile (Lane wird aus dem Board gelesen: questionLane). Vor dem Schreiben prueft Reject, ob der Name im Archiv schon belegt ist (archiveTarget, von Archive mitgenutzt), sodass eine Ablehnung nie auf einem Ticket landet, das auf dem Board bleibt. TUI: Taste R (Board und Detailansicht) oeffnet eine Zeile fuer den Grund, leerer Grund laesst die Zeile offen, esc bricht ab; x archiviert weiter ohne Frage. 'jaira archive' listet unter jeder abgelehnten Datei eine Zeile 'rejected <kind>: <why> (<by>, <datum>)' und im --json ein rejected-Objekt pro Datei; 'jaira show' hat eine rejected-Zeile und in --json immer ein rejected-Objekt. 'jaira restore' entfernt die rejected-* Keys (neues Doc.Remove) bevor die Datei zurueckgeht. rejected-why ist im Merge-Driver ein Prosa-Feld. Agent-Block, AGENTS.md/CLAUDE.md, Skill, README, COMMANDS.md und die Release-Note nennen reject."
round: 3
review-gaps: "1) Die Aenderungen der optimize-Lane (archiveTarget-Extraktion in core/ticket/store.go, toter nil-Zweig in internal/tui/model.go, LoadHeader im Test internal/cli/reject_test.go) sind NICHT committet - sie liegen nur im Worktree; ohne Commit geht der Branch ohne sie raus. 2) Ein mehrzeiliger Grund (--why mit Zeilenumbruch, in der Shell moeglich; die TUI faltet auf eine Zeile) bricht die Einrueckung in 'jaira archive': die zweite Zeile steht ohne Einzug an Spalte 0 und sieht aus wie ein eigener Eintrag. Kosmetisch, keine DoD verletzt. 3) Klein: der CLI-Pfad stempelt Commits (StampCommits) bevor Reject prueft, ob es ueberhaupt geht - scheitert Reject (Archivname belegt, --duplicate-of auf sich selbst), hat das Ticket trotzdem eine geaenderte Datei; 'jaira archive' macht es genauso. --duplicate-of auf sich selbst endet mit Exit 1 statt Exit 2 (Usage). 4) Beobachtet, nicht geprueft: 'jaira restore' entfernt rejected-* nur aus der Datei auf der Platte, nicht vom Ref - das ist das bestehende Verhalten von restore fuer archivierte Tickets und nicht neu."
test-verdict: "pass: go clean -testcache; go test ./... -race RC=0. DoD 1-3, 5 auf einem Scratch-Board mit dem gebauten Binaer gesehen: reject ohne --force RC=3 mit human-Weg, ohne --why RC=2, --duplicate-of ZZZZZZ RC=5, mit --duplicate-of: kind duplicate + related, archive-Liste und show zeigen 'duplicate: gleiches Thema (Alexander Sacharov, 2026-10-19)', restore RC=0 und 0 rejected-Zeilen im Ticket. DoD 4 (Taste R) nur ueber die Modell-Tests (TestRejectOnTheBoardWritesTheReason, ...Cancelled...), kein interaktives Terminal hier. DoD 6: notes/A0C081.md, Agent-Block, Skill im Baum."
review-verdict: "Der Diff erfuellt alle sechs DoD-Punkte; per Hand nachgeprueft: Exit 2/3/5, duplicate+related, archive-Liste, show, restore nimmt die Felder ab, Sonderzeichen und Umlaute im Grund werden korrekt als YAML geschrieben. go vet und go test ./core/... ./internal/... -race gruen. Keine Defekte, die zurueckschicken rechtfertigen. Vor dem Annehmen: die uncommitteten optimize-Aenderungen committen (Gap 1). Gap 2 (mehrzeiliger Grund in der Archiv-Liste) ist kosmetisch - Mensch entscheidet, ob es hier noch rein soll."
review-round: 2
review-check: "0. Zuerst die liegengebliebenen optimize-Aenderungen committen: git add core/ticket/store.go internal/cli/reject_test.go internal/tui/model.go .jaira/tickets/01M59DFMZ1TKMZDZDJHBA0C081-ein-ticket-laesst-sich-mit-grund-ablehnen-jaira.md && git commit -m 'refactor(A0C081): one archive-target check for archive and reject'  1. go build -o /tmp/jaira ./cmd/jaira  2. In einem Wegwerf-Repo: git init, /tmp/jaira init, zwei Tickets mit /tmp/jaira create 'a' --goal g --context c --dod d anlegen (Handles merken: A und B)  3. /tmp/jaira reject A -> Exit 2, Meldung 'reject needs --why'  4. /tmp/jaira reject A --why x -> Exit 3, Meldung nennt 'jaira move A --to human --question'  5. /tmp/jaira reject A --why 'doppelt' --duplicate-of ZZZZZZ --force -> Exit 5  6. /tmp/jaira reject A --why 'doppelt' --duplicate-of B --force -> 'Rejected A (duplicate)'  7. /tmp/jaira archive -> unter der Datei die Zeile 'rejected duplicate: doppelt (<du>, <datum>)'  8. /tmp/jaira show A -> Zeilen 'related B' und 'rejected duplicate: doppelt ...'  9. /tmp/jaira restore <datei aus Schritt 7> und dann grep rejected .jaira/tickets/*A* -> keine Treffer  10. /tmp/jaira oeffnen, B markieren, R druecken, 'ueberholt' tippen, Enter -> Meldung 'Rejected B: ueberholt', B ist vom Board weg; dasselbe in einer offenen Detailansicht, dort bricht esc ab und das Ticket bleibt."
---

# Ein Ticket laesst sich mit Grund ablehnen: jaira reject, und das Archiv sagt warum

## Definition of Done

- [x] 'jaira reject <id> --why <text>' schreibt rejected.kind/why/by/at ins Frontmatter und verschiebt das Ticket nach .jaira/archive/; ohne --why Exit 2
  proof: internal/cli/reject.go newRejectCmd; core/ticket/store.go Store.Reject; tests TestRejectWritesTheReasonAndArchives, TestRejectWithoutWhyIsAUsageError
- [x] --duplicate-of <id> setzt kind=duplicate und verlinkt das Original; eine unbekannte id wird mit Exit 5 abgelehnt
  proof: core/ticket/store.go Store.Reject (LoadAnywhere, related); test TestRejectDuplicateOfLinksTheOriginal, TestRejectDuplicateOfAnUnknownTicketIsNotFound
- [x] Ein Agent wird mit Exit 3 abgewiesen und die Meldung nennt den Weg ueber human --question
  proof: internal/cli/reject.go !force -> ExitValidation needs_human, questionLane; test TestRejectWithoutForceSendsAnAgentToTheHumanLane
- [x] Auf dem Board laesst sich ein Ticket mit einer einzeiligen Begruendung ablehnen
  proof: internal/tui/model.go startReject/confirmReject, Taste R (Board + Detail); tests TestRejectOnTheBoardWritesTheReason, TestRejectCanBeCancelledFromAnOpenTicket
- [x] 'jaira archive' (Liste) und 'jaira show' zeigen rejected-kind und rejected-why; 'jaira restore' nimmt die rejected-* Felder wieder ab
  proof: internal/cli/archive.go archivedRejections; internal/cli/tickets.go rejectionLine/rejectionJSON; core/ticket/store.go unreject; tests TestArchiveListAndShowSayWhyATicketWasRejected, TestAnArchivedTicketIsNotShownAsRejected, TestRestoreTakesTheRejectionOff
- [x] Release-Note core/release/notes/<ID>.md, Agent-Block/jaira-Skill nennen reject; go test ./... -race gruen
  proof: core/release/notes/A0C081.md; core/board/announce.go + CLAUDE.md/AGENTS.md (TestAgentNoteNamesReject); .claude/skills/jaira/SKILL.md 'Taking a ticket off the board'; go test ./... -race exit 0
- [ ] Ein mehrzeiliger Grund wird beim Schreiben auf eine Zeile gefaltet; 'jaira archive' behaelt die Einrueckung
- [ ] --duplicate-of auf das Ticket selbst endet mit Exit 2, geprueft bevor etwas gestempelt oder geschrieben wird

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Schema: Felder rejected-kind/why/by/at in core/ticket/schema.go, Ticket.Rejected, Store.Reject (schreibt Felder, ruft Archive) und Store.Restore entfernt sie
- [x] CLI: internal/cli/reject.go - --why Pflicht (Exit 2), --kind wontfix|obsolete, --duplicate-of (Exit 5 bei unbekannter id, related setzen), ohne --force Exit 3 mit Weg ueber human --question
- [x] TUI: Taste R fragt eine Zeile Grund und ruft Store.Reject; x bleibt archivieren
- [x] Anzeige: jaira archive (Liste) und jaira show zeigen kind und why (Web-Portal: TPCKZB)
- [x] Release-Note, Agent-Block, jaira-Skill; go test ./... -race

## Progress
- **2026-10-19 06:27 · Alexander Sacharov** — parent session: jaira-a7
- **2026-10-19 06:28 · Alexander Sacharov** — decision (worker, settled): Art waehlt --kind wontfix|obsolete, Default wontfix; --duplicate-of setzt duplicate und verbietet --kind. Auf dem Board wird mit Default wontfix abgelehnt - eine Zeile Grund reicht (DoD 4).
- **2026-10-19 06:28 · Alexander Sacharov** — question 1 (an jaira-a7): Der Mechanismus aus VW1WB1 ist 'Interactive' (core/gate/gate.go:395) - nur das Board setzt ihn, jeder CLI-Aufruf gilt als Agent. Damit widersprechen sich DoD 1 (CLI reject klappt) und DoD 3 (Agent Exit 3). Optionen: 1) (Recommended) CLI reject ohne --force -> Exit 3 mit Hinweis auf human --question bzw. das Board; ein Mensch in der Shell nimmt --force, genau wie heute bei move aus human. 2) CLI reject nur ueber das Board, kein --force-Weg - DoD 1 gilt dann fuer den Store-Befehl, den das Board aufruft. 3) CLI reject immer erlaubt, DoD 3 faellt weg - Agenten werden nur per Skill-Regel abgehalten.
- **2026-10-19 06:28 · Alexander Sacharov** — question 2 (an jaira-a7): Welche Taste lehnt auf dem Board ab? 1) (Recommended) eigene Taste R (in der Detailansicht und auf dem Board frei), fragt nach einer Zeile Grund; x bleibt reines Archivieren - Aufraeumen und Ablehnen bleiben getrennt, wie der Kommentar model.go:1532 es fuer x verlangt. 2) x fragt nach einem Grund, leer = nur archivieren - eine Taste weniger, aber jedes Archivieren bekommt eine Frage dazwischen.
- **2026-10-19 06:31 · Alexander Sacharov** — decision (Alex via jaira-a7, 19.10.2026): Q1 -> Option 1: CLI reject ohne --force -> Exit 3, Meldung nennt human --question und das Board; ein Mensch in der Shell nimmt --force, wie bei move aus human. Q2 -> Option 1: eigene Taste R fragt nach einer Zeile Grund; x bleibt reines Archivieren.
- **2026-10-19 06:32 · Alexander Sacharov** — question 3 (an jaira-a7): Kontext/DoD verlangen einen Block 'rejected:' mit kind/why/by/at. Das Format kennt aber nur flache Top-Level-Keys: core/ticket/schema.go:69-72 begruendet, warum outcome-what/why/resolves flach sind (Merge-Driver loest pro Top-Level-Key, ein Feld = eine Zeile), und der Writer (core/ticket/frontmatter.go SetScalar/SetList) kann keine verschachtelten Mappings schreiben. Optionen: 1) (Recommended) flache Keys rejected-kind, rejected-why, rejected-by, rejected-at - wie outcome-*, kein neuer Writer, 'jaira set' und Merge-Driver funktionieren sofort. 2) verschachtelter Block 'rejected:' wie im Ticket - braucht einen neuen Block-Writer im Frontmatter-Code und weicht von der Regel in schema.go ab.
- **2026-10-19 06:33 · Alexander Sacharov** — decision (Alex via jaira-a7): Q3 -> Option 1: flache Keys rejected-kind, rejected-why, rejected-by, rejected-at wie outcome-*. DoD 1/5 'rejected.kind/why/by/at' und 'rejected-Block' sind damit diese vier Keys.
- **2026-10-19 06:35 · Alexander Sacharov** — scope (Teamlead via jaira-a7): Web-Portal-Teil von DoD 5 wandert nach TPCKZB (Web-Filter, parallel auf feat/tpckzb) - zwei Branches an core/web gleichzeitig wuerden kollidieren. TPCKZB zeigt rejected-* an, wenn vorhanden. Hier wird core/web nicht angefasst.
- **2026-10-19 06:36 · Alexander Sacharov** — decision (worker): --duplicate-of verlinkt ueber related (bestehendes Feld, von beiden Seiten sichtbar) statt eines fuenften Keys - Alex hat genau vier rejected-* Keys entschieden. Das Original wird mit LoadAnywhere gesucht, damit ein Duplikat eines schon abgelegten Tickets geht. Doc.Remove (frontmatter.go) ist neu: restore nimmt die Keys ganz ab statt sie leer stehen zu lassen. DoD 1-3 in einem Schritt gebaut, weil sie ein Befehl sind.
- **2026-10-19 06:37 · Alexander Sacharov** — Alex (via jaira-a7) hat DoD 1-3 abgenommen; nicht committen, Alex committet selbst. Weiter mit DoD 4 (Taste R).
- **2026-10-19 06:43 · Alexander Sacharov** — DoD 4: R im Board und in der Detailansicht; leere Begruendung laesst den Prompt offen (ohne Grund = x). Art auf dem Board immer wontfix. Hilfe: x und R teilen eine Zeile, weil TestHelpExplainsEveryCardMark bei 60 Zeilen Hoehe sonst '@name' abschneidet. README/COMMANDS.md schon hier ergaenzt, weil TestEveryCommandIsDocumented sonst rot ist.
- **2026-10-19 06:44 · Alexander Sacharov** — Alex (via jaira-a7) hat DoD 4 abgenommen. Weiter mit DoD 5 (archive-Liste, show).
- **2026-10-19 06:45 · Alexander Sacharov** — DoD 5: archive --json behaelt 'archived' als Namensliste und bekommt daneben 'rejected' (Datei -> kind/why/by/at), damit bestehende Leser nicht brechen. show --json hat jetzt immer einen 'rejected'-Block (leer, wenn nicht abgelehnt), wie 'outcome'.
- **2026-10-19 06:46 · Alexander Sacharov** — Alex (via jaira-a7) hat DoD 5 abgenommen; ab jetzt keine Diff-Stopps pro Punkt mehr, Ticket bis signoff tragen. Commit-Zeilen weiter an jaira-a7, Alex committet.
- **2026-10-19 06:50 · Alexander Sacharov** — critique (Runde 1): 2 Befunde - (1) archive.go archivedRejections -> ticket.LoadHeader statt eigenem ReadFile/ParseDoc/Decode; (2) Store.Reject: Archiv-Ziel vor dem Mutate pruefen, sonst bleibt bei Archive-Fehler ein abgelehntes Ticket auf dem Board. Stehen gelassen: questionLane in reject.go (kein bestehender Helfer fuer die Frage-Lane gefunden); show --json mit immer vorhandenem rejected-Block (wie outcome); TUI-R ohne StampCommits (wie x/archiveSelected).
- **2026-10-19 06:51 · Alexander Sacharov** — in-progress (Runde 2): beide critique-Befunde umgesetzt - archivedRejections nimmt ticket.LoadHeader; Store.Reject prueft das Archiv-Ziel vor dem Mutate (Test TestRejectIntoAnOccupiedArchiveWritesNothing). go test ./internal/... ./core/... -race exit 0.
- **2026-10-19 06:56 · Alexander Sacharov** — critique (Runde 2): beide Befunde aus Runde 1 umgesetzt (4e6d036e) - nichts mehr. Beobachtet, nicht als Befund (bricht keine DoD): Store.Reject ruft jetzt s.resolve vor Mutate; fuer ein Ticket, das nur auf seinem Ref liegt, kommt damit 'not found' statt Mutates 'pull it first' - die CLI laedt vorher mit s.Load und meldet das schon dort.
- **2026-10-19 07:04 · Alexander Sacharov** — review (fresh context): DoD 1-6 erfuellt, per Hand geprueft, Tests -race gruen. Offen fuer den Menschen: optimize-Aenderungen (store.go archiveTarget, model.go, reject_test.go) sind uncommittet - Commit-Zeile steht in review-check Schritt 0. Kosmetisch: mehrzeiliger --why bricht die Einrueckung in 'jaira archive'.
- **2026-10-19 07:32 · Alexander Sacharov** — decision (Alex via jaira-a7): beide review-gaps hier beheben - (1) ein mehrzeiliger --why wird beim Schreiben auf eine Zeile gefaltet, damit 'jaira archive' eingerueckt bleibt; (2) --duplicate-of auf das Ticket selbst ist Exit 2 und wird vor dem Commit-Stempeln geprueft. Eine kombinierte Commit-Zeile (optimize + diese Fixes) an jaira-a7.
