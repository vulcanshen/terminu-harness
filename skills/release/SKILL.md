---
name: release
description: 發一個 terminu-go app 的新版：先升到最新的正式版 terminu-go（app 已經在 beta 上才算 beta）、make check、照 tdp review 紀錄、列出還沒做的需求、整理 CHANGELOG，問過使用者才打 tag、push。review 還有不合的不擋，但要使用者確認要發版還是繼續打磨。user 要發版時使用。
argument-hint: [版號，例：v0.3.0、v0.3.0-beta.1]
---

# /terminu:release

發版的規定在家族的 repo 範本（terminu project）的 project.6（CHANGELOG）與 project.7（發布）。2. 升完 terminu-go 之後看 app 照的
tdp（`check` 第一行）：還有 project/（v1.1.10 以前）就用 `tdp_section` 讀 `tdp project.6`、`tdp project.7`；沒有（v1.2.0-beta.1
起，tdp E6 交給 harness）就讀 `harness project.6`、`harness project.7`。下面寫 `harness project.N` 的地方，照前一種的換成
`tdp project.N`。下面是 harness 加的檢查與順序。哪一步沒過就停，跟使用者說卡在哪。照 tdp 做的 app（不是家族成員）的 CHANGELOG 與發布步驟可以自行決定
（harness project.1、tdp 的 E 章開頭）：app 改過的，6.–9. 照 app 自己的寫法；沒改過的跟家族一樣。

1. **工作目錄乾淨**：`git status --porcelain` 沒有輸出，而且在預設 branch 上。不乾淨就停，問使用者要先 commit 還是先收起來。
2. **升到最新的 terminu-go，`make check` 與 `check` 過**：go.mod 的版本比 newapp「0. 準備」的「要用的 terminu-go」舊
   （`../newapp/SKILL.md`；預設是最新的正式版，app 已經在 beta 上才算 beta，使用者指定的照指定的）時，照 newapp 的「升版」升上去，
   `spec` 改了的 task 照新版做完再往下（user 2026-10-27：app 發版前先升到最新；只改文件的版平常不升，在這裡一起升）。正式版的 app
   不會在這裡被升到 beta。接著 `make check`；有 `release-check` target 就也跑（harness project.7 寫了哪一種警告不算失敗）。再跑
   MCP 工具 `check`，報的問題改到沒有：app 照的 tdp 沒有 project/ 時，家族成員的 repo（範本、目錄、Makefile、CHANGELOG）由它拿
   harness 的範本比；還有 project/ 時，由那一版 terminu-go 的 tdptest 在 `make check` 裡查。
3. **review**：用 Skill 叫 `terminu:review`。它改了 harness-assets 或程式，就 commit（不加 Claude 署名），回到 2。
4. **還有不合的、還沒做的**：
   - review.yaml 還有 `accepted: false` 的：列出來（task、reason），問使用者（AskUserQuestion）：**照樣發版**，或**繼續打磨**
     （用 Skill 叫 `terminu:newapp`，從「6. 繼續」接著做；這次發版停在這裡）。不擋，但要使用者選。
   - `check` 的「待做」還有沒做的 task：一起列給使用者看（spec 有、還沒做的需求），只是告知，不擋。
5. **版號**：使用者給了就用；沒給就照 `## Unreleased` 的改動提一個（有新功能升 minor，只有修正升 patch），問使用者確認。
   格式 `vX.Y.Z`；先發 beta 是 `vX.Y.Z-beta.N`（harness project.7）：N 從 1 起，同一個 vX.Y.Z 已經發過 beta 就加一。發過 beta 之後
   發正式版，版號就是那個 vX.Y.Z。
6. **CHANGELOG**：照 harness project.6 把 `## Unreleased` 的標題改成這一版的標題（`## vX.Y.Z — tdp vA.B.C`，beta 是 `## vX.Y.Z-beta.N — tdp vA.B.C`；
   tdp 版本是 `check` 第一行的那一版）。發正式版而上一個正式版之後發過 beta 的：這一段寫上一個正式版以來的全部改動，把那幾個
   beta 段落寫的改動併進來，beta 自己的段落留著（harness project.6）。`## Unreleased` 沒有內容、也沒有要併的 beta 就停，跟使用者說
   這一版沒有寫下改了什麼。
   **升到新的大版號（v2 以上）**：go.mod 的 module 結尾要是 `/vN`（例 v4 是 `github.com/<帳號>/<app>/v4`；harness project.3），跟這個
   標題放同一個 commit 一起改三樣：go.mod 的 module、所有 import、README 與文件裡的 `go install` 路徑
   （`…/<app>/vN/cmd/<app>@latest`）。`check` 拿 CHANGELOG 最新一版的大版號對 module，分開改中間的 `check` 不會過；已經發過
   v2 以上卻沒帶的，2. 的 `check` 就會指出來，也照這樣改。
   `make check` 與 `check` 過了，commit（不加 Claude 署名）。
7. **問使用者**：給使用者看版號、這一版的 CHANGELOG 段落、review 的結果（合的幾個、使用者選了照樣發的不合項），問要不要打這一版
   的 tag 並 push（AskUserQuestion）。push 之後 `release.yml` 會建置、發 GitHub release、推 Homebrew tap；beta 發成 GitHub 的
   pre-release，不推 tap 與 scoop bucket（harness project.7）。
8. **同意才做**：`git tag <版號>`，push commit 與 tag。沒有 remote 就打完 tag 停下，跟使用者說。不同意就停在這裡（6. 的 commit
   留著），跟使用者說。
9. **發完**：`gh run watch`（或給使用者 Actions 的網址）看 `release.yml` 跑完；`gh release view <版號>` 對照內文是不是 CHANGELOG
   那一段，beta 的 `--json isPrerelease` 要是 true。失敗就把錯誤給使用者看，不自己刪 tag 或重發。
