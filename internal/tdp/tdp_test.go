package tdp

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestGetRetries: a 5xx, a 429 or a network error is tried again, up to
// three tries in all; a 404 is final.
func TestGetRetries(t *testing.T) {
	saved := retryWait
	retryWait = []time.Duration{time.Millisecond, time.Millisecond}
	defer func() { retryWait = saved }()
	calls := 0
	codes := []int{http.StatusServiceUnavailable, http.StatusTooManyRequests, http.StatusOK}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		code := codes[min(calls, len(codes)-1)]
		calls++
		w.WriteHeader(code)
		_, _ = w.Write([]byte("tdp"))
	}))
	defer srv.Close()
	if b, err := get(srv.URL); err != nil || string(b) != "tdp" || calls != 3 {
		t.Errorf("503, 429, 200: %q, %v, %d calls", b, err, calls)
	}
	calls, codes = 0, []int{http.StatusBadGateway}
	if _, err := get(srv.URL); err == nil || calls != 3 {
		t.Errorf("502 each time: %v, %d calls", err, calls)
	}
	calls, codes = 0, []int{http.StatusNotFound}
	if _, err := get(srv.URL); err == nil || calls != 1 {
		t.Errorf("404: %v, %d calls", err, calls)
	}
	srv.Close()
	if _, err := get(srv.URL); err == nil {
		t.Errorf("a closed server: no error")
	}
}

func TestVersionOrder(t *testing.T) {
	order := []string{
		"v0.9.16", "v0.11.5", "v1.0.2",
		"v1.1.0-alpha", "v1.1.0-alpha.1", "v1.1.0-beta.1", "v1.1.0-beta.2", "v1.1.0-beta.10", "v1.1.0-rc.1",
		"v1.1.0", "v1.1.1",
	}
	for i := range order {
		a, ok := semver(order[i])
		if !ok {
			t.Fatalf("%s is not a version", order[i])
		}
		for j := range order {
			b, _ := semver(order[j])
			if less(a, b) != (i < j) {
				t.Errorf("less(%s, %s) = %v", order[i], order[j], !(i < j))
			}
		}
	}
	for _, tag := range []string{"1.0.2", "v1.1", "v1.1.0-", "v1.1.0-beta..1", "v1.x.0", "main"} {
		if _, ok := semver(tag); ok {
			t.Errorf("%q taken as a version", tag)
		}
	}
}

// TestNewest: the newest tag is the newest release, as Go's @latest takes
// it; a beta only while there is no release, and a release after its betas.
func TestNewest(t *testing.T) {
	for _, c := range []struct {
		tags []string
		want string
	}{
		{[]string{"v1.1.10", "v1.2.0-beta.1", "v1.1.9", "main"}, "v1.1.10"},
		{[]string{"v1.2.0-beta.1", "v1.2.0", "v1.2.0-beta.2"}, "v1.2.0"},
		{[]string{"v1.1.10", "v1.1.9", "v1.1.11"}, "v1.1.11"},
		{[]string{"v1.0.0-beta.1", "v1.0.0-beta.2"}, "v1.0.0-beta.2"},
		{[]string{"main", "x"}, ""},
	} {
		if got := newest(c.tags); got != c.want {
			t.Errorf("newest(%q) = %q, want %q", c.tags, got, c.want)
		}
	}
}
