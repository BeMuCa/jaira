import { request } from '@/shared/api';
import type { SyncState } from '@/entities/sync';

export interface EpicFields {
  title: string;
  description: string;
  body: string;
}

// An edit names only the fields it changes, and base is what the editor read:
// a field somebody else changed since answers 409 and nothing is written.
export type EpicEdit = Partial<EpicFields> & { base: EpicFields };

// Rewrites refs/jaira/epics/<name> and pushes it, as 'jaira epic set' does.
export const setEpic = (name: string, edit: EpicEdit) =>
  request<SyncState>(`/api/epics/${encodeURIComponent(name)}`, { method: 'PATCH', body: JSON.stringify(edit) });
