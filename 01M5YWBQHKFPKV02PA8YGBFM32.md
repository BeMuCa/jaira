---
id: 01M5YWBQHKFPKV02PA8YGBFM32
title: "jaira update bringt standardmaessig alles auf Stand - Abschnitt, Rollen, Cockpit, Skill - und 'jaira update diff' zeigt, was zurueckliegt"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Ein einziges 'jaira update' bringt Board-Abschnitt, Rollen-Skills, Cockpit-Mod und jaira-Skill auf den Stand des installierten Binaries; 'jaira update diff' sagt vorher, was zurueckliegt."
context: |-
  Wunsch von berkc am 27.10.2026: 'Jaira update should do all of it by default, from skill to mod to roles. jaira update diff can be used to see which lanes or roles are behind.'
  Heute (0.5.6): 'jaira update' schreibt nur den jaira-Abschnitt in CLAUDE.md/AGENTS.md und die Release-Notes. Rollen und Cockpit kommen nur ueber 'jaira roles install --global|--project'.
  PR #91 (T8A8KM) fuegt 'update --roles', '--lanes' und '--dry-run' hinzu - Rollen sind dort noch opt-in.
  Offen: Lanes auch per Default? Lanes sind geteilte Board-Konfiguration und brauchen einen Commit - Vorschlag: Rollen/Mod/Skill per Default, Lanes nur mit --lanes, 'diff' zeigt beides.
  Gefunden am 23.10.: 'roles install' hielt zwei alte, unveraenderte Cockpit-Dateien (board.ts aus 5bc22fb1, register.tsx aus 2c4c6790) fuer 'von dir bearbeitet', weil ihr Hash in core/shipped/shipped.sha256 fehlt - auch in #91. Vermutung (ungeprueft): scripts/shipped-history.sh geht nicht durch zweite Merge-Eltern. Ohne den Fix wuerde ein Default-Update genau diese Dateien stehen lassen.
definition-of-done: "'jaira update' ohne Flags aktualisiert Abschnitt, Rollen, Cockpit-Mod und jaira-Skill ueberall, wo sie schon installiert sind; bearbeitete Dateien bleiben (exit 3)"
tags:
  - cli
blocked-by: []
related:
  - 01M21PJTWPK4DZ89PRY6T8A8KM
commits: []
created-at: 2026-10-27T14:31:21Z
updated-at: 2026-10-30T21:41:03Z
updated-by: berkc
claimed-by: berk-ca-2129445
claimed-at: 2026-10-30T21:41:03Z
---

# jaira update bringt standardmaessig alles auf Stand - Abschnitt, Rollen, Cockpit, Skill - und 'jaira update diff' zeigt, was zurueckliegt

## Definition of Done

- [ ] 'jaira update' ohne Flags aktualisiert Abschnitt, Rollen, Cockpit-Mod und jaira-Skill ueberall, wo sie schon installiert sind; bearbeitete Dateien bleiben (exit 3)
- [ ] 'jaira update diff' listet je Rolle, Mod-Datei und Lane: aktuell, zurueck (welche Version) oder bearbeitet; schreibt nichts
- [ ] Alte ausgelieferte Versionen aus Merge-Zweigen (board.ts 5bc22fb1, register.tsx 2c4c6790) gelten als ausgeliefert, nicht als bearbeitet
- [ ] Tests; Release-Note

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

