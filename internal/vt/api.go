package vt

import (
	"fmt"
	"hash/fnv"

	"github.com/charmbracelet/x/ansi"
)

// The colours and widths vt.go needs come from terminu-go's color.go and
// width.go (the same commit, GPL-3.0), cut down to what a screen here uses.

// RGB is a 24-bit colour.
type RGB [3]uint8

// Hex is the colour as #rrggbb.
func (c RGB) Hex() string { return fmt.Sprintf("#%02x%02x%02x", c[0], c[1], c[2]) }

type color struct{ RGB RGB }

// Text and Base are what the terminal answers for its own foreground and
// background (OSC 10, 11), and what a cell with no colour of its own is drawn
// in: the family's (Catppuccin Mocha).
var (
	Text = color{RGB{0xcd, 0xd6, 0xf4}}
	Base = color{RGB{0x1e, 0x1e, 0x2e}}
)

var palette16 = [16]RGB{
	{0, 0, 0}, {205, 0, 0}, {0, 205, 0}, {205, 205, 0}, {0, 0, 238}, {205, 0, 205}, {0, 205, 205}, {229, 229, 229},
	{127, 127, 127}, {255, 0, 0}, {0, 255, 0}, {255, 255, 0}, {92, 92, 255}, {255, 0, 255}, {0, 255, 255}, {255, 255, 255},
}

// color256 is the RGB of a 256-colour index: 0–15 from the palette, the rest
// by the standard formula.
func color256(n int) RGB {
	switch {
	case n < 16:
		return palette16[max(n, 0)]
	case n < 232:
		n -= 16
		step := func(v int) uint8 {
			if v == 0 {
				return 0
			}
			return uint8(55 + v*40)
		}
		return RGB{step(n / 36), step(n / 6 % 6), step(n % 6)}
	default:
		v := uint8(8 + (min(n, 255)-232)*10)
		return RGB{v, v, v}
	}
}

// isPowerline reports whether r is a powerline symbol (U+E0A0–U+E0D7), which
// a terminal may draw narrower than the other icons (tdp E5).
func isPowerline(r rune) bool { return r >= 0xE0A0 && r <= 0xE0D7 }

// isIcon reports whether r is in a private use area, where Nerd Font glyphs live.
func isIcon(r rune) bool {
	return (r >= 0xE000 && r <= 0xF8FF) || (r >= 0xF0000 && r <= 0xFFFFD) || (r >= 0x100000 && r <= 0x10FFFD)
}

// runeWidth is how many cells this terminal advances for r.
func (t *vt) runeWidth(r rune) int {
	switch {
	case isPowerline(r):
		return t.powerlineCells
	case isIcon(r):
		return t.iconCells
	}
	return ansi.StringWidth(string(r))
}

// Screen is the terminal a program draws on: it reads what the program
// prints and keeps every cell.
type Screen struct{ t *vt }

// New is a blank screen of cols × rows that gives an icon iconCells cells and
// a powerline symbol powerlineCells (1 or 2).
func New(cols, rows, iconCells, powerlineCells int) *Screen {
	t := newVT(cols, rows)
	t.iconCells, t.powerlineCells = min(max(iconCells, 1), 2), min(max(powerlineCells, 1), 2)
	return &Screen{t}
}

// Write feeds what the program printed.
func (s *Screen) Write(b []byte) { s.t.write(b) }

// TakeReply returns what the terminal answers the program (cursor position,
// device attributes, colours) and forgets it.
func (s *Screen) TakeReply() []byte {
	r := s.t.reply
	s.t.reply = nil
	return r
}

// Resize changes the size, keeping what fits.
func (s *Screen) Resize(cols, rows int) { s.t.resize(cols, rows) }

// Size is the screen's cols and rows.
func (s *Screen) Size() (cols, rows int) { return s.t.w, s.t.h }

// Widths is how many cells an icon and a powerline symbol take.
func (s *Screen) Widths() (icon, powerline int) { return s.t.iconCells, s.t.powerlineCells }

// Cursor is where the cursor is (0-based) and whether the program shows it.
func (s *Screen) Cursor() (col, row int, visible bool) { return s.t.x, s.t.y, s.t.cursorOn }

// AltScreen reports whether the program is on the alternate screen.
func (s *Screen) AltScreen() bool { return s.t.altOn }

// AppCursor reports whether arrow keys are to be sent as ESC O A.
func (s *Screen) AppCursor() bool { return s.t.appCursor }

// Title is the window title the program set.
func (s *Screen) Title() string { return s.t.title }

// Counts is how many characters wrapped to the next line and how many lines
// scrolled off the top since the screen was made: a full-screen app that
// prints a line wider than the terminal wraps.
func (s *Screen) Counts() (wraps, scrolls int) { return s.t.wraps, s.t.scrolls }

// Text is what each row says, trailing blanks cut.
func (s *Screen) Text() []string { return s.t.text() }

// Cell is one cell of the screen.
type Cell struct {
	Text  string // the character; "" for a blank
	Width int    // 1, 2, or 0 for the right half of a wide character
	// Fg and Bg are the colours; FgDefault and BgDefault say the program left
	// them to the terminal (drawn in Text and Base).
	Fg, Bg                                                  RGB
	FgDefault, BgDefault                                    bool
	Bold, Faint, Italic, Underline, Reverse, Hidden, Strike bool
}

// Row is the cells of row y.
func (s *Screen) Row(y int) []Cell {
	line := s.t.cells[y]
	out := make([]Cell, len(line))
	for x, c := range line {
		a := c.attr
		out[x] = Cell{
			Text: c.s, Width: c.w,
			Fg: a.fg.rgbOr(Text.RGB), Bg: a.bg.rgbOr(Base.RGB),
			FgDefault: a.fg.kind == 0, BgDefault: a.bg.kind == 0,
			Bold: a.bold, Faint: a.faint, Italic: a.italic, Underline: a.underline,
			Reverse: a.reverse, Hidden: a.hidden, Strike: a.strike,
		}
	}
	return out
}

// RowSum is a checksum of row y, characters and attributes: it changes when
// anything on the row does.
func (s *Screen) RowSum(y int) uint64 {
	h := fnv.New64a()
	var b []byte
	for _, c := range s.t.cells[y] {
		a := c.attr
		var flags byte
		for i, on := range [...]bool{a.bold, a.faint, a.italic, a.underline, a.reverse, a.hidden, a.strike} {
			if on {
				flags |= 1 << i
			}
		}
		b = append(b[:0], c.s...)
		b = append(b, 0, byte(c.w), flags,
			a.fg.kind, a.fg.idx, a.fg.rgb[0], a.fg.rgb[1], a.fg.rgb[2],
			a.bg.kind, a.bg.idx, a.bg.rgb[0], a.bg.rgb[1], a.bg.rgb[2])
		h.Write(b)
	}
	return h.Sum64()
}
