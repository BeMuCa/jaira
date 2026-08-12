#!/usr/bin/env python3
"""Render a terminal screen — ANSI colours and all — to a PNG.

Screenshots of a TUI are otherwise impossible to produce in a headless
environment, and a plain text block loses exactly the thing worth showing: the
colour that tells you which board is waiting on you.

    python3 scripts/termshot.py in.ans out.png [--cols 100]

Reads ANSI text on a pseudo-terminal capture, keeps SGR colour, discards cursor
movement, and draws the result with a monospace font.
"""
import re
import sys
import unicodedata

from PIL import Image, ImageDraw, ImageFont

FONT = "/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf"
SIZE = 17
PAD = 18
BG = (13, 17, 23)      # a dark terminal, which is what the palette is tuned for
FG = (201, 209, 217)

# xterm 256-colour cube, enough for the codes jaira actually emits.
CUBE = [0, 95, 135, 175, 215, 255]
BASE16 = [
    (0, 0, 0), (205, 49, 49), (13, 188, 121), (229, 229, 16),
    (36, 114, 200), (188, 63, 188), (17, 168, 205), (229, 229, 229),
    (102, 102, 102), (241, 76, 76), (35, 209, 139), (245, 245, 67),
    (59, 142, 234), (214, 112, 214), (41, 184, 219), (255, 255, 255),
]


def xterm(n):
    if n < 16:
        return BASE16[n]
    if n < 232:
        n -= 16
        return (CUBE[n // 36 % 6], CUBE[n // 6 % 6], CUBE[n % 6])
    v = 8 + (n - 232) * 10
    return (v, v, v)


SGR = re.compile(r"\x1b\[([0-9;]*)m")
OTHER_CSI = re.compile(r"\x1b\[[0-9;?]*[A-Za-z]")
OSC = re.compile(r"\x1b\][^\x07\x1b]*(\x07|\x1b\\)")


def cells(text):
    """Yield (char, fg, bg, bold) for one screen."""
    fg, bg, bold = FG, None, False
    text = OSC.sub("", text)
    i = 0
    while i < len(text):
        m = SGR.match(text, i)
        if m:
            for part in (m.group(1) or "0").split(";"):
                p = part or "0"
                if p == "0":
                    fg, bg, bold = FG, None, False
                elif p == "1":
                    bold = True
                elif p == "22":
                    bold = False
            # 38;5;N and 48;5;N need the whole sequence, handled here.
            codes = (m.group(1) or "").split(";")
            for j, c in enumerate(codes):
                if c in ("38", "48") and j + 2 < len(codes) and codes[j + 1] == "5":
                    col = xterm(int(codes[j + 2]))
                    if c == "38":
                        fg = col
                    else:
                        bg = col
                if c in ("38", "48") and j + 4 < len(codes) and codes[j + 1] == "2":
                    col = tuple(int(x) for x in codes[j + 2:j + 5])
                    if c == "38":
                        fg = col
                    else:
                        bg = col
            i = m.end()
            continue
        m = OTHER_CSI.match(text, i)
        if m:
            i = m.end()
            continue
        ch = text[i]
        i += 1
        if ch == "\x1b" or ch == "\r":
            continue
        yield ch, fg, bg, bold


def render(text, out, cols=100):
    font = ImageFont.truetype(FONT, SIZE)
    bold_font = ImageFont.truetype(FONT, SIZE)
    cw = font.getlength("M")
    ch = SIZE + 5

    lines = [[]]
    for c, fg, bg, bold in cells(text):
        if c == "\n":
            lines.append([])
            continue
        lines[-1].append((c, fg, bg, bold))
    while lines and not lines[-1]:
        lines.pop()

    width = int(cw * cols) + PAD * 2
    height = ch * len(lines) + PAD * 2
    img = Image.new("RGB", (width, height), BG)
    d = ImageDraw.Draw(img)

    for row, line in enumerate(lines):
        x = PAD
        y = PAD + row * ch
        for c, fg, bg, bold in line:
            w = cw * (2 if unicodedata.east_asian_width(c) in ("W", "F") else 1)
            if bg:
                d.rectangle([x, y, x + w, y + ch], fill=bg)
            if c != " ":
                d.text((x, y), c, font=bold_font if bold else font, fill=fg)
            x += w
    img.save(out)
    print(f"{out}  {width}x{height}  {len(lines)} rows")


if __name__ == "__main__":
    cols = 100
    if "--cols" in sys.argv:
        cols = int(sys.argv[sys.argv.index("--cols") + 1])
    with open(sys.argv[1], "r", errors="replace") as f:
        render(f.read(), sys.argv[2], cols)
