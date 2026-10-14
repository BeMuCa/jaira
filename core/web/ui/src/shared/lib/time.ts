// "12:04:31", or "never" before the first fetch (the server leaves the field out).
export function clock(iso: string | undefined): string {
  if (!iso) return 'never';
  return new Date(iso).toLocaleTimeString();
}

// "14.10.2026, 08:46" in the reader's locale and time zone.
export function dateTime(iso: string): string {
  return new Date(iso).toLocaleString(undefined, { dateStyle: 'short', timeStyle: 'short' });
}
