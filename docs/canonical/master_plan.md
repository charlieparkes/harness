```
date: 2026-09-04
version: v1.1.1
```

---

# Master Plan

## Purpose

Harness is a local application that manages structured agent workflows. It provides a terminal interface, durable state, workspace isolation, and human control.

This document describes the intended design and architecture. It also records the repository contents and the remaining work.

The functional requirements inform the design. They do not define a fixed product contract.

## Product Direction

Harness converts one task description into a controlled workflow. The main phases are `workspace`, `plan`, `execute`, `review`, and `apply`.

Each phase has one responsibility:

- `workspace` prepares isolated locations for the task targets.
- `plan` creates and refines a structured plan.
- `execute` completes plan steps in dependency order.
- `review` examines the result and requests corrections for defects.
- `apply` writes each accepted workspace result to its target.

The user selects the phase boundaries that require input. Automation can continue through the other boundaries.

Harness stores task and phase state after each material change. A process restart preserves accepted work and does not repeat completed work.

## Design Principles

### Preserve human control

The automation policy determines when Harness pauses. Harness must not perform a destructive or uncertain action beyond this policy.

The TUI shows the current state, available actions, failures, and requests for user input. Background work does not block TUI updates.

### Use durable state

The store package is the source of truth for active and historical data. This data includes tasks, targets, workspaces, and agents.

Memory contains only a temporary view of the durable state.

Harness stores a state change before it starts the next dependent action. This sequence supports restart recovery and prevents duplicate work.

SQLite stores the durable state.

### Isolate file changes

Git targets use worktrees under `$HOME/.harness/worktree/{id}`. Agents do not write to source checkouts.

The workspace phase creates one workspace for each target. Each workspace records its target, source revision, and worktree path.

### Separate deterministic work from agent work

Deterministic code controls state transitions, Git operations, persistence, dependency order, and policy enforcement.

Agents create plans, change task files, and review results. Agent output cannot directly override the workflow state or safety policy.

### Keep architecture boundaries explicit

The repository uses an onion architecture. Each layer defines the interfaces that it needs from another layer.

The `internal/harness` and `internal/store` trees do not depend on unrelated local packages. The command layer connects concrete implementations.

### Make behavior observable

Harness records phase changes, agent conclusions, tool calls, policy rejections, failures, and user decisions.

Activity records support the TUI, recovery, diagnosis, and later review. Secret values must not enter these records.

## Architecture

### Command layer

`cmd/harness` is the composition root. It handles process signals, opens dependencies, constructs the application, and starts the TUI.

The command layer connects these concrete services:

- the SQLite store
- the agent runtime
- Git workspace operations
- the prompt catalog
- the logger

The command layer does not contain workflow rules.

### Harness layer

`internal/harness` contains application rules and service interfaces. It coordinates phases and applies the automation policy.

This layer owns:

- task lifecycle rules
- phase transitions
- pause, resume, cancel, and retry behavior
- plan refinement
- step dependency scheduling
- review loops
- apply authorization

External operations enter this layer through interfaces. Tests use fake implementations of these interfaces.

### Model layer

`internal/harness/model` contains durable domain records and enums. These types do not depend on the TUI or a store implementation.

The complete model requires these concepts:

- `Task`: user intent, status, timestamps, and automation policy
- `Target`: one path affected by a task
- `Workspace`: one isolated location for one target
- `Plan`: a stable identifier and structured revisions for one task
- `Step`: an executable unit with dependencies and status
- `Review`: one review pass, findings, responses, and decision
- `Activity`: an immutable event for a task or phase
- `Phase`: the current workflow position

Each durable model ID starts with a four-letter prefix that is unique to its table. A hyphen separates the prefix from a four-character Crockford Base32 value.

`TaskStatus` describes the broad task lifecycle. Phase status and step status give detailed progress.

### TUI layer

`internal/harness/ui` presents state and sends user intent to the application layer. It does not control workflow decisions.

The task list is the default view. It shows active tasks first and then sorts tasks by their last activity.

The complete task interface requires:

- a selectable task list
- a create action on `c`
- a cancel action on `x`
- a create flow for one or more target paths, a task description, and an automation policy
- a bottom information bar with available keyboard shortcuts
- clear states for pending input, failures, and progress

Keyboard interaction has priority. Mouse interaction is a later interface improvement.

### Terminal command interface

Harness will provide terminal commands for scripts and direct inspection. These commands run outside the interactive TUI.

The no-argument `harness` command continues to start the TUI. A terminal command performs its operation and then exits.

Each durable model provides commands that examine its data. Some models also provide commands that modify or remove data.

The initial command set is:

Task commands:

- `harness task get <task-id>`
- `harness task list`, alias `harness task ls`
- `harness task remove <task-id>`, alias `harness task rm <task-id>`

Target commands:

- `harness target get <target-id>`
- `harness target list`, alias `harness target ls`

Workspace commands:

- `harness workspace get <workspace-id>`
- `harness workspace list`, alias `harness workspace ls`

Plan commands:

- `harness plan get <plan-id> <revision>`
- `harness plan list`, alias `harness plan ls`

Step commands:

- `harness step get <step-id>`
- `harness step list`, alias `harness step ls`

Review commands:

- `harness review get <review-id>`
- `harness review list`, alias `harness review ls`

Activity commands:

- `harness activity get <activity-id>`
- `harness activity list`, alias `harness activity ls`

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

`prompts` contains embedded templates for agent phases. The command layer loads the catalog and gives phase definitions to the harness layer.

Each agent phase has a model, prompt, allowed tools, sandbox policy, and pass limit. A deterministic phase does not require an agent call.

Harness includes prompt versions in the activity history. This record connects each agent result to its instructions.

### Agent provider

`internal/provider` defines a provider-neutral interface for LLM and agent providers. An adapter can connect Harness to a provider such as `cursor-cli`.

The provider interface covers:

- provider identity and capabilities
- model selection
- new and resumed sessions
- prompts and workspace context
- normalized messages and tool calls
- cancellation
- provider errors

Provider adapters translate external commands or APIs into this interface. They do not contain workflow rules, prompts, persistence, or TUI behavior.

Provider dependencies point toward the harness layer. The harness layer does not import `internal/provider`.

The command layer selects a provider adapter. It injects this adapter through the narrow agent interface that the harness layer requires.

### Agent runtime

The agent runtime applies Harness policy to the selected provider. It starts isolated contexts and returns structured results.

The provider boundary keeps the runtime independent of one external command, API, or model vendor.

The runtime must support:

- a new context for each plan refinement pass
- phase-specific models and prompts
- allowed-tool enforcement
- workspace boundaries
- cancellation
- structured results
- tool-call activity

Active runtime state or durable leases provide the running-agent count. The store must not keep this count as a fixed value.

### Git workspace service

The Git workspace service creates, inspects, applies, and cleans worktrees. The harness layer controls when these operations occur.

Workspace creation must not change any source checkout. Apply must occur only after the review result and automation policy permit it.

## Runtime Flow

1. The command opens the SQLite database and starts the TUI.
2. The TUI loads durable tasks and focuses the first available task.
3. The user creates a task, enters one or more target paths, and selects its automation policy.
4. The workflow coordinator starts or pauses at the `workspace` boundary.
5. The workspace service prepares one workspace for each Git target.
6. The planner creates a plan and runs the configured refinement passes.
7. The coordinator makes sure that the plan dependency graph is valid.
8. Executor agents complete ready steps in dependency order against the task workspaces.
9. Reviewer agents examine the complete set of workspace results.
10. Correctable findings return to execution within the configured limits.
11. The apply phase waits for the required authorization.
12. Harness writes each accepted workspace result to its target and marks the task complete.

Harness stores durable state throughout this flow. The TUI reads the same state and reports progress.

## Task and Workflow State

A new task starts in `Ready`. Active phase work changes the task to `Running`.

A user-input boundary changes the task to `Pending Feedback`. Accepted input returns the task to `Ready` or `Running`.

A successful apply changes the task to `Completed`. User rejection changes it to `Canceled`.

If apply fails for a target, Harness writes a `Failure` activity for the apply phase. The user resolves the failure.

Automatic resolution is not in scope.

A system or agent error changes the task to `Failed`. A failed task remains available for inspection and retry.

Phase records contain detailed progress. Task status remains a concise summary for the task list.

## Current Implementation

### Built

The command starts the TUI when `harness` has no subcommand. It handles common process signals and closes the TUI after an interruption.

The task model provides:

- a `TASK-` prefix and a four-character Crockford Base32 value
- title and timestamps
- lifecycle status values
- active-state classification
- active-first and recent-first comparison

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

An activity model and an activity-type enum provide an initial event structure. The logger and UI theme also have concrete implementations.

### Partially built

The command uses the stub store and inserts fake tasks at startup. The displayed tasks are not durable user tasks.

The status view appears above the task list. It reports counts but does not provide the required bottom shortcut information.

The store contract can create tasks but cannot update, cancel, or advance them. The TUI has no task creation or cancellation flow.

The running-agent count exists in the contract, but the stub always returns zero.

The SQLite and migration packages exist as placeholders. They do not open a database or store records.

The prompt package exists as a placeholder. It does not embed phase templates.

The data design describes target, workspace, plan, step, activity, and review concepts. Most of these concepts do not have complete Go models or persistence.

### Not built

The repository does not contain:

- durable SQLite task storage
- task update and cancellation operations
- the create-task modal sequence
- target-list navigation
- automation-policy controls
- a workflow coordinator
- durable phase state
- Git worktree management
- an agent runtime integration
- sandbox and allowed-tool enforcement
- prompt templates
- plan refinement
- plan dependency validation
- step scheduling
- parallel step execution
- review loops
- apply behavior
- complete audit records
- restart recovery for active work
- read-only terminal model commands and task removal

## Remaining Work

### Establish durable state

Complete the task fields and workflow models. Add target records and target-specific workspaces.

Implement SQLite migrations and all required store operations.

Replace fake startup data with the SQLite store. Keep the stub implementation for unit tests.

### Add terminal commands

Add `get`, `list`, and `ls` commands for each durable model. Add `remove` and `rm` for tasks.

Route read commands through the same query interfaces as the TUI.

Print a table with aligned, tab-separated columns as the default output. Add a `--json` flag that prints the same response as JSON.

Add parent, status, and result-limit filters where they apply.

Do not add terminal write commands other than task removal in this stage.

### Complete task interaction

Add the create and cancel actions. Add overlays for the target list, task description, and automation policy.

Move shortcut information to the bottom of the TUI. Preserve the current task-list behavior and styles.

### Add workflow coordination

Implement explicit phase transitions and automation boundaries. Add pause, resume, cancel, retry, and restart behavior.

The coordinator must store durable state before it starts a dependent action.

### Add workspace isolation

Create one Git worktree for each target. Record each target, source revision, and worktree path.

Handle existing worktrees. Record workspace errors as activities.

Do not modify any source checkout during preparation.

### Add controlled agent execution

Define the provider interface and add the first provider adapter. Connect it through the agent interface that the harness layer requires.

Add phase definitions and embedded prompt templates. Apply the model and provider selection from the configuration.

Enforce the sandbox policy and allowed-tool policy before agents perform task work. Record safe activity data for every agent run.

### Add plan refinement

Create a structured plan format. Start each refinement pass in a new agent context.

Give the prior plan to the new context.

Assign the plan ID to the first revision. Store each successful pass as a new row with the same plan ID.

Use the next revision number for each new row. Accept the current revision through an update to its existing row.

Do not increment the revision number for acceptance. Do not refine an accepted plan or create a later revision.

Before execution, reject duplicate step IDs, unknown dependencies, and dependency cycles.

### Add dependency-based execution

Run a step only after its dependencies succeed. If task policy permits parallel work, run independent steps in parallel.

Store the step state and agent results. Resume incomplete execution without another run of successful steps.

### Add review and apply

Run review passes against all workspace results. Return correctable findings to execution within the configured limits.

Define apply behavior and its authorization boundary.

If apply fails for a target, record a `Failure` activity for the apply phase. The user resolves the failure.

Automatic resolution of a partial apply is not in scope.

### Harden completed workflows

Add activity redaction, a retention policy, worktree cleanup, concurrency limits, and conflict handling.

Add mouse interaction only after the keyboard workflow is complete and stable.

## Quality Strategy

Model and coordinator rules use small tests. Store behavior uses the shared store contract.

SQLite and Git worktree behavior use medium tests with temporary resources. Agent orchestration uses deterministic fake agents.

TUI tests cover keyboard input, modal transitions, selection, layout, and displayed state. End-to-end tests cover recovery across process boundaries.

Every top-level test uses the applicable `testsize` assertion.

## Open Design Decisions

- Select the first agent runtime and supported model providers.
- Define behavior for targets that are not Git repositories.
- Define the phase boundaries that always require user authorization.
- Define limits for plan passes, review passes, parallel steps, and concurrent tasks.
- Define worktree retention for completed, canceled, and failed tasks.
- Define activity redaction and retention rules.
- Define conflict behavior when parallel steps change the same files.

## Intended End State

The no-argument command opens a persistent task list. The user can create, inspect, cancel, resume, and complete tasks from the TUI.

Each task has one or more targets. Each Git target uses one isolated worktree.

Plans use fresh-context refinement and valid dependency graphs. Execution follows these graphs, and review controls whether Harness applies a result.

The database and activity history explain each material decision. Process restarts preserve progress and do not repeat completed work.
