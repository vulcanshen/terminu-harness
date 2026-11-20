---
name: newapp
description: 照 terminu design（tdp）用 terminu-go 做一個 TUI app：把使用者的需求寫進 idea.md、對到 tdp 寫成 spec.yaml、照 harness 帶的家族 repo 範本（terminu project）開 repo、一個 task 一個 task 寫並記進 work.yaml、卡住就回報 tdp 或 terminu-go、做完給使用者看。也用來繼續做到一半的 app、加新需求、跟上新版的 tdp。由 /terminu 帶進來，或 user 要做一個 terminu 家族的 app 時使用。
argument-hint: [app 的需求]
---

# 做一個 terminu app

規定一律在 tdp，這裡只有流程。tdp 是宇集：app 的每個決定，都要在 tdp 找得到根據 —— 照規定做，或在 tdp 明文寫「app 可以自行決定」
的範圍裡（tdp P5）。找不到、寫得不清楚、兩處矛盾，回報 tdp，不自己決定。寫程式時照 tdp 的原文，不憑印象；不確定就讀那一節。

harness 的紀錄放在 app 的 `harness-assets/`：`idea.md`、`spec.yaml`、`work.yaml`、`review.yaml`。格式照這個 skill 資料夾裡的
`harness-assets.md`，動手前先讀它。

目前目錄（或使用者指的 repo）已經有 `harness-assets/spec.yaml`，就是**繼續做**：跳到「6. 繼續」。

## 0. 準備

- **要用的 terminu-go**（這裡定一次，「繼續」、「升版」、`terminu:release`、`terminu:analysis` 都照這個）：
  - 預設是最新的正式版，跟 Go 的 `@latest` 一樣不算 beta：`go list -m -f '{{.Version}}' github.com/vulcanshen/terminu-go@latest`。
  - app 的 go.mod 已經是 beta（版本帶 `-`，例 `v1.2.0-beta.1`）：算 beta，拿 `go list -m -versions github.com/vulcanshen/terminu-go`
    的最後一個（照 semver 排，正式版排在它的 beta 後面，所以正式版出來就升到正式版）。
  - 使用者明確指定某一版（例：「升到 v1.2.0-beta.1 試試」）：用那一版。
  - beta 是給試的：家族一輪 beta 只有指定的成員跟，沒指定的不要自己升上去。
- **tdp**：用那一版 terminu-go 對到的版本（terminu-go 那一版 CHANGELOG 的標題 `## vX.Y.Z — tdp vA.B.C`）。之後讀 tdp 一律指定這個
  版本，不讀最新的。
- **tdp**：用 `tdp_file` 讀 `principle/README-zh_TW.md`（P5「明文放行」寫了 tdp 管什麼、哪些 app 可以自行決定）、
  `principle/rules-zh_TW.md`、`components/README-zh_TW.md`。用到哪個元件再讀那個元件的檔
  （`components/layout/panel-zh_TW.md`、`components/dialog/form-zh_TW.md`……）。每一節的標題前面就是它的編號（`panel.5`）。
- **家族的 repo 範本**（terminu project）：tdp v1.2.0-beta.1 起不在 tdp，在 harness（tdp E 章開頭、E6 交給這裡）。用 `tdp_section`
  讀 `harness project`（整頁），範本讀 `harness project/install.sh` 這樣的路徑。引用寫 `harness project.3`、`harness project-13`。
- **terminu-go**：讀那一版的 README 與 examples；`README-zh_TW.md` 這幾節一定要讀，照它寫的用 lib：「文字與顏色」（上色給 span，
  不把顏色碼塞進字串；外面來的字自己畫時用 `Escape`）、「App Log」（使用者執行的每個 operation，只看不改的也算，與背景做的事記進去；tdp note.7 寫了哪些不記）、「讀不了自己的檔時」、「設定」、「測試 app」
  （`tdptest`）、「計時器：`App.Tick`」、「從 app 自己的 popup 執行動作」（`App.RunAction`，terminu-go v1.2.7 起）、「開著時會變的內容」
  （`App.ReadAgain`，terminu-go v1.2.14 起）、「成功的報告」（`App.Report`，terminu-go v1.2.16 起）、「Diff」（`App.OpenDiff`，
  terminu-go v1.2.18 起；panel 裡的差異 `Pager.SetDiff`，v1.2.19 起；v1.2.19 以前這一節叫「比較」）、「占好幾列的項目」
  （`Row.Below`，v1.2.19 起）、「命令列」（v1.2.19 起；`<app> update` 的 `SelfUpdate`，v1.3.0 起）、「項目的附註」
  （`Panel.Note`，v1.3.0 起）、「說明牌」（`Panel.Sign`，v1.3.0 起）。
- **使用者給的參考**（舊專案、文件）先讀完。
- 確認有 `go`、`git`、`gh`（`gh auth status`）。沒有 `gh`，回報那一步改成把 issue 內容給使用者自己貼。

## 1. 問需求，寫 idea.md

- **只問需求，不問設計**：使用者只提需求（要做什麼、給誰用、在哪些平台、哪些事不能發生）；資料模型、一筆設定的單位、畫面與
  panel、按鍵、檔案格式這些設計由 harness 照 tdp 與好的 TUI UX 自己決定，寫進 spec.yaml 的 remark 並附理由（user 2026-10-19：
  「這個應該要 harness 自己決定，使用者只是提出需求」）。
- 需求不清楚時才問，一次問一題（AskUserQuestion），用問句寫明要決定什麼，列 2–4 個選項，推薦的排第一個並寫理由。
- TUI 不是 CLI：別把 CLI 的做法搬進來（例：讓使用者直接編一個設定檔）。用 TUI 給得起的 UX：清單、明細、一筆一筆加與改的表單；
  app 自己的資料存成 app 自己的格式，要給別的程式的檔（例：hosts 檔）是匯出出去的。
- 每一題都對到 tdp：需求落在哪個元件。選元件前先讀那個元件的「用途」節：它寫了什麼該用、什麼不該用（例：設定檔交給 `$EDITOR`，
  不用 textarea）。跟 tdp 衝突的當場說出來，讓使用者選「改需求」或「回報 tdp」。
- **不問 tdp 已經規定的事**：按鍵、外框、顏色、popup 的開關、splash、README 的結構，照 tdp 做。
- 使用者不是 vulcanshen 時，app 是「照 tdp 做的 app」，不是家族成員（tdp Rules E 章開頭）：Homebrew tap、splash 的作者信箱換成
  使用者自己的，README 的 family 段寫「照 terminu design 做」並附 tdp 的連結。家族的 repo 範本只綁家族成員（harness project.1；
  tdp 的 E 章開頭：照 tdp 做的 app 的 repo 自己決定）：3. 照樣拿它當起點，範本裡的 `vulcanshen`（repo 的路徑、tap 的 owner、首頁）
  換成使用者的帳號（harness project.2）；之後使用者要改 repo 的目錄、Makefile、CHANGELOG、發布步驟都可以，tdp 的 E 章照樣守。
  問需求時一起問使用者的帳號、tap 與信箱。
- 使用者講的需求照他的話寫進 `harness-assets/idea.md`（repo 還沒開就先記著，3. 開 repo 時寫進去）。

## 2. spec.yaml

照 idea.md 與 tdp 寫第一份 spec：`id: 1`、`from: 0`、`to` 是 0. 定的 tdp 版本。

- 每個需求切成 task（`harness-assets.md` 的「一個 task 一個需求」），`spec: 1`。
- 每個 task 的 `tdp` 寫它照的 tdp 編號；在放行範圍裡做的引那個放行與範圍內還要守的規定。
- 設計決定寫在 remark 並附理由：畫面的配置（ASCII 圖）、每個 panel 放什麼、`Enter` 做什麼、Space menu 有哪些動作、每個動作開什麼
  popup、確認什麼、資料存什麼、存在哪（tdp E3 的位置）、格式、權限、外部指令、各平台的差別。
- 找不到根據的：照「5. 卡住」回報，remark 記 `{note, issue}`。
- task 的順序就是做的順序：一個一個做，每做完一個看得到東西。

寫完就開始做，不停下來等使用者確認（user 2026-10-19：「不應該問我，你應該先做出來，然後我們在來看結果」）：跟使用者講一遍
spec 的重點，使用者看了成品再提要改的。

## 3. 開 repo

- 放在哪裡問使用者。repo 照 harness project.3 的目錄，加上 `harness-assets/`。
- 範本檔（`family.mk`、`release.yml`、`goreleaser.yaml`、`install.sh`、`uninstall.sh`、`version.go`、`gitignore`）用 `tdp_section`
  讀 `harness project/<檔名>`，照 harness project.2 的表複製，只改表上最後一欄允許的地方。
- `Makefile`、`CHANGELOG.md` 的開頭照 harness project.5、project.6；兩份 README 照 Rules E7；`docs/icon.svg` 照 Rules S5（U 框照 tdp 的
  `principle/icon.svg`，金色的字拼 app 的名字）；`LICENSE`。
- `harness-assets/`：`idea.md`（1. 記的需求）、`spec.yaml`（2. 寫的）、`work.yaml`（`tasks: []`）。
- `.claude/rules/terminu.md`：讓之後在這個 repo 寫程式的人（包括別的 session）照 harness 的規矩，內容照抄這個 skill 資料夾裡的
  `terminu-rule.md`。

- `go mod init github.com/<帳號>/<app>`（家族成員是 `vulcanshen`），`go get github.com/vulcanshen/terminu-go@<0. 定的 tag>`，不用 `go.work`。
- 把建 app 的程式寫成 `newApp() *terminu.App`（main 裡 `Run` 以外的全部，位置照使用者啟動時一樣從環境變數讀），在它所在的
  package 放 `tdp_test.go`：

  ```go
  func TestFirstLaunch(t *testing.T) {
  	tdptest.Check(t, tdptest.App{Name: "<app>", New: newApp})
  }
  ```

  app 會讀寫系統檔或家目錄以外的檔（例：hosts 檔）時，加 `Env: func(home string) map[string]string { … }` 把 app 自己的環境
  變數指到沙盒裡，`Seed` 放初始的檔；`Seed` 放的資料裡有 `tdptest.Tricky` 的值時設 `Tricky: true`（README「測試 app」）。
  `make check` 的 `go test` 就會跑它。
- 跑 `check`，沒有問題再 `git init`、第一個 commit（不加 Claude 署名）。**建 GitHub repo、push 先問使用者。**

## 4. 一個 task 一個 task 寫

照 spec.yaml 的順序，每一個 task：

1. 寫程式與測試。測試一律用 terminu-go 的 `tdptest`，不自己寫 driver：
   - 流程測試用 `tdptest.Start`（每一格畫面它都自動量 L4）；要驗 E5 的包在 `tdptest.EachIconWidth` 裡。app 自己會開的 popup、
     toast、terminal popup 都要有流程測試走到（tdp E6、E5）。task 的每個入口都要有測試走到。
   - 使用者給的資料（名稱、路徑、設定裡的字串）用 `tdptest.Tricky` 餵一次。
   - 第一次啟動的掃描是 `tdp_test.go` 的 `Check`。
   - E6 的 splash 與 `docs/icon.svg` 一致的測試還是自己寫。
   - app 自己的計時器用 `App.Tick`，不用 `tea.Tick`，測試才推得動它。
   - app 自己的 popup 執行 app 的動作（tdp popup.8）用 `App.RunAction`，不自己照 menu 重做一遍。從 app 自己的模式裡打開的 popup，
     動作執行時會離開模式，但不跑 `OnExit`：`OnExit` 原本要做的事（例：暫停），與離開會丟掉使用者要留下的東西時的詢問（tdp K11），
     都在呼叫前自己做。
   - 來源變了、要重新讀進開著的 viewer 或 panel 的內容（例：監看的檔案，tdp note.9）用 `App.ReadAgain`、`Pager.ReadAgain`，不用
     `SetLines` 換掉或重新打開：搜尋、選取、跟著結尾才會照 note.9。這種內容不算串流，不設 `Panel.Stream`（tdp T2）。開著時會多出
     新的行的 `Pager`（`Append` 或 `ReadAgain`）要設 `Follow`，結尾在畫面裡時才會跟著往下（tdp note-90）。viewer 的來源讀不到時
     叫 `App.CannotReadAgain(v, <app 的字>)`：Error toast，或選取模式的 Space menu、`?` 開著時疊在上面的 error popup，都由它跳
     （tdp note-100），app 不另外跳。app 自己重讀時，App Log 只記一次（tdp note-104、note-105）：讀到變了用
     `App.Log(terminu.Info, …)` 記一筆、寫出改了什麼，沒變不記；讀不到只在第一次與換了一個失敗時說（viewer 的給
     `CannotReadAgain` 字，其他的跳 Error toast），同一個失敗（同一句話）再發生不跳、不記，viewer 的叫 `CannotReadAgain(v)`
     不給字；又讀得到時記一筆 INFO，不跳 toast。
   - 自己的檔讀得進來、只是有用不了的值或不合規則的一筆（tdp E3）：先用 `terminu.KeepAside` 備份，成功了才用預設值重寫，那一件
     是 `terminu.MisfitNotice(path, copied, misfits…)`；備份做不出來就不重寫，預設值只用這一次，那一件是
     `terminu.Notice{Kind: terminu.Error, Text: err.Error()}`，後面接 `MisfitNotice(path, "", …)`（沒有那一行，只帶 WARN）。
   - 值指到使用者自己資料裡不在的一筆（例：設定指向使用者建的 profile，那一筆被刪掉或改名），不算用不了的值，不備份、不做
     `MisfitNotice`（tdp E3「指到不在的一筆」）：檔案不動、值留著，app 重寫那個檔時也留著，使用者選了別的一筆才改掉；這一次用
     內建的預設值，有好幾筆的值用還在的那幾筆、一筆都不在才用預設值，那一件是 Warning 的 `terminu.Notice`；畫面上照原本寫的
     顯示這個值，label 畫 Red（一列一個值的 panel 用 `Field.Err`），那一筆不在時它的 select 用 `Selected: -1` 打開，選了哪一筆
     都算改。
   - 值指到 app 外面的東西（例：狀態記著的叢集 namespace）不一樣（tdp E3-41）：拿名字就能用的直接用，不等外面回答；外面回答說
     其中一筆不在了，就少了那一筆繼續，那一件是 Warning 的 `terminu.Notice`，下次寫那個檔時把那一筆拿掉，下次啟動就不會再警告。
     找到之前用不了的（例：畫面開了之後才問得到的種類），找完之前先用預設值、不算不在，找完還是沒有才是那一件（tdp E3-42）。
   - 啟動時上面這幾件，連同位置搬了的那句（tdp E3「換了位置時」），照發生的順序收成 `Notice`，叫 `App.Notify(ns…)`：一件一個
     toast，跟所有 toast 一樣往上疊，App Log 一件一筆、WARN 接在後面（tdp E3-37）。啟動後才知道的（外面回答不在了、找完還是沒有），
     知道時再交給 `App.Notify`，一樣是自己的一個 toast。app 不另外叫 `App.Toast`、`App.Log`、`DidNotFit`。沒有畫面時（在終端機
     或 pipe 跑的指令），改把每個 `Notice` 的 `Lines()`（toast 那一行，再接它的 WARN）一行一行寫到 stderr，每行開頭 `<app>: `；讀不了的檔不叫 `MoveAside`，照它的 `Esc` 做：不動檔案，寫
     `terminu.CannotReadText(path, err)`（可以多一行說怎麼辦），exit status 1。TUI 會繼續的地方，指令也繼續（tdp E3-40、E3-27）。
   - TUI 以外（不在 TUI 裡的畫面，例：外掛畫的側邊欄；印給別的程式顯示的輸出，例：狀態列）：上面這幾件與讀不了的檔都不動檔案，
     不叫 `KeepAside`、`MoveAside`，不重寫，這一次用預設值（讀不了時怎麼繼續由 app 決定）；在那裡讓使用者看到哪個檔有什麼問題
     （例：狀態列最後加一個標記），不是只寫 stderr（tdp E3「TUI 以外」）。下次打開 TUI 時才照上面做。
   - 設定檔裡 app 不讀的 key（打錯字，或改版後不再讀的）不是用不了的值：不備份、不收進 `Notice`、不跳 toast，app 重寫那個檔時
     也留著（照 struct 整個重新寫出來會丟掉它）。每次啟動用 `App.Log(terminu.Info, …)` 記一筆，寫出 key 與檔案：app 知道自己
     以前讀過的寫 `No longer used: …`，其他的寫 `Not read: …`（tdp E3-38）。沒有畫面時，這幾行與改版後不再用的檔
     （`App.NoLongerUsed` 記的，`terminu.NoLongerUsedLines(paths...)`）也寫到 stderr，開頭 `<app>: `（tdp E3-40）。
   - 使用者要讀完、toast 又放不下的成功報告（tdp toast.2）用 `App.Report`，不自己開 viewer，也不另外記 App Log 的 Success。
   - 刪除、Diff、附註的鍵整個家族都一樣（tdp K13）：會永久拿掉東西的動作（刪除、丟掉改動、drop），一項的（cursor 上那一項）放 `x`
     （`[x] Delete`），很多的（標記的那些、整批）放 `X`（例：`[X] Delete marked`），不放 `d`、`D`，也不用其他字母；`x`、`X` 也只
     放這種，還能再打開的（例：關分頁）用別的字母（tdp K13-5）。Diff 模式的鍵是 `D`（`DiffMode{Key: "D"}`，lib 沒有預設的鍵）；
     沒有 Diff 模式的 app 可以把 `D` 給不是刪除的事。`i` 只打開附註（下一條），app 自己的動作都不放 `i`（tdp K13-4）。
   - 一項是什麼、欄位又放不下的（例：連結會去哪、備註、為什麼沒在跑）放那一項的附註（tdp panel.18）：設 `Panel.Note`（app 自己
     的模式裡 `Mode.Note`），回傳 cursor 那一項的附註，頭、`i`、打開的 viewer 都是 lib 的；app 自己的 body 要實作
     `terminu.NoteSpot`。不用 toast 說（tdp toast.1），也不用 `Panel.Inspect`、`Mode.Inspect`（deprecated，之後的 minor 版拿掉）。
   - 只用來說明、沒有操作的 panel（例：遊戲旁邊的 NEXT、圖例）做成說明牌（tdp panel.17）：設 `Panel.Sign`（`terminu.SignLit`
     常亮、`terminu.SignDim` 常暗，由 app 挑），不做成拿得到 focus 的一般 panel；內容在最小畫面（60 × 32）分到的大小裡要放得下。
   - 看差異一律經過 lib（tdp panel.9、note.8）：app 自己的動作要看同一項的兩個版本（例：commit 裡的檔對它的上一個 commit）用
     `App.OpenDiff(p, item)`，不開顯示 patch 的 viewer（tdp note-101）；Diff 模式（`App.DiffModeAction`）裡裝著許多項目的列
     （資料夾、commit）給 `DiffMode.Many`、`DiffMode.Items`（舊的名字 `CompareAction`、`Compare` 是 deprecated，之後的 minor
     版拿掉，不用）。有差異的項目清單、標記、`Enter`、`Esc` 都是 lib 的，app 不自己列。清單有跟著 cursor 的明細 panel、要把
     Diff 模式的差異放在明細裡時（放不放由 app 挑），設 `DiffMode.InDetails`（空狀態裡的名詞給 `Noun`），不自己跟著 cursor 在
     明細裡畫：兩個框、`Esc` 一層一層退、空狀態與錯誤狀態都是 lib 的（tdp panel-185～187）。app 自己的動作要在那個 panel 顯示別的
     東西前，先叫 `Panel.LeaveDiff()`。顯示差異的 panel（例：跟著 cursor、顯示那個檔改了什麼的明細）也一樣：
     用 `Pager.SetDiff(anchor, other)` 照 UNIFIED 畫，不放工具自己的輸出（`git diff` 的標頭與 `@@`）；顯示的東西裝著許多項目
     （新的資料夾、commit）用 `Pager.SetDiffItems(items)`，不只寫名字；panel 的 `Enter` 用 `App.OpenDiff` 打開同一個 Diff 或進到
     同一份清單（tdp note-103、note-106）。
   - app 自己的動作要換掉 panel 列的東西，或把它的 cursor 放到某一項上（例：在別的 panel 移動 cursor，把這個 panel 帶到那一項；
     進到一項裡面），先叫 `List.ClearFilter()`，不管篩選有沒有藏住那一項（tdp panel-180）。Diff 模式與 `App.OpenDiff` 進去前 lib
     自己清；只是 focus 離開，篩選留著。
   - 一個項目要占好幾列時（例：commit 下面的分支與 tag）用 `Row.Below`，不拆成好幾個項目：移動、`N of M`、篩選才會算項目
     （tdp panel.5）。
   - 命令列的指令（tdp E1）照 tgo README「命令列」印：結果寫 stdout；錯誤、警告寫 stderr，一件一行，開頭 `<app>: `，跟同一件事的
     toast 同一句話；用錯（不認得的指令、少了參數）寫 `<app>: ` 接哪裡錯、空一行、用法，exit status 2；沒做完 1。`Run` 回傳的
     error 也照這樣（`fmt.Fprintln(os.Stderr, "<app>:", err)`、`os.Exit(1)`）。不用 `log.Fatal`（前面是時間，不是 `<app>: `），
     也不照 `flag` 套件預設的印法（錯誤沒有 `<app>: `，`-h` 的用法寫到 stderr）。`iconwidth` 印 `terminu.Probe("<app>").Lines()`。
     每個 app 都有 `<app> update`（tdp E1-9、E1-10），用 `terminu.SelfUpdate{App, Version, Owner, Tap, Source}.Run(stdout, stderr)`
     的回傳當 exit status，不自己寫：`Version` 是 `version.Version`，`Owner`、`Tap` 照 tdp E6（家族成員是 `vulcanshen`、
     `vulcanshen/tap`，其他的用使用者自己的帳號與 tap），`Source` 寫從原始碼編的怎麼更新（例：`git pull && make install`）。它下載的
     release 檔名是 harness project/goreleaser.yaml 的命名，改了發布步驟的 app 要留著這個命名。照 tgo 的 showcase 寫成
     `command(args, stdout, stderr) int`，有參數時 `main` 把它的回傳交給 `os.Exit`，測試才拿得到印出的東西與 exit status。
2. `make check` 過。
3. 用 harness 的終端機開起來看（`open`、`keys`、`screen`、`image`），照這個 task 引的 tdp 自己走一遍。app 會寫系統檔或使用者的
   檔案時（例：hosts 檔、家目錄裡的設定），用 app 自己的環境變數（tdp E3 的覆寫）指到暫存目錄再開，不碰真的檔案。
4. 更新 `work.yaml`：這個 task 的 `spec`、`implementation`、`test`（`harness-assets.md` 的寫法）。寫不成測試的，請使用者確認，記
   `user: {date, how}`。過程中改了設計，同時改 spec.yaml 的 remark。
5. 跑 `check`，有問題就改到沒有。
6. 把畫面給使用者看，講做了什麼。
7. commit（不加 Claude 署名），程式與 harness-assets 一起。

## 5. 卡住

- 分四種：tdp 沒寫、寫不清楚、兩處矛盾 → 回報 tdp；terminu-go 做不到或有 bug → 回報 terminu-go；家族的 repo 範本（harness
  project）有問題 → 回報 terminu-harness；app 自己的錯 → 自己修。
- 跟使用者說明卡在哪、暫時打算怎麼做，用 Skill 叫 `terminu:reportgap` 寫 issue，使用者同意才發。
- **其他部分照做**：卡住的地方有合 tdp 的暫時做法就先照它做，程式裡標 `TODO(tdp #N)` 或 `TODO(tgo #N)`，spec.yaml 那個 task 的
  remark 記 `{note: 暫時怎麼做, issue: tdp #N}`；連暫時的做法都沒有，那個 task 先不做，其他照做。tdp 或 terminu-go 更新後回來改掉。

## 6. 繼續

1. 讀 `harness-assets/` 的四個檔，跑 `check`：問題先改掉；「待做」就是接下來要做的。
2. 先查在等的事，不用使用者開口：
   - **新版**：「0. 準備」的「要用的 terminu-go」有沒有比 go.mod 新（每一版 tdp 至少有一版 terminu-go 對到它，對到同一版 tdp 的有好幾版時取最新的；go.mod 是正式版時
     beta 不算）。有的話讀 go.mod 那一版之後每一版
     terminu-go 的 CHANGELOG：有一版改了 lib、或寫了 app 的程式或行為要改，就跟使用者說新版對到哪一版 tdp、tdp 的 CHANGELOG 改了
     什麼（對到的 tdp 跟 go.mod 那一版的一樣時，不講 tdp，改說 terminu-go 多了什麼 API、app 要不要改），問要不要升（AskUserQuestion）。要升就照下面的「升版」做。都只改文件的（例：tdp v1.0.1 只改家族成員 README 的一段），
     只跟使用者說有新版、發版前會升，不問（user 2026-10-27：只改文件的版不用馬上升，`/terminu:release` 發版前升到最新）。
   - **回報**：spec.yaml remark 裡的每個 issue，照當初送的方式看有沒有回應（GitHub issue 用 `gh issue view`；在本機送給別的
     session 的，看那邊的回覆或紀錄）。
   - **TODO**：程式裡還剩哪些 `TODO(tdp #N)`、`TODO(tgo #N)`（`grep -rn 'TODO(t' --include='*.go' .`）。
3. 把查到的講給使用者聽；回報有了答案的，改掉暫時的做法、更新 remark 與 work.yaml。
4. 接著做「待做」的 task（回到 4）。

**升版**：

1. `go get github.com/vulcanshen/terminu-go@<要用的版本>`（「0. 準備」的「要用的 terminu-go」），`go mod tidy`。中間隔了好幾版
   （跨 minor 也一樣）也一次升上去。
   - **只升 terminu-go**：新版對到的 tdp 跟 spec 的 `to` 一樣（例：terminu-go v1.2.7 還是 tdp v1.2.6）時，不開新 spec：跳過 2. 的
     tdp 部分和 3.，只讀 go.mod 那一版之後每一版 terminu-go 的 CHANGELOG，寫了 app 要改的才改，`make check` 過了就 commit。
2. 把 spec 的 `to` 到新版之間 tdp 每一版的 CHANGELOG 統整成一份：同一節改過好幾次只看最後的寫法，加了又拿掉的不算，只改
   terminu-go 畫法、app 不用動的拿掉。每一版的 CHANGELOG 寫了改到的編號。terminu-go 舊版之後每一版的 CHANGELOG 也讀：寫了 app
   要改的（例：「An app that logs its own filter takes that out」）一起放進這一份。
3. 開新一份 spec：id 加一，`from` 改成原本的 `to`，`to` 改成新版。拿統整的那一份對每個 task：引到改了的節的，拿新版的原文對它
   做的事；沒引到、但做的事被新的規定管到的也一樣。新的寫法改到它（紀錄要改，或照紀錄做的程式得改）才照新版改 task（`tdp`、
   remark），`spec` 改成新的 id，4. 會照新版重做；改的地方跟它做的事無關的（例：T1 加了一列管搜尋，引 T1 但不搜尋的 task），
   `spec` 不動：每次 review 都拿新版判全部的 task，加 id 只多一次白做。idea.md 同時有改的也在這一份裡處理。
   新版要每個 app 都有的東西（例：tdp v1.3.0 的 `<app> update`，E1-9）放進管那一塊的 task（例：命令列那個 task），沒有的話開一個
   新的 task，`request` 寫 tdp 要的那件事（例：`tdp E1：<app> update`），review 才判得到。對不到任何 task 的（repo
   範本的檔、README、CHANGELOG 的寫法）不開 task，直接照新版改，`make check` 過了 commit。照 tdp 做的 app 不用跟家族的 repo
   範本改（harness project.1）。
   - **第一次升到沒有 project/ 的 tdp**（v1.2.0-beta.1 起）：家族的 repo 範本搬到 harness，編號不變。spec 裡的 `tdp project.3`、
     `tdp project/install.sh`、`tdp project-13` 改成 `harness project.3`、`harness project/install.sh`、`harness project-13`；家族成員
     的範本檔開頭的註解從 `(tdp project/…)` 改成 `(harness project/…)`，照 `tdp_section` 讀到的 harness 範本重新複製；work.yaml
     裡守 repo 的作保（tdptest 的 repo 子測試，例 `TestTDP/repo`，terminu-go v1.2.0-beta.1 起拿掉了）換成 `check: repo`。`check`
     三樣都會指出來。之後範本改了，是 harness 發版時改的，`check` 一樣會比出來。
4. 回到 4，把 `spec` 改了的 task 照新版做一次。

**使用者加需求、改需求**：寫進 idea.md，照「升版」的 3. 開新一份 spec（`from`、`to` 不變），新增或修改 task，回到 4。

## 7. 完成

- `check` 的「待做」只剩在等回報的 task、`make check` 過 → 把成品給使用者看：做了什麼、怎麼開、還在等的回報與 TODO。
- 驗不了的（例：要一台 Windows 才驗得了）在這個 session 裡問使用者，確認了記進 work.yaml 的 `user`。
- 要發版時用 Skill 叫 `terminu:release`：它先 review，打 tag、push 前問使用者。Homebrew 也先問使用者。
