package lane

import "github.com/BeMuCa/jaira/core/board"

// Facts reads a board's pipeline into the plain form core/board renders into
// the agent note. A nil set gives nil, so a caller without a board writes the
// note it wrote before lanes were in it.
//
// The translation lives here rather than in core/board because that package is
// stdlib-only on purpose, and here rather than in the CLI and the TUI because
// both write the same note: two copies had already drifted once a lane fact
// was added to one of them and not the other.
func (s *Set) Facts() []board.LaneFact {
	if s == nil {
		return nil
	}
	out := make([]board.LaneFact, 0, len(s.Lanes))
	for _, l := range s.Lanes {
		out = append(out, board.LaneFact{
			ID: l.ID, Name: l.Name, Description: l.Description,
			ModelTier: l.ModelTier, Agentic: l.Agentic, Terminal: l.Terminal,
			HumanExit: l.RequiresHumanExit, Question: l.RequiresQuestion,
			Parking: l.RequiresBlockedReason, RejectsTo: l.RejectsTo,
			Produces: l.OutputProduces, Fresh: l.Context == ContextFresh,
		})
	}
	return out
}
