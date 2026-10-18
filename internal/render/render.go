// Package render draws a screen as a PNG, with the fonts it carries: Go Mono
// for text and Symbols Nerd Font Mono for icons and powerline symbols, so the
// picture is the same on any machine, one without fonts too. Box drawing and
// block elements are drawn as shapes, as terminals do.
package render

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"unicode/utf8"

	"github.com/vulcanshen/terminu-harness/internal/vt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

//go:embed fonts/SymbolsNerdFontMono-Regular.ttf
var symbolsTTF []byte

const fontSize = 16

type face struct {
	f    *sfnt.Font
	face font.Face
	buf  sfnt.Buffer
}

func (f *face) has(r rune) bool {
	i, err := f.f.GlyphIndex(&f.buf, r)
	return err == nil && i != 0
}

type faces struct {
	regular, bold, symbols *face
	cellW, cellH, ascent   int
}

func load(ttf []byte) (*face, error) {
	f, err := opentype.Parse(ttf)
	if err != nil {
		return nil, err
	}
	fc, err := opentype.NewFace(f, &opentype.FaceOptions{Size: fontSize, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		return nil, err
	}
	return &face{f: f, face: fc}, nil
}

func loadFaces() (*faces, error) {
	reg, err := load(gomono.TTF)
	if err != nil {
		return nil, err
	}
	bold, err := load(gomonobold.TTF)
	if err != nil {
		return nil, err
	}
	sym, err := load(symbolsTTF)
	if err != nil {
		return nil, err
	}
	m := reg.face.Metrics()
	adv, _ := reg.face.GlyphAdvance('M')
	return &faces{
		regular: reg, bold: bold, symbols: sym,
		cellW:  adv.Ceil(),
		cellH:  (m.Ascent + m.Descent).Ceil() + 2,
		ascent: m.Ascent.Ceil() + 1,
	}, nil
}

// PNG draws the screen, the cursor too when the program shows it.
func PNG(scr *vt.Screen) ([]byte, error) {
	fs, err := loadFaces()
	if err != nil {
		return nil, err
	}
	cols, rows := scr.Size()
	cx, cy, cursorOn := scr.Cursor()
	img := image.NewRGBA(image.Rect(0, 0, cols*fs.cellW, rows*fs.cellH))
	draw.Draw(img, img.Bounds(), image.NewUniform(rgba(vt.Base.RGB)), image.Point{}, draw.Src)
	for y := range rows {
		row := scr.Row(y)
		for x, c := range row {
			if c.Width == 0 {
				continue
			}
			fg, bg := c.Fg, c.Bg
			if c.Reverse {
				fg, bg = bg, fg
			}
			if cursorOn && x == cx && y == cy {
				fg, bg = bg, fg
			}
			if c.Faint {
				fg = mix(fg, bg)
			}
			if c.Hidden {
				fg = bg
			}
			w := max(c.Width, 1)
			cell := image.Rect(x*fs.cellW, y*fs.cellH, (x+w)*fs.cellW, (y+1)*fs.cellH)
			draw.Draw(img, cell, image.NewUniform(rgba(bg)), image.Point{}, draw.Src)
			fs.glyph(img.SubImage(cell).(*image.RGBA), c, rgba(fg))
			if c.Underline {
				line(img, cell.Min.X, cell.Max.Y-2, cell.Max.X, cell.Max.Y-1, rgba(fg))
			}
			if c.Strike {
				mid := cell.Min.Y + fs.cellH/2
				line(img, cell.Min.X, mid, cell.Max.X, mid+1, rgba(fg))
			}
		}
	}
	var b bytes.Buffer
	if err := png.Encode(&b, img); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func rgba(c vt.RGB) color.RGBA { return color.RGBA{c[0], c[1], c[2], 255} }

func mix(a, b vt.RGB) vt.RGB {
	return vt.RGB{uint8((int(a[0]) + int(b[0])) / 2), uint8((int(a[1]) + int(b[1])) / 2), uint8((int(a[2]) + int(b[2])) / 2)}
}

func line(img *image.RGBA, x0, y0, x1, y1 int, c color.RGBA) {
	draw.Draw(img, image.Rect(x0, y0, x1, y1), image.NewUniform(c), image.Point{}, draw.Src)
}

// glyph draws one cell's character into dst, the cell's own rectangle.
func (fs *faces) glyph(dst *image.RGBA, c vt.Cell, fg color.RGBA) {
	if c.Text == "" || c.Text == " " {
		return
	}
	r, _ := utf8.DecodeRuneInString(c.Text)
	if boxDrawing(dst, r, fg) || blockElement(dst, r, fg) {
		return
	}
	f := fs.regular
	if c.Bold {
		f = fs.bold
	}
	switch {
	case f.has(r):
	case fs.regular.has(r):
		f = fs.regular
	case fs.symbols.has(r):
		f = fs.symbols
	default:
		// No font here has it: an outline where the character is.
		b := dst.Bounds().Inset(2)
		for _, e := range []image.Rectangle{
			image.Rect(b.Min.X, b.Min.Y, b.Max.X, b.Min.Y+1), image.Rect(b.Min.X, b.Max.Y-1, b.Max.X, b.Max.Y),
			image.Rect(b.Min.X, b.Min.Y, b.Min.X+1, b.Max.Y), image.Rect(b.Max.X-1, b.Min.Y, b.Max.X, b.Max.Y),
		} {
			draw.Draw(dst, e, image.NewUniform(fg), image.Point{}, draw.Src)
		}
		return
	}
	adv := font.MeasureString(f.face, c.Text)
	bounds := dst.Bounds()
	x := bounds.Min.X + (bounds.Dx()-adv.Ceil())/2
	d := font.Drawer{Dst: dst, Src: image.NewUniform(fg), Face: f.face, Dot: fixed.P(x, bounds.Min.Y+fs.ascent)}
	d.DrawString(c.Text)
}

// boxArms gives each box drawing character (U+2500–U+257F) the weight of
// its arm up, right, down and left: 0 none, 1 light, 2 heavy, 3 double.
// Dashes are drawn solid, arcs as corners; the diagonals are left to a font.
var boxArms = map[rune]string{
	'─': "0101", '━': "0202", '│': "1010", '┃': "2020", '┄': "0101", '┅': "0202", '┆': "1010", '┇': "2020",
	'┈': "0101", '┉': "0202", '┊': "1010", '┋': "2020", '┌': "0110", '┍': "0210", '┎': "0120", '┏': "0220",
	'┐': "0011", '┑': "0012", '┒': "0021", '┓': "0022", '└': "1100", '┕': "1200", '┖': "2100", '┗': "2200",
	'┘': "1001", '┙': "1002", '┚': "2001", '┛': "2002", '├': "1110", '┝': "1210", '┞': "2110", '┟': "1120",
	'┠': "2120", '┡': "2210", '┢': "1220", '┣': "2220", '┤': "1011", '┥': "1012", '┦': "2011", '┧': "1021",
	'┨': "2021", '┩': "2012", '┪': "1022", '┫': "2022", '┬': "0111", '┭': "0112", '┮': "0211", '┯': "0212",
	'┰': "0121", '┱': "0122", '┲': "0221", '┳': "0222", '┴': "1101", '┵': "1102", '┶': "1201", '┷': "1202",
	'┸': "2101", '┹': "2102", '┺': "2201", '┻': "2202", '┼': "1111", '┽': "1112", '┾': "1211", '┿': "1212",
	'╀': "2111", '╁': "1121", '╂': "2121", '╃': "2112", '╄': "2211", '╅': "1122", '╆': "1221", '╇': "2212",
	'╈': "1222", '╉': "2122", '╊': "2221", '╋': "2222", '╌': "0101", '╍': "0202", '╎': "1010", '╏': "2020",
	'═': "0303", '║': "3030", '╒': "0310", '╓': "0130", '╔': "0330", '╕': "0013", '╖': "0031", '╗': "0033",
	'╘': "1300", '╙': "3100", '╚': "3300", '╛': "1003", '╜': "3001", '╝': "3003", '╞': "1310", '╟': "3130",
	'╠': "3330", '╡': "1013", '╢': "3031", '╣': "3033", '╤': "0313", '╥': "0131", '╦': "0333", '╧': "1303",
	'╨': "3101", '╩': "3303", '╪': "1313", '╫': "3131", '╬': "3333", '╭': "0110", '╮': "0011", '╯': "1001",
	'╰': "1100", '╴': "0001", '╵': "1000", '╶': "0100", '╷': "0010", '╸': "0002", '╹': "2000", '╺': "0200",
	'╻': "0020", '╼': "0201", '╽': "1020", '╾': "0102", '╿': "2010",
}

func boxDrawing(dst *image.RGBA, r rune, fg color.RGBA) bool {
	arms, ok := boxArms[r]
	if !ok {
		return false
	}
	b := dst.Bounds()
	cx, cy := b.Min.X+b.Dx()/2, b.Min.Y+b.Dy()/2
	light := max(1, int(math.Round(float64(b.Dx())/9)))
	// bar draws one stroke of an arm th thick, off from the centre line; it
	// reaches past the centre by half a heavy stroke so the arms meet.
	bar := func(arm, off, th int) {
		reach := light
		var rect image.Rectangle
		switch arm {
		case 0:
			rect = image.Rect(cx-th/2+off, b.Min.Y, cx-th/2+off+th, cy+reach)
		case 1:
			rect = image.Rect(cx-reach, cy-th/2+off, b.Max.X, cy-th/2+off+th)
		case 2:
			rect = image.Rect(cx-th/2+off, cy-reach, cx-th/2+off+th, b.Max.Y)
		case 3:
			rect = image.Rect(b.Min.X, cy-th/2+off, cx+reach, cy-th/2+off+th)
		}
		draw.Draw(dst, rect, image.NewUniform(fg), image.Point{}, draw.Src)
	}
	for arm := range 4 {
		switch arms[arm] {
		case '1':
			bar(arm, 0, light)
		case '2':
			bar(arm, 0, light*2)
		case '3':
			bar(arm, -light, light)
			bar(arm, light, light)
		}
	}
	return true
}

func blockElement(dst *image.RGBA, r rune, fg color.RGBA) bool {
	if r < 0x2580 || r > 0x259F {
		return false
	}
	b := dst.Bounds()
	w, h := b.Dx(), b.Dy()
	fill := func(x0, y0, x1, y1 float64) {
		rect := image.Rect(b.Min.X+int(math.Round(x0*float64(w))), b.Min.Y+int(math.Round(y0*float64(h))),
			b.Min.X+int(math.Round(x1*float64(w))), b.Min.Y+int(math.Round(y1*float64(h))))
		draw.Draw(dst, rect, image.NewUniform(fg), image.Point{}, draw.Src)
	}
	shade := func(a uint8) {
		c := fg
		c.A = a
		draw.Draw(dst, b, image.NewUniform(c), image.Point{}, draw.Over)
	}
	n := float64(r - 0x2580)
	switch {
	case r == 0x2580:
		fill(0, 0, 1, 0.5)
	case r >= 0x2581 && r <= 0x2588: // lower eighths up to the full block
		fill(0, 1-n/8, 1, 1)
	case r >= 0x2589 && r <= 0x258F: // left eighths, seven down to one
		fill(0, 0, (8-(n-8))/8, 1)
	case r == 0x2590:
		fill(0.5, 0, 1, 1)
	case r == 0x2591:
		shade(64)
	case r == 0x2592:
		shade(128)
	case r == 0x2593:
		shade(192)
	case r == 0x2594:
		fill(0, 0, 1, 0.125)
	case r == 0x2595:
		fill(0.875, 0, 1, 1)
	default: // quadrants
		quads := map[rune]string{0x2596: "0010", 0x2597: "0001", 0x2598: "1000", 0x2599: "1011", 0x259A: "1001",
			0x259B: "1110", 0x259C: "1101", 0x259D: "0100", 0x259E: "0110", 0x259F: "0111"}
		q := quads[r]
		at := [4][2]float64{{0, 0}, {0.5, 0}, {0, 0.5}, {0.5, 0.5}}
		for i := range 4 {
			if q[i] == '1' {
				fill(at[i][0], at[i][1], at[i][0]+0.5, at[i][1]+0.5)
			}
		}
	}
	return true
}
