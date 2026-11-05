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

The file and its layout are the subject of ticket 3Q92JY. Until that lands, a
plugin reads its token from where 3Q92JY puts it, and stores nothing of its own.

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

which is how `jaira web` and a later `jaira plugins` learn what an installed
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
  "capabilities": ["read", "write", "network", "secrets"]
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

`capabilities` is a closed list:

| Value | Means the plugin… |
|---|---|
| `read` | reads tickets through jaira's read commands |
| `write` | changes tickets through jaira's write commands |
| `network` | talks to a host on the network |
| `secrets` | reads `~/.jaira/secrets.json` |

The list is a declaration, not a sandbox: jaira shows it, it does not enforce
it. A plugin that does more than it declares is removed from the catalogue.

## Releases

- Tag `v<version>` on the plugin's repository, matching `version` in its
  manifest.
- Release assets, one per platform, as goreleaser builds them:
  `jaira-<id>_<os>_<arch>` (`.exe` on Windows) — `linux`, `darwin`,
  `windows` × `amd64`, `arm64` — and `checksums.txt` beside them.
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
