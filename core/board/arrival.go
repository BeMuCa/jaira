package board

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BeMuCa/jaira/core/ticket"
)

// The three values of the board's pr file: how work arrives on the default
// branch. Here and not in core/settings, because the note this package writes
// into CLAUDE.md and AGENTS.md says it too, and settings imports this package,
// not the other way round. settings.PRFor and settings.PRIn read through these.
const (
	// PRPerson: an agent pushes the branch and stops; a person opens the
	// pull request. The default.
	PRPerson = "person"
	// PRAgent: an agent opens the pull request itself once every ticket on
	// the branch is accepted (JG3REW).
	PRAgent = "agent"
	// PRNone: the project takes no pull requests. A ticket still rides on a
	// branch of its own, and once every ticket on it is accepted the branch
	// lands on the default branch with a fast-forward push.
	PRNone = "none"
)

// PRFileName is the board file under .jaira/ that holds one of the three.
const PRFileName = "pr"

// PRWordIn is the first line of a .jaira/pr file's content — or any other
// one-line board file's: settings reads pr-when, release and the rest through
// it — that is neither blank nor a comment, trimmed; "" when there is none. A UTF-8 byte order mark
// in front of it — what an editor on Windows may save — is not part of it.
func PRWordIn(content string) string {
	content = strings.TrimPrefix(content, "\uFEFF")
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		return line
	}
	return ""
}

// PRModeIn reads the content of a .jaira/pr file: its PRWordIn, when that is
// exactly "agent" or "none" — anything else, missing and misspelt included, is
// PRPerson. A typo must never be what lets an agent open a pull request or
// push to the default branch.
func PRModeIn(content string) string {
	if w := PRWordIn(content); w == PRAgent || w == PRNone {
		return w
	}
	return PRPerson
}

// PRModeFor reads the board at root's .jaira/pr off disk the way PRModeIn
// reads its content; PRPerson when root is "" or the file is unreadable.
func PRModeFor(root string) string {
	if root == "" {
		return PRPerson
	}
	b, err := os.ReadFile(filepath.Join(root, ticket.DirName, PRFileName))
	if err != nil {
		return PRPerson
	}
	return PRModeIn(string(b))
}

// arrivalSection is what the note says about how work reaches the default
// branch on this board. PRPerson writes nothing: the block stays byte for byte
// what it was before boards could say this, so a board nobody configured sees
// no change and no drift. The role prompts read the same file at run time
// ('jaira roles list'); this is for an agent that has no jaira role loaded.
func arrivalSection(mode string) string {
	switch mode {
	case PRAgent:
		return "\n\n## How work arrives on this board\n\n" +
			"`.jaira/pr` says `agent`: work rides on a branch of its own and arrives\n" +
			"through a pull request, and an agent may open that pull request itself —\n" +
			"through `jaira-role-pr`, only once `jaira roles pr-gate` finds every ticket\n" +
			"on the branch accepted. Merging and approving a pull request stay a\n" +
			"person's: no agent does either."
	case PRNone:
		return "\n\n## How work arrives on this board\n\n" +
			"`.jaira/pr` says `none`: this project takes no pull requests. Work still\n" +
			"rides on a branch of its own, and once every ticket on it is in done,\n" +
			"`jaira-role-pr` lands it on the default branch — `jaira roles pr-gate`\n" +
			"first, then a fast-forward push. Do not open a pull request, and never\n" +
			"force a push the default branch refused: bring the branch up to date and\n" +
			"push again."
	}
	return ""
}
