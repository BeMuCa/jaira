---
id: 01M1BNSDH7ZSKJ1GEHEMHREQJR
title: Jede TUI-Ansicht bricht lange Zeilen an der Terminalbreite um
status: done
ready: true
creator: BeMuCa
goal: "Fliesstext in jedem TUI-Fenster bricht an der verfuegbaren Breite um, statt rechts abgeschnitten zu werden"
context: "Berk am 31.08.: bei verkleinertem Terminalfenster sind Saetze zum Teil nicht lesbar, und nach rechts scrollen geht nicht. internal/tui hat wrap() (view.go:1236) und wrapHints() (view.go:972), genutzt u.a. in edit.go:192, signoff.go:44, view.go:543/587. Viele Ansichten nutzen stattdessen truncate(): browse.go (206,230,236,238,295,315,318,320), defaultboard.go (176,191,208,217,219), lanes.go (676,687,695,701,730), edit.go (152,167,182), dropboard.go (150). Bewusste Einzeiler (Kartentitel in Board-Spalten, Statuszeilen, Tabellenzellen) duerfen gekuerzt bleiben; Fliesstext (goal, context, outcome, review, notes, question, Fehlermeldungen, Pfade in Meldungen) muss umbrechen. Ein Explore-Audit ueber alle 17 Dateien laeuft; sein Befund gehoert in den Plan."
definition-of-done: In jedem TUI-Fenster ist bei schmalem Terminal jeder Satz vollstaendig lesbar (umgebrochen statt abgeschnitten); bewusste Einzeiler-Kuerzungen bleiben begruendet stehen; Tests fuer schmale Breiten; go test ./... -race gruen
blocked-by: []
commits: []
created-at: 2026-08-31T10:28:08Z
updated-at: 2026-08-31T16:21:34Z
claimed-by: EE-3NX6GL3-2143723
claimed-at: 2026-08-31T10:28:19Z
updated-by: BeMuCa
assignee: BeMuCa
outcome-what: "Runde 1: wrap() bricht ueberlange Woerter hart; wrapLines() erhaelt Absaetze/Einrueckung; 18 Fliesstext-Stellen und 7 Hint-Footer in 9 Dateien umgestellt; edit-Footer kompensiert Hoehe. Runde 2 (Review-Funde): proof-Budgets tragen ihren 13er-Einzug (signoff + renderChecklist, letzteres vorbestehend); wrapLines laesst passende Zeilen unangetastet (git --stat-Spalten intakt); wrap bricht auch unter Mindestbreite hart statt ueberbreit zurueckzugeben; styleLines() stylt zeilenweise, weil lipgloss mehrzeilige Bloecke auf die breiteste Zeile polstert (Titel, proofs, Projektpfade liefen sonst per Polsterung raus); browse-Hoehenbudget kompensiert Hints; home-Meldung je Zeile zentriert; defaultboard-Klammerhinweis eigener Hint-Punkt"
outcome-why: "bei verkleinertem Terminal gibt es kein horizontales Scrollen - ein rechts abgeschnittener Satz ist unlesbar (Berks Meldung vom 31.08.); bewusste Einzeiler (Karten, Spaltenkoepfe, Statusbar, unfokussierte edit-Felder) kuerzen weiter, Begruendung je Stelle in den Notizen"
outcome-resolves: "Screen-Tests messen jetzt Zeilenbreite auf dem ungeklammerten Render (fingen zwei echte Fehler: proof-Off-by-7, lipgloss-Blockpolsterung); 10 Tests in wrap_test.go gruen; alle bestehenden Tests gruen; go test ./... -race komplett gruen bei geleertem Cache; gofmt nur bekanntes tickets.go; Binary neu gebaut, Stempel=HEAD"
review-summary: "wrap() bricht ueberlange Woerter hart (hardBreak) und faellt auch unter der Mindestbreite nie auf ueberbreite Zeilen zurueck; wrapLines() erhaelt Absaetze, laesst passende Zeilen unangetastet (git-stat-Spalten intakt) und begrenzt tiefe Einrueckung; styleLines() stylt zeilenweise gegen lipgloss-Blockpolsterung nach Praefixen; 18 Fliesstext-Stellen + 7 Hint-Footer in 9 TUI-Dateien auf Umbruch umgestellt, proof-Budgets tragen ihren Einzug, browse/edit kompensieren Hint-Hoehe; Screen-Tests messen Zeilenbreite ungeklammert"
review-gaps: "Runde 2: nichts entfernt (kein Duplikat, kein toter Code, kein Fluff, keine Kostenstelle im Fix-Diff). Gelassen wie in Runde 1, plus: styleLines nur an den vier Praefix-Stellen statt ueberall (Polsterung ohne Praefix ist unsichtbar und harmlos); dropboard/settings/home weiter ohne Hoehenkompensation (kein Budget vorhanden, Schnitt unten wie vor der Aenderung); fitwindow-Testtabelle deckt die geaenderten Screens weiterhin nicht ab - als bekannte Luecke notiert statt Scope auszuweiten."
review-check: "1. go build -o ~/.local/bin/jaira ./cmd/jaira - keine Ausgabe. 2. Terminal auf ~50 Spalten, jaira starten: Meldung und Tastenzeile der Projektauswahl brechen mehrzeilig um; das … auf den Board-Zeilen und der Versionszeile ist gewollt (Einzeiler). 3. Board oeffnen, ein Ticket MIT Commits oeffnen (z.B. QPJNQP): Titel, goal, context brechen um, die Commits-Statistik behaelt ihre |-Spalten; bei HREQJR erscheint die Blocker-Zeile 'Before this can start' umgebrochen. 4. Verbotene Aktion ausloesen (m auf eine Lane, die das Gate ablehnt): Fehlermeldung bricht um. 5. a / x / d / S: lange Pfade brechen hart um, Footer-Tasten alle sichtbar. 6. Ticket in signoff bei ~40 Spalten oeffnen (HREQJR hat einen langen proof am DoD-Punkt): jede proof-Zeile vollstaendig lesbar, nichts laeuft rechts raus. 7. e auf einem Ticket: Hint-Zeilen sichtbar, Editorbox laeuft nicht unten raus."
review-verdict: "accept (unabhaengiges Opus-Review, 2. Pass): jeder Fund aus Pass 1 als tatsaechlich gefixt verifiziert, styleLines an genau den noetigen Stellen und nirgends grundlos; gemessen mit lipgloss.Width auf dem ungeklammerten Render bei 30/40/60 Spalten - keine ueberbreite Zeile auf signoff, detail, browse, browse-results, dropboard, settings, projects, message, renderChecklist. Letzter Low-Severity-Fund (wrapLines ohne Indent-Schranke: Ein-Runen-Spalte bei tiefer Einrueckung) danach ebenfalls gefixt und mit Test gepinnt."
---

# Jede TUI-Ansicht bricht lange Zeilen an der Terminalbreite um

## Definition of Done

- [x] In jedem TUI-Fenster ist bei schmalem Terminal jeder Satz vollstaendig lesbar (umgebrochen statt abgeschnitten); bewusste Einzeiler-Kuerzungen bleiben begruendet stehen; Tests fuer schmale Breiten; go test ./... -race gruen
  proof: wrap_test.go (6 Tests) + bestehende Narrow-Tests gruen; go test ./... -race komplett gruen, gofmt nur bekanntes tickets.go

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Audit-Stellen verifizieren: jede Fundstelle lesen, bevor sie geaendert wird
  proof: view.go 784/864/870/1018/1046, edit.go 145-195, signoff.go 18-75, lanes.go 790-805, settings.go 70-90, browse.go 198-245+285-325, dropboard.go 140-195, defaultboard.go 170-222, home.go 375-395 gelesen; Audit bestaetigt
- [x] wrap() bricht ueberlange Woerter (Pfade) hart um - failing Test zuerst
  proof: TestWrapBreaksOverlongWords, TestWrapStillBreaksAtSpaces (wrap_test.go); hardBreak in view.go
- [x] Mehrzeiliger Text: Helfer der Absaetze erhaelt und je Zeile umbricht, mit Test
  proof: wrapLines in view.go; TestWrapLinesKeepsLineStructure
- [x] Fliesstext auf wrap umstellen: view.go 864/870/1018/1046, home.go 385, settings.go 82, lanes.go 799, signoff.go 30/31/66/70, browse.go 206/236, dropboard.go 150/191, defaultboard.go 176/217, edit.go 167
  proof: 9 Dateien umgestellt, git diff 101+/39-; TestRenderMessageWrapsAtNarrowWidth, TestProjectPathsWrapInsteadOfTruncate, TestSignOffWrapsChecklistText
- [x] Bewusste Einzeiler bleiben truncate, je Stelle begruendet (Karten, Spaltenkoepfe, Statuszeile, Menuezeilen, edit.go 182, signoff 64er-Breite)
  proof: Begruendungen als Notiz auf dem Ticket; edit.go 167/182, signoff Statuszeile+64er-Breite, Karten/Spalten/Statusbar unveraendert
- [x] Schmale-Breiten-Tests; go test ./... -race, gofmt, Binary neu bauen
  proof: go test ./... -race (Cache geleert) gruen; gofmt nur tickets.go; Binary 45f78f2 vcs.modified=false

## Progress
- **2026-08-31 10:31 · BeMuCa** — Audit (Explore-Agent, alle 17 tui-Dateien) verifiziert die Richtung: Fliesstext/Meldungen/Pfade sind truncated statt umgebrochen. Zwei Dinge musste man herausfinden: (1) wrap() (view.go:1236) bricht nur an Leerzeichen - ein langer Pfad ist EIN Wort und bleibt ungebrochen; deshalb braucht wrap() harten Umbruch fuer Woerter breiter als die Zeile, sonst bleibt genau der Pfad-Fall unlesbar. (2) wrap() plaettet Zeilenumbrueche (strings.Fields), also braucht mehrzeiliger Text (Lane-Prompt lanes.go:799, git-stat view.go:864) einen Umweg je Zeile. Bewusst NICHT umgestellt: einzeilige Karten/Spalten/Status/Menues (Layout, wrapHints deckt Hints ab), edit.go:182 unfokussierte Felder (Fokus zeigt den vollen Text gewrappt, sonst schiebt ein langer context den Editor aus dem Bild), signoff.go:44 min(w-14,64) (Lesebreite auf breiten Terminals; bei schmalen greift w-14 bereits).
- **2026-08-31 10:39 · BeMuCa** — Bewusst NICHT umgestellt, je Stelle: edit.go editBody-Titelzeile und unfokussierte Feldzeilen bleiben truncate - das Hoehenbudget (overhead = Felder + 6) rechnet mit einer Zeile pro Feld, und der volle Text ist im fokussierten Feld gewrappt sichtbar; signoff.go Statuszeile 'waiting on your sign-off' (26 Zellen, fix); signoff.go min(w-14,64) bleibt - Lesebreite auf breiten Terminals, bei schmalen greift w-14; Karten/Spaltenkoepfe/Statusbar/Menuezeilen/home-Boardzeilen/versionLine bleiben Einzeiler. Alle Tastatur-Hint-Footer (view projects, settings, browse x2, dropboard, defaultboard, home, edit) auf wrapHints umgestellt - jede Taste sichtbar statt abgeschnitten, wie pipeline.go es schon tat. edit.go: Editor-Box gibt pro zusaetzlicher Hint-Zeile eine Zeile ab, damit nichts unten rausschiebt.
- **2026-08-31 10:55 · BeMuCa** — Review (Opus, unabhaengig) lehnte ab - Funde verifiziert und uebernommen: signoff proof-Budget w-6 bei Einzug 13 ergab w+7-Zeilen (der gemeldete Bug selbst, auf dem genannten Screen); renderChecklist proof hat denselben vorbestehenden Off-by-7; wrapLines kollabierte Whitespace auch bei passenden Zeilen (git --stat-Spalten kaputt auf jeder Breite); die zwei neuen Screen-Tests pruefen keine Zeilenbreite und waren dadurch vakuum; wrap gab bei width<=8 ueberbreite Zeilen zurueck; home-Meldung war als Block statt je Zeile zentriert; defaultboard-Hint 'esc back (...)' 45 Zellen unteilbar; browse-Hoehenbudget ignorierte mehrzeilige Hints. Bewusst offen gelassen (review-gaps): fitwindow-Testtabelle deckt die geaenderten Screens weiter nicht ab; dropboard/settings/home ohne Hoehenkompensation (kein Budget vorhanden, Schnitt unten wie vorher).
- **2026-08-31 11:00 · BeMuCa** — Beim Fixen der Review-Funde vom verschaerften Test gefangen, dritte Erkenntnis: lipgloss Style.Render POLSTERT einen mehrzeiligen Block rechts auf seine breiteste Zeile. Nach einem Praefix (Handle+Titel, '      proof:', '    '+Pfad) laeuft die gepolsterte Zeile ueber die Pane hinaus, obwohl wrap() korrekt budgetiert hat. Deshalb styleLines() in view.go: stylt zeilenweise statt als Block. Betroffen waren detail-Titel, signoff-Titel, beide proof-Stellen, Projektpfade. Merkregel: nie ein mehrzeiliges wrap()-Ergebnis durch Style.Render schicken, wenn davor ein Praefix steht.
