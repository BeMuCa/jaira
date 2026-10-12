// What each agent used, one row per session in ~/.jaira/metrics.csv, in the
// columns scripts/lane-eval/metrics.sh writes its metrics.csv in — so the
// two can be read side by side. Pure; register.tsx sums the turns and writes.

export const HEADER =
  'arm,ticket,run,lane,input_tokens,cache_write_tokens,cache_read_tokens,output_tokens,cost_eq,peak_context,duration_ms,tool_uses,loops,verdict'

export type Totals = {
  input: number
  cacheWrite: number
  cacheRead: number
  output: number
  // Context of the last turn: input + cache read + cache write + output.
  peak: number
  tools: number
}

export const emptyTotals = (): Totals => ({ input: 0, cacheWrite: 0, cacheRead: 0, output: 0, peak: 0, tools: 0 })

export type Usage = {
  input_tokens?: number
  cache_creation_input_tokens?: number
  cache_read_input_tokens?: number
  output_tokens?: number
}

export function add(t: Totals, u: Usage): Totals {
  const input = u.input_tokens ?? 0
  const cacheWrite = u.cache_creation_input_tokens ?? 0
  const cacheRead = u.cache_read_input_tokens ?? 0
  const output = u.output_tokens ?? 0
  return {
    input: t.input + input,
    cacheWrite: t.cacheWrite + cacheWrite,
    cacheRead: t.cacheRead + cacheRead,
    output: t.output + output,
    peak: input + cacheWrite + cacheRead + output,
    tools: t.tools,
  }
}

const cell = (v: string | number) => {
  const s = String(v)
  return /[",\n]/.test(s) ? `"${s.replace(/"/g, '""')}"` : s
}

// arm is the dispatch mode a live run happens under ("live" when unknown),
// run the session's id; loops and verdict are lane-eval's and stay empty.
export function row(input: { arm: string; ticket: string; session: string; lane: string; totals: Totals; durationMs: number }): string {
  const t = input.totals
  const costEq = Math.floor(t.input + t.cacheWrite * 1.25 + t.cacheRead * 0.1)
  return [
    input.arm, input.ticket, input.session, input.lane,
    t.input, t.cacheWrite, t.cacheRead, t.output, costEq, t.peak, input.durationMs, t.tools, '', '',
  ].map(cell).join(',')
}

// The file with the row added, its header written once.
export function appended(existing: string | undefined, line: string): string {
  const head = existing === undefined || existing.trim() === '' ? `${HEADER}\n` : existing.endsWith('\n') ? existing : `${existing}\n`
  return `${head}${line}\n`
}

// The rows of one ticket, for its pane: lane and what it used.
export function rowsOf(csv: string, ticket: string): string[] {
  return csv
    .split('\n')
    .slice(1)
    .filter(l => l.split(',')[1] === ticket)
    .map(l => {
      const c = l.split(',')
      return `${c[3] || '—'} · cost_eq ${c[8]} · peak ${Math.round(Number(c[9]) / 1000)}k · ${Math.round(Number(c[10]) / 60_000)} min · ${c[11]} tools`
    })
}
