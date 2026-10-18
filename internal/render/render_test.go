package render

import (
	"bytes"
	"image/png"
	"testing"

	"github.com/vulcanshen/terminu-harness/internal/vt"
)

func TestBoxArmsCoverTheBlock(t *testing.T) {
	for r := rune(0x2500); r <= 0x257F; r++ {
		_, ok := boxArms[r]
		diagonal := r >= 0x2571 && r <= 0x2573
		if ok == diagonal {
			t.Errorf("U+%04X: in the table %v", r, ok)
		}
	}
}

func TestPNG(t *testing.T) {
	s := vt.New(20, 3, 2, 1)
	s.Write([]byte("\x1b[1;38;2;137;180;250m┏━━┓\x1b[0m ab \U000F024B  中 █▌\r\n\x1b[7mrev\x1b[0m"))
	b, err := PNG(s)
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	fs, _ := loadFaces()
	if got := img.Bounds().Dx(); got != 20*fs.cellW {
		t.Errorf("width %d, cell %d", got, fs.cellW)
	}
	// The heavy line's middle is in the colour asked for.
	r, g, bl, _ := img.At(fs.cellW+fs.cellW/2, fs.cellH/2).RGBA()
	if r>>8 != 137 || g>>8 != 180 || bl>>8 != 250 {
		t.Errorf("━ drawn in %d,%d,%d", r>>8, g>>8, bl>>8)
	}
}
