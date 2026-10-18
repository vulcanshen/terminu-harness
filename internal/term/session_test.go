//go:build unix

package term

import (
	"strings"
	"testing"
	"time"

	"github.com/vulcanshen/terminu-harness/internal/vt"
)

func TestSessionTypesAndRecords(t *testing.T) {
	s, err := Start(Options{Argv: []string{"cat"}, Cols: 20, Rows: 4, IconCells: 1, PowerlineCells: 1})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	s.Record()
	if err := s.Type("abc"); err != nil {
		t.Fatal(err)
	}
	s.Settle(time.Now(), 100*time.Millisecond, 2*time.Second)
	frames := s.StopRecording()
	var text []string
	s.View(func(scr *vt.Screen) { text = scr.Text() })
	if text[0] != "abc" {
		t.Errorf("screen %q", text)
	}
	if len(frames) == 0 || frames[len(frames)-1].Changes[0].Text != "abc" {
		t.Errorf("frames %+v", frames)
	}
	if st := s.Close(); st == "" {
		t.Error("no exit status after close")
	}
}

func TestEnvironDropsTheRealTerminal(t *testing.T) {
	t.Setenv("FILU__ICON_WIDTH", "2")
	t.Setenv("TERMINU__POWERLINE_WIDTH", "2")
	t.Setenv("TERM", "tmux-256color")
	env := strings.Join(Environ([]string{"DEMOU__ICON_WIDTH=2"}), "\n")
	for _, gone := range []string{"FILU__ICON_WIDTH", "TERMINU__POWERLINE_WIDTH", "tmux-256color"} {
		if strings.Contains(env, gone) {
			t.Errorf("%s passed on", gone)
		}
	}
	for _, kept := range []string{"TERM=xterm-256color", "COLORTERM=truecolor", "DEMOU__ICON_WIDTH=2"} {
		if !strings.Contains(env, kept) {
			t.Errorf("%s missing", kept)
		}
	}
}
