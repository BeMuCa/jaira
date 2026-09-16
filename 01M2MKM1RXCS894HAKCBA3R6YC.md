---
id: 01M2MKM1RXCS894HAKCBA3R6YC
title: "Die pr-Rolle oeffnet den Pull Request selbst, wenn ein Mensch sie aufruft"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer /jaira-role-pr tippt, bekommt einen offenen Pull Request und keine Befehlszeile zum Einfuegen; ruft ein Agent die Rolle auf, bleibt es beim Pushen."
context: |-
  Heute pusht die Rolle jaira-role-pr den Branch und schreibt eine 'gh pr create'-Zeile aus, die der Mensch danach selbst einfuegt und ausfuehrt.

  Alex am 2026-09-16: 'если я стартую /jaira-role-pr то она должна сразу и создавать, просто без команды'. Sein Aufruf IST die Anweisung, den Pull Request zu oeffnen - der Kopierschritt danach kauft nichts und kostet eine Runde.

  Die Regel dahinter bleibt: wer die Aenderung geschrieben hat, entscheidet nicht, dass sie ankommt. Sie wird nur anders aufgehaengt - nicht mehr an 'kein Agent oeffnet je einen PR', sondern an 'nur der Aufruf eines Menschen oeffnet einen'. Ein Dispatcher oder Teamlead, der die Rolle aufruft, pusht weiter nur und reicht die Zeile zurueck. Merge und Approve bleiben in jedem Fall verboten.

  Der Text ist bereits geschrieben und liegt auf dem lokalen Branch chore/pr-role-opens-the-pr (Commit 5b849ca): core/role/builtin/jaira-role-pr/SKILL.md umgeschrieben, eine Zeile in core/release/NOTES.md. Er ist noch nicht durch critique, optimize oder review gelaufen und traegt keine Ticket-Id im Commit.

  Dabei gefunden und offen: in einem Fork zielt 'gh pr create' nicht selbstverstaendlich auf das richtige Repository - origin ist der Fork, die Pull Requests dieses Boards gehen aber alle nach BeMuCa/jaira. Die Rolle muss das pruefen, bevor sie das Kommando ausfuehrt, sonst steht der Pull Request im falschen Repository und ein Mensch muss ihn von Hand schliessen.
definition-of-done: "core/role/builtin/jaira-role-pr/SKILL.md sagt: Aufruf durch einen Menschen -> pushen und oeffnen; Aufruf durch einen Agenten -> pushen und die Zeile zurueckgeben; merge und approve bleiben in beiden Faellen verboten; die Rolle prueft vor dem Oeffnen, in welches Repository der Pull Request geht; eine Zeile unter ## Unreleased in core/release/NOTES.md; go test ./core/role/... gruen"
tags:
  - docs
  - cli
blocked-by: []
related: []
commits:
  - ea78a3abd48ed2c7568c3bb65671a46d262d6d3b
created-at: 2026-09-16T07:59:07Z
updated-at: 2026-09-16T08:09:43Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-11544
claimed-at: 2026-09-16T07:59:59Z
outcome-what: "core/role/builtin/jaira-role-pr/SKILL.md nach critique-Durchgang 1 umgebaut: (1) 'Which repository it goes to' steht jetzt vor dem Listing, das Listing selbst fragt mit '--repo' das Zielrepository und schreibt den Head als '<owner-of-origin>:<branch>'; (2) die Mensch/Agent-Aufteilung steht nur noch unter '## Open it', der Push-Abschnitt sagt sie nicht mehr; (3) jaira.remote wird als Remote-NAME behandelt, mit 'git remote get-url \"$(git config jaira.remote)\"' als Schritt zum owner/repo; (4) die Leiter hat eine vierte Sprosse fuer 'Fork ohne jaira.remote' -> Parent aus 'gh repo view'. Die NOTES.md-Zeile unter ## Unreleased nennt Fallback, Widerspruch und das Listing im Zielrepository mit."
outcome-why: "critique hatte vier Findings: die Reihenfolge war halb gedreht, sodass das Listing im Fork nach einem PR fragte, der im Parent offen ist - genau der zweite PR, den der Text an anderer Stelle ausschliesst; die Mensch/Agent-Regel stand an zwei Stellen mit zwei verschiedenen Aussagen; ':152 read the owner/repo off it' war nicht ausfuehrbar, weil jaira.remote einen Remote-Namen haelt (core/settings/settings.go:145); und der Normalfall 'jaira.remote gar nicht gesetzt' hatte keine Sprosse."
outcome-resolves: "Definition of Done unveraendert erfuellt und jetzt ohne die vier Widersprueche: Mensch -> pushen und oeffnen, Agent -> pushen und Zeile zurueck (an genau einer Stelle gesagt), merge/approve in beiden Faellen verboten, Zielrepository vor Listing UND Oeffnen geprueft mit ausfuehrbaren Schritten, eine Zeile unter ## Unreleased, go test ./core/role/... gruen."
review-summary: |-
  core/role/builtin/jaira-role-pr/SKILL.md:158 'Skip the next section' zeigt seit der Aufteilung auf die falsche Stelle: zwischen Listing und 'Answering review comments' stehen jetzt ZWEI Abschnitte, 'Write the description' und 'Open it'. Wer nur den naechsten ueberspringt, landet auf '## Open it' und oeffnet den zweiten Pull Request, den :154 gerade verboten hat. Beide Abschnitte namentlich nennen: 'Skip **Write the description** and **Open it**'.
  core/role/builtin/jaira-role-pr/SKILL.md:124 Sprosse 2 verschluckt Sprosse 4: 'A fork, and jaira.remote names a remote' trifft auch zu, wenn jaira.remote auf 'origin' zeigt (core/settings/remotefor_test.go:74 setzt genau das) - der Leser haelt auf Sprosse 2 an, nimmt den Fork als Ziel und erreicht Sprosse 4 nie. Sprosse 2 auf den Fall einschraenken, in dem jaira.remote den Parent nennt, oder die Widerspruchspruefung vor Sprosse 2 ziehen - so wie die Forge-Leiter :67-79 sich gegenseitig ausschliesst.
---

# Die pr-Rolle oeffnet den Pull Request selbst, wenn ein Mensch sie aufruft

## Definition of Done

- [x] core/role/builtin/jaira-role-pr/SKILL.md sagt: Aufruf durch einen Menschen -> pushen und oeffnen; Aufruf durch einen Agenten -> pushen und die Zeile zurueckgeben; merge und approve bleiben in beiden Faellen verboten; die Rolle prueft vor dem Oeffnen, in welches Repository der Pull Request geht; eine Zeile unter ## Unreleased in core/release/NOTES.md; go test ./core/role/... gruen
  proof: core/role/builtin/jaira-role-pr/SKILL.md:10-25 (Mensch oeffnet, Agent reicht zurueck), :101-135 (Zielrepository vor Listing und Oeffnen, 4 Faelle), :137-151 (Listing fragt das Zielrepository), :182-202 (Open it: einzige Stelle der Aufteilung), :218-221 (merge/approve verboten); core/release/NOTES.md:17; go test ./core/role/... ok

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Den Zielrepository-Check in core/role/builtin/jaira-role-pr/SKILL.md vor den create-Befehl ziehen und konkret machen: wie man das Ziel ermittelt (gh repo view / git config jaira.remote), was bei einem Fork gilt, und dass ohne klares Ziel nicht geoeffnet wird
- [x] NOTES.md-Zeile unter ## Unreleased so ergaenzen, dass der Zielrepository-Check darin steht
- [x] go test ./core/role/... laufen lassen, Ticket verschieben und zusammen mit dem Code committen

## Progress
- **2026-09-16 08:02 · Alexander Sacharov** — Der Text aus chore/pr-role-opens-the-pr lag beim Start dieser Lane schon als a553e45 auf feat/A3R6YC (mit Ticket-Id im Subject, auf neuerer Basis als 5b849ca). Diese Lane hat nur noch den Zielrepository-Check gebaut, den der Ticket-Kontext als offen markiert hatte.

Warum der Check eine eigene Ueberschrift bekam und nicht ein Absatz beim create-Befehl blieb: er stand vorher NACH den create-Blocks und sagte nur 'wenn die Basis nicht offensichtlich ist, frag nach'. Ein Modell, das den Block von oben liest, hat das Kommando dann schon getippt. Jetzt endet 'Write the description' mit 'du oeffnest nichts, bevor der naechste Abschnitt geklaert hat, wohin' - die Reihenfolge im Dokument IST die Durchsetzung, es gibt keinen Code, der das erzwingt.

Geprueft und verworfen: 'git rev-parse --abbrev-ref origin/HEAD' als Zielquelle - das nennt nur den Default-Branch des Forks, nicht das Zielrepository. gh/glab plus jaira.remote sind die einzigen zwei Quellen, die das Upstream ueberhaupt benennen; darum die Drei-Faelle-Leiter mit 'widersprechen sich -> nicht oeffnen, fragen'.

go test ./core/role/... prueft nur Installation und Parsing der Rolle, nicht den Inhalt von SKILL.md. Gruen heisst hier also 'nichts kaputt', nicht 'Regel getestet'.
- **2026-09-16 08:04 · Alexander Sacharov** — critique (Durchgang 1, 4 Findings): (1) Die Reihenfolge im Dokument ist noch nicht fertig gedreht. Der Zielrepository-Abschnitt kam vor die create-Befehle, aber 'gh pr list --head' bei :99 laeuft immer noch davor und fragt im Fork den Fork - ein offener Pull Request im Parent wird nicht gefunden, und die Rolle oeffnet einen zweiten, obwohl :108 sagt, sie tue das nie. Der Abschnitt muss vor das Listing, und das Listing braucht das Repository als Flag. (2) :93-94 'if an agent did, the push is where you stop' widerspricht :175-177, wo der Agent noch Beschreibung schreibt und die Zeile zurueckreicht. Die Aufteilung Mensch/Agent gehoert an genau eine Stelle - den create-Befehl - sonst liest jeder Durchlauf zwei verschiedene Regeln. (3) jaira.remote haelt einen Remote-NAMEN, nicht owner/repo (core/settings/settings.go:149) - ':152 read the owner/repo off it' ist so nicht ausfuehrbar. (4) Fall 3 der Leiter hat keine Sprosse fuer 'jaira.remote ist gar nicht gesetzt' - genau der Normalfall auf einem Board, dessen Refs nie geteilt wurden. Nicht geprueft, weil Sache der review-Lane: ob die Flags '--repo' / '--project' bei gh/glab genau so heissen.
- **2026-09-16 08:06 · Alexander Sacharov** — in-progress (Durchgang 2, nach critique) — alle 4 Findings abgearbeitet.

(1) Reihenfolge: 'Which repository it goes to' steht jetzt VOR dem Listing, nicht nur vor den create-Blocks. Das Listing trägt '--repo <owner/repo>'. Dabei gefunden: 'gh pr list --head <branch>' findet einen Cross-Repo-PR im Parent nur, wenn der Head als '<owner-of-origin>:<branch>' geschrieben wird — sonst filtert GitHub auf einen gleichnamigen Branch IM Parent und liefert leer, also genau die Nicht-Findung, die den zweiten PR öffnet. Steht jetzt im Text.

Verworfen: das Owner-Präfix per sed aus 'git remote get-url origin' ableiten. Die Regex deckt SSH, HTTPS und '.git'-Suffix nur halb ab und eine falsche Ableitung liefert wieder leer — also keinen Fehler, sondern denselben doppelten PR. Der Text nennt stattdessen '<owner-of-origin>' als Platzhalter und lässt das Modell es aus der URL lesen, die es zwei Abschnitte vorher ohnehin schon angesehen hat.

(2) Mensch/Agent-Aufteilung: aus dem Push-Abschnitt entfernt, lebt jetzt nur noch unter '## Open it'. Der Kopf des Dokuments (:10-25) sagt die Regel weiterhin als Rahmen — das ist Absicht, nicht die zweite Stelle: er sagt WAS gilt, '## Open it' sagt WO es ausgeführt wird.

(3) jaira.remote: hält einen Remote-NAMEN (core/settings/settings.go:145 RemoteFor). Der Text sagt das jetzt ausdrücklich und gibt 'git remote get-url "$(git config jaira.remote)"' als den Schritt, der daraus ein owner/repo macht.

(4) Leiter hat jetzt 4 Sprossen: kein Fork / Fork mit jaira.remote / Fork OHNE jaira.remote -> Parent aus 'gh repo view' / Widerspruch -> nicht öffnen.

Weiterhin ungeprüft und Sache der review-Lane: ob 'glab mr list --repo' und 'glab mr create --target-project' bei der installierten glab-Version genau so heißen. go test ./core/role/... prüft nur Installation und Parsing, nicht den Inhalt von SKILL.md.
- **2026-09-16 08:09 · Alexander Sacharov** — critique (Durchgang 2, 2 Findings). Die 4 Findings aus Durchgang 1 sind abgearbeitet und werden nicht wieder aufgemacht; beide neuen Findings sind Folgen der Umbauten aus Durchgang 2 selbst.

(1) :158 'Skip the next section' - vor der Aufteilung hiess der Abschnitt 'Write the description, then open it' und war EINER, da stimmte der Verweis. Jetzt sind es zwei, und der Zweig 'One listed' ueberspringt nur den ersten. Das Ergebnis ist exakt der doppelte Pull Request, gegen den vier Zeilen hoeher (:154 'you never open a second one') und der ganze Listing-Abschnitt gebaut sind. Fix: beide Abschnitte beim Namen nennen statt 'the next section' - ein relativer Verweis in einem Dokument, dessen Abschnitte sich verschieben, bricht beim naechsten Umbau wieder.

(2) :122-132 Die Leiter ist nicht disjunkt. Sprosse 2 fragt nur 'ist jaira.remote gesetzt', nicht 'zeigt es auf den Parent'. jaira.remote = origin ist ein gueltiger Zustand (core/settings/remotefor_test.go:74), und auf einem Fork trifft dann Sprosse 2 zu - Ziel wird der Fork, also genau der Fehler, gegen den der Abschnitt existiert. Sprosse 4 faengt es nicht ab, weil eine Leiter auf der ersten passenden Sprosse anhaelt. Die Forge-Leiter :67-79 direkt darueber macht es richtig: ihre Sprossen schliessen sich gegenseitig aus. Fix entweder Sprosse 2 einschraenken ('und der Remote nennt den Parent') oder die Widerspruchspruefung vor Sprosse 2 ziehen.

Bewusst NICHT Finding, weil Sache der review-Lane: ob 'glab mr list --repo' und 'glab mr create --target-project' bei der installierten glab-Version so heissen, und ob 'gh pr list --head owner:branch' sich so verhaelt. Diese Lane fragt nicht, ob es laeuft.

Ebenfalls stehen gelassen: dass Kopf (:10-25) und Boundaries (:222-224) die Mensch/Agent-Regel wiederholen, obwohl :99 und :184 'nur an einer Stelle' behaupten. Durchgang 1 hat das aufgemacht, der Implementierer hat es begruendet (Rahmen vs. Ausfuehrungsort) - das bleibt so.
