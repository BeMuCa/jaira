// The Overview pane the band's Overview button opens: the board at a glance
// (overview.ts reads it), in two tabs. Queue lists in work, queue and put
// back one beneath the other, at full width, each folded at its heading;
// Log is the history of the last days. Only the drawing: register.tsx reads
// the board, keeps tab and folds, and hooks the pane, since every call on
// `$` lives there.

import type { ElementTable } from 'claude-code'
import { clockOf, COLORS, HATS, HISTORY_DAYS, slotOf, SPINNER, type Overview, type Row, type Spot, type Walker } from './overview.ts'
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
// The Workers tab: its figures (null: the board could not be read), where
// each ticket's figures stand on their track, each ticket's colour and hat,
// the frame the animation is at, and the pane's room, in cells.
type WorkersScene = {
  walkers: Walker[] | null
  where: Record<string, Spot>
  colors: Record<string, string>
  hats: Record<string, string>
  frame: number
  columns: number
  rows: number
}

// ── the Workers tab's drawing: rows of cells, so a figure moves by columns ─

type Style = { color?: string; dim?: boolean; bold?: boolean }
type Cell = Style & { ch: string }

// Claude Code's own crab, Clawd, nine cells wide, a spark's cell beside it.
const SPRITE = 10
const HEAD = ' ▐▛███▜▌ '
const ARMS = ['▝▜█████▛▘', '▗▜█████▛▖']
const LEGS = { stand: '  ▘▘ ▝▝  ', step: '  ▝▘ ▘▝  ', sit: '  ▀▀▀▀▀  ' }
// A track's lines: the walls' names, what floats above the hat, the hat,
// head, arms, legs, and the line that says who it is.
const TRACK = 7

function put(row: Cell[], x: number, text: string, style: Style = {}): void {
  ;[...text].forEach((ch, i) => {
    if (x + i >= 0 && x + i < row.length) row[x + i] = { ch, ...style }
  })
}

// The tracks, one per figure, top to bottom with a rule between, sharing
// the room's height. On each, the left wall is the lane its ticket stands
// in, the right one the lane it goes to; the crab stands as far between as
// it has come, and the line beneath says who, doing what, on which ticket.
function trackRows(scene: WorkersScene): Cell[][] {
  const { walkers, where, colors, hats, frame, columns } = scene
  const width = Math.max(columns, SPRITE + 6)
  const n = walkers?.length ?? 0
  const height = Math.max(TRACK, Math.floor((scene.rows - (n - 1)) / Math.max(n, 1)))
  const blank = () => Array.from({ length: width }, () => ({ ch: ' ' }) as Cell)
  const rows: Cell[][] = []
  ;(walkers ?? []).forEach((f, k) => {
    if (k > 0) {
      const rule = blank()
      put(rule, 0, '─'.repeat(width), { dim: true })
      rows.push(rule)
    }
    const track = Array.from({ length: height }, blank)
    const spot = where[f.id] ?? { lane: f.lane, x: f.progress }
    const leaving = spot.lane !== f.lane
    const walking = leaving || spot.x !== f.progress
    const right = leaving ? f.lane : f.next || '—'
    for (const line of track) {
      put(line, 0, '┃', { dim: true })
      put(line, width - 1, '┃', { dim: true })
    }
    // Each wall's name has half the line, so neither ever writes over the other.
    const room = Math.floor((width - 6) / 2)
    put(track[0]!, 2, spot.lane.slice(0, room), { bold: true })
    put(track[0]!, width - 3 - Math.min(right.length, room), right.slice(0, room), { dim: true })
    const x = 2 + Math.round(spot.x * (width - 4 - SPRITE))
    const top = height - 6
    const look: Style = { color: colors[f.id], bold: f.state !== 'idle', dim: f.state === 'idle' }
    if (f.state === 'question') put(track[top]!, x + 3, '(?)', { color: 'warning', bold: true })
    if (f.state === 'idle') put(track[top]!, x + 5, 'z z', { dim: true })
    put(track[top + 1]!, x + 1, HATS[hats[f.id] ?? ''] ?? '', look)
    put(track[top + 2]!, x, HEAD, look)
    put(track[top + 3]!, x, ARMS[f.state === 'working' && !walking ? Math.floor(frame / 2) % 2 : 0]!, look)
    if (f.state === 'working') put(track[top + 2]!, x + 9, SPINNER[frame % SPINNER.length]!, { color: colors[f.id] })
    put(track[top + 4]!, x, f.state === 'idle' && !walking ? LEGS.sit : walking && frame % 2 === 1 ? LEGS.step : LEGS.stand, look)
    const words = `${f.id} · ${w.states[f.state]} · `
    const caption = words + f.title.slice(0, Math.max(0, width - 4 - words.length))
    put(track[top + 5]!, Math.max(2, Math.min(x, width - 2 - caption.length)), caption)
    rows.push(...track)
  })
  return rows
}

// `people`: the names the Log is filtered to; empty shows everybody.
type OverviewState = { tab: OverviewTab; folded: Record<string, boolean>; people: string[] }
type OverviewActions = { close: () => void; tab: (tab: OverviewTab) => void; fold: (key: string) => void; person: (name: string) => void; everybody: () => void }

// `view` null: the board could not be read.
export function drawOverview({ Box, Button, Text }: ElementTable, view: Overview | null, state: OverviewState, act: OverviewActions, scene: WorkersScene) {
  const person = (who: string) => <Text key="who" color={personColor(who)}>{who}</Text>
  const closeButton = <Button key="overview-close" label={w.close} onPress={act.close} />
  if (view === null) return <Box flexDirection="column"><Text dimColor>{w.unread}</Text>{closeButton}</Box>
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
  // The Log filters to the people chosen, several at once, as `u` does on the
  // board; a day with nobody chosen in it drops out.
  const names = [...new Set(view.history.flatMap(d => d.events.map(ev => ev.who)).filter(Boolean))].sort()
  const chosen = (who: string) => state.people.length === 0 || state.people.includes(who)
  const days = view.history.map(d => ({ ...d, events: d.events.filter(ev => chosen(ev.who)) })).filter(d => d.events.length > 0)
  const log = [
    <Box key="people" flexDirection="row" flexWrap="wrap" columnGap={2} marginTop={1}>
      <Button key="overview-people-all" label={w.everybody} variant={state.people.length === 0 ? 'primary' : undefined} dimColor={state.people.length !== 0} onPress={act.everybody} />
      {names.map(name => (
        <Button key={`overview-person-${name}`} label={name} variant={state.people.includes(name) ? 'primary' : undefined} dimColor={!state.people.includes(name)} onPress={() => act.person(name)}>
          <Text color={personColor(name)}>{name}</Text>
        </Button>
      ))}
    </Box>,
    <Box key="history" marginTop={1}><Text bold>{w.history(HISTORY_DAYS)}</Text></Box>,
    days.length === 0 ? <Text key="history-none" dimColor>{w.none}</Text> : null,
    ...days.map(d => (
      <Box key={`day-${d.day}`} flexDirection="column">
        <Text dimColor>── {d.day} ──</Text>
        {d.events.map(ev => (
          <Text key={`${d.day}-${ev.kind}-${ev.id}`} wrap="wrap">
            {clockOf(ev.at)}  <Text bold>{ev.id}</Text><Text dimColor>{SEP}{ev.lane}{SEP}</Text>{ev.title}<Text dimColor>{SEP}</Text>
            {ev.kind === 'created' ? w.created : <Text color={GROUP.parked}>{w.putBack}</Text>}
            {ev.who !== '' ? [w.by, person(ev.who)] : null}
          </Text>
        ))}
      </Box>
    )),
  ]
  // The Workers tab: each row of cells one line, its runs of one style each.
  const workers = () => {
    if (scene.walkers === null) return [<Text key="workers-unread" dimColor>{w.unread}</Text>]
    if (scene.walkers.length === 0) return [<Text key="workers-none" dimColor>{w.noWorkers}</Text>]
    return trackRows(scene).map((row, y) => {
      const runs: { text: string; style: Style }[] = []
      for (const { ch, ...style } of row) {
        const last = runs.at(-1)
        if (last && last.style.color === style.color && last.style.dim === style.dim && last.style.bold === style.bold) last.text += ch
        else runs.push({ text: ch, style })
      }
      return (
        <Text key={`track-${y}`} wrap="truncate">
          {runs.map((r, i) => (r.style.color || r.style.dim || r.style.bold
            ? <Text key={`run-${i}`} color={r.style.color} dimColor={r.style.dim} bold={r.style.bold}>{r.text}</Text>
            : r.text))}
        </Text>
      )
    })
  }
  return (
    <Box flexDirection="column">
      <Box flexDirection="row" columnGap={2}>
        {tabButton('queue', w.tabQueue)}
        {tabButton('log', w.tabLog)}
        {tabButton('workers', w.tabWorkers)}
      </Box>
      {state.tab === 'queue' ? queue : state.tab === 'log' ? log : workers()}
      {closeButton}
    </Box>
  )
}
