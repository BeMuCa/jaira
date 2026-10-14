import { request } from '@/shared/api';
import type { Ticket, TicketDetail } from '../model/types';

export const fetchUnassigned = () => request<Ticket[]>('/api/unassigned');

export const fetchTicket = (id: string) => request<TicketDetail>(`/api/tickets/${encodeURIComponent(id)}`);
