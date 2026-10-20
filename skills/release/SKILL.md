---
name: release
description: 發一個 terminu-go app 的新版：make check、reviewer 看上一個 tag 到現在的改動、整理 CHANGELOG，問過使用者才打 tag、push。user 要發版時使用。
argument-hint: [版號，例：v0.3.0]
---

# /terminu:release

發版的規定在 tdp：先用 `tdp_file` 讀 `project/README-zh_TW.md` 的「CHANGELOG」與「發布」兩節，照那一版做；下面是 harness 加的檢查
與順序。哪一步沒過就停，跟使用者說卡在哪。

1. **工作目錄乾淨**：`git status --porcelain` 沒有輸出，而且在預設 branch 上。不乾淨就停，問使用者要先 commit 還是先收起來。
2. **`make check` 過**；有 `release-check` target 就也跑（tdp 寫了哪一種警告不算失敗）。
3. **review**：範圍是上一個 tag 到 HEAD（`git describe --tags --abbrev=0`；沒有 tag 就是整個 repo）。照 `/terminu:review` 的
   2、3 做：用 Skill 叫 `terminu:review`，參數給範圍（例：`v0.2.0..HEAD`）。修了東西就 commit（不加 Claude 署名），回到 2。
   **還有問題就不往下走。**
4. **版號**：使用者給了就用；沒給就照 `## Unreleased` 的改動提一個（有新功能升 minor，只有修正升 patch），問使用者確認。
   格式 `vX.Y.Z`。
5. **CHANGELOG**：照 tdp 把 `## Unreleased` 的標題改成這一版的標題（`## vX.Y.Z — tdp vA.B.C`，tdp 版本是 app 現在跟的那一版）。
   `## Unreleased` 沒有內容就停，跟使用者說這一版沒有寫下改了什麼。commit（不加 Claude 署名）。
6. **問使用者**：給使用者看版號、這一版的 CHANGELOG 段落、review 的結果，問要不要打 tag `vX.Y.Z` 並 push（AskUserQuestion）。
   push 之後 `release.yml` 會建置、發 GitHub release、推 Homebrew tap。
7. **同意才做**：`git tag vX.Y.Z`，push commit 與 tag。沒有 remote 就打完 tag 停下，跟使用者說。不同意就停在這裡（5. 的 commit
   留著），跟使用者說。
8. **發完**：`gh run watch`（或給使用者 Actions 的網址）看 `release.yml` 跑完；`gh release view vX.Y.Z` 對照內文是不是 CHANGELOG
   那一段。失敗就把錯誤給使用者看，不自己刪 tag 或重發。
