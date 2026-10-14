#!/usr/bin/env bash
# Worktree + Herdr pane + claude + role prompt for one jaira lane.
# Usage: spawn.sh [--no-worktree] [--parent <session>] <slug> <ticket-id> <lane> [repo-root]
# Prints the pane id on stdout; everything else goes to stderr.
set -euo pipefail

usage() {
  cat <<'USAGE'
spawn.sh [--no-worktree] [--parent <session>] <slug> <ticket-id> <lane> [repo-root]

  --no-worktree   Start the worker in the repository directory itself: no
                  worktree, no branch of its own. For a one-lane job that
                  belongs on the branch already checked out. Two workers then
                  share one directory, so run only one writing worker at a
                  time; a read-only one — the critique that runs beside the
                  work in conversational mode — is the exception, because
                  nothing it does can collide. The slug is still required but
                  unused: it names a worktree, and there is none.
                  JAIRA_NO_WORKTREE=1 in the environment does the same.
  --parent <s>    Only with the lane <dispatch>, <worker> or <discuss>: the
                  session of whoever starts it (a teamlead's name, as
                  ListAgents shows it, or its uds: address). A dispatcher or
                  worker hands its questions there instead of asking the person
                  in its own tab; a discuss chat sends it the decision.
  -h, --help      This text.

Five lane names are special. <dispatch> starts the dispatcher itself in the tab
(/jaira-dispatcher <ticket-id>) instead of one lane's worker; <worker> starts
the worker that carries the whole ticket itself (/jaira-worker <ticket-id>),
which is what a teamlead passes when settings.json says "dispatch": "worker".
<explain> and <discuss> open a chat for a person about the ticket — the review
finding explained, or the ticket's open question talked through — and are what
the jaira-cockpit mod's buttons start, with --no-worktree. They start claude on
the board's .jaira/chat-model, else "chat-model" in ~/.jaira/settings.json, else
the account's default — and get the part of the ticket they need in their
prompt. <pr> starts /jaira-role-pr <ticket-id> --by-person: the cockpit's Accept
starts it when a person accepted a ticket whose branch has no pull request yet,
and nothing else may.
Every other lane name starts /jaira-role-lane <ticket-id> <lane>.

Environment: JAIRA_BRANCH_PREFIX (default feat), JAIRA_NO_WORKTREE,
HERDR_BIN_PATH, HERDR_WORKSPACE_ID.
USAGE
}

# Flags before the positionals, and -- ends them: a slug is free text and could
# one day start with a dash.
no_worktree="${JAIRA_NO_WORKTREE:-}"
parent=""
while [ $# -gt 0 ]; do
  case "$1" in
    --no-worktree) no_worktree=1; shift ;;
    --parent)      parent="${2:?--parent needs a session name}"; shift 2 ;;
    -h|--help)     usage; exit 0 ;;
    --)            shift; break ;;
    -*)            echo "unknown flag: $1" >&2; usage >&2; exit 2 ;;
    *)             break ;;
  esac
done

slug="${1:?slug}"; ticket="${2:?ticket id}"; lane="${3:?lane}"
root="${4:-$PWD}"
herdr="${HERDR_BIN_PATH:-herdr}"
[ "${HERDR_ENV:-}" = 1 ] || { echo "not inside a Herdr pane" >&2; exit 1; }

# Every JSON below is read with 'jaira json get', so a jaira from before it —
# the roles reinstalled, the binary not — would fail right after 'tab create'
# and leave an empty tab standing (#79). Ask first, before a worktree or a tab
# is made.
if ! jaira json get ok --default ok <<<'{}' >/dev/null 2>&1; then
  echo "the jaira on PATH ($(command -v jaira || echo none)) has no 'jaira json get':" \
       "install the new jaira first, then run 'jaira roles install' again" >&2
  exit 1
fi

# Beside the repository, never inside it: a worktree under the repo is a second
# copy of the sources on a different branch, and grep -r / find / ls -R walk
# straight into it. Derived, so no path is baked in for one machine.
# The --no-worktree branch is what usage() above describes. Why anybody wants
# it: for a one-lane job on the branch already checked out — a doc line, a
# note, a lane that only reads — a worktree costs a clone, a branch and a merge
# for nothing.
# A dispatcher runs in its ticket's worktree, and the root it passes is that
# worktree. Deriving from it again would nest a second one beside it, with a
# second copy of the ticket the dispatcher never reads.
if [ "$no_worktree" = 1 ] \
  || { [ "$(basename "$(dirname "$root")")" = .worktrees ] && [ "${root%-"$slug"}" != "$root" ]; }; then
  wt="$root"
else
  wt="$(cd "$root/.." && pwd)/.worktrees/$(basename "$root")-$slug"
fi

if [ ! -d "$wt" ]; then
  # feat/, the prefix this board's own branches use. A repository that names
  # them differently sets JAIRA_BRANCH_PREFIX rather than editing this script.
  git -C "$root" worktree add "$wt" -b "${JAIRA_BRANCH_PREFIX:-feat}/$slug" >&2

  # Only a repo that carries a container stack needs its own ports. A repo
  # without a .env has nothing to offset, and must not fail here.
  if [ -f "$root/.env" ]; then
    # Deterministic per-slug port offset, so two workers never share a stack.
    # 0 stays free: it belongs to the main directory's own stack.
    off=$(( ( $(printf '%s' "$slug" | cksum | cut -d' ' -f1) % 40 ) + 1 ))
    cp "$root/.env" "$wt/.env"
    {
      echo
      echo "# worker stack: slug=$slug offset=$off"
      # Derived from the repository, never a name written in here: a prefix
      # baked into this script belongs to one project and silently names every
      # other project's stack after it. Docker takes lower case only, so the
      # name is folded before it is cleaned — a jaira slug is upper case, and
      # `docker compose` refuses the whole stack over a single capital. printf
      # rather than echo: the newline echo appends is a character like any
      # other to tr, and would come back as a trailing underscore.
      echo "COMPOSE_PROJECT_NAME=$(printf '%s_%s' "$(basename "$root")" "$slug" \
        | tr '[:upper:]' '[:lower:]' | tr -c 'a-z0-9_-' '_')"
      echo "HTTP_PORT=$((8080 + off))"
      echo "DB_PORT_HOST=$((5500 + off * 2))"
      echo "DB_PORT_TEST_HOST=$((5501 + off * 2))"
    } >> "$wt/.env"
  fi
fi

# Where the tab's command has to start claude: inside WSL when this script runs
# in WSL and Herdr is the Windows one, natively on Windows when it runs in Git
# Bash, MSYS or Cygwin — they are Windows, whatever /proc/version there says —
# and in a Unix shell everywhere else.
start_on() {
  case "$(uname -s 2>/dev/null)" in
    MINGW*|MSYS*|CYGWIN*) echo windows; return ;;
  esac
  case "$herdr" in /mnt/*|*.exe) ;; *) echo unix; return ;; esac
  if [ -n "${WSL_DISTRO_NAME:-}" ] || grep -qi microsoft /proc/version 2>/dev/null; then
    echo wsl
  else
    echo unix
  fi
}

# The directory as Herdr takes it. A Windows Herdr reads C:\..., not the
# /c/... Git Bash shows; cygpath, which Git Bash and MSYS ship, translates it.
on="$(start_on)"
cwd="$wt"
if [ "$on" = windows ] && command -v cygpath >/dev/null 2>&1; then
  cwd="$(cygpath -w "$wt")"
fi

# Named, because without --workspace Herdr decides for itself where the tab
# lands: the worker can open in a window nobody is looking at, and a worker
# nobody sees is one nobody notices dying. Herdr exports its own workspace into
# every pane it starts, so this is the workspace the caller is sitting in.
ws=()
if [ -n "${HERDR_WORKSPACE_ID:-}" ]; then ws=(--workspace "$HERDR_WORKSPACE_ID"); fi

# A tab per worker, never a split. A split divides the height of one screen: at
# four workers each strip is a few lines, and nobody can read what any of them is
# doing — which is the whole reason a worker gets a surface of its own.
# JSON is read with 'jaira json get', never python3: on Windows python3 is
# missing or the Store's stub, which prints "Python" and exits 49 (issue #79),
# and jq is not everywhere either. jaira is on every machine that has the roles.
pane="$("$herdr" tab create ${ws[@]+"${ws[@]}"} --cwd "$cwd" --label "$ticket/$lane" --no-focus \
  | jaira json get result.root_pane.pane_id)"

# <explain> and <discuss> write a few lines for a person, so they start on the
# chat model 'jaira roles list' names for this board (its .jaira/chat-model,
# else this machine's settings.json) rather than the account's default model,
# and they get the part of the ticket they need typed into their prompt rather
# than reading the whole of it: 'jaira show --json' carries every note and every DoD proof, most
# of a chat's context for nothing. When jaira cannot answer, both stay as they
# were — the account's model, and a prompt that reads the ticket itself.
chat_model=""
brief=""
if [ "$lane" = explain ] || [ "$lane" = discuss ]; then
  chat_model="$(jaira -C "$wt" roles list --json 2>/dev/null \
    | jaira json get chat_model 2>/dev/null)" || chat_model=""
  # It lands in a command line: only a model name passes.
  [[ "$chat_model" =~ ^[][A-Za-z0-9._-]*$ ]] || chat_model=""
  brief="$(jaira -C "$wt" show "$ticket" --brief 2>/dev/null)" || brief=""
fi
# Their request, brief included, runs to thousands of characters, and typed
# into the tab it lost its middle on the way herdr.exe -> wsl.exe -> claude
# (42C3RZ: 1395 characters arrived, the 'jaira note' instruction gone). It goes
# into a file instead, and the tab is typed one short line that points at it.
# The file sits in a folder of its own outside the repository — inside it,
# git status would show it — and claude starts with that folder added, or its
# Read would ask before opening the file. Native Windows claude reads the
# Windows spelling of the path, which cygpath gives.
chat_dir=""
if [ "$lane" = explain ] || [ "$lane" = discuss ]; then
  chat_dir="$(mktemp -d "${TMPDIR:-/tmp}/jaira-chat.XXXXXX")"
  chat_file="$chat_dir/$ticket-$lane.md"
  chat_dir_win="$chat_dir"; chat_file_win="$chat_file"
  if [ "$on" = windows ] && command -v cygpath >/dev/null 2>&1; then
    chat_dir_win="$(cygpath -w "$chat_dir")"; chat_file_win="$(cygpath -w "$chat_file")"
  fi
fi
claude="claude${chat_model:+ --model '$chat_model'}${chat_dir:+ --add-dir '$chat_dir'}"

# The tab's shell runs on the machine Herdr itself runs on. When that is Windows
# and this script runs in WSL, that shell has no /home/... at all: --cwd is
# resolved against Windows and dropped, `cd "$wt"` fails, and claude comes up in
# the Windows home in front of its trust dialog — which the send-keys below would
# then answer on the human's behalf. Cross back into WSL instead: wsl.exe --cd
# sets the directory before any shell starts, so no cd function and no wrong home
# can intervene, and `bash -lic` is what puts claude on PATH.
# Only when this script really runs in WSL: a herdr.exe is just as much the
# Herdr of a native Windows machine, where this runs in Git Bash, the tab is
# already on the right machine, and wsl.exe would start in a distro that cannot
# reach C:\ under that path and leave claude at its trust dialog in / (#79).
# On Windows the tab's shell may be cmd or PowerShell 5.1, which know neither
# && nor single quotes: --cwd above has put the tab in the worktree already, so
# it gets the bare command, which every one of them runs alike. The model name
# needs no quotes — only letters, digits and . _ - [ ] pass.
case "$on" in
  wsl)     start="wsl.exe --cd '$wt' -- bash -lic \"$claude\"" ;;
  windows) start="claude${chat_model:+ --model $chat_model}${chat_dir:+ --add-dir \"$chat_dir_win\"}" ;;
  *)       start="cd '$cwd' && $claude" ;;
esac

"$herdr" pane run "$pane" "$start" >/dev/null

# The state hook reports idle a few seconds after startup.
for _ in $(seq 20); do
  sleep 2
  pj="$("$herdr" pane get "$pane")"
  st="$(jaira json get result.pane.agent --default - <<<"$pj") $(jaira json get result.pane.agent_status --default - <<<"$pj")"
  case "$st" in claude\ idle|claude\ done) break ;; esac
done
# Only the two states the loop breaks on may pass. Anything else — above all
# Herdr's `blocked`, its state for a detected approval dialog — must stop here:
# the send-keys below would answer that dialog on the human's behalf, and this
# role is forbidden from ever doing that.
case "${st:-}" in
  claude\ idle|claude\ done) ;;
  claude\ blocked)
    echo "claude is up in $pane but an approval dialog is waiting:" \
         "report it to the human, let them answer it in that pane," \
         "then start this worker again" >&2
    exit 1 ;;
  *) echo "claude did not come up in $pane: ${st:-none}" >&2; exit 1 ;;
esac

# The lane name "dispatch" starts the dispatcher itself rather than one lane's
# worker — same worktree, same tab, same waiting for claude to come up, so it
# lives here rather than in a second script that drifts from this one.
# "worker" is the same for the role that carries a ticket through every lane
# itself (jaira-worker), the alternative a teamlead starts when the dispatch
# setting asks for it.
# <explain> and <discuss> are no agent's lane at all: the cockpit mod opens
# them for a person, a chat of its own instead of the teamlead's. They type a
# plain request, so the tab answers the person and never works the ticket.
# Without Herdr the cockpit asks its own session the same two requests
# (explainBySubagent and discussInThisChat in jaira-cockpit's hooks/texts.ts):
# a change to one belongs in the other.
if [ "$lane" = dispatch ]; then
  prompt="/jaira-dispatcher $ticket${parent:+ --parent $parent}"
elif [ "$lane" = worker ]; then
  prompt="/jaira-worker $ticket${parent:+ --parent $parent}"
elif [ "$lane" = pr ]; then
  # A person pressed Accept in the cockpit: --by-person tells jaira-role-pr
  # the invocation is theirs, so it opens the pull request instead of handing
  # the create line back. Only that press passes this lane.
  prompt="/jaira-role-pr $ticket --by-person"
elif [ "$lane" = explain ]; then
  src=" (jaira show $ticket --json: review.verdict, review.summary, review.gaps, review.check, dazu die Notes)"
  [ -z "$brief" ] || src=""
  prompt="Erklaere mir in einfachen Worten den Review-Befund zu Ticket $ticket$src: was genau gefunden wurde, warum es zaehlt, welche Optionen es gibt und was der Reviewer empfiehlt. Schreibe danach dieselbe Erklaerung an das Ticket, als eine Note mit genau diesen Zeilen, jede Zeile ein Teil: jaira note $ticket \"explain: <Befund: was gefunden wurde, ein bis zwei Saetze>\nWarum: <warum es zaehlt, ein bis zwei Saetze>\n1) <erste Option, eine Zeile> [give back]\n2) <zweite Option, eine Zeile> [accept]\nEmpfehlung: <Nummer> - <warum, ein Satz>\" - mit echten Zeilenumbruechen im Text, zwei bis vier Optionen mit 1) 2) 3) 4) am Zeilenanfang - dieselben Nummern, mit denen die Person im Cockpit waehlt, die Woerter Warum: und Empfehlung: genau so. Ans Ende jeder Option gehoert [give back], wenn sie heisst, die Arbeit zurueckzugeben und etwas zu aendern, [accept], wenn sie heisst, das Ticket so anzunehmen, und keine Marke, wenn sie keins von beidem ist. Das Cockpit zeigt die Note unter [Read explanation], und die Wahl dort fuehrt genau diesen Schritt aus. Sonst aendere nichts, committe nichts, bewege das Ticket nicht. Antworte in der Sprache, in der ich mit dir spreche."
elif [ "$lane" = discuss ]; then
  src=" (jaira show $ticket --json: das Feld question und die letzten Notes)"
  [ -z "$brief" ] || src=""
  prompt="Ticket $ticket wartet auf eine Entscheidung von mir$src. Erklaere mir die Frage und die Optionen in einfachen Worten und frag mich. Wenn ich entschieden habe: schreibe die Entscheidung in einem Satz mit jaira note $ticket \"decision: <Entscheidung>\" an das Ticket und leere die Frage mit jaira set $ticket question=. Heisst die Entscheidung, dass weitergearbeitet wird, bewege das Ticket selbst zurueck, genau so: jaira move $ticket --to in-progress (ohne --what/--why/--resolves: das Ticket hat sein Ergebnis schon, und die Entscheidung steht in der Note). In die naechste Lane annehmen nur, wenn ich ausdruecklich ja dazu sage.${parent:+ Schicke danach eine Zeile mit der Entscheidung per SendMessage an $parent.} Aendere keinen Code. Der Tab schliesst sich danach von selbst. Antworte in der Sprache, in der ich mit dir spreche."
else
  prompt="/jaira-role-lane $ticket $lane"
fi
if [ -n "$brief" ]; then
  prompt="$prompt Was du vom Ticket brauchst, steht hier - title, goal, question, review und die letzten Notes als JSON; lies das Ticket nicht noch einmal ganz: $brief"
fi
if [ -n "$chat_dir" ]; then
  printf '%s\n' "$prompt" > "$chat_file"
  if [ "$on" = windows ]; then chat_file="$chat_file_win"; fi
  prompt="Lies die Datei $chat_file ganz und tu, was darin steht."
fi
"$herdr" pane send-text "$pane" "$prompt" >/dev/null
sleep 1
"$herdr" pane send-keys "$pane" enter >/dev/null

echo "$pane"
