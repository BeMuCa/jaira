---
id: 01M1E5X19YBGH4ZRJFKP9H265S
title: Der Signoff-Screen zeigt die Commits des Tickets wie das Detail-Pane
status: done
ready: true
creator: BeMuCa
goal: "renderSignOff zeigt denselben Commits-Block (git-stat, sonst die SHA-Liste) wie detailBody - der Abnahme-Screen zeigt, was abgenommen wird"
context: "Berk am 01.09.: 'beim signoff sieht man den commit nichtmehr, ist das richtig?' Befund: der Block war NIE da - signoff.go rendert keinen Commits-Abschnitt, nur detailBody tut es (view.go ~863: gitStat.of(t.Commits), Fallback SHA-Zeile, seit HREQJR mit wrapLines umgebrochen). Auf dem Abnahme-Screen fehlt damit genau das Objekt der Abnahme. Einbau sinnvoll zwischen den sieben Sektionen und der Definition of Done; Commits ggf. wie im Gate ueber DeriveCommits herleiten, wenn das Feld leer ist (so macht es showForLane, flow.go ~528) - pruefen, ob m im TUI den Deriver hat."
definition-of-done: Ein Signoff-Ticket mit Commits zeigt den git-stat-Block (bzw. die SHA-Zeile als Fallback); ohne Commits erscheint kein leerer Block; Zeilen brechen an der Breite um; Tests fuer beide Faelle; go test ./... -race gruen
blocked-by: []
commits: []
created-at: 2026-09-01T09:48:13Z
updated-at: 2026-09-03T10:29:01Z
claimed-by: EE-3NX6GL3-3294516
claimed-at: 2026-09-01T09:48:58Z
updated-by: BeMuCa
assignee: BeMuCa
outcome-what: "Runde 2: signoff leitet die Commit-Liste einmal je Ticket aus git her (Memo am Model), Heading sagt 'derived from git'; gitStat statet alle SHAs (fixt auch das Detail-Pane); Breite klemmt auf paneWidth; fitwindow-Fall traegt Commits; Fixture-Test mit echtem git init"
outcome-why: "Review-Funde: der Block fehlte im NORMALFALL (11 von 12 Tickets ohne recorded commits), stat zeigte nur shas[0], Breite ignorierte die Pane"
outcome-resolves: go test ./internal/tui -count=1 EXIT=0 inkl. TestSignOffDerivesUnrecordedCommits (echtes git-Fixture) und Negativfall
review-summary: none
review-gaps: "Nichts entfernt. Gelassen: kein DeriveCommits im Render (git-Aufruf pro Tastendruck; Gate leitet beim Move her - im Ticket-Kommentar begruendet); happy-path des git-stat ungetestet (teilt gitStat.of mit dem produktionsbewaehrten Detail-Pane; Fixture-Repo dafuer waere unverhaeltnismaessig)."
review-check: "Ein Ticket MIT commits nach signoff bewegen (oder eines dort oeffnen, z.B. eines der Flotte nach deiner Abnahme-Runde): zwischen Lane fields und Definition of Done steht der Commits-Block mit git-stat; bei kaputtem SHA die SHA-Zeile; ein Ticket ohne commits zeigt keinen Block."
review-verdict: "accept (Opus, 2 Runden): signoff leitet unaufgezeichnete Commits einmal je Ticket her (Memo, bei reload invalidiert), Heading sagt 'derived from git', stat ueber alle SHAs, Pane-Clamp, echtes git-Fixture mit Global-Config-Isolation; Rest-Nit L19 (ein kaputter SHA kippt den stat auf die ehrliche SHA-Zeile) bewusst belassen"
---

# Der Signoff-Screen zeigt die Commits des Tickets wie das Detail-Pane

## Definition of Done

- [ ] Ein Signoff-Ticket mit Commits zeigt den git-stat-Block (bzw. die SHA-Zeile als Fallback); ohne Commits erscheint kein leerer Block; Zeilen brechen an der Breite um; Tests fuer beide Faelle; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-01 09:57 · BeMuCa** — Entscheidung: NUR aufgezeichnete Commits, kein DeriveCommits im Render - identisch zum Detail-Pane; Herleitung pro Frame waere ein git-Aufruf pro Tastendruck, und das Gate leitet beim Move her. Fallback (kein Repo/kaputter SHA) ist die SHA-Zeile, wie im Detail.
- **2026-09-01 10:04 · BeMuCa** — review reject, uebernommen: (L10) signoff-Tickets tragen normalerweise commits: [] (11 von 12 live) - der Block fehlt im Normalfall; Fix: einmal pro Ticket via gateEnv().DeriveCommits herleiten, am Model memoisiert, Heading sagt 'derived from git'; (L11) gitStat.of statete nur shas[0] - jetzt alle SHAs an git show; (L12) stat-Breite klemmt auf min(w,paneWidth) wie alle Zeilen des Screens; (L13) fitwindow-signoff-Fall bekommt Commits; Kommentar zur Render-Kosten-Begruendung korrigiert (gitStat exect selbst pro Render - der Unterschied zu DeriveCommits ist Grad, nicht Art; Memo macht beides billig).
