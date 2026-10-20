---
name: terminu
description: terminu harness 的入口：目前目錄有做到一半的 app 時，先查它在等的事再接著做；沒有的話問使用者要做什麼，再帶進對應的流程。user 打 /terminu 時使用。
argument-hint: [想做的事]
---

# /terminu

terminu harness 照 terminu design（tdp）用 terminu-go 做 TUI app。這裡只決定接下來走哪個流程。

1. 目前目錄有 `.local/terminu/progress.md`（做到一半的 app）、`$ARGUMENTS` 也沒說要做別的事：不用問，用 Skill 工具叫
   `terminu:newapp`，它從「6. 繼續」開始，先查在等的事、講給使用者聽，再接著做。
2. `$ARGUMENTS` 說了要做什麼，就直接走那個流程，把使用者的原話一起帶過去。
3. 都沒有的話，問一題（AskUserQuestion）：「要做什麼？」
   - **做一個新的 app**：問需求、開發，做完給你看
   - **回報 tdp 或 terminu-go**：tdp 沒規範、寫不清楚，或 terminu-go 做不到的地方
4. 做新 app、繼續做 → 用 Skill 工具叫 `terminu:newapp`；回報 → 用 Skill 工具叫 `terminu:reportgap`。
