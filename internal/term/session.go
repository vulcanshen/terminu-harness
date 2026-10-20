//go:build unix

// Package term runs a program in a PTY on a screen of the harness's own, so
// the agent can press keys and read what it drew.
package term

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/creack/pty"
	"github.com/vulcanshen/terminu-harness/internal/vt"
)

// Options says what to run and on what kind of terminal.
type Options struct {
	Argv                      []string
	Dir                       string
	Env                       []string // KEY=VALUE, over the harness's own environment
	Cols, Rows                int
	IconCells, PowerlineCells int // 1 or 2
}

// Change is one row that changed, as it reads after the change.
type Change struct {
	Row       int
	Text      string
	StyleOnly bool // the characters stayed, colours or attributes changed
}

// Frame is the screen changing once, At after the recording started.
type Frame struct {
	At      time.Duration
	Changes []Change
}

// Session is one program running on one screen.
type Session struct {
	mu      sync.Mutex
	scr     *vt.Screen
	cmd     *exec.Cmd
	pty     *os.File
	text    []string
	sums    []uint64
	changed time.Time // when the screen last changed
	rec     *recording
	done    chan struct{}
	exitErr error
}

type recording struct {
	start  time.Time
	frames []Frame
}

// mergeWithin is how close two changes come to be one frame (under a frame
// at 120 fps, the most Bubble Tea draws): a program's
// frame can reach the PTY in several reads.
const mergeWithin = 6 * time.Millisecond

// Start runs the program.
func Start(o Options) (*Session, error) {
	if len(o.Argv) == 0 {
		return nil, errors.New("no command")
	}
	cmd := exec.Command(o.Argv[0], o.Argv[1:]...)
	cmd.Dir, cmd.Env = o.Dir, Environ(o.Env)
	f, err := pty.StartWithSize(cmd, &pty.Winsize{Cols: uint16(o.Cols), Rows: uint16(o.Rows)})
	if err != nil {
		return nil, err
	}
	s := &Session{
		scr: vt.New(o.Cols, o.Rows, o.IconCells, o.PowerlineCells),
		cmd: cmd, pty: f, done: make(chan struct{}),
	}
	s.snapshot()
	in := make(chan chunk, 4096)
	go s.read(in)
	go s.process(in)
	go func() {
		s.exitErr = cmd.Wait()
		close(s.done)
	}()
	return s, nil
}

// Environ is the harness's environment without what describes the terminal
// it runs in, with the harness terminal and env over it. Widths of icons set
// for the real terminal (<APP>__ICON_WIDTH, TERMINU__ICON_WIDTH, the powerline
// ones) are left out, so the program measures this one, unless env sets them.
func Environ(env []string) []string {
	var out []string
	for _, kv := range os.Environ() {
		k, _, _ := strings.Cut(kv, "=")
		switch {
		case k == "TERM", k == "COLORTERM", k == "LINES", k == "COLUMNS", k == "TMUX", k == "TMUX_PANE",
			strings.HasSuffix(k, "__ICON_WIDTH"), strings.HasSuffix(k, "__POWERLINE_WIDTH"):
			continue
		}
		out = append(out, kv)
	}
	out = append(out, "TERM=xterm-256color", "COLORTERM=truecolor")
	return append(out, env...)
}

type chunk struct {
	b  []byte
	at time.Time
}

// read takes what the program prints as fast as it comes, so the program
// never waits on a full PTY while the screen is worked out; process does that.
func (s *Session) read(in chan<- chunk) {
	defer close(in)
	buf := make([]byte, 64*1024)
	for {
		n, err := s.pty.Read(buf)
		if n > 0 {
			in <- chunk{append([]byte(nil), buf[:n]...), time.Now()}
		}
		if err != nil {
			return
		}
	}
}

// process feeds the screen everything read so far at once: a frame the
// program wrote in one go becomes one change.
func (s *Session) process(in <-chan chunk) {
	for c := range in {
		batch := []chunk{c}
	more:
		for {
			select {
			case c, ok := <-in:
				if !ok {
					break more
				}
				batch = append(batch, c)
			default:
				break more
			}
		}
		s.feed(batch)
	}
}

func (s *Session) feed(batch []chunk) {
	s.mu.Lock()
	for _, c := range batch {
		s.scr.Write(c.b)
	}
	reply := s.scr.TakeReply()
	now := batch[0].at
	if ch := s.snapshot(); len(ch) > 0 {
		s.changed = now
		if r := s.rec; r != nil {
			at := now.Sub(r.start)
			if n := len(r.frames); n > 0 && at-r.frames[n-1].At < mergeWithin {
				r.frames[n-1].Changes = merge(r.frames[n-1].Changes, ch)
			} else {
				r.frames = append(r.frames, Frame{At: at, Changes: ch})
			}
		}
	}
	s.mu.Unlock()
	if len(reply) > 0 {
		_, _ = s.pty.Write(reply)
	}
}

// snapshot takes the rows' text and sums again and returns the rows that changed.
func (s *Session) snapshot() []Change {
	text := s.scr.Text()
	_, rows := s.scr.Size()
	var ch []Change
	sums := make([]uint64, rows)
	for y := range rows {
		sums[y] = s.scr.RowSum(y)
		if y < len(s.sums) && sums[y] == s.sums[y] {
			continue
		}
		styleOnly := y < len(s.text) && s.text[y] == text[y]
		ch = append(ch, Change{Row: y, Text: text[y], StyleOnly: styleOnly})
	}
	s.text, s.sums = text, sums
	return ch
}

// merge folds later changes into earlier ones of the same frame.
func merge(a, b []Change) []Change {
	for _, c := range b {
		found := false
		for i := range a {
			if a[i].Row == c.Row {
				a[i].StyleOnly = a[i].StyleOnly && c.StyleOnly
				a[i].Text = c.Text
				found = true
				break
			}
		}
		if !found {
			a = append(a, c)
		}
	}
	return a
}

// Settle waits until the screen has not changed for quiet since from (the
// keys went out, say) or since it last changed after that, or max has gone
// by, or the program has ended. It reports whether the screen went quiet.
func (s *Session) Settle(from time.Time, quiet, max time.Duration) bool {
	deadline := time.Now().Add(max)
	for {
		s.mu.Lock()
		last := s.changed
		s.mu.Unlock()
		if from.After(last) {
			last = from
		}
		since := time.Since(last)
		if since >= quiet {
			return true
		}
		if s.Exited() || time.Now().After(deadline) {
			return false
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// WaitChange waits until the screen changes after t, at most max, or the
// program ends.
func (s *Session) WaitChange(t time.Time, max time.Duration) {
	deadline := time.Now().Add(max)
	for time.Now().Before(deadline) && !s.Exited() {
		s.mu.Lock()
		changed := s.changed.After(t)
		s.mu.Unlock()
		if changed {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// Keys sends the keys one by one, gap apart: a terminal sends Esc then j
// typed at once as Alt-j otherwise.
func (s *Session) Keys(keys []string, gap time.Duration) error {
	for i, k := range keys {
		s.mu.Lock()
		app := s.scr.AppCursor()
		s.mu.Unlock()
		b, err := Encode(k, app)
		if err != nil {
			return err
		}
		if i > 0 && gap > 0 {
			time.Sleep(gap)
		}
		if _, err := s.pty.Write(b); err != nil {
			return err
		}
	}
	return nil
}

// Type sends text as typed.
func (s *Session) Type(text string) error {
	_, err := io.WriteString(s.pty, text)
	return err
}

// Record keeps every change of the screen from now until StopRecording.
func (s *Session) Record() {
	s.mu.Lock()
	s.rec = &recording{start: time.Now()}
	s.mu.Unlock()
}

// StopRecording returns the frames since Record.
func (s *Session) StopRecording() []Frame {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.rec == nil {
		return nil
	}
	f := s.rec.frames
	s.rec = nil
	return f
}

// Resize changes the terminal's size and tells the program.
func (s *Session) Resize(cols, rows int) error {
	s.mu.Lock()
	s.scr.Resize(cols, rows)
	s.sums = nil
	s.snapshot()
	s.changed = time.Now()
	s.mu.Unlock()
	return pty.Setsize(s.pty, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})
}

// View reads the screen under the lock.
func (s *Session) View(f func(scr *vt.Screen)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f(s.scr)
}

// Exited reports whether the program has ended.
func (s *Session) Exited() bool {
	select {
	case <-s.done:
		return true
	default:
		return false
	}
}

// ExitStatus describes how the program ended, "" while it runs.
func (s *Session) ExitStatus() string {
	if !s.Exited() {
		return ""
	}
	if s.exitErr == nil {
		return "exit status 0"
	}
	return s.exitErr.Error()
}

// Close hangs up, as closing a terminal window does, and kills the program
// if it is still there a second later.
func (s *Session) Close() string {
	if !s.Exited() && s.cmd.Process != nil {
		_ = s.cmd.Process.Signal(syscall.SIGHUP)
		select {
		case <-s.done:
		case <-time.After(time.Second):
			_ = s.cmd.Process.Kill()
			<-s.done
		}
	}
	s.pty.Close()
	return s.ExitStatus()
}
