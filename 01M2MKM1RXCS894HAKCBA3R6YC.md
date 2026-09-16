---
id: 01M2MKM1RXCS894HAKCBA3R6YC
title: "Die pr-Rolle oeffnet den Pull Request selbst, wenn ein Mensch sie aufruft"
status: optimize
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
updated-at: 2026-09-16T11:16:11Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-35292
claimed-at: 2026-09-16T11:10:23Z
outcome-what: "core/role/builtin/jaira-role-pr/SKILL.md: der glab-Aufruf der Zielrepository-Leiter traegt jetzt '-F json' (:110), ein Halbsatz sagt warum (:118-120), Sprosse 1 liest den Projektpfad aus dieser JSON statt aus der Textausgabe (:139-141), und Sprosse 2 trennt die beiden Forges - GitHub '.parent.owner.login' + '.parent.name', GitLab der forked-from-Eintrag derselben JSON (:144-148)."
outcome-why: "critique Durchgang 6: auf GitLab hatte die Leiter keine Quelle. 'glab repo view <url>' laeuft per Default auf -F text und druckt Beschreibung und README - weder den Fork-Status, nach dem die Leiter verzweigt, noch den Elternteil, den Sprosse 2 lesen soll. Der Durchlauf waere auf GitLab still auf Sprosse 1 (den Fork) gefallen: dieselbe stille Verzweigung, die ca7f53c fuer gh geschlossen hat."
outcome-resolves: "Die DoD-Klausel 'die Rolle prueft vor dem Oeffnen, in welches Repository der Pull Request geht' gilt jetzt auf beiden Forges, nicht nur auf GitHub. go test ./core/role/... gruen."
review-summary: "none"
review-gaps: |-
  Vier Durchgaenge, drei Aenderungen, alle in core/role/builtin/jaira-role-pr/SKILL.md, keine davon verhaltensaendernd.

  Kosten: 'jaira whoami --json' und 'git remote get-url <name>' standen im selben Block wie 'gh/glab repo view' und liefen damit bei JEDEM Aufruf der Rolle. Sprosse 1 (kein Fork) braucht beide nicht - sie beantwortet sich allein aus isFork/nameWithOwner. Die beiden Lookups stehen jetzt hinter Sprosse 1, eingeleitet mit 'Ein Fork hat ein zweites Repository'; der haeufige Nicht-Fork-Pfad spart zwei Kommandos.

  Fluff, gestrichen: 'und du brauchst die Antwort vor dem naechsten Kommando, nicht vor dem letzten' (Kommentar ueber die Reihenfolge im Dokument, nicht ueber die Arbeit); 'Das ist hier, wo die zwei Jobs sich trennen, und der einzige Ort, wo sie das tun' vor dem create-Befehl (der Halbsatz 'nur wenn ein Mensch die Rolle aufgerufen hat' zwei Zeilen weiter sagt dasselbe); der zweite Satz von 'Everything below takes that repository as <owner/repo>' war eine Wiederholung des ersten und ist in ihn hineingezogen.

  Duplikation: keine. grep ueber core/role/builtin und core/ nach 'repo view', 'nameWithOwner', 'target-project', '--repo ' findet die Zielrepository-Leiter nur in dieser einen Datei; keine andere Rolle und kein Go-Code ermittelt ein Zielrepository.

  Stehen gelassen und warum: die lange Begruendung zu 'jaira whoami statt git config jaira.remote' (:133-139) - sie sieht wie Fluff aus, ist aber der einzige Ort, der die Falle benennt, dass ein leerer Config-Wert als 'nichts widerspricht' gelesen wird; die Wiederholung der Mensch/Agent-Regel im Abschnitt Boundaries - Boundaries ist in dieser Datei durchgehend eine Wiederholungsliste, das ist Struktur und nicht diese Aenderung; die Zeile in NOTES.md - sie beschreibt Verhalten, das sich nicht geaendert hat.

  go test ./core/role/... gruen nach der letzten Aenderung. Der Test prueft nur Installation und Parsing, nicht den Text - gruen heisst 'nichts kaputt'.
test-verdict: "fail: das blanke 'gh repo view' in der Zielrepository-Leiter beschreibt in einem Fork-Clone nicht origin, sondern das Upstream — auf genau diesem Board meldet es isFork:false, Sprosse 1 feuert, und der Fork-Fall den das Ticket adressiert wird nie erreicht"
---

# Die pr-Rolle oeffnet den Pull Request selbst, wenn ein Mensch sie aufruft

## Definition of Done

- [x] core/role/builtin/jaira-role-pr/SKILL.md sagt: Aufruf durch einen Menschen -> pushen und oeffnen; Aufruf durch einen Agenten -> pushen und die Zeile zurueckgeben; merge und approve bleiben in beiden Faellen verboten; die Rolle prueft vor dem Oeffnen, in welches Repository der Pull Request geht; eine Zeile unter ## Unreleased in core/release/NOTES.md; go test ./core/role/... gruen
  proof: core/role/builtin/jaira-role-pr/SKILL.md:110,118-120 (glab repo view -F json, Begruendung), :139-141 (Sprosse 1 liest den Pfad aus der JSON), :144-148 (Sprosse 2 je Forge), :204-224 ('Open it'), :242-247 (merge/approve verboten); core/release/NOTES.md:17; go test ./core/role/... ok

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Den Zielrepository-Check in core/role/builtin/jaira-role-pr/SKILL.md vor den create-Befehl ziehen und konkret machen: wie man das Ziel ermittelt (gh repo view / git config jaira.remote), was bei einem Fork gilt, und dass ohne klares Ziel nicht geoeffnet wird
- [x] NOTES.md-Zeile unter ## Unreleased so ergaenzen, dass der Zielrepository-Check darin steht
- [x] go test ./core/role/... laufen lassen, Ticket verschieben und zusammen mit dem Code committen
- [x] Die Zielrepository-Leiter auf 'jaira whoami --json' umstellen statt 'git config jaira.remote' (critique Durchgang 3)

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
- **2026-09-16 10:51 · Alexander Sacharov** — in-progress (Durchgang 4, nach critique 3) — die eine Feststellung aus critique 3 abgearbeitet: die Zielrepository-Leiter liest den Board-Remote jetzt aus 'jaira whoami --json' (.remote / .remote_source) statt aus 'git config jaira.remote'.

Warum whoami und keine eigene Herleitung: 'git config jaira.remote' ist Sprosse 1 von vier in core/settings/settings.go RemoteSourceFor; danach kommen settings.json, das einzige Remote und der laute Fehlschlag. Auf diesem Board ist der Key leer und der Remote trotzdem 'upstream' (verifiziert: 'jaira whoami --json' liefert remote=upstream, remote_source='from settings.json on this machine'). Leer gelesen haette die Rolle 'nichts widerspricht' gelesen und den Check uebersprungen. internal/cli/whoami.go:136-142 sagt genau das schon als Kommentar: wer die Leiter ein zweites Mal ableitet, nennt ein anderes Remote als der Code, der spaeter faellt.

Verworfen: 'git remote get-url "$(jaira whoami --json | jq -r .remote)"'. jq steht in keiner anderen Rolle und ist keine Voraussetzung dieses Projekts; das Modell liest die JSON ohnehin selbst. Im Text steht jetzt 'git remote get-url <the .remote name>'.

Neu in Sprosse 2: der Fall, dass der Name gar kein Remote hier benennt (RemoteSourceFor gibt den konfigurierten Namen auch dann zurueck, wenn es ihn nicht gibt — Sprosse 4 der Go-Leiter). Ohne den Fall haette 'git remote get-url' gefehlt und keine Sprosse getroffen.

Mitgezogen: :63 sagte 'Not off jaira.remote' und haette den Leser genau auf den Key zurueckgeschickt, den der Abschnitt darunter verbietet — heisst jetzt 'the board's remote'. NOTES.md-Zeile nennt whoami und warum nicht den Key.

Weiterhin ungeprueft und Sache der review-Lane: die Flagnamen '--repo' / '--target-project' bei der installierten gh/glab-Version und das Verhalten von 'gh pr list --head owner:branch'.
- **2026-09-16 10:54 · Alexander Sacharov** — critique Durchgang 4: ein Finding, und es ist keine Wiederholung — die Umstellung auf 'jaira whoami' aus Durchgang 3 ist sauber (Felder .remote und .remote_source existieren, internal/cli/whoami.go:66-67; RemoteSourceFor existiert, core/settings/settings.go:178).

Finding: :149-150 'In case 1 leave the flag off'. Das Repository-Flag ist damit sprossenabhaengig, und diese Abhaengigkeit muss der Durchlauf ueber vier Befehle in vier Abschnitten mitnehmen: :160 (gh pr list), :166 (glab mr list), :204 (gh pr create), :210 (glab mr create). Gekauft wird damit nichts ausser einem nicht getippten Flag. Fix: Sprosse 1 (:134-135) sagt statt 'braucht kein Flag', dass das owner/repo dort aus 'nameWithOwner' kommt — die Zeile 'gh repo view --json isFork,parent,nameWithOwner' (:110) holt das Feld ohnehin schon und niemand benutzt es bisher. Dann ist <owner/repo> auf jeder Sprosse gesetzt, :149-150 faellt ersatzlos weg und die vier Befehle stehen unbedingt da.

Warum das in dieser Lane ein Finding ist und nicht Kosmetik: dieses Dokument hat genau einen Fehlermodus — ein Modell liest einen Zweig falsch. Drei der vier bisherigen Durchgaenge haben genau das gefunden (Reihenfolge, 'skip the next section', nicht-disjunkte Leiter). Ein Zweig, der nur Tipparbeit spart, ist in so einem Dokument teurer als das Flag.

Geprueft und NICHT beanstandet: die Reihenfolge Push -> Zielrepository -> Listing -> Beschreibung -> Oeffnen; die Disjunktheit der drei Sprossen; der head-Praefix <owner-of-origin>: bei gh pr list; der Mensch/Agent-Schnitt an genau einer Ausfuehrungsstelle (:199-217) mit Kopf und Boundaries als Rahmen — das hat Durchgang 1 aufgemacht und der Implementierer begruendet, es bleibt.

Erwogen und verworfen: :62-65 sagt 'the board's remote' zwei Abschnitte bevor :101-122 erklaert, was das ist — ein Vorwaertsverweis, aber der Satz dort ist eine Abgrenzung ('nicht davon ablesen') und braucht die Definition nicht. Ebenso verworfen: die Forge wird aus origin ermittelt (:62), das Ziel kann der Parent sein — Forks liegen immer auf derselben Forge, und ein Board-Remote auf einer anderen Forge faellt auf Sprosse 3 und fragt.

Sache der review-Lane, hier bewusst nicht gepruet: ob 'gh pr create --repo' und 'glab mr create --target-project' auf ein Nicht-Fork-Repository angewandt harmlos sind — mein Finding setzt das voraus. Ebenso die Flagnamen selbst und das Verhalten von 'gh pr list --head owner:branch'.
- **2026-09-16 10:55 · Alexander Sacharov** — in-progress (Durchgang 5, nach critique 4) — das eine Finding abgearbeitet: das Repository-Flag ist nicht mehr sprossenabhängig.

Sprosse 1 sagt jetzt, woher das owner/repo kommt (nameWithOwner aus dem 'gh repo view', das drei Zeilen darüber ohnehin schon läuft, bzw. der Pfad aus 'glab repo view') statt 'braucht kein Flag'. Damit hält jede der drei Sprossen ein <owner/repo>, und der Satz danach sagt das ausdrücklich: 'there is no branch where the flag is left off'. Die vier Befehle (gh pr list, glab mr list, gh pr create, glab mr create) stehen jetzt unbedingt mit --repo / --target-project da; kein Durchlauf muss mehr über vier Abschnitte hinweg mitnehmen, auf welcher Sprosse er stand.

Mitgezogen gegenüber dem Finding-Wortlaut: critique nannte nur nameWithOwner (gh). Sprosse 1 gilt aber für beide Forges, also nennt sie auch den Pfad aus 'glab repo view' — sonst hätte GitLab auf Sprosse 1 keine Quelle und der Zweig wäre genau so still kaputt wie der, den das Finding schließt.

NOTES.md nicht angefasst: die Zeile beschreibt den Zielrepository-Check als Ganzes, und von außen ändert sich an ihm nichts — es ist derselbe Check, nur ohne Sonderfall im Text.

Weiterhin ungeprüft und Sache der review-Lane: ob '--repo' / '--target-project' auf ein Nicht-Fork-Repository angewandt harmlos sind (das setzt dieser Umbau voraus), die Flagnamen bei der installierten gh/glab-Version, und das Verhalten von 'gh pr list --head owner:branch'.
- **2026-09-16 10:57 · Alexander Sacharov** — critique Durchgang 5: keine Findings. Nachgeprueft statt geglaubt: 'jaira whoami --json' gibt hier wirklich .remote (upstream) und .remote_source ('from settings.json on this machine') aus, und RemoteSourceFor in core/settings/settings.go:178 hat tatsaechlich die vier Stufen, die SKILL.md:116-120 behauptet - die Begruendung, warum 'git config jaira.remote' zu wenig ist, stimmt also woertlich. Die Leiter bei SKILL.md:134-148 ist disjunkt und total: Sprosse 2 faengt auch den Fall, dass der Remote-Name hier gar keine URL hat (RemoteSourceFor Stufe 4 gibt einen Namen zurueck, den es nicht gibt). Auf Widerspruch mit dem Rest geprueft: jaira-teamlead SKILL.md:90-102 und jaira-dispatcher SKILL.md:191 sagen weiter 'ein Agent oeffnet nie' - das bleibt unter der neuen Regel wahr, weil beide die Rolle als Agent aufrufen. go test ./core/role/... ok. Nicht neu aufgemacht, weil in Durchgang 3/4 schon entschieden: dass 'jaira whoami' nur fuer den Widerspruchsfall (Sprosse 3) ueberhaupt gebraucht wird - Sprosse 1 und 2 landen ohnehin beim Parent.
- **2026-09-16 11:00 · Alexander Sacharov** — optimize: vier Durchgaenge ueber den Diff, nur Prosa gekuerzt, keine Regel entfernt. Duplikat-Suche: 'gh pr create --repo', '--target-project' und 'nameWithOwner' kommen im ganzen Repository nur in core/role/builtin/jaira-role-pr/SKILL.md vor - es gibt keine zweite Zielrepository-Leiter, die mitgepflegt werden muesste. Der whoami-Absatz behauptet Fakten, die stimmen: core/settings/settings.go RemoteSourceFor hat wirklich vier Stufen (Kommentar Zeile 174), und internal/cli/whoami.go:67 setzt wirklich 'remote_source'. Beides geprueft statt geglaubt.

Entfernt wurde nur, was an zwei Stellen stand: der Hinweis im Intro von 'Which repository it goes to', dass eine Abfrage an den Fork den offenen Pull Request upstream nicht sieht - das ist der Punkt des naechsten Abschnitts 'Does it already have one open' und steht dort schon; die Wiederholung der Definition 'die Remote, auf der die Ticket-Refs reisen' in Sprosse 2, die der Forge-Abschnitt zwanzig Zeilen weiter oben schon gibt; und die Begruendung 'so that the run reads one rule rather than two' unter 'Push the branch', die nur erklaert, warum der Satz davor dasteht. Der whoami-Absatz wurde von acht auf sieben Zeilen gezogen, Inhalt unveraendert.

Stehen gelassen und bewusst nicht angefasst: die Mensch/Agent-Regel steht drei Mal (Kopf, 'Open it', Boundaries). Das ist keine Fluff-Wiederholung, sondern die Regel am Ort der Handlung - ein Modell, das erst bei 'Open it' einsteigt, liest den Kopf nicht mehr. Boundaries ist per Bauart eine Wiederholung. 4988 -> 1929 Woerter waere Kuerzen um des Kuerzens willen gewesen; die Leiter selbst bleibt Wort fuer Wort, weil jede Sprosse eine Unterscheidung traegt.
- **2026-09-16 11:04 · Alexander Sacharov** — testing: Gates gruen, aber die Zielrepository-Leiter greift auf einem Fork-Clone nicht.

GATES (gruen, nichts kaputt): go build ./... OK, go vet ./... OK, go test ./... -race -count=1 RC=0 (alle Pakete ok, internal/tui 108.8s). Das neue Binary embedded den Text: 'jaira roles list' des frischen Builds nennt die neue pr-Beschreibung, 'roles install --global --force' schreibt 8 Rollen, die installierte SKILL.md ist byte-identisch zur Quelle.

BEFUND (SKILL.md:107): der Befehl in der Leiter ist 'gh repo view --json isFork,parent,nameWithOwner' ohne Repository-Argument. gh loest das Basis-Repository selbst auf und bevorzugt dabei das Upstream, nicht origin. Auf genau diesem Clone (origin=git@github.com:sashasoft90/jaira.git, upstream=git@github.com:BeMuCa/jaira.git, kein 'gh repo set-default'):

    $ gh repo view --json isFork,parent,nameWithOwner
    {"isFork":false,"nameWithOwner":"BeMuCa/jaira","parent":null}

    $ gh repo view "$(git remote get-url origin)" --json isFork,parent,nameWithOwner
    {"isFork":true,"nameWithOwner":"sashasoft90/jaira","parent":{..."login":"BeMuCa"}}

Folge: Sprosse 1 ('Not a fork') feuert auf einem Fork. Hier faellt das Ergebnis zufaellig richtig aus, weil gh schon aufs Upstream aufgeloest hat — aber Sprosse 2 und Sprosse 3 werden nie erreicht. Der Widerspruchsfall, fuer den das Ticket existiert (Fork, dessen Board-Remote ein drittes Repository nennt -> nicht oeffnen, fragen), kann so nicht ausloesen. Dazu sagt Sprosse 1 woertlich 'the target is origin's own repository: the nameWithOwner the same gh repo view already returned' — das ist auf diesem Clone sachlich falsch, die beiden sind verschiedene Repositories.

FIX (klein, nur Text): das Repository im Befehl benennen, statt gh raten zu lassen —

    gh repo view "$(git remote get-url origin)" --json isFork,parent,nameWithOwner

gh akzeptiert die SSH-URL, oben verifiziert. Fuer glab dasselbe pruefen, 'glab repo view' hat kein Argument in der Leiter.

Nebenbefund am selben Absatz: Sprosse 2 sagt 'the parent from gh repo view'. --json parent liefert kein nameWithOwner, sondern .parent.owner.login und .parent.name — der Leser muss 'BeMuCa/jaira' erst zusammensetzen. Ein Halbsatz dazu erspart die Runde.

DoD-Punkt 1 habe ich auf --doing zurueckgesetzt: die Klausel 'die Rolle prueft vor dem Oeffnen, in welches Repository der Pull Request geht' steht im Text, laeuft aber im Fork nicht. Die uebrigen Klauseln der Zeile (Mensch oeffnet / Agent reicht zurueck / merge+approve verboten / NOTES-Zeile) habe ich im Baum verifiziert und sind erfuellt: SKILL.md:9-25, :214-216, :230-234, NOTES.md:17 (eine Zeile, Format intakt, 4 Eintraege unter ## Unreleased).
- **2026-09-16 11:06 · Alexander Sacharov** — in-progress (Durchgang 6, nach testing) — der Befund aus der testing-Lane abgearbeitet: die Zielrepository-Leiter fragte den falschen Clone.

Der Fehler war nicht die Leiter, sondern ihre Eingabe. 'gh repo view' ohne Argument loest das Basis-Repository selbst auf und bevorzugt dabei das Upstream; auf diesem Fork-Clone (origin=sashasoft90/jaira, upstream=BeMuCa/jaira, kein 'gh repo set-default') antwortet es isFork:false / nameWithOwner=BeMuCa/jaira. Sprosse 1 feuerte damit auf einem Fork und Sprosse 2 und 3 waren unerreichbar — der Widerspruchsfall, fuer den das Ticket existiert, konnte nie ausloesen. Selbst nachgeprueft statt uebernommen: beide Aufrufe hier ausgefuehrt, die Ausgaben unterscheiden sich genau so.

Fix: 'gh repo view "$(git remote get-url origin)"' und dieselbe Form fuer glab. 'glab repo view --help' der hier installierten Version (/home/alex/.local/bin/glab) nennt 'glab repo view git@gitlab.com:user/repo.git' ausdruecklich als Aufrufform — die SSH-URL ist also fuer beide Forges gueltig, nicht nur fuer gh.

Bewusst KEIN '--output json' bei glab dazugeschrieben, obwohl es die Flagliste hat. Ich haette dann Feldnamen fuer Pfad und Fork-Elternteil nennen muessen ('path_with_namespace', 'forked_from_project'), die ich hier gegen keine GitLab-Instanz pruefen kann — ein erfundener Feldname ist in diesem Dokument teurer als der Satz 'die Pfadangabe, die glab repo view ausgegeben hat', der ohne Feldnamen auskommt.

Nebenbefund der testing-Lane mitgenommen: Sprosse 2 sagt jetzt, dass '--json parent' kein fertiges owner/repo liefert, sondern '.parent.owner.login' und '.parent.name', die man selbst mit einem Schraegstrich verbindet.

NOTES.md nicht angefasst: die Zeile beschreibt den Check als Ganzes ('settles the target repository from gh repo view / glab repo view'). Von aussen aendert sich nichts — derselbe Check, nur mit dem Repository benannt statt geraten.
- **2026-09-16 11:09 · Alexander Sacharov** — critique Durchgang 6 (nach testing + in-progress 6): ein Finding, und es ist die GitLab-Haelfte genau des Fehlers, den ca7f53c fuer GitHub geschlossen hat.

FINDING SKILL.md:110 / :139-140 / :144 — auf GitLab hat die Zielrepository-Leiter keine Quelle. 'glab repo view "$(git remote get-url origin)"' laeuft ohne -F, und der Default ist text: 'Display the description and README of a project' (glab repo view --help der hier installierten Version, /home/alex/.local/bin/glab, Flag '-F --output  Format output as: text, json. (text)'). Diese Ausgabe nennt weder den Fork-Status, nach dem die Leiter ueberhaupt verzweigt, noch den Fork-Elternteil, den Sprosse 2 daraus lesen soll. Sprosse 1 sagt woertlich 'or the path glab repo view printed' — text-Output druckt keinen Pfad als Feld. Auf GitLab faellt der Durchlauf damit entweder still auf Sprosse 1 zurueck (also auf origin, den Fork) oder bleibt ohne Antwort stehen. Das ist derselbe stille Zweig, den die testing-Lane fuer gh gefunden hat.

FIX ohne erfundene Feldnamen: '-F json' in den Codeblock bei :110, und Sprosse 1 und 2 sagen 'lies Pfad bzw. Fork-Elternteil aus dieser JSON' statt 'den Pfad, den glab repo view ausgegeben hat'. Damit muss niemand 'path_with_namespace' oder 'forked_from_project' behaupten — der Durchgang-5-Einwand gegen '--output json' war, Feldnamen nicht gegen eine GitLab-Instanz pruefen zu koennen, und der bleibt erfuellt, weil das Modell die JSON selbst liest. Ohne '-F json' liest es dagegen Prosa, in der die Felder gar nicht vorkommen.

GEPRUEFT UND NICHT BEANSTANDET: die gh-Seite von ca7f53c stimmt — 'gh repo view "$(git remote get-url origin)"' ist der richtige Aufruf und '--json parent' liefert wirklich .parent.owner.login / .parent.name, nicht nameWithOwner. Die Reihenfolge Push -> Zielrepository -> Listing -> Beschreibung -> Oeffnen, die Disjunktheit der drei Sprossen, der head-Praefix bei 'gh pr list', der whoami-Absatz und der Mensch/Agent-Schnitt an einer Ausfuehrungsstelle: alles unveraendert in Ordnung, nichts davon wird neu aufgemacht.

SACHE DER REVIEW-LANE, hier bewusst nicht gefragt: ob 'glab mr list --repo' und 'glab mr create --target-project' bei der installierten Version so heissen, und ob 'gh pr list --head owner:branch' sich so verhaelt.

ZUR SCHLEIFE: das ist der sechste critique-Durchgang. Er wiederholt kein frueheres Finding — Durchgang 5 fand nichts, testing hat danach den gh-Aufruf aufgemacht, und dieses Finding ist die GitLab-Seite derselben Korrektur. Wenn der Mensch die Schleife trotzdem beenden will, ist der saubere Schnitt: A3R6YC ohne den GitLab-Zweig weiterziehen und den Zweig als eigenes Ticket fuehren.
- **2026-09-16 11:11 · Alexander Sacharov** — in-progress (Durchgang 7, nach critique 6) — das eine Finding abgearbeitet: die GitLab-Haelfte der Zielrepository-Leiter hatte keine Quelle.

'glab repo view <url>' laeuft per Default auf -F text und druckt Beschreibung und README - weder Fork-Status noch Fork-Elternteil. Selbst nachgeprueft statt uebernommen: 'glab repo view --help' der hier installierten Version (/home/alex/.local/bin/glab) listet '-F --output  Format output as: text, json. (text)' und nennt die Git-URL ausdruecklich als Aufrufform. Der Codeblock bei :110 traegt jetzt '-F json'.

Sprosse 1 sagt 'the project path in the same glab repo view JSON' statt 'the path glab repo view printed'; Sprosse 2 trennt die beiden Forges: GitHub '.parent.owner.login' + '.parent.name' mit Schraegstrich verbunden, GitLab 'read the forked-from project out of the same JSON'.

Bewusst WEITER keine GitLab-Feldnamen ('path_with_namespace', 'forked_from_project') im Text - dieselbe Entscheidung wie in Durchgang 6, und sie traegt jetzt mehr als vorher: mit '-F json' liest das Modell die Felder selbst, vorher haette es Prosa gelesen, in der sie gar nicht vorkommen. Ein erfundener Feldname, gegen keine GitLab-Instanz geprueft, waere in diesem Dokument teurer als die Umschreibung.

NOTES.md nicht angefasst: die Zeile beschreibt den Check als Ganzes ('settles the target repository from gh repo view / glab repo view'). Von aussen aendert sich nichts - derselbe Check, nur auf GitLab mit einer Ausgabe, in der die Felder stehen.
- **2026-09-16 11:13 · Alexander Sacharov** — critique Durchgang 7 (nach in-progress 7): keine Findings. Der Diff 1e4442f beantwortet das Finding aus Durchgang 6 vollstaendig und macht nichts Neues auf.

NACHGEPRUEFT STATT GEGLAUBT: 'glab repo view --help' der hier installierten Version (/home/alex/.local/bin/glab) sagt woertlich 'Display the description and README of a project' und listet '-F --output  Format output as: text, json. (text)' — die Begruendung bei SKILL.md:118-120 stimmt also Wort fuer Wort, und '-F json' bei :110 ist die Flagschreibweise dieser Version. Die Git-URL als Argument steht ebenfalls ausdruecklich in den EXAMPLES, fuer glab wie fuer gh.

ERWOGEN UND VERWORFEN (haette ein Finding werden koennen): die GitLab-Haelfte von 'Does it already have one open' (:176) traegt kein Gegenstueck zum GitHub-Praefix '<owner-of-origin>:<branch>' bei :170 — also scheinbar dieselbe Asymmetrie, die Durchgang 6 in der Leiter gefunden hat. Ist aber keine: 'glab mr list --help' kennt nur '-s --source-branch <name>', und auf GitLab liegt ein Merge Request aus einem Fork im ZIEL-Projekt mit dem Branchnamen des Forks als source_branch. Die Abfrage am Ziel findet ihn also ohne Praefix. Das Praefix ist eine GitHub-Eigenheit, kein fehlender Zweig. Deshalb kein Finding — und deshalb steht es hier, damit der naechste Durchgang es nicht noch einmal aufmacht.

NICHT NEU AUFGEMACHT, weil in frueheren Durchgaengen entschieden: der Feldname fuer den Fork-Elternteil in der glab-JSON bleibt bewusst ungenannt (Durchgang 5/6: nicht gegen eine GitLab-Instanz pruefbar, das Modell liest die JSON selbst); die Mensch/Agent-Regel steht an drei Stellen als Rahmen + Ausfuehrungsort (Durchgang 1); die Reihenfolge Push -> Zielrepository -> Listing -> Beschreibung -> Oeffnen, die Disjunktheit der drei Sprossen und der whoami-Absatz sind seit Durchgang 5 unbeanstandet.

SACHE DER REVIEW-LANE, hier bewusst nicht gefragt: ob 'glab repo view -F json' den Fork-Elternteil ueberhaupt enthaelt, und ob 'glab mr create --target-project' so heisst. Diese Lane fragt nicht, ob es laeuft.
