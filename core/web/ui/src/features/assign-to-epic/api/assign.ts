import { request } from '@/shared/api';
import type { SyncState } from '@/entities/sync';

// Writes the epic's ref and pushes it; answers with the sync record, whose
// push_error says whether the push reached the remote.
export const assignToEpic = (epic: string, id: string) =>
  request<SyncState>(`/api/epics/${encodeURIComponent(epic)}/tickets`, {
    method: 'POST',
    body: JSON.stringify({ id }),
  });
