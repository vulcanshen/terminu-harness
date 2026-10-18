// Command fixture is a small full-screen program the server's tests drive in
// place of a real app: it measures an icon the way a terminu family app does
// (print one, ask where the cursor went), shows the last key it got, plays a
// short animation on a, prints past the edge on w and ends on q.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var out = os.Stdout

func main() {
	if err := stty("raw", "-echo"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer stty("sane")
	fmt.Fprint(out, "\x1b[?1049h\x1b[?25l")
	defer fmt.Fprint(out, "\x1b[?25h\x1b[?1049l")

	icon := probe()
	key, frame := "none", 0
	draw := func() {
		cols, rows := size()
		fmt.Fprint(out, "\x1b[H\x1b[2J")
		fmt.Fprint(out, "\x1b[1;1H\x1b[1;38;2;137;180;250mfixture\x1b[0m")
		fmt.Fprintf(out, "\x1b[2;1Hicon %d", icon)
		fmt.Fprintf(out, "\x1b[3;1Hkey %s", key)
		fmt.Fprintf(out, "\x1b[4;1Hsize %dx%d", cols, rows)
		fmt.Fprintf(out, "\x1b[5;1Hframe %d", frame)
		fmt.Fprintf(out, "\x1b[%d;1H\x1b[38;2;88;91;112m┗%s┛\x1b[0m", rows, strings.Repeat("━", max(cols-2, 0)))
	}
	draw()

	winch := make(chan os.Signal, 1)
	signal.Notify(winch, syscall.SIGWINCH)
	keys := make(chan []byte)
	go func() {
		buf := make([]byte, 64)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				close(keys)
				return
			}
			keys <- append([]byte(nil), buf[:n]...)
		}
	}()
	for {
		select {
		case <-winch:
			draw()
		case b, ok := <-keys:
			if !ok {
				return
			}
			switch string(b) {
			case "q":
				return
			case "a":
				for frame = 1; frame <= 5; frame++ {
					fmt.Fprintf(out, "\x1b[5;1H\x1b[Kframe %d", frame)
					time.Sleep(40 * time.Millisecond)
				}
				frame = 5
			case "w":
				cols, _ := size()
				fmt.Fprintf(out, "\x1b[8;1H%s", strings.Repeat("x", cols+5))
				continue
			default:
				key = fmt.Sprintf("%x", b)
			}
			draw()
		}
	}
}

// probe prints an icon at the start of a line and asks where the cursor is.
func probe() int {
	fmt.Fprint(out, "\r\U000F024B\x1b[6n")
	var reply []byte
	b := make([]byte, 1)
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stdin.Read(b); err != nil {
			return 0
		}
		reply = append(reply, b[0])
		if b[0] == 'R' {
			break
		}
	}
	var row, col int
	if _, err := fmt.Sscanf(string(reply), "\x1b[%d;%dR", &row, &col); err != nil {
		return 0
	}
	return col - 1
}

func stty(args ...string) error {
	cmd := exec.Command("stty", args...)
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// size asks the terminal for its size (TIOCGWINSZ).
func size() (cols, rows int) {
	var ws struct{ rows, cols, x, y uint16 }
	syscall.Syscall(syscall.SYS_IOCTL, 0, uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&ws)))
	return int(ws.cols), int(ws.rows)
}
