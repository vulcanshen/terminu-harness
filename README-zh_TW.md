# terminu-harness

[![License](https://img.shields.io/badge/license-GPL--3.0-blue)](LICENSE)

**語言**：[English](README.md) · 繁體中文

**照 [terminu design](https://github.com/vulcanshen/terminu)（tdp）做 TUI app 的 Claude Code plugin：從你想要什麼一路做到完成的 app，也檢查一個 app 有多少 terminu 的 gene。** tdp 清單上 app 符合的每一條，就是一個 gene。

## 它能做什麼

打 `/terminu`，說你想做什麼：

- **做一個新的 app**：一次問一題，弄清楚 app 要做什麼，寫一份規格給你確認，再照 tdp 的 project 範本用 [terminu-go](https://github.com/vulcanshen/terminu-go) 一塊一塊做出來，每做完一塊就跑給你看。遇到 tdp 沒規範、或 terminu-go 還做不到的地方，它會跟你說明，替 tdp 或 terminu-go 寫好 issue 草稿，你同意了才發；等的期間先用暫時的做法，它們更新後再回來改。做完自動檢查，給你成績單。
- **繼續做到一半的 app**：從上次停的地方接著做，看在等的 issue 有沒有回應，tdp 或 terminu-go 出新版就升級。
- **判一個 app**：在自己帶的終端機裡打開 app，照清單每一條的驗法按鍵、看畫面，有 repo 就讀程式，寫一份成績單：app **通過**（沒有違反）、**待修**（違反都列在它的 `docs/<app>-terminu-fix.md`，在等 tdp 或 terminu-go）還是**不通過**（有違反沒列進去）；每一條的判定（符合、偏離、違反、不適用、判不了）與看到的證據；寫成 `tdp-0.7.4 860 gene` 的 gene 數，看 app 的進度；app 做了但 tdp 還沒規範的東西。任何 TUI app 都能判，不一定是 terminu 家族的成員，拿不拿得到程式碼都行。

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
