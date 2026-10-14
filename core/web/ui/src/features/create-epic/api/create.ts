import { request } from '@/shared/api';
import type { SyncState } from '@/entities/sync';

// Writes refs/jaira/epics/<name> and pushes it, as 'jaira epic create' does;
// the server normalises the name ("Pay Later" becomes "pay-later").
export const createEpic = (name: string, title: string, description: string) =>
  request<SyncState>('/api/epics', { method: 'POST', body: JSON.stringify({ name, title, description }) });
