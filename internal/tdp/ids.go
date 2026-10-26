package tdp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// IDs is what an app can cite in tdp at a version: every ID written in a
// heading of principle, components and project (rules K4, principle P5, terms
// terms.3, sections panel.5), every component and project file by name
// (confirm, input), which stands for the whole file, every project template by
// its path (project/install.sh), and every checklist item by its ID (K9-6),
// which points at the item's source.
type IDs struct {
	Sections  map[string]string   // ID -> the heading's title
	Paths     map[string]string   // ID -> the path of the file it is in
	Files     map[string]string   // name -> the file's path
	Templates map[string]bool     // path of a project template
	Items     map[string][]string // checklist item ID -> its sources
}

// Has reports whether id is a section, a file, a template or a checklist item
// in tdp.
func (ids IDs) Has(id string) bool {
	_, s := ids.Sections[id]
	_, f := ids.Files[id]
	_, i := ids.Items[id]
	return s || f || i || ids.Templates[id]
}

// Files lists the paths of every file in tdp at the tag.
func (s Source) Files(tag string) ([]string, error) {
	tag = Tag(tag)
	if _, ok := semver(tag); !ok {
		return nil, fmt.Errorf("%q is not a vX.Y.Z version", tag)
	}
	cached := ""
	if s.Cache != "" {
		cached = filepath.Join(s.Cache, tag, ".tree")
		if b, err := os.ReadFile(cached); err == nil {
			return strings.Fields(string(b)), nil
		}
	}
	var paths []string
	if s.Repo != "" {
		out, err := exec.Command("git", "-C", s.Repo, "ls-tree", "-r", "--name-only", tag).Output()
		if err != nil {
			return nil, fmt.Errorf("files of %s in %s: %w", tag, s.Repo, err)
		}
		paths = strings.Fields(string(out))
	} else {
		b, err := get("https://api.github.com/repos/" + GitHub + "/git/trees/" + tag + "?recursive=1")
		if err != nil {
			return nil, err
		}
		var tree struct {
			Tree []struct{ Path, Type string }
		}
		if err := json.Unmarshal(b, &tree); err != nil {
			return nil, fmt.Errorf("files of %s at %s: %w", GitHub, tag, err)
		}
		for _, t := range tree.Tree {
			if t.Type == "blob" {
				paths = append(paths, t.Path)
			}
		}
	}
	sort.Strings(paths)
	if cached != "" {
		if os.MkdirAll(filepath.Dir(cached), 0o755) == nil {
			_ = os.WriteFile(cached, []byte(strings.Join(paths, "\n")+"\n"), 0o644)
		}
	}
	return paths, nil
}

// IDs reads the IDs in tdp at the tag from the headings of the Traditional
// Chinese files, the original.
func (s Source) IDs(tag string) (IDs, error) {
	paths, err := s.Files(tag)
	if err != nil {
		return IDs{}, err
	}
	ids := IDs{Sections: map[string]string{}, Paths: map[string]string{}, Files: map[string]string{},
		Templates: map[string]bool{}, Items: map[string][]string{}}
	for _, p := range paths {
		if template(p) {
			ids.Templates[p] = true
		}
		if p == ItemsFile {
			b, err := s.File(tag, p)
			if err != nil {
				return IDs{}, err
			}
			ids.Items = items(b)
		}
		if !citable(p) {
			continue
		}
		b, err := s.File(tag, p)
		if err != nil {
			return IDs{}, err
		}
		for id, title := range headingIDs(b) {
			ids.Sections[id] = title
			ids.Paths[id] = p
		}
		if name := fileName(p); name != "" {
			ids.Files[name] = p
		}
	}
	return ids, nil
}

// citable is a file whose headings carry IDs.
func citable(p string) bool {
	if !strings.HasSuffix(p, "-zh_TW.md") {
		return false
	}
	for _, dir := range []string{"principle/", "components/", "project/"} {
		if strings.HasPrefix(p, dir) {
			return true
		}
	}
	return false
}

// template is a project template: a file of project other than its README.
func template(p string) bool {
	return strings.HasPrefix(p, "project/") && !strings.HasSuffix(p, ".md")
}

// fileName is the name a component or project file is cited by: the file's
// own name, or its folder's for a README (components, input, project).
func fileName(p string) string {
	if !strings.HasPrefix(p, "components/") && !strings.HasPrefix(p, "project/") {
		return ""
	}
	base := strings.TrimSuffix(path.Base(p), "-zh_TW.md")
	if base == "README" {
		return path.Base(path.Dir(p))
	}
	return base
}

var (
	headingRE = regexp.MustCompile(`^#{2,4} (\S+)\s+(.*)$`)
	idRE      = regexp.MustCompile(`^(?:[A-Z][0-9]+|[a-z][a-z0-9-]*\.[0-9]+)$`)
	itemRE    = regexp.MustCompile(`^\| ([A-Za-z][A-Za-z0-9-]*-[0-9]+) \|([^|]*)\|`)
)

// items are the checklist's item IDs and the sources each names.
func items(b []byte) map[string][]string {
	out := map[string][]string{}
	for _, line := range strings.Split(string(b), "\n") {
		if m := itemRE.FindStringSubmatch(line); m != nil {
			out[m[1]] = sources(m[2])
		}
	}
	return out
}

// headingIDs are the IDs at the start of the headings, outside code blocks.
func headingIDs(b []byte) map[string]string {
	ids := map[string]string{}
	inCode := false
	sc := bufio.NewScanner(bytes.NewReader(b))
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "```") {
			inCode = !inCode
			continue
		}
		if inCode {
			continue
		}
		m := headingRE.FindStringSubmatch(line)
		if m != nil && idRE.MatchString(m[1]) {
			ids[m[1]] = m[2]
		}
	}
	return ids
}
