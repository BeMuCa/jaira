import { TYPED_BY_COCKPIT } from './sessions.ts'

// The context and prompt-cache rules: how long the cache stays warm, what the
// band says about it, and when the guard holds a prompt. Pure; register.tsx
// reads the clock and the session's figures and hands them in.

// A subscription keeps the prompt cache an hour, an API key five minutes.
// `$.session.usage().rateLimits` is filled only on a subscription.
export const SUBSCRIPTION_TTL_MS = 60 * 60_000
export const API_KEY_TTL_MS = 5 * 60_000

export function ttlMs(rateLimits: readonly unknown[]): number {
  return rateLimits.length > 0 ? SUBSCRIPTION_TTL_MS : API_KEY_TTL_MS
}

// The subscription windows the status line shows, by the `kind` the engine
// gives them: `5h 36% · 7d 12%`, each only when it has a reading.
const WINDOWS: ReadonlyArray<[kind: string, label: string]> = [['five_hour', '5h'], ['seven_day', '7d']]

export function windowsUsed(rateLimits: ReadonlyArray<{ kind: string; percentUsed: number }>): string {
  return WINDOWS.flatMap(([kind, label]) => {
    const w = rateLimits.find(r => r.kind === kind)
    return w === undefined ? [] : [`${label} ${Math.round(w.percentUsed)}%`]
  }).join(' · ')
}

// Past these the band points at /compact, then at a new session.
export const COMPACT_AT = 100_000
export const NEW_SESSION_AT = 300_000

export type CacheView = {
  // Context size, in thousands of tokens.
  k: number
  // Whole minutes the cache has left, or undefined once it is cold (and
  // before the first answer, when there is nothing cached yet).
  minutesLeft?: number
  hint?: 'compact' | 'new'
}

// Every answer of the model reads the cache and so keeps it warm for another
// TTL: the clock runs from the last answer, not from the last write.
export function cacheView(tokens: number, now: number, lastAnswerAt: number | undefined, ttl: number): CacheView {
  const k = Math.round(tokens / 1000)
  const left = lastAnswerAt === undefined ? 0 : lastAnswerAt + ttl - now
  const hint = tokens >= NEW_SESSION_AT ? 'new' : tokens >= COMPACT_AT ? 'compact' : undefined
  return { k, ...(left > 0 ? { minutesLeft: Math.ceil(left / 60_000) } : {}), ...(hint ? { hint } : {}) }
}

// Who may be held: a person's prompt alone — typed, or through Remote
// Control. Never an agent's message, a notification, a scheduled prompt.
// Herdr types into a pane exactly as a person does (origin `composer` both
// ways), so what is typed for a person is told apart by its words: spawn.sh
// types slash commands (`/jaira-worker <id>`), and the cockpit's own text to a
// worker leads with `[jaira-cockpit]`. Its explain/discuss chats start a new
// session, which has no cache to guard.
const PERSON = new Set(['composer', 'bridge'])

export function isPersonsPrompt(origin: string, text: string): boolean {
  const lead = text.trimStart()
  return PERSON.has(origin) && !lead.startsWith('/') && !lead.startsWith(TYPED_BY_COCKPIT)
}

export type Guard = { coldMinutes: number; k: number; usd: number }

// Whether to hold this prompt, and what to say if so: the cache went cold
// and the context it would write again is at least `threshold` tokens.
export function guard(input: {
  origin: string
  text: string
  tokens: number
  now: number
  lastAnswerAt: number | undefined
  ttl: number
  threshold: number
  usdPerMTok: number
}): Guard | undefined {
  if (!isPersonsPrompt(input.origin, input.text)) return undefined
  if (input.lastAnswerAt === undefined || input.tokens < input.threshold) return undefined
  const coldFor = input.now - (input.lastAnswerAt + input.ttl)
  if (coldFor < 0) return undefined
  return {
    coldMinutes: Math.max(1, Math.round(coldFor / 60_000)),
    k: Math.round(input.tokens / 1000),
    // A cache write costs 1.25 times the input price.
    usd: (input.tokens * 1.25 * input.usdPerMTok) / 1_000_000,
  }
}

export function dollars(usd: number): string {
  return usd < 0.01 ? '<$0.01' : `$${usd.toFixed(2)}`
}
