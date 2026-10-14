import { request } from '@/shared/api';
import type { SyncState } from '../model/types';

export const fetchSync = () => request<SyncState>('/api/sync');
export const syncNow = () => request<SyncState>('/api/sync', { method: 'POST' });
