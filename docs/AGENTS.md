# File Structure

- `todo.md`
- `canonical/`: High quality, human verified documentation
- `apocryphal/`: Incomplete, LLM generated, or poorly-defined documentation

# Instructions

- NEVER modify documents in `canonical/` without *explicit* user acknowledgment.
- NEVER promote an `apocryphal/` document to `canonical/` without *explicit* user acknowledgement.
- When asked to modify existing documentation, follow template `docs/apocryphal/_modify.tmpl.md`.
- When asked to write documentation, follow template `docs/apocryphal/_create.tmpl.md`.
- When asked to promote documentation, follow template `docs/canonical/_standard.tmpl.md`.

# Validation

- No explicit validation steps for docs.
- Do not run `make` for changes that only modify files in `docs/`.