import { request } from '@/shared/api';
import type { Epic } from '../model/types';

export const fetchEpics = () => request<Epic[]>('/api/epics');
