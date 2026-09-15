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
commits:
  - f6ce687c74d229a4d37f9a99a856a72ad1f865f3
  - pending
  - 39c2659 79bb6de
created-at: 2026-09-11T14:58:42Z
updated-at: 2026-09-15T06:54:46Z
updated-by: Alexander Sacharov
assignee: Alexander Sacharov
question: |-
  GitLab-Runde fertig (2026-09-15). Es geht nur noch um deine Annahme.

  Alle vier DoD-Punkte abgehakt und belegt, test-verdict = pass.

  Was diese Runde gebracht hat (f6ce687, 9fc224c):
  - jaira-role-pr/SKILL.md hat eine neue Sektion 'Which forge this repository is on' (:36-67). 'git config jaira.forge' gewinnt immer; sonst entscheidet der Host des Push-Remotes; ist beides unklar, bleibt die Rolle stehen und sagt es, statt einen Befehl zu raten. Damit ist die Wahl waehlbar und nicht geraten - deine drei Bedingungen.
  - Beide Wege sind ausgeschrieben: Auflisten (:77 gh pr list --head / :83 glab mr list --source-branch), Aufmachen (:118 gh pr create / :124 glab mr create -d) und Boundaries (:143-148).
  - DoD 4 habe ich zusaetzlich selbst nachgelesen, weil ein fehlendes 'Never run' auf dem GitLab-Weg das ganze Ticket rueckgaengig machen wuerde, ohne dass man es sieht: :143-144 'never run gh pr create / gh pr merge / gh pr review --approve', :145-146 dasselbe mit glab mr create / merge / approve. Steht auf beiden Wegen.
  - Wortschatz wie von dir verlangt eingegrenzt: :60 sagt dem Leser, 'pull request' auf dem GitLab-Weg als 'merge request' zu lesen - kein globales Ersetzen, der GitHub-Weg redet weiter von pull requests.

  Zwei Dinge, die du wissen solltest, weil sie Urteile sind und keine Mechanik:
  1. critique hat diese Runde einen echten Fehler gefunden: die Forge wurde zuerst aus 'git config jaira.remote' abgeleitet, gepusht wird aber nach 'origin'. Das sind zwei verschiedene Remotes - jaira.remote ist der Remote der Ticket-Refs und loest hier zu 'upstream' auf. Auf diesem Board faellt es nicht auf, weil beide auf GitHub liegen; auf einem Fork mit Remotes auf verschiedenen Forges haette die Rolle den falschen Host untersucht und das falsche Werkzeug gewaehlt. Behoben in 9fc224c, die Begruendung steht als Prosa unter dem Block, damit es niemand 'verbessert'.
  2. testing hat nicht nur die Gates laufen lassen, sondern gegen das echte glab 1.114.0 geprueft: '--source-branch', '-d/--description' und 'glab mr approve' existieren, '--body-file' gibt es auf glab NICHT. Die create-Zeile waere sonst unlauffaehig gewesen - genau der Fehler, den dieses Ticket auf der GitHub-Seite schon einmal hatte.

  Nicht getan, wie angewiesen: kein Pull Request und kein Merge Request aufgemacht, aktualisiert oder gemerged. Der Zweig feat/13VMA8-pr-is-the-humans ist gepusht; das Kommando gibst du.

  Hinweis: review-verdict und review-check sind noch die von gestern und beschreiben einen Baum, den es nicht mehr gibt. Die review-Lane kommt auf diesem Board erst NACH human und schreibt beide dann neu.
outcome-what: "Critique uebersprungen"
outcome-why: "Alex: critique ist auf diesem Ticket zweimal ohne Befund geschlossen, die Aenderung ist ein Satz"
outcome-resolves: "DoD 3 - 'gibt der Remote nichts her und ist nichts gesetzt, sagt die Rolle das, statt den falschen Befehl zu raten' - ist jetzt nicht nur gesagt, sondern durchgesetzt: die Rolle kann den falschen Befehl nicht mehr raten, weil sie vor dem ersten Befehl stehenbleibt. Proof auf DoD 3 aktualisiert. go test ./... -race: Exit 0, kein FAIL; 'roles install --into' aus einer frisch gebauten Binary traegt den neuen Absatz, der go:embed-Pfad ist also mit."
claimed-by: DESKTOP-RFTCH11-19054
claimed-at: 2026-09-15T06:45:11Z
review-summary: "Die Rolle jaira-role-pr spricht jetzt zwei Forges. Neu ist die Sektion 'Which forge this repository is on' (SKILL.md:36-67): zuerst 'git config jaira.forge' - ist es gesetzt, gewinnt es ohne Wenn und Aber; sonst entscheidet der Host von 'git remote get-url origin', also des Remotes, auf den der Branch gepusht wird (github.com -> gh, Host mit 'gitlab' -> glab); gibt der Host nichts her, nennt die Rolle kein Werkzeug, sagt das und schreibt die eine Zeile hin, die es klaert ('git config jaira.forge gitlab'). Dass 'origin' und nicht 'jaira.remote' gelesen wird, steht mit Begruendung im Text (:46-49) - jaira.remote traegt die Ticket-Refs und ist im Fork das Upstream, waehrend der Branch zum Fork geht. Danach sind genau drei Stellen zweisprachig: Auflisten (:77 gh pr list --head / :83 glab mr list --source-branch), Aufmachen (:118 gh pr create --body-file / :124 glab mr create --description \"$(cat ...)\") und die Boundaries (:143-148). Der uebrige Ablauf bleibt einmalig, statt als zweite Kopie zu existieren. Der Wortwechsel ist begrenzt: :65-67 weist an, 'pull request' NUR auf dem GitLab-Weg als 'merge request' zu lesen - der GitHub-Weg redet weiter von Pull Requests. Ausserdem hat 9fc224c die Abfrage der offenen Requests aus der Forge-Sektion in die Push-Sektion verschoben, wo der Branch tatsaechlich genommen wird. Dazu eine Unreleased-Zeile in core/release/NOTES.md:20 und ein Satz in der SKILL-description."
review-gaps: |-
  Ein Befund, klein aber echt, und genau von der Sorte, die kein Test sieht: Zweig 4 der Forge-Leiter (SKILL.md:55-63, Host ist weder github.com noch gitlab-haltig, jaira.forge ungesetzt) laesst den ausfuehrenden Agenten ohne Anweisung fuer den Rest der Datei stehen. Er soll sagen, dass er es nicht entscheiden kann, und 'git config jaira.forge' nennen - aber es steht nirgends, ob er danach anhaelt oder weiterliest. Liest er weiter, steht er bei :71 vor einem unbedingten 'git push -u origin HEAD' und bei :74-84 vor einer Gabel 'auf GitHub ... oder auf GitLab ...', die er per Voraussetzung nicht aufloesen kann. Kein Widerspruch, aber eine Luecke: ein Satz wie 'Stop here and report; the rest of this file needs a settled forge' schliesst sie. Die Definition of Done ist davon nicht verletzt - DoD 3 verlangt nur, dass die Rolle es sagt statt zu raten, und das tut sie.

  Geprueft und in Ordnung: (a) Beide Wege am Stueck gelesen - die Verbots-Regel steht auf beiden, :143-144 gh pr create/merge/review --approve und :145-146 glab mr create/merge/approve, dazu :10-14, :86-87 'Either way you never open one' und :113-115 'You write it; you never run it'. Auf keinem Weg fehlt eines. (b) Der Fehler der frueheren Runde - ein woertliches Stop vor Sektionen, die noch Arbeit schulden - ist weg: :71-72 sagt, wo der Push endet, und liest unmittelbar weiter ('Now ask the forge ...'). Keine zwei widersprechenden Anweisungen mehr gefunden. (c) Der Forge-Host kommt von origin, dem Remote des Pushes, nicht von jaira.remote; die Begruendung steht im Prompt selbst (:46-49) und deckt sich mit dem Push bei :71, der ebenfalls fest auf origin geht - beide Stellen meinen denselben Remote. (d) Kein per Analogie zu gh erfundenes glab-Flag: benutzt werden nur 'mr list --source-branch', 'mr create --title/--description', 'mr merge', 'mr approve'. Kein --state, kein --head, kein --body-file auf der GitLab-Seite. (e) Die Vokabel-Umdeutung ist auf den GitLab-Weg begrenzt (:65-67), nicht global. (f) core/release/NOTES.md:20 ist genau eine Zeile unter ## Unreleased, beginnt mit '- ', nicht umgebrochen, und beschreibt, was der Leser anders tun muss. (g) go test ./... -race: RC=0, kein FAIL - die Prompts haengen per go:embed an der Binary, der Pfad ist also mitgetestet. (h) Nichts im outcome-Text steht ohne Deckung im Diff.
test-verdict: "fail: Gates gruen (go build + go test ./... -race RC=0) und die Funktion laeuft auf allen vier Forge-Sprossen, aber die Proofs von DoD 1, 2 und 4 zeigen seit 9fc224c sechs Zeilen zu frueh - DoD 2 belegt den GitLab-Weg mit den gh-Zeilen :83/:124 statt :89/:130"
review-verdict: "Angenommen, mit einer benannten Luecke. Die vier DoD-Punkte sind am Diff gedeckt, nicht nur am Bericht: beide Werkzeuge sind ausgeschrieben, die Wahl ist per 'git config jaira.forge' waehlbar und nur ersatzweise geraten, der unentscheidbare Host fuehrt zum Stehenbleiben statt zum falschen Befehl, die nie-ausfuehren-Regel steht auf beiden Wegen vollstaendig, und der Wortwechsel zu 'merge request' ist auf den GitLab-Weg begrenzt. go test ./... -race ist gruen (RC=0). Die eine Luecke - Zweig 4 der Forge-Leiter sagt nicht, ob danach weitergelesen oder angehalten wird - reicht meiner Einschaetzung nach nicht fuer eine Ruecksendung nach in-progress: sie trifft nur den selbstgehosteten Fall, in dem die Rolle ohnehin schon zugibt, dass sie es nicht weiss, und sie ist ein Satz. Ob sie hier noch geschlossen oder als Folgeticket gefuehrt wird, ist Alex' Entscheidung, nicht meine."
review-check: "Alle Schritte hier sind vor dem Eintragen einmal so gelaufen; die genannten Ausgaben sind die, die dabei kamen.\n\n1. In den Worktree gehen: cd /home/alex/projects/.worktrees/jaira-13VMA8\n2. Tests: go test ./... -race - erwartet wird Exit-Code 0 und kein FAIL in der Ausgabe.\n3. Binary bauen: go build -o /tmp/jaira-check ./cmd/jaira - erwartet wird Exit-Code 0 und keine Ausgabe.\n4. Prompts aus der frischen Binary herausschreiben: mkdir -p /tmp/rollen && /tmp/jaira-check roles install --into /tmp/rollen - die glab-Zeilen muessen wirklich in der Binary stecken (go:embed), nicht nur im Quellbaum.\n5. Nachsehen, dass die Verbots-Regel auf BEIDEN Wegen steht: grep -n 'never run' /tmp/rollen/jaira-role-pr/SKILL.md - es muessen ZWEI Zeilen kommen, eine mit 'gh pr create' und eine mit 'glab mr create'. Kommt nur eine, ist das Ticket nicht erledigt.\n6. Nachsehen, dass kein gh-Flag auf die GitLab-Seite geraten wurde: grep -n 'glab' /tmp/rollen/jaira-role-pr/SKILL.md - erwartet werden genau vier glab-Kommandos: 'mr list --source-branch', 'mr create --title/--description', 'mr merge', 'mr approve'. Taucht irgendwo '--body-file' oder '--head' hinter einem glab auf, ist es falsch: beide Flags gibt es in glab nicht.\n7. Die Forge-Erkennung von Hand nachstellen. Ein Wegwerf-Repo mit GitLab-Remote:\n   mkdir -p /tmp/fx-gl && cd /tmp/fx-gl && git init -q . && git remote add origin https://gitlab.com/a/b.git\n   git config jaira.forge   -> gibt nichts aus und endet mit Exit 1 (nicht gesetzt)\n   git remote get-url origin   -> https://gitlab.com/a/b.git\n   Nach Regel 3 in SKILL.md:54 enthaelt der Host 'gitlab', also ist das Werkzeug glab.\n8. Denselben Test mit einem Host, der es nicht verraet:\n   mkdir -p /tmp/fx-sh && cd /tmp/fx-sh && git init -q . && git remote add origin git@git.esprit-engineering.de:team/requirementsgenie.git\n   git remote get-url origin   -> git@git.esprit-engineering.de:...\n   Weder 'github.com' noch 'gitlab' im Host: nach SKILL.md:55-63 darf hier KEIN Werkzeug genannt werden, die Rolle sagt, dass sie es nicht entscheiden kann, und nennt 'git config jaira.forge gitlab'.\n9. Dass die Einstellung wirklich sticht: im selben Ordner git config jaira.forge github ausfuehren, dann git config jaira.forge -> github. Ab jetzt gewinnt das laut SKILL.md:51-52 gegen jeden Remote.\n10. Zuletzt die Datei selbst einmal von oben nach unten lesen, wie ein Agent sie ausfuehrt: sed -n '36,153p' core/role/builtin/jaira-role-pr/SKILL.md. Auf keinem der beiden Wege darf ein Schritt ohne Anweisung enden oder zwei sich widersprechende tragen. Die einzige Stelle, an der das noch klemmt, steht in review-gaps: Zweig 4 bei :55-63 sagt nicht, ob danach angehalten oder weitergelesen wird."
merge-conflicts: []
conflict-theirs-question: ""
---

# Jede Aenderung faehrt auf einem Branch und kommt durch einen PR

## Definition of Done

- [x] hinter dem jaira:local-Marker in CLAUDE.md und AGENTS.md steht die Regel: Arbeit laeuft auf einem Branch, das Ticket faehrt in denselben Commits mit, master wird nur durch einen PR erreicht, und das Abnehmen des PRs gehoert dem Maintainer - ein Agent macht ihn auf und merged ihn nie; dieselbe Regel steht im README unter Development, damit sie auch findet, wer nie einen Agenten benutzt; dieser Branch und sein PR sind selbst das erste Beispiel dafuer
  proof: core/role/builtin/jaira-role-pr/SKILL.md:86-92 verzweigt nach dem Push in beide Betriebsarten statt zu stoppen; CLAUDE.md:156-169, AGENTS.md:166-179, README.md:842-851 tragen die Regel wortgleich
- [x] Die Rolle arbeitet auf GitLab wie auf GitHub: sie listet die offenen Merge Requests des aktuellen Zweigs mit 'glab mr list --source-branch' und schreibt dem Menschen eine lauffaehige 'glab mr create'-Zeile aus, so wie sie es auf GitHub mit 'gh pr list' und 'gh pr create' tut. Nachgestellt auf einem Fixture mit einem GitLab-Remote - das requirementsgenie-Board auf git.esprit-engineering.de ist der echte Fall.
  proof: core/role/builtin/jaira-role-pr/SKILL.md:83 (glab mr list --source-branch) und :124 (glab mr create --title/--description); nachgestellt auf einem git-Fixture mit Remote git@git.esprit-engineering.de:team/requirementsgenie.git -> 'tool: gitlab, would run: glab mr list --source-branch feat/X'; Flags gegen glab 1.114.0 --help geprueft
- [x] Welches Werkzeug laeuft, ist waehlbar und nicht nur geraten: aus dem Remote abgeleitet, wenn er es hergibt, und ausdruecklich setzbar, wenn nicht oder wenn der Mensch es anders will. Gibt der Remote nichts her und ist nichts gesetzt, sagt die Rolle das, statt den falschen Befehl zu raten.
  proof: core/role/builtin/jaira-role-pr/SKILL.md:36-68: git config jaira.forge gewinnt immer, sonst der Host von 'git remote get-url origin'; github.com -> gh, Host mit 'gitlab' -> glab, sonst nennt Sprosse 4 kein Werkzeug und haelt dort an - ':65-68' sagt ausdruecklich kein Push und keine Abfrage offener Requests, und nennt dem Menschen 'git config jaira.forge' als das, was die Rolle beim naechsten Start weiterlaufen laesst
- [x] Die Regel steht auf BEIDEN Wegen und stimmt: die Zeile wird ausgeschrieben und nie ausgefuehrt, nichts wird gemerged, nichts freigegeben. Nachgestellt, indem beide Wege gelesen werden - auf keinem darf ein 'Never run' fehlen, und das Wort Merge Request ersetzt pull request nur dort, wo von GitLab die Rede ist.
  proof: core/role/builtin/jaira-role-pr/SKILL.md:143-148: 'never run gh pr create/gh pr merge/gh pr review --approve' UND 'never run glab mr create/glab mr merge/glab mr approve', dazu :113-115 'You write it; you never run it' fuer beide; Wortwechsel nur bei :65-67, nicht global

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
- [x] SKILL.md: Forge-Erkennung in 'Before you push anything' - git config jaira.forge, sonst Remote-Host, sonst stehenbleiben und die Einstellung nennen
- [x] SKILL.md: Abfrage der offenen Requests, die create-Zeile und die Boundaries nennen beide Werkzeuge (gh/glab); 'Merge Request' nur auf der GitLab-Seite
- [x] Beide Wege am Stueck von oben nach unten lesen - auf keinem darf 'Never run'/'never merge'/'never approve' fehlen
- [x] Erkennung auf Fixtures nachstellen: gitlab.com-Remote, github.com-Remote, selbstgehosteter Host + jaira.forge
- [x] NOTES.md: Unreleased-Zeile fuer den GitLab-Weg
- [x] go test ./... -race gruen, dann committen
- [x] critique-Befund 1+2: SKILL.md:43 die Forge vom Push-Remote ablesen (git remote get-url origin), die jaira.remote-Zeile streichen
- [x] critique-Befund 3: den Block 'ist schon einer offen' aus der Forge-Sektion in die Push-Sektion verschieben
- [x] go test ./... -race gruen, dann committen
- [x] Sprosse 4 der Forge-Leiter: die Rolle dort anhalten - kein Push, keine PR-Abfrage - und dem Menschen sagen, dass er git config jaira.forge setzt und neu startet
- [x] NOTES.md: Unreleased-Zeile fuer das Anhalten bei unbekannter Forge
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
- **2026-09-15 05:47 · Alexander Sacharov** — Arbeitsanweisung fuer diese Runde (Dispatcher, 2026-09-15), damit sie auf dem Board steht und nicht in einer Sitzung stirbt.

Worktree /home/alex/projects/.worktrees/jaira-13VMA8, Zweig feat/13VMA8-pr-is-the-humans. KEIN neuer Worktree, KEIN neuer Zweig. /home/alex/projects/jaira und .worktrees/jaira-9ET6NC nicht anfassen.

Neu an diesem Ticket: Alex hat die GitLab-Arbeit hier hineingefaltet statt sie als eigenes Ticket 179VBX zu fuehren (das ist archiviert). Grund: es ist dieselbe Datei. core/role/builtin/jaira-role-pr/SKILL.md hat gerade vier Runden hinter sich, und ein zweites Ticket haette spaeter um dieselben Zeilen gekaempft.

Der Befund: die Rolle ruft 'gh' an :22, :72 und :91-92 auf. Keiner dieser Befehle existiert auf GitLab. Das ist nicht hypothetisch - /home/alex/projects/requirementsgenie hat ein .jaira/-Board, dessen einziger Remote git.esprit-engineering.de ist, und 'glab' ist auf diesem Rechner installiert, sowohl als Windows-exe als auch unter ~/.local/bin/glab.

DoD 2 und 3 in Alex' eigenen Worten: welches Werkzeug laeuft, muss WAEHLBAR sein, nicht bloss geraten. Aus dem Remote abgeleitet, wo der Remote es hergibt; ausdruecklich setzbar, wo er es nicht hergibt oder wo der Mensch es anders will; und wo keins von beidem greift, sagt die Rolle das, statt einen Befehl zu raten, der fehlschlaegt.

DoD 4 ist der Punkt, der nicht durchrutschen darf: die Regel, fuer die dieses Ticket ueberhaupt existiert, muss auf BEIDEN Wegen stehen. Ein zweiter Weg durch den Prompt ist eine zweite Stelle, an der 'die Zeile ausschreiben, nie ausfuehren, nie mergen, nie freigeben' vorhanden UND richtig sein muss. Beide Wege am Stueck von oben nach unten lesen, bevor abgehakt wird - ein fehlendes 'Never run' auf dem GitLab-Weg macht das ganze Ticket rueckgaengig und sieht dabei fertig aus.

Wortschatz: GitLab nennt es Merge Request. 'pull request' NUR dort ersetzen, wo tatsaechlich von GitLab die Rede ist - nicht global, sonst redet der GitHub-Weg von MRs.

Zu beachten, weil es gerade erst erarbeitet wurde: die Modus-Weiche aus den letzten Runden (:24-27 fragt 'gh pr list --head', :42-48 verzweigt danach in genau eine namentlich genannte Sektion) ist das Ergebnis von zwei critique-Runden. Der GitLab-Weg braucht dieselbe Weiche mit 'glab mr list --source-branch', nicht eine zweite, anders gebaute.

KEINEN Pull Request und KEINEN Merge Request oeffnen, aktualisieren oder mergen. Das ueberschreibt jede anderslautende Stelle in CLAUDE.md, AGENTS.md oder README.md.

Hinweis zum Nachbarticket: KSGSKK liegt im selben Worktree und steht nach der dritten critique-Rueckgabe still, bis Alex entscheidet. Es fasst jaira-dispatcher/SKILL.md, jaira-teamlead/SKILL.md und spawn.sh an - nicht jaira-role-pr/SKILL.md. Keine Kollision, aber die Dateien der anderen Baustelle bleiben unangetastet.
- **2026-09-15 05:52 · Alexander Sacharov** — in-progress-Runde 2026-09-15 (GitLab-Weg): die Rolle laeuft jetzt auf beiden Forges, und zwar als EIN Ablauf mit forge-spezifischen Kommandozeilen, nicht als zwei Sektionen.

Warum nicht zwei Sektionen (naheliegend, aber verworfen): eine Kopie des ganzen Ablaufs fuer GitLab heisst, dass jede kuenftige Aenderung an zwei Stellen richtig sein muss - und DoD 4 sagt genau, was dann passiert: auf einem der beiden Wege fehlt irgendwann das 'Never run', und es sieht dabei fertig aus. Die drei Stellen, an denen die Werkzeuge sich wirklich unterscheiden, sind Auflisten (:64-74), Aufmachen (:109-121) und Boundaries (:139-144); ueberall sonst ist der Text derselbe. Darum nennen diese drei Stellen beide Kommandos und der Rest bleibt einmalig.

Die Modus-Weiche aus den vorigen Runden ist unveraendert geblieben - dieselbe Weiche, nur mit zwei moeglichen Abfragekommandos davor. Keine zweite, anders gebaute.

Warum 'git config jaira.forge' und nicht ~/.jaira/settings.json: die Wahl gehoert dem Repository, nicht dem Rechner. Dasselbe Argument, mit dem 'git config jaira.remote' in dieser Version eingefuehrt wurde (NOTES.md ## Unreleased) - wer auf einem Rechner ein GitHub- und ein GitLab-Board hat, kann mit einer maschinenweiten Einstellung nichts anfangen. Der Prompt liest 'jaira.remote' beim Ableiten auch gleich mit, statt blind 'origin' anzunehmen.

Dead end, aufgeschrieben damit es niemand zweimal probiert: den selbstgehosteten Fall automatisch zu erkennen geht nicht sauber. 'git.esprit-engineering.de' traegt weder 'github' noch 'gitlab' im Namen. Denkbar waere, glabs eigene Hosts-Konfiguration (~/.config/glab-cli/config.yml) zu lesen oder probeweise 'glab mr list' laufen zu lassen - beides ist Maschinerie in einem Prompt, und der Probelauf authentifiziert im Zweifel gegen das falsche Projekt. Deshalb bleibt die Rolle dort ausdruecklich stehen und nennt die eine Zeile, die es fuer immer klaert. Genau das verlangt DoD 3 im dritten Satz.

Was ich an glab 1.114.0 nachgesehen und nicht geraten habe:
- 'glab mr list' listet ohne Flag bereits nur offene MRs; ein Gegenstueck zu '--state open' gibt es nicht, '--source-branch' ist das Gegenstueck zu '--head'.
- 'glab mr create' hat kein '--body-file'. '-d/--description' nimmt einen String, darum steht in der ausgeschriebenen Zeile '--description "$(cat <die Beschreibung>)"'. '--source-branch' ist unnoetig: der aktuelle Zweig ist die Vorgabe.
- Freigeben heisst auf GitLab 'glab mr approve', nicht 'review --approve'.

Fixtures: drei git-Repos mit Remote github.com / gitlab.com / git.esprit-engineering.de, je einmal ohne und mit gesetztem jaira.forge - alle vier Zweige der Regel aus :46-58 laufen so, wie der Prompt sie beschreibt. Danach frisch gebaute Binary, 'jaira roles install --into' schreibt die glab-Zeilen wirklich heraus (go:embed nimmt die Datei mit) und 'jaira update' gibt die neue NOTES-Zeile zurueck.

Fuer die naechste review-Runde: review-check Schritt 4 verlangt genau EINEN Treffer fuer 'gh pr create'. Der Check war schon vor dieser Runde ueberholt (drei Treffer), jetzt kommen die 'glab mr create'-Treffer dazu. Der Check gehoert neu geschrieben, die Regel ist unveraendert.
- **2026-09-15 05:55 · Alexander Sacharov** — critique (GitLab-Runde, 2026-09-15): drei Befunde, alle in core/role/builtin/jaira-role-pr/SKILL.md. Der erste ist der, der weh tut.

1. Der Prompt liest die Forge vom falschen Remote ab (:43 gegen :81). :43 nimmt 'git config jaira.remote || echo origin' und leitet daraus gh/glab ab; :81 pusht dann fest nach 'origin'. Das sind zwei verschiedene Dinge: jaira.remote ist der Remote, auf dem die TICKET-REFS liegen, nicht der, auf dem der Branch landet. In genau diesem Repo faellt das auseinander - 'jaira whoami --json' sagt remote=upstream (aus settings.json), gepusht wird nach origin. Dass beide hier zufaellig auf github.com zeigen, verdeckt es nur; auf einem Fork, dessen Board auf einem GitLab-Upstream liegt, waehlt der Prompt glab und pusht nach GitHub. Die Frage, die der Prompt beantworten will, lautet: auf welcher Forge liegt der Zweig, den ich gerade gepusht habe. Also den Host des Remotes lesen, auf den :81 pusht ('git remote get-url origin'), und die jaira.remote-Zeile ersatzlos streichen. jaira.forge als ausdrueckliche Uebersteuerung bleibt richtig und unveraendert.

2. Dieselbe Zeile :43 ist eine zweite, verkuerzte Kopie einer Leiter, die es schon gibt. core/settings/settings.go:145-168 (RemoteFor/RemoteSourceFor) definiert den Board-Remote in vier Schritten - jaira.remote, settings.json, der einzige Remote, sonst laut scheitern - und der Kommentar darueber warnt woertlich davor, die Leiter ein zweites Mal hinzuschreiben, weil dann die Stelle, deren einzige Aufgabe die Wahrheit ueber den Remote ist, einen anderen Namen nennt als der Code. Der Prompt schreibt sie mit zwei von vier Schritten hin. Falls Befund 1 anders entschieden wird und der Board-Remote doch gebraucht wird, ist 'jaira whoami --json' (Feld remote, dazu remote_source im Klartext) der Weg - das ist genau der Befehl, den es dafuer gibt.

3. Sektionsgrenze: :64-77 laesst die Abfrage der offenen Requests in der Sektion 'Which forge this repository is on' laufen, die damit drei Aufgaben traegt. Die Verzweigung wird aber erst in '## Push, then take the branch the listing put you on' (:79) ausgewertet. Den Block dorthin verschieben; die Forge-Sektion klaert dann nur noch das Werkzeug und den Wortwechsel. Kein Verhalten, nur die Stelle - faellt beim Fix von 1 ohnehin an, weil :43 und :81 an dieser Grenze liegen.

Stehen gelassen und warum: (a) Ein Ablauf statt zwei Kopien - richtig entschieden, die Begruendung des Implementierers zu DoD 4 traegt. (b) Die vierfache Wiederholung von 'never open/merge/approve' bleibt, das ist in einem Prompt kein Fluff (schon in der letzten optimize-Runde so entschieden, wird nicht wieder aufgemacht). (c) 'glab mr list' ohne State-Flag und '--description "$(cat ...)"' statt --body-file: gegen glab 1.114.0 nachgesehen, das ist die review-Lane, nicht meine. (d) Der veraltete review-check (erwartet EINEN Treffer fuer 'gh pr create') ist ein Ticketfeld, kein Diff-Befund - steht schon in der Notiz der in-progress-Runde und gehoert der review-Lane.
- **2026-09-15 06:00 · Alexander Sacharov** — in-progress-Runde 2026-09-15 (critique-Rueckgabe, GitLab): die Forge wird jetzt von dem Remote abgelesen, auf den auch gepusht wird.
- Nachgeprueft, nicht angenommen: 'jaira whoami --json' meldet in diesem Repo remote=upstream (remote_source='from settings.json on this machine'), waehrend der Branch nach origin geht. Die alte Zeile haette hier also den Host des Upstream gelesen und den des Push-Ziels nie gesehen - im Fork sind das zwei verschiedene Hosts und im schlimmsten Fall zwei verschiedene Forges.
- Warum nicht 'jaira whoami --json' lesen, wie critique als Ausweichweg anbot: der Prompt braucht den Board-Remote gar nicht. Er braucht den Remote, auf den 'git push -u origin HEAD' schiebt, und der steht wortwoertlich in dieser Zeile. Damit ist auch die zweite Kopie der Remote-Leiter aus core/settings/settings.go:145-168 weg, ohne dass eine dritte Quelle dazukommt.
- Der Grund steht als Prosa unter dem Codeblock, nicht als Kommentar in der Zeile: wer 'git remote get-url origin' sieht, fragt sich genau dann nach jaira.remote, wenn er den Prompt liest - und ohne die drei Zeilen holt die naechste Runde die alte Fassung zurueck.
- Befund 3 (Block verschoben): die Sektion 'Which forge this repository is on' endet jetzt beim Wortwechsel gh/glab. Die Push-Sektion heisst 'Push, then ask which of your two jobs this is' - erst pushen, dann listen, dann verzweigen. Vorher stand die Abfrage vor dem Push und die Verzweigung dahinter, mit 'du branchst erst nach dem Push' als Klammer dazwischen; die Klammer ist ersatzlos weg, weil die Reihenfolge sie jetzt selbst erzaehlt.
- Keine neue NOTES.md-Zeile: die Unreleased-Zeile zum GitLab-Weg sagt 'otherwise the remote host does', und das bleibt nach dieser Korrektur wahr. Ausgeliefert ist der Prompt noch nicht; korrigiert wurde ein Wortlaut, den noch niemand hat.
- **2026-09-15 06:02 · Alexander Sacharov** — critique (GitLab-Runde, vierter Durchgang, 2026-09-15): nichts mehr zu sagen. Alle drei Befunde der letzten Runde sind in 9fc224c behoben und ich habe sie am Baum nachgelesen, nicht am outcome-Text geglaubt: (1) die Forge wird jetzt in SKILL.md:43 von 'git remote get-url origin' abgelesen - demselben Remote, nach dem :71 pusht; (2) die zweite Kopie der Remote-Leiter aus core/settings/settings.go:145-168 ist ersatzlos weg, es bleibt genau ein git-config-Griff (jaira.forge), und der ist eine eigene Einstellung, keine Nachbildung; (3) die Abfrage der offenen Requests steht in '## Push, then ask which of your two jobs this is' (:69-91), wo die Verzweigung auch ausgewertet wird - die Forge-Sektion klaert nur noch Werkzeug und Wortwahl.

Gegengelesen und stehen gelassen, damit die naechste Runde es nicht wieder aufmacht: (a) Die drei Zeilen Prosa :46-49 ('Not off jaira.remote') sind Begruendung und keine Handlung - in einem Prompt sind sie trotzdem richtig, weil jaira.remote an anderer Stelle des Projekts der normale Griff ist (README.md:350-355, core/settings/settings.go:149) und ein Agent ihn sonst fuer die bessere Wahl haelt. Ohne den Satz kommt Befund 1 durch die Hintertuer zurueck. (b) Die Forge-Sektion steht vor dem Push, obwohl die Antwort erst danach gebraucht wird - absichtlich: 'settle it before you type either' ist die Aussage, und die Platzierung der Listung wurde gerade erst entschieden, ich mache sie nicht in der Gegenrichtung wieder auf. (c) Ein Ablauf mit drei forge-spezifischen Stellen statt zwei Sektionen - unveraendert richtig, DoD 4 haengt daran. (d) 'glab mr list' ohne State-Flag und '--description "$(cat ...)"' statt --body-file: gegen glab 1.114.0 nachgesehen, ob es stimmt, ist die Aufgabe der review-Lane, nicht meine. (e) Der veraltete review-check (erwartet EINEN Treffer fuer 'gh pr create') bleibt ein Ticketfeld und kein Diff-Befund; steht seit der in-progress-Notiz auf dem Ticket und gehoert der review-Lane.
- **2026-09-15 06:06 · Alexander Sacharov** — optimize (GitLab-Runde, 2026-09-15): nichts zu entfernen, keine Datei angefasst. Was ich nachgesehen habe, damit die naechste Runde es nicht wiederholt: (1) es gibt im ganzen Repo nur eine Stelle, die eine Forge kennt oder gh/glab tippt - jaira-role-pr/SKILL.md; dispatcher und teamlead tragen die PR-Regel als Prosa ohne Werkzeugnamen. (2) Nichts zeigt auf den alten Aufbau: 'the gh pr list above' ist mit 9fc224c verschwunden, und kein Go-Code liest den Prompttext (core/role/role_test.go:18 kennt nur den Rollennamen, nicht den Inhalt). (3) Die einzige neue Wiederholung ist :72 gegen :86 ('you do not open the pull request' / 'Either way you never open one'), entstanden beim Verschieben der Listung. Sie bleibt stehen: die eine steht am Push, die andere an der Verzweigung, und die Runde davor hat schon entschieden, dass Wiederholung an der Stelle der Handlung in einem Prompt kein Fluff ist. || Offen und nicht meins: NOTES.md haelt zwei Unreleased-Zeilen zu dieser Rolle, und die aeltere ('hands you ... the gh pr create line') stimmt seit heute nur noch auf GitHub. Sie gehoert einer anderen Aenderung; wer sie anfasst, sollte es bewusst tun.
- **2026-09-15 06:17 · Alexander Sacharov** — testing-Lane 2026-09-15 (GitLab-Runde): pass. Was ich wirklich laufen liess, nicht nur gelesen:

Gates: go build ./... RC=0, go test ./... -race RC=0, kein FAIL. core/role 1.034s gruen - das ist das Paket, das die eingebetteten Prompts traegt.

DoD 3 (Werkzeugwahl) habe ich nicht am Text abgenommen, sondern die Regel aus SKILL.md:46-63 auf sechs frischen git-Fixtures durchgespielt: Remote github.com / gitlab.com / git.esprit-engineering.de, je einmal ohne und mit 'git config jaira.forge gitlab'. Alle sechs verhalten sich wie beschrieben, einschliesslich des Falls, der am leichtesten falsch waere: gesetztes jaira.forge=gitlab schlaegt einen github.com-Remote. Der selbstgehostete Host ohne Einstellung bleibt stehen und nennt die eine Zeile, statt zu raten.

glab-Flags gegen das echte Binary geprueft (glab 1.114.0, /home/alex/.local/bin/glab), weil der Prompt sonst Befehle ausschreibt, die auf dem Rechner des Menschen scheitern:
- 'glab mr list' sagt im Hilfetext selbst 'Defaults to open merge requests' - das fehlende Gegenstueck zu --state open ist also richtig und kein Versehen.
- 'glab mr create' hat -d/--description und KEIN --body-file. Die ausgeschriebene Zeile mit --description "$(cat ...)" ist damit die einzige, die laeuft.
- 'glab mr approve' existiert als eigener Befehl. Die Boundary nennt ihn richtig, nicht 'review --approve'.

Funktion: frisch gebaute Binary, 'jaira roles install --into <scratch>' schreibt 8 Rollen, und die herausgeschriebene jaira-role-pr/SKILL.md traegt die glab-Zeilen an :83, :124, :145-146. go:embed nimmt die Aenderung also wirklich mit - ohne diesen Schritt haette ich nur die Quelldatei bestaetigt, nicht das, was ein Nutzer bekommt.

Offen und ausdruecklich NICHT von mir angefasst (Lane-Grenze): review-check Schritt 4 verlangt genau EINEN Treffer fuer 'gh pr create'; im Baum sind es drei, plus jetzt die glab-Gegenstuecke. Der Check ist ueberholt, die Regel dahinter unveraendert. Die review-Lane muss ihn neu schreiben - das steht schon in der in-progress-Notiz und ist hier bestaetigt, nicht behoben.

Ebenfalls nicht angefasst: im Worktree liegen uncommittete Aenderungen des Nachbartickets KSGSKK (jaira-dispatcher/SKILL.md, spawn.sh, NOTES.md). Gehoeren nicht zu dieser Runde, bleiben ungestaged.
- **2026-09-15 06:33 · Alexander Sacharov** — Alex hat die GitLab-Runde am 2026-09-15 angenommen; den Uebergang human -> review hat ein Agent auf seine Anweisung getippt. Entschieden hat der Mensch. review-verdict und review-check stammen noch aus der Runde vom Vortag und muessen in dieser Lane neu erarbeitet werden - seither sind f6ce687 und 9fc224c dazugekommen, die den Forge-Zweig ueberhaupt erst eingefuehrt und dann vom falschen Remote auf den Push-Remote umgestellt haben.
- **2026-09-15 06:38 · Alexander Sacharov** — review-Runde 2026-09-15 zur GitLab-Aenderung (f6ce687 + 9fc224c), am Diff geprueft und nicht am Bericht. Verdict: angenommen, mit einer benannten Luecke.

Die Luecke, damit sie nicht verlorengeht: SKILL.md:55-63, Zweig 4 der Forge-Leiter. Host ist weder github.com noch gitlab-haltig und jaira.forge ist ungesetzt - die Rolle sagt korrekt, dass sie es nicht entscheiden kann, aber sie sagt nicht, ob der Agent danach anhaelt. Liest er weiter, steht er bei :71 vor einem unbedingten Push und bei :74-84 vor einer Gabel gh/glab, die er per Voraussetzung nicht aufloesen kann. Ein Satz an :63 ('Stop here and report; the rest of this file needs a settled forge') schliesst das. Nicht von mir eingebaut - review ist eine Urteils-Lane, keine Aenderungs-Lane, und die DoD verlangt es nicht.

Was ich ausdruecklich verifiziert und nicht angenommen habe, weil beides das Ticket unbemerkt aushebeln wuerde: (1) die nie-ausfuehren-Regel steht vollstaendig auf BEIDEN Wegen - :143-144 fuer gh, :145-146 fuer glab, dazu :86-87 und :113-115. (2) Der Forge-Host kommt aus 'git remote get-url origin', also dem Remote des Pushes, nicht aus jaira.remote; das ist genau das, was 9fc224c repariert hat, und die Begruendung steht im Prompt selbst bei :46-49. Auf diesem Board waere beides GitHub gewesen und der Fehler unsichtbar geblieben.

Der Hinweis der letzten in-progress-Runde stimmte: das alte review-check verlangte genau einen Treffer fuer 'gh pr create' und war ueberholt. Es ist ersetzt - der neue Check zaehlt keine Treffer mehr, sondern verlangt zwei never-run-Zeilen und stellt die Forge-Erkennung auf zwei Wegwerf-Repos nach. Jeder Schritt darin ist einmal gelaufen, bevor er eingetragen wurde.

go test ./... -race: RC=0, kein FAIL. Die glab-Zeilen kommen aus einer frisch gebauten Binary per 'roles install --into' wirklich heraus.

Nicht angefasst: spawn.sh, jaira-dispatcher/SKILL.md, jaira-teamlead/SKILL.md - die gehoeren KSGSKK im selben Worktree.
- **2026-09-15 06:45 · Alexander Sacharov** — Letzte Luecke, hier zu schliessen (Alex entschieden, kein Folgeticket): core/role/builtin/jaira-role-pr/SKILL.md, Sprosse 4 der Forge-Leiter (:55-63). Sie sagt heute nur "nicht raten" und nennt das Kommando, haelt den Agenten aber nicht an: die Ausfuehrung laeuft weiter in den unbedingten Push (:71) und in die gh/glab-Verzweigung (:74-84), die sie nicht aufloesen kann. Genau ein Satz schliesst das, im Ton des umgebenden Prompts: er muss unmissverstaendlich sagen, dass die Rolle hier stehenbleibt (kein Push, keine PR-Abfrage) und was sie dem Menschen sagt - naemlich dass er "git config jaira.forge gitlab|github" setzt und die Rolle danach erneut startet. Sprosse 1 ist genau dieses git config, die Anweisung an den Menschen ist also eine echte. Nur diese Datei plus eine Zeile in core/release/NOTES.md unter ## Unreleased (Zeilen dort nicht umbrechen, fremde Zeilen nicht anfassen). NICHT anfassen: scripts/spawn.sh, jaira-dispatcher/SKILL.md, jaira-teamlead/SKILL.md - daran arbeitet parallel Ticket KSGSKK im selben Worktree. Worktree und Branch stehen schon (feat/13VMA8-pr-is-the-humans); keinen neuen anlegen. Keinen Pull Request oeffnen oder aktualisieren - Alex gibt das Kommando.
- **2026-09-15 06:49 · Alexander Sacharov** — in-progress-Runde 2026-09-15 (letzte Luecke aus review-gaps): Sprosse 4 der Forge-Leiter haelt jetzt an.
- Der Satz sagt nicht nur 'stop', sondern woran er anhaelt ('here, before the push') und warum ('everything below needs the tool you were just unable to name'). Ein blankes 'Stop here' war genau der Fehler einer frueheren Runde dieses Tickets - der Agent stand dann vor Sektionen, die noch Arbeit schuldeten. Deshalb steht die Grenze ausdruecklich vor dem Push und nicht 'am Ende der Datei'.
- Die Anweisung an den Menschen verweist zurueck auf Sprosse 1 ('it is rung 1 above'), damit der Leser sieht, dass die Einstellung beim naechsten Lauf wirklich greift und der Neustart nicht ins selbe Loch faellt.
- NOTES.md: eine neue Zeile statt die bestehende GitLab-Zeile zu aendern. Die alte bleibt wahr ('it names no tool and tells you to set jaira.forge'); das Anhalten ist eine zusaetzliche Verhaltensaenderung, die ein Leser separat merken muss. Zwei Zeilen zur selben Rolle im selben Unreleased-Block sind Absicht.
- Nachgeprueft statt geglaubt: 'roles install --into /tmp/rollen' aus einer frisch gebauten Binary traegt den neuen Absatz - der go:embed-Pfad ist also mit.
- **2026-09-15 06:50 · Alexander Sacharov** — Zwischenfall in diesem Worktree (2026-09-15): die parallel laufende KSGSKK-Sitzung hat mit einem weiten 'git add' meine beiden Dateien mitgenommen. Commit 79bb6de ('chore(KSGSKK): review passes the dispatcher spawn round ...') traegt core/role/builtin/jaira-role-pr/SKILL.md (+6) und core/release/NOTES.md (+1) - die Aenderung dieser Runde - unter einer fremden Ticket-Id.
- Nicht repariert, bewusst: die Historie dieses Branches gehoert gerade zwei Sitzungen gleichzeitig; ein rebase oder amend haette der anderen unter den Haenden weggezogen. Der Inhalt ist richtig und liegt am Baum.
- Folge fuer die Commit-Ableitung: dieser Commit nennt 13VMA8 nicht. Die Ticket-Datei selbst ist in 79bb6de aber ebenfalls geaendert, und die Ableitung nimmt die Vereinigung aus Ticket-Datei-Historie UND Commits, die die Id nennen - 79bb6de faellt also ueber den ersten Weg trotzdem in die Liste.
- Fuer die naechste Runde in einem geteilten Worktree: 'git add <pfade>' statt 'git add -A'. Der Rollen-Prompt jaira-role-lane sagt heute 'git add -A' und ist damit genau die Anleitung, die diesen Zusammenstoss erzeugt hat.
- **2026-09-15 06:51 · Alexander Sacharov** — Fuer den erneuten Testing-Durchgang: was frueher schon festgestellt wurde und nicht verloren gehen darf - (1) kein glab-Flag ist per Analogie zu gh geraten, --body-file existiert dort nicht, geprueft gegen echtes glab 1.114.0; (2) die Never-Run-Regel (kein merge, kein approve) steht auf beiden Pfaden, gh wie glab; (3) die Forge wird am Push-Remote abgelesen, nicht an jaira.remote. Neu zu pruefen ist nur Sprosse 4 der Leiter in core/role/builtin/jaira-role-pr/SKILL.md: haelt sie den Agenten wirklich vor dem Push an, und sagt sie dem Menschen etwas Ausfuehrbares (git config jaira.forge gitlab|github, danach Rolle neu starten). Dazu die neue Zeile unter ## Unreleased in core/release/NOTES.md - eine Zeile, nicht umgebrochen. NICHT anfassen: scripts/spawn.sh, jaira-dispatcher/SKILL.md, jaira-teamlead/SKILL.md (paralleles Ticket KSGSKK im selben Worktree). Keinen Pull Request oeffnen oder aktualisieren.
- **2026-09-15 06:54 · Alexander Sacharov** — testing 2026-09-15: Gates gruen, Funktion gruen, ABER drei von vier DoD-Proofs zeigen auf falsche Zeilen - fail.

Was gruen ist, damit es niemand zweimal prueft:
- go build ./... RC=0; go test ./... -race RC=0, kein FAIL (core/lane 5.055s, internal/wintrap 1.514s, Rest cached).
- DoD 1 inhaltlich erfuellt: CLAUDE.md:156-169 steht hinter dem jaira:local-Marker (:154) vor :170 jaira:end, AGENTS.md:166-179, README.md:842-851 - alle drei tragen die Regel, die beiden Agenten-Dateien wortgleich.
- Forge-Leiter auf sechs Fixtures durchgespielt (Remotes github.com / gitlab.com / git.esprit-engineering.de, je ohne und mit gesetztem jaira.forge): Sprosse 1 gewinnt immer, github.com -> gh, gitlab.com -> glab, git.esprit-engineering.de ohne Einstellung -> kein Werkzeug genannt und Halt vor dem Push. Genau wie SKILL.md:51-69 es beschreibt.
- glab 1.114.0 gegengeprueft: 'mr list -s/--source-branch' existiert, 'mr create -t/--title' und '-d/--description' existieren, '--body-file' gibt es NICHT (die $(cat ...)-Form ist also richtig), 'glab mr approve' existiert.
- Frisch gebaute Binary: 'jaira roles install --into <dir>' schreibt die glab-Zeilen wirklich heraus (go:embed traegt sie mit) - im ausgelieferten SKILL.md stehen :89 glab mr list, :130 glab mr create, :151-152 die glab-Verbote.
- 'merge request' steht genau einmal (:71), also nur auf der GitLab-Seite und nicht global.

Der Befund: 9fc224c ('read the forge off the remote the branch is pushed to') hat den Absatz 'Then stop, here, before the push' (:65-69) eingezogen und alles darunter um sechs Zeilen verschoben. In 8c2324f wurde nur der Proof von DoD 3 nachgezogen, die drei aelteren nicht. Sie zeigen jetzt auf Text, der ihre Aussage nicht traegt:
- DoD 1 proof sagt 'SKILL.md:86-92 verzweigt nach dem Push in beide Betriebsarten' - :86-90 ist der glab-mr-list-Block. Richtig ist :92-97.
- DoD 2 proof sagt ':83 (glab mr list --source-branch) und :124 (glab mr create)' - :83 ist 'gh pr list', :124 ist 'gh pr create'. Der Proof fuer den GitLab-Weg zeigt auf die GitHub-Befehle, liest sich also als sein eigenes Gegenteil. Richtig ist :89 und :130.
- DoD 4 proof sagt ':143-148' fuer die never-run-Zeilen - :143-148 sind die Bullets zum Beantworten von Review-Kommentaren. Richtig ist :149-152. Und ':113-115' fuer 'You write it; you never run it' ist in Wahrheit :119-121.
Nur DoD 3 (:36-68, :65-68) stimmt.

Warum das nicht kosmetisch ist: die Proofs sind genau das, was jaira-role-pr in die PR-Beschreibung unter 'How to check' schreibt. Ein Reviewer, der DoD 2 folgt, landet auf 'gh pr create' und liest es als Beleg fuer GitLab-Unterstuetzung.

Fix, fuenf Minuten: 'jaira dod 13VMA8 1|2|4 --done --proof "..."' mit den oben genannten richtigen Zeilen. Inhaltlich ist nichts zu aendern - nur die Anker. Danach lohnt ein Blick, ob review-check Schritt 4 ('genau EIN Treffer fuer gh pr create') noch stimmt; die vorige Runde hat ihn selbst schon als ueberholt notiert.
