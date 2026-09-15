---
id: 01M2GPAV3W200QAWWGNQ1K9KZS
title: "Ein Board, das es schon gibt, bekommt eine geaenderte Lane nie zu sehen"
status: in-progress
ready: true
creator: Alexander Sacharov
goal: "Eine Korrektur an einer ausgelieferten Lane erreicht auch die Boards, die es schon gibt - ohne dass jemand auf jedem Rechner eine Zeile von Hand loescht."
context: |-
  74VM40 hat 'logbook-on-entry' aus den mitgelieferten Lanes entfernt, weil ein Move nach done 49 fremde fertige Tickets ins Logbuch fegte (Issue #6). Auf jedem Board, das vor dieser Aenderung entstanden ist, passiert das weiterhin.

  Die Ursache steht in core/lane/lane.go:479. Load() liest bei ProjectLanesActive NUR das Lane-Verzeichnis des Projekts und legt die eingebauten Lanes nie darunter. Eine Lane-Datei, die einmal geschrieben wurde, ist damit fuer immer die Wahrheit - auch wenn das Binary sie laengst korrigiert hat.

  Nachgestellt, zweimal auf diesem Board: am 11.09. und wieder am 14.09. hat ein Move nach done fremde fertige Tickets mitgenommen. Beide Male hat ein Mensch die Zeile von Hand aus .jaira/lanes/done.md geloescht - und weil .jaira/lanes/ in gitignore steht, pro Checkout einzeln. Am 14.09. waren es zwei Verzeichnisse.

  Dazu kommt, dass die Release-Notiz das Gegenteil behauptet: sie sagt glatt 'Finishing a ticket no longer files anything' und erklaert danach nur, wie man das alte Verhalten zurueckholt. Wer sie liest und ein aelteres Board hat, glaubt etwas Falsches ueber sein eigenes Werkzeug.

  Der Befund stammt aus der review-Lane von 74VM40 am 2026-09-14, die ihn ausdruecklich als Grund genannt hat, 74VM40 NICHT abzunehmen.

  Die Spannung, die die Plan-Lane aufloesen muss: 'ein Board ist sein Lane-Verzeichnis' ist eine bewusste Entscheidung (Commit 743737f, 27.08.) - eine vom Nutzer geschriebene Lane-Datei IST die Lane, keine Ueberschreibung. Eine Migration darf das nicht stillschweigend umdrehen und darf niemandem eine selbst geaenderte Lane unter den Fuessen wegziehen. Gesucht ist der Weg dazwischen: eine Korrektur erreicht das Board, eine Anpassung des Nutzers ueberlebt sie.
definition-of-done: "Ein Board, dessen done.md noch 'logbook-on-entry: true' traegt, fegt beim naechsten Move nach done keine fremden fertigen Tickets mehr ins Logbuch. Nachgestellt an einem Board-Fixture, das mit der alten Lane-Datei angelegt wurde."
tags:
  - gates
  - cli
blocked-by: []
related:
  - 01M28MHSDBABYVD8785A74VM40
commits: []
created-at: 2026-09-14T19:29:33Z
updated-at: 2026-09-15T07:24:07Z
assignee: "Alexander Sacharov"
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-74981
claimed-at: 2026-09-15T07:18:08Z
outcome-what: "core/lane/corrections.go: benannte, einmalige, feldgenaue Lane-Korrektur im Ladepfad; entfernt 'logbook-on-entry: true' aus einem done.md, das erkennbar die von 2ecc670 ausgelieferte Datei ist, meldet es und merkt sich das in .jaira/lanes/corrections. Eine selbstgeschriebene Lane wird gemeldet, nie editiert. Tests: core/lane/corrections_test.go (6), core/move/oldboard_test.go (end-to-end). Eine Zeile in core/release/NOTES.md unter ## Unreleased."
outcome-why: "Boards von vor 9ad7aa9 fegen beim Move nach done weiterhin fremde fertige Tickets ins Logbuch (Issue #6), weil Load bei ProjectLanesActive nur das Lane-Verzeichnis liest - eine einmal geschriebene Lane-Datei ist fuer immer die Wahrheit. Die Handarbeit-Anweisung der 0.1.4-Notiz greift nicht, weil .jaira/lanes/ gitignored ist und die Zeile pro Checkout einzeln entfernt werden muesste."
outcome-resolves: "DoD 1-6 abgehakt mit Proof; go test ./... -race RC=0; alle drei Faelle zusaetzlich von Hand mit gebautem Binary nachgestellt"
review-summary: "core/lane/corrections.go:applyCorrections meldet die Korrektur nur als Set.Warnings-Eintrag; internal/cli/root.go:291 verwirft alle Lane-Warnungen unter --json, und rund fuenfzehn weitere lane.Load-Aufrufer (internal/cli/mergedriver.go:55, tags.go:320, links.go:43, validate.go:35, lanes.go:110, checklist.go:137, resume.go:84, internal/tui/browse.go:157) lesen .Warnings gar nicht. Die Korrektur ist einmalig und schreibt den Marker im selben Aufruf: faellt sie in einem dieser Aufrufe an - beim Agenten-'jaira next --json' oder im git-Merge-Driver - wird die done.md editiert und niemand erfaehrt es je. Stattdessen den Says-/Skipped-Satz in applyCorrections direkt auf os.Stderr schreiben, so wie nudgeIfStale in internal/cli/update.go:37 es fuer genau diesen Fall schon tut und im Kommentar begruendet: 'it must reach the terminal regardless of --json - stdout is reserved for the payload an agent parses'. Das ist DoD 3: der Test beweist, dass der Satz erzeugt wird, nicht dass er ankommt."
---

# Ein Board, das es schon gibt, bekommt eine geaenderte Lane nie zu sehen

## Definition of Done

- [x] Ein Board, dessen done.md noch 'logbook-on-entry: true' traegt, fegt beim naechsten Move nach done keine fremden fertigen Tickets mehr ins Logbuch. Nachgestellt an einem Board-Fixture, das mit der alten Lane-Datei angelegt wurde.
  proof: core/move/oldboard_test.go:38 TestMoveIntoDoneOnAnOldBoardFilesNothing
- [x] Eine vom Nutzer selbst geaenderte Lane ueberlebt die Migration unveraendert: was er geschrieben hat, wird nicht zurueckgesetzt. Die Entscheidung von 743737f - ein Board ist sein Lane-Verzeichnis - bleibt gueltig.
  proof: core/lane/corrections_test.go:93 TestCorrectionRunsOncePerBoard; core/lane/corrections_test.go:134 TestCorrectionLeavesALaneSomebodyWroteAlone
- [x] Der Nutzer erfaehrt, dass eine seiner Lanes von einer Korrektur betroffen ist, statt es an seinem Verhalten zu merken. Nachgestellt an dem Board-Fixture aus Punkt 1.
  proof: core/lane/corrections.go:144 applyCorrections writes to os.Stderr (correctionsOut:121); core/lane/corrections_test.go:244 TestCorrectionSpeaksOnStderrAndNotOnStdout, :55 TestCorrectionRemovesTheDoorwayFromAnOldBoard; core/move/oldboard_test.go:63
- [x] Die Zeile in core/release/NOTES.md, die 'Finishing a ticket no longer files anything' behauptet, stimmt danach fuer alle Boards - oder sie sagt, fuer welche sie nicht gilt und was zu tun ist.
  proof: core/release/NOTES.md:18
- [x] Eine Zeile in core/release/NOTES.md unter ## Unreleased.
  proof: core/release/NOTES.md:18
- [x] Eine Korrektur fasst nur eine Lane an, die erkennbar die ausgelieferte ist, die sie zu korrigieren behauptet - etwa weil die Datei einer ausgelieferten Fassung entspricht oder sich nur in genau dem Feld unterscheidet, um das es geht. Eine Lane, die ein Mensch selbst geschrieben hat, wird gemeldet und nicht angefasst, auch wenn sie dieselbe id traegt. Nachgestellt mit einem selbstgeschriebenen done.md, das absichtlich logbook-on-entry: true fuehrt: es bleibt unveraendert, und der Mensch erfaehrt davon.
  proof: core/lane/corrections.go:212 correction.recognises; core/lane/corrections_test.go:134 TestCorrectionLeavesALaneSomebodyWroteAlone, :169 TestCorrectionLeavesTodaysLanePlusTheLineAlone

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] read migrateLegacy (core/lane/lane.go:645) and stampCreatorLine (core/lane/share.go:60) - the two precedents for a once-only, surgical edit of a board's lane files
- [x] design the correction record: an embedded list of named corrections (lane id, the defect, the one field it removes, the sentence the user is told) plus an 'applied' marker file beside 'order', so each correction touches a board once and never again
- [x] write the board fixture: a board whose done.md still carries 'logbook-on-entry: true', built the way an old board was
- [x] failing test: a move into done on that fixture files nothing and leaves the other finished tickets standing
- [x] failing test: what the user wrote survives - a changed prompt/description is untouched, and 'logbook-on-entry: true' put back by hand after the correction ran stays put
- [x] failing test: the correction is reported once, naming the file it changed and how to keep the old behaviour
- [x] implement core/lane/corrections.go: the embedded list, dropFrontmatterLine beside stampCreatorLine, applyCorrections(root) called from Load next to migrateLegacy, marker via readIDList/writeIDList
- [x] run go test ./... -race and replay the fixture by hand with a built binary
- [x] one line in core/release/NOTES.md under ## Unreleased: what an older board does now, replacing the 'remove that line by hand' instruction of the closed 0.1.4 section
- [x] critique loop: the correction's report goes to os.Stderr from applyCorrections, not into Set.Warnings, so it survives --json and the callers that drop Warnings

## Progress
- **2026-09-15 07:00 · Alexander Sacharov** — pre-process: warum der Plan so aussieht.

Verworfen: die Builtins beim Laden unter das Board legen. Das dreht 743737f um - eine Lane-Datei waere dann keine Lane mehr, sondern eine Ueberschreibung, und jede vom Nutzer entfernte Zeile kaeme still zurueck.

Verworfen: die Korrektur an 'jaira update' haengen. Der Pfad existiert (internal/cli/update.go, release.Stamped + nudgeIfStale), aber die DoD verlangt 'beim naechsten Move nach done' - also muss die Korrektur im Ladepfad sitzen, nicht in einem Befehl, den jemand aufrufen muss.

Verworfen: die ganze Lane-Datei aus dem Builtin neu schreiben (Export). Das setzt einen geaenderten Prompt oder eine geaenderte Beschreibung mit zurueck. DoD 2 verbietet genau das.

Gewaehlt: eine benannte Korrektur, die genau ein Feld entfernt, einmal pro Board. Vorbild ist migrateLegacy (lane.go:645) - eine Migration, die sich an einer Datei merkt, dass sie gelaufen ist. Und stampCreatorLine (share.go:60), das ausdruecklich eine Zeile einfuegt statt das Frontmatter zu parsen und neu zu schreiben, damit nichts Unbekanntes verlorengeht. Die Korrektur entfernt spiegelbildlich eine Zeile.

Das Einmal-Marker ist der Punkt, an dem DoD 2 haelt: wer 'logbook-on-entry: true' nach der Korrektur bewusst wieder hinschreibt, behaelt es - die Korrektur fragt das Board nicht noch einmal. Format wie 'order' und 'removed': eine Zeile pro Korrektur-Id, gelesen mit readIDList.

Warnung: lanes.Warnings werden in internal/cli/root.go:291 und internal/tui/model.go:351 schon ausgegeben - DoD 3 braucht keinen neuen Kanal, nur den richtigen Satz.

NOTES.md: die 0.1.4-Zeile bleibt stehen. 0.1.4 ist getaggt, also geschlossene Historie, und Alex hat am 15.09. entschieden, sie nicht anzufassen (74VM40). DoD 4 und 5 fallen deshalb auf eine einzige neue Zeile unter ## Unreleased zusammen, die sagt, was ein aelteres Board ab jetzt tut - womit die Handarbeit-Anweisung der alten Zeile gegenstandslos wird.
- **2026-09-15 07:02 · Alexander Sacharov** — Neues Kriterium 6 kam waehrend der in-progress-Lane dazu und verengt den Plan (Form bleibt, die Praezedenzfaelle und Tests bleiben): Eine Korrektur darf nur eine Lane anfassen, die erkennbar die ausgelieferte ist - die Datei entspricht einer ausgelieferten Fassung oder unterscheidet sich von ihr nur in genau dem Feld, das korrigiert wird. Eine selbst geschriebene Lane mit derselben id wird gemeldet, nie editiert. Der 'applied'-Marker deckt das NICHT ab: er verhindert die Wiederholung, nicht dass die Korrektur beim ersten Mal falsch ist. Erkennung muss aus dem Inhalt kommen: 'creator:' taugt nicht dafuer - stampCreatorLine (core/lane/share.go:60) setzt es, aber auf diesem Board tragen es nur critique.md, optimize.md und testing.md (aus einem Katalog uebernommen); die ausgelieferten Lanes haben die Zeile gar nicht. Zusaetzlicher Test: ein handgeschriebenes done.md, das absichtlich logbook-on-entry: true fuehrt, bleibt unveraendert und der Mensch erfaehrt, dass die Korrektur uebersprungen wurde und warum.
- **2026-09-15 07:05 · Alexander Sacharov** — in-progress: Erkennung fuer Kriterium 6 kommt aus dem Dateiinhalt, nicht aus dem geparsten Lane. corrections.go haelt die ausgelieferte done.md von 2ecc670 woertlich als Konstante doneDoorway; korrigiert wird nur eine Datei, die ihr - ohne die logbook-on-entry-Zeile - genau gleicht (Zeilenenden/Trailing-Space normalisiert). Bewusst NICHT aufgenommen: die heutige Builtin-Fassung (9ad7aa9). Sonst waere 'heutiges done.md + von Hand hinzugefuegte Zeile' erkannt worden - genau der Fall, den Kriterium 6 schuetzt, denn diese Kombination hat nie ein Build ausgeliefert. Verworfen: Vergleich ueber lanesEquivalent - ein geparster Vergleich sieht nur die Felder, die dieses Binary kennt, also kaemen ein umgeschriebener Prompt-Text (gleiche Felder, anderer Inhalt zaehlt zwar, aber ein Kommentar oder ein Feld einer neueren Version nicht) durch und die Datei wuerde trotzdem editiert.
- **2026-09-15 07:11 · Alexander Sacharov** — in-progress fertig. Von Hand nachgestellt mit gebautem Binary auf einem Scratch-Board (nicht nur im Test): (a) altes done.md aus 2ecc670, Marker geloescht -> 'jaira list' korrigiert und meldet es; (b) Zeile danach von Hand wieder hineingeschrieben -> bleibt stehen, keine Meldung; (c) selbstgeschriebenes done.md mit derselben Zeile, Marker geloescht -> unveraendert, Meldung 'is not the lane jaira shipped'. Wichtig fuer den naechsten Leser: ein frisch angelegtes Board bekommt den Marker sofort beim ersten Load geschrieben. Ein Test-Fixture, das erst Load() aufruft und danach die alte done.md hinlegt, muss .jaira/lanes/corrections wieder loeschen - sonst laeuft die Korrektur nie und der Test ist gruen aus dem falschen Grund. Genau das macht oldBoard() in corrections_test.go und der Fixture-Aufbau in core/move/oldboard_test.go. Nicht angefasst (ausserhalb dieser Lane): das .jaira/lanes/ dieses Repos traegt die Zeile evtl. noch - die Korrektur greift dort, sobald das gebaute Binary installiert ist, nicht vorher.
- **2026-09-15 07:16 · Alexander Sacharov** — critique: eine Findung, Zustellweg der Meldung. applyCorrections gibt Says/Skipped als Set.Warnings zurueck. internal/cli/root.go:291 unterdrueckt Lane-Warnungen unter --json, und die meisten lane.Load-Aufrufer (mergedriver.go:55, tags.go:320, links.go:43, validate.go:35, lanes.go, checklist.go:137, resume.go:84, tui/browse.go:157) lesen .Warnings ueberhaupt nicht. Weil die Korrektur einmalig ist und den Marker im selben Load schreibt, ist die Meldung danach fuer immer weg - und der erste Load nach einem Update ist auf diesem Board wahrscheinlich ein Agenten-Aufruf mit --json oder der Merge-Driver. Fix: in applyCorrections direkt auf os.Stderr schreiben, Vorbild und Begruendung stehen schon in internal/cli/update.go:37 (nudgeIfStale). Betrifft DoD 3 - der Unit-Test prueft die Erzeugung des Satzes, nicht seine Zustellung.

Geprueft und NICHT beanstandet, damit es nicht noch einmal aufgemacht wird: (1) doneDoorway stimmt byteweise mit core/lane/builtin/50-done.md bei 2ecc670 ueberein, und zwischen 2ecc670 und 9ad7aa9 gab es genau eine Fassung der Datei - die Shipped-Liste ist vollstaendig. (2) Export kopiert die Lane-Datei verbatim (core/lane/share.go:80), der Byte-Vergleich ist also tragfaehig. (3) Der corrections-Marker wird bewusst auch auf frischen Boards geschrieben: genau das schuetzt eine spaeter von Hand gesetzte Zeile davor, entfernt zu werden. (4) Platzierung in Load ist richtig - 'jaira update' waere ein Befehl, den jemand ausfuehren muss, und genau daran ist die 0.1.4-Notiz gescheitert. (5) Die Korrektur-Tabelle mit einer Zeile ist keine Spekulation, sondern eine append-only Migrationsliste.
- **2026-09-15 07:22 · Alexander Sacharov** — critique-Schleife erledigt: der Satz der Korrektur geht jetzt direkt auf os.Stderr (correctionsOut, nil = os.Stderr zur Schreibzeit), applyCorrections gibt nichts mehr zurueck und Load haengt nichts mehr an Warnings. Warum nicht beides: root.go:291 druckt Warnings ohne --json bereits aus, also haette jede nicht-json-Ausgabe die Zeile doppelt gezeigt. Gepruefte Sorge, die kein Problem ist: die TUI ist der einzige Warnings-Leser, der sie selbst rendert (model.go:351) - sie laedt das Board aber in tui.New (model.go:290), bevor tea.NewProgram laeuft, die Zeile landet also vor dem Alternate Screen und nicht mitten in einem Frame. Verworfen: ein eigenes Feld Set.Corrections neben Warnings - dasselbe Zustellproblem, nur mit einem zweiten Kanal, den wieder niemand liest. correctionsOut ist bewusst nil-default statt 'io.Writer = os.Stderr': sonst friert die Variable beim Paket-Init den echten Deskriptor ein und ein Test, der os.Stderr auf eine Pipe legt (Vorbild captureStdio, internal/cli/update_test.go:44), sieht nichts.
