# terminu-harness

[![License](https://img.shields.io/badge/license-GPL--3.0-blue)](LICENSE)

**語言**：[English](README.md) · 繁體中文

**照 [terminu design](https://github.com/vulcanshen/terminu)（tdp）做 TUI app 的 Claude Code plugin：從你想要什麼一路做到完成的 app，也檢查一個 app 有多少 terminu 的 gene。** tdp 清單上 app 符合的每一條，就是一個 gene。

## 它能做什麼

打 `/terminu`，說你想做什麼：

- **做一個新的 app**：一次問一題，弄清楚 app 要做什麼，寫一份規格給你確認，再照 tdp 的 project 範本用 [terminu-go](https://github.com/vulcanshen/terminu-go) 一塊一塊做出來，每做完一塊就跑給你看。遇到 tdp 沒規範、或 terminu-go 還做不到的地方，它會跟你說明，替 tdp 或 terminu-go 寫好 issue 草稿，你同意了才發；等的期間先用暫時的做法，它們更新後再回來改。做完自動檢查，給你成績單。
- **繼續做到一半的 app**：從上次停的地方接著做，看在等的 issue 有沒有回應，tdp 或 terminu-go 出新版就升級。
- **判一個 app**：在自己帶的終端機裡打開 app，照清單每一條的驗法按鍵、看畫面，有 repo 就讀程式，寫一份成績單：app **通過**（沒有違反）還是**不通過**（有任何違反，不管是 app 還是 terminu-go 要修；app 寫明理由的偏離不算）；每一條的判定（符合、偏離、違反、不適用、判不了）與看到的證據；寫成 `tdp-0.8.1 860 gene` 的 gene 數，看 app 的進度；app 做了但 tdp 還沒規範的東西。任何 TUI app 都能判，不一定是 terminu 家族的成員，拿不拿得到程式碼都行。

不需要 tmux 或圖形介面；會離開你這台機器的事（issue、GitHub repo、發版）都先問你；判 app 時不動 app 的 repo。

## 需求

- [Claude Code](https://claude.com/claude-code)
- macOS 或 Linux
- [Go](https://go.dev/dl/)：plugin 第一次啟動時從原始碼建出它的終端機；新的 app 也用 Go 寫
- [GitHub CLI](https://cli.github.com)（`gh`），登入好，才能發 issue；沒有的話它把 issue 內容給你自己貼

## 安裝

在 Claude Code 裡：

```
/plugin marketplace add vulcanshen/terminu-harness
/plugin install terminu@terminu-harness
```

第一次啟動會把終端機建到 `~/.terminu-harness/bin`，要等一下。Claude Code 說 plugin 的 MCP server 沒啟動成功的話，等建好之後跑 `/mcp` 重新連上。

移除：`/plugin uninstall terminu@terminu-harness`，再刪掉 `~/.terminu-harness`。

## 使用方式

```
/terminu
```

它會問你要做什麼。也可以直接說：`/terminu 做一個……的 app`。每個步驟也各有自己的指令：

| 指令 | 做什麼 |
|---|---|
| `/terminu:new-app [app 要做什麼]` | 做新的 app，或繼續做這個資料夾裡的 app |
| `/terminu:judge <app 怎麼啟動> [app 的 repo] [tdp 版本]` | 判一個 app；沒給 repo，要讀程式的條目判「判不了」；沒給版本，照最新的 tdp |
| `/terminu:report-gap [缺了什麼]` | 替 tdp 或 terminu-go 寫 issue 草稿，你同意後發出去 |

判一次會把清單分給十個 judge 同時判，要一段時間，也用掉不少 token。

### 成績單怎麼看

tdp 只寫要查什麼，怎麼判由 harness 決定：

- **結果**：有任何一條違反就**不通過**，不管是 app 自己要修、還是在等 tdp 或 terminu-go；沒有違反就**通過**。違反分兩組列：
  app 要修的，以及列在 app 的 `docs/<app>-terminu-fix.md`、在等 tdp 或 terminu-go 的。
- **每一條的判定**：

| 判定 | 意思 |
|---|---|
| 符合 | 做到這一條 |
| 偏離 | 沒做到，但 app 的 `docs/dev-remarks.md`「偏離 tdp」寫了這一條與原因；不擋，列在成績單最前面 |
| 違反 | 沒做到，也沒寫成偏離；每個 app 都要有的東西整個沒有時，它的每一條細節也算違反 |
| 不適用 | app 沒有這一條講的功能，tdp 也沒要求要有；理由寫出怎麼知道沒有。「都要」的條目不會不適用 |
| 判不了 | 沒辦法確認，寫出原因：tdp 沒寫清楚或 terminu-go 擋住了（回報給 terminu），或其他原因（例：拿不到 repo，要你回應） |

- **gene**：判成符合的條數，寫成 `tdp-<版本> <數> gene`，綁在判的那一版 tdp。用來看進度；不同 app 的功能多寡不同，數字不能互比。
- 從畫面量時間時，跟條目寫的差在兩格畫面（約 34ms）內算符合。
- 同一個 app 用同一版 tdp 與 terminu-go 判過之後，再判只重判上次之後改到的部分，其他沿用上一份成績單。

## 存放位置

| 什麼 | 在哪 |
|---|---|
| 成績單 | `~/.terminu-harness/judge/<app>/<時間>-tdp-<版本>/`：`report.md`、每一條的判定、證據 |
| 新 app 的進度與 issue 草稿 | app repo 的 `.local/terminu/`（git 不收） |
| 新 app 的規格 | app repo 的 `docs/spec.md` |
| 抓下來的 tdp 檔案 | `~/.terminu-harness/cache/tdp/` |
| 建好的終端機 | `~/.terminu-harness/bin/` |

| 環境變數 | 作用 |
|---|---|
| `TERMINU_HARNESS__HOME` | 全部改放到 `~/.terminu-harness` 以外的地方 |
| `TERMINU_HARNESS__OUT` | 成績單改放到別處 |
| `TERMINU_HARNESS__TDP_REPO` | 從本機的 terminu clone 讀 tdp，不從 GitHub 抓 |

## 限制

- plugin 本身不支援 Windows（做出來的 app 可以出 Windows 版）。
- 等 issue 不是一直開著的 session：tdp 或 terminu-go 回應之後，回來打 `/terminu` 選「繼續」。
- 終端機不送滑鼠事件：要用滑鼠驗的條目判「判不了」。
- 從畫面量的時間（動畫、toast 停留多久）先以 33ms 為容許誤差，等 tdp 定下來。

## 相關連結

- [terminu](https://github.com/vulcanshen/terminu)：terminu design（tdp）與它的清單
- [terminu-go](https://github.com/vulcanshen/terminu-go)：terminu 家族 app 用的 Go lib
- [CHANGELOG](CHANGELOG.md)
- [docs/dev-remarks.md](docs/dev-remarks.md)：運作方式，給開發者看

## License

[GPL-3.0](LICENSE)
