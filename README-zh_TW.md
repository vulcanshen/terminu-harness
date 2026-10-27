# terminu-harness

[![License](https://img.shields.io/badge/license-GPL--3.0-blue)](LICENSE)

**語言**：[English](README.md) · 繁體中文

**照 [terminu design](https://github.com/vulcanshen/terminu)（tdp）用 [terminu-go](https://github.com/vulcanshen/terminu-go) 做 TUI app 的 Claude Code plugin：從你想要什麼一路做到完成的 app。你只管功能，畫面與操作照 tdp。**

## 它能做什麼

打 `/terminu`，說你想做什麼：

- **做一個新的 app**：一次問一題，弄清楚 app 要做什麼（只問需求；畫面、按鍵、檔案格式由它照 tdp 決定），把你的話記下來、每個需求對到 tdp，就開始做：照 tdp 的 project 範本用 terminu-go 一個需求一個需求做出來，每一個都跑測試、在它自己帶的終端機裡打開來看，記下怎麼做的、哪些測試守著，做完給你看。
- **tdp 找不到答案時不自己決定**：app 的每個決定都要在 tdp 找得到根據，或在 tdp 明文寫「app 可以自行決定」的範圍裡。tdp 沒寫、寫得不清楚、或 terminu-go 還做不到的地方，它會跟你說明，替 tdp 或 terminu-go 寫好 issue 草稿，你同意了才發；等的期間先用暫時的做法，它們更新後再回來改。
- **繼續做、加需求、跟上新版**：在 app 的資料夾打 `/terminu`，它先查 terminu-go 有沒有新版（每一版 tdp 都有一版 terminu-go 對到它）、在等的 issue 有沒有回應、程式裡還剩哪些暫時的做法，講給你聽；要升版、要加需求都從這裡接著做。
- **接手現成的 terminu app**：先用 tdp 與 terminu-go 做、之後才掛上 plugin 的 app，從它的 README（做什麼）與 git history（每一塊怎麼做的、哪些測試守著）整理出紀錄，之後照同一套流程走。
- **review**：另一個只讀的 reviewer 拿 tdp 對 app 的紀錄判每一個需求的做法，不讀程式碼；要你自己確認的（例：要一台 Windows 才驗得了）每次都請你重新確認。不合的它自己改紀錄或程式修到合，只有要不要回報 tdp 會問你。
- **發版**：跑測試、review、告訴你還有哪些不合與還沒做的需求、照 tdp 整理 CHANGELOG，問過你才打 tag、push。
- **回報 tdp 或 terminu-go**：替 tdp 或 terminu-go 寫 issue 草稿，你同意了才發。

不需要 tmux 或圖形介面；會離開你這台機器的事（issue、GitHub repo、發版）都先問你。

## 需求

- [Claude Code](https://claude.com/claude-code)
- macOS 或 Linux
- [Go](https://go.dev/dl/)：plugin 第一次啟動時從原始碼建出它的終端機與檢查工具；app 也用 Go 寫
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

它會問你要做什麼；目前資料夾有它在管的 app 時，直接接著做；是還沒管過的 terminu app 時，問你要不要接手。也可以直接說：`/terminu 做一個……的 app`。每個步驟也各有自己的指令：

| 指令 | 做什麼 |
|---|---|
| `/terminu:newapp [app 要做什麼]` | 做新的 app，或繼續做、加需求、跟上新版 |
| `/terminu:analysis` | 接手現成的 terminu app：從它的 README 與 git history 整理出紀錄 |
| `/terminu:review [task 編號]` | 照 tdp 判 app 的紀錄；沒給編號就判全部 |
| `/terminu:release [版號]` | 發版：測試、review、整理 CHANGELOG，問過你才打 tag、push |
| `/terminu:reportgap [缺了什麼]` | 替 tdp 或 terminu-go 寫 issue 草稿，你同意後發出去 |

在它管的 app 裡 commit 時，改了程式卻沒更新紀錄（`harness-assets/work.yaml`），commit 會被擋下來，並告訴 Claude 要補紀錄；只是重構、行為沒變時，commit 訊息最後加一行 `Terminu-Work: unchanged`。

## 存放位置

| 什麼 | 在哪 |
|---|---|
| app 的紀錄：你的需求（`idea.md`）、每個需求對到的 tdp（`spec.yaml`）、怎麼做的與哪些測試守著（`work.yaml`）、最近一次 review（`review.yaml`） | app repo 的 `harness-assets/` |
| 給在這個 repo 寫程式的 Claude 的規矩 | app repo 的 `.claude/rules/terminu.md` |
| 抓下來的 tdp 檔案 | `~/.terminu-harness/cache/tdp/` |
| 建好的終端機與檢查工具 | `~/.terminu-harness/bin/` |

| 環境變數 | 作用 |
|---|---|
| `TERMINU_HARNESS__HOME` | 全部改放到 `~/.terminu-harness` 以外的地方 |
| `TERMINU_HARNESS__TDP_REPO` | 從本機的 terminu clone 讀 tdp，不從 GitHub 抓 |

## 限制

- **review 只判紀錄裡寫了的**：程式做了、`harness-assets/work.yaml` 沒寫的，review 看不到，harness 不保證。
- 你在終端機自己打的 commit 不會被擋，下一次 review 才看得到。
- plugin 本身不支援 Windows（做出來的 app 可以出 Windows 版）。
- 等 issue 不是一直開著的 session：tdp 或 terminu-go 回應之後，回到 app 的資料夾打 `/terminu`。

## 相關連結

- [terminu](https://github.com/vulcanshen/terminu)：terminu design（tdp）與它的清單
- [terminu-go](https://github.com/vulcanshen/terminu-go)：terminu 家族 app 用的 Go lib
- [CHANGELOG](CHANGELOG.md)
- [docs/dev-remarks.md](docs/dev-remarks.md)：運作方式，給開發者看

## License

[GPL-3.0](LICENSE)
