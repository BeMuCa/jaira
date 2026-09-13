#!/bin/sh
# A jaira notification hook. Print it with 'jaira hook example'.
#
# jaira runs this script after every 'move' and every 'claim', with the ticket
# in the environment:
#
#   JAIRA_EVENT     "move" or "claim"
#   JAIRA_TICKET    the ticket's full id
#   JAIRA_TITLE     its title
#   JAIRA_STATUS    the lane it is in now
#   JAIRA_ASSIGNEE  who holds it
#   JAIRA_ACTOR     who did it
#   JAIRA_ROOT      the board's directory
#
# jaira kills the script after 5 seconds and ignores whatever it exits with, so
# the delivery has to be quick and must never wait on a host that is down.
#
# The rule this example carries: only a state in which nothing moves without a
# person is worth a sound. That rule has exactly two exceptions here, and both
# are deliberate: a lane that waits on a person (human, signoff), and the lane
# that ends the ticket (done). done waits on nobody -- it rings because a
# finish is the other thing you would rather be told than have to go and look
# for. Every other lane stays silent, or the sound stops meaning anything and
# gets switched off.
#
# To deliver somewhere other than this terminal, replace the two printf lines in
# deliver() below. Keep the request/finished split; that is the whole point.
#
#   ntfy.sh      curl -fsS -H "Title: $2" -d "$3" https://ntfy.sh/your-topic >/dev/null 2>&1
#   desktop      notify-send "$2" "$3"
#   herdr        herdr notification show --title "$2" --body "$3"

set -u

# deliver <tone> <title> <body>. The terminal bell is the one channel every
# machine already has. Read the printed line, not the ringing: "jaira human:"
# against "jaira done:" is what actually tells a request from a finish. The
# bell count -- two rings for a request, one for a finish -- is a hint and no
# more, because plenty of terminals merge or throttle bells that arrive in the
# same breath and sound both of them as one tone. A channel that carries a
# title of its own, like the ntfy or notify-send line above, gets the
# distinction back properly.
deliver() {
	case "$1" in
	request) printf '\a\a%s: %s\n' "$2" "$3" ;;
	*) printf '\a%s: %s\n' "$2" "$3" ;;
	esac
}

# A claim moves nothing and waits for nobody.
[ "${JAIRA_EVENT:-}" = "move" ] || exit 0

# The lanes that belong to a person, and the lane that ends the ticket. The
# names are this board's; change them if yours are named differently.
case "${JAIRA_STATUS:-}" in
human | signoff) tone=request ;;
done) tone=finished ;;
*) exit 0 ;;
esac

title="jaira $JAIRA_STATUS"
body="${JAIRA_TITLE:-${JAIRA_TICKET:-}}"

# Straight to the terminal when there is one, so the bell survives a caller that
# discards stdout — which is exactly what jaira does with it. The probe is a
# subshell because a redirection that fails in one costs a subshell rather than
# the script, and because only opening /dev/tty tells you whether this process
# still has a terminal; the file is there either way.
if ( : >/dev/tty ) 2>/dev/null; then
	deliver "$tone" "$title" "$body" >/dev/tty 2>/dev/null
else
	deliver "$tone" "$title" "$body"
fi

exit 0
