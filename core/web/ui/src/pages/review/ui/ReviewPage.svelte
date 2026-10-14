<script lang="ts">
  import { onMount } from 'svelte';
  import { fetchEpics, type Epic } from '@/entities/epic';
  import { fetchUnassigned, type Ticket } from '@/entities/ticket';
  import { fetchSync, type SyncState } from '@/entities/sync';
  import { assignToEpic } from '@/features/assign-to-epic';
  import { EpicList } from '@/widgets/epic-list';
  import { UnassignedList } from '@/widgets/unassigned-list';
  import { SyncStatus } from '@/widgets/sync-status';

  // The server fetches once a minute; the page rereads on the same beat, so
  // what a teammate changed shows up without a reload.
  const REFRESH_MS = 60_000;

  let epics = $state<Epic[]>([]);
  let unassigned = $state<Ticket[]>([]);
  let sync = $state<SyncState>({});
  let error = $state('');

  async function load() {
    try {
      [epics, unassigned, sync] = await Promise.all([fetchEpics(), fetchUnassigned(), fetchSync()]);
      error = '';
    } catch (e) {
      error = (e as Error).message;
    }
  }

  async function assign(epic: string, id: string) {
    try {
      sync = await assignToEpic(epic, id);
      await load();
    } catch (e) {
      error = (e as Error).message;
    }
  }

  onMount(() => {
    load();
    const t = setInterval(load, REFRESH_MS);
    return () => clearInterval(t);
  });
</script>

<header class="top">
  <h1>jaira</h1>
  <SyncStatus {sync} />
</header>
{#if error}<p class="error" role="alert">{error}</p>{/if}
<main>
  <UnassignedList tickets={unassigned} />
  <EpicList {epics} onassign={assign} />
</main>

<style>
  .top { display: flex; justify-content: space-between; align-items: center; gap: 1rem; flex-wrap: wrap; }
  h1 { margin: 0; font-size: 1.25rem; }
  main { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; margin-top: 1rem; }
  @media (max-width: 720px) { main { grid-template-columns: 1fr; } }
  .error { color: var(--error); }
</style>
