---
name: promote-documentation
version: 1.0.0
description: |
  Instructions on how to promote and codify new or modified documentation.
compatibility: claude-code cursor codex gemini-cli opencode
---

- Move the indicated document from `docs/apocryphal/` to `docs/canonical/`.
- For new documents, pick a simple, one to three word filename, underscore separated.
- For modified documents, replace the existing document, using the existing document name. See `source: {path}` in the staged document to find it's source path in `docs/canonical/`.
- When promoting a document, keep the metadata fields indicated by `docs/canonical/_default.tmpl.md`.