# TODO

## High Priority

- [ ] Define storage pattern for json document driven tables. 
    - [x] How are updates applied to the domain model? Standardize way that agentmodel returned by an LLM response is applied to a domain model. Increment revision on domain model, and set new revisioned value on each field? Should any fields be ignored?
    - [ ] How is the domain model reflected in the store? Likely the domain model is converted to json, and stored in a "document" column in a store table. Represent stored data by defining a `{ModelName}Record` struct for every top level entity stored to the database. Eliminate all links tables in the store. Store everything in top level model aligned tables. Some data from the domain model is pulled out to table columns. For example, ID, revision, timestamps. Also, foreign keys, if foreign key represents another top level model.
- [ ] Align documentation/implementation for agent models, domain models, and store models

## Medium Priority

- [ ] Implement stub store methods and store testsuite.
- [ ] Implement per-phase agent prompt templates using go templates. Templates should inject a required json schema for output based on available `internal/domain/agentmodel` models.
- [ ] Implement tool call system prompt template which will be generated from go code, depending on available tools
- [ ] Implement git repo worktree management. `internal/domain/worktree`, `NewWorktreeManager`
- [ ] Implement agent tool calls. (Approved tool list, pre-defined allowed commands.)
- [ ] Implement agent audit logs "AgentActivity". Log tool calls / reads / writes.

## Low Priority

- [ ] Enable mouse interaction using https://github.com/lrstanley/bubblezone
- [ ] Implement step DAG with parallel agents (show "allow multitasking plan" option?)
- [ ] Auto review tool calls
- [ ] Agent sandboxing? Using containers or another 3rd party tool?

## Complete

- [x] Implement revisioned fields (lineage).
- [x] Write application models.
- [x] Define data model.
- [x] Define master plan.
- [x] Design task list interface.