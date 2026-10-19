package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/vulcanshen/terminu-harness/internal/appgit"
	"github.com/vulcanshen/terminu-harness/internal/run"
	"github.com/vulcanshen/terminu-harness/internal/scan"
)

func (sv *Server) judgeScan(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	r, err := sv.openRun(req.GetString("run", ""))
	if err != nil || r == nil {
		return fail("run: %v", err)
	}
	if r.Meta.Repo == "" {
		return mcp.NewToolResultText("no repo: nothing to scan, every item is the judges'"), nil
	}
	res, err := scan.Scan(r.Meta.Repo, appgit.Tgo(r.Meta.Repo))
	if err != nil {
		return fail("scan: %v", err)
	}
	if res == nil {
		return mcp.NewToolResultText("not written with terminu-go: nothing to scan, every item is the judges'"), nil
	}
	b, _ := json.MarshalIndent(res, "", "  ")
	if err := os.WriteFile(filepath.Join(r.Dir, "scan.json"), b, 0o644); err != nil {
		return fail("%v", err)
	}
	judged, err := r.Judged()
	if err != nil {
		return fail("%v", err)
	}
	var vs []run.Verdict
	for _, it := range r.Items() {
		if _, ok := judged[it.ID]; ok {
			continue
		}
		if na, missing := res.NotApplicable(it.Applies); na {
			vs = append(vs, run.Verdict{ID: it.ID, Verdict: run.NotApply,
				Reason: "程式掃描：app 沒有「" + strings.Join(missing, "」「") + "」：原始碼沒用到 terminu-go 做它的元件，也沒有自己寫的 panel 或 popup"})
		}
	}
	kept, problems := r.AddVerdicts(vs)

	var has, lacks []string
	for f, ok := range res.Has {
		if ok {
			has = append(has, f)
		} else {
			lacks = append(lacks, f)
		}
	}
	sort.Strings(has)
	sort.Strings(lacks)
	var out strings.Builder
	fmt.Fprintf(&out, "scanned %s (terminu-go %s), kept in scan.json\n", r.Meta.Repo, res.Tgo)
	fmt.Fprintf(&out, "has: %s\nhas not: %s\n", strings.Join(has, "、"), strings.Join(lacks, "、"))
	if len(res.Custom) > 0 {
		fmt.Fprintf(&out, "bodies or popups of its own in %s: nothing judged, the judges decide what applies\n", strings.Join(res.Custom, ", "))
	}
	fmt.Fprintf(&out, "judged %d items 不適用; judge_status lists what is left\n", kept)
	if len(problems) > 0 {
		out.WriteString(strings.Join(problems, "\n") + "\n")
	}
	return mcp.NewToolResultText(out.String()), nil
}
