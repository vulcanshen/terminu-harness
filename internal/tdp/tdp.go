// Package tdp gets the files of a version of terminu design (tdp).
package tdp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// GitHub is where tdp lives.
const GitHub = "vulcanshen/terminu"

// Source gets tdp's files at a tag: from a local clone of terminu when Repo is
// set, from GitHub otherwise; Cache keeps what it got.
type Source struct {
	Repo  string
	Cache string
}

var client = &http.Client{Timeout: 30 * time.Second}

// Tag is version as a tag: "0.7.4" and "v0.7.4" are both v0.7.4.
func Tag(version string) string {
	if version == "" || strings.HasPrefix(version, "v") {
		return version
	}
	return "v" + version
}

// Latest is the newest vX.Y.Z tag.
func (s Source) Latest() (string, error) {
	var names []string
	if s.Repo != "" {
		out, err := exec.Command("git", "-C", s.Repo, "tag", "--list", "v*").Output()
		if err != nil {
			return "", fmt.Errorf("tags of %s: %w", s.Repo, err)
		}
		names = strings.Fields(string(out))
	} else {
		b, err := get("https://api.github.com/repos/" + GitHub + "/tags?per_page=100")
		if err != nil {
			return "", err
		}
		var tags []struct{ Name string }
		if err := json.Unmarshal(b, &tags); err != nil {
			return "", fmt.Errorf("tags of %s: %w", GitHub, err)
		}
		for _, t := range tags {
			names = append(names, t.Name)
		}
	}
	best, bestV := "", []int(nil)
	for _, n := range names {
		v, ok := semver(n)
		if ok && (bestV == nil || less(bestV, v)) {
			best, bestV = n, v
		}
	}
	if best == "" {
		return "", fmt.Errorf("no vX.Y.Z tag")
	}
	return best, nil
}

func semver(tag string) ([]int, bool) {
	parts := strings.Split(strings.TrimPrefix(tag, "v"), ".")
	if len(parts) != 3 || !strings.HasPrefix(tag, "v") {
		return nil, false
	}
	v := make([]int, 3)
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, false
		}
		v[i] = n
	}
	return v, true
}

func less(a, b []int) bool {
	for i := range a {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}

// File is the file at path in tdp at the tag.
func (s Source) File(tag, path string) ([]byte, error) {
	tag = Tag(tag)
	if _, ok := semver(tag); !ok {
		return nil, fmt.Errorf("%q is not a vX.Y.Z version", tag)
	}
	path = strings.TrimPrefix(filepath.ToSlash(filepath.Clean(path)), "/")
	if strings.HasPrefix(path, "..") {
		return nil, fmt.Errorf("%q is outside tdp", path)
	}
	cached := ""
	if s.Cache != "" {
		cached = filepath.Join(s.Cache, tag, filepath.FromSlash(path))
		if b, err := os.ReadFile(cached); err == nil {
			return b, nil
		}
	}
	var b []byte
	var err error
	if s.Repo != "" {
		b, err = exec.Command("git", "-C", s.Repo, "show", tag+":"+path).Output()
		if err != nil {
			err = fmt.Errorf("%s:%s in %s: %w", tag, path, s.Repo, err)
		}
	} else {
		b, err = get("https://raw.githubusercontent.com/" + GitHub + "/" + tag + "/" + path)
	}
	if err != nil {
		return nil, err
	}
	if cached != "" {
		if os.MkdirAll(filepath.Dir(cached), 0o755) == nil {
			_ = os.WriteFile(cached, b, 0o644)
		}
	}
	return b, nil
}

func get(url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", url, resp.Status)
	}
	return b, nil
}
