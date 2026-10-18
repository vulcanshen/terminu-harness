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

**界線**：harness 只放流程（skill、agent、MCP server、script、範本），**不寫任何規定**。規定一律在 tdp；judge 判到規範沒寫、寫得模糊
的地方，回報成 tdp 的缺漏，不在 harness 裡自己補一條。

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
- new-app（問需求、照 project 範本開 repo、用 terminu-go 寫）、report-gap（開 issue 回 tdp 或 terminu-go）、align 之後再做。

## 還沒定（跟 user 一題一題定）

- 成績單的格式（檔名、欄位、證據怎麼放）與輸出目錄放哪裡。
- MCP 工具的確切介面（名稱、參數、回傳的格子長什麼樣）。
- 時間軸怎麼取樣（開關動畫、splash 揭露的 250ms、400ms）。
- 畫圖片用的字型（Nerd Font 要不要跟著帶）。
- 怎麼拿到某一版 tdp 的清單（本機 `git show <tag>:checklist/…`，或 GitHub 上的檔案）。
- 終端機模擬器從哪來：terminu-go 的 `vt.go` 沒有對外公開（型別 `vt` 是小寫），要用就照 GPL 複製過來、註明出處；或改用別的模擬器
  （例：charmbracelet 的 x/vt）。
- harness 不是家族 app，tdp 的 project（`family.mk` 等）不套；repo 自己的 Makefile、發布方式要不要照它，跟 user 定。

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
