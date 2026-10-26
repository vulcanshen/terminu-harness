---
name: reportgap
description: 把 tdp 沒寫、寫不清楚、兩處矛盾，或 terminu-go 做不到的地方寫成 issue，給使用者看過、同意後發到 GitHub（vulcanshen/terminu 或 vulcanshen/terminu-go），記進 spec.yaml 那個 task 的 remark。newapp 開發卡住、review 判出 tdp 不清楚、使用者要回報 tdp 或 terminu-go 時使用。
argument-hint: [卡在哪]
---

# 回報給 tdp 或 terminu-go

tdp 是宇集：app 的每個決定都要在 tdp 找得到根據，或在 tdp 明文寫「app 可以自行決定」的範圍裡。找不到的，harness 不自己補規定：
規範的事回報 tdp，lib 的事回報 terminu-go，等它們更新。

## 1. 發到哪

| 卡在 | 發到 |
|---|---|
| tdp 沒寫、寫不清楚、兩處互相矛盾、照字面做會出錯、規定不合這個 app 用 | `vulcanshen/terminu` |
| terminu-go 做不到 tdp 寫的、有 bug、缺 API | `vulcanshen/terminu-go` |
| app 自己的錯 | 不發，自己修 |

## 2. 先查有沒有人回報過

`gh issue list -R <repo> --state all --search "<關鍵字>"`。已經有就不重發，把那張的網址記下來（必要時在那張留言補上這次的情況，
留言前一樣先給使用者看）。

## 3. 寫草稿

用繁體中文（tdp 的原文是繁中）。標題一句話講卡在哪。版本：terminu-go 看 app 的 go.mod，tdp 是它對到的那一版（MCP 工具 `check`
的第一行）。內文：

```
## 卡在哪
<app 名>、在做什麼（spec.yaml 的 task 與 request）、遇到什麼。tdp <版本>，terminu-go <版本>。

## tdp 怎麼寫
引用編號（例 tdp K4、tdp panel.5）並貼原文；tdp 沒寫的，寫找過哪幾節。

## 可能的做法
A …（後果）
B …（後果）

## app 暫時怎麼做
程式裡的標記（TODO(tdp #?) 或 TODO(tgo #?)，發出去之後補上編號）與暫時的做法；沒有合 tdp 的暫時做法時，寫這個 task 先不做。

---
由 terminu-harness 回報
```

## 4. 給使用者看，同意才發

把發到哪、標題、內文給使用者看，問要不要發。改到使用者同意為止；不同意就不發，記下「沒發」與原因。

同意後：`gh issue create -R <repo> --title "<標題>" --body-file <檔案>`。沒有 `gh` 或沒登入時，把標題與內文給使用者自己貼。

使用者指定別的送法時照使用者的（例：先不發 GitHub，在本機用 SendMessage 送給負責 tdp 的 Claude session）。

## 5. 記下來

在 `harness-assets/spec.yaml` 那個 task 的 remark 加一筆 `{note: 卡在哪、暫時怎麼做, issue: tdp #N}`（tgo 的寫 `tgo #N`；GitHub
以外的送法寫送去哪、哪一天）。程式裡的標記把 `#?` 換成編號。跑 `check`。
