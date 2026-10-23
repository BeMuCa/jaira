---
id: 01M5GFB5KWABVZVWWDVNP5ABVR
title: "Hackathon-Boards: Ideen-Pipeline, Auswahl durch den Menschen und Builder-Pipeline, verbunden ueber den Teamlead"
status: ideas
ready: true
creator: berkc
assignee: berkc
goal: "Fuer einen Hackathon liegen fertige Boards bereit: aus dem Thema werden mehrere Umsetzungsideen mit Roadmap, ein Mensch waehlt aus, und ein Builder-Board baut UI und Backend nach festen Bausteinen."
context: |-
  Wunsch von berkc am 22.10.2026 fuer den AI Hub Hackathon Frankfurt Ende Oktober 2026: 'prepare by creating Boards for for example AI SAAS applications. Things that always have to be there: Auth, backend, UI (We should go for Svelte + svelte mcp, and define some collections for frontend things like buttons etc. so we have a fixed design and the ui doesnt bloat.), or even Idea creation board (so i can describe the topic of the hackathon and it creates me multiple sets of implementation ideas + roadmap - We could then even connect multiple boards -> idea pipeline -> hitl, to select which idea to go through or add or remove things -> Builder Pipeline (ui+backend) using the teamlead as the connector node between the boards.)'
  Bausteine:
  - Ideen-Board: Thema -> N Ideen mit Roadmap (Agent) -> HITL: auswaehlen, ergaenzen, streichen.
  - Builder-Board: Auth, Backend, UI (Svelte + Svelte-MCP, feste Komponentensammlung), Deploy.
  - Teamlead als Verbindung: legt aus der gewaehlten Idee Tickets auf dem Builder-Board an. Board-uebergreifende Weitergabe gibt es heute nicht; der Teamlead waehlt nur ein Board pro Ticket (8NHD8Q).
  Recherche zu Sponsoren, Themen und Produktart: in der Sitzung vom 22.10.2026 als Antwort geliefert (siehe Notiz).
  Verwandt: Y2C9QN (Board-Typen), Board-wie-Lanes-Ticket, PP24H1.
definition-of-done: Ein Ideen-Board und ein Builder-Board liegen als Board-Typen bereit
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T00:14:30Z
updated-at: 2026-10-23T21:04:35Z
updated-by: berkc
---

# Hackathon-Boards: Ideen-Pipeline, Auswahl durch den Menschen und Builder-Pipeline, verbunden ueber den Teamlead

## Definition of Done

- [ ] Ein Ideen-Board und ein Builder-Board liegen als Board-Typen bereit
- [ ] Der Teamlead legt aus einer im HITL gewaehlten Idee Tickets auf dem Builder-Board an (follows auf die Idee); einmal durchgespielt
- [ ] Eine Svelte-Komponentensammlung (Buttons, Formulare, Layout) ist festgelegt und der UI-Lane-Prompt verweist darauf

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-22 00:32 · berkc** — Recherche 22.10.2026 (Webquellen: ai-week-frankfurt.com, ai-week-frankfurt.com/hack, ai-hub-frankfurt.de/en/events-en/). Der 'AI Hub Hackathon' ist der AI Week Hack am Fr 30.10.2026 in The Squaire (Flughafen Frankfurt), Teil der AI Week Frankfurt 26.-30.10., veranstaltet von AI Hub Frankfurt. Ablauf: 09:00 Einlass, 10:00 Kick-off, 10:15-16:00 bauen, 16:00 Demos, 17:00 Finale, 18:00 Builders Night - also ~6 Stunden Bauzeit. Zwei Wege: (1) Open Build Floor - ohne festes Thema (Agent, Produkt, Automatisierung, Arbeit oder Spass), jedes Modell/Stack, Credits von Grokbot/Cursor, Glassflow, Wispr Flow, Exa, Netlify, Firecrawl, Render; getragen von AI Hub Frankfurt und Cursor Community Frankfurt; nur Ticket noetig. (2) OpenAI Codex Enterprise Hackathon - nur mit Bewerbung (~100 Teilnehmer), echte Business-Challenges von Condor (Airline), FraAlliance (Joint Venture Fraport + Lufthansa fuer den Betrieb am Flughafen FRA) und KfW IPEX-Bank (Export- und Projektfinanzierung); OpenAI-API-Credits, Codex; getragen von statworx und OpenAI. Welche Firma welche Challenge stellt und Preise: nicht veroeffentlicht. Sponsoren der Woche: STADA, Equinix, Commerzbank, DekaBank, Helaba, Merck, Dataiku, HP & NVIDIA, Metzler, TD SYNNEX, Heraeus, KfW IPEX-Bank, statworx u.a.; Donnerstag ist Banking & Finance. Wahrscheinliche Themen (abgeleitet, nicht bestaetigt): Luftfahrt/Flughafenbetrieb (FraAlliance arbeitet laut Presse 12/2025 an KI fuer Engpaesse im Terminal, Wartezeiten an der Sicherheitskontrolle, Passagierlenkung, digitalem Zwilling, Computer Vision fuer Turnaround-Ereignisse), Airline-Kundenservice und Stoerungsmanagement (Condor: Umbuchung bei Verspaetung, Kundenkommunikation), Finanzierung (KfW IPEX: Dokumentenpruefung, Risiko-/Due-Diligence-Zusammenfassungen, ESG-Pruefung von Projekten). Produktart in 6 Stunden: Agent oder interne Web-App mit Demo, nicht SaaS mit Abrechnung. Vorbereitung daraus: Builder-Board mit Auth-Stub, Backend (API + LLM-Aufruf + Dokumenten-/Datenimport), Svelte-UI mit fester Komponentensammlung, Deploy auf Netlify/Render (Credits vorhanden); Ideen-Board mit Vorlagen fuer die drei Branchen.
