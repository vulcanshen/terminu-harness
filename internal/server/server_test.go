//go:build unix

package server

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// fixture is the test program (testdata/fixture), built once.
var fixture string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "terminu-harness-fixture")
	if err != nil {
		panic(err)
	}
	fixture = filepath.Join(dir, "fixture")
	if out, err := exec.Command("go", "build", "-o", fixture, "./testdata/fixture").CombinedOutput(); err != nil {
		panic(string(out))
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

type tool func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)

// call runs a tool and returns its text, failing the test on a tool error.
func call(t *testing.T, f tool, args map[string]any) (string, *mcp.CallToolResult) {
	t.Helper()
	var req mcp.CallToolRequest
	req.Params.Arguments = args
	res, err := f(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	var text []string
	for _, c := range res.Content {
		if tc, ok := c.(mcp.TextContent); ok {
			text = append(text, tc.Text)
		}
	}
	if res.IsError {
		t.Fatalf("tool error: %s", strings.Join(text, "\n"))
	}
	return strings.Join(text, "\n"), res
}

func TestTerminalTools(t *testing.T) {
	sv := New()
	defer sv.Close()

	// The terminal answers the probe with the icon width asked for.
	out, _ := call(t, sv.open, map[string]any{"command": fixture, "cols": 40, "rows": 10, "icon_width": 2})
	for _, want := range []string{"s1 fixture · 40×10 · icon 2, powerline 1", "alt screen", "0│fixture", "1│icon 2", "4│frame 0"} {
		if !strings.Contains(out, want) {
			t.Errorf("open: no %q in\n%s", want, out)
		}
	}

	// Keys arrive as a terminal sends them, Esc and j as two keys.
	for _, c := range []struct {
		keys []any
		want string
	}{{[]any{"x"}, "key 78"}, {[]any{"Up"}, "key 1b5b41"}, {[]any{"Ctrl-C"}, "key 03"}, {[]any{"Esc", "j"}, "key 6a"}} {
		out, _ = call(t, sv.keys, map[string]any{"session": "s1", "keys": c.keys})
		if !strings.Contains(out, c.want) {
			t.Errorf("keys %v: no %q in\n%s", c.keys, c.want, out)
		}
	}

	out, _ = call(t, sv.styles, map[string]any{"session": "s1", "row": 0, "width": 10})
	if !strings.Contains(out, `0–6 "fixture" fg #89b4fa bg default bold`) {
		t.Errorf("styles:\n%s", out)
	}

	// The animation is recorded frame by frame, about 40ms apart.
	out, _ = call(t, sv.keys, map[string]any{"session": "s1", "keys": []any{"a"}, "record_ms": 500})
	var at []float64
	for _, m := range regexp.MustCompile(`frame \d+ \+([\d.]+)ms`).FindAllStringSubmatch(out, -1) {
		v, _ := strconv.ParseFloat(m[1], 64)
		at = append(at, v)
	}
	if len(at) < 5 {
		t.Fatalf("%d frames:\n%s", len(at), out)
	}
	for i := 1; i < 5; i++ {
		if gap := at[i] - at[i-1]; gap < 30 || gap > 60 {
			t.Errorf("frame %d came %.1fms after the one before", i+1, gap)
		}
	}

	out, _ = call(t, sv.overflow, map[string]any{"session": "s1"})
	if !strings.Contains(out, "0 characters wrapped, 0 rows scrolled") {
		t.Errorf("overflow before: %s", out)
	}
	call(t, sv.keys, map[string]any{"session": "s1", "keys": []any{"w"}})
	out, _ = call(t, sv.overflow, map[string]any{"session": "s1"})
	if !strings.Contains(out, "printed past the edge") {
		t.Errorf("overflow after: %s", out)
	}

	out, _ = call(t, sv.resize, map[string]any{"session": "s1", "cols": 60, "rows": 12})
	if !strings.Contains(out, "size 60x12") || !strings.Contains(out, "60×12") {
		t.Errorf("resize:\n%s", out)
	}

	_, res := call(t, sv.image, map[string]any{"session": "s1"})
	var img *mcp.ImageContent
	for _, c := range res.Content {
		if ic, ok := c.(mcp.ImageContent); ok {
			img = &ic
		}
	}
	if img == nil || img.MIMEType != "image/png" || len(img.Data) == 0 {
		t.Errorf("no image")
	}

	call(t, sv.keys, map[string]any{"session": "s1", "keys": []any{"q"}})
	out, _ = call(t, sv.screen, map[string]any{"session": "s1"})
	if !strings.Contains(out, "the app has ended: exit status 0") {
		t.Errorf("after q:\n%s", out)
	}
	out, _ = call(t, sv.close, map[string]any{"session": "s1"})
	if !strings.Contains(out, "s1 closed: exit status 0") {
		t.Errorf("close: %s", out)
	}
	var req mcp.CallToolRequest
	req.Params.Arguments = map[string]any{"session": "s1"}
	if res, _ := sv.screen(context.Background(), req); !res.IsError {
		t.Error("a closed session still answers")
	}
}

// TestTdpFile reads tdp from a local clone of terminu, when there is one.
func TestTdpFile(t *testing.T) {
	if os.Getenv(EnvTDPRepo) == "" {
		t.Skip("no terminu clone (" + EnvTDPRepo + ")")
	}
	t.Setenv(EnvHome, t.TempDir())
	sv := New()
	out, _ := call(t, sv.tdpFile, map[string]any{"path": "checklist/README-zh_TW.md", "version": "v0.8.3"})
	if !strings.Contains(out, "# terminu checklist") {
		t.Errorf("tdp_file: %.200s", out)
	}
}
