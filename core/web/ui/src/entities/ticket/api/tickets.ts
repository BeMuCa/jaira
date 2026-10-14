import { request } from '@/shared/api';
import type { Ticket } from '../model/types';

export const fetchUnassigned = () => request<Ticket[]>('/api/unassigned');
