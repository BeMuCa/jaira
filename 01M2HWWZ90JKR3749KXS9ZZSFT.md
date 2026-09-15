---
id: 01M2HWWZ90JKR3749KXS9ZZSFT
title: Zwei von drei Commits aendern nur eine Ticket-Datei
status: critique
ready: true
creator: Alexander Sacharov
goal: "Ein Zweig zeigt die Arbeit, nicht die Buchhaltung: wer den Verlauf liest, sieht Aenderungen am Werkzeug und nicht jede Lane, die einen Vermerk hinterlassen hat."
definition-of-done: "Eine Lane, die keinen Code aendert, erzeugt keinen eigenen Commit mehr. Nachgestellt an einem Ticket, das critique, testing und review durchlaeuft: danach steht im Verlauf kein Commit, der nur .jaira/ anfasst."
tags:
  - cli
  - docs
blocked-by: []
related: []
commits:
  - 5164191ae41d9168398545a5d5915974f85ca343
created-at: 2026-09-15T06:43:33Z
updated-at: 2026-09-15T13:32:35Z
assignee: "Alexander Sacharov"
updated-by: Alexander Sacharov
context: |-
  Gemessen am 2026-09-15 auf dem Zweig feat/9ET6NC-per-board-remote: 92 Commits, davon 30 mit Code. Zweiundsechzig fassen nur .jaira/ an - eine Lane, die ihren Vermerk hinterlaesst, ein Uebergang, eine Notiz.

  Das widerspricht der Regel, die das Projekt selbst aufgeschrieben hat. In CLAUDE.md steht, das Ticket faehrt im SELBEN Commit wie der Code, damit ein Leser die Aenderung und ihren Grund an einer Stelle sieht. Gemacht wird das Gegenteil: je Lane ein eigener Commit, der nichts als die Ticket-Datei bewegt.

  Warum es so kam: jede Lane ist ein eigener Worker, und critique, testing und review aendern gar keinen Code. Sie haben nichts, woran sie ihren Vermerk haengen koennten, also committen sie ihn allein.

  Der Punkt, der die Sache entscheidet: diese Commits sind nicht das, was den Zustand bewahrt. Record() (core/refsync/refsync.go:108) stellt JEDE Schreibung - note, move, dod - in die Outbox und schickt sie auf den Ref, unabhaengig von git. Die Commits dienen der Lesbarkeit des Pull Requests, nicht der Haltbarkeit. Wer sie weglaesst, verliert nichts Dauerhaftes.

  Was dabei nicht kaputtgehen darf: jaira leitet die Commit-Liste eines Tickets aus der Vereinigung zweier Quellen ab - der Historie der Ticket-Datei UND der Commits, die seine Id nennen. Faellt die erste Quelle weg, haengt alles daran, dass Commits die Id im Betreff tragen. Heute tun sie das ohnehin ('fix(KSGSKK): ...'), aber aus einer Gewohnheit wird damit eine Bedingung.
claimed-by: DESKTOP-RFTCH11-41016
claimed-at: 2026-09-15T13:03:43Z
outcome-what: "Die Formulierung 'the `jaira logbook <id>` commit' ist in allen sieben Kopien durch 'the commit that files the ticket away with `jaira logbook <id>`' ersetzt: core/board/announce.go:100 (und darueber erzeugt AGENTS.md:71 und CLAUDE.md:98), docs/AGENTS.md:76, .claude/skills/jaira/SKILL.md:281, core/role/builtin/jaira-role-lane/SKILL.md:36, README.md:848, core/release/NOTES.md:17 sowie die handgeschriebenen PR-Abschnitte AGENTS.md:186 und CLAUDE.md:176, die 'jaira update' nicht anfasst. Die Absaetze in docs/AGENTS.md und .claude/skills/jaira/SKILL.md, in die der logbook-Halbsatz im vorigen Durchgang eingefuegt wurde, laufen wieder durchgehend auf 80 Zeichen."
outcome-why: "Befund (2) des zweiten critique-Durchgangs: der alte Wortlaut macht jaira zum Urheber des Commits. Store.Logbook (core/ticket/store.go:328) verschiebt nur die Datei, und .claude/skills/jaira/SKILL.md:289 sagt sechs Zeilen tiefer 'jaira never commits for you' - der Widerspruch stand auf einem Bildschirm. Befund (1): der eingefuegte Halbsatz war nicht neu umbrochen, docs/AGENTS.md:77 lief auf 123 und .claude/skills/jaira/SKILL.md:282 auf 120 Zeichen mitten in einem 80er-Absatz; der naechste Diff dieses Absatzes haette jede Zeile neu gefaerbt."
outcome-resolves: "Beide Befunde des zweiten critique-Durchgangs sind erledigt; kein 'logbook ... commit' in diesem Wortlaut bleibt im Baum (grep leer). DoD 3 und 4 bleiben getickt und werden praeziser: die Regel lautet in allen Quellen gleich und behauptet nichts Falsches ueber das Werkzeug. DoD 1 und 2 bleiben der testing-Lane vorbehalten. go build ./... und go test ./... gruen; der Zweig traegt drei Commits, jeder mit Code und Handle im Betreff, keiner nur mit .jaira/."
review-summary: |-
  docs/AGENTS.md:77 und .claude/skills/jaira/SKILL.md:282: der eingefuegte Halbsatz wurde nicht neu umbrochen - mitten im auf 80 Zeichen umbrochenen Absatz steht jetzt eine Zeile mit 123 bzw. 120 Zeichen ("... carries its final state. Nothing is lost by waiting ..."). Den Absatz in beiden Dateien wieder auf die Breite der Nachbarzeilen umbrechen, sonst faerbt der naechste Diff dieses Absatzes jede Zeile neu.
  .claude/skills/jaira/SKILL.md:281 sagt "the 'jaira logbook <id>' commit", sechs Zeilen darunter sagt Zeile 289 "jaira never commits for you". Die beiden Saetze widersprechen sich in derselben Datei: jaira legt keinen Commit an, Store.Logbook (core/ticket/store.go:328) verschiebt nur die Datei. Stattdessen "the commit that files the ticket away with 'jaira logbook <id>'" schreiben - in dieser Datei zwingend, und derselben Formulierung wegen gleich mit in core/board/announce.go:100, docs/AGENTS.md:76, core/role/builtin/jaira-role-lane/SKILL.md:36, README.md:847 und core/release/NOTES.md:17.
---

# Zwei von drei Commits aendern nur eine Ticket-Datei

## Definition of Done

- [ ] Eine Lane, die keinen Code aendert, erzeugt keinen eigenen Commit mehr. Nachgestellt an einem Ticket, das critique, testing und review durchlaeuft: danach steht im Verlauf kein Commit, der nur .jaira/ anfasst.
- [ ] Die Commit-Liste eines Tickets bleibt vollstaendig, obwohl die Ticket-Datei seltener committet wird. Nachgestellt an einem Ticket, das die Lanes durchlaeuft und danach 'jaira move' in die Endlane erreicht - die abgeleitete Liste nennt jeden Code-Commit, der zu ihm gehoert.
- [x] Die Regel steht dort, wo ein Agent sie liest: im erzeugten jaira-Block und in den Rollen-Prompts, nicht nur in einem Ticket.
  proof: core/board/announce.go:92 (erzeugter Block, Punkt 'a lane that changed no code commits nothing'); core/role/builtin/jaira-role-lane/SKILL.md:32; core/role/builtin/jaira-role-pr/SKILL.md:29; docs/AGENTS.md:72; .claude/skills/jaira/SKILL.md:277; README.md:843
- [x] Eine Zeile in core/release/NOTES.md unter ## Unreleased, falls sich etwas an der Ableitung oder am Verhalten der Befehle aendert.
  proof: core/release/NOTES.md:16 unter ## Unreleased; core/release TestNotes gruen

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] sammeln: jede Stelle, die einem Agenten sagt, die Ticket-Datei zu committen (core/board/announce.go, core/role/builtin/*/SKILL.md, docs/AGENTS.md, README.md, .claude/skills/jaira/SKILL.md)
- [x] entscheiden und als Notiz festhalten: was mit der Ticket-Datei nach der letzten Code-Lane passiert - Ref traegt sie, oder ein Abschluss-Commit
- [x] pruefen: laesst das Gate den Zug in die Endlane noch zu, wenn die Ticket-Datei nie committet wurde und nur die Commit-Nachricht die Id nennt (core/gate/gate.go:322, core/gitrepo/derive.go:19)
- [x] core/board/announce.go: den Commit-Punkt umschreiben - eine Lane ohne Code-Aenderung committet nichts, der Ref haelt den Zustand, jeder Commit nennt die Id
- [x] core/role/builtin/jaira-role-lane/SKILL.md: 'git add -A und commit' durch dieselbe Bedingung ersetzen
- [x] core/role/builtin/jaira-role-pr/SKILL.md: die Vorbedingung 'nothing uncommitted' auf alles ausser der Ticket-Datei verengen
- [x] die nicht erzeugten Kopien nachziehen: docs/AGENTS.md, README.md, .claude/skills/jaira/SKILL.md
- [x] 'jaira update' laufen lassen, damit CLAUDE.md und AGENTS.md den neuen Block tragen; core/board/announce_test.go und core/role/role_test.go gruen
- [x] eine Zeile in core/release/NOTES.md unter ## Unreleased
- [~] am eigenen Ticket nachstellen: nach critique, testing und review zeigt 'git log --name-only origin/master..HEAD' keinen Commit, der nur .jaira/ anfasst

## Progress
- **2026-09-15 13:06 · Alexander Sacharov** — Der Plan ist reine Textarbeit, kein Code: die Regel, die die Buchhaltungs-Commits erzeugt, steht an genau einer Stelle als Quelle - core/role/builtin/jaira-role-lane/SKILL.md sagt heute bedingungslos 'move the ticket, then git add -A and commit'. Der erzeugte Block (core/board/announce.go:87) sagt dasselbe als 'the ticket rides in the same commit as the code'. Eine Lane, die keinen Code anfasst, kann die Regel nur erfuellen, indem sie allein committet - daher die 62 von 92.

Warum kein Code: nichts in jaira committet selbst. Der Zug wird von einem Worker ausgefuehrt, der einen Prompt liest. Ein Gate, das '.jaira/-only-Commit' verbietet, koennte man bauen, aber es wuerde nach der Tat greifen und die Lane in eine Sackgasse setzen, in der sie ihren Vermerk nirgends unterbringt. Die Instruktion ist der Hebel.

Die offene Entscheidung (Plan-Schritt 2): was mit der Ticket-Datei geschieht, nachdem die letzte Code-Lane committet hat. critique, testing und review laufen danach und schreiben Notizen. Zwei Wege:
(a) Der Ref traegt sie, git nie. Die Ticket-Datei bleibt im Arbeitsbaum geaendert; jaira-role-pr muss seine Vorbedingung 'nothing uncommitted' entsprechend verengen, sonst blockiert sie jeden Push. Preis: nach dem Merge steht auf master eine veraltete Ticket-Datei, bis irgendein spaeterer Commit sie mitnimmt. Der Zustand ist nicht verloren - refsync.Pull holt ihn vom Ref zurueck.
(b) Ein einziger Abschluss-Commit vor dem Push. Preis: das ist wieder ein Commit, der nur .jaira/ anfasst, und widerspricht der DoD woertlich.
Ich wuerde (a) nehmen, weil die DoD (b) ausschliesst und der Kontext genau dieses Argument macht: die Commits bewahren nichts, Record() tut es.

Was dabei kippen kann und Schritt 3 prueft: die abgeleitete Commit-Liste (core/gitrepo/derive.go:19) ist die Vereinigung aus Ticket-Datei-Historie und Commits, die die Id nennen. Faellt die erste Quelle fuer die spaeten Lanes weg, haengt alles an der zweiten. Das Nennen der Id im Betreff wird damit von einer Gewohnheit zur Bedingung - und das muss im erzeugten Block stehen, nicht nur hier.
- **2026-09-15 13:07 · Alexander Sacharov** — Befund beim Planen, der die Empfehlung aus der vorigen Notiz einschraenkt: refsync laeuft auf diesem Board NICHT. 'git for-each-ref refs/jaira/*' ist leer, .jaira/ hat keine Outbox, und 'git config jaira.remote' ist nicht gesetzt. Syncer.Usable() gibt dann einen Fehler und Record() (core/refsync/refsync.go:108) kehrt sofort zurueck, ohne etwas zu queuen.

Damit ist die Kernaussage des Kontexts - 'die Commits bewahren nichts, der Ref tut es' - nur auf einem Board mit konfiguriertem Remote wahr. Auf einem Board ohne Remote ist der Commit der einzige Traeger des Ticket-Zustands.

Folge fuer Plan-Schritt 2: Weg (a) - 'der Ref traegt sie, git nie' - darf nicht bedingungslos in den erzeugten Block geschrieben werden. Die Regel muss entweder an 'Board hat ein Remote' geknuepft werden, oder der Verzicht auf den Commit gilt nur fuer Lanes, deren Vermerk eine spaetere Code-Lane ohnehin mitnimmt.

Deshalb committet dieser pre-process-Schritt seine Ticket-Datei doch: sie war untracked, kein Ref haelt sie, und ohne Commit gaebe es auf dem Zweig nichts, was das Ticket mit der Arbeit verbindet - und keine ableitbare Commit-Liste.
- **2026-09-15 13:10 · Alexander Sacharov** — Entscheidung zu Plan-Schritt 2 (Weg (a), mit einer Verschaerfung): eine Lane ohne Code-Aenderung committet nichts. Die Ticket-Datei bleibt im Arbeitsbaum geaendert; der naechste Commit, der Code traegt, nimmt sie mit. Weg (b) - ein Abschluss-Commit vor dem Push - ist ausgeschlossen, weil er woertlich das ist, was die DoD verbietet.

Schritt 3 geprueft, und er entschaerft den Einwand aus der vorigen Notiz: das Gate laesst den Zug in die Endlane zu, auch wenn die Ticket-Datei nie committet wurde. core/gate/gate.go:322 verlangt nur, dass CommitsForTicket etwas findet, und core/gitrepo/derive.go:19 bildet die Vereinigung aus Datei-Historie UND Commits, die den Handle im Betreff nennen. Die zweite Quelle allein genuegt. Die Refusal-Meldung sagt das sogar selbst ('or name <handle> in the commit message'). Damit wird das Nennen der Id im Betreff von einer Gewohnheit zur Bedingung - und genau deshalb steht es jetzt als eigener Punkt im erzeugten Block, nicht als Nebensatz.

Der Preis, den ich bewusst nehme: Notizen, die nach dem letzten Code-Commit entstehen (critique/testing/review am Ende einer Schleife, die nichts mehr zurueckschickt), stehen nicht auf dem Zweig. Auf einem Board mit konfiguriertem Remote holt refsync.Pull sie zurueck; auf diesem Board - ohne Remote, siehe vorige Notiz - bleiben sie im Arbeitsbaum, bis ein spaeterer Commit sie mitnimmt. Das ist begrenzt: die Ticket-Datei ist ab dem ersten Code-Commit getrackt, veraltet sind also die letzten Notizen, nicht das Ticket. Wer das spaeter anders entscheidet, muss zuerst die DoD dieses Tickets aendern - sie schliesst den Abschluss-Commit aus, nicht ich.
- **2026-09-15 13:14 · Alexander Sacharov** — Der pre-process-Commit 4e89f86 ('docs(9ZZSFT): plan the removal of bookkeeping-only commits') fasst nur .jaira/ an und ist damit genau das, was die neue Regel verbietet - pre-process aendert keinen Code. Er war nicht gepusht und hatte kein Upstream, also habe ich ihn per 'git reset --soft HEAD~1' in diesen Commit gefaltet, statt ihn stehen zu lassen und die DoD-Demonstration (Plan-Schritt 10) daran scheitern zu lassen. Das ist zugleich die Probe aufs Exempel: die Ticket-Datei erscheint jetzt zum ersten Mal in dem Commit, der auch Code traegt.

Warum das kein Zustandsverlust war: 'git reset --soft' laesst den Arbeitsbaum unberuehrt, und der Ticket-Zustand liegt in der Datei, nicht im Commit. Wer das nachmacht, muss vorher 'git ls-remote --heads origin <branch>' pruefen - auf einem gepushten Zweig ist das Falten verboten und der Buchhaltungs-Commit bleibt stehen.
- **2026-09-15 13:16 · Alexander Sacharov** — Commit dieser Lane: 239c881 'docs(9ZZSFT): stop a lane without a code change from committing' - Code und Ticket-Datei zusammen, Handle im Betreff. 'git log --name-only master..HEAD' zeigt danach genau diesen einen Commit und keinen, der nur .jaira/ anfasst; der alte pre-process-Commit ist darin aufgegangen.

Fuer critique/testing/review: ihr seid selbst die Nachstellung von DoD 1 und 2. Committet nichts. Wenn ihr am Ende in 'git log --name-only master..HEAD' immer noch nur 239c881 seht, ist DoD 1 bewiesen; wenn 'jaira move' in die Endlane durchgeht und die abgeleitete Liste 239c881 nennt, ist DoD 2 bewiesen. Beide sind absichtlich ungetickt geblieben - sie lassen sich in dieser Lane nicht beweisen, nur vorbereiten.
- **2026-09-15 13:19 · Alexander Sacharov** — critique: vier Befunde, alle mit klarer Korrektur - zurueck nach in-progress.

1. Die Regel steht nicht in allen sechs Quellen gleich. CLAUDE.md:168 und AGENTS.md:178 tragen den Abschnitt 'Work rides on a branch' hinter dem jaira:local-Marker; README.md:843 hat dort den neuen Zusatz bekommen, die beiden anderen nicht. Das outcome-what behauptet 'ueberall gleich' - drei Kopien desselben Absatzes, zwei davon alt.
2. core/role/builtin/jaira-role-pr/SKILL.md:28 fuehrt Punkt 2 mit der abgeloesten Fettzeile an. Ein Agent, der nur die Fettzeilen liest - wofuer sie da sind -, liest die alte Regel; die Einschraenkung kommt erst im Rumpf.
3. 'see the next point' (core/board/announce.go:85) zeigt auf den falschen Punkt. Duenn wird die Datei-Historie durch den dritten Punkt, nicht den zweiten.
4. Der wesentliche: 'the next commit that carries code takes it along' hat fuer die letzten Lanes keinen Adressaten. Findet critique nichts, laufen testing und review ohne weitere Code-Aenderung - dann gibt es keinen naechsten Commit, und der PR zeigt eine Ticket-Datei im in-progress-Stand, ohne review-summary und test-verdict. Genau der Zustand, den die Regel 'der Leser sieht die Aenderung und ihren Grund an einer Stelle' verhindern soll.

Das ist ausdruecklich KEIN Ruf nach dem Abschluss-Commit, den die DoD verbietet. Der Traeger existiert schon: 'jaira logbook <id>' verschiebt die Ticket-Datei nach .jaira/logbook/ und nimmt ihren Endstand in den Commit mit, der diese Verschiebung traegt. Die Regel muss ihn nur benennen, sonst liest sie sich als 'der Endstand bleibt liegen'. Die Note vom 13:10 nimmt den Preis bewusst in Kauf und verweist auf refsync.Pull - die Note vom 13:07 zeigt, dass refsync auf diesem Board nicht laeuft. Damit bleibt logbook der einzige Traeger, und er gehoert in den Text.
- **2026-09-15 13:20 · Alexander Sacharov** — in-progress (zweiter Durchgang): die vier critique-Befunde werden abgearbeitet. Zu Befund 3 weiche ich vom Wortlaut der Korrektur ab: 'see the last point in this list' waere falsch, die Liste endet mit 'jaira logbook'. Ich schreibe stattdessen einen selbstbeschreibenden Verweis auf den Punkt ueber Lanes ohne Code-Aenderung - der bleibt richtig, auch wenn die Liste spaeter umsortiert wird.
- **2026-09-15 13:23 · Alexander Sacharov** — Alle vier critique-Befunde behoben. Was dabei zusaetzlich auffiel und mit erledigt wurde: die NOTES.md-Zeile unter ## Unreleased sagte ebenfalls nur 'der naechste Code-Commit nimmt sie mit' - sie ist dieselbe Aussage fuer den Leser des Releases und haette die Luecke aus Befund 4 nach aussen getragen; der logbook-Halbsatz steht jetzt auch dort. Nicht angefasst, weil ausserhalb dieser Lane: docs/AGENTS.md:85 spricht noch von 'jaira archive <id>', den Befehl gibt es nicht mehr - er heisst logbook. Das ist ein eigenes Ticket wert.
- **2026-09-15 13:27 · Alexander Sacharov** — critique (2. Durchgang): die vier Befunde des ersten Durchgangs sind alle abgearbeitet und werden nicht erneut aufgemacht. Zwei neue, beide aus dem Nachtrag-Commit 9e4d870 selbst: (1) der Halbsatz wurde eingefuegt, ohne den Absatz neu umzubrechen - docs/AGENTS.md:77 (123 Zeichen) und .claude/skills/jaira/SKILL.md:282 (120 Zeichen) stehen in Absaetzen, die sonst auf 80 umbrechen; AGENTS.md, CLAUDE.md, announce.go und README.md sind korrekt umbrochen, betroffen sind nur diese zwei. (2) "the 'jaira logbook <id>' commit" liest sich, als lege jaira den Commit an; Store.Logbook (core/ticket/store.go:328) verschiebt nur die Datei, und .claude/skills/jaira/SKILL.md sagt sechs Zeilen weiter unten ausdruecklich "jaira never commits for you" - dort steht der Widerspruch auf einem Bildschirm.
Nicht aufgemacht und bewusst stehen gelassen: dass dieselbe Regel in sechs handgepflegten Kopien steht. Genau diese Vervielfachung hat die Befunde beider Durchgaenge erzeugt, aber eine billigere Form gibt es nicht - README und docs richten sich an Menschen, die Rollen-SKILLs sind Prompts, nur CLAUDE.md/AGENTS.md werden aus announce.go erzeugt. Das waere ein eigenes Ticket, keine Aenderung an diesem.
- **2026-09-15 13:31 · Alexander Sacharov** — in-progress (3. Durchgang): beide critique-Befunde des 2. Durchgangs behoben. Die Formulierung lautet jetzt ueberall 'the commit that files the ticket away with `jaira logbook <id>`' - das Subjekt ist der Commit, nicht jaira, und der Widerspruch zu '`jaira` never commits for you' sechs Zeilen tiefer in .claude/skills/jaira/SKILL.md ist weg. Angefasst wurden sieben Stellen statt der sechs aus dem Befund: AGENTS.md:187 und CLAUDE.md:176 tragen denselben Satz im handgeschriebenen PR-Abschnitt, der NICHT aus announce.go erzeugt wird - 'jaira update' laesst ihn stehen. Wer die Formulierung spaeter noch einmal aendert, muss diese beiden von Hand mitnehmen.

Beim Umbrechen: es genuegt nicht, die eine geaenderte Zeile zu kuerzen - der Umbruch schiebt sich durch den Rest des Absatzes, und eine kurze Zeile in der Mitte ist derselbe Diff-Schaden wie eine zu lange. In docs/AGENTS.md und .claude/skills/jaira/SKILL.md waren jeweils drei Nachlaeufe noetig, bis der Absatz wieder durchgehend auf 80 lief.
- **2026-09-15 13:32 · Alexander Sacharov** — Falle beim Ablauf 'commit - move --commits $(git rev-parse HEAD) - amend, damit die Ticket-Datei mitfaehrt': das Amend vergibt eine neue Sha, und die im Ticket steht dann auf einem Commit, den es nicht mehr gibt. Korrigiert mit 'jaira set 9ZZSFT commits=<neue sha>'. Die Reihenfolge, die das vermeidet: erst 'git add' von Code UND Ticket-Datei nach dem Move, dann ein einziger Commit, und die Sha danach per 'jaira set' nachtragen - sie laesst sich vorher nicht kennen. Die Ticket-Datei bleibt dadurch immer um genau diesen einen 'set'-Schreibvorgang hinter dem Commit zurueck; sie liegt jetzt geaendert im Arbeitsbaum und faehrt mit dem naechsten Code-Commit, wie die Regel es vorsieht. Die abgeleitete Liste ist davon unberuehrt - sie liest den Handle aus dem Betreff.
