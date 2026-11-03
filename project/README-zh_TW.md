# terminu project

**Language**: [English](README.md) · 繁體中文

terminu project 是家族成員的 repo 長什麼樣 —— 目錄、Makefile、版本、CHANGELOG、發布 —— 與它的範本檔。
[tdp](https://github.com/vulcanshen/terminu) 講 app 應該要怎樣，跟用什麼語言、什麼建置工具無關；這一頁是家族照 Go 與
goreleaser 做的維運，所以放在 terminu-harness —— 產生與發布家族 app 的工具 —— 裡。tdp 的 E 章開頭與
[E6](https://github.com/vulcanshen/terminu/blob/main/principle/rules-zh_TW.md#e6-發布測試與家族資產) 把家族成員的 repo 交給這裡。

- **從 tdp 搬來的**：tdp v1.2.0-beta.1 起（user 2026-11-03）。之前在 tdp 的 `project/`，最後一版是 tdp v1.1.10；節編號
  project.1–8 與清單 `project-` 的編號照舊。
- **照哪一份**：app 照的 tdp（go.mod 的 terminu-go 對到的那一版）還有 `project/` 時（v1.1.10 以前），照那一版 tdp 的，由那一版
  terminu-go 的 tdptest 檢查；沒有時照這裡，由 harness 的 `check` 拿目前裝的這一版 harness 的範本比。
- **引用**：app 的檔案引用 tdp 寫 `tdp` 加編號；引用這一頁寫 `harness` 加編號：節 `harness project.3`、整頁 `harness project`、
  範本 `harness project/install.sh`、清單的條目 `harness project-13`（[`items-zh_TW.md`](items-zh_TW.md)）。這一頁會被抄進
  app 的段落也照這樣寫。

## project.1 什麼是規定

這一頁只綁家族成員（tdp 的 E 章開頭）；對家族成員，這裡的敘述都是規定。
不是規定的只有三種：標成「例」的例子、「為什麼」、標成「實作參考（不是規定）」的段落。寫「app 可以自行
決定」的地方是明文放行，範圍內還要守的規定寫在旁邊（[tdp P5](https://github.com/vulcanshen/terminu/blob/main/principle/README-zh_TW.md#p5-明文放行)）。

範本檔也是規定：照原樣複製，只改「範本檔」表上最後一欄寫的地方。

**照 tdp 做的 app**（不是家族成員）：repo 的目錄、Makefile、版本放在哪、CHANGELOG、發布的步驟，app 可以自行決定；這一頁與範本
可以拿來當起點，複製之後怎麼改都可以。範圍內還要守的是 tdp 的 E 章（例：E1 的 `version`、E6 的 `install.sh` 與 tap、E6 的
畫面測試）。

**為什麼**：這一頁的價值都落在同時維護整個家族的人身上 —— 一個 app 踩過的發版問題修一次，全家族都拿到；家族共用一個
Homebrew tap；跨 repo diff 一下就知道誰走樣。自己做一個 app 用不到這些，綁住也沒有好處（user 2026-10-27）。

## project.2 範本檔

| 檔 | 複製到 | app 自己改的 |
|---|---|---|
| [`family.mk`](family.mk) | `family.mk` | 不改 |
| [`release.yml`](release.yml) | `.github/workflows/release.yml` | 出 Windows 版時，測試的 matrix 加 `windows-latest` |
| [`goreleaser.yaml`](goreleaser.yaml) | `.goreleaser.yaml` | `<app>`、標著 `app:` 的幾行；出 Windows 版時，打開標著 `Windows` 的註解 |
| [`install.sh`](install.sh) | `install.sh` | `<app>`、`WINDOWS`、標著 `app's own` 的段落 |
| [`uninstall.sh`](uninstall.sh) | `uninstall.sh` | `<app>`、標著 `app's own` 的段落 |
| [`version.go`](version.go) | `internal/version/version.go` | `<app>` |
| [`gitignore`](gitignore) | `.gitignore` | `<app>`；app 自己的行接在最後 |

- **照 tdp 做的 app**（不是家族成員，project.1）拿範本當起點時，範本裡的 `vulcanshen`（repo 的路徑、tap 的 owner、首頁）換成
  自己的帳號。
- **標記管到哪裡**：`goreleaser.yaml` 的 `# app:` 註解照原樣留著，管的是緊接在它底下的那一個鍵的值，清單或整塊都算（`ignore`、
  `description`、`caveats`）；註解講的鍵範本裡還沒有時（`dependencies`），在註解那裡加上那個鍵；其他的鍵（例：`test:`）照原樣。
  打開 `Windows` 是標著 `# Windows` 的那一列連同它底下同一塊的註解列，只拿掉行首的 `# `，行尾的 `# Windows` 留著。
  `uninstall.sh` 註解裡的 `<APP>` 不是 `<app>`，照原樣。
- **檢查有沒有走樣**：把 app 的名字換回 `<app>`，跟這裡的範本 diff，只該差上表最後一欄的地方。harness 的 `check` 就是這樣比。
- **範本要改，先改這裡**：發一版 harness（CHANGELOG 寫出改了哪個範本），成員下一次跑 `check`（升版、發版）時照著複製過去，不在
  app 裡改。app 用起來發現範本有問題，回報給 terminu-harness。

**為什麼**：家族本來共用一份 Makefile 範本，抄來抄去就走樣了 —— `fmt` 管的目錄、`check` 跑不跑 `-race`、`make help` 的欄寬、
段名的括號，每個 app 都差一點。同一份檔案，diff 一下就知道哪裡不一樣；改的時候也只有一個地方要改。

## project.3 目錄

```
<app>/
├── cmd/<app>/main.go               main package
├── internal/                       app 自己分的 package
│   └── version/version.go          版本（見「版本」）
├── docs/                           icon.svg、README 的 demo gif（tdp E6、E7）
├── .local/                         git 不收；demos/ 放 VHS 的 tape
├── .github/workflows/release.yml
├── .goreleaser.yaml
├── .gitignore
├── go.mod、go.sum
├── family.mk
├── Makefile
├── install.sh、uninstall.sh
├── CHANGELOG.md
├── README.md、README-zh_TW.md
└── LICENSE
```

- app 本身的 Go 程式只在 `cmd/<app>/` 與 `internal/`。`go install` 的路徑是 `github.com/vulcanshen/<app>/cmd/<app>@latest`
  （照 tdp 做的 app 換成自己的帳號：`github.com/<帳號>/<app>`，tdp 的 E 章開頭）。
- **v2 起，module 的路徑結尾帶大版號**，這是 Go 的規定：kbu v4 是 `github.com/vulcanshen/kbu/v4`，import 與 `go install`
  的路徑（`github.com/vulcanshen/kbu/v4/cmd/kbu@latest`）跟著帶。升到新的大版號（v2 以上）的那一版，在同一個 commit 裡改。
  **為什麼**：不帶的話，Go 不把 v2 以上的 tag 當成這個 module 的版本：kbu 已經是 v4.0.0，`go install …@latest` 裝到的卻是
  v1.7.12，執行檔的 build 資訊也只是 pseudo-version（user 2026-10-28）。
- **開發用的小工具**（Go 寫的也一樣）可以放頂層的 `tools/`，不算 app 本身。
- 圖上沒有的其他頂層目錄與檔案（例：端到端測試），app 可以自行決定。
- **Claude Code 的檔案不進 git**：`CLAUDE.md` 與 `.claude/` 底下的檔案不追蹤（範本 [`gitignore`](gitignore) 排除了），
  harness 的 `.claude/rules/terminu.md` 除外。工作目錄裡可以有。**為什麼**：那是維護者自己的指示與設定，寫滿一台機器的路徑；
  公開的 repo 用不到，還把維護者的路徑公開了出去（user 2026-10-28）。
- **`docs/icon.svg` 是家族 mark 的正本**。Go 的 `go:embed` 讀不到 `docs/`，所以 splash 用的那一份放在程式裡（`internal/` 裡的副本檔，
  或寫死的格子），tdp E6 的測試守住它跟 `docs/icon.svg` 一模一樣。
- 要不要出 Windows 版，app 可以自行決定；出的話另有 `install.ps1`、`uninstall.ps1`，不在範本裡，怎麼寫 app 可以自行決定。

**為什麼**：家族裡大多數 app 本來就是 `cmd/<app>` 加 `internal/`，照它搬的 app 最少。程式只在兩個地方，根目錄只剩 repo 層級的
檔案，打開任何一個家族 repo 都知道東西在哪裡。

## project.4 版本

- 版本放在 `internal/version`（範本 [`version.go`](version.go)）：`Version` 是不帶 v 的版號，`Display()` 印出時補上 v。
- goreleaser 發版時注入 tag 去掉 v 的版號（例：tag `v0.2.0` → `0.2.0`）。`<app> version` 與 splash 印 `Display()`（例：
  `v0.2.0`；tdp E1、S5）。
- 本機的 `make build`、`go build` 不注入，版本是 `dev`。

**為什麼**：`dev` 一眼看得出不是發布版。要比較版本新舊的地方（例：ccdrawu 的更新檢查）用不帶 v 的 `Version`，遇到 `dev`
就知道是本機版、不去比。

## project.5 Makefile

app 的 `Makefile` 只設變數、include [`family.mk`](family.mk)，自己的 target 寫在最後一段：

```make
# <app> — build / run / test / package（家族共用的 target 在 family.mk）
#
# 跑 `make`（或 `make help`）列出所有指令。

BINARY := <app>
include family.mk

##@ <app> 自己的
```

**變數**（在 `include` 之前設）：

| 變數 | 意思 | 預設 |
|---|---|---|
| `BINARY` | 執行檔的名字，也是 `cmd/` 底下那個目錄的名字 | 必填 |
| `CHECK` | 接在 `check` 後面跑的 app 自己的 target（例：ccdrawu 的 `mod-check`） | 空 |
| `GOFMT_DIRS` | gofmt 管的目錄 | `cmd internal` |
| `TAPES` | `make gif` 跑的 VHS tape | `.local/demos/demo.tape` |
| `VHS` | vhs 的路徑 | `vhs` |

**family.mk 的 target**，`make help` 照這個順序列，app 自己的段列在最後：

| 段 | target |
|---|---|
| 編譯（build） | `build`、`install`、`uninstall`、`tidy` |
| 發布（release） | `release-check`、`snapshot` |
| 執行（run） | `run`、`dev` |
| 測試 / 檢查（test） | `test`、`test-race`、`check-lib`、`vet`、`fmt`、`fmt-check`、`check` |
| 打包（package） | `package` |
| 展示（demo） | `gif` |
| 其他 | `clean`、`help` |

- **`check` 是 commit 前跑的**：`fmt-check`、`vet`、`test-race`；有 `go.work` 時加 `check-lib`（關掉 `go.work`，照 `go.mod` 寫的
  版本編譯、測試）；最後是 `CHECK`。
- **`tidy`、`snapshot` 關掉 `go.work`**：CI 上沒有 `go.work`（git 不收），本機要跟它一樣。
- **app 自己的 target 全放在最後一段** `##@ <app> 自己的`，名字不跟 family.mk 的重複。
- family.mk 的說明是繁中，跟家族原本的 Makefile 一樣。

**為什麼**：CI 跑的是 `go test -race`，`check` 跑同一個，本機過了 CI 才會過。`check-lib` 抓的是「本機用 `go.work` 指著還沒發版
的 lib，過了；照 `go.mod` 就過不了」—— 發版前一定要過。app 自己的 target 集中在最後，看 `make help` 就分得出哪些是家族共用的。

## project.6 CHANGELOG

```markdown
# Changelog

Each release names the version of terminu design (tdp) it follows.

## Unreleased

- …

## v0.1.2 — tdp v0.6.11

- …
```

- **每一版的標題** 是 push 的 tag（帶 v），接著 ` — tdp <版本>`：這一版跟的是哪一版 tdp。
- **beta** 也有自己的標題，就是它帶 pre-release 的 tag（`## v0.2.0-beta.1 — tdp v1.1.0-beta.2`）；之後的正式版那一段寫上一個
  正式版以來的全部改動，beta 的也算，沒裝 beta 的人才看得懂。
- **`## Unreleased`**：改動進來就寫在這一段（沒有就開一個，放在最上面）；發版時把這一段的標題改成這一版的標題。
- **tdp v0.7.0 以前發的版本照當時的寫法**，標題不改。
- 用英文寫。release 的內文就是 CHANGELOG 裡這一版的那一段（見「發布」）。

**為什麼**：標題開頭就是 tag，`release.yml` 用 `$2 == tag` 就抓得到那一段，同一段 awk 全家族通用；寫出 tdp 的版本，一眼看得出
每個 app 跟到哪一版規範。`Unreleased` 讓改動進來時就記下，session 中斷、或一版包了好幾個改動時不會漏寫。舊的標題不改：`release.yml` 只抓新 tag 的那一段，舊標題不影響；舊版當時跟的是
哪一版 tdp 多半查不到，硬補會寫錯。

## project.7 發布

1. `make check` 通過。
2. CHANGELOG 的 `## Unreleased` 改成 `## vX.Y.Z — tdp vA.B.C`。
3. commit，打 tag `vX.Y.Z`。
4. push commit 與 tag。[`release.yml`](release.yml) 接著做：
   - 測試：ubuntu 與 macOS 跑 `go test -race`（出 Windows 版的 app 加 Windows，不帶 `-race`）；
   - goreleaser：建置、發 GitHub release、推 Homebrew tap `vulcanshen/homebrew-tap`（照 tdp 做的 app 推自己的 tap；出 Windows 版的加
     scoop bucket）；
   - 用 CHANGELOG 那一段蓋過 release 的內文。goreleaser 推 tap 失敗時這一步照樣做；release 還沒建好就跳過。

- **beta**：正式版之前，app 可以先發 beta：tag `vX.Y.Z-beta.N`（N 從 1 起），步驟一樣。goreleaser 把它發成 GitHub 的
  pre-release，不推 Homebrew tap 與 scoop bucket（`skip_upload: auto`），所以 brew、scoop 與 `install.sh`（抓最新的
  release）拿到的還是正式版；`go install …@vX.Y.Z-beta.N` 裝得到 beta。**為什麼**：v1.0.0 同一天就得補 v1.0.1（user
  2026-10-27：發 stable 之前要先發 beta）；beta 把整個發版流程走一遍，又不動到 brew、scoop 使用者拿到的版本。
- **release notes 寫在 checkout 外面**（`$RUNNER_TEMP`）：寫在 checkout 裡會讓工作目錄變髒，Go 記進執行檔的 build 資訊
  就成了 `vX.Y.Z+dirty`（`go version -m` 看得到），家族 v1.0.0 那一輪的 release 都是這樣。
- **token**：goreleaser 讀 `GITHUB_TOKEN`，發 release 與推 tap 都用它。workflow 在 repo 的 secret `ACTION_TOKEN`（PAT）有設時
  用它，沒設就用內建的 token —— 內建的寫不了 tap。在本機手動發版時用 shell 的 `GITHUB_TOKEN`。
- **`make snapshot`** 在本機試打包，不發布、不推 tap。
- **`make release-check`** 驗證 `.goreleaser.yaml`。goreleaser 把 `brews` 標成 deprecated，`goreleaser check` 因此回失敗；
  `release-check` 照樣印出這個警告，但只有這一種不算失敗。

**為什麼**：Homebrew 用 formula、不用 cask：cask 裝的檔案會被標上 `com.apple.quarantine`，沒簽章的執行檔在 macOS 上會被
Gatekeeper 擋下；formula 不會。一個 token 走到底，手動發版與 workflow 用的是同一個名字。release 內文另外蓋一次，是因為
goreleaser 的 `--release-notes` 不可靠；推 tap 失敗時照樣蓋，是 kbu 遇過 tap 推送 401、release 的內文因此空著。

## project.8 參考實作

demou：家族的操作教學 app，照這一頁從零做。範本檔用起來的問題由它先遇到，回報給 terminu-harness。
