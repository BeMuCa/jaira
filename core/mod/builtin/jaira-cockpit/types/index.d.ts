// Self-contained, as a contract must be: these mirror hooks/board.ts and
// hooks/lang.ts.
export type Waiting = {
  id: string; title: string; lane: string; verdict: string; question: string; acceptTo: string | null
}
export type Counts = { worker: number; review: number; person: number }
export type Lang = 'en' | 'de' | 'ru'

// What the cockpit's band and pane draw from, held by the host for the
// session so a hot reload keeps it.
export type BoardView = { waiting: Waiting[]; counts: Counts }
export type CacheState = { tokens: number; lastAnswerAt: number | null; ttl: number }
export type ThreadTally = { questions: number; stale: number }

declare module 'claude-code' {
  interface PluginState {
    'jaira-cockpit': {
      // null: this project has no board, or no jaira on PATH.
      board: BoardView | null
      cache: CacheState
      threads: ThreadTally
      // The ticket whose "what is missing?" field is open in the band.
      giving: string | null
      // The ticket the pane shows.
      paneTicket: string | null
      lang: Lang
      // Moves once a minute, so the cache countdown redraws.
      minute: number
      // The last thing a button did, shown under the band until the next.
      notice: string | null
    }
  }
}
