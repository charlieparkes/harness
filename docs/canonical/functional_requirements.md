```
date: 2026-09-04
version: v1.0.0
```

---

# Functional Requirements

These requirements are for brainstorming purposes, and are not concrete absolutes which must be implemented. Use as a reference for the designer's intent and thought process.

- Running `harness` with no arguments should default to running its TUI.
- When `harness` starts in TUI, it should open to a list of selectable tasks.
- Task list should be ordered first by "active" tasks (in progress and not canceled), then ordered by last activity.
- Task list should hover/focus the first item in the list, if any tasks are available.
- An information bar at the bottom of the TUI should list available keyboard shortcuts.
- The task list should provide keyboard shortcuts for:
    - c: create
    - x: cancel (currently focused/hovered task)
- `harness` should step through a "workspace, plan, execute, review, apply" workflow, with separate models and prompts defined for each phase of the workflow.
    - workspace: if target path is git repository, create a worktree in $HOME/.harness/worktree/{id}, and save ID to task table
    - plan: using planner agent definition (model, prompt), create plan, iteratively refine plan in configurable number of passes, each pass is a new agent / new context window, and is provided plan created by last iteration
    - execute: using executor agent definition (model, prompt), execute plan steps in DAG order, as defined by plan
    - review: using reviewer agent definition (model, prompt), iteratively review the worktree output
- `harness` should write state of active tasks sqlite database in $HOME/.harness/.
- When asked to create a task, a series of questions asked in overlay/modal windows should be asked.
    1. Target Path: default to path where `harness` was run (directory nav modal, or allow manual entry)
    2. Describe Task: free text input, no size limit
    3. Automation Settings: select which steps of the workflow to run without user input. For example, if the workspace step is marked as "automate", upon creating the task, immediately start the workspace step.


## TODO

- git worktrees
- agent sandboxing
- audit logs (log tool calls, agent conclusions after thinking)
- approved tools list
- auto review tool calls
- step DAG with parallel agents ("allow multitasking plan" option?)