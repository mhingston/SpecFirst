# Production Readiness Philosophy

To be considered production-ready, code must be **Assertive** and **Observable**.

## 1. Assertive Programming (Crash Early)
- **Invariants vs. Errors**: Distinguish between "errors" (network down, bad input) and "impossible states" (null ID on active user).
- **The Crash Rule**: If an invariant is violated, the system is corrupt. Do not return an error. **Panic/Crash immediately.**
- **No Swallow**: Never catch an error without re-throwing or logging with full context.

## 2. Contextual Observability
- **Logs are for Context**: A log message saying "Error" is useless. It must say *who*, *what*, and *state*.
- **The Reproduction Rule**: A log entry must contain enough variable data to reproduce the issue locally without guessing.

## 3. Boring Code
- **Explicit over Clever**: Unroll loops if it makes logic clearer.
- **No Magic**: Avoid implicit state mutations.
