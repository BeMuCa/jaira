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
  // When the question was asked (ms), and the lane the ticket it is about
  // stood in then, read with its title; a ticket that leaves that lane
  // settles the question (CSQXKR).
  askedAt?: number
  aboutLane?: string
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
    ' "question": "<a question the assistant asked the person in this turn that the person must answer, one line, quoted in the assistant\'s own words and language, never translated, or null>",',
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
  const { question: _q, about: _a, aboutTitle: _t, options: _o, askedAt: _at, aboutLane: _l, ...rest } = t
  return { ...rest, lastTurn: turn }
}

// The threads after this turn's filing; now: when it was asked, for a new
// question.
export function file(threads: readonly Thread[], f: Filing, turn: number, now?: number): Thread[] {
  const next = threads.map(t => (f.answered.includes(t.id) ? answered(t, turn) : t))
  const at = f.thread === null ? -1 : next.findIndex(t => t.id === f.thread)
  const base: Thread = at >= 0
    ? next[at]!
    : { id: Math.max(0, ...threads.map(t => t.id)) + 1, title: clip(f.title, 80), lines: [], lastTurn: turn }
  // A new question brings its own ticket and choices; none keeps the open one
  // as it was, unless this very turn answered it.
  const kept = f.answered.includes(base.id) ? answered(base, turn) : base
  // The same question asked again on its thread keeps the moment it was
  // first asked.
  const askedAt = base.question === f.question && base.askedAt !== undefined ? base.askedAt : now
  const asked: Thread = f.question === null
    ? kept
    : {
      ...answered(base, turn),
      question: f.question,
      ...(askedAt !== undefined ? { askedAt } : {}),
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
  return capped(next, updated.id)
}

// At most THREADS threads: the ones touched longest ago go first. Neither
// the thread this turn touched nor one whose question still waits is ever
// dropped (CSQXKR), however many of them there are; nor is the newest, whose
// id the next new thread counts on from.
function capped(threads: Thread[], touched: number): Thread[] {
  const over = threads.length - THREADS
  if (over <= 0) return threads
  const newest = Math.max(...threads.map(t => t.id))
  const drop = new Set(threads.filter(t => t.question === undefined && t.id !== touched && t.id !== newest).sort((a, b) => a.lastTurn - b.lastTurn || a.id - b.id).slice(0, over).map(t => t.id))
  return threads.filter(t => !drop.has(t.id))
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
  askedAt?: number
  // Set once it is settled: answered, dismissed, or its ticket moved on.
  answeredAt?: number
}

// How long a settled question stays in the pane, dimmed.
export const ANSWERED_MS = 24 * 60 * 60 * 1000

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
        // '' : looked up, and jaira had no title for it.
        title: (t.about !== undefined ? t.aboutTitle : undefined) || t.title,
        question: t.question!,
        options: t.options ?? [],
        ...(t.askedAt !== undefined ? { askedAt: t.askedAt } : {}),
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

// The threads with a ticket's title written onto every question about it;
// '' where jaira had none, so it is not looked up again. lane: the lane the
// ticket stands in now, which the question waits on.
export function titled(threads: readonly Thread[], ticket: string, title: string, lane?: string): Thread[] {
  return threads.map(t => (t.question !== undefined && t.about === ticket && t.aboutTitle === undefined
    ? { ...t, aboutTitle: clip(title, 120), ...(lane !== undefined && lane !== '' ? { aboutLane: lane } : {}) }
    : t))
}

// The tickets whose lane open questions wait on.
export function watched(threads: readonly Thread[]): string[] {
  return [...new Set(threads.filter(t => t.question !== undefined && t.about !== undefined && t.aboutLane !== undefined).map(t => t.about!))]
}

// The threads with every question settled whose ticket left the lane it was
// asked in. lanes: each watched ticket's lane now; null where jaira no
// longer knows it; a ticket missing from it is not known yet and stays.
export function moved(threads: readonly Thread[], lanes: ReadonlyMap<string, string | null>, turn: number): Thread[] {
  return threads.map(t => {
    if (t.question === undefined || t.about === undefined || t.aboutLane === undefined || !lanes.has(t.about)) return t
    return lanes.get(t.about) === t.aboutLane ? t : answered(t, turn)
  })
}

// The questions open before a change and gone after it, settled at now:
// answered, dismissed, moved on, or replaced by a new one on their thread.
// A thread missing from after was dropped to keep the list short, not
// answered: its question settles nothing. Their options go: the pane no
// longer offers them.
export function settled(before: readonly Thread[], after: readonly Thread[], now: number): Ask[] {
  const kept = new Set(after.map(t => t.id))
  const still = new Set(after.filter(t => t.question !== undefined).map(t => `${t.id}\n${t.question}`))
  return asksOf(before)
    .filter(a => kept.has(a.thread) && !still.has(`${a.thread}\n${a.question}`))
    .map(a => ({ ...a, options: [], answeredAt: now }))
}

// The settled questions still shown: the last 24 hours, settled last first.
export function recent(done: readonly Ask[], now: number): Ask[] {
  return done.filter(a => a.answeredAt !== undefined && now - a.answeredAt < ANSWERED_MS).sort((a, b) => b.answeredAt! - a.answeredAt!)
}

const two = (n: number) => String(n).padStart(2, '0')

// A moment as the pane says it: 'today 07:12', another day '23.11. 21:40',
// in local time.
export function whenOf(at: number, now: number): string {
  const d = new Date(at)
  const n = new Date(now)
  const time = `${two(d.getHours())}:${two(d.getMinutes())}`
  return d.toDateString() === n.toDateString() ? `today ${time}` : `${two(d.getDate())}.${two(d.getMonth() + 1)}. ${time}`
}

// Every text the pane shows of these questions that a translation changes:
// each question and each option, once.
export function askTexts(asks: readonly Ask[]): string[] {
  return [...new Set(asks.flatMap(a => [a.question, ...a.options]))]
}

// The cheap model's prompt for the pane's Translate: every text in one call,
// back as one JSON array of as many strings in the same order.
export function asksTranslatePrompt(language: string, texts: readonly string[]): string {
  return [
    `Translate each text below into ${language}, plainly. Each is a question an AI assistant put to a person, or one of the answers it offered.`,
    'Keep code identifiers, file paths, commands and ticket handles exactly as written. Add nothing a text does not say.',
    `Answer with a JSON array of exactly ${texts.length} strings, in the same order, and nothing else.`,
    '',
    JSON.stringify(texts),
  ].join('\n')
}

// The threads with this thread's question taken as answered, at this turn.
export function answer(threads: readonly Thread[], thread: number, turn: number): Thread[] {
  return threads.map(t => (t.id === thread ? answered(t, turn) : t))
}
