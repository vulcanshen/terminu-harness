# terminu-harness

terminu 的 harness：一個 Claude Code plugin，放進 marketplace，讓人照 terminu design（tdp）用 terminu-go 做 TUI app，寫程式時就檢查有沒有照
tdp 做。

- repo：`terminu-harness`（本機 `~/Documents/sideproj/terminu-harness`，還沒有 GitHub repo）；plugin 名 `terminu`（skill 打
  `/terminu:newapp` 這種）；skill 名不帶連字號（user 2026-10-20：`new-app` → `newapp`、`report-gap` → `reportgap`）。
- 授權：GPL-3.0（跟家族其他 repo 一樣；可以用 terminu-go 的程式）。
- 規範：`~/Documents/sideproj/terminu`（tdp），照**最新 tag** 做（現在是 **tdp v0.10.2**，2026-10-26 發的；terminu-go 最新是 v0.5.2，對應 tdp v0.10.2，tdptest 從 v0.3.0 起）。tdptest 與 reviewer 的依據是 tdp 的條目清單 `checklist/`（`items-zh_TW.md`、`items.md`，1453 條）。

上面這些都已經跟 user 定案，不要重提。

## 為什麼有它

user 要把 terminu 做成 harness（2026-10-17～18 在 ideas 與 terminu 的 session 談的；摘要在
`~/Documents/sideproj/terminu/.local/harness/ideas-discussion-2026-10-18.md`）：一般人裝了 plugin，skill 問需求、照 tdp 用 terminu-go
做出家族 app；做不到的開 issue 回 tdp，user 審核後更新。harness 分四層：說明（tdp、skill）、限制（terminu-go）、回饋（issue → tdp
→ terminu-go）、**檢查** —— 前三層已經有，缺的是檢查。這個 repo 第一步就是做檢查：judge。2026-10-20 改了方向、2026-10-25 改了
架構（見「方向」）：不再做完才 judge；harness 替 app 記紀錄（`harness-assets/`），review 拿 tdp 對紀錄、不讀程式碼，程式對不對
交給測試。

**界線**：tdp 只放規範（app 該長什麼樣、怎麼操作）；harness 放流程（skill、agent、MCP server、script、範本）**與判斷**（怎麼判、
怎樣算通過）。user（2026-10-19）：「tdp 應該就是單純規範，至於判斷應該在 harness 這邊，如果判斷都在 tdp 了，那乾脆整個 harness
都搬去那邊，因為這是 harness 的核心」。harness 不替 tdp 補 app 的規定：開發時或 review 遇到規範沒寫、寫得模糊的地方，回報成 tdp
的缺漏。2026-10-20 起，程式判得了的檢查放在 terminu-go（見「方向」）；harness 留下流程、專案骨架、紀錄與它的檢查（`check`、commit
hook）、reviewer、reportgap、MCP 終端機，以及決定要 tgo 做哪些檢查、提哪些刁鑽的值。

## 方向（user 2026-10-20 定方向、2026-10-25 定架構，在 harness 的 session）

2026-10-20，user：「我們應該是想辦法把東西作對，而不是做完以後花大把力氣 judge」。對照 harness engineering（模型以外的一切：說明、
限制、回饋；讓 agent 一開始就做對、做錯當場知道），驗證該是寫程式時一直跑、便宜、結果固定的檢查，不是做完之後的大考。

2026-10-25，user 從頭重想了架構（「我想要重新思考我們這個 harness 的架構」），流程與檔案結構是 user 自己寫定的：**用這個 plugin
很快做出符合 tdp 的 terminu app，使用者只管功能**，UI/UX 交給 terminu。harness 只做兩件事：「創建 terminu 成員」與「檢查 app 的
機制有沒有違反 tdp」；程式對不對是測試的事。

**這一節取代下面各節跟 judge 有關的決定。** 2026-10-20 版的細節（reviewer 看 diff、progress.md、「偏離 tdp」、tdptest 怎麼提給
tgo）在 `e36d757` 的 CLAUDE.md。

- **只管用 terminu-go 做的 app**（user 選 A：「terminu harness 當然只看 tdp 家族就好, 就是用 tgo 做的東西」）。取代 2026-10-18 的
  「能判任何 TUI app」。
- **harness 是封閉的，tdp 是宇集**（user 2026-10-25）：app 的每個決定都要在 tdp 找得到根據，或在 tdp 明文寫「app 可以自行決定」的
  範圍裡（範圍內 app 自己處理 UI/UX，甚至可以不用 tgo）。找不到、不清楚、矛盾只能回報 tdp；比較差但合 tdp 的暫時做法照樣回報，
  連暫時做法都沒有，就只有那個功能等。**不再有「偏離」**（user：「連放行的地方都要是明文放行，就是明確寫『app 可以自行決定』」）。
  harness 照這個回報了 tdp：tdp v0.10.0 拿掉了「偏離」與開發者文件的規定，P5 改成明文放行，每一節都給了編號（引用寫 `tdp K4`、
  `tdp panel.5`、整個檔 `tdp confirm`；v0.10.2 加了範本路徑 `tdp project/install.sh` 與清單條目 `tdp K9-6`）。
- **紀錄放 app 的 `harness-assets/`**（user 2026-10-25 定的結構；格式的唯一一份在 `skills/newapp/harness-assets.md`）：
  - `idea.md`：使用者的需求，跟著互動一直更新。
  - `spec.yaml`：需求與 tdp 的對照。`id` 在 idea.md 或 tdp 任一個改了時加一（同時改只加一次；user：兩者是「或」）；`from`、`to`
    是 tdp 版本（第一份 `from: 0`）；`tasks` 每筆有 `id`（整個 app 不重複使用）、`spec`（最後新增或修改它的那份 spec）、`request`、
    `tdp`、`canceled`（原因）、`remark`。只新增、更新，反映現況。開發中的決定（設計、remark、做不到而取消）不加 id。
  - `work.yaml`：spec 的每個 task 一筆，`id`、`spec`（照第幾份 spec 做的）、`implementation`、`test`、`canceled`。只新增、更新。
  - `review.yaml`：最近一次 review，`spec` 與每個 task 的 `accepted`、`reason`，每次覆蓋。
  - 做到哪：spec 有、work 沒有 → 要做；task 的 `spec` 比 work 的大 → 需求改過，要重做。不要 progress.md。
  - **app 照的版本不另外存**：tgo 看 go.mod，tdp 從那一版 tgo 的 CHANGELOG 推（`## v0.5.2 — tdp v0.10.2`）。每一版 tdp 都有一版
    tgo 對到它（user 選 A，2026-10-26），不加「只多放行的小版本當成前一版」的後備規則：放行本身就是 review 的根據。
- **我提的預設**（user 沒有異議）：implementation 寫觸發、經過、讀寫、用到的 tgo 元件與畫面上的字，一個入口一行；test 寫測試名字
  （`check` 拿 `go test -list` 對），寫不成測試的記 `user: {date, how}`、每次 review 請使用者重新確認，再不行才指出 `commit:`；
  review 每次全部重判；reason 開頭標「違反」「沒有測試」「tdp 不清楚」。
- **流程**（user 2026-10-25）：`/terminu` 描述需求 → 記進 idea.md → 依 idea.md 與 tdp 產生、更新 spec.yaml → 照 spec 建置，互動中的
  決定同時更新 spec.yaml 與 work.yaml → review 覆蓋 review.yaml → 發版前跑 review，還有 accepted false 就問使用者「照樣發版還是
  繼續打磨」，不擋。在這幾步之間循環。打 `/terminu` 時先查 tgo 有沒有新版，使用者同意就升。
- **升版**（user）：把 spec 的 `to` 到新版之間 tdp 各版的 CHANGELOG 先統整（同一節以最新的為主、只改 tgo 畫法的拿掉），再跟紀錄比對，
  開新一份 spec。
- **review 不讀程式碼**（user 2026-10-25）：`terminu:reviewer` agent 從乾淨的 context 開始，只讀 harness-assets 與 tdp（`tdp_section`
  照編號讀），判紀錄寫的做法合不合 tdp。用 tgo 做的 app 元件本身一定照 tdp，讀程式碼找問題是白工。**只判紀錄寫了的**：程式做了、
  work.yaml 沒寫的，harness 不保證（user 選 A，寫進 README 與 app 的規則檔）。review 不拿 spec 判需求有沒有做到；發版時列出還沒做的
  需求，不擋。
- **reviewer 什麼時候跑**（user 2026-10-20）：「應該不管是不是 newapp, 有兩種方式可以叫 reviewer, 1 個是 app 發版前， 1 個是
  手動呼叫」：`/terminu:release` 與 `/terminu:review`，都在 app 自己的 session 裡打。
- **`/terminu:release` 包整個發版**（user 選 A）：照 tdp project.6、project.7：工作目錄乾淨、`make check`、review、CHANGELOG 的
  `## Unreleased` 改成 `## vX.Y.Z — tdp <版本>`，問使用者版號與要不要打 tag，同意才 tag、push。
- **commit hook**（user 同意）：harness 管的 app 裡，commit 改了程式、`harness-assets/work.yaml` 沒改就擋；行為沒變的在 commit 訊息
  最後加一行 `Terminu-Work: unchanged`。
- **接手現成的 app：`/terminu:analysis`**（user 2026-10-26：「依照專案的 readme.md 生成 spec.yaml，然後依照 git commit history，
  生成 work.yaml，這樣 review 的基礎就有了，接下來就可以回到 harness 的流程中」）：idea.md 也從 README 寫；測試名字從
  `go test -list` 拿，用 git history 找加進它的 commit 對回 task；不讀程式碼，也不讀舊的開發者文件。原本叫 genwork，也讀 CHANGELOG
  與開發者文件。接手後刪掉舊的開發者文件（tdp 不再管開發者文件，成員的開發者文件在接手時刪掉，user 在 terminu 定的；給使用者看的、
  程式會用到的、拿不準的留著）。
- **app 的規則檔放 app 自己的 `.claude/rules/terminu.md`**（user 2026-10-26，在 terminu 定的），不由 plugin 在 SessionStart 帶進去：
  newapp 開 repo 時放、analysis 接手時放。家族成員的接手在成員自己的 session 裡跑，算成員自己改 `.claude/`。
- **tdptest 全部放在 terminu-go**（user 選 A）。user：「既然 lib 是 tgo, 那 tgo 自己作一個檢測 api, 應該更簡單? harness 直接使用就好，
  這樣 harness 也不用追版本，反正 app 用哪一版的 tgo, 就會跟著哪一版的 tgo test」。terminu-go v0.3.0 起有 `tdptest`（driver、
  `EachIconWidth`、`Check`、`Tricky`／`Sizes`／`Mashes`），app 的測試一律用它；harness 不改 tgo，要加檢查就寫成回報、等 tgo 發版。
  - **harness 提刁鑽的值**：開發 app 時踩到的值（控制字元、ESC 序列、寬字元、超長字串、空值、極小的終端機、一次一串按鍵）經
    reportgap 交給 tgo 收進自己的測試資料；「該長什麼樣」由 tgo 照 tdp 判。user 問過由 harness 連期望結果一起給：那等於 harness
    又在寫檢查，不做。
  - **涵蓋清單**：tgo 的測試標清單編號，產生「哪些條目有測試」的清單；harness 不另外驗 tgo（驗 lib 不是 harness 的事，就像做網站
    的 harness 不驗 React）。
  - 不在測試裡重算版面（那只驗得出 tgo 跟自己一致），只量畫出來的結果。
- **三個環**：內環（寫一個 task → `make check`、看畫面 → 更新 work.yaml、跑 `check` → commit；review 在發版前與手動叫時跑）；外環
  （回報 → tdp、tgo 發版 → harness 跟上 → app 升版、拿掉 TODO）；學習環（同一種錯出現第二次 → 請 tgo 加進 tdptest，或改 API 讓錯的
  寫法寫不出來）。
- **試整套流程時在 app 自己的 session**（user 選 A）：在 app 的資料夾開新的 session、只帶 plugin（例
  `claude --plugin-dir ~/Documents/sideproj/terminu-harness`），從 `/terminu` 開始走；harness 的 session 只改 harness，看那邊卡在
  哪裡再回來修 skill。理由：harness 的 session 知道太多（CLAUDE.md、tdp 的細節、user 講過的話），會不自覺替 skill 寫不清楚的地方
  補上，測不出來。取代 2026-10-18 的「過程中不會開啟新的 session，由你直接進行整個流程」。
- **跟 terminu-go 直接講**（user 2026-10-20：「tgo 那邊你可以直接跟他 session 溝通」）：給 terminu-go 的需求與回報，在本機直接
  用 SendMessage 跟 terminu-go 的 session 談，不用先給 user 看草稿。tdp（terminu 的 session）照舊先給 user 看。
- **terminu 的請求與通知直接照做**（user 2026-10-21：「可以接受 terminu 那邊的操作請求和通知」）：terminu 的 session 請 harness
  做的事（跟上新版、改 harness 的檔、commit）不用再問 user，做完回報 terminu。user 2026-10-26 確認範圍包括改 harness 自己的
  CLAUDE.md。建 GitHub repo、push、打 tag、發布照舊先問 user。
- **拿掉 judge**（user 選 A；2026-10-20：「judge 現在就可以拆了」）：成績單、十批 judge agent、`judge_*` 工具都不要了；MCP 終端機與
  `tdp_file` 留著。程式在 `fb89ad1` 還找得到；舊的成績單留在 `~/.terminu-harness/judge/`。
- **不留 gene**（user 選 A）：app 不數分數，看 review.yaml 裡每個 task 合不合 tdp。

## harness 的全貌（user 2026-10-18，在 harness 的 session）

user：「harness 不應該只是這樣而已」—— judge 只是最後一步。整套應該是：

1. 使用者裝好 plugin，打 **`/terminu`** 就開始跟使用者互動。
2. 要做新 app：**問需求**，最終確認後**開始開發**。
3. 開發中遇到問題**跟使用者反饋**，並**發 issue 到 terminu 的 GitHub**，等 tdp 與 tgo 回應或更新，再繼續開發，直到完成。
4. 完成後**自動 judge**，生成報告。（2026-10-20 改：不再做完才 judge，寫程式時就跑 tdptest 與 reviewer，見「方向」。）

細節（入口怎麼叫、怎麼等 issue、誰寫程式、狀態放哪……）跟 user 一題一題定：

- **卡在 tdp 或 tgo、issue 發出去之後，其他部分照做**（user 選 A）：卡住的地方用暫時的做法，在程式標 `TODO(tdp #N)`
  （tgo 的標 `TODO(tgo #N)`），tdp 或 tgo 更新後回來改掉。理由：issue 可能等好幾天。
- 查過 Claude Code 文件（2026-10-18）：plugin skill 的 frontmatter 設 `name: terminu`，裸的 `/terminu` 也叫得到它（沒有別的
  指令佔這個名字時），要實測；一個 session 沒辦法可靠地等好幾天（Monitor 每次最多 30 分鐘、`/loop` 要 session 開著、routine
  是 research preview 且沒有 issue 事件）。
- 照建議定（user 不在，2026-10-19 檢查）：入口 skill `skills/terminu`（`name: terminu`）只問要做什麼，再叫 `terminu:newapp`
  或 `terminu:judge`；開發由主 session 寫，judge 照舊是獨立的 agent；進度放 app repo 的 `.local/terminu/progress.md`（tdp 的
  gitignore 範本不收 `.local/`），規格放 `docs/spec.md`；issue 一律先給使用者看草稿、同意才發（`terminu:reportgap`）；完成後
  judge，違反的條目問使用者要不要修。等 issue 的方式：使用者回來打 `/terminu` 選「繼續」，harness 查 issue 與新 tag。
  （2026-10-20 改：judge 拿掉；不用選「繼續」，打 `/terminu` 就先查，見「方向」。2026-10-25 起沒有 progress.md 與
  `docs/spec.md`，紀錄在 `harness-assets/`。）

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
- user 看過第一版 spec 後加的需求（2026-10-19）：「對於某一個 profile 異動，或者新增 profile，應該要有我們 hostu 的 format，
  來儲存，類似 kamonu 的作法，設定是一回事，匯出是另一回事，所以 profiles.yaml 應該儲存這些設定，當啟用某一個 profile 的時候，
  才把該 profile 匯出覆蓋到 hosts 檔案，相反的，也應該要能夠把當前 hosts 檔案 parse 成 profile 格式的功能，這樣才能建立「預設」
  profile……在 TUI 中我們可以提供使用者一個 domain 一個 domain 的添加，或者一個 domain 多個 ip 的添加，不是讓使用者直接去編輯
  hosts 檔本身」。
- **2026-10-19 移除了 hostu**（user：「我們的 hostu 應該要移除，等這次 terminu 更新後，重新用更新後的 harness 架構重建，現在
  hostu 已經有形狀了，這種流程不太對了」）：一路跟著規範與判法改出來的 app 測不出整套流程，等 tdp 與 terminu-go 這次更新、
  harness 跟上之後，從上面兩段需求重新做一次。刪了 repo 與開發時寫到真的家目錄的 `~/.config/hostu`（真的 `/etc/hosts` 沒被
  寫過）；留著三次 judge 的成績單 `~/.terminu-harness/judge/hostu/`（terminu 的紀錄有引用），回報的草稿搬到 harness 的
  `.local/hostu-v2/issues/`。做過一次學到的寫在 `docs/dev-remarks.md`。
- **2026-10-19 照更新後的 harness 重做了 hostu**（user：「不應該問我，你應該先做出來，然後我們在來看結果」—— 不等 spec 確認）：
  HEAD 9080664，只在本機；第一輪 judge（tdp v0.8.2）不通過，12 條違反修了，沒有重判（之後 judge 拿掉了）。4 張回報送了 terminu，
  記成 HU11～HU14。進度在 hostu 的 `.local/terminu/progress.md`。

## 已定（user 2026-10-18，在 terminu 的 session）

> 2026-10-20 起被「方向」取代：judge 拿掉、不留 gene、只管用 terminu-go 做的 app。留著當紀錄。

**分數是 gene 的數量**（原本寫在 checklist 的說明頁；tdp v0.8.0 起判法都歸 harness，見「已定（user 2026-10-19）」）：

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

> 2026-10-20 起被「方向」取代。留著當紀錄。

- MCP 終端機、judge 的 skill 與 agent、成績單。
- 第一個判的 app 是 demou（`~/Documents/sideproj/demou`）：小、每次啟動都一樣、是 tdp project 的參考實作。判完得到第一個
  `tdp-0.7.4 N gene`，也順便試清單好不好用。
- **現在先不判任何成員**（user 2026-10-18）：所有成員正在大重構，判了也沒意義。專心把 harness 做完；要開真的 app 驗 harness
  的話，用 harness 自己的測試用 app，不碰成員。2026-10-18 對 demou 跑到一半的那次已經停掉、刪掉。
- newapp（問需求、照 project 範本開 repo、用 terminu-go 寫）、reportgap（開 issue 回 tdp 或 terminu-go）、align 之後再做。

## 已定（user 2026-10-18，在 harness 的 session）

> 2026-10-20：終端機模擬器（`vt.go`）、看畫面只回文字、時間軸、字型、清單從哪裡抓、Makefile、`scripts/mcp.sh`，以及 MCP 終端機的
> 工具（`open`、`keys`、`screen`、`styles`、`image`、`resize`、`overflow`、`close`、`tdp_file`）留著，給開發時看畫面；成績單、
> 分工、judge agent、`judge_*` 工具被「方向」取代。

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
  一次；`judge_finish` 寫成績單，還是沒判到的算判不了。理由：一個 agent 判不完一千多條。
- **judge agent 不給 Write、Edit**：判定、profile、tdp 沒規範的東西都用 MCP 工具記；Bash 只准看東西的指令（寫在 agent 裡）。
- **MCP 工具**：`open`、`keys`、`screen`、`styles`、`image`、`resize`、`overflow`、`close`、`tdp_file`、`judge_start`、
  `judge_profile`、`judge_verdicts`、`judge_gap`、`judge_status`、`judge_finish`（server 名 `harness`）。Go 寫，用
  `mark3labs/mcp-go`（家族的 clerk 也用）、`creack/pty`；只支援 macOS、Linux。
- **repo 自己的 Makefile**（build、test、check、clean），不套 `family.mk`。plugin.json 的 `mcpServers` 跑 `scripts/mcp.sh`（2026-10-19 從根目錄的 `.mcp.json` 搬過去）：有
  `.local/bin/terminu-harness`、而且沒有原始碼比它新就用（2026-10-26 hostu 試用時用到舊檔，加了比新舊），不然就用 Go 從原始碼建到 `~/.terminu-harness/bin/`，所以現在使用者要有 Go。發版要不要改成
  下載 release 的 binary，發版時跟 user 定。
- **時間的暫時判法**照 HN1、HN2 的「暫時」寫在 judge agent 裡，tdp 定案後改（tdp v0.8.0 定了，見下一節）。

## 已定（user 2026-10-19，在 harness 的 session）

> 2026-10-20：「判斷在 harness」的界線與「newapp 只問需求」照舊；其他跟 judge 有關的（判斷線、不適用要有根據、只重判有改到的、
> 按程式碼分、判法與 HN2、HN4、HN7）被「方向」取代：通過改成 `make check` 全過、reviewer 沒有附出處的問題。

- **判斷線在 harness**（見上面的「界線」）：有任何違反就**不通過**（user 選 A，包括在等 terminu-go 的），沒有就通過；偏離不擋、列在
  最前面；違反分開列 app 自己要修的與在等 tdp／terminu-go 的。成績單第一行寫結果，gene 數量看進度。`terminu-harness report <run>`
  照目前的判法重寫成績單。
- **不適用要有根據**（user 問「會有不適用的情況嗎」）：components 規範「用了才要這樣」，不適用留著；但只有 app 真的沒有、tdp 也沒
  要求要有才算，要寫怎麼知道；「都要」的與 tdp 要求必須有的，缺了判違反。
- **只重判有改到的，粗分**（user 選 A）：記 commit 與 terminu-go 版本；改到程式重判「讀檔案」以外的條目，改到其他檔重判「讀檔案」
  的條目，其他沿用上一份成績單；terminu-go 換版、沒有上一份就全部判。
- **判一次太久，按程式碼分**（user 的想法，2026-10-19 做了 `judge_scan`）：有原始碼時掃程式，知道哪些功能用 terminu-go 做、哪些是 app 自己寫的；沒有的
  功能直接判不適用，tgo 做的直接判，app 自己的先看適不適用。
- **newapp 只問需求**：設計由 harness 決定（user：「這個應該要 harness 自己決定，使用者只是提出需求」）。
- **tdp v0.8.0 把判法交給 harness**（user 在 terminu 的 session 選 A，S60）：tdp 只留條目與說明頁的「適用」表；五種判定、gene 怎麼數、
  判斷線、不適用要有根據、判不了要寫原因，以及 HN2（量時間的誤差）、HN4（東西整個沒有時細節怎麼判）、HN7（departure-1 與 fix
  文件）都是 harness 的，寫在 `agents/judge.md` 與 README 的「成績單怎麼看」。HN2、HN7 照原本的暫時判法定下來；**HN4 user 選細節
  全部判違反**（跟「都要的缺了判違反」一致）。HN1、HN5、HN6、HN8 tdp 自己定了，judge 的暫時判法拿掉。

## 還沒定

- 發布：GitHub repo、tag、release 的 binary 怎麼給使用者（見上一節最後）。

## 規範與清單是唯一的依據

- **讀什麼**：terminu 最新 tag 的 `checklist/` 全部（說明頁與條目）、`principle/`、`components/`、`project/`。繁中是原文。
- **參考**：terminu-go（`~/Documents/sideproj/terminu-go`）的 `vt.go`、`scripts/tmux-shot.sh`（tmux 版的「開 app、送鍵、截畫面」，
  harness 不用 tmux，但做法可以參考）；ccdrawu（`~/Documents/sideproj/ccdrawu`）是家族裡已經做成 Claude Code plugin 的 app，
  `.claude-plugin/plugin.json`、`marketplace.json`、`claude plugin validate`／`claude plugin test` 可以照它。
- **不改**：terminu 的規範與清單、terminu-go、其他家族 app。

**盲點**：清單或規範寫不清楚、開發時或 tdptest、reviewer 判不下去、照字面會做錯的地方，不在 harness 裡自己決定，記下來：

記在 `~/Documents/sideproj/terminu/.local/family-fix/harness/README.md` 的「待定」，編號用 **HN**：

```
### HN<n> <一句話標題>
- 狀態：待定
- 種類：清單（條目寫錯、漏拆、適用不對）／規範（盲點）
- 出處：<清單編號> 或 <檔案> 的 <段落>（tdp v<版本>）
- 卡在：要知道什麼才能做或判
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
