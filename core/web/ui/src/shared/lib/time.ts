// "12:04:31", or "never" before the first fetch (the server leaves the field out).
export function clock(iso: string | undefined): string {
  if (!iso) return 'never';
  return new Date(iso).toLocaleTimeString();
}
