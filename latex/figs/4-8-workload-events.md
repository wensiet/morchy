sequenceDiagram
    participant A as Agent
    participant CP as Control Plane API
    participant DB as PostgreSQL

    Note over A,DB: Event Sequence
    A->>CP: 1. POST /events (container_created)
    CP->>DB: 2. INSERT event

    A->>CP: 3. POST /events (healthcheck_success)
    CP->>DB: 4. INSERT event

    A->>CP: 5. POST /events (healthcheck_success)
    CP->>DB: 6. INSERT event

    Note over CP: Status Inference: Status = active
    CP->>DB: 7. Query events for status
    DB-->>CP: 8. Return event history