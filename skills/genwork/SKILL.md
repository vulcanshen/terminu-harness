---
name: genwork
description: 接手一個還沒用 harness 管過的 terminu-go app：從 README、CHANGELOG、開發文件（dev-remarks、設計文件）與 git history 產生 harness-assets（idea.md、spec.yaml、work.yaml），不讀程式碼；把開發文件搬進去後刪掉，做完告訴使用者。/terminu 問了要接手、或 /terminu:review 找不到 harness-assets 時使用。
---

# 接手一個現成的 terminu app

這個 app 是先用 tdp 與 terminu-go 做的，之後才掛上 harness。這裡把它現有的文件整理成 harness 的紀錄（`harness-assets/`），之後就照
harness 的流程 review、加需求、發版。格式照 newapp 的 `harness-assets.md`（這個 skill 資料夾旁邊的 `../newapp/harness-assets.md`），
動手前先讀它。

- **不讀程式碼**：來源只有文件與 git history（加上 `go test -list` 列出的測試名字）。文件裡找不到的寫「不知道」，不去程式裡猜；
  review 會把它們列出來，之後寫程式的人補。
- **不停下來確認**：產生完再告訴使用者，使用者要改就跟 harness 說。

## 1. 先看

- `go.mod` 沒有 `github.com/vulcanshen/terminu-go`：不是 harness 管的 app，跟使用者說，結束。
- 已經有 `harness-assets/`：不用接手，改叫 `terminu:review`。
- 工作目錄不乾淨（`git status --porcelain` 有輸出）：停下來，問使用者要先 commit 還是先收起來。
- **版本**：terminu-go 是 go.mod 裡的版本；tdp 是那一版 terminu-go 的 CHANGELOG 標題寫的（`## vX.Y.Z — tdp vA.B.C`；CHANGELOG 在
  `go list -m -f {{.Dir}} github.com/vulcanshen/terminu-go` 那個目錄）。之後讀 tdp 一律指定這一版。CHANGELOG 沒寫這一版對到哪個
  tdp（terminu-go v0.5.0 以前的舊版可能沒有）：跟使用者說，問要不要先升 terminu-go（`go get github.com/vulcanshen/terminu-go@latest`）。

## 2. 讀來源

有的就讀，沒有就跳過：

- `README-zh_TW.md`（或 `README.md`）：一句話、它能做什麼、特色、快速開始、使用方式與按鍵、設定與資料存放位置、限制。
- `CHANGELOG.md`。
- `docs/dev-remarks.md`：運作方式、設計決定、已否決、已知的牆與未做、偏離 tdp、設計文件導讀、建置與開發、發布。
- 設計文件：dev-remarks「設計文件導讀」指到的，和 `docs/` 裡其他的設計文件（例：`docs/spec.md`、`ui.md`）。
- 以前的 harness 留下的 `.local/terminu/progress.md`（回報與在等的事）。
- `git log --oneline`，要找某個功能是哪個 commit 做的時 `git log --oneline -S<字> ` 或 `--grep`。
- `go test -list . ./...`：測試的名字。

## 3. 寫 harness-assets

- **`idea.md`**：這個 app 要做什麼、給誰用、在哪些平台、哪些事不能發生，從 README 的一句話、它能做什麼、限制，與設計文件裡的
  需求整理。開頭寫一行「由 genwork 從 README、CHANGELOG、開發文件整理（<日期>），不是使用者的原話」。
- **`spec.yaml`**：`id: 1`、`from: 0`、`to` 是 1. 的 tdp 版本。
  - task 從 README 的使用方式與按鍵、設計文件、dev-remarks 的運作方式切（`harness-assets.md` 的「一個 task 一個需求」），`spec: 1`。
  - `tdp`：用 `tdp_section`、`tdp_file` 讀那一版 tdp，引每個 task 照的編號。
  - remark：dev-remarks「設計決定」的決定與理由放進對應的 task；「已知的牆」寫成 remark；progress.md 裡的回報寫成 `{note, issue}`。
  - dev-remarks「已否決，不要重提」的每一項寫成取消的 task（`canceled: <當時的理由>`），以後就不會再被提出來。
  - dev-remarks「偏離 tdp」的每一項：先讀這一版 tdp 有沒有放行或給了方案（tdp v0.10.0 起沒有「偏離」，原本的幾家都改成放行或加了
    方案）。有的話照現在的 tdp 引；沒有的話照實寫進 task，review 會判出來，再讓使用者決定修掉或回報 tdp。
- **`work.yaml`**：每個沒取消的 task 一筆，`spec: 1`。
  - `implementation`：從 README 的按鍵與使用方式、dev-remarks 的運作方式、設計文件寫（觸發、經過、讀寫、用到的 terminu-go 元件）。
    文件裡找不到的部分寫「不知道：<缺什麼>」。
  - `test`：`go test -list` 裡名字看得出是在測這個 task 的，寫上去；看不出來的不硬配。沒有測試、但 git history 找得到做這件事的
    commit 的，寫 `commit: <hash>`。都沒有就空著，review 會判沒有測試。

## 4. 搬完刪掉舊文件

tdp v0.10.0 起不管開發者文件（tdp P5），成員的開發文件搬進 harness-assets 後刪掉（user 在 terminu 定的）：

- 刪 `docs/dev-remarks.md` 與 2. 讀過的設計文件。`docs/icon.svg`、README 用的 gif 留著。
- README 兩份裡連到刪掉的檔的連結拿掉（tdp E7：README 第 10 段只連 CHANGELOG）。
- 「建置與開發」「發布」裡沒有地方放、又值得留的（例：發布踩過的坑），寫進相關 task 的 remark；真的沒有對應的 task，就留在 5. 跟
  使用者說。刪掉的內容 git history 都找得到。
- `.local/terminu/progress.md` 是 git 不收的檔，不刪，跟使用者說它不再用了。

## 5. 驗與 commit

1. 跑 MCP 工具 `check`。格式、id、tdp 編號的問題改到沒有；測試相關的問題留著，交給 review。
2. commit（不加 Claude 署名），訊息寫「Move the developer docs into harness-assets」這類。
3. 告訴使用者：
   - harness-assets 產生了，之後 review 會拿它判；要改內容就跟 harness 說。
   - 幾個 task、幾個取消的、寫了「不知道」的有哪幾個、沒有作保的有哪幾個。
   - 刪了哪些檔，沒有地方放的內容。
4. 接著用 Skill 叫 `terminu:review`（從 /terminu 進來接手的，或是 review 叫你的，都回到 review）。
