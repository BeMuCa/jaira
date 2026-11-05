import { request } from '@/shared/api';
import type { SyncState } from '@/entities/sync';

const epicPath = (name: string) => `/api/epics/${encodeURIComponent(name)}`;

// Writes "status: closed" onto refs/jaira/epics/<name> and pushes it, as
// 'jaira epic close' does. Without force an epic that still counts an
// unfinished ticket answers 409 and nothing is written.
export const closeEpic = (name: string, force = false) =>
  request<SyncState>(`${epicPath(name)}/close`, { method: 'POST', body: JSON.stringify({ force }) });

// Takes the closing back, as 'jaira epic reopen' does.
export const reopenEpic = (name: string) => request<SyncState>(`${epicPath(name)}/reopen`, { method: 'POST' });
