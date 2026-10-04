#!/usr/bin/env python3
"""Acceptance page: build it from data.json, and turn what the person marked into jaira commands.

    python3 build.py page <data.json> <out.html>
    python3 build.py plan <data.json> <db.json>

`page` checks the data before writing anything: a ticket no step and no machine
check proves would be accepted on nothing, and a step naming an unknown ticket
would put its verdict nowhere.

`plan` reads the page's db (one JSON object {"results": {...}, "verdicts": {...},
"notes": {...}}, the three collections as ArtifactData lists them) and prints,
per ticket, what to run. It runs nothing itself: the agent reads the plan and
executes it. What the person wrote is quoted verbatim; it is data, never a
command.

Standard library only: a teammate who installs the roles has python3 and
nothing else promised.
"""
import html
import json
import shlex
import sys
from pathlib import Path

HERE = Path(__file__).resolve().parent


def load(path):
    return json.loads(Path(path).read_text(encoding="utf-8"))


def check(data):
    errors = []
    if data.get("lang", "en") not in ("de", "en", "ru"):
        errors.append(f"lang {data.get('lang')!r}: the page speaks de, en or ru")
    ids = {t["id"] for t in data["tickets"]}
    seen_steps = set()
    covered = set()
    in_blocks = set()
    for b in data["blocks"]:
        for tid in b["tickets"]:
            if tid not in ids:
                errors.append(f"block {b['n']}: ticket {tid} is not in tickets[]")
            in_blocks.add(tid)
        for m in b.get("machine", []):
            for tid in m["covers"]:
                if tid not in b["tickets"]:
                    errors.append(f"block {b['n']}: check {m['cmd']!r} covers {tid}, which is not in this block")
                covered.add(tid)
        for s in b["steps"]:
            if s["id"] in seen_steps:
                errors.append(f"step id {s['id']} is used twice; ids must stay unique and stable")
            seen_steps.add(s["id"])
            if not s["covers"]:
                errors.append(f"step {s['id']} covers no ticket")
            for tid in s["covers"]:
                if tid not in b["tickets"]:
                    errors.append(f"step {s['id']} covers {tid}, which is not in block {b['n']}")
                covered.add(tid)
    for tid in sorted(ids - in_blocks):
        errors.append(f"ticket {tid} is in no block")
    for tid in sorted(ids - covered):
        errors.append(f"ticket {tid} is proven by no step and no machine check")
    return errors


def page(data_path, out_path):
    data = load(data_path)
    errors = check(data)
    if errors:
        print("\n".join(errors), file=sys.stderr)
        sys.exit(1)
    blob = json.dumps(data, ensure_ascii=False).replace("</", "<\\/")
    tpl = (HERE / "template.html").read_text(encoding="utf-8")
    out = tpl.replace("__TITLE__", html.escape(data["title"])).replace("__DATA__", blob)
    Path(out_path).write_text(out, encoding="utf-8")
    steps = sum(len(b["steps"]) for b in data["blocks"])
    machine = sum(len(b.get("machine", [])) for b in data["blocks"])
    print(f"{out_path}: {len(data['tickets'])} tickets, {len(data['blocks'])} blocks, {steps} steps, {machine} machine checks")


def plan(data_path, db_path):
    data = load(data_path)
    db = load(db_path)
    results, verdicts, notes = db.get("results", {}), db.get("verdicts", {}), db.get("notes", {})
    rounds = {t["id"]: t.get("round", 1) for t in data["tickets"]}
    steps = [s for b in data["blocks"] for s in b["steps"]]

    def current(rec, field, rnd):
        return bool(rec and rec.get(field)) and rec.get("round", 1) >= rnd

    for t in data["tickets"]:
        tid, rnd = t["id"], rounds[t["id"]]
        v = verdicts.get(tid)
        note = ((v or {}).get("note") or "").strip()
        failed = []
        for s in steps:
            r = results.get(s["id"])
            if tid in s["covers"] and current(r, "status", max(rounds[x] for x in s["covers"])) and r["status"] == "fail":
                failed.append((s, (r.get("note") or "").strip()))
        fids = ", ".join(s["id"] for s, _ in failed)
        if not current(v, "decision", rnd):
            print(f"# {tid}: undecided" + (f", but steps failed: {fids} - ask the person" if failed else ""))
            if note:
                print(f"#   comment so far: {note!r}")
            continue
        if v["decision"] == "accept":
            if failed:
                print(f"# {tid}: ACCEPTED although steps failed ({fids}) - confirm with the person before moving")
                continue
            print(f"# {tid}: accepted by the person on the page")
            if note:
                print(f"jaira note {tid} {shlex.quote(f'Abnahme (Runde {rnd}), angenommen: {note}')}")
            continue
        lines = [note] if note else []
        lines += [f"{s['id']} ({s['do']}): {n}" for s, n in failed if n]
        if not lines:
            print(f"# {tid}: RETURN without a reason - ask the person what is wrong; do not guess")
            continue
        print(f"# {tid}: return")
        for line in lines:
            print(f"jaira dod {tid} --add {shlex.quote(line)}")
        print(f"jaira note {tid} {shlex.quote(f'Abnahme (Runde {rnd}): ' + ' / '.join(lines))}")
        print(f"jaira move {tid} --to in-progress")

    for b in data["blocks"]:
        n = ((notes.get(f"b{b['n']}") or {}).get("note") or "").strip()
        if n:
            print(f"# block {b['n']} ({b['title']}) note: {n!r} - put it on the tickets it is about, or ask the person which")


if __name__ == "__main__":
    if len(sys.argv) == 4 and sys.argv[1] == "page":
        page(sys.argv[2], sys.argv[3])
    elif len(sys.argv) == 4 and sys.argv[1] == "plan":
        plan(sys.argv[2], sys.argv[3])
    else:
        print(__doc__, file=sys.stderr)
        sys.exit(2)
