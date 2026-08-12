// Command iconpreview renders icon/jAIra.png in several terminal styles so a
// style can be chosen by looking at it rather than by describing it.
//
// This is a scratch tool, deliberately outside the jaira binary: nothing here
// ships until a style is picked.
//
// Run it from the repository root:
//
//	go run ./scripts/iconpreview
package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
)

// The source has no alpha channel — it is a photograph-style export with a dark
// gradient baked in — so "no background" has to be decided by brightness rather
// than by transparency. The glyph is near-white and the ground is around 12%
// luminance, so the gap is wide and the threshold is not delicate.
const keyThreshold = 0.28

func luminance(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	return (0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)) / 65535.0
}

// sample maps a cell of the output grid back to a pixel of the source.
func sample(img image.Image, gx, gy, gw, gh int) color.Color {
	b := img.Bounds()
	x := b.Min.X + gx*b.Dx()/gw
	y := b.Min.Y + gy*b.Dy()/gh
	if x >= b.Max.X {
		x = b.Max.X - 1
	}
	if y >= b.Max.Y {
		y = b.Max.Y - 1
	}
	return img.At(x, y)
}

func rgb(c color.Color) (int, int, int) {
	r, g, b, _ := c.RGBA()
	return int(r >> 8), int(g >> 8), int(b >> 8)
}

// halfBlocks renders two pixel rows per text row using the upper-half block, so
// a cell carries two independent colours and the result is roughly square.
func halfBlocks(img image.Image, cols int, key bool) string {
	b := img.Bounds()
	rows := cols * b.Dy() / (2 * b.Dx())
	gw, gh := cols, rows*2
	var out string
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			top := sample(img, c, r*2, gw, gh)
			bot := sample(img, c, r*2+1, gw, gh)
			tKey := key && luminance(top) < keyThreshold
			bKey := key && luminance(bot) < keyThreshold
			tr, tg, tb := rgb(top)
			br, bg, bb := rgb(bot)
			switch {
			case tKey && bKey:
				out += " "
			case tKey:
				// Only the lower pixel is opaque: draw it as a lower-half block so
				// the terminal's own background shows through the top.
				out += fmt.Sprintf("\x1b[38;2;%d;%d;%dm▄\x1b[0m", br, bg, bb)
			case bKey:
				out += fmt.Sprintf("\x1b[38;2;%d;%d;%dm▀\x1b[0m", tr, tg, tb)
			default:
				out += fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀\x1b[0m",
					tr, tg, tb, br, bg, bb)
			}
		}
		out += "\n"
	}
	return out
}

// silhouette drops colour entirely and paints the glyph in one accent, which is
// the only style that stays legible in a terminal without truecolour.
func silhouette(img image.Image, cols int, ansi string) string {
	b := img.Bounds()
	rows := cols * b.Dy() / (2 * b.Dx())
	var out string
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			top := luminance(sample(img, c, r*2, cols, rows*2)) >= keyThreshold
			bot := luminance(sample(img, c, r*2+1, cols, rows*2)) >= keyThreshold
			ch := " "
			switch {
			case top && bot:
				ch = "█"
			case top:
				ch = "▀"
			case bot:
				ch = "▄"
			}
			if ch == " " {
				out += " "
			} else {
				out += ansi + ch + "\x1b[0m"
			}
		}
		out += "\n"
	}
	return out
}

// ascii is the last resort: no colour, no block glyphs, works over any link.
func asciiArt(img image.Image, cols int) string {
	ramp := []rune(" .:-=+*#%@")
	b := img.Bounds()
	rows := cols * b.Dy() / (2 * b.Dx())
	var out string
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			l := luminance(sample(img, c, r, cols, rows))
			i := int(l * float64(len(ramp)-1))
			if i < 0 {
				i = 0
			}
			if i >= len(ramp) {
				i = len(ramp) - 1
			}
			out += string(ramp[i])
		}
		out += "\n"
	}
	return out
}

// colourASCII is the ramp style with each character tinted by its own pixel.
// The shape comes from the character, the shading from the colour, which reads
// better than either alone at small sizes.
func colourASCII(img image.Image, cols int, key bool) string {
	ramp := []rune(" .:-=+*#%@")
	b := img.Bounds()
	rows := cols * b.Dy() / (2 * b.Dx())
	var out string
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			px := sample(img, c, r, cols, rows)
			l := luminance(px)
			if key && l < keyThreshold {
				out += " "
				continue
			}
			i := int(l * float64(len(ramp)-1))
			if i < 0 {
				i = 0
			}
			if i >= len(ramp) {
				i = len(ramp) - 1
			}
			cr, cg, cb := rgb(px)
			out += fmt.Sprintf("\x1b[38;2;%d;%d;%dm%c\x1b[0m", cr, cg, cb, ramp[i])
		}
		out += "\n"
	}
	return out
}

// crescentOnly masks the image down to the moon.
//
// The crescent, the prompt glyphs and each ray are separate shapes, but at the
// keying threshold they are bridged by the glow around them and merge into one
// blob. Separating them needs a higher threshold — measured at 0.55, where the
// crescent, ">", "_" and the six rays come apart into distinct components — and
// the largest of those is the moon. The mask is then dilated and intersected
// back with the ordinary threshold so the soft antialiased edge survives, which
// a hard 0.55 cut would leave jagged.
func crescentOnly(img image.Image) image.Image {
	const separate = 0.55
	const dilate = 2

	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	hot := make([]bool, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			hot[y*w+x] = luminance(img.At(b.Min.X+x, b.Min.Y+y)) >= separate
		}
	}

	// Flood fill each component, remembering the largest.
	seen := make([]bool, w*h)
	var best []int
	for i := range hot {
		if !hot[i] || seen[i] {
			continue
		}
		var comp []int
		stack := []int{i}
		seen[i] = true
		for len(stack) > 0 {
			p := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			comp = append(comp, p)
			px, py := p%w, p/w
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					nx, ny := px+dx, py+dy
					if nx < 0 || ny < 0 || nx >= w || ny >= h {
						continue
					}
					q := ny*w + nx
					if hot[q] && !seen[q] {
						seen[q] = true
						stack = append(stack, q)
					}
				}
			}
		}
		if len(comp) > len(best) {
			best = comp
		}
	}

	keep := make([]bool, w*h)
	for _, p := range best {
		px, py := p%w, p/w
		for dy := -dilate; dy <= dilate; dy++ {
			for dx := -dilate; dx <= dilate; dx++ {
				nx, ny := px+dx, py+dy
				if nx < 0 || ny < 0 || nx >= w || ny >= h {
					continue
				}
				keep[ny*w+nx] = true
			}
		}
	}

	out := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			src := img.At(b.Min.X+x, b.Min.Y+y)
			if keep[y*w+x] && luminance(src) >= keyThreshold {
				out.Set(x, y, src)
			} else {
				// Pure black reads as background to every style below.
				out.Set(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}
	return out
}

func main() {
	path := "icon/jAIra.png"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "iconpreview: %v\n(run this from the repository root)\n", err)
		os.Exit(1)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "iconpreview: %v\n", err)
		os.Exit(1)
	}

	show := func(n int, title, desc, art string) {
		fmt.Printf("\n\x1b[1m── %d. %s\x1b[0m\n\x1b[2m%s\x1b[0m\n\n%s", n, title, desc, art)
	}

	show(1, "Half-blocks, background keyed out, 36 cols",
		"Full colour, your terminal background shows through. The 'transparent' look.",
		halfBlocks(img, 36, true))

	show(2, "Half-blocks, background kept, 36 cols",
		"The image exactly as exported, gradient and all — a solid tile on the screen.",
		halfBlocks(img, 36, false))

	show(3, "Half-blocks, keyed out, 20 cols",
		"Same as 1, small enough to sit beside a project list rather than above it.",
		halfBlocks(img, 20, true))

	show(4, "Silhouette in the board's accent colour, 28 cols",
		"One colour, no gradient. The only style that survives a terminal without truecolour.",
		silhouette(img, 28, "\x1b[38;5;39m"))

	show(5, "Silhouette in white, 28 cols",
		"As 4, taking the terminal's own foreground instead of an accent.",
		silhouette(img, 28, "\x1b[97m"))

	show(6, "ASCII ramp, 44 cols",
		"No colour and no block glyphs. Works anywhere, including a plain log.",
		asciiArt(img, 44))

	moon := crescentOnly(img)

	show(7, "Colour ASCII ramp, background kept, 60 cols",
		"Style 6 with each character tinted by its own pixel — shape from the glyph, shading from the colour.",
		colourASCII(img, 60, false))

	show(8, "Colour ASCII ramp, background keyed out, 60 cols",
		"As 7, with the gradient dropped so the terminal shows through.",
		colourASCII(img, 60, true))

	show(9, "Half-blocks, background kept, 60 cols",
		"Style 2, larger.",
		halfBlocks(img, 60, false))

	show(10, "Half-blocks, background kept, 80 cols",
		"Style 2, larger still. Needs a wide terminal.",
		halfBlocks(img, 80, false))

	show(11, "Crescent only, keyed out, 36 cols",
		"The moon by itself: prompt glyphs and rays removed by component analysis.",
		halfBlocks(moon, 36, true))

	show(12, "Crescent only, silhouette in accent, 28 cols",
		"The moon by itself, one colour, no gradient.",
		silhouette(moon, 28, "\x1b[38;5;39m"))

	fmt.Printf("\n\x1b[2mSource: %s\x1b[0m\n", path)
}
