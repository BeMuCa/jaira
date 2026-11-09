import { request } from '@/shared/api';
import type { SyncState } from '@/entities/sync';

export interface EpicFields {
  title: string;
  description: string;
  body: string;
}

// An edit names only the fields it changes, and base is what the editor read:
// a field somebody else changed since answers 409 and nothing is written.
// add_external puts keys on, each "<KEY> [<URL>]"; drop_external takes keys
// off. They are changes against what the editor read, never a new list.
export type EpicEdit = Partial<EpicFields> & { base: EpicFields; add_external?: string[]; drop_external?: string[] };

// Rewrites refs/jaira/epics/<name> and pushes it, as 'jaira epic set' does.
export const setEpic = (name: string, edit: EpicEdit) =>
  request<SyncState>(`/api/epics/${encodeURIComponent(name)}`, { method: 'PATCH', body: JSON.stringify(edit) });
