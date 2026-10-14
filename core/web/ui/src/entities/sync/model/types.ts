// What 'jaira web' knows about its last round trips to the remote.
export interface SyncState {
  last_fetch?: string;
  fetch_error?: string;
  last_push?: string;
  push_error?: string;
}
