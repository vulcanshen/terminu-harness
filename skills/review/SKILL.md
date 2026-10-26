---
name: review
description: 照 tdp 判一個 terminu-go app 的 harness-assets：先驗紀錄、請使用者重新確認要人確認的作保，再派 terminu 的 reviewer 拿 tdp 對 work.yaml 的每一個 task，結果寫進 review.yaml，不合的一項一項跟使用者決定怎麼處理。不讀程式碼。user 要 review 時使用；/terminu:release 也照這裡做。
argument-hint: [只看哪幾個 task，例：3 5 8]
---

# /terminu:review

review 拿 tdp 對 app 的紀錄（`harness-assets/`），不讀程式碼：程式對不對是測試的事，review 判的是「紀錄寫的做法合不合 tdp」。
work.yaml 沒寫的，review 看不到。格式見 newapp 的 `harness-assets.md`（這個 skill 資料夾旁邊的 `../newapp/harness-assets.md`）。

## 1. 驗紀錄

1. 沒有 `harness-assets/spec.yaml`（現成的 terminu app，還沒用 harness 管過）：用 Skill 叫 `terminu:genwork` 產生，回來再往下。
2. 跑 MCP 工具 `check`。
   - `go.mod` 讀不到 terminu-go：這不是 harness 管的 app，跟使用者說，結束。
   - 測試、commit、使用者確認的問題（「測試 X 不在 go test -list 裡」「沒有測試」……）：記下是哪些 task，3. 直接判 `沒有測試`。
   - 其他問題（格式、id、spec 對不上、tdp 編號找不到）：紀錄本身錯了，先改 harness-assets，再跑 `check`，到只剩上一條那種問題為止。
3. 記下 `check` 第一行的 tdp 版本：這次照這一版判。

## 2. 請使用者重新確認

work.yaml 裡沒取消的 task，作保是 `user: {date, how}` 的，每一次 review 都請使用者重新確認（AI 不替使用者確認）：

- 用 AskUserQuestion 一次問完：每一題 `multiSelect: true`，選項是「task <id>：<how>（上次 <date>）」，一題最多 4 個，一次最多 4 題，
  多的分幾次問。問句寫「這幾項你剛剛照著確認過了嗎？確認過的勾起來」。
- 勾了的：把 date 改成今天。沒勾的：那個 task 判 `沒有測試：使用者沒有重新確認（<how>）`。
- 使用者沒照著做就勾，是使用者自己負責。

## 3. 派 reviewer

用 Agent 工具派一個 `terminu:reviewer`，prompt 寫：

- app 的 repo（目前目錄）、tdp 版本（1. 記的）
- 要判的 task：work.yaml 裡沒取消的全部；`$ARGUMENTS` 指定了就只判那幾個。1.、2. 已經判成 `沒有測試` 的不用給它。
- 第二次以後：上一次判不合的 task 與原因，與這次改了哪些 task

## 4. 寫 review.yaml

把 reviewer 的結果，加上 1.、2. 判的 `沒有測試`，寫成 `harness-assets/review.yaml`，覆蓋舊的：`spec` 是 spec.yaml 的 id，`tasks`
列 work.yaml 裡每一個沒取消的 task。`$ARGUMENTS` 只指定了幾個時，其他 task 照舊的 review.yaml 留著。

reviewer 說 spec.yaml 少引了哪些 tdp 的，補進那個 task 的 `tdp`（紀錄的事，不算判定）。

跑 `check`，review.yaml 沒有問題。

## 5. 不合的，跟使用者決定

把判的結果給使用者看：合的幾個、不合的照三種分開列（task、reason）。不合的一項一項問（AskUserQuestion，一次一題）：

| reason 開頭 | 選項 |
|---|---|
| 違反 | **修掉**（改程式、測試與 work.yaml）；**回報 tdp**（覺得 tdp 的規定不合用，用 Skill 叫 `terminu:reportgap`，remark 記 issue）；**先不處理** |
| 沒有測試 | **補測試**；**改成使用者確認**（問使用者怎麼確認，記進 `user`）；**先不處理** |
| tdp 不清楚 | **回報 tdp**（`terminu:reportgap`，remark 記 issue）；**先不處理** |

- 選「先不處理」的留在 review.yaml，下一次 review 還是會列出來。
- 修了東西的：改完跑 `check`，再回到 3. 只判改了的 task，更新 review.yaml。

## 6. 結束

說合的、不合的各幾個，不合的裡面使用者選了什麼，改了哪些檔。修了東西不自己 commit，除非使用者要，或是 `/terminu:release` 叫你的。
