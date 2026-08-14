package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/board"
	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/project"
	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/ticket"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Prepare this repository for jaira",
		Long: `Creates .jaira/ with a tickets directory, plus a .gitignore so ephemeral
session and lock state is never committed. Safe to run more than once.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := ticket.At(g.dir)
			if err != nil {
				return err
			}
			created, err := s.Init()
			if err != nil {
				return err
			}
			// Register the board here as well as on first open, so the switcher
			// knows about a project before anyone has launched the TUI in it.
			project.Remember(s.Root)

			// The default board decides which lanes this board starts with. A
			// project that has already scoped its own .jaira/lanes/ — by hand or
			// from an earlier init — is left alone: init must be safe to run more
			// than once, and re-applying the default board over a project's own
			// choices would silently discard them.
			alreadyScoped := lane.ProjectLanesActive(s.Root)
			var materialised []string
			db, _ := lane.LoadDefaultBoard()
			if !alreadyScoped {
				lanesSet, err := lane.Load(s.Root)
				if err != nil {
					return err
				}
				materialised, err = lane.Materialise(s.Root, lanesSet, db)
				if err != nil {
					return err
				}
			}

			// A new board is private: the tickets stay out of git until the user
			// decides to publish them with 'jaira share'. Tell the agent the board
			// is here too. The skill's description already says to use jaira in a
			// repository that has one, but that relies on the model noticing a
			// directory; a line in the file it is handed at the start of every
			// session does not.
			p := board.Prepare(s.Root)
			// A board just prepared by this binary must not immediately be
			// nagged about being out of date.
			_ = release.Stamp(s.StateDir())

			if g.jsonOut {
				return emit(cmd.OutOrStdout(), map[string]any{
					"root": s.Root, "tickets_dir": s.TicketsDir(), "created": created,
					"private": true, "gitignore_written": p.Ignored,
					"state_dir":   s.SessionsDir(),
					"agent_notes": p.Notes,
					"default_board": db.Path, "lanes_written": materialised,
					"lane_warnings": db.Warnings,
				})
			}
			if created {
				fmt.Fprintf(cmd.OutOrStdout(), "Initialized jaira in %s\n", filepath.Join(s.Root, ticket.DirName))
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "jaira already initialized in %s\n", filepath.Join(s.Root, ticket.DirName))
			}
			switch {
			case alreadyScoped:
				fmt.Fprintf(cmd.OutOrStdout(), "This project already scopes its own lanes; the default board was not applied.\n")
			case len(materialised) == 0:
				fmt.Fprintf(cmd.OutOrStdout(), "Using the built-in lanes.\n")
			default:
				fmt.Fprintf(cmd.OutOrStdout(), "Wrote %d lane file(s) from your default board.\n", len(materialised))
			}
			for _, w := range db.Warnings {
				fmt.Fprintf(os.Stderr, "jaira: warning: %s\n", w)
			}
			if p.IgnoreErr != nil {
				fmt.Fprintf(os.Stderr, "jaira: warning: could not write .gitignore: %v\n", p.IgnoreErr)
			} else if p.Ignored {
				fmt.Fprintf(cmd.OutOrStdout(), "This board is private: .jaira/ is gitignored, so nobody else sees it.\n")
			}
			if p.NoteErr != nil {
				fmt.Fprintf(os.Stderr, "jaira: warning: could not update an agent instruction file: %v\n", p.NoteErr)
			} else if line := announceLine(p.Notes); line != "" {
				fmt.Fprint(cmd.OutOrStdout(), line)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\nRun 'jaira share' when you want your team to have it.\n")
			return nil
		},
	}
}

func newCreateCmd() *cobra.Command {
	var (
		title, goalV, dod, contextV, assignee, laneID, tier, body, follows string
		blockedBy                                                          []string
		ready                                                              bool
	)
	cmd := &cobra.Command{
		Use:   "create <title>",
		Short: "Create a ticket",
		Long: `Creates a ticket in the backlog.

Only a title is required — capture should be cheap. A ticket cannot leave the
backlog until it has a goal, a definition of done, the context it came from, and
an assignee, so supplying those now saves a round trip later.

The context is the only record of why this ticket exists — there is no separate
description. Write it for someone who was not in the conversation, reading it
weeks from now: what is wrong today, what triggered it, and what is already known
or already ruled out. Write it as if that reader has mild ADHD — what is wrong
first, short and concrete, names and paths rather than adjectives, no preamble.
If acting on it would need a question answered first, it is not finished. It can
span several lines; they are stored readably, and someone should be able to act
after the first two.`,
		Args: minArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			if title == "" {
				title = strings.Join(args, " ")
			}
			lanes, err := lane.Load(s.Root)
			if err != nil {
				return err
			}
			target := laneID
			if target == "" {
				target = lanes.Default().ID
			}
			if _, ok := lanes.Get(target); !ok {
				return fail(ExitUsage, "no_such_lane", "no lane %q is installed; available: %s", target, strings.Join(lanes.IDs(), ", "))
			}

			// A follows: link is only worth writing if it resolves: a dead
			// reference is worse than none, because it looks like a trail exists.
			// Load accepts an unambiguous prefix, so a handle works here, and it
			// is normalised to the full id so the link matches on both sides.
			followsID := ""
			if strings.TrimSpace(follows) != "" {
				src, err := s.Load(follows)
				if err != nil {
					return err
				}
				followsID = src.ID
			}

			now := time.Now()
			me := identity()
			tplBody, tplFields, tplLists, hasTemplate := templateBody(s, title)
			// Ownership defaults to whoever created the ticket, but a board
			// template naming an assignee is an explicit choice and outranks that
			// default. An explicit --assignee outranks both.
			if assignee == "" {
				if v, ok := tplFields[ticket.FieldAssignee]; ok && strings.TrimSpace(v) != "" {
					assignee = v
				} else {
					assignee = me
				}
			}

			fields := map[string]string{
				ticket.FieldID:        ticket.NewID(now),
				ticket.FieldTitle:     title,
				ticket.FieldStatus:    target,
				ticket.FieldCreator:   me,
				ticket.FieldAssignee:  assignee,
				ticket.FieldGoal:      goalV,
				ticket.FieldContext:   contextV,
				ticket.FieldDoD:       dod,
				ticket.FieldFollows:   followsID,
				ticket.FieldModelTier: tier,
				ticket.FieldCreatedAt: ticket.FormatTime(now),
				ticket.FieldUpdatedAt: ticket.FormatTime(now),
			}
			lists := map[string][]string{
				ticket.FieldBlockedBy: blockedBy,
				ticket.FieldCommits:   nil,
			}
			// `ready` is a derived convenience recording whether the promotion
			// gate is satisfied, so the board can render without re-deriving. It
			// is written up front to keep field order canonical; the gate itself
			// remains the authority on every move.
			// Readiness is recomputed from the written file below; this initial
			// value only avoids an extra write for the common fully-specified case.
			fields[ticket.FieldReady] = fmt.Sprintf("%t",
				goalV != "" && dod != "" && contextV != "" && assignee != "")

			// A new ticket is born with the section headings a human will fill in.
			// The definition of done is a checklist rather than a sentence because
			// that is how acceptance criteria are actually written — and because a
			// ticked box is a human act, which is what lets the terminal lane
			// require evidence a model cannot manufacture.
			if body == "" {
				if hasTemplate {
					body = tplBody
					for k, v := range tplFields {
						// An explicit flag always beats a template default.
						if cur, set := fields[k]; !set || strings.TrimSpace(cur) == "" {
							fields[k] = v
						}
					}
					for k, v := range tplLists {
						if _, set := lists[k]; !set {
							lists[k] = v
						}
					}
				} else {
					db, _ := lane.LoadDefaultBoard()
					body = ticket.NewBody(title, dod, lane.ResolveOptions(lanes, db))
				}
			}
			t, err := s.Create(fields, lists, body)
			if err != nil {
				return err
			}
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), ticketJSON(t, lanes))
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created %s  %s\n", ticket.Handle(t.ID), t.Title)
			if !gate.Ready(t) {
				missing := gate.Violations(nil)
				fmt.Fprintf(cmd.OutOrStdout(), "In %s. Still needed before it can start: %s\n",
					t.Status, strings.Join(missingFields(t), ", "))
				_ = missing
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&title, "title", "", "ticket title (defaults to the positional arguments)")
	f.StringVar(&goalV, "goal", "", "what this ticket is for")
	f.StringVar(&dod, "dod", "", "definition of done: the checkable target")
	f.StringVar(&contextV, "context", "", "why this ticket exists: what is wrong now, what triggered it, what is already known")
	f.StringVar(&assignee, "assignee", "", "human who owns the outcome (defaults to you)")
	f.StringVar(&laneID, "lane", "", "lane to create in (default: backlog)")
	f.StringVar(&tier, "tier", "", "model tier alias for agentic lanes")
	f.StringVar(&body, "body", "", "markdown body")
	f.StringSliceVar(&blockedBy, "blocked-by", nil, "ticket ids that must finish first")
	f.StringVar(&follows, "follows", "", "id of the ticket this one follows on from")
	f.BoolVar(&ready, "ready", false, "unused; readiness is derived from the gate")
	_ = f.MarkHidden("ready")
	return cmd
}

// templateBody returns the board's own ticket template if it has one.
//
// A board that already has a house style for tickets — a Jira export, a team
// convention, a different language — should keep it. jaira only requires the
// frontmatter it manages and a recognisable definition-of-done heading, so the
// rest of the shape is the board's business, not the tool's.
//
// The title placeholder is the first level-one heading, which is also what jaira
// falls back to when a ticket has no title field.
func templateBody(s *ticket.Store, title string) (string, map[string]string, map[string][]string, bool) {
	for _, name := range []string{"template.md", "TEMPLATE.md"} {
		raw, err := os.ReadFile(filepath.Join(s.Root, ticket.DirName, name))
		if err != nil {
			continue
		}
		body := string(raw)
		// A template's frontmatter is the board's defaults — a Jira export's
		// type, labels and epic link, say. Keys jaira manages itself are ignored;
		// everything else is copied onto the new ticket, which is the only reason
		// a template with frontmatter is worth having.
		defaults := map[string]string{}
		listDefaults := map[string][]string{}
		if d, err := ticket.ParseDoc(raw); err == nil {
			body = d.Body()
			for _, k := range d.Keys() {
				if managedField(k) {
					continue
				}
				if v, ok, err := d.Scalar(k); err == nil && ok && strings.TrimSpace(v) != "" {
					defaults[k] = v
					continue
				}
				// Labels and similar collections are as much a board default as
				// a scalar is, and dropping them would lose exactly the fields a
				// Jira export exists to carry.
				if items, err := d.List(k); err == nil && len(items) > 0 {
					listDefaults[k] = items
				}
			}
		}
		out := make([]string, 0, 32)
		replaced := false
		for _, line := range strings.Split(body, "\n") {
			if !replaced && strings.HasPrefix(line, "# ") {
				out = append(out, "# "+title)
				replaced = true
				continue
			}
			out = append(out, line)
		}
		if !replaced {
			out = append([]string{"# " + title, ""}, out...)
		}
		return strings.TrimLeft(strings.Join(out, "\n"), "\n"), defaults, listDefaults, true
	}
	return "", nil, nil, false
}

// managedField reports whether jaira writes this key itself, in which case a
// template must not be able to preset it — an id or a status from a template
// would be wrong on every ticket made from it.
func managedField(k string) bool {
	switch k {
	case ticket.FieldID, ticket.FieldStatus, ticket.FieldReady,
		ticket.FieldCreatedAt, ticket.FieldUpdatedAt, ticket.FieldTitle,
		ticket.FieldCreator, ticket.FieldClaimedBy, ticket.FieldClaimedAt:
		return true
	}
	return false
}

func missingFields(t *ticket.Ticket) []string {
	var out []string
	for _, f := range gate.PromotionFields {
		switch f {
		case ticket.FieldGoal:
			if strings.TrimSpace(t.Goal) == "" {
				out = append(out, f)
			}
		case ticket.FieldDoD:
			if !t.HasDoD() {
				out = append(out, f)
			}
		case ticket.FieldContext:
			if strings.TrimSpace(t.Context) == "" {
				out = append(out, f)
			}
		case ticket.FieldAssignee:
			if strings.TrimSpace(t.Assignee) == "" {
				out = append(out, f)
			}
		}
	}
	return out
}

func newListCmd() *cobra.Command {
	var laneFilter, assigneeFilter, query string
	var actionableOnly bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tickets",
		Args:  noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			env, all, err := loadEnv(s)
			if err != nil {
				return err
			}
			var out []*ticket.Ticket
			for _, t := range all {
				if laneFilter != "" && t.Status != laneFilter {
					continue
				}
				if assigneeFilter != "" && !strings.EqualFold(t.Assignee, assigneeFilter) {
					continue
				}
				if query != "" && !matches(t, query) {
					continue
				}
				if actionableOnly && !gate.Actionable(env, t) {
					continue
				}
				out = append(out, t)
			}
			if g.jsonOut {
				arr := make([]map[string]any, 0, len(out))
				for _, t := range out {
					arr = append(arr, ticketJSON(t, env.Lanes))
				}
				return emit(cmd.OutOrStdout(), map[string]any{"tickets": arr, "count": len(arr)})
			}
			if len(out) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No tickets match.")
				return nil
			}
			printTable(cmd.OutOrStdout(), out, env)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&laneFilter, "lane", "", "only tickets in this lane")
	f.StringVar(&assigneeFilter, "assignee", "", "only tickets assigned to this person")
	f.StringVarP(&query, "query", "q", "", "substring match over title, goal and id")
	f.BoolVar(&actionableOnly, "actionable", false, "only tickets that could be started right now")
	return cmd
}

func matches(t *ticket.Ticket, q string) bool {
	q = strings.ToLower(q)
	for _, field := range []string{t.ID, t.Title, t.Goal, t.Context, t.DoD, t.Assignee, t.Status} {
		if strings.Contains(strings.ToLower(field), q) {
			return true
		}
	}
	return false
}

func printTable(w io.Writer, ts []*ticket.Ticket, env gate.Env) {
	// Group by lane in board order so the text output reads like the board.
	byLane := map[string][]*ticket.Ticket{}
	for _, t := range ts {
		byLane[t.Status] = append(byLane[t.Status], t)
	}
	statuses := make([]string, 0, len(byLane))
	for k := range byLane {
		statuses = append(statuses, k)
	}
	for _, l := range env.Lanes.Columns(statuses) {
		group := byLane[l.ID]
		if len(group) == 0 {
			continue
		}
		label := l.Name
		if l.Unknown {
			label += "  (unrecognized lane)"
		}
		fmt.Fprintf(w, "\n%s\n", label)
		sort.Slice(group, func(i, j int) bool { return group[i].ID < group[j].ID })
		for _, t := range group {
			flags := ""
			if !gate.Ready(t) {
				flags += " [needs spec]"
			}
			if len(t.BlockedBy) > 0 && !gate.Actionable(env, t) {
				flags += " [blocked]"
			}
			fmt.Fprintf(w, "  %-10s %-52s %s%s\n", ticket.Handle(t.ID), truncate(t.Title, 52), t.Assignee, flags)
		}
	}
	fmt.Fprintln(w)
}

func newShowCmd() *cobra.Command {
	var forLane string
	cmd := &cobra.Command{
		Use:   "show <id>",
		Short: "Show one ticket in full",
		Args:  exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			t, err := s.Load(args[0])
			if err != nil {
				return err
			}
			env, _, err := loadEnv(s)
			if err != nil {
				return err
			}
			if forLane != "" {
				return showForLane(cmd, s, env, t, forLane)
			}
			if g.jsonOut {
				j := ticketJSON(t, env.Lanes)
				j["body"] = t.Body
				j["path"] = t.Path
				return emit(cmd.OutOrStdout(), j)
			}
			printDetail(cmd.OutOrStdout(), t, env)
			return nil
		},
	}
	cmd.Flags().StringVar(&forLane, "for-lane", "", "assemble the bounded input a lane's agent should receive")
	return cmd
}

func printDetail(w io.Writer, t *ticket.Ticket, env gate.Env) {
	fmt.Fprintf(w, "%s  %s\n", ticket.Handle(t.ID), t.Title)
	fmt.Fprintf(w, "%s\n", strings.Repeat("─", 64))
	row := func(k, v string) {
		if strings.TrimSpace(v) == "" {
			return
		}
		fmt.Fprintf(w, "%-10s %s\n", k, v)
	}
	row("lane", t.Status)
	row("assignee", t.Assignee)
	row("creator", t.Creator)
	row("executed", t.ExecutedBy)
	row("tier", t.ModelTier)
	row("goal", t.Goal)
	row("context", t.Context)
	if len(t.DoDItems) > 0 {
		fmt.Fprintf(w, "%-10s\n", "done when")
		for _, it := range t.DoDItems {
			fmt.Fprintf(w, "           [%s] %s\n", it.State.Marker(), it.Text)
		}
	} else {
		row("dod", t.DoD)
	}
	if len(t.BlockedBy) > 0 {
		shorts := make([]string, 0, len(t.BlockedBy))
		for _, b := range t.BlockedBy {
			shorts = append(shorts, ticket.Handle(b))
		}
		row("blocked-by", strings.Join(shorts, ", "))
	}
	if t.Follows != "" {
		row("follows", ticket.Handle(t.Follows))
	}
	if t.Question != "" {
		row("question", t.Question)
	}
	if t.Outcome.What != "" || t.Outcome.Why != "" || t.Outcome.Resolves != "" {
		fmt.Fprintf(w, "\nOutcome\n")
		row("  what", t.Outcome.What)
		row("  why", t.Outcome.Why)
		row("  resolves", t.Outcome.Resolves)
	}
	if len(t.Commits) > 0 {
		row("commits", strings.Join(t.Commits, " "))
	}
	if miss := missingFields(t); len(miss) > 0 {
		fmt.Fprintf(w, "\nBefore this can leave the backlog: %s\n", strings.Join(miss, ", "))
	}
	if t.Body != "" {
		fmt.Fprintf(w, "\n%s\n", strings.TrimSpace(t.Body))
	}
}

func newSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <id> <field=value>...",
		Short: "Set ticket fields",
		Long: `Sets one or more frontmatter fields.

Writing a field rewrites only that field's bytes: unrelated fields, comments and
blank lines are left exactly as they were, so a change shows up in git as a
one-line diff.

List fields take a comma-separated value, for example blocked-by=01AAA,01BBB.`,
		Args: minArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			id := args[0]
			assignments := args[1:]
			listFields := map[string]bool{ticket.FieldBlockedBy: true, ticket.FieldCommits: true}

			// A ticket in a lane this installation does not have is read-only, and
			// that has to hold for every mutation path. Enforcing it only in `move`
			// would leave the rule true in name and false in practice.
			cur, err := s.Load(id)
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

			t, err := s.Mutate(id, func(t *ticket.Ticket) error {
				for _, a := range assignments {
					k, v, ok := strings.Cut(a, "=")
					if !ok {
						return fail(ExitUsage, "usage", "expected field=value, got %q", a)
					}
					k = strings.TrimSpace(k)
					if k == ticket.FieldID {
						return fail(ExitUsage, "immutable", "id cannot be changed")
					}
					if k == ticket.FieldStatus {
						return fail(ExitUsage, "use_move", "use 'jaira move' to change the lane, so gates are applied")
					}
					if listFields[k] {
						var items []string
						for _, p := range strings.Split(v, ",") {
							if p = strings.TrimSpace(p); p != "" {
								items = append(items, p)
							}
						}
						if err := t.Doc().SetList(k, items); err != nil {
							return err
						}
						continue
					}
					if err := t.Doc().SetScalar(k, v); err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
			// Keep the derived readiness flag honest after any edit.
			t, err = s.Mutate(t.ID, func(t *ticket.Ticket) error {
				return ticket.SetReady(t.Doc(), gate.Ready(t))
			})
			if err != nil {
				return err
			}
			env, _, err := loadEnv(s)
			if err != nil {
				return err
			}
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), ticketJSON(t, env.Lanes))
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Updated %s\n", ticket.Handle(t.ID))
			if miss := missingFields(t); len(miss) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Still needed before it can start: %s\n", strings.Join(miss, ", "))
			}
			return nil
		},
	}
	return cmd
}

func newLanesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lanes",
		Short: "List the installed lanes",
		Args:  noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			// jaira lanes must work outside a board too — that is what makes it
			// useful for writing a project's first lane file — so the root is
			// best-effort rather than required.
			lanes, err := lane.Load(bestEffortRoot())
			if err != nil {
				return err
			}
			// The default board is a fourth source of lane-shaped trouble — a
			// typo'd id or option — that nothing checked before this task. One
			// warning channel: its findings are appended to the loader's own
			// rather than reported through a second path.
			db, _ := lane.LoadDefaultBoard()
			warnings := append(append([]string{}, lanes.Warnings...), lane.Validate(db, lanes)...)
			if g.jsonOut {
				arr := make([]map[string]any, 0, len(lanes.Lanes))
				for _, l := range lanes.Lanes {
					arr = append(arr, map[string]any{
						"id": l.ID, "name": l.Name, "precedence": l.Precedence,
						"agentic": l.Agentic, "terminal": l.Terminal,
						"model_tier": l.ModelTier, "builtin": l.Builtin,
						"input_requires": l.InputRequires, "output_produces": l.OutputProduces,
						"source": l.Source, "prompt": l.Prompt, "creator": l.Creator,
					})
				}
				return emit(cmd.OutOrStdout(), map[string]any{"lanes": arr, "warnings": warnings})
			}
			w := cmd.OutOrStdout()
			// CREATOR is not a column here: the table is already six columns wide,
			// and 'lanes show <id>' is where it lives instead.
			// RANK, not PREC: column order follows after:, never this number — it
			// is the rank a merge uses to decide which lane wins, nothing more.
			fmt.Fprintf(w, "%-14s %-16s %-6s %-8s %-8s %s\n", "ID", "NAME", "RANK", "AGENTIC", "TIER", "SOURCE")
			for _, l := range lanes.Lanes {
				src := "built-in"
				if !l.Builtin {
					src = l.Source
				}
				fmt.Fprintf(w, "%-14s %-16s %-6d %-8t %-8s %s\n", l.ID, l.Name, l.Precedence, l.Agentic, dash(l.ModelTier), src)
			}
			for _, warn := range warnings {
				fmt.Fprintf(os.Stderr, "jaira: warning: %s\n", warn)
			}
			return nil
		},
		Long: `Lists the installed lanes.

The loop for building a lane or a default board by hand, no TUI required:
  1. 'jaira lanes path' to find where a file belongs.
  2. 'jaira lanes template' (or 'template --board') for its shape.
  3. Write the file with the tools you already have.
  4. 'jaira lanes' to check the result — a bad id or a bad option is a
     warning here, not a silent failure later.
  5. 'jaira lanes use <id>' to put it to work in this project.`,
	}
	cmd.AddCommand(
		newLanesShowCmd(), newLanesPathCmd(), newLanesTemplateCmd(), newLanesSharedCmd(),
		newLanesUseCmd(), newLanesPublishCmd(), newLanesAdoptCmd(), newLanesDefaultCmd(),
		newLanesAddCmd(), newLanesRemoveCmd(), newLanesMoveCmd(),
	)
	return cmd
}

// newLanesSharedCmd lists lanes teammates have published to this project's
// .jaira/shared/ tree. An agent has no way to press the TUI's adopt key, and
// "read what teammates published" is a read operation the CLI already
// promises for everything else.
//
// Shared lanes are never loaded onto the board by this command or by
// lane.Load — listing them is not adopting them.
func newLanesSharedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "shared",
		Short: "List lanes teammates have published to this project",
		Args:  noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			root := bestEffortRoot()
			shared, warnings, err := lane.Shared(root)
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				arr := make([]map[string]any, 0, len(shared))
				for _, sl := range shared {
					arr = append(arr, map[string]any{
						"folder": sl.Folder, "id": sl.Lane.ID, "name": sl.Lane.Name,
						"creator": sl.Lane.Creator, "path": sl.Path,
					})
				}
				return emit(w, map[string]any{"shared": arr, "warnings": warnings})
			}
			if root == "" {
				fmt.Fprintln(w, "not in a project directory; nothing to list")
				return nil
			}
			if len(shared) == 0 {
				fmt.Fprintln(w, "no lanes have been published to this project")
			} else {
				fmt.Fprintf(w, "%-14s %-14s %-16s %s\n", "ID", "FOLDER", "CREATOR", "PATH")
				for _, sl := range shared {
					fmt.Fprintf(w, "%-14s %-14s %-16s %s\n", sl.Lane.ID, sl.Folder, dash(sl.Lane.Creator), sl.Path)
				}
			}
			for _, warn := range warnings {
				fmt.Fprintf(os.Stderr, "jaira: warning: %s\n", warn)
			}
			return nil
		},
	}
}

// newLanesShowCmd prints one lane in full, prompt included — this is the
// surface an agent reads before deciding a lane needs changing, with no
// ticket in hand and nothing to run.
func newLanesShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Show one lane's full contract and prompt",
		Args:  exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			lanes, err := lane.Load(bestEffortRoot())
			if err != nil {
				return err
			}
			l, ok := lanes.Get(args[0])
			if !ok {
				return fail(ExitUsage, "no_such_lane", "no lane %q is installed; available: %s", args[0], strings.Join(lanes.IDs(), ", "))
			}
			src := "built-in"
			if !l.Builtin {
				src = l.Source
			}
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), map[string]any{
					"id": l.ID, "name": l.Name, "precedence": l.Precedence,
					"agentic": l.Agentic, "terminal": l.Terminal,
					"model_tier": l.ModelTier, "builtin": l.Builtin,
					"input_requires": l.InputRequires, "output_produces": l.OutputProduces,
					"source": src, "prompt": l.Prompt, "creator": l.Creator,
					"after": l.After, "description": l.Description, "overrides": l.Overrides,
				})
			}
			w := cmd.OutOrStdout()
			fmt.Fprintf(w, "ID:          %s\n", l.ID)
			fmt.Fprintf(w, "Name:        %s\n", l.Name)
			fmt.Fprintf(w, "Anchor:      %s\n", dash(l.After))
			// Labelled as what it is, not as a position: column order follows
			// Anchor, above, and never this number.
			fmt.Fprintf(w, "Merge rank:  %d\n", l.Precedence)
			fmt.Fprintf(w, "Tier:        %s\n", dash(l.ModelTier))
			fmt.Fprintf(w, "Input:       %s\n", dash(strings.Join(l.InputRequires, ", ")))
			fmt.Fprintf(w, "Output:      %s\n", dash(strings.Join(l.OutputProduces, ", ")))
			fmt.Fprintf(w, "Creator:     %s\n", dash(l.Creator))
			fmt.Fprintf(w, "Source:      %s\n", src)
			if l.Overrides != "" {
				fmt.Fprintf(w, "Overrides:   %s\n", l.Overrides)
			}
			if l.Description != "" {
				fmt.Fprintf(w, "\n%s\n", l.Description)
			}
			if l.Prompt != "" {
				fmt.Fprintf(w, "\n%s\n", l.Prompt)
			}
			return nil
		},
	}
}

// newLanesPathCmd names where a lane file belongs, so writing one is a
// command an agent can run rather than something it has to guess at. It
// names all three locations a lane-shaped file can go: the catalogue, the
// project directory, and the default board — a path command that named only
// two of the three would be a trap for whichever it left out.
func newLanesPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the catalogue, project and default board locations",
		Args:  noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			catalogue := lane.UserLanesDir()
			root := bestEffortRoot()
			w := cmd.OutOrStdout()
			boardPath := lane.DefaultBoardPath()
			_, statErr := os.Stat(boardPath)
			boardExists := statErr == nil

			if root == "" {
				if g.jsonOut {
					return emit(w, map[string]any{
						"catalogue": catalogue, "project": "", "in_project": false, "active": "catalogue",
						"default_board": boardPath, "default_board_exists": boardExists,
					})
				}
				fmt.Fprintf(w, "Catalogue:     %s (active)\n", catalogue)
				fmt.Fprintf(w, "Project:       not in a project directory\n")
				fmt.Fprintf(w, "Default board: %s%s\n", boardPath, existsMark(boardExists))
				return nil
			}

			projDir := lane.ProjectLanesDir(root)
			active := "catalogue"
			if lane.ProjectLanesActive(root) {
				active = "project"
			}
			if g.jsonOut {
				return emit(w, map[string]any{
					"catalogue": catalogue, "project": projDir, "in_project": true, "active": active,
					"default_board": boardPath, "default_board_exists": boardExists,
				})
			}
			mark := func(which string) string {
				if which == active {
					return " (active)"
				}
				return ""
			}
			fmt.Fprintf(w, "Catalogue:     %s%s\n", catalogue, mark("catalogue"))
			fmt.Fprintf(w, "Project:       %s%s\n", projDir, mark("project"))
			fmt.Fprintf(w, "Default board: %s%s\n", boardPath, existsMark(boardExists))
			return nil
		},
	}
}

// existsMark labels a path that has not been written yet, so 'lanes path'
// doubles as "does my default board exist at all" without a second command.
func existsMark(exists bool) string {
	if exists {
		return ""
	}
	return " (not written yet)"
}

// defaultBoardTemplate is the skeleton 'jaira lanes template --board' prints.
// It states the two rules that are not visible from the fields alone: an
// absent or empty lanes: means the built-ins, and options: names entries a
// ticket's Options checklist already has — from the installed lanes'
// requires-option fields, not a fixed list this file invents.
const defaultBoardTemplate = `---
lanes: [backlog, brainstorm, todo, pre-process, in-progress, human, review, signoff, done, blocked]
                                  # which lanes a freshly initialised board gets.
                                  # absent, or an empty list, means the built-ins — not "no lanes".
options: [brainstorm]            # ticket Options to pre-tick on a new ticket. Names come from the
                                  # installed lanes' requires-option fields, never a fixed list.
---

# Default board

Decides which lanes a freshly initialised board gets, and which ticket
Options start ticked. Edit this file directly, or from jaira's home screen
with 'd'.
`

// newLanesTemplateCmd prints a lane skeleton to stdout and nothing else — no
// file is created, no directory is scaffolded. 'jaira lanes template >
// $(jaira lanes path)/my-lane.md' is the whole workflow; --board does the
// same for the default board file.
func newLanesTemplateCmd() *cobra.Command {
	var board bool
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Print a commented lane skeleton to stdout",
		Args:  noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if board {
				fmt.Fprint(cmd.OutOrStdout(), defaultBoardTemplate)
				return nil
			}
			fmt.Fprint(cmd.OutOrStdout(), lane.Template)
			return nil
		},
	}
	cmd.Flags().BoolVar(&board, "board", false, "print a default board skeleton instead of a lane")
	return cmd
}

func dash(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

// checklistJSON renders one checklist for machine consumption.
//
// Each item carries the number 'jaira dod' takes, so an agent acting on a gate
// refusal has the argument in hand rather than having to infer it from array
// position. The state is spelled out as a word, and "done" is kept alongside it
// so a consumer can branch on completion without knowing the state vocabulary.
func checklistJSON(items []ticket.DoDItem) []map[string]any {
	out := make([]map[string]any, 0, len(items))
	for i, it := range items {
		out = append(out, map[string]any{
			"n":     i + 1,
			"text":  it.Text,
			"state": it.State.String(),
			"done":  it.Checked(),
			"proof": it.Proof,
		})
	}
	return out
}

func ticketJSON(t *ticket.Ticket, lanes *lane.Set) map[string]any {
	_, known := lanes.Get(t.Status)
	dodComplete, _ := t.DoDComplete()
	planComplete, _ := t.PlanComplete()
	return map[string]any{
		"plan_items":         checklistJSON(t.PlanItems),
		"dod_items":          checklistJSON(t.DoDItems),
		"plan_complete":      planComplete,
		"dod_complete":       dodComplete,
		"id":                 t.ID,
		"handle":             ticket.Handle(t.ID),
		"title":              t.Title,
		"status":             t.Status,
		"status_known":       known,
		"ready":              gate.Ready(t),
		"creator":            t.Creator,
		"assignee":           t.Assignee,
		"executed_by":        t.ExecutedBy,
		"goal":               t.Goal,
		"context":            t.Context,
		"definition_of_done": t.DoD,
		"blocked_by":         t.BlockedBy,
		"follows":            t.Follows,
		"commits":            t.Commits,
		"model_tier":         t.ModelTier,
		"question":           t.Question,
		"outcome": map[string]string{
			"what": t.Outcome.What, "why": t.Outcome.Why, "resolves": t.Outcome.Resolves,
		},
		"created_at": nonZero(t.CreatedAt),
		"updated_at": nonZero(t.UpdatedAt),
	}
}

func nonZero(t time.Time) any {
	if t.IsZero() {
		return nil
	}
	return ticket.FormatTime(t)
}

func truncate(s string, n int) string {
	if len([]rune(s)) <= n {
		return s
	}
	r := []rune(s)
	return string(r[:n-1]) + "…"
}
