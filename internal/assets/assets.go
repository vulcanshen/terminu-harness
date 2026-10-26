// Package assets reads and checks the files the harness keeps in an app's
// repo, under harness-assets/: idea.md (what the user asked for), spec.yaml
// (each request matched to tdp), work.yaml (how each was built and what backs
// it) and review.yaml (the last review).
package assets

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Dir is the folder in the app's repo.
const Dir = "harness-assets"

// The files in Dir.
const (
	IdeaFile   = "idea.md"
	SpecFile   = "spec.yaml"
	WorkFile   = "work.yaml"
	ReviewFile = "review.yaml"
)

// Spec is spec.yaml: the requests in idea.md matched to a version of tdp. Its
// ID goes up by one each time idea.md or the tdp version changes.
type Spec struct {
	ID    int        `yaml:"id"`
	From  string     `yaml:"from"` // the tdp version before this spec; 0 for the first
	To    string     `yaml:"to"`   // the tdp version this spec follows
	Tasks []SpecTask `yaml:"tasks"`
}

// SpecTask is one request and the tdp it follows.
type SpecTask struct {
	ID       int      `yaml:"id"`   // unique in the app, never reused
	Spec     int      `yaml:"spec"` // the spec that last added or changed it
	Request  string   `yaml:"request"`
	TDP      []string `yaml:"tdp"`      // "tdp K4", "tdp panel.5", "tdp confirm"
	Canceled string   `yaml:"canceled"` // why it was taken out
	Remark   []Remark `yaml:"remark"`
}

// Remark is a design decision or a gap in tdp. Written as plain text, it is a
// note; a gap reported to tdp or terminu-go also names the issue.
type Remark struct {
	Note  string `yaml:"note"`
	Issue string `yaml:"issue"` // e.g. "tdp #12", "tgo #5"
}

func (r *Remark) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind == yaml.ScalarNode {
		r.Note = n.Value
		return nil
	}
	type plain Remark
	return n.Decode((*plain)(r))
}

// Work is work.yaml: how each task was built.
type Work struct {
	Tasks []WorkTask `yaml:"tasks"`
}

// WorkTask is how a task was built and what backs it.
type WorkTask struct {
	ID             int       `yaml:"id"`
	Spec           int       `yaml:"spec"` // the spec it was built to
	Implementation string    `yaml:"implementation"`
	Test           []Backing `yaml:"test"`
	Canceled       bool      `yaml:"canceled"`
}

// Backing is one thing that backs a task: a test's name, a commit (when no
// test can show it), or the user's own check.
type Backing struct {
	Test   string     `yaml:"test"`
	Commit string     `yaml:"commit"`
	User   *UserCheck `yaml:"user"`
}

func (b *Backing) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind == yaml.ScalarNode {
		b.Test = n.Value
		return nil
	}
	type plain Backing
	return n.Decode((*plain)(b))
}

// UserCheck is how the user checked it, and when they last did.
type UserCheck struct {
	Date string `yaml:"date"` // YYYY-MM-DD
	How  string `yaml:"how"`
}

// Review is review.yaml: the last review, of every task in work.yaml.
type Review struct {
	Spec  int          `yaml:"spec"`
	Tasks []ReviewTask `yaml:"tasks"`
}

// ReviewTask is whether a task follows tdp.
type ReviewTask struct {
	ID       int    `yaml:"id"`
	Accepted bool   `yaml:"accepted"`
	Reason   string `yaml:"reason"` // starts with 違反, 沒有測試 or tdp 不清楚
}

// Files is what is in Dir; a file that is not there is nil (Idea is "").
type Files struct {
	Idea   string
	Spec   *Spec
	Work   *Work
	Review *Review
}

// Load reads Dir in the repo. A missing file is left nil; a file that is
// there but does not parse is an error.
func Load(repo string) (Files, error) {
	dir := filepath.Join(repo, Dir)
	var f Files
	if b, err := os.ReadFile(filepath.Join(dir, IdeaFile)); err == nil {
		f.Idea = string(b)
	} else if !os.IsNotExist(err) {
		return f, err
	}
	var err error
	if f.Spec, err = load[Spec](dir, SpecFile); err != nil {
		return f, err
	}
	if f.Work, err = load[Work](dir, WorkFile); err != nil {
		return f, err
	}
	if f.Review, err = load[Review](dir, ReviewFile); err != nil {
		return f, err
	}
	return f, nil
}

func load[T any](dir, name string) (*T, error) {
	b, err := os.ReadFile(filepath.Join(dir, name))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	v := new(T)
	if err := yaml.Unmarshal(b, v); err != nil {
		return nil, fmt.Errorf("%s/%s: %w", Dir, name, err)
	}
	return v, nil
}
