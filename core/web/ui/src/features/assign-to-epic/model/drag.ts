// A ticket travels in a drag as its id and the epic it is dragged out of
// (empty from the left column), under a type of our own, so a drop of
// anything else — text, a file, a link — is not mistaken for a ticket.
export const TICKET_MIME = 'application/x-jaira-ticket';

export interface DraggedTicket {
  id: string;
  from: string;
}

export function startTicketDrag(e: DragEvent, id: string, from = '') {
  e.dataTransfer?.setData(TICKET_MIME, JSON.stringify({ id, from } satisfies DraggedTicket));
  if (e.dataTransfer) e.dataTransfer.effectAllowed = 'move';
}

export function carriesTicket(e: DragEvent): boolean {
  return e.dataTransfer?.types.includes(TICKET_MIME) ?? false;
}

export function droppedTicket(e: DragEvent): DraggedTicket | null {
  try {
    const t = JSON.parse(e.dataTransfer?.getData(TICKET_MIME) ?? '');
    return typeof t?.id === 'string' && t.id ? { id: t.id, from: String(t.from ?? '') } : null;
  } catch {
    return null;
  }
}
