package tdp

import (
	"fmt"
	"strings"
)

// ItemsFile is the checklist, whose source column names the section an item
// comes from.
const ItemsFile = "checklist/items-zh_TW.md"

// Section is a part of tdp, or of harness's project/, an app cites: its text
// and the checklist items that come from it.
type Section struct {
	Harness  bool // harness's project/, not tdp
	ID       string
	Path     string
	Text     string   // from its heading to the next heading; the whole file for a file
	Children []string // the IDs of the sections under it
	Items    []string // checklist rows whose source is it (for a file, any of its sections)
}

// Sections are the parts of tdp at the tag that ids cite: rules, principle,
// terms and section IDs (K4, P5, terms.3, panel.5), file names (confirm),
// template paths (project/install.sh, while tdp had them) or checklist items
// (K9-6). An item comes as its row, followed by the sections it names as its
// source. IDs written with the harness prefix (harness project.3, harness
// project/install.sh, harness project-13) are harness's project/; they come
// after tdp's.
func (s Source) Sections(tag string, ids []string) ([]Section, error) {
	var ours, theirs []string // tdp's, harness's
	for _, c := range ids {
		if inHarness, id := Cited(c); inHarness {
			theirs = append(theirs, id)
		} else {
			ours = append(ours, id)
		}
	}
	var out []Section
	if len(ours) > 0 {
		t, err := s.tree(tag)
		if err != nil {
			return nil, err
		}
		secs, err := t.sections(ours, func(id string) error {
			if HarnessIDs().Has(id) {
				return fmt.Errorf("%q is not in %s: from tdp v1.2.0-beta.1 on the family's repo templates are terminu-harness's, cited as %q", id, t.at, HarnessPrefix+id)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		out = append(out, secs...)
	}
	if len(theirs) > 0 {
		secs, err := harnessTree().sections(theirs, nil)
		if err != nil {
			return nil, err
		}
		for i := range secs {
			secs[i].Harness = true
		}
		out = append(out, secs...)
	}
	return out, nil
}

// sections are the parts of t that ids cite. An ID not in t is an error:
// miss may give a better one.
func (t tree) sections(ids []string, miss func(id string) error) ([]Section, error) {
	all, err := t.ids()
	if err != nil {
		return nil, err
	}
	items, err := t.file(t.items)
	if err != nil {
		return nil, err
	}
	var out []Section
	seen := map[string]bool{}
	for i := 0; i < len(ids); i++ {
		id := strings.TrimSpace(ids[i])
		if seen[id] {
			continue
		}
		seen[id] = true
		if srcs, ok := all.Items[id]; ok {
			out = append(out, Section{ID: id, Path: t.items, Text: itemRow(items, id)})
			ids = append(ids, srcs...)
			continue
		}
		if all.Templates[id] {
			b, err := t.file(id)
			if err != nil {
				return nil, err
			}
			out = append(out, Section{ID: id, Path: id, Text: string(b)})
			continue
		}
		if p, ok := all.Files[id]; ok {
			b, err := t.file(p)
			if err != nil {
				return nil, err
			}
			out = append(out, Section{ID: id, Path: p, Text: string(b), Items: itemRows(items, func(src string) bool {
				return strings.HasPrefix(src, id+".")
			})})
			continue
		}
		p, ok := all.Paths[id]
		if !ok {
			if miss != nil {
				if err := miss(id); err != nil {
					return nil, err
				}
			}
			return nil, fmt.Errorf("%q is not in %s", id, t.at)
		}
		b, err := t.file(p)
		if err != nil {
			return nil, err
		}
		text, children := sectionText(b, id)
		out = append(out, Section{ID: id, Path: p, Text: text, Children: children, Items: itemRows(items, func(src string) bool {
			return src == id
		})})
	}
	return out, nil
}

// sectionText is the text from the heading of id to the next heading, and the
// IDs of the headings under it before one at its level or above.
func sectionText(b []byte, id string) (string, []string) {
	lines := strings.Split(string(b), "\n")
	start, level := -1, 0
	inCode := false
	var text []string
	var children []string
	for i, line := range lines {
		if strings.HasPrefix(line, "```") {
			inCode = !inCode
		}
		heading := !inCode && !strings.HasPrefix(line, "```") && headingRE.MatchString(line)
		hl := 0
		if heading {
			hl = len(line) - len(strings.TrimLeft(line, "#"))
		}
		switch {
		case start < 0:
			if m := headingRE.FindStringSubmatch(line); heading && m[1] == id {
				start, level = i, hl
				text = append(text, line)
			}
		case heading && hl <= level:
			return strings.TrimRight(strings.Join(text, "\n"), "\n"), children
		case heading:
			if m := headingRE.FindStringSubmatch(line); idRE.MatchString(m[1]) {
				children = append(children, m[1])
			}
			if len(children) == 0 {
				text = append(text, line)
			}
		case len(children) == 0:
			text = append(text, line)
		}
	}
	return strings.TrimRight(strings.Join(text, "\n"), "\n"), children
}

// itemRows are the checklist's table rows with a source that passes keep.
func itemRows(items []byte, keep func(src string) bool) []string {
	var rows []string
	for _, line := range strings.Split(string(items), "\n") {
		if !strings.HasPrefix(line, "| ") {
			continue
		}
		cols := strings.Split(line, "|")
		if len(cols) <= 3 {
			continue
		}
		for _, src := range sources(cols[2]) {
			if keep(src) {
				rows = append(rows, line)
				break
			}
		}
	}
	return rows
}

// itemRow is the checklist's table row of the item id.
func itemRow(items []byte, id string) string {
	for _, line := range strings.Split(string(items), "\n") {
		if strings.HasPrefix(line, "| "+id+" |") {
			return line
		}
	}
	return ""
}

// sources are the sections a source column names: several are separated by
// 、, and what follows （ or " (" only narrows one down (（圖）, (diagram)).
func sources(col string) []string {
	var out []string
	for _, s := range strings.Split(col, "、") {
		if i := strings.Index(s, "（"); i >= 0 {
			s = s[:i]
		}
		if i := strings.Index(s, " ("); i >= 0 {
			s = s[:i]
		}
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}
