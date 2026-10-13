// Self-contained, as a contract must be: these mirror hooks/board.ts.
export type Waiting = {
  id: string; title: string; lane: string; verdict: string; question: string; acceptTo: string | null; assignee: string; explainable: boolean
}
export type Counts = { worker: number; review: number; person: number }

// What the cockpit's band and pane draw from, held by the host for the
// session so a hot reload keeps it.
export type BoardView = { waiting: Waiting[]; theirs: Waiting[]; counts: Counts }
export type CacheState = { tokens: number; lastAnswerAt: number | null; ttl: number }
export type ThreadTally = { questions: number; stale: number }
// What the Overview pane draws: these mirror hooks/overview.ts.
export type OverviewRow = { id: string; lane: string; title: string; who: string; reason: string }
export type OverviewEvent = OverviewRow & { kind: 'created' | 'parked'; at: string }
export type Overview = {
  working: OverviewRow[]; queue: OverviewRow[]; parked: OverviewRow[]; history: { day: string; events: OverviewEvent[] }[]
}

declare module 'claude-code' {
  interface PluginState {
    'jaira-cockpit': {
      // null: this project has no board, or no jaira on PATH.
      board: BoardView | null
      cache: CacheState
      threads: ThreadTally
      // The ticket the Give back pane asks "what is missing?" about.
      giving: string | null
      // The ticket the pane shows.
      paneTicket: string | null
      // Moves once a minute, so the cache countdown redraws.
      minute: number
      // The last thing a button did, shown under the band until the next.
      notice: string | null
      // What the message pane shows: an error, or an explanation, whole.
      message: { title: string; text: string } | null
      // The explanation the explanation pane shows, with its ticket, and
      // the ticket the own-answer pane asks about.
      explaining: { id: string; text: string; decided?: string } | null
      answering: string | null
      // The explanation an Explain chat wrote onto each ticket the band shows,
      // and the decision that answered it, when one did.
      explained: Record<string, { text: string; decided?: string }>
      // The decisions Discuss chats wrote back, by ticket, until read.
      decided: Record<string, string>
      // What the person used: the five-hour and seven-day windows, or cost.
      used: string
      // Whether the band shows every waiting ticket, not the first three.
      expanded: boolean
      // Whether the band shows the teammates' waiting tickets one by one.
      theirsOpen: boolean
      // The tickets whose Accept, Give back or Explain was pressed and has
      // not yet done its part, each with when, so a late end frees only its own.
      busy: Record<string, number>
      // The board the Overview pane shows, read when its button was pressed;
      // null when that read failed.
      overview: Overview | null
      // The Overview pane's tab and its folded sections, kept for the
      // session across closing and opening the pane.
      overviewTab: 'queue' | 'log'
      overviewFolded: Record<string, boolean>
    }
  }
}
