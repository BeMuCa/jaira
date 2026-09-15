---
id: 01M2HWWZ90JKR3749KXS9ZZSFT
title: Zwei von drei Commits aendern nur eine Ticket-Datei
status: in-progress
ready: true
creator: Alexander Sacharov
goal: "Ein Zweig zeigt die Arbeit, nicht die Buchhaltung: wer den Verlauf liest, sieht Aenderungen am Werkzeug und nicht jede Lane, die einen Vermerk hinterlassen hat."
definition-of-done: "Eine Lane, die keinen Code aendert, erzeugt keinen eigenen Commit mehr. Nachgestellt an einem Ticket, das critique, testing und review durchlaeuft: danach steht im Verlauf kein Commit, der nur .jaira/ anfasst."
tags:
  - cli
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-09-15T06:43:33Z
updated-at: 2026-09-15T13:20:27Z
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
outcome-what: "Die Commit-Regel in allen sechs Quellen umgeschrieben, die sie einem Agenten sagen. Neu und ueberall gleich: eine Lane, die keinen Code geaendert hat, committet gar nichts und laesst die Ticket-Datei im Arbeitsbaum stehen; der naechste Commit, der Code traegt, nimmt sie mit. Als eigener Punkt daneben: jeder Commit nennt den Handle im Betreff - das ist jetzt die Bedingung fuer die abgeleitete Commit-Liste, nicht mehr eine Gewohnheit. Betroffen: core/board/announce.go (erzeugter Block, dadurch CLAUDE.md und AGENTS.md nach 'jaira update'), core/role/builtin/jaira-role-lane/SKILL.md (das unbedingte 'git add -A und commit' ist weg, ersetzt durch die Fallunterscheidung und durch 'adde mit Pfad, nie -A'), core/role/builtin/jaira-role-pr/SKILL.md (die Vorbedingung heisst jetzt 'nothing uncommitted but the ticket file', und Punkt 2 sagt ausdruecklich, dass eine geaenderte Ticket-Datei beim Push kein Fehler ist), docs/AGENTS.md, README.md, .claude/skills/jaira/SKILL.md. Eine Zeile in core/release/NOTES.md unter ## Unreleased. Kein Code-Verhalten geaendert - jaira committet nie selbst, der Hebel ist die Instruktion."
outcome-why: "Auf feat/9ET6NC-per-board-remote fassten 62 von 92 Commits nur .jaira/ an: je Lane ein Vermerk, ein Uebergang, eine Notiz. Das kam aus einer Instruktion, die bedingungslos formuliert war - 'move the ticket, then git add -A and commit'. Eine Lane wie critique, testing oder review aendert keinen Code und konnte die Regel nur erfuellen, indem sie allein committete. Wer den Zweig liest, sieht dann die Buchhaltung und nicht die Arbeit."
outcome-resolves: "DoD 3 (die Regel steht dort, wo ein Agent sie liest) ist erfuellt: sie steht im erzeugten Block (core/board/announce.go:92) und in beiden Rollen-Prompts, nicht nur im Ticket. DoD 4 ist erfuellt: core/release/NOTES.md:16. DoD 1 und 2 sind bewusst noch offen - sie verlangen die Nachstellung an einem Ticket, das critique, testing und review DURCHLAUFEN hat, und diese Lanes kommen erst noch; sie sind von der testing-Lane zu ticken. Vorbereitet ist beides: der pre-process-Commit 4e89f86, der nur .jaira/ anfasste, ist in diesen Commit gefaltet, so dass der Zweig jetzt genau einen Commit traegt, der Code und Ticket-Datei zusammen fuehrt und den Handle im Betreff nennt. Dass die Ableitung ohne Datei-Historie traegt, ist an der Quelle geprueft: core/gate/gate.go:322 verlangt nur ein nichtleeres Ergebnis von core/gitrepo/derive.go:19, und dessen zweite Quelle - Commits, die den Handle nennen - genuegt allein. go build ./... und go test ./... gruen."
review-summary: |-
  CLAUDE.md:168 und AGENTS.md:178 (Abschnitt 'Work rides on a branch', hinter dem jaira:local-Marker) sagen weiter nur 'the ticket rides in the same commits as the code'; README.md:843 hat den neuen Zusatz bekommen, diese beiden nicht - denselben Satz ('It rides with the code and never alone: a lane that changed no code ... commits nothing at all') dort ergaenzen.
  core/role/builtin/jaira-role-pr/SKILL.md:28 fuehrt Punkt 2 weiter mit der unbedingten Fassung an ('**The ticket rides in the same commits as the code.**'), die jede andere Quelle abgelegt hat; der Rumpf relativiert sie erst danach - die Fettzeile auf '**The ticket rides with the code, never on its own.**' aendern.
  core/board/announce.go:85 (und die erzeugten Kopien CLAUDE.md:82, AGENTS.md:55) verweisen mit 'see the next point' auf den Punkt 'the ticket rides with the code' - duenn macht die Datei-Historie aber erst der Punkt danach; 'see the last point in this list' schreiben.
  core/board/announce.go:98 (gleichlautend docs/AGENTS.md:75, .claude/skills/jaira/SKILL.md:280) sagt 'the next commit that carries code takes it along' und hat fuer die letzten Lanes keine Antwort: laeuft critique/testing/review nach dem letzten Code-Commit durch, gibt es keinen naechsten - der Zweig geht mit einer Ticket-Datei im in-progress-Stand in den PR. Den Schritt benennen, der sie doch traegt: der 'jaira logbook <id>'-Commit verschiebt die Datei und nimmt ihren Endstand mit - als Halbsatz an 'takes it along' anhaengen.
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
