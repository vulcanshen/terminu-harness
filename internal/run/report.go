package run

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vulcanshen/terminu-harness/internal/tdp"
)

// ReportFile is the report card, for people.
const ReportFile = "report.md"

// Finish writes the report card and returns its summary. Items no judge got
// to count as 判不了, with that as the reason.
func (r *Run) Finish() (string, error) {
	judged, err := r.Judged()
	if err != nil {
		return "", err
	}
	gaps, err := r.Gaps()
	if err != nil {
		return "", err
	}
	verdictOf := func(it tdp.Item) Verdict {
		if v, ok := judged[it.ID]; ok {
			return v
		}
		return Verdict{ID: it.ID, Verdict: CannotTell, Reason: "沒有 judge 判到這一條"}
	}
	count := map[string]int{}
	unjudged := 0
	var violated []string
	for _, it := range r.items {
		v := verdictOf(it).Verdict
		count[v]++
		if v == Violates {
			violated = append(violated, it.ID)
		}
		if _, ok := judged[it.ID]; !ok {
			unjudged++
		}
	}
	genes := Genes(r.Meta.TDP, count[Meets])
	outcome, unlisted := r.Outcome(violated)

	var s strings.Builder
	fmt.Fprintf(&s, "# %s 的成績單\n\n**%s**　%s\n\n", r.Meta.App, outcome, genes)
	s.WriteString("- 判斷線：沒有違反是「通過」；違反都列在 app 的 `docs/<app>-terminu-fix.md`（在等 tdp 或 terminu-go）是「待修」；" +
		"有違反沒列進去是「不通過」。gene 數量看進度。\n")
	if len(unlisted) > 0 {
		fmt.Fprintf(&s, "- 沒列進 fix 文件的違反（%d）：%s\n", len(unlisted), strings.Join(unlisted, "、"))
	}
	fmt.Fprintf(&s, "- app：`%s`\n", r.Meta.Command)
	if r.Meta.Repo != "" {
		fmt.Fprintf(&s, "- repo：`%s`\n", r.Meta.Repo)
	} else {
		s.WriteString("- repo：沒有，要讀程式或讀檔案的條目判「判不了」\n")
	}
	fmt.Fprintf(&s, "- 照 tdp %s 判，清單 %d 條；%s 開始，%s 寫成\n", r.Meta.TDP, r.Total(),
		r.Meta.Started.Format("2006-01-02 15:04"), time.Now().Format("2006-01-02 15:04"))
	s.WriteString("- judge 認識到的 app：[profile.md](profile.md)；證據的檔案在 `evidence/`\n\n")

	s.WriteString("## 判定\n\n| 判定 | 條數 |\n|---|---|\n")
	for _, v := range Verdicts {
		label := v
		switch v {
		case Meets:
			label = "符合（gene）"
		case CannotTell:
			if unjudged > 0 {
				label = fmt.Sprintf("判不了（含沒判到的 %d 條）", unjudged)
			}
		}
		fmt.Fprintf(&s, "| %s | %d |\n", label, count[v])
	}
	fmt.Fprintf(&s, "| 合計 | %d |\n\n", r.Total())

	for _, v := range []string{Violates, Departs, CannotTell} {
		fmt.Fprintf(&s, "## %s（%d）\n\n", v, count[v])
		if count[v] == 0 {
			s.WriteString("沒有。\n\n")
			continue
		}
		s.WriteString("| 編號 | 檢查什麼 | 為什麼 | 證據 |\n|---|---|---|---|\n")
		for _, it := range r.items {
			if vd := verdictOf(it); vd.Verdict == v {
				fmt.Fprintf(&s, "| %s | %s | %s | %s |\n", it.ID, cell(it.Check), cell(vd.Reason), r.links(vd.Evidence))
			}
		}
		s.WriteString("\n")
	}

	fmt.Fprintf(&s, "## tdp 沒規範的東西（%d）\n\n", len(gaps))
	if len(gaps) == 0 {
		s.WriteString("沒有。\n\n")
	} else {
		s.WriteString("app 做了、tdp 還沒規範的東西：不算 gene 也不扣，回報成 tdp 的缺漏。\n\n| 什麼 | 在哪 | 證據 |\n|---|---|---|\n")
		for _, g := range gaps {
			fmt.Fprintf(&s, "| %s | %s | %s |\n", cell(g.What), cell(g.Where), r.links(g.Evidence))
		}
		s.WriteString("\n")
	}

	s.WriteString("## 全部條目\n")
	section := ""
	for _, it := range r.items {
		if it.Section != section {
			section = it.Section
			fmt.Fprintf(&s, "\n### %s\n\n| 編號 | 判定 | 為什麼 | 證據 |\n|---|---|---|---|\n", section)
		}
		vd := verdictOf(it)
		fmt.Fprintf(&s, "| %s | %s | %s | %s |\n", it.ID, vd.Verdict, cell(vd.Reason), r.links(vd.Evidence))
	}
	if err := os.WriteFile(filepath.Join(r.Dir, ReportFile), []byte(s.String()), 0o644); err != nil {
		return "", err
	}

	var sum strings.Builder
	fmt.Fprintf(&sum, "%s：%s，%s\n", r.Meta.App, outcome, genes)
	if len(unlisted) > 0 {
		fmt.Fprintf(&sum, "  沒列進 fix 文件的違反：%s\n", strings.Join(unlisted, "、"))
	}
	for _, v := range Verdicts {
		fmt.Fprintf(&sum, "  %s %d\n", v, count[v])
	}
	if unjudged > 0 {
		fmt.Fprintf(&sum, "  （判不了裡有 %d 條沒有 judge 判到）\n", unjudged)
	}
	fmt.Fprintf(&sum, "  tdp 沒規範的東西 %d 項\n成績單：%s\n", len(gaps), filepath.Join(r.Dir, ReportFile))
	return sum.String(), nil
}

// cell makes text fit in one table cell.
func cell(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "|", `\|`)
	s = strings.ReplaceAll(s, "\r\n", "<br>")
	return strings.ReplaceAll(s, "\n", "<br>")
}

// links points to the evidence files.
func (r *Run) links(ids []string) string {
	var out []string
	for _, id := range ids {
		if strings.Contains(id, "/") { // carried from an earlier run: a path to it
			out = append(out, fmt.Sprintf("[%s](%s)", strings.TrimSuffix(filepath.Base(id), filepath.Ext(id)), id))
			continue
		}
		m, _ := filepath.Glob(filepath.Join(r.Dir, "evidence", id+".*"))
		if len(m) == 0 {
			out = append(out, id)
			continue
		}
		out = append(out, fmt.Sprintf("[%s](evidence/%s)", id, filepath.Base(m[0])))
	}
	return strings.Join(out, " ")
}
