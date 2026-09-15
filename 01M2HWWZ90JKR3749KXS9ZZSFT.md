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
updated-at: 2026-09-15T13:11:38Z
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
outcome-what: "Plan geschrieben: zehn Schritte, die die Commit-Regel an ihrer Quelle aendern - core/board/announce.go und core/role/builtin/jaira-role-lane/SKILL.md - statt ein Gate zu bauen."
outcome-why: "Die Buchhaltungs-Commits entstehen aus einer Instruktion, die bedingungslos formuliert ist ('git add -A und commit'). Eine Lane ohne Code-Aenderung kann sie nur erfuellen, indem sie allein committet. Der Hebel ist der Text, nicht der Code - jaira committet nie selbst."
outcome-resolves: "Wie die Aenderung gemacht wird, steht fest; die offene Entscheidung ueber die Ticket-Datei nach der letzten Code-Lane ist als Plan-Schritt 2 benannt und in einer Notiz mit Empfehlung (a: der Ref traegt sie) hinterlegt."
---

# Zwei von drei Commits aendern nur eine Ticket-Datei

## Definition of Done

- [ ] Eine Lane, die keinen Code aendert, erzeugt keinen eigenen Commit mehr. Nachgestellt an einem Ticket, das critique, testing und review durchlaeuft: danach steht im Verlauf kein Commit, der nur .jaira/ anfasst.
- [ ] Die Commit-Liste eines Tickets bleibt vollstaendig, obwohl die Ticket-Datei seltener committet wird. Nachgestellt an einem Ticket, das die Lanes durchlaeuft und danach 'jaira move' in die Endlane erreicht - die abgeleitete Liste nennt jeden Code-Commit, der zu ihm gehoert.
- [ ] Die Regel steht dort, wo ein Agent sie liest: im erzeugten jaira-Block und in den Rollen-Prompts, nicht nur in einem Ticket.
- [ ] Eine Zeile in core/release/NOTES.md unter ## Unreleased, falls sich etwas an der Ableitung oder am Verhalten der Befehle aendert.

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] sammeln: jede Stelle, die einem Agenten sagt, die Ticket-Datei zu committen (core/board/announce.go, core/role/builtin/*/SKILL.md, docs/AGENTS.md, README.md, .claude/skills/jaira/SKILL.md)
- [x] entscheiden und als Notiz festhalten: was mit der Ticket-Datei nach der letzten Code-Lane passiert - Ref traegt sie, oder ein Abschluss-Commit
- [x] pruefen: laesst das Gate den Zug in die Endlane noch zu, wenn die Ticket-Datei nie committet wurde und nur die Commit-Nachricht die Id nennt (core/gate/gate.go:322, core/gitrepo/derive.go:19)
- [x] core/board/announce.go: den Commit-Punkt umschreiben - eine Lane ohne Code-Aenderung committet nichts, der Ref haelt den Zustand, jeder Commit nennt die Id
- [ ] core/role/builtin/jaira-role-lane/SKILL.md: 'git add -A und commit' durch dieselbe Bedingung ersetzen
- [ ] core/role/builtin/jaira-role-pr/SKILL.md: die Vorbedingung 'nothing uncommitted' auf alles ausser der Ticket-Datei verengen
- [ ] die nicht erzeugten Kopien nachziehen: docs/AGENTS.md, README.md, .claude/skills/jaira/SKILL.md
- [ ] 'jaira update' laufen lassen, damit CLAUDE.md und AGENTS.md den neuen Block tragen; core/board/announce_test.go und core/role/role_test.go gruen
- [ ] eine Zeile in core/release/NOTES.md unter ## Unreleased
- [ ] am eigenen Ticket nachstellen: nach critique, testing und review zeigt 'git log --name-only origin/master..HEAD' keinen Commit, der nur .jaira/ anfasst

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
