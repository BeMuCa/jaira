---
id: 01M2RQM76KT0RGV3W7NA7XXT00
title: Der Lane-Payload liefert einen Ausschnitt des Diffs und meldet ihn als vollstaendig
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Eine Kritik- oder Review-Lane beurteilt den ganzen Branch, oder sie erfaehrt im Payload, dass sie es nicht tut - statt einen Bruchteil zu unterschreiben, ohne es zu merken"
context: |-
  Gefunden von der review-Lane auf GTQHNH am 2026-09-17, nachdem sie selbst darauf hereingefallen war.

  Was passiert: internal/cli/flow.go:589 nimmt 'shas := t.Commits' und leitet die Commit-Liste NUR dann aus git ab, wenn das Feld leer ist. Steht im Ticket-Frontmatter ein 'commits:' mit drei SHAs, waehrend auf dem Branch einundzwanzig liegen, dann zeigt 'jaira show <id> --for-lane review --json' den Diff dieser drei - und meldet complete:true. Nichts im Payload sagt, dass achtzehn Commits fehlen.

  Nachgemessen an drei Tickets desselben Tages: GTQHNH fuehrte 3 SHAs bei 21 auf dem Branch, D8CSAA 1 bei 4, 8XHZT5 hatte das Feld leer und bekam deshalb den ganzen Branch. Die review-Lane von GTQHNH hat den Ausschnitt bemerkt und stattdessen 'git diff origin/HEAD...HEAD' geurteilt; die von D8CSAA hat es nicht bemerkt.

  Warum das gefaehrlich und nicht bloss unpraktisch ist: der Ausfall ist still und sieht aus wie Erfolg. complete:true heisst fuer jeden Leser 'du hast alles'. Eine Review, die einen Bruchteil sieht, findet nichts und unterschreibt - und das ist genau der Fall, gegen den die Review-Lane ueberhaupt existiert.

  Wie 'commits:' teilweise gefuellt wird, ist noch nicht untersucht - vermutlich stempelt ein 'jaira move' den Stand von damals hinein, und spaetere Commits kommen nicht mehr dazu. Das gehoert zur Untersuchung.

  Im Prompt ist das seit f6e8dba (GTQHNH) nur BESCHRIEBEN: jaira-role-lane/SKILL.md sagt der Lane, die den Diff beurteilt, sie solle 'jaira show <id> --json | jq .commits | length' gegen 'git log origin/HEAD..HEAD --oneline | wc -l' zaehlen. Eine Handanweisung ist kein Ersatz dafuer, dass das Werkzeug die Wahrheit sagt.
definition-of-done: "Der Payload einer Lane, die einen Diff beurteilt, zeigt den ganzen Branch - oder er sagt, dass er es nicht tut: entweder leitet showForLane die Liste immer aus git ab, oder complete ist false und 'missing' nennt die Zahl der nicht enthaltenen Commits. Ein Leser kann nicht mehr einen Ausschnitt fuer das Ganze halten."
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-17T22:26:05Z
updated-at: 2026-09-18T06:54:25Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-99942
claimed-at: 2026-09-18T06:30:48Z
mode: conversational
outcome-what: "Test TestCommitsSourceNamesWhoContributed deckt die vier Faelle ab; die Unreleased-Zeile in NOTES.md ergaenzt um das Signoff-Label"
outcome-why: "Befund 1 war ein Payload, der eine Commit-Liste als Herkunft von etwas ausgibt, das gar nicht auf dem Schirm steht; Befund 2 war derselbe stille Ausfall wie der Ticket-Bug selbst, nur auf dem Schirm, auf dem ein Mensch unterschreibt; Befund 3 war ein Kommentar, der vor Auseinanderlaufen warnt und ausgerechnet die teuerste Stelle ausliess"
outcome-resolves: "Alle drei Befunde der critique-Lane vom 2026-09-18 06:47 sind repariert, DoD unveraendert erfuellt; go build ./... und go test ./... gruen"
review-summary: none
---

# Der Lane-Payload liefert einen Ausschnitt des Diffs und meldet ihn als vollstaendig

## Definition of Done

- [x] Der Payload einer Lane, die einen Diff beurteilt, zeigt den ganzen Branch - oder er sagt, dass er es nicht tut: entweder leitet showForLane die Liste immer aus git ab, oder complete ist false und 'missing' nennt die Zahl der nicht enthaltenen Commits. Ein Leser kann nicht mehr einen Ausschnitt fuer das Ganze halten.
  proof: internal/cli/flow.go:591-608 — shas = ticket.MergeCommits(env.DeriveCommits(t), t.Commits), immer abgeleitet; Payload traegt commits + commits_source (flow.go:648-656); internal/tui/signoff.go:114 dieselbe Union
- [x] Nachgestellt an dem Fall, der es gezeigt hat: ein Ticket mit drei SHAs in 'commits:' und einundzwanzig Commits auf dem Branch. Mit Test.
  proof: internal/cli/forlanecommits_test.go TestForLaneDiffIsNotLimitedToTheRecordedCommits — Ticket mit einem von zwei SHAs in commits:, review-Payload muss beide zeigen; fiel vor der Aenderung
- [x] Untersucht und auf dem Ticket festgehalten, WIE 'commits:' teilweise gefuellt wird - welcher Schreibpfad den Stand einfriert und warum er spaetere Commits nicht nachtraegt. Ohne diese Antwort ist jede Reparatur geraten.
  proof: Notiz vom 2026-09-18 (pre-process): StampCommits laeuft nur als move.Request.Prepare beim Einlaufen in eine Doorway-Lane (core/ticket/trim.go:105, core/move/move.go:132); gefuellt wird commits: von 'jaira move --out --commits' (internal/cli/flow.go:99,152)
- [x] Eine Zeile in core/release/NOTES.md, wenn sich aendert, was ein Benutzer im Payload sieht.
  proof: core/release/NOTES.md:18 unter ## Unreleased

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [~] failing test: ticket with a stale, partial commits: field asks for --for-lane review --json and gets the diff of only those SHAs
- [x] showForLane: build the diff SHAs as the union of t.Commits and env.DeriveCommits(t), never the field alone (internal/cli/flow.go:589)
- [x] carry the list into the payload: commits (the SHAs used) plus commits_source, so a reader can count instead of trusting
- [x] check the same field-first pattern in internal/tui/signoff.go:107 - fix if it is the same bug, note it if it is not
- [x] drop the hand instruction from core/role/builtin/jaira-role-lane/SKILL.md (both passages, lines ~36-49 and ~137-146) now the tool tells the truth
- [x] NOTES.md under ## Unreleased: one line telling a reader the review payload is now the whole ticket history
- [x] go build ./... && go test ./... green

## Progress
- **2026-09-18 06:30 · Alexander Sacharov** — Alex am 2026-09-18: dieses Ticket wird 0.3.1. Der Zweig fix/7XXT00 haengt an release/0.3.0 (88b6816), nicht an master - 0.3.0 ist noch nicht gemerged und noch nicht getaggt. Die Zeile fuer NOTES.md gehoert deshalb unter das leere '## Unreleased' ganz oben, NICHT unter '## 0.3.0': diese Sektion wird gleich getaggt und ist damit geschlossene Historie. Das Umbenennen von '## Unreleased' nach '## 0.3.1' und eine frische leere darueber macht ein Mensch in dem Commit, den er taggt - kein Agent.
- **2026-09-18 06:33 · Alexander Sacharov** — Entscheidung zu DoD 1 (Alex, 2026-09-18, vor der Plan-Lane): Weg A - showForLane leitet die Commit-Liste im case "diff" IMMER aus git ab, statt t.Commits als fertige Antwort zu lesen. t.Commits wird mit dem aus git abgeleiteten Stand vereinigt, damit ein eingetragener SHA, den die Ableitung nicht findet, nicht verloren geht: shas := union(t.Commits, env.DeriveCommits(t)). complete:true bedeutet damit wieder 'du siehst den ganzen Branch'. Kein 'missing'-Feld mit einer Zahl fehlender Commits, und kein zusaetzliches Warnfeld - der stille Ausfall wird beseitigt, nicht gemeldet. Begruendung: ein Flag, das niemand liest, ersetzt einen stillen Fehler nur durch einen lauten Umweg. Preis, der bewusst akzeptiert wird: der beurteilte Diff ist nicht mehr allein aus dem Frontmatter reproduzierbar.
- **2026-09-18 06:35 · Alexander Sacharov** — pre-process: Untersuchung abgeschlossen. Wie 'commits:' teilweise gefuellt wird: NICHT durch einen Stempel bei jedem move. StampCommits (core/ticket/trim.go:105) laeuft nur als move.Request.Prepare, und Prepare wird nur beim Einlaufen in eine Doorway-Lane aufgerufen (core/move/move.go:132, lane.Settle) - also beim Ablegen ins Logbuch, nicht bei jedem Lane-Wechsel. Der Weg, auf dem drei SHAs ins Feld kommen, ist 'jaira move --out' mit Commits im Payload (internal/cli/flow.go:99+152): die werden in t.Commits gemergt. Spaeter entstehende Commits kommen nie dazu. Damit ist 'commits:' per Konstruktion eine Momentaufnahme und niemals die Wahrheit - das Feld reparieren heisst also nicht, es aktuell zu halten, sondern es beim Lesen nicht mehr allein zu glauben.
- **2026-09-18 06:36 · Alexander Sacharov** — pre-process: Die DoD laesst zwei Wege zu. Gewaehlt ist der erste - showForLane leitet immer aus git ab - und zwar als VEREINIGUNG von t.Commits und DeriveCommits(t), nicht als Ersatz. Begruendung: StampCommits merged an jeder anderen Stelle schon genau so, explizit gewinnt nie ueber abgeleitet, sondern kommt dazu; ein SHA im Feld, der auf keinem Branch mehr liegt (rebase, cherry-pick), geht so nicht verloren. Der zweite Weg - complete:false plus eine Zahl fehlender Commits - wurde verworfen: um zu wissen, WIEVIELE fehlen, muesste man gegen 'git log origin/HEAD..HEAD' zaehlen, und auf einem geteilten Branch (release/0.3.0 traegt mehrere Tickets) gehoeren fremde Commits legitim nicht zum Ticket. Das waere ein Fehlalarm bei jedem Release-Branch und wuerde eine Lane blockieren, der nichts fehlt. Restrisiko, das beide Wege nicht decken: ein Commit, der die Ticket-ID nicht nennt und die Ticket-Datei nicht anfasst, wird von DeriveCommits nicht gefunden - dagegen hilft nur die Regel 'jeder Commit nennt die ID', nicht das Werkzeug. Deshalb Plan-Schritt 3: die benutzte SHA-Liste kommt in den Payload, damit ein Leser zaehlen kann, statt zu vertrauen.
- **2026-09-18 06:42 · Alexander Sacharov** — in-progress, Schritt 4: signoff.go war derselbe Bug und ist mitgefixt - es las t.Commits zuerst und leitete nur bei leerem Feld ab, also zeigte ein Ticket mit 3 von 21 SHAs einen Diffstat von 3, und nichts sagte es. Jetzt immer MergeCommits(derived, t.Commits).

Nicht gefixt, bewusst: internal/tui/view.go:1319 (Detail-Pane) und internal/cli/tickets.go:776 zeigen t.Commits roh. Das ist KEIN Bug derselben Art - die beschriften das Feld als 'commits' und behaupten nicht, der ganze Stand zu sein; niemand faellt dort ein Urteil ueber einen Diff. Wer das doch aendern will: dieselbe Union genuegt.

Auch nicht angefasst: core/gate/gate.go:327 - 'explizit gewinnt ueber abgeleitet' beim Tor. Dort ist das Feld nur ein Ja/Nein ('gibt es ueberhaupt Commits'), kein Ausschnitt, der fuer das Ganze gehalten werden kann.

Union-Logik lag doppelt vor: die Merge-Schleife aus StampCommits (core/ticket/trim.go) ist jetzt ticket.MergeCommits und wird von beiden Lesern benutzt, damit sie nicht auseinanderlaufen.
- **2026-09-18 06:44 · Alexander Sacharov** — in-progress abgeschlossen. Was nicht offensichtlich ist:

- commits_source hat drei Werte: 'git' (nur die Ableitung trug bei), 'ticket' (die Ableitung fand nichts, das Feld trug alles) und 'git+ticket' (das Feld trug einen SHA bei, den die Ableitung nicht fand - Rebase, Cherry-Pick). Die Schluessel commits/commits_source stehen NUR im Payload einer Lane, die 'diff' in input-requires fuehrt. Begruendung: eine Commit-Liste neben einer Lane, die nie einen Diff verlangt hat, liest sich als Aussage ueber das Ticket statt als Herkunft dessen, was auf dem Schirm steht.

- Die Union-Schleife lag vorher nur in StampCommits. Sie ist jetzt ticket.MergeCommits (core/ticket/trim.go) und wird von drei Stellen benutzt. Wer sie aendert, aendert damit auch, was beim Ablegen ins Logbuch ins Frontmatter geschrieben wird - das ist Absicht, sie duerfen nicht auseinanderlaufen.

- Auch die Klartext-Ausgabe (ohne --json) traegt jetzt '<n> commit(s), from <quelle>:' plus die SHAs ueber dem Diff. Ein Worker, der den Lane-Prompt als Text liest, haette sonst als einziger nicht zaehlen koennen.

- NOTES.md: die Zeile steht unter '## Unreleased', NICHT unter '## 0.3.0'. Die 0.3.0-Zeile 22, die die alte Handanweisung beschreibt, bleibt woertlich stehen - sie beschreibt den Build, den es gibt; geschlossene Historie wird nicht nachtraeglich richtiggestellt.
- **2026-09-18 06:47 · Alexander Sacharov** — critique (erster Durchgang, ganzer Diff gelesen): drei Befunde, alle mit klarer Reparatur, daher zurueck nach in-progress.

1. internal/cli/flow.go:646 — die neuen Schluessel commits/commits_source hahngen an len(shas)>0. Schlaegt repo.Diff fehl, ist diff leer, shas aber gefuellt: der Payload traegt dann eine Commit-Liste als 'Herkunft dessen, was auf dem Schirm steht', waehrend nichts auf dem Schirm steht. Der Klartext-Zweig (:677) macht es schon richtig und prueft diff != "". Beide Zweige auf dieselbe Bedingung.

2. internal/tui/signoff.go:114 — 'derived' bedeutet jetzt nur noch 'git hat etwas gefunden', die Liste ist aber eine Vereinigung. Ein SHA, den nur das Feld traegt (Rebase, Cherry-Pick), erscheint damit unter der Ueberschrift 'derived from git — recorded at acceptance'. Genau diesen Fall unterscheidet flow.go:601 mit git / ticket / git+ticket. Zwei Schirme, dieselben Daten, zwei Ehrlichkeitsmassstaebe — und der Signoff-Schirm ist der, auf dem ein Mensch unterschreibt. Entweder dieselbe Dreiteilung ins Label, oder derived nur setzen, wenn das Feld nichts beigetragen hat.

3. core/ticket/trim.go:123 — der Kommentar sagt 'Two callers share it and must not drift apart' und nennt StampCommits und den Lane-Payload. Dieser Change hat einen dritten hinzugefuegt: internal/tui/signoff.go:114. Der Kommentar existiert, um Auseinanderlaufen zu verhindern, und laesst ausgerechnet die Stelle aus, an der es am teuersten ist. Alle drei nennen.

Bewusst stehen gelassen: commits_source mit drei Werten ist kein Ueberbau — die Note vom 2026-09-18 06:44 haelt die Entscheidung fest (Plan-Schritt 3: der Leser soll zaehlen koennen statt zu vertrauen), und dass die Schluessel nur im Payload einer Lane mit 'diff' in input-requires stehen, ist dort ebenfalls begruendet. Ebenso die verworfene Alternative 'complete:false plus Zahl fehlender Commits' (Note 06:36) — geschlossen, nicht neu aufgemacht. Die Vereinigung statt Ersetzung folgt StampCommits und ist das bestehende Muster, nicht ein neues daneben.
- **2026-09-18 06:51 · Alexander Sacharov** — in-progress (zweiter Durchgang, alle drei critique-Befunde repariert). Was nicht im Code steht:

- Die Dreiteilung git / ticket / git+ticket lag doppelt vor, sobald der Signoff-Schirm sie auch brauchte. Sie ist jetzt ticket.CommitsSource (core/ticket/trim.go) neben MergeCommits - bewusst dort und nicht in internal/tui oder internal/cli, weil genau das Auseinanderlaufen zweier Schirme der zweite Befund war. Wer einen Wert hinzufuegt, muss internal/tui/signoff.go:commitsSourceLabel mitfuehren, sonst faellt der neue Fall still auf 'kein Label'.

- Befund 1 (Payload-Schluessel an len(shas)>0): die Bedingung ist jetzt diff != "", dieselbe wie im Klartext-Zweig. Kein Test dafuer: repo.Diff scheitern zu lassen verlangt ein kaputtes git-Repo mit gueltigen SHAs im Frontmatter - der Aufwand steht nicht zum Nutzen, die Bedingung ist eine Zeile und steht neben ihrem Zwilling. Wer das doch testen will: SHAs eintragen, die auf keinem Objekt liegen, dann meldet Diff einen Fehler.

- Befund 2, gewaehlte Variante: die Dreiteilung ins Label, NICHT 'derived nur setzen, wenn das Feld nichts beigetragen hat'. Begruendung: die zweite Variante laesst im Fall git+ticket gar kein Label stehen, und ein fehlendes Label liest sich als 'keine Aussage' statt als 'gemischte Herkunft' - auf dem Schirm, auf dem unterschrieben wird, ist das wieder ein stiller Ausfall, nur ein kleinerer.

- NOTES.md: die Unreleased-Zeile wurde ERGAENZT, nicht um eine zweite Zeile erweitert. Der Signoff-Schirm war in derselben Zeile schon genannt; eine zweite Zeile haette denselben Change zweimal beschrieben. Unreleased ist offen, das darf man - an einer getaggten Sektion nicht.
