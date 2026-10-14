import { request } from '@/shared/api';
import type { SyncState } from '@/entities/sync';

const tickets = (epic: string) => `/api/epics/${encodeURIComponent(epic)}/tickets`;

// Writes the epic's ref and pushes it; with from set it takes the ticket out
// of that epic too, a move. Answers with the sync record, whose push_error
// says whether the push reached the remote; a failed write throws with the
// server's message, which names the step that failed.
export const assignToEpic = (epic: string, id: string, from = '') =>
  request<SyncState>(tickets(epic), { method: 'POST', body: JSON.stringify({ id, from }) });

export const removeFromEpic = (epic: string, id: string) =>
  request<SyncState>(`${tickets(epic)}/${encodeURIComponent(id)}`, { method: 'DELETE' });
