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
updated-at: 2026-09-15T05:07:23Z
updated-by: Alexander Sacharov
assignee: Alexander Sacharov
question: "Nichts blockiert: testing hat mit pass bestaetigt, dass kein Rollen-Prompt mehr erlaubt, einen PR aufzumachen, dass die drei Dokumentationsstellen und die Prompts dasselbe sagen, und dass jaira-role-pr seine uebrige Arbeit behalten hat. Beim Pruefen kam heraus, dass das Frontmatter-Feld dieses Tickets noch die alte Regel trug, waehrend das Kaestchen im Rumpf schon die neue hatte - ich habe es angeglichen und daraus Ticket NYW4M7 gemacht, weil es der dritte Fall an einem Tag war. Du musst hier nur sagen, ob du die Arbeit annimmst."
outcome-what: "Die doppelte Begruendung aus der neuen PR-Sektion gefaltet - CLAUDE.md und AGENTS.md sagen sie jetzt in der README-Formulierung"
outcome-why: "Absatz 2 war woertlich der Satz, den der generierte jaira-Block 40 Zeilen darueber schon traegt (core/board/announce.go:88-91)"
outcome-resolves: "Die Regel steht unveraendert an allen drei Stellen, nur ohne die Wiederholung; go test ./... gruen"
claimed-by: DESKTOP-RFTCH11-3086
claimed-at: 2026-09-15T05:06:05Z
review-summary: "Der Diff schreibt an drei Dokumentationsstellen dieselbe Regel auf und dreht danach die ausgelieferten Rollen-Prompts darauf um. CLAUDE.md:156-169 und AGENTS.md:166-179 (hinter dem jaira:local-Marker) und README.md:842-851 (unter Development) sagen jetzt wortgleich: nichts landet direkt auf master, die Aenderung faehrt auf einem eigenen Branch, das Ticket faehrt in denselben Commits mit, und 'der Pull Request gehoert dem Maintainer von dem Moment an, in dem er existiert' - ein Agent pusht seinen Branch und hoert dort auf. Die frueheren zwei Absaetze Begruendung sind zu einem Halbsatz gefaltet, damit die drei Kopien gleich lauten. || core/role/builtin/jaira-role-pr/SKILL.md ist von 'Open it, answer it, never accept it' zu 'Push it, hand it over, never open or accept it' umgeschrieben: Beschreibung, Einleitung und Boundaries sagen jetzt 'never gh pr create' zusaetzlich zu merge und approve, die Sektion 'Before you open anything' heisst 'Before you push anything' und endet mit 'git push -u origin HEAD', und die PR-Beschreibung wird nicht mehr aufgemacht sondern als fertiger Text an den Menschen zurueckgegeben. Push auf einen offenen PR und das Beantworten von Review-Kommentaren bleibt ausdruecklich Aufgabe der Rolle (SKILL.md:59-71). || core/role/builtin/jaira-teamlead/SKILL.md:79-81 verbietet dem Teamlead jetzt auch das Aufmachen, und der Tab-Schluss bei :90-93 haengt nicht mehr am offenen PR sondern am gepushten Branch. || core/release/NOTES.md bekommt eine Unreleased-Zeile dafuer."
review-gaps: "Vier Befunde, keiner davon ein Grund zurueckzuschicken, der erste aber vor dem Signoff zu beheben. || 1) AUSFUEHRBARKEIT, der ernsteste: core/role/builtin/jaira-role-pr/SKILL.md:36-37 sagt 'Then git push -u origin HEAD and stop.' - und danach kommen noch zwei Sektionen, die Arbeit verlangen (die Beschreibung zurueckgeben, ab :39; drei Zeilen berichten, :79-80). Ein Prompt wird ausgefuehrt, nicht gelesen: ein Agent, der von oben nach unten arbeitet, trifft genau in dem Moment, in dem sein Branch gepusht ist, auf ein woertliches 'stop' und hat keine Anweisung, die ihn weiterschickt. Entweder er hoert wirklich auf und der Mensch bekommt die versprochene fertige Beschreibung nie, oder er improvisiert. Ein Wort repariert es ('and stop there - do not open the pull request'), oder der Push wandert hinter die Beschreibungssektion. || 2) Die Rolle hat zwei Betriebsarten - erster Push ohne PR (:36) und Push auf einen PR, den jemand aufgemacht hat (:12, :61) - aber nirgends steht, woran sie erkennt, in welcher sie ist. Die Checkliste 'Before you push anything' (:18-22) fragt git, nie 'gh pr view' oder 'gh pr list'. An der Stelle, an der es zaehlt, muss der Agent raten. || 3) :41 verlangt 'so opening it is one command and no thinking', nennt dieses eine Kommando aber nicht - und da 'Never gh pr create' als Boundary danebensteht, ist nicht klar, ob der Agent es dem Menschen wenigstens hinschreiben darf. Er darf es nicht ausfuehren; das sollte dastehen. || 4) core/release/NOTES.md, Unreleased-Zeile: 'jaira roles install --force' ist so nicht lauffaehig - nachgestellt, das Kommando antwortet 'choose exactly one of --project, --global or --into'. Ausserdem braucht es --force nur, wer die Datei selbst editiert hat; sonst genuegt ein normales 'jaira roles install --project'. Die Zeile ist die einzige, die ein Nutzer je zu sehen bekommt, und wer sie abtippt bekommt einen Fehler. || Geprueft und in Ordnung: kein 'gh pr create' und keine Erlaubnis zum Aufmachen mehr irgendwo in core/role/builtin (grep ueber alle SKILL.md); jaira-role-pr hat Push und Kommentar-Beantwortung behalten; der Tab-Schluss des Dispatchers haengt an der human-Lane, nicht am PR, war also nie am Aufmachen verankert und brauchte keine Aenderung; go test ./... -race gruen. || Ausdruecklich nicht gewertet (gehoert KSGSKK): jaira-dispatcher/SKILL.md:157 und die Transport-Passagen."
test-verdict: "pass: alle sieben ausgelieferten Prompts sagen jetzt dasselbe wie die Dokumentation - kein 'may open' und keine andere Formulierung von 'mach den PR auf' mehr in core/role/builtin; jaira-role-pr behaelt Ordner- und frontmatter-Namen (core/role/role_test.go:18) und kann weiter zu einem offenen PR pushen und Review-Kommentare beantworten; NOTES.md traegt eine einzeilige Unreleased-Zeile; go test ./... -race gruen, RC=0"
review-verdict: "Die Definition of Done ist erfuellt: die Regel steht wortgleich in CLAUDE.md, AGENTS.md und README.md, und die ausgelieferten Rollen-Prompts widersprechen ihr nicht mehr - nachgeprueft, nicht uebernommen. go test ./... -race gruen. Kein Defekt, der die Arbeit zurueckschickt. || Ich bin aber nicht rundum ueberzeugt, und sage das lieber als es aufzurunden: jaira-role-pr/SKILL.md:36-37 laesst einen Agenten genau an der Stelle stehen, an der die neue Regel greift ('git push and stop', danach noch zwei Sektionen mit Arbeit). Das ist kein Dokumentationsfehler, das ist die eine Stelle, an der sich entscheidet, ob die Regel im Betrieb haelt. Ein Wort behebt es. Meine Empfehlung: annehmen, diese eine Zeile vorher oder als Folgeticket nachziehen."
review-check: "1. cd /home/alex/projects/.worktrees/jaira-13VMA8 -- alles Folgende laeuft dort. || 2. sed -n '840,852p' README.md -- du musst zwei Absaetze sehen: 'Nothing lands on master directly' und darunter 'The pull request belongs to the maintainer from the moment it exists.' || 3. sed -n '156,170p' CLAUDE.md und sed -n '164,180p' AGENTS.md -- beide muessen denselben zweiten Absatz tragen, Wort fuer Wort wie im README, und in CLAUDE.md muss er vor der Zeile <!-- jaira:end --> stehen. Steht er dahinter, ueberlebt er die naechste Regeneration nicht. || 4. grep -rn 'gh pr create' core/role/builtin/ -- genau ein Treffer, in jaira-role-pr/SKILL.md, und davor muss 'Never' stehen. || 5. grep -rni 'you may open\\|contributor.s job' core/role/builtin/ -- muss leer bleiben. Kommt hier etwas zurueck, traegt ein ausgeliefertes Prompt noch die alte Regel. || 6. sed -n '59,71p' core/role/builtin/jaira-role-pr/SKILL.md -- die Sektion 'Answering review comments' muss noch da sein. Die Rolle soll das Aufmachen verlieren, nicht ihre uebrige Arbeit. || 7. Jetzt der Punkt, an dem ich haenge: cat -n core/role/builtin/jaira-role-pr/SKILL.md und lies Zeile 36 bis 41 am Stueck, so wie ein Agent sie ausfuehrt. Zeile 36-37 sagt 'Then git push -u origin HEAD and stop.' Zeile 39 faengt eine neue Sektion an, die noch Arbeit verlangt. Entscheide, ob du dem Agenten zutraust, nach dem Wort 'stop' weiterzulesen. Wenn nein, ist das die eine Zeile, die noch zu aendern ist. || 8. go test ./... -race -- laeuft rund zwei Minuten und muss mit 'ok' pro Paket enden, ohne FAIL. Die Prompts stecken per go:embed in der Binary, der Lauf deckt also ab, dass die geaenderten SKILL.md-Dateien noch eingebettet werden. || 9. jaira roles install --force -- gibt 'choose exactly one of --project, --global or --into' aus. Genau dieses Kommando steht als Rat in core/release/NOTES.md unter ## Unreleased; wer es abtippt, bekommt diesen Fehler. Entscheide, ob die Zeile vor dem Release korrigiert wird."
merge-conflicts: []
conflict-theirs-question: ""
---

# Jede Aenderung faehrt auf einem Branch und kommt durch einen PR

## Definition of Done

- [x] hinter dem jaira:local-Marker in CLAUDE.md und AGENTS.md steht die Regel: Arbeit laeuft auf einem Branch, das Ticket faehrt in denselben Commits mit, master wird nur durch einen PR erreicht, und das Abnehmen des PRs gehoert dem Maintainer - ein Agent macht ihn auf und merged ihn nie; dieselbe Regel steht im README unter Development, damit sie auch findet, wer nie einen Agenten benutzt; dieser Branch und sein PR sind selbst das erste Beispiel dafuer
  proof: CLAUDE.md und AGENTS.md tragen die Regel hinter dem jaira:local-Marker, README unter Development; dieser Branch mit seinem PR ist das erste Beispiel

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

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
