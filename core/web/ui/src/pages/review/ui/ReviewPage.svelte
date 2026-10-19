<script lang="ts">
  import { onMount, untrack } from 'svelte';
  import { fetchEpics, type Epic } from '@/entities/epic';
  import { fetchTicket, fetchUnassigned, type Ticket, type TicketDetail } from '@/entities/ticket';
  import { fetchSync, type SyncState } from '@/entities/sync';
  import { assignToEpic, removeFromEpic } from '@/features/assign-to-epic';
  import { createEpic } from '@/features/create-epic';
  import { FilterBar, choices, fromQuery, matches, narrows, toQuery } from '@/features/filter-tickets';
  import { EpicList } from '@/widgets/epic-list';
  import { UnassignedList } from '@/widgets/unassigned-list';
  import { SyncStatus } from '@/widgets/sync-status';
  import { TicketDetails } from '@/widgets/ticket-details';

  // The server fetches once a minute; the page rereads on the same beat, so
  // what a teammate changed shows up without a reload.
  const REFRESH_MS = 60_000;

  let epics = $state<Epic[]>([]);
  let unassigned = $state<Ticket[]>([]);
  let sync = $state<SyncState>({});
  let error = $state('');
  let opened = $state<TicketDetail | null>(null);
  let dialog = $state<HTMLDialogElement>();
  // The filter starts from the address and is written back to it, so a link
  // opens the same list.
  let filter = $state(fromQuery(location.search));

  const options = $derived(choices([...unassigned, ...epics.flatMap((e) => e.tickets)]));
  const shown = $derived(unassigned.filter((t) => matches(t, filter)));
  // An epic none of whose tickets match has nothing to show while a filter
  // narrows the page; unfiltered, every epic stays, empty ones too.
  const shownEpics = $derived(
    epics
      .map((e) => ({ ...e, tickets: e.tickets.filter((t) => matches(t, filter)) }))
      .filter((e) => !narrows(filter) || e.tickets.length > 0),
  );

  $effect(() => {
    const q = toQuery(filter);
    if (q !== location.search) history.replaceState(null, '', location.pathname + q + location.hash);
  });

  // Loads at once, and again whenever the switch flips: it changes what the
  // server is asked for.
  $effect(() => {
    void filter.all;
    untrack(load);
  });

  async function load() {
    try {
      [epics, unassigned, sync] = await Promise.all([fetchEpics(), fetchUnassigned(filter.all), fetchSync()]);
      error = '';
    } catch (e) {
      error = (e as Error).message;
    }
  }

  // A write that failed still reloads: half of a move may have landed, and
  // the page has to show where the ticket actually is beside the message.
  async function change(write: () => Promise<SyncState>) {
    let failed = '';
    try {
      sync = await write();
    } catch (e) {
      failed = (e as Error).message;
    }
    await load();
    if (failed) error = failed;
  }

  async function create(name: string, title: string, description: string) {
    sync = await createEpic(name, title, description);
    await load();
  }

  async function open(id: string) {
    try {
      opened = await fetchTicket(id);
      dialog?.showModal();
    } catch (e) {
      error = (e as Error).message;
    }
  }

  onMount(() => {
    const t = setInterval(load, REFRESH_MS);
    return () => clearInterval(t);
  });
</script>

<header class="top">
  <h1>jaira</h1>
  <SyncStatus {sync} />
</header>
<FilterBar bind:filter lanes={options.lanes} tags={options.tags} people={options.people} />
{#if error}<p class="error" role="alert">{error}</p>{/if}
<main>
  <UnassignedList
    tickets={shown}
    empty={narrows(filter) ? 'No ticket in no epic matches the filter.' : 'Every open ticket belongs to an epic.'}
    onopen={open}
    onremove={(epic, id) => change(() => removeFromEpic(epic, id))}
  />
  <EpicList epics={shownEpics} onopen={open} oncreate={create} onassign={(epic, id, from) => change(() => assignToEpic(epic, id, from))} />
</main>

<dialog bind:this={dialog} onclose={() => (opened = null)}>
  {#if opened}
    <TicketDetails ticket={opened} onclose={() => dialog?.close()} />
  {/if}
</dialog>

<style>
  .top { display: flex; justify-content: space-between; align-items: center; gap: 1rem; flex-wrap: wrap; }
  h1 { margin: 0; font-size: 1.25rem; }
  main { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; margin-top: 1rem; }
  @media (max-width: 720px) { main { grid-template-columns: 1fr; } }
  .error { color: var(--error); }
  dialog {
    width: min(48rem, calc(100vw - 2rem));
    max-height: calc(100vh - 4rem);
    padding: 1rem 1.25rem;
    color: inherit;
    background: var(--panel);
    border: 1px solid var(--line);
    border-radius: 10px;
  }
  dialog::backdrop { background: rgb(0 0 0 / 0.4); }
</style>
