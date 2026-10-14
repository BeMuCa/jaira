import DOMPurify from 'dompurify';
import { Marked } from 'marked';

const escapeHTML = (s: string) =>
  s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');

// Raw HTML in a note is shown as the text it is, never rendered: a note is
// written by whoever can write the board, and the page runs on loopback with
// the right to change epics.
const md = new Marked({
  gfm: true,
  breaks: true,
  renderer: { html: ({ text }) => escapeHTML(text) },
});

// A link in a note opens beside the page instead of replacing it.
DOMPurify.addHook('afterSanitizeAttributes', (node) => {
  if (node.tagName === 'A' && node.hasAttribute('href')) {
    node.setAttribute('target', '_blank');
    node.setAttribute('rel', 'noopener noreferrer');
  }
});

// renderMarkdown turns a note into HTML that is safe to insert. Escaping the
// raw HTML above is the first wall; DOMPurify is the second and also strips
// javascript: links and event attributes that markdown itself can produce.
export function renderMarkdown(text: string): string {
  return DOMPurify.sanitize(md.parse(text, { async: false }) as string);
}
