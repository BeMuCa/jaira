---
id: 01M5ZJ5JPVKP4H0FDE5Q9T77A0
title: "Ein beratender PII-Check meldet persoenliche Daten in den neuen Zeilen eines PR, ohne ihn zu blockieren"
status: backlog
ready: false
creator: BeMuCa
goal: "Ein eigener Workflow .github/workflows/pii.yaml laeuft auf pull_request, scannt nur die hinzugefuegten Zeilen des PR mit DataFog und schreibt Treffer in die Job-Zusammenfassung; er faellt nie rot (continue-on-error), blockiert nichts und ersetzt die Denylist aus T22Z2C nicht."
context: |-
  Was fehlt: nichts im Repo prueft PR-Diffs auf E-Mail-Adressen, Telefonnummern, Steuer-IDs oder Kartennummern; die Denylist (T22Z2C) kennt nur die vier bekannten Firmen-Strings.
  Gewollt (Entscheidung 27.10.2026): beratend, nie blockierend - warnen, damit ein Mensch vor dem Merge hinschaut.
  Stand DataFog (PyPI, 27.10.2026): Version 4.9.0, Python >= 3.10, pip install datafog, CLI datafog scan-text '<text>' und --locale de fuer deutsche Kennungen; Python-API datafog.scan(text, engine='regex').entities. Vor dem Schreiben des Schritts den README-Stand noch einmal lesen - die CLI-Flags koennen sich geaendert haben.
  Aufbau: actions/checkout mit fetch-depth 0, actions/setup-python, pip install datafog, git diff base...head nur hinzugefuegte Zeilen (grep '^+' ohne '+++'), DataFog darueber, Ausgabe nach $GITHUB_STEP_SUMMARY, continue-on-error: true.
  Trigger: pull_request; pull_request_target nur, wenn T22Z2C mit B entschieden wird - und auch dann nie PR-Code ausfuehren, nur den Diff als Daten lesen.
  Nicht gewollt: ein Required Check, ein Secret, oder ein Scan des ganzen Baums.
definition-of-done: "Die Datei .github/workflows/pii.yaml existiert und laeuft auf pull_request; der Job hat continue-on-error: true"
tags:
  - ci
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:52:29Z
updated-at: 2026-10-27T20:52:29Z
---

# Ein beratender PII-Check meldet persoenliche Daten in den neuen Zeilen eines PR, ohne ihn zu blockieren

## Definition of Done

- [ ] Die Datei .github/workflows/pii.yaml existiert und laeuft auf pull_request; der Job hat continue-on-error: true
- [ ] Gescannt werden nur hinzugefuegte Zeilen des PR-Diffs (base...head), nicht der ganze Baum
- [ ] Treffer stehen in der Job-Zusammenfassung mit Datei, Zeile und Typ; ohne Treffer steht dort 'keine'
- [ ] Auf einem Test-PR mit einer erfundenen E-Mail-Adresse zeigt die Zusammenfassung den Treffer, und der PR bleibt gruen
- [ ] ci.yaml bleibt unveraendert

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

