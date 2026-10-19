//go:build unix

// Package server is the harness's MCP server: a terminal the judge runs an
// app in, tdp's files, and the run that keeps the verdicts.
package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/vulcanshen/terminu-harness/internal/render"
	"github.com/vulcanshen/terminu-harness/internal/run"
	"github.com/vulcanshen/terminu-harness/internal/tdp"
	"github.com/vulcanshen/terminu-harness/internal/term"
	"github.com/vulcanshen/terminu-harness/internal/vt"
)

// The environment the harness reads (named the way tdp E2 names an app's).
const (
	EnvHome    = "TERMINU_HARNESS__HOME"     // where the harness keeps its things: ~/.terminu-harness
	EnvOut     = "TERMINU_HARNESS__OUT"      // where runs go: <home>/judge
	EnvTDPRepo = "TERMINU_HARNESS__TDP_REPO" // a local clone of terminu to read tdp from, rather than GitHub
)

// batches is how many judges share the checklist.
const batches = 10

type entry struct {
	s                 *term.Session
	name              string
	run               *run.Run
	wraps, scrolls    int
	cols, rows        int
	iconW, powerlineW int
}

// Server holds the open sessions and runs.
type Server struct {
	mu       sync.Mutex
	sessions map[string]*entry
	runs     map[string]*run.Run
	next     int
	src      tdp.Source
	home     string
}

// Home is where the harness keeps its things.
func Home() string {
	if h := os.Getenv(EnvHome); h != "" {
		return h
	}
	dir, err := os.UserHomeDir()
	if err != nil {
		return ".terminu-harness"
	}
	return filepath.Join(dir, ".terminu-harness")
}

// New is a server with nothing open.
func New() *Server {
	home := Home()
	return &Server{
		sessions: map[string]*entry{},
		runs:     map[string]*run.Run{},
		src:      tdp.Source{Repo: os.Getenv(EnvTDPRepo), Cache: filepath.Join(home, "cache", "tdp")},
		home:     home,
	}
}

// Close ends every program still open.
func (sv *Server) Close() {
	sv.mu.Lock()
	defer sv.mu.Unlock()
	for id, e := range sv.sessions {
		e.s.Close()
		delete(sv.sessions, id)
	}
}

// MCP is the server's tools on an MCP server.
func (sv *Server) MCP(version string) *server.MCPServer {
	s := server.NewMCPServer("terminu-harness", version, server.WithToolCapabilities(false))
	session := mcp.WithString("session", mcp.Required(), mcp.Description("The session id open returned, e.g. s1."))
	runDir := mcp.WithString("run", mcp.Required(), mcp.Description("The run directory judge_start returned."))

	s.AddTool(mcp.NewTool("open",
		mcp.WithDescription("Start an app in a new terminal of the harness's own (a PTY and an emulator; no tmux, no desktop) and return its first screen. "+
			"The terminal answers what a terminu family app asks at startup (cursor position, device attributes, colours) and gives an icon "+
			"and a powerline symbol the cells you say. TERM is xterm-256color, COLORTERM truecolor; icon width variables of the real terminal "+
			"(*__ICON_WIDTH, *__POWERLINE_WIDTH) are not passed on unless you set them in env. The app keeps running until close."),
		mcp.WithString("command", mcp.Required(), mcp.Description("The program to run: a name on PATH or a path.")),
		mcp.WithArray("args", mcp.WithStringItems(), mcp.Description("Its arguments.")),
		mcp.WithString("cwd", mcp.Description("Where to run it; the harness's own directory when left out.")),
		mcp.WithNumber("cols", mcp.Description("Columns, 120 when left out.")),
		mcp.WithNumber("rows", mcp.Description("Rows, 40 when left out.")),
		mcp.WithNumber("icon_width", mcp.Description("Cells a Nerd Font icon takes on this terminal: 1 (default) or 2.")),
		mcp.WithNumber("powerline_width", mcp.Description("Cells a powerline symbol (U+E0A0–U+E0D7) takes: 1 (default) or 2.")),
		mcp.WithArray("env", mcp.WithStringItems(), mcp.Description("Extra environment, KEY=VALUE each.")),
		mcp.WithString("run", mcp.Description("A run directory from judge_start: every screen, recording and image of this session is kept there as evidence.")),
		mcp.WithNumber("wait_ms", mcp.Description("How long at most to wait for the screen to settle before returning it, 5000 when left out.")),
	), sv.open)

	s.AddTool(mcp.NewTool("keys",
		mcp.WithDescription("Press keys in a session and return the screen once it settles (no change for 150ms, at most wait_ms). "+
			"Keys are written the checklist's way: Enter, Esc, Space, Tab, Shift-Tab, Backspace, Delete, Up, Down, Left, Right, Home, End, "+
			"PgUp, PgDn, F1–F12, Ctrl-C, Alt-x, Shift-Down, or one character (j, ?, V, 1). They go gap_ms apart, so Esc then j stays two keys. "+
			"text is typed after the keys, as is. With record_ms the screen is recorded from just before the first key for that long instead: "+
			"every change the app drew, with its time in ms and the rows that changed (rows marked 'colours only' changed colour or attributes, "+
			"not characters). Record with no keys to watch time pass (a toast's countdown, an animation already running)."),
		session,
		mcp.WithArray("keys", mcp.WithStringItems(), mcp.Description("The keys, in order.")),
		mcp.WithString("text", mcp.Description("Text to type after the keys.")),
		mcp.WithNumber("gap_ms", mcp.Description("Time between keys, 60 when left out.")),
		mcp.WithNumber("record_ms", mcp.Description("Record this long and return the frames.")),
		mcp.WithNumber("wait_ms", mcp.Description("How long at most to wait for the screen to settle, 2000 when left out.")),
	), sv.keys)

	s.AddTool(mcp.NewTool("screen",
		mcp.WithDescription("The session's screen as text, row by row with row numbers (0-based), the cursor, the size, and whether the app has ended. "+
			"A wide character (CJK, an icon given 2 cells) is one character in the text but two columns; use styles for exact columns and colours."),
		session,
	), sv.screen)

	s.AddTool(mcp.NewTool("styles",
		mcp.WithDescription("Colours and attributes of a block of the screen: for each row, runs of cells drawn alike, with their columns (0-based), "+
			"text, foreground, background (#rrggbb, or default when the app left it to the terminal) and bold, faint, italic, underline, reverse, strike."),
		session,
		mcp.WithNumber("row", mcp.Required(), mcp.Description("The first row, 0-based.")),
		mcp.WithNumber("height", mcp.Description("How many rows, 1 when left out.")),
		mcp.WithNumber("col", mcp.Description("The first column, 0 when left out.")),
		mcp.WithNumber("width", mcp.Description("How many columns, to the end of the row when left out.")),
	), sv.styles)

	s.AddTool(mcp.NewTool("image",
		mcp.WithDescription("The session's screen drawn as a PNG, with the fonts the harness carries (Go Mono, Symbols Nerd Font Mono; box drawing "+
			"and blocks drawn as shapes). A character no font has is an outline. Use it to look at the layout as a whole; read exact text and "+
			"colours with screen and styles."),
		session,
	), sv.image)

	s.AddTool(mcp.NewTool("resize",
		mcp.WithDescription("Change the terminal's size (the app gets SIGWINCH) and return the screen once it settles."),
		session,
		mcp.WithNumber("cols", mcp.Required(), mcp.Description("Columns.")),
		mcp.WithNumber("rows", mcp.Required(), mcp.Description("Rows.")),
		mcp.WithNumber("wait_ms", mcp.Description("How long at most to wait for the screen to settle, 2000 when left out.")),
	), sv.resize)

	s.AddTool(mcp.NewTool("overflow",
		mcp.WithDescription("Check that the app never printed past the edge since open or the last overflow call: how many characters wrapped "+
			"to the next row and how many rows scrolled off the top. A full-screen app whose every row is exactly as wide as the terminal "+
			"neither wraps nor scrolls."),
		session,
	), sv.overflow)

	s.AddTool(mcp.NewTool("close",
		mcp.WithDescription("End the session: hang up as closing a terminal window does, kill the app a second later if it is still there, and say how it ended."),
		session,
	), sv.close)

	s.AddTool(mcp.NewTool("tdp_file",
		mcp.WithDescription("A file of terminu design (tdp) at a version, e.g. principle/rules-zh_TW.md, components/layout/popup-zh_TW.md, "+
			"project/family.mk, principle/icon.svg, checklist/README-zh_TW.md. Traditional Chinese is the original."),
		mcp.WithString("path", mcp.Required(), mcp.Description("The path in the terminu repo.")),
		mcp.WithString("version", mcp.Description("A tdp version like v0.7.4; the run's version with run, the newest tag otherwise.")),
		mcp.WithString("run", mcp.Description("A run directory: use the version it judges by.")),
	), sv.tdpFile)

	s.AddTool(mcp.NewTool("judge_start",
		mcp.WithDescription("Start judging an app against a version of tdp: makes the run directory, keeps the checklist in it and cuts the "+
			"items into batches, one file each, for the judges. Returns the run directory and the batches."),
		mcp.WithString("app_command", mcp.Required(), mcp.Description("How the app is started, e.g. demou or ./demou.")),
		mcp.WithString("repo", mcp.Description("The app's repo, when there is one: items to check by reading code or files need it.")),
		mcp.WithString("tdp_version", mcp.Description("The tdp version to judge by, the newest tag when left out.")),
		mcp.WithString("out_dir", mcp.Description("Where runs go, ~/.terminu-harness/judge when left out.")),
	), sv.judgeStart)

	s.AddTool(mcp.NewTool("judge_profile",
		mcp.WithDescription("Keep what the judge learned about the app (profile.md in the run): its help and version, screens, panels, popups, "+
			"the features the checklist's 適用 column names and which the app has, and how to reach each surface. Replaces the earlier one."),
		runDir,
		mcp.WithString("text", mcp.Required(), mcp.Description("The profile, in Markdown.")),
	), sv.judgeProfile)

	s.AddTool(mcp.NewTool("judge_verdicts",
		mcp.WithDescription("Keep verdicts on items. verdict is one of 符合 偏離 違反 不適用 判不了; reason says why in Traditional Chinese "+
			"(for 判不了 what is missing); for 判不了, cause says why: tdp (tdp says nothing clear), tgo (terminu-go can't), other "+
			"(no repo, a special environment); evidence lists evidence ids the terminal tools returned. A later verdict on an item replaces an "+
			"earlier one. Returns how many were kept and what was wrong with the rest."),
		runDir,
		mcp.WithArray("verdicts", mcp.Required(), mcp.Items(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":       map[string]any{"type": "string", "description": "The item number, e.g. K4-2."},
				"verdict":  map[string]any{"type": "string", "enum": run.Verdicts},
				"reason":   map[string]any{"type": "string"},
				"cause":    map[string]any{"type": "string", "enum": run.Causes},
				"evidence": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			},
			"required": []string{"id", "verdict", "reason"},
		}), mcp.Description("The verdicts.")),
	), sv.judgeVerdicts)

	s.AddTool(mcp.NewTool("judge_gap",
		mcp.WithDescription("Keep one thing the app does that tdp does not cover yet: no gene, no penalty, reported back to tdp."),
		runDir,
		mcp.WithString("what", mcp.Required(), mcp.Description("What the app does, in Traditional Chinese.")),
		mcp.WithString("where", mcp.Description("Where in the app.")),
		mcp.WithArray("evidence", mcp.WithStringItems(), mcp.Description("Evidence ids.")),
	), sv.judgeGap)

	s.AddTool(mcp.NewTool("judge_status",
		mcp.WithDescription("How far the run has come: how many items have each verdict, and batch by batch the items no one has judged yet."),
		runDir,
	), sv.judgeStatus)

	s.AddTool(mcp.NewTool("judge_finish",
		mcp.WithDescription("Write the report card (report.md in the run) and return its summary with the gene count. Items no one judged count as 判不了."),
		runDir,
	), sv.judgeFinish)
	return s
}

func fail(format string, a ...any) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultError(fmt.Sprintf(format, a...)), nil
}

func (sv *Server) entry(req mcp.CallToolRequest) (string, *entry, error) {
	id := req.GetString("session", "")
	sv.mu.Lock()
	defer sv.mu.Unlock()
	e, ok := sv.sessions[id]
	if !ok {
		return id, nil, fmt.Errorf("no session %q", id)
	}
	return id, e, nil
}

func (sv *Server) openRun(dir string) (*run.Run, error) {
	if dir == "" {
		return nil, nil
	}
	dir = filepath.Clean(dir)
	sv.mu.Lock()
	defer sv.mu.Unlock()
	if r, ok := sv.runs[dir]; ok {
		return r, nil
	}
	r, err := run.Open(dir)
	if err != nil {
		return nil, err
	}
	sv.runs[dir] = r
	return r, nil
}

func ms(req mcp.CallToolRequest, key string, def int) time.Duration {
	return time.Duration(req.GetInt(key, def)) * time.Millisecond
}

func (sv *Server) open(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	command := req.GetString("command", "")
	args := req.GetStringSlice("args", nil)
	if len(args) == 0 && strings.Contains(command, " ") {
		if _, err := exec.LookPath(command); err != nil {
			f := strings.Fields(command)
			command, args = f[0], f[1:]
		}
	}
	r, err := sv.openRun(req.GetString("run", ""))
	if err != nil {
		return fail("%v", err)
	}
	cols, rows := req.GetInt("cols", 120), req.GetInt("rows", 40)
	if cols < 1 || rows < 1 || cols > 1000 || rows > 1000 {
		return fail("a %d×%d terminal", cols, rows)
	}
	iconW, plW := req.GetInt("icon_width", 1), req.GetInt("powerline_width", 1)
	if iconW < 1 || iconW > 2 || plW < 1 || plW > 2 {
		return fail("icon_width and powerline_width are 1 or 2")
	}
	started := time.Now()
	s, err := term.Start(term.Options{
		Argv: append([]string{command}, args...), Dir: req.GetString("cwd", ""), Env: req.GetStringSlice("env", nil),
		Cols: cols, Rows: rows, IconCells: iconW, PowerlineCells: plW,
	})
	if err != nil {
		return fail("start %s: %v", command, err)
	}
	sv.mu.Lock()
	sv.next++
	id := "s" + strconv.Itoa(sv.next)
	e := &entry{s: s, name: filepath.Base(command), run: r, cols: cols, rows: rows, iconW: iconW, powerlineW: plW}
	sv.sessions[id] = e
	sv.mu.Unlock()
	// A slow start draws nothing for a while: wait for the first change, then
	// for the screen to settle.
	wait := ms(req, "wait_ms", 5000)
	s.WaitChange(started, wait)
	s.Settle(started, 300*time.Millisecond, wait-time.Since(started))
	return mcp.NewToolResultText(sv.screenText(id, e, "")), nil
}

func (sv *Server) keys(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, e, err := sv.entry(req)
	if err != nil {
		return fail("%v", err)
	}
	keys, text := req.GetStringSlice("keys", nil), req.GetString("text", "")
	recordFor := ms(req, "record_ms", 0)
	if recordFor > 0 {
		e.s.Record()
	}
	if err := e.s.Keys(keys, ms(req, "gap_ms", 60)); err != nil {
		e.s.StopRecording()
		return fail("%v", err)
	}
	if text != "" {
		if err := e.s.Type(text); err != nil {
			e.s.StopRecording()
			return fail("%v", err)
		}
	}
	if recordFor <= 0 {
		e.s.Settle(time.Now(), 150*time.Millisecond, ms(req, "wait_ms", 2000))
		return mcp.NewToolResultText(sv.screenText(id, e, "")), nil
	}
	time.Sleep(recordFor)
	frames := e.s.StopRecording()
	return mcp.NewToolResultText(sv.screenText(id, e, framesText(frames, recordFor))), nil
}

func framesText(frames []term.Frame, d time.Duration) string {
	var b strings.Builder
	fmt.Fprintf(&b, "recorded %dms from just before the first key: %d frames\n", d.Milliseconds(), len(frames))
	for i, f := range frames {
		rows := make([]int, len(f.Changes))
		for j, c := range f.Changes {
			rows[j] = c.Row
		}
		fmt.Fprintf(&b, "frame %d +%.1fms, %d rows changed\n", i+1, float64(f.At.Microseconds())/1000, len(rows))
		for _, c := range f.Changes {
			mark := ""
			if c.StyleOnly {
				mark = "  (colours only)"
			}
			fmt.Fprintf(&b, "  %2d│%s%s\n", c.Row, c.Text, mark)
		}
	}
	b.WriteString("the screen at the end of the recording:\n")
	return b.String()
}

// screenText is the screen with a header, kept as evidence when the session
// belongs to a run; before goes above it (a recording).
func (sv *Server) screenText(id string, e *entry, before string) string {
	var body strings.Builder
	e.s.View(func(scr *vt.Screen) {
		cols, rows := scr.Size()
		x, y, visible := scr.Cursor()
		icon, pl := scr.Widths()
		cursor := "shown"
		if !visible {
			cursor = "hidden"
		}
		fmt.Fprintf(&body, "%s %s · %d×%d · icon %d, powerline %d · cursor row %d col %d (%s)", id, e.name, cols, rows, icon, pl, y, x, cursor)
		if scr.AltScreen() {
			body.WriteString(" · alt screen")
		}
		if st := e.s.ExitStatus(); st != "" {
			fmt.Fprintf(&body, " · the app has ended: %s", st)
		}
		body.WriteString("\n")
		width := len(strconv.Itoa(rows - 1))
		for i, line := range scr.Text() {
			fmt.Fprintf(&body, "%*d│%s\n", width, i, line)
		}
	})
	out := before + body.String()
	if e.run != nil {
		eid, path := e.run.NewEvidence(".txt")
		if err := os.WriteFile(path, []byte(out), 0o644); err == nil {
			out = "evidence " + eid + "\n" + out
		}
	}
	const most = 60000
	if len(out) > most {
		out = out[:most] + "\n… cut here; the whole of it is in the evidence file\n"
	}
	return out
}

func (sv *Server) screen(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, e, err := sv.entry(req)
	if err != nil {
		return fail("%v", err)
	}
	return mcp.NewToolResultText(sv.screenText(id, e, "")), nil
}

func (sv *Server) styles(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, e, err := sv.entry(req)
	if err != nil {
		return fail("%v", err)
	}
	row0, height := req.GetInt("row", 0), req.GetInt("height", 1)
	col0, width := req.GetInt("col", 0), req.GetInt("width", 0)
	var b strings.Builder
	var bad error
	e.s.View(func(scr *vt.Screen) {
		cols, rows := scr.Size()
		if row0 < 0 || row0 >= rows || height < 1 || col0 < 0 || col0 >= cols {
			bad = fmt.Errorf("row %d, height %d, col %d outside a %d×%d screen", row0, height, col0, cols, rows)
			return
		}
		end := cols
		if width > 0 {
			end = min(col0+width, cols)
		}
		fmt.Fprintf(&b, "%s styles, rows %d–%d, cols %d–%d\n", id, row0, min(row0+height, rows)-1, col0, end-1)
		for y := row0; y < min(row0+height, rows); y++ {
			fmt.Fprintf(&b, "row %d:\n", y)
			writeRuns(&b, scr.Row(y), col0, end)
		}
	})
	if bad != nil {
		return fail("%v", bad)
	}
	out := b.String()
	if e.run != nil {
		eid, path := e.run.NewEvidence(".txt")
		if os.WriteFile(path, []byte(out), 0o644) == nil {
			out = "evidence " + eid + "\n" + out
		}
	}
	return mcp.NewToolResultText(out), nil
}

// writeRuns writes the cells [from, to) of a row as runs drawn alike.
func writeRuns(b *strings.Builder, row []vt.Cell, from, to int) {
	style := func(c vt.Cell) string {
		color := func(rgb vt.RGB, def bool) string {
			if def {
				return "default"
			}
			return rgb.Hex()
		}
		s := "fg " + color(c.Fg, c.FgDefault) + " bg " + color(c.Bg, c.BgDefault)
		for _, a := range []struct {
			on   bool
			name string
		}{{c.Bold, "bold"}, {c.Faint, "faint"}, {c.Italic, "italic"}, {c.Underline, "underline"}, {c.Reverse, "reverse"}, {c.Hidden, "hidden"}, {c.Strike, "strike"}} {
			if a.on {
				s += " " + a.name
			}
		}
		return s
	}
	start, cur := from, ""
	var text strings.Builder
	flush := func(end int) {
		if end > start {
			fmt.Fprintf(b, "  %d–%d %q %s\n", start, end-1, text.String(), cur)
		}
	}
	for x := from; x < to; x++ {
		c := row[x]
		st := style(c)
		if c.Width == 0 && x > from {
			continue // the right half of a wide character goes with it
		}
		if st != cur {
			flush(x)
			start, cur = x, st
			text.Reset()
		}
		switch {
		case c.Width == 0:
		case c.Text == "":
			text.WriteString(" ")
		default:
			text.WriteString(c.Text)
		}
	}
	flush(to)
}

func (sv *Server) image(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, e, err := sv.entry(req)
	if err != nil {
		return fail("%v", err)
	}
	var png []byte
	e.s.View(func(scr *vt.Screen) { png, err = render.PNG(scr) })
	if err != nil {
		return fail("draw: %v", err)
	}
	text := id + " screen"
	if e.run != nil {
		eid, path := e.run.NewEvidence(".png")
		if os.WriteFile(path, png, 0o644) == nil {
			text = "evidence " + eid + "\n" + text
		}
	}
	return mcp.NewToolResultImage(text, base64.StdEncoding.EncodeToString(png), "image/png"), nil
}

func (sv *Server) resize(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, e, err := sv.entry(req)
	if err != nil {
		return fail("%v", err)
	}
	cols, rows := req.GetInt("cols", 0), req.GetInt("rows", 0)
	if cols < 1 || rows < 1 || cols > 1000 || rows > 1000 {
		return fail("a %d×%d terminal", cols, rows)
	}
	if err := e.s.Resize(cols, rows); err != nil {
		return fail("%v", err)
	}
	e.s.Settle(time.Now(), 150*time.Millisecond, ms(req, "wait_ms", 2000))
	return mcp.NewToolResultText(sv.screenText(id, e, "")), nil
}

func (sv *Server) overflow(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, e, err := sv.entry(req)
	if err != nil {
		return fail("%v", err)
	}
	var wraps, scrolls int
	e.s.View(func(scr *vt.Screen) { wraps, scrolls = scr.Counts() })
	sv.mu.Lock()
	dw, ds := wraps-e.wraps, scrolls-e.scrolls
	e.wraps, e.scrolls = wraps, scrolls
	sv.mu.Unlock()
	verdict := "no row went past the edge"
	if dw > 0 || ds > 0 {
		verdict = "the app printed past the edge"
	}
	return mcp.NewToolResultText(fmt.Sprintf("%s: %d characters wrapped, %d rows scrolled off the top since the last check: %s", id, dw, ds, verdict)), nil
}

func (sv *Server) close(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, e, err := sv.entry(req)
	if err != nil {
		return fail("%v", err)
	}
	status := e.s.Close()
	sv.mu.Lock()
	delete(sv.sessions, id)
	sv.mu.Unlock()
	return mcp.NewToolResultText(fmt.Sprintf("%s closed: %s", id, status)), nil
}

func (sv *Server) tdpFile(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	version := req.GetString("version", "")
	if version == "" && req.GetString("run", "") != "" {
		r, err := sv.openRun(req.GetString("run", ""))
		if err != nil {
			return fail("%v", err)
		}
		version = r.Meta.TDP
	}
	if version == "" {
		v, err := sv.src.Latest()
		if err != nil {
			return fail("the newest tdp: %v", err)
		}
		version = v
	}
	b, err := sv.src.File(version, req.GetString("path", ""))
	if err != nil {
		return fail("%v", err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("tdp %s %s\n\n%s", tdp.Tag(version), req.GetString("path", ""), b)), nil
}

func (sv *Server) judgeStart(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	command := strings.TrimSpace(req.GetString("app_command", ""))
	if command == "" {
		return fail("app_command is empty")
	}
	tag := tdp.Tag(req.GetString("tdp_version", ""))
	if tag == "" {
		v, err := sv.src.Latest()
		if err != nil {
			return fail("the newest tdp: %v", err)
		}
		tag = v
	}
	md, err := sv.src.File(tag, tdp.ItemsPath)
	if err != nil {
		return fail("the checklist of tdp %s: %v", tag, err)
	}
	base := req.GetString("out_dir", "")
	if base == "" {
		base = os.Getenv(EnvOut)
	}
	if base == "" {
		base = filepath.Join(sv.home, "judge")
	}
	repo := req.GetString("repo", "")
	if repo != "" {
		if abs, err := filepath.Abs(repo); err == nil {
			repo = abs
		}
	}
	r, err := run.Start(base, command, repo, tag, md, batches)
	if err != nil {
		return fail("%v", err)
	}
	sv.mu.Lock()
	sv.runs[r.Dir] = r
	sv.mu.Unlock()
	carried, err := carryOver(r, base, repo, tag)
	if err != nil {
		return fail("%v", err)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "run %s\ntdp %s, %d items, %d batches:\n", r.Dir, tag, r.Total(), len(r.Meta.Batches))
	if carried != "" {
		b.WriteString(carried)
	}
	for _, bt := range r.Meta.Batches {
		fmt.Fprintf(&b, "  %s %d items: %s — %s\n", bt.ID, bt.Count, filepath.Join(r.Dir, bt.File), strings.Join(bt.Sections, "、"))
	}
	fmt.Fprintf(&b, "checklist README: tdp_file checklist/README-zh_TW.md with run\n")
	return mcp.NewToolResultText(b.String()), nil
}

func (sv *Server) judgeProfile(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	r, err := sv.openRun(req.GetString("run", ""))
	if err != nil || r == nil {
		return fail("run: %v", err)
	}
	if err := r.SetProfile(req.GetString("text", "")); err != nil {
		return fail("%v", err)
	}
	return mcp.NewToolResultText("kept " + filepath.Join(r.Dir, "profile.md")), nil
}

func (sv *Server) judgeVerdicts(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	r, err := sv.openRun(req.GetString("run", ""))
	if err != nil || r == nil {
		return fail("run: %v", err)
	}
	raw, err := json.Marshal(req.GetArguments()["verdicts"])
	if err != nil {
		return fail("verdicts: %v", err)
	}
	var vs []run.Verdict
	if err := json.Unmarshal(raw, &vs); err != nil {
		return fail("verdicts: %v", err)
	}
	kept, problems := r.AddVerdicts(vs)
	out := fmt.Sprintf("kept %d of %d", kept, len(vs))
	if len(problems) > 0 {
		out += "\n" + strings.Join(problems, "\n")
	}
	return mcp.NewToolResultText(out), nil
}

func (sv *Server) judgeGap(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	r, err := sv.openRun(req.GetString("run", ""))
	if err != nil || r == nil {
		return fail("run: %v", err)
	}
	g := run.Gap{What: req.GetString("what", ""), Where: req.GetString("where", ""), Evidence: req.GetStringSlice("evidence", nil)}
	if err := r.AddGap(g); err != nil {
		return fail("%v", err)
	}
	return mcp.NewToolResultText("kept"), nil
}

func (sv *Server) judgeStatus(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	r, err := sv.openRun(req.GetString("run", ""))
	if err != nil || r == nil {
		return fail("run: %v", err)
	}
	s, err := r.Status()
	if err != nil {
		return fail("%v", err)
	}
	return mcp.NewToolResultText(s), nil
}

func (sv *Server) judgeFinish(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	r, err := sv.openRun(req.GetString("run", ""))
	if err != nil || r == nil {
		return fail("run: %v", err)
	}
	s, err := r.Finish()
	if err != nil {
		return fail("%v", err)
	}
	return mcp.NewToolResultText(s), nil
}
