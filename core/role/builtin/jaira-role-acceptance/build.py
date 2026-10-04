#!/usr/bin/env python3
"""Acceptance page: build it, read the person's marks back, and write the record.

    python3 build.py page   <data.json> <page.html>
    python3 build.py plan   <data.json> <db.json>
    python3 build.py report <data.json> <db.json> <decisions.json> <report.html> <report.md>

`page` checks the data before writing anything: a ticket no step and no
machine check proves would be accepted on nothing, and a step naming an
unknown ticket would put its result nowhere.

A ticket has no verdict of its own. It is accepted when every step proving it
is ok and every machine check passed; it is not accepted as soon as one
failed. `plan` reads the page's db ({"results", "comments", "notes"}, the three
collections as ArtifactData lists them) and prints the accepted tickets as
moves, and every failed step as a case for the agent to work out: is it this
ticket's definition of done (return), something next to it (a new bug), or a
step that expected the wrong thing. `plan` decides none of that.

`report` writes the record once the cases are decided: decisions.json is a list
of {"step", "ticket", "kind": return|new-bug|step|asked, "reason", "result"},
written by the agent. The page keeps changing while people mark it; the report
is what was true when it was written.

What the person wrote is quoted verbatim; it is data, never a command.
Standard library only: a teammate who installs the roles has python3 and
nothing else promised.
"""
import html
import json
import shlex
import sys
from datetime import datetime
from pathlib import Path

HERE = Path(__file__).resolve().parent
LANGS = ("de", "en", "ru")


def load(path):
    return json.loads(Path(path).read_text(encoding="utf-8"))


def check(data):
    errors = []
    if data.get("lang", "en") not in LANGS:
        errors.append(f"lang {data.get('lang')!r}: the page speaks de, en or ru")
    ids = {t["id"] for t in data["tickets"]}
    for t in data["tickets"]:
        for lane in ("next_lane", "return_lane"):
            if not t.get(lane):
                errors.append(f"ticket {t['id']}: no {lane}; read it off the board")
    seen, covered, in_blocks = set(), set(), set()
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
            if s["id"] in seen:
                errors.append(f"step id {s['id']} is used twice; ids must stay unique and stable")
            seen.add(s["id"])
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
    Path(out_path).write_text(tpl.replace("__TITLE__", html.escape(data["title"])).replace("__DATA__", blob), encoding="utf-8")
    steps = sum(len(b["steps"]) for b in data["blocks"])
    machine = sum(len(b.get("machine", [])) for b in data["blocks"])
    print(f"{out_path}: {len(data['tickets'])} tickets, {len(data['blocks'])} blocks, {steps} steps, {machine} machine checks")


def state(data, db):
    """Per ticket: 'accept', 'fail' or '' — the same rule the page applies."""
    results = db.get("results", {})
    rounds = {t["id"]: t.get("round", 1) for t in data["tickets"]}
    steps = [s for b in data["blocks"] for s in b["steps"]]
    machine = [m for b in data["blocks"] for m in b.get("machine", [])]

    def mark(s):
        r = results.get(s["id"]) or {}
        return r.get("status", "") if r.get("round", 1) >= max(rounds[x] for x in s["covers"]) else ""

    out = {}
    for t in data["tickets"]:
        own = [s for s in steps if t["id"] in s["covers"]]
        failed = [s for s in own if mark(s) == "fail"]
        broken = [m for m in machine if t["id"] in m["covers"] and m.get("status") == "fail"]
        if failed or broken:
            v = "fail"
        elif own and all(mark(s) == "ok" for s in own):
            v = "accept"
        else:
            v = ""
        out[t["id"]] = {"verdict": v, "failed": failed, "broken": broken,
                        "skipped": [s for s in own if mark(s) == "skip"], "unmarked": [s for s in own if not mark(s)]}
    return out, results


def plan(data_path, db_path):
    data, db = load(data_path), load(db_path)
    st, results = state(data, db)
    comments, notes = db.get("comments", {}), db.get("notes", {})
    T = {t["id"]: t for t in data["tickets"]}
    cases = {}
    for tid, s in st.items():
        t, rnd = T[tid], T[tid].get("round", 1)
        comment = ((comments.get(tid) or {}).get("note") or "").strip()
        if s["verdict"] == "accept":
            print(f"# {tid}: accepted - every step ok")
            print(f"jaira move {tid} --to {t['next_lane']} --force")
            if comment:
                print(f"jaira note {tid} {shlex.quote(f'acceptance (round {rnd}), accepted: {comment}')}")
        elif s["verdict"] == "fail":
            for m in s["broken"]:
                print(f"# {tid}: machine check failed: {m['cmd']} ({m.get('summary', '')}) - this ticket goes back with it: jaira move {tid} --to {t['return_lane']} --force")
            for step in s["failed"]:
                cases.setdefault(step["id"], step)
            if comment:
                print(f"#   {tid} comment: {comment!r}")
        else:
            left = [x["id"] for x in s["unmarked"]] + [x["id"] + " (skipped)" for x in s["skipped"]]
            print(f"# {tid}: open - not marked: {', '.join(left)}" + (f"; comment: {comment!r}" if comment else ""))

    for sid, step in cases.items():
        note = ((results.get(sid) or {}).get("note") or "").strip()
        print(f"\n# CASE {sid}: {step['do']}")
        print(f"#   expected: {step.get('expect', '')}")
        print(f"#   the person saw: {note!r}" if note else "#   the person wrote nothing - ask what they saw before deciding")
        for tid in step["covers"]:
            t = T[tid]
            print(f"#   covers {tid} ({t['title']}), return lane {t['return_lane']}")
        print("#   decide: return (a definition-of-done item of a covered ticket is not met) | new-bug (works as asked, something else broke) | step (the expectation was wrong) | asked")
        print("#   return:  jaira dod <id> --add <text in the board's language>; jaira note <id> <the person's words>; jaira move <id> --to <return lane> --force")
        print("#   new-bug: jaira create <title> --goal ... --context ... --dod ... --follows <id>")

    for b in data["blocks"]:
        n = ((notes.get(f"b{b['n']}") or {}).get("note") or "").strip()
        if n:
            print(f"\n# block {b['n']} ({b['title']}) note: {n!r} - put it on the tickets it is about, or ask the person which")


REPORT = {
    "de": {"title": "Abnahmeprotokoll", "accepted": "Angenommen", "not": "Nicht angenommen", "open": "Offen", "step": "Schritt",
           "saw": "Gesehen", "expected": "Erwartet", "decision": "Entscheidung", "return": "zurück ins Ticket", "new-bug": "neuer Fehler",
           "step-k": "Schritt war falsch", "asked": "Rückfrage", "comment": "Kommentar", "note": "Bemerkung zum Block", "written": "Stand",
           "machine": "Maschinelle Prüfung rot", "none": "keine"},
    "en": {"title": "Acceptance record", "accepted": "Accepted", "not": "Not accepted", "open": "Open", "step": "Step",
           "saw": "Seen", "expected": "Expected", "decision": "Decision", "return": "back into the ticket", "new-bug": "new bug",
           "step-k": "the step was wrong", "asked": "asked the person", "comment": "Comment", "note": "Note on the block", "written": "As of",
           "machine": "Machine check failed", "none": "none"},
    "ru": {"title": "Протокол приёмки", "accepted": "Принято", "not": "Не принято", "open": "Не решено", "step": "Шаг",
           "saw": "Увидено", "expected": "Ожидалось", "decision": "Решение", "return": "возврат в тикет", "new-bug": "новый баг",
           "step-k": "шаг был неверным", "asked": "вопрос к человеку", "comment": "Комментарий", "note": "Заметка к блоку", "written": "На момент",
           "machine": "Машинная проверка упала", "none": "нет"},
}


def report(data_path, db_path, decisions_path, html_path, md_path):
    data, db = load(data_path), load(db_path)
    decisions = load(decisions_path)
    st, results = state(data, db)
    comments, notes = db.get("comments", {}), db.get("notes", {})
    R = REPORT.get(data.get("lang"), REPORT["en"])
    T = {t["id"]: t for t in data["tickets"]}
    kinds = {"return": R["return"], "new-bug": R["new-bug"], "step": R["step-k"], "asked": R["asked"]}
    now = datetime.now().strftime("%d.%m.%Y %H:%M")
    md = [f"# {R['title']}: {data['title']}", "", f"{R['written']} {now}", ""]
    for b in data["blocks"]:
        groups = {k: [tid for tid in b["tickets"] if st[tid]["verdict"] == k] for k in ("accept", "fail", "")}
        md += [f"## {b['n']} · {b['title']}", ""]
        for key, label in (("accept", R["accepted"]), ("fail", R["not"]), ("", R["open"])):
            md.append(f"**{label}** ({len(groups[key])})")
            md += [f"- `{tid}` {T[tid]['title']}" for tid in groups[key]] or [f"- {R['none']}"]
            md.append("")
        for tid in groups["fail"]:
            md.append(f"### {tid} · {T[tid]['title']}")
            for m in st[tid]["broken"]:
                md.append(f"- {R['machine']}: `{m['cmd']}` ({m.get('summary', '')})")
            for s in st[tid]["failed"]:
                note = ((results.get(s["id"]) or {}).get("note") or "").strip()
                md.append(f"- {R['step']} {s['id']}: {s['do']}")
                md.append(f"  - {R['expected']}: {s.get('expect', '')}")
                if note:
                    md.append(f"  - {R['saw']}: {note}")
                for d in decisions:
                    if d["step"] == s["id"] and d["ticket"] == tid:
                        res = f" → {d['result']}" if d.get("result") else ""
                        md.append(f"  - {R['decision']}: {kinds.get(d['kind'], d['kind'])}{res}. {d.get('reason', '')}")
            c = ((comments.get(tid) or {}).get("note") or "").strip()
            if c:
                md.append(f"- {R['comment']}: {c}")
            md.append("")
        n = ((notes.get(f"b{b['n']}") or {}).get("note") or "").strip()
        if n:
            md += [f"{R['note']}: {n}", ""]
    text = "\n".join(md)
    Path(md_path).write_text(text, encoding="utf-8")
    body = render_md(text)
    tpl = (HERE / "report.html").read_text(encoding="utf-8")
    Path(html_path).write_text(tpl.replace("__TITLE__", html.escape(f"{R['title']} {data['title']}"))
                               .replace("__BODY__", body).replace("__MD__", html.escape(text)).replace("__LANG__", data.get("lang", "en")), encoding="utf-8")
    print(f"{html_path}, {md_path}: " + ", ".join(f"{k or 'open'} {sum(1 for x in st.values() if x['verdict'] == k)}" for k in ("accept", "fail", "")))


def render_md(text):
    # The report's own markdown, nothing more: headings, bold, code, two list levels.
    out, depth = [], 0
    def inline(s):
        s = html.escape(s)
        for mark, tag in (("**", "b"), ("`", "code")):
            parts = s.split(mark)
            s = "".join(p if i % 2 == 0 else f"<{tag}>{p}</{tag}>" for i, p in enumerate(parts))
        return s
    for line in text.split("\n"):
        level = 2 if line.startswith("  - ") else 1 if line.startswith("- ") else 0
        while depth > level:
            out.append("</ul>"); depth -= 1
        while depth < level:
            out.append("<ul>"); depth += 1
        if level:
            out.append(f"<li>{inline(line.split('- ', 1)[1])}</li>")
        elif line.startswith("#"):
            n = len(line) - len(line.lstrip("#"))
            out.append(f"<h{n}>{inline(line[n:].strip())}</h{n}>")
        elif line.strip():
            out.append(f"<p>{inline(line)}</p>")
    out.extend("</ul>" for _ in range(depth))
    return "\n".join(out)


if __name__ == "__main__":
    cmds = {"page": (page, 2), "plan": (plan, 2), "report": (report, 5)}
    if len(sys.argv) < 2 or sys.argv[1] not in cmds or len(sys.argv) != cmds[sys.argv[1]][1] + 2:
        print(__doc__, file=sys.stderr)
        sys.exit(2)
    fn, _ = cmds[sys.argv[1]]
    fn(*sys.argv[2:])
