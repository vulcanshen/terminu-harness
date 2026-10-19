# terminu-harness

terminu 的 harness：一個 Claude Code plugin，放進 marketplace，讓人照 terminu design（tdp）做 TUI app、檢查 app 有多少 terminu 的 gene。

- repo：`terminu-harness`（本機 `~/Documents/sideproj/terminu-harness`，還沒有 GitHub repo）；plugin 名 `terminu`（skill 打
  `/terminu:judge` 這種）。
- 授權：GPL-3.0（跟家族其他 repo 一樣；可以用 terminu-go 的程式）。
- 規範：`~/Documents/sideproj/terminu`（tdp），照**最新 tag** 做（現在是 **tdp v0.7.4**）。judge 的依據是 tdp 的條目清單
  `checklist/`（`items-zh_TW.md`、`items.md`，1302 條）。

上面這些都已經跟 user 定案，不要重提。

## 為什麼有它

user 要把 terminu 做成 harness（2026-10-17～18 在 ideas 與 terminu 的 session 談的；摘要在
`~/Documents/sideproj/terminu/.local/harness/ideas-discussion-2026-10-18.md`）：一般人裝了 plugin，skill 問需求、照 tdp 用 terminu-go
做出家族 app；做不到的開 issue 回 tdp，user 審核後更新。harness 分四層：說明（tdp、skill）、限制（terminu-go）、回饋（issue → tdp
→ terminu-go）、**檢查** —— 前三層已經有，缺的是檢查。這個 repo 第一步就是做檢查：judge。

**界線**：tdp 只放規範（app 該長什麼樣、怎麼操作）；harness 放流程（skill、agent、MCP server、script、範本）**與判斷**（怎麼判、
怎樣算通過）。user（2026-10-19）：「tdp 應該就是單純規範，至於判斷應該在 harness 這邊，如果判斷都在 tdp 了，那乾脆整個 harness
都搬去那邊，因為這是 harness 的核心」。harness 不替 tdp 補 app 的規定：judge 判到規範沒寫、寫得模糊的地方，回報成 tdp 的缺漏。

## harness 的全貌（user 2026-10-18，在 harness 的 session）

user：「harness 不應該只是這樣而已」—— judge 只是最後一步。整套應該是：

1. 使用者裝好 plugin，打 **`/terminu`** 就開始跟使用者互動。
2. 要做新 app：**問需求**，最終確認後**開始開發**。
3. 開發中遇到問題**跟使用者反饋**，並**發 issue 到 terminu 的 GitHub**，等 tdp 與 tgo 回應或更新，再繼續開發，直到完成。
4. 完成後**自動 judge**，生成報告。

細節（入口怎麼叫、怎麼等 issue、誰寫程式、狀態放哪……）跟 user 一題一題定：

- **卡在 tdp 或 tgo、issue 發出去之後，其他部分照做**（user 選 A）：卡住的地方用暫時的做法，在程式標 `TODO(tdp #N)`
  （tgo 的標 `TODO(tgo #N)`），tdp 或 tgo 更新後回來改掉。理由：issue 可能等好幾天。
- 查過 Claude Code 文件（2026-10-18）：plugin skill 的 frontmatter 設 `name: terminu`，裸的 `/terminu` 也叫得到它（沒有別的
  指令佔這個名字時），要實測；一個 session 沒辦法可靠地等好幾天（Monitor 每次最多 30 分鐘、`/loop` 要 session 開著、routine
  是 research preview 且沒有 issue 事件）。
- 照建議定（user 不在，2026-10-19 檢查）：入口 skill `skills/terminu`（`name: terminu`）只問要做什麼，再叫 `terminu:new-app`
  或 `terminu:judge`；開發由主 session 寫，judge 照舊是獨立的 agent；進度放 app repo 的 `.local/terminu/progress.md`（tdp 的
  gitignore 範本不收 `.local/`），規格放 `docs/spec.md`；issue 一律先給使用者看草稿、同意才發（`terminu:report-gap`）；完成後
  judge，違反的條目問使用者要不要修。等 issue 的方式：使用者回來打 `/terminu` 選「繼續」，harness 查 issue 與新 tag。

## 測試：用 harness 做 hostu（user 2026-10-18）

- user 給的需求：hostu 管理本地 hosts 檔，master-detail（側邊欄是存著的設定清單，detail 是那個 hosts 檔的內容），一次只有一個
  生效、要有 icon 標示，可以編輯、新增、切換生效的那一個；參考舊的 CLI `~/Documents/sideproj/hostfile`；支援 Windows、macOS、
  Linux。
- user：「專案可以反覆刪除和創建，直到完成，過程中不會開啟新的 session，由你直接進行整個流程，直到開發完成」；之後設了目標
  「請依照指示執行直到此專案及測試專案開發完成」，明天（2026-10-19）來檢查。
- user 不在的期間：需求由我照推薦選項定、寫進 `docs/spec.md` 標成「待 user 確認」，照目標先開發；**對外的動作不做**：不建
  hostu 的 GitHub repo、不 push、不發版。規範盲點照舊記 HN。
- user（2026-10-18）：「issue 發給 terminu 的 session，直接本地端通知，還不走 github issue 那套」「hostu 先不 push 都留在本地」。
  草稿留在 hostu 的 `.local/terminu/issues/`，內容用 SendMessage 送給 terminu 的 session（2026-10-18 送了三張）。
- hostu 放在 `~/Documents/sideproj/hostu`。它是新做的 app，不是正在重構的成員；完成後照流程 judge。

## 已定（user 2026-10-18，在 terminu 的 session）

**分數是 gene 的數量**（tdp v0.7.3 起寫在 checklist 的說明頁）：

- 清單的每一條就是一個 gene；app 判成「符合」的一條，就有一個 gene。
- 只數有幾個、不算比例，寫成 `tdp-<版本> <數> gene`（例：`tdp-0.7.4 860 gene`），綁在判的那一版 tdp。
- 判定五種：符合、偏離、違反、不適用、判不了；只有符合算 gene。
- app 做了 tdp 沒規範的東西，不算也不扣，另外列出來回報成 tdp 的缺漏。

**judge 怎麼跑**：

- **能判任何 TUI app**：被判的不一定是家族成員、不一定用 Go、不一定拿得到程式碼。judge 從外面操作 app。
- **實際把 app 跑起來看**：照每一條的「怎麼驗」按鍵、看畫面；畫面看不到的（讀程式、讀檔案）拿得到 repo 才判，拿不到判「判不了」。
- **AI judge**：一個獨立的 agent，只讀、不寫程式，跟寫 app 的那一方分開。先從 `<app> --help`、`<app> version` 與 app 裡的說明理解
  app（tdp 的 app 就是 `Space` 與 `?`），列出 app 有哪些功能，決定哪些條目適用（checklist 說明頁的「適用」）。
- **終端機由 harness 自己帶**：一個 PTY 加終端機模擬器的小程式，做成 **MCP server**。工具大致是：開 app（大小、icon 與 powerline
  符號各佔幾格、環境變數）、送按鍵、看畫面的格子（每一格的字、前景色、背景色、粗體、游標）、看畫成的圖片、改大小、關掉；app 在
  server 裡一直開著。**不假設機器上有 tmux 或圖形介面**（user：「執行 judge 的機器不一定有 tmux」）。
- **固定步驟的條目用程式判**：MCP server 另外給幾個檢查工具（例：每一列剛好等於終端機寬度、某一格的色碼、範本 diff）。先做少數，
  判的時候發現常用的再加。
- **只手動判**（`/terminu:judge`），不用 hook 自動擋 PR。
- **每次全部判**：一份成績單是那一次對那一版 tdp 的完整結果。
- **成績單寫在 judge 自己的輸出目錄**，不動 app 的 repo：每一條的判定與證據（畫面的格子或圖片）、gene 數、tdp 沒規範的東西。

## 第一步：judge，拿 demou 試（user 定的）

- MCP 終端機、judge 的 skill 與 agent、成績單。
- 第一個判的 app 是 demou（`~/Documents/sideproj/demou`）：小、每次啟動都一樣、是 tdp project 的參考實作。判完得到第一個
  `tdp-0.7.4 N gene`，也順便試清單好不好用。
- **現在先不判任何成員**（user 2026-10-18）：所有成員正在大重構，判了也沒意義。專心把 harness 做完；要開真的 app 驗 harness
  的話，用 harness 自己的測試用 app，不碰成員。2026-10-18 對 demou 跑到一半的那次已經停掉、刪掉。
- new-app（問需求、照 project 範本開 repo、用 terminu-go 寫）、report-gap（開 issue 回 tdp 或 terminu-go）、align 之後再做。

## 已定（user 2026-10-18，在 harness 的 session）

- **終端機模擬器照 GPL 複製 terminu-go 的 `vt.go`**（連同它量寬度的 `width.go`），註明出處；不用 charmbracelet 的 x/vt。理由：
  icon 與 powerline 各佔幾格是 vt.go 本來就有的設定，量法跟家族 app 同一套，判 E5 才準；它會回 CPR、DA1、調色盤。x/vt 把 icon
  一律算一格，要 fork 才能判 E5。代價：非家族 app 用到 vt.go 沒做的控制序列時畫面可能不對，遇到再補；terminu-go 改 vt.go 時
  harness 自己跟。
- **看畫面預設只回文字**（一列一行，加游標位置）；要看顏色、粗體時另外指定一個範圍，只回那一塊的樣式段（從第幾欄、連續幾格、
  前景、背景、粗體……）。理由：大部分條目看文字就判得出來，顏色只看少數條目的一小塊。demou 一個 120×40 的畫面樣式有 269 段，
  每次全回要多約 5k token；每一格一筆（4800 筆、約 7 萬 token）不考慮。
- **時間軸只在畫面變了的時候記**：app 每次輸出、畫面真的變了，記一格並標上時間（ms）；送鍵時可以要求「接著錄 N ms」，回傳每一格
  的時間與跟上一格比較變了哪幾列。不用固定間隔拍快照（會漏格、會重複、資料量大）。精度受 app 自己的畫面更新率限制
  （HN1、HN2）。

下面這些 user 說「都依照你 recommended 的，一路做到結束，我們再來調整」（2026-10-18），照建議定，做完再一起調：

- **畫圖片的字型由 harness 自己帶**：Go Mono（`golang.org/x/image/font/gofont`）畫字，Symbols Nerd Font Mono（MIT，放在
  `internal/render/fonts/`）畫 icon 與 powerline 符號，框線與方塊字元自己畫；哪個字型都沒有的字畫成外框。理由：跑 judge 的
  機器不一定有字型，自己帶才每台畫得一樣。
- **清單從 GitHub 抓那一版 tag 的檔案**（`raw.githubusercontent.com`，不給版本就用最新的 tag），存在
  `~/.terminu-harness/cache/tdp/`；設 `TERMINU_HARNESS__TDP_REPO` 指到本機的 terminu clone 時改用 `git show`。理由：裝 plugin
  的人沒有 terminu 的 clone。
- **成績單**寫在 `~/.terminu-harness/judge/<app>/<時間>-tdp-<版本>/`（`TERMINU_HARNESS__OUT` 可以換）：`report.md`（給人看）、
  `verdicts.jsonl`、`gaps.jsonl`、`profile.md`、`evidence/`（畫面的文字、錄影、樣式、圖片，編號 `eNNNN`）、`batches/`、清單的
  副本。gene 由程式數。理由：user 多半在 app 的 repo 裡打 `/terminu:judge`，寫在工作目錄就動到 app 的 repo。
- **分工**：一個 judge 先認識 app、存 profile；清單照章切成 10 批（不拆章），每批一個 judge 同時判；`judge_status` 補沒判到的
  一次；`judge_finish` 寫成績單，還是沒判到的算判不了。理由：一個 agent 判不完 1302 條。
- **judge agent 不給 Write、Edit**：判定、profile、tdp 沒規範的東西都用 MCP 工具記；Bash 只准看東西的指令（寫在 agent 裡）。
- **MCP 工具**：`open`、`keys`、`screen`、`styles`、`image`、`resize`、`overflow`、`close`、`tdp_file`、`judge_start`、
  `judge_profile`、`judge_verdicts`、`judge_gap`、`judge_status`、`judge_finish`（server 名 `harness`）。Go 寫，用
  `mark3labs/mcp-go`（家族的 clerk 也用）、`creack/pty`；只支援 macOS、Linux。
- **repo 自己的 Makefile**（build、test、check、clean），不套 `family.mk`。plugin.json 的 `mcpServers` 跑 `scripts/mcp.sh`（2026-10-19 從根目錄的 `.mcp.json` 搬過去）：有
  `.local/bin/terminu-harness` 就用，沒有就用 Go 從原始碼建到 `~/.terminu-harness/bin/`，所以現在使用者要有 Go。發版要不要改成
  下載 release 的 binary，發版時跟 user 定。
- **時間的暫時判法**照 HN1、HN2 的「暫時」寫在 judge agent 裡，tdp 定案後改。

## 已定（user 2026-10-19，在 harness 的 session）

- **判斷線在 harness**（見上面的「界線」）：通過＝沒有違反；待修＝違反都列在 app 的 `docs/<app>-terminu-fix.md`；不通過＝有違反沒
  列進去。成績單第一行寫結果，gene 數量看進度。`terminu-harness report <run>` 照目前的判法重寫成績單。
- **只重判有改到的，粗分**（user 選 A）：記 commit 與 terminu-go 版本；改到程式重判「讀檔案」以外的條目，改到其他檔重判「讀檔案」
  的條目，其他沿用上一份成績單；terminu-go 換版、沒有上一份就全部判。
- **判一次太久，按程式碼分**（user 的想法，還沒做）：有原始碼時掃程式，知道哪些功能用 terminu-go 做、哪些是 app 自己寫的；沒有的
  功能直接判不適用，tgo 做的直接判，app 自己的先看適不適用。
- **new-app 只問需求**：設計由 harness 決定（user：「這個應該要 harness 自己決定，使用者只是提出需求」）。

## 還沒定

- 發布：GitHub repo、tag、release 的 binary 怎麼給使用者（見上一節最後）。

## 規範與清單是唯一的依據

- **讀什麼**：terminu 最新 tag 的 `checklist/` 全部（說明頁與條目）、`principle/`、`components/`、`project/`。繁中是原文。
- **參考**：terminu-go（`~/Documents/sideproj/terminu-go`）的 `vt.go`、`scripts/tmux-shot.sh`（tmux 版的「開 app、送鍵、截畫面」，
  harness 不用 tmux，但做法可以參考）；ccdrawu（`~/Documents/sideproj/ccdrawu`）是家族裡已經做成 Claude Code plugin 的 app，
  `.claude-plugin/plugin.json`、`marketplace.json`、`claude plugin validate`／`claude plugin test` 可以照它。
- **不改**：terminu 的規範與清單、terminu-go、其他家族 app。

**盲點**：清單或規範寫不清楚、judge 判不下去、照字面會判錯的地方，不在 harness 裡自己決定，記下來：

記在 `~/Documents/sideproj/terminu/.local/family-fix/harness/README.md` 的「待定」，編號用 **HN**：

```
### HN<n> <一句話標題>
- 狀態：待定
- 種類：清單（條目寫錯、漏拆、適用不對）／規範（盲點）
- 出處：<清單編號> 或 <檔案> 的 <段落>（tdp v0.7.4）
- 卡在：judge 需要知道什麼才能判
- 可能：A …；B …（各附一句後果）
- 暫時：判「判不了」，或暫時照 A 判
```

寫紀錄前先重讀檔案（terminu 的 session 也會寫它）。

## repo 一開始要做的

1. `git init`、`LICENSE`（GPL-3.0 全文，可以從 terminu-go 複製）。
2. `README.md` 與 `README-zh_TW.md`：README 是工具介紹，開發者的東西放 `docs/dev-remarks.md`。
3. `CHANGELOG.md`。
4. `.claude-plugin/plugin.json` 與 `marketplace.json`（照 ccdrawu）。
5. 建 GitHub repo、push、打 tag、發布都先問 user。

## 跟 user 合作的方式

- 一次問一題，用問句寫明要決定什麼，列 A／B 選項附上建議與理由。不要寫一串直述句再問「可以嗎」。
- 先講流程與操作（judge 會做什麼、user 看到什麼），再寫程式。
- 提案前先查規範與清單怎麼說。
- user 看不懂時，用一個具體的例子重講。
- 只改這個 repo 與自己的紀錄。
- 一律用繁體中文對話；commit 不加任何 Claude 署名。
