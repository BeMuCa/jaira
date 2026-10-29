import type { Ticket } from '@/entities/ticket';

export interface Epic {
  name: string;
  title: string;
  description: string;
  // body is the markdown below title and description as the file has it, for
  // the editor; text is the same as a person reads it, without the heading
  // that repeats the title and without comments.
  body: string;
  text: string;
  done: number;
  // In signoff: finished, waiting for a person to accept it; not in done.
  signoff: number;
  total: number;
  tickets: Ticket[];
}
