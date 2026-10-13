// The Overview pane the band's Overview button opens: the board at a glance
// (overview.ts reads it), in two tabs. Queue lists in work, queue and put
// back one beneath the other, at full width, each folded at its heading;
// Log is the history of the last days. Only the drawing: register.tsx reads
// the board, keeps tab and folds, and hooks the pane, since every call on
// `$` lives there.

import type { ElementTable } from 'claude-code'
import { clockOf, HISTORY_DAYS, type Overview, type Row } from './overview.ts'
import { overviewWords as w } from './texts.ts'

const SEP = ' | '

// Each group its own colour, at its heading and at every row's handle, so the
// three read apart at a glance. Theme keys: they follow the person's theme.
const GROUP = { working: 'claude', queue: 'suggestion', parked: 'warning' } as const

// A person's colour, derived from the name and nothing else: the same person
// reads the same in both tabs and after every reopening. ANSI names, apart
// from the group colours, which the terminal's own scheme keeps readable.
const PEOPLE = ['cyan', 'magenta', 'green', 'cyanBright', 'magentaBright', 'greenBright']

function personColor(name: string): string {
  let h = 0x811c9dc5
  for (const c of name) h = Math.imul(h ^ c.codePointAt(0)!, 0x01000193)
  return PEOPLE[(h >>> 0) % PEOPLE.length]!
}

type OverviewTab = 'queue' | 'log'
// `people`: the names the Log is filtered to; empty shows everybody.
type OverviewState = { tab: OverviewTab; folded: Record<string, boolean>; people: string[] }
type OverviewActions = { close: () => void; tab: (tab: OverviewTab) => void; fold: (key: string) => void; person: (name: string) => void; everybody: () => void }

// `view` null: the board could not be read.
export function drawOverview({ Box, Button, Text }: ElementTable, view: Overview | null, state: OverviewState, act: OverviewActions) {
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
  return (
    <Box flexDirection="column">
      <Box flexDirection="row" columnGap={2}>
        {tabButton('queue', w.tabQueue)}
        {tabButton('log', w.tabLog)}
      </Box>
      {state.tab === 'queue' ? queue : log}
      {closeButton}
    </Box>
  )
}
