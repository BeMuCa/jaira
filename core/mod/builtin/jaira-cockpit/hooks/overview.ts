// The board at a glance, for the Overview pane: what is in work, what could
// be started next, what was put back, and the story of the last days — what
// was created and what was put back on which day, and by whom. Read off the
// board alone (`jaira lanes`, `jaira list`, `jaira next --all`), never out of
// a session: so it is the same after a night, a restart or another session.
// register.tsx runs the commands; this file only reads what they printed.

import { isPersons, oneLine, parse, PERSONS_TABS, WORKING, type Agent, type Lane } from './board.ts'
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
  question?: string
  next_lane?: string
  plan_items?: Item[]
  dod_items?: Item[]
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

// ── the Workers tab: every agent on the board on a track of its own ──────

// How an agent reads: at work, resting with its tab open, or standing with
// a question its ticket carries for a person.
export type WorkerState = 'working' | 'idle' | 'question'
// One figure: an agent, or a ticket waiting on a person's answer with no
// agent on it, with what its track shows — the lane the ticket is in, the
// lane it goes to next, its title, and how far it has come in its lane.
export type Walker = { id: string; state: WorkerState; lane: string; next: string; title: string; progress: number }

type Item = { done?: boolean; state?: string }

// How far a ticket has come in its lane, 0 to 1: its plan ticked off, or
// where it has no plan its definition of done; nothing to count, none. An
// item superseded counts neither way.
export function progressOf(plan: readonly Item[] | undefined, dod: readonly Item[] | undefined): number {
  for (const items of [plan, dod]) {
    const live = (items ?? []).filter(i => i.state !== 'superseded')
    if (live.length > 0) return live.filter(i => i.done === true).length / live.length
  }
  return 0
}

// Every figure, in board order of its ticket's lane, then by handle. The
// lanes a figure stands in are the ones something happens in: every agent's
// lane and every person's; backlog, todo, blocked and done are not.
// A Herdr tab spawn.sh labelled `<id>/<kind>` is an agent, at work while
// Herdr says so and resting otherwise; a subagent of this session whose
// description or name carries the id is one too, resting between turns and
// gone once it ended. A Discuss or Explain tab is the person's own chat, not
// an agent. A ticket whose question is open stands with it, agent or not.
export function workersOf(lanesJson: string, listJson: string, tabs: readonly Tab[], agents: readonly Agent[]): Walker[] {
  const order = (parse<{ lanes: Lane[] }>(lanesJson)?.lanes ?? [])
    .filter(l => l.terminal !== true && (l.agentic || isPersons(l)))
    .map(l => l.id)
  const tickets = (parse<{ tickets: Ticket[] }>(listJson)?.tickets ?? [])
    .filter(t => order.includes(t.status))
    .sort((a, b) => order.indexOf(a.status) - order.indexOf(b.status) || a.handle.localeCompare(b.handle))
  const walkers: Walker[] = []
  for (const t of tickets) {
    const asks = oneLine(t.question) !== ''
    const figure = (working: boolean): Walker => ({
      id: t.handle,
      state: asks ? 'question' : working ? 'working' : 'idle',
      lane: t.status,
      next: t.next_lane ?? '',
      title: oneLine(t.title),
      progress: progressOf(t.plan_items, t.dod_items),
    })
    const before = walkers.length
    for (const tab of tabs) {
      const [id = '', kind = ''] = tab.label.split('/')
      if (id === t.handle && kind !== '' && !PERSONS_TABS.has(kind)) walkers.push(figure(tab.agent_status === 'working'))
    }
    for (const a of agents) {
      if (`${a.description} ${a.name ?? ''}`.includes(t.handle) && (WORKING.has(a.status) || a.status === 'idle')) walkers.push(figure(WORKING.has(a.status)))
    }
    if (walkers.length === before && asks) walkers.push(figure(false))
  }
  return walkers
}

// The colours a person's name in the Log and a ticket's figures wear. ANSI
// names, apart from the group colours, which the terminal's own scheme
// keeps readable.
export const COLORS = ['cyan', 'magenta', 'green', 'cyanBright', 'magentaBright', 'greenBright']

// The hats a ticket's figures wear, one row of seven cells each.
export const HATS: Record<string, string> = {
  crown: '▗▙▟▙▟▙▖',
  topHat: '▄▐███▌▄',
  cap: '▗▄███▄▄',
  helmet: '▗█████▖',
  beanie: '▗▄▓▓▓▄▖',
  propeller: '  ═╪═  ',
  wizard: '   ▲   ',
  bow: '  ▶●◀  ',
}

// A name's slot among `n`, derived from the name and nothing else.
export function slotOf(name: string, n: number): number {
  let h = 0x811c9dc5
  for (const c of name) h = Math.imul(h ^ c.codePointAt(0)!, 0x01000193)
  return (h >>> 0) % n
}

// Each ticket's pick from `choices`, a colour or a hat. One that already
// has one keeps it while it stays on the board; a new one starts at its
// name's slot and moves on to the next choice no other ticket has. Only
// with more tickets than choices do two share one, and then by name alone.
export function pickOf(ids: readonly string[], kept: Readonly<Record<string, string>>, choices: readonly string[]): Record<string, string> {
  const picks: Record<string, string> = {}
  for (const id of ids) if (kept[id] !== undefined && choices.includes(kept[id]!)) picks[id] = kept[id]!
  for (const id of [...ids].sort()) {
    if (picks[id] !== undefined) continue
    const used = new Set(Object.values(picks))
    const start = slotOf(id, choices.length)
    const free = choices.map((_, i) => choices[(start + i) % choices.length]!).find(c => !used.has(c))
    picks[id] = free ?? choices[start]!
  }
  return picks
}

// Where a ticket's figures stand on their track: the lane whose wall is on
// the left, and how far towards the right wall, 0 to 1.
export type Spot = { lane: string; x: number }

// How far a figure walks in one frame, as a share of its track.
const WALK_STEP = 1 / 40

// One frame on from `where`. A ticket not seen before stands where its
// progress puts it. One whose progress moved walks towards it. One whose
// ticket changed lanes walks on to the right wall first; there the walls
// move along — the old right one is now the left — and it starts again at
// the left. A ticket with no figure any more drops out.
export function stepOf(where: Readonly<Record<string, Spot>>, walkers: readonly Walker[]): Record<string, Spot> {
  const next: Record<string, Spot> = {}
  for (const w of walkers) {
    if (next[w.id] !== undefined) continue
    const at = where[w.id]
    if (at === undefined) next[w.id] = { lane: w.lane, x: w.progress }
    else if (at.lane !== w.lane) next[w.id] = at.x >= 1 ? { lane: w.lane, x: 0 } : { lane: at.lane, x: at.x + WALK_STEP > 1 - WALK_STEP / 2 ? 1 : at.x + WALK_STEP }
    else next[w.id] = { lane: w.lane, x: Math.abs(w.progress - at.x) <= WALK_STEP ? w.progress : at.x + Math.sign(w.progress - at.x) * WALK_STEP }
  }
  return next
}

// A working figure's spark turns through these, one each frame.
export const SPINNER = ['✻', '✳', '✶', '✢']
