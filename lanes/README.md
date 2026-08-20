# Catalogue lanes

Lane files shipped with the repository that are **not** built-ins. A built-in is
compiled into the binary and appears on every board; these do not. They sit here
to be read, copied and argued with, and they reach a board only because someone
adopted them on purpose.

To use one:

```bash
jaira lanes adopt lanes/critique.md   # into your catalogue (~/.jaira/lanes)
jaira lanes add critique              # onto this project's board
```

`jaira lanes show <id>` prints the whole contract and prompt once it is
installed. Read it before adopting: adopting means agreeing to run whatever the
prompt says, at whatever model tier it declares.

## What is here

| Lane | Sits | Sends work back to | For |
|---|---|---|---|
| `critique` | after implementing | implementing | Judging whether this is the right implementation, not whether it works |
| `optimise` | after critique | implementing | Removing duplication, dead code and fluff the change left behind |

Together they make a loop: implementing writes it, critique says what is wrong
with the approach, implementing fixes it, and that repeats until critique has
nothing left to say. Then optimise strips what is not needed, and only then does
the work reach review. A decision that is genuinely the user's goes to the HITL
lane on the way, rather than being taken by whichever agent noticed it.

The loop is not enforced. `rejects-to:` declares the back edge so an agent
reading the board can see it, and a backwards move was always allowed — the gate
only checks a move that advances. What stops a bad change is the lane's prompt
and the gate on the lane it is trying to reach, not a state machine.

`core/lane/shipped_test.go` loads everything in this directory on every CI run,
so a lane file here that does not parse fails the build rather than failing the
person who adopts it.
