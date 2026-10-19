---
name: new-app
description: 照 terminu design（tdp）用 terminu-go 做一個 TUI app：問需求、寫規格給使用者確認、照 tdp 的 project 範本開 repo、一塊一塊寫、卡住就回報 tdp 或 terminu-go、做完自動 judge 並報告。也用來繼續做到一半的 app。由 /terminu 帶進來，或 user 要做一個 terminu 家族的 app 時使用。
argument-hint: [app 的需求]
---

# 做一個 terminu app

規定一律在 tdp，這裡只有流程。寫程式時照 tdp 的原文，不憑印象；不確定就讀那個元件的檔。

目前目錄（或使用者指的 repo）已經有 `.local/terminu/progress.md`，就是**繼續做**：跳到「6. 繼續」。

## 0. 準備

- **tdp**：用 `tdp_file`（不給版本就是最新 tag）讀 `principle/README-zh_TW.md`、`principle/rules-zh_TW.md`、
  `components/README-zh_TW.md`、`project/README-zh_TW.md`、`checklist/README-zh_TW.md`（「適用」表就是功能的清單）。用到哪個元件
  再讀那個元件的檔（`components/layout/panel-zh_TW.md`、`components/input/textarea-zh_TW.md`……）。
- **terminu-go**：最新的 tag（`gh api repos/vulcanshen/terminu-go/tags --jq '.[0].name'`），讀它那一版的 README 與 examples。
- **使用者給的參考**（舊專案、文件）先讀完。
- 確認有 `go`、`git`、`gh`（`gh auth status`）。沒有 `gh`，回報那一步改成把 issue 內容給使用者自己貼。

## 1. 問需求

- **只問需求，不問設計**：使用者只提需求（要做什麼、給誰用、在哪些平台、哪些事不能發生）；資料模型、一筆設定的單位、畫面與
  panel、按鍵、檔案格式這些設計由 harness 照 tdp 與好的 TUI UX 自己決定，寫進 spec 並附理由（user 2026-10-19：「這個應該要
  harness 自己決定，使用者只是提出需求」）。
- 需求不清楚時才問，一次問一題（AskUserQuestion），用問句寫明要決定什麼，列 2–4 個選項，推薦的排第一個並寫理由。
- TUI 不是 CLI：別把 CLI 的做法搬進來（例：讓使用者直接編一個設定檔）。用 TUI 給得起的 UX：清單、明細、一筆一筆加與改的表單；
  app 自己的資料存成 app 自己的格式，要給別的程式的檔（例：hosts 檔）是匯出出去的。
- 每一題都對到 tdp：需求落在哪個元件（panel、明細 panel、viewer、textarea、confirm、toast、terminal popup……）。選元件前先讀那個
  元件的「用途」段：它寫了什麼該用、什麼不該用（例：設定檔交給 `$EDITOR`，不用 textarea）。跟 tdp 衝突的
  當場說出來，讓使用者選「改需求」或「回報 tdp」；tdp 沒規範的記下來，開發時回報。
- **不問 tdp 已經規定的事**：按鍵、外框、顏色、popup 的開關、splash、README 與 dev-remarks 的結構，照 tdp 做。
- tdp 綁著 vulcanshen 的規定（module 路徑、Homebrew tap、splash 的信箱，見 tdp 的 HN3）：使用者不是 vulcanshen 時先問怎麼辦。

## 2. 規格與最終確認

寫 `docs/spec.md`（繁中），至少有：

- 一句話定位、它能做什麼
- 版面：每個畫面的 ASCII 圖，每個 panel 的內容、`Enter` 做什麼、Space menu 有哪些動作
- popup：每個動作打開什麼、確認什麼
- 資料：存什麼、存在哪（tdp E3 的位置）、格式
- 權限、外部指令、各平台的差別
- tdp 沒規範、要回報的地方
- 里程碑：一塊一塊、每塊做完看得到東西

給使用者看，問「照這份開始做嗎？」，改到使用者確認。**確認前不寫程式。**

## 3. 開 repo

- 放在哪裡問使用者。repo 照 tdp `project/README` 的目錄。
- 範本檔（`family.mk`、`release.yml`、`goreleaser.yaml`、`install.sh`、`uninstall.sh`、`version.go`、`gitignore`）用 `tdp_file`
  抓同一版，照「範本檔」表複製，只改表上最後一欄允許的地方。
- `Makefile`、`CHANGELOG.md` 的開頭照 project 頁；兩份 README 照 Rules E7 的 12 段；`docs/dev-remarks.md` 照 E7 的八節，
  「設計文件導讀」指到 `docs/spec.md`；`docs/icon.svg` 照 Rules S5（U 框照 tdp 的 `principle/icon.svg`，金色的字拼 app 的名字）；
  `LICENSE`。
- `go mod init github.com/vulcanshen/<app>`，`go get github.com/vulcanshen/terminu-go@<最新 tag>`，不用 `go.work`。
- `git init`，第一個 commit（不加 Claude 署名）。**建 GitHub repo、push 先問使用者。**

## 4. 一塊一塊寫

照 spec 的里程碑，每一塊：

1. 寫程式與測試（tdp 要的測試：L4 跨尺寸、E5 兩種寬度、E6-5 splash 與 `docs/icon.svg` 一致）。
2. `make check` 過。
3. 用 harness 的終端機開起來看（`open`、`keys`、`screen`、`image`），照這一塊相關的清單條目自己走一遍。app 會寫系統檔或使用者的
   檔案時（例：hosts 檔、家目錄裡的設定），用 app 自己的環境變數（tdp E3 的覆寫）指到暫存目錄再開，不碰真的檔案。
4. 把畫面給使用者看，講做了什麼。
5. commit，更新 `.local/terminu/progress.md`。

`progress.md`（git 不收）寫：照的 tdp 與 terminu-go 版本；spec 的里程碑勾到哪；「回報」：每張 issue 的網址、程式裡的標記、
在等什麼；下一步。

## 5. 卡住

- 分三種：tdp 沒規範或寫不清楚 → 回報 tdp；terminu-go 做不到或有 bug → 回報 terminu-go；app 自己的錯 → 自己修。
- 跟使用者說明卡在哪、暫時打算怎麼做，用 Skill 叫 `terminu:report-gap` 寫 issue，使用者同意才發。
- **其他部分照做**：卡住的地方先用暫時的做法，程式裡標 `TODO(tdp #N)` 或 `TODO(tgo #N)`，記進 progress.md；tdp 或 terminu-go
  更新後回來改掉。

## 6. 繼續

1. 讀 `.local/terminu/progress.md` 與 `docs/spec.md`。
2. 「回報」裡的每一張：照當初送的方式看有沒有回應（GitHub issue 用 `gh issue view <網址>`；在本機送給別的 session 的，看那邊的
   回覆或紀錄）。tdp 或 terminu-go 出了新 tag：terminu-go 用 `go get` 升級；tdp 新版先讀
   它的 CHANGELOG 看改了什麼。把對應的 `TODO(tdp #N)`、`TODO(tgo #N)` 改掉，記進 progress.md。
3. 接著做沒做完的里程碑（回到 4）。

## 7. 完成

- 里程碑都做完、`make check` 過、剩下的 TODO 都在等 issue → 用 Skill 叫 `terminu:judge`，app 指令是建好的執行檔，repo 是這個
  repo。app 會寫系統檔或使用者的檔案時，給 judge 一個包裝指令：每次啟動開一個新的暫存目錄，把 `HOME` 與 app 的環境變數指進去，
  並在補充裡寫明（judge 會轉給每個 judge agent）。
- 報告給使用者：gene 數、違反的條目、判不了的原因、tdp 沒規範的東西、還在等的 issue。
- 違反的條目問使用者要不要修；修完再判一次。
- 打 tag、發布、Homebrew 都先問使用者。
