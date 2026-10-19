import type { Place } from './types';

// A filed ticket has left the board: it is read, never sorted.
export const isFiled = (place: Place) => place === 'logbook' || place === 'archive';
