---
name: terminu
description: terminu harness 的入口：目前目錄有 harness 在管的 app 時，先查它在等的事（新版、回報、TODO）再接著做；是還沒用 harness 管過的 terminu-go app 時，問要不要接手；都不是就問使用者要做什麼，再帶進對應的流程。user 打 /terminu 時使用。
argument-hint: [想做的事]
---

# /terminu

terminu harness 照 terminu design（tdp）用 terminu-go 做 TUI app，紀錄放在 app 的 `harness-assets/`。這裡只決定接下來走哪個流程。

1. `$ARGUMENTS` 說了要做什麼，就直接走那個流程（4. 的對照），把使用者的原話一起帶過去。
2. 目前目錄有 `harness-assets/spec.yaml`（harness 在管的 app）：不用問，用 Skill 工具叫 `terminu:newapp`，它從「6. 繼續」開始：
   先查新版、回報、TODO，講給使用者聽，再接著做。
3. 目前目錄是用 terminu-go 做的 app（`go.mod` 有 `github.com/vulcanshen/terminu-go`），但沒有 `harness-assets/`：問一題
   （AskUserQuestion）「這個 app 還沒用 harness 管過，要接手嗎？」
   - **接手**（推薦）：從 README 與 git history 產生 harness-assets，之後照 tdp review、加需求、發版
   - **做別的**：照 4. 問要做什麼
4. 都沒有的話，問一題（AskUserQuestion）：「要做什麼？」
   - **做一個新的 app**：問需求、開發，做完給你看
   - **review**：照 tdp 判 app 的紀錄
   - **發版**：review、整理 CHANGELOG，問過你才打 tag
   - **回報 tdp 或 terminu-go**：tdp 沒寫、寫不清楚，或 terminu-go 做不到的地方
5. 用 Skill 工具叫對應的 skill：做新 app、繼續做、加需求 → `terminu:newapp`；接手 → `terminu:analysis`，做完接著
   `terminu:review`；review → `terminu:review`；發版 → `terminu:release`；回報 → `terminu:reportgap`。
