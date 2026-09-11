#!/bin/bash
# Try everything the git-ref feature does, in a sandbox of its own. Touches
# nothing outside $R.
#   bash scripts/try-everything.sh        (J=/path/to/jaira picks another binary)
set -u
J=${J:-/tmp/jaira}
R=$(mktemp -d); export JAIRA_HOME=$R/home
say() { echo; echo "=== $*"; }
refs() { git -C "$1" ls-remote origin 'refs/jaira/tickets/*' 2>/dev/null | wc -l; }
files() { ls "$1"/.jaira/tickets 2>/dev/null | wc -l; }
id_of() { grep -o '"id": "[^"]*"' | head -1 | cut -d'"' -f4; }

say "0. three clones of one bare repo: ada, grace, carol"
git init -q --bare $R/board.git
for n in ada grace carol; do
  git clone -q $R/board.git $R/$n 2>/dev/null
  git -C $R/$n config user.name $n; git -C $R/$n config user.email $n@x.test
  $J -C $R/$n init >/dev/null
done
echo "   sandbox: $R"

say "1. ada files a ticket for grace — and gets no file"
A=$($J -C $R/ada create "session cookie dropped on 302" --goal "the cookie survives the OAuth round-trip" \
      --context "reported in chat while debugging Safari logouts" --dod "survives the redirect" --assignee grace --json | id_of)
echo "   files at ada: $(files $R/ada)   refs on the remote: $(refs $R/ada)     (expect 0 and 1)"

say "2. grace did nothing but read"
$J -C $R/grace fetch --quiet
echo "   ^ '@you new ref-only' = assigned to you, and in nobody's branch here"

say "3. listing and showing work without a file, writing does not"
$J -C $R/grace list | sed -n '2,3p'
$J -C $R/grace note $A "trying to write"; echo "   ^ exit=$? (expect 3)"

say "4. carol tries to take a ticket that is not hers"
$J -C $R/carol fetch --quiet >/dev/null; $J -C $R/carol pull $A; echo "   exit=$? files at carol: $(files $R/carol)   (expect 3 and 0)"

say "5. grace takes his own — the file appears at HIS clone only"
$J -C $R/grace pull $A | head -1
echo "   grace: $(files $R/grace)   ada: $(files $R/ada)   carol: $(files $R/carol)     (expect 1 0 0)"

say "6. now grace can write"
$J -C $R/grace note $A "started"; echo "   exit=$?   (expect 0)"

say "7. grace hands the ticket back"
$J -C $R/grace release $A | head -2
echo "   files at grace: $(files $R/grace)   (expect 0)"

say "8. so carol can take it"
$J -C $R/carol fetch --quiet >/dev/null; $J -C $R/carol pull $A | head -1

say "9. ada steals it — loudly, and with a trace on the ticket"
$J -C $R/ada fetch --quiet >/dev/null; $J -C $R/ada pull $A --steal | head -2
$J -C $R/ada show $A | grep -A2 'Progress' | tail -2

say "10. offline: the write stays here and waits"
# count only ada's queue: the sandbox shares one JAIRA_HOME across the clones
outbox_ada() { ls $R/home/state/ada-*/outbox/ 2>/dev/null | wc -l; }
git -C $R/ada remote set-url origin https://127.0.0.1:1/nope.git
$J -C $R/ada note $A "written with no route" >/dev/null 2>&1
echo "   queued at ada: $(outbox_ada)   (expect 1)"
git -C $R/ada remote set-url origin $R/board.git
$J -C $R/ada note $A "the network is back" >/dev/null 2>&1
echo "   left after the network came back: $(outbox_ada)   (expect 0)"

say "11. background fetch: carol only reads, and still learns about new work"
B=$($J -C $R/ada create "a second ticket" --goal g --context "to prove the background fetch" --dod d --json | id_of)
rm -f $R/home/state/*/fetch.json
/usr/bin/time -f "   list took %e s" $J -C $R/carol list >/dev/null
sleep 1; echo "   tickets carol sees now: $($J -C $R/carol list --json 2>/dev/null | grep -c '"id"')   (expect 2)"

say "12. the snapshot: the board as files on a branch of its own"
$J -C $R/ada snapshot
git -C $R/ada ls-tree -r --name-only jaira/board
echo "   running it again:"; $J -C $R/ada snapshot

say "13. finish a ticket — the ref SURVIVES being archived"
$J -C $R/ada share >/dev/null 2>&1; git -C $R/ada add -A .jaira .gitignore >/dev/null 2>&1
git -C $R/ada commit -qm "share" >/dev/null 2>&1; git -C $R/ada push -q -u origin master 2>/dev/null
$J -C $R/ada archive $A >/dev/null
echo "   refs after archive: $(refs $R/ada)   (expect: still there)"
$J -C $R/ada snapshot | tail -1
echo "   after a snapshot with nothing landed: $(refs $R/ada)   (expect: still there)"

say "14. land it in master — the snapshot reaps the ref"
git -C $R/ada add -A .jaira >/dev/null; git -C $R/ada commit -qm "archive" >/dev/null; git -C $R/ada push -q origin master
$J -C $R/ada snapshot | tail -1
echo "   refs: $(refs $R/ada)   (expect: one fewer)"
echo "   and the ticket is still in the snapshot:"; git -C $R/ada ls-tree -r --name-only jaira/board | head -3

say "15. the board's history is plain git"
git -C $R/ada log --oneline jaira/board | head -4

say "16. two branches that touched the same ticket (add/add)"
cd $R/ada
git checkout -q -b work; $J pull $B >/dev/null 2>&1; $J set $B assignee=ada >/dev/null
git add -A .jaira >/dev/null; git commit -qm "work on $B" >/dev/null
git checkout -q master; git checkout -q -b other; git checkout -q work -- .jaira/tickets 2>/dev/null
F=$(ls .jaira/tickets | head -1)
[ -n "$F" ] && sed -i 's/^status: .*/status: todo/' .jaira/tickets/$F && git add -A .jaira >/dev/null && git commit -qm "other" >/dev/null
git checkout -q work; git merge --no-edit other 2>&1 | tail -1
echo "   conflict markers: $(grep -c '<<<<<<<' .jaira/tickets/$F 2>/dev/null | head -1)   (expect 0)"
cd - >/dev/null

say "17. a board with no remote behaves as it always did"
mkdir -p $R/plain && $J -C $R/plain init >/dev/null
$J -C $R/plain create "no remote here" --goal g --context "a directory that is not a repository" --dod d | tail -1
echo "   files: $(files $R/plain)   (expect 1 — the file stays)"

echo; echo "=== sandbox kept at: $R"
echo "    remove it with:  rm -rf $R"
