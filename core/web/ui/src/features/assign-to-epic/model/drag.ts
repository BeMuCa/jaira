// A ticket travels in a drag as its id under a type of our own, so a drop of
// anything else — text, a file, a link — is not mistaken for a ticket.
export const TICKET_MIME = 'application/x-jaira-ticket';

export function startTicketDrag(e: DragEvent, id: string) {
  e.dataTransfer?.setData(TICKET_MIME, id);
  if (e.dataTransfer) e.dataTransfer.effectAllowed = 'move';
}

export function carriesTicket(e: DragEvent): boolean {
  return e.dataTransfer?.types.includes(TICKET_MIME) ?? false;
}

export function droppedTicket(e: DragEvent): string {
  return e.dataTransfer?.getData(TICKET_MIME) ?? '';
}
