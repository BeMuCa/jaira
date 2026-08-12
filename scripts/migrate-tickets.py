#!/usr/bin/env python3
"""Migrate requirementsgenie tickets (copies only) into jaira's format.

Adds `id:` and `status: backlog` to the frontmatter, leaves everything else
byte-for-byte untouched (including the German body and the Jira-style keys),
and renames the file to `<ulid>-<slug>.md`. Relies on jaira's own H1 fallback
for the title rather than duplicating it into frontmatter.
"""
import os
import re
import sys
import time
import random

SRC = "/tmp/jtest-sync/tickets-src"
DST = "/tmp/jtest-sync/repoG/.jaira/tickets"

CROCKFORD = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

def encode_b32(value, length):
    chars = []
    for _ in range(length):
        chars.append(CROCKFORD[value & 0x1F])
        value >>= 5
    return "".join(reversed(chars))

_last_ms = int(time.time() * 1000)

def new_ulid():
    global _last_ms
    _last_ms += 1  # keep them strictly increasing / unique across the batch
    ts_part = encode_b32(_last_ms, 10)
    rand_part = encode_b32(random.getrandbits(80), 16)
    return ts_part + rand_part

def slugify(title):
    out = []
    last_dash = True
    for ch in title.lower():
        if ch.isalnum():
            out.append(ch)
            last_dash = False
        else:
            if not last_dash:
                out.append("-")
                last_dash = True
    s = "".join(out).strip("-")
    if len(s) > 48:
        s = s[:48].strip("-")
    return s or "untitled"

def h1_title(body):
    for line in body.splitlines():
        if line.startswith("# "):
            return line[2:].strip()
    return ""

def migrate_one(path):
    with open(path, "r", encoding="utf-8") as f:
        raw = f.read()

    lines = raw.split("\n")
    assert lines[0] == "---", f"{path}: does not start with ---"
    end = None
    for i in range(1, len(lines)):
        if lines[i] == "---":
            end = i
            break
    assert end is not None, f"{path}: no closing ---"

    fm_lines = lines[1:end]
    body = "\n".join(lines[end + 1:])

    tid = new_ulid()
    # Prepend id/status as the first two frontmatter lines; nothing else moves.
    new_fm = [f"id: {tid}", "status: backlog"] + fm_lines

    new_raw = "---\n" + "\n".join(new_fm) + "\n---\n" + body

    title = h1_title(body)
    slug = slugify(title)
    new_name = f"{tid}-{slug}.md"
    new_path = os.path.join(DST, new_name)

    with open(new_path, "w", encoding="utf-8") as f:
        f.write(new_raw)
    return path, new_path, title


def main():
    os.makedirs(DST, exist_ok=True)
    count = 0
    for name in sorted(os.listdir(SRC)):
        if name == "TEMPLATE.md":
            continue
        if not name.endswith(".md"):
            continue
        src_path = os.path.join(SRC, name)
        old, new, title = migrate_one(src_path)
        print(f"{os.path.basename(old)} -> {os.path.basename(new)}  [{title}]")
        count += 1
    print(f"migrated {count} tickets")


if __name__ == "__main__":
    main()
