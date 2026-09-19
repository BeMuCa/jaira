# Catalogue lanes

Lane files shipped with the repository that are **not** built-ins. A built-in is
compiled into the binary and can be installed without a network; these do not
travel in the binary and reach a board only because someone fetched them on
purpose.

To use one — no clone needed:

```bash
jaira lanes market                       # what is here, fetched from GitHub
jaira lanes market adopt secrets-scan    # into your catalogue (~/.jaira/lanes)
jaira lanes add secrets-scan             # onto this board
```

From a clone, `jaira lanes adopt lanes/secrets-scan.md` does the same as the
second line.

`critique`, `optimize` and `testing` are **not** here any more: they travel
inside the binary, so `jaira lanes add critique` installs one with no network
and no adopting. They still stand outside the selection a fresh board starts
with — `jaira init` writes the same ten lanes it always did, and the foot of
`jaira lanes` names the ones this board has not installed.

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
| `secrets-scan` | after implementing | implementing | Catching a credential that reached a commit — keys, tokens, private keys, a tracked `.env` |
| `changelog-writer` | after review | — | Writing the one changelog line for whoever installs the release, rather than for the next agent |

`secrets-scan` belongs ahead of the review loop the binary carries: it is the
cheapest check on the board and the only one whose miss cannot be taken back — a
pushed credential is already public. `changelog-writer` belongs behind it, after review, once what actually
shipped is settled; it writes one field on the ticket and no file, so two tickets
in it in parallel do not conflict.

**Where it lands.** The `Sits` column above and each lane's `after:` field say
where the lane is *meant* to go, and `jaira lanes add` puts it there: the lane is
inserted into `.jaira/lanes/order` just behind the lane its `after:` names, and
only the added id moves — the rest of your column order is left exactly as it is.
The chain is followed through lanes this board has not installed, so a lane whose
anchor is itself uninstalled still lands in the flow rather than at the end. An
anchor nothing in the chain can resolve parks the lane before the first terminal
lane, and says so — a warning on stderr, `warnings` in `--json`; no `after:` at
all parks it there too and says nothing, that being a choice rather than an
omission. The success line names the lane it landed after, because a chain
resolved through lanes you do not have puts it somewhere you never typed.

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
