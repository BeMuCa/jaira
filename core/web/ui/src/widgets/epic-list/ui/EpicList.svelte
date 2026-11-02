<script lang="ts">
  import { EpicProgress, type Epic } from '@/entities/epic';
  import { TicketCard } from '@/entities/ticket';
  import { DropTarget, startTicketDrag } from '@/features/assign-to-epic';
  import { addQuestion, closeQuestion, visible } from '@/features/close-epic';
  import { CreateEpicForm } from '@/features/create-epic';
  import { EditEpicForm, type EpicEdit } from '@/features/edit-epic';
  import { renderMarkdown } from '@/shared/lib';

  let { epics, expanded = false, onassign, onclose, onreopen, oncreate, onedit, onopen }: {
    epics: Epic[];
    // Open every epic's ticket list: a filter's matches must not hide folded.
    // Unset again, each epic is back as the person left it.
    expanded?: boolean;
    // from is the epic the ticket was dragged out of, empty from the left.
    // force is set only after the person confirmed a drop into a closed epic.
    onassign: (epic: string, id: string, from: string, force: boolean) => void;
    // force is set only after the person confirmed closing with unfinished tickets.
    onclose: (epic: string, force: boolean) => void;
    onreopen: (epic: string) => void;
    oncreate: (name: string, title: string, description: string) => Promise<void>;
    // Resolves when the edit is written; throws with the server's message.
    onedit: (name: string, edit: EpicEdit) => Promise<void>;
    onopen: (id: string) => void;
  } = $props();

  // Which epics the person unfolded by hand, by name. Only a toggle made
  // while nothing is forced open counts: the page's own opening is not theirs.
  let unfolded = $state<Record<string, boolean>>({});
  // The one epic whose editor is open, by name.
  let editing = $state('');
  // Closed epics are finished features: hidden unless the person asks.
  let showClosed = $state(false);
  const closedCount = $derived(epics.filter((e) => e.closed).length);
  const listed = $derived(visible(epics, showClosed));

  // A question answered no writes nothing; answered yes, the write goes with
  // force, past the refusal the question stood for.
  function drop(epic: Epic, id: string, from: string) {
    if (from === epic.name) return;
    const q = addQuestion(epic);
    if (q && !confirm(q)) return;
    onassign(epic.name, id, from, q !== null);
  }

  function close(epic: Epic) {
    const q = closeQuestion(epic);
    if (q && !confirm(q)) return;
    onclose(epic.name, q !== null);
  }
</script>

<section>
  <h2>Epics <small>{listed.length}</small></h2>
  {#if closedCount > 0}
    <label class="closed-toggle"><input type="checkbox" bind:checked={showClosed} /> Show closed ({closedCount})</label>
  {/if}
  <CreateEpicForm {oncreate} />
  {#if listed.length === 0}
    <p class="empty">{closedCount > 0 ? 'No open epics.' : 'No epics yet.'}</p>
  {/if}
  {#each listed as epic (epic.name)}
    <DropTarget ondropticket={(t) => drop(epic, t.id, t.from)}>
      <article class:closed={epic.closed}>
        <header>
          <h3>{epic.title}{#if epic.closed} <span class="badge">closed</span>{/if}</h3>
          <EpicProgress done={epic.done} signoff={epic.signoff} total={epic.total} />
        </header>
        {#if editing === epic.name}
          <EditEpicForm
            {epic}
            oncancel={() => (editing = '')}
            onsave={async (edit) => {
              await onedit(epic.name, edit);
              editing = '';
            }}
          />
        {:else}
          {#if epic.description}<p class="desc">{epic.description}</p>{/if}
          {#if epic.external?.length}
            <p class="keys">
              {#each epic.external as x (x.key)}
                <!-- link is set by the server only for an http(s) url. -->
                {#if x.link}<a href={x.link} target="_blank" rel="noopener noreferrer">{x.key}</a>{:else}<span title={x.url}>{x.key}</span>{/if}
              {/each}
            </p>
          {/if}
          {#if epic.text}
            <!-- Folded by default: the description line says what it is, the text says it at length. -->
            <details class="text">
              <summary>Text</summary>
              <!-- renderMarkdown escapes raw HTML and sanitises what is left. -->
              <div class="md">{@html renderMarkdown(epic.text)}</div>
            </details>
          {/if}
          <div class="actions">
            <button class="edit" onclick={() => (editing = epic.name)}>Edit</button>
            {#if epic.closed}
              <button class="edit" onclick={() => onreopen(epic.name)}>Reopen</button>
            {:else}
              <button class="edit" onclick={() => close(epic)}>Close</button>
            {/if}
          </div>
        {/if}
        <details
          open={expanded || !!unfolded[epic.name]}
          ontoggle={(e) => !expanded && (unfolded[epic.name] = e.currentTarget.open)}
        >
          <summary>{epic.tickets.length} tickets</summary>
          <ul>
            {#each epic.tickets as ticket (ticket.id)}
              <!-- Only a member line can be dragged: a follow-up has none to move. -->
              <TicketCard
                {ticket}
                {onopen}
                draggable={ticket.direct}
                ondragstart={ticket.direct ? (e) => startTicketDrag(e, ticket.id, epic.name) : undefined}
              />
            {/each}
          </ul>
        </details>
      </article>
    </DropTarget>
  {/each}
</section>

<style>
  section { display: grid; gap: 0.75rem; align-content: start; }
  article { padding: 0.75rem; border: 1px solid var(--line); border-radius: 8px; background: var(--panel); }
  header { display: grid; gap: 0.4rem; }
  h3 { margin: 0; font-size: 1rem; }
  .desc { margin: 0.5rem 0; white-space: pre-line; }
  .keys { display: flex; flex-wrap: wrap; gap: 0.4rem; margin: 0.3rem 0; font-size: 0.85rem; }
  .keys > * { border: 1px solid var(--line); border-radius: 4px; padding: 0 0.35rem; }
  ul { display: grid; gap: 0.3rem; padding: 0; margin: 0.4rem 0 0; }
  .text { margin: 0.4rem 0; }
  .md { margin-top: 0.3rem; overflow-wrap: anywhere; }
  .md :global(:is(h1, h2, h3)) { font-size: 0.95rem; margin: 0.6rem 0 0.3rem; }
  .md :global(p) { margin: 0.3rem 0; }
  .actions { display: flex; gap: 0.4rem; }
  .closed { opacity: 0.7; }
  .badge { font-size: 0.75rem; font-weight: normal; color: var(--muted); border: 1px solid var(--line); border-radius: 4px; padding: 0 0.3rem; }
  .closed-toggle { color: var(--muted); font-size: 0.85rem; }
  .edit { font: inherit; font-size: 0.8rem; color: var(--muted); background: none; border: 1px solid var(--line); border-radius: 4px; padding: 0.1rem 0.5rem; cursor: pointer; margin-bottom: 0.3rem; }
  summary { cursor: pointer; color: var(--muted); font-size: 0.85rem; }
  small { color: var(--muted); font-weight: normal; }
  .empty { color: var(--muted); }
</style>
