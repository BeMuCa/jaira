export interface Ticket {
  id: string;
  handle: string;
  title: string;
  lane: string;
  done: boolean;
  // false: in the epic only because it follows, or is a child of, a member.
  direct: boolean;
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
  goal: string;
  context: string;
  dod: Check[];
  notes: Note[];
}
