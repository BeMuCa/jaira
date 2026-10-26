import type { Lane } from './types';

// The board's lanes, loaded once by the page; every badge reads its colour
// from here.
const board = $state<{ lanes: Lane[] }>({ lanes: [] });

export function setLanes(lanes: Lane[]) {
  board.lanes = lanes;
}

// Hues (oklch) for the lanes where agents work, taken in the board's order,
// so a lane keeps its colour for as long as the board keeps its order. They
// alternate so neighbouring lanes never look alike, and stay clear of the
// three hues a lane's kind owns: red for parking, amber for a person's lane,
// green for the terminal one. The lightness comes from the theme
// (--lane-l, --lane-c in app/styles.css), so both themes read.
const WORK_HUES = [240, 165, 290, 110, 340, 215, 190, 265, 315];
const KIND_HUES = { parking: 25, human: 70, terminal: 145 } as const;

const oklch = (hue: number, chroma = 'var(--lane-c)') => `oklch(var(--lane-l) ${chroma} ${hue})`;

// The colour of a lane; a lane this board does not have is muted.
export function laneColour(id: string): string {
  let work = 0;
  for (const l of board.lanes) {
    if (l.id === id) {
      if (l.kind !== 'work') return oklch(KIND_HUES[l.kind]);
      // Past the palette a tenth lane is grey, then the hues come round again.
      const slot = work % (WORK_HUES.length + 1);
      return slot === WORK_HUES.length ? oklch(250, '0.02') : oklch(WORK_HUES[slot]);
    }
    if (l.kind === 'work') work++;
  }
  return 'var(--muted)';
}

// Lane ids in the board's order, unknown ones after them by name.
export function byBoardOrder(ids: string[]): string[] {
  const at = (id: string) => {
    const i = board.lanes.findIndex((l) => l.id === id);
    return i < 0 ? board.lanes.length : i;
  };
  return [...ids].sort((a, b) => at(a) - at(b) || a.localeCompare(b));
}
