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
  // The question the assistant put to the person and nobody answered yet,
  // the ticket it is about, that ticket's title once looked up, and the
  // choices it named; all four go together.
  question?: string
  about?: string
  aboutTitle?: string
  options?: string[]
  lastTurn: number
}

const OPTIONS = 9

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
    : threads.map(t => `#${t.id} ${t.title}${t.ticket ? ` [ticket ${t.ticket}]` : ''}${t.question ? ` [open question${t.about ? ` on ${t.about}` : ''}: ${t.question}]` : ''}`).join('\n')
  return [
    'You sort one turn of a conversation between a person and an AI assistant into threads.',
    'A thread is one subject the person raised. Reply with ONE JSON object and nothing else:',
    '{"thread": <number of an existing thread, or null for a new one>, "title": "<thread title, at most 8 words, in English>",',
    ' "line": "<what happened in this turn, one short line, in English>",',
    withTickets
      ? ' "ticket": "<the 6-character jaira ticket id the assistant created or linked for this thread in this turn, or null>",'
      : ' "ticket": null,',
    ' "question": "<a question the assistant asked the person in this turn that the person must answer, one line, or null>",',
    withTickets
      ? ' "about": "<the 6-character jaira ticket id that question is about, or null>",'
      : ' "about": null,',
    ' "options": [<the answers that question offered the person to choose from, each a short string as the assistant worded it; [] when it offered none>],',
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
  about: string | null
  options: string[]
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
  const id = (v: unknown) => {
    const s = str(v)
    return s !== null && /^[0-9A-Za-z]{6}$/.test(s) ? s.toUpperCase() : null
  }
  const options = Array.isArray(raw.options)
    ? raw.options.map(str).filter((o): o is string => o !== null).map(o => clip(o, 160)).slice(0, OPTIONS)
    : []
  return {
    thread: typeof raw.thread === 'number' ? raw.thread : null,
    title: str(raw.title) ?? '…',
    line: str(raw.line) ?? '',
    ticket: id(raw.ticket),
    question: str(raw.question),
    about: id(raw.about),
    options,
    answered: Array.isArray(raw.answered) ? raw.answered.filter((n): n is number => typeof n === 'number') : [],
  }
}

// A thread whose question is answered: the question goes, and with it what
// belonged to it.
export function answered(t: Thread, turn: number): Thread {
  const { question: _q, about: _a, aboutTitle: _t, options: _o, ...rest } = t
  return { ...rest, lastTurn: turn }
}

// The threads after this turn's filing.
export function file(threads: readonly Thread[], f: Filing, turn: number): Thread[] {
  const next = threads.map(t => (f.answered.includes(t.id) ? answered(t, turn) : t))
  const at = f.thread === null ? -1 : next.findIndex(t => t.id === f.thread)
  const base: Thread = at >= 0
    ? next[at]!
    : { id: Math.max(0, ...threads.map(t => t.id)) + 1, title: clip(f.title, 80), lines: [], lastTurn: turn }
  // A new question brings its own ticket and choices; none keeps the open one
  // as it was, unless this very turn answered it.
  const kept = f.answered.includes(base.id) ? answered(base, turn) : base
  const asked: Thread = f.question === null
    ? kept
    : {
      ...answered(base, turn),
      question: f.question,
      ...(f.about !== null ? { about: f.about } : {}),
      ...(f.options.length > 0 ? { options: f.options } : {}),
    }
  const updated: Thread = {
    ...asked,
    lines: f.line === '' ? base.lines : [...base.lines, clip(f.line, 160)].slice(-LINES),
    lastTurn: turn,
    ...(f.ticket !== null ? { ticket: f.ticket } : {}),
  }
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

// One question waiting for the person, as the 'Waiting for you' pane lists it.
export type Ask = {
  thread: number
  // The ticket it is about: the one named with the question, else the
  // thread's own; undefined where neither is known.
  ticket?: string
  title: string
  question: string
  options: string[]
}

// Every open question, the one asked last first.
export function asksOf(threads: readonly Thread[]): Ask[] {
  return [...threads]
    .filter(t => t.question !== undefined)
    .sort((a, b) => b.lastTurn - a.lastTurn || b.id - a.id)
    .map(t => {
      const ticket = t.about ?? t.ticket
      return {
        thread: t.id,
        ...(ticket !== undefined ? { ticket } : {}),
        title: (t.about !== undefined ? t.aboutTitle : undefined) ?? t.title,
        question: t.question!,
        options: t.options ?? [],
      }
    })
}

// What a pick sends into the session: the ticket, or the thread's title where
// there is none, and the choice - 'VTSF24: yes'.
export function answerLine(a: Ask, choice: string): string {
  return `${a.ticket ?? a.title}: ${choice}`
}

// The tickets a question names whose title is not looked up yet.
export function untitled(threads: readonly Thread[]): string[] {
  return [...new Set(threads.filter(t => t.question !== undefined && t.about !== undefined && t.aboutTitle === undefined).map(t => t.about!))]
}

// The threads with a ticket's title written onto every question about it.
export function titled(threads: readonly Thread[], ticket: string, title: string): Thread[] {
  return threads.map(t => (t.question !== undefined && t.about === ticket && t.aboutTitle === undefined ? { ...t, aboutTitle: clip(title, 120) } : t))
}

// The threads with this thread's question taken as answered, at this turn.
export function answer(threads: readonly Thread[], thread: number, turn: number): Thread[] {
  return threads.map(t => (t.id === thread ? answered(t, turn) : t))
}
