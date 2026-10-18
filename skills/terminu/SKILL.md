---
name: terminu
description: terminu harness 的入口：問使用者要做什麼（做一個新的 TUI app、繼續做到一半的 app、判一個 app 有多少 terminu gene），再帶進對應的流程。user 打 /terminu 時使用。
argument-hint: [想做的事]
---

# /terminu

terminu harness 照 terminu design（tdp）用 terminu-go 做 TUI app，做完用 tdp 的清單檢查。這裡只決定接下來走哪個流程。

1. `$ARGUMENTS` 已經說了要做什麼，就直接走那個流程，把使用者的原話一起帶過去。
2. 沒說的話，問一題（AskUserQuestion）：「要做什麼？」
   - **做一個新的 app**：問需求、確認、開發、完成後檢查
   - **繼續做到一半的 app**：目前目錄有 `.local/terminu/progress.md` 時，把這個排第一個並標成推薦
   - **判一個 app**：數一個現成的 TUI app 有多少 terminu gene
3. 做新 app、繼續做 → 用 Skill 工具叫 `terminu:new-app`。
4. 判 app → 用 Skill 工具叫 `terminu:judge`。
