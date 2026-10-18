package term

import "testing"

func TestEncode(t *testing.T) {
	cases := []struct {
		key       string
		appCursor bool
		want      string
	}{
		{"j", false, "j"},
		{"中", false, "中"},
		{"Enter", false, "\r"},
		{"Esc", false, "\x1b"},
		{"Space", false, " "},
		{"Tab", false, "\t"},
		{"Shift-Tab", false, "\x1b[Z"},
		{"Up", false, "\x1b[A"},
		{"Up", true, "\x1bOA"},
		{"Shift-Down", true, "\x1b[1;2B"},
		{"Ctrl-C", false, "\x03"},
		{"ctrl+c", false, "\x03"},
		{"C-a", false, "\x01"},
		{"Alt-x", false, "\x1bx"},
		{"Alt-X", false, "\x1bX"},
		{"F1", false, "\x1bOP"},
		{"F5", false, "\x1b[15~"},
		{"PgDn", false, "\x1b[6~"},
		{"Delete", false, "\x1b[3~"},
		{"Backspace", false, "\x7f"},
		{"Shift-V", false, "V"},
	}
	for _, c := range cases {
		got, err := Encode(c.key, c.appCursor)
		if err != nil || string(got) != c.want {
			t.Errorf("%q: %q, %v; want %q", c.key, got, err, c.want)
		}
	}
	for _, bad := range []string{"Hyper-x", "Foo", "Ctrl-", "F13"} {
		if _, err := Encode(bad, false); err == nil {
			t.Errorf("%q encoded", bad)
		}
	}
}
