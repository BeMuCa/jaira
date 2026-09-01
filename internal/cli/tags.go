package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/tag"
	"github.com/BeMuCa/jaira/core/ticket"
)

// normalizeTags puts written names into the one stored form, dropping repeats,
// and returns a line per name it had to change — because a tag filed under a
// different name than the one typed is exactly the kind of silent difference
// that ends up as two tags for one topic.
func normalizeTags(raw []string) (names []string, notes []string, err error) {
	seen := map[string]bool{}
	for _, r := range raw {
		name, changed, err := tag.Normalize(r)
		if err != nil {
			return nil, nil, fail(ExitUsage, "bad_tag", "%v", err)
		}
		if changed {
			notes = append(notes, fmt.Sprintf("Filed %q as %q.", r, name))
		}
		if seen[name] {
			continue
		}
		seen[name] = true
		names = append(names, name)
	}
	return names, notes, nil
}

// tagsLockName keys the store lock this registry is written under. One name for
// the whole file, because the write is a read-modify-write of every line in it.
const tagsLockName = "tags"

// registerTags gives every name a colour line in .jaira/tags, and reports which
// names the board had never seen. colour is the one explicitly asked for, or -1
// to let the registry pick a free palette colour.
//
// An explicit colour is applied whether or not the name is new: "jaira tag <id>
// ui --color 40" reads as "ui is colour 40", and refusing to recolour an
// existing tag would leave hand-editing the file as the only way to do it.
//
// The whole Load→Set→Save runs under the store's lock. Two sessions tagging at
// the same moment both read the file, both add their own line to what they read,
// and the second save then writes a file assembled before the first one's line
// existed — one tag lost, silently. An atomic write does not fix that: it makes
// each write whole, not the pair of them ordered. The lock is what orders them.
func registerTags(s *ticket.Store, names []string, colour int) (added, reused []string, err error) {
	if len(names) == 0 {
		return nil, nil, nil
	}
	unlock, err := s.Lock(tagsLockName)
	if err != nil {
		return nil, nil, err
	}
	defer unlock()

	root := s.Root
	reg, err := tag.Load(root)
	if err != nil {
		return nil, nil, err
	}
	dirty := false
	for _, name := range names {
		_, known := reg.Colour(name)
		switch {
		case !known:
			c := colour
			if c < 0 {
				c = reg.Assign(name)
			}
			reg.Set(name, c)
			added = append(added, name)
			dirty = true
		case colour >= 0:
			reg.Set(name, colour)
			reused = append(reused, name)
			dirty = true
		default:
			reused = append(reused, name)
		}
	}
	if !dirty {
		return added, reused, nil
	}
	if err := reg.Save(root); err != nil {
		return nil, nil, err
	}
	return added, reused, nil
}

// strOrEmpty renders a nil slice as an empty JSON array rather than null, so a
// caller can iterate the field without a nil check.
func strOrEmpty(v []string) []string {
	if v == nil {
		return []string{}
	}
	return v
}

// tagCount is one row of 'jaira tags': a name, the colour the board gives it,
// and how many tickets still on the board wear it.
type tagCount struct {
	name   string
	colour int
	known  bool
	open   int
}

// countTags pairs every tag the board knows with how many open tickets carry it.
//
// Two sources, deliberately: the registry, so a name with no tickets left is
// still offered for reuse rather than disappearing the moment its last ticket
// closes; and the tickets themselves, so a tag set by hand or by an older
// version shows up even though nothing ever gave it a colour.
//
// A ticket in the terminal lane is not counted. "open" has to mean outstanding
// work or the number is not worth printing — and the logbook and the archive are
// not on the board at all, so they never reach this.
func countTags(all []*ticket.Ticket, lanes *lane.Set, reg *tag.Registry) []tagCount {
	open := map[string]int{}
	seen := map[string]bool{}
	for _, t := range all {
		if l, ok := lanes.Get(t.Status); ok && l.Terminal {
			continue
		}
		// One ticket counts once per tag, however many ways the field spells it.
		// "ui" and "UI" on one ticket are one subject — and that exact pair is
		// what merge=union produces, since it unions the raw strings without
		// knowing they normalise to the same name.
		counted := map[string]bool{}
		for _, raw := range t.Tags {
			name, _, err := tag.Normalize(raw)
			if err != nil || counted[name] {
				continue
			}
			counted[name] = true
			open[name]++
			seen[name] = true
		}
	}
	for _, name := range reg.Names() {
		seen[name] = true
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	out := make([]tagCount, 0, len(names))
	for _, name := range names {
		c, known := reg.Colour(name)
		out = append(out, tagCount{name: name, colour: c, known: known, open: open[name]})
	}
	return out
}

// colourable reports whether w is a terminal, so the swatch is drawn with ANSI
// only where ANSI means something. Piped or captured output stays plain text:
// escape codes in a file an agent parses are noise, and stdout is that agent's
// channel.
func colourable(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	st, err := f.Stat()
	return err == nil && st.Mode()&os.ModeCharDevice != 0
}

// swatch is the colour shown as itself. Two blocks rather than one, because a
// single cell of colour is hard to name at a glance.
func swatch(colour int, known, colours bool) string {
	if !known {
		return "  "
	}
	if !colours {
		return "██"
	}
	return fmt.Sprintf("\x1b[38;5;%dm██\x1b[0m", colour)
}

func newTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "List the tags this board already uses",
		Long: `Lists every tag this board knows, with its colour and how many open tickets
carry it.

Read this before tagging anything. A board with "ui", "frontend" and "gui" on it
has three names for one subject and no way to filter by it, which is the failure
this listing exists to prevent: reuse a name that is already here rather than
inventing a synonym.

Colours live in ` + "`.jaira/tags`" + `, one "name: <ansi256>" line per tag, shared by the
whole board and editable by hand. A tag with no line there is shown without a
colour; nothing is written by this command, since a listing that edited the board
would race every other session reading it. 'jaira tag' is what gives a name a
colour.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			env, all, err := loadEnv(s)
			if err != nil {
				return err
			}
			reg, err := tag.Load(s.Root)
			if err != nil {
				return err
			}
			rows := countTags(all, env.Lanes, reg)

			w := cmd.OutOrStdout()
			if g.jsonOut {
				arr := make([]map[string]any, 0, len(rows))
				for _, r := range rows {
					var colour any
					if r.known {
						colour = r.colour
					}
					// "color", matching --color: one spelling on the machine
					// surface, whatever the prose around it says.
					arr = append(arr, map[string]any{
						"name": r.name, "color": colour, "open": r.open,
					})
				}
				return emit(w, map[string]any{
					"tags": arr, "count": len(arr), "file": tag.Path(s.Root),
				})
			}
			if len(rows) == 0 {
				fmt.Fprintf(w, "No tags on this board yet. 'jaira tag <id> <name>' starts the vocabulary.\n")
				return nil
			}
			colours := colourable(w)
			fmt.Fprintf(w, "\nReuse one of these rather than inventing a synonym for the same subject.\n\n")
			for _, r := range rows {
				// The number beside the swatch, not instead of it: piped output
				// has no colour to show, and a reader hand-editing .jaira/tags
				// needs the value rather than the block.
				value := "  -"
				if r.known {
					value = fmt.Sprintf("%3d", r.colour)
				}
				line := fmt.Sprintf("  %s %-20s %s  %3d open",
					swatch(r.colour, r.known, colours), r.name, value, r.open)
				if !r.known {
					line += "   no colour yet"
				}
				fmt.Fprintln(w, line)
			}
			fmt.Fprintf(w, "\nColours: %s\n", tag.Path(s.Root))
			return nil
		},
	}
	return cmd
}

func newTagCmd() *cobra.Command {
	var colour int
	cmd := &cobra.Command{
		Use:   "tag <id> <name>...",
		Short: "Add tags to a ticket",
		Long: `Adds one or more topic tags to a ticket.

Run 'jaira tags' first. A name this board already uses is reused and said to be
reused; a name it has never seen is new, and is given a free colour from jaira's
palette in ` + "`.jaira/tags`" + ` so the whole board renders it the same way.

Names are stored as lowercase kebab — "My UI" is filed as "my-ui", and you are
told so. A name with anything else in it is refused rather than trimmed down,
because a quietly shortened name is a second name for one subject.

--color <0-255> picks the colour instead of leaving it to the palette, and
recolours a tag that already has one. It takes exactly one name, since one
colour cannot be meant for several.`,
		Args: minArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			names, notes, err := normalizeTags(args[1:])
			if err != nil {
				return err
			}
			// Whether --color was given, not whether its value looks unset: a
			// negative colour is a mistake to report, not a request for the
			// palette, and only the flag itself can tell the two apart.
			chosen := -1
			if cmd.Flags().Changed("color") {
				if !tag.ValidColour(colour) {
					return fail(ExitUsage, "bad_colour", "--color takes an ANSI-256 colour, 0 to 255, not %d", colour)
				}
				if len(names) != 1 {
					return fail(ExitUsage, "usage", "--color takes exactly one tag name, got %d", len(names))
				}
				chosen = colour
			}

			// A ticket in a lane this installation does not have is read-only on
			// every mutation path, tagging included — the same rule 'jaira set'
			// enforces, for the same reason.
			cur, err := s.Load(args[0])
			if err != nil {
				return err
			}
			lanes, err := lane.Load(s.Root)
			if err != nil {
				return err
			}
			if _, known := lanes.Get(cur.Status); !known && cur.Status != "" {
				return &codedError{
					code:   ExitValidation,
					reason: "unknown_lane",
					message: fmt.Sprintf(
						"%s sits in unrecognized lane %q and is read-only; install that lane file first",
						ticket.Handle(cur.ID), cur.Status),
				}
			}

			t, err := s.Mutate(args[0], func(t *ticket.Ticket) error {
				// Existing entries are kept exactly as the file has them. This
				// command was asked to add a tag, not to rewrite the ones already
				// there, and a name that only differs by case is caught by the
				// duplicate check below without touching the line it came from.
				have := map[string]bool{}
				merged := append([]string{}, t.Tags...)
				for _, raw := range t.Tags {
					if n, _, err := tag.Normalize(raw); err == nil {
						have[n] = true
					}
				}
				for _, name := range names {
					if have[name] {
						continue
					}
					have[name] = true
					merged = append(merged, name)
				}
				return t.Doc().SetList(ticket.FieldTags, merged)
			})
			if err != nil {
				return err
			}
			added, reused, err := registerTags(s, names, chosen)
			if err != nil {
				return err
			}

			w := cmd.OutOrStdout()
			if g.jsonOut {
				// The new-versus-reused signal is the point of this command, and
				// --json is the surface the skill tells agents to use — so it
				// travels beside the ticket rather than only in the prose an
				// agent never sees. Always arrays, never null: an agent
				// branching on length must not have to handle both.
				j := ticketJSON(t, lanes)
				j["tags_new"] = strOrEmpty(added)
				j["tags_reused"] = strOrEmpty(reused)
				return emit(w, j)
			}
			for _, n := range notes {
				fmt.Fprintln(w, n)
			}
			fmt.Fprintf(w, "%s tags: %s\n", ticket.Handle(t.ID), strings.Join(t.Tags, " "))
			if len(reused) > 0 {
				fmt.Fprintf(w, "Reused: %s — already on this board.\n", strings.Join(reused, ", "))
			}
			if len(added) > 0 {
				fmt.Fprintf(w, "New: %s. Run 'jaira tags' before naming a tag, and reuse one for the same subject.\n",
					strings.Join(added, ", "))
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&colour, "color", -1, "ANSI-256 colour (0-255) for the tag, instead of a free palette colour")
	return cmd
}
