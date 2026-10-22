---
id: 01M5H9HS15YA6ZTGREC4T22Z2C
title: "CI und ein lokaler Hook lehnen Commits ab, die Arbeitgeber- oder Kundennamen tragen"
status: in-progress
ready: true
creator: BeMuCa
goal: "Kein Commit, keine Datei und keine Ticket-Notiz im Repo nennt Berks Arbeitgeber oder Kunden - weder im Inhalt noch in Author-/Committer-E-Mail - und die Pipeline weist einen PR ab, der es doch tut."
context: |-
  Am 2026-10-22 fiel auf: die globale git-E-Mail auf Berks Rechner lag auf einer Firmen-Domain, und 927 Commits auf master tragen sie. Dazu nennt core/role/builtin/jaira-role-pr/SKILL.md:84 einen Firmen-GitLab-Host, und fuenf Logbuch-/Archiv-Tickets haben Firmennamen in Notizen.
  Die E-Mail fuer dieses Repo ist seit dem 2026-10-22 die GitHub-noreply-Adresse (git config --local). Die offenen Branches wurden umgeschrieben; master wird am Sonntag 2026-10-26 umgeschrieben.
  Die Liste der verbotenen Strings darf selbst nicht im Repo liegen - sonst ist sie die Erwaehnung. Entschieden: ein GitHub-Repository-Secret haelt die Regex-Liste; die CI liest sie. Lokal liest ein pre-commit/pre-push-Hook dieselbe Liste aus einer Datei unter ~/.config/.
  Verworfen: DataFog oder ein anderer PII-Scanner - vier bekannte Strings sind ein Regex, kein ML-Problem.
definition-of-done: "Ein CI-Job in .github/workflows/ci.yaml prueft bei jedem PR den Diff gegen master und Author-/Committer-E-Mail jedes Commits gegen die Regex-Liste aus dem Secret MENTION_DENYLIST und schlaegt bei einem Treffer mit Zeile und Commit fehl; fehlt das Secret, schlaegt der Job mit einer klaren Meldung fehl statt still zu bestehen"
tags:
  - ci
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T07:52:29Z
updated-at: 2026-10-22T07:58:45Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-4193148
claimed-at: 2026-10-22T07:55:54Z
---

# CI und ein lokaler Hook lehnen Commits ab, die Arbeitgeber- oder Kundennamen tragen

## Definition of Done

- [ ] Ein CI-Job in .github/workflows/ci.yaml prueft bei jedem PR den Diff gegen master und Author-/Committer-E-Mail jedes Commits gegen die Regex-Liste aus dem Secret MENTION_DENYLIST und schlaegt bei einem Treffer mit Zeile und Commit fehl; fehlt das Secret, schlaegt der Job mit einer klaren Meldung fehl statt still zu bestehen
- [ ] scripts/check-mentions.sh macht dieselbe Pruefung lokal, liest die Liste aus $XDG_CONFIG_HOME/jaira/mention-denylist (Fallback ~/.config/jaira/mention-denylist), und docs/COMMANDS.md sagt, wie man es als pre-commit- und pre-push-Hook einhaengt
- [ ] core/role/builtin/jaira-role-pr/SKILL.md:84 nennt keinen Firmen-Host mehr, sondern ein generisches Beispiel
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

