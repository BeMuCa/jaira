// The board at a glance, for the Overview pane: what is in work, what could
// be started next, what was put back, and the story of the last days — what
// was created and what was put back on which day, and by whom. Read off the
// board alone (`jaira lanes`, `jaira list`, `jaira next --all`), never out of
// a session: so it is the same after a night, a restart or another session.
// register.tsx runs the commands; this file only reads what they printed.

import { isPersons, type Lane } from './board.ts'

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

function parse<T>(json: string): T | undefined {
  try {
    return JSON.parse(json) as T
  } catch {
    return undefined
  }
}

const oneLine = (s: string | undefined) => (s ?? '').replace(/\s+/g, ' ').trim()

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

// In work: a ticket in an agent's lane that jaira next does not offer, so
// somebody holds it. In the queue: what jaira next offers, short of what was
// put back and what waits for a person (the band shows those). Put back:
// every ticket in backlog or blocked. Each list the last touched first. The history:
// every ticket created since `days` ago, and every one put back since then on
// a later day than it was created — jaira keeps no time a ticket changed
// lanes, so its last change stands for it.
export function overviewOf(lanesJson: string, listJson: string, nextJson: string, now: number, days = HISTORY_DAYS): Overview {
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
    else if (lane.agentic) working.push(row(t, t.assignee))
  }
  const since = now - days * 86_400_000
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
