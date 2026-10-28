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

// Latest is the newest vX.Y.Z tag; a pre-release (v1.1.0-beta.1) counts, and
// comes before its release.
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
	best, bestV := "", version{}
	for _, n := range names {
		v, ok := semver(n)
		if ok && (best == "" || less(bestV, v)) {
			best, bestV = n, v
		}
	}
	if best == "" {
		return "", fmt.Errorf("no vX.Y.Z tag")
	}
	return best, nil
}

// version is a vX.Y.Z tag and its pre-release, the dot-separated part after
// a "-" (v1.1.0-beta.1).
type version struct {
	core [3]int
	pre  []string
}

func semver(tag string) (version, bool) {
	var v version
	core, pre, hasPre := strings.Cut(strings.TrimPrefix(tag, "v"), "-")
	parts := strings.Split(core, ".")
	if len(parts) != 3 || !strings.HasPrefix(tag, "v") {
		return v, false
	}
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return v, false
		}
		v.core[i] = n
	}
	if hasPre {
		v.pre = strings.Split(pre, ".")
		for _, id := range v.pre {
			if id == "" {
				return v, false
			}
		}
	}
	return v, true
}

// less is semver's precedence: a pre-release comes before its release; the
// parts of pre-releases compare one by one, numbers as numbers and before
// words, and the shorter comes first when one is the start of the other.
func less(a, b version) bool {
	for i := range a.core {
		if a.core[i] != b.core[i] {
			return a.core[i] < b.core[i]
		}
	}
	if len(a.pre) == 0 || len(b.pre) == 0 {
		return len(a.pre) > len(b.pre)
	}
	for i := 0; i < len(a.pre) && i < len(b.pre); i++ {
		x, y := a.pre[i], b.pre[i]
		if x == y {
			continue
		}
		nx, ex := strconv.Atoi(x)
		ny, ey := strconv.Atoi(y)
		switch {
		case ex == nil && ey == nil:
			return nx < ny
		case ex == nil || ey == nil:
			return ex == nil
		}
		return x < y
	}
	return len(a.pre) < len(b.pre)
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
