#!/usr/bin/env python3
"""Facts for a logbook summary: what finished, what is under way, and when work happened.

    python3 collect.py [--from YYYY-MM-DD] [--to YYYY-MM-DD] [--gap 90] [--lead-in 30]
                       [--author <email or name>]... [--repo <path>]

Defaults: the last full calendar week (Monday to Sunday), the git user of the
repository, a new session after 90 minutes without a commit, and 30 minutes
counted before a session's first commit, because a commit closes work and does
not start it.

Prints one JSON object. It only reads: git history, `jaira list --json` and
the files under .jaira/logbook/. Hours are an estimate from commit times and
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


def frontmatter(path):
    # Ticket files are YAML written by jaira: one key per line, strings double-quoted
    # with JSON-compatible escapes, lists as "  - item" lines. That subset is all we read.
    out, key = {}, None
    lines = path.read_text(encoding="utf-8").split("\n")
    if not lines or lines[0] != "---":
        return out
    for line in lines[1:]:
        if line == "---":
            break
        if line.startswith("  - ") and key:
            if not isinstance(out.get(key), list):
                out[key] = []
            out[key].append(line[4:].strip())
            continue
        m = re.match(r"^([\w-]+):\s?(.*)$", line)
        if not m:
            continue
        key, val = m.group(1), m.group(2)
        if val.startswith('"'):
            try:
                val = json.loads(val)
            except json.JSONDecodeError:
                val = val.strip('"')
        elif val == "[]":
            val = []
        out[key] = val
    return out


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
               "--format=%H%x1f%aI%x1f%an%x1f%ae%x1f%s"] + [f"--author={x}" for x in authors], repo)
    commits = []
    for line in log.splitlines():
        h, when, name, mail, subject = line.split("\x1f")
        t = dt.datetime.fromisoformat(when)
        m = SUBJECT.match(subject)
        ids = ID.findall(m.group(1)) if m else []
        commits.append({"sha": h[:10], "at": t, "subject": subject, "ids": ids})
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
        day_rows.append({
            "date": d.isoformat(),
            "weekday": d.strftime("%a"),
            "sessions": [f"{se['start'].strftime('%H:%M')}-{se['end'].strftime('%H:%M')}" for se in ss],
            "commits": len(cs),
            "hours_estimate": hours,
            "by_ticket": by_ticket,
        })
        d += dt.timedelta(days=1)

    # Finished: logbook entries filed inside the window (the folder carries the day).
    done = []
    for f in sorted((repo / ".jaira" / "logbook").glob("*/*.md")):
        m = re.search(r"(\d{8})$", f.parent.name)
        if not m:
            continue
        filed = dt.datetime.strptime(m.group(1), "%Y%m%d").date()
        if not (start <= filed <= end):
            continue
        fm = frontmatter(f)
        done.append({k: fm.get(k) for k in ("id", "title", "goal", "outcome-what", "outcome-why", "tags", "follows")}
                    | {"handle": (fm.get("id") or "")[-6:], "filed": filed.isoformat(), "commits": len(fm.get("commits") or [])})

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
        "days": day_rows,
        "done": done,
        "active": active,
    }, ensure_ascii=False, indent=1, default=str))


if __name__ == "__main__":
    main()
