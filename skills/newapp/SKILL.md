---
name: newapp
description: 照 terminu design（tdp）用 terminu-go 做一個 TUI app：把使用者的需求寫進 idea.md、對到 tdp 寫成 spec.yaml、照 tdp 的 project 範本開 repo、一個 task 一個 task 寫並記進 work.yaml、卡住就回報 tdp 或 terminu-go、做完給使用者看。也用來繼續做到一半的 app、加新需求、跟上新版的 tdp。由 /terminu 帶進來，或 user 要做一個 terminu 家族的 app 時使用。
argument-hint: [app 的需求]
---

# 做一個 terminu app

規定一律在 tdp，這裡只有流程。tdp 是宇集：app 的每個決定，都要在 tdp 找得到根據 —— 照規定做，或在 tdp 明文寫「app 可以自行決定」
的範圍裡（tdp P5）。找不到、寫得不清楚、兩處矛盾，回報 tdp，不自己決定。寫程式時照 tdp 的原文，不憑印象；不確定就讀那一節。

harness 的紀錄放在 app 的 `harness-assets/`：`idea.md`、`spec.yaml`、`work.yaml`、`review.yaml`。格式照這個 skill 資料夾裡的
`harness-assets.md`，動手前先讀它。

目前目錄（或使用者指的 repo）已經有 `harness-assets/spec.yaml`，就是**繼續做**：跳到「6. 繼續」。

## 0. 準備

- **版本**：terminu-go 用最新的 tag（`gh api repos/vulcanshen/terminu-go/tags --jq '.[0].name'`）；tdp 用那一版 terminu-go 對到的版本
  （terminu-go 那一版 CHANGELOG 的標題 `## vX.Y.Z — tdp vA.B.C`）。之後讀 tdp 一律指定這個版本，不讀最新的。
- **tdp**：用 `tdp_file` 讀 `principle/README-zh_TW.md`（P5「明文放行」寫了 tdp 管什麼、哪些 app 可以自行決定）、
  `principle/rules-zh_TW.md`、`components/README-zh_TW.md`、`project/README-zh_TW.md`。用到哪個元件再讀那個元件的檔
  （`components/layout/panel-zh_TW.md`、`components/dialog/form-zh_TW.md`……）。每一節的標題前面就是它的編號（`panel.5`）。
- **terminu-go**：讀那一版的 README 與 examples；`README-zh_TW.md` 這幾節一定要讀，照它寫的用 lib：「文字與顏色」（上色給 span，
  不把顏色碼塞進字串；外面來的字自己畫時用 `Escape`）、「App Log」（使用者執行的每個 operation，只看不改的也算，與背景做的事記進去；tdp note.7 寫了哪些不記）、「讀不了自己的檔時」、「設定」、「測試 app」
  （`tdptest`）、「計時器：`App.Tick`」。
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
- 使用者不是 vulcanshen 時，app 是「照 tdp 做的 app」，不是家族成員（tdp Rules E 章開頭）：repo 的路徑、Homebrew tap、splash
  的作者信箱換成使用者自己的，README 的 family 段寫「照 terminu design 做」並附連結、不列家族成員。問需求時一起問使用者的
  帳號、tap 與信箱。
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

- 放在哪裡問使用者。repo 照 tdp project.3 的目錄，加上 `harness-assets/`。
- 範本檔（`family.mk`、`release.yml`、`goreleaser.yaml`、`install.sh`、`uninstall.sh`、`version.go`、`gitignore`）用 `tdp_file`
  抓同一版，照 project.2 的表複製，只改表上最後一欄允許的地方。
- `Makefile`、`CHANGELOG.md` 的開頭照 project 頁；兩份 README 照 Rules E7；`docs/icon.svg` 照 Rules S5（U 框照 tdp 的
  `principle/icon.svg`，金色的字拼 app 的名字）；`LICENSE`。
- `harness-assets/`：`idea.md`（1. 記的需求）、`spec.yaml`（2. 寫的）、`work.yaml`（`tasks: []`）。
- `.claude/rules/terminu.md`：讓之後在這個 repo 寫程式的人（包括別的 session）照 harness 的規矩，內容照抄這個 skill 資料夾裡的
  `terminu-rule.md`。

- `go mod init github.com/vulcanshen/<app>`，`go get github.com/vulcanshen/terminu-go@<0. 定的 tag>`，不用 `go.work`。
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
2. `make check` 過。
3. 用 harness 的終端機開起來看（`open`、`keys`、`screen`、`image`），照這個 task 引的 tdp 自己走一遍。app 會寫系統檔或使用者的
   檔案時（例：hosts 檔、家目錄裡的設定），用 app 自己的環境變數（tdp E3 的覆寫）指到暫存目錄再開，不碰真的檔案。
4. 更新 `work.yaml`：這個 task 的 `spec`、`implementation`、`test`（`harness-assets.md` 的寫法）。寫不成測試的，請使用者確認，記
   `user: {date, how}`。過程中改了設計，同時改 spec.yaml 的 remark。
5. 跑 `check`，有問題就改到沒有。
6. 把畫面給使用者看，講做了什麼。
7. commit（不加 Claude 署名），程式與 harness-assets 一起。

## 5. 卡住

- 分三種：tdp 沒寫、寫不清楚、兩處矛盾 → 回報 tdp；terminu-go 做不到或有 bug → 回報 terminu-go；app 自己的錯 → 自己修。
- 跟使用者說明卡在哪、暫時打算怎麼做，用 Skill 叫 `terminu:reportgap` 寫 issue，使用者同意才發。
- **其他部分照做**：卡住的地方有合 tdp 的暫時做法就先照它做，程式裡標 `TODO(tdp #N)` 或 `TODO(tgo #N)`，spec.yaml 那個 task 的
  remark 記 `{note: 暫時怎麼做, issue: tdp #N}`；連暫時的做法都沒有，那個 task 先不做，其他照做。tdp 或 terminu-go 更新後回來改掉。

## 6. 繼續

1. 讀 `harness-assets/` 的四個檔，跑 `check`：問題先改掉；「待做」就是接下來要做的。
2. 先查在等的事，不用使用者開口：
   - **新版**：terminu-go 有沒有比 go.mod 新的 tag（每一版 tdp 都有一版 terminu-go 對到它）。有的話跟使用者說新版對到哪一版
     tdp、tdp 的 CHANGELOG 改了什麼，問要不要升（AskUserQuestion）。要升就照下面的「升版」做。
   - **回報**：spec.yaml remark 裡的每個 issue，照當初送的方式看有沒有回應（GitHub issue 用 `gh issue view`；在本機送給別的
     session 的，看那邊的回覆或紀錄）。
   - **TODO**：程式裡還剩哪些 `TODO(tdp #N)`、`TODO(tgo #N)`（`grep -rn 'TODO(t' --include='*.go' .`）。
3. 把查到的講給使用者聽；回報有了答案的，改掉暫時的做法、更新 remark 與 work.yaml。
4. 接著做「待做」的 task（回到 4）。

**升版**：

1. `go get github.com/vulcanshen/terminu-go@<新的 tag>`，`go mod tidy`。中間隔了好幾版（跨 minor 也一樣）也一次升到新的 tag。
2. 把 spec 的 `to` 到新版之間 tdp 每一版的 CHANGELOG 統整成一份：同一節改過好幾次只看最後的寫法，加了又拿掉的不算，只改
   terminu-go 畫法、app 不用動的拿掉。每一版的 CHANGELOG 寫了改到的編號。terminu-go 舊版之後每一版的 CHANGELOG 也讀：寫了 app
   要改的（例：「An app that logs its own filter takes that out」）一起放進這一份。
3. 開新一份 spec：id 加一，`from` 改成原本的 `to`，`to` 改成新版。拿統整的那一份對每個 task：引到改了的節、或做的事被新的規定
   管到的，照新版改 task（`tdp`、remark），`spec` 改成新的 id。idea.md 同時有改的也在這一份裡處理。對不到任何 task 的（project
   的檔、README、CHANGELOG 的寫法）不開 task，直接照新版改，`make check` 過了 commit。
4. 回到 4，把 `spec` 改了的 task 照新版做一次。

**使用者加需求、改需求**：寫進 idea.md，照「升版」的 3. 開新一份 spec（`from`、`to` 不變），新增或修改 task，回到 4。

## 7. 完成

- `check` 的「待做」只剩在等回報的 task、`make check` 過 → 把成品給使用者看：做了什麼、怎麼開、還在等的回報與 TODO。
- 驗不了的（例：要一台 Windows 才驗得了）在這個 session 裡問使用者，確認了記進 work.yaml 的 `user`。
- 要發版時用 Skill 叫 `terminu:release`：它先 review，打 tag、push 前問使用者。Homebrew 也先問使用者。
