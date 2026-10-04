#!/usr/bin/env python3
"""Logbook summary page: build it from summary.json, after checking the text has no code in it.

    python3 build.py page <summary.json> <out.html> <out.md>

The reader is a project lead. A file path, a commit hash or a function name in
the text means it was written for a developer; `page` refuses it and names the
line, so the text is rewritten rather than published. The internal ids under
`appendix` are exempt — they are there on purpose.

Standard library only: the roles promise python3 and nothing else.
"""
import html
import json
import re
import sys
from pathlib import Path

HERE = Path(__file__).resolve().parent
CODE = [
    (re.compile(r"\b[\w-]+/[\w./-]+\.\w{1,5}\b"), "a file path"),
    (re.compile(r"\b[\w-]+\.(?:go|py|ts|js|svelte|md|yml|yaml|json|sql|rs)\b"), "a file name"),
    (re.compile(r"\b[0-9a-f]{7,40}\b"), "a commit hash"),
    (re.compile(r"\b[a-z]+_[a-z_]+\b"), "a snake_case name"),
    (re.compile(r"\b[a-z]+[A-Z]\w*\("), "a function call"),
    (re.compile(r"`"), "code formatting"),
]


def texts(node, where="page"):
    # Every string the reader sees, wherever it sits. The appendix holds the
    # internal ids on purpose, the hours cells are figures, and lang is a code.
    if isinstance(node, str):
        yield where, node
    elif isinstance(node, list):
        for i, x in enumerate(node):
            yield from texts(x, f"{where}[{i}]")
    elif isinstance(node, dict):
        for k, v in node.items():
            if k in ("appendix", "lang", "status") or (k == "rows" and where.endswith("hours")):
                continue
            yield from texts(v, f"{where}.{k}")


def check(data):
    errors = []
    for where, text in texts(data):
        for rx, what in CODE:
            m = rx.search(text or "")
            if m:
                errors.append(f"{where}: {what} ({m.group(0)!r}) - a project lead reads this; say what it means instead")
    return errors


def to_md(data):
    out = [f"# {data['title']}", ""]
    if data.get("subtitle"):
        out += [data["subtitle"], ""]
    for f in data.get("figures", []):
        out.append(f"- **{f['label']}:** {f['value']}" + (f" ({f['note']})" if f.get("note") else ""))
    out.append("")
    h = data.get("hours")
    if h and h.get("rows"):
        out += [f"## {h.get('title', 'Hours')}", "", "| " + " | ".join(h["columns"]) + " |", "|" + "---|" * len(h["columns"])]
        out += ["| " + " | ".join(str(c) for c in r["cells"]) + " |" for r in h["rows"]]
        if h.get("note"):
            out += ["", h["note"]]
        out.append("")
    for s in data.get("sections", []):
        out += [f"## {s['title']}", ""]
        if s.get("intro"):
            out += [s["intro"], ""]
        for it in s.get("items", []):
            out += [f"### {it['title']}" + (f" ({it['status_label']})" if it.get("status_label") else ""), "", it.get("copy") or it.get("body", ""), ""]
            for c in it.get("comments", []):
                out += [f"**{c['date']}**", "", c["text"], ""]
    if data.get("appendix"):
        out += [f"## {data['appendix']['title']}", ""] + [f"- {x}" for x in data["appendix"]["lines"]]
    return "\n".join(out).rstrip() + "\n"


def page(src, out_html, out_md):
    data = json.loads(Path(src).read_text(encoding="utf-8"))
    errors = check(data)
    if errors:
        print("\n".join(errors), file=sys.stderr)
        sys.exit(1)
    md = to_md(data)
    data["full_copy"] = md
    blob = json.dumps(data, ensure_ascii=False).replace("</", "<\\/")
    tpl = (HERE / "template.html").read_text(encoding="utf-8")
    Path(out_html).write_text(tpl.replace("__TITLE__", html.escape(data["title"])).replace("__DATA__", blob), encoding="utf-8")
    Path(out_md).write_text(md, encoding="utf-8")
    print(f"{out_html}, {out_md}: {sum(len(s.get('items', [])) for s in data.get('sections', []))} items")


if __name__ == "__main__":
    if len(sys.argv) == 5 and sys.argv[1] == "page":
        page(*sys.argv[2:])
    else:
        print(__doc__, file=sys.stderr)
        sys.exit(2)
