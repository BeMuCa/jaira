// The board at a glance, for the Overview pane: what is in work, what could
// be started next, what was put back, and the story of the last days — what
// was created and what was put back on which day, and by whom. Read off the
// board alone (`jaira lanes`, `jaira list`, `jaira next --all`), never out of
// a session: so it is the same after a night, a restart or another session.
// register.tsx runs the commands; this file only reads what they printed.

import { isPersons, oneLine, openOf, parse, PERSONS_TABS, WORKING, type Agent, type Item, type Lane } from './board.ts'
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

// `lanes`: the board's lanes in its own order, so the Log's lane filter
// reads as the board does.
export type Overview = { working: Row[]; queue: Row[]; parked: Row[]; history: Day[]; lanes: string[] }

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
  return { working, queue, parked, history, lanes: [...lanes.keys()] }
}

// ── the Workers tab: every agent on the board at one stone ───────────────

// How an agent reads: at work, resting with its tab open, resting while
// its ticket carries a question for a person, or held at an approval dialog
// in its tab (Herdr's `blocked`) — which stops it whatever else is true.
export type WorkerState = 'working' | 'idle' | 'question' | 'blocked'
// One agent: the ticket it carries, what kind of agent it is (the part of
// its Herdr tab's label after the slash, or `subagent`), the lane its ticket
// is in, its title, and how many definition-of-done items that ticket still
// has open — what it adds to the stone.
export type Worker = { id: string; role: string; state: WorkerState; lane: string; title: string; open: number }

// What a subagent is called on the Workers tab, having no tab label.
export const SUBAGENT = 'subagent'

// Every agent on the board, in board order of its ticket's lane, then by
// handle and kind. The lanes an agent is counted in are the ones something
// happens in: every agent's lane and every person's; backlog, todo, blocked
// and done are not. A Herdr tab spawn.sh labelled `<id>/<kind>` is an agent,
// at work while Herdr says so and resting otherwise; a subagent of this
// session whose description or name carries the id is one too, resting
// between turns and gone once it ended. A Discuss or Explain tab is the
// person's own chat, not an agent. A ticket with an open question and no
// agent is nobody at work: the band already lists it as waiting.
export function workersOf(lanesJson: string, listJson: string, tabs: readonly Tab[], agents: readonly Agent[]): Worker[] {
  const order = (parse<{ lanes: Lane[] }>(lanesJson)?.lanes ?? [])
    .filter(l => l.terminal !== true && (l.agentic || isPersons(l)))
    .map(l => l.id)
  const tickets = (parse<{ tickets: Ticket[] }>(listJson)?.tickets ?? [])
    .filter(t => order.includes(t.status))
    .sort((a, b) => order.indexOf(a.status) - order.indexOf(b.status) || a.handle.localeCompare(b.handle))
  const workers: Worker[] = []
  for (const t of tickets) {
    const asks = oneLine(t.question) !== ''
    const worker = (role: string, status: string): Worker => ({
      id: t.handle,
      role,
      state: status === 'blocked' ? 'blocked' : status === 'working' ? 'working' : asks ? 'question' : 'idle',
      lane: t.status,
      title: oneLine(t.title),
      open: openOf(t.dod_items),
    })
    const mine: Worker[] = []
    for (const tab of tabs) {
      const [id = '', kind = ''] = tab.label.split('/')
      if (id === t.handle && kind !== '' && !PERSONS_TABS.has(kind)) mine.push(worker(kind, tab.agent_status ?? ''))
    }
    for (const a of agents) {
      if (`${a.description} ${a.name ?? ''}`.includes(t.handle) && (WORKING.has(a.status) || a.status === 'idle')) mine.push(worker(SUBAGENT, WORKING.has(a.status) ? 'working' : 'idle'))
    }
    workers.push(...mine.sort((a, b) => a.role.localeCompare(b.role)))
  }
  return workers
}

// One agent's key: its ticket and its kind. Two agents of one kind on one
// ticket share it, and with it the time they were first seen.
export const workerKey = (w: Pick<Worker, 'id' | 'role'>) => `${w.id}/${w.role}`

// The stone's open work: every ticket with an agent at it counted once,
// however many agents stand at it.
export function stoneOf(workers: readonly Worker[]): number {
  const open = new Map(workers.map(w => [w.id, w.open]))
  return [...open.values()].reduce((a, b) => a + b, 0)
}

// How big the stone is drawn for `open` items, in cells, within a box of
// `maxW` × `maxH`: nothing open leaves a pebble; then it grows with the
// square root of what is open — twice the work reads larger, not twice as
// wide — and never past the box. Never smaller for more open work.
export function stoneSize(open: number, maxW: number, maxH: number): { w: number; h: number } {
  const fit = (w: number, h: number) => ({ w: Math.max(1, Math.min(w, maxW)), h: Math.max(1, Math.min(h, maxH)) })
  if (open <= 0) return fit(3, 1)
  const s = Math.sqrt(open)
  return fit(4 + Math.round(4 * s), 1 + Math.round(s))
}

// How long an agent has been seen, in plain units: under a minute, minutes,
// then hours and minutes.
export function durationOf(ms: number): string {
  const min = Math.floor(Math.max(0, ms) / 60_000)
  if (min < 1) return '<1m'
  if (min < 60) return `${min}m`
  return `${Math.floor(min / 60)}h${two(min % 60)}`
}

// One filter row of the Log, people or lanes: the values chosen and the values
// left out. Both empty: the row filters nothing.
export type Filter = { only: string[]; not: string[] }
export const NO_FILTER: Filter = { only: [], not: [] }

// A chip's next state on a press: off, then chosen, then left out, then off.
export function cycled(f: Filter, x: string): Filter {
  if (f.only.includes(x)) return { only: f.only.filter(y => y !== x), not: [...f.not, x] }
  if (f.not.includes(x)) return { only: f.only, not: f.not.filter(y => y !== x) }
  return { only: [...f.only, x], not: f.not }
}

// Whether a value passes a row: one of the chosen, where any are, and none of
// those left out.
export function passes(f: Filter, x: string): boolean {
  return (f.only.length === 0 || f.only.includes(x)) && !f.not.includes(x)
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
// has one keeps it while it stays on the board and no other ticket holds
// it too; a new one, or the later by name of two sharing one, starts at its
// name's slot and moves on to the next choice no other ticket has. Only
// with more tickets than choices do two share one, and then by name alone.
export function pickOf(ids: readonly string[], kept: Readonly<Record<string, string>>, choices: readonly string[]): Record<string, string> {
  const picks: Record<string, string> = {}
  const held = new Set<string>()
  for (const id of [...ids].sort()) {
    const pick = kept[id]
    if (pick === undefined || !choices.includes(pick) || held.has(pick)) continue
    picks[id] = pick
    held.add(pick)
  }
  for (const id of [...ids].sort()) {
    if (picks[id] !== undefined) continue
    const used = new Set(Object.values(picks))
    const start = slotOf(id, choices.length)
    const free = choices.map((_, i) => choices[(start + i) % choices.length]!).find(c => !used.has(c))
    picks[id] = free ?? choices[start]!
  }
  return picks
}

// A working agent's spark turns through these, one each frame.
export const SPINNER = ['✻', '✳', '✶', '✢']
