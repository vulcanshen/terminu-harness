---
name: review
description: 叫 terminu 的 reviewer 看一個 terminu-go app 的一段 diff：照 tdp 與 app 的 spec 找用法上的問題，每個問題附出處；有問題就修，再看一次。user 要 review 目前的改動時使用；/terminu:release 也照這裡做。
argument-hint: [範圍，例：main..HEAD、v0.2.0..HEAD]
---

# /terminu:review

## 1. 範圍

- `$ARGUMENTS` 給了範圍就用它：`git diff <範圍>`。
- 沒給：還沒 commit 的改動（`git diff HEAD`，加上 `git status --porcelain` 裡沒追蹤的新檔），加上這個 branch 還沒 push 的 commit
  （`git diff @{upstream}...HEAD`；沒有 upstream 就從最新的 tag 起算，沒有 tag 就從第一個 commit）。
- 範圍是空的：說「沒有要看的改動」，結束。

## 2. 派 reviewer

用 Agent 工具派一個 `terminu:reviewer`，prompt 寫：

- app 的 repo（目前目錄）
- 範圍，與讀 diff 的指令（例：`git diff v0.2.0..HEAD`；有沒追蹤的新檔就列出來）
- 第二次以後：上一次的問題原文，與這次修正的範圍

## 3. 看結果、修

1. 把 reviewer 的問題、建議、tdp 不清楚的照原樣給使用者看。
2. 問題一條一條修（app 自己的錯）。修完再派一次 reviewer，只看修正的那一段（2. 的「第二次以後」）。**最多來回兩次。**
3. 兩次之後還在的，或你不同意 reviewer 的，一條一條問使用者（AskUserQuestion，一次一題）：
   - **修掉**
   - **寫成偏離**：寫進 `docs/dev-remarks.md` 的「偏離 tdp」，附理由（tdp departure-1）
   - **回報 tdp**：用 Skill 叫 `terminu:reportgap`；程式標 `TODO(tdp #N)`，列進 `docs/<app>-terminu-fix.md`
4. tdp 不清楚的：問使用者要不要交給 `terminu:reportgap`。

## 4. 結束

說還剩幾個問題（0 才算過）、建議有幾個、改了哪些檔。修了東西不自己 commit，除非使用者要，或是 `/terminu:release` 叫你的。
