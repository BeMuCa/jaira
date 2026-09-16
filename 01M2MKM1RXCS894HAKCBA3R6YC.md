---
id: 01M2MKM1RXCS894HAKCBA3R6YC
title: "Die pr-Rolle oeffnet den Pull Request selbst, wenn ein Mensch sie aufruft"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer /jaira-role-pr tippt, bekommt einen offenen Pull Request und keine Befehlszeile zum Einfuegen; ruft ein Agent die Rolle auf, bleibt es beim Pushen."
context: |-
  Heute pusht die Rolle jaira-role-pr den Branch und schreibt eine 'gh pr create'-Zeile aus, die der Mensch danach selbst einfuegt und ausfuehrt.

  Alex am 2026-09-16: 'если я стартую /jaira-role-pr то она должна сразу и создавать, просто без команды'. Sein Aufruf IST die Anweisung, den Pull Request zu oeffnen - der Kopierschritt danach kauft nichts und kostet eine Runde.

  Die Regel dahinter bleibt: wer die Aenderung geschrieben hat, entscheidet nicht, dass sie ankommt. Sie wird nur anders aufgehaengt - nicht mehr an 'kein Agent oeffnet je einen PR', sondern an 'nur der Aufruf eines Menschen oeffnet einen'. Ein Dispatcher oder Teamlead, der die Rolle aufruft, pusht weiter nur und reicht die Zeile zurueck. Merge und Approve bleiben in jedem Fall verboten.

  Der Text ist bereits geschrieben und liegt auf dem lokalen Branch chore/pr-role-opens-the-pr (Commit 5b849ca): core/role/builtin/jaira-role-pr/SKILL.md umgeschrieben, eine Zeile in core/release/NOTES.md. Er ist noch nicht durch critique, optimize oder review gelaufen und traegt keine Ticket-Id im Commit.

  Dabei gefunden und offen: in einem Fork zielt 'gh pr create' nicht selbstverstaendlich auf das richtige Repository - origin ist der Fork, die Pull Requests dieses Boards gehen aber alle nach BeMuCa/jaira. Die Rolle muss das pruefen, bevor sie das Kommando ausfuehrt, sonst steht der Pull Request im falschen Repository und ein Mensch muss ihn von Hand schliessen.
definition-of-done: "core/role/builtin/jaira-role-pr/SKILL.md sagt: Aufruf durch einen Menschen -> pushen und oeffnen; Aufruf durch einen Agenten -> pushen und die Zeile zurueckgeben; merge und approve bleiben in beiden Faellen verboten; die Rolle prueft vor dem Oeffnen, in welches Repository der Pull Request geht; eine Zeile unter ## Unreleased in core/release/NOTES.md; go test ./core/role/... gruen"
tags:
  - docs
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-16T07:59:07Z
updated-at: 2026-09-16T08:05:01Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-11544
claimed-at: 2026-09-16T07:59:59Z
outcome-what: "core/role/builtin/jaira-role-pr/SKILL.md: der Zielrepository-Check steht jetzt als eigener Abschnitt 'Which repository it goes to' VOR dem Oeffnen, mit drei Faellen (kein Fork / Fork -> Parent bzw. jaira.remote / Widerspruch -> nicht oeffnen, fragen); die create-Befehle tragen --repo bzw. --target-project; die NOTES.md-Zeile unter ## Unreleased nennt den Check mit. Der Rollentext 'Mensch oeffnet, Agent reicht die Zeile zurueck, merge und approve nie' lag bereits als a553e45 auf dem Branch."
outcome-why: "Ohne den Check zielt 'gh pr create' in einem Fork auf origin, also den Fork - und dieses Board schickt alle Pull Requests nach BeMuCa/jaira. Ein Pull Request im falschen Repository sieht niemand und ein Mensch muss ihn von Hand schliessen. Der Absatz dazu stand vorher nach den create-Blocks: wer von oben liest, hat das Kommando dann schon getippt."
outcome-resolves: "Definition of Done vollstaendig: Mensch -> pushen und oeffnen, Agent -> pushen und Zeile zurueck, merge/approve in beiden Faellen verboten, Zielrepository vor dem Oeffnen geprueft, eine Zeile unter ## Unreleased, go test ./core/role/... gruen."
review-summary: |-
  core/role/builtin/jaira-role-pr/SKILL.md:96-106 lists existing pull requests before the target repository is settled; in a fork 'gh pr list --head' queries the fork, an open pull request on the parent is not found, and the role opens a second one — which :108 says it never does. Move the 'Which repository it goes to' section (:136-155) ahead of the listing and pass the settled repository to it: 'gh pr list --repo <owner/repo> --head …' and 'glab mr list --repo <path> --source-branch …'.
  core/role/builtin/jaira-role-pr/SKILL.md:93-94 ('if an agent did, the push is where you stop') contradicts :175-177, where an agent-invoked run still writes the description and hands back a filled-in create line. Delete the clause at :93-94 — the push section only needs the push. The human/agent split takes effect in exactly one place, the create command at :175, and the boundary at :197 restates it.
  core/role/builtin/jaira-role-pr/SKILL.md:145,152 'read the owner/repo off jaira.remote' cannot be done as written: jaira.remote holds a remote NAME, not owner/repo (core/settings/settings.go:149, core/gitref/gitref.go:139 — the example value is 'upstream'). Make case 2 take the target from 'gh repo view --json parent' and use 'git remote get-url "$(git config jaira.remote)"' only where a URL is actually needed.
  core/role/builtin/jaira-role-pr/SKILL.md:153-155 case 3 has no rung for an unset jaira.remote. On a board whose refs were never shared 'git config jaira.remote' prints nothing, and the ladder reads that as neither case 2 nor case 3 — or as a disagreement, and refuses to open. Say it explicitly: jaira.remote unset means case 2, the parent is the answer; case 3 needs both sides actually named.
---

# Die pr-Rolle oeffnet den Pull Request selbst, wenn ein Mensch sie aufruft

## Definition of Done

- [x] core/role/builtin/jaira-role-pr/SKILL.md sagt: Aufruf durch einen Menschen -> pushen und oeffnen; Aufruf durch einen Agenten -> pushen und die Zeile zurueckgeben; merge und approve bleiben in beiden Faellen verboten; die Rolle prueft vor dem Oeffnen, in welches Repository der Pull Request geht; eine Zeile unter ## Unreleased in core/release/NOTES.md; go test ./core/role/... gruen
  proof: core/role/builtin/jaira-role-pr/SKILL.md:9-27 (Mensch oeffnet, Agent reicht zurueck), :136 (Zielrepository vor dem Oeffnen), :178-183 (merge/approve verboten); core/release/NOTES.md:17; go test ./core/role/... ok

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Den Zielrepository-Check in core/role/builtin/jaira-role-pr/SKILL.md vor den create-Befehl ziehen und konkret machen: wie man das Ziel ermittelt (gh repo view / git config jaira.remote), was bei einem Fork gilt, und dass ohne klares Ziel nicht geoeffnet wird
- [x] NOTES.md-Zeile unter ## Unreleased so ergaenzen, dass der Zielrepository-Check darin steht
- [x] go test ./core/role/... laufen lassen, Ticket verschieben und zusammen mit dem Code committen

## Progress
- **2026-09-16 08:02 · Alexander Sacharov** — Der Text aus chore/pr-role-opens-the-pr lag beim Start dieser Lane schon als a553e45 auf feat/A3R6YC (mit Ticket-Id im Subject, auf neuerer Basis als 5b849ca). Diese Lane hat nur noch den Zielrepository-Check gebaut, den der Ticket-Kontext als offen markiert hatte.

Warum der Check eine eigene Ueberschrift bekam und nicht ein Absatz beim create-Befehl blieb: er stand vorher NACH den create-Blocks und sagte nur 'wenn die Basis nicht offensichtlich ist, frag nach'. Ein Modell, das den Block von oben liest, hat das Kommando dann schon getippt. Jetzt endet 'Write the description' mit 'du oeffnest nichts, bevor der naechste Abschnitt geklaert hat, wohin' - die Reihenfolge im Dokument IST die Durchsetzung, es gibt keinen Code, der das erzwingt.

Geprueft und verworfen: 'git rev-parse --abbrev-ref origin/HEAD' als Zielquelle - das nennt nur den Default-Branch des Forks, nicht das Zielrepository. gh/glab plus jaira.remote sind die einzigen zwei Quellen, die das Upstream ueberhaupt benennen; darum die Drei-Faelle-Leiter mit 'widersprechen sich -> nicht oeffnen, fragen'.

go test ./core/role/... prueft nur Installation und Parsing der Rolle, nicht den Inhalt von SKILL.md. Gruen heisst hier also 'nichts kaputt', nicht 'Regel getestet'.
- **2026-09-16 08:04 · Alexander Sacharov** — critique (Durchgang 1, 4 Findings): (1) Die Reihenfolge im Dokument ist noch nicht fertig gedreht. Der Zielrepository-Abschnitt kam vor die create-Befehle, aber 'gh pr list --head' bei :99 laeuft immer noch davor und fragt im Fork den Fork - ein offener Pull Request im Parent wird nicht gefunden, und die Rolle oeffnet einen zweiten, obwohl :108 sagt, sie tue das nie. Der Abschnitt muss vor das Listing, und das Listing braucht das Repository als Flag. (2) :93-94 'if an agent did, the push is where you stop' widerspricht :175-177, wo der Agent noch Beschreibung schreibt und die Zeile zurueckreicht. Die Aufteilung Mensch/Agent gehoert an genau eine Stelle - den create-Befehl - sonst liest jeder Durchlauf zwei verschiedene Regeln. (3) jaira.remote haelt einen Remote-NAMEN, nicht owner/repo (core/settings/settings.go:149) - ':152 read the owner/repo off it' ist so nicht ausfuehrbar. (4) Fall 3 der Leiter hat keine Sprosse fuer 'jaira.remote ist gar nicht gesetzt' - genau der Normalfall auf einem Board, dessen Refs nie geteilt wurden. Nicht geprueft, weil Sache der review-Lane: ob die Flags '--repo' / '--project' bei gh/glab genau so heissen.
