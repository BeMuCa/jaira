<script lang="ts">
  import { EpicProgress, type Epic } from '@/entities/epic';
  import { TicketCard } from '@/entities/ticket';
  import { DropTarget, startTicketDrag } from '@/features/assign-to-epic';
  import { CreateEpicForm } from '@/features/create-epic';
  import { EditEpicForm, type EpicEdit } from '@/features/edit-epic';
  import { renderMarkdown } from '@/shared/lib';

  let { epics, expanded = false, onassign, oncreate, onedit, onopen }: {
    epics: Epic[];
    // Open every epic's ticket list: a filter's matches must not hide folded.
    // Unset again, each epic is back as the person left it.
    expanded?: boolean;
    // from is the epic the ticket was dragged out of, empty from the left.
    onassign: (epic: string, id: string, from: string) => void;
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
</script>

<section>
  <h2>Epics <small>{epics.length}</small></h2>
  <CreateEpicForm {oncreate} />
  {#if epics.length === 0}
    <p class="empty">No epics yet.</p>
  {/if}
  {#each epics as epic (epic.name)}
    <DropTarget ondropticket={(t) => t.from !== epic.name && onassign(epic.name, t.id, t.from)}>
      <article>
        <header>
          <h3>{epic.title}</h3>
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
          {#if epic.text}
            <!-- Folded by default: the description line says what it is, the text says it at length. -->
            <details class="text">
              <summary>Text</summary>
              <!-- renderMarkdown escapes raw HTML and sanitises what is left. -->
              <div class="md">{@html renderMarkdown(epic.text)}</div>
            </details>
          {/if}
          <button class="edit" onclick={() => (editing = epic.name)}>Edit</button>
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
  ul { display: grid; gap: 0.3rem; padding: 0; margin: 0.4rem 0 0; }
  .text { margin: 0.4rem 0; }
  .md { margin-top: 0.3rem; overflow-wrap: anywhere; }
  .md :global(:is(h1, h2, h3)) { font-size: 0.95rem; margin: 0.6rem 0 0.3rem; }
  .md :global(p) { margin: 0.3rem 0; }
  .edit { font: inherit; font-size: 0.8rem; color: var(--muted); background: none; border: 1px solid var(--line); border-radius: 4px; padding: 0.1rem 0.5rem; cursor: pointer; margin-bottom: 0.3rem; }
  summary { cursor: pointer; color: var(--muted); font-size: 0.85rem; }
  small { color: var(--muted); font-weight: normal; }
  .empty { color: var(--muted); }
</style>
