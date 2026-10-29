import { request } from '@/shared/api';
import type { Ticket, TicketDetail } from '../model/types';

// all: every ticket in no epic that ever was, filed ones included.
export const fetchUnassigned = (all = false) => request<Ticket[]>(all ? '/api/unassigned?all=1' : '/api/unassigned');

export const fetchTicket = (id: string) => request<TicketDetail>(`/api/tickets/${encodeURIComponent(id)}`);
