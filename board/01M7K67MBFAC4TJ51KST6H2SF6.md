---
id: 01M7K67MBFAC4TJ51KST6H2SF6
title: "CI prüft Go-Code mit golangci-lint (errcheck, staticcheck) und govulncheck"
status: ideas
ready: false
creator: Alexander Sacharov
goal: Jeder PR auf BeMuCa/jaira läuft durch golangci-lint (mindestens errcheck und staticcheck) und govulncheck; ein Fund lässt den CI-Lauf rot werden
context: |-
  Was fehlt: .github/workflows/ci.yaml prüft heute nur mit go test -race. Kein Linter, keine Prüfung der Abhängigkeiten auf bekannte Schwachstellen.
  Warum jetzt: KCRS3W fand Fehler der Art "Fehler verschluckt, Erfolg gemeldet". Genau das findet errcheck mechanisch, ohne Modell und ohne Bughunt-Lauf.
  govulncheck ist das offizielle Go-Werkzeug: meldet nur Schwachstellen in Abhängigkeiten, deren Code jaira wirklich aufruft.
  Abgrenzung: ersetzt die Rolle jaira-role-bughunt (KCRS3W) nicht. Linter fangen das Billige und Mechanische, Bughunt die Logik (Nebenläufigkeit, Fehlerpfade).
  Erwogen und verworfen: Wiz (kommerzielle Cloud-Sicherheitsplattform, jaira hat keine Cloud). CodeQL ist eine mögliche spätere Ergänzung (kostenlos für öffentliche Repos), nicht Teil dieses Tickets.
  Erster Lauf wird Funde im bestehenden Code bringen: entweder beheben oder in .golangci.yml begründet ausnehmen, nicht pauschal abschalten.
  CI gehört dem Kollegen (BeMuCa/jaira): vor dem Merge abstimmen.
definition-of-done: ci.yaml führt golangci-lint mit errcheck und staticcheck aus; ein Fund macht den Lauf rot
tags:
  - ci
blocked-by: []
related: []
commits: []
created-at: 2026-11-16T22:04:23Z
updated-at: 2026-11-16T22:04:23Z
---

# CI prüft Go-Code mit golangci-lint (errcheck, staticcheck) und govulncheck

## Definition of Done

- [ ] ci.yaml führt golangci-lint mit errcheck und staticcheck aus; ein Fund macht den Lauf rot
- [ ] ci.yaml führt govulncheck ./... aus; ein Fund macht den Lauf rot
- [ ] Funde des ersten Laufs sind behoben oder in .golangci.yml mit Begründung ausgenommen
- [ ] Mit dem Kollegen abgestimmt, bevor es upstream geht

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

