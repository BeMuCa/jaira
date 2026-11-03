import type { Epic } from './types';

// The name shown small beside an epic's title: the name is what a person
// types into 'jaira epic ...', the title what they read. An epic whose title
// is its name shows it once, so this is empty then.
export const nameBeside = (e: Pick<Epic, 'name' | 'title'>): string => (e.title && e.title !== e.name ? e.name : '');
