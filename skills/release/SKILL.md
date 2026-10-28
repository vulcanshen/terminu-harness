---
name: release
description: 發一個 terminu-go app 的新版：先升到最新的 terminu-go、make check、照 tdp review 紀錄、列出還沒做的需求、整理 CHANGELOG，問過使用者才打 tag、push。review 還有不合的不擋，但要使用者確認要發版還是繼續打磨。user 要發版時使用。
argument-hint: [版號，例：v0.3.0、v0.3.0-beta.1]
---

# /terminu:release

發版的規定在 tdp：先用 `tdp_file` 讀 `project/README-zh_TW.md` 的 project.6（CHANGELOG）與 project.7（發布），照 app 照的那一版 tdp
（MCP 工具 `check` 第一行的 tdp 版本）做；下面是 harness 加的檢查與順序。哪一步沒過就停，跟使用者說卡在哪。照 tdp 做的 app
（不是家族成員）的 CHANGELOG 與發布步驟可以自行決定（tdp project.1）：app 改過的，6.–9. 照 app 自己的寫法；沒改過的跟家族一樣。

1. **工作目錄乾淨**：`git status --porcelain` 沒有輸出，而且在預設 branch 上。不乾淨就停，問使用者要先 commit 還是先收起來。
2. **升到最新的 terminu-go，`make check` 過**：go.mod 的版本不是最新的 tag（`gh api repos/vulcanshen/terminu-go/tags --jq '.[0].name'`）
   時，照 newapp 的「升版」（`../newapp/SKILL.md`）升上去，`spec` 改了的 task 照新版做完再往下（user 2026-10-27：app 發版前先升到
   最新；只改文件的版平常不升，在這裡一起升）。接著 `make check`；有 `release-check` target 就也跑（tdp 寫了哪一種警告不算失敗）。
3. **review**：用 Skill 叫 `terminu:review`。它改了 harness-assets 或程式，就 commit（不加 Claude 署名），回到 2。
4. **還有不合的、還沒做的**：
   - review.yaml 還有 `accepted: false` 的：列出來（task、reason），問使用者（AskUserQuestion）：**照樣發版**，或**繼續打磨**
     （用 Skill 叫 `terminu:newapp`，從「6. 繼續」接著做；這次發版停在這裡）。不擋，但要使用者選。
   - `check` 的「待做」還有沒做的 task：一起列給使用者看（spec 有、還沒做的需求），只是告知，不擋。
5. **版號**：使用者給了就用；沒給就照 `## Unreleased` 的改動提一個（有新功能升 minor，只有修正升 patch），問使用者確認。
   格式 `vX.Y.Z`；先發 beta 是 `vX.Y.Z-beta.N`（tdp project.7）：N 從 1 起，同一個 vX.Y.Z 已經發過 beta 就加一。發過 beta 之後
   發正式版，版號就是那個 vX.Y.Z。
6. **CHANGELOG**：照 tdp 把 `## Unreleased` 的標題改成這一版的標題（`## vX.Y.Z — tdp vA.B.C`，beta 是 `## vX.Y.Z-beta.N — tdp vA.B.C`；
   tdp 版本是 `check` 第一行的那一版）。發正式版而上一個正式版之後發過 beta 的：這一段寫上一個正式版以來的全部改動，把那幾個
   beta 段落寫的改動併進來，beta 自己的段落留著（tdp project.6）。`## Unreleased` 沒有內容、也沒有要併的 beta 就停，跟使用者說
   這一版沒有寫下改了什麼。
   **升到新的大版號（v2 以上）**：go.mod 的 module 結尾要是 `/vN`（例 v4 是 `github.com/<帳號>/<app>/v4`；tdp project.3），跟這個
   標題放同一個 commit 一起改三樣：go.mod 的 module、所有 import、README 與文件裡的 `go install` 路徑
   （`…/<app>/vN/cmd/<app>@latest`）。tdptest 拿 CHANGELOG 最新一版的大版號對 module，分開改中間的 `make check` 不會過；已經發過
   v2 以上卻沒帶的，2. 的 `make check` 就會不過，也照這樣改。
   `make check` 過了，commit（不加 Claude 署名）。
7. **問使用者**：給使用者看版號、這一版的 CHANGELOG 段落、review 的結果（合的幾個、使用者選了照樣發的不合項），問要不要打這一版
   的 tag 並 push（AskUserQuestion）。push 之後 `release.yml` 會建置、發 GitHub release、推 Homebrew tap；beta 發成 GitHub 的
   pre-release，不推 tap 與 scoop bucket（tdp project.7）。
8. **同意才做**：`git tag <版號>`，push commit 與 tag。沒有 remote 就打完 tag 停下，跟使用者說。不同意就停在這裡（6. 的 commit
   留著），跟使用者說。
9. **發完**：`gh run watch`（或給使用者 Actions 的網址）看 `release.yml` 跑完；`gh release view <版號>` 對照內文是不是 CHANGELOG
   那一段，beta 的 `--json isPrerelease` 要是 true。失敗就把錯誤給使用者看，不自己刪 tag 或重發。
