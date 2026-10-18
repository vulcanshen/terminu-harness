# terminu-harness

[![License](https://img.shields.io/badge/license-GPL--3.0-blue)](LICENSE)

**語言**：[English](README.md) · 繁體中文

**給照 [terminu design](https://github.com/vulcanshen/terminu)（tdp）做的 TUI app 用的 Claude Code plugin。** 它判一個 app 有多少 terminu 的 gene：tdp 清單上 app 符合的每一條，就是一個 gene。

## 它能做什麼

`/terminu:judge` 在自己帶的終端機裡打開 app，讀它的說明，照清單每一條寫的驗法按鍵、看畫面。要讀程式或其他檔案的條目，拿得到 app 的 repo 時才判。判完寫一份成績單：

- 每一條的判定（符合、偏離、違反、不適用、判不了）與看到的證據：畫面的文字、動畫的錄影、顏色、圖片
- 照那一版 tdp 數出的 gene，寫成 `tdp-0.7.4 860 gene`
- app 做了、tdp 還沒規範的東西，回報給 tdp

它能判任何 TUI app：不一定是 terminu 家族的成員、不管用什麼語言寫、拿不拿得到程式碼都行。不需要 tmux 或圖形介面，也不動 app 的 repo。

## 需求

- [Claude Code](https://claude.com/claude-code)
- macOS 或 Linux
- [Go](https://go.dev/dl/)：plugin 第一次啟動時，從原始碼建出它的終端機

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
/terminu:judge <app 怎麼啟動> [app 的 repo] [tdp 版本]
```

例：`/terminu:judge demou ~/src/demou v0.7.4`。沒給 repo，要讀程式的條目判「判不了」；沒給版本，照最新的 tdp 判。

它先認識 app（說明、`Space`、`?`、每個畫面與 popup），再把清單分給十個 judge 同時判，最後告訴你 gene 數、app 違反了什麼、哪些判不了、成績單在哪。判完一次要一段時間，也用掉不少 token。

## 存放位置

| 什麼 | 在哪 |
|---|---|
| 成績單 | `~/.terminu-harness/judge/<app>/<時間>-tdp-<版本>/`：`report.md`、每一條的判定、證據 |
| 抓下來的 tdp 檔案 | `~/.terminu-harness/cache/tdp/` |
| 建好的終端機 | `~/.terminu-harness/bin/` |

| 環境變數 | 作用 |
|---|---|
| `TERMINU_HARNESS__HOME` | 全部改放到 `~/.terminu-harness` 以外的地方 |
| `TERMINU_HARNESS__OUT` | 成績單改放到別處 |
| `TERMINU_HARNESS__TDP_REPO` | 從本機的 terminu clone 讀 tdp，不從 GitHub 抓 |

## 限制

- 不支援 Windows。
- 終端機不送滑鼠事件：要用滑鼠驗的條目判「判不了」。
- 從畫面量的時間（動畫、toast 停留多久）先以 33ms 為容許誤差，等 tdp 定下來。

## 相關連結

- [terminu](https://github.com/vulcanshen/terminu)：terminu design（tdp）與它的清單
- [terminu-go](https://github.com/vulcanshen/terminu-go)：terminu 家族 app 用的 Go lib
- [CHANGELOG](CHANGELOG.md)
- [docs/dev-remarks.md](docs/dev-remarks.md)：運作方式，給開發者看

## License

[GPL-3.0](LICENSE)
