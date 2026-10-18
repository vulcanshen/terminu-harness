package term

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Encode is what a terminal sends for one key, written the way the checklist
// writes keys: `Enter`, `Esc`, `Space`, `Tab`, `Shift-Tab`, `Up`, `Ctrl-C`,
// `Alt-x`, `F1`, or a single character. appCursor is the program's cursor key
// mode (DECCKM): arrows as ESC O A rather than ESC [ A.
func Encode(key string, appCursor bool) ([]byte, error) {
	if utf8.RuneCountInString(key) == 1 {
		return []byte(key), nil
	}
	name := strings.ToLower(key)
	var shift, alt, ctrl bool
	for {
		rest, ok := cutModifier(name)
		if !ok {
			break
		}
		switch name[0] {
		case 's':
			shift = true
		case 'a', 'm':
			alt = true
		case 'c':
			ctrl = true
		}
		name = rest
	}
	// The key the modifiers go on keeps its case: Alt-X is ESC X.
	base := key[len(key)-len(name):]
	if name == "" {
		return nil, fmt.Errorf("no key in %q", key)
	}
	if b, ok := special(name, shift, alt, ctrl, appCursor); ok {
		return b, nil
	}
	if utf8.RuneCountInString(base) != 1 {
		return nil, fmt.Errorf("unknown key %q", key)
	}
	b := []byte(base)
	if ctrl {
		c, ok := ctrlByte(base[0])
		if !ok || len(base) != 1 {
			return nil, fmt.Errorf("no Ctrl for %q", key)
		}
		b = []byte{c}
	} else if shift && len(base) == 1 {
		b = []byte(strings.ToUpper(base))
	}
	if alt {
		b = append([]byte{0x1b}, b...)
	}
	return b, nil
}

// cutModifier takes a leading ctrl-, alt-, shift- (or c-, a-, m-, s-, with -
// or +) off name.
func cutModifier(name string) (string, bool) {
	for _, m := range []string{"ctrl", "control", "alt", "meta", "shift", "c", "a", "m", "s"} {
		for _, sep := range []string{"-", "+"} {
			if p := m + sep; strings.HasPrefix(name, p) && len(name) > len(p) {
				return name[len(p):], true
			}
		}
	}
	return name, false
}

func ctrlByte(c byte) (byte, bool) {
	switch {
	case c >= 'a' && c <= 'z':
		return c - 'a' + 1, true
	case c >= 'A' && c <= 'Z':
		return c - 'A' + 1, true
	case c == ' ' || c == '@' || c == '2':
		return 0, true
	case c == '[':
		return 0x1b, true
	case c == '\\':
		return 0x1c, true
	case c == ']':
		return 0x1d, true
	case c == '^':
		return 0x1e, true
	case c == '_' || c == '/':
		return 0x1f, true
	}
	return 0, false
}

func special(name string, shift, alt, ctrl, appCursor bool) ([]byte, bool) {
	mod := 1
	if shift {
		mod++
	}
	if alt {
		mod += 2
	}
	if ctrl {
		mod += 4
	}
	esc := func(s string) []byte { return []byte("\x1b" + s) }
	withAlt := func(b []byte) []byte {
		if alt {
			return append([]byte{0x1b}, b...)
		}
		return b
	}
	// cursor keys: CSI A, CSI 1;m A with modifiers, SS3 A in cursor key mode.
	cursor := func(c string) []byte {
		switch {
		case mod > 1:
			return esc(fmt.Sprintf("[1;%d%s", mod, c))
		case appCursor:
			return esc("O" + c)
		}
		return esc("[" + c)
	}
	tilde := func(n int) []byte {
		if mod > 1 {
			return esc(fmt.Sprintf("[%d;%d~", n, mod))
		}
		return esc(fmt.Sprintf("[%d~", n))
	}
	switch name {
	case "enter", "return":
		return withAlt([]byte{'\r'}), true
	case "tab":
		if shift {
			return esc("[Z"), true
		}
		return withAlt([]byte{'\t'}), true
	case "backtab":
		return esc("[Z"), true
	case "esc", "escape":
		return withAlt([]byte{0x1b}), true
	case "space":
		if ctrl {
			return []byte{0}, true
		}
		return withAlt([]byte{' '}), true
	case "backspace", "bs":
		if ctrl {
			return withAlt([]byte{0x08}), true
		}
		return withAlt([]byte{0x7f}), true
	case "up":
		return cursor("A"), true
	case "down":
		return cursor("B"), true
	case "right":
		return cursor("C"), true
	case "left":
		return cursor("D"), true
	case "home":
		return cursor("H"), true
	case "end":
		return cursor("F"), true
	case "insert", "ins":
		return tilde(2), true
	case "delete", "del":
		return tilde(3), true
	case "pgup", "pageup":
		return tilde(5), true
	case "pgdn", "pagedown":
		return tilde(6), true
	}
	if n, ok := fkey(name); ok {
		if n <= 4 {
			c := string(rune('P' + n - 1))
			if mod > 1 {
				return esc(fmt.Sprintf("[1;%d%s", mod, c)), true
			}
			return esc("O" + c), true
		}
		codes := map[int]int{5: 15, 6: 17, 7: 18, 8: 19, 9: 20, 10: 21, 11: 23, 12: 24}
		return tilde(codes[n]), true
	}
	return nil, false
}

func fkey(name string) (int, bool) {
	var n int
	if _, err := fmt.Sscanf(name, "f%d", &n); err != nil || fmt.Sprintf("f%d", n) != name || n < 1 || n > 12 {
		return 0, false
	}
	return n, true
}
