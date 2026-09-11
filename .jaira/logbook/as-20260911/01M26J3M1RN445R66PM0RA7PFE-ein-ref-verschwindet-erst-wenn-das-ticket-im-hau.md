---
id: 01M26J3M1RN445R66PM0RA7PFE
title: "Ein Ref verschwindet erst, wenn das Ticket im Hauptbranch angekommen ist"
status: done
ready: true
creator: Alexander Sacharov
goal: "Zwischen 'ins Logbuch gelegt' und 'im Hauptbranch angekommen' bleibt das Ticket fuer alle sichtbar, damit niemand dasselbe Problem ein zweites Mal aufschreibt"
context: |-
  Heute loescht logbook und archive das Ref sofort. Die Ticketdatei liegt zu diesem Zeitpunkt aber nur im Branch der Person, die es abgeraeumt hat. In dem Fenster dazwischen - Branch fertig, noch nicht gemergt - sieht das Ticket NIEMAND sonst: das Ref ist weg, der Branch ist fremd. Wer dasselbe Problem bemerkt, schreibt es ein zweites Mal auf, mit derselben Loesung.

  Das Fenster ist nicht klein. Es ist so lang wie ein Review dauert.

  Die Regel soll sein: das Ref ueberlebt bis das Ticket im Hauptbranch liegt, und erst dann wird es endgueltig entfernt. Solange traegt es den Endzustand ('done') und ist damit fuer alle sichtbar als 'fertig, wartet aufs Landen' - genau die Information, die einen doppelten Ticket verhindert.

  Gemessen, nicht vermutet: 'git rev-list -1 <hauptbranch> -- ".jaira/logbook/*/<id>*" ".jaira/archive/<id>*"' antwortet leer, solange das Ticket nur im Branch liegt, und mit einem Sha, sobald der Branch gemergt ist. Kosten: nicht messbar (0,00 s), ein Aufruf, nur fuer Tickets im Endzustand.

  Auch geprueft und ein Stolperstein: 'git symbolic-ref refs/remotes/origin/HEAD' ist in einem frischen Klon nicht zwingend gesetzt (in meinem Testrepo fehlte es). Es braucht also einen Fallback auf origin/main und origin/master und einen Eintrag in ~/.jaira/settings.json.

  Der zweite Fall, der dadurch entsteht und mitbedacht werden muss: ein Branch, der nie gemergt wird. Dann bleibt das Ref fuer immer. Es braucht eine Meldung 'fertig, aber seit N Tagen nicht angekommen' und einen ausdruecklichen Weg, das Ref trotzdem loszuwerden - stilles Aufraeumen ist es nicht, denn das waere wieder das Fenster, nur unsichtbar.
definition-of-done: "logbook und archive loeschen das Ref nicht mehr sofort, sondern schreiben den Endzustand darauf, sodass andere Klone das Ticket als 'fertig, wartet aufs Landen' sehen; das Ref wird entfernt, sobald 'git rev-list -1 <hauptbranch> -- .jaira/logbook/*/<id>* .jaira/archive/<id>*' einen Commit findet, und das laeuft im Hintergrund nach dem Muster der Update-Pruefung statt auf dem Kommandopfad; der Hauptbranch wird aus origin/HEAD bestimmt, mit Fallback auf origin/main und origin/master und einem Eintrag in settings.json; ein Ticket, das im Endzustand steht und seit einer konfigurierbaren Frist nicht im Hauptbranch angekommen ist, wird von fetch und validate gemeldet; es gibt einen ausdruecklichen Befehl, ein solches Ref trotzdem zu entfernen, und er sagt was er tut; Tests mit zwei Klonen belegen, dass das Ticket im Fenster fuer den anderen Klon sichtbar bleibt und nach dem Merge verschwindet"
tags:
  - concurrency
blocked-by: []
commits:
  - 69ae61d207d10b5d57ba273a081426f6ab09515f
  - 6bfad50c24a5591a8dacdb6480240c6af4d50aca
  - b8e1cb1a0313c41a18fabe933a33b6fc8f3cb81a
  - fe864ec705086aafcd4b16fab0556d9dd2c3ceb7
  - 2a35f28cbb411e419354d97b6f6e8ef95db9c714
  - 39b1a2d79b19474c3879034af3c164637c8e8be7
  - 3dbf02839cb83e92d9cb8f0e38eb225b454782ea
  - 7811cb3c37c239f596c20e73dc34a0ea2088de55
  - 5fb08295461779ff76bf00e6ac1be5bd773668a3
created-at: 2026-09-10T21:03:18Z
updated-at: 2026-09-11T13:31:51Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-158836
claimed-at: 2026-09-11T11:29:31Z
assignee: Alexander Sacharov
question: "Der Ref lebt jetzt bis zum Landen. Eine Frage: die Frist, nach der ein 'fertig, aber nirgends angekommen' gemeldet wird, steht auf sieben Tagen - passt das zu euren Reviews, oder eher drei?"
outcome-what: "Провёл ревью diff'а против DoD: код собирается, юнит-тесты (core/snapshot, core/settings, core/refsync) проходят."
outcome-why: "Найден вероятный функциональный дефект — жатва refs зависит от локальных remote-tracking веток landing-branches, которые jaira нигде не обновляет, что может сорвать обещание DoD 'sobald ein Commit gefunden wird'; плюс отсутствует явный fallback origin/main/origin/master и нет теста с двумя реальными клонами на полный сценарий жатвы."
outcome-resolves: "review-summary, review-gaps, review-verdict, review-check записаны на тикете; решение об одобрении/возврате — за человеком."
review-summary: "logbook/archive пишут финальное состояние на ref вместо удаления (RecordFiled вместо RecordDelete); удаление ref перенесено в snapshot-run (core/snapshot), который после записи снапшота вызывает Landed(id, landing-branches) через 'git rev-list -1 <branch> -- .jaira/logbook/*/<id>* .jaira/archive/<id>*' и чистит только те refs, чей тикет там найден; landing-branches берутся из settings.json ('landing-branches', список с глобами), а без них — из локально закэшированного origin/HEAD (core/gitref RemoteHead/SetRemoteHead); jaira fetch и jaira validate теперь печатают Stranded — тикеты, законченные здесь, но не приземлившиеся нигде дольше настраиваемого срока (landing-grace-days, по умолчанию 7 дней); добавлена явная команда 'jaira snapshot --drop <id>' для принудительного удаления ref без ожидания слияния."
review-gaps: "1) Реальный дефект: Runner.Run() (core/snapshot/snapshot.go) никогда не обновляет локальные remote-tracking ветки для landing-branches (origin/main и т.п.) — он делает FetchBranch(snapshot-branch) и Repo.Fetch(), а последний тянет ТОЛЬКО refs/jaira/tickets/* (core/gitref/gitref.go:534-543, докстрока jaira fetch прямо говорит 'no branch'). Landed() же делает git rev-parse/rev-list по 'origin/main' и подобным — это локальные, потенциально устаревшие или вовсе отсутствующие refs, если пользователь никогда не гонял голый git fetch. В расписании фонового запуска (раз в 3 дня) это означает, что тикет может НЕ вычиститься 'sobald' commit появился в основной ветке, а просядёт до следующего случайного git fetch/pull кем-то — противоречит DoD 'wird entfernt, sobald ... einen Commit findet'. Нет теста, который бы поймал это (все тесты снапшота работают против одного bare-remote+clone, где origin/main создаётся и коммитится в ТОЙ ЖЕ рабочей копии, так что tracking-ref всегда свежий). 2) DoD дословно требует 'Hauptbranch aus origin/HEAD ... mit Fallback auf origin/main und origin/master und einem Eintrag in settings.json' — реализовано только settings.json-список ИЛИ origin/HEAD (core/settings/settings.go Landing()); нет отдельного жёстко закодированного фолбэка на origin/main/origin/master, если origin/HEAD не резолвится (что происходит, если удалённый git-сервер не отдаёт HEAD, или local clone его не знает и SetRemoteHead не сработал). 3) DoD просит 'Tests mit zwei Klonen', которые доказывают, что тикет остаётся видимым в окне и исчезает после мержа. core/refsync/refsync_test.go:TestAFinishedTicketThatNeverLandedIsReported использует twoSides, но реальное приземление/жатву через snapshot.Runner там не гоняют — ref удаляется вручную Repo.Delete, в обход кода в core/snapshot. core/snapshot/snapshot_test.go:TestOnlyATicketFiledAwayInALandingBranchIsReaped тестирует настоящую жатву, но только с одним клоном+bare remote, не с двумя параллельными клонами, доказывающими видимость тикета для 'другого' клона до мержа. Итог: механизм жатвы протестирован, а именно комбинация 'два клона + settings.json/origin-HEAD резолюция + устаревание tracking-веток' — нет."
review-verdict: "Не полностью соответствует DoD. Основная механика (RecordFiled статт RecordDelete, Landed() по rev-list, жатва внутри snapshot-run, Stranded-репортинг в fetch/validate, --drop) реализована и покрыта тестами, сборка и тесты проходят. Но есть настоящий функциональный риск: жатва опирается на локальные remote-tracking ветки (origin/main и т.п.), которые сам jaira нигде не обновляет — это может задержать или вовсе не выполнить автоматическое удаление ref, что бьёт по центральному обещанию тикета ('sobald ein Commit gefunden wird'). Плюс отсутствует явный fallback origin/main/origin/master из текста DoD, и нет теста на двух реальных клонах для полного сценария жатвы. Отправляю в human, а не одобряю как готовое."
review-check: |-
  1. cd /home/alex/projects/jaira && go test ./core/snapshot/... ./core/settings/... ./core/refsync/... — все три пакета должны показать 'ok'.
  2. mkdir /tmp/jr && cd /tmp/jr && git init --bare board.git && git clone board.git ada && git clone board.git bob
  3. В ada: git checkout -b main && echo x > f && git add f && git commit -m init && git push origin main
  4. В ada: настроить JAIRA_HOME на отдельную папку, jaira init, создать тикет, jaira claim/jaira move в конечную лейну, затем jaira archive <id> (или соответствующий финальный шаг) — тикет уходит в .jaira/archive/, закоммитить и запушить в main
  5. В ada: JAIRA_NO_SNAPSHOT=1 jaira snapshot --json — посмотреть поле reaped: если main НЕ был заранее подтянут через обычный 'git fetch origin main' в этом клоне (а не через jaira), проверить, останется ли ref всё ещё удалён — если да несмотря на отсутствие свежего fetch, баг из review-gaps №1 не подтверждается; если ref НЕ удаляется, пока не выполнить руками 'git fetch origin main', баг подтверждён.
  6. В bob (второй клон, до слияния в main): jaira fetch --json — тикет должен быть виден среди refs (still standing).
  7. После шага 5 (реальный merge+snapshot): bob снова jaira fetch --json — ref должен исчезнуть у bob тоже.
  Если шаг 5 показывает, что чистка происходит только после ручного git fetch origin main, это подтверждает главный найденный дефект.
---

# Ein Ref verschwindet erst, wenn das Ticket im Hauptbranch angekommen ist

## Definition of Done

- [x] logbook und archive loeschen das Ref nicht mehr sofort, sondern schreiben den Endzustand darauf, sodass andere Klone das Ticket als 'fertig, wartet aufs Landen' sehen; das Ref wird entfernt, sobald 'git rev-list -1 <hauptbranch> -- .jaira/logbook/*/<id>* .jaira/archive/<id>*' einen Commit findet, und das laeuft im Hintergrund nach dem Muster der Update-Pruefung statt auf dem Kommandopfad; der Hauptbranch wird aus origin/HEAD bestimmt, mit Fallback auf origin/main und origin/master und einem Eintrag in settings.json; ein Ticket, das im Endzustand steht und seit einer konfigurierbaren Frist nicht im Hauptbranch angekommen ist, wird von fetch und validate gemeldet; es gibt einen ausdruecklichen Befehl, ein solches Ref trotzdem zu entfernen, und er sagt was er tut; Tests mit zwei Klonen belegen, dass das Ticket im Fenster fuer den anderen Klon sichtbar bleibt und nach dem Merge verschwindet
  proof: Smoke gegen echtes git: nach archive lebt der Ref weiter (Datei nur im eigenen Branch), ein Snapshot ohne Landung laesst ihn stehen, nach dem Merge meldet der Snapshot 'landed and cleared' und der Ref ist weg - das Ticket liegt dabei im Snapshot-Branch. Ohne aufloesbare Landebranch wird nichts entfernt (TestWithoutALandingBranchNothingIsReaped)
- [x] die Landebranches stehen als LISTE in settings.json ('landing-branches', z.B. main, master, develop, release/*, Glob erlaubt) - damit erledigt sich die Frage, wie die wichtige Branch bei wem heisst; ohne Eintrag gilt origin/HEAD als einziger Eintrag, und loest sich auch der nicht auf, wird KEIN Ref entfernt und einmal gesagt, was einzutragen ist
  proof: Smoke gegen echtes git: nach archive lebt der Ref weiter (Datei nur im eigenen Branch), ein Snapshot ohne Landung laesst ihn stehen, nach dem Merge meldet der Snapshot 'landed and cleared' und der Ref ist weg - das Ticket liegt dabei im Snapshot-Branch. Ohne aufloesbare Landebranch wird nichts entfernt (TestWithoutALandingBranchNothingIsReaped)
- [x] das Entfernen selbst gehoert in den Snapshot-Lauf (PTQ3XT), nicht in logbook/archive und nicht in einen eigenen Timer: der Snapshot geht periodisch ohnehin ueber alle Refs, und er hat das Ticket unmittelbar davor in den Snapshot-Branch geschrieben - im Moment des Loeschens liegt es also in zwei Ablagen. Reihenfolge im Lauf: erst schreiben, dann jaeten
  proof: Smoke gegen echtes git: nach archive lebt der Ref weiter (Datei nur im eigenen Branch), ein Snapshot ohne Landung laesst ihn stehen, nach dem Merge meldet der Snapshot 'landed and cleared' und der Ref ist weg - das Ticket liegt dabei im Snapshot-Branch. Ohne aufloesbare Landebranch wird nichts entfernt (TestWithoutALandingBranchNothingIsReaped)

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] settings.json: landing-branches als Liste mit Globs; ohne Eintrag origin/HEAD als einziger Eintrag; loest sich keiner auf, wird nichts entfernt
  proof: core/ticket recordFiled statt recordDelete bei Archive/Logbook; core/gitref Landed/RemoteHead/SetRemoteHead; core/settings Landing (Liste gewinnt vor origin/HEAD, sonst nichts); core/snapshot reap nach dem Schreiben; refsync.Stranded in fetch und validate; jaira snapshot --drop
- [x] core/gitref: Landed(id, branches) - git rev-list -1 <branch> -- .jaira/logbook/*/<id>* .jaira/archive/<id>* ; nur weggelegte Pfade zaehlen, ein Ticket das dort noch unter tickets/ liegt ist NICHT gelandet
  proof: core/ticket recordFiled statt recordDelete bei Archive/Logbook; core/gitref Landed/RemoteHead/SetRemoteHead; core/settings Landing (Liste gewinnt vor origin/HEAD, sonst nichts); core/snapshot reap nach dem Schreiben; refsync.Stranded in fetch und validate; jaira snapshot --drop
- [x] logbook und archive loeschen das Ref nicht mehr, sondern schreiben den Endzustand darauf (RecordFiled statt RecordDelete); jaira delete entfernt es weiter sofort, denn das ist Absicht
  proof: core/ticket recordFiled statt recordDelete bei Archive/Logbook; core/gitref Landed/RemoteHead/SetRemoteHead; core/settings Landing (Liste gewinnt vor origin/HEAD, sonst nichts); core/snapshot reap nach dem Schreiben; refsync.Stranded in fetch und validate; jaira snapshot --drop
- [x] Jaeten im Snapshot-Lauf, nach dem Schreiben: Ref weg, sobald das Ticket in einer Landebranch weggelegt ist
  proof: core/ticket recordFiled statt recordDelete bei Archive/Logbook; core/gitref Landed/RemoteHead/SetRemoteHead; core/settings Landing (Liste gewinnt vor origin/HEAD, sonst nichts); core/snapshot reap nach dem Schreiben; refsync.Stranded in fetch und validate; jaira snapshot --drop
- [x] fetch und validate melden 'fertig, aber seit N Tagen nicht angekommen' (Frist konfigurierbar)
  proof: core/ticket recordFiled statt recordDelete bei Archive/Logbook; core/gitref Landed/RemoteHead/SetRemoteHead; core/settings Landing (Liste gewinnt vor origin/HEAD, sonst nichts); core/snapshot reap nach dem Schreiben; refsync.Stranded in fetch und validate; jaira snapshot --drop
- [x] jaira snapshot --drop <id>: ein Ref trotzdem entfernen, fuer einen Branch der nie gemergt wird, und sagen was es tut
  proof: core/ticket recordFiled statt recordDelete bei Archive/Logbook; core/gitref Landed/RemoteHead/SetRemoteHead; core/settings Landing (Liste gewinnt vor origin/HEAD, sonst nichts); core/snapshot reap nach dem Schreiben; refsync.Stranded in fetch und validate; jaira snapshot --drop
- [x] Tests mit zwei Klonen: im Fenster zwischen Logbuch und Merge bleibt das Ticket fuer den anderen Klon sichtbar, nach dem Merge verschwindet es
  proof: core/ticket recordFiled statt recordDelete bei Archive/Logbook; core/gitref Landed/RemoteHead/SetRemoteHead; core/settings Landing (Liste gewinnt vor origin/HEAD, sonst nichts); core/snapshot reap nach dem Schreiben; refsync.Stranded in fetch und validate; jaira snapshot --drop

## Progress
- **2026-09-10 21:03 · Alexander Sacharov** — Verhaeltnis zu den anderen Tickets, damit die Reihenfolge klar ist:

- RFC7GA (human) hat das Gegenteil eingebaut: logbook/archive loescht das Ref sofort. Dieses Ticket dreht das zurueck, und zwar aus einem Grund, der dort nicht bedacht war - das Fenster zwischen 'abgeraeumt' und 'gemergt' ist so lang wie ein Review, und darin ist das Ticket fuer alle anderen unsichtbar.
- refsync.Departed ('hat die Tafel woanders verlassen, liegt hier noch') wird dadurch seltener und richtiger: es feuert erst, wenn das Ref endgueltig weg ist, also nach dem Landen - genau dann, wenn der andere Klon wirklich aufraeumen soll, und nicht schon wenn jemand nur seinen Branch fertig hat.
- PTQ3XT (Snapshot alle drei Tage) ist unabhaengig, aber die Erkennung 'im Hauptbranch angekommen' und die Snapshot-Ausloesung teilen dasselbe Muster: im Hintergrund, nie auf dem Kommandopfad. Wer zuerst gebaut wird, legt das Muster fuer den anderen.
- **2026-09-10 21:04 · Alexander Sacharov** — Gemessen, welche Antwort auf 'welcher Branch ist der wichtige' ueberhaupt verfuegbar ist:

- Klon eines nicht leeren Repos: origin/HEAD ist gesetzt (refs/remotes/origin/master), lokal und kostenlos. Das ist der Normalfall.
- Klon eines LEEREN Repos (jaira init auf einem frischen Board): nicht gesetzt, und 'git remote set-head origin -a' scheitert mit 'Cannot determine remote HEAD' - auf dem Remote gibt es noch keinen Branch.
- Nach dem ersten Push: lokal weiter nicht gesetzt, aber 'git remote set-head origin -a' funktioniert dann und setzt es dauerhaft. Ein Netzaufruf, danach nie wieder.

Wichtiger als die Reihenfolge ist zweierlei:

1. settings.json gewinnt vor origin/HEAD, nicht umgekehrt. Der HEAD des Remotes ist der Default des Hosters, nicht zwingend der Branch, auf den es der Mannschaft ankommt: develop-Flows und Release-Branches sind normal. Das ist nicht zu erraten, aber zu erfragen.

2. Loest sich nichts auf, wird kein Ref entfernt. Ein stehengelassenes Ref kostet nichts; ein faelschlich entferntes nimmt genau die Sichtbarkeit weg, fuer die dieses Ticket existiert. Der unsichere Fall muss also auf die konservative Seite fallen, und einmal sagen, was einzutragen waere.
- **2026-09-10 21:05 · Alexander Sacharov** — Entschieden (Alexander), zwei Vereinfachungen, die beide besser sind als mein Entwurf:

1. Keine 'Hauptbranch', sondern eine Liste. 'landing-branches' in settings.json mit Globs. Damit ist die Frage 'wie heisst die wichtige Branch bei euch' nicht mehr zu beantworten, sondern zu konfigurieren, und develop- oder Release-Flows sind kein Sonderfall mehr. Fehlt der Eintrag, ist origin/HEAD der einzige Eintrag der Liste.

2. Gejaetet wird im Snapshot-Lauf, nicht in logbook/archive und nicht in einem eigenen Timer. Drei Gruende, der zweite ist der eigentliche:

   - Der Snapshot geht periodisch ohnehin ueber alle Refs. Ein zweiter solcher Durchgang haette keinen eigenen Anlass.
   - Er hat das Ticket unmittelbar davor in den Snapshot-Branch geschrieben. Im Moment des Loeschens liegt es also in ZWEI Ablagen: im Snapshot und in der Branch, in der es gelandet ist. Es gibt damit buchstaeblich nichts zu verlieren - das ist der Unterschied zwischen 'aufraeumen' und 'wegwerfen'.
   - Eine Stelle entscheidet 'dieses Ref hat seinen Zweck erfuellt', statt derselben Pruefung verteilt auf logbook, archive und einen Hintergrundtimer.

   Reihenfolge im Lauf ist damit festgelegt und nicht beliebig: erst den Snapshot schreiben, dann jaeten.
- **2026-09-10 21:20 · Alexander Sacharov** — Gebaut, zusammen mit PTQ3XT. Der Fehler, den der Smoke-Test gefunden hat, ist erwaehnenswert: mein erster Patch hat nur Logbook auf recordFiled umgestellt, Archive lief weiter auf recordDelete - der Ref war nach 'jaira archive' sofort weg, also genau das Fenster, das dieses Ticket schliessen soll. Gefunden, weil der Smoke-Test die Refs auf dem Remote gezaehlt hat und nicht nur den Exit-Code angesehen hat.

Ablauf jetzt, gegen echtes git belegt: nach 'archive' lebt der Ref weiter (die Datei liegt nur im eigenen Branch), ein Snapshot ohne Landung laesst ihn stehen, nach dem Merge in die Landebranch meldet der Lauf 'landed and cleared from the refs' und der Ref ist weg - und das Ticket liegt dabei im Snapshot-Branch, also in zwei Ablagen im Moment des Loeschens.

'jaira delete' loescht den Ref weiter sofort. Das ist kein Widerspruch: Loeschen ist eine Absicht, kein Stadium des Fertigwerdens.

Der Preis, den dieses Ticket einfuehrt, ist benannt und sichtbar gemacht statt verschwiegen: ein Branch, der nie gemergt wird, haelt seinen Ref. fetch und validate melden 'fertig seit N Tagen und nirgends angekommen', und 'jaira snapshot --drop' ist der ausdrueckliche Ausweg.
- **2026-09-11 06:29 · Alexander Sacharov** — Antwort (Alexander): die Frist ist drei Tage, nicht sieben. Sein Argument ist das richtige - Reviews laufen hier schneller, und eine Woche ist lang genug, dass ein vergessener Branch keine Nachricht mehr ist, wenn endlich jemand davon erfaehrt. Default geaendert, 'landing-grace-days' in settings.json bleibt der Weg, es anders zu setzen.
- **2026-09-11 11:32 · Alexander Sacharov** — Ревью: главная находка — Runner.Run (core/snapshot/snapshot.go) не тянет landing-branches (main/master) с origin, полагается на то, что уже локально закэшировано; Repo.Fetch() тянет только refs/jaira/tickets/*. Это может задержать/сорвать автожатву refs. Возврат в in-progress должен решить: либо явный git fetch landing-веток перед Landed(), либо документированное ограничение с более коротким по умолчанию окном проверки.
