---
id: 01M26DTCJ9EZ631S0NFVNQCKF5
title: "Ein add/add-Merge verliert die andere Seite, weil die leere Basis abgelehnt wird"
status: done
ready: true
creator: Alexander Sacharov
goal: "Zwei Branches, die dieselbe Ticketdatei anlegen, mergen feldweise wie zwei Branches, die sie aendern - ohne dass eine Seite still verschwindet"
context: |-
  Belegt, nicht vermutet: git ruft den Merge-Driver bei einem add/add-Konflikt mit einer LEEREN Basisdatei auf. core/merge.Merge parst die Basis mit ticket.ParseDoc, das lehnt sie ab ('file does not begin with ---'), der Driver beendet mit 1, git markiert den Pfad als konfliktbehaftet und laesst 'ours' stehen.

  Der Schaden ist still: in der Datei stehen KEINE Konfliktmarker. Es sieht aus wie ein erfolgreicher Merge, aber Lane und updated-by der anderen Seite sind weg.

  Reproduktion (gemessen mit dem gebauten Binary): geteiltes Board, Ticket auf Branch ada anlegen und committen, dieselbe Datei auf Branch berk anlegen mit status: todo und updated-by: berk, dann 'git merge berk'. Ergebnis: 'CONFLICT (add/add)', 0 Konfliktmarker, status: backlog - berks Seite verloren. Direkt nachgestellt mit 'jaira merge-driver <leere-datei> ours theirs' - exit 1 mit genau dieser Meldung.

  Heute ist der Fall schwer zu treffen, weil zwei Tickets verschiedene Ids und damit verschiedene Dateinamen haben. Er wird zum Normalfall, sobald ein per Ref empfangenes Ticket lokal als Datei landet (Frage 1 an 8566KF): dann existiert dieselbe Datei auf zwei Branches. Der Fix gehoert trotzdem hierher und nicht dorthin - der Bug besteht auch ohne das.
definition-of-done: "core/merge.Merge behandelt eine leere oder nur aus Leerraum bestehende Basis als 'kein gemeinsamer Vorfahre' statt sie abzulehnen; ein add/add-Merge zweier Branches an derselben Ticketdatei laeuft ohne Konfliktmarker durch, behaelt die weitere Lane und vereinigt Listen beider Seiten; Prosafelder, die sich auf beiden Seiten unterscheiden, bleiben ein echter Konflikt und landen in conflict-theirs-<field>; ein Test gegen echtes git deckt genau den add/add-Fall ab"
tags:
  - concurrency
blocked-by: []
commits:
  - 68aae28834ff19a76ff59b0010771567258ae872
  - 55f39efb8d5f562cb5a7415cb69f14a91d2b18b7
  - 3dbf02839cb83e92d9cb8f0e38eb225b454782ea
  - d10cde734ec25ee527f5a5d72e2ce88c2c9c0667
created-at: 2026-09-10T19:48:21Z
updated-at: 2026-09-11T13:31:39Z
claimed-by: DESKTOP-RFTCH11-147473
claimed-at: 2026-09-11T11:24:06Z
updated-by: Alexander Sacharov
assignee: Alexander Sacharov
question: "Fix steht und ist gegen echtes git belegt. Bleibt die Frage aus 8566KF: soll ein per Ref empfangenes Ticket lokal als Datei landen (dann ist es eine normale Karte in list, show, claim, move) oder als read-only Karte ohne Datei angezeigt werden?"
outcome-what: "ревью подтвердило diff против DoD, дефектов не найдено"
outcome-why: "второй проход модели пройден, ждём человека для приёмки"
outcome-resolves: "все 4 пункта DoD проверены построчно и тестами, go build и оба новых теста PASS"
review-summary: "core/merge.Merge (core/merge/merge.go:87-91) заменяет пустую или состоящую только из пробелов базу на минимальный валидный документ '---\\n---\\n' перед ParseDoc, вместо того чтобы отклонять её ошибкой. Раз база теперь парсится как документ без полей, дальше срабатывают уже существующие правила трёхстороннего слияния как для поля, не унаследованного ни одной стороной: mergeStatus сравнивает через Precedence и берёт более дальнюю lane, mergeList объединяет списки (tags/blocked-by/commits), а mergeScalar для прозы (goal/context/dod/...) при расхождении обеих сторон создаёт настоящий конфликт (conflict-theirs-<field>). Добавлены unit-тест в core/merge/merge_test.go (эмулирует пустую/пробельную базу напрямую) и интеграционный тест в internal/cli/mergebranches_test.go, который реально собирает бинарь, создаёт add/add конфликт через git и проверяет итоговый файл."
review-gaps: none
review-verdict: "Диф закрывает все 4 пункта DoD: пустая/пробельная база принимается (core/merge/merge.go:87-91), add/add merge проходит без маркеров конфликта с сохранением дальней lane и объединением списков (подтверждено TestAnEmptyBaseIsAnAbsentAncestorNotABrokenFile и TestTwoBranchesThatBothCreateTheTicketStillMergeFieldAware), расходящаяся проза остаётся настоящим конфликтом в conflict-theirs-<field> (тест на поле goal), и тест против настоящего git есть. Локально прогнал go build ./... (чисто) и оба новых теста — оба PASS. Изменение точечное, использует уже существующую инфраструктуру слияния без новых спецслучаев в driver'е, как и было заявлено в outcome. Дефектов не нашёл."
review-check: "1. cd /home/alex/projects/jaira 2. go test ./core/merge/... -run TestAnEmptyBaseIsAnAbsentAncestorNotABrokenFile -v — должен быть PASS, это прямой тест на пустую/пробельную базу и на конфликт по прозе (поле goal) 3. go test ./internal/cli/... -run TestTwoBranchesThatBothCreateTheTicketStillMergeFieldAware -v — собирает реальный бинарь jaira, создаёт два branch (ada и berk), которые оба создают один и тот же ticket-файл (add/add конфликт), мержит их через настоящий git merge и проверяет итоговый файл 4. В выводе теста ищите: 'status: todo' в итоговом файле (дальняя lane выиграла, хотя у ady был более новый timestamp), оба тега 'concurrency' и 'cli' присутствуют (списки объединены), и отсутствие '<<<<<<<' (конфликтных маркеров) — все эти проверки уже есть внутри теста как t.Errorf/t.Fatalf, так что просто смотрите на PASS/FAIL"
---

# Ein add/add-Merge verliert die andere Seite, weil die leere Basis abgelehnt wird

## Definition of Done

- [x] core/merge.Merge behandelt eine leere oder nur aus Leerraum bestehende Basis als 'kein gemeinsamer Vorfahre' statt sie abzulehnen; ein add/add-Merge zweier Branches an derselben Ticketdatei laeuft ohne Konfliktmarker durch, behaelt die weitere Lane und vereinigt Listen beider Seiten; Prosafelder, die sich auf beiden Seiten unterscheiden, bleiben ein echter Konflikt und landen in conflict-theirs-<field>; ein Test gegen echtes git deckt genau den add/add-Fall ab
  proof: core/merge/merge.go (leere/Leerraum-Basis -> emptyDoc); core/merge/merge_test.go TestAnEmptyBaseIsAnAbsentAncestorNotABrokenFile (vier leere Formen, weitere Lane, Listen-Union, Prosa bleibt Konflikt auf 'goal'); internal/cli/mergebranches_test.go TestTwoBranchesThatBothCreateTheTicketStillMergeFieldAware (echter git merge des add/add-Falls, 0 Marker)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-10 19:51 · Alexander Sacharov** — Ursache genau: git ruft den Merge-Driver beim add/add-Konflikt mit einer leeren %O-Datei auf. ParseDoc lehnt sie ab, der Driver endet mit 1, git laesst 'ours' liegen und markiert den Pfad - ohne Konfliktmarker in der Datei. Das ist der stille Teil: es liest sich wie ein sauberer Merge.

Fix: eine leere Basis ist ein ABWESENDER Vorfahre, kein kaputtes File. Sie wird durch ein leeres Ticket ('---\n---\n') ersetzt, das kleinste, das ParseDoc annimmt. Damit greifen alle bestehenden Regeln so, wie sie es fuer ein Feld tun, das keine Seite geerbt hat: weitere Lane gewinnt, Listen vereinigen, Prosa mit zwei verschiedenen Saetzen bleibt ein echter Konflikt. Keine neue Regel, keine Sonderbehandlung im Driver.

Warum nicht im Driver abfangen: der Driver ist nur der Aufrufer. Die Frage 'was heisst keine gemeinsame Basis' ist eine Merge-Regel und gehoert zu den anderen, sonst hat der zweite Aufrufer (das Board via refsync.Reconcile, das denselben Fall hat, wenn ein Ref noch keinen Parent-Commit hat) dieselbe Luecke wieder.

Handprobe nach dem Fix, gemessen: 'git merge berk' laeuft durch (1 file changed), 0 Konfliktmarker, status: todo obwohl adas Zeitstempel neuer war, tags: cli + concurrency. updated-by bleibt ada - das ist die dokumentierte Regel fuer uebrige Skalare (neueres updated-at gewinnt), nicht ein Verlust.

Damit ist Option A aus Frage 1 an 8566KF (empfangenes Ticket landet als Datei) nicht mehr durch diesen Bug blockiert.
