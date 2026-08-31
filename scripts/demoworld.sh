#!/usr/bin/env sh
#
# Build the demo world the README screenshots are taken from.
#
# The screenshots in docs/img/ have to show a board that looks like somebody's
# actual work: several projects, tickets spread over the lanes, one waiting on a
# person. That world used to be assembled by hand, which is why the board
# screenshot went stale and could not be rebuilt. Hence this script.
#
#     scripts/demoworld.sh [dir]        # default: a fresh mktemp -d
#
# It never touches your own board. The world lives entirely under <dir>, and
# JAIRA_HOME, JAIRA_LANES_DIR, JAIRA_DEFAULT_BOARD and JAIRA_USER are pointed
# there, so neither your project registry nor your lane catalogue is read or
# written. It builds jaira from this checkout rather than using the jaira on
# your PATH, so the world matches the code the screenshot is of.
#
# Rendering is the second half of the recipe and is left to you on purpose:
# regenerating all four images when only one is stale is how the other three
# would drift. The script prints these lines with <dir> filled in.
#
#     export JAIRA_HOME=<dir>/home JAIRA_LANES_DIR=<dir>/home/lanes
#     export JAIRA_DEFAULT_BOARD=<dir>/home/default-board.md JAIRA_USER=Demo
#     go run ./scripts/shotgen <dir>/checkout-service board 150 20 \
#         | python3 scripts/termshot.py /dev/stdin docs/img/board.png --cols 150
#
# Views: home, board, pipeline, signoff, edit. The committed images are 150
# columns for board and pipeline, 110 for home and signoff. termshot.py needs
# Pillow and DejaVu Sans Mono.
set -eu

repo=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
dir=${1:-$(mktemp -d 2>/dev/null || mktemp -d -t jaira-demo)}
mkdir -p "$dir"
dir=$(CDPATH= cd -- "$dir" && pwd -P)

JAIRA_HOME="$dir/home"
JAIRA_LANES_DIR="$dir/home/lanes"
JAIRA_DEFAULT_BOARD="$dir/home/default-board.md"
JAIRA_USER=Demo
export JAIRA_HOME JAIRA_LANES_DIR JAIRA_DEFAULT_BOARD JAIRA_USER
mkdir -p "$JAIRA_LANES_DIR"

jaira="$dir/bin/jaira"
(cd "$repo" && go build -o "$jaira" ./cmd/jaira)

# A board wants a git repository: that is where the identity comes from, and a
# ticket's derived commit list needs a history to look at.
newboard() {
	mkdir -p "$dir/$1"
	(
		cd "$dir/$1"
		git init -q .
		git config user.name Demo
		git config user.email demo@example.com
		echo "$1" >README.md
		git add README.md
		git -c commit.gpgsign=false commit -qm "init"
		"$jaira" init >/dev/null
	)
}

# Create a ticket and print its handle. Everything after a create addresses the
# ticket by handle, since the ids are generated and cannot be written down here.
mk() {
	b=$1
	shift
	"$jaira" -C "$dir/$b" create --json "$@" |
		sed -n 's/.*"handle": "\([A-Z0-9]*\)".*/\1/p' | head -1
}

newboard checkout-service
newboard mobile-app
newboard infra-tooling

co="$dir/checkout-service"
j() { "$jaira" -C "$co" "$@" >/dev/null; }

# Backlog: captured but not specified, and marked as needing the brainstorm
# step. It deliberately has no goal: that is what brainstorm produces, and a
# ticket missing one is what the "○ spec" marker on the card is reporting.
backlog=$(mk checkout-service "Decide SameSite for the session cookie" --mine \
	--dod "the setting is chosen, written down, and the redirect keeps the session" \
	--context "The session cookie is set without SameSite, so browsers apply Lax.
The payment provider posts the customer back cross-site and the session is gone.
Strict breaks the same redirect. Nobody has decided between Lax plus a same-site
redirect hop and None plus Secure.")
j dod "$backlog" --option brainstorm

# Todo: specified, waiting to be picked up.
todo=$(mk checkout-service "Retry the vendor webhook with backoff" --mine \
	--goal "A failed webhook delivery is retried instead of dropped" \
	--dod "three retries with growing delay, and the give-up is logged once" \
	--context "The vendor posts order updates to /hooks/vendor. A 500 from our
side is not retried by them and not retried by us, so the order stays unpaid
until somebody notices. Their docs say they deliver once and do not retry.")
j move "$todo" --to todo

# Implementing: an agent is on it, half the plan done.
impl=$(mk checkout-service "Fix session cookie loss on the payment redirect" --mine \
	--goal "The customer comes back from the payment provider still logged in" \
	--dod "the return leg keeps the session, covered by a test" \
	--context "Coming back from the provider the customer lands on an empty cart
and a login form. The session cookie is not sent on the cross-site POST back.
Reproduced against the sandbox, and unrelated to the cart code.")
j move "$impl" --to todo
j dod "$impl" --plan --add "reproduce the drop against the sandbox" \
	--add "set the cookie on the return leg and cover it with a test"
j move "$impl" --to in-progress
j dod "$impl" 1 --plan --done --proof "session_test.go: TestReturnLegKeepsSession reproduces it"
j dod "$impl" 2 --plan --doing

# Human Review: implemented, judged by a model, waiting on a person.
sign=$(mk checkout-service "Rate limit the login endpoint" --mine \
	--goal "Stop credential stuffing without locking out real users" \
	--dod "429 above 100 req/min per IP, covered by a test" \
	--context "The login endpoint has no limit at all. The access log shows one
address trying 40 passwords a minute against real addresses. A per-account lock
was ruled out: it hands anyone a way to lock a customer out.")
j move "$sign" --to todo
j move "$sign" --to in-progress
j dod "$sign" 1 --done --proof "ratelimit_test.go: TestBurstIsCapped"
j move "$sign" --to review \
	--what "added a token bucket per client IP" \
	--why "the endpoint had no limit at all" \
	--resolves "bursts above 100/min now get 429 with Retry-After, verified by TestBurstIsCapped"
j set "$sign" review-summary="a token-bucket limiter per client IP, returning 429 with Retry-After above 100/min" \
	review-gaps="none" \
	review-verdict="the diff matches the criteria; no defects found" \
	review-check="go test ./... passes"
j move "$sign" --to signoff

# Blocked: parked on somebody outside the team.
blocked=$(mk checkout-service "Move the payment sandbox to the new tenant" --mine \
	--goal "The sandbox runs under the tenant the contract names" \
	--dod "sandbox orders settle under the new tenant id" \
	--context "The sandbox still points at the tenant from the trial contract.
The provider has to issue credentials for the new one. The request is in, and
there is nothing to do here until they answer.")
j move "$blocked" --to blocked --reason "provider has not issued the new tenant credentials"

# A second board with one ticket and a third with none, so the launcher and the
# board's project row have something to switch between.
mk mobile-app "Cache the product list offline" --mine \
	--goal "The catalogue opens without a connection" \
	--dod "a cold start with no network shows the last catalogue" \
	--context "In the tunnel the catalogue is an empty list and a spinner. The
response is already JSON we could keep. No offline store exists in the app yet." >/dev/null

# The project row is ordered by when each board was last opened. Writing the
# registry here, rather than relying on the order of the calls above, keeps that
# row identical on every run: those timestamps are only second-accurate, so ties
# would reorder the row at random.
cat >"$JAIRA_HOME/projects.json" <<EOF
[
  {"root": "$dir/checkout-service", "name": "checkout-service", "last_open": "2026-08-31T12:00:03Z"},
  {"root": "$dir/mobile-app", "name": "mobile-app", "last_open": "2026-08-31T12:00:02Z"},
  {"root": "$dir/infra-tooling", "name": "infra-tooling", "last_open": "2026-08-31T12:00:01Z"}
]
EOF

cat <<EOF

Demo world built in $dir

  export JAIRA_HOME=$dir/home JAIRA_LANES_DIR=$dir/home/lanes
  export JAIRA_DEFAULT_BOARD=$dir/home/default-board.md JAIRA_USER=Demo
  go run ./scripts/shotgen $dir/checkout-service board 150 20 \\
      | python3 scripts/termshot.py /dev/stdin docs/img/board.png --cols 150

Or open it: $jaira -C $dir/checkout-service board
EOF
