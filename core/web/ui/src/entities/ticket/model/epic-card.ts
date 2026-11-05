import type { Ticket } from './types';

// What a card inside an epic offers about leaving it. A member has a line in
// the epic, so it can be taken out — wherever it lies, filed away included,
// as 'jaira epic rm' takes it. A follow-up has no line of its own: it names
// the member it is counted through, the one whose line would have to go.
export type EpicCardExit = { removable: true } | { removable: false; via: string };

export function epicCardExit(t: Pick<Ticket, 'direct' | 'via'>): EpicCardExit {
  if (t.direct) return { removable: true };
  return { removable: false, via: t.via ?? '' };
}

// A ticket's handle: the last six characters of its id, as the CLI prints it.
export const handleOf = (id: string) => id.slice(-6);
