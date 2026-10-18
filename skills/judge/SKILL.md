---
name: judge
description: 照 terminu design（tdp）的條目清單判一個 TUI app 有多少 terminu gene，寫成績單。user 打 /terminu:judge、要判某個 app 照 tdp 做得怎樣、要數 gene 時使用。
argument-hint: <app 的指令> [app 的 repo] [tdp 版本]
---

# /terminu:judge

照 tdp 的清單把一個 TUI app 全部判一次，寫一份成績單：每一條的判定與證據、`tdp-<版本> <數> gene`、tdp 沒規範的東西。
judge 從外面操作 app（harness 自己帶的終端機），任何 TUI app 都能判；拿得到 repo 時，要讀程式、讀檔案的條目也一起判。

## 1. 弄清楚要判什麼

從 `$ARGUMENTS` 拿這三樣，缺的問 user（一次問一題）：

- **app 怎麼啟動**：指令，例 `demou`、`./demou`。必要。
- **app 的 repo**：本機路徑，沒有就不給。現在的工作目錄看起來就是那個 app 的 repo 時，問 user 要不要用它。
- **照哪一版 tdp**：不給就用最新的 tag。
- **怎麼啟動、操作 app 的補充**（有的話）：例如要先設環境變數、要用包裝指令、不能碰真的系統檔。下面每個 judge 的 prompt 都原樣
  附上這段。

## 2. 開一個 run

呼叫 `judge_start`（`app_command`、`repo`、`tdp_version`）。它回傳 run 目錄與各批條目的檔案。告訴 user：照哪一版 tdp、幾條、
分幾批，成績單會寫在哪個目錄。

## 3. 認識 app

用 Agent 工具派一個 `terminu:judge`，prompt 寫明：

- 這次做「認識 app」
- run 目錄、app 的指令、repo（或「沒有 repo」）、在哪個目錄啟動 app（目前的工作目錄）、使用者的補充

等它回來。它把 profile 存在 run 目錄的 `profile.md`。

## 4. 判每一批

在**同一則訊息**裡，每一批派一個 `terminu:judge`（讓它們同時跑），prompt 寫明：

- 這次做「判一批條目」，批次編號與那一批的檔案路徑
- run 目錄、app 的指令、repo（或「沒有 repo」）、在哪個目錄啟動 app、使用者的補充

全部回來之後呼叫 `judge_status`。有沒判到的條目，就把那幾條的編號交給一個 `terminu:judge` 再判一次（只補一次）。

## 5. 寫成績單

呼叫 `judge_finish`，把結果告訴 user：

- gene 數（`tdp-<版本> <數> gene`）與各種判定的條數
- 違反的條目：列出編號與一句話原因（多的話列前十條，其他看成績單）
- 判不了的原因分類（拿不到 repo、清單不清楚、沒判到……）
- judge 回報「清單不清楚」的條目：建議回報給 tdp
- tdp 沒規範的東西有幾項
- 成績單的路徑（`report.md`）

不要改 app 的 repo，也不要替 app 開 issue；要不要回報 tdp 由 user 決定。
