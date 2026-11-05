# Plugins

A jaira plugin is a program of its own, written and released by its author in
a repository of their own. jaira stores no plugin code. This document is the
standard such a program follows, and the record of why it looks this way.

**Status:** this is the standard, agreed on 2026-11-05. jaira does not start
plugins yet — `jaira <id>` dispatch, install and the catalogue lookup are a
later ticket. A plugin written to this document today is what that ticket
will run.

Templates and the catalogue:

- Template repository: <https://github.com/panda-bobr/jaira-plugin-template> —
  press **Use this template**; do not fork it.
- Catalogue: `plugins.json` in <https://github.com/panda-bobr/jaira-releases>.
- Schemas in this repository:
  [`docs/schema/jaira-plugin.schema.json`](schema/jaira-plugin.schema.json)
  for a plugin's manifest, and
  [`docs/schema/plugins-entry.schema.json`](schema/plugins-entry.schema.json)
  for one entry of the catalogue.

## Why plugins, and why this shape

Decided by the maintainer on 2026-11-05.

- **One program per plugin, in its author's repository.** The model is
  Obsidian's: a plugin lives where its author works on it, ships its own
  releases, and jaira only knows where to find it.
- **Not Go's `plugin` package.** It does not work on Windows, the plugin must
  be built with exactly the toolchain and dependency versions jaira was built
  with, and it needs cgo — which breaks the single static binary every
  teammate installs.
- **No Jira built into jaira.** Every company uses Jira differently; building
  one company's way in makes jaira grow for everybody and still fit only one.
  The `external` field (below) is the whole of what jaira itself gives an
  adapter.
- **The catalogue is a whitelist.** A plugin is listed because a person merged
  its entry, not because it exists.
- **Later:** a website with a database and registration replaces the pull
  request and mirrors into `plugins.json` by a bot, as Obsidian does today.
  The format of `plugins.json` stays the same when that happens, so nothing
  that reads it has to change.

## Naming

A plugin has an **id**: lower-case letters, digits and single dashes, starting
with a letter — `jira`, `youtrack`, `gitlab-issues`. The pattern is in the
schema.

- The program is called `jaira-<id>` (`jaira-<id>.exe` on Windows).
- jaira's own commands win: an id that names a built-in command (`list`,
  `move`, `web`, …) can never be called and is refused by the catalogue.

## Invocation

```
jaira <id> [args...]
```

starts

```
~/.jaira/plugins/<id>/jaira-<id> [args...]
```

- Every argument after the id is passed on unchanged.
- stdin, stdout and stderr are the plugin's own; jaira adds nothing to them.
- jaira exits with the plugin's exit code.
- The working directory is the one `jaira` was started in.

### Environment

jaira sets three variables before it starts the plugin:

| Variable | Value |
|---|---|
| `JAIRA_ROOT` | the repository root — the directory that holds `.jaira/`. The same meaning the variable already has for hooks |
| `JAIRA_BOARD` | the board the call works on: what `--board` named, else the repository's default board |
| `JAIRA_VERSION` | the version of the running `jaira`, e.g. `0.5.3` |

A plugin compares `JAIRA_VERSION` with its own `minJairaVersion` and exits 3
with a message when jaira is too old.

## Writing: through the CLI only

A plugin never writes under `.jaira/` itself — not `.jaira/tickets/`, not a
ticket ref. It calls the `jaira` binary:

```bash
jaira create "<title>" --goal ... --context ... --dod ... --json
jaira set <id> external=jira:PROJ-123 --json
jaira note <id> "<text>" --json
```

Why: the file format, the merge driver, the ticket refs and the gates are
jaira's. A plugin that writes the file itself bypasses the gates, and breaks
the day the format changes. Through the CLI it gets the same checks a person
gets, and exit codes it can branch on.

Reading is free: `jaira show <id> --json`, `jaira list --json` and the other
read commands are the stable way; parsing the markdown files directly works
but follows the format wherever it goes.

**Today `show --json` and `list --json` do not carry the `external` field.**
A plugin that needs its key has to read it from the ticket file until ticket
8Z3VRF adds it to the JSON output.

## The `external` field

A ticket that mirrors something in another system says so in its frontmatter
field `external`:

```yaml
external: "jira:PROJ-123"
```

- The value is `<id>:<key>` — the plugin's id, a colon, and the key the other
  system uses. `jira:PROJ-123`, `youtrack:APP-42`.
- The field is reserved for exactly this (`core/ticket/schema.go`,
  `FieldExternal`). jaira never interprets it and never rewrites it; on a
  merge, whichever side changed it wins, and both changing it is a conflict.
- Set it with `jaira set <id> external=<id>:<key>`.
- A plugin touches only values that start with its own `<id>:`. A value with
  another prefix belongs to another plugin.

## Secrets

Tokens and passwords live in `~/.jaira/secrets.json` and nowhere else —
not in `settings.json` (written world-readable), not in the repository, not in
the plugin's own folder. The file is written with mode `0600` in a `0700`
directory: readable by its owner only. That is not encryption: every program
of that user, and every backup of the disk, can read it.

The file and its layout are the subject of ticket 3Q92JY; until that ships
the file does not exist. A plugin stores no token of its own anywhere else in
the meantime.

## Events

A command is how a person or an agent calls a plugin. An event is how jaira
calls it on its own: after a write succeeds, every plugin that subscribed to
that write hears about it — a Jira plugin moves the issue when the ticket
moves, without anybody typing `jaira jira push`.

A plugin subscribes in its manifest:

```json
"events": ["create", "move", "claim", "set"]
```

The list is closed; these four are all there are.

| Event | After |
|---|---|
| `create` | `jaira create` filed a ticket |
| `move` | `jaira move` put a ticket into another lane |
| `claim` | `jaira claim` took a ticket |
| `set` | `jaira set` changed a ticket's fields |

After the write has succeeded, jaira runs, for every subscribed plugin in
order of their ids:

```
jaira-<id> event <name>
```

- **stdin** is one JSON object: what `jaira show <id> --json` prints after the
  change, with two keys added at the top level — `actor`, who made the change,
  and `from_lane`, the lane the ticket left (empty for every event but `move`).
- **Environment**: the same `JAIRA_ROOT`, `JAIRA_BOARD` and `JAIRA_VERSION` as
  a command gets.
- **Timeout**: 5 seconds, the same cap jaira's hooks have
  (`core/hook/hook.go`). A plugin still running after that is stopped.
- **Failure never fails the command.** The write has already happened; a
  plugin that exits non-zero, times out or is missing changes nothing about
  the command's own result or exit code. Its stdout and stderr are discarded,
  as a hook's are, so they never mix into JSON an agent is parsing.
- **A plugin that cannot reach its target buffers by itself** — jaira does not
  retry and does not queue. The next event, or a `jaira <id>` command, is the
  plugin's chance to send what it missed.

Events are separate from hooks: `core/hook` stays the place for a user's own
scripts, events are for installed plugins.

## Web

A plugin can add pages of its own to `jaira web` — a Jira plugin, say, that
shows the board the way its team reads Jira.

### Pages

The manifest names them:

```json
"web": {
  "pages": [
    { "path": "index.html", "title": "Jira" }
  ]
}
```

- The files ship in the release as `jaira-<id>_web.tar.gz`, a tarball of a
  `web/` folder, and are unpacked into `~/.jaira/plugins/<id>/web/`.
- `jaira web` serves that folder under `/plugins/<id>/` — the page above is
  `/plugins/jira/index.html` — and shows every page as a tab or link with its
  `title`. `path` is relative to `web/`.
- A page reads the board through `jaira web`'s existing JSON API (`/api/...`),
  the same one jaira's own page uses.

### Its own data

For data only the plugin has, `jaira web` passes a request on, CGI-style:

```
/plugins/<id>/api/<path>   →   jaira-<id> web <METHOD> /<path>
```

- The request body is on stdin; the query string is in `JAIRA_QUERY`.
- The response is JSON on stdout.
- The HTTP status comes from the exit code: `0` → 200, `2` → 400, `3` → 403,
  `5` → 404, anything else → 500.
- The same host and origin checks `jaira web` applies to its own API apply
  here: a write from a foreign origin never reaches the plugin.

One process per request, started and ended by `jaira web`. A plugin runs no
server of its own and opens no port: ports, process lifetime and how Windows
handles background processes would all become `jaira web`'s problem, for
every plugin.

A plugin with pages or an `api` declares the capability `web`.

## Exit codes

The same contract as jaira's own (see [COMMANDS.md](COMMANDS.md#exit-codes)),
because agents branch on them:

| Code | Meaning |
|---|---|
| 0 | success |
| 1 | unexpected error |
| 2 | usage error — bad flags or arguments |
| 3 | refused — the operation was understood and declined (a gate, a missing token, jaira too old) |
| 4 | unresolved dependencies |
| 5 | not found — no such ticket, issue or key |

A refusal says on stderr what to do instead.

## The manifest

Every plugin carries `jaira-plugin.json` at the root of its repository, and
prints the same object for

```
jaira-<id> manifest --json
```

which is how `jaira web` learns what an installed
plugin offers without reading its repository.

```json
{
  "id": "jira",
  "name": "Jira",
  "version": "1.2.0",
  "minJairaVersion": "0.6.0",
  "commands": [
    { "name": "pull", "summary": "create or update tickets from a JQL query" },
    { "name": "push", "summary": "write a ticket's lane back to its Jira issue" }
  ],
  "capabilities": ["read", "write", "network", "secrets", "web"],
  "events": ["move"],
  "web": {
    "pages": [{ "path": "index.html", "title": "Jira" }]
  }
}
```

| Field | Meaning |
|---|---|
| `id` | the plugin's id, as in [Naming](#naming) |
| `name` | what a person reads |
| `version` | semver without a `v`; the release tag is `v<version>` |
| `minJairaVersion` | the oldest jaira the plugin works with |
| `commands` | what `jaira <id> <name>` does, one line each |
| `capabilities` | what the plugin does, shown to a person before install |
| `events` | optional — the writes it wants to hear about, see [Events](#events) |
| `web` | optional — its pages in `jaira web`, see [Web](#web) |

`capabilities` is a closed list:

| Value | Means the plugin… |
|---|---|
| `read` | reads tickets through jaira's read commands |
| `write` | changes tickets through jaira's write commands |
| `network` | talks to a host on the network |
| `secrets` | reads `~/.jaira/secrets.json` |
| `web` | adds pages or an `api` to `jaira web` |

The list is a declaration, not a sandbox: jaira shows it, it does not enforce
it. A plugin that does more than it declares is removed from the catalogue.

## Releases

- Tag `v<version>` on the plugin's repository, matching `version` in its
  manifest.
- Release assets, one per platform, as goreleaser builds them:
  `jaira-<id>_<os>_<arch>` (`.exe` on Windows) — `linux`, `darwin`,
  `windows` × `amd64`, `arm64` — and `checksums.txt` beside them.
- A plugin with pages adds `jaira-<id>_web.tar.gz`, the `web/` folder, once
  for all platforms.
- An update is a new release. It needs no pull request to the catalogue: the
  catalogue names the repository, not a version.

The template repository builds exactly this with its goreleaser config.

## Into the catalogue

1. Fork <https://github.com/panda-bobr/jaira-releases>.
2. Add one entry to `plugins.json`:

   ```json
   {
     "id": "jira",
     "name": "Jira",
     "author": "Jane Doe",
     "description": "Mirror Jira issues into tickets and back.",
     "repo": "https://github.com/jane/jaira-jira"
   }
   ```

   All five fields are required; the schema is
   [`docs/schema/plugins-entry.schema.json`](schema/plugins-entry.schema.json).
   `repo` is the plugin's GitHub repository, where its releases are.
3. Open a pull request. CI checks the entry against the schema, that the id is
   not taken and not a jaira command, and that the repository has a release
   with a `jaira-plugin.json` whose `id` matches.
4. A maintainer merges it. From then on every release of the plugin is
   available without another pull request.
