# Catalogue lanes

Lane files shipped with the repository that are **not** built-ins. A built-in is
compiled into the binary and appears on every board; these do not. They sit here
to be read, copied and argued with, and they reach a board only because someone
adopted them on purpose.

To use one — no clone needed:

```bash
jaira lanes market                    # what is here, fetched from GitHub
jaira lanes market adopt critique     # into your catalogue (~/.jaira/lanes)
jaira lanes add critique              # onto this board
```

From a clone, `jaira lanes adopt lanes/critique.md` does the same as the
second line.

## Adding yours

Open a pull request that adds **one file** here, `<id>.md`, in the shape
`jaira lanes template` prints: frontmatter with at least `id`, `name`, `after`
and `precedence`; `agentic: true` with a `model-tier` if a model works it, and
then `input-requires` / `output-produces` for what it reads and must write back;
the markdown body is the prompt. Add a row to the table below saying where it
sits, where it sends work back to, and what it is for.

Two things get your lane merged: it parses — `core/lane/shipped_test.go` loads
every file in this directory on every CI run, so run `go test ./core/lane/`
before you push — and the prompt says when the lane is *done*, so a loop it
sits in can end.

`jaira lanes show <id>` prints the whole contract and prompt once it is
installed. Read it before adopting: adopting means agreeing to run whatever the
prompt says, at whatever model tier it declares.

## What is here

| Lane | Sits | Sends work back to | For |
|---|---|---|---|
| `secrets-scan` | after implementing, once you move the column there | implementing | Catching a credential that reached a commit — keys, tokens, private keys, a tracked `.env` |
| `critique` | after implementing | implementing | Judging whether this is the right implementation, not whether it works |
| `optimize` | after critique | implementing | Removing duplication, dead code and fluff the change left behind |
| `changelog-writer` | after review | — | Writing the one changelog line for whoever installs the release, rather than for the next agent |

`critique` and `optimize` together make a loop: implementing writes it, critique says what is wrong
with the approach, implementing fixes it, and that repeats until critique has
nothing left to say. Then optimize strips what is not needed, and only then does
the work reach review. A decision that is genuinely the user's goes to the HITL
lane on the way, rather than being taken by whichever agent noticed it.

`secrets-scan` belongs ahead of that loop: it is the cheapest check on the board
and the only one whose miss cannot be taken back — a pushed credential is already
public. `changelog-writer` belongs behind it, after review, once what actually
shipped is settled; it writes one field on the ticket and no file, so two tickets
in it in parallel do not conflict.

**Belongs, not lands.** The `Sits` column above and each lane's `after:` field say
where the lane is *meant* to go; neither moves a column. `jaira lanes add` appends
the lane as the last line of `.jaira/lanes/order`, so a freshly adopted lane is the
rightmost column on the board until you move that line — `after:` is only consulted
when there is no order file at all, and `jaira init` writes one. Rearranging that
file is the adopter's job, and it is a one-line edit.

The loop is not enforced. `rejects-to:` declares the back edge so an agent
reading the board can see it, and a backwards move was always allowed — the gate
only checks a move that advances. A lane may declare two of them —
`rejects-to: [in-progress, human]` — when its rejection is not always the same
kind of thing: a flaw goes back to be implemented, a decision goes to a person.
What stops a bad change is the lane's prompt and the gate on the lane it is
trying to reach, not a state machine.

`core/lane/shipped_test.go` loads everything in this directory on every CI run,
so a lane file here that does not parse fails the build rather than failing the
person who adopts it.
