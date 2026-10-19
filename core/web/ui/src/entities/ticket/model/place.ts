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
