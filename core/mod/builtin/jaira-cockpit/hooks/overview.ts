// The board at a glance, for the Overview pane: what is in work, what could
// be started next, what was put back, and the story of the last days — what
// was created and what was put back on which day, and by whom. Read off the
// board alone (`jaira lanes`, `jaira list`, `jaira next --all`), never out of
// a session: so it is the same after a night, a restart or another session.
// register.tsx runs the commands; this file only reads what they printed.

import { isPersons, oneLine, parse, PERSONS_TABS, type Agent, type Lane } from './board.ts'
import type { Tab } from './sessions.ts'

export type Row = {
  id: string
  lane: string
  title: string
  // Who holds it, created it or put it back; '' for nobody.
  who: string
  // Why a blocked ticket waits; '' elsewhere.
  reason: string
}

export type Event = Row & { kind: 'created' | 'parked'; at: string }
export type Day = { day: string; events: Event[] }

export type Overview = { working: Row[]; queue: Row[]; parked: Row[]; history: Day[] }

// How far back the history reaches: a weekend and a day off on top of the
// night a person slept through.
export const HISTORY_DAYS = 7

// The lanes a ticket is put back into. No lane says so of itself.
const PARKED = new Set(['backlog', 'blocked'])

type Ticket = {
  handle: string
  title: string
  status: string
  assignee?: string
  creator?: string
  created_at?: string
  updated_at?: string
  updated_by?: string
  blocked_reason?: string
}

const two = (n: number) => String(n).padStart(2, '0')

// The local calendar day of an instant, as YYYY-MM-DD.
export function dayOf(at: string): string {
  const d = new Date(at)
  return `${d.getFullYear()}-${two(d.getMonth() + 1)}-${two(d.getDate())}`
}

// The local time of day of an instant, as HH:MM.
export function clockOf(at: string): string {
  const d = new Date(at)
  return `${two(d.getHours())}:${two(d.getMinutes())}`
}

// A missing or unreadable time is none: never within the history.
const time = (at: string) => (Number.isNaN(Date.parse(at)) ? -Infinity : Date.parse(at))

// In work: a ticket in an agent's lane that jaira next does not offer and
// somebody is assigned to; one next leaves out for want of readiness or for
// an open blocked_by, with nobody on it, is in nobody's work. In the queue:
// what jaira next offers, short of what was put back and what waits for a
// person (the band shows those). Put back: every ticket in backlog or
// blocked. Each list the last touched first. The history: every ticket
// created in the last HISTORY_DAYS, and every one put back since then on a
// later day than it was created — jaira keeps no time a ticket changed
// lanes, so its last change stands for it.
export function overviewOf(lanesJson: string, listJson: string, nextJson: string, now: number): Overview {
  const lanes = new Map((parse<{ lanes: Lane[] }>(lanesJson)?.lanes ?? []).map(l => [l.id, l]))
  const tickets = parse<{ tickets: Ticket[] }>(listJson)?.tickets ?? []
  const startable = new Set((parse<{ tickets: { handle: string }[] }>(nextJson)?.tickets ?? []).map(t => t.handle))
  const row = (t: Ticket, who: string | undefined): Row => ({
    id: t.handle,
    lane: t.status,
    title: oneLine(t.title),
    who: oneLine(who),
    reason: t.status === 'blocked' ? oneLine(t.blocked_reason) : '',
  })
  const working: Row[] = []
  const queue: Row[] = []
  const parked: Row[] = []
  const newestFirst = (a: string, b: string) => time(b) - time(a)
  for (const t of [...tickets].sort((a, b) => newestFirst(a.updated_at ?? '', b.updated_at ?? ''))) {
    const lane = lanes.get(t.status)
    if (PARKED.has(t.status)) parked.push(row(t, t.updated_by || t.creator))
    else if (lane === undefined || lane.terminal === true || isPersons(lane)) continue
    else if (startable.has(t.handle)) queue.push(row(t, t.assignee))
    else if (lane.agentic && (t.assignee ?? '').trim() !== '') working.push(row(t, t.assignee))
  }
  const since = now - HISTORY_DAYS * 86_400_000
  const events: Event[] = []
  for (const t of tickets) {
    const created = t.created_at ?? ''
    if (time(created) >= since) events.push({ ...row(t, t.creator), kind: 'created', at: created })
    const updated = t.updated_at ?? ''
    if (PARKED.has(t.status) && time(updated) >= since && dayOf(updated) !== dayOf(created)) {
      events.push({ ...row(t, t.updated_by || t.creator), kind: 'parked', at: updated })
    }
  }
  events.sort((a, b) => newestFirst(a.at, b.at))
  const history: Day[] = []
  for (const e of events) {
    const day = dayOf(e.at)
    if (history.at(-1)?.day !== day) history.push({ day, events: [] })
    history.at(-1)!.events.push(e)
  }
  return { working, queue, parked, history }
}

// ── the Workers tab: every agent on the board as a figure at a station ────

// How an agent reads: at work, resting with its tab open, or standing with
// a question its ticket carries for a person.
export type WorkerState = 'working' | 'idle' | 'question'
// One figure: an agent, or a ticket waiting on a person's answer with no
// agent on it.
export type Walker = { id: string; state: WorkerState }
// The stations in board order, and who stands at each, by the lane of the
// ticket the figure works on.
export type Workers = { stations: string[]; at: Record<string, Walker[]> }

// A subagent at work: running a turn, about to, or held on its own work.
const RUNNING = new Set(['pending', 'running', 'waiting'])

// The stations are the lanes something happens in: every agent's lane and
// every person's, in board order; backlog, todo, blocked and done are not.
// A Herdr tab spawn.sh labelled `<id>/<kind>` is an agent, at work while
// Herdr says so and resting otherwise; a subagent of this session whose
// description or name carries the id is one too, resting between turns and
// gone once it ended. A Discuss or Explain tab is the person's own chat, not
// an agent. A ticket whose question is open stands with it, agent or not.
// An agent on a ticket in no station's lane is not drawn.
export function workersOf(lanesJson: string, listJson: string, tabs: readonly Tab[], agents: readonly Agent[]): Workers {
  const stations = (parse<{ lanes: Lane[] }>(lanesJson)?.lanes ?? [])
    .filter(l => l.terminal !== true && (l.agentic || isPersons(l)))
    .map(l => l.id)
  const tickets = parse<{ tickets: (Ticket & { question?: string })[] }>(listJson)?.tickets ?? []
  const at: Record<string, Walker[]> = Object.fromEntries(stations.map(s => [s, []]))
  for (const t of tickets) {
    const here = at[t.status]
    if (here === undefined) continue
    const asks = oneLine(t.question) !== ''
    const state = (working: boolean): WorkerState => (asks ? 'question' : working ? 'working' : 'idle')
    const walkers: Walker[] = []
    for (const tab of tabs) {
      const [id = '', kind = ''] = tab.label.split('/')
      if (id === t.handle && kind !== '' && !PERSONS_TABS.has(kind)) walkers.push({ id, state: state(tab.agent_status === 'working') })
    }
    for (const a of agents) {
      if (!`${a.description} ${a.name ?? ''}`.includes(t.handle)) continue
      if (RUNNING.has(a.status) || a.status === 'idle') walkers.push({ id: t.handle, state: state(RUNNING.has(a.status)) })
    }
    if (walkers.length === 0 && asks) walkers.push({ id: t.handle, state: 'question' })
    here.push(...walkers)
  }
  return { stations, at }
}

// The colours a person's name in the Log and a ticket's figures wear. ANSI
// names, apart from the group colours, which the terminal's own scheme
// keeps readable.
export const COLORS = ['cyan', 'magenta', 'green', 'cyanBright', 'magentaBright', 'greenBright']

// A name's slot among `n`, derived from the name and nothing else.
export function slotOf(name: string, n: number): number {
  let h = 0x811c9dc5
  for (const c of name) h = Math.imul(h ^ c.codePointAt(0)!, 0x01000193)
  return (h >>> 0) % n
}

// Each ticket's colour. One that already has one keeps it while it stays on
// the board; a new one starts at its name's slot and moves on to the next
// colour no other ticket wears. Only with more tickets than colours do two
// share one, and then by name alone.
export function colorsOf(ids: readonly string[], kept: Readonly<Record<string, string>>): Record<string, string> {
  const colors: Record<string, string> = {}
  for (const id of ids) if (kept[id] !== undefined && COLORS.includes(kept[id]!)) colors[id] = kept[id]!
  for (const id of [...ids].sort()) {
    if (colors[id] !== undefined) continue
    const used = new Set(Object.values(colors))
    const start = slotOf(id, COLORS.length)
    const free = COLORS.map((_, i) => COLORS[(start + i) % COLORS.length]!).find(c => !used.has(c))
    colors[id] = free ?? COLORS[start]!
  }
  return colors
}

// Where each ticket's figures stand, as a station's index: the station its
// lane is for a ticket not seen before, and one station nearer it for one
// that walks — so a ticket that changed lanes crosses every station between.
export function stepOf(where: Readonly<Record<string, number>>, workers: Workers): Record<string, number> {
  const next: Record<string, number> = {}
  workers.stations.forEach((station, i) => {
    for (const { id } of workers.at[station] ?? []) {
      const from = where[id] ?? i
      next[id] = from === i ? i : from + Math.sign(i - from)
    }
  })
  return next
}

// A working figure turns through these, one each frame.
export const SPINNER = ['✻', '✳', '✶', '✢']
