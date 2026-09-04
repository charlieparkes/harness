---
name: modify-documentation
version: 1.0.0
description: |
  Instructions on how to modify existing documentation.
compatibility: claude-code cursor codex gemini-cli opencode
---

# Modifying apocryphal document
1. Update document in-place. Do not change the version or any other metadata.

# Modifying canonical document
1. Create a new file in `docs/apocryphal/` based on `_modify.tmpl.md` as `{YYYY-MM-DD}-{original document name}.md`.
2. Choose the next version based on scope of requested change. Apocryphal document version is incremented only once at time of file creation.
3. Copy contents of source file into {document contents}.
4. Update contents, as instructed.
