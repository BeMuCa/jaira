<script lang="ts">
  import { onMount } from 'svelte';
  import { fetchEpics, type Epic } from '@/entities/epic';
  import { fetchTicket, fetchUnassigned, type Ticket, type TicketDetail } from '@/entities/ticket';
  import { fetchSync, type SyncState } from '@/entities/sync';
  import { assignToEpic, removeFromEpic } from '@/features/assign-to-epic';
  import { createEpic } from '@/features/create-epic';
  import { EpicList } from '@/widgets/epic-list';
  import { UnassignedList } from '@/widgets/unassigned-list';
  import { SyncStatus } from '@/widgets/sync-status';
  import { TicketDetails } from '@/widgets/ticket-details';
  import { clock } from '@/shared/lib';

  // The server fetches once a minute; the page rereads on the same beat, so
  // what a teammate changed shows up without a reload.
  const REFRESH_MS = 60_000;
  // How often the page asks whether 'jaira web' still answers. Short, so a
  // page whose command is gone says so before anybody drags onto it.
  const HEARTBEAT_MS = 10_000;

  let epics = $state<Epic[]>([]);
  let unassigned = $state<Ticket[]>([]);
  let sync = $state<SyncState>({});
  let error = $state('');
  // When the page last failed to reach the command; null while it answers.
  let lostSince = $state<string | null>(null);
  let opened = $state<TicketDetail | null>(null);
  let dialog = $state<HTMLDialogElement>();

  async function load() {
    try {
      [epics, unassigned, sync] = await Promise.all([fetchEpics(), fetchUnassigned(), fetchSync()]);
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

  // On the way back the page rereads everything: what it shows is from
  // before the gap, and the command may be a new one on a fresh board.
  async function heartbeat() {
    try {
      sync = await fetchSync();
    } catch {
      lostSince ??= new Date().toISOString();
      return;
    }
    if (lostSince) {
      lostSince = null;
      await load();
    }
  }

  onMount(() => {
    load();
    const t = setInterval(load, REFRESH_MS);
    const h = setInterval(heartbeat, HEARTBEAT_MS);
    return () => {
      clearInterval(t);
      clearInterval(h);
    };
  });
</script>

<header class="top">
  <h1>jaira</h1>
  <SyncStatus {sync} />
</header>
{#if lostSince}
  <p class="lost" role="alert">
    jaira web has not answered since {clock(lostSince)} — is the command still running? What you see may be out of date; retrying every {HEARTBEAT_MS / 1000} s.
  </p>
{/if}
{#if error}<p class="error" role="alert">{error}</p>{/if}
<main>
  <UnassignedList tickets={unassigned} onopen={open} onremove={(epic, id) => change(() => removeFromEpic(epic, id))} />
  <EpicList {epics} onopen={open} oncreate={create} onassign={(epic, id, from) => change(() => assignToEpic(epic, id, from))} />
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
  .lost {
    margin: 1rem 0 0;
    padding: 0.6rem 0.9rem;
    color: var(--error);
    border: 1px solid var(--error);
    border-radius: 8px;
  }
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
