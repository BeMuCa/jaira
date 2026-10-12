// Conversation threads: a person jumps between subjects, and five turns later
// nobody knows whether an idea became a ticket, was forgotten, or got a
// question back that was never answered. After each of the person's turns a
// cheap model files the turn under a thread and says where the thread stands.
// Pure; register.tsx makes the model call and keeps the threads.

export type Thread = {
  id: number
  title: string
  // One line per turn the thread got, oldest first, the last few kept.
  lines: string[]
  ticket?: string
  // The question the assistant put to the person and nobody answered yet.
  question?: string
  lastTurn: number
}

export const STALE_TURNS = 5
const LINES = 6
const THREADS = 30

export type Mark = '->' | '?' | '~' | ' '

// -> a ticket holds it; ? a question waits for the person; ~ left lying:
// STALE_TURNS turns with neither; blank: still being talked about.
export function markOf(t: Thread, turn: number): Mark {
  if (t.question !== undefined) return '?'
  if (t.ticket !== undefined) return '->'
  return turn - t.lastTurn >= STALE_TURNS ? '~' : ' '
}

export function tally(threads: readonly Thread[], turn: number): { questions: number; stale: number } {
  let questions = 0
  let stale = 0
  for (const t of threads) {
    const m = markOf(t, turn)
    if (m === '?') questions++
    else if (m === '~') stale++
  }
  return { questions, stale }
}

const clip = (s: string, n: number) => (s.length <= n ? s : `${s.slice(0, n - 1)}…`)

// The request to the cheap model: the threads so far, the person's message,
// the assistant's answer to it. `withTickets` false on a project without a
// board: then no thread can stand on a ticket.
export function classifyPrompt(threads: readonly Thread[], human: string, answer: string, withTickets: boolean): string {
  const known = threads.length === 0
    ? '(none yet)'
    : threads.map(t => `#${t.id} ${t.title}${t.ticket ? ` [ticket ${t.ticket}]` : ''}${t.question ? ` [open question: ${t.question}]` : ''}`).join('\n')
  return [
    'You sort one turn of a conversation between a person and an AI assistant into threads.',
    'A thread is one subject the person raised. Reply with ONE JSON object and nothing else:',
    '{"thread": <number of an existing thread, or null for a new one>, "title": "<thread title, at most 8 words, in English>",',
    ' "line": "<what happened in this turn, one short line, in English>",',
    withTickets
      ? ' "ticket": "<the 6-character jaira ticket id the assistant created or linked for this thread in this turn, or null>",'
      : ' "ticket": null,',
    ' "question": "<a question the assistant asked the person in this turn that the person must answer, one line, or null>",',
    ' "answered": [<numbers of threads whose open question the person\'s message answers>]}',
    '',
    'Threads so far:',
    known,
    '',
    'The person wrote:',
    clip(human, 4000),
    '',
    'The assistant answered:',
    clip(answer, 6000),
  ].join('\n')
}

export type Filing = {
  thread: number | null
  title: string
  line: string
  ticket: string | null
  question: string | null
  answered: number[]
}

export function parseFiling(text: string): Filing | undefined {
  const json = /\{[\s\S]*\}/.exec(text)?.[0]
  if (json === undefined) return undefined
  let raw: Record<string, unknown>
  try {
    raw = JSON.parse(json) as Record<string, unknown>
  } catch {
    return undefined
  }
  const str = (v: unknown) => (typeof v === 'string' && v.trim() !== '' ? v.trim() : null)
  const ticket = str(raw.ticket)
  return {
    thread: typeof raw.thread === 'number' ? raw.thread : null,
    title: str(raw.title) ?? '…',
    line: str(raw.line) ?? '',
    ticket: ticket !== null && /^[0-9A-Za-z]{6}$/.test(ticket) ? ticket.toUpperCase() : null,
    question: str(raw.question),
    answered: Array.isArray(raw.answered) ? raw.answered.filter((n): n is number => typeof n === 'number') : [],
  }
}

// The threads after this turn's filing.
export function file(threads: readonly Thread[], f: Filing, turn: number): Thread[] {
  const next = threads.map(t => (f.answered.includes(t.id) ? { ...t, question: undefined, lastTurn: turn } : t))
  const at = f.thread === null ? -1 : next.findIndex(t => t.id === f.thread)
  const base: Thread = at >= 0
    ? next[at]!
    : { id: Math.max(0, ...threads.map(t => t.id)) + 1, title: clip(f.title, 80), lines: [], lastTurn: turn }
  const updated: Thread = {
    ...base,
    lines: f.line === '' ? base.lines : [...base.lines, clip(f.line, 160)].slice(-LINES),
    lastTurn: turn,
    ...(f.ticket !== null ? { ticket: f.ticket } : {}),
    question: f.question ?? (f.answered.includes(base.id) ? undefined : base.question),
  }
  if (updated.question === undefined) delete updated.question
  if (at >= 0) next[at] = updated
  else next.push(updated)
  return next.slice(-THREADS)
}

export function listThreads(threads: readonly Thread[], turn: number, title: string, none: string): string {
  if (threads.length === 0) return none
  const rows = [...threads].reverse().map(t => {
    const mark = markOf(t, turn)
    const stand = mark === '->' ? ` -> ${t.ticket}` : mark === '?' ? ` ? ${t.question}` : ''
    const last = t.lines.at(-1)
    return `${mark} #${t.id} ${t.title}${stand}${last ? `\n     ${last}` : ''}`
  })
  return `${title}\n${rows.join('\n')}`
}
