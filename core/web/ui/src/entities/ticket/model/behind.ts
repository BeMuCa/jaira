// What the page says about a ticket whose file in this checkout is older than
// its ref: the card and the details show the ref, and this says why.
export function behindRef(handle: string): string {
  return `The ticket file in this checkout is older than its ref; this shows the ref. 'jaira pull ${handle}' brings the file up to it.`;
}
