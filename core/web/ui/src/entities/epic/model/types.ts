import type { Ticket } from '@/entities/ticket';

export interface Epic {
  name: string;
  title: string;
  description: string;
  done: number;
  // In signoff: finished, waiting for a person to accept it; not in done.
  signoff: number;
  total: number;
  tickets: Ticket[];
}
