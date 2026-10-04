#!/usr/bin/env python3
"""Facts for a logbook summary: what finished, what is under way, and when work happened.

    python3 collect.py [--from YYYY-MM-DD] [--to YYYY-MM-DD] [--gap 90] [--lead-in 30]
                       [--author <email or name>]... [--repo <path>]

Defaults: the last full calendar week (Monday to Sunday), the git user of the
repository, a new session after 90 minutes without a commit, and 30 minutes
counted before a session's first commit, because a commit closes work and does
not start it.

Prints one JSON object. Hours per ticket (`hours_by_ticket`, per day and for
the period) share each day's estimate out by the time between commits; they
add up to the day's hours exactly. It only reads: git history, `jaira list --json` and
`jaira logbook` and `jaira show`. Hours are an estimate from commit times and
are labelled so; a day with no commits has no hours here, whatever happened on
it.

Standard library only: the roles promise python3 and nothing else.
"""
import argparse
import datetime as dt
import json
import re
import subprocess
from collections import defaultdict
from pathlib import Path

ID = re.compile(r"\b([0-9A-Z]{6})\b")
SUBJECT = re.compile(r"^\w+\(([^)]*)\)!?:\s*(.*)$")


def run(args, cwd):
    return subprocess.run(args, cwd=cwd, capture_output=True, text=True, check=False).stdout


def last_week(today):
    monday = today - dt.timedelta(days=today.weekday() + 7)
    return monday, monday + dt.timedelta(days=6)


def split_hours(sessions, lead, day_hours):
    """A day's hours shared out over the tickets its commits name.

    Inside a session, the time since the previous commit belongs to the ticket
    of the commit that ended it; the first commit gets the lead-in. A commit
    naming two tickets splits its time evenly, one naming none counts as "-".
    The day's rounded hours are then dealt out in quarter hours by largest
    remainder, so the tickets of a day add up to exactly that day's hours.
    """
    minutes = defaultdict(float)
    for se in sessions:
        prev = se["start"] - lead
        for c in se["commits"]:
            span = (c["at"] - prev).total_seconds() / 60
            ids = c["ids"] or ["-"]
            for i in ids:
                minutes[i] += span / len(ids)
            prev = c["at"]
    whole = sum(minutes.values())
    if not whole:
        return {}
    quarters = round(day_hours * 4)
    shares = {i: m / whole * quarters for i, m in minutes.items()}
    out = {i: int(x) for i, x in shares.items()}
    for i in sorted(shares, key=lambda i: shares[i] - out[i], reverse=True)[:quarters - sum(out.values())]:
        out[i] += 1
    return {i: q / 4 for i, q in sorted(out.items(), key=lambda kv: -kv[1]) if q}


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--from", dest="start")
    ap.add_argument("--to", dest="end")
    ap.add_argument("--gap", type=int, default=90)
    ap.add_argument("--lead-in", type=int, default=30)
    ap.add_argument("--author", action="append")
    ap.add_argument("--repo", default=".")
    a = ap.parse_args()

    repo = Path(run(["git", "rev-parse", "--show-toplevel"], a.repo).strip() or a.repo)
    start, end = last_week(dt.date.today())
    if a.start:
        start = dt.date.fromisoformat(a.start)
    if a.end:
        end = dt.date.fromisoformat(a.end)
    authors = a.author or [run(["git", "config", "user.email"], repo).strip()]

    # Commits: every local branch, no merges, by the given authors.
    log = run(["git", "log", "--branches", "--no-merges", "--date=iso-strict",
               f"--since={start}T00:00:00", f"--until={end}T23:59:59",
               "--format=%aI%x1f%s"] + [f"--author={x}" for x in authors], repo)
    commits = []
    for line in log.splitlines():
        when, subject = line.split("\x1f")
        t = dt.datetime.fromisoformat(when)
        m = SUBJECT.match(subject)
        ids = ID.findall(m.group(1)) if m else []
        commits.append({"at": t, "subject": subject, "ids": ids})
    commits.sort(key=lambda c: c["at"])

    # Sessions run over the whole period, so work past midnight stays one
    # session; a session and its commits count on the day it started.
    gap, lead = dt.timedelta(minutes=a.gap), dt.timedelta(minutes=a.lead_in)
    sessions = []
    for c in commits:
        if sessions and c["at"] - sessions[-1]["end"] <= gap:
            sessions[-1]["end"] = c["at"]
            sessions[-1]["commits"].append(c)
        else:
            sessions.append({"start": c["at"], "end": c["at"], "commits": [c]})
    by_day = defaultdict(list)
    for se in sessions:
        by_day[se["start"].date()].append(se)

    day_rows, total = [], 0.0
    ticket_total = defaultdict(float)
    d = start
    while d <= end:
        ss = by_day.get(d, [])
        cs = [c for se in ss for c in se["commits"]]
        minutes = sum(((se["end"] - se["start"]) + lead).total_seconds() / 60 for se in ss)
        hours = round(minutes / 30) / 2
        total += hours
        by_ticket = defaultdict(list)
        for c in cs:
            for i in c["ids"] or ["-"]:
                by_ticket[i].append(c["subject"])
        hours_by_ticket = split_hours(ss, lead, hours)
        for i, h in hours_by_ticket.items():
            ticket_total[i] += h
        day_rows.append({
            "date": d.isoformat(),
            "weekday": d.strftime("%a"),
            "sessions": [f"{se['start'].strftime('%H:%M')}-{se['end'].strftime('%H:%M')}" for se in ss],
            "commits": len(cs),
            "hours_estimate": hours,
            "by_ticket": by_ticket,
            "hours_by_ticket": hours_by_ticket,
        })
        d += dt.timedelta(days=1)

    # Finished: logbook entries filed inside the window. The folder name
    # carries the day; jaira itself reads the file.
    listing = json.loads(run(["jaira", "logbook", "--since", "0", "--json"], repo) or '{"logbook": []}')
    done = []
    for rel in listing.get("logbook") or []:
        m = re.search(r"(\d{8})/[0-9A-Z]{20}([0-9A-Z]{6})-", rel)
        if not m:
            continue
        filed = dt.datetime.strptime(m.group(1), "%Y%m%d").date()
        if not (start <= filed <= end):
            continue
        t = json.loads(run(["jaira", "show", m.group(2), "--json"], repo) or "{}")
        done.append({k: t.get(k) for k in ("handle", "title", "goal", "outcome", "tags", "follows")}
                    | {"filed": filed.isoformat(), "commits": len(t.get("commits") or [])})

    # Under way: tickets on the board that moved inside the window.
    board = json.loads(run(["jaira", "list", "--json"], repo) or '{"tickets": []}')["tickets"]
    active = []
    for t in board:
        upd = (t.get("updated_at") or "")[:10]
        if not upd or not (start.isoformat() <= upd <= end.isoformat()):
            continue
        active.append({k: t.get(k) for k in ("handle", "title", "goal", "status", "tags", "follows", "milestones")}
                      | {"dod": f"{sum(1 for x in t.get('dod_items') or [] if x.get('done'))}/{len(t.get('dod_items') or [])}"})

    print(json.dumps({
        "period": {"from": start.isoformat(), "to": end.isoformat()},
        "authors": authors,
        "session_rule": {"gap_minutes": a.gap, "lead_in_minutes": a.lead_in},
        "hours_total_estimate": total,
        "hours_by_ticket": dict(sorted(ticket_total.items(), key=lambda kv: -kv[1])),
        "days": day_rows,
        "done": done,
        "active": active,
    }, ensure_ascii=False, indent=1, default=str))


if __name__ == "__main__":
    main()
