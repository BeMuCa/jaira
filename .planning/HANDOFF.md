# Handoff — 2026-08-31, Abend

Stand nach der Backlog-Flotte; dieses File ist das Gedächtnis. Der vorige
Handoff (Morgen desselben Tags) ist in der git-History dieses Files.

## Wo alles steht

- **master gepusht** bis einschließlich des 0.1.1-NOTES-Blocks. ~17 Commits
  heute: Schema-Spec, HREQJR (Wrap überall), dann die Flotte
  (DNAEPN, MFD7P3, MR3CN8, A1TZ4N, YBC0MT, TQXBY5, 88H1P4, GEC3TK).
  `go test ./... -race` (Cache geleert) grün auf dem Push-Stand;
  `gofmt -l` nur das bekannte `internal/cli/tickets.go`.
- **Binary** `~/.local/bin/jaira` neu gebaut; nach dem NOTES-Commit einmal
  nachziehen, falls der Stempel nicht HEAD ist.
- Zwei Spitzen-Commits wurden vor dem Push **re-splittet** (geteilter Index
  zweier paralleler Agents hatte sie verheddert); Baum-Identität bewiesen.

## Was Berk jetzt entscheidet (sein Inbox-Stand)

1. **signoff: 18 Tickets** warten auf Abnahme im Board — `--force` kommt
   nicht mehr aus der Lane heraus (`gate.go:377`, nur `req.Interactive`).
2. **human: BDV0HM** — Doku (SKILL.md:275, docs/AGENTS.md:157) verspricht
   den `--force`-Weg auf Zuruf, das Gate verweigert ihn. Welche Seite folgt?
3. **human: CD9TCB** — drei Optionen (durch 743737f erledigt erklären /
   Doku-Feld `overrides:` / Warnung bewusst reaktivieren).
4. **human: 7ZQ0ZN** — NOTES-0.1.1-Block ist auf master; Release = zwei
   Befehle: `git tag v0.1.1 && git push origin v0.1.1`.
5. **Schema-Spec-Review**: `docs/superpowers/specs/2026-08-31-schema-design.md`
   (`1e8fc57`) — Cut 1 startet erst auf sein Go. Offene Mini-Frage dazu:
   Entscheidung 10 nennt `verdict` nicht in der Reviewer-Reihenfolge; der
   Code zeigt es zwischen gaps und check.
6. **backlog, neu**: „lanes add setzt neben den Anker statt ans Ende" —
   Design-Frage (order.go ist Invarianten-Territorium), bewusst nicht gefixt;
   README sagt inzwischen ehrlich „Belongs, not lands".
7. **DoD-Klausel 3 von 88H1P4** ungetickt: frische Session + `jaira hook
   print`-Snippet in settings.json → arbeitet nach „go" bis signoff. Nur er
   kann das prüfen.

## Zurückgestellt, mit Grund

- **NFJCTK, B4MGTP, FCMP17**: stecken wörtlich im Schema-Spec (Cuts 1/4/6
  bzw. eigener Cut) — einzeln bauen hieße den abgestimmten Plan zerlegen.
- **88H1P4 Mechanismus c** (Block-Satz für Subagents): wartet auf sein Go.

## Arbeitsweise der Flotte (falls wiederholt)

Sequenziell je Ticket: Subagent implementiert und fährt das Board selbst
(claim → lanes → critique), Koordinator macht critique/optimize am selbst
gelesenen Diff, unabhängiges Opus-Review in Batches, Funde als Schleife
zurück. Lektionen (auch im Memory `jaira-code-learnings`): geteilter
git-Index verschluckt fremde staged Files → file-genau adden und notfalls
re-splitten; transiente FAILs ohne Testnamen = paralleler Schreiber
mid-compile; lipgloss polstert mehrzeilige Blöcke (styleLines).

## Ungeprüft / bekannt offen

- `human`-Lane weiterhin nicht exit-gated (nur signoff trägt
  requires-human-exit) — unverändert aus dem Morgen-Handoff, nicht ticketiert.
- Die grüne Lücke aus Berks Screenshot: Code kann sie nicht erzeugen
  (fitWindow-Repro leftover=0 auf jeder Breite); Diagnose: laufender Prozess
  mit veralteter Breite (Resize-Event verloren). Sein Wiggle-Test steht aus.
