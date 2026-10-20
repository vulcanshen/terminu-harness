# terminu-harness

[![License](https://img.shields.io/badge/license-GPL--3.0-blue)](LICENSE)

**語言**：[English](README.md) · 繁體中文

**照 [terminu design](https://github.com/vulcanshen/terminu)（tdp）用 [terminu-go](https://github.com/vulcanshen/terminu-go) 做 TUI app 的 Claude Code plugin：從你想要什麼一路做到完成的 app。**

## 它能做什麼

打 `/terminu`，說你想做什麼：

- **做一個新的 app**：一次問一題，弄清楚 app 要做什麼（只問需求；畫面、按鍵、檔案格式由它照 tdp 決定），寫成規格就開始做：照 tdp 的 project 範本用 terminu-go 一塊一塊做出來，每一塊都跑測試、在它自己帶的終端機裡打開來看，做完給你看。遇到 tdp 沒規範、或 terminu-go 還做不到的地方，它會跟你說明，替 tdp 或 terminu-go 寫好 issue 草稿，你同意了才發；等的期間先用暫時的做法，它們更新後再回來改。
- **繼續做到一半的 app**：在 app 的資料夾打 `/terminu`，它先查在等的 issue 有沒有回應、tdp 或 terminu-go 有沒有新版、程式裡還剩哪些暫時的做法，講給你聽，再從上次停的地方接著做。
- **review 改動**：另一個只讀的 reviewer 照 tdp 與 app 的規格看一段改動，每個問題都附出處（tdp 的條目或規格的哪一句），有問題就修了再看一次。
- **發版**：跑測試、review 上一版之後的改動、照 tdp 整理 CHANGELOG，問過你才打 tag、push。
- **回報 tdp 或 terminu-go**：替 tdp 或 terminu-go 寫 issue 草稿，你同意了才發。

不需要 tmux 或圖形介面；會離開你這台機器的事（issue、GitHub repo、發版）都先問你。

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

它會問你要做什麼；目前資料夾有做到一半的 app 時，直接接著做。也可以直接說：`/terminu 做一個……的 app`。每個步驟也各有自己的指令：

| 指令 | 做什麼 |
|---|---|
| `/terminu:newapp [app 要做什麼]` | 做新的 app，或繼續做這個資料夾裡的 app |
| `/terminu:review [範圍]` | review 一段改動；沒給範圍就看還沒 commit 的改動與還沒 push 的 commit |
| `/terminu:release [版號]` | 發版：測試、review 上一版之後的改動、整理 CHANGELOG，問過你才打 tag、push |
| `/terminu:reportgap [缺了什麼]` | 替 tdp 或 terminu-go 寫 issue 草稿，你同意後發出去 |

## 存放位置

| 什麼 | 在哪 |
|---|---|
| 新 app 的進度與 issue 草稿 | app repo 的 `.local/terminu/`（git 不收） |
| 新 app 的規格 | app repo 的 `docs/spec.md` |
| 抓下來的 tdp 檔案 | `~/.terminu-harness/cache/tdp/` |
| 建好的終端機 | `~/.terminu-harness/bin/` |

| 環境變數 | 作用 |
|---|---|
| `TERMINU_HARNESS__HOME` | 全部改放到 `~/.terminu-harness` 以外的地方 |
| `TERMINU_HARNESS__TDP_REPO` | 從本機的 terminu clone 讀 tdp，不從 GitHub 抓 |

## 限制

- plugin 本身不支援 Windows（做出來的 app 可以出 Windows 版）。
- 等 issue 不是一直開著的 session：tdp 或 terminu-go 回應之後，回到 app 的資料夾打 `/terminu`。

## 相關連結

- [terminu](https://github.com/vulcanshen/terminu)：terminu design（tdp）與它的清單
- [terminu-go](https://github.com/vulcanshen/terminu-go)：terminu 家族 app 用的 Go lib
- [CHANGELOG](CHANGELOG.md)
- [docs/dev-remarks.md](docs/dev-remarks.md)：運作方式，給開發者看

## License

[GPL-3.0](LICENSE)
