// "12:04:31", or "never" before the first fetch (the server leaves the field out).
export function clock(iso: string | undefined): string {
  if (!iso) return 'never';
  return new Date(iso).toLocaleTimeString();
}

// "14.10.2026, 08:46" in the reader's locale and time zone.
export function dateTime(iso: string): string {
  return new Date(iso).toLocaleString(undefined, { dateStyle: 'short', timeStyle: 'short' });
}

// "14.10.26" in the reader's locale; empty when the server sent no date.
export function date(iso: string | undefined): string {
  if (!iso) return '';
  return new Date(iso).toLocaleDateString(undefined, { dateStyle: 'short' });
}

// How long ago iso was, coarse: "3d" from a day on, "5h" below that, "<1h"
// under an hour. Empty when the server sent no date.
export function age(iso: string | undefined, now: number = Date.now()): string {
  if (!iso) return '';
  const hours = Math.floor((now - new Date(iso).getTime()) / 3_600_000);
  if (hours >= 24) return `${Math.floor(hours / 24)}d`;
  if (hours >= 1) return `${hours}h`;
  return '<1h';
}
