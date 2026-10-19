package scan

import (
	"os"
	"testing"

	"github.com/vulcanshen/terminu-harness/internal/tdp"
)

// TestARealApp scans the app at TERMINU_HARNESS__SCAN_REPO against tdp
// v0.7.4 from the terminu clone at TERMINU_HARNESS__TDP_REPO, and says how
// many items the scan judges not applicable.
func TestARealApp(t *testing.T) {
	app, terminu := os.Getenv("TERMINU_HARNESS__SCAN_REPO"), os.Getenv("TERMINU_HARNESS__TDP_REPO")
	if app == "" || terminu == "" {
		t.Skip("no app or terminu clone to try")
	}
	r, err := Scan(app, "v0.1.13")
	if err != nil || r == nil {
		t.Fatalf("scan: %v", err)
	}
	b, err := tdp.Source{Repo: terminu}.File("v0.7.4", tdp.ItemsPath)
	if err != nil {
		t.Fatal(err)
	}
	secs, _ := tdp.Parse(string(b))
	n, total := 0, 0
	for _, s := range secs {
		for _, it := range s.Items {
			total++
			if na, _ := r.NotApplicable(it.Applies); na {
				n++
			}
		}
	}
	t.Logf("custom %v; not applicable by the scan: %d of %d", r.Custom, n, total)
}
