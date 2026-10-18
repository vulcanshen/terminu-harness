# terminu-harness

[![License](https://img.shields.io/badge/license-GPL--3.0-blue)](LICENSE)

**語言**：[English](README.md) · 繁體中文

**給照 [terminu design](https://github.com/vulcanshen/terminu)（tdp）做的 TUI app 用的 Claude Code plugin。** 它判一個 app 有多少 terminu 的 gene：tdp 清單上 app 符合的每一條，就是一個 gene。

## 它能做什麼

`/terminu:judge` 在自己帶的終端機裡打開 app，讀它的說明，照清單每一條寫的驗法按鍵、看畫面。要讀程式或其他檔案的條目，拿得到 app 的 repo 時才判。判完寫一份成績單：

- 每一條的判定（符合、偏離、違反、不適用、判不了）與看到的證據
- 照那一版 tdp 數出的 gene，寫成 `tdp-0.7.4 860 gene`
- app 做了、tdp 還沒規範的東西，回報給 tdp

它能判任何 TUI app：不一定是 terminu 家族的成員、不管用什麼語言寫、拿不拿得到程式碼都行。不需要 tmux 或圖形介面，也不動 app 的 repo。

## 需求

- [Claude Code](https://claude.com/claude-code)

## 安裝

在 Claude Code 裡：

```
/plugin marketplace add vulcanshen/terminu-harness
/plugin install terminu@terminu-harness
```

## 使用方式

```
/terminu:judge
```

## 相關連結

- [terminu](https://github.com/vulcanshen/terminu)：terminu design（tdp）與它的清單
- [terminu-go](https://github.com/vulcanshen/terminu-go)：terminu 家族 app 用的 Go lib
- [CHANGELOG](CHANGELOG.md)

## License

[GPL-3.0](LICENSE)
