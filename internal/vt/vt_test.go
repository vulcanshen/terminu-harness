package vt

// The first tests come from terminu-go's vt_test.go (the same commit as
// vt.go), without the ones on render.

import (
	"strings"
	"testing"
)

const glyphFolder = "\U000F024B" // nf-md-folder

func vtAfter(w, h int, out string) *vt {
	s := New(w, h, 1, 1)
	s.Write([]byte(out))
	return s.t
}

func TestVTPrintAndWrap(t *testing.T) {
	v := vtAfter(5, 3, "abcdefg\r\nxy")
	if got := strings.Join(v.text(), "|"); got != "abcde|fg|xy" {
		t.Errorf("wrap: %q", got)
	}
	if v.wraps != 1 {
		t.Errorf("wraps %d", v.wraps)
	}
	// Writing the last column holds the wrap until the next character.
	v = vtAfter(5, 2, "abcde\r\nz")
	if got := strings.Join(v.text(), "|"); got != "abcde|z" {
		t.Errorf("a full line then CR LF: %q", got)
	}
	if v.wraps != 0 {
		t.Errorf("a line exactly as wide wrapped %d", v.wraps)
	}
}

func TestVTWideCharacters(t *testing.T) {
	s := New(6, 2, 2, 1)
	s.Write([]byte("中文" + glyphFolder + "x"))
	v := s.t
	if got := strings.Join(v.text(), "|"); got != "中文"+glyphFolder+"|x" {
		t.Errorf("CJK and a two-cell icon: %q", got)
	}
	if v.x != 1 || v.y != 1 {
		t.Errorf("cursor at %d,%d", v.x, v.y)
	}
	// Combining marks join the character before.
	v = vtAfter(4, 1, "e\U00000301x")
	if got := v.text()[0]; got != "e\U00000301x" {
		t.Errorf("combining: %q", got)
	}
}

func TestVTIconWidths(t *testing.T) {
	for _, c := range []struct{ icon, pl, want int }{{1, 1, 3}, {2, 1, 4}, {2, 2, 5}} {
		s := New(10, 1, c.icon, c.pl)
		s.Write([]byte(glyphFolder + "x"))
		if s.t.x != c.want {
			t.Errorf("icon %d powerline %d: cursor at %d, want %d", c.icon, c.pl, s.t.x, c.want)
		}
	}
	// A reset keeps the widths.
	s := New(10, 1, 2, 2)
	s.Write([]byte("\x1bc" + glyphFolder))
	if s.t.x != 2 {
		t.Errorf("after RIS an icon took %d", s.t.x)
	}
}

func TestVTCursorAndErase(t *testing.T) {
	v := vtAfter(6, 3, "aaaaaa\r\nbbbbbb\r\ncccccc\x1b[2;3H\x1b[K\x1b[3;1H\x1b[1K\x1b[1;4HX")
	if got := strings.Join(v.text(), "|"); got != "aaaXaa|bb| ccccc" {
		t.Errorf("CUP, EL: %q", got)
	}
	v = vtAfter(4, 3, "1111\r\n2222\r\n3333\x1b[2;1H\x1b[J")
	if got := strings.Join(v.text(), "|"); got != "1111||" {
		t.Errorf("ED: %q", got)
	}
}

func TestVTScrollRegionAndHistory(t *testing.T) {
	v := vtAfter(3, 3, "a\r\nb\r\nc\r\nd")
	if got := strings.Join(v.text(), "|"); got != "b|c|d" || len(v.history) != 1 || v.scrolls != 1 {
		t.Errorf("scroll: %q, history %d, scrolls %d", got, len(v.history), v.scrolls)
	}
	// A region: only lines 2–3 move.
	v = vtAfter(3, 3, "a\r\nb\r\nc\x1b[2;3r\x1b[3;1H\nd")
	if got := strings.Join(v.text(), "|"); got != "a|c|d" {
		t.Errorf("region: %q", got)
	}
}

func TestVTAltScreen(t *testing.T) {
	v := vtAfter(4, 2, "main\x1b[?1049h\x1b[Halt")
	if v.text()[0] != "alt" {
		t.Errorf("alt screen: %q", v.text())
	}
	v.write([]byte("\x1b[?1049l"))
	if v.text()[0] != "main" {
		t.Errorf("back to main: %q", v.text())
	}
}

func TestVTAnswers(t *testing.T) {
	v := vtAfter(10, 5, "\x1b[3;4H\x1b[6n\x1b[c\x1b]4;1;?\x1b\\")
	r := string(v.reply)
	if !strings.HasPrefix(r, "\x1b[3;4R\x1b[?62;22c\x1b]4;1;rgb:") {
		t.Errorf("replies %q", r)
	}
}

func TestVTSplitSequences(t *testing.T) {
	v := New(10, 2, 1, 1).t
	for _, part := range []string{"\x1b", "[3", "1m", "re", "\xe4\xb8", "\xad"} {
		v.write([]byte(part))
	}
	if v.text()[0] != "re中" {
		t.Errorf("text %q", v.text()[0])
	}
	if c := v.cells[0][0]; c.attr.fg != (vtColor{kind: 1, idx: 1}) {
		t.Errorf("SGR split across writes: %+v", c.attr)
	}
}

func TestVTColours(t *testing.T) {
	s := New(10, 1, 1, 1)
	s.Write([]byte("\x1b[38;2;1;2;3;48;5;196mA\x1b[38:2::4:5:6mB\x1b[0;7mC"))
	row := s.Row(0)
	if row[0].Fg.Hex() != "#010203" || row[0].Bg.Hex() != "#ff0000" || row[1].Fg.Hex() != "#040506" {
		t.Errorf("colours %+v %+v", row[0], row[1])
	}
	if !row[2].Reverse || !row[2].FgDefault || !row[2].BgDefault {
		t.Errorf("reverse on the default colours: %+v", row[2])
	}
}

func TestVTResize(t *testing.T) {
	v := vtAfter(4, 4, "a\r\nb\r\nc\r\nd")
	v.resize(6, 2)
	if got := strings.Join(v.text(), "|"); got != "c|d" || v.y != 1 {
		t.Errorf("shrinking keeps the cursor's lines: %q at %d", got, v.y)
	}
}

func TestRowSum(t *testing.T) {
	s := New(4, 2, 1, 1)
	s.Write([]byte("ab"))
	before0, before1 := s.RowSum(0), s.RowSum(1)
	s.Write([]byte("\r\x1b[1mab"))
	if s.RowSum(0) == before0 {
		t.Error("a style change kept the sum")
	}
	if s.RowSum(1) != before1 {
		t.Error("an untouched row changed its sum")
	}
}
