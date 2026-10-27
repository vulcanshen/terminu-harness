---
name: analysis
description: 接手一個還沒用 harness 管過的 terminu-go app：從 README 寫 idea.md 與 spec.yaml、從 git commit history 寫 work.yaml，讓 review 有依據，之後照 harness 的流程走。寫紀錄不讀程式碼（terminu-go 不是最新版時先升，那一段要改程式）；刪掉舊的開發者文件，放上 app 的規則檔，做完告訴使用者。/terminu 問了要接手、或 /terminu:review 找不到 harness-assets 時使用。
---

# 接手一個現成的 terminu app

這個 app 是先用 tdp 與 terminu-go 做的，之後才掛上 harness。這裡替它補出 harness 的紀錄（`harness-assets/`）：功能從 README 來，
怎麼做的、哪些測試作保從 git commit history 來。有了紀錄，review 就有依據，之後照 harness 的流程 review、加需求、發版。格式照
newapp 的 `harness-assets.md`（這個 skill 資料夾旁邊的 `../newapp/harness-assets.md`），動手前先讀它。

- **來源只有 README 與 git history**（加上 `go test -list` 列出的測試名字）。寫紀錄不讀程式碼（先升 terminu-go 時例外：那一段要
  改程式），也不讀舊的開發者文件：接手時的 spec
  只要寫出現在有哪些功能，每個做法的根據是引的 tdp 編號，不是文件裡寫的理由。兩樣都找不到的寫「不知道」，review 會列出來，之後
  寫程式的人補。
- **不停下來確認**：產生完再告訴使用者，使用者要改就跟 harness 說。

## 1. 先看

- `go.mod` 沒有 `github.com/vulcanshen/terminu-go`：不是 harness 管的 app，跟使用者說，結束。
- 已經有 `harness-assets/`：不用接手，改叫 `terminu:review`。
- 工作目錄不乾淨（`git status --porcelain` 有輸出）：停下來，問使用者要先 commit 還是先收起來。
- **版本**：terminu-go 是 go.mod 裡的版本；tdp 是那一版 terminu-go 的 CHANGELOG 標題寫的（`## vX.Y.Z — tdp vA.B.C`；CHANGELOG 在
  `go list -m -f {{.Dir}} github.com/vulcanshen/terminu-go` 那個目錄）。之後讀 tdp 一律指定這一版。
- **先升 terminu-go**：go.mod 的版本不是最新的 tag（`gh api repos/vulcanshen/terminu-go/tags --jq '.[0].name'`）時，問使用者要不要
  先升（AskUserQuestion，推薦升）：spec 照哪一版 tdp 寫，就要引那一版的編號，照舊版寫完馬上升版等於做兩次。tdp 早於 v0.10.0、或
  CHANGELOG 沒寫對到哪一版 tdp 的，一定要先升：那時 components 的節還沒有編號，spec 引不了。使用者不升而 tdp 早於 v0.10.0：跟
  使用者說接手不了，結束。升（這一段要改程式，analysis 只有這裡讀程式）：
  1. `go get github.com/vulcanshen/terminu-go@<最新的 tag>`、`go mod tidy`。
  2. 用 `tdp_file` 讀新版 tdp 的 `CHANGELOG.md`，把舊版之後到新版的每一版統整成一份，照 newapp 的「升版」2.（`../newapp/SKILL.md`）：
     同一節改過好幾次只看最後的寫法，加了又拿掉的不算，只改 terminu-go 畫法、app 不用動的拿掉。
  3. `make check` 過：不過就照 terminu-go 舊版之後每一版的 CHANGELOG 與失敗訊息修到過。
  4. 統整剩下的是測試蓋不到、app 自己要改的事（app 給的字、README、CHANGELOG、project 的檔、app 自己的 popup 與流程）：一條一條
     讀相關的程式與檔案，沒照新版的就改，改完再跑一次 `make check`。改不了的、tdp 寫不清楚的記下來，6. 跟使用者說。
  5. commit（不加 Claude 署名），訊息寫升到哪一版，列出 4. 改了哪幾條。

## 2. 從 README 寫 idea.md 與 spec.yaml

讀 `README-zh_TW.md`（沒有就讀 `README.md`）：一句話、它能做什麼、特色、快速開始、使用方式與按鍵、設定與資料存放位置、限制。

- **`idea.md`**：這個 app 要做什麼、給誰用、在哪些平台、哪些事不能發生，從 README 的一句話、它能做什麼、限制整理。開頭寫一行
  「由 analysis 從 README 整理（<日期>），不是使用者的原話」。
- **`spec.yaml`**：`id: 1`、`from: 0`、`to` 是 1. 的 tdp 版本。
  - task 從 README 寫的功能、使用方式與按鍵、設定與資料切（`harness-assets.md` 的「一個 task 一個需求」），`request` 照 README 的
    說法寫，`spec: 1`。只寫現在有的功能，沒有取消的 task。
  - `tdp`：用 `tdp_section`、`tdp_file` 讀那一版 tdp，引每個 task 照的編號；在放行範圍裡做的引那個放行與範圍內還要守的規定。
    README 寫的做法跟 tdp 不合的，照樣引它該照的編號，review 會判出來。
  - remark：README 寫得出來的設計（畫面的配置、每個 panel 放什麼、資料存在哪）寫進對應的 task。
  - tdp 找不到根據的（tdp 沒寫這件事）：`tdp` 留空，remark 寫 `note: 找不到 tdp 根據：<缺什麼>`，6. 跟使用者說。

## 3. 從 git history 寫 work.yaml

每個 task 一筆，`spec: 1`。

- **找做這個 task 的 commit**：`git log --reverse --format='%h %ad %s' --date=short` 看整段歷史；找某個功能用
  `git log --oneline -i --grep='<詞>'`，或用畫面上的字 `git log --oneline -S'<字>'`（只看列出來的 commit 與訊息，不讀 diff）；
  要看 commit 的完整訊息用 `git show -s --format=%B <hash>`。
- **`implementation`**：從 README 的按鍵、使用方式與 commit 訊息寫（觸發、經過、讀寫、用到的 terminu-go 元件、畫面上的字）。
  兩樣都找不到的部分寫「不知道：<缺什麼>」。
- **`test`**：
  1. `go test -list . ./...` 列出測試的名字。
  2. 每個測試用 `git log --reverse --format='%h %s' -S'func <TestX>(' -- '*_test.go'` 找加進它的 commit；那個 commit 做的是哪個
     task，這個測試就寫在那個 task。名字一看就知道在測哪個 task 的也算。對不上的不硬配。
  3. 沒有測試、但找得到做這個 task 的 commit：寫 `commit: <hash>`。
  4. 都沒有就空著，review 會判沒有測試。

## 4. 放規則檔

`.claude/rules/terminu.md`：讓之後在這個 repo 寫程式的人（包括別的 session）照 harness 的規矩，內容照抄 newapp 的
`terminu-rule.md`（`../newapp/terminu-rule.md`）。已經有這個檔就換成這一份。接手以後它就是 harness 管的 app，跟新 app 一樣。

## 5. 刪掉舊的開發者文件

tdp v0.10.0 起不管開發者文件（tdp P5），harness 管的 app 改用 harness-assets，舊的開發者文件接手時刪掉（user 在 terminu 定的）。
內容 git history 都找得到。

- **只刪開發者文件**：`docs/dev-remarks.md`，與 `docs/` 裡寫給開發者看的設計文件（例：`spec.md`、`ui.md`、`ux.md`、
  `function.md`、`<app>-ui-design.md`）。看檔名與開頭幾行判斷是不是。
- **給使用者看的、程式會用到的都留著**：
  - `docs/icon.svg`（splash 的正本，tdp E6）、README 用的 gif、`social-preview.png` 一定留著。
  - README 連到的檔都留著（dev-remarks 除外）。例：README 連到的範例資料夾是給使用者的。
  - Makefile 或程式會產生、會讀的檔留著。刪之前用 `grep -rn '<檔名>' Makefile *.mk cmd internal` 看有沒有提到它（只找檔名，不讀
    程式）；有提到就不刪。
  - 拿不準是不是開發者文件的，不刪，6. 跟使用者說。
- README 兩份裡連到刪掉的檔的連結拿掉。
- `.local/terminu/progress.md`（以前的 harness 留下的，git 不收）不刪，跟使用者說它不再用了。

## 6. 驗與 commit

1. 跑 MCP 工具 `check`。格式、id、tdp 編號的問題改到沒有；測試相關的問題留著，交給 review；「沒有 tdp 根據」留到回報了 tdp。
2. commit（不加 Claude 署名），harness-assets、規則檔、刪掉的文件一起，訊息寫「Take the app over with the terminu harness」這類。
3. 告訴使用者：
   - 先升了 terminu-go 的：升到哪一版、為了新版 tdp 改了哪幾條、改不了的卡在哪；tdp 寫不清楚的用 Skill 叫 `terminu:reportgap`
     寫回報（它給使用者看草稿，同意才發）。
   - harness-assets 產生了，之後 review 會拿它判；要改內容就跟 harness 說。
   - 幾個 task；寫了「不知道」的有哪幾個；沒有作保的有哪幾個。
   - tdp 找不到根據的有哪幾個：用 Skill 叫 `terminu:reportgap` 寫回報（它給使用者看草稿，同意才發），發了把 issue 補進那個 task 的
     remark。
   - 刪了哪些檔、拿不準所以留著的檔、progress.md 不再用了。
4. 接著用 Skill 叫 `terminu:review`（從 /terminu 進來接手的，或是 review 叫你的，都回到 review）。
