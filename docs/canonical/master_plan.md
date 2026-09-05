```
date: 2026-09-05
version: v1.2.1
```

---

# Master Plan

## Purpose

Harness is a local application that manages structured agent workflows. It provides a terminal interface, durable state, worktree isolation, and human control.

This document describes the intended design and architecture. It also records the repository contents and the remaining work.

The functional requirements inform the design. They do not define a fixed product contract.

## Product Direction

Harness converts one task description into a controlled workflow. The main phases are `Worktree`, `Plan`, `Execute`, `Review`, and `Apply`.

Each phase has one responsibility:

- `Worktree` prepares isolated locations for Git targets.
- `Plan` creates, reviews, and refines a structured plan.
- `Execute` completes plan steps in dependency order.
- `Review` examines the result and requests corrections for defects.
- `Apply` writes each reviewed worktree result to its `ReadWrite` target.

The user selects the phase boundaries that require input. Automation can continue through the other boundaries.

Harness stores task and phase state after each material change. A process restart preserves accepted work and does not repeat completed work.

## Design Principles

### Preserve human control

The automation policy determines the pause boundaries. Harness must not perform a destructive or uncertain action beyond this policy.

The TUI shows the current state, available actions, failures, and requests for user input. Background work does not block TUI updates.

### Use durable state

The store package is the source of truth for active and historical data. This data includes tasks, targets, worktrees, plans, agents, and agent activities.

Task and phase rows store current workflow status. They do not retain a history of status changes. Immutable agent activities retain the history of agent work.

Memory contains only a temporary view of the durable state. Active runtime state or durable leases determine which agents are running.

Harness stores a state change before it starts the next dependent action. This sequence supports restart recovery and prevents duplicate work.

SQLite stores the durable state.

### Isolate file changes

Each target has a type and a purpose. Target types are `Unspecified`, `Generic`, and `Git`. Target purposes are `Unspecified`, `Read`, and `ReadWrite`.

A `Read` target can be `Generic` or `Git`. A `ReadWrite` target must be `Git`.

The `Worktree` phase creates one worktree under `$HOME/.harness/worktree/{id}` for each Git target. Each worktree records its target, source revision, and worktree path.

Generic targets have no worktree. Harness provides them to agents as read-only context.

Agents do not write to source checkouts. Execution can change only worktrees for `ReadWrite` targets.

### Separate deterministic work from agent work

Deterministic code controls state transitions, Git operations, persistence, dependency order, and policy enforcement.

Agents create and review plans, change task files, and review results. Agent output cannot directly override the workflow state or safety policy.

### Keep architecture boundaries explicit

The repository uses an onion architecture. Each layer defines the interfaces that it needs from another layer.

The `internal/harness` and `internal/store` trees do not depend on unrelated local packages. The command layer connects concrete implementations.

### Make behavior observable

Harness records each ephemeral agent run and its immutable agent activities. Agent activities include starts, tool calls, reads, writes, errors, and completions.

Task, phase, worktree, step, plan, and review rows expose their current state. The Worktree and Apply phase `Error` fields expose redacted aggregate failures for their corresponding deterministic operations.

These records support the TUI, recovery, diagnosis, and later review. Secret values must not enter them.

## Architecture

### Command layer

`cmd/harness` is the composition root. It handles process signals, opens dependencies, constructs the application, and starts the TUI.

The command layer connects these concrete services:

- the SQLite store
- the agent runtime
- Git worktree operations
- the prompt catalog
- the logger

The command layer does not contain workflow rules.

### Harness layer

`internal/harness` contains application rules and service interfaces. It coordinates phases and applies the automation policy.

This layer owns:

- task lifecycle rules
- phase transitions
- pause, resume, cancel, and retry behavior
- plan review and refinement
- step dependency scheduling
- review loops
- apply authorization

External operations enter this layer through interfaces. Tests use fake implementations of these interfaces.

### Model layer

`internal/harness/model` contains durable domain records and enums. These types do not depend on the TUI or a store implementation.

The complete model requires these concepts:

- `Task`: user intent, current status, timestamps, and automation policy
- `Target`: one path available to a task, with its type and purpose
- `Worktree`: one isolated location for one Git target
- `Agent`: one durable record of an ephemeral agent run, with its task, role, and model
- `AgentActivity`: one immutable event for one agent run
- `Plan`: a stable identifier and structured revisions for one task
- `Step`: an executable unit with dependencies and current status
- `StepDependency`: one dependency edge between two steps
- `StepWorktree`: one link from a step to a writable worktree
- `Review`: one execution-review pass, findings, responses, and decision
- `Question`: one question that can receive an answer
- `PlanQuestion`: one link from a question to a plan revision
- `ReviewQuestion`: one link from a question to an execution review
- `Phase`: the current state of one workflow phase

There is no apply table and no plan-review table.

Each durable entity ID starts with a four-letter prefix that is unique to its table. A hyphen separates the prefix from an eight-character Crockford Base32 value. Relationship tables use composite keys instead of durable entity IDs.

`TaskStatus` describes the broad task lifecycle. Phase status and step status give detailed current progress. These rows do not retain a history of status changes.

### TUI layer

`internal/harness/ui` presents state and sends user intent to the application layer. It does not control workflow decisions.

The task list is the default view. It shows active tasks first and then sorts tasks by `UpdatedAt`.

The complete task interface requires:

- a selectable task list
- a create action on `c`
- a cancel action on `x`
- a create flow for one or more target paths, each target's type and purpose, a task description, and an automation policy
- a bottom information bar with available keyboard shortcuts
- clear states for pending input, failures, and progress

Keyboard interaction has priority. Mouse interaction is a later interface improvement.

### Terminal command interface

Harness will provide terminal commands for scripts and direct inspection. These commands run outside the interactive TUI.

The no-argument `harness` command continues to start the TUI. A terminal command performs its operation and then exits.

Each durable entity model provides commands that examine its data. Some models also provide commands that modify or remove data.

The initial command set is:

Task commands:

- `harness task get <task-id>`
- `harness task list`, alias `harness task ls`
- `harness task remove <task-id>`, alias `harness task rm <task-id>`

Target commands:

- `harness target get <target-id>`
- `harness target list`, alias `harness target ls`

Worktree commands:

- `harness worktree get <worktree-id>`
- `harness worktree list`, alias `harness worktree ls`

Agent commands:

- `harness agent get <agent-id>`
- `harness agent list`, alias `harness agent ls`

Agent activity commands:

- `harness agent-activity get <agent-activity-id>`
- `harness agent-activity list`, alias `harness agent-activity ls`

Plan commands:

- `harness plan get <plan-id> <revision>`
- `harness plan list`, alias `harness plan ls`

Step commands:

- `harness step get <step-id>`
- `harness step list`, alias `harness step ls`

Review commands:

- `harness review get <review-id>`
- `harness review list`, alias `harness review ls`

Question commands:

- `harness question get <question-id>`
- `harness question list`, alias `harness question ls`

Phase commands:

- `harness phase get <phase-id>`
- `harness phase list`, alias `harness phase ls`

List commands can use read-only filters for parent IDs, status, and result limits.

By default, each command prints a table with aligned, tab-separated columns. This table format matches Docker CLI table output.

The `--json` flag prints the same response as JSON.

The command layer sends each request through harness-layer query interfaces. Terminal commands do not read SQLite tables directly.

Task removal is the only terminal write operation in this command set. Later designs can add create, update, cancel, retry, and apply commands.

### Store layer

`internal/store` contains persistence implementations. The shared store tests define the required behavior for each implementation.

The stub store remains useful for unit tests. The SQLite store provides application persistence under `$HOME/.harness/`.

The SQLite implementation requires schema migrations. Migrations must preserve task history and support application upgrades.

The SQLite driver must support the current `CGO_ENABLED=0` build.

### Prompt catalog

`prompts` contains embedded templates for agent roles. The command layer loads the catalog and gives role definitions to the harness layer.

Each agent role has a model, prompt, prompt version, allowed tools, sandbox policy, and pass limit. A deterministic phase does not require an agent call.

For each agent run, Harness stores the selected model on the `Agent` row. The `Start` agent activity records the prompt identifier and version, provider configuration, and other safe run configuration in redacted details. This record connects each agent result to its instructions.

### Agent provider

`internal/provider` defines a provider-neutral interface for LLM and agent providers. An adapter can connect Harness to a provider such as `cursor-cli`.

The provider interface covers:

- provider identity and capabilities
- model selection
- new and resumed sessions
- prompts and worktree context
- normalized messages and tool calls
- cancellation
- provider errors

Provider adapters translate external commands or APIs into this interface. They do not contain workflow rules, prompts, persistence, or TUI behavior.

The provider layer depends on the harness layer. The harness layer does not import `internal/provider`.

The command layer selects a provider adapter. It injects this adapter through the narrow agent interface that the harness layer requires.

### Agent runtime

The agent runtime applies Harness policy to the selected provider. It starts isolated contexts and returns structured results.

One `Agent` row represents one ephemeral run. Each run has a `TaskID`, a role, and a model. Harness does not reuse an agent context for a later run.

The supported roles are:

- `Planner`: creates or refines a plan revision
- `PlanReviewer`: reviews the current draft revision of the plan
- `Executor`: runs one plan step
- `Reviewer`: reviews the complete execution result

The provider boundary keeps the runtime independent of one external command, API, or model vendor.

The runtime must support:

- a new context and new `Agent` record for every run
- role-specific models, versioned prompts, and pass limits
- allowed-tool enforcement
- writable worktree boundaries and read-only target context
- cancellation
- structured results
- immutable agent activities

Active runtime state or durable leases provide the running-agent count. The count must not derive from the number of `Agent` rows.

### Git worktree service

The Git worktree service creates, inspects, applies, and cleans worktrees. The harness layer controls the timing of these operations.

Worktree creation must not change any source checkout. The service creates one worktree for every Git target. This rule includes Git targets with purpose `Read`.

Apply changes only worktrees whose targets have purpose `ReadWrite`. Apply must occur only after the review result and automation policy permit it.

If one or more target operations fail during worktree preparation or apply, Harness writes one redacted aggregate summary to the corresponding phase `Error` field. Each update replaces the previous aggregate. Harness does not create a generic `Failure` activity.

## Runtime Flow

1. The command opens the SQLite database and starts the TUI.
2. The TUI loads durable tasks and focuses the first available task.
3. The user creates a task. The user enters one or more targets. The user selects a type and purpose for each target. The user enters a description. Then the user selects an automation policy.
4. The workflow coordinator starts or pauses at the `Worktree` boundary.
5. The worktree service prepares one worktree for each Git target. Generic targets remain read-only context without worktrees.
6. A fresh Planner agent creates the first draft plan revision.
7. A fresh PlanReviewer agent reviews the current draft revision.
8. If the PlanReviewer reports correctable findings and the configured pass limit permits another pass, a fresh Planner agent receives the draft and findings. The Planner creates the next revision. Then plan review repeats against that current draft.
9. When refinement stops, plan acceptance occurs at its separate policy or user boundary. A PlanReviewer conclusion does not accept the plan.
10. Before acceptance, the coordinator makes sure that the current plan dependency graph is valid.
11. Executor agents complete ready steps in dependency order. They can write only to linked worktrees for `ReadWrite` targets.
12. Reviewer agents examine the complete set of worktree results.
13. Correctable execution-review findings return to execution within the configured review-pass limits.
14. The `Apply` phase waits for the required authorization.
15. Harness applies each reviewed worktree result only to its `ReadWrite` target and marks the task complete.

Harness stores durable state throughout this flow. The TUI reads the same state and reports progress.

Each Planner, PlanReviewer, Executor, or Reviewer invocation creates a new `Agent` row. The agent activities record the prompt version, tool use, safe inputs and outputs, errors, and completion. The current task and phase status remain on their own rows.

## Task and Workflow State

A new task starts in `Ready`. Active phase work changes the task to `Running`.

A user-input boundary changes the task to `PendingFeedback`. Accepted input returns the task to `Ready` or `Running`.

A successful apply changes the task to `Completed`. User rejection changes it to `Canceled`.

If worktree preparation or apply fails for one or more targets, Harness sets the corresponding phase state. Harness writes one redacted aggregate summary to that phase's `Error` field. Each update replaces the previous aggregate. The user resolves the failure.

Automatic resolution is not in scope. Harness does not create a generic `Failure` activity for these failures.

A system or agent error changes the applicable current state to `Failed`. An agent error also creates an immutable `Error` agent activity for that run. A failed task remains available for inspection and retry.

Phase records contain detailed current progress. Task status remains a concise current summary for the task list. Task and phase records do not retain status-change history.

## Current Implementation

### Built

When `harness` has no subcommand, the command starts the TUI. It handles common process signals and closes the TUI after an interruption.

The current task model provides:

- a `TASK-` prefix and a four-character Crockford Base32 value
- title and timestamps
- lifecycle status values
- active-state classification
- active-first and recent-first comparison

The four-character Task value is the current implementation. The intended durable ID format uses four-letter prefixes and eight-character Crockford Base32 values.

The TUI provides:

- a Bubble Tea application shell
- a default task-list view
- keyboard list navigation
- first-item selection
- selection preservation across refreshes
- periodic task refresh
- task status styles and a running animation
- a status view with task and agent counts
- quit behavior on `q` and `ctrl+c`

The harness layer defines a small store contract. It supports task lists, task lookup, task creation, task counts, and running-agent counts.

The stub store implements this contract. Shared store tests define the current contract behavior.

A current `Activity` model and activity-type enum provide an initial legacy structure for events. This structure is incomplete relative to the intended agent-scoped, immutable `AgentActivity` design. The logger and UI theme also have concrete implementations.

### Partially built

The command uses the stub store and inserts fake tasks at startup. The displayed tasks are not durable user tasks.

The status view appears above the task list. It reports counts but does not provide the required bottom shortcut information.

The store contract can create tasks but cannot update, cancel, or advance them. The TUI has no task creation or cancellation flow.

The running-agent count exists in the contract, but the stub always returns zero. No runtime-state or durable-lease implementation supplies this count.

The SQLite and migration packages exist as placeholders. They do not open a database or store records.

The prompt package exists as a placeholder. It does not embed role templates or record prompt versions.

The data design describes target, worktree, agent, agent activity, plan, step, dependency, worktree-link, review, question, question-link, and phase concepts. Most of these concepts do not have complete Go models or persistence.

### Not built

The repository does not contain:

- durable SQLite task storage
- task update and cancellation operations
- the create-task modal sequence
- target-list navigation
- target type and purpose controls
- automation-policy controls
- a workflow coordinator
- durable phase state
- Git worktree management
- an agent runtime integration
- durable `Agent` and immutable `AgentActivity` storage
- sandbox and allowed-tool enforcement
- prompt templates and prompt-version recording
- plan creation, PlanReviewer review, and refinement
- plan dependency validation
- step scheduling
- parallel step execution
- execution-review loops
- apply behavior
- complete redacted agent audit records
- aggregate Worktree and Apply phase error reporting
- restart recovery for active work
- read-only terminal model commands and task removal

## Remaining Work

### Establish durable state

Complete the task fields and workflow models.

Add these durable-state records:

- targets with type and purpose
- one worktree for each Git target
- agents
- immutable agent activities
- plans
- steps and relationships
- reviews
- questions and relationships
- phases

Implement the intended four-letter prefixes and eight-character Crockford Base32 values for durable entity IDs.

Before durable compatibility depends on the Task value, migrate the current four-character format.

Implement SQLite migrations and all required store operations.

Replace fake startup data with the SQLite store. Keep the stub implementation for unit tests.

Derive the running-agent count from active runtime state or durable leases. Do not derive it from `Agent` row count.

### Add terminal commands

Add `get`, `list`, and `ls` commands for each durable entity model. Include `agent` and `agent-activity`.

Do not add generic activity commands for the intended model.

Add `remove` and `rm` for tasks.

Route read commands through the same query interfaces as the TUI.

Print a table with aligned, tab-separated columns as the default output. Add a `--json` flag that prints the same response as JSON.

Add parent, status, and result-limit filters where they apply.

Do not add terminal write commands other than task removal in this stage.

### Complete task interaction

Add the create action. Add the cancel action.

Add overlays for the target list, each target's type and purpose, task description, and automation policy.

Make sure that a `Read` target is `Generic` or `Git`. Make sure that a `ReadWrite` target is `Git`.

Move shortcut information to the bottom of the TUI. Preserve the current task-list behavior and styles.

### Add workflow coordination

Implement explicit phase transitions and automation boundaries. Add pause, resume, cancel, retry, and restart behavior.

Before the coordinator starts a dependent action, store durable state. Task and phase rows retain current status without status-change history.

### Add worktree isolation

Create one Git worktree for each Git target. Record each target, source revision, and worktree path.

For a Generic target, do not create a worktree. Provide Generic target paths as read-only context.

Handle existing worktrees. Store one redacted aggregate summary of per-target worktree failures in the Worktree phase `Error` field. Replace the previous aggregate on each update.

During preparation, do not modify any source checkout.

### Add controlled agent execution

Define the provider interface. Add the first provider adapter. Connect the adapter through the agent interface that the harness layer requires.

Add role definitions and embedded, versioned prompt templates for Planner, PlanReviewer, Executor, and Reviewer.

Apply the model and provider selection from the configuration.

Create one `Agent` row for every ephemeral run. Record the task, role, and model on this row.

Record the selected prompt identifier, prompt version, and safe run configuration in the run's `Start` agent activity.

Before agents perform task work, enforce these controls:

- the sandbox policy
- the writable-worktree boundary
- the read-only context boundary
- the allowed-tool policy

Record redacted, immutable agent activities for every agent run.

### Add plan refinement

Create a structured plan format. Start each Planner and PlanReviewer run in a new agent context.

The first Planner creates revision 1. A PlanReviewer reviews only the current draft revision.

Store the redacted conclusion and correctable findings in the `Complete` agent activity details. Do not add a plan-review table.

If correctable findings remain and another configured pass is permitted, give the current draft and findings to a fresh Planner context.

Store each successful refinement as a new row with the same plan ID and the next revision number. Then use a fresh PlanReviewer.

Tell the PlanReviewer to review the current draft.

Assign the plan ID to the first revision. Use the next revision number for each new row. Accept the current revision through an update to its existing row.

Plan acceptance remains a separate automation-policy or user boundary. A PlanReviewer does not accept a plan.

Do not increment the revision number for acceptance. Do not refine an accepted plan. Do not create a later revision.

Before acceptance, reject duplicate step IDs, unknown dependencies, and dependency cycles.

### Add dependency-based execution

After its dependencies succeed, run the step. If task policy permits parallel work, run independent steps in parallel.

Link steps only to worktrees whose targets have purpose `ReadWrite`. Executors can read other task context but cannot change `Read` targets.

Store the step state and agent results. Resume incomplete execution without another run of successful steps.

### Add review and apply

Run execution-review passes against all worktree results. Within the configured limits, return correctable findings to execution.

Define apply behavior and its authorization boundary. Apply only worktrees whose targets have purpose `ReadWrite`. Do not change `Read` targets.

If apply fails for one or more targets, store one redacted aggregate summary in the Apply phase `Error` field. Replace the previous aggregate on each update. Do not create a generic `Failure` activity.

The user resolves the failure. Automatic resolution of a partial apply is not in scope.

### Harden completed workflows

Add agent-activity redaction, a retention policy, worktree cleanup, concurrency limits, durable leases where needed, and conflict handling.

After the keyboard workflow is complete and stable, add mouse interaction.

## Quality Strategy

Model and coordinator rules use small tests. Store behavior uses the shared store contract.

SQLite and Git worktree behavior use medium tests with temporary resources. Agent orchestration uses deterministic fake agents.

TUI tests cover keyboard input, modal transitions, selection, layout, and displayed state. End-to-end tests cover recovery across process boundaries.

Every top-level test uses the applicable `testsize` assertion.

## Open Design Decisions

- Select the first agent runtime and supported model providers.
- Define the phase boundaries that always require user authorization.
- Define limits for plan-refinement passes, PlanReviewer passes, execution-review passes, parallel steps, concurrent agents, and concurrent tasks.
- Define worktree retention for completed, canceled, and failed tasks.
- Define agent-activity redaction and retention rules.
- When parallel steps change the same files, define conflict behavior.
- If leases provide the running-agent count, define durable-lease ownership, expiry, and recovery.
- When a task resumes after a catalog upgrade, define prompt-version compatibility.

## Intended End State

The no-argument command opens a persistent task list. The user can create, inspect, cancel, resume, and complete tasks from the TUI.

Each task has one or more targets with an explicit type and purpose. Each Git target uses one isolated worktree. Generic targets remain read-only context without worktrees. Only `ReadWrite` Git targets can receive changes.

Plans use fresh Planner and PlanReviewer contexts, current-draft review, configured refinement limits, a separate acceptance boundary, and valid dependency graphs. Execution follows these graphs, and execution review controls whether Harness can continue to apply authorization.

Each ephemeral agent run has a durable `Agent` record and immutable `AgentActivity` history. Prompt and model selection remain attributable to that run. Task and phase rows show current state without status-change history.

Worktree and apply failures appear as replaceable, redacted aggregate summaries on their corresponding phase records. Process restarts preserve progress and do not repeat completed work.
