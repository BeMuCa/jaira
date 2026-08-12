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

USAGE = """usage: migrate-tickets.py <src-tickets-dir> <dst-tickets-dir>

Reads every *.md except TEMPLATE.md from src, writes a migrated copy into dst.
Never modifies src. dst is created if absent.

  python3 scripts/migrate-tickets.py ~/git/PROJECT/tickets /tmp/mig/.jaira/tickets
"""

CROCKFORD = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

def encode_b32(value, length):
    chars = []
    for _ in range(length):
        chars.append(CROCKFORD[value & 0x1F])
        value >>= 5
    return "".join(reversed(chars))

_used_ms = set()

def new_ulid(when_ms):
    """A ULID whose timestamp reflects when the ticket was actually written.

    ULIDs sort lexicographically by their timestamp, so using each file's mtime
    makes `ls` show tickets in authoring order rather than in whatever order this
    script happened to process them. Collisions within a millisecond are nudged
    forward so every id stays unique.
    """
    while when_ms in _used_ms:
        when_ms += 1
    _used_ms.add(when_ms)
    return encode_b32(when_ms, 10) + encode_b32(random.getrandbits(80), 16)

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

def migrate_one(path, dst):
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

    tid = new_ulid(int(os.path.getmtime(path) * 1000))
    # Prepend id/status as the first two frontmatter lines; nothing else moves.
    new_fm = [f"id: {tid}", "status: backlog"] + fm_lines

    new_raw = "---\n" + "\n".join(new_fm) + "\n---\n" + body

    title = h1_title(body)
    slug = slugify(title)
    new_name = f"{tid}-{slug}.md"
    new_path = os.path.join(dst, new_name)

    with open(new_path, "w", encoding="utf-8") as f:
        f.write(new_raw)
    return path, new_path, title


def main():
    if len(sys.argv) != 3:
        print(USAGE, file=sys.stderr)
        return 2
    src, dst = sys.argv[1], sys.argv[2]
    if os.path.abspath(src) == os.path.abspath(dst):
        print("refusing to migrate a directory onto itself", file=sys.stderr)
        return 2
    os.makedirs(dst, exist_ok=True)
    count = 0
    for name in sorted(os.listdir(src)):
        if name == "TEMPLATE.md" or not name.endswith(".md"):
            continue
        old, new, title = migrate_one(os.path.join(src, name), dst)
        print(f"{os.path.basename(old)} -> {os.path.basename(new)}  [{title}]")
        count += 1
    print(f"migrated {count} tickets into {dst}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
