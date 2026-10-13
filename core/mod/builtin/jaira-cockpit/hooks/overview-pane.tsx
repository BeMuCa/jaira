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

type OverviewTab = 'queue' | 'log'
type OverviewState = { tab: OverviewTab; folded: Record<string, boolean> }
type OverviewActions = { close: () => void; tab: (tab: OverviewTab) => void; fold: (key: string) => void }

// `view` null: the board could not be read.
export function drawOverview({ Box, Button, Text }: ElementTable, view: Overview | null, state: OverviewState, act: OverviewActions) {
  const closeButton = <Button key="overview-close" label={w.close} onPress={act.close} />
  if (view === null) return <Box flexDirection="column"><Text dimColor>{w.unread}</Text>{closeButton}</Box>
  const tabButton = (tab: OverviewTab, label: string) => (
    <Button key={`overview-tab-${tab}`} label={label} variant={state.tab === tab ? 'primary' : undefined} dimColor={state.tab !== tab} onPress={() => act.tab(tab)} />
  )
  // A section folds at its heading; open until the person folds it.
  const section = (key: string, heading: string, rows: Row[], extra: (r: Row) => string) => {
    const open = state.folded[key] !== true
    const label = `${open ? '▾' : '▸'} ${heading} (${rows.length})`
    return (
      <Box key={key} flexDirection="column" marginTop={1}>
        <Box flexDirection="row">
          <Button key={`overview-fold-${key}`} label={label} onPress={() => act.fold(key)} />
        </Box>
        {open && rows.length === 0 ? <Text dimColor>{w.none}</Text> : null}
        {open
          ? rows.map(r => (
            <Box key={`${key}-${r.id}`} flexDirection="column">
              <Text wrap="wrap"><Text bold>{r.id}</Text>{SEP}{r.lane}{SEP}{r.title}</Text>
              {extra(r) !== '' ? <Text dimColor wrap="wrap">  {extra(r)}</Text> : null}
            </Box>
          ))
          : null}
      </Box>
    )
  }
  const queue = [
    section('working', w.working, view.working, r => r.who),
    section('queue', w.queue, view.queue, r => r.who),
    section('parked', w.parked, view.parked, r => [r.who, r.reason].filter(Boolean).join(': ')),
  ]
  const log = [
    <Box key="history" marginTop={1}><Text bold>{w.history(HISTORY_DAYS)}</Text></Box>,
    view.history.length === 0 ? <Text key="history-none" dimColor>{w.none}</Text> : null,
    ...view.history.map(d => (
      <Box key={`day-${d.day}`} flexDirection="column">
        <Text dimColor>── {d.day} ──</Text>
        {d.events.map(ev => (
          <Text key={`${d.day}-${ev.kind}-${ev.id}`} wrap="wrap">
            {clockOf(ev.at)}  <Text bold>{ev.id}</Text>{SEP}{ev.lane}{SEP}{ev.title}{SEP}{ev.kind === 'created' ? w.createdBy(ev.who) : w.parkedBy(ev.who)}
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
