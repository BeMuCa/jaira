// The Overview pane the band's Overview button opens: the board at a glance
// (overview.ts reads it), in two tabs. Queue lists in work, queue and put
// back one beneath the other, at full width, each folded at its heading;
// Log is the history of the last days. Only the drawing: register.tsx reads
// the board, keeps tab and folds, and hooks the pane, since every call on
// `$` lives there.

import type { ElementTable } from 'claude-code'
import { clockOf, COLORS, durationOf, HATS, HISTORY_DAYS, passes, slotOf, SPINNER, stoneOf, stoneSize, workerKey, type Event, type Filter, type Overview, type Row, type Worker } from './overview.ts'
import { overviewWords as w } from './texts.ts'

const SEP = ' | '

// Each group its own colour, at its heading and at every row's handle, so the
// three read apart at a glance. Theme keys: they follow the person's theme.
const GROUP = { working: 'claude', queue: 'suggestion', parked: 'warning' } as const

// A person's colour, derived from the name and nothing else: the same person
// reads the same in both tabs and after every reopening.
function personColor(name: string): string {
  return COLORS[slotOf(name, COLORS.length)]!
}

export type OverviewTab = 'queue' | 'log' | 'workers'
// The Workers tab: its agents (null: the board could not be read), each
// ticket's colour and hat, when each agent was first seen (by workerKey),
// the first line of each ticket's last report, the time now, how many
// frames ago the stone last lost a chunk (-1: none falling), the frame the
// animation is at, and the pane's room, in cells.
type WorkersScene = {
  workers: Worker[] | null
  colors: Record<string, string>
  hats: Record<string, string>
  since: Record<string, number>
  reports: Record<string, string>
  now: number
  falling: number
  frame: number
  columns: number
  rows: number
}

// ── the Workers tab's drawing: rows of cells, so a figure stands by columns ─

type Style = { color?: string; dim?: boolean; bold?: boolean }
type Cell = Style & { ch: string }

// Claude Code's own crab, Clawd, nine cells wide and four high with its hat,
// a spark's cell on the side it faces the stone. With more agents than fit,
// a small one: four cells wide, two high.
type Sprite = { w: number; h: number }
const BIG: Sprite = { w: 10, h: 4 }
const SMALL: Sprite = { w: 5, h: 2 }
const HEAD = ' ▐▛███▜▌ '
const ARMS = ['▝▜█████▛▘', '▗▜█████▛▖']
const LEGS = { stand: '  ▘▘ ▝▝  ', sit: '  ▀▀▀▀▀  ' }
const MINI_HEAD = '▐▛▜▌'
const MINI_LEGS = ['▝▘▝▘', '▗▘▝▖']
// A chunk falls this many frames once the stone lost one.
const FALL_FRAMES = 6

function put(row: Cell[], x: number, text: string, style: Style = {}): void {
  ;[...text].forEach((ch, i) => {
    if (x + i >= 0 && x + i < row.length) row[x + i] = { ch, ...style }
  })
}

// Where each agent stands: on the ground on both sides of the stone, the
// nearest places first, left and right in turn, then a row higher. The
// sprite is the big crab while every agent fits, else the small one; with
// more agents than even that leaves room for, the rest are only in the lines.
type Place = { x: number; y: number; left: boolean }
function placesOf(n: number, width: number, height: number, stone: { x: number; w: number }, sprite: Sprite): Place[] {
  const leftRoom = Math.max(0, stone.x - 1)
  const rightRoom = Math.max(0, width - (stone.x + stone.w) - 1)
  const perLeft = Math.floor(leftRoom / sprite.w)
  const perRight = Math.floor(rightRoom / sprite.w)
  const levels = Math.floor(height / sprite.h)
  const places: Place[] = []
  for (let level = 0; level < levels && places.length < n; level++) {
    for (let j = 0; j < Math.max(perLeft, perRight) && places.length < n; j++) {
      const y = height - (level + 1) * sprite.h
      if (j < perLeft) places.push({ x: stone.x - 1 - (j + 1) * sprite.w, y, left: true })
      if (j < perRight && places.length < n) places.push({ x: stone.x + stone.w + 1 + j * sprite.w, y, left: false })
    }
  }
  return places
}

// The picture: every agent at one stone, on a ground line, as high as it
// needs to be and no higher than `room`. The stone's size follows the open
// work of every ticket with an agent at it (stoneSize); the line above it
// says how much. Empty when there is no room for even the small crab.
function sceneRows(scene: WorkersScene, workers: Worker[], room: number): Cell[][] {
  const width = Math.max(scene.columns, 20)
  const ground = room - 1
  if (ground < SMALL.h + 1) return []
  const open = stoneOf(workers)
  // The stone keeps a third of the width and one line above it for its count.
  const size = stoneSize(open, Math.floor(width / 3), ground - 1)
  const stone = { x: Math.floor((width - size.w) / 2), w: size.w }
  const big = placesOf(workers.length, width, ground, stone, BIG)
  const sprite = big.length === workers.length ? BIG : SMALL
  const places = sprite === BIG ? big : placesOf(workers.length, width, ground, stone, SMALL)
  // Only as high as the tallest of stone, its count and the figures.
  const top = Math.min(ground - size.h - 1, ...places.map(p => p.y))
  const blank = () => Array.from({ length: width }, () => ({ ch: ' ' }) as Cell)
  const rows = Array.from({ length: room }, blank)
  put(rows[ground]!, 0, '─'.repeat(width), { dim: true })
  const stoneTop = ground - size.h
  for (let y = stoneTop; y < ground; y++) {
    const line = size.h > 1 && y === stoneTop && size.w > 2 ? `▟${'█'.repeat(size.w - 2)}▙` : size.h === 1 ? '▄'.repeat(size.w) : '█'.repeat(size.w)
    put(rows[y]!, stone.x, line, { dim: true })
  }
  const count = w.stoneOpen(open)
  put(rows[stoneTop - 1]!, stone.x + Math.floor((size.w - count.length) / 2), count, { bold: true })
  if (scene.falling >= 0 && scene.falling < FALL_FRAMES) {
    const y = Math.min(ground - 1, stoneTop + scene.falling)
    put(rows[y]!, stone.x + size.w, '▞', { dim: true })
  }
  workers.forEach((f, k) => {
    const at = places[k]
    if (at === undefined) return
    const look: Style = { color: scene.colors[f.id], bold: f.state !== 'idle', dim: f.state === 'idle' }
    const working = f.state === 'working'
    const swing = Math.floor(scene.frame / 2) % 2
    // The sprite on the side away from the stone, the spark's cell towards it.
    const body = at.left ? at.x : at.x + 1
    const spark = at.left ? at.x + sprite.w - 1 : at.x
    const mark = working
      ? { ch: SPINNER[scene.frame % SPINNER.length]!, style: { color: scene.colors[f.id] } }
      : f.state === 'question' ? { ch: '?', style: { color: 'warning', bold: true } } : { ch: 'z', style: { dim: true } }
    if (sprite === BIG) {
      put(rows[at.y]!, body + 1, HATS[scene.hats[f.id] ?? ''] ?? '', look)
      put(rows[at.y + 1]!, body, HEAD, look)
      put(rows[at.y + 2]!, body, ARMS[working ? swing : 0]!, look)
      put(rows[at.y + 3]!, body, working ? LEGS.stand : LEGS.sit, look)
      put(rows[at.y + 1]!, spark, mark.ch, mark.style)
    } else {
      put(rows[at.y]!, body, MINI_HEAD, look)
      put(rows[at.y + 1]!, body, MINI_LEGS[working ? swing : 0]!, look)
      put(rows[at.y]!, spark, mark.ch, mark.style)
    }
  })
  const hidden = workers.length - places.length
  if (hidden > 0) put(rows[top]!, 0, w.moreWorkers(hidden), { dim: true })
  return rows.slice(Math.max(0, top - 1))
}

// One line per agent, beneath the picture: its ticket in the ticket's
// colour, lane, kind, state, how long it has been seen, and the first line
// of its ticket's last report — the title, dimmed, where none came yet.
function workerLines(Text: ElementTable['Text'], scene: WorkersScene, workers: Worker[]) {
  const pad = (pick: (f: Worker) => string) => Math.max(...workers.map(f => pick(f).length))
  const lane = pad(f => f.lane)
  const role = pad(f => f.role)
  const state = pad(f => w.states[f.state])
  return workers.map((f, k) => {
    const seen = scene.since[workerKey(f)]
    const dur = seen === undefined ? '' : durationOf(scene.now - seen)
    const report = scene.reports[f.id] ?? ''
    return (
      <Text key={`worker-${k}`} wrap="truncate">
        <Text bold color={scene.colors[f.id]}>{f.id}</Text>
        {`  ${f.lane.padEnd(lane)}  ${f.role.padEnd(role)}  `}
        <Text color={f.state === 'question' ? 'warning' : undefined} dimColor={f.state === 'idle'}>{w.states[f.state].padEnd(state)}</Text>
        {`  ${dur.padStart(5)}  `}
        {report !== '' ? report : <Text dimColor>{f.title}</Text>}
      </Text>
    )
  })
}

// The sign each kind of Log line starts with. Only single code points the
// terminal draws two cells wide by default (Emoji_Presentation): a sign that
// needs a variation selector (↩️) is one cell in one terminal and two in
// another, and shifts every column after it.
const KIND_SIGN: Record<Event['kind'], string> = { created: '🆕', parked: '🔙' }

// `people`: the names the Log is filtered to, `lanes` the lanes, each chosen
// and left out; both empty shows all.
type OverviewState = { tab: OverviewTab; folded: Record<string, boolean>; people: Filter; lanes: Filter }
type OverviewActions = {
  close: () => void
  tab: (tab: OverviewTab) => void
  fold: (key: string) => void
  person: (name: string) => void
  everybody: () => void
  lane: (id: string) => void
  allLanes: () => void
}

// The head of a pane: one row of buttons as the terminal draws every Button,
// `[ label ]`, the tabs (if any) on the left and Close at the right edge, so
// each is one line high and hit anywhere on its brackets. The host's own ×
// stays beside it; a mod cannot size it, so Close is the one to aim at. The
// row never wraps: the Overview budgets it one row (OVERVIEW_ROWS), and its
// widest head, the three tabs and Close, fits an 80-column terminal.
export function headRow({ Box, Button }: Pick<ElementTable, 'Box' | 'Button'>, key: string, close: { key: string; label: string; onPress: () => void }, left: unknown[] = []) {
  return (
    <Box key={key} flexDirection="row" flexWrap="nowrap" justifyContent="space-between">
      <Box key={`${key}-tabs`} flexDirection="row" flexWrap="nowrap" columnGap={1} flexShrink={0}>{left}</Box>
      <Box key={`${key}-close`} flexShrink={0}><Button key={close.key} role="dismiss" label={close.label} onPress={close.onPress} /></Box>
    </Box>
  )
}

// `view` null: the board could not be read.
export function drawOverview({ Box, Button, Text }: ElementTable, view: Overview | null, state: OverviewState, act: OverviewActions, scene: WorkersScene) {
  const person = (who: string) => <Text key="who" color={personColor(who)}>{who}</Text>
  const close = { key: 'overview-close', label: w.close, onPress: act.close }
  if (view === null) return <Box flexDirection="column">{headRow({ Box, Button }, 'overview-top-row', close)}<Text dimColor>{w.unread}</Text></Box>
  // The open tab is drawn as the primary button, in the theme's accent; the
  // others are dimmed.
  const tabButton = (tab: OverviewTab, label: string) => (
    <Button key={`overview-tab-${tab}`} label={label} variant={state.tab === tab ? 'primary' : undefined} dimColor={state.tab !== tab} onPress={() => act.tab(tab)} />
  )
  // A section folds at its heading; open until the person folds it.
  const section = (key: keyof typeof GROUP, heading: string, rows: Row[]) => {
    const open = state.folded[key] !== true
    const label = `${open ? '▾' : '▸'} ${heading} (${rows.length})`
    return (
      <Box key={key} flexDirection="column" marginTop={1}>
        <Box flexDirection="row">
          <Button key={`overview-fold-${key}`} label={label} onPress={() => act.fold(key)}><Text bold color={GROUP[key]}>{label}</Text></Button>
        </Box>
        {open && rows.length === 0 ? <Text dimColor>{w.none}</Text> : null}
        {open
          ? rows.map(r => (
            <Box key={`${key}-${r.id}`} flexDirection="column">
              <Text wrap="wrap"><Text bold color={GROUP[key]}>{r.id}</Text><Text dimColor>{SEP}{r.lane}{SEP}</Text>{r.title}</Text>
              {r.who !== '' || r.reason !== ''
                ? <Text wrap="wrap">  {r.who !== '' ? person(r.who) : null}{r.reason !== '' ? <Text dimColor>{r.who !== '' ? ': ' : ''}{r.reason}</Text> : null}</Text>
                : null}
            </Box>
          ))
          : null}
      </Box>
    )
  }
  const queue = [
    section('working', w.working, view.working),
    section('queue', w.queue, view.queue),
    section('parked', w.parked, view.parked),
  ]
  // The Log filters to the people chosen and to the lanes chosen — the lane a
  // ticket stands in now — several of each at once, as `u` does on the board,
  // and leaves out the ones left out; both rows hold together, and a day left
  // with nothing drops out.
  const events = view.history.flatMap(d => d.events)
  // A button for every value in today's events and for every value chosen or
  // left out: one with no event left (the ticket moved on, the events aged
  // out) keeps its button, so the filter it sets can still be undone.
  const set = (f: Filter) => [...f.only, ...f.not]
  const names = [...new Set([...events.map(ev => ev.who), ...set(state.people)].filter(Boolean))].sort()
  const present = new Set([...events.map(ev => ev.lane), ...set(state.lanes)].filter(Boolean))
  // Board order; a lane the board no longer names (or a view kept from before
  // lanes were read) goes last, so no event is left without its button.
  const order = view.lanes ?? []
  const laneIds = [...order.filter(l => present.has(l)), ...[...present].filter(l => !order.includes(l)).sort()]
  const chosen = (ev: Event) => passes(state.people, ev.who) && passes(state.lanes, ev.lane)
  const days = view.history.map(d => ({ ...d, events: d.events.filter(chosen) })).filter(d => d.events.length > 0)
  // A chip cycles on each press: off (dimmed), chosen (primary), left out
  // (`✕ name` in the theme's error colour), off again. All is lit only while
  // its row neither chooses nor leaves out anything, and clears both.
  const filterRow = (key: string, prefix: string, allKey: string, f: Filter, all: () => void, items: string[], pick: (id: string) => void, look: (id: string) => string | undefined) => {
    const idle = f.only.length === 0 && f.not.length === 0
    return (
      <Box key={key} flexDirection="row" flexWrap="wrap" columnGap={2} marginTop={1}>
        <Button key={allKey} label={w.everybody} variant={idle ? 'primary' : undefined} dimColor={!idle} onPress={all} />
        {items.map(id => {
          const on = f.only.includes(id)
          const out = f.not.includes(id)
          const label = out ? `✕ ${id}` : id
          return (
            <Button key={`${prefix}-${id}`} label={label} variant={on ? 'primary' : undefined} dimColor={!on && !out} onPress={() => pick(id)}>
              <Text color={out ? 'error' : look(id)}>{label}</Text>
            </Button>
          )
        })}
      </Box>
    )
  }
  const log = [
    filterRow('people', 'overview-person', 'overview-people-all', state.people, act.everybody, names, act.person, personColor),
    filterRow('lanes', 'overview-lane', 'overview-lanes-all', state.lanes, act.allLanes, laneIds, act.lane, () => undefined),
    <Box key="history" marginTop={1}><Text bold>{w.history(HISTORY_DAYS)}</Text></Box>,
    days.length === 0 ? <Text key="history-none" dimColor>{w.none}</Text> : null,
    ...days.map(d => (
      <Box key={`day-${d.day}`} flexDirection="column">
        <Text dimColor>── {d.day} ──</Text>
        {d.events.map(ev => (
          <Text key={`${d.day}-${ev.kind}-${ev.id}`} wrap="wrap">
            {KIND_SIGN[ev.kind]} {clockOf(ev.at)}  <Text bold>{ev.id}</Text><Text dimColor>{SEP}{ev.lane}{SEP}</Text>{ev.title}<Text dimColor>{SEP}</Text>
            {ev.kind === 'created' ? w.created : <Text color={GROUP.parked}>{w.putBack}</Text>}
            {ev.who !== '' ? [w.by, person(ev.who)] : null}
          </Text>
        ))}
      </Box>
    )),
  ]
  // The Workers tab: the picture, each row of cells one line of its runs of
  // one style each, and one line per agent beneath it — both within the
  // pane's rows, so nothing has to be scrolled to.
  const workers = () => {
    const all = scene.workers
    if (all === null) return [<Text key="workers-unread" dimColor>{w.unread}</Text>]
    if (all.length === 0) return [<Text key="workers-none" dimColor>{w.noWorkers}</Text>]
    const picture = sceneRows(scene, all, scene.rows - all.length).map((row, y) => {
      const runs: { text: string; style: Style }[] = []
      for (const { ch, ...style } of row) {
        const last = runs.at(-1)
        if (last && last.style.color === style.color && last.style.dim === style.dim && last.style.bold === style.bold) last.text += ch
        else runs.push({ text: ch, style })
      }
      return (
        <Text key={`scene-${y}`} wrap="truncate">
          {runs.map((r, i) => (r.style.color || r.style.dim || r.style.bold
            ? <Text key={`run-${i}`} color={r.style.color} dimColor={r.style.dim} bold={r.style.bold}>{r.text}</Text>
            : r.text))}
        </Text>
      )
    })
    return [...picture, ...workerLines(Text, scene, all)]
  }
  return (
    <Box flexDirection="column">
      {headRow({ Box, Button }, 'overview-top-row', close, [tabButton('queue', w.tabQueue), tabButton('log', w.tabLog), tabButton('workers', w.tabWorkers)])}
      {state.tab === 'queue' ? queue : state.tab === 'log' ? log : workers()}
    </Box>
  )
}
