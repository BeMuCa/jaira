---
id: 01M6H73DR7R7BTZXV2AA3Q92JY
title: jaira bewahrt API-Tokens nur lokal in ~/.jaira/secrets.json mit Rechten 0600 auf; jaira web verwaltet sie auf einer Einstellungsseite
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Ein PM hinterlegt API-Tokens (zuerst fuer Jira) im Browser, und sie verlassen nie seinen Rechner: in keinem Repo, keiner Antwort der Seite und keinem Log."
context: |-
  Was heute fehlt: es gibt keinen Ort fuer Tokens. ~/.jaira/settings.json ist keiner: es wird mit 0644 geschrieben (core/settings/settings.go:340), also lesbar fuer jeden Benutzer der Maschine.
  Teil der Aufteilung von CMN1YT (jaira web als Zentrale fuer PMs = Leute ohne Claude und ohne CLI); dort stehen Brainstorm und Antworten des Maintainers.
  Entschieden (Maintainer): zuerst eine Datei mit Rechten 0600; Schluesselbund und Secret-Manager spaeter. Tokens bleiben nur beim PM und werden nirgends geteilt.
  Fuer den Maintainer - was ist eine '0600-Datei'? Unter Linux und macOS hat jede Datei drei Rechte-Gruppen: Besitzer, Gruppe, alle anderen. 0600 heisst: der Besitzer darf lesen und schreiben (6), Gruppe und alle anderen gar nichts (0, 0). Der Ordner bekommt 0700 (nur der Besitzer darf hinein). Das ist dieselbe Regel, die ssh fuer ~/.ssh/id_* verlangt. Wichtig: das ist KEINE Verschluesselung. Wer als dieser Benutzer angemeldet ist, jedes Programm dieses Benutzers und jedes Backup der Platte kann den Token lesen. Unter Windows wirkt chmod nicht so; dort schuetzt nur, dass die Datei im eigenen Benutzerprofil liegt (Zugriffsliste nur fuer den Benutzer).
  Fuer den Maintainer - wie Firmen Schluessel ueblicherweise ablegen: (a) im Schluesselbund des Betriebssystems - macOS Keychain, Windows Credential Manager, unter Linux der Secret Service (GNOME Keyring/KWallet); verschluesselt, mit dem Login entsperrt; so speichert z.B. git seine Zugangsdaten ueber einen Credential-Helper. (b) in einem zentralen Secret-Manager (ein Tresor-Dienst der Firma oder ihres Cloud-Anbieters, oder ein Team-Passwortmanager); Tokens werden zur Laufzeit geholt, oft als Umgebungsvariable, zentral widerrufbar und protokolliert. (c) kurzlebige Tokens per SSO/OAuth statt fester API-Tokens. Fallstrick: WSL und Server haben oft keinen Secret Service - deshalb bleibt die 0600-Datei auch spaeter der Rueckfall.
  Ausgeschlossen hier: Schluesselbund, Secret-Manager, Umgebungsvariable - spaeter, eigenes Ticket.
  Bekannt: jaira web lehnt schon fremde Hosts und Schreibzugriffe mit fremdem Origin ab (core/web/web.go:293-327); die Token-Endpunkte muessen dahinter liegen.
definition-of-done: "Tokens liegen in ~/.jaira/secrets.json, getrennt von settings.json; Datei 0600, Ordner 0700, auch wenn die Datei vorher mit anderen Rechten existierte; Test prueft die Rechte"
tags:
  - web
blocked-by:
  - NEYWJH
related: []
follows: 01M67K6Z3NSBTBXAK4EVCMN1YT
commits: []
created-at: 2026-11-03T17:25:23Z
updated-at: 2026-11-03T17:25:23Z
---

# jaira bewahrt API-Tokens nur lokal in ~/.jaira/secrets.json mit Rechten 0600 auf; jaira web verwaltet sie auf einer Einstellungsseite

## Definition of Done

- [ ] Tokens liegen in ~/.jaira/secrets.json, getrennt von settings.json; Datei 0600, Ordner 0700, auch wenn die Datei vorher mit anderen Rechten existierte; Test prueft die Rechte
- [ ] Einstellungsseite: Token anlegen (Anbieter, Adresse, Token), auflisten, loeschen; die Liste zeigt nur 'gesetzt, endet auf ...' mit den letzten 4 Zeichen; Test
- [ ] Kein Token steht in einer Antwort der Seite, im Web-Log oder unter dem Repo-Ordner; Test schreibt einen Token und durchsucht Antworten und Log danach
- [ ] docs/COMMANDS.md erklaert Ablageort, 0600 und dass es keine Verschluesselung ist; eine Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

