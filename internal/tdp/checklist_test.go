package tdp

import (
	"os"
	"testing"
)

const sample = "# items\n\n## rules — K Core key\n\n| 編號 | 出處 | 檢查什麼 | 適用 | 怎麼驗 |\n|---|---|---|---|---|\n" +
	"| K1-1 | K1 | a | 都要 | press `j\\|k` |\n| K1-2 | K1 | b | 表單 | c |\n\n## input/file-picker\n\n" +
	"| 編號 | 出處 | 檢查什麼 | 適用 | 怎麼驗 |\n|---|---|---|---|---|\n| file-picker-51 | file-picker §metadata | d | finder、file-picker | e |\n"

func TestParse(t *testing.T) {
	secs, err := Parse(sample)
	if err != nil {
		t.Fatal(err)
	}
	if len(secs) != 2 || len(secs[0].Items) != 2 || len(secs[1].Items) != 1 {
		t.Fatalf("%+v", secs)
	}
	if got := secs[0].Items[0].HowTo; got != "press `j|k`" {
		t.Errorf("escaped pipe: %q", got)
	}
	if got := secs[1].Items[0]; got.ID != "file-picker-51" || got.Section != "input/file-picker" {
		t.Errorf("%+v", got)
	}
}

// TestParseTheChecklist reads tdp's own items when a clone of terminu is at
// hand (TERMINU_HARNESS__TDP_REPO): the count is what its README says.
func TestParseTheChecklist(t *testing.T) {
	repo := os.Getenv("TERMINU_HARNESS__TDP_REPO")
	if repo == "" {
		t.Skip("no terminu clone")
	}
	b, err := Source{Repo: repo}.File("v0.8.1", ItemsPath)
	if err != nil {
		t.Fatal(err)
	}
	secs, err := Parse(string(b))
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, s := range secs {
		total += len(s.Items)
	}
	if total != 1391 {
		t.Errorf("%d items in v0.8.1, want 1391", total)
	}
	batches := Batches(secs, 10)
	sum := 0
	for _, b := range batches {
		sum += len(b.Items)
		t.Logf("%s %d %v", b.ID, len(b.Items), b.Sections)
	}
	if sum != total || len(batches) < 8 || len(batches) > 12 {
		t.Errorf("%d batches, %d items", len(batches), sum)
	}
}
