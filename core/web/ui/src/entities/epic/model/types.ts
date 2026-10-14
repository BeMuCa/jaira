import type { Ticket } from '@/entities/ticket';

export interface Epic {
  name: string;
  title: string;
  description: string;
  done: number;
  total: number;
  tickets: Ticket[];
}
