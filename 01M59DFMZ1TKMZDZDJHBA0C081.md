---
id: 01M59DFMZ1TKMZDZDJHBA0C081
title: "Ein Ticket laesst sich mit Grund ablehnen: jaira reject, und das Archiv sagt warum"
status: review
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
updated-at: 2026-10-19T08:20:25Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-35137
claimed-at: 2026-10-19T08:18:23Z
mode: conversational
outcome-what: "Neuer Befehl jaira reject <id> --why [--kind wontfix|obsolete] [--duplicate-of] [--force] schreibt rejected-kind/why/by/at und archiviert ueber Store.Archive; Taste R auf dem Board; archive-Liste und show zeigen den Grund; restore nimmt die Felder ab (Doc.Remove)."
outcome-why: "Ein abgelehntes Ticket sah im Archiv genauso aus wie ein weggeraeumtes; Monate spaeter liess sich nicht sagen, ob es Duplikat, verworfen oder ueberholt war."
outcome-resolves: "DoD 1-3: internal/cli/reject.go + Store.Reject, Tests in internal/cli/reject_test.go; DoD 4: modeReject in internal/tui/model.go, internal/tui/reject_test.go; DoD 5: archivedRejections/rejectionLine, unreject; DoD 6: release note, agent block, Skill, go test ./... -race gruen."
review-summary: "jaira reject <id> --why <text> --force schreibt rejected-kind/why/by/at (flache Keys) aufs Ticket, stempelt vorher die Commits und archiviert es; ohne --force Exit 3 mit Weg ueber human --question, ohne --why Exit 2, --duplicate-of wird vor jedem Schreiben aufgeloest (Exit 5 unbekannt, Exit 2 sich selbst) und setzt kind=duplicate plus related. Store.Reject faltet den Grund auf eine Zeile und prueft das Archiv-Ziel (neuer Helfer archiveTarget, den Archive mitbenutzt) bevor es mutiert. Taste R auf Board und Detailansicht fragt eine Zeile Grund, stempelt jetzt wie die CLI die Commits (Rueckgabe aus signoff) und lehnt mit wontfix ab; x archiviert weiter nur. jaira archive zeigt eine eingerueckte rejected-Zeile pro Datei (und rejected in --json), jaira show eine rejected-Zeile, restore nimmt die Keys ab. Release-Note, Agent-Block und Skill nennen reject."
round: 4
review-gaps: "Blockierend: nichts. (1) Die Aenderungen aus optimize, Runde 3 und Runde 4 (store.go, reject.go, reject_test.go, model.go, tui/reject_test.go, Release-Note) sind weiter uncommittet - ohne Commit fehlen sie im Merge, siehe review-check Schritt 0. (2) Kosmetisch, nicht DoD: restore laesst den related-Eintrag aus --duplicate-of stehen; R stempelt mit dem beim Druecken gehaltenen t (Mutate liest aber frisch, MergeCommits vereinigt, also kein Verlust). Der entfernte nil-Guard in confirmReject ist sicher: modeReject wird nur ueber startReject betreten, das rejecting setzt."
test-verdict: "pass (Runde 4): go clean -testcache; go test ./... -race RC=0; TestRejectOnTheBoardStampsCommits gruen und mit auskommentierter StampCommits-Zeile rot (FAIL) - die Pruefung haengt wirklich an der Aenderung. DoD 1-8 wie Runde 3."
review-verdict: "approve: DoD 1-9 im Diff belegt; go vet und go test ./... -race gruen; TestRejectOnTheBoardStampsCommits ohne die StampCommits-Zeile rot gesehen; per Binary nachgespielt: Selbst-Duplikat Exit 2, mehrzeiliger Grund als 'a b c' eingerueckt im Archiv, show zeigt rejected und den gestempelten Commit."
review-round: 4
review-check: "0. Committen: git add core/release/notes/A0C081.md core/ticket/store.go internal/cli/reject.go internal/cli/reject_test.go internal/tui/model.go internal/tui/reject_test.go .jaira/tickets/01M59DFMZ1TKMZDZDJHBA0C081-ein-ticket-laesst-sich-mit-grund-ablehnen-jaira.md && git commit -m 'fix(A0C081): R on the board stamps commits, reasons fold to one line, a self-duplicate is exit 2'  1. go build -o /tmp/jaira ./cmd/jaira  2. in einem leeren Test-Repo: git init, /tmp/jaira init, /tmp/jaira create Test --goal g --context c --dod d --assignee du  3. git commit --allow-empty -m 'fix(<ID>): tried'  4. /tmp/jaira reject <ID> --why x -> Exit 3, Meldung nennt human --question  5. /tmp/jaira reject <ID> --why x --duplicate-of <ID> --force -> 'cannot be a duplicate of itself', Exit 2, Ticket unveraendert  6. /tmp/jaira reject <ID> --why $'a\\nb' --force -> 'Rejected <ID> (wontfix)'  7. /tmp/jaira archive -> unter dem Dateinamen eine eingerueckte Zeile 'rejected wontfix: a b (...)'  8. /tmp/jaira show <ID> -> Zeile rejected und Zeile commits mit dem Commit aus Schritt 3  9. Board: /tmp/jaira oeffnen, auf einer Karte R druecken, Grund tippen, Enter -> Meldung 'Rejected ...', Karte ist weg; /tmp/jaira archive zeigt den Grund"
---

# Ein Ticket laesst sich mit Grund ablehnen: jaira reject, und das Archiv sagt warum

## Definition of Done

- [x] 'jaira reject <id> --why <text>' schreibt rejected.kind/why/by/at ins Frontmatter und verschiebt das Ticket nach .jaira/archive/; ohne --why Exit 2
  proof: internal/cli/reject.go newRejectCmd; core/ticket/store.go Store.Reject; tests TestRejectWritesTheReasonAndArchives, TestRejectWithoutWhyIsAUsageError
- [x] --duplicate-of <id> setzt kind=duplicate und verlinkt das Original; eine unbekannte id wird mit Exit 5 abgelehnt
  proof: internal/cli/reject.go --duplicate-of via LoadAnywhere vor dem Stempeln (Exit 5); core/ticket/store.go Store.Reject setzt duplicate + related; tests TestRejectDuplicateOfLinksTheOriginal, TestRejectDuplicateOfAnUnknownTicketIsNotFound
- [x] Ein Agent wird mit Exit 3 abgewiesen und die Meldung nennt den Weg ueber human --question
  proof: internal/cli/reject.go !force -> ExitValidation needs_human, questionLane; test TestRejectWithoutForceSendsAnAgentToTheHumanLane
- [x] Auf dem Board laesst sich ein Ticket mit einer einzeiligen Begruendung ablehnen
  proof: internal/tui/model.go startReject/confirmReject, Taste R (Board + Detail); tests TestRejectOnTheBoardWritesTheReason, TestRejectCanBeCancelledFromAnOpenTicket
- [x] 'jaira archive' (Liste) und 'jaira show' zeigen rejected-kind und rejected-why; 'jaira restore' nimmt die rejected-* Felder wieder ab
  proof: internal/cli/archive.go archivedRejections; internal/cli/tickets.go rejectionLine/rejectionJSON; core/ticket/store.go unreject; tests TestArchiveListAndShowSayWhyATicketWasRejected, TestAnArchivedTicketIsNotShownAsRejected, TestRestoreTakesTheRejectionOff
- [x] Release-Note core/release/notes/<ID>.md, Agent-Block/jaira-Skill nennen reject; go test ./... -race gruen
  proof: core/release/notes/A0C081.md; core/board/announce.go + CLAUDE.md/AGENTS.md (TestAgentNoteNamesReject); .claude/skills/jaira/SKILL.md 'Taking a ticket off the board'; go test ./... -race exit 0
- [x] Ein mehrzeiliger Grund wird beim Schreiben auf eine Zeile gefaltet; 'jaira archive' behaelt die Einrueckung
  proof: core/ticket/store.go Store.Reject faltet why (strings.Fields); internal/cli/reject.go ebenso; test TestRejectFoldsAMultiLineReason
- [x] --duplicate-of auf das Ticket selbst endet mit Exit 2, geprueft bevor etwas gestempelt oder geschrieben wird
  proof: internal/cli/reject.go: LoadAnywhere + orig.ID == t.ID -> ExitUsage vor StampCommits; test TestRejectAsADuplicateOfItselfIsAUsageError
- [x] R auf dem Board stempelt vor dem Ablehnen die Commits wie 'jaira reject', damit ein abgelehntes Ticket seine commits-Liste hat
  proof: internal/tui/model.go confirmReject: StampCommits(t, m.gateEnv().DeriveCommits) vor Store.Reject; test TestRejectOnTheBoardStampsCommits (ohne die Zeile rot)

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
- **2026-10-19 07:33 · Alexander Sacharov** — in-progress (Runde 3, review-gaps): Grund wird in Store.Reject gefaltet (gilt fuer CLI und TUI) und in der CLI vor der Ausgabe; --duplicate-of wird in der CLI vor StampCommits aufgeloest (Exit 5 unbekannt, Exit 2 sich selbst), der Selbst-Check in Store.Reject bleibt fuer andere Aufrufer. DoD-2-Proof nachgezogen, weil die Exit-5-Pruefung von Store.Reject in die CLI gewandert ist. go clean -testcache; go test ./... -race RC=0.
- **2026-10-19 07:34 · Alexander Sacharov** — critique (Runde 3): liest nur die Antwort auf die review-gaps. Beide Gaps behoben, Fix bricht keine DoD - none.
- **2026-10-19 07:39 · Alexander Sacharov** — review (fresh context, Runde 2): approve. DoD 1-8 im Diff belegt; go vet, go test ./... -race gruen; end-to-end per Binary nachgespielt (Exit 3/2/5/2/0, archive, show, restore). Nicht blockierend: R auf dem Board stempelt keine Commits (wie x), CLI stempelt vor einem spaeter scheiternden Reject (wie archive), restore laesst related aus --duplicate-of stehen. Offen: optimize+Runde-3-Aenderungen uncommittet - Commit-Zeile in review-check Schritt 0.
- **2026-10-19 07:40 · Alexander Sacharov** — explain: Review sagt approve: DoD 1-8 sind belegt, Tests und End-to-end laufen gruen. Offen sind nur kleine Punkte und die noch nicht committeten Aenderungen aus optimize und Runde 3.
Warum: Ohne Commit fehlen store.go, reject.go, model.go und die Release-Note im Merge. Die kleinen Punkte (R stempelt keine Commits, restore laesst related stehen) brechen keine DoD und verhalten sich wie bei x und archive.
1) Ticket annehmen, vorher die Aenderungen mit der Zeile aus review-check Schritt 0 committen [accept]
2) Zurueckgeben: R auf dem Board soll auch Commits stempeln, damit abgelehnte Tickets eine commits-Liste haben [give back]
Empfehlung: 1 - die Luecken sind nicht blockierend und bei x und archive schon genauso, ein eigener Umbau waere Scope-Wachstum.
- **2026-10-19 08:15 · Alexander Sacharov** — decision: 2) Zurueckgeben: R auf dem Board soll auch Commits stempeln, damit abgelehnte Tickets eine commits-Liste haben
- **2026-10-19 08:15 · Alexander Sacharov** — Given back from signoff: 2) Zurueckgeben: R auf dem Board soll auch Commits stempeln, damit abgelehnte Tickets eine commits-Liste haben
- **2026-10-19 08:15 · Alexander Sacharov** — cockpit: given back to the old worker - tab alive, cache warm for 25 more min, context 236k
- **2026-10-19 08:17 · Alexander Sacharov** — in-progress (Runde 4, aus signoff zurueck): confirmReject stempelt best-effort wie die CLI und wie x es nicht tut; x bleibt unveraendert (nicht Teil der Rueckgabe). Test TestRejectOnTheBoardStampsCommits mit echtem git-Repo, ohne die Zeile rot gesehen. go clean -testcache; go test ./... -race RC=0.
- **2026-10-19 08:17 · Alexander Sacharov** — critique (Runde 4): eine Zeile, gleicher Aufruf wie settlePrepare/CLI - none.
- **2026-10-19 08:20 · Alexander Sacharov** — review (fresh context, Runde 3): approve. Runde-4-Rueckgabe (R stempelt Commits) im Diff, Test ohne die Zeile rot gesehen, -race gruen, CLI per Binary nachgespielt. Offen nur: Aenderungen seit 4e6d036e uncommittet - Commit-Zeile in review-check Schritt 0.
