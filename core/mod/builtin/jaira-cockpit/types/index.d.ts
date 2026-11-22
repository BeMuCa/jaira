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
  working: OverviewRow[]; queue: OverviewRow[]; parked: OverviewRow[]; history: { day: string; events: OverviewEvent[] }[]; lanes: string[]
}

// What the Workers tab draws: these mirror hooks/overview.ts.
export type Worker = { id: string; role: string; state: 'working' | 'idle' | 'question'; lane: string; title: string; open: number }

declare module 'claude-code' {
  interface PluginState {
    'jaira-cockpit': {
      // null: this project has no board, or no jaira on PATH.
      board: BoardView | null
      cache: CacheState
      threads: ThreadTally
      // The ticket the Give back pane asks "what is missing?" about.
      giving: string | null
      // The review notes that pane lists, each with what the person chose
      // for it (fix now, follow-up, accept).
      choosing: { id: string; items: { text: string; cls: 'blocker' | 'edge' | 'limit' | 'elsewhere'; choice: 'now' | 'later' | 'accept'; earlier?: 'now' | 'later' | 'accept'; follow?: string }[] } | null
      // The ticket the pane shows.
      paneTicket: string | null
      // Moves once a minute, so the cache countdown redraws.
      minute: number
      // A newer jaira release the band offers Update and Dismiss for, with
      // the version installed; null when there is none or it was dismissed.
      release: { installed: string; latest: string } | null
      // What the update pane shows before anything is written: the versions
      // and what 'jaira update --roles --lanes --dry-run' said.
      upgrade: { from: string; to: string; preview: string } | null
      // The last thing a button did, shown under the band until the next.
      notice: string | null
      // What the message pane shows: an error, or an explanation, whole.
      message: { title: string; text: string } | null
      // The explanation the explanation pane shows, with its ticket, and
      // the ticket the own-answer pane asks about.
      explaining: { id: string; text: string; decided?: string } | null
      answering: string | null
      // The ticket the Verdict pane shows: its goal, verdict and every note.
      verdictTicket: string | null
      // The review notes of each ticket in the person's language, original
      // text to the text shown; a note without an entry is shown as written.
      rephrased: Record<string, Record<string, string>>
      // The newest bug hunt of each ticket the band shows, as jaira's brief
      // says it stands: ok, fail, unknown or skipped.
      hunts: Record<string, 'ok' | 'fail' | 'unknown' | 'skipped'>
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
      // The Overview pane's tab, its folded sections and the people its Log
      // is filtered to, kept for the session across closing and opening it.
      overviewTab: 'queue' | 'log' | 'workers'
      overviewFolded: Record<string, boolean>
      overviewPeople: string[]
      overviewLanes: string[]
      // The Workers tab: its agents, read while it is open; each ticket's
      // colour and hat, kept while it stays on the board; when each agent
      // (<id>/<kind>) was first seen; the first line of each ticket's last
      // report; the frame the stone last lost a chunk (-1: never); and the
      // frame the animation is at.
      workers: Worker[] | null
      workersColors: Record<string, string>
      workersHats: Record<string, string>
      workersSince: Record<string, number>
      workersReports: Record<string, string>
      workersFell: number
      workersFrame: number
    }
  }
}
