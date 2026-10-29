// One lane of the board, in the board's order. kind is what the colour
// marks beside the position.
export interface Lane {
  id: string;
  name: string;
  kind: 'work' | 'human' | 'terminal' | 'parking';
}
