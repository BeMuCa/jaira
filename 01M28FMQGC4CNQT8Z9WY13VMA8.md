---
id: 01M28FMQGC4CNQT8Z9WY13VMA8
title: Jede Aenderung faehrt auf einem Branch und kommt durch einen PR
status: in-progress
ready: true
creator: Alexander Sacharov
goal: "Es steht als Regel des Projekts geschrieben, dass Arbeit auf einem Branch mit ihrem Ticket faehrt und ueber einen PR ankommt - und dass das Pruefen dieses PRs dem Maintainer gehoert, nicht dem, der ihn aufmacht"
context: |-
  Heute gilt die Regel schon, aber nur als Gewohnheit: die Arbeit an den Ref-Tickets lief auf einem Branch, das Ticket ritt in denselben Commits mit, und der PR (#9) ging ans Original. Nirgends steht das aufgeschrieben.

  Was fehlt, konkret: der generierte jaira-Block sagt 'das Ticket faehrt im selben Commit wie der Code', aber nicht, dass dieser Commit auf einem Branch liegt und master nur durch einen PR erreicht. Ein Agent, der die Anweisungen liest, darf daraus schliessen, direkt auf master zu committen - was heute niemand tut, aber nichts verbietet.

  Und die zweite Haelfte, die genauso fehlt: wer den PR aufmacht, prueft ihn nicht selbst ab. Das Abnehmen gehoert dem Maintainer. Das Projekt sagt dasselbe schon einmal an anderer Stelle - 'a review agent cannot certify its own work' - nur eben nicht ueber PRs.

  Der Platz dafuer ist der Bereich hinter dem jaira:local-Marker in CLAUDE.md und AGENTS.md: alles dahinter ueberlebt die naechste Regeneration des Blocks. Dort steht schon die NOTES.md-Regel, also gehoert diese daneben.
definition-of-done: "hinter dem jaira:local-Marker in CLAUDE.md und AGENTS.md steht die Regel: Arbeit laeuft auf einem Branch, das Ticket faehrt in denselben Commits mit, master wird nur durch einen PR erreicht, und der PR gehoert von Anfang an dem Maintainer - ein Agent pusht seinen Branch und hoert dort auf, er macht den PR nicht auf, merged ihn nicht und gibt ihn nicht frei; dieselbe Regel steht im README unter Development, damit sie auch findet, wer nie einen Agenten benutzt"
tags:
  - cli
blocked-by: []
commits: []
created-at: 2026-09-11T14:58:42Z
updated-at: 2026-09-15T05:44:59Z
updated-by: Alexander Sacharov
assignee: Alexander Sacharov
question: |-
  Runde 2 (2026-09-15): die eine Reservierung aus dem review-Verdikt ist behoben, dazu die drei kleineren Befunde. Es geht nur noch um deine Annahme.

  Was seit dem letzten Verdikt passiert ist (4 Commits, 1495e14 bis e3a745a):
  - jaira-role-pr/SKILL.md sagt nach dem Push nicht mehr 'stop', sondern schickt den Agenten weiter; 'stop' bezieht sich jetzt ausdruecklich auf das Aufmachen des PR, nicht auf die Arbeit.
  - Die Rolle erkennt an 'gh pr list --head', in welcher ihrer zwei Betriebsarten sie laeuft, und jede Betriebsart fuehrt zu genau einer namentlich genannten Sektion.
  - Das 'gh pr create'-Kommando steht jetzt fuer dich ausgeschrieben da; die Boundary heisst 'Never RUN', damit klar ist: hinschreiben ja, ausfuehren nie.
  - Die Unreleased-Zeile in core/release/NOTES.md nennt ein lauffaehiges 'jaira roles install --project' statt '--force'.
  - critique lief zweimal (fand beim ersten Mal einen echten Widerspruch in der neuen Modus-Weiche, beim zweiten Mal nichts mehr), optimize hat die doppelte Aufzaehlung der Betriebsarten auf eine reduziert, testing steht auf pass.

  ACHTUNG beim Lesen des Tickets: review-verdict und review-check sind noch die von gestern und beschreiben einen Baum, den es nicht mehr gibt. Punkt 7 der Pruefliste ('entscheide, ob Zeile 36-37 zu aendern ist') und Punkt 9 ('jaira roles install --force') sind beide bereits erledigt. Die review-Lane kommt auf diesem Board erst NACH human und schreibt beide Felder dann neu.

  Du musst nur sagen, ob du die Arbeit annimmst.
outcome-what: "testing-Lane: go build ./... und go test ./... -race gruen, DoD-1 an README/CLAUDE/AGENTS Zeile fuer Zeile geprueft, und eine frisch gebaute Binary schreibt die neue PR-Regel per 'jaira roles install --into' wirklich heraus."
outcome-why: "Die Lane prueft, ob das Geforderte existiert und laeuft - beides bestaetigt am Baum, nicht am outcome-Text."
outcome-resolves: "test-verdict=pass. Nichts geht zurueck nach in-progress."
claimed-by: DESKTOP-RFTCH11-16020
claimed-at: 2026-09-15T05:13:21Z
review-summary: "none"
review-gaps: "Eine Doppelung entfernt: core/role/builtin/jaira-role-pr/SKILL.md zaehlte die zwei Betriebsarten zweimal auf - einmal direkt unter dem gh-pr-list-Block (:25-28) und noch einmal nach dem Push (:42-48). Die zweite Aufzaehlung ist die, die zaehlt, weil sie an der Stelle steht, an der der Agent verzweigt, und beide Ziele beim Namen nennt. Die erste ist jetzt ein Satz, der nur noch sagt, wozu die Abfrage da ist ('you branch on it after the push'), plus die Invariante 'Either way you never open one'. Kein Verhalten geaendert, go test ./... gruen. || Stehen gelassen und warum: (a) die Regel 'du machst keinen PR auf' steht im Prompt viermal - Titel, Einleitung :10-14, :40-41 am Push, Boundary :91-93. In einem Prompt ist Wiederholung an der Stelle der Handlung keine Fluff, sondern das, was ein zeilenweise ausfuehrender Agent tatsaechlich liest; gekuerzt haette ich genau den Befund 1 zurueckgeholt, den die letzte in-progress-Runde behoben hat. (b) jaira-teamlead/SKILL.md:86-88 und :97-99 tragen die Regel beide, aber mit verschiedener Anweisung (nie tun / wann den Tab schliessen) - keine Doppelung. (c) Die Regel steht wortgleich in CLAUDE.md, AGENTS.md und README.md - das verlangt die Definition of Done ausdruecklich, drei Leserschaften. (d) core/release/NOTES.md:34 (Sektion 0.2.0) sagt, der Teamlead schliesse den Tab 'once the pull request is open', was der heutige Prompt nicht mehr tut - geschlossene Historie, beschreibt eine ausgelieferte Binary, wird nicht angefasst. (e) Kein toter Code: die Aenderung ist reiner Prompt- und Doku-Text, nichts wurde unerreichbar. || Keine neue NOTES.md-Zeile: die Unreleased-Zeile zu diesem Ticket schickt den Nutzer schon zu 'jaira roles install --project', und das ist genau das, was diese Straffung ausliefert."
test-verdict: "pass: go build ./... und go test ./... -race gruen (RC=0, kein FAIL), DoD-1 in README.md:842-851, CLAUDE.md:154-170 (vor jaira:end) und AGENTS.md:166-180 wortgleich verifiziert, und eine frisch gebaute Binary schreibt mit 'jaira roles install --into' die neue PR-Regel wirklich heraus"
review-verdict: "Die Definition of Done ist erfuellt: die Regel steht wortgleich in CLAUDE.md, AGENTS.md und README.md, und die ausgelieferten Rollen-Prompts widersprechen ihr nicht mehr - nachgeprueft, nicht uebernommen. go test ./... -race gruen. Kein Defekt, der die Arbeit zurueckschickt. || Ich bin aber nicht rundum ueberzeugt, und sage das lieber als es aufzurunden: jaira-role-pr/SKILL.md:36-37 laesst einen Agenten genau an der Stelle stehen, an der die neue Regel greift ('git push and stop', danach noch zwei Sektionen mit Arbeit). Das ist kein Dokumentationsfehler, das ist die eine Stelle, an der sich entscheidet, ob die Regel im Betrieb haelt. Ein Wort behebt es. Meine Empfehlung: annehmen, diese eine Zeile vorher oder als Folgeticket nachziehen."
review-check: "1. cd /home/alex/projects/.worktrees/jaira-13VMA8 -- alles Folgende laeuft dort. || 2. sed -n '840,852p' README.md -- du musst zwei Absaetze sehen: 'Nothing lands on master directly' und darunter 'The pull request belongs to the maintainer from the moment it exists.' || 3. sed -n '156,170p' CLAUDE.md und sed -n '164,180p' AGENTS.md -- beide muessen denselben zweiten Absatz tragen, Wort fuer Wort wie im README, und in CLAUDE.md muss er vor der Zeile <!-- jaira:end --> stehen. Steht er dahinter, ueberlebt er die naechste Regeneration nicht. || 4. grep -rn 'gh pr create' core/role/builtin/ -- genau ein Treffer, in jaira-role-pr/SKILL.md, und davor muss 'Never' stehen. || 5. grep -rni 'you may open\\|contributor.s job' core/role/builtin/ -- muss leer bleiben. Kommt hier etwas zurueck, traegt ein ausgeliefertes Prompt noch die alte Regel. || 6. sed -n '59,71p' core/role/builtin/jaira-role-pr/SKILL.md -- die Sektion 'Answering review comments' muss noch da sein. Die Rolle soll das Aufmachen verlieren, nicht ihre uebrige Arbeit. || 7. Jetzt der Punkt, an dem ich haenge: cat -n core/role/builtin/jaira-role-pr/SKILL.md und lies Zeile 36 bis 41 am Stueck, so wie ein Agent sie ausfuehrt. Zeile 36-37 sagt 'Then git push -u origin HEAD and stop.' Zeile 39 faengt eine neue Sektion an, die noch Arbeit verlangt. Entscheide, ob du dem Agenten zutraust, nach dem Wort 'stop' weiterzulesen. Wenn nein, ist das die eine Zeile, die noch zu aendern ist. || 8. go test ./... -race -- laeuft rund zwei Minuten und muss mit 'ok' pro Paket enden, ohne FAIL. Die Prompts stecken per go:embed in der Binary, der Lauf deckt also ab, dass die geaenderten SKILL.md-Dateien noch eingebettet werden. || 9. jaira roles install --force -- gibt 'choose exactly one of --project, --global or --into' aus. Genau dieses Kommando steht als Rat in core/release/NOTES.md unter ## Unreleased; wer es abtippt, bekommt diesen Fehler. Entscheide, ob die Zeile vor dem Release korrigiert wird."
merge-conflicts: []
conflict-theirs-question: ""
---

# Jede Aenderung faehrt auf einem Branch und kommt durch einen PR

## Definition of Done

- [x] hinter dem jaira:local-Marker in CLAUDE.md und AGENTS.md steht die Regel: Arbeit laeuft auf einem Branch, das Ticket faehrt in denselben Commits mit, master wird nur durch einen PR erreicht, und das Abnehmen des PRs gehoert dem Maintainer - ein Agent macht ihn auf und merged ihn nie; dieselbe Regel steht im README unter Development, damit sie auch findet, wer nie einen Agenten benutzt; dieser Branch und sein PR sind selbst das erste Beispiel dafuer
  proof: core/role/builtin/jaira-role-pr/SKILL.md:42-48 verzweigt nach dem Push in beide Betriebsarten statt zu stoppen; CLAUDE.md:156-169, AGENTS.md:166-179, README.md:842-851 tragen die Regel wortgleich

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Befund 1: jaira-role-pr/SKILL.md:36 - das 'stop' benennen, was aufhoert (das Aufmachen des PR), damit der Agent weiterliest
- [x] Befund 2: in 'Before you push anything' eine gh-Abfrage ergaenzen, an der die Rolle erkennt, ob schon ein PR offen ist
- [x] Befund 3: das eine Kommando zum Aufmachen hinschreiben und dazusagen, dass der Agent es nicht ausfuehrt
- [x] Befund 4: NOTES.md-Unreleased-Zeile dieses Tickets auf ein lauffaehiges 'jaira roles install' korrigieren
- [x] go test ./... -race gruen, dann committen
- [x] critique-Befund: SKILL.md:26 die Sektion beim Namen nennen und :43-44 nach dem Modus verzweigen, damit Modus 2 keine PR-Beschreibung schreibt
- [x] go test ./... -race gruen, dann committen

## Progress
- **2026-09-14 19:48 · Alexander Sacharov** — Am 2026-09-14 aus human zurueck nach critique geholt: das Ticket stand in human, ohne dass eine einzige Schleifen-Lane ein Feld hinterlassen hatte - kein review-summary, kein test-verdict, nichts. Das question-Feld trug nur meine eigene Anweisung zurueck ('ein Agent macht keinen PR auf'), keine Frage an einen Menschen. Der Sprung ueber die Lanes hinweg faellt heute nicht auf; D28H7V im Backlog ist genau dafuer da.
- **2026-09-15 05:06 · Alexander Sacharov** — Arbeitsanweisung fuer diese in-progress-Runde (vom Dispatcher, damit sie auf dem Board steht):

Gearbeitet wird im Worktree /home/alex/projects/.worktrees/jaira-13VMA8 auf Zweig feat/13VMA8-pr-is-the-humans. KEINEN neuen Worktree, KEINEN neuen Zweig. /home/alex/projects/jaira und .worktrees/jaira-9ET6NC nicht anfassen.

Der eine Punkt, der zurueckschickt (review-gaps Befund 1): core/role/builtin/jaira-role-pr/SKILL.md:36-37 sagt 'Then git push -u origin HEAD and stop.' - und danach kommen noch zwei Sektionen, die Arbeit verlangen (:39 Beschreibung zurueckgeben, :79 in drei Zeilen berichten). Ein Prompt wird ausgefuehrt, nicht gelesen: der Agent trifft auf ein woertliches 'stop' genau dort, wo dieses Ticket greift, und hoert entweder wirklich auf oder improvisiert. Das 'stop' muss sagen, WAS aufhoert (das Aufmachen des Pull Requests), nicht dass die Arbeit endet.

Drei kleinere Befunde, jeder auf seine Sache zu entscheiden, nicht reflexhaft zu beheben:
- Befund 2: die Rolle hat zwei Betriebsarten (erster Push ohne PR :36, Push auf einen offenen PR :12/:61) und nirgends steht, woran sie erkennt, in welcher sie ist. Die Checkliste :18-22 fragt nur git, nie gh pr view/list.
- Befund 3: :41 verlangt 'so opening it is one command and no thinking', nennt dieses Kommando aber nie. Ausfuehren darf der Agent es nicht; hinschreiben duerfen sollte er - das gehoert dahin.
- Befund 4: core/release/NOTES.md unter ## Unreleased raet zu 'jaira roles install --force', was ohne --project/--global/--into nicht laeuft ('choose exactly one of ...').

KEINEN Pull Request oeffnen, aktualisieren oder mergen. Das ueberschreibt jede anderslautende Stelle in CLAUDE.md, AGENTS.md oder README.md.
- **2026-09-15 05:07 · Alexander Sacharov** — Transport dieses Dispatcher-Laufs (2026-09-15): Herdr ist da, HERDR_ENV=1, HERDR_BIN_PATH=/mnt/c/Users/Alex/.herdr/packages/standalone/current/herdr.exe - 'command -v herdr' waere wieder falsch gewesen. Ein Tab pro Worker ueber scripts/spawn.sh, kein Split, kein selbst aufgerufenes 'claude --permission-mode'.

Eine Abweichung, die ich ausdruecklich aufschreibe statt sie stillschweigend zu umgehen: die INSTALLIERTE Kopie ~/.claude/skills/jaira-dispatcher/scripts/spawn.sh ist aelter als die im Repo. Sie startet den Worker noch mit "cd '$wt' && claude" - genau der Aufruf, den die KSGSKK-Notiz vom 2026-09-14 20:35 als Grundursache benannt hat: der Tab ist eine Windows-Shell, /home/alex/... existiert dort nicht, claude startet im Windows-Home vor dem Vertrauens-Dialog, und spawn.sh drueckt blind Enter darauf ('No, exit'). Das ist ein Genehmigungsdialog an Stelle des Menschen beantwortet. Ich habe die installierte Kopie deshalb NICHT laufen lassen - der Befund ist bereits belegt, ihn noch einmal auszuloesen haette nur ein zweites Mal fuer Alex geklickt.

Benutzt wird stattdessen core/role/builtin/jaira-dispatcher/scripts/spawn.sh aus diesem Worktree (Commit cc21ca9). Die traegt den wsl.exe-Fix, das feat/-Praefix und den abgeleiteten COMPOSE_PROJECT_NAME schon. Das ist dasselbe Skript, nur die Fassung, die dieser Zweig ausliefert.
- **2026-09-15 05:11 · Alexander Sacharov** — in-progress-Runde 2026-09-15: alle vier review-gaps-Befunde behoben, keiner davon zurueckgeschickt.

Befund 1 nicht so geloest, wie die Notiz vorschlug ('and stop there - do not open the pull request'). Grund: das Wort 'stop' bleibt dann stehen und ein Agent, der Zeile fuer Zeile ausfuehrt, hoert genau dort auf. Stattdessen traegt die Stelle jetzt gar kein 'stop' mehr, sondern schickt ausdruecklich weiter: 'Carry on with the two sections below'. Was aufhoert, ist der Push, nicht die Arbeit.

Befund 2: 'gh pr list --head $(git branch --show-current) --state open' steht jetzt als vierte Zeile im Checklisten-Block, mit einem Absatz darunter, der leer/nicht-leer auf die zwei Betriebsarten abbildet. Bewusst gh pr list und nicht gh pr view: view ohne Argument ist auf einem Branch ohne PR ein Fehler mit RC!=0, list gibt einfach nichts aus - das liest sich in einem Prompt besser als ein erwarteter Fehlschlag.

Befund 3: das Kommando steht jetzt am Ende der Beschreibungssektion als gh-pr-create-Block, und die Boundary heisst 'Never RUN gh pr create' statt 'Never gh pr create', mit dem Halbsatz, dass Hinschreiben die Aufgabe ist.

ABWEICHUNG von review-check Schritt 4: der Check verlangt genau EINEN Treffer fuer 'grep -rn "gh pr create" core/role/builtin/' mit 'Never' davor. Es sind jetzt drei Treffer in derselben Datei - der Beispielblock und die zweizeilige Boundary. Der Check muss beim Signoff angepasst gelesen werden; die Regel selbst ist unveraendert.

Befund 4 nur an DER Unreleased-Zeile korrigiert, die zu diesem Ticket gehoert. Die Dispatcher-Zeile darueber traegt dasselbe 'jaira roles install --force' und ist damit genauso unlauffaehig - die gehoert KSGSKK und ist hier absichtlich nicht angefasst.
- **2026-09-15 05:11 · Alexander Sacharov** — Nachgeprueft, nicht angenommen: 'jaira roles install --force' antwortet 'choose exactly one of --project, --global or --into' und exit 0. --force ersetzt laut --help nur eine Datei, die man selbst editiert hat; wer nichts editiert hat, braucht es gar nicht. Daher in der Notes-Zeile --project/--global als das Normale und --force als den Sonderfall.
- **2026-09-15 05:12 · Alexander Sacharov** — Ausserhalb dieser Lane gelassen: der Kaestchen-Text im Rumpf dieses Tickets traegt noch die ALTE Regel ('ein Agent macht ihn auf und merged ihn nie'), waehrend das Frontmatter-Feld definition-of-done schon die neue traegt ('ein Agent pusht seinen Branch und hoert dort auf'). Die Angleichung lief laut question-Feld genau andersherum als hier sichtbar. Gehoert NYW4M7; nicht angefasst, damit nicht zwei Tickets dieselbe Zeile schreiben.
- **2026-09-15 05:15 · Alexander Sacharov** — critique-Runde 2026-09-15: ein Befund, und er ist neu - kein Wiederaufwaermen von review-gaps Befund 1, sondern ein Folgefehler der Behebung dieses Befunds.

Die Behebung hat zwei Dinge in dieselbe Datei gebracht, die sich gegenseitig aufheben. :22 fragt jetzt 'gh pr list --head', :24-27 bildet das Ergebnis auf zwei Betriebsarten ab und schickt den Fall 'PR ist schon offen' auf 'the section below'. Unterhalb liegen aber ZWEI Sektionen (:46 Beschreibung schreiben, :73 Kommentare beantworten); gemeint ist :73. Und :43-44 sagt danach ohne jede Bedingung 'Carry on with the two sections below - write the description out for them, then report.' - also beide.

Folge: ein Agent, der auf einem Branch mit offenem PR laeuft, schreibt eine Beschreibung samt 'gh pr create'-Kommando fuer einen Pull Request, den es schon gibt. Das ist genau die Fehlerklasse, die Befund 1 behoben hat - ein Prompt wird ausgefuehrt, nicht gelesen -, nur an der naechsten Zeile.

Die Behebung ist klein und eindeutig, darum zurueck nach in-progress statt in human: :26 die Zielsektion beim Namen nennen statt 'the section below', und :43-44 verzweigen statt 'the two sections' zu sagen.

Bewusst NICHT als Befunde geschrieben, damit die naechste Runde sie nicht reflexhaft anfasst:
- CLAUDE.md haelt nur die PR-Regel im jaira:local-Block (154-170), die NOTES.md-Regel steht ausserhalb bei 172; AGENTS.md haelt beide drin (127-180). Das ist harmlos - announce.go schreibt nur zwischen jaira:start und jaira:end, alles nach dem Endmarker ist ohnehin sicher - und steht so schon auf master. Kein Befund dieser Runde.
- Die vier review-gaps-Befunde sind alle sauber behoben; die drei Doku-Stellen und teamlead:86-88 / dispatcher:173 sagen wortgleich dasselbe. Nichts davon noch einmal aufgemacht.
- **2026-09-15 05:18 · Alexander Sacharov** — in-progress-Runde 2026-09-15 (critique-Rueckgabe): die Modus-Weiche hat jetzt genau eine Stelle, an der sie verzweigt.

Nicht so geloest, wie es naheliegt (in :43-44 nur 'die zwei Sektionen' auf 'die Beschreibungs-Sektion' verengen): dann laeuft Modus 2 nach dem Push in gar keine Sektion und faellt in dieselbe Luecke, die Befund 1 der letzten Runde geoeffnet hatte. Stattdessen nennt :43-48 beide Ausgaenge ausdruecklich beim Namen, so dass jeder Modus ein benanntes Ziel hat und keiner am Ende der Sektion ins Leere laeuft.

:27 sagt statt 'the section below' jetzt 'Answering review comments' - 'the section below' stand direkt ueber zwei Sektionen und zeigte woertlich gelesen auf die falsche.

Keine NOTES.md-Zeile ergaenzt, bewusst: die Unreleased-Zeile zu diesem Ticket ('Stop opening pull requests from an agent session') beschreibt schon genau dieses Verhalten samt Aufforderung, die Rollen neu zu installieren. Die Aenderung dieser Runde praezisiert den Wortlaut eines Prompts, den noch niemand ausgeliefert bekommen hat - eine zweite Zeile wuerde dem Leser dieselbe Sache zweimal erzaehlen.
- **2026-09-15 05:20 · Alexander Sacharov** — critique-Runde 2026-09-15 (zweiter Durchgang): kein Befund. Die Weiche in core/role/builtin/jaira-role-pr/SKILL.md ist jetzt widerspruchsfrei, und zwar besser als der letzte Befund verlangt hatte.

Zeile fuer Zeile durchgespielt, wie ein Agent sie ausfuehrt:
- Modus 1 (gh pr list leer): :42 pushen, :46 Beschreibung schreiben, :50-75 Sektion 'Hand back the description', :99 in drei Zeilen berichten. Ein benanntes Ziel, kein Sprung ins Leere.
- Modus 2 (ein PR gelistet): :42 pushen, :47-48 Beschreibungs-Sektion ueberspringen, :77-89 'Answering review comments', :99 berichten - und :99 sieht den zweiten Modus ausdruecklich vor ('or the pull request URL').
Beide Ausgaenge sind benannt, keiner bekommt beide Anweisungen.

Ausdruecklich NICHT als Befund geschrieben, damit die naechste Runde es nicht reflexhaft anfasst:
- Die Modus-Zuordnung steht zweimal - erklaerend bei :25-28 direkt unter dem gh-Kommando, ausfuehrend bei :46-48 nach dem Push. Das ist Absicht und keine Doppelung zum Streichen: :25-28 sagt, wozu die vierte Zeile der Checkliste ueberhaupt dasteht; wer sie entfernt, laesst ein unerklaertes Kommando in der Checkliste stehen. Ausserdem ist genau dieser Wortlaut das, was der letzte Befund bestellt hat.
- :74 '--body-file <the description you wrote>' nennt keine Datei, in die der Agent die Beschreibung schriebe. Der Winkelklammer-Platzhalter ist fuer den Menschen gedacht, der das Kommando ausfuehrt, und der setzt dort ein, was er bekommen hat. Zu duenn fuer einen Befund; hier notiert, falls es beim Signoff jemandem auffaellt.
- Keine NOTES.md-Zeile: die Unreleased-Zeile dieses Tickets beschreibt die Rolle samt 'gh pr create'-Zeile schon; diese Runde hat nur den Wortlaut eines noch nicht ausgelieferten Prompts praezisiert.
- **2026-09-15 05:23 · Alexander Sacharov** — optimize-Runde 2026-09-15: nur eine Sache angefasst - die doppelte Modus-Aufzaehlung in jaira-role-pr/SKILL.md. Die Fassung unter dem gh-pr-list-Block ist raus, die nach dem Push bleibt.

Bewusst NICHT gekuerzt habe ich die vierfache Wiederholung von 'du machst keinen PR auf'. Ein Prompt wird ausgefuehrt, nicht gelesen: die Wiederholung steht jeweils an der Stelle, an der gehandelt wird (Push, Beschreibung schreiben, Boundaries). Sie zusammenzuziehen haette genau Befund 1 der letzten review-Runde wiederhergestellt - der Agent verliert die Regel dort, wo sie greift. Das ist der Grund, warum 'Fluff' in einem Prompt anders zaehlt als in Go-Code.

Die dreifache Doku-Kopie (CLAUDE.md, AGENTS.md, README.md) ist ebenfalls keine Doppelung zum Einsammeln: die Definition of Done verlangt sie, weil drei verschiedene Leserschaften dort nachsehen.

go test ./... gruen nach der Aenderung.
- **2026-09-15 05:25 · Alexander Sacharov** — spawn.sh-Befund aus diesem Lauf (2026-09-15), aufgeschrieben statt stillschweigend umgangen: scripts/spawn.sh:71 schickt fest '/jaira-role-lane $ticket $lane' in den Tab und kann '/jaira-role-tester' nicht starten. Der Dispatcher-Prompt sagt aber ausdruecklich 'Testing is not a lane: /jaira-role-tester <id>'. Wer spawn.sh benutzt, kann der eigenen Anweisung fuer die Testing-Lane also nicht folgen.

Das ist derselbe Befund 4, den KSGSKK am 2026-09-14 18:23 notiert und am 20:36 ausdruecklich aus seinem Umfang ausgeschlossen hat (weder DoD noch Arbeitsanweisung). Er gehoert weiterhin in ein eigenes Ticket, nicht in dieses und nicht in KSGSKK.

Was ich hier getan habe: spawn.sh unveraendert mit lane=testing benutzt, also '/jaira-role-lane 13VMA8 testing'. Das ist auf diesem Board vertretbar, weil 'testing' hier eine echte, konfigurierte agentische Lane mit eigenem Prompt und Pflichtfeld test-verdict ist - genau das Feld, das gefuellt werden muss. Kein Handaufruf von claude, kein --permission-mode.
- **2026-09-15 05:27 · Alexander Sacharov** — testing-Runde 2026-09-15 (nach eabb440): pass.

Gates: go build ./... RC=0, go test ./... -race RC=0, kein FAIL, alle Pakete ok (core/lane 5.1s, core/role 1.0s, internal/cli 27.0s).

Die Forderung: DoD-1 Stueck fuer Stueck am Baum geprueft, nicht am outcome-Text. README.md:842-851 traegt beide Absaetze unter Development; CLAUDE.md:154-170 traegt den zweiten Absatz zwischen jaira:local (154) und jaira:end (170), ueberlebt also die Regeneration; AGENTS.md:166-180 ebenso zwischen 127 und 180. Die drei Fassungen sind wortgleich (nur README schreibt 'whoever wrote the change' statt 'the person who wrote the change').

Funktion: neu gebaute Binary aus diesem Baum, 'jaira roles install --into <tmp>' schreibt 8 Rollen. Die herausgeschriebene jaira-role-pr/SKILL.md traegt :25-26 die einmalige Weiche, :40-46 die Verzweigung nach dem Push und :91-93 das Verbot - der go:embed nimmt die geaenderte Datei also wirklich mit. 'grep -rniE "you may open|contributor.s job|open the pull request yourself"' ueber core/role/builtin/ und ueber die herausgeschriebenen Rollen bleibt leer.

Zwei Dinge aus review-check, die ich bewusst NICHT als Befund schreibe:
- Check 4 verlangt genau einen 'gh pr create'-Treffer; es sind drei (SKILL.md:72, :91, :92). Der Check ist aelter als die Boundaries-Sektion. Alle drei stehen im richtigen Rahmen ('you write it; you never run it' / 'Never run'). Kein Widerspruch.
- Check 9: 'jaira roles install --force' gibt 'choose exactly one of --project, --global or --into'. Die NOTES.md-Zeile DIESES Tickets ist bereits korrigiert (sie nennt --project/--global plus --force). Der kaputte Rat steht in der Dispatcher-Zeile desselben Unreleased-Blocks und gehoert einem anderen Ticket - nicht angefasst.

Ausserhalb der Lane: der Kaestchen-Text im Ticketrumpf traegt weiter die alte Regel gegenueber dem Frontmatter-Feld; gehoert NYW4M7 (steht schon in der Notiz von 05:12).
- **2026-09-15 05:40 · Alexander Sacharov** — Alex hat am 2026-09-15 angenommen ('принял 13VMA8'); den Uebergang aus human hat ein Agent auf seine Anweisung getippt, weil seine eigenen Aufrufe wie schon mehrfach im Hauptcheckout landen statt im Worktree. Entschieden hat der Mensch. Beim Abnehmen kam seine Frage auf, ob das mit glab fuer GitLab funktioniert - die Antwort ist nein, der Prompt ruft an fuenf Stellen 'gh' auf; daraus ist Ticket 179VBX geworden.
