---
id: 01M8BQF1BNAA3HVPZEXXNZJQ51
title: "Marketplace: was es schon gibt, wie andere es machen, und eine Einschaetzung"
status: backlog
ready: true
creator: Berk
assignee: Berk
goal: "Wer ueber den Marketplace fuer Lanes und Boards entscheidet, hat in einem Ticket den Stand, die Vorbilder und eine Empfehlung vor sich."
context: |-
  Wozu: Informationsgrundlage fuer das Diskussions-Ticket zum Marketplace (Berk und Alex). Kein Bauauftrag.

  WAS ES SCHON GIBT (im Code geprueft, 23.11.2026)
  - 'jaira lanes market' listet die Lanes aus dem Ordner lanes/ des jaira-Repos (BeMuCa/jaira), fest auf den Tag der laufenden Version (core/market/market.go). Eingebaut am 27.08.2026, Ticket YM7QSA.
  - 'jaira lanes market adopt <id>' holt eine Lane in den persoenlichen Katalog ~/.jaira/lanes - das ist schon das 'per ID installieren'.
  - Heute liegen dort nur 2 Lanes: changelog-writer und secrets-scan.
  - CI prueft lanes/: core/lane/shipped_test.go laedt den Ordner, scripts/ci-code-changed.sh zaehlt lanes/*.md als Code.
  - Eine Lane traegt creator: (core/lane/lane.go:168), gestempelt beim Veroeffentlichen und Uebernehmen.
  - Verwandt: AN1XR5 (jaira plugin install aus einem Katalog, Obsidian-artig, mit Pruefsummen); 4G07X9 (signoff: Boards teilen und uebernehmen wie Lanes).

  WIE ANDERE ES MACHEN (aus dem Gedaechtnis, Quellen nicht geprueft; HR3ZK3 hat Fables Recherche)
  - Obsidian: der Plugin-Code bleibt im Repo des Autors mit Releases. Ein zentrales Repo haelt nur eine Liste (community-plugins.json). Aufnahme per PR auf diese Liste plus Review. Installiert wird aus dem Release des Autors.
  - Claude-Code-Plugins: ein Marketplace ist jedes git-Repo mit einer marketplace.json, eingebunden per '/plugin marketplace add owner/repo'. Es gibt viele Marketplaces; eine Firma betreibt einfach ihren eigenen.
  - Homebrew-Taps und GitHub Actions: kein zentraler Server, man verweist direkt auf owner/repo@version.
  - Firmen: eigener Index oder Spiegel-Repo, eine Erlaubnisliste genehmigter Erweiterungen, Versionen auf einen Commit festgenagelt (z.B. private Marketplaces bei VS Code).

  EINSCHAETZUNG (Claude, 23.11.2026)
  - Gut: git-Repo als Ablage, CI laesst nur bestimmte Dateien rein, Installation per ID, Website als GitHub Pages aus dem Repo erzeugt. Alles ohne Server.
  - Anders machen: ein Index plus Autoren-Repos statt alles in einem Repo. Der Index haelt nur ID -> Repo + Commit; Lane oder Board bleibt beim Autor wie bei Obsidian. Das zentrale Repo bleibt klein, Autoren behalten ihre Versionen, eine Firma setzt ihre eigene Index-URL.
  - Hochladen ueber die Website erstmal weglassen: braucht ein Backend oder GitHub-Login und bricht 'kein Server'. Ein Link, der auf GitHub einen PR oeffnet, leistet fast dasselbe.
  - Sicherheit ist der Kern: eine Lane traegt einen Prompt (Prompt-Injection), seit ENCFVZ auch script: (Code-Ausfuehrung). Vor dem Uebernehmen muss jaira den Inhalt zeigen und den Commit festnageln; eine Lane mit Script braucht eine extra Bestaetigung.
  - Kleinster erster Schritt: lanes market auf Boards erweitern (4G07X9), eine statische Seite, die den Ordner listet. Danach Index und Autoren-Repos.
definition-of-done: Beim Diskussions-Ticket zum Marketplace gelesen und besprochen
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-11-26T10:47:18Z
updated-at: 2026-11-26T10:47:43Z
updated-by: Berk
---

# Marketplace: was es schon gibt, wie andere es machen, und eine Einschaetzung

## Definition of Done

- [ ] Beim Diskussions-Ticket zum Marketplace gelesen und besprochen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

