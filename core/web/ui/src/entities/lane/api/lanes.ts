import { request } from '@/shared/api';
import type { Lane } from '../model/types';

export const fetchLanes = () => request<Lane[]>('/api/lanes');
