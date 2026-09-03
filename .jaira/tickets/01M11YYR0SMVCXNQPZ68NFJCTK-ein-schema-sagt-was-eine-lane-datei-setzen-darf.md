---
id: 01M11YYR0SMVCXNQPZ68NFJCTK
title: "Ein Schema sagt, was eine Lane-Datei setzen darf"
status: backlog
ready: false
creator: BeMuCa
goal: "Wer eine Lane schreibt, hat eine Stelle, die jedes Feld nennt, was es bedeutet, was erlaubt ist und was fehlt"
context: |-
  Berk am 27.08.: ein Schema fuer die Lane-Erstellung - Eingabe, Ausgabe, Beschreibung fuer das Modell, wann die Lane aufgerufen wird ('wenn das Modell eine Frage an den Nutzer hat'), was danach kommt (critique -> in-progress bei Kritik), der Prompt zusaetzlich zur Eingabe.

  Was heute da ist, verstreut: core/lane/lane.go parse() liest ~20 Frontmatter-Felder (id, name, after, precedence, agentic, model-tier, input-requires, output-produces, requires-question, requires-outcome, requires-human-exit, requires-option, rejects-to, requires-specified, terminal, description, creator ...) und der Markdown-Body ist der Prompt. 'jaira lanes template' druckt ein kommentiertes Skelett. 'jaira lanes show' zeigt den Vertrag einer Lane. Validierung: nur ein Teil (agentic braucht model-tier, rejects-to muss aufloesen).

  Was fehlt: eine einzige, vollstaendige Beschreibung - welche Felder es gibt, Typ, Pflicht/optional, Wechselwirkungen (requires-question ist heute nur ein Eingangs-Gate, kein Ausgangs-Gate: siehe die Luecke bei human). Ein unbekanntes Feld wird still ignoriert - ein Tippfehler in 'output-produces' faellt nicht auf.

  Vorschlag: das Schema als Tabelle in docs/LANES.md plus 'jaira lanes template' daraus erzeugt, und validate warnt bei unbekannten Feldern. Kein JSON-Schema-Werkzeug - kein Leser dafuer.
definition-of-done: "docs/LANES.md nennt jedes Frontmatter-Feld mit Typ, Pflicht und Bedeutung; 'jaira lanes template' und die Tabelle stimmen ueberein (Test); 'jaira validate' warnt bei einem unbekannten Feld in einer Lane-Datei"
blocked-by: []
commits: []
created-at: 2026-08-27T15:55:56Z
updated-at: 2026-08-27T15:55:56Z
---

# Ein Schema sagt, was eine Lane-Datei setzen darf

## Definition of Done

- [ ] docs/LANES.md nennt jedes Frontmatter-Feld mit Typ, Pflicht und Bedeutung; 'jaira lanes template' und die Tabelle stimmen ueberein (Test); 'jaira validate' warnt bei einem unbekannten Feld in einer Lane-Datei

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

