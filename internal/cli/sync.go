package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/berk/jaira/core/gate"
	"github.com/berk/jaira/core/ticket"
)

// Task mirrors the shape of Claude Code's structured task tools
// (TaskCreate/TaskUpdate/TaskGet/TaskList). Only the fields jaira uses are
// modelled, and unknown fields are ignored, because this is an external API that
// has already changed once — it replaced the older flat TodoWrite list — and will
// change again. Keeping the adapter this thin is what makes the next change a
// small patch rather than a redesign.
type Task struct {
	ID          string         `json:"id"`
	Subject     string         `json:"subject"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	Owner       string         `json:"owner"`
	BlockedBy   []string       `json:"blockedBy"`
	Metadata    map[string]any `json:"metadata"`
}

type taskPayload struct {
	Tasks []Task `json:"tasks"`
}

// taskMap remembers which ticket mirrors which task.
//
// Identity is the whole problem here. It is kept out of the repository because it
// is per-session scratch state, and keyed by the task's own ID where one is
// available — with a content hash only as a fallback for a task that predates
// having been mirrored.
type taskMap struct {
	ByTaskID map[string]string `json:"by_task_id"`
	ByHash   map[string]string `json:"by_hash"`
	// Digest of the last payload, so an unchanged sync writes nothing at all.
	LastDigest string `json:"last_digest"`
}

const metadataTicketKey = "jaira_id"

func loadTaskMap(s *ticket.Store) (*taskMap, string) {
	path := filepath.Join(s.SessionsDir(), "task-map.json")
	m := &taskMap{ByTaskID: map[string]string{}, ByHash: map[string]string{}}
	if b, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(b, m)
		if m.ByTaskID == nil {
			m.ByTaskID = map[string]string{}
		}
		if m.ByHash == nil {
			m.ByHash = map[string]string{}
		}
	}
	return m, path
}

func (m *taskMap) save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

func taskHash(t Task) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(t.Subject) + "\x00" + strings.TrimSpace(t.Description)))
	return hex.EncodeToString(sum[:8])
}

func newSyncTasksCmd() *cobra.Command {
	var file string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "sync-tasks",
		Short: "Mirror an agent's task list into the backlog",
		Long: `Reads a JSON task list on stdin and mirrors it into the backlog.

Intended to be driven from a PostToolUse hook on the task tools. The payload is
{"tasks":[{"id","subject","description","status","owner","blockedBy","metadata"}]}.

Three properties make this safe to run on every tool call:

  • It only writes ticket files. It never calls a tool or a model, so it cannot
    itself trigger another round of syncing.
  • It is idempotent. A task already mapped to a ticket updates that ticket
    rather than creating a second one, and an unchanged payload writes nothing —
    so a board→tasks→board round trip settles instead of oscillating.
  • A task disappearing from the list never deletes its ticket. Task lists are
    working memory; tickets are the durable record.

Mirrored tickets land in the backlog behind the promotion gate. Lane movement is
deliberately not mirrored: letting an external status push a ticket into the
pipeline would route around the gate that makes the pipeline worth having.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			var raw []byte
			if file != "" {
				raw, err = os.ReadFile(file)
			} else {
				raw, err = io.ReadAll(cmd.InOrStdin())
			}
			if err != nil {
				return err
			}
			raw = []byte(strings.TrimSpace(string(raw)))
			if len(raw) == 0 {
				return fail(ExitUsage, "usage", "no payload on stdin; expected {\"tasks\":[…]}")
			}

			var p taskPayload
			if err := json.Unmarshal(raw, &p); err != nil {
				// Tolerate a bare array, which is what some callers send.
				var bare []Task
				if err2 := json.Unmarshal(raw, &bare); err2 != nil {
					return fail(ExitUsage, "bad_payload", "could not parse the task payload: %v", err)
				}
				p.Tasks = bare
			}

			m, mapPath := loadTaskMap(s)
			digest := digestOf(p.Tasks)
			if digest == m.LastDigest {
				// Nothing changed. Writing anyway would produce a git diff for no
				// reason, which is how a sync loop becomes visible noise.
				if g.jsonOut {
					return emit(cmd.OutOrStdout(), map[string]any{
						"created": 0, "updated": 0, "unchanged": true,
					})
				}
				fmt.Fprintln(cmd.OutOrStdout(), "Nothing changed.")
				return nil
			}

			me := identity()
			var created, updated []string
			now := time.Now()

			for _, task := range p.Tasks {
				if strings.TrimSpace(task.Subject) == "" {
					continue
				}
				id := m.resolve(s, task)

				if id == "" {
					if dryRun {
						created = append(created, task.Subject)
						continue
					}
					tid := ticket.NewID(now)
					now = now.Add(time.Millisecond) // keep ULIDs ordered within a burst
					owner := task.Owner
					if owner == "" {
						owner = me
					}
					fields := map[string]string{
						ticket.FieldID:        tid,
						ticket.FieldTitle:     task.Subject,
						ticket.FieldStatus:    "backlog",
						ticket.FieldReady:     "false",
						ticket.FieldCreator:   me,
						ticket.FieldAssignee:  owner,
						ticket.FieldGoal:      strings.TrimSpace(task.Description),
						ticket.FieldCreatedAt: ticket.FormatTime(now),
						ticket.FieldUpdatedAt: ticket.FormatTime(now),
						"source":              "agent-task",
						"source-task-id":      task.ID,
					}
					lists := map[string][]string{
						ticket.FieldBlockedBy: nil,
						ticket.FieldCommits:   nil,
					}
					t, err := s.Create(fields, lists, "")
					if err != nil {
						return err
					}
					m.remember(task, t.ID)
					created = append(created, ticket.Handle(t.ID))
					continue
				}

				if dryRun {
					updated = append(updated, ticket.Handle(id))
					continue
				}
				changed, err := m.update(s, id, task)
				if err != nil {
					return err
				}
				if changed {
					updated = append(updated, ticket.Handle(id))
				}
			}

			// blocked-by is resolved in a second pass, once every task in the
			// payload has a ticket to point at.
			if !dryRun {
				if err := m.linkDependencies(s, p.Tasks); err != nil {
					return err
				}
				m.LastDigest = digest
				if err := m.save(mapPath); err != nil {
					return err
				}
			}

			if g.jsonOut {
				return emit(cmd.OutOrStdout(), map[string]any{
					"created": created, "updated": updated, "unchanged": false,
				})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Mirrored %d task(s): %d created, %d updated.\n",
				len(p.Tasks), len(created), len(updated))
			if len(created) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(),
					"New tickets are in the backlog and need a definition of done before they can start.\n")
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "read the payload from a file instead of stdin")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "report what would change without writing")
	return cmd
}

func digestOf(tasks []Task) string {
	parts := make([]string, 0, len(tasks))
	for _, t := range tasks {
		parts = append(parts, strings.Join([]string{
			t.ID, t.Subject, t.Description, t.Status, t.Owner,
			strings.Join(t.BlockedBy, ","),
		}, "\x1f"))
	}
	sort.Strings(parts)
	sum := sha256.Sum256([]byte(strings.Join(parts, "\x1e")))
	return hex.EncodeToString(sum[:12])
}

// sourceField marks a ticket as having been created by the sync, which is what
// makes it eligible to be updated by a later sync.
const sourceField = "source"
const sourceAgentTask = "agent-task"

// syncOwned reports whether a ticket was created by mirroring a task.
//
// Sync may only modify what it created. A payload arrives from outside and is not
// trustworthy about which tickets it owns: without this check, a task carrying
// another ticket's id in its metadata could rename or rewrite a ticket a person
// authored by hand, and a hand-edited goal could be clobbered by whatever an agent
// happened to put in a task description.
func syncOwned(t *ticket.Ticket) bool {
	v, _, err := t.Doc().Scalar(sourceField)
	return err == nil && v == sourceAgentTask
}

// resolve finds the ticket already mirroring a task.
//
// A candidate is accepted only if it is sync-owned. An id claimed in task
// metadata is treated as a hint to verify, never as authority.
func (m *taskMap) resolve(s *ticket.Store, task Task) string {
	accept := func(id string) string {
		if id == "" {
			return ""
		}
		t, err := s.Load(id)
		if err != nil || !syncOwned(t) {
			return ""
		}
		return id
	}
	if task.Metadata != nil {
		if v, ok := task.Metadata[metadataTicketKey]; ok {
			if id, ok := v.(string); ok {
				if got := accept(id); got != "" {
					return got
				}
			}
		}
	}
	if task.ID != "" {
		if got := accept(m.ByTaskID[task.ID]); got != "" {
			return got
		}
	}
	return accept(m.ByHash[taskHash(task)])
}

func (m *taskMap) remember(task Task, ticketID string) {
	if task.ID != "" {
		m.ByTaskID[task.ID] = ticketID
	}
	m.ByHash[taskHash(task)] = ticketID
}

// update refreshes the mirrored fields of an existing ticket, and reports
// whether anything actually changed.
func (m *taskMap) update(s *ticket.Store, id string, task Task) (bool, error) {
	cur, err := s.Load(id)
	if err != nil {
		return false, err
	}
	if !syncOwned(cur) {
		// Defence in depth: resolve should never have returned this ticket.
		m.remember(task, id)
		return false, nil
	}
	desc := strings.TrimSpace(task.Description)
	titleSame := cur.Title == task.Subject
	// A task description may only SEED an empty goal. Once any goal exists —
	// whether a person wrote it or an earlier sync did — an external list does not
	// get to replace it. The earlier version of this rule protected only tickets
	// that already passed the promotion gate, which left a half-written goal, the
	// most vulnerable state, wide open.
	goalSame := desc == "" || cur.Goal != ""
	if titleSame && goalSame {
		m.remember(task, id)
		return false, nil
	}
	_, err = s.Mutate(id, func(t *ticket.Ticket) error {
		if !titleSame {
			if err := t.Doc().SetScalar(ticket.FieldTitle, task.Subject); err != nil {
				return err
			}
		}
		if !goalSame {
			if err := t.Doc().SetScalar(ticket.FieldGoal, desc); err != nil {
				return err
			}
		}
		return ticket.SetReady(t.Doc(), gate.Ready(t))
	})
	if err != nil {
		return false, err
	}
	m.remember(task, id)
	return true, nil
}

// linkDependencies translates task-level blocking into ticket blocked-by, so the
// two systems agree about ordering.
func (m *taskMap) linkDependencies(s *ticket.Store, tasks []Task) error {
	for _, task := range tasks {
		if len(task.BlockedBy) == 0 {
			continue
		}
		id := m.resolve(s, task)
		if id == "" {
			continue
		}
		var deps []string
		for _, b := range task.BlockedBy {
			if dep, ok := m.ByTaskID[b]; ok {
				deps = append(deps, dep)
			}
		}
		if len(deps) == 0 {
			continue
		}
		cur, err := s.Load(id)
		if err != nil {
			continue
		}
		merged := append([]string{}, cur.BlockedBy...)
		var added bool
		for _, d := range deps {
			if !contains(merged, d) {
				merged = append(merged, d)
				added = true
			}
		}
		if !added {
			continue
		}
		if _, err := s.Mutate(id, func(t *ticket.Ticket) error {
			return t.Doc().SetList(ticket.FieldBlockedBy, merged)
		}); err != nil {
			return err
		}
	}
	return nil
}

// newTasksCmd emits what an agent should put in its own task list.
//
// The reverse direction cannot be enforced by a hook — a hook can inject context
// or block a call, but it cannot issue a tool call on the model's behalf. So this
// direction is a convention the skill teaches, supported by a command that makes
// following it trivial.
func newTasksCmd() *cobra.Command {
	var laneFilter string
	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "Emit the board as a task list for an agent to adopt",
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
			type outTask struct {
				Subject     string            `json:"subject"`
				Description string            `json:"description"`
				Status      string            `json:"status"`
				Owner       string            `json:"owner,omitempty"`
				Metadata    map[string]string `json:"metadata"`
			}
			var out []outTask
			for _, t := range all {
				if l, ok := env.Lanes.Get(t.Status); ok && l.Terminal {
					continue
				}
				if laneFilter != "" && t.Status != laneFilter {
					continue
				}
				status := "pending"
				if t.Status == "in-progress" {
					status = "in_progress"
				}
				desc := t.Goal
				if t.DoD != "" {
					desc = strings.TrimSpace(desc + "\nDone when: " + t.DoD)
				}
				out = append(out, outTask{
					Subject:     t.Title,
					Description: desc,
					Status:      status,
					Owner:       t.Assignee,
					// Carrying the ticket id back means the next sync matches on
					// identity rather than guessing from the text.
					Metadata: map[string]string{metadataTicketKey: t.ID},
				})
			}
			if g.jsonOut || true {
				// This command exists to be machine-read, so it always emits JSON.
				return emit(cmd.OutOrStdout(), map[string]any{"tasks": out, "count": len(out)})
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&laneFilter, "lane", "", "only tickets in this lane")
	return cmd
}
