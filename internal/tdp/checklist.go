package tdp

import (
	"fmt"
	"regexp"
	"strings"
)

// Item is one item of the checklist: one gene.
type Item struct {
	ID      string // 編號
	Source  string // 出處
	Check   string // 檢查什麼
	Applies string // 適用
	HowTo   string // 怎麼驗
	Section string // the heading it is under
}

// Section is one heading of the checklist and its items.
type Section struct {
	Name  string
	Items []Item
}

var idPattern = regexp.MustCompile(`^[A-Za-z0-9]+(-[A-Za-z0-9]+)*-[0-9]+$`)

// Parse reads the items file: every table row under a ## heading whose first
// cell is an item number.
func Parse(md string) ([]Section, error) {
	var secs []Section
	seen := map[string]bool{}
	for n, line := range strings.Split(md, "\n") {
		if name, ok := strings.CutPrefix(line, "## "); ok {
			secs = append(secs, Section{Name: strings.TrimSpace(name)})
			continue
		}
		if !strings.HasPrefix(line, "|") || len(secs) == 0 {
			continue
		}
		cells := cells(line)
		if len(cells) == 0 || !idPattern.MatchString(cells[0]) {
			continue
		}
		if len(cells) != 5 {
			return nil, fmt.Errorf("line %d: %d cells, want 5", n+1, len(cells))
		}
		if seen[cells[0]] {
			return nil, fmt.Errorf("line %d: %s twice", n+1, cells[0])
		}
		seen[cells[0]] = true
		s := &secs[len(secs)-1]
		s.Items = append(s.Items, Item{
			ID: cells[0], Source: cells[1], Check: cells[2], Applies: cells[3], HowTo: cells[4], Section: s.Name,
		})
	}
	out := secs[:0]
	for _, s := range secs {
		if len(s.Items) > 0 {
			out = append(out, s)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no items")
	}
	return out, nil
}

// cells splits a table row on the pipes that are not escaped; an escaped
// pipe reads as a pipe.
func cells(line string) []string {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "|")
	if strings.HasSuffix(line, "|") && !strings.HasSuffix(line, "\\|") {
		line = line[:len(line)-1]
	}
	var out []string
	var b strings.Builder
	for i := 0; i < len(line); i++ {
		switch {
		case line[i] == '\\' && i+1 < len(line) && line[i+1] == '|':
			b.WriteByte('|')
			i++
		case line[i] == '|':
			out = append(out, strings.TrimSpace(b.String()))
			b.Reset()
		default:
			b.WriteByte(line[i])
		}
	}
	return append(out, strings.TrimSpace(b.String()))
}

// Batch is a run of whole sections, judged together by one judge.
type Batch struct {
	ID       string
	Sections []string
	Items    []Item
}

// Batches cuts the sections into about n batches of whole sections, in order:
// a batch takes the next section while it stays within a quarter over the
// average.
func Batches(secs []Section, n int) []Batch {
	total := 0
	for _, s := range secs {
		total += len(s.Items)
	}
	limit := (total + n - 1) / n * 5 / 4
	var out []Batch
	for _, s := range secs {
		if len(out) == 0 || len(out[len(out)-1].Items)+len(s.Items) > limit {
			out = append(out, Batch{ID: fmt.Sprintf("b%d", len(out)+1)})
		}
		b := &out[len(out)-1]
		b.Sections = append(b.Sections, s.Name)
		b.Items = append(b.Items, s.Items...)
	}
	return out
}
