// "12:04:31", or "never" for a zero or missing timestamp.
export function clock(iso: string | undefined): string {
  if (!iso) return 'never';
  const d = new Date(iso);
  if (Number.isNaN(d.getTime()) || d.getFullYear() < 2) return 'never';
  return d.toLocaleTimeString();
}
