```
date: 2026-09-07
version: v1.1.0
```

---

# Application Design

## Principals

This application follows an onion architecture. 

- The domain should never import from the presentation (tui) or store layers, but these other layers may import from domain.
- When the domain needs to interact with something outside of itself, it should define an interface. For example, define a store interface, rather than importing from the store package.

## Diagram

```mermaid
%%{init: {"flowchart": {"curve": "monotoneX", "nodeSpacing": 80, "rankSpacing": 100}}}%%
flowchart TD
    subgraph ENTRY["Executable"]
        CMD["cmd/harness<br/>CLI startup and dependency wiring"]
    end

    subgraph PRESENTATION["Presentation"]
        TUI["internal/tui<br/>Terminal UI"]
        STYLES["internal/tui/styles<br/>Styling"]
    end

    subgraph CORE["Domain"]
        DOMAIN["internal/domain<br/>Business logic"]
        MODEL["internal/domain/model<br/>Application models"]
        AGENTMODEL["internal/domain/agentmodel<br/>Agent API contracts"]
        PROMPT["internal/domain/prompt<br/>Agent prompt templates"]
    end

    subgraph STORE["Store"]
        STUB["internal/store/stub<br/>In-memory implementation"]
        SQLITE["internal/store/sqlite<br/>SQLite implementation"]
        STORETEST["internal/store/storetest<br/>Store test suite"]
    end

    ENTRY --> PRESENTATION
    ENTRY --> CORE

    CORE --> STORE

    TUI --> STYLES

    DOMAIN --> MODEL
    DOMAIN --> AGENTMODEL
    DOMAIN --> PROMPT
```
