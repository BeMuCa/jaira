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
# regenerating every image when only one is stale is how the others would
# drift. The script prints these lines with <dir> filled in.
#
#     export JAIRA_HOME=<dir>/home JAIRA_LANES_DIR=<dir>/home/lanes
#     export JAIRA_DEFAULT_BOARD=<dir>/home/default-board.md JAIRA_USER=Demo
#     export JAIRA_NO_UPDATE_CHECK=1
#     <dir>/bin/shotgen <dir>/checkout-service board 170 26 z right right \
#         | python3 scripts/termshot.py /dev/stdin <dir>/board-main.png --cols 170
#     cp <dir>/board-main.png docs/img/board-main.png     # once it looks right
#
# The render lands beside the world, not in docs/img/, so that trying a
# different world does not dirty a committed asset on every attempt. Copying it
# over is the separate, deliberate step. JAIRA_NO_UPDATE_CHECK=1 keeps the
# corner at "jaira <version>": with the check on, it asks GitHub whether a
# newer release exists and the line changes from run to run.
#
# Views: home, board, pipeline, signoff, edit, roadmap; the words after the
# size are keys pressed in order. The committed images were made with:
#
#     board-main.png         board 170 26 z right right
#     board-links.png        board 170 26 z right right right right l
#     launcher.png           home 110 24
#     milestone-roadmap.png  roadmap 100 24
#
# and pipeline.png at 150 columns, signoff.png at 110. z narrows the empty
# lanes, so all ten lanes fit at 170 columns; without it, or narrower, they do
# not all fit and the board scrolls sideways, which can cut a title in half.
# right right rests the cursor on a card without links: the cursor card's link
# lines are drawn across the narrow lanes' upright titles. Four rights reach
# the card being implemented, and l opens its Links window.
#
# termshot.py needs Pillow and DejaVu Sans Mono.
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

# shotgen renders the screenshots. It carries the newest release's version, set
# the way release builds set it, so the corner of every image reads "jaira
# 0.5.1" and not "jaira dev", which is a build nobody can install.
version=$(ls "$repo/core/release/notes" | grep -E '^[0-9]+\.[0-9]+\.[0-9]+$' |
	sort -t. -k1,1n -k2,2n -k3,3n | tail -1)
(cd "$repo" && go build -ldflags "-X github.com/BeMuCa/jaira/core/release.Current=$version" \
	-o "$dir/bin/shotgen" ./scripts/shotgen)

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
#
# The check belongs here rather than at each of the six call sites: extracting
# the handle is a pipeline, so its status is sed's and not jaira's, and set -eu
# would let a refused create through with an empty handle and then fail on some
# later command with an error naming the wrong thing. sh has no pipefail to lean
# on, so the output is captured first and judged on its own.
mk() {
	b=$1
	title=$2
	shift
	out=$("$jaira" -C "$dir/$b" create --json "$@") || {
		echo "demoworld: jaira create refused \"$title\" on $b" >&2
		exit 1
	}
	handle=$(printf '%s\n' "$out" |
		sed -n 's/.*"handle": "\([A-Z0-9]*\)".*/\1/p' | head -1)
	[ -n "$handle" ] || {
		echo "demoworld: no handle in the create output for \"$title\" on $b" >&2
		exit 1
	}
	printf '%s\n' "$handle"
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

# Finished tickets, links and milestones on checkout-service, so the Milestone
# Roadmap has progress to fill and the Links window has something to list.

# Carry a ticket to done: tick, a commit naming it (the move into the terminal
# lane wants one), model review, accept. signoff is passed by: an agent may not
# move work out of it, and this is demo data.
finish() {
	j move "$1" --to todo
	j move "$1" --to in-progress
	j dod "$1" 1 --done --proof "$2"
	git -C "$co" -c commit.gpgsign=false commit -q --allow-empty -m "fix($1): $3"
	j move "$1" --to review --what "$3" --why "$4" --resolves "$5"
	j set "$1" review-summary="$3" review-gaps="none" \
		review-verdict="the diff matches the criteria" review-check="go test ./... passes"
	j move "$1" --to done
}

d1=$(mk checkout-service "Show the order total with tax on the cart page" --mine \
	--goal "The cart shows what the customer will pay" \
	--dod "the cart total includes tax, covered by a test" \
	--context "The cart showed the net total and the checkout the gross one.")
finish "$d1" "cart_test.go: TestTotalIncludesTax" "show the gross total in the cart" \
	"cart and checkout disagreed" "the cart total now includes tax"

d2=$(mk checkout-service "Return 404 for an unknown order id" --mine \
	--goal "An unknown order id is a 404, not a 500" \
	--dod "GET /orders/<unknown> returns 404, covered by a test" \
	--context "A typo in an order link showed the error page and paged on-call.")
finish "$d2" "orders_test.go: TestUnknownOrderIs404" "map a missing order to 404" \
	"a missing row was treated as a server error" "unknown ids return 404"

d3=$(mk checkout-service "Log the provider payment id with each order" --mine \
	--goal "Support can find an order from the provider's payment id" \
	--dod "every paid order stores the provider payment id" \
	--context "Support gets payment ids from the provider and cannot search orders by them.")
finish "$d3" "orders_test.go: TestPaymentIdStored" "store the provider payment id" \
	"support could not match payments to orders" "paid orders carry the payment id"

d4=$(mk checkout-service "Send the receipt mail after the payment settles" --mine \
	--goal "The receipt goes out once the money is there" \
	--dod "the receipt is sent on the settle webhook, not on checkout" \
	--context "Receipts went out on checkout, also for payments that later failed.")
finish "$d4" "receipt_test.go: TestReceiptOnSettle" "send the receipt on settle" \
	"failed payments got a receipt" "receipts follow the settle event"

# Around the ticket being implemented: one follows it, one waits on it, one is
# related to it. --blocked-by gets the full id: create stores it as typed, and
# the blocker's Links window finds a full id only, not a handle.
impl_id=$("$jaira" -C "$co" show "$impl" --json | "$jaira" json get id)
mk checkout-service "Cover the payment return leg in the end-to-end suite" --mine --follows "$impl" \
	--goal "The redirect bug cannot come back unnoticed" \
	--dod "the end-to-end suite pays in the sandbox and checks the session" \
	--context "The fix is covered by a unit test only; the real redirect is not tested." >/dev/null
wait=$(mk checkout-service "Remove the login retry on the return page" --mine --blocked-by "$impl_id" \
	--goal "The return page no longer asks for a second login" \
	--dod "the retry code is gone and the return page keeps the session" \
	--context "A retry was added as a stopgap for the lost session. It goes once the fix lands.")
rel=$(mk checkout-service "Write the cookie settings into the ops runbook" --mine --related "$impl" \
	--goal "Ops knows which cookie flags are set and why" \
	--dod "the runbook lists every cookie with its flags" \
	--context "Ops asked which cookie flags production sets after the redirect incident.")

# Milestones on the roadmap: three in a chain, and one in a swimlane of its own.
j milestone create cart-basics "$d1" "$d2" --after "" \
	--goal "The cart shows correct totals and fails cleanly"
j milestone create payments-v1 "$impl" "$sign" "$blocked" "$d3" --after cart-basics \
	--goal "Customers pay and come back logged in"
j milestone create webhooks "$todo" "$d4" --after payments-v1 \
	--goal "Order updates from the vendor are never lost"
j milestone create sessions "$backlog" "$wait" "$rel" --after "" \
	--goal "Sessions survive cross-site redirects"

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
  export JAIRA_NO_UPDATE_CHECK=1
  $dir/bin/shotgen $dir/checkout-service board 170 26 z right right \\
      | python3 scripts/termshot.py /dev/stdin $dir/board-main.png --cols 170
  cp $dir/board-main.png docs/img/board-main.png     # once it looks right

Or open it: $jaira -C $dir/checkout-service board
EOF
