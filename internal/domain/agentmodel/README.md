# agentmodel

A series of models which define the acceptable JSON output from an agent.

These models are effectively an API contract.

A harness-initiated agent that outputs anything other than one of these schemas will be outright rejected.

Agent responses will be sanitized and verified before storage.