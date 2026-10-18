# terminu-harness 開發者備忘

terminu 的 harness：一個 Claude Code plugin，照 terminu design（tdp）用 terminu-go 做 TUI app，做完照 tdp 的清單檢查。harness 只放流程，不寫規定；規定一律在 tdp。

## 運作方式

- **入口** `/terminu`（`skills/terminu`，frontmatter `name: terminu`，所以裸的 `/terminu` 叫得到；實測過）：只問要做新 app、繼續做、還是判 app，再用 Skill 叫 `terminu:new-app` 或 `terminu:judge`。
- **new-app**（`skills/new-app`，主 session 照著做）：準備（讀 tdp、terminu-go 最新 tag）→ 一次一題問需求 → `docs/spec.md` 給使用者確認 → 照 tdp project 範本開 repo → 照里程碑一塊一塊寫，每塊 `make check`、用 harness 的終端機開起來看、給使用者看、commit → 卡住叫 `terminu:report-gap` 寫 issue（同意才發），程式標 `TODO(tdp #N)`／`TODO(tgo #N)` 先用暫時的做法 → 做完叫 `terminu:judge`。進度在 app repo 的 `.local/terminu/progress.md`，「繼續」從它接著做，並查 issue 與新 tag。
- **report-gap**（`skills/report-gap`）：判斷發到 tdp 還是 terminu-go、查重複、照固定格式寫草稿、使用者同意才 `gh issue create`。

judge 的流程：

```
/terminu:judge（skill，主 session）
  ├─ judge_start ──────────────► run 目錄：清單副本、batches/b1..b10.md、run.json
  ├─ judge agent「認識 app」 ──► open / keys / screen … ──► profile.md
  ├─ judge agent × 10「判一批」► open / keys / styles … ──► verdicts.jsonl、gaps.jsonl、evidence/
  ├─ judge_status ─────────────► 沒判到的再派一次
  └─ judge_finish ─────────────► report.md（gene 由程式數）
```

- **MCP server**（`cmd/terminu-harness`，子指令 `mcp`）：plugin 的 `.mcp.json` 跑 `scripts/mcp.sh`，server 名 `harness`，工具在
  Claude Code 裡叫 `mcp__plugin_terminu_harness__<工具>`。
- **終端機**（`internal/term`、`internal/vt`）：`creack/pty` 開 app，`vt` 是從 terminu-go 複製的模擬器。讀 PTY 與處理畫面分兩個
  goroutine：一個只管讀（app 寫 PTY 不會被卡住），一個把讀到的全部一次餵給畫面，再用每一列的 checksum 找出變了哪幾列。錄影記的
  就是這些變化，時間取那一批第一塊輸出讀到的時刻；6ms 內的變化併成一格。
- **等畫面穩定**：從送鍵之後、或之後最後一次變化起算，150ms 沒變就回傳（`open` 先等到第一次變化再等 300ms）。從上一次變化
  起算會在 app 還沒反應前就回傳。
- **證據**：session 帶 `run` 時，`screen`、`keys`、`styles`、`image` 每次的結果都存成 `evidence/eNNNN.*`，回傳編號給 judge 引用。
- **清單**（`internal/tdp`）：解析 `checklist/items-zh_TW.md` 的表格（`\|` 是表格裡的 `|`），照章切批，不拆章。
- **成績單**（`internal/run`）：判定、tdp 沒規範的東西都是 append 的 JSON lines，同一條後寫的蓋掉先寫的；`report.md` 由
  `judge_finish` 從它們寫出來。

## 設計決定

- **模擬器複製 terminu-go 的 `vt.go`**（GPL，檔頭註明出處與 commit），不用 charmbracelet 的 x/vt：icon 與 powerline 各佔幾格是
  vt.go 本來就有的設定，量法跟家族 app 同一套，判 E5 才準；x/vt 把 icon 一律算一格。改了的地方：寬度跟著每個畫面、數折行與
  捲動、拿掉 render（harness 直接讀格子）。terminu-go 改 vt.go 時要手動跟。
- **看畫面預設只回文字**，顏色另外指定範圍查：demou 一個 120×40 的畫面樣式有 269 段，每次全回多約 5k token。
- **錄影只在畫面變的時候記**，不固定間隔拍：不漏格、不重複、資料量小。
- **字型自己帶**（Go Mono、Symbols Nerd Font Mono），框線與方塊字元自己畫：跑 judge 的機器不一定有字型；家族的粗框 `━┃┏` 不在
  Go Mono 裡。
- **judge agent 沒有 Write、Edit**：判定、profile 都經 MCP 工具記，agent 從頭到尾只讀。
- **環境變數**：開 app 時拿掉真終端機的 `TERM`、`*__ICON_WIDTH`、`*__POWERLINE_WIDTH`，app 才會探測 harness 這個終端機；judge 要測
  E5-5 這種條目時自己用 `env` 設。
- **時間的暫時判法**（33ms、splash 的 tick 只看順序）寫在 `agents/judge.md`，對應 tdp 那邊的 HN1、HN2；tdp 定案後改。

## 已知的牆與未做

- 只支援 macOS、Linux（`creack/pty`，沒有 ConPTY）。
- 不送滑鼠事件。
- Bubble Tea v1 一秒最多畫 60 格，tick 比一格短的動畫從畫面上看不出每個 tick（HN1）。
- 使用者要有 Go：第一次啟動時 `scripts/mcp.sh` 從原始碼建 binary，可能超過 Claude Code 等 MCP server 的時間。
- 等 issue 不是一直開著的 session：Claude Code 沒有可靠地等好幾天的方式（Monitor 每次最多 30 分鐘、`/loop` 要 session 開著、routine 沒有 issue 事件），所以靠使用者回來打 `/terminu` 選「繼續」。
- 還沒做：align（tdp 新版時重判、列出要修的）；judge 找到的 tdp 沒規範的東西還沒自動轉成 report-gap 的草稿。

## 建置與開發

- `make build`：建到 `.local/bin/terminu-harness`，`scripts/mcp.sh` 優先用它。
- `make test`；設 `TERMINU_HARNESS__TDP_REPO=~/Documents/sideproj/terminu` 時也驗真的清單（v0.7.4 要是 1302 條）。
- server 的測試開 `internal/server/testdata/fixture`（測試時才建的小程式，不靠外部套件）代替真的 app：它照家族 app 的方式探測
  icon 寬度、顯示收到的鍵、按 `a` 播 5 格動畫、按 `w` 印超過邊界、按 `q` 結束。驗 MCP 工具用它，不拿家族成員來試。
- `make check`：gofmt、go vet、測試、`claude plugin validate`。
- 不裝 plugin 試：`claude --plugin-dir <這個 repo>`；headless 判一次：
  `claude -p --plugin-dir <repo> --permission-mode auto --allowedTools "mcp__plugin_terminu_harness__*" Agent Read Grep Glob "/terminu:judge <app> <repo>"`。

## 發布

還沒發過。GitHub repo、tag、release 的 binary 怎麼給使用者，發版時跟 user 定。
