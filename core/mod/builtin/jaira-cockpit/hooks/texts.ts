// Every text the cockpit shows a person, in English whatever language the
// person speaks with Claude: one set of words, the same on every machine.

export const t = {
  waiting: 'Waiting for you',
  accept: 'Accept',
  giveBack: 'Give back',
  explain: 'Explain',
  discuss: 'Discuss',
  more: (n: number) => `+${n} more`,
  reasonLabel: 'What is missing?',
  reasonPlaceholder: 'goes onto the ticket as a note',
  accepted: (id: string, lane: string) => `${id} accepted: ${lane}`,
  returnedTo: (id: string, who: 'old' | 'new' | 'none') =>
    who === 'none' ? `${id} back in in-progress; no worker started (the ticket's note says why)` : `${id} given back to ${who === 'old' ? 'the worker that carried it' : 'a new worker'}`,
  failed: (what: string) => `${what} failed`,
  close: 'Close',
  bandTitle: 'jaira-cockpit',
  noRoom: 'No room for a pane here: widen the terminal and press Give back again.',
  readExplanation: 'Read explanation',
  explanationOf: (id: string) => `Explanation of ${id}`,
  needsHerdr: 'A chat of its own needs Herdr (HERDR_ENV=1).',
  opened: (id: string) => `${id}: opened in a Herdr tab of its own`,
  context: (k: number) => `context ${k}k`,
  cacheLeft: (min: number) => `cache ${min} min left`,
  cacheCold: 'cache cold',
  hintCompact: '/compact advised',
  hintNew: 'start a new session',
  guardQuestion: (k: number, min: number, usd: string) => `The prompt cache went cold ${min} min ago. Sending writes ${k}k tokens to it again (about ${usd}). What now?`,
  guardSend: 'Send',
  guardCompact: 'Compact first',
  guardClear: 'Start over (/clear)',
  guardCancel: 'Cancel',
  guardHeader: 'Cold cache',
  threadsOpen: (q: number, s: number) => `threads: ${q} open question(s), ${s} left lying`,
  threadsNone: 'No conversation threads yet.',
  threadsTitle: 'Conversation threads (? open question to you, ~ left lying, -> ticket):',
  status: (w: number, r: number, p: number) => `${w} worker · ${r} review · ${p} waiting for you`,
  usage: 'session',
  paneTitle: (id: string) => `Ticket ${id}`,
  paneReports: 'Reports',
  paneUsage: 'Usage per agent',
  paneNothing: 'nothing recorded yet',
  usageText:
    'Usage: /cockpit <subcommand>\n\n  report [id]   full text of the agent reports shortened for a ticket, newest first; without an id, the tickets that have some\n  threads       the conversation threads of this session and where each stands\n  ticket <id>   a pane with the ticket\'s reports and what its agents used',
}

// The noise filter's own words: the digest line and /cockpit report.
export const filterWords = {
  // Leads every digest, so the model never takes the line for the agent's
  // own words.
  marker: '[shortened by jaira-cockpit]',
  question: 'question',
  fullReport: 'full report',
  none: 'No agent reports kept for this project yet.',
  noneFor: (id: string) => `No reports kept for ${id} in this project.`,
  list: 'Tickets with kept reports in this project, most recent first:',
  hint: '/cockpit report <id> shows one.',
  count: (n: number) => `${n} report(s)`,
  last: 'last',
  of: 'of',
}
