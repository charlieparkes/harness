---
name: my-plan
version: 1.0.0
description: |
  Define consistent standard for writing plans.
compatibility: claude-code cursor codex gemini-cli opencode
---

The plan should match this structure exactly.

```md
# Summary

{single sentence or short paragraph}

## Assumptions

- {assumption}

## Steps

### S1 - {step title}

{step description}

#### Dependencies

- {step ID this step depends on}

#### Files

- {affected file}

#### Verification

- {command or process to verify step completion}

## Risks

- {risk}

## Questions

- {question}

## Definition of Done

{definition of done}
```


- Every step must have verifiable completion criteria in its `verification` array (commands, tests, or concrete checks).
- The plan must identify tests or validation commands that prove the work is done, using only the permitted programs above.
- State assumptions explicitly in `assumptions`.
- List material unresolved questions in `questions` when you cannot answer them from the repository.
- Step ids must be unique. Dependencies must reference existing step ids only.
- Use `definition_of_done` for overall acceptance criteria beyond individual step verification.
- Plan the smallest thing that satisfies the task. Before adding a step, check in this order whether the work is needed at all, already exists in this repository, is covered by the standard library, or is covered by an already-installed dependency. Do not plan a step for work one of those already does.
- Do not plan abstractions, indirection layers, or dependencies the task did not ask for. Fewest steps, fewest files, fewest new dependencies. The exception is structure the repository's own architecture requires: where a layering rule means a dependency has to be inverted (a domain package that must not import the store or the api, for example), plan the interface that layering needs. Follow the boundary the repository already draws rather than inventing one, and name it in `assumptions` when it is why a step exists.
- When the task is a bug, plan the fix at the cause: one change in the shared function every caller reaches, not one change per calling site.
- Do not plan steps whose product is documentation, comments, or boilerplate nobody asked for.
- Non-trivial logic gets one verification that fails if the logic breaks, in the style of tests this repository already has. A trivial one-liner does not need its own test.
