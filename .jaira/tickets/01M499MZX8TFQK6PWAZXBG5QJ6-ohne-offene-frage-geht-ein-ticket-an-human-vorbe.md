---
id: 01M499MZX8TFQK6PWAZXBG5QJ6
title: Ohne offene Frage geht ein Ticket an human vorbei direkt in review
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Ticket betritt die Lane human nur, wenn es eine Frage an einen Menschen trägt; sonst geht es von testing direkt nach review, und der Mensch nimmt es in signoff ab."
context: |-
  Was falsch ist: auf diesem Board steht human fest in der Reihenfolge zwischen testing und review (.jaira/lanes/order). Jedes Ticket hält dort an, auch wenn es nichts zu entscheiden gibt. Der Mensch wird zweimal gefragt: einmal in human, einmal in signoff.
  Zweite Folge, am 2026-10-06 an V3S0MW und HE3WWW gesehen: review-check schreibt erst die Lane review, die NACH human kommt. In human hat ein Ticket deshalb nie einen review-check, und /jaira-role-acceptance kann dort kein einziges Ticket abnehmen. Alex musste beide mit 'jaira move --to review' selbst herausholen.
  Was schon da ist: human.md trägt 'requires-question: true' — gedacht ist die Lane also als 'Agent hängt an einer Entscheidung', nicht als Pflichtstation. Und es gibt 'requires-option' (brainstorm.md, pre-process.md): eine Lane, die nur betreten wird, wenn das Ticket die Option ankreuzt.
  'jaira move <id> --to review --dry-run' aus testing ist heute schon erlaubt (an V3S0MW geprüft); die Lücke ist also nicht der Gate, sondern dass Lane-Prompts, Dispatcher und Rollen den Weg über human vorschreiben.
  Zu prüfen: testing.md-Prompt, jaira-dispatcher ('bis es in einer human-Lane sitzt'), jaira-role-lane, der Agent-Block in core/board/announce.go ('Loop'-Zeilen und Lane-Liste), und ob D28H7V ('Ein Sprung über eine Lane hinweg fällt auf') den Sprung testing->review dann als Fehler meldet.
  Entscheidung von Alex (2026-10-06): human nur bei offener Frage; die Abnahme durch einen Menschen passiert in signoff.
definition-of-done: "Ein Ticket ohne Frage geht aus testing nach review, ohne human zu betreten; mit --question geht es nach human wie heute"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-10-06T19:04:37Z
updated-at: 2026-10-07T06:07:58Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-10262
claimed-at: 2026-10-07T06:05:41Z
outcome-what: "Agent-Block, jaira-role-lane und Dispatcher sagen: human nur mit offener Frage, sonst next_lane (review); die Lane-Zeile einer Frage-Lane sagt 'without one, work passes it by'."
outcome-why: "Agenten lasen human als Station und parkten jedes fertige Ticket dort — ohne review-check, doppelte Abnahme durch den Menschen."
outcome-resolves: "Ein Ticket ohne Frage geht testing -> review -> signoff; der Mensch nimmt einmal ab, mit review-check."
---

# Ohne offene Frage geht ein Ticket an human vorbei direkt in review

## Definition of Done

- [x] Ein Ticket ohne Frage geht aus testing nach review, ohne human zu betreten; mit --question geht es nach human wie heute
  proof: core/lane/next.go:38 (Next überspringt RequiresQuestion, core/lane/next_test.go); Agent-Block und Rollen folgen jetzt next_lane
- [x] Der testing-Prompt, der Dispatcher und der Agent-Block beschreiben diesen Weg; keiner schickt ein Ticket ohne Frage nach human
  proof: core/board/announce.go (Lane-Zeile + Absatz 'next_lane'), core/role/builtin/jaira-role-lane/SKILL.md:283, core/role/builtin/jaira-dispatcher/SKILL.md:379
- [x] Ein Ticket in signoff hat einen review-check, sodass /jaira-role-acceptance es abnehmen kann
  proof: review kommt vor signoff und schreibt review-check; V3S0MW/HE3WWW am 2026-10-06 so abgenommen
- [x] Ein Sprung testing->review wird nicht als übersprungene Lane gemeldet
  proof: jaira move --dry-run testing->review: 'would be allowed', keine Warnung; Sprung-Erkennung (D28H7V) existiert noch nicht
- [x] core/release/NOTES.md hat eine Zeile unter Unreleased
  proof: core/release/NOTES.md:17

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-07 06:07 · Alexander Sacharov** — Die Mechanik war schon richtig: core/lane/next.go überspringt Lanes mit requires-question, next_lane aus testing ist review. In human landeten Tickets, weil Agent-Block ('until it sits in a human lane'), jaira-role-lane ('review and human are a person's lanes' — review ist hier agentisch) und Dispatcher human wie eine Station lasen. Geändert wurden nur diese Texte plus die Lane-Zeile für Frage-Lanes. Hinweis für D28H7V hinterlegt.
