# terminu project：條目

**Language**: [English](items.md) · 繁體中文

[terminu project](README-zh_TW.md) 的條目，從 tdp v1.1.10 的清單搬來，編號照舊（引用寫 `harness project-13`）。每一欄的意思與適用照 tdp 的
[checklist](https://github.com/vulcanshen/terminu/blob/main/checklist/README-zh_TW.md)；只查家族成員（project.1）。

| 編號 | 出處 | 檢查什麼 | 適用 | 怎麼驗 |
|---|---|---|---|---|
| project-1 | project.2 | `family.mk` 照範本原樣複製到 `family.mk`，不改 | 都要 | 讀檔案：跟 harness 的 `project/family.mk` diff，沒有差 |
| project-2 | project.2 | `release.yml` 照範本複製到 `.github/workflows/release.yml`，只有出 Windows 版時在測試的 matrix 加 `windows-latest` | 都要 | 讀檔案：跟 harness 的範本 diff，只差 matrix 的 `windows-latest`（出 Windows 版時） |
| project-3 | project.2 | `goreleaser.yaml` 照範本複製到 `.goreleaser.yaml`，只改 `<app>`、標著 `app:` 的幾行、出 Windows 版時打開標著 `Windows` 的註解 | 都要 | 讀檔案：把 app 名換回 `<app>` 跟 harness 的範本 diff，只差這些地方 |
| project-4 | project.2 | `install.sh` 照範本複製，只改 `<app>`、`WINDOWS`、標著 `app's own` 的段落 | 都要 | 讀檔案：把 app 名換回 `<app>` 跟 harness 的範本 diff，只差這些地方 |
| project-5 | project.2 | `uninstall.sh` 照範本複製，只改 `<app>`、標著 `app's own` 的段落 | 都要 | 讀檔案：把 app 名換回 `<app>` 跟 harness 的範本 diff，只差這些地方 |
| project-6 | project.2 | `version.go` 照範本複製到 `internal/version/version.go`，只改 `<app>` | 都要 | 讀檔案：把 app 名換回 `<app>` 跟 harness 的範本 diff，沒有其他差 |
| project-7 | project.2 | `gitignore` 照範本複製到 `.gitignore`，只改 `<app>`，app 自己的行接在最後 | 都要 | 讀檔案：把 app 名換回 `<app>` 跟 harness 的範本 diff，只差最後加的行 |
| project-8 | project.2 | 出 Windows 版時，`release.yml` 測試的 matrix 加了 `windows-latest` | 出 Windows 版 | 讀檔案：`.github/workflows/release.yml` 的 matrix |
| project-9 | project.2 | 出 Windows 版時，`.goreleaser.yaml` 標著 `Windows` 的註解都打開 | 出 Windows 版 | 讀檔案：`.goreleaser.yaml` |
| project-10 | project.2 | 出 Windows 版時，`install.sh` 的 `WINDOWS` 設成 `yes`（範本註解） | 出 Windows 版 | 讀檔案：`install.sh` 的 `WINDOWS` |
| project-11 | project.3 | app 本身的 Go 程式只在 `cmd/<app>/` 與 `internal/`；開發用的小工具可以放頂層的 `tools/` | 都要 | 讀檔案：`tools/` 以外的 `.go` 檔都在這兩個目錄底下 |
| project-12 | project.3（圖） | main package 在 `cmd/<app>/main.go` | 都要 | 讀檔案 |
| project-13 | project.3 | `go.mod` 的 module 是 `github.com/<帳號>/<app>`，v2 起結尾帶最新的 release 的大版號 `/vN`（家族成員的帳號是 `vulcanshen`，照 tdp 做的 app 是自己的帳號） | 都要 | 讀檔案：`go.mod` 的 module 是 `github.com/<帳號>/<app>`，最新的 release 是 vN.x.y、N ≥ 2 時結尾帶 `/vN`；`go list -m <module>@latest` 是最新的 release |
| project-14 | project.3（圖） | README 的 demo gif 放在 `docs/` | 都要 | 讀檔案：README 引用的 gif 在 `docs/` 底下 |
| project-15 | project.3（圖） | 根目錄有 `go.mod`、`go.sum` | 都要 | 讀檔案 |
| project-16 | project.3（圖） | 根目錄有 `LICENSE` | 都要 | 讀檔案 |
| project-18 | project.3 | 家族 mark 的正本是 `docs/icon.svg`；splash 用的那一份在程式裡（`internal/` 裡的副本檔，或寫死的格子） | 都要 | 讀檔案：有 `docs/icon.svg`；讀程式：splash 畫的是 `internal/` 裡的副本檔或寫死的格子 |
| project-19 | project.3 | 出 Windows 版的 app 另有 `install.ps1`、`uninstall.ps1` | 出 Windows 版 | 讀檔案 |
| project-37 | project.3 | 沒有追蹤任何 `CLAUDE.md`，`.claude/` 底下也只追蹤 `.claude/rules/terminu.md` | 都要 | 跑 `git ls-files \| grep -E '(^\|/)CLAUDE\.md$\|^\.claude/'`：除了 `.claude/rules/terminu.md` 什麼都沒有 |
| project-20 | project.4 | `<app> version` 與 splash 的版本印的是 `internal/version` 的 `Display()` | 都要 | 讀程式：兩處都呼叫 `Display()` |
| project-21 | project.4 | 發布版的版本是 tag 去掉 v 注入的，印出來帶 v（例：tag `v0.2.0` 印 `v0.2.0`） | 都要 | 讀檔案：`.goreleaser.yaml` 把 tag 去掉 v 注入、`internal/version` 印出時帶 v；發過版的，裝發布版跑 `<app> version`，印的跟 tag 一樣 |
| project-22 | project.4 | 本機 `make build`、`go build` 不注入版本，版本是 `dev` | 都要 | 本機 `make build` 後跑 `./<app> version`，印 `dev` |
| project-23 | project.5 | `Makefile` 只設變數、`include family.mk`，app 自己的 target 全寫在最後一段 `##@ <app> 自己的` | 都要 | 讀檔案；跑 `make help`，家族共用的段照 family.mk 的順序、app 自己的段在最後 |
| project-24 | project.5 | 變數在 `include family.mk` 之前設 | 都要 | 讀檔案 |
| project-25 | project.5 | `BINARY` 設成 app 的名字，也是 `cmd/` 底下那個目錄的名字 | 都要 | 讀檔案 |
| project-26 | project.5 | Makefile 開頭的兩行註解照範本（`<app> — build / run / test / package（家族共用的 target 在 family.mk）`、「跑 make（或 make help）列出所有指令。」） | 都要 | 讀檔案的開頭 |
| project-27 | project.5 | app 自己的 target 名字不跟 family.mk 的重複 | 都要 | 讀檔案：比對兩邊的 target 名 |
| project-28 | project.6 | `CHANGELOG.md` 開頭照範本：`# Changelog`，接著 `Each release names the version of terminu design (tdp) it follows.` | 都要 | 讀檔案的開頭 |
| project-29 | project.6 | 每一版的標題是 push 的 tag（帶 v）接 ` — tdp <版本>` | 發過版 | 讀檔案：tdp v0.7.0 之後發的每一版標題 |
| project-30 | project.6 | 還沒發版的改動寫在最上面的 `## Unreleased` | 都要 | 讀檔案與最新 tag 之後的 commit：改動都記在最上面的 `## Unreleased` |
| project-31 | project.6 | tdp v0.7.0 以前發的版本，標題照當時的寫法、不改 | 發過版 | 讀檔案：`git log -p CHANGELOG.md`，舊版的標題沒被改寫 |
| project-32 | project.6 | CHANGELOG 用英文寫 | 都要 | 讀檔案 |
| project-33 | project.6 | GitHub release 的內文就是 CHANGELOG 裡這一版的那一段 | 發過版 | `gh release view <tag>` 對照 CHANGELOG 那一段 |
| project-34 | project.7 | 發版的 tag 是 `vX.Y.Z`，beta 是 `vX.Y.Z-beta.N` | 發過版 | `git tag` |
| project-35 | project.7 | 現在的 repo 上 `make check` 通過 | 都要 | 在現在的 repo 跑 `make check` |
| project-36 | project.7 | 出 Windows 版的發布時另推 scoop bucket | 出 Windows 版 | 看 scoop bucket 有這一版 |
| project-38 | project.7 | 發過的 beta 在 GitHub 是 pre-release，Homebrew tap（與 scoop bucket）上還是最新的正式版 | 發過版 | `gh release view <beta 的 tag> --json isPrerelease`；看 tap 上 formula 的版本 |
| project-39 | project.7 | 最新一版 release 的執行檔，build 資訊裡沒有 `+dirty` | 發過版 | 下載最新一版的執行檔跑 `go version -m`：main module 的版本就是 tag，`vcs.modified=false` |
