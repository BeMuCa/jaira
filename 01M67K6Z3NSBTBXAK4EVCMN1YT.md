---
id: 01M67K6Z3NSBTBXAK4EVCMN1YT
title: "jaira web wird Beobachtungs- und Steuerzentrale fuer PMs ohne Claude: Onboarding mit Anbieterwahl, lokale Tokens, Projekte aus Repos, Tabs fuer Tickets, Aktive, Milestones/Epics, Schreiben"
status: ideas
ready: true
creator: berkc
assignee: berkc
goal: "Ein Projektmanager ohne Claude und ohne CLI-Kenntnisse verfolgt und steuert die Arbeit seiner Entwickler komplett im Browser, und Code und Ticket-Tracking (jaira, Jira, GitHub/GitLab Issues) passen zusammen."
context: |-
  Wunsch von berkc am 31.10.2026 (Punkte 2-5, wortgetreu gespeichert in der Claude-Memory prompt-2026-10-31-web-pm-and-lanes).
  Web-Nutzer heissen ab jetzt PM.
  - Nur git noetig; Infos an bestimmte Entwickler evtl. ueber Git-Refs weiterreichen.
  - curl-Onboarding: Wahl git / Jira / andere grosse Anbieter (YouTrack, GitLab/GitHub Issues, Epics, Projects), dann verbinden (API-Token bzw. git). Fehlt git: Link zu git ausgeben, oder Pfad zu git einfuegen bzw. in eine Umgebungsvariable schreiben. Danach oeffnet das Web.
  - Web-Settings: der PM legt weitere Tokens an; sie bleiben nur bei ihm, werden nirgends geteilt (Sicherheit zuerst).
  - Neues Projekt: Anbieter waehlen; bei git alle gefundenen Repos zur Auswahl oder Clone-Link einfuegen (Default, wenn keins gefunden). Das Projekt wird angelegt, das Repo geholt, und es zeigt Tickets, Milestone-Roadmap mit Strassennamen und die uebrigen Milestones.
  - Begriffe: Milestone = Epic. jaira-Milestones muessen sich als Epics nach Jira uebertragen lassen; im Web muss klar sein, dass Milestones Epics sind.
  - Tabs: 'Alle Tickets'; 'Aktiv' - kleine Avatare je Nutzer am Board, darunter seine Tickets, Farbe nach Stand: backlog grau, in Arbeit orange, in review orange-gruen, done gruen, blockiert rot; 'Milestones/Epics' mit Jira-(o.ae.)-Verbindung; 'Schreiben' - Tickets und Milestones anlegen, auch direkt aus einem Jira-Ticket uebernehmen.
  - Der PM weist Tickets und Milestones zu (ueber den festen jaira-Namen, siehe Ticket zum festen Namen).
  - Start von jaira web gibt eine schoen gestaltete Ausgabe, die den Server vorstellt.
  - Punkt 5 der Nachricht bricht ab: 'nach einer Zuweisung durch Kollegen oder PM siehst du es auf dem Board als ...' - Ende fehlt.
  Heute: jaira web gibt es (Filter, Assignee auf Karten, Lane-Farben, Log - ZDBXN1, TPCKZB, HEY6T7).
definition-of-done: "Brainstorm (starkes Modell): Funktionsliste, Aufteilung in Tickets, Jira-Mapping (Epics/Tasks/Subtasks vs. Strassenname/Milestone/Ticket), Token-Ablage, was fehlt"
tags:
  - cli
blocked-by: []
related:
  - 01M5MNWJF240306ZVDA41K3RA8
commits: []
created-at: 2026-10-30T23:44:38Z
updated-at: 2026-10-31T01:25:20Z
updated-by: berkc
---

# jaira web wird Beobachtungs- und Steuerzentrale fuer PMs ohne Claude: Onboarding mit Anbieterwahl, lokale Tokens, Projekte aus Repos, Tabs fuer Tickets, Aktive, Milestones/Epics, Schreiben

## Definition of Done

- [ ] Brainstorm (starkes Modell): Funktionsliste, Aufteilung in Tickets, Jira-Mapping (Epics/Tasks/Subtasks vs. Strassenname/Milestone/Ticket), Token-Ablage, was fehlt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-30 23:49 · berkc** — Brainstorm CMN1YT (31.10.2026, starkes Modell). Grundlage: jaira 0.5.7, jaira web --help, core/web/web.go:273-286, core/milestone, core/project, core/settings, core/identity, docs/COMMANDS.md, Punkte 2-5 der Nachricht.

HEUTE (geprueft)
- jaira web: lokale Seite auf 127.0.0.1, Epics rechts, Tickets ohne Epic links, Drag&Drop, Epic anlegen/bearbeiten, Filter, Ticket lesen. Tickets werden dort NICHT bearbeitet, keine Milestones, keine Zuweisung (Routen core/web/web.go:273-286).
- Projekte: core/project kennt schon Repo-Suche (Discover, 2 Ebenen tief) und eine gemerkte Projektliste (project.go Load/Remember). Das ist die Basis fuer 'Neues Projekt -> gefundene Repos'.
- Zuweisung an Entwickler: laeuft schon ueber Git-Refs. Tickets reisen auf refs/jaira/tickets/*, jaira fetch meldet ein neu zugewiesenes Ticket per Desktop-Benachrichtigung (COMMANDS.md:164). Es fehlt nur ein Web-Weg zum Schreiben und ein eindeutiger Name (DVPAPE).
- Name: core/identity/identity.go nimmt JAIRA_USER, dann git user.name, dann $USER; Aliasse global in ~/.jaira/identity. Daher BeMuCa/berkc doppelt.
- Strassenname: eine Roadmap-Swimlane hat keinen eigenen Namen. Sie ist der Baum unter einem Top-Milestone (after: [], core/milestone/layout.go:29). 'Strassenname' = Name des Top-Milestones.
- Installer scripts/install.sh (123 Zeilen) ist nicht interaktiv; er installiert Binary + Rollen.

1) EPICS VS MILESTONES
- Epic = WAS (ein Feature). Nur auf refs/jaira/epics/<name>, keine Datei, Unfertiges bleibt im Epic, Folge-Tickets (follows/parent) zaehlen automatisch mit.
- Milestone = WANN (eine Runde Arbeit). Datei .jaira/milestones/<name>-<ID>.md + Ref, Unfertiges wandert in die naechste Runde, hat Roadmap (after:), wird per logbook abgelegt.
- Ueberlappung: beides sind Ticket-Mengen mit Fortschritt 'X von Y'. Sonst zwei Achsen: ein Ticket ist in Feature 'login' UND in Runde 'q4'.
- Urteil: NICHT zusammenlegen. Wer sie vereint, verliert entweder 'Unfertiges zieht weiter' oder 'Folge-Tickets zaehlen mit'. Fuer den PM im Web aber EIN Tab 'Milestones & Epics' mit einem Satz Erklaerung je Typ.
- Jira-Mapping (Vorschlag als Default):
  jaira Epic -> Jira Epic (gleiche Semantik: Feature, Unfertiges bleibt).
  jaira Milestone -> Jira Fix Version (oder Sprint, falls Scrum-Board). Gleiche Semantik: Runde, Unfertiges zieht weiter.
  Roadmap-Swimlane (Strassenname) -> Jira Initiative (nur Premium/Plans) oder Label; ohne Premium: kein Gegenstueck, nur Label.
  jaira Ticket -> Story/Task/Bug (Typ aus Tag oder Default Task).
  DoD-Punkte -> Checkliste im Beschreibungstext, NICHT Subtasks (Subtasks haetten eigenen Status, jaira-DoD nicht).
  follows/parent/blocked-by -> Jira Issue Links (relates/blocks).
  Lane -> Jira Status per Zuordnungstabelle (backlog->To Do, in-progress..testing->In Progress, review/signoff->In Review, done->Done).
- GitHub: Milestone -> GitHub Milestone (1:1, sogar gleicher Begriff). Epic -> Eltern-Issue mit Sub-Issues. Ticket -> Issue. Lane -> Status-Feld eines Projects-Boards. Swimlane -> Label oder Project-Ansicht.
- GitLab: Milestone -> Milestone (oder Iteration). Epic -> Epic (nur Premium, Gruppenebene), sonst Eltern-Issue/Label. Lane -> Board-Label.
- YouTrack: Epic -> Issue-Typ Epic, Milestone -> Sprint im Agile Board.
- Vorschlag 'Strassenname = Jira-Epic, Milestones darunter = Tickets': nur teilweise richtig. Richtig: eine Swimlane ist ein langes Thema, also eher Epic/Initiative als Ticket. Falsch: Milestones enthalten Tickets; macht man Milestones zu Jira-Tickets, muessen die jaira-Tickets Subtasks werden. Subtasks koennen keine Kinder haben, verlieren eigene Workflows und Jira-Boards zeigen sie schlecht. Ausserdem hat die Swimlane keinen eigenen Namen (s.o.). Kompromiss: Mapping pro Verbindung einstellbar ('Milestone als Epic' als zweite Variante fuer Firmen ohne Fix Versions), Default wie oben.
- Wunsch 'Milestone = Epic fuer PMs' widerspricht der heutigen Definition in jaira epic --help ('Nor is it a milestone'). Muss entschieden werden (Frage 1 unten).

2) PM NUR MIT GIT
- Onboarding nicht im curl-Skript, sondern im Browser: install.sh bleibt klein und nicht interaktiv (curl|sh mit Fragen ueber /dev/tty bricht unter Windows/CI/WSL leicht), am Ende 'jaira web' starten. Die Seite fragt beim ersten Start: Name (DVPAPE), Anbieter (git / Jira / GitHub Issues / GitLab Issues / YouTrack), Verbindung.
- Fehlt git: jaira web erkennt es (exec.LookPath), zeigt Link git-scm.com/downloads und ein Feld 'Pfad zu git'. Gespeichert als 'git-path' in ~/.jaira/settings.json; JAIRA_GIT als Env-Alternative. Keine Shell-Profile anfassen.
- Git-Zugang: KEIN eigener Token. jaira ruft git auf, git nutzt seinen Credential-Helper (Git Credential Manager unter Windows/macOS). Kein Geheimnis in jaira.
- API-Tokens (Jira, GitHub-API, GitLab, YouTrack): Stufe 1 eine Datei ~/.jaira/secrets.json mit Rechten 0600 (Ordner 0700), getrennt von settings.json, ausserhalb jedes Repos, nie in einem Commit, nie in einer API-Antwort der Seite (die zeigt nur 'gesetzt, endet auf ...ab12'), nie im Web-Log (HEY6T7-Log filtern). Unter Windows liegt die Datei im Benutzerprofil (ACL nur fuer den Nutzer). Stufe 2 optional OS-Schluesselbund (macOS Keychain, Windows Credential Manager, Linux Secret Service) mit Fallback auf die Datei, weil WSL und Server oft keinen Secret Service haben. Bibliothek dafuer noch nicht geprueft.
- Zuweisung an Entwickler: der PM schreibt assignee auf den Ticket-Ref und pusht; beim naechsten Fetch des Entwicklers kommt die vorhandene Desktop-Benachrichtigung. Punkt 5 (abgebrochen) als Vorschlag: auf dem Board traegt das Ticket bis zum ersten Oeffnen ein Zeichen 'neu fuer dich von <Name>'.
- Fester Name (DVPAPE) ist Voraussetzung: Zuweisung vergleicht Strings. Der PM bekommt seinen Namen im Browser, gespeichert lokal pro Projekt. Die Avatare im Tab 'Aktiv' brauchen genau diese Namensliste.
- Haken: der PM braucht Schreibrecht auf das Remote (Push von refs/jaira/*). Ein PM mit nur Leserecht kann beobachten, aber nicht zuweisen.

3) TABS
- Alle Tickets: heutige Liste + Filter, ergaenzt um Lane-Spalten-Ansicht.
- Aktiv: Avatar (Initialen + Farbe aus Namen, kein Bilddienst) je Person, die offene Tickets hat; darunter seine Tickets. Farben: backlog/ideas/todo grau, brainstorm..testing orange, human/review/signoff orange-gruen, done gruen, Parkspur (heute 'blocked', in diesem Checkout geloescht) oder blocked-by offen rot. Zuordnung aus der Lane-Reihenfolge ableiten, nicht hart verdrahten, weil jedes Board eigene Lanes hat.
- Milestones/Epics: Roadmap mit Swimlanes (Strassennamen = Top-Milestone), uebrige Milestones, Epics mit Fortschritt; je Eintrag Status der externen Verbindung (verknuepft mit JIRA-123, zuletzt abgeglichen).
- Schreiben: Ticket anlegen (Titel, Ziel, Kontext, DoD, Assignee - dieselben Pflichtfelder wie jaira create), Milestone anlegen, Zuweisen. 'Aus Jira uebernehmen': Key eingeben -> Felder vorbefuellt, Jira-Key als Feld 'external' am Ticket.
- Weitere Ideen, nach Wert sortiert:
  1. 'Braucht dich': Tickets in human/signoff mit offener Frage - der PM beantwortet sie im Browser (das ist der eigentliche Steuerhebel).
  2. Zuweisen per Drag auf einen Avatar.
  3. Haengt fest: Tickets ohne Aenderung seit N Tagen, claimed aber still.
  4. Ticket -> Code: Commit-Liste und Branch/PR-Link je Ticket (jaira kennt die Commits schon).
  5. Wochenbericht/Logbook-Zusammenfassung als Seite (gibt es als Rolle jaira-role-logbook-summary).
  6. Tickets bearbeiten (Ziel, Kontext, DoD) statt nur lesen.
  7. Mehrere Projekte in einer Seite umschalten (core/project).
  8. Schoene Startausgabe im Terminal (Wunsch aus Punkt 2) - billig.
  Spaeter/kein Muss: Burndown, Zeitschaetzung, Kommentare, Mentions.

4) WAS DRAUSSEN BLEIBT
- Kein Server, kein Konto, keine Datenbank: jaira web bleibt pro PM lokal auf Loopback. Kein Hosting fuer mehrere PMs.
- Keine Zwei-Wege-Dauer-Synchronisation mit Jira im Kern. Die ist der Teil, der Projekte wie paca gross macht (Konflikte, Webhooks, Feldzuordnungen).
- Keine Rechteverwaltung, keine Rollen ausser dem was git-Schreibrecht schon regelt.
- Keine Avatare von Gravatar o.ae. (Netzzugriff, Datenschutz).
- In eine optionale Komponente (eigenes Binary 'jaira-connect' oder Plugin nach 1YTBPX): Jira/GitHub/GitLab/YouTrack-Clients, Mapping-Matrix, Import/Export. Der Kern kennt nur ein Feld 'external: <anbieter>:<key>' am Ticket/Milestone/Epic. So bleibt das Binary klein und startet sofort.

FOLGE-TICKETS (Bauordnung)
1. DVPAPE (existiert): fester jaira-Name pro Projekt, lokal und gitignoriert; Voraussetzung fuer alles mit Personen.
2. Web schreibt Tickets: Ticket anlegen und Assignee setzen im Browser, ueber denselben Pfad wie jaira create/assign (Ref schreiben + push).
3. Web-Tab 'Aktiv': Avatare je Person mit Ticketliste, Farbe aus Lane-Stufe abgeleitet.
4. Web-Tab 'Milestones': Roadmap mit Swimlanes und Milestone-Liste, Milestone anlegen/zuweisen.
5. Web-Tabs-Geruest: Alle Tickets / Aktiv / Milestones & Epics / Schreiben, mit Erklaersatz Epic vs Milestone.
6. 'Braucht dich': offene Fragen aus human/signoff im Web beantworten.
7. Neu-zugewiesen-Markierung auf Board und Web (Punkt 5): Zeichen bis zum ersten Oeffnen, wer hat zugewiesen.
8. Web-Onboarding ohne CLI: Name, git-Pruefung mit Link/Pfadfeld, Projekt aus gefundenen Repos oder Clone-Link (core/project).
9. Lokaler Token-Speicher: ~/.jaira/secrets.json 0600, Einstellungsseite, nie in Antworten/Logs.
10. Feld 'external' an Ticket/Milestone/Epic im Kern (nur der Verweis, kein Client).
11. Optionale Komponente jaira-connect: Jira zuerst - Ticket aus Jira-Key importieren, Epic/Milestone als Jira Epic/Fix Version exportieren (einmalig, nicht dauernd).
12. Schoene Startausgabe von jaira web (klein, jederzeit einschiebbar).

ALTERNATIVEN
A) Alles in jaira web im Kern-Binary, inklusive Jira/GitHub/GitLab/YouTrack-Clients und Dauer-Sync. Kosten: geschaetzt mehrere tausend Zeilen Go plus vier API-Clients, Token-Handling im Kern, Binary waechst, jede API-Aenderung eines Anbieters wird ein jaira-Bug. Bricht 'kleiner als paca'.
B) PM-Funktionen nur ueber Jira/GitHub selbst: jaira spiegelt Tickets einseitig nach Jira, der PM arbeitet in Jira, jaira web bleibt wie heute. Kosten: kleinster Kern, aber der PM sieht keine Lanes/Agentenstufen, Zuweisungen in Jira kommen nicht zurueck, und 'nur git noetig' ist aufgegeben.

EMPFEHLUNG
Mittelweg: Kern bekommt nur, was mit git allein geht (Tickets 1-10: Name, Schreiben, Tabs, Fragen beantworten, Onboarding, lokaler Token-Speicher, external-Feld). Anbieter-Anbindung als optionale Komponente (11), zuerst nur Jira und nur einmaliger Import/Export, kein Dauer-Sync. Epics und Milestones bleiben getrennt; Default-Mapping Epic->Jira Epic, Milestone->Fix Version/GitHub Milestone. Start mit 1-3: danach kann ein PM schon zuweisen und sehen, wer woran sitzt.

OFFENE FRAGEN
1. Milestone = Epic: sollen PMs sie wirklich als dasselbe sehen, oder reicht ein Tab mit beiden und klarer Erklaerung? (Empfehlung: getrennt lassen.)
2. Wenn Jira: Milestone als Fix Version, als Sprint oder als Epic? Nutzt eure Firma Jira Premium (Initiativen)?
3. Darf der PM Tickets nur anlegen/zuweisen oder auch Lanes bewegen (z.B. signoff annehmen)? signoff ist heute 'ein Mensch' - ist der PM dieser Mensch?
4. Hat der PM Push-Recht auf das Repo? Sonst ist er nur Beobachter.
5. Punkt 5 bricht ab: 'siehst du es auf dem Board als ...' - als was? (Vorschlag: Markierung bis zum ersten Oeffnen.)
6. Sync mit Jira einmalig (Import/Export) oder dauerhaft in beide Richtungen? Dauerhaft ist der groesste Kostentreiber.
7. Onboarding im Browser statt im curl-Skript ok?
8. Token-Speicher: reicht zuerst die 0600-Datei, oder ist OS-Schluesselbund Pflicht von Anfang an?
- **2026-10-31 01:25 · berkc** — Antworten Maintainer 31.10.2026: (1) Milestone vs Epic: gleiche Form (Beschreibung + Liste + X von Y); Vorschlag Roads (KQQD1H) = Epics, Milestones liegen auf einer Road - ein Begriff weniger. (7) Onboarding im Browser: ja. (8) Token-Speicher: zuerst Datei mit Rechten 0600 (nur der eigene Benutzer darf lesen, keine Verschluesselung); Firmen: OS-Schluesselbund oder Umgebungsvariable aus einem Secret-Manager - spaeter. (6) Jira nur auf Anfrage, einmalig, wegen API-Quoten; kein Dauer-Sync. Mapping Epic/Milestone -> Jira soll im Web einstellbar sein (Default wie vorgeschlagen). Alles in einen Milestone. Offen: 2, 3, 4, 5.
