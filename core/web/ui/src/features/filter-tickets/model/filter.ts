import type { Ticket } from '@/entities/ticket';

// Which date the from-to range reads.
export type DateField = 'created' | 'updated' | 'filed';

// What the page filters by. Every field set narrows the list further: the
// filters are ANDed, as 'jaira list' ANDs its flags.
export interface Filter {
  // Substring of the title, the goal or the id, any case.
  q: string;
  lane: string;
  tag: string;
  // Creator, assignee or whoever changed it last.
  person: string;
  date: DateField;
  // yyyy-mm-dd, both ends included; empty is open.
  from: string;
  to: string;
  // Filed tickets too, not only the board.
  all: boolean;
}

export const emptyFilter = (): Filter => ({ q: '', lane: '', tag: '', person: '', date: 'created', from: '', to: '', all: false });

const DATES: DateField[] = ['created', 'updated', 'filed'];

// The filter lives in the address, so a link shows the same list.
export function fromQuery(search: string): Filter {
  const p = new URLSearchParams(search);
  const f = emptyFilter();
  f.q = p.get('q') ?? '';
  f.lane = p.get('lane') ?? '';
  f.tag = p.get('tag') ?? '';
  f.person = p.get('person') ?? '';
  const date = p.get('date') as DateField;
  if (DATES.includes(date)) f.date = date;
  f.from = p.get('from') ?? '';
  f.to = p.get('to') ?? '';
  f.all = p.get('all') === '1';
  return f;
}

// Only what is set goes into the address; an unfiltered page keeps a bare one.
export function toQuery(f: Filter): string {
  const p = new URLSearchParams();
  for (const k of ['q', 'lane', 'tag', 'person', 'from', 'to'] as const) {
    if (f[k]) p.set(k, f[k]);
  }
  if ((f.from || f.to) && f.date !== 'created') p.set('date', f.date);
  if (f.all) p.set('all', '1');
  const s = p.toString();
  return s ? `?${s}` : '';
}

// Whether anything narrows the list; the switch only widens it.
export const narrows = (f: Filter) => !!(f.q || f.lane || f.tag || f.person || f.from || f.to);

// The local day of an RFC 3339 time, comparable as a string with the
// range's yyyy-mm-dd.
function day(iso: string | undefined): string {
  if (!iso) return '';
  const d = new Date(iso);
  if (isNaN(d.getTime())) return '';
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
}

const dateOf = (t: Ticket, f: DateField) => day(f === 'created' ? t.created_at : f === 'updated' ? t.updated_at : t.filed_at);

export function matches(t: Ticket, f: Filter): boolean {
  if (f.q) {
    const q = f.q.toLowerCase();
    if (![t.title, t.goal ?? '', t.handle, t.id].some((s) => s.toLowerCase().includes(q))) return false;
  }
  if (f.lane && t.lane !== f.lane) return false;
  if (f.tag && !(t.tags ?? []).includes(f.tag)) return false;
  if (f.person && ![t.creator, t.assignee, t.updated_by].includes(f.person)) return false;
  if (f.from || f.to) {
    // A ticket with no such date (never filed, or filed before it was
    // recorded) is outside every range.
    const d = dateOf(t, f.date);
    if (!d || (f.from && d < f.from) || (f.to && d > f.to)) return false;
  }
  return true;
}

// The values the dropdowns offer: what the tickets on the page carry.
export function choices(tickets: Ticket[]) {
  const lanes = new Set<string>();
  const tags = new Set<string>();
  const people = new Set<string>();
  for (const t of tickets) {
    if (t.lane) lanes.add(t.lane);
    for (const g of t.tags ?? []) tags.add(g);
    for (const p of [t.creator, t.assignee, t.updated_by]) if (p) people.add(p);
  }
  const sorted = (s: Set<string>) => [...s].sort((a, b) => a.localeCompare(b));
  return { lanes: sorted(lanes), tags: sorted(tags), people: sorted(people) };
}
