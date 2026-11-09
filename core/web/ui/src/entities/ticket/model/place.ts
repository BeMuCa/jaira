import type { Place } from './types';

// A filed ticket has left the board: it is read, never sorted.
export const isFiled = (place: Place) => place === 'logbook' || place === 'archive';

// Where a ticket lies, as the page says it.
export const PLACE_LABELS: Record<Place, string> = {
  board: 'on the board',
  ref: 'only on its ref',
  logbook: 'in the logbook',
  archive: 'in the archive',
};

// The lane a badge shows. The logbook holds only finished work, so its
// tickets say archive; an archived ticket keeps the lane it left the board
// from, which is what tells a dropped one from a finished one.
export const badgeLane = (t: { lane: string; place?: Place }) => (t.place === 'logbook' ? 'archive' : t.lane);

// How a ticket left the board, or that it has not: finished, put away
// unfinished, or rejected with 'jaira reject'. Each one looks different on a
// card.
export type Fate = 'open' | 'done' | 'left' | 'rejected';

export function fate(t: { done: boolean; place?: Place; rejected?: string }): Fate {
  if (t.place === 'archive' && t.rejected) return 'rejected';
  if (t.done) return 'done';
  return t.place && isFiled(t.place) ? 'left' : 'open';
}
