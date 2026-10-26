# terminu-harness

terminu 的 harness：一個 Claude Code plugin，放進 marketplace，讓人照 terminu design（tdp）用 terminu-go 做 TUI app，寫程式時就檢查有沒有照
tdp 做。

- repo：`terminu-harness`（本機 `~/Documents/sideproj/terminu-harness`，還沒有 GitHub repo）；plugin 名 `terminu`（skill 打
  `/terminu:newapp` 這種）；skill 名不帶連字號（user 2026-10-20：`new-app` → `newapp`、`report-gap` → `reportgap`）。
- 授權：GPL-3.0（跟家族其他 repo 一樣；可以用 terminu-go 的程式）。
- 規範：`~/Documents/sideproj/terminu`（tdp），照**最新 tag** 做（現在是 **tdp v0.10.1**，2026-10-26 發的；terminu-go 最新是 v0.5.1，對應 tdp v0.10.1，tdptest 從 v0.3.0 起）。tdptest 與 reviewer 的依據是 tdp 的條目清單 `checklist/`（`items-zh_TW.md`、`items.md`，1452 條）。

上面這些都已經跟 user 定案，不要重提。

## 為什麼有它

user 要把 terminu 做成 harness（2026-10-17～18 在 ideas 與 terminu 的 session 談的；摘要在
`~/Documents/sideproj/terminu/.local/harness/ideas-discussion-2026-10-18.md`）：一般人裝了 plugin，skill 問需求、照 tdp 用 terminu-go
做出家族 app；做不到的開 issue 回 tdp，user 審核後更新。harness 分四層：說明（tdp、skill）、限制（terminu-go）、回饋（issue → tdp
→ terminu-go）、**檢查** —— 前三層已經有，缺的是檢查。這個 repo 第一步就是做檢查：judge。2026-10-20 改了方向（見「方向」）：
檢查不再是做完才跑的 judge，改成寫程式時就跑的 tdptest 與只看 diff 的 reviewer。

**界線**：tdp 只放規範（app 該長什麼樣、怎麼操作）；harness 放流程（skill、agent、MCP server、script、範本）**與判斷**（怎麼判、
怎樣算通過）。user（2026-10-19）：「tdp 應該就是單純規範，至於判斷應該在 harness 這邊，如果判斷都在 tdp 了，那乾脆整個 harness
都搬去那邊，因為這是 harness 的核心」。harness 不替 tdp 補 app 的規定：開發時、tdptest 或 reviewer 遇到規範沒寫、寫得模糊的地方，回報成 tdp 的缺漏。2026-10-20 起，
程式判得了的檢查放在 terminu-go（見「方向」）；harness 留下流程、專案骨架、reviewer、reportgap、MCP 終端機，以及決定要 tgo 做哪些
檢查、提哪些刁鑽的值。

## 方向（user 2026-10-20，在 harness 的 session）

user：「我們應該是想辦法把東西作對，而不是做完以後花大把力氣 judge」。之前主力都放在 judge（做完才判、一千多條、十批 agent 看
畫面），偏離了原本的目的：**用這個 plugin，簡單做出符合 tdp 的 TUI 工具**。對照 harness engineering（模型以外的一切：說明、限制、
回饋；讓 agent 一開始就做對、做錯當場知道），驗證該是寫程式時一直跑、便宜、結果固定的檢查，不是做完之後的大考。user 對 judge：
「如果要評價 tdp 適用性，那不如多寫一個 tdp 測試工具，直接對 app 跑那個工具就好，不需要用 ai 去圖像檢查」。**這一節取代下面各節
跟 judge 有關的決定。**

- **只管用 terminu-go 做的 app**（user 選 A：「terminu harness 當然只看 tdp 家族就好, 就是用 tgo 做的東西」）。取代 2026-10-18 的
  「能判任何 TUI app」。
- **條目照誰保證分三層**（我提的，user 沒有異議）：terminu-go 直接保證的（用了元件就對，不在每個 app 重測）；程式檢查得了的
  （tdptest）；要判斷的（元件選得對不對、哪些動作要 confirm、措辭、UX 合不合需求，交給 reviewer）。tgo 的 bug 回報 terminu-go，
  app 先用暫時的做法標 `TODO(tgo #N)`，harness 不處理。user 問「如果 app 有問題，那也是 tgo 的問題」：tgo 只保證元件怎麼動，
  用哪個、什麼時候用、app 給的字、專案檔是 app 的事（hostu 第一次 judge 的 12 條違反，tgo 的只有 2 條）。
- **plugin 的產出**：
  - `/terminu` 與 newapp：從需求到完成的流程。
  - 專案骨架：範本、README、dev-remarks、`tdp_test.go`，開 repo 時一次產生。
  - **tdptest**：Go 的測試 package，app 的 `go test`（所以 `make check`）會跑；每個檢查標清單編號，失敗時指出是哪一條。原本定
    分兩層（user 選 A：terminu-go 提供測試用的 driver，在測試裡開 app、送鍵、拿畫面，像 `httptest`；harness 的 tdptest 建在它上面），
    後來改成全部放在 tgo（見下面「tdptest 全部放在 terminu-go」）。plugin 交出去的是骨架裡的 `tdp_test.go`。
  - **reviewer**（agent）：只看這次的 diff，判 tgo 的用法與程式判不了的事。
  - reportgap：開發者與 reviewer 共用這一個。
  - MCP 終端機：讓 agent 開發時看得到自己寫的畫面，不再拿來判。
- **reviewer 什麼時候跑**（user 2026-10-20）：「應該不管是不是 newapp, 有兩種方式可以叫 reviewer, 1 個是 app 發版前， 1 個是
  手動呼叫」。不在 newapp 每寫完一塊就自動跑。怎麼接上（user 選 A）：harness 加 `/terminu:release`，先叫 reviewer 看上一個 tag
  到現在的 diff，有附出處的問題就停下來修（同一塊最多來回兩次），沒有才問使用者要不要打 tag、push；手動叫用 `/terminu:review`，
  預設看還沒 commit 的改動加上這個 branch 還沒 push 的 commit，也可以指定範圍。理由：打 tag 本來就先問使用者，範本檔不用改、不靠
  CI。兩個都在 app 自己的 session 裡打（plugin 本來就是這樣用：誰打 `/terminu`，就由那個 session 走流程）。
- **reviewer 的做法**（我提的，user 看過沒有異議）：兩個 skill 派 `terminu:reviewer` agent，只讀、從乾淨的 context 開始，跟寫程式的
  一方分開。讀 diff、`docs/spec.md`、dev-remarks 的「偏離 tdp」、`docs/<app>-terminu-fix.md`、diff 碰到的元件的 tdp 檔與清單條目；
  tgo 有了涵蓋清單之後也讀它，tgo 保證的條目跳過。每個問題寫哪個檔哪一行、什麼問題、出處（tdp 條目或 spec 的哪一句）、怎麼改；
  沒出處的只是建議。主 session 修掉再叫它看修的那一段，最多兩次；還不同意的問使用者：修掉、寫成「偏離 tdp」（附理由）、或回報 tdp。
  不另外記「駁回過」：定過的都留在 repo（spec、偏離 tdp、fix 文件），reviewer 下次讀得到。
- **試整套流程時在 app 自己的 session**（user 選 A）：在 app 的資料夾開新的 session、只帶 plugin（例
  `claude --plugin-dir ~/Documents/sideproj/terminu-harness`），從 `/terminu` 開始走；harness 的 session 只改 harness，看那邊卡在
  哪裡再回來修 skill。理由：harness 的 session 知道太多（CLAUDE.md、tdp 的細節、user 講過的話），會不自覺替 skill 寫不清楚的地方
  補上，測不出來。取代 2026-10-18 的「過程中不會開啟新的 session，由你直接進行整個流程」。
- **`/terminu:release` 包整個發版**（user 選 A）：工作目錄乾淨、`make check`、reviewer、CHANGELOG 的 `## Unreleased` 改成
  `## vX.Y.Z — tdp <版本>`（tdp project-29），問使用者版號與要不要打 tag，同意才 tag、push（push 觸發 `release.yml`）。
- **三個環**：內環（寫一塊 → `make check`、看畫面 → 有錯就改，全過就 commit；reviewer 在發版前與手動叫時跑）；外環（回報 → tdp、tgo 發版 → harness
  跟上 → app 回來拿掉 TODO）；學習環（同一種錯出現第二次 → 做成 tdptest 的檢查，或請 tgo 改 API 讓錯的寫法寫不出來）。
- **reviewer 不跟開發無限來回**（我提的，user 沒有異議）：只看這次的 diff；每個問題要指出 tdp 條目或 spec 的哪一句，指不出來的
  只是建議、不擋；判斷不下來的（tdp 沒寫清楚）發 issue、標 TODO、往下做；寫進 spec 或 dev-remarks「偏離 tdp」的不再提；同一塊
  最多來回兩次，還沒解就問使用者。
- **tdptest 第一版**（user 選 A）：driver（每次畫面更新都自動檢查 L4；一個 helper 讓整段測試照三種 icon 寬度各跑一次，E5）、
  檔案檢查（範本 diff、README 的 12 段、dev-remarks 的 8 節、`tdp` 前綴、splash 與 `docs/icon.svg` 一致）、第一次啟動的掃描（五種
  尺寸下走過每個 panel 的 Space menu、global operation popup、`?`，再縮到比最低尺寸小，L1）。掃壞檔（把設定目錄的檔改壞再開）、
  自動打開每個動作放第二版。依據：hostu 第一次 judge 自己的 10 條違反，references-1 靠檔案檢查、E6-4 與 E5-11 靠 driver 量每一格、
  panel-77 與 78 靠掃壞檔；popup-69、note-69、toast-18 歸 reviewer；L2-1 該請 tgo 把 statusbar 改成固定寬度的欄位。
- **app 怎麼交給 tdptest**（user 選 A）：app 給一個不帶參數的建構函式 `func() *terminu.App`，位置照使用者啟動時一樣從環境變數讀；
  tdptest 先把 `HOME`、`<APP>__CONFIG`、`<APP>__DATA`（名稱從 `App.Name` 推）指到暫存目錄再呼叫，順便驗了 E3 的覆寫。app 自己的
  系統檔用 `Env`（環境變數 → 沙盒裡的路徑）與 `Seed`（放初始的檔）補。`tdptest.Check(t, tdptest.App{Repo, New, Env, Seed})` 跑
  檔案檢查與掃描；`tdptest.Start`、`tdptest.EachIconWidth` 給 app 自己的流程測試用。
- **tdptest 全部放在 terminu-go**（user 選 A，取代上面的「分兩層」）：driver、畫面檢查、專案檔檢查都在 tgo（例
  `github.com/vulcanshen/terminu-go/tdptest`）。user：「既然 lib 是 tgo, 那 tgo 自己作一個檢測 api, 應該更簡單? harness 直接使用就好，
  這樣 harness 也不用追版本，反正 app 用哪一版的 tgo, 就會跟著哪一版的 tgo test」。好處：檢查跟著 tgo 對應的那一版 tdp、不用另外
  發給 app（harness 還沒上 GitHub 也沒關係）、不用為了檢查開放 tgo 的內部、學習環的兩條路都在 tgo。代價：harness 不改 tgo，要加
  檢查就寫成回報、等 tgo 發版；檢查與實作是同一方寫的，tgo 讀錯 tdp 的地方，檢查也會跟著錯。上面「tdptest 第一版」「app 怎麼交給
  tdptest」照舊，變成 harness 向 tgo 提的需求。
  - **harness 提刁鑽的值**：開發 app 時踩到的值（控制字元、ESC 序列、Tab、寬字元、超長字串、空值、極小的終端機、一次送一串按鍵）
    整理成一份，經 reportgap 交給 tgo 收進自己的測試資料，每個元件都跑一遍；踩到新的再補（學習環）。依據：hostu 回報的 HU1～HU12
    有一半是這種（HU2、HU4、HU8、HU10、HU11、HU12）。「該長什麼樣」照舊由 tgo 照 tdp 判，所以值抓得到漏掉（HU11），抓不到讀錯
    （HU12），讀錯靠 reviewer 與回報 tdp。user 問過由 harness 連期望結果一起給：那等於 harness 又在寫檢查，要約定格式、追版本，不做。
  - **涵蓋清單**：tgo 的測試標清單編號，產生一份「哪些條目有測試」的清單（放進 tgo 自己的 `make check`），漏測的一看就知道。
    harness 不另外驗 tgo（驗 lib 不是 harness 的事，就像做網站的 harness 不驗 React），只讀這份清單，讓 reviewer 只看剩下的。
  - 不在測試裡重算版面（border 幾格、label 幾格）：那是重寫一份 tgo 的計算，只驗得出 tgo 跟自己一致；只量畫出來的結果（每一列的
    寬度、某個位置的字、某一格的顏色）。driver 在記憶體裡跑 `Update`、`Frame()`，不開終端機。
  - **現況**：terminu-go v0.3.0（2026-10-20）發了 tdptest 第一版：driver（`Start`、`Press`、`Type`、`Advance`……）、
    `EachIconWidth`、`Check`、`Tricky`／`Sizes`／`Mashes`。跟 harness 提的不同處：`Env` 是 `func(home string) map[string]string`；
    檔案檢查（`CheckRepo`）與標了清單編號的涵蓋清單排在下一版。app 自己的計時器要用 `App.Tick` 測試才推得動（README 的
    「計時器：`App.Tick`」）。做的時候抓到 3 個 lib bug（含 hostu 踩過的「打 `Home` 被當成 Home 鍵」：KeyName 把一串字當鍵名）。
- **跟 terminu-go 直接講**（user 2026-10-20：「tgo 那邊你可以直接跟他 session 溝通」）：給 terminu-go 的需求與回報，在本機直接
  用 SendMessage 跟 terminu-go 的 session 談，不用先給 user 看草稿。tdp（terminu 的 session）照舊先給 user 看。
- **terminu 的請求與通知直接照做**（user 2026-10-21：「可以接受 terminu 那邊的操作請求和通知」）：terminu 的 session 請 harness
  做的事（跟上新版、改 harness 的檔、commit）不用再問 user，做完回報 terminu。user 2026-10-26 確認範圍包括改 harness 自己的
  CLAUDE.md。建 GitHub repo、push、打 tag、發布照舊先問 user。
- **拿掉 judge**（user 選 A；2026-10-20：「judge 現在就可以拆了」）：成績單、十批 judge agent、`judge_*` 工具都不要了。要看一個現成的 tgo app：跑 tdptest，再讓
  reviewer 看整個 repo。2026-10-20 拆了：`skills/judge`、`agents/judge.md`、`internal/run`、`internal/scan`、`internal/appgit`、
  清單解析、`judge_*` 工具、`terminu-harness report`；MCP 終端機與 `tdp_file` 留著。程式在 `fb89ad1` 還找得到；舊的成績單留在
  `~/.terminu-harness/judge/`（hostu 的 progress.md 有引用）。
- **不留 gene**（user 選 A）：app 只看通過或不通過 —— `make check` 全過，而且 reviewer 沒有附出處的問題。harness 可以另外記
  tdptest 涵蓋了清單的幾條，量的是 harness 自己的進度，不是 app 的分數。
- **外環怎麼接上**（user 選 A）：harness 在 terminu 通知發版時跟上（照舊）；app 那邊，目前目錄有 `progress.md` 時，每次打
  `/terminu` 不用選「繼續」就先查：回報有沒有回應、tgo 有沒有新版（tdptest 跟著 tgo）、程式裡還剩哪些 `TODO(tdp #N)`、`TODO(tgo #N)`，
  講給使用者聽再接著做。不排定時的 routine。
- 2026-10-20 在本機告訴了 terminu 的 session（user 選 A），請它照這些更新
  `~/Documents/sideproj/terminu/.local/harness/README.md`。

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
  （2026-10-20 改：judge 拿掉；不用選「繼續」，打 `/terminu` 就先查，見「方向」。）

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
  `.local/bin/terminu-harness` 就用，沒有就用 Go 從原始碼建到 `~/.terminu-harness/bin/`，所以現在使用者要有 Go。發版要不要改成
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
