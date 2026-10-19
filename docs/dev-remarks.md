# terminu-harness 開發者備忘

terminu 的 harness：一個 Claude Code plugin，照 terminu design（tdp）用 terminu-go 做 TUI app，做完照 tdp 的清單檢查。harness 放流程與判法（怎麼判、怎樣算通過），不寫 app 的規定；規定一律在 tdp。

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
  └─ judge_finish ─────────────► report.md（結果與 gene 都由程式算）
```

- **MCP server**（`cmd/terminu-harness`，子指令 `mcp`）：plugin.json 的 `mcpServers` 跑 `scripts/mcp.sh`（不放根目錄的 `.mcp.json`：在這個 repo 裡開 Claude Code 時它會被當成專案的 MCP 設定，`${CLAUDE_PLUGIN_ROOT}` 是空的就啟動失敗），server 名 `harness`，工具在
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

- **判斷線在 harness**（user 2026-10-19：tdp 只放規範，判斷是 harness 的核心）：有任何一條違反就「不通過」，不管是 app 自己要修、
  還是在等 tdp 或 terminu-go（user 選的：有違反就是有東西壞了）；沒有就「通過」。偏離（寫明理由、刻意不照做，tdp 允許的）不擋，
  列在成績單最前面。違反照 app 的 `docs/<app>-terminu-fix.md`（tdp E7）分開列：列在裡面的是在等 tdp 或 terminu-go，沒列的是 app
  自己要修（程式找編號，`K1-3` 不會配到 `K1-30`）。gene 數量留著看進度：它跟 app 的功能多寡綁在一起，不同 app 的數字比不了。
- **不適用要有根據**（user 2026-10-19 問「會有不適用的情況嗎」）：components 規範的是「用了這個元件時要這樣」，不是每個 app 都要
  有，所以不適用留著；但只有 app 真的沒有、tdp 也沒要求要有才算，理由要寫怎麼知道沒有；「都要」的條目、tdp 要求某種情況必須有的，
  缺了判違反。寫在 `agents/judge.md`。
- **只重判有改到的**（user 2026-10-19，粗分）：開 run 時記下 app 的 commit 與 terminu-go 版本；找同一個 app、同一版 tdp、判完的上一份成績單（同一版 terminu-go、有記 commit），看兩個 commit 之間（加上還沒 commit 的）改了哪些檔：改到程式（`.go`、`go.mod`、`cmd/`、`internal/`）就重判「讀檔案」以外的條目，改到其他檔就重判「讀檔案」的條目，其他沿用，理由前面標「沿用 <commit>」、證據連回上一份成績單。細分（只重判改動碰到的功能）不做：改到共用的程式時會漏判。
- **按程式碼分（`judge_scan`）**（user 2026-10-19）：有 repo、用 terminu-go 寫的 app，掃非測試的 `.go`，看用了哪些 terminu-go 元件（`terminu.DatePicker`、`OpenTerminal(`、`Details:`……），對到清單「適用」表的功能；「適用」寫的功能全部確定沒有的條目直接判不適用、理由寫出來。掃描看不出來的功能（多 panel、存設定……）與「都要」不判；app 有自己畫的 body 或 popup 時一律不判（它可能用自己的方式做了那個功能）。hostu 試過：1302 條裡 493 條。tdp v0.8.0 加的：模式分成 panel 上（`terminu.Mode`、`EnterMode(`，與 PTY）與 popup 裡（viewer 的選取、reorder 的 Drag、terminal 的歷史與選取，元件自己帶的）；有預覽的 finder、有檔案預覽看一般 finder 與 file-picker；發過版看 repo 有沒有 tag（沒有 tag 就一定沒發過）。terminu-go v0.2.0 加的：位置列看 `List.Location`；error popup 也看 `ErrorAbout(`、不限第一個參數是字串（`err.Error()` 不算）；confirm 也看 `AskFor:`；loading 也看 `LoadSelect(`、`LoadChecks(`、`LoadingFrame`；color picker 算有 select（清單說明頁）；自己畫的東西不限參數名，popup 的 body（`Draw(w, h int, layer Color)`）、`App.Popup(` 與 `App.Canvas` 都算。zoom、工作清單沒有對應的 terminu-go API，不判。不適用的理由照功能寫怎麼知道的（`scan.Why`）。
- **`terminu-harness report <run>`**：用目前這版 harness 的判法重寫一個 run 的成績單，不用重判。

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
- **harness 定的判法**寫在 `agents/judge.md`：時間差在兩格畫面內（HN2）、每個 app 都要有的東西整個沒有時細節也全部判違反
  （HN4，user 2026-10-19 選的）、列在 fix 文件的不算漏寫偏離（HN7）。不寫的話，十批 judge 各判各的（用測試程式跑完整流程時看到
  過）。原本還有 HN1、HN5、HN6、HN8 的暫時判法，tdp v0.8.0 自己定了，拿掉。
- **會寫系統檔的 app 用包裝指令判**：judge 會在 app 裡按會寫檔的鍵。判 hostu 時給 judge 一個包裝指令：每次啟動在暫存資料夾開新的
  `HOME` 與 hosts 檔（系統檔的複本），使用者的 `XDG_CONFIG_HOME` 指向真的家目錄時拿掉；judge skill 把使用者的這段補充原樣轉給
  每個 judge。new-app 也照這樣交代。

## 已知的牆與未做

- 只支援 macOS、Linux（`creack/pty`，沒有 ConPTY）。
- 不送滑鼠事件。
- Bubble Tea v1 一秒最多畫 60 格，tick 比一格短的動畫從畫面上看不出每個 tick（HN1）。
- 使用者要有 Go：第一次啟動時 `scripts/mcp.sh` 從原始碼建 binary，可能超過 Claude Code 等 MCP server 的時間。
- 等 issue 不是一直開著的 session：Claude Code 沒有可靠地等好幾天的方式（Monitor 每次最多 30 分鐘、`/loop` 要 session 開著、routine 沒有 issue 事件），所以靠使用者回來打 `/terminu` 選「繼續」。
- 判一次很貴：一千三百多條、11 個 agent。測試程式（幾乎沒有功能）約 22 美元；hostu 約 74 美元、30 分鐘。
- 還沒做：align（tdp 新版時重判、列出要修的）；judge 找到的 tdp 沒規範的東西，judge skill 會問要不要交給 report-gap，還沒有自動
  轉成草稿。

## 用 hostu 試整套流程（2026-10-18）

new-app 照 user 給的需求做了 hostu（只在本機，2026-10-19 移除；三次 judge 的成績單還在 `~/.terminu-harness/judge/hostu/`）。學到的：

- judge 抓得到 new-app 自己寫錯的地方：spec 的 D5 原本用 textarea 編 hosts 檔，judge 照 tdp textarea-1 判違反（設定檔要交給
  `$EDITOR`），改了。new-app 問需求時要先讀用到的元件的「用途」段，不只看「長相」。
- 開發中用 harness 的終端機開 app，抓到 L4 違反（Pager 照原樣畫 Tab），畫面測試量不出來（`terminu.Width` 把 `\t` 算 0 格）。
- 卡在 terminu-go 的條目多（第一次 judge 的 20 條違反裡 8 條），照流程標 TODO、列進 app 的 `docs/<app>-terminu-fix.md`、回報。
- 第一版直接編 hosts 檔，user 要的是 TUI 的做法：app 有自己的格式存設定（`profiles.yaml`，一個 domain 一筆、可以有多個位址），
  啟用時才匯出成系統檔，也能把現有的系統檔解析回來。new-app 決定設計時，管理某種檔案的 app 先想「自己的格式＋匯出」。
- 試流程的 app 不要一路跟著改：hostu 做了兩版、判了三次，中間 tdp 的判法與 harness 都改過，最後的樣子是一路修出來的，測不出
  「從需求做到通過」的整套流程（user：「現在 hostu 已經有形狀了，這種流程不太對了」）。規範或 harness 大改之後，從需求重做一次。

## 建置與開發

- `make build`：建到 `.local/bin/terminu-harness`，`scripts/mcp.sh` 優先用它。
- `make test`；設 `TERMINU_HARNESS__TDP_REPO=~/Documents/sideproj/terminu` 時也驗真的清單（v0.8.2 要是 1394 條）。
- server 的測試開 `internal/server/testdata/fixture`（測試時才建的小程式，不靠外部套件）代替真的 app：它照家族 app 的方式探測
  icon 寬度、顯示收到的鍵、按 `a` 播 5 格動畫、按 `w` 印超過邊界、按 `q` 結束。驗 MCP 工具用它，不拿家族成員來試。
- `make check`：gofmt、go vet、測試、`claude plugin validate`。
- 不裝 plugin 試：`claude --plugin-dir <這個 repo>`。這個 flag 只管那一次啟動，`--resume`、`-c` 接回舊的 session 時也要再帶，
  沒帶就沒有 plugin，也沒有 harness 的 MCP（2026-10-19 在 harness 的 session 踩到：`claude --resume -c` 開的 session 沒有 harness
  的工具）。`claude --plugin-dir <repo> mcp list` 看得到 `plugin:terminu:harness … Connected` 就是連得上。headless 判一次：
  `claude -p --plugin-dir <repo> --permission-mode auto --allowedTools "mcp__plugin_terminu_harness__*" Agent Read Grep Glob "/terminu:judge <app> <repo>"`。

## 發布

還沒發過。GitHub repo、tag、release 的 binary 怎麼給使用者，發版時跟 user 定。
