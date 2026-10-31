---
id: 01M67K7GEZF694985EHTENCFVZ
title: "Neue Lane-Arten: eine Lane, die nur einen Skill startet, und Automations-Lanes (Artefakte, E-Mail, API-Aufrufe, Telegram senden/empfangen, Video-Upload)"
status: ideas
ready: true
creator: berkc
assignee: berkc
goal: "Eine Lane kann statt eines Prompts einen Skill ausfuehren oder eine Automation anstossen, so dass Boards ueber Code hinaus Arbeit erledigen - wie ein N8N-Ablauf."
context: |-
  Wunsch von berkc am 31.10.2026 (Punkt 25): 'einen neuen Board/Lane-Typ, wo man nur einen Skill reinhaengt und der dann gestartet wird, wenn das Ticket da landet. z.B. Ponytail-Skill soll etwas zum Code optimieren sein. Eine Lane, die nur Ponytail + vllt einen Mini-Prompt hat.' Weitere Ideen: Artefakt-Lanes (Bilder, HTML als Output), Lanes, die E-Mails versenden, bestimmte APIs aufrufen, eine Telegram-Nachricht an ihn senden oder eine von ihm erwarten, Videos per API hochladen. 'Ich denk in Richtung N8N-Automation.'
  Heute: eine Lane ist Prompt + Tier + Ein-/Ausgabevertrag; Artefakte fehlen (eigenes altes Ticket zu .jaira/artifacts/). Eine Lane, die auf eine Nachricht wartet, waere eine menschliche Lane mit externem Kanal.
  Brainstorm durch ein starkes Modell gewuenscht.
definition-of-done: "Brainstorm: Lane-Arten mit Format, Sicherheit (Secrets, Seiteneffekte), was im Binary bleibt und was ein Skill ist"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-30T23:44:56Z
updated-at: 2026-10-31T01:25:44Z
updated-by: berkc
---

# Neue Lane-Arten: eine Lane, die nur einen Skill startet, und Automations-Lanes (Artefakte, E-Mail, API-Aufrufe, Telegram senden/empfangen, Video-Upload)

## Definition of Done

- [ ] Brainstorm: Lane-Arten mit Format, Sicherheit (Secrets, Seiteneffekte), was im Binary bleibt und was ein Skill ist

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-30 23:49 · berkc** — Brainstorm (starkes Modell, 31.10.2026). Gelesen: .jaira/lanes/*.md, core/lane/lane.go (parse), core/lane/template.go, core/lane/share.go, core/hook/hook.go, core/settings/settings.go, jaira-worker SKILL.md, --for-lane JSON, FCMP17 (Artefakte), ATDG74 (Schwarm).

Ausgangslage heute:
- Eine Lane ist Frontmatter + Prompt-Body. lane.go verlangt fuer agentic: true einen Prompt-Body und ein model-tier.
- input-requires / output-produces meinen Ticket-Felder; das Gate verweigert den Move, solange ein Output fehlt.
- jaira fuehrt nichts aus: ein Agent liest --for-lane und arbeitet. Einzige Ausnahme: der Hook (core/hook) - ein Skript aus ~/.jaira/settings.json, pro Person, 5 s Timeout, bei move/claim. Das ist das Vorbild fuer alles unten.
- Geteilte Lanes (.jaira/shared/) sind absichtlich inert, weil committeter Prompt fremder Text ist (T-5-02). Dieselbe Regel muss fuer alles gelten, was etwas ausfuehrt.

(1) Skill-Lane
- Neues Feld: skill: ponytail. Der Prompt-Body wird optional und ist dann der Mini-Prompt, der zusaetzlich gilt.
- lane.go: agentic verlangt Prompt-Body ODER skill. model-tier bleibt Pflicht.
- jaira loest den Namen auf (.claude/skills/<name>/SKILL.md, dann ~/.claude/skills/<name>/SKILL.md) und gibt im --for-lane JSON skill und skill_path aus. So kann auch ein Agent ohne Skill-Tool die Datei lesen - die CLI bleibt fuer jeden bash-faehigen Agenten nutzbar.
- Nicht gefunden: jaira lanes / validate warnt; --for-lane setzt missing: [skill ponytail]; die Lane wird nicht still uebersprungen.
- Worker/Dispatcher/role-lane: ein Schritt mehr in "Read the lane's brief": steht skill im Payload, den Skill mit Ticket-ID + Mini-Prompt + bounded input starten, dann wie immer output-produces setzen und moven.
- Ausgabevertrag: unveraendert output-produces. Ein Skill kennt jaira-Felder nicht - der Mini-Prompt sagt, wohin das Ergebnis gehoert (z.B. "schreib die Befunde nach review-gaps"). Ohne output-produces faellt der Default auf note: mindestens eine jaira note der Lane, sonst Move verweigert.
- Kombinierbar mit rejects-to: ponytail findet etwas -> zurueck nach in-progress, wie optimize heute.

(2) Automations-Lanes
Was in jaira gehoert und was nicht:
- Ins Binary: nur das Benennen, das Ausloesen, die Quittung. Kein HTTP-Client, kein SMTP, kein Telegram-SDK, kein Secret-Speicher, kein Scheduler.
- Alles, was nach aussen spricht, ist ein Skript oder ein Skill der Person.
Format: die Lane nennt eine Absicht, nicht ein Programm: action: telegram-ask. Welches Skript dahinter steht, steht pro Person in ~/.jaira/settings.json unter actions: {"telegram-ask": "~/bin/tg-ask.sh"} - genau wie der Hook heute. Grund: ein committetes Skript, das jaira beim Move ausfuehrt, ist Code-Ausfuehrung per Pull Request auf dem Rechner jedes Teammitglieds. Die Repo-Datei sagt "hier wird eine Telegram-Frage gestellt", der Rechner der Person entscheidet, ob und womit.
Fehlt die Zuordnung bei einer Person: die Lane ist fuer sie nicht ausfuehrbar; next --per-lane zeigt sie als "nicht eingerichtet", das Ticket wartet (nicht still durchreichen - eine E-Mail, die nie rausging, darf nicht wie gesendet aussehen).
Ausfuehren: kein Daemon, also loest ein Agent oder eine Person aus: jaira run <id>. jaira ruft das Skript mit Umgebungsvariablen wie beim Hook (JAIRA_TICKET, JAIRA_LANE, ...) plus nur den Feldern aus input-requires (JAIRA_FIELD_GOAL usw.) und den Artefakt-Pfaden. Nur was die Lane deklariert, verlaesst den Rechner.
Secrets:
- Nie im Repo, nie in der Lane-Datei, nie im Ticket, nie in einer Note. jaira kennt keinen Token und reicht keinen weiter.
- Das Skript liest seine eigenen Secrets (Umgebung, ~/.config, pass, OS-Keychain). jaira kann das nicht pruefen, nur dokumentieren.
- jaira run gibt stdout des Skripts nicht ungefiltert ins Ticket: nur eine Zeile JAIRA_RESULT=... bzw. eine Datei als Artefakt. Sonst landet ein versehentlich geloggter Token im Ticket und damit in git.
Seiteneffekte und Idempotenz:
- Quittung auf dem Ticket: action-<lane>: pending <key> vor dem Aufruf, action-<lane>: done <key> <ergebnis> danach. key = ticket-id + lane + Durchlauf-Nummer (wie oft das Ticket die Lane betreten hat).
- jaira run verweigert (Exit 3), wenn fuer denselben key schon done steht; --again nur mit Absicht.
- pending ohne done (Absturz zwischen Senden und Quittung): jaira weiss nicht, ob gesendet wurde -> keine automatische Wiederholung, sondern Frage an die Person. Ehrlich: genau-einmal ist ohne Server unmoeglich; jaira garantiert hoechstens-einmal pro key, das Skript bekommt JAIRA_IDEMPOTENCY_KEY und kann ihn beim Anbieter mitschicken (Message-ID bei E-Mail, Idempotency-Key bei APIs).
- Ein Ruecksprung (critique -> in-progress -> wieder hier) ist ein neuer Durchlauf und darf wieder senden - das ist gewollt, z.B. eine neue Demo.
- Die Quittung ist ein Ticket-Feld, wird also mit git geteilt: ein Teammitglied sieht, dass schon gesendet wurde, und sendet nicht doppelt.
Warte-Lane (Telegram-Antwort) gegen heutige human-Lane:
- Heute: human verlangt --question und requires-human-exit; die Person liest im Board und moved selbst.
- Vorschlag: keine neue Lane-Art, sondern eine human-Lane mit action: telegram-ask. Beim Betreten (oder per jaira run) schickt das Skript die offene Frage.
- Rueckweg: jaira pollt nicht. Ein Listener der Person (Telegram-Bot, cron, systemd) ruft auf ihrem Rechner jaira note <id> "<antwort>" und jaira move <id> --to <rejects-to> auf. Damit handelt die Person selbst (ihr Rechner, ihre Identitaet) und requires-human-exit bleibt wahr.
- Pflicht im Beispielskript: nur Nachrichten aus der eigenen chat-id annehmen, sonst kann jeder, der den Bot findet, Tickets bewegen.
- Unterschied zur heutigen human-Lane also nur: wie die Frage hinausgeht und von wo die Antwort kommt. Der Gate-Vertrag bleibt gleich.

(3) Artefakte (baut auf FCMP17 auf)
- Ablage: .jaira/artifacts/<handle>/<ticket-kurz>/<lane>-<name>.<ext> fuer geteilte, .jaira/artifacts/local/... gitignoriert fuer grosse oder private.
- Die Lane deklariert output-artifacts: [demo.html]; das Gate prueft beim Move, dass die Datei existiert.
- Die Folge-Lane deklariert input-artifacts: [demo.html]; --for-lane gibt Pfade, nie Inhalte (bounded input bleibt bounded).
- Groessengrenze fuer geteilte Artefakte (Vorschlag 1 MB, Text/HTML/PNG/SVG). Videos und alles Grosse nur local/ oder extern; dann steht auf dem Ticket nur ein Link-Feld (artifact-url-<name>).
- Geteiltes Artefakt = fremder Text, der als Eingabe in ein Modell laeuft; dieselbe Warnung wie bei shared/ Lanes.
- Eine action-Lane bekommt Artefakt-Pfade als JAIRA_ARTIFACT_<NAME>; so laedt ein Upload-Skript das Video hoch, ohne dass jaira je davon weiss.

(4) Weitere Lane-Ideen, nach Nutzen sortiert
1. Skill-Lanes fuer vorhandene Skills: ponytail/simplify nach in-progress, security-review mit requires-option. Hoechster Nutzen, fast keine Kosten.
2. Research-Lane vor pre-process mit skill: jaira-role-research (requires-option: research) - die Rolle gibt es schon.
3. Telegram/ntfy-Frage aus der human-Lane (action + Listener) - die Person kann unterwegs entscheiden; heute schlaeft ein Ticket, bis jemand ins Board schaut.
4. Demo-Artefakt-Lane vor signoff: Screenshot oder HTML-Seite, die zeigt, was gebaut wurde; signoff liest sie statt des Diffs. Passt zu jaira-role-acceptance.
5. Release-Notes-Lane (projektspezifisch): schreibt core/release/notes/<TICKET>.md, wenn die Aenderung client-facing ist - heute vergisst man das am leichtesten.
6. PR-Lane mit skill: jaira-role-pr: pusht den Branch und gibt die gh-Zeile zurueck (oeffnet nie selbst, wie CLAUDE.md verlangt).
7. CI-Warte-Lane: wartet auf gruene GitHub Actions; der "Listener" ist ein gh run watch, der jaira move ausfuehrt. Selbes Muster wie die Telegram-Warte-Lane.
8. Docs-Lane: aktualisiert docs/ und README fuer client-facing Aenderungen, requires-option: docs.
9. Benachrichtigung bei Ankunft in signoff (action: notify). Ueberschneidet sich mit dem Hook - nur bauen, wenn der Hook pro Lane filtern koennen soll.
10. Schwarm-Lane (ATDG74) als Variante der Skill-Lane: swarm: 5 an critique/review. Nur lesende Lanes. Eigenes Ticket, nicht hier.
11. Upload-Lane (Video/Bild per API) - geringster Nutzen fuer ein Code-Board; geht mit action + Artefakt ohne Extra-Code, also kein eigenes Feature.

(5) Was draussen bleibt, damit es kleiner als paca bleibt
- Kein Daemon, kein Scheduler, kein Polling in jaira. Ausloeser ist immer eine Session, ein jaira run oder das Skript der Person.
- Keine Netzwerk-Clients und keine Anbieter-Integrationen im Binary (kein SMTP, Telegram, YouTube, OpenAI).
- Kein Secret-Speicher, keine verschluesselten Felder.
- Keine Ablauf-Sprache: keine Bedingungen, Verzweigungen, Schleifen, Variablen in Lane-Dateien ausser dem, was rejects-to schon kann. Das waere N8N nachgebaut.
- Kein grafischer Ablauf-Editor, kein Marketplace fuer Actions.
- Keine Skripte aus dem Repo, die jaira ausfuehrt. Repo nennt, Person fuehrt aus.
- Keine Binaerdateien ueber der Grenze in git.

Format-Skizze einer Lane-Datei
Skill-Lane:
---
id: ponytail
after: in-progress
agentic: true
model-tier: strong
skill: ponytail
rejects-to: in-progress
input-requires: [goal, diff, notes]
output-produces: [review-gaps]
---
# Prompt
Nur den Diff dieses Tickets optimieren. Befunde nach review-gaps, je Zeile einer. Nichts gefunden: review-gaps="keine".

Automations-Lane (Frage per Telegram):
---
id: ask-phone
after: in-progress
agentic: false
requires-question: true
requires-human-exit: true
action: telegram-ask
input-requires: [title, question]
---
Dazu pro Person in ~/.jaira/settings.json: "actions": {"telegram-ask": "~/bin/jaira-tg-ask.sh"}

Artefakt-Lane:
---
id: demo
after: testing
agentic: true
model-tier: cheap
output-artifacts: [demo.html]
---
# Prompt
Baue eine HTML-Seite, die zeigt, was dieses Ticket geaendert hat ...

Folge-Tickets in Bau-Reihenfolge
1. skill-Feld: lane.go (Prompt ODER skill), Aufloesung skill_path, --for-lane-Feld, lanes template, validate-Warnung, jaira-worker/dispatcher/role-lane lesen es, Release-Note. Klein, sofort nuetzlich, keine Sicherheitsfrage (Skill liegt ohnehin lokal).
2. Artefakte = FCMP17 mit den Ergaenzungen oben (Pfadschema pro Ticket, Groessengrenze, local/ vs. Handle-Ordner, Pfade statt Inhalte).
3. action-Feld + actions-Zuordnung in ~/.jaira/settings.json + jaira run <id> + Quittung (pending/done, key, Exit 3 bei Wiederholung) + "nicht eingerichtet" in next --per-lane.
4. Warte-Muster: human-Lane mit action; Doku + Beispielskript (Telegram, chat-id-Pruefung) unter docs/, nicht im Binary. Kein neuer jaira-Code ausser vielleicht action beim Betreten automatisch ausloesen (offene Frage).
5. Beispiel-Lanes (ponytail, research, demo, release-notes) als shared/-Vorlagen zum Adoptieren.

Zwei Alternativen mit Kosten
A) "Alles ist ein Skill": nur Schritt 1 + 2. Automation (Mail, Telegram, Upload) macht ein Skill, der eigene Skripte aufruft; jaira weiss nichts davon. Kosten: klein (ein Feld, ein JSON-Feld, Skill-Text, Artefakte aus FCMP17). Schwaeche: keine Quittung - ob die Mail doppelt rausgeht, haengt an der Disziplin des Agenten; ein abgebrochener Lauf wiederholt sie. Kein Lauf ohne Agent (ein Agent muss die Lane oeffnen, Tokens fuer eine Mail).
B) "Action-Lanes im Binary" (Schritte 1-4): zusaetzlich action-Feld, Zuordnung pro Person, jaira run mit Quittung. Kosten: mittel - ein neuer Befehl, Settings-Erweiterung, Gate-Regel, Tests fuer Wiederholung/Absturz, Doku; geschaetzt ein Vielfaches von A. Gewinn: hoechstens-einmal pro Durchlauf, laeuft ohne Modell (cron, Listener), Quittung im git sichtbar fuer das Team.
(Verworfen: eine N8N-artige Ablauf-Engine in jaira - widerspricht git-only, kein Daemon, kleiner als paca.)

Empfehlung
A jetzt bauen (Schritt 1, dann FCMP17). Ponytail als erste echte Skill-Lane benutzen. B erst, wenn eine konkrete Automation (vermutlich die Telegram-Frage) laeuft und eine doppelte oder verlorene Nachricht wirklich wehgetan hat - dann Schritt 3 in kleinster Form (Quittung + jaira run, sonst nichts). Die Warte-Lane braucht in beiden Faellen keinen neuen Lane-Typ, nur ein Skript der Person, das jaira note + jaira move aufruft.

Offene Fragen an berkc
1. Was ist "Ponytail" genau - ein vorhandener Skill (bei dir nicht in .claude/skills oder ~/.claude/skills gefunden) oder einer, der erst geschrieben werden soll? Soll er nur Befunde liefern (wie optimize, zurueck nach in-progress) oder selbst Code aendern?
2. Darf eine Lane eine action beim Betreten automatisch ausloesen (jaira move ruft das Skript, wie heute der Hook), oder immer nur explizit per jaira run? Automatisch ist bequemer, aber dann sendet jeder Move eines Teammitglieds mit derselben Zuordnung.
3. Fehlt einer Person die actions-Zuordnung: Ticket wartet (Vorschlag) oder Lane wird fuer sie uebersprungen?
4. Sollen Artefakte ueberhaupt in git (Handle-Ordner) oder grundsaetzlich nur local/ + Link? Groessengrenze 1 MB ok?
5. Telegram-Antwort: reicht "Antwort wird Note, Ticket geht nach rejects-to", oder soll eine Antwort wie "ok" auch weiter nach vorne moven duerfen (dann entscheidet ein Bot, was heute nur eine Person darf)?
6. A oder B - oder A jetzt und B als Ticket im Backlog?
- **2026-10-31 01:25 · berkc** — Antworten Maintainer 31.10.2026: Ponytail = github.com/DietrichGebert/ponytail, ein Skill 'lazy senior dev' (minimaler Code, /ponytail-review, -audit, -debt) - passt als Skill-Lane neben optimize. Skript pro Person: der Maintainer will ein mitgeliefertes Skript (z.B. Telegram-API-Aufruf), das der Agent nutzt - keine Zuordnung pro Person noetig, solange das Skript keine Secrets enthaelt (Token aus lokalem Speicher, CMN1YT). Quittung gegen Doppelsenden: entfaellt vorerst, der Text auf dem Ticket soll ohnehin weg (eigenes Ticket). Entscheidung: Variante A (alles ist ein Skill, der Agent ruft die API).
