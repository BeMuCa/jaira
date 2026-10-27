export interface Ticket {
  id: string;
  handle: string;
  title: string;
  lane: string;
  done: boolean;
  // false: in the epic only because it follows, or is a child of, a member.
  direct: boolean;
  // Where it lies; a filed ticket (logbook, archive) is shown, never dragged.
  place: Place;
  goal?: string;
  tags?: string[];
  creator?: string;
  assignee?: string;
  updated_by?: string;
  created_at?: string;
  updated_at?: string;
  // The day it left the board, when that is known.
  filed_at?: string;
  // The kind of no, for a ticket 'jaira reject' archived.
  rejected?: string;
}

export type Place = 'board' | 'ref' | 'logbook' | 'archive';

// Why nobody will do a ticket, as 'jaira reject' wrote it.
export interface Rejection {
  kind: string;
  why: string;
  by?: string;
  at?: string;
}

// One item of the definition of done.
export interface Check {
  text: string;
  done: boolean;
}

// One progress note; text is markdown as its author wrote it.
export interface Note {
  at: string;
  author: string;
  text: string;
}

// A ticket opened on the page: what 'jaira show' says about it, read only.
export interface TicketDetail {
  id: string;
  handle: string;
  title: string;
  lane: string;
  place: Place;
  goal: string;
  context: string;
  dod: Check[];
  notes: Note[];
  rejected?: Rejection;
}
