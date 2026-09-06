# lineage

## Revisioned Fields

### RevisionedField `lineage.RevisionedField[T]`

| Field    | Go type                | JSON field | Notes              |
| -------- | ---------------------- | ---------- | ------------------ |
| `Values` | `[]RevisionedValue[T]` | `values`   | Omitted when empty |


### RevisionedValue `lineage.RevisionedValue[T]`

| Field      | Go type | JSON field | Notes                                               |
| ---------- | ------- | ---------- | --------------------------------------------------- |
| `Revision` | `int64` | `revision` | Parent document revision at which the value was set |
| `Value`    | `T`     | `value`    | Omitted when it has the zero value                  |


