<script lang="ts">
  import { dateTime, renderMarkdown } from '@/shared/lib';
  import type { Note } from '../model/types';

  let { note }: { note: Note } = $props();

  // renderMarkdown escapes raw HTML and sanitises what is left, which is what
  // makes {@html} safe here.
  const html = $derived(renderMarkdown(note.text));
</script>

<li class="note">
  <header><time datetime={note.at}>{dateTime(note.at)}</time> · <span class="author">{note.author}</span></header>
  <div class="md">{@html html}</div>
</li>

<style>
  .note { list-style: none; padding: 0.5rem 0; border-top: 1px solid var(--line); }
  header { font-size: 0.8rem; color: var(--muted); }
  .author { font-weight: 600; }
  .md :global(p) { margin: 0.3rem 0; }
  .md :global(pre) { overflow-x: auto; padding: 0.5rem; background: var(--card); border-radius: 4px; }
  .md :global(code) { font-family: ui-monospace, monospace; font-size: 0.9em; }
  .md { overflow-wrap: anywhere; }
</style>
