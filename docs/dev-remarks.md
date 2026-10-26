# terminu-harness 開發者備忘

terminu 的 harness：一個 Claude Code plugin，照 terminu design（tdp）用 terminu-go 做 TUI app，寫程式時就做對。harness 放流程，不寫 app 的規定；規定一律在 tdp，程式判得了的檢查放在 terminu-go。

## 運作方式

- **紀錄**：harness 管的 app 在 repo 頂層放 `harness-assets/`：`idea.md`（使用者的需求）、`spec.yaml`（每個需求對到的 tdp、設計
  決定；idea.md 或 tdp 改了 id 加一）、`work.yaml`（怎麼做的、哪些測試作保）、`review.yaml`（最近一次 review）。格式的唯一一份在
  `skills/newapp/harness-assets.md`，其他 skill 都指到它。
- **入口** `/terminu`（`skills/terminu`，frontmatter `name: terminu`，所以裸的 `/terminu` 叫得到；實測過）：有
  `harness-assets/spec.yaml` 就叫 `terminu:newapp` 接著做；是 terminu-go 的 app 但沒有 harness-assets，問要不要接手（analysis）；
  都不是就問要做新 app、review、發版還是回報。
- **newapp**（`skills/newapp`，主 session 照著做）：讀 app 的 tgo 對到的那一版 tdp → 一次一題問需求，寫進 idea.md → 對 tdp 寫 spec.yaml
  就開始做（不等確認）→ 照 tdp project 範本開 repo（多放 `.claude/rules/terminu.md`，內容在 `skills/newapp/terminu-rule.md`）→ 一個 task 一個 task 寫，每個 `make check`、
  用 harness 的終端機開起來看、記進 work.yaml、跑 `check`、給使用者看、commit → 卡住叫 `terminu:reportgap`，remark 記 issue，程式標
  `TODO(tdp #N)`。「繼續」先查 tgo 新版（要升就統整中間各版 tdp 的 CHANGELOG、開新一份 spec）、remark 裡的 issue、TODO。
- **analysis**（`skills/analysis`）：接手先用 tdp 與 tgo 做的 app。從 README 寫 idea.md 與 spec.yaml，從 git commit history 寫
  work.yaml（測試名字從 `go test -list` 拿，用 `git log -S'func TestX('` 找加進它的 commit，對回 task），不讀程式碼也不讀舊的開發者
  文件；找不到的寫「不知道」。放 `.claude/rules/terminu.md`，刪掉舊的開發者文件（只刪開發者文件），接著 review。
- **review**（`skills/review`）：`check` → 使用者重新確認 `user` 作保 → 派 `terminu:reviewer`（`agents/reviewer.md`，只讀 harness-assets
  與 tdp，不讀程式碼，工具只有 Read、Grep、Glob、`tdp_section`、`tdp_file`）判每一個 task → 寫 review.yaml → 不合的一項一項問使用者。
- **release**（`skills/release`）：照 tdp project.7：工作目錄乾淨、`make check`、review（還有不合的問使用者照樣發還是繼續打磨）、
  `## Unreleased` 改成 `## vX.Y.Z — tdp vA.B.C`、問使用者才 tag 與 push。
- **reportgap**（`skills/reportgap`）：判斷發到 tdp 還是 terminu-go、查重複、照固定格式寫草稿、使用者同意才 `gh issue create`，記進
  spec.yaml 那個 task 的 remark。
- **檢查**（`internal/assets`，MCP 工具 `check`、命令列 `terminu-harness check`）：四個檔彼此對不對得上；app 照的 tdp 版本從 go.mod 的
  terminu-go 讀它 CHANGELOG 的 `## vX — tdp vY` 推（每一版 tdp 都有一版 tgo 對到它），spec.yaml 的 tdp 編號拿那一版 tdp 標題裡的編號對；
  work.yaml 的測試拿 `go test -list` 對、commit 拿 git 對。另外列出待做。
- **commit hook**（`hooks/hooks.json` → `scripts/hook-commit.sh` → `terminu-harness hook-commit`，`internal/hook`）：PreToolUse 掛在
  Bash 上。shell 先濾掉不提到 commit 的指令；harness 管的 repo 裡，commit 改了程式、work.yaml 沒改、訊息沒有 `Terminu-Work: unchanged`
  就 exit 2 擋下並說明。
- **MCP server**（`cmd/terminu-harness`，子指令 `mcp`）：plugin.json 的 `mcpServers` 跑 `scripts/mcp.sh`（不放根目錄的 `.mcp.json`：在這個 repo 裡開 Claude Code 時它會被當成專案的 MCP 設定，`${CLAUDE_PLUGIN_ROOT}` 是空的就啟動失敗），server 名 `harness`，工具在
  Claude Code 裡叫 `mcp__plugin_terminu_harness__<工具>`：`open`、`keys`、`screen`、`styles`、`image`、`resize`、`overflow`、`close`、`tdp_file`、
  `tdp_section`（照編號回那一節的原文與「出處」是它的清單條目）、`check`。
- **終端機**（`internal/term`、`internal/vt`）：`creack/pty` 開 app，`vt` 是從 terminu-go 複製的模擬器。讀 PTY 與處理畫面分兩個
  goroutine：一個只管讀（app 寫 PTY 不會被卡住），一個把讀到的全部一次餵給畫面，再用每一列的 checksum 找出變了哪幾列。錄影記的
  就是這些變化，時間取那一批第一塊輸出讀到的時刻；6ms 內的變化併成一格。
- **等畫面穩定**：從送鍵之後、或之後最後一次變化起算，150ms 沒變就回傳（`open` 先等到第一次變化再等 300ms）。從上一次變化
  起算會在 app 還沒反應前就回傳。
- **tdp 的檔案**（`internal/tdp`）：從 GitHub 抓那一版 tag 的檔案，存在 `~/.terminu-harness/cache/tdp/`；設
  `TERMINU_HARNESS__TDP_REPO` 時用本機 clone 的 `git show`。

## 設計決定

- **harness 是封閉的，tdp 是宇集**（user 2026-10-25）：app 的每個決定都要在 tdp 找得到根據，或在 tdp 明文寫「app 可以自行決定」的
  範圍裡；找不到、不清楚、矛盾只能回報 tdp。tdp v0.10.0 照這個拿掉了「偏離」、開發者文件的規定，每一節都給了編號。
- **review 不讀程式碼，拿 tdp 對紀錄**（user 2026-10-25）：用 tgo 做的 app 元件本身一定照 tdp，讀程式碼找問題是白工；harness 要管的
  是 app 的機制有沒有違反 tdp。程式對不對交給測試（work.yaml 的作保）。代價：work.yaml 沒寫的 review 看不到，harness 不保證。
- **reviewer 跟寫程式的一方分開**：agent 從乾淨的 context 開始，看得到的只有紀錄與 tdp。不另外記「駁回過的問題」：使用者選了先不處理
  的留在 review.yaml，下一次照樣列出來。
- **app 照的 tdp 版本不另外存**：從 go.mod 的 terminu-go 推。tdp 只多一句放行的小版本也要有自己的 tgo 版本（user 選 A，2026-10-26），
  因為放行本身就是 review 的根據。
- **接手只看 README 與 git history**（user 2026-10-26）：接手時的 spec 只要寫出現在有哪些功能，那是 README 本來就寫的；做法的根據是
  引的 tdp 編號，不是舊文件裡寫的理由，所以不讀 dev-remarks 與設計文件。原本叫 genwork、也讀那些文件，改成 analysis。
- **app 的規則檔放 app 自己的 `.claude/rules/terminu.md`**（user 2026-10-26，在 terminu 定的），不由 plugin 在 SessionStart 帶進去：
  newapp 開 repo 時放，analysis 接手時放。家族成員的接手在成員自己的 session 裡跑，算成員自己改 `.claude/`。
- **寫程式時就做對，不做完才判**（user 2026-10-20）：原本的 judge（做完才判、一千三百多條、11 個 agent 看畫面）拿掉了。程式判得
  了的檢查放在 terminu-go 的測試 package，app 用哪一版 tgo 就跑哪一版的檢查；要判斷的交給 reviewer（發版前與手動叫）。harness 提
  要 tgo 做哪些檢查與「刁鑽的值」。judge 的程式與文件在 `fb89ad1` 還找得到。
- **模擬器複製 terminu-go 的 `vt.go`**（GPL，檔頭註明出處與 commit），不用 charmbracelet 的 x/vt：icon 與 powerline 各佔幾格是
  vt.go 本來就有的設定，量法跟家族 app 同一套；x/vt 把 icon 一律算一格。改了的地方：寬度跟著每個畫面、數折行與捲動、拿掉
  render（harness 直接讀格子）。terminu-go 改 vt.go 時要手動跟。
- **看畫面預設只回文字**，顏色另外指定範圍查：demou 一個 120×40 的畫面樣式有 269 段，每次全回多約 5k token。
- **錄影只在畫面變的時候記**，不固定間隔拍：不漏格、不重複、資料量小。
- **字型自己帶**（Go Mono、Symbols Nerd Font Mono），框線與方塊字元自己畫：跑 harness 的機器不一定有字型；家族的粗框 `━┃┏` 不在
  Go Mono 裡。
- **環境變數**：開 app 時拿掉真終端機的 `TERM`、`*__ICON_WIDTH`、`*__POWERLINE_WIDTH`，app 才會探測 harness 這個終端機；要測
  別的寬度時自己用 `env` 設。
- **會寫系統檔的 app 用沙盒開**：用 app 自己的環境變數（tdp E3 的覆寫）把 `HOME` 與 app 的檔指到暫存目錄，不碰真的檔案。

## 已知的牆與未做

- 只支援 macOS、Linux（`creack/pty`，沒有 ConPTY）。
- 終端機不送滑鼠事件。
- 使用者要有 Go：第一次啟動時 `scripts/mcp.sh` 從原始碼建 binary，可能超過 Claude Code 等 MCP server 的時間。
- 等 issue 不是一直開著的 session：Claude Code 沒有可靠地等好幾天的方式（Monitor 每次最多 30 分鐘、`/loop` 要 session 開著、routine 沒有 issue 事件），所以靠使用者回到 app 的資料夾打 `/terminu`。
- 使用者在終端機自己打的 commit 不經過 hook，下一次 review 才看得到。
- 還沒做：terminu-go 下一版有了 `CheckRepo` 之後，newapp 的骨架改用 `CheckRepo`（E6 的 splash 測試不用 app 自己寫）。

## 用 hostu 試整套流程（2026-10-18～19）

newapp 照 user 給的需求做了 hostu（只在本機；第一版 2026-10-19 移除，同一天照更新後的 harness 重做）。學到的：

- newapp 問需求時要先讀用到的元件的「用途」段，不只看「長相」：spec 原本用 textarea 編 hosts 檔，tdp 的 textarea 寫明設定檔要
  交給 `$EDITOR`。
- 開發中用 harness 的終端機開 app，抓到 L4 違反（Pager 照原樣畫 Tab），畫面測試量不出來（`terminu.Width` 把 `\t` 算 0 格）。
- 第一版直接編 hosts 檔，user 要的是 TUI 的做法：app 有自己的格式存設定（`profiles.yaml`，一個 domain 一筆、可以有多個位址），
  啟用時才匯出成系統檔，也能把現有的系統檔解析回來。newapp 決定設計時，管理某種檔案的 app 先想「自己的格式＋匯出」。
- 試流程的 app 不要一路跟著改：最後的樣子是一路修出來的，測不出「從需求做到完成」的整套流程（user：「現在 hostu 已經有形狀了，
  這種流程不太對了」）。規範或 harness 大改之後，從需求重做一次。
- 重做的那一版判過一次（tdp v0.8.2）：12 條違反裡 10 條是 app 自己的，2 條是 terminu-go 的；這份結果是 2026-10-20 決定「檢查放
  terminu-go、判斷交給 reviewer」的依據（CLAUDE.md「方向」）。

## 建置與開發

- `make build`：建到 `.local/bin/terminu-harness`，`scripts/mcp.sh` 優先用它。有原始碼比它新時不用，改從原始碼建：2026-10-26 hostu
  試用時用到 09:46 建的舊檔，比認得範本路徑的 `d5889fa` 早，check 把 `tdp project/install.sh` 報成找不到。改了 server 的程式之後，
  已經開著的 session 要重新連 MCP（`/mcp`）才換成新的。
- `make test`；設 `TERMINU_HARNESS__TDP_REPO=~/Documents/sideproj/terminu` 時也驗從本機 clone 讀 tdp 的檔案。
- server 的測試開 `internal/server/testdata/fixture`（測試時才建的小程式，不靠外部套件）代替真的 app：它照家族 app 的方式探測
  icon 寬度、顯示收到的鍵、按 `a` 播 5 格動畫、按 `w` 印超過邊界、按 `q` 結束。驗 MCP 工具用它，不拿家族成員來試。
- `make check`：gofmt、go vet、測試、`claude plugin validate`。
- 不裝 plugin 試：`claude --plugin-dir <這個 repo>`。這個 flag 只管那一次啟動，`--resume`、`-c` 接回舊的 session 時也要再帶，
  沒帶就沒有 plugin，也沒有 harness 的 MCP（2026-10-19 在 harness 的 session 踩到：`claude --resume -c` 開的 session 沒有 harness
  的工具）。`claude --plugin-dir <repo> mcp list` 看得到 `plugin:terminu:harness … Connected` 就是連得上。

## 發布

還沒發過。GitHub repo、tag、release 的 binary 怎麼給使用者，發版時跟 user 定。
