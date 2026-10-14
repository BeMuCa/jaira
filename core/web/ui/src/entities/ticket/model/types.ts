export interface Ticket {
  id: string;
  handle: string;
  title: string;
  lane: string;
  done: boolean;
  // false: in the epic only because it follows, or is a child of, a member.
  direct: boolean;
}
