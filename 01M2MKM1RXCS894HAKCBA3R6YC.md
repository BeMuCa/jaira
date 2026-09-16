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
updated-at: 2026-09-16T10:49:16Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-89868
claimed-at: 2026-09-16T10:49:16Z
outcome-what: "core/role/builtin/jaira-role-pr/SKILL.md nach critique-Durchgang 2 an zwei Stellen korrigiert: (1) der Zweig 'One listed' im Listing-Abschnitt nennt jetzt beide zu ueberspringenden Abschnitte namentlich ('Skip **Write the description** and **Open it**') statt 'the next section'; (2) die Zielrepository-Leiter ist von 4 auf 3 Sprossen umgebaut und disjunkt - Sprosse 2 fragt 'nennt jaira.remote ein drittes Repository' statt 'ist jaira.remote gesetzt', sodass unset, =origin und =parent auf dieselbe Antwort (den Parent) fallen und nur ein echter Widerspruch auf Sprosse 3 landet. Davor steht jetzt die Anweisung, Parent und jaira.remote-URL beide aufzuloesen, bevor die Leiter gelesen wird. Die NOTES.md-Zeile nennt den Fallback entsprechend weiter gefasst."
outcome-why: "critique-Durchgang 2 hatte zwei Findings, beide Folgen der Umbauten aus Durchgang 2: der relative Verweis 'the next section' zeigte nach der Aufteilung in zwei Abschnitte auf '## Open it' und haette genau den zweiten Pull Request geoeffnet, den vier Zeilen hoeher verboten ist; und die Leiter hielt auf Sprosse 2 an, sobald jaira.remote ueberhaupt gesetzt war - bei jaira.remote=origin (ein gueltiger Zustand, core/settings/remotefor_test.go:74) also im Fork, also genau der Fehler, gegen den der Abschnitt existiert."
outcome-resolves: "Definition of Done unveraendert erfuellt: Mensch -> pushen und oeffnen, Agent -> pushen und Zeile zurueck, merge/approve verboten, Zielrepository vor Listing und Oeffnen geprueft - jetzt mit einer Leiter, die jeden Zustand genau einmal trifft, und einem Sprungziel, das beim naechsten Umbau nicht still falsch wird. Eine Zeile unter ## Unreleased, go test ./core/role/... gruen."
review-summary: "core/role/builtin/jaira-role-pr/SKILL.md:112,119 liest das Zielrepository-Gegenstueck mit 'git config jaira.remote' selbst aus - das ist nur Stufe 1 der vierstufigen Leiter in core/settings/settings.go:145-200 (RemoteFor). Auf genau diesem Board ist jaira.remote ungesetzt und das Board-Remote trotzdem 'upstream' (aus settings.json): 'jaira whoami --json' antwortet {\"remote\":\"upstream\",\"remote_source\":\"from settings.json on this machine\"}, die Rolle sieht leer und faellt auf Sprosse 2 durch, also ohne jeden Gegencheck. internal/cli/whoami.go:136-142 schreibt genau diese Regel schon auf: wer das Remote neu herleitet, nennt womoeglich ein anderes als der Code, der scheitert. Stattdessen: 'jaira whoami --json' lesen, '.remote' nehmen, mit 'git remote get-url <name>' in owner/repo aufloesen, und den Leiter-Text von 'jaira.remote ist ungesetzt' auf 'jaira whoami nennt kein drittes Repository' umstellen."
---

# Die pr-Rolle oeffnet den Pull Request selbst, wenn ein Mensch sie aufruft

## Definition of Done

- [x] core/role/builtin/jaira-role-pr/SKILL.md sagt: Aufruf durch einen Menschen -> pushen und oeffnen; Aufruf durch einen Agenten -> pushen und die Zeile zurueckgeben; merge und approve bleiben in beiden Faellen verboten; die Rolle prueft vor dem Oeffnen, in welches Repository der Pull Request geht; eine Zeile unter ## Unreleased in core/release/NOTES.md; go test ./core/role/... gruen
  proof: core/role/builtin/jaira-role-pr/SKILL.md:122-136 (disjunkte Leiter), :162-163 (beide Abschnitte namentlich), :182-202 (Mensch/Agent an einer Stelle), :218-221 (merge/approve verboten); core/release/NOTES.md:17; go test ./core/role/... ok

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
- **2026-09-16 08:11 · Alexander Sacharov** — in-progress (Durchgang 3, nach critique 2) — beide Findings abgearbeitet.

(1) :158 nennt jetzt beide Abschnitte beim Namen ('Skip **Write the description** and **Open it**') statt 'the next section'. Bewusst namentlich und nicht relativ: die Abschnitte dieses Dokuments haben sich in zwei Durchgaengen zweimal verschoben, ein relativer Verweis bricht beim naechsten Umbau erneut und bricht dabei still.

(2) Die Zielrepository-Leiter hat jetzt 3 statt 4 Sprossen und ist disjunkt. Der Fehler war nicht die fehlende Sprosse, sondern die Frage: Sprosse 2 fragte 'ist jaira.remote gesetzt', und 'gesetzt' ist kein Kriterium fuer 'zeigt aufs Upstream' — jaira.remote=origin ist ein gueltiger Zustand (core/settings/remotefor_test.go:74). Sprosse 2 fragt jetzt 'nennt jaira.remote ein DRITTES Repository', und die Faelle unset / =origin / =parent fallen alle auf dieselbe Antwort (Parent), weil sie dieselbe Antwort haben. Dadurch verschwindet die alte Sprosse 3 als eigener Fall.

Verworfen: Sprosse 2 nur einschraenken auf 'und der Remote nennt den Parent'. Dann haette 'jaira.remote=origin' gar keine Sprosse mehr getroffen und die Leiter waere am Ende ohne Antwort dagestanden — ein Loch statt einer Ueberlappung.

Neu im Text davor: 'Resolve both the parent and that URL before you read the ladder'. Ohne das liest ein Modell die Leiter von oben und beantwortet Sprosse 2, bevor es den Parent ueberhaupt abgefragt hat.

NOTES.md-Zeile mitgezogen: sagte 'falls back to the fork's parent when jaira.remote is unset' — das war nach der Umstellung zu eng.

Weiterhin ungeprueft und Sache der review-Lane: die Flagnamen '--repo' / '--target-project' bei der installierten gh/glab-Version und das Verhalten von 'gh pr list --head owner:branch'.
- **2026-09-16 08:13 · Alexander Sacharov** — critique Durchgang 3: eine Feststellung, und sie betrifft die Quelle, nicht die Form der Leiter.

Die Leiter selbst ist nach Durchgang 2 disjunkt - daran ist nichts mehr. Falsch ist, WORAUS Sprosse 2 und 3 ihre Antwort ziehen: 'git config jaira.remote' ist nur Stufe 1 der Leiter, die jaira selbst faehrt (core/settings/settings.go:145-200, RemoteFor - danach kommen settings.json, das einzige Remote, und der laute Fehlschlag). Auf diesem Board ist jaira.remote ungesetzt und das Board-Remote ist trotzdem 'upstream', aus settings.json. Die Rolle liest leer, liest das als 'nichts widerspricht' und oeffnet ohne Gegencheck - hier faellt es nicht auf, weil parent und upstream dasselbe Repository sind, aber der Check ist dann schlicht nicht gelaufen.

internal/cli/whoami.go:136-142 hat genau dafuer schon einen Kommentar: 'a command that re-derived it could name a different one than the code that fails'. Die Rolle ist der zweite Ableiter, vor dem der Kommentar warnt.

Fix: 'jaira whoami --json' statt 'git config jaira.remote'; Feld '.remote' ist der Remote-NAME, 'git remote get-url <name>' macht owner/repo daraus. Der Leiter-Wortlaut muss mit: Sprosse 2 fragt dann nicht mehr 'jaira.remote ist ungesetzt oder nennt kein drittes Repository', sondern 'jaira whoami nennt kein drittes Repository'. 'remote_source' aus derselben Ausgabe ist ausserdem der Satz, den man beim Nachfragen auf Sprosse 3 zitieren kann - er sagt, woher der Name kam.

Geprueft und NICHT beanstandet: die Reihenfolge Push -> Zielrepository -> Listing -> Beschreibung -> Oeffnen, das Sprungziel in 'One listed', der head-Praefix <owner-of-origin>: bei gh pr list, und der Mensch/Agent-Schnitt, der nur an einer Stelle steht. jaira.forge ist eine Einstellung, die kein Go-Code liest - die stand aber schon vor diesem Ticket da und gehoert nicht in diesen Diff.
- **2026-09-16 08:14 · Alexander Sacharov** — Dispatcher: critique hat dreimal zurueckgeschickt, damit ist die Schleifengrenze erreicht und die Entscheidung liegt beim Menschen. Die drei Runden waren inhaltlich verschieden und keine wiederholte sich: (1) Zielrepository-Check stand nach den create-Blocks und jaira.remote wurde als owner/repo statt als Remote-Name gelesen; (2) das Listing fragte den Fork statt das Zielrepository, und 'skip the next section' zeigte nach der Umstellung auf 'Open it'; (3) 'git config jaira.remote' ist nur Sprosse 1 der vierstufigen Leiter in core/settings/settings.go:145-200 - auf diesem Board ist jaira.remote leer und das Board-Remote trotzdem 'upstream' aus settings.json, also faellt die Rolle ohne Gegencheck durch. Der vorgeschlagene Ausweg ist 'jaira whoami --json' statt eigener Herleitung; internal/cli/whoami.go:136-142 schreibt genau das schon vor. Offene Frage an den Menschen: eine vierte Runde laufen lassen, oder den Punkt als eigenes Ticket fuehren und A3R6YC ohne ihn weiterziehen.
