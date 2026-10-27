# harness-assets：harness 在 app 裡的紀錄

app 的 repo 頂層放 `harness-assets/`，四個檔。tdp 不管開發者文件（tdp P5），這幾個檔由 harness 定；tdp project.3 允許圖上沒有的
頂層目錄。寫完、改完用 MCP 工具 `check` 驗：四個檔對不對得上、tdp 編號在不在那一版 tdp 裡、測試在不在 `go test -list` 裡。

| 檔 | 是什麼 | 誰寫 |
|---|---|---|
| `idea.md` | 使用者描述的需求本身，跟著互動一直更新 | 照使用者的話寫，不寫設計 |
| `spec.yaml` | 每個需求對到的 tdp 根據、設計決定、取消的原因 | harness 照 idea.md 與 tdp 寫 |
| `work.yaml` | 每個需求實際怎麼做的、用什麼作保 | 寫程式的一方，跟程式同一個 commit |
| `review.yaml` | 最近一次 review 的結果，每次覆蓋 | `/terminu:review` |

spec.yaml、work.yaml 只新增、更新，反映 app 現在的樣子；舊的樣子在 git history。

## idea.md

使用者講的需求，用使用者的話，一段一段寫：要做什麼、給誰用、在哪些平台、哪些事不能發生、參考什麼。使用者之後加的、改的也寫進來。
不寫畫面、按鍵、資料格式這些設計（那是 spec.yaml 的 remark）。

## spec.yaml

```yaml
id: 3            # 第幾份 spec。idea.md 改了、或 tdp 升版，id 加一（同時發生只加一次）
from: v0.10.0    # 上一份 spec 照的 tdp；第一份寫 0；只改 idea.md 時跟 to 一樣
to: v0.10.1      # 這份 spec 照的 tdp：go.mod 的 terminu-go 對到的那一版（check 會比對）
tasks:
  - id: 1                  # 整個 app 一路遞增，不重複使用，取消了也不回收
    spec: 1                # 這個 task 最後一次新增或修改是在第幾份 spec
    request: 側邊是存著的 profile 清單，每個 profile 一列，寫名字與有幾個 domain   # idea.md 裡的那個需求
    tdp: [tdp panel.1, tdp panel.5]   # 這個需求照的 tdp：tdp 加編號
    remark:                # 設計決定（附理由）、tdp 的缺漏、發出去的回報
      - 名字前面的 icon 標出生效的那一個（tdp P5：領域的符號 app 可以自行決定）
  - id: 2
    spec: 2
    request: 匯出時看得到進度
    tdp: []
    remark:
      - note: tdp 沒有進度條；先用 Info toast 寫進度，等 tdp 回覆
        issue: tdp #12
  - id: 3
    spec: 3
    request: 可以釘選 profile
    tdp: [tdp P5]
    canceled: 使用者說用不到（2026-10-26）
```

- **開新一份 spec（id 加一）的時機**：idea.md 改了一次（不管加了幾個需求），或 tdp 升版。新的需求加成新的 task，改到的 task 更新內容、
  `spec` 改成新的 id；tdp 升版時，`from` 改成原本的 `to`，`to` 改成新版。
- **開發中的決定不加 id**：設計決定、remark、做不到而取消，都在現在這份 spec 裡改。
- **一個 task 一個需求**：從使用者看得到的事切；一個畫面、一個動作、一種資料、一個背景會做的事都可以是一個 task。畫面的配置
  （有幾個 panel、各放什麼）也是 task，版面的 ASCII 圖寫在它的 remark。
- **tdp 引用**：一律 `tdp` 加編號（tdp 的「引用」）。
  - rules 用規則編號：`tdp K4`。principle 用 `tdp P5`，術語用 `tdp terms.3`。
  - components 與 project 用節編號：`tdp panel.5`、`tdp project.3`；README 用資料夾名：`tdp components.5`、`tdp input.2`。
  - 一個元件整個都用到時引整個檔：`tdp confirm`、`tdp form`。
  - project 的範本檔寫路徑：`tdp project/install.sh`，等於那一份範本。
  - 要比一整節精確時引清單的條目：`tdp K9-6`、`tdp note-78`，指的是那一條的出處，依據還是規範的文字。
  - 在 tdp「app 可以自行決定」的範圍裡做的，引那個放行（領域熱鍵、領域的符號與顏色、功能邏輯都是 `tdp P5`），範圍內還要守的規定
    也一起引（例：熱鍵還要引 `tdp M5`、`tdp K12`）。
  - 編號寫在 tdp 每一節的標題前面。
- **找不到根據**（tdp 沒寫、寫得不清楚、兩處矛盾）：不自己決定，用 `terminu:reportgap` 回報，在 remark 記 `{note, issue}`。
  有合 tdp 的暫時做法就先照它做，note 寫暫時怎麼做；`tdp` 欄位可以先空著，check 看到 remark 有 issue 就不報。
- **取消**：寫 `canceled: <原因>`，task 留著不刪；work.yaml 那一筆也標 `canceled: true`。

## work.yaml

```yaml
tasks:
  - id: 1
    spec: 1                # 這筆實作照的是第幾份 spec
    implementation: |
      [1] PROFILES 是一般清單（tgo List）：NAME、DOMAINS 兩欄，DOMAINS 是數字靠右；生效的那一個名稱前面畫 nf-md-check_circle（Green）
      順序是建立的順序；空的時候寫空狀態「No profiles」
      讀 $XDG_CONFIG_HOME/hostu/profiles.yaml
    test: [TestProfilesList, TestProfilesEmpty]
  - id: 4
    spec: 2
    implementation: |
      Space menu「刪除」與 [D]：都開 confirm「Delete office?」，Enter 刪、Esc 不動；刪完 Success toast「Deleted office」
      寫 profiles.yaml
    test:
      - TestDeleteFromMenuAsks
      - TestDeleteKeyAsks
      - user: {date: 2026-10-26, how: 在 Windows 上刪一個 profile，看 hosts 檔沒有被動到}
      - commit: 3f2a1c9
```

- **implementation 寫到 review 判得出來**：review 不讀程式碼，只拿這裡對 tdp。寫從哪裡觸發（哪個 panel、哪個鍵、Space menu 的哪一項、
  啟動時、計時）、經過什麼（開什麼 popup、問什麼、成功與失敗給什麼回饋、畫面上的字）、讀寫什麼（哪個檔、哪個位置、外部指令）、
  用了 terminu-go 的哪些元件。一個 task 有好幾個入口時，每個入口各寫一行。tgo 本身就會做的（Space menu 怎麼開、`?` 怎麼畫）不用寫。
  analysis 從 README 與 git history 找不到的部分寫「不知道：<缺什麼>」，review 會判成沒有測試，等寫程式的一方看程式補上。
- **test 是作保**，三種：
  - 測試的名字（`TestDeleteKeyAsks`，子測試寫 `TestX/sub`）：check 拿 `go test -list` 對。每個入口都要有測試走到。
  - `user: {date, how}`：寫不成測試的（例：要一台 Windows），寫使用者怎麼確認的與日期。每次 review 都請使用者重新確認，確認了就更新日期。
  - `commit: <hash>`：連使用者確認都做不到時，指出做這件事的 commit 當紀錄。
  - 沒有作保的 task，review 判 `沒有測試`。
- **review 只判寫在這裡的**：程式做了、這裡沒寫的，review 看不到，harness 不保證。
- **做到哪**：spec.yaml 有、這裡沒有的 task 還沒做；task 的 `spec` 比這裡的 `spec` 大，需求改過、要改。check 會列出來。

## review.yaml

```yaml
spec: 3          # 這次 review 對的 spec
tasks:
  - id: 1
    accepted: true
  - id: 4
    accepted: false
    reason: 違反：confirm 的問句不在最後（tdp confirm.2「問句放最後，下面接著就是 hint」）
  - id: 5
    accepted: false
    reason: 沒有測試：[D] 那個入口沒有測試走到
  - id: 6
    accepted: false
    reason: tdp 不清楚：背景工作的進度要用哪個元件，tdp 沒寫
```

- 每次 review 覆蓋整個檔，列 work.yaml 裡每一個沒取消的 task。
- `accepted: false` 的 reason 開頭寫是哪一種：`違反`、`沒有測試`、`tdp 不清楚`，接「：」。三種的處理不一樣：違反要修掉或回報 tdp；
  沒有測試要補測試或請使用者確認；tdp 不清楚要回報 tdp。
- 改紀錄、改程式、補測試處理得了的，review 自己處理，不問使用者；還沒處理完的（在等 tdp 的回報、使用者沒重新確認、使用者不發
  回報的），下一次 review 還是 false，一直列著。
