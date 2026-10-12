// The Overview pane the band's Overview button opens: the board at a glance
// (overview.ts reads it), drawn as three columns like lanes and the history
// of the last days beneath them. Only the drawing: register.tsx reads the
// board and hooks the pane, since every call on `$` lives there.

import { clockOf, HISTORY_DAYS, type Overview, type Row } from './overview.ts'
import { overviewWords as w } from './texts.ts'

const SEP = ' | '

// The elements the surface draws with, as $.ui.resolve hands them over.
type Table = { Box: any; Button: any; Text: any }

// `view` null: the board could not be read.
export function drawOverview({ Box, Button, Text }: Table, view: Overview | null, close: () => void) {
  const closeButton = <Button key="overview-close" label={w.close} onPress={close} />
  if (view === null) return <Box flexDirection="column"><Text dimColor>{w.unread}</Text>{closeButton}</Box>
  const column = (key: string, heading: string, rows: Row[], extra: (r: Row) => string) => (
    <Box key={key} flexDirection="column" width="33%">
      <Text bold>{heading} ({rows.length})</Text>
      {rows.length === 0 ? <Text dimColor>{w.none}</Text> : null}
      {rows.map(r => (
        <Box key={`${key}-${r.id}`} flexDirection="column">
          <Text wrap="truncate-end"><Text bold>{r.id}</Text>{SEP}{r.lane}{SEP}{r.title}</Text>
          {extra(r) !== '' ? <Text dimColor wrap="truncate-end">  {extra(r)}</Text> : null}
        </Box>
      ))}
    </Box>
  )
  return (
    <Box flexDirection="column">
      <Box flexDirection="row" columnGap={2}>
        {column('working', w.working, view.working, r => r.who)}
        {column('queue', w.queue, view.queue, r => r.who)}
        {column('parked', w.parked, view.parked, r => [r.who, r.reason].filter(Boolean).join(': '))}
      </Box>
      <Box marginTop={1}><Text bold>{w.history(HISTORY_DAYS)}</Text></Box>
      {view.history.length === 0 ? <Text dimColor>{w.none}</Text> : null}
      {view.history.map(d => (
        <Box key={`day-${d.day}`} flexDirection="column">
          <Text dimColor>── {d.day} ──</Text>
          {d.events.map(ev => (
            <Text key={`${d.day}-${ev.kind}-${ev.id}`} wrap="truncate-end">
              {clockOf(ev.at)}  <Text bold>{ev.id}</Text>{SEP}{ev.lane}{SEP}{ev.title}{SEP}{ev.kind === 'created' ? w.createdBy(ev.who) : w.parkedBy(ev.who)}
            </Text>
          ))}
        </Box>
      ))}
      {closeButton}
    </Box>
  )
}
